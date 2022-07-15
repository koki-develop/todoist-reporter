package reporter

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/koki-develop/todoist-reporter/pkg/todoist"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type slackAPI interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

type Reporter struct {
	slackAPI slackAPI
}

func New(slackToken string) *Reporter {
	return &Reporter{
		slackAPI: slack.New(slackToken),
	}
}

func (r *Reporter) ReportDaily(channel string, completed, wip, waiting todoist.Items, labels todoist.Labels) error {
	var blocks []slack.Block

	// 完了したタスク
	blocks = append(blocks, r.newHeaderBlock(":white_check_mark: 完了したタスク"))
	for labelname, items := range r.groupItemsByLabel(completed, labels) {
		if len(items) == 0 {
			continue
		}

		tpl, err := template.New("done").Parse("*{{ .label }}*\n{{ range .items }}• {{ .Content }}\n{{ end }}")
		if err != nil {
			return errors.WithStack(err)
		}
		b := new(bytes.Buffer)
		if err := tpl.Execute(b, map[string]interface{}{"label": labelname, "items": items}); err != nil {
			return errors.WithStack(err)
		}
		blocks = append(blocks, r.newMarkdownBlock(b.String()))
	}
	blocks = append(blocks, slack.NewDividerBlock())

	// 進行中のタスク
	blocks = append(blocks, r.newHeaderBlock(":man-running: 進行中のタスク"))
	for labelname, items := range r.groupItemsByLabel(wip, labels) {
		if len(items) == 0 {
			continue
		}
		rows := []string{fmt.Sprintf("*%s*", labelname)}
		r.appendWipItemsToRows(&rows, items, 0)
		blocks = append(blocks, r.newMarkdownBlock(strings.Join(rows, "\n")))
	}
	blocks = append(blocks, slack.NewDividerBlock())

	// 待ちのタスク
	blocks = append(blocks, r.newHeaderBlock(":clock3: 待ちのタスク"))
	for labelname, items := range r.groupItemsByLabel(waiting, labels) {
		if len(items) == 0 {
			continue
		}

		tpl, err := template.New("waiting").Parse("*{{ .label }}*\n{{ range .items }}• {{ .Content }}\n{{ end }}")
		if err != nil {
			return errors.WithStack(err)
		}
		b := new(bytes.Buffer)
		if err := tpl.Execute(b, map[string]interface{}{"label": labelname, "items": items}); err != nil {
			return errors.WithStack(err)
		}
		blocks = append(blocks, r.newMarkdownBlock(b.String()))

	}
	blocks = append(blocks, slack.NewDividerBlock())

	// 投稿
	if _, _, err := r.slackAPI.PostMessage(channel, slack.MsgOptionText("デイリーレポート", false), slack.MsgOptionBlocks(blocks...)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Reporter) newHeaderBlock(text string) *slack.HeaderBlock {
	return slack.NewHeaderBlock(&slack.TextBlockObject{Type: slack.PlainTextType, Text: text})
}

func (r *Reporter) newMarkdownBlock(text string) *slack.SectionBlock {
	return &slack.SectionBlock{
		Type: slack.MBTSection,
		Text: &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		},
	}
}

func (r *Reporter) groupItemsByLabel(items todoist.Items, labels todoist.Labels) map[string]todoist.Items {
	m := map[string]todoist.Items{}
	for _, label := range labels {
		m[label.Name] = items.FilterByLabelID(label.ID)
	}
	return m
}

func (r *Reporter) appendWipItemsToRows(rows *[]string, items todoist.Items, depth int) {
	prefix := strings.Repeat("\t", depth)
	for _, item := range items {
		if item.Completed {
			*rows = append(*rows, fmt.Sprintf("%s• (完了) ~%s~", prefix, item.Content))
		} else {
			*rows = append(*rows, fmt.Sprintf("%s• %s", prefix, item.Content))
		}
		if len(item.Children) > 0 {
			r.appendWipItemsToRows(rows, item.Children, depth+1)
		}
	}
}
