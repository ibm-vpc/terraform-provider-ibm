---
subcategory: "Catalog Management"
layout: "ibm"
page_title: "IBM : cm_offering_instance"
description: |-
  Manages cm_offering_instance.
---

# ibm_cm_offering_instance

Create, modify, or delete an `ibm_cm_offering_instance` resources. You can manage the settings for all catalogs across your account. Management tasks include setting the visibility of the IBM Cloud catalog and controlling access to products in the public catalog and private catalogs for users in your account. For more information, about managing catalog, refer to [catalog management settings](https://cloud.ibm.com/docs/account?topic=account-account-getting-started).


## Example usage

```terraform
resource "ibm_cm_offering_instance" "cm_offering_instance" {
  catalog_id = "catalog_id"
  label = "placeholder"
  kind_format = "operator"
  version = "placeholder"
  cluster_id = "placeholder"
  cluster_region = "placeholder"
  cluster_namespaces = [ "placeholder", "placeholder2" ]
  cluster_all_namespaces = false
}
```

## Argument reference
Review the argument reference that you can specify for your resource. 

- `catalog_id` - (Required, String) The catalog ID an instance  is created.
- `cluster_id` - (Required, String) The cluster ID.
- `cluster_region` - (Required, String) The cluster region for example, `us-south`.
- `cluster_namespaces`- (Required, List) The list of target namespaces to install into.
- `cluster_all_namespaces`- (Required, Bool) Designate to install into all namespaces.
- `kind_format` - (Required, String) The format an instance such as `helm`, `operator`, `ova`. **Note** Currently the only supported formate is `operator`.
- `label` - (Required, String) The label for this instance.
- `offering_id` - (Required, String) The offering ID an instance is created .
- `version` - (Required, String) The version an instance was installed from (but not from the version ID).


## Attribute reference
In addition to all argument references list, you can access the following attribute references after your resource is created. 

- `crn` - (String) The platform CRN for an instance.
- `id` - (String) The unique identifier of the `cm_offering_instance`.
- `url` - (String) The URL reference to an object.
- `_rev` - (String) The cloudant revision of this object
- `schematics_workspace_id` - (String) The ID of the schematics workspace used to install this offering, if applicable
