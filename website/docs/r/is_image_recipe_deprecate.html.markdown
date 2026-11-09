---
layout: "ibm"
page_title: "IBM : ibm_is_image_recipe_deprecate"
description: |-
  Deprecates an IBM Cloud VPC Image Recipe.
subcategory: "VPC infrastructure"
---

# ibm_is_image_recipe_deprecate

Provides support to deprecate a VPC Image Recipe. Deprecating an image recipe transitions its `status` to `deprecated` and sets `deprecation_at` to the current date and time. A deprecated image recipe can still be used to build images, but signals that it should no longer be extended.

For more information, see [IBM Cloud Docs: Image Builder](https://cloud.ibm.com/docs/vpc).

## Example Usage

```terraform
resource "ibm_is_image_recipe" "example" {
  name             = "my-image-recipe"
  base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
}

resource "ibm_is_image_recipe_deprecate" "example" {
  image_recipe = ibm_is_image_recipe.example.id
}
```

## Argument Reference

You can specify the following arguments for this resource:

* `image_recipe` - (Required, Forces new resource, String) The unique identifier of the image recipe to deprecate.

## Attribute Reference

After your resource is created, you can read values from the following attributes:

* `id` - (String) The unique identifier of the image recipe.
* `base_image` - (String) The ID of the base image from which the image will be built.
* `created_at` - (String) The date and time that the image recipe was created.
* `crn` - (String) The CRN for this image recipe.
* `deprecation_at` - (String) The deprecation date and time (UTC) for this image recipe.
* `href` - (String) The URL for this image recipe.
* `image_components` - (List of Strings) The IDs of the image components included in this image recipe.
* `name` - (String) The unique name for this image recipe.
* `obsolescence_at` - (String) The obsolescence date and time (UTC) for this image recipe.
* `resource_group` - (String) The ID of the resource group for this image recipe.
* `resource_type` - (String) The resource type (`image_recipe`).
* `status` - (String) The status of the image recipe (`deprecated`).
* `status_reasons` - (List) The reasons for the current status (if any).
  Nested schema for `status_reasons`:
  * `code` - (String) A reason code for this status.
  * `message` - (String) An explanation of the status reason.
  * `more_info` - (String) A link to documentation about this status reason.
* `superseded_at` - (String) The date and time that this image recipe was superseded by a newer version (if applicable).
* `updated_at` - (String) The date and time that the image recipe was last updated.
* `version` - (String) A generated content fingerprint for this version of the image recipe.
* `visibility` - (String) The visibility of this image recipe (`private`, `public`).

## Import

You can import the `ibm_is_image_recipe_deprecate` resource by using the image recipe ID:

```bash
$ terraform import ibm_is_image_recipe_deprecate.example r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8
```
