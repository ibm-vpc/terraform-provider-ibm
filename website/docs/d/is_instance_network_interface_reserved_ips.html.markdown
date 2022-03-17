---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : reserved_ips"
description: |-
  Lists all the info in reserved IP for Instance network interface.
---

# ibm\_is_instance_network_interface_reserved_ips

Import the details of all the Reserved IPs in a network interface of an instance as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm__is_instance_network_interface_reserved_ips" "data_reserved_ips" {
  instance_id = ibm_is_instance.test_instance.id
  network_interface_id = ibm_is_instance.test_instance.network_interfaces.0.id
}
```

## Argument Reference

The following arguments are supported as inputs/request params:

* `instance_id` - (Required, string) The id for the instance.
* `network_interface_id` - (Required, string) The id for the network interface.


## Attribute Reference

The following attributes are exported as output/response:

* `id` - The id for the all the reserved ID (current timestamp)
* `limit` - The number of reserved IPs to list
* `reserved_ips` - The collection of all the reserved IPs in the network inetrface
   - `address` - The IP bound for the reserved IP
   - `auto_delete` - If reserved ip shall be deleted automatically
   - `created_at` - The date and time that the reserved IP was created
   - `href` - The URL for this reserved IP
   - `reserved_ip` - The unique identifier for this reserved IP
   - `name` - The user-defined or system-provided name for this reserved IP
   - `owner` - The owner of a reserved IP, defining whether it is managed by the user or the provider
   - `resource_type` - The resource type
   - `target` - The id for the target for the reserved IP.

* `sort` - The keyword on which all the reserved IPs are sorted
* `instance_id` - The id of the instance for the reserved IP
* `network_interface_id` - The id of the network interface for the reserved IP
* `total_count` - The number of reserved IP in the network interface of the instance
