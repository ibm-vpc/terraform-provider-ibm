---
layout: "ibm"
page_title: "IBM : ibm_is_public_address_range_profiles"
description: |-
  Get information about PublicAddressRangeProfileCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_public_address_range_profiles

Provides a read-only data source to retrieve information about a PublicAddressRangeProfileCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_public_address_range_profiles" "is_public_address_range_profiles" {
}
```


## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the PublicAddressRangeProfileCollection.
* `profiles` - (List) A page of public address range profiles.
Nested schema for **profiles**:
	* `family` - (String) The product family this public address range profile belongs to.- `provider`: The public IP addresses in the public address range with this profile are   owned by the provider.- `user`: The public IP addresses in the public address range with this profile are   owned by the user.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `provider`, `user`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `href` - (String) The URL for this public address range profile.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)\/public_address_range\/profiles\/([^?#]+)$/`.
	* `ip_version` - (String) The IP version for public address ranges with this profile:- `ipv4`: An IPv4 public address range.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `ipv4`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `name` - (String) The globally unique name for this public address range profile.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `public_address_range_profile`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

