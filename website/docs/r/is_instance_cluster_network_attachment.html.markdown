---
layout: "ibm"
page_title: "IBM : ibm_is_instance_cluster_network_attachment"
description: |-
  Manages InstanceClusterNetworkAttachment.
subcategory: "VPC infrastructure"
---

# ibm_is_instance_cluster_network_attachment

Create, update, and delete InstanceClusterNetworkAttachments with this resource.

## Example Usage

```hcl
resource "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
  before {
		id = "0717-fb880975-db45-4459-8548-64e3995ac213"
  }
  cluster_network_interface {
		name = "my-cluster-network-interface"
		primary_ip {
			address = "10.1.0.6"
			name = "my-cluster-network-subnet-reserved-ip"
		}
		subnet {
			id = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		}
  }
  instance_id = "instance_id"
  name = "my-instance-network-attachment"
}
```

## Argument Reference

You can specify the following arguments for this resource.

- `before` - (Optional, List) The instance cluster network attachment that is immediately before. If absent, this is thelast instance cluster network attachment.
	Nested schema for **before**:
	- `href` - (Required, String) The URL for this instance cluster network attachment.
	- `id` - (Required, String) The unique identifier for this instance cluster network attachment.
	- `name` - (Computed, String) The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.
	- `resource_type` - (Computed, String) The resource type.
- `cluster_network_interface` - (Required, List) The cluster network interface for this instance cluster network attachment.
Nested schema for **cluster_network_interface**:
	- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		- `more_info` - (Computed, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `href` - (Required, String) The URL for this cluster network interface.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	- `id` - (Required, String) The unique identifier for this cluster network interface.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	- `name` - (Required, String) The name for this cluster network interface. The name is unique across all interfaces in the cluster network.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	- `primary_ip` - (Required, List) The primary IP for this cluster network interface.
	Nested schema for **primary_ip**:
		- `address` - (Required, String) The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
		- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (Computed, String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		- `href` - (Required, String) The URL for this cluster network subnet reserved IP.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		- `id` - (Required, String) The unique identifier for this cluster network subnet reserved IP.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		- `name` - (Required, String) The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		- `resource_type` - (Computed, String) The resource type.
		  * Constraints: Allowable values are: `cluster_network_subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `cluster_network_interface`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `subnet` - (Required, List)
	Nested schema for **subnet**:
		- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (Computed, String) Link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		- `href` - (Required, String) The URL for this cluster network subnet.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^https:\/\/([^\/?#]*)([^?#]*)\/cluster_networks\/[-0-9a-z_]+\/subnets\/[-0-9a-z_]+$/`.
		- `id` - (Required, String) The unique identifier for this cluster network subnet.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		- `name` - (Computed, String) The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		- `resource_type` - (Computed, String) The resource type.
		  * Constraints: Allowable values are: `cluster_network_subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
- `instance_id` - (Required, Forces new resource, String) The virtual server instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
- `name` - (Optional, String) The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

- `id` - The unique identifier of the InstanceClusterNetworkAttachment.
- `href` - (String) The URL for this instance cluster network attachment.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
- `instance_cluster_network_attachment_id` - (String) The unique identifier for this instance cluster network attachment.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
- `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).
	Nested schema for **lifecycle_reasons**:
	- `code` - (String) A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future. Allowable values are: `internal_error`, `resource_suspended_by_provider`. 
	- `message` - (String) An explanation of the reason for this lifecycle state.
	- `more_info` - (String) Link to documentation about the reason for this lifecycle state.
- `lifecycle_state` - (String) The lifecycle state of the instance cluster network attachment. Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
- `resource_type` - (String) The resource type. Allowable values are: `instance_cluster_network_attachment`.


## Import

You can import the `ibm_is_instance_cluster_network_attachment` resource by using `id`.
The `id` property can be formed from `instance_id`, and `instance_cluster_network_attachment_id` in the following format:

<pre>
&lt;instance_id&gt;/&lt;instance_cluster_network_attachment_id&gt;
</pre>
- `instance_id`: A string. The virtual server instance identifier.
- `instance_cluster_network_attachment_id`: A string in the format `0717-fb880975-db45-4459-8548-64e3995ac213`. The unique identifier for this instance cluster network attachment.

# Syntax
<pre>
$ terraform import ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment &lt;instance_id&gt;/&lt;instance_cluster_network_attachment_id&gt;
</pre>
