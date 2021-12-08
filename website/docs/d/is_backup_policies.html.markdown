---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_backup_policies"
description: |-
  Get information about backup policies.
---

# ibm\_is_backup_policies

Provides a read-only data source for BackupPolicyCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_backup_policies" "is_backup_policies" {
	name = "my-backup-policy"
}
```

## Attribute Reference
In addition to all argument reference list, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the BackupPolicyCollection.
* `backup_policies` - Collection of backup policies. 

Nested `backup_policies` blocks have the following structure:
- `created_at` - The date and time that the backup policy was created.
- `crn` - The CRN for this backup policy.
- `href` - The URL for this backup policy.
- `id` - The unique identifier for this backup policy.
- `lifecycle_state` - The lifecycle state of the backup policy.
- `match_resource_types` - A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.
- `match_user_tags` - The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.
- `name` - The unique user-defined name for this backup policy.
- `plans` - The plans for the backup policy. 
- `resource_type` - The type of resource referenced.
- `resource_group` - The resource group for this backup policy. 
* `total_count` - The total number of resources across all pages.


Nested `plans` blocks have the following structure:
- `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
- `href` - The URL for this backup policy plan.
- `id` - The unique identifier for this backup policy plan.
- `name` - The unique user-defined name for this backup policy plan.
- `resource_type` - The type of resource referenced.

Nested `deleted` blocks have the following structure:
- `more_info` - Link to documentation about deleted resources.
		

Nested `resource_group` blocks have the following structure:
- `href` - The URL for this resource group.
- `id` - The unique identifier for this resource group.
- `name` - The user-defined name for this resource group.



