package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	SlackToken              string `json:"slack_token"`
	SlackChannel            string `json:"slack_channel"`
	TodoistToken            string `json:"todoist_token"`
	TodoistLabelIDs         []int  `json:"todoist_label_ids"`
	TodoistProjectID        int    `json:"todoist_project_id"`
	TodoistWipSectionID     int    `json:"todoist_wip_section_id"`
	TodoistWaitingSectionID int    `json:"todoist_waiting_section_id"`
}

func Load() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return &cfg, nil
}
