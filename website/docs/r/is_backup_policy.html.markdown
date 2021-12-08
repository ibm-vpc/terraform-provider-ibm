---
layout: "ibm"
page_title: "IBM : is_backup_policy"
description: |-
  Manages BackupPolicy.
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_backup_policy

Provides a resource for BackupPolicy. This allows BackupPolicy to be created, updated and deleted.

## Example Usage

```hcl
resource "is_backup_policy" "is_backup_policy" {
  name = "my-backup-policy"
}
```

## Argument Reference

The following arguments are supported:

* `match_resource_types` - (Optional, List) A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.
  * Constraints: The minimum length is `1` item.
* `match_user_tags` - (Optional, List) The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.
  * Constraints: The minimum length is `1` item.
* `name` - (Optional, string) The user-defined name for this backup policy. Names must be unique within the region this backup policy resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
* `plans` - (Optional, List) The prototype objects for backup plans to be created for this backup policy.
  * Constraints: The maximum length is `1` item. The minimum length is `1` item.
  * `attach_user_tags` - (Optional, []interface{}) User tags to attach to each resource created by this plan. If unspecified, no user tags will be attached.
  * `copy_user_tags` - (Optional, bool) Indicates whether to copy the source's user tags to the created resource.
    * Constraints: The default value is `true`.
  * `cron_spec` - (Required, string) The cron specification for the backup schedule.
    * Constraints: The maximum length is `63` characters. The minimum length is `9` characters. The value must match regular expression `/^((((\\d+,)+\\d+|([\\d\\*]+(\/|-)\\d+)|\\d+|\\*) ?){5,7})$/`
  * `deletion_trigger` - (Optional, BackupPolicyPlanDeletionTriggerPrototype) 
  * `name` - (Optional, string) The user-defined name for this backup policy plan. Names must be unique within the backup policy this plan resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
    * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
* `resource_group` - (Optional, List) The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.
  * `id` - (Optional, string) The unique identifier for this resource group.
    * Constraints: The value must match regular expression `/^[0-9a-f]{32}$/`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the BackupPolicy.
* `created_at` - The date and time that the backup policy was created.
* `crn` - The CRN for this backup policy.
* `href` - The URL for this backup policy.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
* `lifecycle_state` - The lifecycle state of the backup policy.
  * Constraints: Allowable values are: deleting, failed, pending, stable, updating, waiting, suspended
* `resource_type` - The type of resource referenced.
  * Constraints: Allowable values are: backup_policy

## Import

You can import the `is_backup_policy` resource by using `id`. The unique identifier for this backup policy.

```
$ terraform import is_backup_policy.is_backup_policy 0fe9e5d8-0a4d-4818-96ec-e99708644a58
```
