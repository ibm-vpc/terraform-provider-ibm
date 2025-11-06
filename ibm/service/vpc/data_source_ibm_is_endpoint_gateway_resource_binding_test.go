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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsEndpointGatewayResourceBindingDataSourceBasic(t *testing.T) {
	endpointGatewayResourceBindingEndpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsEndpointGatewayResourceBindingDataSourceConfigBasic(endpointGatewayResourceBindingEndpointGatewayID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "endpoint_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "is_endpoint_gateway_resource_binding_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "service_endpoint"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "target.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "type"),
				),
			},
		},
	})
}

func TestAccIBMIsEndpointGatewayResourceBindingDataSourceAllArgs(t *testing.T) {
	endpointGatewayResourceBindingEndpointGatewayID := fmt.Sprintf("tf_endpoint_gateway_id_%d", acctest.RandIntRange(10, 100))
	endpointGatewayResourceBindingName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsEndpointGatewayResourceBindingDataSourceConfig(endpointGatewayResourceBindingEndpointGatewayID, endpointGatewayResourceBindingName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "endpoint_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "is_endpoint_gateway_resource_binding_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_reasons.0.code"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_reasons.0.message"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_reasons.0.more_info"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "service_endpoint"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "target.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance", "type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsEndpointGatewayResourceBindingDataSourceConfigBasic(endpointGatewayResourceBindingEndpointGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_endpoint_gateway_resource_binding" "is_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}

		data "ibm_is_endpoint_gateway_resource_binding" "is_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance.endpoint_gateway_id
			is_endpoint_gateway_resource_binding_id = ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance.is_endpoint_gateway_resource_binding_id
		}
	`, endpointGatewayResourceBindingEndpointGatewayID)
}

func testAccCheckIBMIsEndpointGatewayResourceBindingDataSourceConfig(endpointGatewayResourceBindingEndpointGatewayID string, endpointGatewayResourceBindingName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_endpoint_gateway_resource_binding" "is_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = "%s"
			name = "%s"
			target {
				crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"
			}
		}

		data "ibm_is_endpoint_gateway_resource_binding" "is_endpoint_gateway_resource_binding_instance" {
			endpoint_gateway_id = ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance.endpoint_gateway_id
			is_endpoint_gateway_resource_binding_id = ibm_is_endpoint_gateway_resource_binding.is_endpoint_gateway_resource_binding_instance.is_endpoint_gateway_resource_binding_id
		}
	`, endpointGatewayResourceBindingEndpointGatewayID, endpointGatewayResourceBindingName)
}

func TestDataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTarget)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.DataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetCRNToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.EndpointGatewayResourceBindingTargetCRN)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue")

	result, err := vpc.DataSourceIBMIsEndpointGatewayResourceBindingEndpointGatewayResourceBindingTargetCRNToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
