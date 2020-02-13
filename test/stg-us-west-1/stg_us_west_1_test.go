package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Test the Terraform module in examples/stg-us-west-1
func TestStgUswest1(t *testing.T) {
	t.Parallel()

	// Create state for passing data between test stages
	// https://github.com/gruntwork-io/terratest#iterating-locally-using-test-stages
	exampleFolder := test_structure.CopyTerraformFolderToTemp(
		t,
		"../../",
		"examples/stg-us-west-1",
	)

	// At the end of the test, `terraform destroy` the created resources
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.Destroy(t, terraformOptions)
	})

	// Deploy the tested infrastructure
	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := configureTerraformOptions(t, exampleFolder)

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)

		// Run `terraform init` and `terraform apply` and fail if there are errors
		terraform.InitAndApply(t, terraformOptions)

		// Run `terraform output` to get the value of an output variable
		vpcID := terraform.Output(t, terraformOptions, "vpc_id")
		defaultSecurityGroupID := terraform.Output(t, terraformOptions, "default_security_group_id")
		privateCIDRBlocks := terraform.Output(t, terraformOptions, "public_cidr_blocks")
		publicCIDRBlocks := terraform.Output(t, terraformOptions, "private_cidr_blocks")
		privateSubnetIDs := terraform.Output(t, terraformOptions, "public_subnet_ids")
		publicSubnetIDs := terraform.Output(t, terraformOptions, "private_subnet_ids")

		// Save the VPC ID for the validation stage
		test_structure.SaveString(t, exampleFolder, "vpcID", vpcID)
		test_structure.SaveString(t, exampleFolder, "defaultSecurityGroupID", defaultSecurityGroupID)
		test_structure.SaveString(t, exampleFolder, "publicCIDRBlocks", publicCIDRBlocks)
		test_structure.SaveString(t, exampleFolder, "privateCIDRBlocks", privateCIDRBlocks)
		test_structure.SaveString(t, exampleFolder, "publicSubnetIDs", publicSubnetIDs)
		test_structure.SaveString(t, exampleFolder, "privateSubnetIDs", privateSubnetIDs)

	})

	// Validate the test infrastructure
	test_structure.RunTestStage(t, "validate", func() {
		testVpcValid(t, exampleFolder)
		testIAMBaselineValid(t, exampleFolder)
	})
}

func configureTerraformOptions(t *testing.T, exampleFolder string) *terraform.Options {

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"environment": "staging",
			"region":      "us-west-1",
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": "us-west-1",
		},
	}

	return terraformOptions
}

func testVpcValid(t *testing.T, exampleFolder string) {

	// Load outputs for validation
	vpcID := test_structure.LoadString(t, exampleFolder, "vpcID")
	defaultSecurityGroupID := test_structure.LoadString(t, exampleFolder, "defaultSecurityGroupID")
	privateCIDRBlocks := test_structure.LoadString(t, exampleFolder, "privateCIDRBlocks")
	publicCIDRBlocks := test_structure.LoadString(t, exampleFolder, "publicCIDRBlocks")
	privateSubnetIDs := test_structure.LoadString(t, exampleFolder, "privateSubnetIDs")
	publicSubnetIDs := test_structure.LoadString(t, exampleFolder, "publicSubnetIDs")

	// The VPC can be found by ID
	region := "us-west-1"
	maxRetries := 10
	timeBetweenRetries := 1 * time.Second
	description := fmt.Sprintf("Awaiting creation of VPC: %s", vpcID)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := aws.GetVpcByIdE(t, vpcID, region)
		if err != nil {
			return "", fmt.Errorf("Expected VPC %s to be present in %s", vpcID, region)
		}
		return "", nil
	})
	vpc, err := aws.GetVpcByIdE(t, vpcID, region)
	assert.NoError(t, err)

	// The VPC has valid characteristics
	assert.Regexp(t, "^vpc-[[:alnum:]]+$", vpc.Id)
	assert.True(t, len(vpc.Subnets) > 0)

	// The module has valid outputs
	assert.NotEmpty(t, defaultSecurityGroupID)
	assert.NotEmpty(t, publicCIDRBlocks)
	assert.NotEmpty(t, privateCIDRBlocks)
	assert.NotEmpty(t, publicSubnetIDs)
	assert.NotEmpty(t, privateSubnetIDs)

}

func testIAMBaselineValid(t *testing.T, exampleFolder string) {
	cmd := shell.Command{
		Command: "aws",
		Args:    []string{"iam", "get-role", "--role-name", "IAM-Support"},
	}

	out := shell.RunCommandAndGetOutput(t, cmd)
	assert.Contains(t, out, ":role/IAM-Support")
}
