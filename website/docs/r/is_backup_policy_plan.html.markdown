---
layout: "ibm"
page_title: "IBM : ibm_is_backup_policy_plan"
description: |-
  Manages BackupPolicyPlan.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_backup_policy_plan

Provides a resource for BackupPolicyPlan. This allows BackupPolicyPlan to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_is_backup_policy_plan" "is_backup_policy_plan" {
  backup_policy_id = "backup_policy_id"
  cron_spec = "*/5 1,2,3 * * *"
  name = "my-backup-policy"
}
```

## Argument Reference

Review the argument reference that you can specify for your resource.

* `active` - (Optional, Boolean) Indicates whether the plan is active.
* `attach_user_tags` - (Optional, List) User tags to attach to each resource created by this plan. If unspecified, no user tags will be attached.
* `backup_policy_id` - (Required, Forces new resource, String) The backup policy identifier.
* `copy_user_tags` - (Optional, Boolean) Indicates whether to copy the source's user tags to the created resource.
  * Constraints: The default value is `true`.
* `cron_spec` - (Required, String) The cron specification for the backup schedule.
  * Constraints: The maximum length is `63` characters. The minimum length is `9` characters. The value must match regular expression `/^((((\\d+,)+\\d+|([\\d\\*]+(\/|-)\\d+)|\\d+|\\*) ?){5,7})$/`
* `deletion_trigger` - (Optional, List) 
Nested scheme for **deletion_trigger**:
	* `delete_after` - (Optional, Integer) The number of days to keep the backup.
	* `delete_after_backup_count` - (Optional, Integer) The number of latest backup to be retained.
* `name` - (Optional, String) The user-defined name for this backup policy plan. Names must be unique within the backup policy this plan resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the BackupPolicyPlan.
* `created_at` - (Required, String) The date and time that the backup policy plan was created.
* `href` - (Required, String) The URL for this backup policy plan.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
* `lifecycle_state` - (Required, String) The lifecycle state of this backup policy plan.
  * Constraints: Allowable values are: deleting, failed, pending, stable, updating, waiting, suspended
* `resource_type` - (Required, String) The type of resource referenced.
  * Constraints: Allowable values are: backup_policy_plan

## Import

You can import the `ibm_is_backup_policy_plan` resource by using `id`.
The `id` property can be formed from `backup_policy_id`, and `id` in the following format:

```
<backup_policy_id>/<id>
```
* `backup_policy_id`: A string. The backup policy identifier.
* `id`: A string. The backup policy plan identifier.

# Syntax
```
$ terraform import ibm_is_backup_policy_plan.is_backup_policy_plan <backup_policy_id>/<id>
```
