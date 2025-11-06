// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.107.1-41b0fbd0-20250825-080732
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsVirtualEndpointGatewayResourceBindingsDataSourceBasic(t *testing.T) {
	endpointGatewayResourceBindingEndpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingsDataSourceConfigBasic(endpointGatewayResourceBindingEndpointGatewayID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "endpoint_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.#"),
				),
			},
		},
	})
}

func TestAccIBMIsVirtualEndpointGatewayResourceBindingsDataSourceAllArgs(t *testing.T) {
	endpointGatewayResourceBindingEndpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))
	endpointGatewayResourceBindingName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualEndpointGatewayResourceBindingsDataSourceConfig(endpointGatewayResourceBindingEndpointGatewayID, endpointGatewayResourceBindingName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "endpoint_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.name", endpointGatewayResourceBindingName),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.service_endpoint"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_endpoint_gateway_resource_bindings.is_virtual_endpoint_gateway_resource_bindings_instance", "resource_bindings.0.type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingsDataSourceConfigBasic(endpointGatewayResourceBindingEndpointGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_virtual_endpoint_gateway_resource_binding" "is_virtual_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}

		data "ibm_is_virtual_endpoint_gateway_resource_bindings" "is_virtual_endpoint_gateway_resource_bindings_instance" {
			endpoint_gateway_id = ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance.endpoint_gateway_id
		}
	`, endpointGatewayResourceBindingEndpointGatewayID)
}

func testAccCheckIBMIsVirtualEndpointGatewayResourceBindingsDataSourceConfig(endpointGatewayResourceBindingEndpointGatewayID string, endpointGatewayResourceBindingName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_virtual_endpoint_gateway_resource_binding" "is_virtual_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			name = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}

		data "ibm_is_virtual_endpoint_gateway_resource_bindings" "is_virtual_endpoint_gateway_resource_bindings_instance" {
			endpoint_gateway_id = ibm_is_virtual_endpoint_gateway_resource_binding.is_virtual_endpoint_gateway_resource_binding_instance.endpoint_gateway_id
		}
	`, endpointGatewayResourceBindingEndpointGatewayID, endpointGatewayResourceBindingName)
}

func TestDataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		endpointGatewayResourceBindingLifecycleReasonModel := make(map[string]interface{})
		endpointGatewayResourceBindingLifecycleReasonModel["code"] = "resource_suspended_by_provider"
		endpointGatewayResourceBindingLifecycleReasonModel["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		endpointGatewayResourceBindingLifecycleReasonModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		endpointGatewayResourceBindingTargetModel := make(map[string]interface{})
		endpointGatewayResourceBindingTargetModel["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		model := make(map[string]interface{})
		model["created_at"] = "2025-10-19T11:59:46Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/endpoint_gateways/r006-7610ebfb-f5dc-4d42-bc18-287d47f7a5b0/resource_bindings/r006-a7ba95b6-a254-47e4-b129-10593df8a373"
		model["id"] = "r006-a7ba95b6-a254-47e4-b129-10593df8a373"
		model["lifecycle_reasons"] = []map[string]interface{}{endpointGatewayResourceBindingLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["name"] = "my-resource-binding"
		model["resource_type"] = "endpoint_gateway_resource_binding"
		model["service_endpoint"] = "bucket-27200-lwx4cfvcue.s3.direct.us-south.cloud-object-storage.appdomain.cloud"
		model["target"] = []map[string]interface{}{endpointGatewayResourceBindingTargetModel}
		model["type"] = "weak"

		assert.Equal(t, result, model)
	}

	endpointGatewayResourceBindingLifecycleReasonModel := new(vpcv1.EndpointGatewayResourceBindingLifecycleReason)
	endpointGatewayResourceBindingLifecycleReasonModel.Code = core.StringPtr("resource_suspended_by_provider")
	endpointGatewayResourceBindingLifecycleReasonModel.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	endpointGatewayResourceBindingLifecycleReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	endpointGatewayResourceBindingTargetModel := new(vpcv1.EndpointGatewayResourceBindingTargetCRN)
	endpointGatewayResourceBindingTargetModel.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	model := new(vpcv1.EndpointGatewayResourceBinding)
	model.CreatedAt = CreateMockDateTime("2025-10-19T11:59:46Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/endpoint_gateways/r006-7610ebfb-f5dc-4d42-bc18-287d47f7a5b0/resource_bindings/r006-a7ba95b6-a254-47e4-b129-10593df8a373")
	model.ID = core.StringPtr("r006-a7ba95b6-a254-47e4-b129-10593df8a373")
	model.LifecycleReasons = []vpcv1.EndpointGatewayResourceBindingLifecycleReason{*endpointGatewayResourceBindingLifecycleReasonModel}
	model.LifecycleState = core.StringPtr("stable")
	model.Name = core.StringPtr("my-resource-binding")
	model.ResourceType = core.StringPtr("endpoint_gateway_resource_binding")
	model.ServiceEndpoint = core.StringPtr("bucket-27200-lwx4cfvcue.s3.direct.us-south.cloud-object-storage.appdomain.cloud")
	model.Target = endpointGatewayResourceBindingTargetModel
	model.Type = core.StringPtr("weak")

	result, err := vpc.DataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTarget)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.DataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingTargetCRNToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTargetCRN)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.DataSourceIBMIsVirtualEndpointGatewayResourceBindingsEndpointGatewayResourceBindingTargetCRNToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
