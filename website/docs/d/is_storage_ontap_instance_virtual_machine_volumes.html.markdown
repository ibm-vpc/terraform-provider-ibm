---
layout: "ibm"
page_title: "IBM : ibm_is_storage_ontap_instance_virtual_machine_volumes"
description: |-
  Get information about StorageOntapInstanceStorageVirtualMachineVolumeCollection
subcategory: "ontap"
---

# ibm_is_storage_ontap_instance_virtual_machine_volumes

Provides a read-only data source to retrieve information about a StorageOntapInstanceStorageVirtualMachineVolumeCollection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_storage_ontap_instance_virtual_machine_volumes" "is_storage_ontap_instance_virtual_machine_volumes" {
	storage_ontap_instance_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_ontap_instance_id
	storage_virtual_machine_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_virtual_machine_id
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `name` - (Optional, String) Filters the collection to resources with a `name` property matching the exact specified name.
* `storage_ontap_instance_id` - (Required, Forces new resource, String) The storage ontap instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `storage_virtual_machine_id` - (Required, Forces new resource, String) The storage virtual machine identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the StorageOntapInstanceStorageVirtualMachineVolumeCollection.
* `first` - (List) A link to the first page of resources.
Nested schema for **first**:
	* `href` - (String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `limit` - (Integer) The maximum number of resources that can be returned by the request.
  * Constraints: The maximum value is `100`. The minimum value is `1`.

* `next` - (List) A link to the next page of resources. This property is present for all pagesexcept the last page.
Nested schema for **next**:
	* `href` - (String) The URL for a page of resources.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `total_count` - (Integer) The total number of resources across all pages.
  * Constraints: The minimum value is `0`.

* `volumes` - (List) Collection of storage volumes.
Nested schema for **volumes**:
	* `capacity` - (Integer) The capacity of the storage volume (in gigabytes).
	  * Constraints: The maximum value is `16000`. The minimum value is `10`.
	* `cifs_share` - (List) The named access point that enables CIFS clients to view, browse, and manipulatefiles on this storage volumeThis will be present when `security_style` is `mixed` or `windows`.
	Nested schema for **cifs_share**:
		* `access_control_list` - (List) The access control list for the CIFS share.
		Nested schema for **access_control_list**:
			* `permission` - (String) The permission granted to users matching this access control list entry.
			  * Constraints: Allowable values are: `change`, `full_control`, `no_access`, `read`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
			* `users` - (List) The users matching this access control list entry.
			  * Constraints: The list items must match regular expression `/^([a-zA-Z0-9~!@#$%^&*()\\\\\\\\\\\\-_. ]+)+$/`. The maximum length is `100` items. The minimum length is `1` item.
		* `mount_path` - (String) The SMB/CIFS mount point for the storage volume.
		  * Constraints: The maximum length is `256` characters. The minimum length is `1` character.
		* `name` - (String) The share name registered in Active Directory that SMB/CIFS clients use to mount the share. The name is unique within the Active Directory domain.
		  * Constraints: The maximum length is `80` characters. The minimum length is `1` character. The value must match regular expression `/^[!@#$%&'_\\\\\\-.~(){}a-zA-Z0-9][!@#$%&'_\\\\\\-.~(){}a-zA-Z0-9 ]{0,79}$/`.
	* `created_at` - (String) The date and time that the storage volume was created.
	* `enable_storage_efficiency` - (Deprecated, Boolean) Indicates whether storage efficiency is enabled for the storage volume.If `true`, data-deduplication, compression and other efficiencies for space-management are enabled for this volume.
	* `export_policy` - (List) The NFS export policy for the storage volume.This will be present when `security_style` is `mixed` or `unix`.
	Nested schema for **export_policy**:
		* `mount_path` - (String) The NFS mount point for the storage volume.
		  * Constraints: The maximum length is `256` characters. The minimum length is `1` character.
		* `rules` - (List) The NFS export policy rules for this storage volume.Only NFS clients included in the rules will access the volume, and only according to the specified access controls and NFS protocol versions.
		  * Constraints: The maximum length is `100` items. The minimum length is `1` item.
		Nested schema for **rules**:
			* `access_control` - (String) The access control that is provided to clients that match this rule.
			  * Constraints: Allowable values are: `none`, `read_only`, `read_write`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
			* `clients` - (List) The clients that match this rule. Every client in the list of `clients` is unique.
			Nested schema for **clients**:
				* `address` - (String) The IP address of the NFS client.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
				  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
				* `cidr_block` - (String) The CIDR block containing IP addresses of the NFS clients. The CIDR block `0.0.0.0/0` matches all client addresses.This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.
				  * Constraints: The maximum length is `18` characters. The minimum length is `9` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`.
				* `domain_suffix` - (String) The domain names suffixes of the NFS clients.
				  * Constraints: The maximum length is `255` characters. The minimum length is `4` characters. The value must match regular expression `/^(\\.)([A-Za-z0-9-]*\\.)+[A-Za-z]{2,63}\\.?$/`.
				* `hostname` - (String) The hostname of the NFS client.
				  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9]$/`.
			* `index` - (Integer) The unique index for this rule. Rules are applied from lowest to highest index.
			* `is_superuser` - (Boolean) If `true`, clients matching this rule that request super-user access are honored. Otherwise, clients are mapped to the anonymous user.
			* `nfs_version` - (List) The NFS versions that is provided to clients that match this rule.
			  * Constraints: Allowable list items are: `nfs3`, `nfs4`.
	* `health_reasons` - (List) The reasons for the current storage volume health_state (if any):- `primary_node_down`: The storage volume is experiencing higher latency due to  the primary node being unavailable, and I/O being routed to the secondary node.- `volume_unavailable`: The storage volume is unavailable as both the primary and secondary nodes are down.- `internal_error`: Internal error (contact IBM support).
	  * Constraints: The minimum length is `0` items.
	Nested schema for **health_reasons**:
		* `code` - (String) A snake case string succinctly identifying the reason for this health state.
		  * Constraints: Allowable values are: `internal_error`, `primary_node_down`, `volume_unavailable`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `message` - (String) An explanation of the reason for this health state.
		* `more_info` - (String) Link to documentation about the reason for this health state.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `health_state` - (String) The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
	  * Constraints: Allowable values are: `degraded`, `faulted`, `inapplicable`, `ok`.
	* `href` - (String) The URL for this storage volume.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this storage volume.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `junction_path` - (String) The path clients can use to mount or access this storage volume. The path is case insensitive and is unique within a storage virtual machine.
	  * Constraints: The maximum length is `255` characters. The value must match regular expression `/^\/([^?\\\\*#><|"]*)[^?\/\\\\*#><|"]{1,255}$/`.
	* `lifecycle_state` - (String) The lifecycle state of the storage volume.
	  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.
	* `name` - (String) The name for this storage volume. The name is unique across all storage volumes in the storage virtual machine.
	  * Constraints: The maximum length is `203` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z][0-9a-zA-Z_]{0,202}$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `storage_ontap_instance_storage_virtual_machine_volume`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `security_style` - (String) The security style for the storage volume:- `unix`: NFS clients can access the storage volume.- `windows`: SMB/CIFS clients can access the storage volume.- `mixed`: Both SMB/CIFS and NFS clients can access the storage volume.- `none`: No clients can access the volume.
	  * Constraints: Allowable values are: `mixed`, `none`, `unix`, `windows`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `storage_efficiency` - (String) The storage efficiency mode used for this storage volume.- `disabled`: storage efficiency methods will not be used- `enabled`: data-deduplication, compression and other methods will be used.
	  * Constraints: Allowable values are: `disabled`, `enabled`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `type` - (String) The type of the storage volume.
	  * Constraints: Allowable values are: `data_protection`, `read_write`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

