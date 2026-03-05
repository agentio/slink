package slink

import (
	"encoding/json"
)

type Blob struct {
	LexiconTypeID string `json:"$type,omitempty"`
	Ref           Link   `json:"ref,omitempty"`
	MimeType      string `json:"mimeType,omitempty"`
	Size          int64  `json:"size"`
}

type Link struct {
	LexiconLink string `json:"$link"`
}

type Bytes struct {
	Bytes string `json:"$bytes"`
}

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
