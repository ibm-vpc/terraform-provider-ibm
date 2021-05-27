---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_dedicated_hosts"
description: |-
  Get information about DedicatedHostCollection
---

# ibm\_is_dedicated_hosts

Provides a read-only data source for DedicatedHostCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_dedicated_hosts" "is_dedicated_hosts" {
	host_group = "1e09281b-f177-46fb-baf1-bc152b2e391a"
}
```

## Argument Reference

The following arguments are supported:

* `host_group` - (Optional, string) The unique identifier for the dedicated host group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the DedicatedHostCollection.
* `dedicated_hosts` - Collection of dedicated hosts. Nested `dedicated_hosts` blocks have the following structure:
	* `available_memory` - The amount of memory in gibibytes that is currently available for instances.
	* `available_vcpu` - The available VCPU for the dedicated host. Nested `available_vcpu` blocks have the following structure:
		* `architecture` - The VCPU architecture.
		* `count` - The number of VCPUs assigned.
	* `created_at` - The date and time that the dedicated host was created.
	* `crn` - The CRN for this dedicated host.
	* `host_group` - The unique identifier of the dedicated host group this dedicated host is in.
	* `disks` - Collection of the dedicated host's disks. Nested `disks` blocks have the following structure:
		* `available` - The remaining space left for instance placement in GB (gigabytes).
		* `created_at` - The date and time that the disk was created.
		* `href` - The URL for this disk.
		* `id` - The unique identifier for this disk.
		* `instance_disks` - Instance disks that are on this dedicated host disk. Nested `instance_disks` blocks have the following structure:
			* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
				* `more_info` - Link to documentation about deleted resources.
			* `href` - The URL for this instance disk.
			* `id` - The unique identifier for this instance disk.
			* `name` - The user-defined name for this disk.
			* `resource_type` - The resource type.
		* `interface_type` - The disk interface used for attaching the diskThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.
		* `lifecycle_state` - The lifecycle state of this dedicated host disk.
		* `name` - The user-defined or system-provided name for this disk.
		* `provisionable` - Indicates whether this dedicated host disk is available for instance disk creation.
		* `resource_type` - The type of resource referenced.
		* `size` - The size of the disk in GB (gigabytes).
		* `supported_instance_interface_types` - The instance disk interfaces supported for this dedicated host disk.
	* `href` - The URL for this dedicated host.
	* `id` - The unique identifier for this dedicated host.
	* `instance_placement_enabled` - If set to true, instances can be placed on this dedicated host.
	* `instances` - Array of instances that are allocated to this dedicated host. Nested `instances` blocks have the following structure:
		* `crn` - The CRN for this virtual server instance.
		* `deleted` - If present, this property indicates the referenced resource has been deleted and providessome supplementary information. Nested `deleted` blocks have the following structure:
			* `more_info` - Link to documentation about deleted resources.
		* `href` - The URL for this virtual server instance.
		* `id` - The unique identifier for this virtual server instance.
		* `name` - The user-defined name for this virtual server instance (and default system hostname).
	* `lifecycle_state` - The lifecycle state of the dedicated host resource.
	* `memory` - The total amount of memory in gibibytes for this host.
	* `name` - The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.
	* `profile` - The profile this dedicated host uses. Nested `profile` blocks have the following structure:
		* `href` - The URL for this dedicated host.
		* `name` - The globally unique name for this dedicated host profile.
	* `provisionable` - Indicates whether this dedicated host is available for instance creation.
	* `resource_group` - The resource group for this dedicated host. Nested `resource_group` blocks have the following structure:
		* `href` - The URL for this resource group.
		* `id` - The unique identifier for this resource group.
		* `name` - The user-defined name for this resource group.
	* `resource_type` - The type of resource referenced.
	* `socket_count` - The total number of sockets for this host.
	* `state` - The administrative state of the dedicated host.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the dedicated host on which the unexpected property value was encountered.
	* `supported_instance_profiles` - Array of instance profiles that can be used by instances placed on this dedicated host. Nested `supported_instance_profiles` blocks have the following structure:
		* `href` - The URL for this virtual server instance profile.
		* `name` - The globally unique name for this virtual server instance profile.
	* `vcpu` - The total VCPU of the dedicated host. Nested `vcpu` blocks have the following structure:
		* `architecture` - The VCPU architecture.
		* `count` - The number of VCPUs assigned.
	* `zone` - The globally unique name of the zone this dedicated host resides in.
* `total_count` - The total number of resources across all pages.

