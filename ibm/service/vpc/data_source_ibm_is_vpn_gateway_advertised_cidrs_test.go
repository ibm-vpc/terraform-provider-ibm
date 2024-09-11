// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
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

func TestAccIBMIsVPNGatewayAdvertisedCidrsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayAdvertisedCidrsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "vpn_gateway"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "advertised_cidrs.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_advertised_cidrs.is_vpn_gateway_advertised_cidrs_instance", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayAdvertisedCidrsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_gateway_advertised_cidrs" "is_vpn_gateway_advertised_cidrs_instance" {
			vpn_gateway = "vpn_gateway_id"
		}
	`)
}

func TestDataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionFirstToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpn_gateways/ddf51bec-3424-11e8-b467-0ed5f89f718b/advertised_cidrs?limit=20"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNGatewayAdvertisedCIDRCollectionFirst)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpn_gateways/ddf51bec-3424-11e8-b467-0ed5f89f718b/advertised_cidrs?limit=20")

	result, err := vpc.DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionFirstToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionNextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpn_gateways/ddf51bec-3424-11e8-b467-0ed5f89f718b/advertised_cidrs?start=ffd653466e284937896724b2dd044c9c&limit=20"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNGatewayAdvertisedCIDRCollectionNext)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpn_gateways/ddf51bec-3424-11e8-b467-0ed5f89f718b/advertised_cidrs?start=ffd653466e284937896724b2dd044c9c&limit=20")

	result, err := vpc.DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionNextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
