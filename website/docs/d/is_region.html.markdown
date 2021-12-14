---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : region"
description: |-
  Manages IBM Cloud region.
---

# ibm_is_region
Retrieve information about a VPC Generation 2 Compute region. For more information, about managing IBM Cloud region, see [creating a VPC in a different region](https://cloud.ibm.com/docs/vpc?topic=vpc-creating-a-vpc-in-a-different-region).

## Example usage

```terraform

data "ibm_is_region" "ds_region" {
  name = "us-south"
}

```

```terraform

data "ibm_is_region" "default_region" {
}

```

## Argument reference
Review the argument references that you can specify for your data source. 

- `name` - (Optional, String) The name of the region. If no `name` is provided then default region `name` is taken from the provider block.

## Attribute reference
In addition to the argument reference list, you can access the following attribute references after your data source is created. 

- `status` - (String) The status of the region.
- `endpoint` - (String) The endpoint of the region.
