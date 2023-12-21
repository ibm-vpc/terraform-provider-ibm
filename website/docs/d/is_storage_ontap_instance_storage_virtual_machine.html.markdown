---
layout: "ibm"
page_title: "IBM : ibm_is_storage_ontap_instance_storage_virtual_machine"
description: |-
  Get information about StorageOntapInstanceStorageVirtualMachine
subcategory: "ontap"
---

# ibm_is_storage_ontap_instance_storage_virtual_machine

Provides a read-only data source to retrieve information about a StorageOntapInstanceStorageVirtualMachine. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_storage_ontap_instance_storage_virtual_machine" "is_storage_ontap_instance_storage_virtual_machine" {
	id = "id"
	storage_ontap_instance_id = "storage_ontap_instance_id"
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `id` - (Required, Forces new resource, String) The storage virtual machine identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
* `storage_ontap_instance_id` - (Required, Forces new resource, String) The storage ontap instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the StorageOntapInstanceStorageVirtualMachine.
* `active_directory` - (List) The Active Directory service this storage virtual machine is joined to.If absent, this storage virtual machine is not joined to an Active Directory service.
Nested schema for **active_directory**:
	* `administrators_group` - (String) The name of the domain group whose members have been granted administrative privileges for this storage virtual machine.
	  * Constraints: The maximum length is `256` characters. The minimum length is `1` character.
	* `dns_ips` - (List) The IP addresses of the Active Directory DNS servers or domain controllers.
	  * Constraints: The list items must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`. The maximum length is `3` items. The minimum length is `1` item.
	* `domain_name` - (String) The fully qualified domain name of the self-managed Active Directory.
	  * Constraints: The maximum length is `255` characters. The minimum length is `4` characters. The value must match regular expression `/^((?=[A-Za-z0-9-]{1,63}\\.)[A-Za-z0-9-]*\\.)+[A-Za-z]{2,63}\\.?$/`.
	* `domain_password_credential` - (List) The password credential for the Active Directory domain.
	Nested schema for **domain_password_credential**:
		* `crn` - (String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `netbios_name` - (String) The name of the Active Directory computer object that will be created for the storage virtual machine.
	  * Constraints: The maximum length is `15` characters. The minimum length is `1` character. The value must match regular expression `/^[\\S]{1,15}$/`.
	* `organizational_unit_distinguished_name` - (String) The distinguished name of the organizational unit within the self-managed Active Directory.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character.
	* `username` - (String) The username that this storage virtual machine will use when joining the Active Directory domain. This username will be the same as the username credential used in `domain_password_credential`.
	  * Constraints: The maximum length is `256` characters. The minimum length is `1` character. The value must match regular expression `/^([a-zA-Z0-9~!@#$%^&*()\\\\\\\\\\\\-_. ]+)+$/`.

* `admin_credentials` - (List) The credentials used for the administrator to access the storage virtual machine of thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.
Nested schema for **admin_credentials**:
	* `http` - (List) The security certificate credential for ONTAP REST API access for the storage virtualmachine administrator.
	Nested schema for **http**:
		* `crn` - (String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `password` - (List) The password credential for the storage virtual machine administrator.If present, this password credential is used by the storage virtual machineadministrator for both ONTAP CLI SSH access and ONTAP REST API access.If absent, the storage virtual machine is not accessible through either the ONTAP CLIor ONTAP REST API using password-based authentication.
	Nested schema for **password**:
		* `crn` - (String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
	* `ssh` - (List) The public key credential for ONTAP CLI based ssh login for the storage virtualmachine administrator.
	Nested schema for **ssh**:
		* `crn` - (String) The CRN for this credential.
		  * Constraints: The maximum length is `512` characters. The minimum length is `9` characters.
		* `resource_type` - (String) The resource type.
		  * Constraints: Allowable values are: `credential`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

* `created_at` - (String) The date and time that the storage virtual machine was created.

* `endpoints` - (List) The data and management endpoints for this storage virtual machine.
Nested schema for **endpoints**:
	* `management` - (List) The NetApp management endpoint for this storage virtual machine. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.
	Nested schema for **management**:
		* `ipv4_address` - (String) The unique IP address of an endpoint.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `nfs` - (List) The Network File System (NFS) protocol endpoint for this storage virtual machine.If absent, NFS is not enabled on this storage virtual machine.
	Nested schema for **nfs**:
		* `ipv4_address` - (String) The unique IP address of an endpoint.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `smb` - (List) The Server Message Block (SMB) protocol endpoint for this storage virtual machine.If absent, SMB is not enabled on this storage virtual machine.
	Nested schema for **smb**:
		* `ipv4_address` - (String) The unique IP address of an endpoint.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.

* `href` - (String) The URL for this storage virtual machine.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

* `lifecycle_state` - (String) The lifecycle state of the storage virtual machine.
  * Constraints: Allowable values are: `deleting`, `failed`, `pending`, `stable`, `suspended`, `updating`, `waiting`.

* `name` - (String) The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.
  * Constraints: The maximum length is `40` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.

* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `storage_ontap_instance_storage_virtual_machine`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.

