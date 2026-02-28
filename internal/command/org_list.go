package command

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
	"github.com/snyk-playground/broker-moc/internal/command/output"
)

type orgListOpts struct {
	bma     *app.BrokerMOCApp
	groupID string
	format  string
	output  string
}

func newCmdOrgList(bma *app.BrokerMOCApp) *cobra.Command {
	opts := &orgListOpts{bma: bma}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List orgs",
		Long:  "List all accessible organizations.",

		Aliases: []string{"ls"},

		RunE: func(cmd *cobra.Command, _ []string) error {
			return runOrgList(cmd.Context(), opts)
		},
	}
	cmd.Flags().StringVarP(&opts.groupID, "group-id", "", "", "filter organization by group id")
	cmd.Flags().StringVarP(&opts.format, "format", "f", "table", "output format (json, yaml, table)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write output to file instead of stdout")

	return cmd
}

func runOrgList(ctx context.Context, opts *orgListOpts) error {
	client := opts.bma.APIClient
	log := opts.bma.Logger

	var orgs []output.Org

	s := newSpinner(ctx, "Querying for accessible organizations...")
	if opts.bma.Config.Logging.Level == "debug" {
		s.Stop()
	}

	log.Debug("Listing accessible organizations")
	orgsAPI, errf := client.Orgs.AllAccessibleOrgs(ctx, nil)
	for orgAPI, resp := range orgsAPI {
		log.Debug("Got organization", "org", orgAPI, "snyk_request_id", resp.SnykRequestID)

		if opts.groupID != "" {
			if orgAPI.Attributes.GroupID != opts.groupID {
				log.Debug("Skip organization because it's not part of the group", "group_id", opts.groupID)
				continue
			}
		}

		orgs = append(orgs, output.Org{
			ID:        orgAPI.ID,
			Name:      orgAPI.Attributes.Name,
			Slug:      orgAPI.Attributes.Slug,
			GroupID:   orgAPI.Attributes.GroupID,
			CreatedAt: orgAPI.Attributes.CreatedAt,
		})
	}
	if err := errf(); err != nil {
		s.Stop()
		return fmt.Errorf("unable to list accessible organizations: %w", err)
	}
	s.Stop()

	f, err := output.NewFormatter(output.Format(opts.format))
	if err != nil {
		return err
	}
	result, err := f.Format(output.OrgsView{Orgs: orgs})
	if err != nil {
		return err
	}

	if opts.output != "" {
		return os.WriteFile(opts.output, []byte(result), 0644)
	}
	fmt.Println(result)

	return nil
}
