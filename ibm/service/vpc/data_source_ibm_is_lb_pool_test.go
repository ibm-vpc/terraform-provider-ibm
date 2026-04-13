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

func TestAccIBMIsLbPoolDataSourceBasic(t *testing.T) {
	vpcname := fmt.Sprintf("tflbp-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflbpc-name-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tfcreate%d", acctest.RandIntRange(10, 100))
	poolName := fmt.Sprintf("tflbpoolc%d", acctest.RandIntRange(10, 100))
	alg1 := "round_robin"
	protocol1 := "http"
	delay1 := "45"
	retries1 := "5"
	timeout1 := "15"
	healthType1 := "http"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsLbPoolDataSourceConfigBasic(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name, poolName, alg1, protocol1, delay1, retries1, timeout1, healthType1),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "lb"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "identifier"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "algorithm"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "health_monitor.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "provisioning_status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool", "proxy_protocol"),
				),
			},
		},
	})
}

func TestAccIBMIsLbPoolDataSource_mTLS(t *testing.T) {
	vpcname := fmt.Sprintf("tflbp-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflbpc-name-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tfcreate%d", acctest.RandIntRange(10, 100))
	poolName := fmt.Sprintf("tflbpoolc%d", acctest.RandIntRange(10, 100))
	alg1 := "round_robin"
	protocol1 := "https"
	delay1 := "45"
	retries1 := "5"
	timeout1 := "15"
	healthType1 := "https"

	// Example CRNs - replace with actual values from your test environment
	clientCertCRN := "crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"
	serverCACRN := "crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:3f897bbc-53gb-74fe-25g9-e77e663g5ee6"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsLbPoolDataSourceConfigmTLS(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name, poolName, alg1, protocol1, delay1, retries1, timeout1, healthType1, clientCertCRN, serverCACRN),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "lb"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "identifier"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_pool.is_lb_pool_mtls", "protocol", protocol1),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "client_authentication.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "client_authentication.0.certificate_instance.#"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_pool.is_lb_pool_mtls", "client_authentication.0.certificate_instance.0.crn", clientCertCRN),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "server_authentication.#"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_pool.is_lb_pool_mtls", "server_authentication.0.verify_certificate", "true"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb_pool.is_lb_pool_mtls", "server_authentication.0.certificate_authority.#"),
					resource.TestCheckResourceAttr("data.ibm_is_lb_pool.is_lb_pool_mtls", "server_authentication.0.certificate_authority.0.crn", serverCACRN),
				),
			},
		},
	})
}

func testAccCheckIBMIsLbPoolDataSourceConfigBasic(vpcname, subnetname, zone, cidr, name, poolName, algorithm, protocol, delay, retries, timeout, healthType string) string {
	return testAccCheckIBMISLBPoolConfig(vpcname, subnetname, zone, cidr, name, poolName, algorithm, protocol, delay, retries, timeout, healthType) + fmt.Sprintf(`
        data "ibm_is_lb_pool" "is_lb_pool" {
            lb = "${ibm_is_lb.testacc_LB.id}"
            identifier = "${element(split("/",ibm_is_lb_pool.testacc_lb_pool.id),1)}"
        }
    `)
}

func testAccCheckIBMIsLbPoolDataSourceConfigmTLS(vpcname, subnetname, zone, cidr, name, poolName, algorithm, protocol, delay, retries, timeout, healthType, clientCertCRN, serverCACRN string) string {
	return testAccCheckIBMISLBPoolConfig(vpcname, subnetname, zone, cidr, name, poolName, "round_robin", "http", delay, retries, timeout, "http") + fmt.Sprintf(`
	
	resource "ibm_is_lb_pool" "testacc_lb_pool_mtls" {
		name = "%s-mtls"
		lb = ibm_is_lb.testacc_LB.id
		algorithm = "%s"
		protocol = "%s"
		health_delay = %s
		health_retries = %s
		health_timeout = %s
		health_type = "%s"
		client_authentication {
			certificate_instance = "%s"
		}
		server_authentication {
			verify_certificate = true
			certificate_authority = "%s"
		}
	}

	data "ibm_is_lb_pool" "is_lb_pool_mtls" {
		lb = ibm_is_lb.testacc_LB.id
		identifier = element(split("/", ibm_is_lb_pool.testacc_lb_pool_mtls.id), 1)
	}
	`, poolName, algorithm, protocol, delay, retries, timeout, healthType, clientCertCRN, serverCACRN)
}
