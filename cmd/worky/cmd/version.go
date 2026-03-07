package cmd

import "runtime/debug"

// Version is set at build time via -ldflags "-X github.com/davideimola/worky/cmd/worky/cmd.Version=x.y.z".
// If not set, it falls back to the module version embedded by the Go toolchain (e.g. via go install).
var Version = func() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}()
