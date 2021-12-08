---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_backup_policy_plan"
description: |-
  Get information about backup policy plan.
---

# ibm\_is_backup_policy_plan

Provides a read-only data source for BackupPolicyPlan. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_backup_policy_plan" "is_backup_policy_plan" {
	backup_policy_id = "backup_policy_id"
	identifier = "id"
}
```

## Argument Reference
Review the argument references that you can specify for your data source. 

- `backup_policy_id` - (Required, string) The backup policy identifier.
- `identifier` - (Optional, string) The backup policy plan identifier.
- `name` - (Optional, string) The unique user-defined name for backup policy..

## Attribute Reference
In addition to all argument reference list, you can access the following attribute references after your data source is created.

- `id` - The unique identifier of the BackupPolicyPlan.
- `active` - Indicates whether the plan is active.
- `attach_user_tags` - User tags to attach to each resource created by this plan.
- `copy_user_tags` - Indicates whether to copy the source's user tags to the created resource.
- `created_at` - The date and time that the backup policy plan was created.
- `cron_spec` - The cron specification for the backup schedule.
- `deletion_trigger`  Nested `deletion_trigger` blocks have the following structure:
	- `delete_after` - The number of days to keep the backup.
- `href` - The URL for this backup policy plan.
- `lifecycle_state` - The lifecycle state of this backup policy plan.
- `resource_type` - The type of resource referenced.

