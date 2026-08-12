---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : bare_metal_server_storage_access_secret_rotate"
description: |-
  Manages IBM bare metal server storage access secret rotation.
---

# ibm\_is_bare_metal_server_storage_access_secret_rotate

Rotate the NVMe-oF (NVMe over TCP) storage access secret for a Bare Metal Server. The secret is a DH-HMAC-CHAP key used to authenticate and encrypt `nvme-tcp` connections between the bare metal server and its attached block storage volumes. The rotated secret is returned encrypted with a customer-supplied SSH RSA public key. This is a one time action resource. For multiple rotations, multiple `ibm_is_bare_metal_server_storage_access_secret_rotate` resources need to be used. For more information, about managing VPC Bare Metal Server, see [About Bare Metal Servers for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-bare-metal-servers).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

**Note:** 
The API enforces a rate limit of **one rotation per hour** per bare metal server. Attempting a second rotation within that window will be rejected. Destroying the resource only removes the Terraform state — it does not invalidate the secret on the server.

## Example Usage

In the following example, you can rotate the storage access secret of a Bare Metal Server without specifying a key:

```terraform
resource "ibm_is_bare_metal_server_storage_access_secret_rotate" "example" {
  bare_metal_server = ibm_is_bare_metal_server.example.id
}
```

## Example Usage with key

In the following example, you can rotate the storage access secret of a Bare Metal Server with an explicit SSH RSA key:

```terraform
resource "ibm_is_bare_metal_server_storage_access_secret_rotate" "example" {
  bare_metal_server = ibm_is_bare_metal_server.example.id
  key               = ibm_is_ssh_key.example.id
}
```

## Argument Reference

Review the argument references that you can specify for your resource. 


- `bare_metal_server` - (Required, Forces new resource, String) The unique identifier of the bare metal server whose storage access secret will be rotated.
- `key` - (Optional, String) The unique identifier of an SSH RSA key to use for encrypting the new storage access secret. If omitted, the server's existing SSH RSA initialization key is used. Changing this value after creation triggers a new rotation.


## Attribute Reference

In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of the bare metal server.
- `created_at` - (String) The date and time the storage access secret was originally created.
- `encrypted_secret` - (String, Sensitive) The storage access secret, encrypted using the SSH RSA public key identified by `public_key` and returned as a base64-encoded string. Present only when `status` is `active`.
- `public_key` - (String) The fingerprint of the SSH RSA public key used to encrypt the storage access secret.
- `rotated_at` - (String) The date and time the storage access secret was last rotated. Populated after the first rotation.
- `status` - (String) The current status of the storage access secret. Allowable values are: [ **active**, **updating** ]

## Import

The `ibm_is_bare_metal_server_storage_access_secret_rotate` resource can be imported by using the bare metal server ID.

**Syntax**

```sh
terraform import ibm_is_bare_metal_server_storage_access_secret_rotate.<name> <bare_metal_server_id>
```

**Example**

```sh
terraform import ibm_is_bare_metal_server_storage_access_secret_rotate.example 0716-7be49970-816a-912e-3aca-14c7fab5fda2
```

**Note:**
On import, the `key` argument and `encrypted_secret` attribute are not restored from the API as neither is returned by `GET /bare_metal_servers/{id}/storage_access`. Add `key` back to your configuration manually after importing if needed.
