---
layout: "ibm"
page_title: "IBM : ibm_is_floating_ip_profile"
description: |-
  Get information about FloatingIPProfile
subcategory: "Virtual Private Cloud API"
---

# ibm_is_floating_ip_profile

Provides a read-only data source to retrieve information about a FloatingIPProfile. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_floating_ip_profile" "is_floating_ip_profile" {
	name = "user-ipv4"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `name` - (Required, Forces new resource, String) The floating IP profile name.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the FloatingIPProfile.
* `family` - (String) The product family this floating IP profile belongs to.- `provider`: The floating IP with this profile is owned by the provider.- `user`: The floating IP with this profile is owned by the user.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `provider`, `user`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `href` - (String) The URL for this floating IP profile.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)\/floating_ip\/profiles\/([^?#]+)$/`.
* `ip_version` - (String) The IP version for floating IPs with this profile:- `ipv4`: An IPv4 floating IP.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `ipv4`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `floating_ip_profile`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

