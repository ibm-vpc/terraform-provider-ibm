---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : load balancer"
description: |-
  Manages IBM load balancers.
---

# ibm_is_lbs
Retrieve information of an existing IBM VPC load balancers as a read-only data source. For more information, about VPC load balancer, see [load balancers for VPC overview](https://cloud.ibm.com/docs/vpc?topic=vpc-nlb-vs-elb).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example usage

```terraform
data "ibm_is_lbs" "example" {
}
```


## Attribute reference
Review the attribute references that you can access after you retrieve your data source. 

- `load_balancers` - (List) The Collection of load balancers.

	Nested scheme for `load_balancers`:
	- `access_mode` - (String) The access mode for this load balancer. One of **private**, **public**, **private_path**.
	- `access_tags`  - (String) Access management tags associated for the load balancer.
	- `attached_load_balancer_pool_members` - (List) The load balancer pool members attached to this load balancer.
		Nested scheme for `members`:
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
			Nested scheme for `deleted`:
			- `more_info` - (String) Link to documentation about deleted resources.
		- `href` - (String) The URL for this load balancer pool member.
		- `id` - (String) The unique identifier for this load balancer pool member.
	- `availability` - (String) The availability of this load balancer
	- `id` - (String) The unique identifier of the load balancer.
	- `instance_groups_supported` - (Boolean) Indicates whether this load balancer supports instance groups.
	- `created_at` - (String) The date and time this load balancer was created.
	- `crn` - (String) The load balancer's CRN.
	- `dns` - (List) The DNS configuration for this load balancer.

		Nested scheme for `dns`:
		- `instance_crn` - (String) The CRN of the DNS instance associated with the DNS zone
		- `zone_id` - (String) The unique identifier of the DNS zone.
	- `failsafe_policy_actions` - (List) The supported `failsafe_policy.action` values for this load balancer's pools. Allowable list items are: [ `bypass`, `drop`, `fail`, `forward` ]. A load balancer failsafe policy action:

		- `bypass`: Bypasses the members and sends requests directly to their destination IPs.
		- `drop`: Drops requests.
		- `fail`: Fails requests with an HTTP 503 status code.
		- `forward`: Forwards requests to the target pool.

	- `name` - (String) Name of the load balancer.
	- `subnets` - (List) The subnets this load balancer is part of.

		Nested scheme for `subnets`:
		- `crn` - (String) The CRN for the subnet.
		- `id` - (String) The unique identifier for this subnet.
		- `href` - (String) The URL for this subnet.
		- `name` - (String) The user-defined name for this subnet.
	- `hostname` - (String) The Fully qualified domain name assigned to this load balancer.
	- `listeners` - (List) The listeners of this load balancer.

		Nested scheme for `listeners`:
		- `id` - (String) The unique identifier for this load balancer listener.
		- `href` - (String) The listener's canonical URL.
	- `operating_status` - (String) The operating status of this load balancer.
	- `pools` - (List) The pools of this load balancer.

		Nested scheme for `pools`:
		- `href` - (String) The pool's canonical URL.
		- `id` - (String) The unique identifier for this load balancer pool.
		- `name` - (String) The user-defined name for this load balancer pool.
	- `profile` - (List) The profile to use for this load balancer.

		Nested scheme for `profile`:
		- `family` - (String) The product family this load balancer profile belongs to.
		- `href` - (String) The URL for this load balancer profile.
		- `name` - (String) The name for this load balancer profile.
	- `private_ip` - (List) The private IP addresses assigned to this load balancer as reserved IP references. Will be empty if `is_public` is `true`.

		Nested scheme for `private_ip`:
		- `address` - (String) The IP address. If the address has not yet been selected, the value will be `0.0.0.0`. This property may expand to support IPv6 addresses in the future.
		- `href` - (String) The URL for this reserved IP.
		- `name` - (String) The name for this reserved IP. The name is unique across all reserved IPs in a subnet.
		- `reserved_ip` - (String) The unique identifier for this reserved IP.
		- `resource_type` - (String) The resource type. Always `subnet_reserved_ip`.
	- `address_mode` - (String) The address mode for this load balancer. One of `static` (IPs remain unchanged throughout the life of the load balancer, horizontal scaling disabled) or `dynamic` (IPs may change during maintenance).
	- `private_ips` - (List of String) The private IP addresses assigned to this load balancer. Same as `private_ip.[].address`. Will be empty if `is_public` is `true`.
	- `provisioning_status` - (String) The provisioning status of this load balancer. Possible values are: **active**, **create_pending**, **delete_pending**, **failed**, **maintenance_pending**, **update_pending**, **migrate_pending**.
	- `public_ips` - (List of String) The public IP addresses assigned to this load balancer. Will be empty if `is_public` is `false`.
	- `public_ip` - (List) The public IP address details assigned to this load balancer. Each entry is either a floating IP reference or a plain IP address.

		Nested scheme for `public_ip`:
		- `address` - (String) The globally unique IP address. This property may expand to support IPv6 addresses in the future.
		- `crn` - (String) The CRN for this floating IP. Present only when the public IP is a floating IP.
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.

			Nested scheme for `deleted`:
			- `more_info` - (String) Link to documentation about deleted resources.
		- `href` - (String) The URL for this floating IP. Present only when the public IP is a floating IP.
		- `id` - (String) The unique identifier for this floating IP. Present only when the public IP is a floating IP.
		- `name` - (String) The name for this floating IP. The name is unique across all floating IPs in the region. Present only when the public IP is a floating IP.
	- `resource_group` - (String) The resource group id, where the load balancer is created.
	- `route_mode` - (Bool) Indicates whether route mode is enabled for this load balancer.
	- `source_ip_session_persistence_supported` - (Boolean) Indicates whether this load balancer supports source IP session persistence.
	- `status` - (String) The status of the load balancers.
	- `type` - (String) The type of the load balancer.
	- `tags` - (String) Tags associated with the load balancer.
	- `udp_supported`- (Bool) Indicates whether this load balancer supports UDP.
