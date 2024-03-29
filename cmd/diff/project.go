package diff

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/infracost/external/providers"
	"github.com/kaytu-io/pennywise/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg"
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
		tfVarFiles := flags.ReadStringArrayFlag(cmd, "terraform-var-file")
		if jsonPath != nil {
			err := tfPlanJsonDiff(classic, *jsonPath, compareTo, usage, pkg.DefaultServerAddress)
			if err != nil {
				return err
			}
		} else {
			err := terraformProjectDiff(classic, projectPath, compareTo, usage, pkg.DefaultServerAddress, tfVarFiles)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func tfPlanJsonDiff(classic bool, jsonPath string, compareToId string, usage usagePackage.Usage, ServerClientAddress string) error {
	if classic {
		return fmt.Errorf("classic view not available for diff")
	}
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
	modularShowDiff := schema.ModularStateDiff{
		Resources: stateDiff.Resources,
		PriorCost: stateDiff.PriorCost,
		NewCost:   stateDiff.NewCost,
	}
	err = outputDiff.ShowStateCosts(&modularShowDiff)
	if err != nil {
		return err
	}
	return nil
}

func terraformProjectDiff(classic bool, projectPath string, compareToId string, usage usagePackage.Usage, ServerClientAddress string, tfVarFiles []string) error {
	if classic {
		return fmt.Errorf("classic view not available for diff")
	}
	var project *schema.ModuleDef
	var err error
	if providers.IsTerragruntNestedDir(projectPath, 5) {
		fmt.Println("terragrunt project...")
		project, err = hcl.ParseTerragruntProject(projectPath, usage)
	} else {
		project, err = hcl.ParseHclResources(projectPath, usage, tfVarFiles)
	}
	if err != nil {
		return err
	}
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}

	var compareTo *schema.SubmissionV2
	if compareToId == "" {
		compareTo, err = schema.GetLatestSubmissionV2()
		if err != nil {
			return err
		}
	} else {
		compareTo, err = schema.ReadSubmissionFileV2(compareToId)
		if err != nil {
			return err
		}
	}

	sub, err := schema.CreateSubmissionV2(*project)
	if err != nil {
		return err
	}
	err = sub.StoreAsFile()
	if err != nil {
		return err
	}

	req := schema.SubmissionsDiffV2{
		Current:   *sub,
		CompareTo: *compareTo,
	}

	stateDiff, err := serverClient.GetSubmissionsDiffV2(req)
	if err != nil {
		return err
	}
	err = outputDiff.ShowStateCosts(stateDiff)
	if err != nil {
		return err
	}

	return nil
}
