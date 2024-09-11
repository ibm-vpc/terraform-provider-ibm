---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_gateway_route_reports"
description: |-
  Get information about VPNRouteReportCollection
subcategory: "VPC infrastructure"
---

# ibm_is_vpn_gateway_route_reports

Provides a read-only data source to retrieve information about a VPNRouteReportCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_vpn_gateway_route_reports" "is_vpn_gateway_route_reports" {
	vpn_gateway_id = ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance.vpn_gateway_id
}
```

## Argument Reference

You can specify the following arguments for this data source.

- `vpn_gateway` - (Required, Forces new resource, String) The VPN gateway identifier.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

- `id` - The unique identifier of the VPNRouteReportCollection.
- `route_reports` - (List) Collection of VPN gateway route reports.
	Nested schema for **route_reports**:
	- `created_at` - (String) The date and time that this route report was created.
	- `id` - (String) The unique identifier for this route report.
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
	- `status` - (String) Route report status. The list of enumerated values for this property may expand in the future. Code and processes using this field must tolerate unexpected values.- `pending`: generic routing encapsulation tunnel attached- `complete`: generic routing encapsulation tunnel detached. Allowable values are: `complete`, `pending`
	- `updated_at` - (String) The date and time that this route report was updated.

