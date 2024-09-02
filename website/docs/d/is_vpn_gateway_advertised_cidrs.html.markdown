---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_gateway_advertised_cidrs"
description: |-
  Get information about VPNGatewayAdvertisedCIDRCollection
subcategory: "VPC infrastructure"
---

# ibm_is_vpn_gateway_advertised_cidrs

Provides a read-only data source to retrieve information about a VPNGatewayAdvertisedCIDRCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_vpn_gateway_advertised_cidrs" "is_vpn_gateway_advertised_cidrs" {
	vpn_gateway_id = "vpn_gateway_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

- `vpn_gateway` - (Required, Forces new resource, String) The VPN gateway identifier.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

- `id` - The unique identifier of the VPNGatewayAdvertisedCIDRCollection.
- `advertised_cidrs` - (List) The additional CIDRs advertised through any enabled routing protocol (for example, BGP). The routing protocol will advertise routes with these CIDRs and VPC prefixes as route destinations.
- `first` - (List) A link to the first page of resources.
	Nested schema for **first**:
	- `href` - (String) The URL for a page of resources.
- `limit` - (Integer) The maximum number of resources that can be returned by the request.
- `next` - (List) A link to the next page of resources. This property is present for all pagesexcept the last page.
	Nested schema for **next**:
	- `href` - (String) The URL for a page of resources.
- `total_count` - (Integer) The total number of resources across all pages.

