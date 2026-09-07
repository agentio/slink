package generate

import (
	"fmt"
	"os"

	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var keyfile string
	var iss string
	var sub string
	var aud []string
	var lxm string
	var htm string
	var htu string
	var typ string
	var nonce bool
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a signed authorization token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			keybytes, err := os.ReadFile(keyfile)
			if err != nil {
				return err
			}
			claims := map[string]any{}
			if len(aud) > 0 {
				claims["aud"] = aud
			}
			if iss != "" {
				claims["iss"] = iss
			}
			if sub != "" {
				claims["sub"] = sub
			}
			if htm != "" {
				claims["htm"] = htm
			}
			if htu != "" {
				claims["htu"] = htu
			}
			if lxm != "" {
				claims["lxm"] = lxm
			}
			if nonce {
				claims["nonce"] = slink.NewJwtID()
			}
			tok, err := slink.GenerateAuthToken(keybytes, claims, typ)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(tok))
			return nil
		},
	}
	cmd.Flags().StringVar(&keyfile, "key", "", "key file")
	cmd.Flags().StringArrayVar(&aud, "aud", nil, "audience")
	cmd.Flags().StringVar(&iss, "iss", "", "issuer")
	cmd.Flags().StringVar(&sub, "sub", "", "subject")
	cmd.Flags().StringVar(&lxm, "lxm", "", "lexicon method")
	cmd.Flags().StringVar(&htm, "htm", "", "http method")
	cmd.Flags().StringVar(&htu, "htu", "", "http url")
	cmd.Flags().StringVar(&typ, "typ", "jwt", "type")
	cmd.Flags().BoolVar(&nonce, "nonce", false, "include a nonce")
	return cmd
}
