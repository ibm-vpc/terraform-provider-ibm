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
data "is_share" "is_share" {
	id = "id"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, string) The file share identifier.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the Share.
* `created_at` - The date and time that the file share is created.

* `crn` - The CRN for this share.

* `encryption` - The type of encryption used for this file share.

* `encryption_key` - The key used to encrypt this file share. Nested `encryption_key` blocks have the following structure:
	* `crn` - The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.

* `href` - The URL for this share.

* `iops` - The maximum input/output operation performance bandwidth per second for the file share.

* `lifecycle_state` - The lifecycle state of the file share.

* `name` - The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.

* `profile` - The profile this file share uses. Nested `profile` blocks have the following structure:
	* `href` - The URL for this share profile.
	* `name` - The globally unique name for this share profile.
	* `resource_type` - The resource type.

* `resource_group` - The resource group for this file share. Nested `resource_group` blocks have the following structure:
	* `href` - The URL for this resource group.
	* `id` - The unique identifier for this resource group.
	* `name` - The user-defined name for this resource group.

* `resource_type` - The type of resource referenced.

* `size` - The size of the file share rounded up to the next gigabyte.

* `targets` - Mount targets for the file share. Nested `targets` blocks have the following structure:
	* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
		* `more_info` - Link to documentation about deleted resources.
	* `href` - The URL for this share target.
	* `id` - The unique identifier for this share target.
	* `name` - The user-defined name for this share target.
	* `resource_type` - The type of resource referenced.

* `zone` - The zone this file share will reside in. Nested `zone` blocks have the following structure:
	* `href` - The URL for this zone.
	* `name` - The globally unique name for this zone.

