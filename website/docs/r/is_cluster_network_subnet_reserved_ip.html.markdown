---
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network_subnet_reserved_ip"
description: |-
  Manages ClusterNetworkSubnetReservedIP.
subcategory: "VPC infrastructure"
---

# ibm_is_cluster_network_subnet_reserved_ip

Create, update, and delete ClusterNetworkSubnetReservedIPs with this resource.

## Example Usage

```hcl
resource "ibm_is_cluster_network_subnet_reserved_ip" "is_cluster_network_subnet_reserved_ip_instance" {
  address = "192.168.3.4"
  cluster_network_id = "cluster_network_id"
  cluster_network_subnet_id = "cluster_network_subnet_id"
  name = "my-cluster-network-subnet-reserved-ip"
}
```

## Argument Reference

You can specify the following arguments for this resource.

- `address` - (Optional, String) The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.
  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
- `cluster_network_id` - (Required, Forces new resource, String) The cluster network identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
- `cluster_network_subnet_id` - (Required, Forces new resource, String) The cluster network subnet identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
- `name` - (Optional, String) The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

- `id` - The unique identifier of the ClusterNetworkSubnetReservedIP.
- `auto_delete` - (Boolean) Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.
- `created_at` - (String) The date and time that the cluster network subnet reserved IP was created.
- `href` - (String) The URL for this cluster network subnet reserved IP.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
- `cluster_network_subnet_reserved_ip_id` - (String) The unique identifier for this cluster network subnet reserved IP.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
- `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **lifecycle_reasons**:
	- `code` - (String) A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `internal_error`, `resource_suspended_by_provider`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `message` - (String) An explanation of the reason for this lifecycle state.
	- `more_info` - (String) Link to documentation about the reason for this lifecycle state.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
- `lifecycle_state` - (String) The lifecycle state of the cluster network subnet reserved IP.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `owner` - (String) The owner of the cluster network subnet reserved IPThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `provider`, `user`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `cluster_network_subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `target` - (List) The target this cluster network subnet reserved IP is bound to.If absent, this cluster network subnet reserved IP is provider-owned or unbound.
Nested schema for **target**:
	- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		- `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `href` - (String) The URL for this cluster network interface.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `id` - (String) The unique identifier for this cluster network interface.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	- `name` - (String) The name for this cluster network interface. The name is unique across all interfaces in the cluster network.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	- `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_interface`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

- `etag` - ETag identifier for ClusterNetworkSubnetReservedIP.

## Import

You can import the `ibm_is_cluster_network_subnet_reserved_ip` resource by using `id`.
The `id` property can be formed from `cluster_network_id`, `cluster_network_subnet_id`, and `cluster_network_subnet_reserved_ip_id` in the following format:

<pre>
&lt;cluster_network_id&gt;/&lt;cluster_network_subnet_id&gt;/&lt;cluster_network_subnet_reserved_ip_id&gt;
</pre>
- `cluster_network_id`: A string. The cluster network identifier.
- `cluster_network_subnet_id`: A string. The cluster network subnet identifier.
- `cluster_network_subnet_reserved_ip_id`: A string in the format `6d353a0f-aeb1-4ae1-832e-1110d10981bb`. The unique identifier for this cluster network subnet reserved IP.

# Syntax
<pre>
$ terraform import ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip &lt;cluster_network_id&gt;/&lt;cluster_network_subnet_id&gt;/&lt;cluster_network_subnet_reserved_ip_id&gt;
</pre>
