---
layout: "ibm"
page_title: "IBM : ibm_is_backup_policy_job"
description: |-
  Get information about BackupPolicyJob
subcategory: "Virtual Private Cloud API"
---

# ibm_is_backup_policy_job

Provides a read-only data source for BackupPolicyJob. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_backup_policy_job" "is_backup_policy_job" {
	backup_policy_id = "backup_policy_id"
	id = "id"
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

* `backup_policy_id` - (Required, Forces new resource, String) The backup policy identifier.
* `id` - (Required, Forces new resource, String) The backup policy job identifier.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the BackupPolicyJob.
* `backup_info` - (Required, List) The snapshot created by this backup policy job (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).
Nested scheme for **backup_info**:
	* `crn` - (Required, String) The CRN for this snapshot.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
	Nested scheme for **deleted**:
		* `deleted_at` - (Required, String) The date and time that the reference resource was deleted.
		* `final_name` - (Required, String) The user-defined name the referenced resource had at the time it was deleted.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `href` - (Required, String) The URL for this snapshot.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `id` - (Required, String) The unique identifier for this snapshot.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`
	* `name` - (Required, String) The user-defined name for this snapshot.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: snapshot

* `backup_resource_type` - (Required, String) A resource type this backup policy job applied to.

* `completed_at` - (Required, String) The date and time that the backup policy job was completed.

* `href` - (Required, String) The URL for this backup policy job.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

* `plan_id` - (Required, String) The unique identifier for this backup policy plan.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`

* `resource_type` - (Required, String) The type of resource referenced.
  * Constraints: Allowable values are: backup_policy_job

* `source_volume` - (Required, List) The source volume this backup was created from (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).
Nested scheme for **source_volume**:
	* `crn` - (Required, String) The CRN for this volume.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
	Nested scheme for **deleted**:
		* `deleted_at` - (Required, String) The date and time that the reference resource was deleted.
		* `final_name` - (Required, String) The user-defined name the referenced resource had at the time it was deleted.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `href` - (Required, String) The URL for this volume.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `id` - (Required, String) The unique identifier for this volume.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`
	* `name` - (Required, String) The unique user-defined name for this volume.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: volume

* `started_at` - (Required, String) The date and time that the backup policy job was started.

* `status` - (Required, String) The status of the backup policy job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the volume on which the unexpected property value was encountered.
  * Constraints: Allowable values are: running, failed, completed

* `status_reasons` - (Required, List) The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.
Nested scheme for **status_reasons**:
	* `code` - (Required, String) A snake case string succinctly identifying the status reason.
	  * Constraints: Allowable values are: backup_failedThe value must match regular expression `/^[a-z]+(_[a-z]+)*$/`
	* `message` - (Required, String) An explanation of the status reason.
	* `more_info` - (Optional, String) Link to documentation about this status reason.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

