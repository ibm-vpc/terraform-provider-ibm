---
layout: "ibm"
page_title: "IBM : is_share"
sidebar_current: "docs-ibm-resource-is-share"
description: |-
  Manages Share.
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_share

Provides a resource for Share. This allows Share to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_is_share" "is_share" {
  name = "my-share"
  size = 200
  profile = "tier-3iops"
  zone = "us-south-2"
}
```

## Argument Reference

The following arguments are supported:

- `encryption_key` - (Optional, String) The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.
- `initial_owner_gid` - (Optional, int) The initial group identifier for the file share.
- `initial_owner_uid` - (Optional, int) The initial user identifier for the file share.
- `iops` - (Optional, int) The maximum input/output operation performance bandwidth per second for the file share.
- `name` - (Required, string) The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.
- `profile` - (Required, string) The globally unique name for this share profile.
- `resource_group` - (Optional, string) The unique identifier for this resource group. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.
- `size` - (Required, int) The size of the file share rounded up to the next gigabyte.
- `share_target_prototype` - (Optional, List) Share targets for the file share.
  - `name` - (Optional, string) The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
  - `vpc` - (Required, string) The VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.
- `zone` - (Required, string) The globally unique name for this zone.
- `access_tags`  - (Optional, List of Strings) The list of access management tags to attach to the share. **Note** For more information, about creating access tags, see [working with tags](https://cloud.ibm.com/docs/account?topic=account-tag).
- `tags`  - (Optional, List of Strings) The list of user tags to attach to the share.

## Attribute Reference

The following attributes are exported:

- `id` - The unique identifier of the Share.
- `created_at` - The date and time that the file share is created.
- `crn` - The CRN for this share.
- `encryption` - The type of encryption used for this file share.
- `href` - The URL for this share.
- `lifecycle_state` - The lifecycle state of the file share.
- `resource_type` - The type of resource referenced.
- `encryption_key` - (Optional, String) The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.
- `iops` - The maximum input/output operation performance bandwidth per second for the file share.
- `resource_group` - (Optional, string) The unique identifier for this resource group. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.
- `share_targets` - Mount targets for the file share. Nested `targets` blocks have the following structure:
	- `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
		- `more_info` - Link to documentation about deleted resources.
	- `href` - The URL for this share target.
	- `id` - The unique identifier for this share target.
	- `name` - The user-defined name for this share target.
	- `resource_type` - The type of resource referenced.
- `access_tags`  - (String) Access management tags associated to the share.
- `tags`  - (String) User tags associated for to the share.