package list

import (
	"context"
	"fmt"

	"github.com/burntcarrot/gointelowl/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type JobListOptions struct {
	jobData []internal.Job
	table   [][]string
}

func NewCmdJobsList() *cobra.Command {
	opts := JobListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	return cmd
}

func (opts *JobListOptions) Run() error {
	ctx := context.Background()
	err := opts.getData(ctx)
	if err != nil {
		return err
	}

	opts.showJobs()

	return nil
}

func (opts *JobListOptions) getData(ctx context.Context) error {
	client, err := internal.NewClient("testAPIKey")
	if err != nil {
		return err
	}

	opts.jobData, err = client.GetJobs(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (opts *JobListOptions) showJobs() {
	opts.table = make([][]string, len(opts.jobData)+1)

	// add header
	opts.table[0] = []string{"ID", "Name", "Type", "Tags", "Analyzers Called", "Process Time", "Status"}

	// populate table
	for index, job := range opts.jobData {
		id := fmt.Sprint(job.ID)
		observableName := job.ObservableName
		observableClassification := job.ObservableClassification
		analyzersExecuted := job.NoOfAnalyzersExecuted
		processTime := fmt.Sprint(job.ProcessTime)
		status := job.Status
		tags := job.Tags

		tagString := ""
		for _, t := range tags {
			tagString += t.Label + ", "
		}

		opts.table[index+1] = []string{id, observableName, observableClassification, tagString, analyzersExecuted, processTime, status}
	}

	pterm.DefaultTable.WithHasHeader().WithSeparator("\t").WithData(opts.table).Render()
}
