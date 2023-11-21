---
layout: "ibm"
page_title: "IBM : ibm_is_dynamic_route_server_peer"
description: |-
  Manages DynamicRouteServerPeer.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_dynamic_route_server_peer

Create, update, and delete DynamicRouteServerPeers with this resource.

## Example Usage

```hcl
resource "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
  asn = 64513
  bfd {
		mode = "asynchronous"
		role = "active"
		sessions {
			source_ip {
				address = "192.168.3.4"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				name = "my-reserved-ip"
				resource_type = "subnet_reserved_ip"
			}
			state = "admin_down"
		}
  }
  dynamic_route_server_id = "dynamic_route_server_id"
  ip {
		address = "192.168.3.4"
  }
  name = "my-dynamic-route-server-peer"
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `asn` - (Required, Integer) The autonomous system number (ASN) for this dynamic route server peer.
* `bfd` - (Optional, List) The bidirectional forwarding detection (BFD) configuration for this dynamic route serverpeer.
Nested schema for **bfd**:
	* `mode` - (Computed, String) The bidirectional forwarding detection operating mode on this peer.
	  * Constraints: Allowable values are: `asynchronous`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `role` - (Required, String) The bidirectional forwarding detection role in session initialization.
	  * Constraints: Allowable values are: `active`, `disabled`, `passive`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `sessions` - (Required, List) The sessions for this bidirectional forwarding detection for this peer.
	  * Constraints: The maximum length is `3` items. The minimum length is `1` item.
	Nested schema for **sessions**:
		* `source_ip` - (Required, List) The source IP of the dynamic route server used to establish bidirectional forwardingdetection session with this dynamic route server peer.
		Nested schema for **source_ip**:
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
		* `state` - (Required, String) The current bidirectional forwarding detection session state as seen by this dynamic route server.
		  * Constraints: Allowable values are: `admin_down`, `down`, `init`, `up`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `dynamic_route_server_id` - (Required, Forces new resource, String) The dynamic route server identifier.
* `ip` - (Required, List) The IP address of this dynamic route server peer.The peer IP must be in a subnet in the VPC this dynamic route server is serving.
Nested schema for **ip**:
	* `address` - (Required, String) The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
* `name` - (Optional, String) The name for this dynamic route server peer. The name is unique across all peers for the dynamic route server.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the DynamicRouteServerPeer.
* `authentication_enabled` - (Boolean) Indicates whether TCP MD5 authentication key is configured and enabled in this dynamic route server peer.
* `created_at` - (String) The date and time that the dynamic route server peer was created.
* `dynamic_route_server_peer_id` - (String) The unique identifier for this dynamic route server peer.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `href` - (String) The URL for this dynamic route server peer.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_state` - (String) The lifecycle state of the dynamic route server peer.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `dynamic_route_server_peer`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `sessions` - (List) The sessions for this dynamic route server peer.
  * Constraints: The maximum length is `3` items. The minimum length is `1` item.
Nested schema for **sessions**:
	* `established_at` - (String) The date and time that the BGP session was established.This property will be present only when the session `state` is `established`.
	* `source_ip` - (List) The source IP of the dynamic route server used to establish routing protocol with thisdynamic route server peer.
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
	* `state` - (String) The state of the routing protocol with this dynamic route server peer.
	  * Constraints: Allowable values are: `active`, `connect`, `established`, `idle`, `initializing`, `open_confirm`, `open_sent`.

* `etag` - ETag identifier for DynamicRouteServerPeer.

## Import

You can import the `ibm_is_dynamic_route_server_peer` resource by using `id`.
The `id` property can be formed from `dynamic_route_server_id`, and `id` in the following format:

```
<dynamic_route_server_id>/<id>
```
* `dynamic_route_server_id`: A string. The dynamic route server identifier.
* `id`: A string. The dynamic route server peer identifier.

# Syntax
```
$ terraform import ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer <dynamic_route_server_id>/<id>
```
