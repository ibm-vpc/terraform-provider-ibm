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

func TestAccIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "authorized_cidr_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "allocations.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidr_allocations" "is_public_address_range_authorized_cidr_allocations_instance" {
			authorized_cidr_id = "authorized_cidr_id"
			allocations[].resource_type = "floating_ip"
		}
	`)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["address"] = "203.0.113.1"
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::floating-ip:r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/floating_ips/r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["id"] = "r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["name"] = "my-floating-ip"
		model["resource_type"] = "floating_ip"
		model["cidr"] = "testString"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItem)
	model.Address = core.StringPtr("203.0.113.1")
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::floating-ip:r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/floating_ips/r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.ID = core.StringPtr("r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.Name = core.StringPtr("my-floating-ip")
	model.ResourceType = core.StringPtr("floating_ip")
	model.CIDR = core.StringPtr("testString")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["address"] = "203.0.113.1"
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::floating-ip:r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/floating_ips/r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["id"] = "r006-f45e0d90-12a8-4460-8210-290ff2ab75cd"
		model["name"] = "my-floating-ip"
		model["resource_type"] = "floating_ip"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReference)
	model.Address = core.StringPtr("203.0.113.1")
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::floating-ip:r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/floating_ips/r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.ID = core.StringPtr("r006-f45e0d90-12a8-4460-8210-290ff2ab75cd")
	model.Name = core.StringPtr("my-floating-ip")
	model.ResourceType = core.StringPtr("floating_ip")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["cidr"] = "testString"
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::public-address-range:r006-a4841334-b584-4293-938e-3bc63b4a5b6a"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/public_address_ranges/r006-a4841334-b584-4293-938e-3bc63b4a5b6a"
		model["id"] = "r006-a4841334-b584-4293-938e-3bc63b4a5b6a"
		model["name"] = "my-public-address-range"
		model["resource_type"] = "public_address_range"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReference)
	model.CIDR = core.StringPtr("testString")
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::public-address-range:r006-a4841334-b584-4293-938e-3bc63b4a5b6a")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/public_address_ranges/r006-a4841334-b584-4293-938e-3bc63b4a5b6a")
	model.ID = core.StringPtr("r006-a4841334-b584-4293-938e-3bc63b4a5b6a")
	model.Name = core.StringPtr("my-public-address-range")
	model.ResourceType = core.StringPtr("public_address_range")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
