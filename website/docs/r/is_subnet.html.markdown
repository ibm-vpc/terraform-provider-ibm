---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : subnet"
description: |-
  Manages IBM subnet.
---

# ibm\_is_subnet

Provides a subnet resource. This allows subnet to be created, updated, and cancelled.


## Example Usage

```terraform
resource "ibm_is_vpc" "testacc_vpc" {
  name = "test"
}

resource "ibm_is_vpc_routing_table" "test_cr_route_table1" {
  name   = "test-cr-route-table1"
  vpc    = data.ibm_is_vpc.testacc_vpc.id
}


resource "ibm_is_subnet" "testacc_subnet" {
  name            = "test_subnet"
  vpc             = ibm_is_vpc.testacc_vpc.id
  zone            = "us-south-1"
  ipv4_cidr_block = "192.168.0.0/1"
  routing_table   = ibm_is_vpc_routing_table.test_cr_route_table1.routing_table  

  //User can configure timeouts
  timeouts {
    create = "90m"
    delete = "30m"
  }
}
```

## Timeouts

ibm_is_subnet provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 10 minutes) Used for creating Instance.
* `update` - (Default 10 minutes) Used for creating Instance.
* `delete` - (Default 10 minutes) Used for deleting Instance.

## Argument Reference

The following arguments are supported:


* `ipv4_cidr_block` - (Optional, Forces new resource, string)   The IPv4 range of the subnet.
    **NOTE**: Conflicts with `total_ipv4_address_count`
* `total_ipv4_address_count` - (Optional, Forces new resource, string) The total number of IPv4 addresses.
    **NOTE**: Conflicts with `ipv4_cidr_block` and one of `ipv4_cidr_block`, `total_ipv4_address_count` is mandatory.
* `ip_version` - (Optional, Forces new resource, string) The Ip Version. The default is `ipv4`.
* `name` - (Required, string) The name of the subnet.
* `network_acl` - (Optional, string) The ID of the network ACL for the subnet.
* `public_gateway` - (Optional, string) The ID of the public-gateway for the subnet.
* `vpc` - (Required, Forces new resource, string) The vpc id.
* `zone` - (Required, Forces new resource, string) The subnet zone name.
* `routing_table` - (Optional, string) The routing table identifier that is associated with the subnet. 
* `resource_group` - (Optional, Forces new resource, string) The resource group ID where the Subnet to be created (This argument is supported only for Generation `2` infrastructure)
* `tags` - (Optional, list(string)) Tags associated with the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the subnet.
* `ipv6_cidr_block` - The IPv6 range of the subnet.
* `status` - The status of the subnet.
* `available_ipv4_address_count` - The total number of available IPv4 addresses.
* `crn` - The CRN of subnet.

## Import

ibm_is_subnet can be imported using ID, eg

```
$ terraform import ibm_is_subnet.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
