package output

// Integration is a combination of a projection of the Snyk API's integration data
// (snyk.BrokerIntegration), containing only the fields needed for display in CLI commands,
// and operation results during integrate/disconnect actions.
type Integration struct {
	ID             string `json:"id" yaml:"id"`
	ConnectionID   string `json:"connection_id" yaml:"connection_id"`
	ConnectionType string `json:"connection_type" yaml:"connection_type"`
	OrgID          string `json:"organization_id" yaml:"organization_id"`
	OrgName        string `json:"organization_name" yaml:"organization_name"`
	TenantID       string `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`

	Status       string `json:"status" yaml:"status"`
	ErrorMessage string `json:"error_message,omitempty" yaml:"error_message,omitempty"`
}

type IntegrationsView struct {
	Integrations []Integration `json:"integrations" yaml:"integrations"`
}

func (v IntegrationsView) Headers() []any {
	return []any{"ID", "ORGANIZATION NAME", "CONNECTION ID"}
}

func (v IntegrationsView) Rows() [][]any {
	rows := make([][]any, 0, len(v.Integrations))
	for _, i := range v.Integrations {
		rows = append(rows, []any{
			i.ID,
			i.OrgName,
			i.ConnectionID,
		})
	}
	return rows
}
