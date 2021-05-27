---

subcategory: "Classic infrastructure"
layout: "ibm"
page_title: "IBM: network_private_ip"
description: |-
  Manages IBM Network Public IP.
---

# ibm\_network_public_ip

Provides a public IP resource to route between servers. This allows public IPs to be created, updated, and deleted. Public IPs are not restricted to routing within the same data center.

For additional details, see the [IBM Cloud Classic Infrastructure (SoftLayer) API docs](http://sldn.softlayer.com/reference/services/SoftLayer_Network_Subnet_IpAddress_Global) and [public IP address overview](https://knowledgelayer.softlayer.com/learning/global-ip-addresses).

## Example Usage

```terraform
resource "ibm_network_public_ip" "test_public_ip " {
    routes_to = "119.81.82.163"
    notes     = "public ip notes"
}
```

## Timeouts
ibm_network_public_ip provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 30 minutes) Used for Creating Instance.

## Argument Reference

The following arguments are supported:

* `routes_to` - (Required, string) The destination IP address that the public IP routes traffic through. The destination IP address can be a public IP address of IBM resources in the same account, such as a public IP address of a VM or public virtual IP addresses of NetScaler VPXs.
* `notes` - (Optional, string) Descriptive text to associate with the public IP instance.
* `tags` - (Optional, array of strings) Tags associated with the public IP instance.  

  **NOTE**: `Tags` are managed locally and not stored on the IBM Cloud service endpoint at this moment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the public IP.
* `ip_address` - The address of the public IP.
