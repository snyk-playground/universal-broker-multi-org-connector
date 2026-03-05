package command

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
)

func newCmdVersion(bma *app.BrokerMOCApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version information of broker-moc",

		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("broker-moc %v (%v/%v)\n", bma.Version, runtime.GOOS, runtime.GOARCH)
		},
	}
	return cmd
}
