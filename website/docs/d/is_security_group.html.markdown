---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : security_group"
description: |-
  Reads IBM Cloud security group.
---

# ibm_is_security_group
Retrieve information about a security group as a read-only data source. For more information, about managing IBM Cloud security group , see [about security group](https://cloud.ibm.com/docs/vpc?topic=vpc-using-security-groups).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example usage
The following example allows to create a different types of protocol rules `ALL`, `ICMP`, `UDP`, `TCP` and read the security group.

```terraform
resource "ibm_is_vpc" "example" {
  name = "example-vpc"
}

resource "ibm_is_security_group" "example" {
  name = "example-sg"
  vpc  = ibm_is_vpc.example.id
}

resource "ibm_is_security_group_rule" "example" {
  group     = ibm_is_security_group.example.id
  direction = "inbound"
  remote    = "127.0.0.1"
}

resource "ibm_is_security_group_rule" "example" {
  group     = ibm_is_security_group.example.id
  direction = "inbound"
  remote    = "127.0.0.1"
  # Deprecated block: replaced with 'protocol', 'code', and 'type' arguments
  # icmp {
  #   code = 20
  #   type = 30
  # }
  protocol  = "icmp"
  code      = 20
  type      = 30
}

resource "ibm_is_security_group_rule" "example" {
  group     = ibm_is_security_group.example.id
  direction = "inbound"
  remote    = "127.0.0.1"
  # Deprecated block: replaced with 'protocol', 'port_min', and 'port_max' arguments
  # udp {
  #   port_min = 805
  #   port_max = 807
  # }
  protocol  = "udp"
  port_min = 805
  port_max = 807
}

resource "ibm_is_security_group_rule" "example" {
  group     = ibm_is_security_group.example.id
  direction = "inbound"
  remote    = "127.0.0.1"
  # Deprecated block: replaced with 'protocol', 'port_min', and 'port_max' arguments
  # tcp {
  #  port_min = 8080
  #  port_max = 8080
  # }
  protocol  = "tcp"
  port_min = 8080
  port_max = 8080
}

data "ibm_is_security_group" "example" {
  name = ibm_is_security_group.example.name
}

data "ibm_is_security_group" "examplevpc" {
  name = ibm_is_security_group.example.name
  vpc  = ibm_is_vpc.example.id
}
```

## Argument reference
Review the argument references that you can specify for your resource. 

- `name` - (Required, String) The name of the security group.
- `vpc` - (Optional, String) The identifier of the vpc where this security group resides. (Useful when two security groups have same name across different VPCs)
- `vpc_name` - (Optional, String) The name of the vpc where this security group resides. (Useful when two security groups have same name across different VPCs)
- `resource_group` - (Optional, String) The identifier of the resource group where this security group resides.

## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created. 

- `access_tags`  - (List) Access management tags associated for the security group.
- `crn` - The CRN of the security group.
- `id` - (String) The ID of the security group.
- `rules` - (List of Objects) The rules associated with security group. Each rule has following attributes.

  Nested scheme for `rules`:
  - `rule_id` - (String) ID of the rule.
  - `direction` - (String) Direction of traffic to enforce, either `inbound` or `outbound`.
  - `local` - (String) The local IP address or range of local IP addresses to which this rule will allow inbound traffic (or from which, for outbound traffic). A CIDR block of `0.0.0.0/0` allows traffic to all local IPv4 addresses; `::/0` allows traffic to all local IPv6 addresses. Accepts an IP address or a CIDR block (IPv4 or IPv6).
  - `ip_version` - (String) The IP version to enforce. Supported values are `ipv4` and `ipv6`.
  - `protocol` - (String) The protocol to enforce. Supported values for IPv4 rules: `icmp`, `tcp`, `udp`, `any`, `icmp_tcp_udp`, and others. Supported additional values for IPv6 rules: `ipv6_icmp` (protocol 58), `ipv6_hop_opt` (protocol 0), `ipv6_route` (protocol 43), `ipv6_frag` (protocol 44), `ipv6_no_next` (protocol 59), `ipv6_dest_opts` (protocol 60), `ipv6_mobility` (protocol 135).
  - `type` - (String) The ICMP traffic type to allow.
  - `name` - (String) The name for this security group rule. The name must not be used by another rule in the security group.
  - `code` - (String) The ICMP traffic code to allow. Also used for `ipv6_icmp`.
  - `port_max` - (Integer) The TCP/UDP port range that includes the maximum bound.
  - `port_min` - (Integer) The TCP/UDP port range that includes the minimum bound.
  - `remote` - (String) Security group ID, an IP address, a CIDR block (IPv4 or IPv6), or a single security group identifier.
- `tags` - Tags associated with the security group.
  


