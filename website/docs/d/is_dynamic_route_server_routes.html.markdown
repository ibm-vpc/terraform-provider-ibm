---
layout: "ibm"
page_title: "IBM : ibm_is_dynamic_route_server_routes"
description: |-
  Get information about DynamicRouteServerRouteCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_dynamic_route_server_routes

Provides a read-only data source to retrieve information about a DynamicRouteServerRouteCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_dynamic_route_server_routes" "is_dynamic_route_server_routes" {
	dynamic_route_server_id = "dynamic_route_server_id"
	sort = "created_at"
	type = "learned"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `dynamic_route_server_id` - (Required, Forces new resource, String) The dynamic route server identifier.
* `peer_id` - (Optional, String) Filters the collection to dynamic route server routes with `peer.id` matching the specified value.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `sort` - (Optional, String) Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-destination` sorts the collection by the `destination` property in descending order.
  * Constraints: The default value is `-created_at`. Allowable values are: `created_at`, `destination`.
* `type` - (Optional, String) Filters the collection to dynamic route server routes with `type` matching the specified value.
  * Constraints: Allowable values are: `learned`, `redistributed_service_routes`, `redistributed_subnets`, `redistributed_user_routes`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the DynamicRouteServerRouteCollection.
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

* `routes` - (List) Collection of dynamic route server routes.
Nested schema for **routes**:
	* `as_path` - (List) The ordered sequence of autonomous systems that network packets will traverse to get to this dynamic route server, per the rules defined in [RFC 4271](https://www.rfc-editor.org/rfc/rfc4271).
	  * Constraints: The maximum length is `64` items. The minimum length is `0` items.
	* `created_at` - (String) The date and time that the route was created.
	* `destination` - (String) The destination of the route. Each learned route must have a unique combination of`destination`, `source_ip`, and `next_hop`. Similarly, each redistributed route must have a unique combination of `destination` and `next_hop`. The learned route is not added to the VPC routing table if its `destination` is the same as the redistributed route.
	  * Constraints: The maximum length is `18` characters. The minimum length is `9` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`.
	* `href` - (String) The URL for this dynamic route server route.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this dynamic route server route.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `next_hop` - (List) The next hop packets will be routed to.
	Nested schema for **next_hop**:
		* `address` - (String) The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `peer` - (List)
	Nested schema for **peer**:
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this dynamic route server peer.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this dynamic route server peer.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this dynamic route server peer. The name is unique across all peers for the dynamic route server.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `dynamic_route_server_peer`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `dynamic_route_server_route`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `source_ip` - (List) The source IP of the dynamic route server used to establish routing protocol withthis dynamic route server peer.This property will be present only when the route `type` is `learned`.
	Nested schema for **source_ip**:
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
	* `type` - (String) The type of this route:- `learned`: route was learned from the dynamic route server peer via the routing  protocol. The learned route was evaluated based on [the best route selection  algorithm](https://cloud.ibm.com/docs/vpc?topic=drs-best-route-selection) to determine if  it was added to the VPC routing table- `redistributed_subnets`: route was redistributed to the dynamic route  server peer via the routing protocol, and route's destination is a subnet IP CIDR  block- `redistributed_user_routes`: route was redistributed to the dynamic route server  peer via the routing protocol, and it is from a VPC routing table with the `origin`  set as `user`- `redistributed_service_routes`: route was redistributed to the dynamic route server  peer via the routing protocol, and it is from a VPC routing table with the `origin`  set as `service`The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered.
	  * Constraints: Allowable values are: `learned`, `redistributed_service_routes`, `redistributed_subnets`, `redistributed_user_routes`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `total_count` - (Integer) The total number of resources across all pages.
  * Constraints: The minimum value is `0`.

