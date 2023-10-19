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

func TestAccIBMIsVPCDnsResolutionBindingResourceBasic(t *testing.T) {
	vpcname1 := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	enable_hub1 := true
	vpcname2 := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	bindingname := fmt.Sprintf("tf-vpc-dns-binding-%d", acctest.RandIntRange(10, 100))
	enable_hub2 := false
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPCDnsResolutionBindingResourceConfigBasic(vpcname1, vpcname2, bindingname, enable_hub1, enable_hub2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc_id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "created_at"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "href"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "name"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "resource_type"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc.#"),
				),
			},
		},
	})
}
func TestAccIBMIsVPCDnsResolutionBindingResourceForceDelete(t *testing.T) {
	name1 := fmt.Sprintf("terraformvpcuat-%d", acctest.RandIntRange(10, 100))
	name2 := fmt.Sprintf("terraformvpcuat-%d", acctest.RandIntRange(10, 100))
	subnet1 := fmt.Sprintf("terraformsubnet-%d", acctest.RandIntRange(10, 100))
	subnet2 := fmt.Sprintf("terraformsubnet-%d", acctest.RandIntRange(10, 100))
	subnet3 := fmt.Sprintf("terraformsubnet-%d", acctest.RandIntRange(10, 100))
	subnet4 := fmt.Sprintf("terraformsubnet-%d", acctest.RandIntRange(10, 100))
	resourecinstance := fmt.Sprintf("terraformresource-%d", acctest.RandIntRange(10, 100))
	resolver1 := fmt.Sprintf("terraformresolver-%d", acctest.RandIntRange(10, 100))
	resolver2 := fmt.Sprintf("terraformresolver-%d", acctest.RandIntRange(10, 100))
	binding := fmt.Sprintf("terraformbinding-%d", acctest.RandIntRange(10, 100))
	enableHubTrue := true
	enableHubFalse := false
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPCDnsResolutionBindingForceDeleteResourceConfig(name1, name2, subnet1, subnet2, subnet3, subnet4, resourecinstance, resolver1, resolver2, binding, enableHubTrue, enableHubFalse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc_id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "created_at"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "href"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "name"),
					resource.TestCheckResourceAttr("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "name", binding),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "resource_type"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc.#"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "name", name1),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "name", name2),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "dns.0.enable_hub", fmt.Sprintf("%t", enableHubTrue)),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "dns.0.resolver.0.type", "system"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "dns.0.enable_hub", fmt.Sprintf("%t", enableHubFalse)),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "dns.0.resolver.0.type", "system"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsVPCDnsResolutionBindingForceDeleteUpdateResourceConfig(name1, name2, subnet1, subnet2, subnet3, subnet4, resourecinstance, resolver1, resolver2, binding, enableHubTrue, enableHubFalse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc_id"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "created_at"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "href"),
					resource.TestCheckResourceAttr("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "name", binding),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "name"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "resource_type"),
					resource.TestCheckResourceAttrSet("ibm_is_vpc_dns_resolution_binding.is_vpc_dns_resolution_binding", "vpc.#"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "name", name1),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "name", name2),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "dns.0.resolution_binding_count", "1"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "dns.0.enable_hub", fmt.Sprintf("%t", enableHubTrue)),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_true", "dns.0.resolver.0.type", "system"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "dns.0.enable_hub", fmt.Sprintf("%t", enableHubFalse)),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "dns.0.resolver.0.type", "delegated"),
					resource.TestCheckResourceAttr(
						"ibm_is_vpc.hub_false_delegated", "dns.0.resolution_binding_count", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPCDnsResolutionBindingForceDeleteResourceConfig(vpcname, vpcname2, subnetname1, subnetname2, subnetname3, subnetname4, resourceinstance, resolver1, resolver2, bindingname string, enableHub, enablehubfalse bool) string {
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		is_default	   =  true
	}
	
	resource ibm_is_vpc hub_true {
		name = "%s"
		dns {
			enable_hub = %t
		}
	}
	
	resource ibm_is_vpc hub_false_delegated {
		name = "%s"
		dns {
			enable_hub = %t
		}
	}
	
	resource "ibm_is_subnet" "hub_true_sub1" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_true.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_true_sub2" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_true.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_false_delegated_sub1" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_false_delegated.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_false_delegated_sub2" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_false_delegated.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_resource_instance" "dns-cr-instance" {
		name		   		=  "%s"
		resource_group_id  	=  data.ibm_resource_group.rg.id
		location           	=  "global"
		service		   		=  "dns-svcs"
		plan		   		=  "standard-dns"
	}
	resource "ibm_dns_custom_resolver" "test_hub_true" {
		name		   		=  "%s"
		instance_id 	   	=  ibm_resource_instance.dns-cr-instance.guid
		description	   		=  "new test CR - TF"
		high_availability  	=  true
		enabled 	   		=  true
		locations {
				subnet_crn  = ibm_is_subnet.hub_true_sub1.crn
				enabled	 = true
		}
		locations {
				subnet_crn  = ibm_is_subnet.hub_true_sub2.crn
				enabled	 = true
		}
	}
	resource "ibm_dns_custom_resolver" "test_hub_false_delegated" {
		name		   		=  "%s"
		instance_id 	   	=  ibm_resource_instance.dns-cr-instance.guid
		description	   		=  "new test CR - TF"
		high_availability  	=  true
		enabled 	   		=  true
		locations {
				subnet_crn  = ibm_is_subnet.hub_false_delegated_sub1.crn
				enabled	 = true
		}
		locations {
				subnet_crn  = ibm_is_subnet.hub_false_delegated_sub2.crn
				enabled	 = true
		}
	}
	
	resource ibm_is_vpc_dns_resolution_binding is_vpc_dns_resolution_binding {
		name = "%s"
		vpc_id=  ibm_is_vpc.hub_false_delegated.id
		vpc {
			id = ibm_is_vpc.hub_true.id
		}
		force_delete = true
	}
	`, vpcname, enableHub, vpcname2, enablehubfalse, subnetname1, acc.ISZoneName, subnetname2, acc.ISZoneName, subnetname3, acc.ISZoneName, subnetname4, acc.ISZoneName, resourceinstance, resolver1, resolver2, bindingname)
}
func testAccCheckIBMIsVPCDnsResolutionBindingForceDeleteUpdateResourceConfig(vpcname, vpcname2, subnetname1, subnetname2, subnetname3, subnetname4, resourceinstance, resolver1, resolver2, bindingname string, enableHub, enablehubfalse bool) string {
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		is_default	   =  true
	}
	
	resource ibm_is_vpc hub_true {
		name = "%s"
		dns {
			enable_hub = %t
		}
	}
	
	resource ibm_is_vpc hub_false_delegated {
		name = "%s"
		dns {
			enable_hub = %t
			resolver {
				type = "delegated"
				vpc_id = ibm_is_vpc.hub_true.id
			}
		}
	}
	
	resource "ibm_is_subnet" "hub_true_sub1" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_true.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_true_sub2" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_true.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_false_delegated_sub1" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_false_delegated.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_is_subnet" "hub_false_delegated_sub2" {
		name		   				=  "%s"
		vpc      	   				=  ibm_is_vpc.hub_false_delegated.id
		zone		   				=  "%s"
		total_ipv4_address_count 	= 16
	}
	resource "ibm_resource_instance" "dns-cr-instance" {
		name		   		=  "%s"
		resource_group_id  	=  data.ibm_resource_group.rg.id
		location           	=  "global"
		service		   		=  "dns-svcs"
		plan		   		=  "standard-dns"
	}
	resource "ibm_dns_custom_resolver" "test_hub_true" {
		name		   		=  "%s"
		instance_id 	   	=  ibm_resource_instance.dns-cr-instance.guid
		description	   		=  "new test CR - TF"
		high_availability  	=  true
		enabled 	   		=  true
		locations {
				subnet_crn  = ibm_is_subnet.hub_true_sub1.crn
				enabled	 = true
		}
		locations {
				subnet_crn  = ibm_is_subnet.hub_true_sub2.crn
				enabled	 = true
		}
	}
	resource "ibm_dns_custom_resolver" "test_hub_false_delegated" {
		name		   		=  "%s"
		instance_id 	   	=  ibm_resource_instance.dns-cr-instance.guid
		description	   		=  "new test CR - TF"
		high_availability  	=  true
		enabled 	   		=  true
		locations {
				subnet_crn  = ibm_is_subnet.hub_false_delegated_sub1.crn
				enabled	 = true
		}
		locations {
				subnet_crn  = ibm_is_subnet.hub_false_delegated_sub2.crn
				enabled	 = true
		}
	}
	
	resource ibm_is_vpc_dns_resolution_binding is_vpc_dns_resolution_binding {
		name = "%s"
		vpc_id=  ibm_is_vpc.hub_false_delegated.id
		vpc {
			id = ibm_is_vpc.hub_true.id
		}
		force_delete = true
	}
	
	`, vpcname, enableHub, vpcname2, enablehubfalse, subnetname1, acc.ISZoneName, subnetname2, acc.ISZoneName, subnetname3, acc.ISZoneName, subnetname4, acc.ISZoneName, resourceinstance, resolver1, resolver2, bindingname)

}
func testAccCheckIBMIsVPCDnsResolutionBindingResourceConfigBasic(vpcname1, vpcname2, bindingname string, enablehub1, enablehub2 bool) string {
	return fmt.Sprintf(`
	resource ibm_is_vpc testacc_vpc1 {
		name = "%s"
		dns {
			enable_hub = %t
		}
	}
	resource ibm_is_vpc testacc_vpc2 {
		name = "%s"
		dns {
			enable_hub = %t
		}
	}
	resource ibm_is_vpc_dns_resolution_binding is_vpc_dns_resolution_binding {
		name = "%s"
		vpc_id=  ibm_is_vpc.testacc_vpc2.id
		vpc {
			id = ibm_is_vpc.testacc_vpc1.id
		}
	}
	`, vpcname1, enablehub1, vpcname2, enablehub2, bindingname)
}
