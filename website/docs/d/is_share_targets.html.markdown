---
layout: "ibm"
page_title: "IBM : is_share_targets"
sidebar_current: "docs-ibm-datasource-is-share-targets"
description: |-
  Get information about ShareTargetCollection
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_share_targets

Provides a read-only data source for ShareTargetCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "is_share_targets" "is_share_targets" {
	share_id = "share_id"
}
```

## Argument Reference

The following arguments are supported:

* `share_id` - (Required, string) The file share identifier.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the ShareTargetCollection.
* `targets` - Collection of share targets. Nested `targets` blocks have the following structure:
	* `created_at` - The date and time that the share target was created.
	* `href` - The URL for this share target.
	* `id` - The unique identifier for this share target.
	* `lifecycle_state` - The lifecycle state of the mount target.
	* `mount_path` - The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.
	* `name` - The user-defined name for this share target.
	* `resource_type` - The type of resource referenced.
	* `security_groups` - Collection of security groups. Nested `security_groups` blocks have the following structure:
		* `crn` - The security group's CRN.
		* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
			* `more_info` - Link to documentation about deleted resources.
		* `href` - The security group's canonical URL.
		* `id` - The unique identifier for this security group.
		* `name` - The user-defined name for this security group. Names must be unique within the VPC the security group resides in.
	* `subnet` - The subnet associated with this file share target. Nested `subnet` blocks have the following structure:
		* `crn` - The CRN for this subnet.
		* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
			* `more_info` - Link to documentation about deleted resources.
		* `href` - The URL for this subnet.
		* `id` - The unique identifier for this subnet.
		* `name` - The user-defined name for this subnet.
		* `resource_type` - The resource type.
	* `vpc` - The VPC to which this share target is allowing to mount the file share. Nested `vpc` blocks have the following structure:
		* `crn` - The CRN for this VPC.
		* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
			* `more_info` - Link to documentation about deleted resources.
		* `href` - The URL for this VPC.
		* `id` - The unique identifier for this VPC.
		* `name` - The unique user-defined name for this VPC.
		* `resource_type` - The resource type.

