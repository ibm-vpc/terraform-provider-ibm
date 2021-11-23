---
subcategory: "VPC infrastructure"
page_title: "IBM : ibm_is_lb_listener"
description: |-
  Get information about LoadBalancerListener
---

# ibm_is_lb_listener

Provides a read-only data source for LoadBalancerListener.For more information, about load balancer listener, see [working with listeners](https://cloud.ibm.com/docs/vpc?topic=vpc-nlb-listeners).

## Example Usage

```terraform
data "ibm_is_lb_listener" "example" {
	listener_id = ibm_is_lb_listener.example.listener_id
	lb = ibm_is_lb.example.id
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

- `listener_id` - (Required, String) The listener identifier.
- `lb` - (Required, String) The load balancer identifier.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

- `id` - The unique identifier of the LoadBalancerListener.
- `accept_proxy_protocol` - (Required, Boolean) If set to `true`, this listener will accept and forward PROXY protocol information. Supported by load balancers in the `application` family (otherwise always `false`). Additional restrictions:- If this listener has `https_redirect` specified, its `accept_proxy_protocol` value must  match the `accept_proxy_protocol` value of the `https_redirect` listener.- If this listener is the target of another listener's `https_redirect`, its  `accept_proxy_protocol` value must match that listener's `accept_proxy_protocol` value.

- `certificate_instance` - (Optional, List) The certificate instance used for SSL termination. It is applicable only to `https`protocol.
Nested scheme for `certificate_instance`:
	- `crn` - (Required, String) The CRN for this certificate instance.

- `connection_limit` - (Optional, Integer) The connection limit of the listener.

- `created_at` - (Required, String) The date and time that this listener was created.

- `default_pool` - (Optional, List) The default pool associated with the listener.
Nested scheme for `default_pool`:
	- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
	Nested scheme for `deleted`:
		- `more_info` - (Required, String) Link to documentation about deleted resources.
	- `href` - (Required, String) The pool's canonical URL.
	- `id` - (Required, String) The unique identifier for this load balancer pool.
	- `name` - (Required, String) The user-defined name for this load balancer pool.

- `href` - (Required, String) The listener's canonical URL.

- `https_redirect` - (Optional, List) If specified, the target listener that requests are redirected to.
Nested scheme for `https_redirect`:
	- `http_status_code` - (Required, Integer) The HTTP status code for this redirect.
	- `listener` - (Required, List)
	Nested scheme for `listener`:
		- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
		Nested scheme for `deleted`:
			- `more_info` - (Required, String) Link to documentation about deleted resources.
		- `href` - (Required, String) The listener's canonical URL.
		- `id` - (Required, String) The unique identifier for this load balancer listener.
	- `uri` - (Optional, String) The redirect relative target URI.

- `policies` - (Optional, List) The policies for this listener.
Nested scheme for `policies`:
	- `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
	Nested scheme for `deleted`:
		- `more_info` - (Required, String) Link to documentation about deleted resources.
	- `href` - (Required, String) The listener policy's canonical URL.
	- `id` - (Required, String) The policy's unique identifier.

- `port` - (Required, Integer) The listener port number, or the inclusive lower bound of the port range. Each listener in the load balancer must have a unique `port` and `protocol` combination.

- `port_max` - (Required, Integer) The inclusive upper bound of the range of ports used by this listener.Only load balancers in the `network` family support more than one port per listener.

- `port_min` - (Required, Integer) The inclusive lower bound of the range of ports used by this listener.Only load balancers in the `network` family support more than one port per listener.

- `protocol` - (Required, String) The listener protocol. Load balancers in the `network` family support `tcp`. Load balancers in the `application` family support `tcp`, `http`, and `https`. Each listener in the load balancer must have a unique `port` and `protocol` combination.

- `provisioning_status` - (Required, String) The provisioning status of this listener.

