package remove

import (
	"context"
	"errors"
	"fmt"

	"github.com/agentio/slink/gen/xrpc"
	"github.com/agentio/slink/pkg/froda"
	"github.com/agentio/slink/pkg/pretty"
	"github.com/agentio/slink/pkg/resolve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cursor string
	var limit int64
	var reverse bool
	cmd := &cobra.Command{
		Use:   "remove ATURI",
		Short: "Remove documents matching an AT URI.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			aturi, err := resolve.ATUriFromString(args[0])
			if err != nil {
				return err
			}
			if err = aturi.ResolveAuthority(); err != nil {
				return err
			}
			if aturi.RKey != "" {
				// Remove a specified record.
				response, err := deleteRecord(cmd.Context(), aturi)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(response))
			} else if aturi.Collection != "" {
				// Remove all records that are returned by listing a collection with the flags below.
				pdsurl, err := aturi.ATProtoPDSURL()
				if err != nil {
					return err
				}
				c := froda.NewClientWithOptions(froda.ClientOptions{
					Host: pdsurl,
				})
				response, err := xrpc.ComATProtoRepoListRecords(cmd.Context(), c,
					aturi.Collection,
					cursor,
					limit,
					aturi.Authority,
					reverse)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(response))
				for _, record := range response.Records {
					record_aturi, err := resolve.ATUriFromString(record.Uri)
					if err != nil {
						return err
					}
					response, err := deleteRecord(cmd.Context(), record_aturi)
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(response))
				}
			} else {
				return errors.New("this tool is unable to remove an account")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cursor, "cursor", "", "cursor to use when listing records")
	cmd.Flags().Int64Var(&limit, "limit", 20, "number of records to return in lists")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "reverse list order when listing records")
	return cmd
}

func deleteRecord(ctx context.Context, aturi *resolve.ATURI) (*xrpc.ComATProtoRepoDeleteRecord_Output, error) {
	pdsurl, err := aturi.ATProtoPDSURL()
	if err != nil {
		return nil, err
	}
	c := froda.NewClientWithOptions(froda.ClientOptions{
		Host: pdsurl,
	})
	return xrpc.ComATProtoRepoDeleteRecord(ctx, c, &xrpc.ComATProtoRepoDeleteRecord_Input{
		Collection: aturi.Collection,
		Repo:       aturi.Authority,
		Rkey:       aturi.RKey,
	})
}
