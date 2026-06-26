---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_instance_rescue"
description: |-
  Get information about IBM Cloud VPC Instance Rescue Mode configuration.
---

# ibm_is_instance_rescue

Provides a read-only data source to retrieve information about an instance's rescue mode configuration. For more information, about instance rescue mode, see [Rescuing a virtual server instance](https://cloud.ibm.com/docs/vpc?topic=vpc-vsi_is_rescue).

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
data "ibm_is_instance_rescue" "example" {
  instance_id = ibm_is_instance.example.id
}

output "rescue_image_id" {
  value = data.ibm_is_instance_rescue.example.image[0].id
}

output "rescue_volume_attachment_id" {
  value = data.ibm_is_instance_rescue.example.rescue_volume_attachment[0].id
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `instance_id` - (Required, String) The virtual server instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the instance rescue configuration (same as `instance_id`).

* `image` - (List) The image used for rescuing the instance.

  Nested schema for **image**:
  * `crn` - (String) The CRN for this image.
  * `deleted` - (List) If present, this property indicates the referenced resource has been deleted.
  
    Nested schema for **deleted**:
    * `more_info` - (String) A link to documentation about deleted resources.
  
  * `href` - (String) The URL for this image.
  * `id` - (String) The unique identifier for this image.
  * `name` - (String) The name for this image. The name is unique across all images in the region.
  * `remote` - (List) If present, this property indicates that the resource associated with this reference is remote.
  
    Nested schema for **remote**:
    * `account` - (List) If present, this property indicates that the referenced resource is remote to this account.
    
      Nested schema for **account**:
      * `id` - (String) The unique identifier for this account.
      * `resource_type` - (String) The resource type.
    
    * `region` - (List) If present, this property indicates that the referenced resource is remote to this region.
    
      Nested schema for **region**:
      * `href` - (String) The URL for this region.
      * `name` - (String) The globally unique name for this region.
  
  * `resource_type` - (String) The resource type.

* `keys` - (List) The public SSH keys used at initialization for the rescue volume.

  Nested schema for **keys**:
  * `crn` - (String) The CRN for this key.
  * `deleted` - (List) If present, this property indicates the referenced resource has been deleted.
  
    Nested schema for **deleted**:
    * `more_info` - (String) A link to documentation about deleted resources.
  
  * `fingerprint` - (String) The fingerprint for this key.
  * `href` - (String) The URL for this key.
  * `id` - (String) The unique identifier for this key.
  * `name` - (String) The name for this key. The name is unique across all keys in the region.

* `password` - (List) The administrator password for rescue mode access (Windows instances only).

  Nested schema for **password**:
  * `encrypted_password` - (String) The administrator password at rescue, encrypted using `encryption_key`, and returned base64-encoded.
  * `encryption_key` - (List) The public SSH key used to encrypt the administrator password.
  
    Nested schema for **encryption_key**:
    * `crn` - (String) The CRN for this key.
    * `deleted` - (List) If present, this property indicates the referenced resource has been deleted.
    
      Nested schema for **deleted**:
      * `more_info` - (String) A link to documentation about deleted resources.
    
    * `fingerprint` - (String) The fingerprint for this key.
    * `href` - (String) The URL for this key.
    * `id` - (String) The unique identifier for this key.
    * `name` - (String) The name for this key.

* `rescue_volume_attachment` - (List) The rescue volume attachment for this instance.

  Nested schema for **rescue_volume_attachment**:
  * `deleted` - (List) If present, this property indicates the referenced resource has been deleted.
  
    Nested schema for **deleted**:
    * `more_info` - (String) A link to documentation about deleted resources.
  
  * `device` - (List) The configuration for the volume as a device in the instance operating system.
  
    Nested schema for **device**:
    * `id` - (String) A unique identifier for the device which is exposed to the instance operating system.
  
  * `href` - (String) The URL for this volume attachment.
  * `id` - (String) The unique identifier for this volume attachment.
  * `name` - (String) The name for this volume attachment.
  * `volume` - (List) The attached volume.
  
    Nested schema for **volume**:
    * `crn` - (String) The CRN for this volume.
    * `deleted` - (List) If present, this property indicates the referenced resource has been deleted.
    
      Nested schema for **deleted**:
      * `more_info` - (String) A link to documentation about deleted resources.
    
    * `href` - (String) The URL for this volume.
    * `id` - (String) The unique identifier for this volume.
    * `name` - (String) The name for this volume.
    * `resource_type` - (String) The resource type.
