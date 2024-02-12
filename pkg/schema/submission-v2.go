package schema

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/pkg"
	"github.com/sony/sonyflake"
	"os"
	"path/filepath"
	"time"
)

// SubmissionV2 to store and track resources and usage data for each run
type SubmissionV2 struct {
	ID         string    `json:"id"`
	Version    string    `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	ProjectId  string    `json:"project_id"`
	RootModule ModuleDef `json:"root_modules"`
}

func (s *SubmissionV2) GetResources() []ResourceDef {
	var resources []ResourceDef
	resources = append(resources, getModuleResources(s.RootModule)...)

	return resources
}

func getModuleResources(module ModuleDef) []ResourceDef {
	var resources []ResourceDef
	resources = append(resources, module.Resources...)
	for _, childModule := range module.ChildModules {
		resources = append(resources, getModuleResources(childModule)...)
	}
	return resources
}

// CreateSubmissionV2 creates a new version2 submission to store resources and usage data
func CreateSubmissionV2(module ModuleDef) (*SubmissionV2, error) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := sf.NextID()
	if err != nil {
		return nil, err
	}
	return &SubmissionV2{
		ID:         fmt.Sprintf("submission-%d", id),
		Version:    "2.0.0",
		CreatedAt:  time.Now(),
		RootModule: module,
	}, nil
}

// StoreAsFile stores the submission as a file in .pennywise/submissions directory
func (s *SubmissionV2) StoreAsFile() error {
	jsonData, err := json.MarshalIndent(*s, "", "  ")
	if err != nil {
		return err
	}

	submissionsDir := filepath.Join(pkg.PennywiseDir, "submissions")
	err = os.MkdirAll(submissionsDir, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(submissionsDir, fmt.Sprintf("%s.json", s.ID))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
