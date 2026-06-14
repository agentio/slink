package handle

import (
	"fmt"

	"github.com/agentio/slink/pkg/resolve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle HANDLE",
		Short: "Resolve a handle by looking up and returning the corresponding DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			did, err := resolve.Handle(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", did)
			return nil
		},
	}
	return cmd
}
