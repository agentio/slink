package lexicons

import (
	"github.com/agentio/slink/pkg/fetch"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "lexicons HANDLE",
		Short: "Fetch all of the lexicons associated with a handle to a local directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fetch.FetchLexicons(cmd.Context(), args[0], dir)
		},
	}
	cmd.Flags().StringVarP(&dir, "directory", "d", "lexicons-local", "local lexicons directory")
	return cmd
}
