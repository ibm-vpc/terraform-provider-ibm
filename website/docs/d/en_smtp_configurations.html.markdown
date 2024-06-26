---
subcategory: 'Event Notifications'
layout: 'ibm'
page_title: 'IBM : ibm_en_smtp_configurations'
description: |-
  List all the SMTP Configurations
---

# ibm_en_smtp_configurations

Provides a read-only data source for SMTP Configurations. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example usage

```terraform
data "ibm_en_smtp_configurations" "smtp_config_list" {
  instance_guid = ibm_resource_instance.en_terraform_test_resource.guid
}
```

## Argument reference

Review the argument reference that you can specify for your data source.

- `instance_guid` - (Required, Forces new resource, String) Unique identifier for IBM Cloud Event Notifications instance.

- `search_key` - (Optional, String) Filter the SMTP Configuration by name.

## Attribute reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

- `id` - The unique identifier of the `smtp_config_list`.

- `smtp_configurations` - (List) List of SMTP Configurations.

  - `id` - (String) Autogenerated SMTP Configuration ID.

  - `name` - (String) Name of the SMTP Configuration.

  - `description` - (String) SMTP description.

  - `domain` - (String) Domain Name.

  - `updated_at` - (Stringr) Created time.

- `total_count` - (Integer) Total number of SMTP configurations.
