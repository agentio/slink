package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateSubscribe(s *strings.Builder, defname string, def *Def) {
	fmt.Fprintf(s, "const %s_Description = \"%s\"\n\n", defname, def.Description)
	if def.Message != nil {
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		fmt.Fprintf(s, "// The handler is passed a reader that returns bytes.\n")
		fmt.Fprintf(s, "// Use this with your favorite CBOR decoder.\n")
		fmt.Fprintf(s, "type %s_Handler func(m io.Reader) error\n\n", defname)
		params += fmt.Sprintf(", fn %s_Handler", defname)
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client%s) (error) {\n", defname, params)
		fmt.Fprintf(s, "params := map[string]any{}\n")
		if paramsok {
			for _, parameterName := range sortedPropertyNames(def.Parameters.Properties) {
				if parameterName == "cursor" {
					cursorProperty := def.Parameters.Properties[parameterName]
					if cursorProperty.Type == "string" {
						fmt.Fprintf(s, "if %s != \"\" {params[\"%s\"] = %s}\n", parameterName, parameterName, parameterName)
					} else {
						fmt.Fprintf(s, "if %s >= 0 {params[\"%s\"] = %s}\n", parameterName, parameterName, parameterName)
					}
				} else {
					fmt.Fprintf(s, "params[\"%s\"] = %s,\n", parameterName, parameterName)
				}
			}
		}
		fmt.Fprintf(s, "return c.Subscribe(ctx, \"%s\", params, fn)\n", lexicon.Id)
		fmt.Fprintf(s, "}\n\n")
		lexicon.generateUnion(s, defname+"_Message", def.Message.Schema.Refs)
	} else {
		fmt.Fprintf(s, "// FIXME skipping subscription with no message %+v\n", def)
	}
}
