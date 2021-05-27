---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : load balancer"
description: |-
  Manages IBM load balancer.
---

# ibm\_is_lb

Provides a load balancer resource. This allows load balancer to be created, updated, and cancelled.


## Example Usage

In the following example, you can create an application load balancer:

```terraform
resource "ibm_is_lb" "lb" {
  name    = "loadbalancer1"
  subnets = ["04813493-15d6-4150-9948-6cc646cb67f2"]
}

```

In the following example, you can create a network load balancer:

```terraform
resource "ibm_is_lb" "lb" {
  name    = "loadbalancer1"
  subnets = ["04813493-15d6-4150-9948-6cc646cb67f2"]
  profile = "network-fixed"
}

```

## Timeouts

ibm_is_lb provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 30 minutes) Used for creating Instance.
* `delete` - (Default 30 minutes) Used for deleting Instance.

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) Name of the loadbalancer.
* `subnets` - (Required, list) ID of the subnets to provision this load balancer.
* `profile` - (Required for network load balancer, Forces new resource, string) The profile to use for this load balancer. Supported value `network-fixed`.  
* `type` - (Optional, Forces new resource, string) The type of the load balancer. Default value `public`. Supported values `public` and  `private`.
* `resource_group` - (Optional, Forces new resource, string) The resource group where the load balancer to be created.
* `tags` - (Optional, array of strings) Tags associated with the load balancer.
* `logging` - (Optional, bool) Enable or disable datapath logging for this load balancer. If unspecified, datapath logging is disabled. This is applicable only for application load balancer. One of: false, true.
* `security_groups` - (Optional, list) The security groups to use for this load balancer.This is applicable only for application load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the load balancer.
* `public_ips` - The public IP addresses assigned to this load balancer.
* `private_ips` - The private IP addresses assigned to this load balancer.
* `status` - The status of load balancer.
* `operating_status` - The operating status of this load balancer.
* `hostname` - Fully qualified domain name assigned to this load balancer.
* `security_groups_supported` - Indicates whether this load balancer supports security groups.



## Import

ibm_is_lb can be imported using lbID, eg

```
$ terraform import ibm_is_lb.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
