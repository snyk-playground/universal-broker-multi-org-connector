package command

import (
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func newCmdConnection(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection <command>",
		Short: "Manage connections",
		Long:  "Work with Broker connections",
	}

	cmd.AddCommand(newCmdConnectionIntegrate(bma))
	cmd.AddCommand(newCmdConnectionDisconnect(bma))

	return cmd
}
