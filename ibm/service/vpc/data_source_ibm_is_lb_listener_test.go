// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsLbListenerDataSourceBasic(t *testing.T) {
	vpcname := fmt.Sprintf("tflblis-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflblis-subnet-%d", acctest.RandIntRange(10, 100))
	lbname := fmt.Sprintf("tflblis%d", acctest.RandIntRange(10, 100))
	protocol1 := "http"
	port1 := "8080"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsLbListenerDataSourceConfigBasic(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, lbname, port1, protocol1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "listener_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "lb"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "accept_proxy_protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "port"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "port_max"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "port_min"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener", "provisioning_status"),
				),
			},
		},
	})
}

func TestAccIBMIsLbListenerDataSource_ClientAuth(t *testing.T) {
	vpcname := fmt.Sprintf("tflblis-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflblis-subnet-%d", acctest.RandIntRange(10, 100))
	lbname := fmt.Sprintf("tflblis%d", acctest.RandIntRange(10, 100))
	protocol := "https"
	port := "443"

	// Example CRNs must be replaced
	certCRN := "crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"
	caCRN := "crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:3f897bbc-53gb-74fe-25g9-e77e663g5ee6"
	crlContent := "-----BEGIN X509 CRL-----\nMIIBpjCBjwIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdUZXN0IENBGA8y\nMDI0MDEwMTAwMDAwMFoYDzIwMjUwMTAxMDAwMDAwWjAkMCICEQDExample1RevokedCert\nGA8yMDI0MDEwMTAwMDAwMFqgDjAMMAoGA1UdFAQDAgEBMA0GCSqGSIb3DQEBCwUA\nA4IBAQCExample+CRLSignature==\n-----END X509 CRL-----"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsLbListenerDataSourceConfigClientAuth(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, lbname, port, protocol, certCRN, caCRN, crlContent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener_mtls", "listener_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener_mtls", "lb"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener_mtls", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_listener.is_lb_listener_mtls", "protocol", protocol),
					resource.TestCheckResourceAttr("data.ibm_is_lb_listener.is_lb_listener_mtls", "port", port),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener_mtls", "client_authentication.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_listener.is_lb_listener_mtls", "client_authentication.0.certificate_authority.#"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_listener.is_lb_listener_mtls", "client_authentication.0.certificate_authority.0.crn", caCRN),
					resource.TestCheckResourceAttr("data.ibm_is_lb_listener.is_lb_listener_mtls", "client_authentication.0.certificate_revocation_list", crlContent),
				),
			},
		},
	})
}

func testAccCheckIBMIsLbListenerDataSourceConfigBasic(vpcname, subnetname, zone, cidr, lbname, port, protocol string) string {
	return testAccCheckIBMISLBListenerConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, lbname, port, protocol) + fmt.Sprintf(`

	data "ibm_is_lb_listener" "is_lb_listener" {
		lb = "${ibm_is_lb.testacc_LB.id}"
		listener_id = ibm_is_lb_listener.testacc_lb_listener.listener_id
	}
	`)
}

func testAccCheckIBMIsLbListenerDataSourceConfigClientAuth(vpcname, subnetname, zone, cidr, lbname, port, protocol, certCRN, caCRN, crlContent string) string {
	return testAccCheckIBMISLBListenerConfig(vpcname, subnetname, zone, cidr, lbname, "8080", "http") + fmt.Sprintf(`
	
	resource "ibm_is_lb_listener" "testacc_lb_listener_mtls" {
		lb       = ibm_is_lb.testacc_LB.id
		port     = %s
		protocol = "%s"
		certificate_instance = "%s"
		client_authentication {
			certificate_authority = "%s"
			certificate_revocation_list = "%s"
		}
	}

	data "ibm_is_lb_listener" "is_lb_listener_mtls" {
		lb = ibm_is_lb.testacc_LB.id
		listener_id = ibm_is_lb_listener.testacc_lb_listener_mtls.listener_id
	}
	`, port, protocol, certCRN, caCRN, crlContent)
}
