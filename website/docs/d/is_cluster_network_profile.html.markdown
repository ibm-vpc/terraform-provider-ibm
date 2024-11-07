---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network_profile"
description: |-
  Get information about ClusterNetworkProfile
---

# ibm_is_cluster_network_profile

Provides a read-only data source to retrieve information about a ClusterNetworkProfile. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_cluster_network_profile" "example" {
	name = "name"
}
```

## Argument Reference

You can specify the following arguments for this data source.

- `name` - (Required, Forces new resource, String) The cluster network profile name.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

- `id` - The unique identifier of the ClusterNetworkProfile.
- `family` - (String) The product family this cluster network profile belongs to.
- `href` - (String) The URL for this cluster network profile.
- `resource_type` - (String) The resource type.
- `supported_instance_profiles` - (List) The instance profiles that support this cluster network profile.
	Nested schema for **supported_instance_profiles**:
	- `href` - (String) The URL for this virtual server instance profile.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `name` - (String) The globally unique name for this virtual server instance profile.
	- `resource_type` - (String) The resource type.
- `zones` - (List) Zones in this region that support this cluster network profile.
	Nested schema for **zones**:
	- `href` - (String) The URL for this zone.
	- `name` - (String) The globally unique name for this zone.

