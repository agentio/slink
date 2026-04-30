package lexica

import (
	"os"
	"sort"
	"strings"

	"charm.land/log/v2"
	"golang.org/x/tools/imports"
)

func capitalize(s string) string {
	if s == "atproto" {
		return "ATProto"
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func symbolForID(id string) string {
	if false { // I think the code looks fine if we don't do this.
		id = strings.TrimPrefix(id, "com.atproto.") // put these symbols in the top-level namespace
	}
	var s strings.Builder
	for _, part := range strings.Split(id, ".") {
		s.WriteString(capitalize(part))
	}
	return s.String()
}

func sortedDefNames(defs map[string]*Def) []string {
	var defnames []string
	for defname := range defs {
		defnames = append(defnames, defname)
	}
	sort.Strings(defnames)
	return defnames
}

func sortedPropertyNames(properties map[string]Property) []string {
	var propnames []string
	for propname := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)
	return propnames
}

func writeFormattedFile(filename string, body string) error {
	log.Debugf("writing %s", filename)
	formatted, err := imports.Process(filename, []byte(body), nil)
	if err != nil {
		log.Errorf("failed to run goimports: %v\n%s", err, body)
		os.WriteFile(filename, []byte(body), 0644)
		return nil
	}
	return os.WriteFile(filename, []byte(formatted), 0644)
}
