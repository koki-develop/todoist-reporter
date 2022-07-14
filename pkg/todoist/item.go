package todoist

import "github.com/koki-develop/todoist-reporter/pkg/util"

type Item struct {
	ID        int      `json:"id"`
	ParentID  *int     `json:"parent_id"`
	ProjectID int      `json:"project_id"`
	SectionID int      `json:"section_id"`
	Content   string   `json:"content"`
	LabelIDs  []int    `json:"labels"`
	Due       *ItemDue `json:"due"`
}

type ItemDue struct {
	Date string `json:"date"`
}

type Items []*Item

func (items Items) FilterByProjectID(id int) Items {
	var rtn Items
	for _, item := range items {
		if item.ProjectID == id {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) FilterByLabelIDs(ids []int) Items {
	var rtn Items
	for _, item := range items {
		for _, id := range ids {
			if util.Contains(item.LabelIDs, id) {
				rtn = append(rtn, item)
				continue
			}
		}
	}
	return rtn
}
