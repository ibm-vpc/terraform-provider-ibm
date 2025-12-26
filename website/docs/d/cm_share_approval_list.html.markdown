---
layout: "ibm"
page_title: "IBM : ibm_cm_share_approval_list"
description: |-
  Get information about cm_share_approval_list
subcategory: "Catalog Management"
---

# ibm_cm_share_approval_list

Provides a read-only data source to retrieve information about catalog share approval lists. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_cm_share_approval_list" "cm_share_approval_list" {
  object_kind     = "offering"
  approval_state  = "approved"
}
```

## Example Usage with Limit

```hcl
data "ibm_cm_share_approval_list" "cm_share_approval_list_limited" {
  object_kind     = "offering"
  approval_state  = "pending"
  limit           = 50
}
```

## Example Usage with Enterprise Context

```hcl
data "ibm_cm_share_approval_list" "cm_share_approval_list_enterprise" {
  object_kind     = "offering"
  approval_state  = "approved"
  enterprise_id   = "-ent-enterprise123"
  limit           = 100
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `object_kind` - (Required, String) The object kind for share approval. Options are `offering`, `vpe`, `proxy_source`, or `preset_configuration`.
  * Constraints: Allowable values are: `offering`, `vpe`, `proxy_source`, `preset_configuration`.
* `approval_state` - (Required, String) The approval state to filter by. Options are `approved`, `pending`, or `rejected`.
  * Constraints: Allowable values are: `approved`, `pending`, `rejected`.
* `enterprise_id` - (Optional, String) Enterprise or enterprise account group ID to view or manage requests for the enterprise. Prefix with `-ent-` for an enterprise and `-entgrp-` for an account group.
* `limit` - (Optional, Integer) Number of results to return in the query. Default is 100.
  * Constraints: The default value is `100`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the cm_share_approval_list.
* `total_count` - (Integer) The total count of resources that match the request.
* `resource_count` - (Integer) The number of resources returned in this response.
* `resources` - (List) A list of share approval access records.
  
  Nested schema for **resources**:
  * `id` - (String) Unique ID of the access record.
  * `account` - (String) Account ID.
  * `account_type` - (Integer) Account type (normal account or enterprise).
  * `target_account` - (String) Object's owner's account.
  * `target_kind` - (String) Entity type.
  * `created` - (String) Date and time the access record was created.
  * `approval_state` - (String) Approval state of the access record.