package fetch

import (
	"context"
	"fmt"

	"github.com/agentio/slink/gen/xrpc"
	"github.com/agentio/slink/pkg/froda"
	"github.com/agentio/slink/pkg/pretty"
	"github.com/agentio/slink/pkg/resolve"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var loglevel string
	var cursor string
	var limit int64
	var reverse bool
	cmd := &cobra.Command{
		Use:   "fetch ATURI",
		Short: "Fetch the document associated with an AT URI.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(loglevel); err != nil {
				return err
			}
			aturi, err := resolve.ATUriFromString(args[0])
			if err != nil {
				return err
			}
			if err = aturi.ResolveAuthority(); err != nil {
				return err
			}
			if aturi.RKey != "" {
				response, err := FetchRecord(cmd.Context(), aturi)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(response))
			} else if aturi.Collection != "" {
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
			} else {
				pdsurl, err := aturi.ATProtoPDSURL()
				if err != nil {
					// if we can't find a PDS URL, just return the DID doc.
					didDoc, err := resolve.Did(cmd.Context(), aturi.Authority)
					if err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(didDoc))
				}
				// otherwise, return the output of com.atproto.repo.describeRepo
				c := froda.NewClientWithOptions(froda.ClientOptions{
					Host: pdsurl,
				})
				response, err := xrpc.ComATProtoRepoDescribeRepo(cmd.Context(), c, aturi.Authority)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(response))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&loglevel, "log", "l", "warn", "log level (debug, info, warn, error, fatal)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "cursor to use when listing records")
	cmd.Flags().Int64Var(&limit, "limit", 20, "number of records to return in lists")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "reverse list order when listing records")
	return cmd
}

func FetchRecord(ctx context.Context, aturi *resolve.ATURI) (*xrpc.ComATProtoRepoGetRecord_Output, error) {
	pdsurl, err := aturi.ATProtoPDSURL()
	if err != nil {
		return nil, err
	}
	c := froda.NewClientWithOptions(froda.ClientOptions{
		Host: pdsurl,
	})
	return xrpc.ComATProtoRepoGetRecord(ctx, c, "", aturi.Collection, aturi.Authority, aturi.RKey)
}
