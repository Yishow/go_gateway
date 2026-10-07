//go:build windows && desktop

package desktop

// DesktopBuildDefault distinguishes the artifact default from explicit runtime flags.
//
//nolint:revive // The name distinguishes the artifact default from runtime mode.
func DesktopBuildDefault() bool { return true }
