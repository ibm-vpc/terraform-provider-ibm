---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_server_clients"
description: |-
  Get information about VPNServerClientCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_vpn_server_clients

Provides a read-only data source for VPNServerClientCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_vpn_server_clients" "is_vpn_server_clients" {
	vpn_server_id = "vpn_server_id"
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

* `vpn_server_id` - (Required, Forces new resource, String) The VPN server identifier.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the VPNServerClientCollection.
* `clients` - (Required, List) Collection of VPN clients.
Nested scheme for **clients**:
	* `client_ip` - (Required, List) The IP address assigned to this VPN client from `client_ip_pool`.
	Nested scheme for **client_ip**:
		* `address` - (Required, String) The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
		  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`
	* `common_name` - (Optional, String) The common name of client certificate that the VPN client provided when connecting to the server.This property will be present only when the `certificate` client authentication method is enabled on the VPN server.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character.
	* `created_at` - (Required, String) The date and time that the VPN client was created.
	* `disconnected_at` - (Optional, String) The date and time that the VPN client was disconnected.
	* `href` - (Required, String) The URL for this VPN client.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `id` - (Required, String) The unique identifier for this VPN client.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`
	* `remote_ip` - (Required, List) The remote IP address of this VPN client.
	Nested scheme for **remote_ip**:
		* `address` - (Required, String) The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
		  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`
	* `remote_port` - (Required, Integer) The remote port of this VPN client.
	  * Constraints: The maximum value is `65535`. The minimum value is `1`.
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: vpn_server_client
	* `status` - (Required, String) The status of the VPN client:- `connected`: the VPN client is `connected` to this VPN server.- `disconnected`: the VPN client is `disconnected` from this VPN server.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the VPN client on which the unexpected property value was encountered.
	  * Constraints: Allowable values are: connected, disconnected
	* `username` - (Optional, String) The username that this VPN client provided when connecting to the VPN server.This property will be present only when  the`username` client authentication method is enabled on the VPN server.

* `first` - (Required, List) A link to the first page of resources.
Nested scheme for **first**:
	* `href` - (Required, String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

* `limit` - (Required, Integer) The maximum number of resources that can be returned by the request.
  * Constraints: The maximum value is `100`. The minimum value is `1`.

* `next` - (Optional, List) A link to the next page of resources. This property is present for all pagesexcept the last page.
Nested scheme for **next**:
	* `href` - (Required, String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

* `total_count` - (Required, Integer) The total number of resources across all pages.
  * Constraints: The minimum value is `0`.

