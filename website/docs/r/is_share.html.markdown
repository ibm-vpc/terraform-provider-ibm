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
resource "is_share" "is_share" {
  iops = 100
  name = "my-share"
  size = 200
}
```

## Argument Reference

The following arguments are supported:

* `encryption_key` - (Optional, List) The key to use for encrypting this file share.If no encryption key is provided, the share will not be encrypted.
  * `crn` - (Optional, string) The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.
* `initial_owner` - (Optional, List) The owner assigned to the file share at creation. Subsequent changes to the ownermust be performed by a virtual server instance that has mounted the file share.
  * `gid` - (Optional, int) The initial group identifier for the file share.
  * `uid` - (Optional, int) The initial user identifier for the file share.
* `iops` - (Optional, int) The maximum input/output operation performance bandwidth per second for the file share.
* `name` - (Optional, string) The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.
* `profile` - (Optional, List) Identifies a share profile by a unique property.
  * `name` - (Optional, string) The globally unique name for this share profile.
  * `href` - (Optional, string) The URL for this share profile.
* `resource_group` - (Optional, List) The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.
  * `id` - (Optional, string) The unique identifier for this resource group.
* `size` - (Optional, int) The size of the file share rounded up to the next gigabyte.
* `targets` - (Optional, List) Share targets for the file share.
  * `name` - (Optional, string) The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
  * `subnet` - (Optional, SubnetIdentity) The subnet associated with this file share target.Only virtual server instances in the same VPC as this subnetwill be allowed to mount the file share.In the future, this property may be required and used to assignan IP address for the file share target.
  * `vpc` - (Required, ShareTargetPrototypeVpc) The VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.
* `zone` - (Optional, List) The zone this file share will reside in.
  * `name` - (Optional, string) The globally unique name for this zone.
  * `href` - (Optional, string) The URL for this zone.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the Share.
* `created_at` - The date and time that the file share is created.
* `crn` - The CRN for this share.
* `encryption` - The type of encryption used for this file share.
* `href` - The URL for this share.
* `lifecycle_state` - The lifecycle state of the file share.
* `resource_type` - The type of resource referenced.
