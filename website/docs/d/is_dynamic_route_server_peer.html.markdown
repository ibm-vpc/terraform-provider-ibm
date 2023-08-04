---
layout: "ibm"
page_title: "IBM : ibm_is_dynamic_route_server_peer"
description: |-
  Get information about DynamicRouteServerPeer
subcategory: "Virtual Private Cloud API"
---

# ibm_is_dynamic_route_server_peer

Provides a read-only data source to retrieve information about a DynamicRouteServerPeer. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer" {
	dynamic_route_server_id = ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer.dynamic_route_server_id
	id = "id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `dynamic_route_server_id` - (Required, Forces new resource, String) The dynamic route server identifier.
* `id` - (Required, Forces new resource, String) The dynamic route server peer identifier.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the DynamicRouteServerPeer.
* `asn` - (Integer) The autonomous system number (ASN) for this dynamic route server peer.

* `authentication_enabled` - (Boolean) Indicates whether TCP MD5 authentication key is configured and enabled in this dynamic route server peer.

* `bfd` - (List) The bidirectional forwarding detection (BFD) configuration for this dynamic route serverpeer.
Nested schema for **bfd**:
	* `mode` - (String) The bidirectional forwarding detection operating mode on this peer.
	  * Constraints: Allowable values are: `asynchronous`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `role` - (String) The bidirectional forwarding detection role in session initialization.
	  * Constraints: Allowable values are: `active`, `disabled`, `passive`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `sessions` - (List) The sessions for this bidirectional forwarding detection for this peer.
	  * Constraints: The maximum length is `3` items. The minimum length is `1` item.
	Nested schema for **sessions**:
		* `source_ip` - (List) The source IP of the dynamic route server used to establish bidirectional forwardingdetection session with this dynamic route server peer.
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
		* `state` - (String) The current bidirectional forwarding detection session state as seen by this dynamic route server.
		  * Constraints: Allowable values are: `admin_down`, `down`, `init`, `up`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `created_at` - (String) The date and time that the dynamic route server peer was created.

* `dynamic_route_server_peer_id` - (String) The unique identifier for this dynamic route server peer.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

* `href` - (String) The URL for this dynamic route server peer.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `ip` - (List) The IP address of this dynamic route server peer.The peer IP must be in a subnet in the VPC this dynamic route server is serving.
Nested schema for **ip**:
	* `address` - (String) The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.

* `lifecycle_state` - (String) The lifecycle state of the dynamic route server peer.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.

* `name` - (String) The name for this dynamic route server peer. The name is unique across all peers for the dynamic route server.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

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

