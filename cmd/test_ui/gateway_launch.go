package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/config"
	"go-gateway/internal/desktop"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/web"
)

const logOutputFile = "file"

var errLaunchHandled = errors.New("launch canceled or handed to verified owner")

func prepareHost(h *gatewayHost) error {
	if h.startupCtx.Err() != nil {
		return errLaunchHandled
	}
	exe, err := os.Executable()
	if err != nil {
		return apphost.NewFault("startup.path_unusable", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return apphost.NewFault("startup.path_unusable", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return apphost.NewFault("startup.path_unusable", err)
	}
	if h.launch.Mode != apphost.Desktop {
		h.cfg, err = config.Load()
		if err != nil {
			return apphost.NewFault("startup.config_invalid", err)
		}
	}
	h.data, err = apphost.ResolveData(apphost.DataInput{Mode: h.launch.Mode, Executable: exe, WorkingDir: cwd, Environment: environmentValues()})
	if err != nil {
		return err
	}
	if h.data.ConfirmCreate && !desktop.ConfirmCreateDB(h.data.DatabasePath) {
		return errLaunchHandled
	}
	if h.launch.Mode == apphost.Desktop {
		h.cfg, err = config.LoadFromDirectory(h.data.ConfigRoot)
		if err != nil {
			return apphost.NewFault("startup.config_invalid", err)
		}
		if err = os.Setenv("GATEWAY_DB_PATH", h.data.DatabasePath); err != nil {
			return apphost.NewFault("startup.config_invalid", err)
		}
	}
	if h.cfg.Log.Output != "console" && h.cfg.Log.Output != logOutputFile {
		return apphost.NewFault("startup.config_invalid", errors.New("LOG_OUTPUT must be console or file"))
	}
	h.owner, err = desktop.AcquireOwner(h.data.DatabasePath, h.launch.Mode == apphost.Desktop, func() {
		if err := h.openSetup(); err != nil {
			h.emit("http.server_error", nil)
		}
	})
	if err != nil {
		if errors.Is(err, desktop.ErrAlreadyRunning) {
			if h.launch.Mode == apphost.Desktop && desktop.RequestOpenSetup(h.data.DatabasePath) == nil {
				return errLaunchHandled
			}
			return apphost.NewFault("startup.owner_busy", err)
		}
		if errors.Is(err, desktop.ErrHardlinkedDatabase) {
			return apphost.NewFault("startup.database_linked", err)
		}
		return apphost.NewFault("startup.path_unusable", err)
	}
	if h.startupCtx.Err() != nil {
		return errLaunchHandled
	}
	identity := h.owner.Identity()
	if h.launch.Mode == apphost.Desktop || h.cfg.Log.Output == logOutputFile {
		path := os.Getenv("LOG_FILE")
		if path == "" {
			path = diagnostics.DefaultPath(h.data.DataRoot, identity.ID)
		} else if !filepath.IsAbs(path) {
			path = filepath.Join(h.data.DataRoot, path)
		}
		h.sink, err = diagnostics.OpenFileSink(diagnostics.FileOptions{Path: path, DatabasePath: identity.Path, DatabaseID: identity.ID})
		if err != nil {
			return apphost.NewFault("startup.log_unavailable", err)
		}
	}
	h.logMu.Lock()
	h.logs, err = diagnostics.New(diagnostics.Options{Level: h.cfg.Log.Level, File: h.sink})
	h.logMu.Unlock()
	if err != nil {
		return apphost.NewFault("startup.config_invalid", err)
	}
	if h.startupCtx.Err() != nil {
		return errLaunchHandled
	}
	h.installLogging(h.cfg.Log.Output == "console")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	h.emit("startup.begin", nil)
	if h.launch.Mode == apphost.Desktop {
		factory := h.newShell
		if factory == nil {
			factory = desktop.New
		}
		h.stateMu.Lock()
		h.shell = factory(desktop.Options{BuildMode: string(h.launch.Mode), DataPath: identity.Path, ShutdownDone: h.shutdownDone, OnQuit: h.requestQuit, OnSessionEnd: h.requestQuit, OnFault: h.recordShellFailure, OnForceQuit: func() { os.Exit(2) }})
		h.stateMu.Unlock()
		if err = h.shell.Start(); err != nil {
			return apphost.NewFault("startup.tray_failed", err)
		}
		h.trayReady = true
		// A usable control surface exists before any potentially stalled initial flush.
		ctx, cancel := context.WithTimeout(h.startupCtx, 5*time.Second)
		err = h.logs.Flush(ctx)
		cancel()
		if err != nil {
			return apphost.NewFault("startup.log_unavailable", err)
		}
	}
	if h.startupCtx.Err() != nil {
		return errLaunchHandled
	}
	lc := net.ListenConfig{}
	h.listener, err = lc.Listen(h.startupCtx, "tcp", h.cfg.GetServerAddr())
	if err != nil {
		return apphost.NewFault("startup.bind_failed", err)
	}
	if err = web.ValidateStaticFiles(staticFiles); err != nil {
		return apphost.NewFault("startup.assets_missing", err)
	}
	return nil
}
func environmentValues() map[string]string {
	values := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	return values
}
func runCommand(args []string) int {
	launch, err := apphost.Parse(args, desktop.DesktopBuildDefault())
	if err != nil {
		if desktop.DesktopBuildDefault() && !slices.Contains(args, "--headless") {
			desktop.ShowError(desktop.ErrorInfo{Code: "startup.config_invalid", Message: "命令列選項無效。", NextStep: "使用--help查看支援的選項。"})
		} else {
			log.Print(err)
		}
		return 2
	}
	if launch.Informational() {
		text := fmt.Sprintf("Go Gateway\nVersion: %s\nCommit: %s\nMode: %s\n", desktop.Version, desktop.Commit, launch.Mode)
		if launch.Help {
			text += "Usage: gateway [--headless] [--help] [--version]\nExplicit --headless disables tray, native dialogs and automatic browser opening.\n"
		}
		if launch.Mode == apphost.Desktop {
			desktop.ShowInfo("Go Gateway", text)
		} else {
			fmt.Print(text)
		}
		return 0
	}
	if err = runGateway(launch); err != nil {
		if launch.Mode != apphost.Desktop && os.Getenv("LOG_OUTPUT") != logOutputFile {
			log.Print(err)
		}
		return 1
	}
	return 0
}
