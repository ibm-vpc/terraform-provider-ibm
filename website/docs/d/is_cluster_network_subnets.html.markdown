---
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network_subnets"
description: |-
  Get information about ClusterNetworkSubnetCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_cluster_network_subnets

Provides a read-only data source to retrieve information about a ClusterNetworkSubnetCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_cluster_network_subnets" "is_cluster_network_subnets" {
	cluster_network_id = "cluster_network_id"
	sort = "name"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `cluster_network_id` - (Required, Forces new resource, String) The cluster network identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `name` - (Optional, String) Filters the collection to resources with a `name` property matching the exact specified name.
* `sort` - (Optional, String) Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-created_at` sorts the collection by the `created_at` property in descending order, and the value `name` sorts it by the `name` property in ascending order.
  * Constraints: The default value is `-created_at`. Allowable values are: `created_at`, `name`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the ClusterNetworkSubnetCollection.
* `subnets` - (List) Collection of cluster network subnets.
Nested schema for **subnets**:
	* `available_ipv4_address_count` - (Integer) The number of IPv4 addresses in this cluster network subnet that are not in use, and have not been reserved by the user or the provider.
	* `created_at` - (String) The date and time that the cluster network subnet was created.
	* `href` - (String) The URL for this cluster network subnet.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^https:\/\/([^\/?#]*)([^?#]*)\/cluster_networks\/[-0-9a-z_]+\/subnets\/[-0-9a-z_]+$/`.
	* `id` - (String) The unique identifier for this cluster network subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `ip_version` - (String)
	  * Constraints: Allowable values are: `ipv4`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `ipv4_cidr_block` - (String) The IPv4 range of this cluster network subnet, expressed in CIDR format.
	  * Constraints: The maximum length is `18` characters. The minimum length is `9` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`.
	* `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).
	  * Constraints: The minimum length is `0` items.
	Nested schema for **lifecycle_reasons**:
		* `code` - (String) A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
		  * Constraints: Allowable values are: `internal_error`, `resource_suspended_by_provider`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `message` - (String) An explanation of the reason for this lifecycle state.
		* `more_info` - (String) Link to documentation about the reason for this lifecycle state.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `lifecycle_state` - (String) The lifecycle state of the cluster network subnet.
	  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
	* `name` - (String) The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `total_ipv4_address_count` - (Integer) The total number of IPv4 addresses in this cluster network subnet.Note: This is calculated as 2<sup>(32 - prefix length)</sup>. For example, the prefix length `/24` gives:<br> 2<sup>(32 - 24)</sup> = 2<sup>8</sup> = 256 addresses.

