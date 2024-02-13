package schema

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/pkg"
	"github.com/sony/sonyflake"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

// ReadSubmissionFileV2 Reads a submission from a file
func ReadSubmissionFileV2(id string) (*SubmissionV2, error) {
	submissionsDir := filepath.Join(pkg.PennywiseDir, "submissions")

	jsonFilePath := filepath.Join(submissionsDir, id+".json")

	fileInfo, err := os.Stat(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("rrror checking JSON file: %v", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("file %s is not a regular file", jsonFilePath)
	}

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var submission SubmissionV2
	err = json.Unmarshal(jsonData, &submission)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &submission, nil
}

func getAllSubmissionsV2() ([]SubmissionV2, error) {
	submissionsDir := filepath.Join(pkg.PennywiseDir, "submissions")

	files, err := ioutil.ReadDir(submissionsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading submissions directory: %v", err)
	}

	var submissions []SubmissionV2

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		jsonFilePath := filepath.Join(submissionsDir, file.Name())

		jsonData, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading JSON file %s: %v", jsonFilePath, err)
		}

		var submission SubmissionV2
		err = json.Unmarshal(jsonData, &submission)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON file %s: %v", jsonFilePath, err)
		}

		submissions = append(submissions, submission)
	}

	// Sort submissions by CreatedAt in descending order
	sort.Slice(submissions, func(i, j int) bool {
		return submissions[i].CreatedAt.After(submissions[j].CreatedAt)
	})

	return submissions, nil
}

func GetLatestSubmissionV2() (*SubmissionV2, error) {
	submissions, err := getAllSubmissionsV2()
	if err != nil {
		return nil, err
	}

	if len(submissions) == 0 {
		return nil, fmt.Errorf("no submissions found")
	}

	return &submissions[0], nil
}
