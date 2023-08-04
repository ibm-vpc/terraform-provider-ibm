---
layout: "ibm"
page_title: "IBM : ibm_is_dynamic_route_servers"
description: |-
  Get information about DynamicRouteServerCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_dynamic_route_servers

Provides a read-only data source to retrieve information about a DynamicRouteServerCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_dynamic_route_servers" "is_dynamic_route_servers" {
	sort = "name"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `name` - (Optional, String) Filters the collection to resources with a `name` property matching the exact specified name.
* `resource_group_id` - (Optional, String) Filters the collection to resources with a `resource_group.id` property matching the specified identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `sort` - (Optional, String) Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-created_at` sorts the collection by the `created_at` property in descending order, and the value `name` sorts it by the `name` property in ascending order.
  * Constraints: The default value is `-created_at`. Allowable values are: `created_at`, `name`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the DynamicRouteServerCollection.
* `dynamic_route_servers` - (List) Collection of dynamic route servers.
Nested schema for **dynamic_route_servers**:
	* `asn` - (Integer) The local autonomous system number (ASN) for this dynamic route server.
	* `created_at` - (String) The date and time that the dynamic route server was created.
	* `crn` - (String) The CRN for this dynamic route server.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `health_state` - (String) The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
	  * Constraints: Allowable values are: `degraded`, `faulted`, `inapplicable`, `ok`.
	* `href` - (String) The URL for this dynamic route server.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this dynamic route server.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `ips` - (List) The reserved IPs bound to this dynamic route server.
	Nested schema for **ips**:
		* `address` - (String) The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this reserved IP.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this reserved IP.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this reserved IP. The name is unique across all reserved IPs in a subnet.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `lifecycle_state` - (String) The lifecycle state of the dynamic route server.
	  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
	* `name` - (String) The name for this dynamic route server. The name is unique across all dynamic route servers in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `redistribute_service_routes` - (Boolean) Indicates whether all service routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `service`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the routeAdditionally, the CIDRs `161.26.0.0/16` (IBM services) and `166.8.0.0/14` (Cloud Service Endpoints) will also be redistributed to all peers through the routing protocol.
	* `redistribute_subnets` - (Boolean) Indicates whether subnets meet the following conditions will be redistributed through the routing protocol to all peers as route destinations:- The subnet is attached to a routing table in the VPC this dynamic route server is  serving.- The routing table's `accept_routes_from` property includes the value  `dynamic_route_server`The routing protocol will redistribute routes with these subnets as route destinations.
	* `redistribute_user_routes` - (Boolean) Indicates whether all user routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `user`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the route.
	* `resource_group` - (List) The resource group for this dynamic route server.
	Nested schema for **resource_group**:
		* `href` - (String) The URL for this resource group.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this resource group.
		  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
		* `name` - (String) The name for this resource group.
		  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `dynamic_route_server`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `security_groups` - (List) The security groups targeting this dynamic route server.
	  * Constraints: The minimum length is `1` item.
	Nested schema for **security_groups**:
		* `crn` - (String) The security group's CRN.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The security group's canonical URL.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this security group.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this security group. The name is unique across all security groups for the VPC.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `vpc` - (List) The VPC this dynamic route server resides in.
	Nested schema for **vpc**:
		* `crn` - (String) The CRN for this VPC.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this VPC.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this VPC.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this VPC. The name is unique across all VPCs in the region.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `first` - (List) A link to the first page of resources.
Nested schema for **first**:
	* `href` - (String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `limit` - (Integer) The maximum number of resources that can be returned by the request.
  * Constraints: The maximum value is `100`. The minimum value is `1`.

* `next` - (List) A link to the next page of resources. This property is present for all pagesexcept the last page.
Nested schema for **next**:
	* `href` - (String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `total_count` - (Integer) The total number of resources across all pages.
  * Constraints: The minimum value is `0`.

