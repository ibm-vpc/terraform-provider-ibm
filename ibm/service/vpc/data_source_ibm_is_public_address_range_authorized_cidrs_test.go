// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.111.0-1bfb72c2-20260206-185521
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsPublicAddressRangeAuthorizedCidrsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCidrsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidrs.is_public_address_range_authorized_cidrs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidrs.is_public_address_range_authorized_cidrs_instance", "authorized_cidrs.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCidrsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidrs" "is_public_address_range_authorized_cidrs_instance" {
			allocation.profile_family = "user"
			availability_mode = "regional"
		}
	`)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		publicAddressRangeAuthorizedCIDRAllocationModel := make(map[string]interface{})
		publicAddressRangeAuthorizedCIDRAllocationModel["count"] = int(2)
		publicAddressRangeAuthorizedCIDRAllocationModel["profile_family"] = "user"

		zoneReferenceModel := make(map[string]interface{})
		zoneReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		zoneReferenceModel["name"] = "us-south-1"

		model := make(map[string]interface{})
		model["allocation"] = []map[string]interface{}{publicAddressRangeAuthorizedCIDRAllocationModel}
		model["availability_mode"] = "regional"
		model["cidr"] = "testString"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/public_address_range/authorized_cidrs/r006-ce2c481c-58cb-473d-b92d-92a79469eb0c"
		model["id"] = "r006-ce2c481c-58cb-473d-b92d-92a79469eb0c"
		model["ip_version"] = "ipv4"
		model["name"] = "ibm-house-effect-normal-round-yellow"
		model["network_prefix_length"] = int(24)
		model["resource_type"] = "public_address_range_authorized_cidr"
		model["zone"] = []map[string]interface{}{zoneReferenceModel}

		assert.Equal(t, result, model)
	}

	publicAddressRangeAuthorizedCIDRAllocationModel := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocation)
	publicAddressRangeAuthorizedCIDRAllocationModel.Count = core.Int64Ptr(int64(2))
	publicAddressRangeAuthorizedCIDRAllocationModel.ProfileFamily = core.StringPtr("user")

	zoneReferenceModel := new(vpcv1.ZoneReference)
	zoneReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	zoneReferenceModel.Name = core.StringPtr("us-south-1")

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDR)
	model.Allocation = publicAddressRangeAuthorizedCIDRAllocationModel
	model.AvailabilityMode = core.StringPtr("regional")
	model.CIDR = core.StringPtr("testString")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/public_address_range/authorized_cidrs/r006-ce2c481c-58cb-473d-b92d-92a79469eb0c")
	model.ID = core.StringPtr("r006-ce2c481c-58cb-473d-b92d-92a79469eb0c")
	model.IPVersion = core.StringPtr("ipv4")
	model.Name = core.StringPtr("ibm-house-effect-normal-round-yellow")
	model.NetworkPrefixLength = core.Int64Ptr(int64(24))
	model.ResourceType = core.StringPtr("public_address_range_authorized_cidr")
	model.Zone = zoneReferenceModel

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRAllocationToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["count"] = int(2)
		model["profile_family"] = "user"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocation)
	model.Count = core.Int64Ptr(int64(2))
	model.ProfileFamily = core.StringPtr("user")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRAllocationToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCidrsZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCidrsZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
