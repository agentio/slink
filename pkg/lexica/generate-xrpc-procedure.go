package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateProcedure(s *strings.Builder, defname string, def *Def) {
	fmt.Fprintf(s, "const %s_Description = \"%s\"\n\n", defname, def.Description)
	if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStructAndDependencies(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, def.Input.Schema.Nullable, false, "")
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
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client, input *%s_Input) (*%s_Output"+", error) {\n", defname, defname, defname)
		fmt.Fprintf(s, "var output %s_Output\n", defname)
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Procedure, \"%s\", \"%s\", nil, input, &output); err != nil {\n", def.Input.Encoding, lexicon.Id)
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input == nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStructAndDependencies(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, def.Output.Schema.Nullable, false, "")
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client) (*%s_Output, error) {\n", defname, defname)
		fmt.Fprintf(s, "var output %s_Output\n", defname)
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Procedure, \"\", \"%s\", nil, nil, &output); err != nil {\n", lexicon.Id)
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output == nil {
		lexicon.generateStructAndDependencies(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, def.Input.Schema.Nullable, false, "")
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client, input *%s_Input) error {\n", defname, defname)
		fmt.Fprintf(s, "return c.Do(ctx, slink.Procedure, \"%s\", \"%s\", nil, input, nil)\n", def.Input.Encoding, lexicon.Id)
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input == nil && def.Output == nil {
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client) error {\n", defname)
		fmt.Fprintf(s, "return c.Do(ctx, slink.Procedure, \"\", \"%s\", nil, nil, nil)\n", lexicon.Id)
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
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
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client, input io.Reader) (*%s_Output, error) {\n", defname, defname)
		fmt.Fprintf(s, "var output %s_Output\n", defname)
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Procedure, \"%s\", \"%s\", nil, input, &output); err != nil {\n", def.Input.Encoding, lexicon.Id)
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil && def.Output == nil {
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client, input io.Reader) error {\n", defname)
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Procedure, \"%s\", \"%s\", nil, input, nil); err != nil {\n", def.Input.Encoding, lexicon.Id)
		fmt.Fprintf(s, "return err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else {
		fmt.Fprintf(s, "// FIXME skipping procedure with unhandled types\n")
	}
}
