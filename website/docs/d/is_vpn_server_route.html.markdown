---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_server_route"
description: |-
  Get information about VPNServerRoute
subcategory: "Virtual Private Cloud API"
---

# ibm_is_vpn_server_route

Provides a read-only data source for VPNServerRoute. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_vpn_server_route" "is_vpn_server_route" {
	id = "id"
	vpn_server_id = "vpn_server_id"
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

* `id` - (Required, Forces new resource, String) The VPN route identifier.
* `vpn_server_id` - (Required, Forces new resource, String) The VPN server identifier.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the VPNServerRoute.
* `action` - (Required, String) The action to perform with a packet matching the VPN route:- `translate`: translate the source IP address to one of the private IP addresses of the VPN server.- `deliver`: deliver the packet into the VPC.- `drop`: drop the packetThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the VPN route on which the unexpected property value was encountered.
  * Constraints: Allowable values are: translate, deliver, drop

* `created_at` - (Required, String) The date and time that the VPN route was created.

* `destination` - (Required, String) The destination for this VPN route in the VPN server. If an incoming packet does not match any destination, it will be dropped.
  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`

* `href` - (Required, String) The URL for this VPN route.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

* `lifecycle_state` - (Required, String) The lifecycle state of the VPN route.
  * Constraints: Allowable values are: deleting, failed, pending, stable, updating, waiting, suspended

* `name` - (Required, String) The user-defined name for this VPN route.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`

* `resource_type` - (Required, String) The resource type.
  * Constraints: Allowable values are: vpn_server_route

