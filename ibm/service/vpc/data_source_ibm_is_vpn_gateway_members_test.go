// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.114.0-a902401e-20260427-192904
 */

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsVPNGatewayMembersDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayMembersDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "members.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_members.is_vpn_gateway_members_instance", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayMembersDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_gateway_members" "is_vpn_gateway_members_instance" {
			vpn_gateway_id = "vpn_gateway_id"
		}
	`)
}

func TestDataSourceIBMIsVPNGatewayMembersPageLinkToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "testString"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.PageLink)
	model.Href = core.StringPtr("testString")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersPageLinkToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersVPNGatewayMemberCollectionItemToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		vpnGatewayMemberHealthReasonModel := make(map[string]interface{})
		vpnGatewayMemberHealthReasonModel["code"] = "cannot_reserve_ip_address"
		vpnGatewayMemberHealthReasonModel["message"] = "IP address exhaustion (release addresses on the VPN's subnet)."
		vpnGatewayMemberHealthReasonModel["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-vpn-health"

		vpnGatewayMemberLifecycleReasonModel := make(map[string]interface{})
		vpnGatewayMemberLifecycleReasonModel["code"] = "resource_suspended_by_provider"
		vpnGatewayMemberLifecycleReasonModel["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		vpnGatewayMemberLifecycleReasonModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		subnetReferenceModel := make(map[string]interface{})
		subnetReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		subnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["name"] = "my-subnet"
		subnetReferenceModel["resource_type"] = "subnet"

		reservedIPReferenceVPNGatewayContextModel := make(map[string]interface{})
		reservedIPReferenceVPNGatewayContextModel["address"] = "192.168.3.4"
		reservedIPReferenceVPNGatewayContextModel["deleted"] = []map[string]interface{}{deletedModel}
		reservedIPReferenceVPNGatewayContextModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		reservedIPReferenceVPNGatewayContextModel["id"] = "0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		reservedIPReferenceVPNGatewayContextModel["name"] = "my-reserved-ip"
		reservedIPReferenceVPNGatewayContextModel["resource_type"] = "subnet_reserved_ip"
		reservedIPReferenceVPNGatewayContextModel["subnet"] = []map[string]interface{}{subnetReferenceModel}

		ipModel := make(map[string]interface{})
		ipModel["address"] = "192.168.3.4"

		model := make(map[string]interface{})
		model["health_reasons"] = []map[string]interface{}{vpnGatewayMemberHealthReasonModel}
		model["health_state"] = "ok"
		model["id"] = "0717-9dce0ab3-a719-4582-ad71-00abfaae43ed"
		model["lifecycle_reasons"] = []map[string]interface{}{vpnGatewayMemberLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["private_ip"] = []map[string]interface{}{reservedIPReferenceVPNGatewayContextModel}
		model["public_ip"] = []map[string]interface{}{ipModel}
		model["role"] = "active"

		assert.Equal(t, result, model)
	}

	vpnGatewayMemberHealthReasonModel := new(vpcv1.VPNGatewayMemberHealthReason)
	vpnGatewayMemberHealthReasonModel.Code = core.StringPtr("cannot_reserve_ip_address")
	vpnGatewayMemberHealthReasonModel.Message = core.StringPtr("IP address exhaustion (release addresses on the VPN's subnet).")
	vpnGatewayMemberHealthReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-vpn-health")

	vpnGatewayMemberLifecycleReasonModel := new(vpcv1.VPNGatewayMemberLifecycleReason)
	vpnGatewayMemberLifecycleReasonModel.Code = core.StringPtr("resource_suspended_by_provider")
	vpnGatewayMemberLifecycleReasonModel.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	vpnGatewayMemberLifecycleReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	subnetReferenceModel := new(vpcv1.SubnetReference)
	subnetReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Deleted = deletedModel
	subnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Name = core.StringPtr("my-subnet")
	subnetReferenceModel.ResourceType = core.StringPtr("subnet")

	reservedIPReferenceVPNGatewayContextModel := new(vpcv1.ReservedIPReferenceVPNGatewayContext)
	reservedIPReferenceVPNGatewayContextModel.Address = core.StringPtr("192.168.3.4")
	reservedIPReferenceVPNGatewayContextModel.Deleted = deletedModel
	reservedIPReferenceVPNGatewayContextModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	reservedIPReferenceVPNGatewayContextModel.ID = core.StringPtr("0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	reservedIPReferenceVPNGatewayContextModel.Name = core.StringPtr("my-reserved-ip")
	reservedIPReferenceVPNGatewayContextModel.ResourceType = core.StringPtr("subnet_reserved_ip")
	reservedIPReferenceVPNGatewayContextModel.Subnet = subnetReferenceModel

	ipModel := new(vpcv1.IP)
	ipModel.Address = core.StringPtr("192.168.3.4")

	model := new(vpcv1.VPNGatewayMemberCollectionItem)
	model.HealthReasons = []vpcv1.VPNGatewayMemberHealthReason{*vpnGatewayMemberHealthReasonModel}
	model.HealthState = core.StringPtr("ok")
	model.ID = core.StringPtr("0717-9dce0ab3-a719-4582-ad71-00abfaae43ed")
	model.LifecycleReasons = []vpcv1.VPNGatewayMemberLifecycleReason{*vpnGatewayMemberLifecycleReasonModel}
	model.LifecycleState = core.StringPtr("stable")
	model.PrivateIP = reservedIPReferenceVPNGatewayContextModel
	model.PublicIP = ipModel
	model.Role = core.StringPtr("active")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersVPNGatewayMemberCollectionItemToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersVPNGatewayMemberHealthReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "cannot_reserve_ip_address"
		model["message"] = "IP address exhaustion (release addresses on the VPN's subnet)."
		model["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-vpn-health"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNGatewayMemberHealthReason)
	model.Code = core.StringPtr("cannot_reserve_ip_address")
	model.Message = core.StringPtr("IP address exhaustion (release addresses on the VPN's subnet).")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-vpn-health")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersVPNGatewayMemberHealthReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersVPNGatewayMemberLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNGatewayMemberLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersVPNGatewayMemberLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersReservedIPReferenceVPNGatewayContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		subnetReferenceModel := make(map[string]interface{})
		subnetReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		subnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["name"] = "my-subnet"
		subnetReferenceModel["resource_type"] = "subnet"

		model := make(map[string]interface{})
		model["address"] = "192.168.3.4"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["id"] = "0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["name"] = "my-reserved-ip"
		model["resource_type"] = "subnet_reserved_ip"
		model["subnet"] = []map[string]interface{}{subnetReferenceModel}

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	subnetReferenceModel := new(vpcv1.SubnetReference)
	subnetReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Deleted = deletedModel
	subnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Name = core.StringPtr("my-subnet")
	subnetReferenceModel.ResourceType = core.StringPtr("subnet")

	model := new(vpcv1.ReservedIPReferenceVPNGatewayContext)
	model.Address = core.StringPtr("192.168.3.4")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.ID = core.StringPtr("0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.Name = core.StringPtr("my-reserved-ip")
	model.ResourceType = core.StringPtr("subnet_reserved_ip")
	model.Subnet = subnetReferenceModel

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersReservedIPReferenceVPNGatewayContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersSubnetReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["name"] = "my-subnet"
		model["resource_type"] = "subnet"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.SubnetReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.Name = core.StringPtr("my-subnet")
	model.ResourceType = core.StringPtr("subnet")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayMembersIPToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["address"] = "192.168.3.4"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.IP)
	model.Address = core.StringPtr("192.168.3.4")

	result, err := vpc.DataSourceIBMIsVPNGatewayMembersIPToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
