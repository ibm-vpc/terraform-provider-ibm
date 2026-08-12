---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : bare_metal_server_storage_access"
description: |-
  Manages IBM bare metal server storage access.
---

# ibm\_is_bare_metal_server_storage_access

Import the storage access configuration for an existing IBM Cloud Bare Metal Server as a read-only data source. The storage access secret is a DH-HMAC-CHAP key used to authenticate and encrypt `nvme-tcp` connections between the bare metal server and its attached block storage volumes. For more information, about managing VPC Bare Metal Server, see [About Bare Metal Servers for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-about-bare-metal-servers).

**Note:**
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example Usage

```terraform
data "ibm_is_bare_metal_server_storage_access" "example" {
  bare_metal_server = ibm_is_bare_metal_server.example.id
}
```

## Argument Reference

Review the argument references that you can specify for your data source.

- `bare_metal_server` - (Required, String) The unique identifier of the bare metal server.

## Attribute Reference

In addition to the argument reference list, you can access the following attribute references after your data source is created.

- `created_at` - (String) The date and time the storage access secret was originally created.
- `encrypted_secret` - (String, Sensitive) The storage access secret, encrypted using the SSH RSA public key identified by `public_key` and returned as a base64-encoded string. Present only when `status` is `active`.
- `public_key` - (String) The fingerprint of the SSH RSA public key used to encrypt the storage access secret.
- `rotated_at` - (String) The date and time the storage access secret was last rotated. Populated after the first rotation.
- `status` - (String) The current status of the storage access secret. Allowable values are: [ **active**, **updating** ]
