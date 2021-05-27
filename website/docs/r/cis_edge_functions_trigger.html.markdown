---

subcategory: "Internet services"
layout: "ibm"
page_title: "IBM: ibm_cis_edge_functions_triggers"
description: |-
  Provides a IBM CIS Edge Functions Action resource.
---

# ibm_cis_edge_functions_trigger

Provides a IBM CIS Edge Functions Trigger resource. This resource is associated with an IBM Cloud Internet Services instance and a CIS Domain resource. It allows to create, update, delete Edge Functions Trigger of a domain of a CIS instance

## Example Usage

```terraform
# Add a Edge Functions Trigger to the domain
resource "ibm_cis_edge_functions_trigger" "test_trigger" {
  cis_id      = ibm_cis_edge_functions_action.test_action.cis_id
  domain_id   = ibm_cis_edge_functions_action.test_action.domain_id
  action_name = ibm_cis_edge_functions_action.test_action.action_name
  pattern_url = "example.com/*"
}
```

## Argument Reference

The following arguments are supported:

- `cis_id` - (Required,string) The ID of the CIS service instance
- `domain_id` - (Required,string) The ID of the domain to add the edge functions action.
- `pattern_url` - (Required,string) The domain name pattern on which the edge function action trigger should be executed.
- `action_name` - (Optional, string) The action name of the edge functions action which the trigger associates to. If it is not specified, then the trigger will be disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The Action ID. It is a combination of <`trigger_id`>,<`domain_id`>,<`cis_id`> attributes concatenated with ":".
- `pattern_url` - The Route pattern. It is a domain name which the action will be performed.
- `trigger_id` - The Route ID of action trigger.
- `action_name` - The Edge Functions Action Script name.
- `request_limit_fail_open` - The Action request limit fail open

## Import

The `ibm_cis_edge_functions_trigger` resource can be imported using the `id`. The ID is formed from the `Edge Functions Trigger Route ID`, the `Domain ID` of the domain and the `CRN` (Cloud Resource Name) concatentated using a `:` character.

The Domain ID and CRN will be located on the **Overview** page of the Internet Services instance under the **Domain** heading of the UI, or via using the `ibmcloud cis` CLI commands.

- **Domain ID** is a 32 digit character string of the form: `9caf68812ae9b3f0377fdf986751a78f`

- **CRN** is a 120 digit character string of the form: `crn:v1:bluemix:public:internet-svcs:global:a/4ea1882a2d3401ed1e459979941966ea:31fa970d-51d0-4b05-893e-251cba75a7b3::`

- **Edge Functions Trigger Route ID** is a 32 digit character string of the form : `48996f0da6ed76251b475971b097205c`.

```
$ terraform import ibm_cis_edge_functions_trigger.test_trigger <trigger_id>:<domain-id>:<crn>

$ terraform import ibm_cis_edge_functions_trigger.test_trigger 48996f0da6ed76251b475971b097205c:9caf68812ae9b3f0377fdf986751a78f:crn:v1:bluemix:public:internet-svcs:global:a/4ea1882a2d3401ed1e459979941966ea:31fa970d-51d0-4b05-893e-251cba75a7b3::
```
