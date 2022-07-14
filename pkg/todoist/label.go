package todoist

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
