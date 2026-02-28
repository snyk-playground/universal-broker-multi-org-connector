package output

import "time"

// Group is a projection of the Snyk API's group data (snyk.Group), containing
// only the fields needed for display in CLI commands. This decouples
// the command's output from the underlying SDK's data structures.
type Group struct {
	ID        string    `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name"`
	Slug      string    `json:"slug" yaml:"slug"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

type GroupsView struct {
	Groups []Group `json:"groups" yaml:"groups"`
}

func (v GroupsView) Headers() []any {
	return []any{"NAME", "GROUP ID", "SLUG", "CREATED"}
}

func (v GroupsView) Rows() [][]any {
	rows := make([][]any, 0, len(v.Groups))
	for _, g := range v.Groups {
		rows = append(rows, []any{
			g.Name,
			g.ID,
			g.Slug,
			FormatTime(g.CreatedAt),
		})
	}
	return rows
}
