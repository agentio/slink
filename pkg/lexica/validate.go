package lexica

import (
	"slices"
	"strings"

	"charm.land/log/v2"
)

func (lexicon *Lexicon) Validate(path string) error {
	if lexicon.Lexicon != 1 {
		log.Warnf("unexpected value for lexicon version: %d [%s]", lexicon.Lexicon, path)
	}
	expected := strings.ReplaceAll(lexicon.Id, ".", "/") + ".json"
	if !strings.HasSuffix(path, expected) {
		log.Warnf("lexicon file name should be %s [%s]", expected, path)
	}
	defCount := len(lexicon.Defs)
	if defCount == 0 {
		log.Warnf("lexicon file has no defs [%s]", path)
	}
	for k, v := range lexicon.Defs {
		v.Validate(path + ":" + k)
	}
	return nil
}

func (def *Def) Validate(path string) error {
	switch def.Type {
	case "boolean":
	case "integer":
	case "string":
	case "bytes":
	case "cid-link":
	case "blob":
	case "array":
	case "object":
	case "params":
	case "permission":
	case "token":
	case "ref":
	case "union":
	case "unknown":
	case "record":
	case "query":
		return def.ValidateQuery(path)
	case "procedure":
		return def.ValidateProcedure(path)
	case "subscription":
	case "permission-set":
	default:
		log.Warnf("def has unrecognized type: %s [%s]", def.Type, path)
	}
	return nil
}

func (def *Def) ValidateQuery(path string) error {
	if def.Parameters != nil {
		for name, value := range def.Parameters.Properties {
			if value.Type == "integer" {
				if value.Default == nil && !slices.Contains(def.Required, name) {
					log.Warnf("integer parameter %s has no default and is not required [%s]", name, path)
				}
			}
		}
	}
	return nil
}

func (def *Def) ValidateProcedure(path string) error {
	if def.Input != nil && def.Input.Encoding == "application/json" {
		for name, value := range def.Input.Schema.Properties {
			if value.Type == "integer" {
				if value.Default == nil && !slices.Contains(def.Input.Schema.Required, name) {
					log.Warnf("integer input field %s has no default and is not required [%s]", name, path)
				}
			}
		}
	}
	return nil
}
