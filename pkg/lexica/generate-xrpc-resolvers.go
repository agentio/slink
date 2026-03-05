package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) resolveUnionFieldType(defname, propname string) string {
	return capitalize(defname) + "_" + capitalize(propname)
}

func (lexicon *Lexicon) resolveItemsType(defname, propname string, items *Items) string {
	switch items.Type {
	case "string":
		return "string"
	case "unknown":
		return "any"
	case "cid-link":
		return "slink.Link"
	case "ref":
		return lexicon.resolveRefType(items.Ref)
	case "union":
		return "*" + lexicon.resolveUnionFieldType(defname, propname) + "_Elem"
	default:
	}
	return "/* FIXME defaulting on unsupported items type: " + items.Type + " */ string"
}

func (lexicon *Lexicon) resolveRef(ref string) string {
	if ref == "#main" {
		return lexicon.Id
	}
	if ref[0] == '#' {
		return lexicon.Id + ref
	}
	return ref
}

func (lexicon *Lexicon) resolveRefType(ref string) string {
	if ref[0] == '#' {
		typename := symbolForID(lexicon.Id) + "_" + capitalize(ref[1:])
		return "*" + typename
	} else {
		parts := strings.Split(ref, "#")
		if len(parts) == 2 || len(parts) == 1 {
			id := parts[0]
			var tag string
			if len(parts) == 2 {
				tag = parts[1]
			} else {
				tag = "main"
			}
			var refType string
			refLexicon := LookupLexicon(id)
			if refLexicon != nil {
				refDef := refLexicon.Lookup(tag)
				if refDef != nil {
					refType = refDef.Type
				}
			}
			name := symbolForID(id)
			if tag != "main" {
				name += "_" + capitalize(tag)
			}
			if refType == "array" {
				refDef := refLexicon.Lookup(tag)
				if refDef != nil {
					refType = refDef.Type
				}
				if refDef.Items.Type == "string" {
					return "[]" + name + "_Elem"
				} else {
					return "[]*" + name + "_Elem"
				}
			}
			return "*" + name
		} else {
			return "/* FIXME defaulting on unparseable ref " + fmt.Sprintf("%+v", ref) + " */ string"
		}
	}
}
