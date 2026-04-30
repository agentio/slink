package resolve

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"charm.land/log/v2"
)

func Handle(ctx context.Context, handle string) (string, error) {
	// first check DNS
	records, err := net.LookupTXT(fmt.Sprintf("_atproto.%s", handle))
	if err == nil {
		for _, rec := range records {
			if strings.HasPrefix(rec, "did=") {
				did := strings.Split(rec, "did=")[1]
				if did != "" {
					log.Info("Found DID with DNS")
					return did, nil
				}
			}
		}
	}
	// if that didn't work, check the .well-known/atproto-did path
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://%s/.well-known/atproto-did", handle),
		nil,
	)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return "", fmt.Errorf("unable to resolve %s", handle)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
