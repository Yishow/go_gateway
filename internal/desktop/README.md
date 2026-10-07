# Desktop platform boundary

The gateway owns readiness, acquisition state, shutdown ordering and force exit.
This package owns the native notification icon, fixed-route browser actions,
local dialogs and Windows database ownership. It does not start a collector.

## Ownership and local handoff

After the caller resolves the intended database and confirms implicit creation,
`AcquireOwner` reserves a missing zero-length file and retains its actual handle
without delete sharing. An exclusive `LockFileEx` byte at
`0x7ffffffffffffffe` is beyond SQLite's locking range and maximum database extent;
locking it neither writes nor extends the file. The actual file lock handles
case, symlink and junction aliases across sessions and users. Hardlink identity
is recognized for duplicate-owner checks, but a fresh writable owner rejects a
file with multiple names because SQLite WAL/SHM recovery is pathname-based.
Preserve the original database together with its journal; never move/delete
sidecars or choose another hardlink to bypass the check. Resolve any hardlink
ambiguity with a verified recovery/backup procedure before restarting. Process
death releases OS ownership. A PID file is never an ownership authority.

`FileIdInfo` supplies the full 128-bit file ID plus 64-bit volume identity for
opaque namespaces; unsupported identity providers fail closed. Before opening
a file, native paths must be plain absolute drive/UNC paths. ADS, device
namespaces, drive-relative paths, reserved device names and ambiguous trailing
dots/spaces are rejected rather than normalized into another file. Keep ownership
until cleanup truly completes. A shutdown timeout does not release the lock.

The owner also pins every canonical ancestor directory, root first, with no
delete or write sharing. Reparse components are rejected; the file identity is
rechecked after the whole namespace is pinned. Pin failures release every handle
and fail closed. These are filesystem-object guards, not changes to OS settings;
privileged changes to volume/share mappings are outside this local guard.

Only desktop owners create IPC. A named pipe has an explicit current-user DACL,
rejects remote clients, and authenticates user SID plus session on both sides.
The secondary also checks actual executable identity. Only a fixed one-byte
setup action is accepted; the protocol does not take URLs or command strings.
Client reads and acknowledgements are bounded, including idle/malformed clients.

## Native shell and shutdown

The tray uses an OS-thread-affine message loop with a hidden top-level window.
Its icon is embedded in the Go executable. The fixed menu provides setup, logs,
read-only process/acquisition states, local version/data-source information, and
confirmed exit. Keyboard context menus and double-click setup are supported.
`TaskbarCreated` re-registers the icon without restarting services; a window-only
message filter permits that one recovery message from normal Explorer when the
gateway has elevated integrity.

Browser opening validates a numeric-host/localhost HTTP URL with an explicit
port and a fixed product path, then uses `ShellExecuteW` from a COM STA. It does
not launch a command shell. A bounded one-inflight worker handles each tray
browser open, so a stalled shell handler cannot block quit or timeout handling.
Only safe error notifications are posted to a still-live native window. New-database, exit and force-exit dialogs default to
Cancel. Force exit is a distinct timeout confirmation with incomplete/unknown
and non-durable data warnings. Selected database paths appear only in the native
version dialog. Error dialogs advertise a diagnostic path only when saved.

A native shell failure invokes the fixed OnFault callback before requesting
shutdown, allowing the host to record managed diagnostics and a nonzero result.

A session-shutdown query does not stop services, because another app may veto it.
A committed session end invokes the same coordinator, suppresses force prompts,
and waits up to four seconds for `Options.ShutdownDone`. That channel must close
after service cleanup but before `Shell.Close`, avoiding a UI-thread deadlock.
The OS may still terminate earlier; this path is best effort.

## Builds and native acceptance

Console remains the default. `make build BUILD_MODE=desktop TARGET_GOOS=windows`
or `powershell -File scripts/build.ps1 -Mode desktop` produces
`bin/gateway-desktop.exe`. Console produces `bin/test-ui.exe`. Only the explicit
Windows desktop target sets the `desktop` tag and GUI subsystem. Both builds
require frontend build and embedded-asset synchronization first.

Cross-compilation verifies source/link compatibility only. Run native tests on
Windows and manually verify tray keyboard/overflow behavior, Explorer restart,
modal cancellation/closure, session shutdown, process death, aliases, SID/session
isolation, standalone embedded UI and GUI-headless output/wait/exit behavior in
cmd, PowerShell and the target supervisor. Linux test success is not Windows
runtime acceptance, and no field-device or production-database claim is implied.

## References

- [Windows file locking](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-lockfileex)
- [File identity](https://learn.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-file_id_info)
- [Named pipe security](https://learn.microsoft.com/en-us/windows/win32/ipc/named-pipe-security-and-access-rights)
- [Notification area](https://learn.microsoft.com/en-us/windows/win32/shell/taskbar)
- [Session end](https://learn.microsoft.com/en-us/windows/win32/shutdown/wm-endsession)

- [SQLite multiple-link corruption warning](https://www.sqlite.org/howtocorrupt.html#_multiple_links_to_the_same_file)
