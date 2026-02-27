package lexica

import (
	"fmt"
	"slices"
	"strings"
)

func (lexicon *Lexicon) generateStructAndDependencies(s *strings.Builder, defname, description string, properties map[string]Property, required []string, isRecord bool, name string) {
	if isRecord {
		fmt.Fprintf(s, "const %s_Description = \"%s\"\n", defname, description)
	}
	lexicon.generateStruct(s, defname, properties, required, isRecord, name)
	lexicon.generateDependencies(s, defname, properties)
}

func (lexicon *Lexicon) generateStruct(s *strings.Builder, defname string, properties map[string]Property, required []string, isRecord bool, name string) {
	if isRecord {
		fmt.Fprintf(s, "// %s is a record with Lexicon type %s#%s\n", defname, lexicon.Id, name)
	}
	fmt.Fprintf(s, "type %s struct {\n", defname)
	if isRecord {
		fmt.Fprintf(s, "LexiconTypeID string `json:\"$type,omitempty\"`\n")
	}
	for _, propertyName := range sortedPropertyNames(properties) {
		required := slices.Contains(required, propertyName)
		property := properties[propertyName]
		switch property.Type {
		case "boolean":
			if required {
				fmt.Fprintf(s, "%s bool `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *bool `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "integer":
			if required {
				fmt.Fprintf(s, "%s int64 `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *int64 `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "string":
			if required {
				fmt.Fprintf(s, "%s string `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *string `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "array":
			itemstype := lexicon.resolveItemsType(defname, propertyName, property.Items)
			fmt.Fprintf(s, "%s []%s `json:\"%s,omitempty\"`\n", capitalize(propertyName), itemstype, propertyName)
		case "ref":
			reftype := lexicon.resolveRefType(property.Ref)
			fmt.Fprintf(s, "%s %s `json:\"%s,omitempty\"`\n", capitalize(propertyName), reftype, propertyName)
		case "unknown":
			if required {
				fmt.Fprintf(s, "%s any `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *any `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "blob":
			if required {
				fmt.Fprintf(s, "%s *slink.Blob `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *slink.Blob `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			fmt.Fprintf(s, "%s *%s `json:\"%s,omitempty\"`\n", capitalize(propertyName), uniontype, propertyName)
		case "bytes":
			if required {
				fmt.Fprintf(s, "%s []byte `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *[]byte `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
			if propertyName == "blocks" {
				// add an extra field for expansions of "blocks" (CAR) content
				fmt.Fprintf(s, "%s any `json:\"%s,omitempty\"`\n", "BlocksCARContent", "blocksCARContent")
			}
		case "cid-link":
			if required {
				fmt.Fprintf(s, "%s string `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *string `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		default:
			fmt.Fprintf(s, "// FIXME skipping unsupported property type %s %s required=%t %+v\n", propertyName, property.Type, required, property)
		}
	}
	fmt.Fprintf(s, "}\n\n")
}

func (lexicon *Lexicon) generateDependencies(s *strings.Builder, defname string, properties map[string]Property) {
	propertyNames := sortedPropertyNames(properties)
	for _, propertyName := range propertyNames {
		property := properties[propertyName]
		switch property.Type {
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			lexicon.generateUnion(s, uniontype, property.Refs)
		case "array":
			if property.Items.Type == "union" {
				uniontype := lexicon.resolveUnionFieldType(defname, propertyName) + "_Elem"
				lexicon.generateUnion(s, uniontype, property.Items.Refs)
			}
		}
	}
}

func (lexicon *Lexicon) generateUnion(s *strings.Builder, uniontype string, refs []string) {
	uniontypeinner := uniontype + "_Wrapper"

	// Describe the union with a comment.
	fmt.Fprintf(s, "// %s is a union with these possible values:\n", uniontype)
	for _, ref := range refs {
		ref = lexicon.resolveRef(ref)
		fmt.Fprintf(s, "// - %s (%s)\n", lexicon.unionElementType(ref), ref)
	}
	// Define the union type.
	fmt.Fprintf(s, "type %s struct {\n", uniontype)
	fmt.Fprintf(s, "Wrapper %s", uniontypeinner)
	fmt.Fprintf(s, "}\n\n")
	// Define the inner wrapper type.
	fmt.Fprintf(s, "// Value wrappers must conform to %s\n", uniontypeinner)
	fmt.Fprintf(s, "type %s interface {\n", uniontypeinner)
	fmt.Fprintf(s, "is%s()", uniontype)
	fmt.Fprintf(s, "}\n\n")
	// Define wrappers for each possible value type.
	for _, ref := range refs {
		ref = lexicon.resolveRef(ref)
		wrappertype := lexicon.unionElementWrapperType(uniontype, ref)
		wrappedtype := "*" + lexicon.unionElementType(ref)
		fmt.Fprintf(s, "// %s wraps values of type %s\n", wrappertype, wrappedtype)
		fmt.Fprintf(s, "type %s struct {\n", wrappertype)
		fmt.Fprintf(s, "Value %s\n", wrappedtype)
		fmt.Fprintf(s, "}\n\n")
		fmt.Fprintf(s, "func (t %s) is%s() {}\n\n", wrappertype, uniontype)
	}

	fmt.Fprintf(s, "func (u %s) MarshalJSON() ([]byte, error) {\n", uniontype)
	fmt.Fprintf(s, "switch v := u.Wrapper.(type) {\n")
	for _, ref := range refs {
		ref = lexicon.resolveRef(ref)
		wrappertype := lexicon.unionElementWrapperType(uniontype, ref)
		fmt.Fprintf(s, "case %s:\n", wrappertype)
		fmt.Fprintf(s, "return slink.MarshalWithLexiconType(\"%s\", v.Value)\n", ref)
	}
	fmt.Fprintf(s, "default:\n")
	fmt.Fprintf(s, "	return nil, fmt.Errorf(\"unsupported type %%T\", v)\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "}\n")

	fmt.Fprintf(s, "func (u *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
	fmt.Fprintf(s, "switch slink.LexiconTypeFromJSONBytes(data) {\n")
	for _, ref := range refs {
		ref = lexicon.resolveRef(ref)
		wrappertype := lexicon.unionElementWrapperType(uniontype, ref)
		fmt.Fprintf(s, "case \"%s\":\n", ref)
		fmt.Fprintf(s, "var v %s\n", lexicon.unionElementType(ref))
		fmt.Fprintf(s, "if err := json.Unmarshal(data, &v); err != nil {return err}\n")
		fmt.Fprintf(s, "u.Wrapper = %s{Value: &v}\n", wrappertype)
		fmt.Fprintf(s, "return nil\n")
	}
	fmt.Fprintf(s, "default:\n")
	fmt.Fprintf(s, "return nil\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "}\n\n")
}

func (lexicon *Lexicon) unionElementType(ref string) string {
	parts := strings.Split(ref, "#")
	var id, tag string
	switch len(parts) {
	case 1:
		id = parts[0]
		tag = "main"
	case 2:
		id = parts[0]
		tag = parts[1]
	default:
		return "string /* FIXME defaulting on unparsable union field ref " + fmt.Sprintf("%+v", ref) + " */"
	}
	name := symbolForID(id)
	if tag != "main" {
		name += "_" + capitalize(tag)
	}
	return name
}

func (lexicon *Lexicon) unionElementWrapperType(uniontype, ref string) string {
	return uniontype + "__" + lexicon.unionElementType(ref)
}
