---
layout: "ibm"
page_title: "IBM : ibm_is_instance_network_attachment"
description: |-
  Manages InstanceByNetworkAttachment.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_instance_network_attachment

Create, update, and delete InstanceByNetworkAttachments with this resource.

## Example Usage

```terraform
resource "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
  instance_id = "instance_id"
  virtual_network_interface {
		crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		name = "my-virtual-network-interface"
		resource_type = "virtual_network_interface"
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `instance_id` - (Required, Forces new resource, String) The virtual server instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `name` - (Optional, String) The name for this instance network attachment. The name is unique across all network attachments for the instance.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `virtual_network_interface` - (Required, List) The virtual network interface for this instance network attachment.
Nested schema for **virtual_network_interface**:
	* `crn` - (Required, String) The CRN for this virtual network interface.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `href` - (Required, String) The URL for this virtual network interface.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^https:\/\/([^\/?#]*)([^?#]*)\/virtual_network_interfaces\/[-0-9a-z_]+$/`.
	* `id` - (Required, String) The unique identifier for this virtual network interface.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `virtual_network_interface`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the InstanceByNetworkAttachment.
* `created_at` - (String) The date and time that the instance network attachment was created.
* `href` - (String) The URL for this instance network attachment.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `instance_by_network_attachment_id` - (String) The unique identifier for this instance network attachment.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `lifecycle_state` - (String) The lifecycle state of the instance network attachment.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `port_speed` - (Integer) The port speed for this instance network attachment in Mbps.
* `primary_ip` - (List) The primary IP address of the virtual network interface for the instance networkattachment.
Nested schema for **primary_ip**:
	* `address` - (String) The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this reserved IP.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this reserved IP.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this reserved IP. The name is unique across all reserved IPs in a subnet.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `subnet_reserved_ip`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `instance_network_attachment`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `subnet` - (List) The subnet of the virtual network interface for the instance network attachment.
Nested schema for **subnet**:
	* `crn` - (String) The CRN for this subnet.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this subnet.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this subnet. The name is unique across all subnets in the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `type` - (String) The instance network attachment type.
  * Constraints: Allowable values are: `primary`, `secondary`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.


## Import

You can import the `ibm_is_instance_network_attachment` resource by using `id`.
The `id` property can be formed from `instance_id`, and `id` in the following format:

```
<instance_id>/<id>
```
* `instance_id`: A string. The virtual server instance identifier.
* `id`: A string. The instance network attachment identifier.

# Syntax
```
$ terraform import ibm_is_instance_network_attachment.is_instance_network_attachment <instance_id>/<id>
```
