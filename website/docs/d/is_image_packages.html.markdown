---
layout: "ibm"
page_title: "IBM : ibm_is_image_packages"
description: |-
  Get information about the packages installed in an image
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_packages

Provides a read-only data source to retrieve the packages installed in an image built from an image recipe. The packages are sorted by ascending `name`. The package list is available only after a completed and successful image job. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_image_packages" "is_image_packages" {
	image_id = "r006-4f6a8a6c-1e6e-4c2e-9d8a-0f5f5d5a1b2c"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `image_id` - (Required, Forces new resource, String) The image identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `name` - (Optional, String) Filters the collection to packages with a `name` property matching the exact specified name.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the data source. It is the same as `image_id`.
* `packages` - (List) A list of packages installed in the image.
Nested schema for **packages**:
	* `name` - (String) The package name.
	  * Constraints: The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9][a-zA-Z0-9._+-]*$/`.
	* `version` - (String) The package version.
	  * Constraints: The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9][a-zA-Z0-9:.+~_-]*$/`.
