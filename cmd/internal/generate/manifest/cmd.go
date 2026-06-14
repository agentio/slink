package manifest

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var inputs []string
	var cmd = &cobra.Command{
		Use:   "manifest MANIFEST",
		Short: "Generate a list of dependencies in a Lexicon",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalog := lexica.NewCatalog()
			for _, input := range inputs {
				if err := catalog.Load(input, false /* skip lint */); err != nil {
					return err
				}
			}
			m, err := lexica.BuildManifest(args[0])
			if err != nil {
				return err
			}
			slink.Write(cmd.OutOrStdout(), "-", m)
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{"lexicons"}, "input directory")
	return cmd
}
