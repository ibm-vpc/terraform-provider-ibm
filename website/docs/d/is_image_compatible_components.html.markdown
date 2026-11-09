---
layout: "ibm"
page_title: "IBM : ibm_is_image_compatible_components"
description: |-
  Get information about image components compatible with an image
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_compatible_components

Provides a read-only data source to retrieve information about image components that are compatible with a specific base image based on the image component's `allowed_use.image` property. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_images" "images" {
  visibility = "public"
}

data "ibm_is_image_compatible_components" "compatible_components" {
  image_id = data.ibm_is_images.images.images.0.id
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `image_id` - (Required, String) The unique identifier of the image to retrieve compatible components for.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - (String) The image identifier.
* `image_components` - (List) A page of image components compatible with this image.
	Nested schema for **image_components**:
	* `crn` - (String) The CRN for this image component.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.
		Nested schema for **deleted**:
		* `more_info` - (String) A link to documentation about deleted resources.
	* `href` - (String) The URL for this image component.
	* `id` - (String) The unique identifier for this image component.
	* `name` - (String) The unique name for this image component.
	* `resource_type` - (String) The resource type.
