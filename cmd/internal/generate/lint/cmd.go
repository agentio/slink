package lint

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var inputs []string
	var cmd = &cobra.Command{
		Use:   "lint",
		Short: "Check a directory of Lexicon files for possible problems",
		RunE: func(cmd *cobra.Command, args []string) error {
			catalog := lexica.NewCatalog()
			for _, input := range inputs {
				if err := catalog.Load(input, false /* skip lint */); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{"lexicons"}, "input directory")
	return cmd
}
