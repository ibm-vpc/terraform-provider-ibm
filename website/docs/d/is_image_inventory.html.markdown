---
layout: "ibm"
page_title: "IBM : ibm_is_image_inventory"
description: |-
  Get information about ImageInventory
subcategory: "Virtual Private Cloud API"
---

# ibm_is_image_inventory

Provides a read-only data source to retrieve information about an ImageInventory. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_image_inventory" "is_image_inventory" {
	is_image_inventory_id = "is_image_inventory_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `is_image_inventory_id` - (Required, Forces new resource, String) The image identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the ImageInventory.
* `created_at` - (String) The date and time when image inventory was created.
* `packages` - (List) 
  * Constraints: The maximum length is `1000` items. The minimum length is `0` items.
Nested schema for **packages**:
	* `name` - (String) The package name.
	  * Constraints: The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9][a-zA-Z0-9._+-]*$/`.
	* `version` - (String) The package version.
	  * Constraints: The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9][a-zA-Z0-9:.+~_-]*$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `image_inventory`. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

