---
subcategory: "Internet services"
layout: "ibm"
page_title: "IBM: ibm_cis_edge_functions_triggers"
description: |-
  Get information on an IBM Cloud Internet Services Edge Function Triggers.
---

# ibm_cis_edge_functions_triggers

Imports a read only copy of an existing Internet Services Edge Function Triggers resource.

## Example Usage

```terraform
data "ibm_cis_edge_functions_triggers" "test_triggers" {
    cis_id    = data.ibm_cis.cis.id
    domain_id = data.ibm_cis_domain.cis_domain.domain_id
}
```

## Argument Reference

The following arguments are supported:

- `cis_id` - (Required,string) The ID of the CIS service instance
- `domain_id` - (Required,string) The ID of the domain to add the edge functions triggers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `pattern_url` - The Route pattern. It is a domain name which the action will be performed.
- `trigger_id` - The Route ID of action trigger.
- `action_name` - The Action Script for execution.
- `request_limit_fail_open` - The Action request limit fail open
