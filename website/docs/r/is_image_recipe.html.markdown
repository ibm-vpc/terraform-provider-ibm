---
layout: "ibm"
page_title: "IBM : ibm_is_image_recipe"
description: |-
  Manages ImageRecipe.
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_recipe

Create, update, and delete ImageRecipes with this resource.

## Example Usage

```hcl
resource "ibm_is_image_recipe" "is_image_recipe_instance" {
  base_image {
		crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		href = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		id = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  }
  deprecation_at = "2026-01-02T03:04:05.006Z"
  image_components {
		crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		id = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		name = "my-image-component"
		resource_type = "image_component"
  }
  name = "my-image-recipe"
  obsolescence_at = "2026-01-02T03:04:05.006Z"
  resource_group {
		id = "fee82deba12e4c0fb69c3b09d1f12345"
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `base_image` - (Required, List) The base image that is extended when building an image from this image recipe.
Nested schema for **base_image**:
	* `crn` - (Required, String) The CRN for this image.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Computed, String) A link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this image.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this image.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this image. The name is unique across all images in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `remote` - (Optional, List) If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.
	Nested schema for **remote**:
		* `account` - (Optional, List) If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.
		Nested schema for **account**:
			* `id` - (Computed, String) The unique identifier for this account.
			  * Constraints: Length must be `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
			* `resource_type` - (Computed, String) The resource type.
			  * Constraints: Allowable values are: `account`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `region` - (Optional, List) If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.
		Nested schema for **region**:
			* `href` - (Computed, String) The URL for this region.
			  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
			* `name` - (Computed, String) The globally unique name for this region.
			  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `image`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `deprecation_at` - (Optional, String) The deprecation date and time (UTC) for this image recipe.If absent, no deprecation date and time has been set.
* `image_components` - (Required, List) The image components included in this image recipe.
  * Constraints: The maximum length is `20` items. The minimum length is `1` item.
Nested schema for **image_components**:
	* `crn` - (Required, String) The CRN for this component.
	  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) A link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this component.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this component.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The unique name for this image component.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: `image_component`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `name` - (Optional, String) The unique name for this image recipe.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `obsolescence_at` - (Optional, String) The obsolescence date and time (UTC) for this image recipe.If absent, no obsolescence date and time has been set.
* `resource_group` - (Optional, List) The resource group for this image recipe.
Nested schema for **resource_group**:
	* `href` - (Computed, String) The URL for this resource group.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this resource group.
	  * Constraints: Length must be `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
	* `name` - (Computed, String) The name for this resource group.
	  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the ImageRecipe.
* `created_at` - (String) The date and time that the image recipe was created.
* `crn` - (String) The CRN for this image recipe.
  * Constraints: The maximum length is `512` characters. The minimum length is `17` characters. The value must match regular expression `/^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\\.\/]*$/`.
* `href` - (String) The URL for this image recipe.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `image_recipe`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `status` - (String) The status of an image recipe.- `available`: The recipe can be used to create images.- `degraded`:  The recipe is degraded, and cannot be used.- `deleting`: The recipe is being deleted, and can no longer be used.- `deprecated`: The recipe is deprecated, but can still be used.- `obsolete`: The recipe can no longer be used to create images.- `pending`: The recipe is being created, and is not ready to be used.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
  * Constraints: Allowable values are: `available`, `degraded`, `deleting`, `deprecated`, `obsolete`, `pending`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `status_reasons` - (List) The reasons for the current status (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **status_reasons**:
	* `code` - (String) A reason code for this status:- `base_image_unusable`: The base image for this recipe is in an `obsolete` state or  has been deleted. Image creation from this recipe will fail until the recipe no  longer has an unusable base image.- `image_component_unusable`: One or more image components referenced by this image  recipe have a `status` of `degraded`, `obsolete`, or have been deleted. Image  creation from this recipe will fail until the recipe no longer has unusable image  components.- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `base_image_unusable`, `image_component_unusable`, `internal_error`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `message` - (String) An explanation of the status reason.
	  * Constraints: The maximum length is `512` characters. The minimum length is `1` character. The value must match regular expression `/^[ -~\\n\\r\\t]*$/`.
	* `more_info` - (String) A link to documentation about this status reason.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `superseded_at` - (String) The date and time that this image recipe was superseded by a newer version.If absent, the recipe has not been superseded by a newer version.
* `updated_at` - (String) The date and time that the image recipe was last updated.
* `version` - (String) A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.
  * Constraints: The maximum length is `64` characters. The minimum length is `8` characters. The value must match regular expression `/^[0-9a-f]+$/`.
* `visibility` - (String) The visibility of this image recipe.- `private`: Visible only to this account- `public`: Visible to all accounts.
  * Constraints: Allowable values are: `private`, `public`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.


## Import

You can import the `ibm_is_image_recipe` resource by using `id`. The unique identifier for this image recipe.

# Syntax
<pre>
$ terraform import ibm_is_image_recipe.is_image_recipe &lt;id&gt;
</pre>

# Example
```
$ terraform import ibm_is_image_recipe.is_image_recipe r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8
```
