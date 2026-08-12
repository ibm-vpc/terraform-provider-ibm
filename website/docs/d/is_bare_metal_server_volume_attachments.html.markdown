---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : bare_metal_server_volume_attachments"
description: |-
  Get information about IBM bare metal server volume attachments.
---

# ibm\_is_bare_metal_server_volume_attachments

Import the details of all existing IBM Cloud Bare Metal Server volume attachments as a read-only data source. For more information, about managing VPC Bare Metal Server, see [About Bare Metal Servers for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-bare-metal-servers).

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
data "ibm_is_bare_metal_server_volume_attachments" "example" {
  bare_metal_server = ibm_is_bare_metal_server.example.id
}
```

## Argument Reference

Review the argument references that you can specify for your data source.

- `bare_metal_server` - (Required, String) The unique identifier of the bare metal server.

## Attribute Reference

In addition to the argument reference list, you can access the following attribute references after your data source is created.

- `volume_attachments` - (List) Collection of volume attachments for the bare metal server.

  Nested scheme for `volume_attachments`:
  - `bandwidth` - (Integer) The maximum bandwidth (in megabits per second) for the volume when attached to this bare metal server.
  - `created_at` - (String) The date and time that the volume attachment was created.
  - `delete_volume_on_bare_metal_server_delete` - (Boolean) Indicates whether deleting the bare metal server will also delete the attached volume. This property must be false if the volume's `attachment_mode` is `multiple`.
  - `device` - (String) A unique identifier for the device which is exposed to the bare metal server operating system. This property may be absent if the status of the volume attachment is not available.
  - `href` - (String) The URL for this bare metal server volume attachment.
  - `ips` - (List of String) The IP addresses for connecting to the volume using nvme_tcp.
  - `nvme_qualified_name` - (String) The NVMe Qualified Name (NQN) of the subsystem. This unique identifier is used by the bare metal server to establish a connection to the volume over the nvme-tcp protocol. The NQN must be used when configuring the NVMe initiator on the bare metal server to access the attached volume.
  - `protocol` - (String) The protocol used for this volume attachment: `nvme_tcp` (Non-Volatile Memory Express (NVMe) over TCP/IP, which allows bare metal servers to connect to volumes over the network using the NVMe protocol).
  - `status` - (String) The status of this volume attachment. Allowable values are: **attaching** (volume attachment is being initialized and not yet usable), **available** (volume attachment is usable, connection to the volume can be established from the server's operating system), **detaching** (volume attachment is being removed), **unusable** (volume attachment is unusable due to the underlying volume state).
  - `status_reason` - (List) The reasons for the current status (if any).

    Nested scheme for `status_reason`:
    - `code` - (String) The status reason code: `volume_encryption_key_deleted` (The key associated with the data volume attached to the bare metal server is deleted).
    - `message` - (String) An explanation of the status reason.
    - `more_info` - (String) A link to documentation about this status reason.
  - `type` - (String) The type of volume attachment.
  - `volume` - (List) The attached volume.

    Nested scheme for `volume`:
    - `volume` - (String) The unique identifier for this volume.
    - `volume_name` - (String) The name of the attached volume.
    - `volume_crn` - (String) The CRN of the attached volume.
    - `volume_href` - (String) The URL of the attached volume.
  - `volume_attachment_id` - (String) The unique identifier for this bare metal server volume attachment.
  - `name` - (String) The name for this volume attachment. The name is unique across all volume attachments on the bare metal server.
