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

func TestAccIBMIsPublicAddressRangeAuthorizedCIDRDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "is_public_address_range_authorized_cidr_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "allocation.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "availability_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "cidr"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidr" "is_public_address_range_authorized_cidr_instance" {
			id = "id"
		}
	`)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRPublicAddressRangeAuthorizedCIDRAllocationToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["count"] = int(2)
		model["profile_family"] = "user"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocation)
	model.Count = core.Int64Ptr(int64(2))
	model.ProfileFamily = core.StringPtr("user")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRPublicAddressRangeAuthorizedCIDRAllocationToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
