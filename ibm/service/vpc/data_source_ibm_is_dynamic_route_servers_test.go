// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsDynamicRouteServersDataSourceBasic(t *testing.T) {
	dynamicRouteServerAsn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServersDataSourceConfigBasic(dynamicRouteServerAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.#"),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.asn", dynamicRouteServerAsn),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "total_count"),
				),
			},
		},
	})
}

func TestAccIBMIsDynamicRouteServersDataSourceAllArgs(t *testing.T) {
	dynamicRouteServerAsn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	dynamicRouteServerName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	dynamicRouteServerRedistributeServiceRoutes := "false"
	dynamicRouteServerRedistributeSubnets := "true"
	dynamicRouteServerRedistributeUserRoutes := "true"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServersDataSourceConfig(dynamicRouteServerAsn, dynamicRouteServerName, dynamicRouteServerRedistributeServiceRoutes, dynamicRouteServerRedistributeSubnets, dynamicRouteServerRedistributeUserRoutes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "resource_group_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "sort"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.#"),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.asn", dynamicRouteServerAsn),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.name", dynamicRouteServerName),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.redistribute_service_routes", dynamicRouteServerRedistributeServiceRoutes),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.redistribute_subnets", dynamicRouteServerRedistributeSubnets),
					resource.TestCheckResourceAttr("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.redistribute_user_routes", dynamicRouteServerRedistributeUserRoutes),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "dynamic_route_servers.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "next.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_servers.is_dynamic_route_servers", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsDynamicRouteServersDataSourceConfigBasic(dynamicRouteServerAsn string) string {
	return fmt.Sprintf(`
		resource "ibm_is_dynamic_route_server" "is_dynamic_route_server_instance" {
			asn = %s
			ips {
				address = "192.168.3.4"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				name = "my-reserved-ip"
				resource_type = "subnet_reserved_ip"
			}
			vpc {
				crn = "crn:v1:bluemix:public:is:us-south:a/123456::vpc:4727d842-f94f-4a2d-824a-9bc9b02c523b"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "4727d842-f94f-4a2d-824a-9bc9b02c523b"
				name = "my-vpc"
				resource_type = "vpc"
			}
		}

		data "ibm_is_dynamic_route_servers" "is_dynamic_route_servers_instance" {
			name = ibm_is_dynamic_route_server.is_dynamic_route_server.name
			resource_group.id = "resource_group.id"
			sort = "name"
		}
	`, dynamicRouteServerAsn)
}

func testAccCheckIBMIsDynamicRouteServersDataSourceConfig(dynamicRouteServerAsn string, dynamicRouteServerName string, dynamicRouteServerRedistributeServiceRoutes string, dynamicRouteServerRedistributeSubnets string, dynamicRouteServerRedistributeUserRoutes string) string {
	return fmt.Sprintf(`
		resource "ibm_is_dynamic_route_server" "is_dynamic_route_server_instance" {
			asn = %s
			ips {
				address = "192.168.3.4"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				name = "my-reserved-ip"
				resource_type = "subnet_reserved_ip"
			}
			name = "%s"
			redistribute_service_routes = %s
			redistribute_subnets = %s
			redistribute_user_routes = %s
			resource_group {
				href = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
				id = "fee82deba12e4c0fb69c3b09d1f12345"
				name = "my-resource-group"
			}
			security_groups {
				crn = "crn:v1:bluemix:public:is:us-south:a/123456::security-group:be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/security_groups/be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				id = "be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				name = "my-security-group"
			}
			vpc {
				crn = "crn:v1:bluemix:public:is:us-south:a/123456::vpc:4727d842-f94f-4a2d-824a-9bc9b02c523b"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "4727d842-f94f-4a2d-824a-9bc9b02c523b"
				name = "my-vpc"
				resource_type = "vpc"
			}
		}

		data "ibm_is_dynamic_route_servers" "is_dynamic_route_servers_instance" {
			name = ibm_is_dynamic_route_server.is_dynamic_route_server.name
			resource_group.id = "resource_group.id"
			sort = "name"
		}
	`, dynamicRouteServerAsn, dynamicRouteServerName, dynamicRouteServerRedistributeServiceRoutes, dynamicRouteServerRedistributeSubnets, dynamicRouteServerRedistributeUserRoutes)
}
