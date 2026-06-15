package fetch

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
	"github.com/agentio/slink/gen/xrpc"
	"github.com/agentio/slink/pkg/froda"
	"github.com/agentio/slink/pkg/pretty"
	"github.com/agentio/slink/pkg/resolve"
)

func FetchLexicons(ctx context.Context, handle, dir string) error {
	did, err := resolve.Handle(ctx, handle)
	if err != nil {
		return err
	}
	doc, err := resolve.Did(ctx, did)
	if err != nil {
		return err
	}
	var pds string
	for _, service := range doc.Service {
		if service.ID == "#atproto_pds" {
			pds = service.ServiceEndpoint
		}
	}
	var cursor string
	limit := int64(20) // maybe make this configurable
	client := froda.NewClientWithOptions(froda.ClientOptions{Host: pds})
	for {
		results, err := xrpc.ComATProtoRepoListRecords(ctx, client, "com.atproto.lexicon.schema", cursor, limit, handle, false)
		if err != nil {
			return err
		}
		for _, r := range results.Records {
			value := r.Value
			id := value.(map[string]any)["id"].(string)
			path := filepath.Join(dir, strings.ReplaceAll(id, ".", "/")+".json")
			log.Infof("%s", path)
			dir := filepath.Dir(path)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
			err = os.WriteFile(path, []byte(pretty.JSONValue(value)), 0644)
			if err != nil {
				return err
			}
		}
		if results.Cursor == nil {
			break
		}
		cursor = *results.Cursor
	}
	return nil
}
