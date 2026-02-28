package command

import (
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func newCmdOrgList(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{}

	return cmd
}
