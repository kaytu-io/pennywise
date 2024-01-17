package hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// TerraformProject represents a terraform project for the parser
type TerraformProject struct {
	Directory    string
	Files        []*hcl.File
	Blocks       []Block
	MappedBlocks map[string]interface{}
	Context      *hcl.EvalContext
	Diags        Diags
}

// newTerraformProject creates a new terraform project object for a directory
func newTerraformProject(dir string) *TerraformProject {
	ctx := &hcl.EvalContext{}
	ctx.Variables = make(map[string]cty.Value)
	ctx.Functions = ContextFunctions
	return &TerraformProject{
		Directory: dir,
		Context:   ctx,
		Diags:     Diags{Name: dir, Type: TfProjectDiag},
	}
}

// resetDiags reset all diags (uses in all retries except the last one)
func (tp *TerraformProject) resetDiags() {
	for _, b := range tp.Blocks {
		b.Diags = Diags{Name: b.Name, Type: BlockDiag}
	}
	tp.Diags = Diags{Name: tp.Directory, Type: TfProjectDiag}
}

// getTerraformFiles finds all terraform files in a directory recursively
func getTerraformFiles(path string) ([]*hcl.File, error) {
	hclParser := hclparse.NewParser()
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var files []*hcl.File
	for _, info := range fileInfos {
		if info.IsDir() {
			childFiles, err := getTerraformFiles(filepath.Join(path, info.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, childFiles...)
		}
		if strings.HasSuffix(info.Name(), ".tf") {
			parseFunc := hclParser.ParseHCLFile
			file, diags := parseFunc(filepath.Join(path, info.Name()))
			if diags.HasErrors() {
				return nil, diags
			}
			files = append(files, file)
		}

		if strings.HasSuffix(info.Name(), ".tf.json") {
			parseFunc := hclParser.ParseJSONFile
			file, diags := parseFunc(filepath.Join(path, info.Name()))
			if diags.HasErrors() {
				return nil, diags
			}
			files = append(files, file)
		}

	}
	return files, nil
}

// FindFiles finds and stores terraform files in a directory
func (tp *TerraformProject) FindFiles() error {
	files, err := getTerraformFiles(tp.Directory)
	if err != nil {
		return err
	}
	tp.Files = files
	return nil
}

// FindProjectBlocks finds and stores blocks in terraform files
func (tp *TerraformProject) FindProjectBlocks() error {
	var totalBlocks []Block
	for _, file := range tp.Files {
		myBlocks, err := getFileBlocks(tp.Context, file)
		if err != nil {
			return err
		}
		totalBlocks = append(totalBlocks, myBlocks...)
	}
	tp.Blocks = totalBlocks
	return nil
}

// makeProjectMapStructure returns a map structure of a terraform project blocks and values
func (tp *TerraformProject) makeProjectMapStructure() map[string]interface{} {
	ctxVariableMap := make(map[string]interface{})
	var mapStructure map[string]interface{}
	var oldMapStructure map[string]interface{}
	var retry int
	for retry < 50 {
		tp.resetDiags()
		mapStructure = make(map[string]interface{})
		for _, b := range tp.Blocks {
			forEachItems, err := b.checkForEach()
			if err != nil {
				b.Diags.Errors = append(b.Diags.Errors, fmt.Errorf("error while getting for each: %s", err))
			}
			if forEachItems == nil {
				blockMapStructure := b.makeMapStructure(b.Name, tp.Context)
				mapStructure[b.Name] = blockMapStructure
				ctxVariableMap = updateBlockCtxVariableMap(ctxVariableMap, b)
			} else {
				for key, eachItems := range forEachItems {
					ctx := tp.Context
					ctx.Variables["each"] = cty.ObjectVal(map[string]cty.Value{"value": eachItems})
					blockMapStructure := b.makeMapStructure(fmt.Sprintf("%s[\"%s\"]", b.Name, key), ctx)
					mapStructure[fmt.Sprintf("%s[\"%s\"]", b.Name, key)] = blockMapStructure
					ctxVariableMap = updateBlockCtxVariableMap(ctxVariableMap, b, key)
				}
			}
			tp.Diags.ChildDiags = append(tp.Diags.ChildDiags, b.Diags)
		}
		tp.Context.Variables = makeCtxVariables(ctxVariableMap)
		if mapsEqual(oldMapStructure, mapStructure) {
			break
		}
		oldMapStructure = mapStructure
		retry++
	}

	return mapStructure
}

// mapsEqual checks if two maps are equal
func mapsEqual(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		if value2, ok := map2[key]; ok {
			if !reflect.DeepEqual(value1, value2) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
