---

subcategory: "Classic infrastructure"
layout: "ibm"
page_title: "IBM : hardware firewall shared"
description: |-
  Manages rules for IBM Firewall shared.
---

# ibm\_hardware\_firewall\_shared

Provides a firewall in IBM. One firewall protects one public VLAN and provides in-bound network packet filtering. 

<!-- You can order or find firewalls in the IBM Cloud infrastructure customer portal by navigating to **Network > IP Management > VLANs** and clicking the **Gateway/Firewall** column. -->

For more information about how to configure a firewall, see the [docs](https://knowledgelayer.softlayer.com/procedure/configure-hardware-firewall).

## Example Usage

```terraform
resource "ibm_hardware_firewall_shared" "test_firewall" {
    firewall_type="100MBPS_HARDWARE_FIREWALL"
    hardware_instance_id="12345678"
    
}
```

## Timeouts

ibm_hardware_firewall_shared provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 10 minutes) Used for creating shared firewall.

## Argument Reference

The following arguments are supported:

* `firewall_type` - (Required, string) Specifies whether it needs to be of particular speed. Firewall type is in between [10MBPS_HARDWARE_FIREWALL, 20MBPS_HARDWARE_FIREWALL,100MBPS_HARDWARE_FIREWALL, 1000MBPS_HARDWARE_FIREWALL, 200MBPS_HARDWARE_FIREWALL, 2000MBPS_HARDWARE_FIREWALL]
* `virtual_instance_id` - (Optional, string) Specifies the id of particular guest on which firewall shared is to be deployed.**NOTE**: This is conflicting parameter with hardware_instance_id.
* `hardware_instance_id` - (Optional, string) Specifies the id of particular guest on which firewall shared is to be deployed.**NOTE**: This is conflicting parameter with virtual_instance_id.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

 * `id` - The unique identifier of the hardware firewall.
