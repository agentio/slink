// The slink command calls APIs described with Lexicon.
package main

import (
	"os"
	"strings"

	"charm.land/log/v2"
	"github.com/agentio/slink/cmd/internal/fetch"
	"github.com/agentio/slink/cmd/internal/generate"
	"github.com/agentio/slink/cmd/internal/remove"
	"github.com/agentio/slink/cmd/internal/resolve"
	"github.com/agentio/slink/cmd/internal/token"
	"github.com/agentio/slink/gen/call"
	"github.com/agentio/slink/gen/check"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	var logLevel string
	cmd := &cobra.Command{
		Use: "slink",
		Long: strings.Join(
			[]string{
				``,
				`"Perhaps we’ve shaken him off at last, the miserable slinker!"`,
				``,
				`A tool for working with the AT Protocol.`,
				``,
				`Environment Variables:`,
				`  SLINK_HOST sets the target host (e.g. "https://public.api.bsky.app").`,
				`  SLINK_AUTH sets the authorization header (e.g. "Bearer XXXX").`,
				`  SLINK_ATPROTOPROXY sets the atproto-proxy header.`,
				`  SLINK_PROXYSESSION sets the proxy-session header (used by IO).`,
				`  SLINK_USERDID sets the user-did header (used by IO).`,
			}, "\n"),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ll, err := log.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			log.SetLevel(ll)
			return nil
		},
	}
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(check.Cmd())
	cmd.AddCommand(fetch.Cmd())
	cmd.AddCommand(generate.Cmd())
	cmd.AddCommand(remove.Cmd())
	cmd.AddCommand(resolve.Cmd())
	cmd.AddCommand(token.Cmd())
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
