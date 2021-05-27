---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : zones"
description: |-
  Manages IBM Cloud Zones.
---

# ibm\_is_zones

Import the details of an existing IBM Cloud zones in a particular region as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```terraform
data "ibm_is_zones" "ds_zones" {
  region = "us-south"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, string) The name of the region.
* `status` - (Optional, string) Filter the list by status of zones.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - The list of zones in a region.