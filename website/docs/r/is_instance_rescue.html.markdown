---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_instance_rescue"
description: |-
  Manages IBM Cloud VPC Instance Rescue Mode.
---

# ibm_is_instance_rescue

Provides a resource to manage instance rescue mode configuration. This resource allows you to put a virtual server instance into rescue mode for troubleshooting and recovery purposes. For more information, about instance rescue mode, see [Rescuing a virtual server instance](https://cloud.ibm.com/docs/vpc?topic=vpc-vsi_is_rescue).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Instance Rescue Workflow

Instance rescue mode allows you to troubleshoot and recover a virtual server instance that has become inaccessible. The workflow is as follows:

1. **Stop the Instance**: The instance must be in a `stopped` state before entering rescue mode.
2. **Enter Rescue Mode**: Create this resource to configure rescue mode with a rescue image, SSH keys, and optionally a rescue volume attachment.
3. **Start the Instance**: After rescue mode is configured, start the instance. It will boot from the rescue image.
4. **Perform Recovery**: Access the instance using SSH with the configured keys. The original boot volume is attached as a secondary volume for recovery operations.
5. **Exit Rescue Mode**: Delete this resource to exit rescue mode.
6. **Restart the Instance**: Stop and start the instance to boot normally from the original boot volume.

## Example Usage

### Basic Instance Rescue

```terraform
# Stop the instance first
resource "ibm_is_instance_action" "stop_instance" {
  instance = ibm_is_instance.example.id
  action   = "stop"
}

# Configure rescue mode
resource "ibm_is_instance_rescue" "example" {
  instance_id = ibm_is_instance.example.id
  
  image {
    id = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  }
  
  keys {
    id = ibm_is_ssh_key.example.id
  }
  
  rescue_volume_attachment {
    name                             = "rescue-volume-attachment"
    delete_volume_on_instance_delete = true
    volume {
      name    = "rescue-volume"
      profile = "general-purpose"
    }
  }
  
  depends_on = [ibm_is_instance_action.stop_instance]
}

# Start the instance in rescue mode
resource "ibm_is_instance_action" "start_rescue" {
  instance = ibm_is_instance.example.id
  action   = "start"
  
  depends_on = [ibm_is_instance_rescue.example]
}
```

### Rescue Mode with Multiple SSH Keys and User Data

```terraform
resource "ibm_is_instance_rescue" "example" {
  instance_id = ibm_is_instance.example.id
  
  image {
    id = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
  }
  
  keys {
    id = ibm_is_ssh_key.key1.id
  }
  
  keys {
    id = ibm_is_ssh_key.key2.id
  }
  
  user_data = <<-EOT
    #!/bin/bash
    echo "Rescue mode initialized"
    # Add your recovery scripts here
  EOT
  
  rescue_volume_attachment {
    name                             = "rescue-volume-attachment"
    delete_volume_on_instance_delete = true
    volume {
      name    = "rescue-volume"
      profile = "general-purpose"
    }
  }
}
```

### Rescue Mode with Custom Volume Configuration

```terraform
resource "ibm_is_instance_rescue" "example" {
  instance_id = ibm_is_instance.example.id
  
  image {
    id = data.ibm_is_image.rescue_image.id
  }
  
  keys {
    id = ibm_is_ssh_key.example.id
  }
  
  rescue_volume_attachment {
    name                             = "rescue-volume-attachment"
    delete_volume_on_instance_delete = true
    volume {
      name           = "rescue-volume"
      profile        = "custom"
      capacity       = 100
      iops           = 3000
      encryption_key = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:5437653b-c4b1-447f-9646-b2a2a4cd6179"
      resource_group = data.ibm_resource_group.example.id
      user_tags      = ["env:rescue", "owner:ops"]
    }
  }
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `instance_id` - (Required, Forces new resource, String) The virtual server instance identifier.
  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.

* `image` - (Required, Forces new resource, List) The image to use for rescuing the instance.

  Nested schema for **image**:
  * `id` - (Optional, String) The unique identifier for this image.
    * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
  * `crn` - (Optional, String) The CRN for this image.
  * `href` - (Optional, String) The URL for this image.

* `keys` - (Optional, Forces new resource, List) The public SSH keys used at initialization for the rescue volume.
  * Constraints: The maximum length is `10` items. The minimum length is `0` items.

  Nested schema for **keys**:
  * `id` - (Optional, String) The unique identifier for this key.
    * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
  * `crn` - (Optional, String) The CRN for this key.
  * `href` - (Optional, String) The URL for this key.
  * `fingerprint` - (Optional, String) The fingerprint for this key.

* `user_data` - (Optional, Forces new resource, String) User data to be made available when setting up the rescue volume.

* `rescue_volume_attachment` - (Required, Forces new resource, List) The rescue volume attachment for this instance.

  Nested schema for **rescue_volume_attachment**:
  * `delete_volume_on_instance_delete` - (Optional, Boolean) Indicates whether to delete the rescue volume when the instance is deleted. Default: `true`.
  * `name` - (Optional, String) The name for this volume attachment. The name is unique across all volume attachments on the instance. If unspecified, the name will be a hyphenated list of randomly-selected words.
    * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
  * `volume` - (Optional, List) The volume prototype or reference for the rescue volume.
  
    Nested schema for **volume**:
    * `id` - (Optional, String) The unique identifier for an existing volume.
    * `name` - (Optional, String) The name for this volume. The name is unique across all volumes in the region.
      * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
    * `profile` - (Optional, String) The profile name for this volume (e.g., `general-purpose`, `5iops-tier`, `10iops-tier`, `custom`).
    * `capacity` - (Optional, Integer) The capacity of the volume in gigabytes. Must be at least the image minimum_provisioned_size.
      * Constraints: Range: 10-250 GB for storage_generation 1, 10-32000 GB for storage_generation 2.
    * `iops` - (Optional, Integer) The maximum I/O operations per second (IOPS) for the volume. Applicable only for custom profile volumes.
    * `encryption_key` - (Optional, String) The CRN of the Key Management Service root key to use for encryption.
    * `resource_group` - (Optional, String) The resource group ID for this volume.
    * `user_tags` - (Optional, Set of String) User tags for this volume.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the instance rescue configuration (same as `instance_id`).

* `password` - (List) The administrator password for rescue mode access.

  Nested schema for **password**:
  * `encrypted_password` - (String) The administrator password at rescue, encrypted using `encryption_key`, and returned base64-encoded.
    * Constraints: The maximum length is `172` characters. The minimum length is `4` characters. The value must match regular expression `/^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4})$/`.
  * `encryption_key` - (List) The public SSH key used to encrypt the administrator password.
  
    Nested schema for **encryption_key**:
    * `crn` - (String) The CRN for this key.
    * `fingerprint` - (String) The fingerprint for this key.
    * `href` - (String) The URL for this key.
    * `id` - (String) The unique identifier for this key.
    * `name` - (String) The name for this key.

* `rescue_volume_attachment` - (List) The rescue volume attachment details.

  Nested schema for **rescue_volume_attachment**:
  * `href` - (String) The URL for this volume attachment.
  * `id` - (String) The unique identifier for this volume attachment.
  * `device` - (List) The configuration for the volume as a device in the instance operating system.
  
    Nested schema for **device**:
    * `id` - (String) A unique identifier for the device which is exposed to the instance operating system.

## Import

You can import the `ibm_is_instance_rescue` resource by using `instance_id`.

**Syntax**

```
$ terraform import ibm_is_instance_rescue.example <instance_id>
```

**Example**

```
$ terraform import ibm_is_instance_rescue.example r006-12345678-1234-1234-1234-123456789012
