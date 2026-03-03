package command

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/snyk-playground/broker-moc/internal/app"
	"github.com/snyk-playground/broker-moc/internal/command/output"
)

type groupListOpts struct {
	bma    *app.BrokerMOCApp
	format string
	output string
}

func newCmdGroupList(bma *app.BrokerMOCApp) *cobra.Command {
	opts := &groupListOpts{bma: bma}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Long:  "List all available groups which a user is a member of.",

		Aliases: []string{"ls"},

		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGroupList(cmd.Context(), opts)
		},
	}
	cmd.Flags().StringVarP(&opts.format, "format", "f", "table", "output format (json, yaml, table)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write output to file instead of stdout")

	return cmd
}

func runGroupList(ctx context.Context, opts *groupListOpts) error {
	client := opts.bma.APIClient
	log := opts.bma.Logger

	var groups []output.Group

	s := newSpinner(ctx, "Querying for groups...")
	if opts.bma.Config.Logging.Level == "debug" {
		s.Stop()
	}

	// fetch groups from API and converts group API objects to format objects
	log.Debug("Listing all available groups")
	groupsAPI, errf := client.Groups.All(ctx, nil)
	for groupAPI := range groupsAPI {
		groupID := groupAPI.ID
		groupAPI, resp, err := client.Groups.Get(ctx, groupID)
		if err != nil {
			s.Stop()
			return fmt.Errorf("unable to get group with id '%s': %w", groupID, err)
		}
		log.Debug("Got group with enriched properties", "group", groupAPI, "snyk_request_id", resp.SnykRequestID)

		groups = append(groups, output.Group{
			ID:        groupAPI.ID,
			Name:      groupAPI.Attributes.Name,
			Slug:      groupAPI.Attributes.Slug,
			CreatedAt: groupAPI.Attributes.CreatedAt,
		})
	}
	if err := errf(); err != nil {
		s.Stop()
		return fmt.Errorf("unable to list all groups: %w", err)
	}
	s.Stop()

	f, err := output.NewFormatter(output.Format(opts.format))
	if err != nil {
		return err
	}
	result, err := f.Format(output.GroupsView{Groups: groups})
	if err != nil {
		return err
	}

	if opts.output != "" {
		return os.WriteFile(opts.output, []byte(result), 0600)
	}
	fmt.Println(result)

	return nil
}
