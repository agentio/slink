package resolve

import (
	"github.com/agentio/slink/cmd/internal/resolve/did"
	"github.com/agentio/slink/cmd/internal/resolve/handle"
	"github.com/agentio/slink/cmd/internal/resolve/now"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve atproto identifiers",
	}
	cmd.AddCommand(handle.Cmd())
	cmd.AddCommand(did.Cmd())
	cmd.AddCommand(now.Cmd())
	return cmd
}
