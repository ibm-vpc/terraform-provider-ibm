---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : bare_metal_server_volume_attachment"
description: |-
  Manages IBM bare metal server volume attachment.
---

# ibm\_is_bare_metal_server_volume_attachment

Create, update, and delete a volume attachment for an IBM Cloud Bare Metal Server. The volume must use the `sdp` profile to be compatible with bare metal server attachment. Volume attachments allow bare metal servers to access block storage volumes over the NVMe-over-TCP protocol. For more information, about managing VPC Bare Metal Server, see [About Bare Metal Servers for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-bare-metal-servers).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example Usage

In the following example, you can attach an existing volume to a bare metal server:

```terraform
resource "ibm_is_bare_metal_server_volume_attachment" "example" {
  bare_metal_server                         = ibm_is_bare_metal_server.example.id
  volume                                    = ibm_is_volume.example.id
  name                                      = "example-bms-vol-att"
  delete_volume_on_bare_metal_server_delete = false
  delete_volume_on_attachment_delete        = false
}
```

## Example Usage with new inline volume

In the following example, you can create a new volume and attach it to a bare metal server in one step:

```terraform
resource "ibm_is_bare_metal_server_volume_attachment" "example" {
  bare_metal_server                         = ibm_is_bare_metal_server.example.id
  name                                      = "example-bms-vol-att"
  profile                                   = "sdp"
  capacity                                  = 10000
  iops                                      = 10000
  delete_volume_on_bare_metal_server_delete = false
  delete_volume_on_attachment_delete        = true
}
```

## Argument Reference

Review the argument references that you can specify for your resource. 

- `bare_metal_server` - (Required, Forces new resource, String) The unique identifier of the bare metal server.
- `name` - (Optional, String) The name for this volume attachment. The name must be unique across all volume attachments on the bare metal server.
- `volume` - (Optional, Forces new resource, String) The unique identifier of an existing volume to attach. Conflicts with `capacity`, `profile`, `iops`, `encryption_key`, `attachment_mode`, `resource_group`, `user_tags`, `allowed_use`, `source_snapshot`, `bandwidth`, and `volume_name`.
- `capacity` - (Optional, Integer) The capacity of the volume in gigabytes. At least one of `volume`, `capacity`, or `source_snapshot` must be specified. Conflicts with `volume`.
- `profile` - (Optional, String) The profile name for the volume. Bare metal server volume attachments support only the `sdp` profile. Conflicts with `volume`.
- `iops` - (Optional, Integer) The maximum I/O operations per second for the volume. Applicable only to the `sdp` profile. Conflicts with `volume`.
- `bandwidth` - (Optional, Integer) The maximum bandwidth (in megabits per second) for the volume. Conflicts with `volume`.
- `encryption_key` - (Optional, Forces new resource, String) The CRN of the encryption key (BYOK) for the volume. Conflicts with `volume`.
- `attachment_mode` - (Optional, Forces new resource, String) The attachment mode of the volume. Allowable values are: **single**, **multiple**. Conflicts with `volume`.
- `resource_group` - (Optional, Forces new resource, String) The resource group ID for the new volume. Conflicts with `volume`.
- `user_tags` - (Optional, Set of String) The user tags to associate with the volume. Conflicts with `volume`.
- `allowed_use` - (Optional, Forces new resource, List) The usage constraints to be matched against requested instance or bare metal server properties to determine compatibility. Conflicts with `volume`.

  Nested scheme for `allowed_use`:
  - `api_version` - (Optional, String) The API version with which to evaluate the expressions.
  - `bare_metal_server` - (Optional, String) The expression that must be satisfied by the properties of a bare metal server provisioned using this volume.
  - `instance` - (Optional, String) The expression that must be satisfied by the properties of a virtual server instance provisioned using this volume.
- `source_snapshot` - (Optional, Forces new resource, String) The unique identifier of the snapshot to use as the source for the new volume. At least one of `volume`, `capacity`, or `source_snapshot` must be specified. Conflicts with `volume`.
- `delete_volume_on_bare_metal_server_delete` - (Optional, Boolean) If set to `true`, deleting the bare metal server also deletes the attached volume.
- `delete_volume_on_attachment_delete` - (Optional, Boolean) If set to `true`, deleting this attachment also deletes the volume. Default is `true`.
- `protocol` - (Optional, String) The protocol to use for this volume attachment. Allowable values are: **nvme_tcp**.
- `volume_name` - (Optional, String) The name of the new volume to create inline with the attachment.

## Attribute Reference

In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of the bare metal server volume attachment, in the format `<bare_metal_server_id>/<volume_attachment_id>`.
- `volume_attachment_id` - (String) The unique identifier for this bare metal server volume attachment.
- `href` - (String) The URL for this bare metal server volume attachment.
- `status` - (String) The status of this volume attachment. Allowable values are: **attaching**, **available**, **detaching**, **unusable**.
- `status_reason` - (List) The reasons for the current status (if any).

  Nested scheme for `status_reason`:
  - `code` - (String) The status reason code.
  - `message` - (String) An explanation of the status reason.
  - `more_info` - (String) A link to documentation about this status reason.
- `type` - (String) The type of volume attachment.
- `created_at` - (String) The date and time that the volume attachment was created.
- `device` - (String) A unique identifier for the device which is exposed to the bare metal server operating system.
- `nvme_qualified_name` - (String) The NVMe Qualified Name (NQN) of the subsystem. This unique identifier is used by the bare metal server to establish a connection to the volume over the nvme-tcp protocol.
- `ips` - (List of String) The IP addresses for connecting to the volume using nvme_tcp.
- `volume_crn` - (String) The CRN of the attached volume.
- `volume_href` - (String) The URL of the attached volume.
- `volume_deleted` - (String) Link to documentation about deleted resources, if the volume has been deleted.
- `source_snapshot_crn` - (String) The CRN of the source snapshot for the attached volume.

## Import

The `ibm_is_bare_metal_server_volume_attachment` resource can be imported by using the bare metal server ID and the volume attachment ID.

**Syntax**

```sh
terraform import ibm_is_bare_metal_server_volume_attachment.<name> <bare_metal_server_id>/<volume_attachment_id>
```

**Example**

```sh
terraform import ibm_is_bare_metal_server_volume_attachment.example 0716-7be49970-816a-912e-3aca-14c7fab5fda2/0716-af1a75e0-4f1f-4862-9a11-5a4c1a0b3e1f
```
