package lexica

import (
	"fmt"
	"strings"

	"charm.land/log/v2"
)

func (lexicon *Lexicon) generateDef(s *strings.Builder, name string, def *Def) {
	prefix := symbolForID(lexicon.Id)
	var defname string
	if name == "main" {
		defname = prefix
	} else {
		defname = prefix + "_" + capitalize(name)
	}
	switch def.Type {
	case "query":
		lexicon.generateQuery(s, defname, def)
	case "procedure":
		lexicon.generateProcedure(s, defname, def)
	case "object":
		lexicon.generateStructAndDependencies(s, defname, def.Description, def.Properties, def.Required, def.Nullable, true, name)
	case "string":
		fmt.Fprintf(s, "type %s string\n", defname)
	case "record":
		lexicon.generateStructAndDependencies(s, defname, def.Description, def.Record.Properties, def.Record.Required, def.Record.Nullable, true, name)
	case "array":
		if def.Items.Type == "union" {
			uniontype := defname + "_Elem"
			lexicon.generateUnion(s, uniontype, def.Items.Refs)
		} else {
			fmt.Fprintf(s, "// FIXME: ungenerated array %+v\n", def)
		}
	case "token":
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "const %s string = \"%s\"\n\n", defname, name)
	case "permission-set":
		fmt.Fprintf(s, "// CHECKME skipping permission set %s\n", defname)
	case "subscription":
		lexicon.generateSubscribe(s, defname, def)
	default:
		log.Warnf("skipping %s.%s (type %s)", lexicon.Id, name, def.Type)
	}
}
