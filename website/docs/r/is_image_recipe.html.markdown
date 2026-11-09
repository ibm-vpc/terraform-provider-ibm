---
layout: "ibm"
page_title: "IBM : ibm_is_image_recipe"
description: |-
  Manages an IBM Cloud VPC Image Recipe.
subcategory: "VPC infrastructure"
---

# ibm_is_image_recipe

Create, update, or delete an IBM Cloud VPC Image Recipe. For more information, see [IBM Cloud Docs: Image Builder](https://cloud.ibm.com/docs/vpc).

## Example Usage

### Basic Usage

```terraform
resource "ibm_is_image_recipe" "example" {
  base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
}
```

### Full Example

```terraform
resource "ibm_is_image_recipe" "example" {
  name             = "my-image-recipe"
  base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  image_components = [
    "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8",
    "r006-83c38c6d-a5c1-59cc-c065-6cfdd802cdc9"
  ]
  resource_group   = "fee82deba12e4c0fb69c3b09d1f12345"
  deprecation_at   = "2029-06-01T00:00:00.000Z"
  obsolescence_at  = "2029-12-01T00:00:00.000Z"
}
```

## Argument Reference

You can specify the following arguments for this resource:

* `base_image` - (Required, Forces new resource, String) The base image from which the image will be built. The base image must meet the usage constraints of the image components specified in the image recipe, and the base image visibility must be `public`.
* `image_components` - (Required, List of Strings) The image components to include in this image recipe (1 to 20 items).
* `name` - (Optional, String) The unique name for this image recipe. Must be 1 to 63 lowercase alphanumeric characters or hyphens.
* `resource_group` - (Optional, Forces new resource, String) The resource group ID for this image recipe. Uses the account default if omitted.
* `deprecation_at` - (Optional, String) The deprecation date and time (UTC, RFC 3339 format) to set for this image recipe. The date and time must not be in the past, and must be earlier than `obsolescence_at` if specified. Specify `null` to remove an existing deprecation date and time.
* `obsolescence_at` - (Optional, String) The obsolescence date and time (UTC, RFC 3339 format) to set for this image recipe. The date and time must not be in the past, and must be later than `deprecation_at` if specified. Specify `null` to remove an existing obsolescence date and time.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes:

* `id` - (String) The unique identifier of this image recipe.
* `created_at` - (String) The date and time that the image recipe was created.
* `crn` - (String) The CRN for this image recipe.
* `href` - (String) The URL for this image recipe.
* `resource_type` - (String) The resource type (`image_recipe`).
* `status` - (String) The status of the image recipe (`available`, `degraded`, `deleting`, `deprecated`, `obsolete`, `pending`).
* `status_reasons` - (List) The reasons for the current status (if any).
  Nested schema for `status_reasons`:
  * `code` - (String) A reason code for this status (`base_image_unusable`, `image_component_unusable`, `internal_error`).
  * `message` - (String) An explanation of the status reason.
  * `more_info` - (String) A link to documentation about this status reason.
* `superseded_at` - (String) The date and time that this image recipe was superseded by a newer version (if applicable).
* `updated_at` - (String) The date and time that the image recipe was last updated.
* `version` - (String) A generated content fingerprint for this version of the image recipe.
* `visibility` - (String) The visibility of this image recipe (`private`, `public`).

## Import

You can import the `ibm_is_image_recipe` resource by using the image recipe ID:

```bash
$ terraform import ibm_is_image_recipe.example r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8
```
