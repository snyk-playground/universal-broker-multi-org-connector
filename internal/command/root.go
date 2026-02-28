package command

import (
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func NewRootCmd(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker-moc",
		Short: "Broker MOC CLI",
		Long: `
A tool to manage and integrate Snyk Universal Broker connections across
multiple Snyk organizations.

The Multi-Org Connector (MOC) simplifies the process of maintaining broker
connections at scale, ensuring consistent setup and reducing manual overhead.
`,

		SilenceUsage: true,
	}
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// commands
	cmd.AddCommand(newCmdGroup(bma))
	cmd.AddCommand(newCmdOrg(bma))

	return cmd
}
