---
layout: "ibm"
page_title: "IBM : is_share_profiles"
sidebar_current: "docs-ibm-datasource-is-share-profiles"
description: |-
  Get information about ShareProfileCollection
subcategory: "Virtual Private Cloud API"
---

# ibm\_is_share_profiles

Provides a read-only data source for ShareProfileCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "is_share_profiles" "is_share_profiles" {
}
```

## Argument Reference

The following arguments are supported:


## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the ShareProfileCollection.
* `first` - A link to the first page of resources. Nested `first` blocks have the following structure:
	* `href` - The URL for a page of resources.

* `limit` - The maximum number of resources that can be returned by the request.

* `next` - A link to the next page of resources. This property is present for all pagesexcept the last page. Nested `next` blocks have the following structure:
	* `href` - The URL for a page of resources.

* `profiles` - Collection of share profiles. Nested `profiles` blocks have the following structure:
	* `family` - The product family this share profile belongs to.
	* `href` - The URL for this share profile.
	* `name` - The globally unique name for this share profile.
	* `resource_type` - The resource type.

* `total_count` - The total number of resources across all pages.

