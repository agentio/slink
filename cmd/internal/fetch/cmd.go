package fetch

import (
	"github.com/agentio/slink/cmd/internal/fetch/aturi"
	"github.com/agentio/slink/cmd/internal/fetch/lexicons"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fetch",
	}
	cmd.AddCommand(aturi.Cmd())
	cmd.AddCommand(lexicons.Cmd())
	return cmd
}
