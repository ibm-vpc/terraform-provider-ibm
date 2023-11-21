// Copyright IBM Corp. 2023 All Rights Reserved.
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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsDynamicRouteServerBasic(t *testing.T) {
	var conf vpcv1.DynamicRouteServer
	asn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	asnUpdate := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsDynamicRouteServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerConfigBasic(asn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsDynamicRouteServerExists("ibm_is_dynamic_route_server.is_dynamic_route_server", conf),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "asn", asn),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerConfigBasic(asnUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "asn", asnUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsDynamicRouteServerAllArgs(t *testing.T) {
	var conf vpcv1.DynamicRouteServer
	asn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	redistributeServiceRoutes := "false"
	redistributeSubnets := "true"
	redistributeUserRoutes := "true"
	asnUpdate := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	redistributeServiceRoutesUpdate := "true"
	redistributeSubnetsUpdate := "false"
	redistributeUserRoutesUpdate := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsDynamicRouteServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerConfig(asn, name, redistributeServiceRoutes, redistributeSubnets, redistributeUserRoutes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsDynamicRouteServerExists("ibm_is_dynamic_route_server.is_dynamic_route_server", conf),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "asn", asn),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "name", name),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_service_routes", redistributeServiceRoutes),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_subnets", redistributeSubnets),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_user_routes", redistributeUserRoutes),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerConfig(asnUpdate, nameUpdate, redistributeServiceRoutesUpdate, redistributeSubnetsUpdate, redistributeUserRoutesUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "asn", asnUpdate),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_service_routes", redistributeServiceRoutesUpdate),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_subnets", redistributeSubnetsUpdate),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server.is_dynamic_route_server", "redistribute_user_routes", redistributeUserRoutesUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_dynamic_route_server.is_dynamic_route_server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsDynamicRouteServerConfigBasic(asn string) string {
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
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
		}
	`, asn)
}

func testAccCheckIBMIsDynamicRouteServerConfig(asn string, name string, redistributeServiceRoutes string, redistributeSubnets string, redistributeUserRoutes string) string {
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
				id = "fee82deba12e4c0fb69c3b09d1f12345"
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
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
		}
	`, asn, name, redistributeServiceRoutes, redistributeSubnets, redistributeUserRoutes)
}

func testAccCheckIBMIsDynamicRouteServerExists(n string, obj vpcv1.DynamicRouteServer) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getDynamicRouteServerOptions := &vpcv1.GetDynamicRouteServerOptions{}

		getDynamicRouteServerOptions.SetID(rs.Primary.ID)

		dynamicRouteServer, _, err := vpcClient.GetDynamicRouteServer(getDynamicRouteServerOptions)
		if err != nil {
			return err
		}

		obj = *dynamicRouteServer
		return nil
	}
}

func testAccCheckIBMIsDynamicRouteServerDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_dynamic_route_server" {
			continue
		}

		getDynamicRouteServerOptions := &vpcv1.GetDynamicRouteServerOptions{}

		getDynamicRouteServerOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetDynamicRouteServer(getDynamicRouteServerOptions)

		if err == nil {
			return fmt.Errorf("DynamicRouteServer still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for DynamicRouteServer (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
