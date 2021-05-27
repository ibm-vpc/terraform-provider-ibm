---
subcategory: "Power Systems"
layout: "ibm"
page_title: "IBM: pi_image"
description: |-
  Manages an image in the Power Virtual Server Cloud.
---

# ibm\_pi_image

Import the details of an existing IBM Power Virtual Server Cloud image as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_pi_image" "ds_image" {
  pi_image_name        = "7200-03-03"
  pi_cloud_instance_id = "49fba6c9-23f8-40bc-9899-aca322ee7d5b"
}
```
## Notes:
* Please find [supported Regions](https://cloud.ibm.com/apidocs/power-cloud#endpoint) for endpoints.
* If a Power cloud instance is provisioned at `lon04`, The provider level attributes should be as follows:
  * `region` - `lon`
  * `zone` - `lon04`
  Example Usage:
  ```terraform
    provider "ibm" {
      region    =   "lon"
      zone      =   "lon04"
    }
  ```
## Argument Reference

The following arguments are supported:

* `pi_image_name` - (Required, string) The name of the image.
* `pi_cloud_instance_id` - (Required, string) The GUID of the service instance associated with the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this image.
* `state` - The state for this image.
* `size` - The size of the image.
* `architecture` - The architecture for this image.
* `operatingsystem` - The operating system for this image.
* `storage_type`  - The storage type for this image
