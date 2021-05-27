---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : vpc"
description: |-
  Manages IBM virtual private cloud.
---

# ibm\_is_vpc

Import the details of an existing IBM Virtual Private cloud as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```terraform
resource "ibm_is_vpc" "testacc_vpc" {
  name = "test"
}

data "ibm_is_vpc" "ds_vpc" {
  name = "test"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the VPC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `crn` - The CRN of VPC.
* `status` - The status of VPC.
* `default_network_acl` - ID of the default network ACL.
* `default_security_group` - The unique identifier of the VPC default security group.
* `default_routing_table` - The unique identifier of the VPC default routing table.
* `classic_access` - Indicates whether this VPC is connected to Classic Infrastructure.
* `resource_group` - The resource group ID where the VPC created.
* `tags` - Tags associated with the instance.
* `cse_source_addresses` - A list describing the cloud service endpoint source ip adresses and zones. The nested cse_source_addresses block have the following structure:
  * `address` - Ip Address of the cloud service endpoint.
  * `zone_name` - Zone associated with the IP Address.
* `subnets` - A list of subnets attached to VPC. The nested subnets block have the following structure:
  * `name` - Name of the subnet.
  * `id` - ID of the subnet.
  * `status` -  Status of the subnet.
  * `zone` -  Zone of the subnet.
  * `total_ipv4_address_count` - Total IPv4 addresses under the subnet.
  * `available_ipv4_address_count` - Available IPv4 addresses available for the usage in the subnet.
* `security_group` - A list of security groups attached to VPC. The nested security group block has the following structure:
  * `group_id` - Security group ID.
  * `group_name` - Name of the security group.
  * `rules` -  Set of rules attached to a security group
    * `rule_id` - ID of the rule
    * `direction` - Direction of the traffic either inbound or outbound
    * `ip_version` - ip version either ipv4 or ipv6
    * `remote` - Security group id, an IP address, a CIDR block, or a single security group identifier.
    * `type` - The ICMP traffic type to allow.
    * `code` - The ICMP traffic code to allow.
    * `port_min` - The inclusive lower bound of TCP port range. 
    * `port_max` - The inclusive upper bound of TCP port range. 
* `default_network_acl_name` - The name of the default network acl.
* `default_security_group_name` - The name of the default security group.
* `default_routing_table_name` - The name of the default routing table.
