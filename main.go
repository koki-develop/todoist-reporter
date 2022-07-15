package main

import (
	"fmt"
	"os"
	"time"

	"github.com/koki-develop/todoist-reporter/pkg/config"
	"github.com/koki-develop/todoist-reporter/pkg/reporter"
	"github.com/koki-develop/todoist-reporter/pkg/todoist"
)

func fatal(err error) {
	fmt.Printf("%+v\n", err)
	os.Exit(1)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fatal(err)
	}

	cl := todoist.New(cfg.TodoistToken)

	r, err := cl.GetResources([]string{"items", "labels"})
	if err != nil {
		fatal(err)
	}

	compitems, err := cl.GetCompletedItems(cfg.TodoistProjectID, time.Now().AddDate(0, 0, -1))
	if err != nil {
		fatal(err)
	}
	r.Items = append(r.Items, compitems...)

	r.Items.Organize()
	r.Items = r.Items.FilterByProjectID(cfg.TodoistProjectID)
	r.Items = r.Items.FilterByLabelIDs(cfg.TodoistLabelIDs)
	r.Labels = r.Labels.FilterByIDs(cfg.TodoistLabelIDs)

	completed := r.Items.FilterOnlyCompleted()
	wip := r.Items.FilterBySectionID(cfg.TodoistWipSectionID).FilterOnlyIncompleted()
	waiting := r.Items.FilterBySectionID(cfg.TodoistWaitingSectionID).FilterOnlyIncompleted().FilterOnlyRoot()

	rpt := reporter.New(cfg.SlackToken)
	if err := rpt.ReportDaily(cfg.SlackChannel, completed, wip, waiting, r.Labels); err != nil {
		fatal(err)
	}
	return
}
