package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type TerraformProject struct {
	Directory    string
	Files        []*hcl.File
	Blocks       []Block
	MappedBlocks map[string]interface{}
	Context      *hcl.EvalContext
	Diags        Diags

	logger *zap.Logger
}

func NewTerraformProject(dir string, logger *zap.Logger) *TerraformProject {
	ctx := &hcl.EvalContext{}
	ctx.Variables = make(map[string]cty.Value)
	ctx.Functions = ContextFunctions
	return &TerraformProject{
		Directory: dir,
		Context:   ctx,
		logger:    logger,
		Diags:     Diags{Name: dir, Type: TfProjectDiag},
	}
}

func (tp *TerraformProject) resetDiags() {
	for _, b := range tp.Blocks {
		b.Diags = Diags{Name: b.Name, Type: BlockDiag}
	}
	tp.Diags = Diags{Name: tp.Directory, Type: TfProjectDiag}
}

func getFiles(path string) ([]*hcl.File, error) {
	hclParser := hclparse.NewParser()
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var files []*hcl.File
	for _, info := range fileInfos {
		if info.IsDir() {
			childFiles, err := getFiles(filepath.Join(path, info.Name()))
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

func (tp *TerraformProject) FindFiles() error {
	files, err := getFiles(tp.Directory)
	if err != nil {
		return err
	}
	tp.Files = files
	return nil
}

func (tp *TerraformProject) ParseProjectBlocks() error {
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

func (tp *TerraformProject) MakeProjectMapStructure() (map[string]interface{}, error) {
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
				blockMapStructure, err := b.makeMapStructure(b.Name, tp.Context)
				if err != nil {
					return nil, err
				}
				mapStructure[b.Name] = blockMapStructure
				ctxVariableMap = makeBlockCtxVariableMap(ctxVariableMap, b)
			} else {
				for key, eachItems := range forEachItems {
					ctx := tp.Context
					ctx.Variables["each"] = cty.ObjectVal(map[string]cty.Value{"value": eachItems})
					clonedBlock := b.cloneBlock(key)
					blockMapStructure, err := b.makeMapStructure(clonedBlock.Name, ctx)
					if err != nil {
						return nil, err
					}
					mapStructure[clonedBlock.Name] = blockMapStructure
					ctxVariableMap = makeBlockCtxVariableMap(ctxVariableMap, b, key)
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

	return mapStructure, nil
}

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
