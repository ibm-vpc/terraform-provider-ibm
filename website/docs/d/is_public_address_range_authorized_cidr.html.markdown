---
layout: "ibm"
page_title: "IBM : ibm_is_public_address_range_authorized_cidr"
description: |-
  Get information about PublicAddressRangeAuthorizedCIDR
subcategory: "Virtual Private Cloud API"
---

# ibm_is_public_address_range_authorized_cidr

Provides a read-only data source to retrieve information about a PublicAddressRangeAuthorizedCIDR. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_public_address_range_authorized_cidr" "is_public_address_range_authorized_cidr" {
	is_public_address_range_authorized_cidr_id = "is_public_address_range_authorized_cidr_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `is_public_address_range_authorized_cidr_id` - (Required, Forces new resource, String) The public address range authorized CIDR identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the PublicAddressRangeAuthorizedCIDR.
* `allocation` - (List) 
Nested schema for **allocation**:
	* `count` - (Integer) The number of resources allocated from this public address range authorized CIDR.
	  * Constraints: The minimum value is `0`.
	* `profile_family` - (String) The profile `family` for resources allocated from this public address range authorized CIDR.- `provider`: The resources allocated from this authorized CIDR will have a profile  with a `family` value of `provider`.- `user`: The resources allocated from this authorized CIDR will have a profile with  a `family` value of `user`.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `provider`, `user`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `availability_mode` - (String) The availability mode of the public address range authorized CIDR:- `regional`: Resources allocated from the authorized CIDR can reside in any zone in the  region.- `zonal`: Resources allocated from the authorized CIDR must reside in the authorized  CIDR's `zone`.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `regional`, `zonal`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `cidr` - (String) The public IP address block for the public address range authorized CIDR, expressed in CIDR format.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 address blocks in the future.
* `href` - (String) The URL for this public address range authorized CIDR.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `ip_version` - (String) The IP version for this public address range authorized CIDR:- `ipv4`: An IPv4 public address range authorized CIDR.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `ipv4`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `name` - (String) The name for this public address range authorized CIDR. The name is unique across all public address range authorized CIDRs in the region.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `network_prefix_length` - (Integer) The network prefix length for this public address range authorized CIDR.
  * Constraints: The maximum value is `128`. The minimum value is `1`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `public_address_range_authorized_cidr`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `zone` - (List) The zone for this public address range authorized CIDR. Resources allocated from thisauthorized CIDR must reside in this zone.
Nested schema for **zone**:
	* `href` - (String) The URL for this zone.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `name` - (String) The globally unique name for this zone.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

