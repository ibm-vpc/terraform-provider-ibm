---
layout: "ibm"
page_title: "IBM : ibm_is_dynamic_route_server"
description: |-
  Manages DynamicRouteServer.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_dynamic_route_server

Create, update, and delete DynamicRouteServers with this resource.

## Example Usage

```hcl
resource "ibm_is_dynamic_route_server" "is_dynamic_route_server_instance" {
  asn = 64512
  ips {
		address = "192.168.3.4"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		name = "my-reserved-ip"
		resource_type = "subnet_reserved_ip"
  }
  name = "my-dynamic-route-server"
  resource_group {
		href = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		id = "fee82deba12e4c0fb69c3b09d1f12345"
		name = "my-resource-group"
  }
  security_groups {
		crn = "crn:v1:bluemix:public:is:us-south:a/123456::security-group:be5df5ca-12a0-494b-907e-aa6ec2bfa271"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/v1/security_groups/be5df5ca-12a0-494b-907e-aa6ec2bfa271"
		id = "be5df5ca-12a0-494b-907e-aa6ec2bfa271"
		name = "my-security-group"
  }
  vpc {
		crn = "crn:v1:bluemix:public:is:us-south:a/123456::vpc:4727d842-f94f-4a2d-824a-9bc9b02c523b"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/4727d842-f94f-4a2d-824a-9bc9b02c523b"
		id = "4727d842-f94f-4a2d-824a-9bc9b02c523b"
		name = "my-vpc"
		resource_type = "vpc"
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `asn` - (Required, Integer) The local autonomous system number (ASN) for this dynamic route server.
* `ips` - (Required, List) The reserved IPs bound to this dynamic route server.
Nested schema for **ips**:
	* `address` - (Required, String) The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this reserved IP.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this reserved IP.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this reserved IP. The name is unique across all reserved IPs in a subnet.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: `subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `name` - (Optional, String) The name for this dynamic route server. The name is unique across all dynamic route servers in the region.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `redistribute_service_routes` - (Optional, Boolean) Indicates whether all service routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `service`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the routeAdditionally, the CIDRs `161.26.0.0/16` (IBM services) and `166.8.0.0/14` (Cloud Service Endpoints) will also be redistributed to all peers through the routing protocol.
* `redistribute_subnets` - (Optional, Boolean) Indicates whether subnets meet the following conditions will be redistributed through the routing protocol to all peers as route destinations:- The subnet is attached to a routing table in the VPC this dynamic route server is  serving.- The routing table's `accept_routes_from` property includes the value  `dynamic_route_server`The routing protocol will redistribute routes with these subnets as route destinations.
* `redistribute_user_routes` - (Optional, Boolean) Indicates whether all user routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `user`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the route.
* `resource_group` - (Optional, List) The resource group for this dynamic route server.
Nested schema for **resource_group**:
	* `href` - (Computed, String) The URL for this resource group.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this resource group.
	  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
	* `name` - (Computed, String) The name for this resource group.
	  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
* `security_groups` - (Optional, List) The security groups targeting this dynamic route server.
  * Constraints: The minimum length is `1` item.
Nested schema for **security_groups**:
	* `crn` - (Required, String) The security group's CRN.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The security group's canonical URL.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this security group.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this security group. The name is unique across all security groups for the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `vpc` - (Required, List) The VPC this dynamic route server resides in.
Nested schema for **vpc**:
	* `crn` - (Required, String) The CRN for this VPC.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this VPC.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this VPC.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this VPC. The name is unique across all VPCs in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the DynamicRouteServer.
* `created_at` - (String) The date and time that the dynamic route server was created.
* `crn` - (String) The CRN for this dynamic route server.
  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
* `health_state` - (String) The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
  * Constraints: Allowable values are: `degraded`, `faulted`, `inapplicable`, `ok`.
* `href` - (String) The URL for this dynamic route server.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_state` - (String) The lifecycle state of the dynamic route server.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `dynamic_route_server`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `etag` - ETag identifier for DynamicRouteServer.

## Import

You can import the `ibm_is_dynamic_route_server` resource by using `id`. The unique identifier for this dynamic route server.

# Syntax
```
$ terraform import ibm_is_dynamic_route_server.is_dynamic_route_server <id>
```

# Example
```
$ terraform import ibm_is_dynamic_route_server.is_dynamic_route_server r006-d7cc5196-9864-48c4-82d8-3f30da41fcc5
```
