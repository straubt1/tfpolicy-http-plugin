// main.go

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/hashicorp/terraform-policy-plugin-framework/policy-plugin/plugins"

	packer "github.com/hashicorp/hcp-sdk-go/clients/cloud-packer-service/stable/2023-01-01/client/packer_service"
	packermodels "github.com/hashicorp/hcp-sdk-go/clients/cloud-packer-service/stable/2023-01-01/models"
)

func main() {

	plugins.RegisterFunction("env", envFunction)
	plugins.RegisterFunction("debug", debugGetHCPConfig)
	plugins.RegisterFunction("get_buckets", getBucketsFunction)

	plugins.Serve()
}

func envFunction(envVarName string, fallback ...string) (string, error) {
	val, ok := os.LookupEnv(envVarName)
	if ok && val != "" {
		return val, nil
	}
	if len(fallback) > 0 {
		return fallback[0], nil
	}
	return "", nil
}

func debugGetHCPConfig() (string, error) {

	envVars := map[string]string{
		"HCP_CLIENT_ID":       os.Getenv("HCP_CLIENT_ID"),
		"HCP_CLIENT_SECRET":   os.Getenv("HCP_CLIENT_SECRET"),
		"HCP_ORGANIZATION_ID": os.Getenv("HCP_ORGANIZATION_ID"),
		"HCP_PROJECT_ID":      os.Getenv("HCP_PROJECT_ID"),
	}

	// Check for empty values and return error if any are missing
	for k, v := range envVars {
		if v == "" {
			return "", fmt.Errorf("environment variable %s is not set or is empty", k)
		}
	}

	// Serialize the envVars map to a string
	var sb strings.Builder
	for k, v := range envVars {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func getBucketsFunction(bucketName string) (*packermodels.HashicorpCloudPacker20230101Version, error) {
	// Construct HCP config
	orgID, _ := os.LookupEnv("HCP_ORGANIZATION_ID")
	projID, _ := os.LookupEnv("HCP_PROJECT_ID")
	hcpConfig, err := config.NewHCPConfig(
		config.FromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}
	// if orgID != "" || projID != "" || hcpConfig.OrganizationID != "" || hcpConfig.ProjectID != "" {
	// 	log.Fatal("Organization ID and Project ID must be set in the environment variables HCP_ORGANIZATION_ID and HCP_PROJECT_ID, or in the HCPConfig.")
	// }

	// Construct HTTP client config
	httpclientConfig := httpclient.Config{
		HCPConfig: hcpConfig,
	}

	// Initialize SDK http client
	cl, err := httpclient.New(httpclientConfig)
	if err != nil {
		log.Fatal(err)
	}

	packerClient := packer.New(cl, nil)

	buckets, err := listPackerBuckets(packerClient, orgID, projID)
	if err != nil {
		log.Fatal(err)
	}
	for _, bucket := range buckets {
		versions, err := listPackerBucketVersions(packerClient, orgID, projID, bucket.Name)
		if err != nil {
			log.Fatal(err)
		}
		println("Bucket Name:", bucket.Name, "Versions Count:", len(versions))

		return versions[0], nil
	}

	return nil, nil
}

func listPackerBucketVersions(packerClient packer.ClientService, orgID, projID, bucketName string) ([]*packermodels.HashicorpCloudPacker20230101Version, error) {
	listParams := packer.NewPackerServiceListVersionsParams()
	listParams.LocationOrganizationID = orgID
	listParams.LocationProjectID = projID
	listParams.BucketName = bucketName
	listParams.SortingOrderBy = []string{"name"}

	resp, err := packerClient.PackerServiceListVersions(listParams, nil)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Versions, nil
}

// listPackerBuckets retrieves a list of Packer bucket names from HCP
func listPackerBuckets(packerClient packer.ClientService, orgID, projID string) ([]*packermodels.HashicorpCloudPacker20230101Bucket, error) {
	listParams := packer.NewPackerServiceListBucketsParams()
	listParams.LocationOrganizationID = orgID
	listParams.LocationProjectID = projID
	listParams.SortingOrderBy = []string{"name"}

	resp, err := packerClient.PackerServiceListBuckets(listParams, nil)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Buckets, nil
}
