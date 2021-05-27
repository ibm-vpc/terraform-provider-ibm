---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : lb_pool"
description: |-
  Manages IBM load balancer pool.
---

# ibm\_is_lb_pool

Provides a load balancer pool resource. This allows load balancer pool to be created, updated, and cancelled.


## Example Usage

In the following example, you can create a load balancer pool:

```terraform
resource "ibm_is_lb_pool" "testacc_pool" {
  name           = "test_pool"
  lb             = "addfd-gg4r4-12345"
  algorithm      = "round_robin"
  protocol       = "http"
  health_delay   = 60
  health_retries = 5
  health_timeout = 30
  health_type    = "http"
  proxy_protocol = "v1"
}

```

In the following example, you can create a load balancer pool with `https` protocol:

```terraform
resource "ibm_is_lb_pool" "testacc_pool" {
  name           = "test_pool"
  lb             = "addfd-gg4r4-12345"
  algorithm      = "round_robin"
  protocol       = "https"
  health_delay   = 60
  health_retries = 5
  health_timeout = 30
  health_type    = "https"
  proxy_protocol = "v1"
}

```

## Timeouts

ibm_is_lb_pool provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 10 minutes) Used for creating Instance.
* `update` - (Default 10 minutes) Used for updating Instance.
* `delete` - (Default 10 minutes) Used for deleting Instance.

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the pool
* `lb` - (Required, Forces new resource, string)  The load balancer unique identifier.
* `algorithm` - (Required, string) The load balancing algorithm. Enumeration type: round_robin, weighted_round_robin, least_connections
* `protocol` - (Required, string) The pool protocol. Enumeration type: http, https, tcp
* `health_delay` - (Required, int) The health check interval in seconds. Interval must be greater than timeout value
* `health_retries` - (Required, int) The health check max retries
* `health_timeout` - (Required, int) The health check timeout in seconds
* `health_type` - (Required, string) The pool protocol. Enumeration type: http, https, tcp
* `health_monitor_url` - (Optional, string) The health check url. This option is applicable only to http type of --health-type
* `health_monitor_port` - (Optional, int) The health check port number
* `session_persistence_type` - (Optional, string) The session persistence type, Enumeration type: source_ip
* `proxy_protocol` - (Otpional, string) The PROXY protocol setting for this pool. Supported by load balancers in the application family otherwise disabled. Valid values: disabled, v1, v2.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the load balancer pool.The id is composed of \<lb_id\>/\<pool_id\>.
* `provisioning_status` - The status of load balancer pool.
* `pool_id`  - Id of the load balancer Pool
* `related_crn` - The crn of the load balancer resource.

## Import

ibm_is_lb_pool can be imported using lbID and poolID, eg

```
$ terraform import ibm_is_lb_pool.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb
```
