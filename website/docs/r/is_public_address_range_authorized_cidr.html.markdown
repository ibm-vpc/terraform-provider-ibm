---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_public_address_range_authorized_cidr"
description: |-
  Manages IBM public address range authorized CIDR.
---

# ibm_is_public_address_range_authorized_cidr

Create, update, and delete a public address range authorized CIDR. For more information, see [public address ranges](https://cloud.ibm.com/docs/vpc?topic=vpc-par-creating&interface=ui).

**Note:**
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example Usage

```terraform
resource "ibm_is_public_address_range_authorized_cidr" "example" {
  ip_version            = "ipv6"
  availability_mode     = "zonal"
  zone                  = "us-south-1"
  network_prefix_length = 64
  name                  = "example-authorized-cidr"
}
```

## Argument Reference

You can specify the following arguments for this resource.

- `ip_version` - (Required, Forces new resource, String) The IP version for this public address range authorized CIDR. Currently only `ipv6` is supported.
- `availability_mode` - (Required, Forces new resource, String) The availability mode of the public address range authorized CIDR. Currently only `zonal` is supported.
- `zone` - (Required, Forces new resource, String) The globally unique name of the zone this public address range authorized CIDR will reside in.
- `network_prefix_length` - (Required, Forces new resource, Integer) The network prefix length for this public address range authorized CIDR. Currently only `64` is supported.
- `name` - (Optional, String) The name for this public address range authorized CIDR. The name must not be used by another public address range authorized CIDR in the region. Names beginning with `ibm-` are reserved for provider-managed resources, and are not allowed.

## Attribute Reference

In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of the PublicAddressRangeAuthorizedCIDR.
- `allocation` - (List) The allocation for this public address range authorized CIDR.

  Nested schema for `allocation`:
  - `count` - (Integer) The number of resources allocated from this public address range authorized CIDR.
  - `profile_family` - (String) The profile family for resources allocated from this public address range authorized CIDR.
- `cidr` - (String) The public IP address block for the public address range authorized CIDR, expressed in CIDR format.
- `crn` - (String) The CRN for this public address range authorized CIDR.
- `href` - (String) The URL for this public address range authorized CIDR.
- `lifecycle_reasons` - (List) The reasons for the current lifecycle_state (if any).

  Nested schema for `lifecycle_reasons`:
  - `code` - (String) A reason code for this lifecycle state.
  - `message` - (String) An explanation of the reason for this lifecycle state.
  - `more_info` - (String) A link to documentation about the reason for this lifecycle state.
- `lifecycle_state` - (String) The lifecycle state of the public address range authorized CIDR.
- `resource_group` - (List) The resource group for this public address range authorized CIDR.

  Nested schema for `resource_group`:
  - `href` - (String) The URL for this resource group.
  - `id` - (String) The unique identifier for this resource group.
  - `name` - (String) The user-defined name for this resource group.
- `resource_type` - (String) The resource type.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the `ibm_is_public_address_range_authorized_cidr` resource by using `id`.
The `id` property can be formed from the public address range authorized CIDR id. For example:

```terraform
import {
  to = ibm_is_public_address_range_authorized_cidr.example
  id = "<id>"
}
```

Using `terraform import`. For example:

```console
% terraform import ibm_is_public_address_range_authorized_cidr.example <id>
```
