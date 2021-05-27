---

subcategory: "Internet services"
layout: "ibm"
page_title: "IBM: ibm_cis_edge_functions_action"
description: |-
  Provides a IBM CIS Edge Functions Action resource.
---

# ibm_cis_edge_functions_action

Provides a IBM CIS Edge Functions Action resource. This resource is associated with an IBM Cloud Internet Services instance and a CIS Domain resource. It allows to create, update, delete Edge Functions Action of a domain of a CIS instance

## Example Usage

```terraform
# Add a Edge Functions Action to the domain
resource "ibm_cis_edge_functions_action" "test_action" {
  cis_id      = data.ibm_cis.cis.id
  domain_id   = data.ibm_cis_domain.cis_domain.domain_id
  action_name = "sample-script"
  script      = file("./script.js")
}
```

## Argument Reference

The following arguments are supported:

- `cis_id` - (Required,string) The ID of the CIS service instance
- `domain_id` - (Required,string) The ID of the domain to add the edge functions action.
- `action_name` - (Required,string) The Action Name of the edge functions action.
- `script` - (Required, string) The script of the edge functions action.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The Action ID. It is a combination of <`action_name`>,<`domain_id`>,<`cis_id`> attributes concatenated with ":".

## Import

The `ibm_cis_edge_functions_action` resource can be imported using the `id`. The ID is formed from the `Edge Functions Action Name/Script Name`, the `Domain ID` of the domain and the `CRN` (Cloud Resource Name) concatentated using a `:` character.

The Domain ID and CRN will be located on the **Overview** page of the Internet Services instance under the **Domain** heading of the UI, or via using the `ibmcloud cis` CLI commands.

- **Domain ID** is a 32 digit character string of the form: `9caf68812ae9b3f0377fdf986751a78f`

- **CRN** is a 120 digit character string of the form: `crn:v1:bluemix:public:internet-svcs:global:a/4ea1882a2d3401ed1e459979941966ea:31fa970d-51d0-4b05-893e-251cba75a7b3::`

- **Edge Functions Action Name/Script Name** is a string : `sample_script`.

```
$ terraform import ibm_cis_edge_functions_action.test_action <action_name>:<domain-id>:<crn>

$ terraform import ibm_cis_edge_functions_action.test_action sample_script:9caf68812ae9b3f0377fdf986751a78f:crn:v1:bluemix:public:internet-svcs:global:a/4ea1882a2d3401ed1e459979941966ea:31fa970d-51d0-4b05-893e-251cba75a7b3::
```