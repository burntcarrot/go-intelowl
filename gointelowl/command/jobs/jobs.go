package jobs

import (
	"github.com/burntcarrot/gointelowl/command/jobs/list"
	"github.com/spf13/cobra"
)

func NewCmdJobs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Show jobs",
	}

	cmd.AddCommand(list.NewCmdJobsList())
	return cmd
}
