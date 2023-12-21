---
layout: "ibm"
page_title: "IBM : ibm_is_storage_ontap_instances"
description: |-
  Get information about StorageOntapInstanceCollection
subcategory: "ontap"
---

# ibm_is_storage_ontap_instances

Provides a read-only data source to retrieve information about a StorageOntapInstanceCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```terraform
data "ibm_is_storage_ontap_instances" "is_storage_ontap_instances" {
	lifecycle_state = ibm_is_storage_ontap_instance.is_storage_ontap_instance.lifecycle_state
}
```

## Argument Reference

You can specify the following arguments for this data source.

- `lifecycle_state` - (Optional, String) Filters the collection to storage ontap instances with a `lifecycle_state` property matching the specified value.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
- `resource_group_id` - (Optional, String) Filters the collection to resources with a `resource_group.id` property matching the specified identifier.


## Attribute Reference

After your data source is created, you can read values from the following attributes.

- `id` - The unique identifier of the StorageOntapInstanceCollection.
- `first` - (List) A link to the first page of resources.
Nested schema for **first**:
	- `href` - (String) The URL for a page of resources.
	  

- `limit` - (Integer) The maximum number of resources that can be returned by the request.
  * Constraints: The maximum value is `100`. The minimum value is `1`.

- `next` - (List) A link to the next page of resources. This property is present for all pagesexcept the last page.
Nested schema for **next**:
	- `href` - (String) The URL for a page of resources.
	  

- `storage_ontap_instances` - (List) Collection of storage ontap instances.
Nested schema for **storage_ontap_instances**:
	- `address_prefix` - (List) An address prefix in the VPC which will be used to allocate `endpoints` for thisstorage ontap instance and its storage virtual machines.
	Nested schema for **address_prefix**:
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `href` - (String) The URL for this address prefix.
		  
		- `id` - (String) The unique identifier for this address prefix.
		
		- `name` - (String) The name for this address prefix. The name is unique across all address prefixes for the VPC.
		  
	- `admin_credentials` - (List) The credentials used (from Secrets Manager) for the cluster administrator to access thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.
	Nested schema for **admin_credentials**:
		- `http` - (List) The security certificate credential for ONTAP REST API access for the clusteradministrator of the storage ontap instance.
		Nested schema for **http**:
			- `crn` - (String) The CRN for this credential.
			  
			- `resource_type` - (String) The resource type.
			  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		- `password` - (List) The password credential for the cluster administrator of the storage ontap instance.If present, this password credential is used by the cluster administrator for bothONTAP CLI SSH access and ONTAP REST API access. If absent, the storage ontapinstance is not accessible through either the ONTAP CLI or ONTAP REST API usingpassword-based authentication.
		Nested schema for **password**:
			- `crn` - (String) The CRN for this credential.
			  
			- `resource_type` - (String) The resource type.
			  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		- `ssh` - (List) The public key credential for ONTAP CLI SSH access for the cluster administratorof the storage ontap instance.
		Nested schema for **ssh**:
			- `crn` - (String) The CRN for this credential.
			  
			- `resource_type` - (String) The resource type.
			  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `capacity` - (Integer) The capacity to use for the storage ontap instance (in terabytes). Volumes in this storage ontap instance will be allocated from this capacity.
	  * Constraints: The maximum value is `64`. The minimum value is `1`.
	- `created_at` - (String) The date and time that the storage ontap instance was created.
	- `crn` - (String) The CRN for this storage ontap instance.
	  
	- `encryption` - (String) The type of encryption used on the storage ontap instance.
	  * Constraints: Allowable values are: `provider_managed`, `user_managed`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `encryption_key` - (List) The root key used to wrap the data encryption key for the storage ontap instance.This property will be present for storage ontap instance with an `encryption` type of`user_managed`.
	Nested schema for **encryption_key**:
		- `crn` - (String) The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Services Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.
		  
	- `endpoints` - (List) The endpoints for this storage ontap instance.
	Nested schema for **endpoints**:
		- `inter_cluster` - (List) The NetApp SnapMirror management endpoint for this storage ontap instance.
		Nested schema for **inter_cluster**:
			- `ipv4_address` - (String) The unique IP address of an endpoint.
			  
		- `management` - (List) The NetApp management endpoint for this storage ontap instance. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.
		Nested schema for **management**:
			- `ipv4_address` - (String) The unique IP address of an endpoint.
			  
	- `health_reasons` - (List) The reasons for the current storage ontap instance health_state (if any):- `cluster_down`: This storage ontap instance is unavailable as both the  primary and secondary nodes are unavailable.- `failback_unavailable`: The capability to failback is unavailable. The secondary node  continues to be available.- `failover_unavailable`: The capability to failover is unavailable. The primary  node continues to be available without any performance impact to clients.- `internal_error`: Internal error (contact IBM support).- `maintenance_in_progress`: A planned maintenance activity is in progress.- `primary_node_down`: The primary node is unavailable, and I/O has failed over to  the secondary node. Clients running in the same zone as the primary node may  experience higher access latency.- `secondary_node_down`: The secondary node is unavailable. Therefore, the capability  to failover is unavailable.
	  * Constraints: The minimum length is `0` items.
	Nested schema for **health_reasons**:
		- `code` - (String) A snake case string succinctly identifying the reason for this health state.
		  * Constraints: Allowable values are: `cluster_down`, `failback_unavailable`, `failover_unavailable`, `internal_error`, `maintenance_in_progress`, `primary_node_down`, `secondary_node_down`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		- `message` - (String) An explanation of the reason for this health state.
		- `more_info` - (String) Link to documentation about the reason for this health state.
		  
	- `health_state` - (String) The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
	  * Constraints: Allowable values are: `degraded`, `faulted`, `inapplicable`, `ok`.
	- `href` - (String) The URL for this storage ontap instance.
	  
	- `id` - (String) The unique identifier for this storage ontap instance.
	
	- `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.
	  * Constraints: The minimum length is `0` items.
	Nested schema for **lifecycle_reasons**:
		- `code` - (String) A snake case string succinctly identifying the reason for this lifecycle state.
		  * Constraints: Allowable values are: `resource_suspended_by_provider`. The value must match regular expression `/^[a-z]+(_[a-z]+)*$/`.
		- `message` - (String) An explanation of the reason for this lifecycle state.
		- `more_info` - (String) Link to documentation about the reason for this lifecycle state.
		  
	- `lifecycle_state` - (String) The lifecycle state of the storage ontap instance.
	  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
	- `name` - (String) The name for this storage ontap instance. The name is unique across all storage ontap instances in the region.
	  * Constraints: The maximum length is `40` characters. The minimum length is `7` characters. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	- `primary_subnet` - (List) The subnet where the primary Cloud Volumes ONTAP node is provisioned in.
	Nested schema for **primary_subnet**:
		- `crn` - (String) The CRN for this subnet.
		  
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `id` - (String) The unique identifier for this subnet.
		
		- `name` - (String) The name for this subnet. The name is unique across all subnets in the VPC.
		  
		- `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `resource_group` - (List) The resource group for this storage ontap instance.
	Nested schema for **resource_group**:
		- `id` - (String) The unique identifier for this resource group.
		  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
		- `name` - (String) The name for this resource group.
		  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
	- `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `storage_ontap_instance`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `routing_tables` - (List) The VPC routing tables for this storage ontap instance.
	  * Constraints: The maximum length is `4` items. The minimum length is `1` item.
	Nested schema for **routing_tables**:
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `href` - (String) The URL for this routing table.
		  
		- `id` - (String) The unique identifier for this routing table.
		
		- `name` - (String) The name for this routing table. The name is unique across all routing tables for the VPC.
		  
		- `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `routing_table`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `secondary_subnet` - (List) The subnet where the secondary Cloud Volumes ONTAP node is provisioned in.
	Nested schema for **secondary_subnet**:
		- `crn` - (String) The CRN for this subnet.
		  
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `id` - (String) The unique identifier for this subnet.
		
		- `name` - (String) The name for this subnet. The name is unique across all subnets in the VPC.
		  
		- `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `security_groups` - (List) The security groups for this storage ontap instance.
	  * Constraints: The minimum length is `1` item.
	Nested schema for **security_groups**:
		- `crn` - (String) The security group's CRN.
		  
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `id` - (String) The unique identifier for this security group.
		
		- `name` - (String) The name for this security group. The name is unique across all security groups for the VPC.
		  
	- `storage_virtual_machines` - (List) The storage virtual machines for this storage ontap instance.
	  * Constraints: The minimum length is `1` item.
	Nested schema for **storage_virtual_machines**:
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `href` - (String) The URL for this storage virtual machine.
		  
		- `id` - (String) The unique identifier for this storage virtual machine.
		
		- `name` - (String) The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.
		  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
		- `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `storage_ontap_instance_storage_virtual_machine`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	- `vpc` - (List) The VPC this storage ontap instance resides in.
	Nested schema for **vpc**:
		- `crn` - (String) The CRN for this VPC.
		  
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
		Nested schema for **deleted**:
			- `more_info` - (String) Link to documentation about deleted resources.
			  
		- `id` - (String) The unique identifier for this VPC.
		
		- `name` - (String) The name for this VPC. The name is unique across all VPCs in the region.
		  
		- `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

- `total_count` - (Integer) The total number of resources across all pages.
  * Constraints: The minimum value is `0`.

