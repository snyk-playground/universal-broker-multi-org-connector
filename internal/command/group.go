package command

import (
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func newCmdGroup(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group <command>",
		Short: "Manage groups",
		Long:  "Work with Snyk groups.",
	}

	cmd.AddCommand(newCmdGroupList(bma))

	return cmd
}
