---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : Image"
description: |-
  Manages IBM Cloud Infrastructure Images.
---

# ibm\_is_image

Import the details of an existing IBM Cloud Infrastructure image as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```terraform

data "ibm_is_image" "ds_image" {
  name = "centos-7.x-amd64"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the image.
* `visibility` - (Optional, string) The visibility of the image. Accepted values are `public` or `private`.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this image.
* `crn` - The CRN for this image.
* `checksum` - The SHA256 Checksum for this image
* `os` - The name of the operating system.
* `status` - The status of this image.
* `architecture` - The architecture for this image.
* `encryption` - The type of encryption used on the image.
* `encryption_key` - The CRN of the Key Protect Root Key or Hyper Protect Crypto Service Root Key for this resource.

