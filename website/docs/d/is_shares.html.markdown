---
layout: "ibm"
page_title: "IBM : is_shares"
sidebar_current: "docs-ibm-datasource-is-shares"
description: |-
  Get information about ShareCollection
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_shares

Provides a read-only data source for ShareCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_shares" "is_shares" {
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional, string) The unique user-defined name for this file share to filter the collection.
- `resource_group` - (Optional, string) The unique identifier for this resource group to filter the collection.

## Attribute Reference

The following attributes are exported:

- `shares` - Collection of file shares. Nested `shares` blocks have the following structure:
	- `created_at` - The date and time that the file share is created.
	- `crn` - The CRN for this share.
	- `encryption` - The type of encryption used for this file share.
	- `encryption_key` - The CRN of the key used to encrypt this file share. Nested `encryption_key` blocks have the following structure:
	- `href` - The URL for this share.
	- `id` - The unique identifier for this file share.
	- `iops` - The maximum input/output operation performance bandwidth per second for the file share.
	- `lifecycle_state` - The lifecycle state of the file share.
	- `name` - The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.
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
- `total_count` - The total number of resources across all pages.

