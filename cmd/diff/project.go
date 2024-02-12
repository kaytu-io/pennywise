package diff

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/infracost/external/providers"
	"github.com/kaytu-io/pennywise/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cmd/flags"
	outputDiff "github.com/kaytu-io/pennywise/pkg/output/diff"
	"github.com/kaytu-io/pennywise/pkg/parser/hcl"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	usagePackage "github.com/kaytu-io/pennywise/pkg/usage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var projectCommand = &cobra.Command{
	Use:   "project",
	Short: `Shows the costs by parsing a project resources.`,
	Long:  `Shows the costs by parsing a project resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		usagePath := flags.ReadStringOptionalFlag(cmd, "usage")
		var usage usagePackage.Usage
		if usagePath != nil {
			usageFile, err := os.Open(*usagePath)
			if err != nil {
				return fmt.Errorf("error while reading usage file %s", err)
			}
			defer usageFile.Close()

			ext := filepath.Ext(*usagePath)
			switch ext {
			case ".json":
				err = json.NewDecoder(usageFile).Decode(&usage)
			case ".yaml", ".yml":
				err = yaml.NewDecoder(usageFile).Decode(&usage)
			default:
				return fmt.Errorf("unsupported file format %s for usage file", ext)
			}
			if err != nil {
				return fmt.Errorf("error while parsing usage file %s", err)
			}

		} else {
			usage = usagePackage.Usage{}
		}

		classic := flags.ReadBooleanFlag(cmd, "classic")
		compareTo := flags.ReadStringFlag(cmd, "compare-to")

		jsonPath := flags.ReadStringOptionalFlag(cmd, "json-path")
		projectPath := flags.ReadStringFlag(cmd, "project-path")
		if jsonPath != nil {
			err := tfPlanJsonDiff(classic, *jsonPath, compareTo, usage, DefaultServerAddress)
			if err != nil {
				return err
			}
		} else {
			err := terraformProjectDiff(classic, projectPath, compareTo, usage, DefaultServerAddress)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func tfPlanJsonDiff(classic bool, jsonPath string, compareToId string, usage usagePackage.Usage, ServerClientAddress string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	resources, err := terraform.ParseTerraformPlanJson(file, usage)
	if err != nil {
		return err
	}
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}

	var compareTo *schema.Submission
	if compareToId == "" {
		compareTo, err = schema.GetLatestSubmission()
		if err != nil {
			return err
		}
	} else {
		compareTo, err = schema.ReadSubmissionFile(compareToId)
		if err != nil {
			return err
		}
	}

	sub, err := schema.CreateSubmission(resources)
	if err != nil {
		return err
	}
	err = sub.StoreAsFile()
	if err != nil {
		return err
	}

	req := schema.SubmissionsDiff{
		Current:   *sub,
		CompareTo: *compareTo,
	}
	stateDiff, err := serverClient.GetSubmissionsDiff(req)
	if err != nil {
		return err
	}

	err = outputDiff.ShowStateCosts(stateDiff)
	if err != nil {
		return err
	}
	return nil
}

func terraformProjectDiff(classic bool, projectPath string, compareToId string, usage usagePackage.Usage, ServerClientAddress string) error {
	var projects []hcl.ParsedProject
	var err error
	if providers.IsTerragruntNestedDir(projectPath, 5) {
		fmt.Println("terragrunt project...")
		projects, err = hcl.ParseTerragruntProject(projectPath, usage)
	} else {
		projects, err = hcl.ParseHclResources(projectPath, usage)
	}
	if err != nil {
		return err
	}
	for _, p := range projects {
		fmt.Println(p.Directory)
		fmt.Println("======================")
		serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
		if err != nil {
			return err
		}

		var compareTo *schema.Submission
		if compareToId == "" {
			compareTo, err = schema.GetLatestSubmission()
			if err != nil {
				return err
			}
		} else {
			compareTo, err = schema.ReadSubmissionFile(compareToId)
			if err != nil {
				return err
			}
		}

		sub, err := schema.CreateSubmission(p.GetResources())
		if err != nil {
			return err
		}
		err = sub.StoreAsFile()
		if err != nil {
			return err
		}

		req := schema.SubmissionsDiff{
			Current:   *sub,
			CompareTo: *compareTo,
		}

		stateDiff, err := serverClient.GetSubmissionsDiff(req)
		if err != nil {
			return err
		}
		err = outputDiff.ShowStateCosts(stateDiff)
		if err != nil {
			return err
		}
	}
	return nil
}
