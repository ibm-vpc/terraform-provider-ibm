---
subcategory: "Power Systems"
layout: "ibm"
page_title: "IBM: pi_cloud_instance"
description: |-
  Provides data for a cloud instance  in the IBM Power Virtual Server Cloud.
---

# ibm\_pi_cloud_instance

Import the details of an existing IBM Power Virtual Server Cloud Instance as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_pi_cloud_instance" "ds_cloud_instance" {
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

* `pi_cloud_instance_id` - (Required, string) The GUID of the service instance associated with the account

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this tenant.
* `creation_date` - The date on which the tenant was created.
* `enabled` - Indicates whether the tenant is enabled.
* `tenant_name` - The name of the tenant.
* `capabilities` - Lists the Capabilities for this cloud instance.
* `total_instances` - The count of lpars that belong to this specific cloud instance
* `total_memory_consumed` - The total memory consumed by this service instance
* `total_processors_consumed` - The total processors consumed by this service instance
* `total_ssd_storage_consumed` - The total SSD Storage consumed by this service instance
* `total_standard_storage_consumed` - The total Standard Storage consumed by this service instance

  
