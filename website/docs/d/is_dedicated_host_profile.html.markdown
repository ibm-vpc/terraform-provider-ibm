---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_dedicated_host_profile"
description: |-
  Get information about DedicatedHostProfile
---

# ibm\_is_dedicated_host_profile

Provides a read-only data source for DedicatedHostProfile. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_dedicated_host_profile" "is_dedicated_host_profile" {
	name = "dh2-56x464"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The globally unique name for this virtual server instance profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the DedicatedHostProfile.
* `class` - The product class this dedicated host profile belongs to.
* `disks` - Collection of the dedicated host profile's disks. Nested `disks` blocks have the following structure:
	* `interface_type`  Nested `interface_type` blocks have the following structure:
		* `type` - The type for this profile field.
		* `value` - The interface of the disk for a dedicated host with this profileThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.
	* `quantity` - The number of disks of this type for a dedicated host with this profile. Nested `quantity` blocks have the following structure:
		* `type` - The type for this profile field.
		* `value` - The value for this profile field.
	* `size` - The size of the disk in GB (gigabytes). Nested `size` blocks have the following structure:
		* `type` - The type for this profile field.
		* `value` - The size of the disk in GB (gigabytes).
	* `supported_instance_interface_types`  Nested `supported_instance_interface_types` blocks have the following structure:
		* `type` - The type for this profile field.
		* `value` - The instance disk interfaces supported for a dedicated host with this profile.
* `family` - The product family this dedicated host profile belongs to.

* `href` - The URL for this dedicated host.

* `memory`  Nested `memory` blocks have the following structure:
	* `type` - The type for this profile field.
	* `value` - The value for this profile field.
	* `default` - The default value for this profile field.
	* `max` - The maximum value for this profile field.
	* `min` - The minimum value for this profile field.
	* `step` - The increment step value for this profile field.
	* `values` - The permitted values for this profile field.

* `socket_count`  Nested `socket_count` blocks have the following structure:
	* `type` - The type for this profile field.
	* `value` - The value for this profile field.
	* `default` - The default value for this profile field.
	* `max` - The maximum value for this profile field.
	* `min` - The minimum value for this profile field.
	* `step` - The increment step value for this profile field.
	* `values` - The permitted values for this profile field.

* `supported_instance_profiles` - Array of instance profiles that can be used by instances placed on dedicated hosts with this profile Nested`supported_instance_profiles` blocks have the following structure:
	* `href` - The URL for this virtual server instance profile.
	* `name` - The globally unique name for this virtual server instance profile.

* `vcpu_architecture`  Nested `vcpu_architecture` blocks have the following structure:
	* `type` - The type for this profile field.
	* `value` - The VCPU architecture for a dedicated host with this profile.

* `vcpu_count`  Nested `vcpu_count` blocks have the following structure:
	* `type` - The type for this profile field.
	* `value` - The value for this profile field.
	* `default` - The default value for this profile field.
	* `max` - The maximum value for this profile field.
	* `min` - The minimum value for this profile field.
	* `step` - The increment step value for this profile field.
	* `values` - The permitted values for this profile field.

