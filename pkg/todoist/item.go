package todoist

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
