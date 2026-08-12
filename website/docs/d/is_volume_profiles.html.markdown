---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : Volume profiles"
description: |-
  Manages IBM Cloud virtual server volume profiles.
---

# ibm_is_volume_profiles
Retrieve information of an existing IBM Cloud VSI. For more information, about the volumes and profiles, see [block storage profiles](https://cloud.ibm.com/docs/vpc?topic=vpc-block-storage-profiles).

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

data "ibm_is_volume_profiles" "example" {
}

```

## Attribute reference
You can access the following attribute references after your data source is created. 

- `profiles` - (List)  Lists all server volume profiles in the region.

  Nested scheme for `profiles`:
  - `adjustable_capacity_states` - (List) 
    Nested schema for **adjustable_capacity_states**:
    - `type` - (String) The type for this profile field.
    - `values` - (List) The attachment states that support adjustable capacity for a volume with this profile. Allowable list items are: `attached`, `unattached`, `unusable`. 
  - `adjustable_iops_states` - (List) 
    Nested schema for **adjustable_iops_states**:
    - `type` - (String) The type for this profile field.
    - `values` - (List) The attachment states that support adjustable IOPS for a volume with this profile. Allowable list items are: `attached`, `unattached`, `unusable`.
	- `name` - (String) The name of the virtual server volume profile.
	- `family` - (String) The family of the virtual server volume profile.
	- `supported_attachment_modes` - (List) The supported attachment modes for a volume with this profile. Allowable list items are: `single`, `multiple`. `single`: the volume can only be attached to a single bare metal server or instance concurrently. `multiple`: the volume can be attached to multiple bare metal servers concurrently (nvme_tcp only).
	- `supported_attachment_protocols` - (List) The supported attachment protocols for a volume with this profile. Allowable list items are: `virtio_blk`, `nvme_tcp`. `virtio_blk`: VirtIO block device for virtual server instances. `nvme_tcp`: NVMe over TCP/IP for bare metal servers.

