package todoist

import "github.com/koki-develop/todoist-reporter/pkg/util"

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Labels []*Label

func (labels Labels) FindByID(id int) *Label {
	for _, label := range labels {
		if label.ID == id {
			return label
		}
	}
	return nil
}

func (labels Labels) FilterByIDs(ids []int) Labels {
	var rtn Labels
	for _, label := range labels {
		if util.Contains(ids, label.ID) {
			rtn = append(rtn, label)
		}
	}
	return rtn
}
