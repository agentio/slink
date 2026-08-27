package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateQuery(s *strings.Builder, defname string, def *Def) {
	fmt.Fprintf(s, "const %s_Description = \"%s\"\n\n", defname, def.Description)
	if def.Output != nil && def.Output.Encoding == "application/json" {
		if def.Output.Schema.Type == "ref" {
			if def.Output.Schema.Ref[0] == '#' {
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", defname+"_"+capitalize(def.Output.Schema.Ref[1:]))
			} else {
				parts := strings.Split(def.Output.Schema.Ref, "#")
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", symbolForID(parts[0])+"_"+capitalize(parts[1]))
			}
		} else {
			lexicon.generateStructAndDependencies(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, def.Output.Schema.Nullable, false, "")
		}
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client%s) (*%s_Output"+", error) {\n", defname, params, defname)
		fmt.Fprintf(s, "var output %s_Output\n", defname)
		fmt.Fprintf(s, "params := map[string]any{\n")
		if paramsok {
			for _, parameterName := range sortedPropertyNames(def.Parameters.Properties) {
				fmt.Fprintf(s, "\"%s\":%s,\n", parameterName, sanitize(parameterName))
			}
		}
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Query, \"\", \"%s\", params, nil, &output); err != nil {\n", lexicon.Id)
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Output != nil {
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "// Returns []byte containing %s.\n", def.Output.Encoding)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client%s) ([]byte, error) {\n", defname, params)
		fmt.Fprintf(s, "params := map[string]any{\n")
		if paramsok {
			for _, parameterName := range sortedPropertyNames(def.Parameters.Properties) {
				fmt.Fprintf(s, "\"%s\":%s,\n", parameterName, sanitize(parameterName))
			}
		}
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "var output []byte\n")
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Query, \"\", \"%s\", params, nil, &output); err != nil {\n", lexicon.Id)
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else {
		fmt.Fprintf(s, "// FIXME skipping query with no output %+v\n", def)
	}
}

func parseQueryParameters(parameters *Parameters) (string, bool) {
	var parms []string
	propertyNames := sortedPropertyNames(parameters.Properties)
	for _, propertyName := range propertyNames {
		propertyValue := parameters.Properties[propertyName]
		declaration := sanitize(propertyName) + " "
		switch propertyValue.Type {
		case "integer":
			declaration += "int64"
		case "string":
			declaration += "string"
		case "boolean":
			declaration += "bool"
		case "array":
			if propertyValue.Items.Type == "string" {
				declaration += "[]string"
			} else {
				return "/* FIXME failing on unsupported parameter array value type: " + propertyValue.Items.Type + " */", false
			}
		default:
			return "/* FIXME failing on unsupported parameter value type: " + propertyValue.Type + "*/", false
		}
		parms = append(parms, declaration)
	}
	return ", " + strings.Join(parms, ", "), true
}

func sanitize(n string) string {
	if n == "type" {
		return "type_"
	}
	return n
}
