---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_placement_group"
description: |-
  Manages PlacementGroup.
---

# ibm\_is_placement_group

Provides a resource for PlacementGroup. This allows PlacementGroup to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_is_placement_group" "is_placement_group" {
  strategy = "host_spread"
  name = "my-placement-group"
}
```

## Argument Reference

The following arguments are supported:

* `strategy` - (Required, string) The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sources. The enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.
* `name` - (Optional, string) The unique user-defined name for this placement group. If unspecified, the name will be a hyphenated list of randomly-selected words.
* `resource_group` - (Optional, string) The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the PlacementGroup.
* `created_at` - The date and time that the placement group was created.
* `crn` - The CRN for this placement group.
* `href` - The URL for this placement group.
* `lifecycle_state` - The lifecycle state of the placement group.
* `resource_type` - The resource type.

## Import

ibm_is_placement_group can be imported using ID, eg

```
$ terraform import ibm_is_placement_group.example d7bec597-4726-451f-8a63-e62e6f19c32c
```