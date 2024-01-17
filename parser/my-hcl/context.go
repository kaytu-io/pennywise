package my_hcl

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
)

// makeCtxVariables converts a map structure of context to acceptable type for hcl sdk
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

// updateBlockCtxVariableMap updates context variable map by getting a new block
// additional labels are used for blocks having for_each expression
func updateBlockCtxVariableMap(ctxVariableMap map[string]interface{}, b Block, additionalLabels ...string) map[string]interface{} {
	blockType := GetBlockTypeByType(b.Type)
	if blockType.name == "resource" {
		ctxVariableMap = updateCtxVariableMapByLabels(ctxVariableMap, b, append(b.Labels, additionalLabels...))
	} else {
		if len(b.Labels)+len(additionalLabels) > 0 {
			if _, ok := ctxVariableMap[blockType.refName]; !ok {
				ctxVariableMap[blockType.refName] = make(map[string]interface{})
			}
			ctxVariableMap[blockType.refName] = updateCtxVariableMapByLabels(ctxVariableMap[blockType.refName].(map[string]interface{}), b, append(b.Labels, additionalLabels...))
		} else {
			ctxVariableMap[blockType.refName] = b.CtxVariable
		}
	}
	return ctxVariableMap
}

// updateCtxVariableMapByLabels updates context variable map by getting a new block and its labels
func updateCtxVariableMapByLabels(ctxVariableMap map[string]interface{}, b Block, labels []string) map[string]interface{} {
	key := labels[0]
	if len(labels) > 1 {
		if _, ok := ctxVariableMap[key]; !ok {
			ctxVariableMap[key] = make(map[string]interface{})
		} else if _, ok := ctxVariableMap[key].(map[string]interface{}); !ok {
			ctxVariableMap[key] = make(map[string]interface{})
		}
		ctxVariableMap[key] = updateCtxVariableMapByLabels(ctxVariableMap[key].(map[string]interface{}), b, labels[1:])
	} else {
		ctxVariableMap[key] = b.CtxVariable
	}
	return ctxVariableMap
}
