package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"charm.land/log/v2"
)

type DidDocument struct {
	ID                 string `json:"id"`
	VerificationMethod []struct {
		ID                 string `json:"id"`
		Type               string `json:"type"`
		Controller         string `json:"controller"`
		PublicKeyMultibase string `json:"publicKeyMultibase"`
	} `json:"verificationMethod"`
	Service []struct {
		ID              string `json:"id"`
		Type            string `json:"type"`
		ServiceEndpoint string `json:"serviceEndpoint"`
	} `json:"service"`
	AlsoKnownAs []string `json:"alsoKnownAs"`
}

func DidBytes(ctx context.Context, did string) ([]byte, error) {
	var url string
	if strings.HasPrefix(did, "did:plc:") {
		url = fmt.Sprintf("https://plc.directory/%s", did)
	} else if strings.HasPrefix(did, "did:web:") {
		url = fmt.Sprintf("https://%s/.well-known/did.json", strings.TrimPrefix(did, "did:web:"))
	} else {
		return nil, fmt.Errorf("%s is not a valid did", did)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("%s is not in the PLC registry", did)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("%s", string(b))
	return b, nil
}

func Did(ctx context.Context, did string) (*DidDocument, error) {
	b, err := DidBytes(ctx, did)
	if err != nil {
		return nil, err
	}
	var didDocument DidDocument
	err = json.Unmarshal(b, &didDocument)
	if err != nil {
		return nil, err
	}
	return &didDocument, nil
}
