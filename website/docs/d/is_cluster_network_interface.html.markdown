---
layout: "ibm"
page_title: "IBM : ibm_is_cluster_network_interface"
description: |-
  Get information about ClusterNetworkInterface
subcategory: "Virtual Private Cloud API"
---

# ibm_is_cluster_network_interface

Provides a read-only data source to retrieve information about a ClusterNetworkInterface. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_cluster_network_interface" "is_cluster_network_interface" {
	cluster_network_id = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.cluster_network_id
	is_cluster_network_interface_id = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.is_cluster_network_interface_id
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `cluster_network_id` - (Required, Forces new resource, String) The cluster network identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `is_cluster_network_interface_id` - (Required, Forces new resource, String) The cluster network interface identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the ClusterNetworkInterface.
* `allow_ip_spoofing` - (Boolean) Indicates whether source IP spoofing is allowed on this cluster network interface. If `false`, source IP spoofing is prevented on this cluster network interface. If `true`, source IP spoofing is allowed on this cluster network interface.

* `auto_delete` - (Boolean) Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.

* `created_at` - (String) The date and time that the cluster network interface was created.

* `enable_infrastructure_nat` - (Boolean) If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.

* `href` - (String) The URL for this cluster network interface.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **lifecycle_reasons**:
	* `code` - (String) A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `internal_error`, `resource_suspended_by_provider`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `message` - (String) An explanation of the reason for this lifecycle state.
	* `more_info` - (String) Link to documentation about the reason for this lifecycle state.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `lifecycle_state` - (String) The lifecycle state of the cluster network interface.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.

* `mac_address` - (String) The MAC address of the cluster network interface. May be absent if`lifecycle_state` is `pending`.
  * Constraints: The maximum length is `17` characters. The minimum length is `17` characters. The value must match regular expression `/^([0-9A-F]{2}:){5}[0-9A-F]{2}$/`.

* `name` - (String) The name for this cluster network interface. The name is unique across all interfaces in the cluster network.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

* `primary_ip` - (List) The cluster network subnet reserved IP for this cluster network interface.
Nested schema for **primary_ip**:
	* `address` - (String) The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this cluster network subnet reserved IP.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this cluster network subnet reserved IP.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `protocol_state_filtering_mode` - (String) The protocol state filtering mode used for this cluster network interface.Protocol state filtering monitors each network connection flowing over this cluster network interface, and drops any packets that are invalid based on the current connection state and protocol. See [About Protocol State Filtering](https://cloud.ibm.com/docs/__TBD__) for more information.
  * Constraints: The default value is `enabled`. Allowable values are: `disabled`, `enabled`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `cluster_network_interface`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `subnet` - (List) 
Nested schema for **subnet**:
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this cluster network subnet.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^https:\/\/([^\/?#]*)([^?#]*)\/cluster_networks\/[-0-9a-z_]+\/subnets\/[-0-9a-z_]+$/`.
	* `id` - (String) The unique identifier for this cluster network subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `target` - (List) The target of this cluster network interface.If absent, this cluster network interface is not attached to a target.The resources supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
Nested schema for **target**:
	* `href` - (String) The URL for this instance cluster network attachment.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this instance cluster network attachment.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `instance_cluster_network_attachment`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `vpc` - (List) The VPC this cluster network interface resides in.
Nested schema for **vpc**:
	* `crn` - (String) The CRN for this VPC.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this VPC.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this VPC.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this VPC. The name is unique across all VPCs in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `zone` - (List) The zone this cluster network interface resides in.
Nested schema for **zone**:
	* `href` - (String) The URL for this zone.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `name` - (String) The globally unique name for this zone.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

