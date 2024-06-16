package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Description: "Print current version",
		Action:      versionCmd,
	}
}

func versionCmd(ctx *cli.Context) error {
	// base := fmt.Sprintf("mountain; version %s; ", appName, version)
	base := ""
	d, ok := debug.ReadBuildInfo()
	fmt.Println(d.Main.Version, d.Main.Sum)
	if !ok {
		fmt.Println(base)
		return nil
	}

	base += fmt.Sprintf("%s; ", d.GoVersion)
	settings := make(map[string]string, len(d.Settings))
	for _, setting := range d.Settings {
		switch setting.Key {
		case "vcs.revision":
			base += fmt.Sprintf("revision %s; ", setting.Value[:8])
		case "vcs.time":
			base += fmt.Sprintf("time %s;", setting.Value)
		default:
			settings[setting.Key] = setting.Value
		}
	}

	fmt.Println(base)
	fmt.Println("build info:")
	for k, v := range settings {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
}
