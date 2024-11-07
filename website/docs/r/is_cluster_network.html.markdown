---
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network"
description: |-
  Manages ClusterNetwork.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_cluster_network

Create, update, and delete ClusterNetworks with this resource.

## Example Usage

```hcl
resource "ibm_is_cluster_network" "is_cluster_network_instance" {
  name = "my-cluster-network"
  profile {
		href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
		name = "h100"
		resource_type = "cluster_network_profile"
  }
  resource_group {
		href = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		id = "fee82deba12e4c0fb69c3b09d1f12345"
		name = "my-resource-group"
  }
  subnet_prefixes {
		allocation_policy = "auto"
		cidr = "10.0.0.0/24"
  }
  vpc {
		crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		id = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		name = "my-vpc"
		resource_type = "vpc"
  }
  zone {
		href = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		name = "us-south-1"
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `name` - (Optional, String) The name for this cluster network. The name must not be used by another cluster network in the region.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `profile` - (Required, List) The profile for this cluster network.
Nested schema for **profile**:
	* `href` - (Required, String) The URL for this cluster network profile.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `name` - (Required, String) The globally unique name for this cluster network profile.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_profile`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `resource_group` - (Optional, List) The resource group for this cluster network.
Nested schema for **resource_group**:
	* `href` - (Computed, String) The URL for this resource group.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this resource group.
	  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
	* `name` - (Computed, String) The name for this resource group.
	  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
* `subnet_prefixes` - (Optional, List) The IP address ranges available for subnets for this cluster network.
  * Constraints: The maximum length is `1` item. The minimum length is `1` item.
Nested schema for **subnet_prefixes**:
	* `allocation_policy` - (Required, String) The allocation policy for this subnet prefix:- `auto`: Subnets created by total count in this cluster network can use this prefix.
	  * Constraints: Allowable values are: `auto`.
	* `cidr` - (Required, String) The CIDR block for this prefix.
	  * Constraints: The maximum length is `18` characters. The minimum length is `9` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`.
* `vpc` - (Required, List) The VPC this cluster network resides in.
Nested schema for **vpc**:
	* `crn` - (Required, String) The CRN for this VPC.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Computed, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this VPC.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this VPC.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this VPC. The name is unique across all VPCs in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `zone` - (Required, List) The zone this cluster network resides in.
Nested schema for **zone**:
	* `href` - (Required, String) The URL for this zone.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `name` - (Required, String) The globally unique name for this zone.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the ClusterNetwork.
* `created_at` - (String) The date and time that the cluster network was created.
* `crn` - (String) The CRN for this cluster network.
  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
* `href` - (String) The URL for this cluster network.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **lifecycle_reasons**:
	* `code` - (String) A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `internal_error`, `resource_suspended_by_provider`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `message` - (String) An explanation of the reason for this lifecycle state.
	* `more_info` - (String) Link to documentation about the reason for this lifecycle state.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_state` - (String) The lifecycle state of the cluster network.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `cluster_network`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `etag` - ETag identifier for ClusterNetwork.

## Import

You can import the `ibm_is_cluster_network` resource by using `id`. The unique identifier for this cluster network.

# Syntax
<pre>
$ terraform import ibm_is_cluster_network.is_cluster_network &lt;id&gt;
</pre>

# Example
```
$ terraform import ibm_is_cluster_network.is_cluster_network 0767-da0df18c-7598-4633-a648-fdaac28a5573
```
