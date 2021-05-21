---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_placement_groups"
description: |-
  Get information about PlacementGroupCollection
---

# ibm\_is_placement_groups

Provides a read-only data source for PlacementGroupCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "is_placement_groups" "is_placement_groups" {
}
```

## Argument Reference

The following arguments are supported:


## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the PlacementGroupCollection.
* `first` - A link to the first page of resources. Nested `first` blocks have the following structure:
	* `href` - The URL for a page of resources.

* `limit` - The maximum number of resources that can be returned by the request.

* `next` - A link to the next page of resources. This property is present for all pagesexcept the last page. Nested `next` blocks have the following structure:
	* `href` - The URL for a page of resources.

* `placement_groups` - Collection of placement groups. Nested `placement_groups` blocks have the following structure:
	* `created_at` - The date and time that the placement group was created.
	* `crn` - The CRN for this placement group.
	* `href` - The URL for this placement group.
	* `id` - The unique identifier for this placement group.
	* `lifecycle_state` - The lifecycle state of the placement group.
	* `name` - The user-defined name for this placement group.
	* `resource_group` - The unique identifier of the resource group for this placement group. 
	* `resource_type` - The resource type.
	* `strategy` - The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.

* `total_count` - The total number of resources across all pages.

