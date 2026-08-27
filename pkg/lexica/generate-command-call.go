package lexica

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"charm.land/log/v2"
	"github.com/iancoleman/strcase"
)

func (catalog *Catalog) GenerateCallCommands(root string) error {
	os.RemoveAll(root)
	os.MkdirAll(root, 0755)
	var wg sync.WaitGroup
	for _, lexicon := range catalog.Lexicons {
		wg.Go(func() {
			lexicon.generateCallCommands(root)
		})
	}
	wg.Wait()
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsDir() {
			wg.Go(func() {
				catalog.generateInternalCommand(path, "Call XRPC methods")
			})
		}
		return nil
	})
	wg.Wait()
	return nil
}

func (lexicon *Lexicon) generateCallCommands(root string) {
	for defname, def := range lexicon.Defs {
		if ManifestIncludes(lexicon.Id, defname) {
			switch def.Type {
			case "query":
				lexicon.generateCallCommandForDef(root, defname, def)
			case "procedure":
				lexicon.generateCallCommandForDef(root, defname, def)
			}
		}
	}
}

func (lexicon *Lexicon) generateCallCommandForDef(root, defname string, def *Def) {
	if defname != "main" {
		log.Errorf("Can't generate call command for %s %s", lexicon.Id, defname)
	}
	defname = symbolForID(lexicon.Id)
	id := strings.Replace(lexicon.Id, ".", "-", 2) // merge initial segments of the lexicon id
	dirname := strings.ToLower(root + "/" + strings.ReplaceAll(id, ".", "/"))
	os.MkdirAll(dirname, 0755)
	filename := dirname + "/cmd.go"
	parts := strings.Split(id, ".")
	lastpart := parts[len(parts)-1]
	packagename := strings.ToLower(lastpart)
	commandname := strcase.ToKebab(lastpart)
	commandname = strings.ReplaceAll(commandname, "-v-", "-v")
	handlerName := symbolForID(lexicon.Id)

	s := &strings.Builder{}
	packageComment(s, packagename)
	fmt.Fprintf(s, "package %s // %s\n\n", packagename, lexicon.Id)
	fmt.Fprintf(s, "import \"github.com/spf13/cobra\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/froda\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/slink\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/gen/xrpc\"\n")
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	fmt.Fprintf(s, "var _output string\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyVariable := sanitize(propertyName)
			propertyValue := def.Parameters.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "var %s string\n", propertyVariable)
			case "integer":
				fmt.Fprintf(s, "var %s int64\n", propertyVariable)
			case "boolean":
				fmt.Fprintf(s, "var %s bool\n", propertyVariable)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "var %s []string\n", propertyVariable)
				} else {
					fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyVariable, propertyValue)
				}
			default:
				fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyVariable, propertyValue)
			}
		}
	} else if def.Type == "procedure" && def.Input != nil {
		for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
			propertyVariable := sanitize(propertyName)
			propertyValue := def.Input.Schema.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "var %s string\n", propertyVariable)
			case "integer":
				fmt.Fprintf(s, "var %s int64\n", propertyVariable)
			case "boolean":
				fmt.Fprintf(s, "var %s bool\n", propertyVariable)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "var %s []string\n", propertyVariable)
				} else {
					fmt.Fprintf(s, "var %s string // (this should be a filename)\n", propertyVariable)
				}
			case "unknown", "ref", "union":
				fmt.Fprintf(s, "var %s string // (this should be a filename)\n", propertyVariable)
			default:
				fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyVariable, propertyValue)
			}
		}
	}
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s\",\n", commandname)
	fmt.Fprintf(s, "Short: slink.TruncateShort(xrpc.%s_Description),\n", handlerName)
	fmt.Fprintf(s, "Long: xrpc.%s_Description,\n", handlerName)
	fmt.Fprintf(s, "Args: cobra.NoArgs,\n")
	fmt.Fprintf(s, "RunE: func(cmd *cobra.Command, args []string) error {\n")
	if def.Type == "query" {
		fmt.Fprintf(s, "client := froda.NewClient()\n")
		fmt.Fprintf(s, "response, err := xrpc.%s(\n", handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		if def.Parameters != nil {
			for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
				propertyVariable := sanitize(propertyName)
				propertyValue := def.Parameters.Properties[propertyName]
				switch propertyValue.Type {
				case "string":
					fmt.Fprintf(s, "%s,\n", propertyVariable)
				case "integer":
					fmt.Fprintf(s, "%s,\n", propertyVariable)
				case "boolean":
					fmt.Fprintf(s, "%s,\n", propertyVariable)
				case "array":
					if propertyValue.Items.Type == "string" {
						fmt.Fprintf(s, "%s,\n", propertyVariable)
					}
				default:
				}
			}
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		fmt.Fprintf(s, "return slink.Write(cmd.OutOrStdout(), _output, response)\n")
	} else if def.Type == "procedure" && (def.Input == nil || def.Input.Encoding == "application/json") {
		fmt.Fprintf(s, "var err error\n")
		if def.Input != nil {
			for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
				propertyValue := def.Input.Schema.Properties[propertyName]
				if propertyValue.Type == "unknown" ||
					propertyValue.Type == "ref" ||
					propertyValue.Type == "union" ||
					(propertyValue.Type == "array" && propertyValue.Items.Type != "string") {
					fmt.Fprintf(s, "%s_value, err := slink.ReadJSONFile(%s)\n", propertyName, propertyName)
					fmt.Fprintf(s, "if err != nil {return err}\n")
				}
			}
		}
		fmt.Fprintf(s, "client := froda.NewClient()\n")
		resultIfNeeded := ""
		assignment := "="
		if def.Output != nil {
			resultIfNeeded = "response, "
			assignment = ":="
		}
		fmt.Fprintf(s, "%serr %s xrpc.%s(\n", resultIfNeeded, assignment, handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		if def.Input != nil {
			fmt.Fprintf(s, "&xrpc.%s_Input{\n", handlerName)
			for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
				propertyVariable := sanitize(propertyName)
				propertyValue := def.Input.Schema.Properties[propertyName]
				switch propertyValue.Type {
				case "string":
					if !slices.Contains(def.Input.Schema.Required, propertyName) {
						fmt.Fprintf(s, "%s: slink.CastStringToPointer(%s),\n", capitalize(propertyName), propertyVariable)
					} else {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyVariable)
					}
				case "integer":
					if !slices.Contains(def.Input.Schema.Required, propertyName) {
						fmt.Fprintf(s, "%s: slink.CastInt64ToPointer(%s),\n", capitalize(propertyName), propertyVariable)
					} else {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyVariable)
					}
				case "boolean":
					if !slices.Contains(def.Input.Schema.Required, propertyName) {
						fmt.Fprintf(s, "%s: slink.CastBoolToPointer(%s),\n", capitalize(propertyName), propertyVariable)
					} else {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyVariable)
					}
				case "array":
					if propertyValue.Items.Type == "string" {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyVariable)
					} else {
						itemstype := lexicon.resolveItemsType(defname+"_Input", propertyName, propertyValue.Items)
						fmt.Fprintf(s, "%s: slink.CastAnyToArray[xrpc.%s](%s_value),\n", capitalize(propertyName), itemstype[1:], propertyName)
					}
				case "unknown":
					fmt.Fprintf(s, "%s: %s_value,\n", capitalize(propertyName), propertyVariable)
				case "ref":
					reftype := lexicon.resolveRefType(propertyValue.Ref)
					if reftype[0] == '*' {
						fmt.Fprintf(s, "%s: slink.CastAnyToStruct[xrpc.%s](%s_value),\n", capitalize(propertyName), reftype[1:], propertyName)
					} else {
						fmt.Fprintf(s, "%s: slink.CastAnyToArray[xrpc.%s](%s_value),\n", capitalize(propertyName), reftype[3:], propertyName)
					}
				case "union":
					uniontype := lexicon.resolveUnionFieldType(defname+"_Input", propertyName)
					fmt.Fprintf(s, "%s: slink.CastAnyToStruct[xrpc.%s](%s_value),\n", capitalize(propertyName), uniontype, propertyName)
				default:
				}
			}
			fmt.Fprintf(s, "},\n")
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		if def.Output == nil {
			fmt.Fprintf(s, "return nil\n")
		} else {
			fmt.Fprintf(s, "return slink.Write(cmd.OutOrStdout(), _output, response)\n")
		}
	} else if def.Type == "procedure" && def.Input != nil {
		fmt.Fprintf(s, "client := froda.NewClient()\n")
		resultIfNeeded := ""
		if def.Output != nil {
			resultIfNeeded = "response, "
		}
		fmt.Fprintf(s, "%serr := xrpc.%s(\n", resultIfNeeded, handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		fmt.Fprintf(s, "cmd.InOrStdin(),\n")
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		if def.Output == nil {
			fmt.Fprintf(s, "return nil\n")
		} else {
			fmt.Fprintf(s, "return slink.Write(cmd.OutOrStdout(), _output, response)\n")
		}
	} else {
		fmt.Fprintf(s, "return errors.New(\"unimplemented\")")
	}
	fmt.Fprintf(s, "},\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "cmd.Flags().StringVarP(&_output, \"output\", \"o\", \"\", \"output destination\")\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyVariable := sanitize(propertyName)
			propertyValue := def.Parameters.Properties[propertyName]
			flagName := strcase.ToKebab(propertyName)
			description := propertyName
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s\")\n", propertyVariable, flagName, description)
			case "integer":
				fmt.Fprintf(s, "cmd.Flags().Int64Var(&%s, \"%s\", %d, \"%s\")\n", propertyVariable, flagName, int64Value(propertyValue.Default), description)
			case "boolean":
				fmt.Fprintf(s, "cmd.Flags().BoolVar(&%s, \"%s\", %t, \"%s\")\n", propertyVariable, flagName, boolValue(propertyValue.Default), description)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "cmd.Flags().StringArrayVar(&%s, \"%s\", nil, \"%s\")\n", propertyVariable, flagName, description)
				} else {
					fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyVariable, propertyValue)
				}
			default:
				fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyVariable, propertyValue)
			}
		}
	} else if def.Type == "procedure" && def.Input != nil {
		for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
			propertyVariable := sanitize(propertyName)
			propertyValue := def.Input.Schema.Properties[propertyName]
			flagName := strcase.ToKebab(propertyName)
			description := propertyName
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s\")\n", propertyVariable, flagName, description)
			case "integer":
				fmt.Fprintf(s, "cmd.Flags().Int64Var(&%s, \"%s\", %d, \"%s\")\n", propertyVariable, flagName, int64Value(propertyValue.Default), description)
			case "boolean":
				fmt.Fprintf(s, "cmd.Flags().BoolVar(&%s, \"%s\", %t, \"%s\")\n", propertyVariable, flagName, boolValue(propertyValue.Default), description)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "cmd.Flags().StringArrayVar(&%s, \"%s\", nil, \"%s\")\n", propertyVariable, flagName, description)
				} else {
					fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s (name of a json file)\")\n", propertyVariable, flagName, description)
				}
			case "unknown", "ref", "union":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s (name of a json file)\")\n", propertyVariable, flagName, description)
			default:
				fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyVariable, propertyValue)
			}
		}
	}
	fmt.Fprintf(s, "return cmd\n")
	fmt.Fprintf(s, "}\n")
	if false { // append lexicon source to generated file
		lexicon.generateSourceComment(s)
	}
	if err := writeFormattedFile(filename, s.String()); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
}

func int64Value(v any) int64 {
	switch v := v.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 1
	}
}

func boolValue(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	default:
		return false
	}
}
