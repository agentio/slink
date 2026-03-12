package slink

import (
	"encoding/base64"
	"encoding/json"
)

// Blob represents the Lexicon "blob" type (https://atproto.com/specs/data-model#blob-type).
type Blob struct {
	LexiconTypeID string `json:"$type,omitempty"`
	Ref           Link   `json:"ref,omitempty"`
	MimeType      string `json:"mimeType,omitempty"`
	Size          int64  `json:"size"`
	Cid           string `json:"cid,omitempty"` // deprecated legacy blob format (see docs linked above).
}

// Link represents the Lexicon "cid-link" type (https://atproto.com/specs/data-model#link).
type Link struct {
	LexiconLink string `json:"$link"`
}

// Bytes represents the Lexicon "bytes" type (https://atproto.com/specs/data-model#bytes).
type Bytes struct {
	Bytes []byte
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	s := base64.RawStdEncoding.EncodeToString(b.Bytes)
	v := map[string]any{
		"$bytes": s,
	}
	return json.Marshal(v)
}

func (b *Bytes) UnmarshalJSON(data []byte) (err error) {
	var v struct {
		Bytes string `json:"$bytes"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	b.Bytes, err = base64.RawStdEncoding.DecodeString(v.Bytes)
	return
}

// LexiconTypeFromJSONBytes extracts the lexicon type from an otherwise-unparsed value.
func LexiconTypeFromJSONBytes(data []byte) string {
	type TypedRecord struct {
		LexiconTypeID string `json:"$type"`
	}
	var record TypedRecord
	err := json.Unmarshal(data, &record)
	if err != nil {
		return ""
	}
	return record.LexiconTypeID
}

// MarshalWithLexiconType marshals an object, adding a specified type.
func MarshalWithLexiconType(t string, v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var temp map[string]any
	err = json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	temp["$type"] = t
	return json.Marshal(temp)
}
