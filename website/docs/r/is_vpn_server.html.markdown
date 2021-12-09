---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : VPN-Server"
description: |-
  Manages IBM VPN Server.
---

# ibm_is_vpn_server

Provides a resource for VPNServer. This allows VPNServer to be created, updated and deleted.

## Example Usage
The following example creates a VPN Server:

```terraform
resource "ibm_is_vpn_server" "is_vpn_server" {
  certificate_crn = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:00406b5615f95dba9bf7c2ab52bb3083"
  client_authentication {
    method    = "certificate"
    client_ca = "crn:v1:staging:public:cloudcerts:us-south:a/efe5afc483594adaa8325e2b4d1290df:86f62739-f3a8-42ac-abea-f23255965983:certificate:6a85a87d01dd5a4268a8bca16cb998eb"
  }
  client_ip_pool         = "10.5.0.0/21"
  subnets                = [ibm_is_subnet.subnet1.id]
  client_dns_server_ips  = ["192.168.3.4"]
  client_idle_timeout    = 2800
  enable_split_tunneling = false
  name                   = "my_vpn_server"
  port                   = 443
  protocol               = "udp"
}
```

## Timeouts
The `ibm_is_vpn_server` resource provides the following [Timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create**: The creation of the VPN server is considered `failed` when no response is received for 10 minutes. 
- **update**: The update of the VPN server is considered `failed` when no response is received for 10 minutes. 
- **delete**: The deletion of the VPN server is considered `failed` when no response is received for 10 minutes. 

## Argument Reference
Review the argument references that you can specify for your resource. 

* `certificate` - (Required, String) The certificate CRN instance for this VPN server.
* `client_authentication` - (Required, List) The methods used to authenticate VPN clients to this VPN server.
* `client_dns_server_ips` - (Optional, List) The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered, the DNS server addresses that will be provided to VPN clients connected to this VPN server.
* `client_idle_timeout` - (Optional, Integer) The seconds a VPN client can be idle before this VPN server will disconnect it.   Specify `0` to prevent the server from disconnecting idle clients.
  * Constraints: The maximum value is `28800`. The minimum value is `0`.
* `client_ip_pool` - (Required, String) The VPN client IPv4 address pool, expressed in CIDR format. The request must not overlap with any existing address prefixes in the VPC or any of the following reserved address ranges:  - `127.0.0.0/8` (IPv4 loopback addresses)  - `161.26.0.0/16` (IBM services)  - `166.8.0.0/14` (Cloud Service Endpoints)  - `169.254.0.0/16` (IPv4 link-local addresses)  - `224.0.0.0/4` (IPv4 multicast addresses)The prefix length of the client IP address pool's CIDR must be between`/9` (8,388,608 addresses) and `/22` (1024 addresses). A CIDR block that contains twice the number of IP addresses that are required to enable the maximum number of concurrent connections is recommended.
  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$/`
* `enable_split_tunneling` - (Optional, Boolean) Indicates whether the split tunneling is enabled on this VPN server.
  * Constraints: The default value is `false`.
* `name` - (Optional, String) The user-defined name for this VPN server. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the VPC this VPN server is serving.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$/`
* `port` - (Optional, Integer) The port number to use for this VPN server.
  * Constraints: The maximum value is `65535`. The minimum value is `1`.
* `protocol` - (Optional, String) The transport protocol to use for this VPN server.
  * Constraints: The default value is `udp`. Allowable values are: udp, tcp
- `resource_group` - (Optional, Forces new resource, String) The resource group where the VPN gateway to be created.
* `security_groups` - (Optional, List) The security groups to use for this VPN server. If unspecified, the VPC's default security group is used.
* `subnets` - (Required, List) The subnets to provision this VPN server in.  Use subnets in different zones for high availability.



  * Constraints: The minimum length is `1` item.
Nested scheme for **client_authentication**:
	* `method` - (Required, String) The type of authentication.
	  * Constraints: Allowable values are: certificate, username
	* `identity_provider` - (Optional, String) The type of identity provider to be used by VPN client.The type of identity provider to be used by the VPN client.- `iam`: IBM identity and access managementThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered.
		  * Constraints: Allowable values are: iam
	* `client_ca` - (Optional, String) The CRN for this certificate instance, the certificate instance to use for the VPN client certificate authority (CA).
	* `crl` - (Optional, String) The certificate revocation list contents, encoded in PEM format.


  * Constraints: The maximum length is `2` items. The minimum length is `1` item.
Nested scheme for **subnets**:
	* `id` - (Optional, String) The unique identifier for this subnet.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`
	* `crn` - (Optional, String) The CRN for this subnet.
	* `href` - (Optional, String) The URL for this subnet.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the VPNServer.
* `vpn_server` - The unique identifier of the VPNServer.
* `client_auto_delete` - (Required, Boolean) If set to `true`, disconnected VPN clients will be automatically deleted after the `client_auto_delete_timeout` time has passed.
* `client_auto_delete_timeout` - (Required, Integer) Hours after which disconnected VPN clients will be automatically deleted. If `0`, disconnected VPN clients will be deleted immediately.
  * Constraints: The maximum value is `24`. The minimum value is `0`.
* `created_at` - (Required, String) The date and time that the VPN server was created.
* `crn` - (Required, String) The CRN for this VPN server.
* `health_state` - (Required, String) The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.
  * Constraints: Allowable values are: ok, degraded, faulted, inapplicable
* `hostname` - (Required, String) Fully qualified domain name assigned to this VPN server.
* `href` - (Required, String) The URL for this VPN server.
  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
* `lifecycle_state` - (Required, String) The lifecycle state of the VPN server.
  * Constraints: Allowable values are: deleting, failed, pending, stable, updating, waiting, suspended
* `private_ips` - (Required, List) The reserved IPs bound to this VPN server.
* `resource_type` - (Required, String) The type of resource referenced.
  * Constraints: Allowable values are: vpn_server
* `version` - Version of the VPNServer.

Nested scheme for **private_ips**:
	* `address` - (Required, String) The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.
	  * Constraints: The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`
	* `deleted` - (Optional, List) If present, this property indicates the referenced resource has been deleted and providessome supplementary information.
	Nested scheme for **deleted**:
		* `more_info` - (Required, String) Link to documentation about deleted resources.
		  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `href` - (Required, String) The URL for this reserved IP.
	  * Constraints: The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`
	* `id` - (Required, String) The unique identifier for this reserved IP.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`
	* `name` - (Required, String) The user-defined or system-provided name for this reserved IP.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`
	* `resource_type` - (Required, String) The resource type.
	  * Constraints: Allowable values are: subnet_reserved_ip

## Import

You can import the `ibm_is_vpn_server` resource by using `id`. The unique identifier for this VPN server.

# Syntax
```
$ terraform import ibm_is_vpn_server.is_vpn_server <id>
```

# Example
```
$ terraform import ibm_is_vpn_server.is_vpn_server r006-d7cc5196-9864-48c4-82d8-3f30da41fcc5
```
