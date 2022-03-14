package command

import (
	"github.com/burntcarrot/gointelowl/command/jobs"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gointelowl <command> <subcommand> [flags]",
		Short: "IntelOwl CLI",
	}

	cmd.AddCommand(jobs.NewCmdJobs())

	return cmd
}

func Execute() error {
	cmd := NewCmdRoot()
	return cmd.Execute()
}
