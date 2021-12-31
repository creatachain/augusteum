package main

import (
	"os"
	"path/filepath"

	cmd "github.com/creatachain/augusteum/cmd/augusteum/commands"
	"github.com/creatachain/augusteum/cmd/augusteum/commands/debug"
	cfg "github.com/creatachain/augusteum/config"
	"github.com/creatachain/augusteum/libs/cli"
	nm "github.com/creatachain/augusteum/node"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GenValidatorCmd,
		cmd.InitFilesCmd,
		cmd.ProbeUpnpCmd,
		cmd.LightCmd,
		cmd.ReplayCmd,
		cmd.ReplayConsoleCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		cmd.TestnetFilesCmd,
		cmd.ShowNodeIDCmd,
		cmd.GenNodeKeyCmd,
		cmd.VersionCmd,
		debug.DebugCmd,
		cli.NewCompletionCmd(rootCmd, true),
	)

	// NOTE:
	// Users wishing to:
	//	* Use an external signer for their validators
	//	* Supply an in-proc msm app
	//	* Supply a genesis doc file from another source
	//	* Provide their own DB implementation
	// can copy this file and use something other than the
	// DefaultNewNode function
	nodeFunc := nm.DefaultNewNode

	// Create & start node
	rootCmd.AddCommand(cmd.NewRunNodeCmd(nodeFunc))

	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", cfg.DefaultAugusteumDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
