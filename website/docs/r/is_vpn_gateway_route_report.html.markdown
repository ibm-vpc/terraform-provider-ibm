---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_gateway_route_report"
description: |-
  Manages VPNRouteReport.
subcategory: "VPC infrastructure"
---

# ibm_is_vpn_gateway_route_report

Create, update, and delete VPNRouteReports with this resource.

## Example Usage

```terraform
resource "ibm_is_vpn_gateway_route_report" "is_vpn_gateway_route_report_instance" {
  vpn_gateway = var.vpn_gateway_id
}
```

## Argument Reference

You can specify the following arguments for this resource.

- `vpn_gateway` - (Required, Forces new resource, String) The VPN gateway identifier.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

- `id` - The unique identifier of the VPNRouteReport.
- `created_at` - (String) The date and time that this route report was created.
- `route_report_id` - (String) The unique identifier for this route report.
- `routes` - (List) The routes of this report.
	Nested schema for **routes**:
	- `as_path` - (List) AS path numbers of this route.
	- `best_path` - (Boolean) Indicates whether this route is best path.
	- `next_hops` - (List) Next hop list of this route.
		Nested schema for **next_hops**:
		- `address` - (String) A unicast IP address, which must not be any of the following values:- `0.0.0.0` (the sentinel IP address)- `224.0.0.0` to `239.255.255.255` (multicast IP addresses)- `255.255.255.255` (the broadcast IP address)This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.
		- `used` - (Boolean) The next hop is used for traffic forward.
	- `peer` - (String) The peer of this route.
	- `prefix` - (String) The destination of this route.
	- `valid` - (Boolean) Indicates whether this route is valid.
	- `weight` - (Integer) Weight of this route.
- `status` - (String) Route report status. The list of enumerated values for this property may expand in the future. Code and processes using this field must tolerate unexpected values.- `pending`: generic routing encapsulation tunnel attached- `complete`: generic routing encapsulation tunnel detached. Allowable values are: `complete`, `pending`.
- `updated_at` - (String) The date and time that this route report was updated.


## Import

You can import the `ibm_is_vpn_gateway_route_report` resource by using `id`.
The `id` property can be formed from `vpn_gateway`, and `route_report_id` in the following format:

<pre>
&lt;vpn_gateway_id&gt;/&lt;is_vpn_gateway_route_report_id&gt;
</pre>
- `vpn_gateway`: A string. The VPN gateway identifier.
- `route_report_id`: A string in the format `ddf51bec-3424-11e8-b467-0ed5f89f718b`. The unique identifier for this route report.

# Syntax
<pre>
$ terraform import ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report &lt;vpn_gateway&gt;/&lt;route_report_id&gt;
</pre>
