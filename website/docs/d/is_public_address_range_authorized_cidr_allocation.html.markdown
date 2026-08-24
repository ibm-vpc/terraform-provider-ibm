---
layout: "ibm"
page_title: "IBM : ibm_is_public_address_range_authorized_cidr_allocation"
description: |-
  Get information about a PublicAddressRangeAuthorizedCIDRAllocation
subcategory: "Virtual Private Cloud API"
---

# ibm_is_public_address_range_authorized_cidr_allocation

Provides a read-only data source to retrieve information about a single allocation from a public address range authorized CIDR. An allocation is either a floating IP or a public address range that has been allocated from the authorized CIDR. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_public_address_range_authorized_cidr_allocation" "example" {
  authorized_cidr_id            = "r134-7be42030-e392-43b0-9ae8-2a8f2798c6f1"
  authorized_cidr_allocation_id = "r134-3a8d6a3c-1d2e-4f5b-9c0a-7e8f9a0b1c2d"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `authorized_cidr_id` - (Required, Forces new resource, String) The public address range authorized CIDR identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `authorized_cidr_allocation_id` - (Required, Forces new resource, String) The public address range authorized CIDR allocation identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the PublicAddressRangeAuthorizedCIDRAllocation.
* `address` - (String) The globally unique IP address. Present when the allocation `resource_type` is `floating_ip`.
  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
* `cidr` - (String) The public IP address block for this public address range, expressed in CIDR format. Present when the allocation `resource_type` is `public_address_range`. This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 address blocks in the future.
* `crn` - (String) The CRN for this allocation.
  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\.\/]*$/`.
* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.
Nested schema for **deleted**:
	* `more_info` - (String) A link to documentation about deleted resources.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `href` - (String) The URL for this allocation.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `name` - (String) The name for this allocation. The name is unique across all allocations in the region.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `floating_ip`, `public_address_range`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
