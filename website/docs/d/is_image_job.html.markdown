---
layout: "ibm"
page_title: "IBM : ibm_is_image_job"
description: |-
  Get information about is_image_job
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_job

Provides a read-only data source to retrieve information about an is_image_job. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_image_job" "is_image_job" {
	is_image_job_id = "is_image_job_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `is_image_job_id` - (Required, Forces new resource, String) The image job identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the is_image_job.
* `completed_at` - (String) The date and time that the image job was completed.If absent, the image job has not yet completed.
* `created_at` - (String) The date and time that the image job was created.
* `delete_image_job_on_image_delete` - (Boolean) Indicates whether deleting the image will also delete the image job that created it.
* `href` - (String) The URL for this image job.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `image` - (List) 
Nested schema for **image**:
	* `crn` - (String) The CRN for this image.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) A link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this image.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this image.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this image. The name is unique across all images in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `remote` - (List) If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.
	Nested schema for **remote**:
		* `account` - (List) If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.
		Nested schema for **account**:
			* `id` - (String) The unique identifier for this account.
			  * Constraints: Length must be `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
			* `resource_type` - (String) The resource type.
			  * Constraints: Allowable values are: `account`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `region` - (List) If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.
		Nested schema for **region**:
			* `href` - (String) The URL for this region.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
			* `name` - (String) The globally unique name for this region.
			  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `image`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `image_recipe` - (List) The image recipe and version used by this image job.
Nested schema for **image_recipe**:
	* `crn` - (String) The CRN for this image recipe.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) A link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The URL for this image recipe.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this image recipe.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The unique name for this image recipe.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `image_recipe`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `version` - (String) A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.
	  * Constraints: The maximum length is `64` characters. The minimum length is `8` characters. The value must match regular expression `/^[0-9a-f]+$/`.
* `image_template` - (List) The template for creating an image from the image job.
Nested schema for **image_template**:
	* `name` - (String) The name for this image. The name must not be used by another image in the region. Names starting with `ibm-` are reserved for system-provided images, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$/`.
	* `resource_group` - (List) The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) will be used.
	Nested schema for **resource_group**:
		* `id` - (String) The unique identifier for this resource group.
		  * Constraints: Length must be `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
* `name` - (String) The name for this image job. The name is unique across all image jobs in the region.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `image_job`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `started_at` - (String) The date and time that the image job started running.If absent, the image job has not yet started.
* `status` - (String) The status of this image job:- `deleting`: The image job is being deleted- `failed`: The image job could not be completed successfully. The customer may  retry by creating a new image job for the same image.- `queued`: The image job has been accepted and is waiting to be scheduled- `running`: The image job is in progress- `succeeded`: The image job was completed successfullyThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `deleting`, `failed`, `queued`, `running`, `succeeded`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `status_reasons` - (List) The reasons for the current status (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **status_reasons**:
	* `code` - (String) A reason code for the image job status:The following reason codes may be set when the image job `status` is `failed`:- `failed_component_installation`: one or more image components failed to install while  building the image- `failed_component_validation`: one or more image components failed validation- `failed_image_capture`: failed to capture the image from the source volume- `failed_timeout`: the image job exceeded the configured `timeout`- `failed_worker_provisioning`: the worker used to run the image job failed to  provision- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `failed_component_installation`, `failed_component_validation`, `failed_image_capture`, `failed_timeout`, `failed_worker_provisioning`, `internal_error`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `message` - (String) An explanation of the status reason.
	  * Constraints: The maximum length is `512` characters. The minimum length is `1` character. The value must match regular expression `/^[ -~\\n\\r\\t]*$/`.
	* `more_info` - (String) A link to documentation about this status reason.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `timeout` - (Integer) The maximum time (in minutes) allowed for this image job to complete.The time begins when the job moves to a `status` of `running`. If the timeout is reached, the image and worker will be deleted.
  * Constraints: The maximum value is `480`. The minimum value is `30`.
* `type` - (String) The type of image job:- `build`: build an image from an image recipeThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `build`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `worker` - (List) The worker that executes this job.
Nested schema for **worker**:
	* `instance` - (List) The instance resource that is the worker for this job. If absent, the instance has notbeen created yet.
	Nested schema for **instance**:
		* `crn` - (String) The CRN for this virtual server instance.
		  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) A link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this virtual server instance.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^https:\/\/([^\/?#]*)([^?#]*)\/instances\/[-0-9a-z_]+$/`.
		* `id` - (String) The unique identifier for this virtual server instance.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this virtual server instance. The name is unique across all virtual server instances in the region.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `profile` - (List) The [profile](https://cloud.ibm.com/docs/vpc?topic=vpc-profiles) used to build the image.
	Nested schema for **profile**:
		* `href` - (String) The URL for this virtual server instance profile.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `name` - (String) The globally unique name for this virtual server instance profile.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `instance_profile`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `subnet` - (List)
	Nested schema for **subnet**:
		* `crn` - (String) The CRN for this subnet.
		  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) A link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this subnet.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this subnet.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The name for this subnet. The name is unique across all subnets in the VPC.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `subnet`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

