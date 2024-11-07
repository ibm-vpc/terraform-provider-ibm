---
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network_profile"
description: |-
  Get information about ClusterNetworkProfile
subcategory: "VPC infrastructure"
---

# ibm_is_cluster_network_profile

Provides a read-only data source to retrieve information about a ClusterNetworkProfile. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_cluster_network_profile" "is_cluster_network_profile" {
	name = "h100"
}
```

## Argument Reference

You can specify the following arguments for this data source.

- `name` - (Required, Forces new resource, String) The cluster network profile name.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

- `id` - The unique identifier of the ClusterNetworkProfile.
- `family` - (String) The product family this cluster network profile belongs to.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `vela`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `href` - (String) The URL for this cluster network profile.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
- `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `cluster_network_profile`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `supported_instance_profiles` - (List) The instance profiles that support this cluster network profile.
  * Constraints: The maximum length is `100` items. The minimum length is `1` item.
Nested schema for **supported_instance_profiles**:
	- `href` - (String) The URL for this virtual server instance profile.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `name` - (String) The globally unique name for this virtual server instance profile.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	- `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `instance_profile`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `zones` - (List) Zones in this region that support this cluster network profile.
  * Constraints: The maximum length is `10` items. The minimum length is `1` item.
Nested schema for **zones**:
	- `href` - (String) The URL for this zone.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `name` - (String) The globally unique name for this zone.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

