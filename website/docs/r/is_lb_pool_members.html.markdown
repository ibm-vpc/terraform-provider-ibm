---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : lb_pool_member"
description: |-
  Manages IBM load balancer pool member.
---

# ibm_is_lb_pool_member
Replace the existing members of the load balancer pool with new members created from the collection of member prototype for a VPC load balancer. For more information, about load balancer listener pool member, see [Creating managed pools and instance groups](https://cloud.ibm.com/docs/vpc?topic=vpc-lbaas-integration-with-instance-groups).


**Note:** 
  This resource will replace all the existing pool members in the load balancer. 

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example usage

### Sample to create a load balancer pool member for application load balancer.

```terraform
resource "ibm_is_lb_pool_members" "example2" {
  lb        = ibm_is_lb.example.id
  pool = element(split("/", ibm_is_lb_pool.testacc_pool6.id), 1)
  members {
    port   = 9545
    target = "10.80.10.5"
    weight = 61
  }
  members {
    port   = 7446
    target = "10.70.10.8"
    weight = 75
  }
  members {
    port   = 447
    target = "10.80.10.6"
    weight = 95
  }
  members {
    port   = 9563
    target = "10.90.10.4"
    weight = 91
  }
}
```

### Sample to create a load balancer pool member for network load balancer.

```terraform
resource "ibm_is_lb_pool_members" "example2" {
  lb        = ibm_is_lb.example.id
  pool = element(split("/", ibm_is_lb_pool.testacc_pool6.id), 1)
  members {
    port   = 9545
    target = ibm_is_instance.example.id
    weight = 61
  }
  members {
    port   = 7446
    target = ibm_is_instance.example1.id
    weight = 75
  }
  members {
    port   = 447
    target = ibm_is_instance.example1.id
    weight = 95
  }
  members {
    port   = 9563
    target = ibm_is_instance.example1.id
    weight = 91
  }
}
```

## Timeouts
The `ibm_is_lb_pool_members` resource provides the following [Timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - (Default 10 minutes) Used for creating Instance.
- **update** - (Default 10 minutes) Used for updating Instance.
- **delete** - (Default 10 minutes) Used for deleting Instance.


## Argument reference
Review the argument references that you can specify for your resource. 

- `lb` - (Required, Forces new resource, String) The load balancer unique identifier.
- `pool` - (Required, Forces new resource, String) The load balancer pool unique identifier.
- `members` - (List) The member prototype objects for this pool.

  Nested scheme for `members`:
    - `port`- (Required, Integer) The port number of the application running in the server member.
    - `target` - (Required, String) ID of Virtual server instance or IPv4 IP address. Load balancers in the network family support virtual server instances. Load balancers in the application family support IP addresses. 
    - `weight` - (Optional, Integer) Weight of the server member. This option takes effect only when the load-balancing algorithm of its belonging pool is `weighted_round_robin`, Minimum allowed weight is `0` and Maximum allowed weight is `100`. Default: 50, Weight of the server member. Applicable only if the pool algorithm is weighted_round_robin.


## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `created_at` - (String) The date and time that this member was created
- `id` - (String) The unique identifier of the load balancer pool member.
- `members` - (List) The member prototype objects for this pool.

  Nested scheme for `members`:
    - `port`- (Integer) The port number of the application running in the server member.
    - `target` - (String) ID of Virtual server instance or IP address. Load balancers in the network family support virtual server instances. Load balancers in the application family support IP addresses. 
    - `weight` - (Integer) Weight of the server member. This option takes effect only when the load-balancing algorithm of its belonging pool is `weighted_round_robin`, Minimum allowed weight is `0` and Maximum allowed weight is `100`. Default: 50, Weight of the server member. Applicable only if the pool algorithm is weighted_round_robin.
    - `provisioning_status` - (String) The provisioning status of this member
    - `health` - (String) Health of the server member in the pool.
    - `href` - (String) The member’s canonical URL.
    - `member` - (String) LB pool member ID


## Import
The `ibm_is_lb_pool_members` resource can be imported by using the load balancer ID, pool ID, pool member ID.

**Syntax**

```
$ terraform import ibm_is_lb_pool_members.example <loadbalancer_ID>/<pool_ID>/<pool_member_ID>
```

**Example**

```
$ terraform import ibm_is_lb_pool_member.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb/gfe6651a-bc0a-5538-8h8a-b0770bbf32cc
```
