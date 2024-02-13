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

type SubmissionsDiff struct {
	Current   Submission `json:"current"`
	CompareTo Submission `json:"compare_to"`
}

type SubmissionsDiffV2 struct {
	Current   SubmissionV2 `json:"current"`
	CompareTo SubmissionV2 `json:"compare_to"`
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

// ReadSubmissionFile Reads a submission from a file
func ReadSubmissionFile(id string) (*Submission, error) {
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

	var submission Submission
	err = json.Unmarshal(jsonData, &submission)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &submission, nil
}

func getAllSubmissions() ([]Submission, error) {
	submissionsDir := filepath.Join(pkg.PennywiseDir, "submissions")

	files, err := ioutil.ReadDir(submissionsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading submissions directory: %v", err)
	}

	var submissions []Submission

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		jsonFilePath := filepath.Join(submissionsDir, file.Name())

		jsonData, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading JSON file %s: %v", jsonFilePath, err)
		}

		var submission Submission
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

func GetLatestSubmission() (*Submission, error) {
	submissions, err := getAllSubmissions()
	if err != nil {
		return nil, err
	}

	if len(submissions) == 0 {
		return nil, fmt.Errorf("no submissions found")
	}

	return &submissions[0], nil
}
