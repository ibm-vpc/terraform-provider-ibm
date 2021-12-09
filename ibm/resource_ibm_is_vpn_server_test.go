// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsVPNServerBasic(t *testing.T) {
	var conf vpcv1.VPNServer
	clientIPPool := "10.5.0.0/21"
	clientIdleTimeout := fmt.Sprintf("%d", acctest.RandIntRange(0, 28800))
	enableSplitTunneling := "true"
	name := fmt.Sprintf("tfname%d", acctest.RandIntRange(10, 100))
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 65535))
	protocol := "udp"

	clientIPPoolUpdate := "10.6.0.0/21"
	clientIdleTimeoutUpdate := fmt.Sprintf("%d", acctest.RandIntRange(0, 28800))
	enableSplitTunnelingUpdate := "false"
	nameUpdate := fmt.Sprintf("tfname%d", acctest.RandIntRange(10, 100))
	portUpdate := fmt.Sprintf("%d", acctest.RandIntRange(1, 65535))
	protocolUpdate := "tcp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIsVPNServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerConfigBasic(clientIPPool, clientIdleTimeout, enableSplitTunneling, name, port, protocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVPNServerExists("ibm_is_vpn_server.is_vpn_server", conf),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_ip_pool", clientIPPool),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_idle_timeout", clientIdleTimeout),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "enable_split_tunneling", enableSplitTunneling),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "name", name),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "port", port),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "protocol", protocol),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerConfigBasic(clientIPPoolUpdate, clientIdleTimeoutUpdate, enableSplitTunnelingUpdate, nameUpdate, portUpdate, protocolUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_ip_pool", clientIPPoolUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_idle_timeout", clientIdleTimeoutUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "enable_split_tunneling", enableSplitTunnelingUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "port", portUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "protocol", protocolUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsVPNServerAllArgs(t *testing.T) {
	var conf vpcv1.VPNServer
	clientIPPool := "10.5.0.0/21"
	clientIdleTimeout := fmt.Sprintf("%d", acctest.RandIntRange(0, 28800))
	enableSplitTunneling := "true"
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 65535))
	protocol := "udp"

	clientIPPoolUpdate := "10.6.0.0/21"
	clientIdleTimeoutUpdate := fmt.Sprintf("%d", acctest.RandIntRange(0, 28800))
	enableSplitTunnelingUpdate := "false"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	portUpdate := fmt.Sprintf("%d", acctest.RandIntRange(1, 65535))
	protocolUpdate := "tcp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIsVPNServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerConfig(clientIPPool, clientIdleTimeout, enableSplitTunneling, name, port, protocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVPNServerExists("ibm_is_vpn_server.is_vpn_server", conf),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_ip_pool", clientIPPool),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_idle_timeout", clientIdleTimeout),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "enable_split_tunneling", enableSplitTunneling),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "name", name),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "port", port),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "protocol", protocol),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerConfig(clientIPPoolUpdate, clientIdleTimeoutUpdate, enableSplitTunnelingUpdate, nameUpdate, portUpdate, protocolUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_ip_pool", clientIPPoolUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "client_idle_timeout", clientIdleTimeoutUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "enable_split_tunneling", enableSplitTunnelingUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "port", portUpdate),
					resource.TestCheckResourceAttr("ibm_is_vpn_server.is_vpn_server", "protocol", protocolUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_vpn_server.is_vpn_server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsVPNServerConfigBasic(clientIPPool string, clientIdleTimeout string, enableSplitTunneling string, name string, port string, protocol string) string {
	return fmt.Sprintf(`
		
	resource "ibm_is_vpn_server" "is_vpn_server" {

			certificate_crn = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:00406b5615f95dba9bf7c2ab52bb3083"
			client_authentication {
				method = "certificate"
				client_ca_crn = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:6a85a87d01dd5a4268a8bca16cb998eb"
			}
			client_ip_pool = "%s"
			subnets = ["0726-61b2f53f-1e95-42a7-94ab-55de8f8cbdd5"]
			client_dns_server_ips = ["192.168.3.4"]
			client_idle_timeout = %s
			enable_split_tunneling = %s
			name = "%s"
			port = %s
			protocol = "%s"
		}
	`, clientIPPool, clientIdleTimeout, enableSplitTunneling, name, port, protocol)
}

func testAccCheckIBMIsVPNServerConfig(clientIPPool string, clientIdleTimeout string, enableSplitTunneling string, name string, port string, protocol string) string {
	return fmt.Sprintf(`

		resource "ibm_is_vpn_server" "is_vpn_server" {
			certificate_crn = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:00406b5615f95dba9bf7c2ab52bb3083"
			client_authentication {
				method = "certificate"
				client_ca_crn = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:6a85a87d01dd5a4268a8bca16cb998eb"
			}
			client_ip_pool = "%s"
			subnets = ["0726-61b2f53f-1e95-42a7-94ab-55de8f8cbdd5"]
			client_dns_server_ips = ["192.168.3.4"]
			client_idle_timeout = %s
			enable_split_tunneling = %s
			name = "%s"
			port = %s
			protocol = "%s"
			security_groups = ["r134-bd0f8527-c45c-496e-8d34-63bda6dd829b"]
	`, clientIPPool, clientIdleTimeout, enableSplitTunneling, name, port, protocol)
}

func testAccCheckIBMIsVPNServerExists(n string, obj vpcv1.VPNServer) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		sess, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
		getVPNServerOptions := &vpcv1.GetVPNServerOptions{}

		getVPNServerOptions.SetID(rs.Primary.ID)

		vpnServer, _, err := sess.GetVPNServer(getVPNServerOptions)
		if err != nil {
			return err
		}

		obj = *vpnServer
		return nil
	}
}

func testAccCheckIBMIsVPNServerDestroy(s *terraform.State) error {
	sess, _ := testAccProvider.Meta().(ClientSession).VpcV1API()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_vpn_server" {
			continue
		}

		getVPNServerOptions := &vpcv1.GetVPNServerOptions{}

		getVPNServerOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := sess.GetVPNServer(getVPNServerOptions)

		if err == nil {
			return fmt.Errorf("VPNServer still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for VPNServer (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

// resource "ibm_is_vpc" "testacc_vpc" {
// 	name = "sunithavpc"
// }

// resource "ibm_is_subnet" "testacc_subnet" {
// 	name = "sunithasubnet"
// 	vpc = ibm_is_vpc.testacc_vpc.id
// 	zone = "us-south-1"
// 	ipv4_cidr_block = "10.240.0.0/24"
// 	tags = ["Tag1", "tag2"]
// }
