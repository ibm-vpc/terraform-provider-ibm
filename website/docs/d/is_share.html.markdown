---
layout: "ibm"
page_title: "IBM : is_share"
sidebar_current: "docs-ibm-datasource-is-share"
description: |-
  Get information about Share
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_share

Provides a read-only data source for Share. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
resource "ibm_is_vpc" "vpc" {
  name = "my-vpc"
}
resource "ibm_is_share" "is_share" {
  name = "my-share"
  size = 200
  profile = "tier-3iops"
  zone = "us-south-2"
}

data "ibm_is_share" "is_share" {
	share = ibm_is_share.is_share.id
}

data "ibm_is_share" "is_share_with_name" {
	name = ibm_is_share.is_share.name
}
```

## Argument Reference

The following arguments are supported:

- `share` - (Optional, String) The file share identifier.
- `name` - (Optional, String) The file share name
**Note** One of the aurgument is mandatory

## Attribute Reference

The following attributes are exported:

- `created_at` - The date and time that the file share is created.
- `crn` - The CRN for this share.
- `encryption` - The type of encryption used for this file share.
- `encryption_key` - The CRN of the key used to encrypt this file share. Nested `encryption_key` blocks have the following structure:
- `href` - The URL for this share.
- `iops` - The maximum input/output operation performance bandwidth per second for the file share.
- `lifecycle_state` - The lifecycle state of the file share.
- `name` - The unique user-defined name for this file share.
- `profile` - The name of the profile this file share uses.
- `resource_group` - The ID of the resource group for this file share.
- `resource_type` - The type of resource referenced.
- `size` - The size of the file share rounded up to the next gigabyte.
- `share_targets` - Mount targets for the file share. Nested `targets` blocks have the following structure:
	- `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
		- `more_info` - Link to documentation about deleted resources.
	- `href` - The URL for this share target.
	- `id` - The unique identifier for this share target.
	- `name` - The user-defined name for this share target.
	- `resource_type` - The type of resource referenced.
- `zone` - The name of the zone this file share will reside in.
- `access_tags`  - (String) Access management tags associated to the share.
- `tags`  - (String) User tags associated for to the share.

