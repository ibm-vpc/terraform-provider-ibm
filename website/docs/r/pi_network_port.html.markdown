---

subcategory: "Power Systems"
layout: "ibm"
page_title: "IBM: pi_network_port"
description: |-
  Manages an Network Port  in the Power Virtual Server Cloud. A network port is equivalent to reserving an ip in the subnet
  which can be used . When the port is created the status will be "DOWN".
  This network port however is not attached to the instance. 
---

# ibm\_pi_network_port

Provides an network_port resource. This allows network_port to be created or updated in the Power Virtual Server Cloud.

## Example Usage

In the following example, you can create an network_port:

```terraform
resource "ibm_pi_network_port" "test-network-port" {
    pi_network_name             = "Zone1-CFN"
    pi_cloud_instance_id  = "51e1879c-bcbe-4ee1-a008-49cdba0eaf60"
    pi_network_port_description         = "IP Reserved for Oracle RAC "
}
```

## Notes:
* Please find [supported Regions](https://cloud.ibm.com/apidocs/power-cloud#endpoint) for endpoints.
* If a Power cloud instance is provisioned at `lon04`, The provider level attributes should be as follows:
  * `region` - `lon`
  * `zone` - `lon04`
  Example Usage:
  ```terraform
    provider "ibm" {
      region    =   "lon"
      zone      =   "lon04"
    }
  ```
## Timeouts

ibm_pi_network_port provides the following [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 60 minutes) Used for creating a network_port.
* `delete` - (Default 60 minutes) Used for deleting a network_port.

## Argument Reference

The following arguments are supported:

* `pi_network_name` - (Required, string) The name of the PI Network.
* `pi_cloud_instance_id` - (Required, string) The GUID of the service instance associated with the account
* `pi_network_port_description` - (Optional, string) The description for the Network Port
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the instance.The id is composed of \<power_network_port_id\>/\<id\>/<pi_network_name>.
* `ipaddress` - The unique identifier of the instance.
* `macaddress` - The macaddress of the port
* `status` - The status of the port
* `portid` - The id of the port .
* `public_ip` - The public ip associated with the port

## Import

ibm_pi_network_port can be imported using `power_instance_id`, `port_id` and `pi_network_name` eg

```
$ terraform import ibm_pi_network_port.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb/network-name
```
