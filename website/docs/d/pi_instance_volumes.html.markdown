---

subcategory: "Power Systems"
layout: "ibm"
page_title: "IBM: pi_instance_volumes"
description: |-
  Manages Instance volumes in the Power Virtual Server Cloud.
---

# ibm\_pi_instance_volumes

Import the details of existing IBM Power Virtual Server Cloud instance volumes as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_pi_instance_volumes" "ds_volumes" {
  pi_instance_name     = "volume_1"
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

* `pi_instance_name` - (Required, string) The name of the instance whose volumes to retrieve.
* `pi_cloud_instance_id` - (Required, string) The GUID of the service instance associated with the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:
* `boot_volume_id` - The unique identifier of the boot volume.
* `instance_volumes` - List of volumes attached to instance.
  * `id` - The ID of this volume.
  * `type` - The disk type for this volume.
  * `state` - The state of the volume.
  * `bootable` - If this volume is bootable or not.
  * `size` - The size of this volume.
  * `shareable` - If this volume is shareable or not.
  * `href` - The href of this volume.
  * `name` - The name of this volume.
