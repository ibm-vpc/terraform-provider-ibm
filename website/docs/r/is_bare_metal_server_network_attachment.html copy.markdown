---
layout: "ibm"
page_title: "IBM : ibm_is_bare_metal_server_network_attachment"
description: |-
  Manages is_bare_metal_server_network_attachment.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_bare_metal_server_network_attachment

Create, update, and delete is_bare_metal_server_network_attachments with this resource.

## Example Usage

```hcl
resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
  allow_to_float = false
  allowed_vlans = 4
  bare_metal_server_id = "bare_metal_server_id"
  interface_type = "hipersocket"
  name = "my-bare-metal-server-network-attachment"
  virtual_network_interface {
		crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
		name = "my-virtual-network-interface"
		resource_type = "virtual_network_interface"
  }
  vlan = 4
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `allow_to_float` - (Optional, Boolean) Indicates if the bare metal server network attachment can automatically float to any other server within the same `resource_group`. The bare metal server network attachment will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to bare metal server network attachments with `vlan` interface type.
* `allowed_vlans` - (Optional, List) Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) attachment.
  * Constraints: The minimum length is `0` items.
* `bare_metal_server_id` - (Required, Forces new resource, String) The bare metal server identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `interface_type` - (Required, String) The network attachment's interface type:- `hipersocket`: a virtual network device that provides high-speed TCP/IP connectivity  within a `s390x` based system- `pci`: a physical PCI device which can only be created or deleted when the bare metal  server is stopped  - Has an `allowed_vlans` property which controls the VLANs that will be permitted    to use the PCI attachment  - Cannot directly use an IEEE 802.1q VLAN tag.- `vlan`: a virtual device, used through a `pci` device that has the `vlan` in its  array of `allowed_vlans`.  - Must use an IEEE 802.1q tag.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.
  * Constraints: Allowable values are: `hipersocket`, `pci`, `vlan`.
* `name` - (Optional, String) The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `virtual_network_interface` - (Required, List) The virtual network interface for this bare metal server network attachment.
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
* `vlan` - (Optional, Integer) Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this attachment.
  * Constraints: The maximum value is `4094`. The minimum value is `1`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the is_bare_metal_server_network_attachment.
* `bare_metal_server_network_attachment_id` - (String) The unique identifier for this bare metal server network attachment.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `created_at` - (String) The date and time that the bare metal server network attachment was created.
* `href` - (String) The URL for this bare metal server network attachment.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_state` - (String) The lifecycle state of the bare metal server network attachment.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `port_speed` - (Integer) The port speed for this bare metal server network attachment in Mbps.
* `primary_ip` - (List) The primary IP address of the virtual network interface for the bare metal servernetwork attachment.
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
  * Constraints: Allowable values are: `bare_metal_server_network_attachment`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `subnet` - (List) The subnet of the virtual network interface for the bare metal server networkattachment.
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
* `type` - (String) The bare metal server network attachment type.
  * Constraints: Allowable values are: `primary`, `secondary`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.


## Import

You can import the `ibm_is_bare_metal_server_network_attachment` resource by using `id`.
The `id` property can be formed from `bare_metal_server_id`, and `id` in the following format:

```
<bare_metal_server_id>/<id>
```
* `bare_metal_server_id`: A string. The bare metal server identifier.
* `id`: A string. The bare metal server network attachment identifier.

# Syntax
```
$ terraform import ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment <bare_metal_server_id>/<id>
```
