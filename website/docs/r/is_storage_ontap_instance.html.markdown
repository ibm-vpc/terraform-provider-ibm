---
layout: "ibm"
page_title: "IBM : ibm_is_storage_ontap_instance"
description: |-
  Manages StorageOntapInstance.
subcategory: "ontap"
---

# ibm_is_storage_ontap_instance

Create, update, and delete StorageOntapInstances with this resource.

## Example Usage

```hcl
resource "ibm_is_storage_ontap_instance" "is_storage_ontap_instance_instance" {
  address_prefix {
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/8e454ead-0db7-48ac-9a8b-2698d8c470a7/address_prefixes/1a15dca5-7e33-45e1-b7c5-bc690e569531"
		id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
		name = "my-address-prefix-1"
  }
  admin_credentials {
		http {
			crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
			resource_type = "credential"
		}
		password {
			crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
			resource_type = "credential"
		}
		ssh {
			crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
			resource_type = "credential"
		}
  }
  capacity = 10
  encryption_key {
		crn = "crn:v1:bluemix:public:kms:us-south:a/123456:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"
  }
  name = "ibmss-my-storage-ontap-instance"
  primary_subnet {
		crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		name = "my-subnet"
		resource_type = "subnet"
  }
  resource_group {
		id = "fee82deba12e4c0fb69c3b09d1f12345"
		name = "my-resource-group"
  }
  routing_tables {
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/982d72b7-db1b-4606-afb2-ed6bd4b0bed1/routing_tables/6885e83f-03b2-4603-8a86-db2a0f55c840"
		id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
		name = "my-routing-table-1"
		resource_type = "routing_table"
  }
  secondary_subnet {
		crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		name = "my-subnet"
		resource_type = "subnet"
  }
  security_groups {
		crn = "crn:v1:bluemix:public:is:us-south:a/123456::security-group:be5df5ca-12a0-494b-907e-aa6ec2bfa271"
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		id = "be5df5ca-12a0-494b-907e-aa6ec2bfa271"
		name = "my-security-group"
  }
  storage_virtual_machines {
		deleted {
			more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
		}
		href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/storage_ontap_instances/r134-d7cc5196-9864-48c4-82d8-3f30da41ffc5/storage_virtual_machines/r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
		id = "r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
		name = "my-storage-virtual-machine"
		resource_type = "storage_ontap_instance_storage_virtual_machine"
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `address_prefix` - (Required, List) An address prefix in the VPC which will be used to allocate `endpoints` for thisstorage ontap instance and its storage virtual machines.
Nested schema for **address_prefix**:
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this address prefix.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this address prefix.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this address prefix. The name is unique across all address prefixes for the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `admin_credentials` - (Optional, List) The credentials used (from Secrets Manager) for the cluster administrator to access thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.
Nested schema for **admin_credentials**:
	* `http` - (Optional, List) The security certificate credential for ONTAP REST API access for the clusteradministrator of the storage ontap instance.
	Nested schema for **http**:
		* `crn` - (Required, String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (Computed, String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `password` - (Optional, List) The password credential for the cluster administrator of the storage ontap instance.If present, this password credential is used by the cluster administrator for bothONTAP CLI SSH access and ONTAP REST API access. If absent, the storage ontapinstance is not accessible through either the ONTAP CLI or ONTAP REST API usingpassword-based authentication.
	Nested schema for **password**:
		* `crn` - (Required, String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (Computed, String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `ssh` - (Optional, List) The public key credential for ONTAP CLI SSH access for the cluster administratorof the storage ontap instance.
	Nested schema for **ssh**:
		* `crn` - (Required, String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (Computed, String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `capacity` - (Required, Integer) The capacity to use for the storage ontap instance (in terabytes). Volumes in this storage ontap instance will be allocated from this capacity.
  * Constraints: The maximum value is `64`. The minimum value is `1`.
* `encryption_key` - (Optional, List) The root key used to wrap the data encryption key for the storage ontap instance.This property will be present for storage ontap instance with an `encryption` type of`user_managed`.
Nested schema for **encryption_key**:
	* `crn` - (Required, String) The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Services Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
* `name` - (Optional, String) The name for this storage ontap instance. The name is unique across all storage ontap instances in the region.
  * Constraints: The maximum length is `40` characters. The minimum length is `7` characters. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `primary_subnet` - (Optional, List) The subnet where the primary Cloud Volumes ONTAP node is provisioned in.
Nested schema for **primary_subnet**:
	* `crn` - (Required, String) The CRN for this subnet.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this subnet. The name is unique across all subnets in the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `resource_group` - (Optional, List) The resource group for this storage ontap instance.
Nested schema for **resource_group**:
	* `id` - (Required, String) The unique identifier for this resource group.
	  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[0-9a-f]{32}$/`.
	* `name` - (Computed, String) The name for this resource group.
	  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-_ ]+$/`.
* `routing_tables` - (Optional, List) The VPC routing tables for this storage ontap instance.
  * Constraints: The maximum length is `4` items. The minimum length is `1` item.
Nested schema for **routing_tables**:
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this routing table.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this routing table.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this routing table. The name is unique across all routing tables for the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: `routing_table`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `secondary_subnet` - (Optional, List) The subnet where the secondary Cloud Volumes ONTAP node is provisioned in.
Nested schema for **secondary_subnet**:
	* `crn` - (Required, String) The CRN for this subnet.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Computed, String) The name for this subnet. The name is unique across all subnets in the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Computed, String) The resource type.
	  * Constraints: Allowable values are: `subnet`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `security_groups` - (Optional, List) The security groups for this storage ontap instance.
  * Constraints: The minimum length is `1` item.
Nested schema for **security_groups**:
	* `crn` - (Required, String) The security group's CRN.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this security group.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this security group. The name is unique across all security groups for the VPC.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
* `storage_virtual_machines` - (Required, List) The storage virtual machines for this storage ontap instance.
  * Constraints: The minimum length is `1` item.
Nested schema for **storage_virtual_machines**:
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (Required, String) The URL for this storage virtual machine.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (Required, String) The unique identifier for this storage virtual machine.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (Required, String) The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.
	  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: `storage_ontap_instance_storage_virtual_machine`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the StorageOntapInstance.
* `created_at` - (String) The date and time that the storage ontap instance was created.
* `crn` - (String) The CRN for this storage ontap instance.
  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
* `encryption` - (String) The type of encryption used on the storage ontap instance.
  * Constraints: Allowable values are: `provider_managed`, `user_managed`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `endpoints` - (List) The endpoints for this storage ontap instance.
Nested schema for **endpoints**:
	* `inter_cluster` - (List) The NetApp SnapMirror management endpoint for this storage ontap instance.
	Nested schema for **inter_cluster**:
		* `ipv4_address` - (String) The unique IP address of an endpoint.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `management` - (List) The NetApp management endpoint for this storage ontap instance. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.
	Nested schema for **management**:
		* `ipv4_address` - (String) The unique IP address of an endpoint.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
* `health_reasons` - (List) The reasons for the current storage ontap instance health_state (if any):- `cluster_down`: This storage ontap instance is unavailable as both the  primary and secondary nodes are unavailable.- `failback_unavailable`: The capability to failback is unavailable. The secondary node  continues to be available.- `failover_unavailable`: The capability to failover is unavailable. The primary  node continues to be available without any performance impact to clients.- `internal_error`: Internal error (contact IBM support).- `maintenance_in_progress`: A planned maintenance activity is in progress.- `primary_node_down`: The primary node is unavailable, and I/O has failed over to  the secondary node. Clients running in the same zone as the primary node may  experience higher access latency.- `secondary_node_down`: The secondary node is unavailable. Therefore, the capability  to failover is unavailable.
  * Constraints: The minimum length is `0` items.
Nested schema for **health_reasons**:
	* `code` - (String) A snake case string succinctly identifying the reason for this health state.
	  * Constraints: Allowable values are: `cluster_down`, `failback_unavailable`, `failover_unavailable`, `internal_error`, `maintenance_in_progress`, `primary_node_down`, `secondary_node_down`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `message` - (String) An explanation of the reason for this health state.
	* `more_info` - (String) Link to documentation about the reason for this health state.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `health_state` - (String) The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
  * Constraints: Allowable values are: `degraded`, `faulted`, `inapplicable`, `ok`.
* `href` - (String) The URL for this storage ontap instance.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_reasons` - (List) The reasons for the current `lifecycle_state` (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.
  * Constraints: The minimum length is `0` items.
Nested schema for **lifecycle_reasons**:
	* `code` - (String) A snake case string succinctly identifying the reason for this lifecycle state.
	  * Constraints: Allowable values are: `resource_suspended_by_provider`. The value must match regular expression `/^[a-z]+(_[a-z]+)*$/`.
	* `message` - (String) An explanation of the reason for this lifecycle state.
	* `more_info` - (String) Link to documentation about the reason for this lifecycle state.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `lifecycle_state` - (String) The lifecycle state of the storage ontap instance.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `storage_ontap_instance`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `vpc` - (List) The VPC this storage ontap instance resides in.
Nested schema for **vpc**:
	* `crn` - (String) The CRN for this VPC.
	  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this VPC.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this VPC. The name is unique across all VPCs in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `vpc`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `etag` - ETag identifier for StorageOntapInstance.

## Import

You can import the `ibm_is_storage_ontap_instance` resource by using `id`. The unique identifier for this storage ontap instance.

# Syntax
```
$ terraform import ibm_is_storage_ontap_instance.is_storage_ontap_instance <id>
```

# Example
```
$ terraform import ibm_is_storage_ontap_instance.is_storage_ontap_instance r134-d7cc5196-9864-48c4-82d8-3f30da41ffc5
```
