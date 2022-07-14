package todoist

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Labels []*Label
