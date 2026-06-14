package check

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var inputs []string
	var output string
	var manifest string
	var cmd = &cobra.Command{
		Use:   "check",
		Short: "Generate a command-line interface to check records in a directory of Lexicon files",
		RunE: func(cmd *cobra.Command, args []string) error {
			catalog := lexica.NewCatalog()
			for _, input := range inputs {
				if err := catalog.Load(input, false /* skip lint */); err != nil {
					return err
				}
			}
			if manifest != "" {
				_, err := lexica.BuildManifest(manifest)
				if err != nil {
					return err
				}
			}
			if err := catalog.GenerateCheckCommands(output); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{"lexicons"}, "input directory")
	cmd.Flags().StringVarP(&output, "output", "o", "gen/check", "output directory")
	cmd.Flags().StringVarP(&manifest, "manifest", "m", "", "manifest")
	return cmd
}
