---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_backup_policy"
description: |-
  Get information about backup policy.
---

# ibm\_is_backup_policy

Provides a read-only data source for BackupPolicy. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_backup_policy" "is_backup_policy" {
	identifier = "id"
}
```

## Argument Reference
Review the argument references that you can specify for your data source. 

- `identifier` - (Optional, string) The backup policy identifier.
- `name` - (Optional, string) The unique user-defined name for backup policy..


## Attribute Reference
In addition to all argument reference list, you can access the following attribute references after your data source is created.

- `id` - The unique identifier of the BackupPolicy.
- `created_at` - The date and time that the backup policy was created.
- `crn` - The CRN for this backup policy.
- `href` - The URL for this backup policy.
- `lifecycle_state` - The lifecycle state of the backup policy.
- `match_resource_types` - A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.
- `match_user_tags` - The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.
- `name` - The unique user-defined name for this backup policy.
- `plans` - The plans for the backup policy.
- `resource_group` - The resource group for this backup policy.
- `resource_type` - The type of resource referenced.

Nested `plans` blocks have the following structure:
- `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. 
- `href` - The URL for this backup policy plan.
- `id` - The unique identifier for this backup policy plan.
- `name` - The unique user-defined name for this backup policy plan.
- `resource_type` - The type of resource referenced.

Nested `deleted` blocks have the following structure:
- `more_info` - Link to documentation about deleted resources.


Nested `resource_group` blocks have the following structure:
- `id` - The unique identifier for this resource group.
- `name` - The user-defined name for this resource group.


