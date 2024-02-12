package schema

import (
	"encoding/json"
	"fmt"
	"github.com/sony/sonyflake"
	"os"
	"path/filepath"
	"time"
)

type SubmissionsDiff struct {
	Current   Submission `json:"current"`
	CompareTo Submission `json:"compare_to"`
}

// Submission to store and track resources and usage data for each run
type Submission struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	ProjectId string        `json:"project_id"`
	Resources []ResourceDef `json:"resources"`
}

// CreateSubmission creates a new submission to store resources and usage data
func CreateSubmission(resources []ResourceDef) (*Submission, error) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := sf.NextID()
	if err != nil {
		return nil, err
	}
	return &Submission{
		ID:        fmt.Sprintf("submission-%d", id),
		CreatedAt: time.Now(),
		Resources: resources,
	}, nil
}

// StoreAsFile stores the submission as a file in .pennywise/submissions directory
func (s *Submission) StoreAsFile() error {
	jsonData, err := json.MarshalIndent(*s, "", "  ")
	if err != nil {
		return err
	}

	pennywiseDir := ".pennywise"
	submissionsDir := filepath.Join(pennywiseDir, "submissions")
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

// ReadSubmissionFile Reads a submission from a file
func ReadSubmissionFile(id string) (*Submission, error) {
	pennywiseDir := ".pennywise"
	submissionsDir := filepath.Join(pennywiseDir, "submissions")

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

	var submission Submission
	err = json.Unmarshal(jsonData, &submission)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &submission, nil
}
