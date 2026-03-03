package output

import "time"

// Org is a projection of the Snyk API's org data (snyk.Organization), containing
// only the fields needed for display in CLI commands. This decouples the command's
// output from the underlying SDK's data structures.
type Org struct {
	ID        string    `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name"`
	Slug      string    `json:"slug" yaml:"slug"`
	GroupID   string    `json:"group_id" yaml:"group_id"`
	TenantID  string    `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

type OrgsView struct {
	Orgs []Org `json:"organizations" yaml:"organizations"`
}

func (v OrgsView) Headers() []any {
	return []any{"NAME", "ORGANIZATION ID", "GROUP ID"}
}

func (v OrgsView) Rows() [][]any {
	rows := make([][]any, 0, len(v.Orgs))
	for _, o := range v.Orgs {
		rows = append(rows, []any{
			o.Name,
			o.ID,
			o.GroupID,
		})
	}
	return rows
}
