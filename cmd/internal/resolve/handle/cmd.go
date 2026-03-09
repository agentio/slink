package handle

import (
	"fmt"

	"github.com/agentio/slink/pkg/resolve"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var loglevel string
	cmd := &cobra.Command{
		Use:   "handle HANDLE",
		Short: "Resolve a handle by looking up and returning the corresponding DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(loglevel); err != nil {
				return err
			}
			did, err := resolve.Handle(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", did)
			return nil
		},
	}
	cmd.Flags().StringVarP(&loglevel, "log", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
