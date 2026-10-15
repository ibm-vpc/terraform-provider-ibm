---
layout: "ibm"
page_title: "IBM : ibm_is_image_components"
description: |-
  Get information about ImageComponentCollection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_components

Provides a read-only data source to retrieve information about an ImageComponentCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_image_components" "is_image_components" {
	is_image_components_id = "is_image_components_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `is_image_components_id` - (Required, Forces new resource, String) The image identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the ImageComponentCollection.
* `image_components` - (List) A page of image components.
Nested schema for **image_components**:
	* `allowed_use` - (List)
	Nested schema for **allowed_use**:
		* `api_version` - (String) The API version with which to evaluate the expressions.
		  * Constraints: Length must be `10` characters. The value must match regular expression `/^\\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])$/`.
		* `image` - (String) The expression that must be satisfied by the properties of an image that uses this image component.The expression follows [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md), but does not support built-in functions and macros. In addition, the following variables are supported, corresponding to `Image` properties:- `operating_system.architecture` (string): The operating system supported system  architecture- `operating_system.family` (string): The operating system family- `operating_system.version` (string): The operating system version.
		  * Constraints: The maximum length is `1000` characters. The minimum length is `1` character. The value must match regular expression `/^([a-zA-Z_][a-zA-Z0-9_]*|[-+*\/%]|&&|\\|\\||!|==|!=|<|<=|>|>=|~|\\bin\\b|\\(|\\)|\\[|\\]|,|\\.|"|'|"|'|\\s+|\\d+)+$/`.
	* `created_at` - (String) The date and time that the image component was created.
	* `crn` - (String) The CRN for this component.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `dependencies` - (List) The dependencies for this image component. An image recipe that uses this component must contain its required dependencies.
	  * Constraints: The maximum length is `19` items. The minimum length is `0` items.
	Nested schema for **dependencies**:
		* `crn` - (String) The CRN for this component.
		  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
		* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			* `more_info` - (String) A link to documentation about deleted resources.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `href` - (String) The URL for this component.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this component.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
		* `name` - (String) The unique name for this image component.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `image_component`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `deprecation_at` - (String) The deprecation date and time for this image component.If absent, no deprecation date and time has been set.
	* `href` - (String) The URL for this component.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this component.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The unique name for this image component.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `obsolescence_at` - (String) The obsolescence date and time for this component.If absent, no obsolescence date and time has been set.
	* `resource_group` - (List) The resource group for this component.
	Nested schema for **resource_group**:
		* `href` - (String) The URL for this resource group.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
		* `id` - (String) The unique identifier for this resource group.
		  * Constraints: Length must be `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
		* `name` - (String) The name for this resource group.
		  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `image_component`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `software` - (List)
	Nested schema for **software**:
		* `name` - (String) The name of the software package.
		  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `software_package`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `version` - (String) The version of the software package.
		  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[A-Za-z0-9][A-Za-z0-9._+\\-:~]*$/`.
	* `status` - (String) The status of this image component- `available`: image component is ready to be used.- `degraded`: image component is degraded, and cannot be used.- `deleting`: image component is being deleted.- `deprecated`: image component is deprecated, but can still be used.- `obsolete`: image component is obsolete, and can no longer be used.- `pending`: image component is being created and is not ready to be used.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `available`, `degraded`, `deleting`, `deprecated`, `obsolete`, `pending`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `status_reasons` - (List) The reasons for the current status (if any).
	  * Constraints: The minimum length is `0` items.
	Nested schema for **status_reasons**:
		* `code` - (String) A reason code for this status:- `dependent_image_component_unusable`: One or more dependent image components  referenced by this image component have a status of `degraded`, `obsolete`, or  `deleting` or have been deleted. Image recipe creation with this image component  will fail until the image component no longer has unusable dependent image components.- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
		  * Constraints: Allowable values are: `dependent_image_component_unusable`, `internal_error`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `message` - (String) An explanation of the status reason.
		  * Constraints: The maximum length is `512` characters. The minimum length is `1` character. The value must match regular expression `/^[ -~\\n\\r\\t]*$/`.
		* `more_info` - (String) A link to documentation about this status reason.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `visibility` - (String) The visibility of this image component.- `private`: Visible only to this account- `public`: Visible to all accounts.
	  * Constraints: Allowable values are: `private`, `public`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

