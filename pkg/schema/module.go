package schema

// ModuleDef is a single module definition.
type ModuleDef struct {
	Address      string        `json:"address"`
	ChildModules []ModuleDef   `json:"child_modules"`
	Resources    []ResourceDef `json:"resources"`
}
