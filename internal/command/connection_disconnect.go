package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
	"github.com/snyk-playground/broker-moc/internal/command/output"
)

type connectionDisconnectOpts struct {
	bma          *app.BrokerMOCApp
	connectionID string
	dryRun       bool
	input        string
	format       string
	output       string
}

func newCmdConnectionDisconnect(bma *app.BrokerMOCApp) *cobra.Command {
	opts := &connectionDisconnectOpts{bma: bma}

	cmd := &cobra.Command{
		Use:   "disconnect <connection-id>",
		Short: "Disconnect connection",
		Long:  "Disconnect integrated broker connection from multiple organizations",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.connectionID = args[0]
			return runConnectionDisconnect(cmd.Context(), opts)
		},
	}
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "print organizations without performing disconnect")
	cmd.Flags().StringVarP(&opts.input, "input", "i", "", "input file with organizations to disconnect (yaml or json)")
	cmd.Flags().StringVarP(&opts.format, "format", "f", "yaml", "output format (json, yaml)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write output to file instead of stdout")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func runConnectionDisconnect(ctx context.Context, opts *connectionDisconnectOpts) error {
	client := opts.bma.APIClient
	log := opts.bma.Logger

	// load input file with integrations and map to the integrations
	var connectedIntegrations output.IntegrationsView
	if err := readIntegrationsFromInputFile(opts.input, &connectedIntegrations); err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	integrationCount := len(connectedIntegrations.Integrations)
	fmt.Println("=> Starting disconnect process...")
	fmt.Println("   integrations: ", integrationCount)
	fmt.Println("   connection ID:", opts.connectionID)
	fmt.Println("------------------------------------------------------------")

	if opts.dryRun {
		fmt.Println()
		fmt.Println("=> NOTE: dry-run flag is set to true, integrations will be printed without performing an disconnect")
		fmt.Println()

		f, _ := output.NewFormatter(output.FormatTable)
		result, err := f.Format(connectedIntegrations)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	}

	i := 0
	tpl := "Disconnect integrations... (%d/%d)"
	s := newSpinner(ctx, fmt.Sprintf(tpl, i, integrationCount))
	if opts.bma.Config.Logging.Level == "debug" {
		s.Stop()
	}

	var integrations []output.Integration
	success := 0
	failure := 0

	for _, connectedIntegration := range connectedIntegrations.Integrations {
		i++
		s.UpdateMessage(fmt.Sprintf(tpl, i, integrationCount))

		log.Debug("Deleting brokered integration",
			"tenant_id", connectedIntegration.TenantID,
			"connection_id", opts.connectionID,
			"org_id", connectedIntegration.OrgID,
			"org_name", connectedIntegration.OrgName,
			"integration_id", connectedIntegration.ID,
		)
		resp, err := client.Brokers.DeleteIntegration(ctx, connectedIntegration.TenantID, opts.connectionID, connectedIntegration.OrgID, connectedIntegration.ID)
		if err != nil {
			snykRequestID := "unknown"
			if resp != nil {
				snykRequestID = resp.SnykRequestID
			}
			log.Debug("Failed to delete brokered integration",
				"tenant_id", connectedIntegration.TenantID,
				"connection_id", opts.connectionID,
				"org_id", connectedIntegration.OrgID,
				"org_name", connectedIntegration.OrgName,
				"integration_id", connectedIntegration.ID,
				"error", err,
				"snyk_request_id", snykRequestID,
			)
			integrations = append(integrations, output.Integration{
				ID:             connectedIntegration.ID,
				ConnectionID:   opts.connectionID,
				ConnectionType: connectedIntegration.ConnectionType,
				OrgID:          connectedIntegration.OrgID,
				OrgName:        connectedIntegration.OrgName,
				TenantID:       connectedIntegration.TenantID,
				Status:         "error",
				ErrorMessage:   err.Error(),
			})
			failure++
			continue
		}
		log.Debug("Deleted brokered integration", "tenant_id", connectedIntegration.TenantID,
			"connection_id", opts.connectionID,
			"org_id", connectedIntegration.OrgID,
			"org_name", connectedIntegration.OrgName,
			"integration_id", connectedIntegration.ID,
		)

		// write success here
		integrations = append(integrations, output.Integration{
			ID:             connectedIntegration.ID,
			ConnectionID:   connectedIntegration.ConnectionID,
			ConnectionType: connectedIntegration.ConnectionType,
			OrgID:          connectedIntegration.OrgID,
			OrgName:        connectedIntegration.OrgName,
			TenantID:       connectedIntegration.TenantID,
			Status:         "disconnected",
		})
		success++
	}
	s.Stop()

	f, err := output.NewFormatter(output.Format(opts.format))
	if err != nil {
		return err
	}
	result, err := f.Format(output.IntegrationsView{Integrations: integrations})
	if err != nil {
		return err
	}

	fmt.Println("=> Summary")
	fmt.Println("   successful disconnections:", success)
	fmt.Println("   failed disconnections:    ", failure)
	fmt.Println("------------------------------------------------------------")

	if opts.output != "" {
		return os.WriteFile(opts.output, []byte(result), 0644)
	}
	fmt.Println(result)

	return nil
}

func readIntegrationsFromInputFile(path string, view *output.IntegrationsView) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return json.Unmarshal(data, view)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, view)
	default:
		return fmt.Errorf("unsupported file format: %s (expected 'json', 'yaml' or 'yml')", ext)
	}
}
