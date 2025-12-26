---
layout: "ibm"
page_title: "IBM : ibm_cm_share_approval"
description: |-
  Manages cm_share_approval.
subcategory: "Catalog Management"
---

# ibm_cm_share_approval

Create, update, and delete catalog share approvals with this resource. This resource allows you to manage approval states for accounts that want to access your catalog offerings.

## Example Usage

```hcl
resource "ibm_cm_share_approval" "cm_share_approval_instance" {
  object_kind     = "offering"
  approval_state  = "approved"
  account_ids     = [
    "-acct-c3b6f7aa72dc448ab4e205026613e0c0",
    "-ent-enterprise123"
  ]
}
```

## Example Usage with Enterprise Context

```hcl
resource "ibm_cm_share_approval" "cm_share_approval_enterprise" {
  object_kind     = "offering"
  approval_state  = "approved"
  account_ids     = [
    "-acct-account123",
    "-entgrp-group456"
  ]
  enterprise_id   = "-ent-enterprise789"
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `object_kind` - (Required, Forces new resource, String) The object kind for share approval. Options are `offering`, `vpe`, `proxy_source`, or `preset_configuration`.
  * Constraints: Allowable values are: `offering`, `vpe`, `proxy_source`, `preset_configuration`.
* `approval_state` - (Required, String) The approval state. Options are `approved`, `pending`, or `rejected`.
  * Constraints: Allowable values are: `approved`, `pending`, `rejected`.
* `account_ids` - (Required, List) List of account IDs to set approval for. Each account ID must be prefixed with one of the following:
  * `-acct-` for regular accounts
  * `-ent-` for enterprise accounts
  * `-entgrp-` for enterprise account groups
* `enterprise_id` - (Optional, String) Enterprise or enterprise account group ID to view or manage requests for the enterprise. Prefix with `-ent-` for an enterprise and `-entgrp-` for an account group.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the cm_share_approval. The ID is in the format `{object_kind}/{approval_state}`.

## Import

You can import the `ibm_cm_share_approval` resource by using `id`. The ID is in the format `{object_kind}/{approval_state}`.

# Syntax
<pre>
$ terraform import ibm_cm_share_approval.cm_share_approval <object_kind>/<approval_state>
</pre>

# Example
```
$ terraform import ibm_cm_share_approval.cm_share_approval offering/approved