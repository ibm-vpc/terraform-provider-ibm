---
layout: "ibm"
page_title: "IBM : is_share_target"
sidebar_current: "docs-ibm-resource-is-share-target"
description: |-
  Manages ShareTarget.
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_share_target

Provides a resource for ShareTarget. This allows ShareTarget to be created, updated and deleted.

## Example Usage

```hcl
resource "is_share_target" "is_share_target" {
  share_id = "share_id"
  vpc = { example: "object" }
  name = "my-share-target"
  subnet = {"id":"7ec86020-1c6e-4889-b3f0-a15f2e50f87e"}
}
```

## Argument Reference

The following arguments are supported:

* `share_id` - (Required, string) The file share identifier.
* `vpc` - (Required, List) The VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.
  * `id` - (Optional, string) The unique identifier for this VPC.
  * `crn` - (Optional, string) The CRN for this VPC.
  * `href` - (Optional, string) The URL for this VPC.
* `name` - (Optional, string) The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
* `subnet` - (Optional, List) The subnet associated with this file share target.Only virtual server instances in the same VPC as this subnetwill be allowed to mount the file share.In the future, this property may be required and used to assignan IP address for the file share target.
  * `id` - (Optional, string) The unique identifier for this subnet.
  * `crn` - (Optional, string) The CRN for this subnet.
  * `href` - (Optional, string) The URL for this subnet.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the ShareTarget.
* `created_at` - The date and time that the share target was created.
* `href` - The URL for this share target.
* `lifecycle_state` - The lifecycle state of the mount target.
* `mount_path` - The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.
* `resource_type` - The type of resource referenced.
* `security_groups` - Collection of security groups.
