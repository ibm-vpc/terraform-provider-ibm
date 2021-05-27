---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : vpc"
description: |-
  Manages IBM virtual private cloud.
---

# ibm\_is_vpc

Provides a vpc resource. This allows VPC to be created, updated, and cancelled.


## Example Usage

In the following example, you can create a VPC:

```terraform
resource "ibm_is_vpc" "testacc_vpc" {
  name = "test"
}

```

## Timeouts

ibm_is_vpc provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 10 minutes) Used for Creating vpc.
* `delete` - (Default 10 minutes) Used for Deleting vpc.


## Argument Reference

The following arguments are supported:

* `default_network_acl` - (Deprecated, string) ID of the default network ACL.
* `address_prefix_management` - (Optional, string) Indicates whether a default address prefix should be automatically created for each zone in this VPC. Value `auto`, `manual`. Default value `auto`.
* `classic_access` -(Optional, bool) Indicates whether this VPC should be connected to Classic Infrastructure. If true, This VPC's resources will have private network connectivity to the account's Classic Infrastructure resources. Only one VPC on an account may be connected in this way.
* `name` - (Required, string) The name of the VPC.
* `resource_group` - (Optional, Forces new resource, string) The resource group ID where the VPC to be created
* `tags` - (Optional, array of strings) Tags associated with the instance.
* `default_network_acl_name` - (Optional, string) The name of the default network acl.
* `default_security_group_name` - (Optional, string) The name of the default security group.
* `default_routing_table_name` - (Optional, string) The name of the default routing table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the VPC.
* `crn` - The CRN of VPC.
* `status` - The status of VPC.
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


## Import

ibm_is_vpc can be imported using ID, eg

```
$ terraform import ibm_is_vpc.example d7bec597-4726-451f-8a63-e62e6f19c32c
```