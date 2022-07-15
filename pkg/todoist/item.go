package todoist

import "github.com/koki-develop/todoist-reporter/pkg/util"

type Item struct {
	ID        int    `json:"id"`
	ParentID  *int   `json:"parent_id"`
	ProjectID int    `json:"project_id"`
	SectionID int    `json:"section_id"`
	Content   string `json:"content"`
	LabelIDs  []int  `json:"labels"`
	Completed bool

	Children Items
}

type Items []*Item

func (items Items) FilterOnlyIncompleted() Items {
	var rtn Items
	for _, item := range items {
		if !item.Completed {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) FilterOnlyCompleted() Items {
	var rtn Items
	for _, item := range items {
		if item.Completed {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) FilterOnlyRoot() Items {
	var rtn Items
	for _, item := range items {
		if item.ParentID == nil {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) FilterByProjectID(id int) Items {
	var rtn Items
	for _, item := range items {
		if item.ProjectID == id {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) FilterByLabelID(id int) Items {
	var rtn Items
	for _, item := range items {
		if util.Contains(item.LabelIDs, id) {
			rtn = append(rtn, item)
			continue
		}
	}
	return rtn
}

func (items Items) FilterByLabelIDs(ids []int) Items {
	var rtn Items
	for _, item := range items {
		if util.ContainsAny(item.LabelIDs, ids) {
			rtn = append(rtn, item)
			continue
		}
	}
	return rtn
}

func (items Items) FilterBySectionID(id int) Items {
	var rtn Items
	for _, item := range items {
		if item.SectionID == id {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (items Items) Organize() {
	for _, item := range items {
		for _, child := range items {
			if child.ParentID == nil {
				continue
			}
			if *child.ParentID != item.ID {
				continue
			}
			item.Children = append(item.Children, child)
		}
	}

	items = items.FilterOnlyRoot()
}
