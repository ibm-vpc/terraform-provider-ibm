---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : vpc-routing-tables"
description: |-
  Manages IBM IS VPC Routing tables.
---

# ibm\_is_vpc_routing_table

Provides a vpc routing tables resource. This allows vpc routing tables to be created, updated, and cancelled.


## Example Usage

```terraform
resource "ibm_is_vpc" "testacc_vpc" {
  name = "testvpc"
}
resource "ibm_is_vpc_routing_table" "test_ibm_is_vpc_routing_table" {
	vpc = ibm_is_vpc.testacc_vpc.id
	name = "routTabletest"
	route_direct_link_ingress = true
	route_transit_gateway_ingress = false
        route_vpc_zone_ingress = false
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, string) The routing table name.
* `vpc` - (Required, Forces new resource, string) The vpc id. 
* `route_direct_link_ingress` - (Optional,boolean) if set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC. For this to succeed, the VPC must not already have a routing table with this property set to true.
* `route_transit_gateway_ingress` - (Optional,boolean) If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC. For this to succeed, the VPC must not already have a routing table with this property set to true.
* `route_vpc_zone_ingress` - (Optional,boolean) f set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC. For this to succeed, the VPC must not already have a routing table with this property set to true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this routing table. The id is composed of \<vpc_id\>/\<vpc_route_table_id\>
* `href` - The URL for this routing table.
* `is_default` - Indicates whether this is the default routing table for this VPC
* `lifecycle_state` - The lifecycle state of the routing table
* `resource_type` - The resource type
* `routing_table` - The generated routing table identifier
* `subnets` - The subnets to which this routing table is attached
  * `id` - The unique identifier for this subnet
  * `name` - The user-defined name for this subnet
* `routes` - 	The routes for this routing table.
  * `id` - The unique identifier for this route
  * `name` - The user-defined name for this route


## Import

ibm_is_vpc_routing_table can be imported using VPC ID and VPC Route table ID, eg

```
$ terraform import ibm_is_vpc_routing_table.example 56738c92-4631-4eb5-8938-8af9211a6ea4/fc2667e0-9e6f-4993-a0fd-cabab477c4d1
```
