---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : instance reinitialize"
description: |-
  Manages IBM instance reinitialization.
---

# ibm_is_instance_reinitialize

Reinitialize a virtual server instance for VPC with a new image, boot volume, or snapshot. Reinitialization replaces the instance's boot volume and resets the operating system while preserving the instance ID, network interfaces, and other configuration. For more information, see [Reinitializing a virtual server instance](https://cloud.ibm.com/docs/vpc?topic=vpc-managing-virtual-server-instances#reinitialize-vs).

**Important Notes:**
- The instance must be stopped before reinitialization (use `auto_stop = true` for automatic handling)
- The old boot volume will be deleted and replaced with a new one
- All data on the boot volume will be lost
- SSH keys, user data, and trusted profile settings can be updated during reinitialization
- The instance will automatically start after reinitialization completes

**Note:** 
VPC infrastructure services are regional specific based endpoints. By default, targets `us-south`. Please make sure to target the right region in the provider block as shown in the `provider.tf` file if VPC service is created in a region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example Usage

### Basic Reinitialization with Image

```terraform
resource "ibm_is_vpc" "example" {
  name = "example-vpc"
}

resource "ibm_is_subnet" "example" {
  name            = "example-subnet"
  vpc             = ibm_is_vpc.example.id
  zone            = "us-south-1"
  ipv4_cidr_block = "10.240.0.0/24"
}

resource "ibm_is_ssh_key" "example" {
  name       = "example-ssh"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
}

resource "ibm_is_instance" "example" {
  name    = "example-instance"
  image   = "r006-original-image-id"
  profile = "bx2-2x8"

  primary_network_interface {
    subnet = ibm_is_subnet.example.id
  }

  vpc  = ibm_is_vpc.example.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.example.id]

  lifecycle {
    ignore_changes = [boot_volume]
  }
}

resource "ibm_is_instance_reinitialize" "example" {
  instance   = ibm_is_instance.example.id
  image      = "r006-new-image-id"
  keys       = [ibm_is_ssh_key.example.id]
  auto_stop  = true
  auto_start = true
}
```

### Reinitialization with Triggers

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  instance   = ibm_is_instance.example.id
  image      = var.current_image_id
  keys       = [ibm_is_ssh_key.example.id]
  auto_stop  = true
  auto_start = true

  triggers = {
    image_version = var.image_version
    deployment_id = var.deployment_id
  }
}
```

### Reinitialization with User Data

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  instance   = ibm_is_instance.example.id
  image      = "r006-new-image-id"
  keys       = [ibm_is_ssh_key.example.id]
  user_data  = file("${path.module}/cloud-init.yaml")
  auto_stop  = true
  auto_start = true
}
```

### Reinitialization with Boot Volume Configuration

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  instance   = ibm_is_instance.example.id
  image      = "r006-new-image-id"
  keys       = [ibm_is_ssh_key.example.id]
  auto_stop  = true

  boot_volume_attachment {
    name     = "my-boot-volume"
    capacity = 100
    profile  = "general-purpose"
  }
}
```

### Reinitialization with Existing Boot Volume

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  instance    = ibm_is_instance.example.id
  boot_volume = "r006-existing-volume-id"
  keys        = [ibm_is_ssh_key.example.id]
  auto_stop   = true
}
```

### Reinitialization with Snapshot

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  instance      = ibm_is_instance.example.id
  boot_snapshot = "r006-snapshot-id"
  keys          = [ibm_is_ssh_key.example.id]
  auto_stop     = true
}
```

### Conditional Reinitialization with Count

```terraform
resource "ibm_is_instance_reinitialize" "example" {
  count      = var.trigger_reinit ? 1 : 0
  instance   = ibm_is_instance.example.id
  image      = var.new_image_id
  keys       = [ibm_is_ssh_key.example.id]
  auto_stop  = true
}
```

### Fleet Management

```terraform
resource "ibm_is_instance_reinitialize" "fleet" {
  for_each   = toset(var.instance_ids)
  instance   = each.value
  image      = var.fleet_image_id
  keys       = var.fleet_ssh_keys
  auto_stop  = true

  triggers = {
    deployment_id = var.deployment_id
  }
}
```

## Argument Reference

Review the argument references that you can specify for your resource.

### Required Arguments

- `instance` - (Required, Forces new resource, String) The instance identifier to reinitialize.

### Boot Source (One Required)

You must specify exactly one of the following boot sources:

- `image` - (Optional, Forces new resource, String) The image ID to use for reinitialization. Conflicts with `boot_volume` and `boot_snapshot`.
- `boot_volume` - (Optional, Forces new resource, String) An existing boot volume ID to attach. Conflicts with `image` and `boot_snapshot`.
- `boot_snapshot` - (Optional, Forces new resource, String) A snapshot ID to create the boot volume from. Conflicts with `image` and `boot_volume`.

### Optional Configuration Arguments

- `keys` - (Optional, Forces new resource, List of Strings) List of SSH key IDs to configure for the instance. These keys will be available for authentication after reinitialization.
- `user_data` - (Optional, Forces new resource, String) User data to make available when setting up the virtual server instance. This is typically used for cloud-init configuration.
- `default_trusted_profile` - (Optional, Forces new resource, List) The default trusted profile configuration. Maximum items: 1.

  Nested `default_trusted_profile`:
  - `target` - (Required, String) The trusted profile ID.
  - `auto_link` - (Optional, Boolean) If set to `true`, the system will create a link to the specified target during instance creation. Default: `false`.

- `boot_volume_attachment` - (Optional, Forces new resource, List) Boot volume configuration when using image-based reinitialization. Maximum items: 1.

  Nested `boot_volume_attachment`:
  - `name` - (Optional, String) The name for the boot volume.
  - `capacity` - (Optional, Integer) The capacity of the boot volume in GB.
  - `profile` - (Optional, String) The profile to use for the boot volume. Default: `general-purpose`.
  - `encryption_key` - (Optional, String) The CRN of the encryption key to use for encrypting the boot volume.

### Automation Arguments

- `auto_stop` - (Optional, Boolean) If set to `true`, the instance will be automatically stopped before reinitialization if it is running. If `false` and the instance is running, reinitialization will fail. Default: `false`.
- `auto_start` - (Optional, Boolean) If set to `true`, the instance will automatically start after reinitialization (this is the default API behavior). Default: `true`.
- `triggers` - (Optional, Forces new resource, Map of Strings) Arbitrary map of values that, when changed, will trigger reinitialization. This is useful for implementing trigger-based reinitialization workflows.

### Timeouts

- `create` - (Default 30 minutes) Used for reinitialization operation.
- `update` - (Default 30 minutes) Used when triggers change.
- `delete` - (Default 10 minutes) Used for cleanup (no-op).

## Attribute Reference

In addition to all argument reference lists, you can access the following attribute references after your resource is created.

- `id` - (String) The unique identifier of the instance (same as the instance ID).
- `status` - (String) The current status of the instance after reinitialization.
- `boot_volume_attachment_id` - (String) The ID of the new boot volume attachment created during reinitialization.
- `reinitialized_at` - (String) Timestamp of when the reinitialization was performed (RFC3339 format).

## Import

The `ibm_is_instance_reinitialize` resource can be imported using the instance ID.

```terraform
import {
  to = ibm_is_instance_reinitialize.example
  id = "<instance_id>"
}
```

Using `terraform import`:

```console
$ terraform import ibm_is_instance_reinitialize.example <instance_id>
```

## Usage Notes

### Lifecycle Management

It is recommended to add a lifecycle rule to the instance resource to ignore boot volume changes:

```terraform
resource "ibm_is_instance" "example" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [boot_volume]
  }
}
```

This prevents Terraform from detecting the boot volume change as a drift after reinitialization.

### Trigger Patterns

The `triggers` argument provides flexible reinitialization control:

**Time-based triggers:**
```terraform
triggers = {
  timestamp = timestamp()  # Triggers on every apply
}
```

**Version-based triggers:**
```terraform
triggers = {
  image_version = var.image_version  # Triggers when version changes
}
```

**Multiple conditions:**
```terraform
triggers = {
  image_version = var.image_version
  config_hash   = md5(file("config.yaml"))
}
```

### Data Loss Warning

Reinitialization **permanently deletes** the existing boot volume and all its data. Ensure you have backups of any important data before reinitializing an instance.

### Instance State Requirements

- The instance must be in a stopped state before reinitialization
- Use `auto_stop = true` to automatically handle this requirement
- If `auto_stop = false` and the instance is running, the operation will fail

### Capacity Considerations

After reinitialization, the instance must be rescheduled on available infrastructure. This may fail for:
- High-demand instance profiles (e.g., GPU instances)
- Capacity-constrained zones
- Dedicated host placements

## Related Resources

- [ibm_is_instance](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_instance)
- [ibm_is_instance_action](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_instance_action)
- [ibm_is_image](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_image)
- [ibm_is_volume](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_volume)
- [ibm_is_snapshot](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_snapshot)