package command

import (
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func newCmdOrg(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org <command>",
		Short: "Manage orgs",
		Long:  "Work with Snyk organizations.",
	}

	cmd.AddCommand(newCmdOrgList(bma))

	return cmd
}
