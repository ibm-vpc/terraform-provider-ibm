---
subcategory: "Internet services"
layout: "ibm"
page_title: "IBM: ibm_cis_ip_addresses"
description: |-
  List the IP addresses used by name servers by Cloud Internet Services. Required for setting whitelist addresses for internet facing application ports.
---

# ibm_cis_ip_addresses

Import the lists of all IP addresses used by the CIS proxy. The CIS proxy uses only addresses from this list, for both client-to-proxy and proxy-to-origin communication. You can then reference the IP addresses by interpolation to configure firewalls, network ACLs and Security Groups to white list these addresses.

## Example Usage

```terraform
data "ibm_cis_ip_addresses" "ip_addresses" {}
```

## Argument Reference

No arguments are required. All CIS instances on an account use the same range of name servers.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `ipv4_cidrs` - The ipv4 address ranges used by CIS for name servers. To be whitelisted by the service user.
- `ipv6_cidrs` - The ipv6 address ranges used by CIS for name servers. To be whitelisted by the service user.
