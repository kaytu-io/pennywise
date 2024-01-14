package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
	tp.MappedBlocks = makeMappedBlocks(tp.Blocks)
	return nil
}

func makeMappedBlocks(blocks []Block) map[string]interface{} {
	mappedBlocks := make(map[string]interface{})
	for _, b := range blocks {
		if len(b.Labels) == 0 {
			mappedBlocks[b.Type] = b
		} else {
			if mappedBlocks[b.Type] == nil {
				mappedBlocks[b.Type] = make(map[string]interface{})
			}
			mappedBlocks[b.Type] = makeMappedBlockItem(mappedBlocks[b.Type].(map[string]interface{}), b.Labels, b)
		}
	}
	return mappedBlocks
}

func makeMappedBlockItem(mappedBlocks map[string]interface{}, labels []string, block Block) map[string]interface{} {
	if len(labels) == 1 {
		mappedBlocks[labels[0]] = block
		return mappedBlocks
	} else {
		if mappedBlocks[labels[0]] == nil {
			mappedBlocks[labels[0]] = make(map[string]interface{})
		}
		mappedBlocks[labels[0]] = makeMappedBlockItem(mappedBlocks[labels[0]].(map[string]interface{}), labels[1:], block)
		return mappedBlocks
	}
}

func (tp *TerraformProject) MakeProjectMapStructure() (map[string]interface{}, error) {
	mapStructure := make(map[string]interface{})
	ctxVariableMap := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		for _, b := range tp.Blocks {
			var blockName string
			if len(b.Labels) > 0 {
				blockName = fmt.Sprintf("%s.%s", b.Type, strings.Join(b.Labels, "."))
			} else {
				blockName = b.Type
			}
			blockMapStructure, err := b.makeMapStructure(tp.MappedBlocks)
			if err != nil {
				return nil, err
			}
			mapStructure[blockName] = blockMapStructure
			ctxVariableMap = tp.makeBlockCtxVariableMap(ctxVariableMap, b)
		}
		tp.Context.Variables = makeCtxVariables(ctxVariableMap)
	}

	return mapStructure, nil
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

func (tp *TerraformProject) makeBlockCtxVariableMap(ctxVariableMap map[string]interface{}, b Block) map[string]interface{} {
	blockType := GetBlockTypeByType(b.Type)
	if blockType.name == "resource" {
		ctxVariableMap = tp.makeCtxVariableMapByLabels(ctxVariableMap, b, b.Labels)
	} else {
		if len(b.Labels) > 0 {
			if _, ok := ctxVariableMap[blockType.refName]; !ok {
				ctxVariableMap[blockType.refName] = make(map[string]interface{})
			}
			ctxVariableMap[blockType.refName] = tp.makeCtxVariableMapByLabels(ctxVariableMap[blockType.refName].(map[string]interface{}), b, b.Labels)
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
		}
		ctxVariableMap[key] = tp.makeCtxVariableMapByLabels(ctxVariableMap[key].(map[string]interface{}), b, labels[1:])
	} else {
		ctxVariableMap[key] = b.CtxVariable
	}
	return ctxVariableMap
}
