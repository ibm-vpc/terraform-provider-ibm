---
layout: "ibm"
page_title: "IBM : ibm_is_public_address_range_authorized_cidr_allocations"
description: |-
  Get information about PublicAddressRangeAuthorizedCIDRAllocationCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_public_address_range_authorized_cidr_allocations

Provides a read-only data source to retrieve information about a PublicAddressRangeAuthorizedCIDRAllocationCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_public_address_range_authorized_cidr_allocations" "is_public_address_range_authorized_cidr_allocations" {
	allocations_resource_type = "floating_ip"
	authorized_cidr_id = "authorized_cidr_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `allocations_resource_type` - (Optional, String) Filters the collection to resources with an item in the `allocations` property with a`resource_type` property matching the specified value.
  * Constraints: The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `authorized_cidr_id` - (Required, Forces new resource, String) The public address range authorized CIDR identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the PublicAddressRangeAuthorizedCIDRAllocationCollection.
* `allocations` - (List) The floating IPs and public address ranges allocated from this public address range authorized CIDR.
  * Constraints: The minimum length is `0` items.
Nested schema for **allocations**:
	* `address` - (String) The globally unique IP address.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `cidr` - (String) The public IP address block for this public address range, expressed in CIDR format.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 address blocks in the future.
	* `crn` - (String) The CRN for this floating IP.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) A link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this floating IP.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this floating IP.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this floating IP. The name is unique across all floating IPs in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `floating_ip`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

