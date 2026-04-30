package lexica

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
)

func (lexicon *Lexicon) generateXRPCSourceFile(root string) {
	filename := xrpcFileName(root, lexicon.Id)
	if filename == "" {
		return
	}
	packagename := xrpcPackageName(root)
	if packagename == "" {
		return
	}
	s := &strings.Builder{}
	packageComment(s, packagename)
	fmt.Fprintf(s, "package %s // %s\n\n", packagename, lexicon.Id)
	fmt.Fprintf(s, "import \"encoding/json\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/slink\"\n")
	var generations int
	for _, name := range sortedDefNames(lexicon.Defs) {
		def := lexicon.Defs[name]
		if ManifestIncludes(lexicon.Id, name) {
			lexicon.generateDef(s, name, def)
			generations++
		}
	}
	if false { // append lexicon source to generated file
		lexicon.generateSourceComment(s)
	}
	if generations == 0 {
		return
	}
	if err := writeFormattedFile(filename, s.String()); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
}

func xrpcFileName(root, id string) string {
	base := strings.ReplaceAll(id, ".", "-")
	os.MkdirAll(root, 0755)
	filename := root + "/" + base + ".go"
	return filename
}

func xrpcPackageName(root string) string {
	return filepath.Base(root)
}
