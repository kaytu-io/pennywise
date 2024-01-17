package submission

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise-server/schema"
	"github.com/sony/sonyflake"
	"os"
	"path/filepath"
	"time"
)

// Submission to store and track resources and usage data for each run
type Submission struct {
	ID        string               `json:"id"`
	CreatedAt time.Time            `json:"created_at"`
	Resources []schema.ResourceDef `json:"resources"`
}

// CreateSubmission creates a new submission to store resources and usage data
func CreateSubmission(resources []schema.ResourceDef) (*Submission, error) {
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
