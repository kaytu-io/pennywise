package schema

import (
	"bytes"
	"fmt"
	"github.com/sony/sonyflake"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
	"time"
)

const (
	ConfigPath = "pennywise_projects_config.yaml"
)

type Project struct {
	ID             string              `yaml:"id"`
	Name           string              `yaml:"name"`
	CreatedAt      time.Time           `yaml:"created_at"`
	Directory      string              `yaml:"directory"`
	Description    string              `yaml:"description"`
	Tags           map[string][]string `yaml:"tags"`
	LastSubmission time.Time           `yaml:"last_submission"`
	SubmissionsIds []string            `yaml:"submissions_ids"`
}

// CreateProject creates a new project
func CreateProject(name, directory, description string, tags map[string][]string) (*Project, error) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := sf.NextID()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = directory
	}
	return &Project{
		ID:          fmt.Sprintf("project-%d", id),
		Name:        name,
		CreatedAt:   time.Now(),
		Directory:   directory,
		Description: description,
		Tags:        tags,
	}, nil
}

// AddSubmission add a new submission to the project
func (p *Project) AddSubmission(sub Submission) {
	p.SubmissionsIds = append(p.SubmissionsIds, sub.ID)
	p.LastSubmission = sub.CreatedAt
}

func GetProjects() ([]Project, error) {
	yamlData, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewBufferString(string(yamlData)))
	var projects []Project
	for {
		var p Project
		if err := decoder.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("project decode failed: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func WriteProjectsConfig(projects []Project) error {
	var projectContents []string
	for _, p := range projects {
		projectContent, err := yaml.Marshal(p)
		if err != nil {
			return err
		}
		projectContents = append(projectContents, string(projectContent))
	}
	err := os.WriteFile(ConfigPath, []byte(strings.Join(projectContents, "---\n")), 0644)
	if err != nil {
		return err
	}
	return nil
}
