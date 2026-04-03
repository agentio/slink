package resolve

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

type ATURI struct {
	Authority  string
	Collection string
	RKey       string
}

func ATUriFromString(s string) (*ATURI, error) {
	a, err := recordATUriFromString(s)
	if err == nil {
		return a, nil
	}
	a, err = authorityATUriFromString(s)
	if err == nil {
		return a, nil
	}
	return nil, errors.New("invalid at:// uri")
}

func recordATUriFromString(s string) (*ATURI, error) {
	re := regexp.MustCompile(`^at://([a-zA-Z0-9._:-]+)/([a-zA-Z0-9.-]+)/([a-zA-Z0-9._~%-]+)$`)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return nil, errors.New("invalid at:// record uri")
	}
	authority := m[1]
	collection := m[2]
	rkey := m[3]
	return &ATURI{
		Authority:  authority,
		Collection: collection,
		RKey:       rkey,
	}, nil
}

func authorityATUriFromString(s string) (*ATURI, error) {
	re := regexp.MustCompile(`^at://([a-zA-Z0-9._:-]+)$`)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return nil, errors.New("invalid at:// authority uri")
	}
	authority := m[1]
	return &ATURI{
		Authority: authority,
	}, nil
}

func (a *ATURI) ResolveAuthority() error {
	if strings.HasPrefix(a.Authority, "did:") {
		return nil
	}
	did, err := Handle(context.TODO(), a.Authority)
	if err != nil {
		return err
	}
	a.Authority = did
	return nil
}

func (a *ATURI) ATProtoPDSURL() (string, error) {
	did := a.Authority
	didDoc, err := Did(context.TODO(), did)
	if err != nil {
		return "", err
	}
	for _, s := range didDoc.Service {
		if s.ID == "#atproto_pds" {
			return s.ServiceEndpoint, nil
		}
	}
	return "", errors.New("no atproto_pds is defined")
}
