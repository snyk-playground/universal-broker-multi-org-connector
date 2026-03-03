package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/snyk-playground/broker-moc/internal/app"
	"github.com/snyk-playground/broker-moc/internal/command/output"
)

type connectionIntegrateOpts struct {
	bma            *app.BrokerMOCApp
	connectionID   string
	connectionType string
	dryRun         bool
	input          string
	format         string
	output         string
	tenantID       string
}

func newCmdConnectionIntegrate(bma *app.BrokerMOCApp) *cobra.Command {
	opts := &connectionIntegrateOpts{bma: bma}

	cmd := &cobra.Command{
		Use:   "integrate <connection-id>",
		Short: "Integrate connection",
		Long:  "Integrate broker connection against multiple organizations",
		Args:  cobra.ExactArgs(1),

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// bind tenant id from flag or environment variable
			if err := viper.BindPFlag("tenant-id", cmd.Flags().Lookup("tenant-id")); err != nil {
				return err
			}
			if err := viper.BindEnv("tenant-id", "SNYK_TENANT_ID"); err != nil {
				return err
			}

			opts.tenantID = viper.GetString("tenant-id")

			// validate if it's empty or not
			if opts.tenantID == "" {
				return fmt.Errorf("required flag(s) \"tenant-id\" not set (use --tenant-id or SNYK_TENANT_ID environment variable)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.connectionID = args[0]
			return runConnectionIntegrate(cmd.Context(), opts)
		},
	}
	cmd.Flags().StringVar(&opts.connectionType, "connection-type", "", "integration type (e.g. github, gitlab)")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "print organizations without performing integration")
	cmd.Flags().StringVarP(&opts.input, "input", "i", "", "input file with organizations to integrate (yaml or json)")
	cmd.Flags().StringVarP(&opts.format, "format", "f", "yaml", "output format (json, yaml)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write output to file instead of stdout")
	cmd.Flags().StringVar(&opts.tenantID, "tenant-id", "", "tenant id")
	_ = cmd.MarkFlagRequired("connection-type")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func runConnectionIntegrate(ctx context.Context, opts *connectionIntegrateOpts) error {
	client := opts.bma.APIClient
	log := opts.bma.Logger

	// load input file with orgs and map to the orgs
	var orgs output.OrgsView
	if err := readOrganizationsFromInputFile(opts.input, &orgs); err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	orgCount := len(orgs.Orgs)
	fmt.Println("=> Starting integration process...")
	fmt.Println("   organizations:  ", orgCount)
	fmt.Println("   tenant ID:      ", opts.tenantID)
	fmt.Println("   connection ID:  ", opts.connectionID)
	fmt.Println("   connection type:", opts.connectionType)
	fmt.Println("------------------------------------------------------------")

	if opts.dryRun {
		fmt.Println()
		fmt.Println("=> NOTE: dry-run flag is set to true, organizations will be printed without performing an integration")
		fmt.Println()

		f, _ := output.NewFormatter(output.FormatTable)
		result, err := f.Format(orgs)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	}

	i := 0
	tpl := "Integrate organizations... (%d/%d)"
	s := newSpinner(ctx, fmt.Sprintf(tpl, i, orgCount))
	if opts.bma.Config.Logging.Level == "debug" {
		s.Stop()
	}

	var integrations []output.Integration
	success := 0
	failure := 0

	for _, org := range orgs.Orgs {
		i++
		s.UpdateMessage(fmt.Sprintf(tpl, i, orgCount))

		log.Debug("Creating brokered integration", "connection_id", opts.connectionID, "connection_type", opts.connectionType, "org_id", org.ID)
		brokerIntegration, _, err := client.Brokers.CreateIntegration(ctx, opts.tenantID, opts.connectionID, org.ID, &snyk.BrokerIntegrationCreateRequest{
			Type: snyk.BrokerConnectionType(opts.connectionType),
		})
		if err != nil {
			log.Debug("Failed to create brokered integration",
				"connection_id", opts.connectionID,
				"connection_type", opts.connectionType,
				"org_id", org.ID,
				"org_name", org.Name,
				"error", err,
			)
			integrations = append(integrations, output.Integration{
				ID:             "<none>",
				ConnectionID:   opts.connectionID,
				ConnectionType: opts.connectionType,
				OrgID:          org.ID,
				OrgName:        org.Name,
				TenantID:       org.TenantID,
				Status:         "error",
				ErrorMessage:   err.Error(),
			})
			failure++
			continue
		}
		log.Debug("Created brokered integration", "broker_integration", brokerIntegration)

		// write success here
		integrations = append(integrations, output.Integration{
			ID:             brokerIntegration.ID,
			ConnectionID:   opts.connectionID,
			ConnectionType: opts.connectionType,
			OrgID:          brokerIntegration.OrgID,
			OrgName:        org.Name,
			TenantID:       org.TenantID,
			Status:         "connected",
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
	fmt.Println("   successful integrations: ", success)
	fmt.Println("   failed integrations:     ", failure)
	fmt.Println("------------------------------------------------------------")

	if opts.output != "" {
		return os.WriteFile(opts.output, []byte(result), 0600)
	}
	fmt.Println(result)

	return nil
}

func readOrganizationsFromInputFile(path string, view *output.OrgsView) error {
	safePath := filepath.Clean(path)
	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(safePath))
	switch ext {
	case ".json":
		return json.Unmarshal(data, view)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, view)
	default:
		return fmt.Errorf("unsupported file format: %s (expected 'json', 'yaml' or 'yml')", ext)
	}
}
