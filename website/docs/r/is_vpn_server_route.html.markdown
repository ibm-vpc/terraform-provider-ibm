---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_vpn_server_route"
description: |-
  Manages IBM VPN Server Route.
---

# ibm_is_vpn_server_route

Provides a resource for VPNServerRoute. This allows VPNServerRoute to be created, updated and deleted.

## Example Usage

```terraform
resource "ibm_is_vpn_server_route" "is_vpn_server_route" {
  vpn_server_id = ibm_is_vpn_server.is_vpn_server.vpn_server
  destination   = "172.16.0.0/16"
  action        = "translate"
  name          = "my_vpn_server_route"
}
```

## Argument Reference

Review the argument reference that you can specify for your resource.

* `action` - (Optional, String) The action to perform with a packet matching the VPN route:- `translate`: translate the source IP address to one of the private IP addresses of the VPN server, then deliver the packet to target.- `deliver`: deliver the packet to the target.- `drop`: drop the packetThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the VPN route on which the unexpected property value was encountered.
  * Constraints: The default value is `deliver`. Allowable values are: translate, deliver, drop
* `destination` - (Required, String) The destination to use for this VPN route in the VPN server. Must be unique within the VPN server. If an incoming packet does not match any destination, it will be dropped.
  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`
* `name` - (Optional, String) The user-defined name for this VPN route. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the VPN server the VPN route resides in.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$/`
* `vpn_server_id` - (Required, Forces new resource, String) The VPN server identifier.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the VPNServerRoute.
* `vpn_route` - The identifier of the VPNServerRoute.
* `created_at` - (Required, String) The date and time that the VPN route was created.
* `href` - (Required, String) The URL for this VPN route.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
* `lifecycle_state` - (Required, String) The lifecycle state of the VPN route.
  * Constraints: Allowable values are: deleting, failed, pending, stable, updating, waiting, suspended
* `resource_type` - (Required, String) The resource type.
  * Constraints: Allowable values are: vpn_server_route

## Import

You can import the `ibm_is_vpn_server_route` resource by using `id`.
The `id` property can be formed from `vpn_server_id`, and `vpn_route` in the following format:

```
<vpn_server_id>/<id>
```
* `vpn_server_id`: A string. The VPN server identifier.
* `vpn_route`: A string. The VPN route identifier.

# Syntax
```
$ terraform import ibm_is_vpn_server_route.is_vpn_server_route <vpn_server_id>/<id>
```
