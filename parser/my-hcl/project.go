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
	}
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
		myBlocks, err := getFileBlocks(tp.logger, tp.Context, file)
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
		mapStructure = make(map[string]interface{})
		for _, b := range tp.Blocks {
			var blockName string
			if len(b.Labels) > 0 {
				blockName = fmt.Sprintf("%s.%s", b.Type, strings.Join(b.Labels, "."))
			} else {
				blockName = b.Type
			}
			forEachItems, err := b.checkForEach()
			if err != nil {
				return nil, err
			}
			if forEachItems == nil {
				blockMapStructure, err := b.makeMapStructure(blockName, tp.Context)
				if err != nil {
					return nil, err
				}
				mapStructure[blockName] = blockMapStructure
				ctxVariableMap = tp.makeBlockCtxVariableMap(ctxVariableMap, b)
			} else {
				for key, eachItems := range forEachItems {
					ctx := tp.Context
					ctx.Variables["each"] = cty.ObjectVal(map[string]cty.Value{"value": eachItems})
					blockMapStructure, err := b.makeMapStructure(fmt.Sprintf("%s[%s]", blockName, key), ctx)
					if err != nil {
						return nil, err
					}
					mapStructure[fmt.Sprintf("%s[%s]", blockName, key)] = blockMapStructure
					ctxVariableMap = tp.makeBlockCtxVariableMap(ctxVariableMap, b, key)
				}
			}
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
			// Check if the values are equal
			if !reflect.DeepEqual(value1, value2) {
				return false
			}
		} else {
			// Key not present in the second map
			return false
		}
	}

	return true
}

func makeCtxVariables(ctxVariableMap map[string]interface{}) map[string]cty.Value {
	ctxValuesMap := make(map[string]cty.Value)
	for key, value := range ctxVariableMap {
		if valuesMap, ok := value.(map[string]interface{}); ok {
			ctxValuesMap[key] = cty.ObjectVal(makeCtxVariables(valuesMap))
		} else if ctyValue, ok := value.(cty.Value); ok {
			ctxValuesMap[key] = ctyValue
		} else {
			fmt.Println("unknown value")
		}
	}
	return ctxValuesMap
}

func (tp *TerraformProject) makeBlockCtxVariableMap(ctxVariableMap map[string]interface{}, b Block, additionalLabels ...string) map[string]interface{} {
	blockType := GetBlockTypeByType(b.Type)
	if blockType.name == "resource" {
		ctxVariableMap = tp.makeCtxVariableMapByLabels(ctxVariableMap, b, append(b.Labels, additionalLabels...))
	} else {
		if len(b.Labels)+len(additionalLabels) > 0 {
			if _, ok := ctxVariableMap[blockType.refName]; !ok {
				ctxVariableMap[blockType.refName] = make(map[string]interface{})
			}
			ctxVariableMap[blockType.refName] = tp.makeCtxVariableMapByLabels(ctxVariableMap[blockType.refName].(map[string]interface{}), b, append(b.Labels, additionalLabels...))
		} else {
			ctxVariableMap[blockType.refName] = b.CtxVariable
		}
	}
	return ctxVariableMap
}

func (tp *TerraformProject) makeCtxVariableMapByLabels(ctxVariableMap map[string]interface{}, b Block, labels []string) map[string]interface{} {
	key := labels[0]
	if len(labels) > 1 {
		if _, ok := ctxVariableMap[key]; !ok {
			ctxVariableMap[key] = make(map[string]interface{})
		} else if _, ok := ctxVariableMap[key].(map[string]interface{}); !ok {
			ctxVariableMap[key] = make(map[string]interface{})
		}
		ctxVariableMap[key] = tp.makeCtxVariableMapByLabels(ctxVariableMap[key].(map[string]interface{}), b, labels[1:])
	} else {
		ctxVariableMap[key] = b.CtxVariable
	}
	return ctxVariableMap
}
