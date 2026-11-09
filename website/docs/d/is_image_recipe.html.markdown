---
layout: "ibm"
page_title: "IBM : ibm_is_image_recipe"
description: |-
  Get information about an IBM Cloud VPC Image Recipe.
subcategory: "VPC infrastructure"
---

# ibm_is_image_recipe

Provides a read-only data source to retrieve information about an IBM Cloud VPC Image Recipe by its ID or unique name.

## Example Usage

### Retrieve by ID

```terraform
data "ibm_is_image_recipe" "example" {
  image_recipe = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
}
```

### Retrieve by Name

```terraform
data "ibm_is_image_recipe" "example" {
  name = "my-image-recipe"
}
```

## Argument Reference

You can specify one of the following arguments for this data source (exactly one is required):

* `image_recipe` - (Optional, String) The image recipe identifier. Exactly one of `image_recipe` or `name` must be specified.
* `name` - (Optional, String) The unique name for this image recipe. Exactly one of `image_recipe` or `name` must be specified.

## Attribute Reference

After your data source is read, you can reference the following attributes:

* `id` - (String) The unique identifier of this image recipe.
* `base_image` - (List) The base image that is extended when building an image from this image recipe.
  Nested schema for `base_image`:
  * `crn` - (String) The CRN for this image.
  * `deleted` - (List) If present, indicates the referenced resource has been deleted.
    Nested schema for `deleted`:
    * `more_info` - (String) A link to documentation about deleted resources.
  * `href` - (String) The URL for this image.
  * `id` - (String) The unique identifier for this image.
  * `name` - (String) The name for this image.
  * `remote` - (List) If present, indicates the resource is remote.
    Nested schema for `remote`:
    * `account` - (List) Identifies the owning account.
      Nested schema for `account`:
      * `id` - (String) The unique identifier for this account.
      * `resource_type` - (String) The resource type.
    * `region` - (List) Identifies the native region.
      Nested schema for `region`:
      * `href` - (String) The URL for this region.
      * `name` - (String) The globally unique name for this region.
  * `resource_type` - (String) The resource type (`image`).
* `created_at` - (String) The date and time that the image recipe was created.
* `crn` - (String) The CRN for this image recipe.
* `deprecation_at` - (String) The deprecation date and time (UTC) for this image recipe (if set).
* `href` - (String) The URL for this image recipe.
* `image_components` - (List) The image components included in this image recipe.
  Nested schema for `image_components`:
  * `crn` - (String) The CRN for this component.
  * `deleted` - (List) If present, indicates the component has been deleted.
    Nested schema for `deleted`:
    * `more_info` - (String) A link to documentation about deleted resources.
  * `href` - (String) The URL for this component.
  * `id` - (String) The unique identifier for this component.
  * `name` - (String) The unique name for this component.
  * `resource_type` - (String) The resource type (`image_component`).
* `name` - (String) The unique name for this image recipe.
* `obsolescence_at` - (String) The obsolescence date and time (UTC) for this image recipe (if set).
* `resource_group` - (List) The resource group for this image recipe.
  Nested schema for `resource_group`:
  * `href` - (String) The URL for this resource group.
  * `id` - (String) The unique identifier for this resource group.
  * `name` - (String) The name for this resource group.
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
