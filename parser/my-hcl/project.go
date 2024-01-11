package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"os"
	"path/filepath"
	"strings"
)

type TerraformProject struct {
	Directory    string
	Files        []*hcl.File
	Blocks       []Block
	MappedBlocks map[string]interface{}
}

func NewTerraformProject(dir string) *TerraformProject {
	return &TerraformProject{
		Directory: dir,
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
		myBlocks, err := getFileBlocks(file)
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
	for _, b := range tp.Blocks {
		var blockName string
		if len(b.Labels) > 0 {
			blockName = fmt.Sprintf("%s.%s", b.Type, strings.Join(b.Labels, "."))
		} else {
			blockName = b.Type
		}
		blockMapStructure, err := b.MakeMapStructure(tp.MappedBlocks)
		if err != nil {
			return nil, err
		}
		mapStructure[blockName] = blockMapStructure
	}
	return mapStructure, nil
}
