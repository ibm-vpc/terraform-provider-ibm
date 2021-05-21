---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : instance"
description: |-
  Manages IBM IS Instance.
---

# ibm\_is_instance

Provides a instance resource. This allows instance to be created, updated, and cancelled.


## Example Usage

```hcl
resource "ibm_is_vpc" "testacc_vpc" {
  name = "testvpc"
}

resource "ibm_is_subnet" "testacc_subnet" {
  name            = "testsubnet"
  vpc             = ibm_is_vpc.testacc_vpc.id
  zone            = "us-south-1"
  ipv4_cidr_block = "10.240.0.0/24"
}

resource "ibm_is_ssh_key" "testacc_sshkey" {
  name       = "testssh"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
}

resource "ibm_is_instance" "testacc_instance" {
  name    = "testinstance"
  image   = "7eb4e35b-4257-56f8-d7da-326d85452591"
  profile = "bc1-2x8"

  primary_network_interface {
    subnet = ibm_is_subnet.testacc_subnet.id
    primary_ipv4_address = "10.240.0.6"
    allow_ip_spoofing = true
  }

  network_interfaces {
    name   = "eth1"
    subnet = ibm_is_subnet.testacc_subnet.id
    allow_ip_spoofing = false
  }

  vpc  = ibm_is_vpc.testacc_vpc.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.testacc_sshkey.id]

  //User can configure timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

```

Here is an example of creating virtual server instance with security group, security group rule. Here, the security group, security group rule, and virtual server instance must be created sequentially as security group rule depends on security group creation and virtual server instance depends on security group, security group rule creation. The sequential creation of resources like security group, security rule, virtual server instance is achieved using "depends_on" attribute. You can find more information about depends_on attribute in [terraform documentation](https://www.terraform.io/docs/configuration/resources.html). Creating security group, security group rule, virtual server instance without depends_on attribute will create the resources in parallel and virtual server instance creation may fail with "Error: The security group to attach to is not available" as security group or security group rule creation is not complete and security group may be in Pending state.         

```hcl

resource "ibm_is_vpc" "testacc_vpc" {
    name = "test"
}

resource "ibm_is_security_group" "testacc_security_group" {
    name = "test"
    vpc = ibm_is_vpc.testacc_vpc.id
}

resource "ibm_is_security_group_rule" "testacc_security_group_rule_all" {
    group = ibm_is_security_group.testacc_security_group.id
    direction = "inbound"
    remote = "127.0.0.1"
    depends_on = [ibm_is_security_group.testacc_security_group]
 }

 resource "ibm_is_security_group_rule" "testacc_security_group_rule_icmp" {
    group = ibm_is_security_group.testacc_security_group.id
    direction = "inbound"
    remote = "127.0.0.1"
    icmp {
        code = 20
        type = 30
    }
    depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_all]

 }

 resource "ibm_is_security_group_rule" "testacc_security_group_rule_udp" {
    group = ibm_is_security_group.testacc_security_group.id
    direction = "inbound"
    remote = "127.0.0.1"
    udp {
        port_min = 805
        port_max = 807
    }
    depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_icmp]
 }

 resource "ibm_is_security_group_rule" "testacc_security_group_rule_tcp" {
    group = ibm_is_security_group.testacc_security_group.id
    direction = "outbound"
    remote = "127.0.0.1"
    tcp {
        port_min = 8080
        port_max = 8080
    }
    depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_udp]
 }

resource "ibm_is_instance" "testacc_instance" {
  name    = "testinstance"
  image   = "7eb4e35b-4257-56f8-d7da-326d85452591"
  profile = "bc1-2x8"

  primary_network_interface {
    subnet = ibm_is_subnet.testacc_subnet.id
    security_groups = [ibm_is_security_group.testacc_security_group.id]
  }

  vpc  = ibm_is_vpc.testacc_vpc.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.testacc_sshkey.id]
  depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_tcp]

  //User can configure timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

data "ibm_resource_group" "default" {
  name = "Default" ///give your resource grp
}

resource "ibm_is_dedicated_host_group" "dh_group01" {
  family = "compute"
  class = "cx2"
  zone = "us-south-1"
  name = "my-dh-group-01"
  resource_group = data.ibm_resource_group.default.id
}

resource "ibm_is_dedicated_host" "is_dedicated_host" {
  profile = "bx2d-host-152x608"
  name = "my-dedicated-host-01"
	host_group = ibm_is_dedicated_host_group.dh_group01.id
  resource_group = data.ibm_resource_group.default.id
}

// Example to provision instance in a dedicated host
resource "ibm_is_instance" "testacc_instance1" {
  name    = "testinstance1"
  image   = "7eb4e35b-4257-56f8-d7da-326d85452591"
  profile = "cx2-2x4"

  primary_network_interface {
    subnet = ibm_is_subnet.testacc_subnet.id
    security_groups = [ibm_is_security_group.testacc_security_group.id]
  }
  dedicated_host = ibm_is_dedicated_host.is_dedicated_host.id
  vpc  = ibm_is_vpc.testacc_vpc.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.testacc_sshkey.id]
  depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_tcp]

  //User can configure timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

// Example to provision instance in a dedicated host that belongs to the provided dedicated host group 
resource "ibm_is_instance" "testacc_instance2" {
  name    = "testinstance2"
  image   = "7eb4e35b-4257-56f8-d7da-326d85452591"
  profile = "cx2-2x4"

  primary_network_interface {
    subnet = ibm_is_subnet.testacc_subnet.id
    security_groups = [ibm_is_security_group.testacc_security_group.id]
  }
  dedicated_host_group = ibm_is_dedicated_host_group.dh_group01.id
  vpc  = ibm_is_vpc.testacc_vpc.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.testacc_sshkey.id]
  depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_tcp]

  //User can configure timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}


```  

## Timeouts

ibm_is_instance provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default 30 minutes) Used for creating Instance.
* `update` - (Default 30 minutes) Used for updating Instance or while attaching it with volume attachments or interfaces.
* `delete` - (Default 30 minutes) Used for deleting Instance.

## Argument Reference

The following arguments are supported:

* `name` - (Optional, string) The instance name.
* `vpc` - (Required, Forces new resource, string) The vpc id. 
* `zone` - (Required, Forces new resource, string) Name of the zone. 
* `profile` - (Required, string) The profile name. 
  * * Updating profile requires instance to be in stopped status, running instance will be stopped on update profile action.
* `dedicated_host` - (Optional, string, ForceNew) The placement restrictions to use for the virtual server instance. Unique Identifier of the Dedicated Host where the instance will be placed
* `dedicated_host_group` - (Optional, string, ForceNew) The placement restrictions to use for the virtual server instance. Unique Identifier of the Dedicated Host Group where the instance will be placed
* `image` - (Optional, string) ID of the image.
    * **NOTE**: Required with `zone`, `primary_network_interface` and one of `boot_volume.source_snapshot`, `image`, `source_template` is mandatory. 
* `boot_volume` - (Optional, list) A block describing the boot volume of this instance.  
`boot_volume` block have the following structure:
  * `name` - (Optional, string) The name of the boot volume.
  * `encryption` -(Optional, string) The encryption of the boot volume.
  * `source_snapshot` -(Optional, string) Snapshot ID of the boot volume.
    * **NOTE**: Required with `zone`, `primary_network_interface`, conflicts with `volume` and one of `boot_volume.source_snapshot`, `image`, `source_template` is mandatory.  
  * `volume` -(Optional, string) ID the boot volume.
    * **NOTE**: Required with `zone`, `primary_network_interface`, conflicts with `source_snapshot` and one of `boot_volume.volume`, `image`, `source_template` is mandatory.    
* `keys` - (Required, list) Comma separated IDs of ssh keys.  
* `primary_network_interface` - (Optional, list) A nested block describing the primary network interface of this instance. We can have only one primary network interface.
Nested `primary_network_interface` block have the following structure:
  * `name` - (Optional, string) The name of the network interface.
  * `port_speed` - (Deprecated, int) Speed of the network interface.
  * `primary_ipv4_address` - (Optional, Forces new resource, string) The IPV4 address of the interface
  * `subnet` -  (Required, string) ID of the subnet.
  * `security_groups` - (Optional, list) Comma separated IDs of security groups.
  * `allow_ip_spoofing` - (Optional, bool) Indicates whether IP spoofing is allowed on this interface. If false, IP spoofing is prevented on this interface. If true, IP spoofing is allowed on this interface.
* `network_interfaces` - (Optional, Forces new resource, list) A nested block describing the additional network interface of this instance.
Nested `network_interfaces` block have the following structure:
  * `name` - (Optional, string) The name of the network interface.
  * `primary_ipv4_address` - (Optional, Forces new resource, string) The IPV4 address of the interface
  * `subnet` -  (Required, string) ID of the subnet.
  * `security_groups` - (Optional, list) Comma separated IDs of security groups.
  * `allow_ip_spoofing` - (Optional, bool) Indicates whether IP spoofing is allowed on this interface. If false, IP spoofing is prevented on this interface. If true, IP spoofing is allowed on this interface.
* `volumes` - (Optional, list) Comma separated IDs of volumes. 
* `snapshots` - (Optional, list) Comma separated IDs of snapshots to be attached as data volume.
* `boot_volume_type` - (Optional, string) Type of the boot volume.
  * One of [`image`, `template`, `volume`]
  * Default : `image`
* `source_template` - (Optional, string) ID of the source template.  
    * **NOTE**: One of `boot_volume.source_snapshot`, `image`, `source_template` is mandatory.   
* `auto_delete_volume` - (Optional, string) If set to true, automatically deletes volumes attached to the instance.  
**Note** Setting this argument may bring some inconsistency in volume resources since the volumes will be destroyed along with instances.
* `user_data` - (Optional, string) User data to transfer to the server instance.
* `resource_group` - (Optional, Forces new resource, string) The resource group ID for this instance.
* `tags` - (Optional, array of strings) Tags associated with the instance.
* `force_recovery_time` - (Optional, int) Define timeout (in minutes), to force the is_instance to recover from a perpetual "starting" state, during provisioning; similarly, to force the is_instance to recover from a perpetual "stopping" state, during deprovisioning.  **Note**: the force_recovery_time is used to retry multiple times until timeout.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the instance.
* `memory` - Memory of the instance.
* `status` - Status of the instance.
* `vcpu` - A nested block describing the VCPU configuration of this instance.
Nested `vcpu` blocks have the following structure:
  * `architecture` - The architecture of the instance.
  * `count` - The number of VCPUs assigned to the instance.
* `gpu` - A nested block describing the gpu of this instance.
Nested `gpu` blocks have the following structure:
  * `cores` - The cores of the gpu.
  * `count` - Count of the gpu.
  * `manufacture` - Manufacture of the gpu.
  * `memory` - Memory of the gpu.
  * `model` - Model of the gpu.
* `primary_network_interface` - A nested block describing the primary network interface of this instance.
Nested `primary_network_interface` blocks have the following structure:
  * `id` - The id of the network interface.
  * `name` - The name of the network interface.
  * `subnet` -  ID of the subnet.
  * `security_groups` -  List of security groups.
  * `primary_ipv4_address` - The primary IPv4 address.
  * `allow_ip_spoofing` - Indicates whether IP spoofing is allowed on this interface.
* `network_interfaces` - A nested block describing the additional network interface of this instance.
Nested `network_interfaces` blocks have the following structure:
  * `id` - The id of the network interface.
  * `name` - The name of the network interface.
  * `subnet` -  ID of the subnet.
  * `security_groups` -  List of security groups.
  * `primary_ipv4_address` - The primary IPv4 address.
  * `allow_ip_spoofing` - Indicates whether IP spoofing is allowed on this interface.
* `boot_volume` - A nested block describing the boot volume.
Nested `boot_volume` blocks have the following structure:
  * `name` - The name of the boot volume.
  * `size` -  Capacity of the volume in GB.
  * `iops` -  Input/Output Operations Per Second for the volume.
  * `profile` - The profile of the volume.
  * `encryption` - The encryption of the boot volume.
* `volume_attachments` - A nested block describing the volume attachments.  
Nested `volume_attachments` block have the following structure:
  * `id` - The id of the volume attachment
  * `name` -  The name of the volume attachment
  * `volume_id` - The id of the volume attachment's volume
  * `volume_name` -  The name of the volume attachment's volume
  * `volume_crn` -  The CRN of the volume attachment's volume
* `disks` - Collection of the instance's disks. Nested `disks` blocks have the following structure:
	* `created_at` - The date and time that the disk was created.
	* `href` - The URL for this instance disk.
	* `id` - The unique identifier for this instance disk.
	* `interface_type` - The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.
	* `name` - The user-defined name for this disk.
	* `resource_type` - The resource type.
	* `size` - The size of the disk in GB (gigabytes).
## Import

ibm_is_instance can be imported using instanceID, eg

```
$ terraform import ibm_is_instance.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
