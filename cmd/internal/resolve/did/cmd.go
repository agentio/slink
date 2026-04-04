package did

import (
	"fmt"

	"github.com/agentio/slink/pkg/pretty"
	"github.com/agentio/slink/pkg/resolve"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var loglevel string
	var pds bool
	cmd := &cobra.Command{
		Use:   "did DID",
		Short: "Resolve a DID by looking up and returning the corresponding document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(loglevel); err != nil {
				return err
			}
			didDocument, err := resolve.Did(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if pds {
				for _, s := range didDocument.Service {
					if s.ID == "#atproto_pds" {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n", s.ServiceEndpoint)
						return nil
					}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", pretty.JSONValue(didDocument))
			return nil
		},
	}
	cmd.Flags().StringVarP(&loglevel, "log", "l", "warn", "log level (debug, info, warn, error, fatal)")
	cmd.Flags().BoolVar(&pds, "pds", false, "just return the pds of the DID")
	return cmd
}
