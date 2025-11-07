// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsVirtualEndpointGatewayResourceBindingBasic(t *testing.T) {
	var conf vpcv1.EndpointGatewayResourceBinding
	endpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingConfigBasic(endpointGatewayID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVirtualEndpointGatewayResourceBindingExists("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", "endpoint_gateway_id", endpointGatewayID),
				),
			},
		},
	})
}

func TestAccIBMIsVirtualEndpointGatewayResourceBindingAllArgs(t *testing.T) {
	var conf vpcv1.EndpointGatewayResourceBinding
	endpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingConfig(endpointGatewayID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVirtualEndpointGatewayResourceBindingExists("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", "endpoint_gateway_id", endpointGatewayID),
					resource.TestCheckResourceAttr("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingConfig(endpointGatewayID, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", "endpoint_gateway_id", endpointGatewayID),
					resource.TestCheckResourceAttr("ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingConfigBasic(endpointGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_virtual_endpoint_gateway_resource_binding" "is_virtual_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}
	`, endpointGatewayID)
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingConfig(endpointGatewayID string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_virtual_endpoint_gateway_resource_binding" "is_virtual_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			name = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}
	`, endpointGatewayID, name)
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingExists(n string, obj vpcv1.EndpointGatewayResourceBinding) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getEndpointGatewayResourceBindingOptions := &vpcv1.GetEndpointGatewayResourceBindingOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getEndpointGatewayResourceBindingOptions.SetEndpointGatewayID(parts[0])
		getEndpointGatewayResourceBindingOptions.SetID(parts[1])

		endpointGatewayResourceBinding, _, err := vpcClient.GetEndpointGatewayResourceBinding(getEndpointGatewayResourceBindingOptions)
		if err != nil {
			return err
		}

		obj = *endpointGatewayResourceBinding
		return nil
	}
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_virtual_endpoint_gateway_resource_binding" {
			continue
		}

		getEndpointGatewayResourceBindingOptions := &vpcv1.GetEndpointGatewayResourceBindingOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getEndpointGatewayResourceBindingOptions.SetEndpointGatewayID(parts[0])
		getEndpointGatewayResourceBindingOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetEndpointGatewayResourceBinding(getEndpointGatewayResourceBindingOptions)

		if err == nil {
			return fmt.Errorf("EndpointGatewayResourceBinding still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for EndpointGatewayResourceBinding (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTarget)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.ResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetCRNToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTargetCRN)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.ResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetCRNToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.ResourceIBMIsVirtualEndpointGatewayResourceBindingEndpointGatewayResourceBindingLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsVirtualEndpointGatewayResourceBindingMapToEndpointGatewayResourceBindingTargetPrototype(t *testing.T) {
	checkResult := func(result vpcv1.EndpointGatewayResourceBindingTargetPrototypeIntf) {
		model := new(vpcv1.EndpointGatewayResourceBindingTargetPrototype)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

	result, err := vpc.ResourceIBMIsVirtualEndpointGatewayResourceBindingMapToEndpointGatewayResourceBindingTargetPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsVirtualEndpointGatewayResourceBindingMapToEndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.EndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN) {
		model := new(vpcv1.EndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

	result, err := vpc.ResourceIBMIsVirtualEndpointGatewayResourceBindingMapToEndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}
