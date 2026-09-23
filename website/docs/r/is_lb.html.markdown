---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : load balancer"
description: |-
  Manages IBM load balancer.
---

# ibm_is_lb
Create, update, or delete a VPC Load Balancer. For more information, about VPC load balancer, see [load balancers for VPC overview](https://cloud.ibm.com/docs/vpc?topic=vpc-nlb-vs-elb).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```


## Example usage
An example to create an application load balancer.

```terraform
resource "ibm_is_lb" "example" {
  name    = "example-load-balancer"
  subnets = [ibm_is_subnet.example.id, ibm_is_subnet.example1.id]
  address_mode = "static"
  public_ips   = [ibm_is_floating_ip.example.id, ibm_is_floating_ip.example1.id]
  private_ips  = [ibm_is_subnet_reserved_ip.example.id, ibm_is_subnet_reserved_ip.example1.id]
}

```

An example to create a network load balancer.

```terraform
resource "ibm_is_lb" "example" {
  name    = "example-load-balancer"
  subnets = [ibm_is_subnet.example.id]
  profile = "network-fixed"
}

```

An example to create a load balancer with private DNS.

```terraform
resource "ibm_is_lb" "example" {
  name    = "example-load-balancer"
  subnets = [ibm_is_subnet.example.id]
  profile = "network-fixed"
  dns   {
    instance_crn = "crn:v1:staging:public:dns-svcs:global:a/exxxxxxxxxxxxx-xxxxxxxxxxxxxxxxx:5xxxxxxx-xxxxx-xxxxxxxxxxxxxxx-xxxxxxxxxxxxxxx::"
    zone_id = "bxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
  }
}

```
## An example to create a private path load balancer.
```terraform
resource "ibm_is_lb" "example" {
  name    = "example-load-balancer"
  subnets = [ibm_is_subnet.example.id]
  profile = "network-private-path"
  type = "private_path"
}
```

## Timeouts
The `ibm_is_lb` resource provides the following [Timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - (Default 30 minutes) Used for creating Instance.
- **delete** - (Default 30 minutes) Used for deleting Instance.


## Argument reference
Review the argument references that you can specify for your resource.

- `address_mode` - (Optional, String) The address mode to use for this load balancer. Supported values are `static` and `dynamic`. If `static`, customer-provided public or private IPs can be specified via `public_ips` and `private_ips`. If unset, defaults to `dynamic`.
- `access_tags`  - (Optional, List of Strings) A list of access management tags to attach to the load balancer.

  ~> **Note:** 
  **&#x2022;** You can attach only those access tags that already exists.</br>
  **&#x2022;** For more information, about creating access tags, see [working with tags](https://cloud.ibm.com/docs/account?topic=account-tag&interface=ui#create-access-console).</br>
  **&#x2022;** You must have the access listed in the [Granting users access to tag resources](https://cloud.ibm.com/docs/account?topic=account-access) for `access_tags`</br>
  **&#x2022;** `access_tags` must be in the format `key:value`.
- `dns` - (Optional, List) The DNS configuration for this load balancer.

  Nested scheme for `dns`:
  - `instance_crn` - (Required, String) The CRN of the DNS instance associated with the DNS zone
  - `zone_id` - (Required, String) The unique identifier of the DNS zone.
- `failsafe_policy_actions` - (List) The supported `failsafe_policy.action` values for this load balancer's pools. Allowable list items are: [ `bypass`, `drop`, `fail`, `forward` ]. 
    A load balancer failsafe policy action:
    - `bypass`: Bypasses the members and sends requests directly to their destination IPs.
    - `drop`: Drops requests.
    - `fail`: Fails requests with an HTTP 503 status code.
    - `forward`: Forwards requests to the target pool.

- `logging`- (Optional, Bool) Enable or disable datapath logging for the load balancer. This is applicable only for application load balancer. Supported values are **true** or **false**. Default value is **false**.
- `name` - (Required, String) The name of the VPC load balancer.
- `profile` - (Optional, Forces new resource, String) For a Network Load Balancer, this attribute is required for network and private path load balancers. Should be set to  `network-private-path` for private path load balancers and `network-fixed` for a network load balancer. For Application Load Balancer, profile is not a required attribute.
- `resource_group` - (Optional, Forces new resource, String) The resource group where the load balancer to be created.
- `route_mode` - (Optional, Forces new resource, Bool) Indicates whether route mode is enabled for this load balancer.

  ~> **NOTE:** Currently, `route_mode` enabled is supported only by private network load balancers.
- `security_groups`  (Optional, List) A list of security groups to use for this load balancer. This option is supported for both application and network load balancers.
- `subnets` - (Required, List) List of the subnets IDs to connect to the load balancer.

  ~> **NOTE:** 
  The subnets must be in the same `VPC`. The load balancer's `availability` will depend on the availability of the `zones` the specified subnets reside in. The load balancer must be in the `application` family for `updating subnets`. Load balancers in the `network` family allow only `one subnet` to be specified.

- `private_ips` - (Optional, List of String) The reserved IP IDs to assign as private IP addresses to this load balancer. Only applicable when `address_mode` is `static`.
- `public_ips` - (Optional, List of String) The floating IP IDs to assign as public IP addresses to this load balancer. Only applicable when `address_mode` is `static`.
- `tags` (Optional, Array of Strings) A list of tags that you want to add to your load balancer. Tags can help you find the load balancer more easily later.
- `type` - (Optional, Forces new resource, String) The type of the load balancer. Default value is `public`. Supported values are `public`, `private` and `private_path`.


## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `access_mode` - (String) The access mode for this load balancer. One of **private**, **public**, **private_path**.
- `attached_load_balancer_pool_members` - (List) The load balancer pool members attached to this load balancer.

  Nested scheme for `members`:
  - `deleted` - (List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.

    Nested scheme for `deleted`:
    - `more_info` - (String) Link to documentation about deleted resources.
  - `href` - (String) The URL for this load balancer pool member.
  - `id` - (String) The unique identifier for this load balancer pool member.
- `availability` - (String) The availability of this load balancer
- `crn` - (String) The CRN for this load balancer.
- `hostname` - (String) The fully qualified domain name assigned to this load balancer.
- `id` - (String) The unique identifier of the load balancer.
- `instance_groups_supported` - (Boolean) Indicates whether this load balancer supports instance groups.
- `operating_status` - (String) The operating status of this load balancer.
- `address_mode` - (String) The address mode for this load balancer. One of `static` (IPs remain unchanged throughout the life of the load balancer, horizontal scaling disabled) or `dynamic` (IPs may change during maintenance).
- `public_ips` - (List) The public IP addresses assigned to this load balancer. Will be empty if `is_public` is `false`.
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
- `private_ip` - (List) The private IP addresses assigned to this load balancer as reserved IP references.

  Nested scheme for `private_ip`:
  - `address` - (String) The IP address. If the address has not yet been selected, the value will be `0.0.0.0`. This property may expand to support IPv6 addresses in the future.
  - `href` - (String) The URL for this reserved IP.
  - `reserved_ip` - (String) The unique identifier for this reserved IP.
  - `name` - (String) The name for this reserved IP. The name is unique across all reserved IPs in a subnet.

- `private_ips` - (List) The private IP addresses assigned to this load balancer. Same as `private_ip.[].address`.
- `status` - (String) The status of the load balancer.
- `security_groups_supported`- (Bool) Indicates if this load balancer supports security groups.
- `source_ip_session_persistence_supported` - (Boolean) Indicates whether this load balancer supports source IP session persistence.
- `udp_supported`- (Bool) Indicates whether this load balancer supports UDP.


## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the `ibm_is_lb` resource by using `id`.
The `id` property can be formed from `load balancer ID`. For example:

```terraform
import {
  to = ibm_is_lb.example
  id = "<lb_ID>"
}
```

Using `terraform import`. For example:

```console
% terraform import ibm_is_lb.example <lb_ID>
```