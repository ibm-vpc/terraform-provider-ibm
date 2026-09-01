---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : public_address_range"
description: |-
  Manages IBM public address range.
---

# ibm_is_public_address_range

Create, update, and delete a public address range. For more information, see [creating public address range](https://cloud.ibm.com/docs/vpc?topic=vpc-par-creating&interface=ui).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example Usage
The following example shows how you can create a public address range for a vpc that are located in a specific zone.

```terraform
resource "ibm_is_public_address_range" "public_address_range_instance" {
  ipv4_address_count = "16"
  name               = "example-public-address-range"
  resource_group {
    id = "11caaa983d9c4beb82690daab18717e9"
  }
  target {
    vpc {
      id = ibm_is_vpc.testacc_vpc.id
    }
    zone {
      name = "us-south-3"
    }
  }
}
```

An example shows how you can create a public address range using a specific CIDR block.

```terraform
resource "ibm_is_public_address_range" "public_address_range_cidr" {
  cidr = "192.0.2.0/24"
  name = "example-public-address-range-cidr"
  resource_group {
    id = "11caaa983d9c4beb82690daab18717e9"
  }
}
```

An example shows how you can create public address range not attached to vpc and zone

```terraform
resource "ibm_is_public_address_range" "public_address_range_instance" {
  ipv4_address_count = "16"
  name               = "example-public-address-range"
  resource_group {
    id = "11caaa983d9c4beb82690daab08717e9"
  }
}
```

An example shows how you can create an IPv6 public address range bound to an authorized CIDR and a virtual network interface target.

```terraform
resource "ibm_is_public_address_range_authorized_cidr" "example_authorized_cidr" {
  ip_version            = "ipv6"
  availability_mode     = "zonal"
  zone                  = "us-south-1"
  network_prefix_length = 64
  name                  = "example-authorized-cidr"
}

resource "ibm_is_virtual_network_interface" "example_vni" {
  name   = "example-vni"
  subnet = ibm_is_subnet.example.id
}

resource "ibm_is_public_address_range" "public_address_range_ipv6" {
  authorized_cidr       = ibm_is_public_address_range_authorized_cidr.example_authorized_cidr.id
  network_prefix_length = 64
  name                  = "example-public-address-range-ipv6"
  resource_group {
    id = "11caaa983d9c4beb82690daab18717e9"
  }
  target {
    virtual_network_interface {
      id = ibm_is_virtual_network_interface.example_vni.id
    }
    zone {
      name = "us-south-1"
    }
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

- `authorized_cidr` - (Optional, Forces new resource, String) The ID of the public address range authorized CIDR from which this public address range is allocated. Required when creating an IPv6 public address range. Mutually exclusive with `cidr` and `ipv4_address_count`.
- `cidr` - (Optional, Forces new resource, String) The public IPv4 range, expressed in CIDR format. Mutually exclusive with `ipv4_address_count`. If not specified, an IP block is automatically allocated.

  ~> **Note:** Exactly one of `cidr` or `ipv4_address_count` must be specified for IPv4 public address ranges.
- `ipv4_address_count` - (Optional, Computed, Integer) The number of IPv4 addresses in this public address range. Mutually exclusive with `cidr`.

  ~> **Note:** Exactly one of `cidr` or `ipv4_address_count` must be specified for IPv4 public address ranges.
- `name` - (Optional, String) The name for this public address range. The name is unique across all public address ranges in the region.
- `network_prefix_length` - (Optional, Forces new resource, Integer) The network prefix length of the CIDR block to allocate from the authorized CIDR. Required when `authorized_cidr` is specified.
- `resource_group` - (Optional, List) The resource group for this public address range.
    
	Nested schema for `resource_group`:
	- `id` - (Required, String) The unique identifier for this resource group.
- `target` - (Optional, List) The target this public address range is bound to. If absent, this public address range is not bound to a target.
    
	Nested schema for `target`:
	- `virtual_network_interface` - (Optional, List) The virtual network interface this public address range is bound to. Used for IPv6 public address ranges. Specify `id` to bind to an existing virtual network interface.

		Nested schema for `virtual_network_interface`:
		- `id` - (Optional, String) The unique identifier for this virtual network interface.
	- `vpc` - (Optional, List) The VPC this public address range is bound to. If present, any of the below value must be specified.
	    
		Nested schema for `vpc`:
		- `crn` - (Optional, String) The CRN for this VPC.
		- `href` - (Optional, String) The URL for this VPC.
		- `id` - (Optional, String) The unique identifier for this VPC.
	- `zone` - (Optional, List) The zone this public address range resides in. If present, any of the below value must be specified.
	    
		Nested schema for `zone`:
		- `href` - (Optional, String) The URL for this zone.
		- `name` - (Optional, String) The globally unique name for this zone.

## Attribute Reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - The unique identifier of the PublicAddressRange.
- `authorized_cidr` - (String) When set as an argument, the ID of the authorized CIDR from which this public address range is allocated. When read as an attribute, the authorized CIDR reference fields are available as sub-attributes.

  Nested schema for `authorized_cidr` (computed attributes):
  - `cidr` - (String) The CIDR block of the authorized CIDR.
  - `crn` - (String) The CRN for this authorized CIDR.
  - `href` - (String) The URL for this authorized CIDR.
  - `id` - (String) The unique identifier for this authorized CIDR.
  - `name` - (String) The name for this authorized CIDR.
  - `resource_type` - (String) The resource type.
- `cidr` - (String) The public IP range, expressed in CIDR format.
- `created_at` - (String) The date and time that the public address range was created.
- `crn` - (String) The CRN for this public address range.
- `href` - (String) The URL for this public address range.
- `ip_version` - (String) The IP version for this public address range (`ipv4` or `ipv6`).
- `lifecycle_state` - (String) The lifecycle state of the public address range.
- `network_prefix_length` - (Integer) The network prefix length of the CIDR block.
- `profile` - (List) The profile for this public address range.

  Nested schema for `profile`:
  - `href` - (String) The URL for this public address range profile.
  - `name` - (String) The globally unique name for this public address range profile.
  - `resource_type` - (String) The resource type.
- `resource_type` - (String) The resource type.
- `resource_group` - (List) The resource group for this public address range.
    
	Nested schema for `resource_group`:
	- `href` - (String) The URL for this resource group.
	- `id` - (String) The unique identifier for this resource group.
	- `name` - (String) The name for this resource group.
- `target` - (List) The target this public address range is bound to. If absent, this public address range is not bound to a target.
    
	Nested schema for `target`:
	- `virtual_network_interface` - (List) The virtual network interface this public address range is bound to (IPv6 only).

		Nested schema for `virtual_network_interface`:
		- `crn` - (String) The CRN for this virtual network interface.
		- `href` - (String) The URL for this virtual network interface.
		- `id` - (String) The unique identifier for this virtual network interface.
		- `name` - (String) The name for this virtual network interface.
		- `resource_type` - (String) The resource type.
	- `vpc` - (List) The VPC this public address range is bound to.
	    
		Nested schema for `vpc`:
		- `crn` - (String) The CRN for this VPC.
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.
			
			Nested schema for `deleted`:
			- `more_info` - (Computed, String) Link to documentation about deleted resources.
		- `href` - (String) The URL for this VPC.
		- `id` - (String) The unique identifier for this VPC.
		- `name` - (Computed, String) The name for this VPC. The name is unique across all VPCs in the region.
		- `resource_type` - (Computed, String) The resource type.
	- `zone` - (List) The zone this public address range resides in.
	    
		Nested schema for `zone`:
		- `href` - (String) The URL for this zone.
		- `name` - (String) The globally unique name for this zone.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the `ibm_is_public_address_range` resource by using `id`.
The `id` property can be formed using the public_address_range id. For example:

```terraform
import {
  to = ibm_is_public_address_range.example
  id = "<id>"
}
```

Using `terraform import`. For example:

```console
% terraform import ibm_is_public_address_range.example <id>
```