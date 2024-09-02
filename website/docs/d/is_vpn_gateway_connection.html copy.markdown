---
layout: "ibm"
page_title: "IBM : ibm_is_vpn_gateway_connection"
description: |-
  Get information about is_vpn_gateway_connection
subcategory: "Virtual Private Cloud API"
---

# ibm_is_vpn_gateway_connection

Provides a read-only data source to retrieve information about an is_vpn_gateway_connection. You can then reference the fields of the data source in other resources within the same configuration by using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_vpn_gateway_connection" "is_vpn_gateway_connection" {
	is_vpn_gateway_connection_id = ibm_is_vpn_gateway_connection.is_vpn_gateway_connection_instance.is_vpn_gateway_connection_id
	vpn_gateway_id = ibm_is_vpn_gateway_connection.is_vpn_gateway_connection_instance.vpn_gateway_id
}
```

## Argument Reference

You can specify the following arguments for this data source.

* `is_vpn_gateway_connection_id` - (Required, Forces new resource, String) The VPN gateway connection identifier.
* `vpn_gateway_id` - (Required, Forces new resource, String) The VPN gateway identifier.

## Attribute Reference

After your data source is created, you can read values from the following attributes.

* `id` - The unique identifier of the is_vpn_gateway_connection.
* `admin_state_up` - (Boolean) If set to false, the VPN gateway connection is shut down.
* `authentication_mode` - (String) The authentication mode. Only `psk` is currently supported.
  * Constraints: Allowable values are: `psk`.
* `created_at` - (String) The date and time that this VPN gateway connection was created.
* `dead_peer_detection` - (List) The Dead Peer Detection settings.
Nested schema for **dead_peer_detection**:
	* `action` - (String) Dead Peer Detection actions.
	  * Constraints: Allowable values are: `clear`, `hold`, `none`, `restart`.
	* `interval` - (Integer) Dead Peer Detection interval in seconds.
	  * Constraints: The maximum value is `86399`. The minimum value is `1`.
	* `timeout` - (Integer) Dead Peer Detection timeout in seconds. Must be at least the interval.
	  * Constraints: The maximum value is `86399`. The minimum value is `2`.
* `establish_mode` - (String) The establish mode of the VPN gateway connection:- `bidirectional`: Either side of the VPN gateway can initiate IKE protocol   negotiations or rekeying processes.- `peer_only`: Only the peer can initiate IKE protocol negotiations for this VPN gateway   connection. Additionally, the peer is responsible for initiating the rekeying process   after the connection is established. If rekeying does not occur, the VPN gateway   connection will be brought down after its lifetime expires.
  * Constraints: Allowable values are: `bidirectional`, `peer_only`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `href` - (String) The VPN connection's canonical URL.
  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `ike_policy` - (List) The IKE policy. If absent, [auto-negotiation isused](https://cloud.ibm.com/docs/vpc?topic=vpc-using-vpn&interface=ui#ike-auto-negotiation-phase-1).
Nested schema for **ike_policy**:
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The IKE policy's canonical URL.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this IKE policy.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this IKE policy. The name is unique across all IKE policies in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `ike_policy`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `ipsec_policy` - (List) The IPsec policy. If absent, [auto-negotiation isused](https://cloud.ibm.com/docs/vpc?topic=vpc-using-vpn&interface=ui#ipsec-auto-negotiation-phase-2).
Nested schema for **ipsec_policy**:
	* `deleted` - (List) If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.
	Nested schema for **deleted**:
		* `more_info` - (String) Link to documentation about deleted resources.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `href` - (String) The IPsec policy's canonical URL.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
	* `id` - (String) The unique identifier for this IPsec policy.
	  * Constraints: The maximum length is `64` characters. The minimum length is `1` character. The value must match regular expression `/^[-0-9a-z_]+$/`.
	* `name` - (String) The name for this IPsec policy. The name is unique across all IPsec policies in the region.
	  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$/`.
	* `resource_type` - (String) The resource type.
	  * Constraints: Allowable values are: `ipsec_policy`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `local` - (List) 
Nested schema for **local**:
	* `ike_identities` - (List) The local IKE identities.A VPN gateway in static route mode consists of two members in active-active mode. The first identity applies to the first member, and the second identity applies to the second member.
	  * Constraints: The maximum length is `2` items. The minimum length is `2` items.
	Nested schema for **ike_identities**:
		* `type` - (String) The IKE identity type.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
		  * Constraints: Allowable values are: `fqdn`, `hostname`, `ipv4_address`, `key_id`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `value` - (String) The IKE identity FQDN value.
		  * Constraints: The maximum length is `255` characters. The minimum length is `4` characters. The value must match regular expression `/^((?=[A-Za-z0-9-]{1,63}\\.)[A-Za-z0-9-]*\\.)+[A-Za-z]{2,63}\\.?$/`.
* `mode` - (String) The mode of the VPN gateway.
  * Constraints: Allowable values are: `policy`, `route`.
* `name` - (String) The name for this VPN gateway connection. The name is unique across all connections for the VPN gateway.
  * Constraints: The maximum length is `63` characters. The minimum length is `1` character. The value must match regular expression `/^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$/`.
* `peer` - (List) 
Nested schema for **peer**:
	* `address` - (String) The IP address of the peer VPN gateway for this connection.
	  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(?!(0(\\.0){3})|(255(\\.255){3})|(22[4-9]|23[0-9](\\.[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]){3}))([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])){3}$/`.
	* `fqdn` - (String) The FQDN of the peer VPN gateway for this connection.
	  * Constraints: The maximum length is `255` characters. The minimum length is `4` characters. The value must match regular expression `/^((?=[A-Za-z0-9-]{1,63}\\.)[A-Za-z0-9-]*\\.)+[A-Za-z]{2,63}\\.?$/`.
	* `ike_identity` - (List) The peer IKE identity.
	Nested schema for **ike_identity**:
		* `type` - (String) The IKE identity type.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
		  * Constraints: Allowable values are: `fqdn`, `hostname`, `ipv4_address`, `key_id`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
		* `value` - (String) The IKE identity FQDN value.
		  * Constraints: The maximum length is `255` characters. The minimum length is `4` characters. The value must match regular expression `/^((?=[A-Za-z0-9-]{1,63}\\.)[A-Za-z0-9-]*\\.)+[A-Za-z]{2,63}\\.?$/`.
	* `type` - (String) Indicates whether `peer.address` or `peer.fqdn` is used.
	  * Constraints: Allowable values are: `address`, `fqdn`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `psk` - (String) The pre-shared key.
  * Constraints: The value must match regular expression `/^(?=[\\-\\+\\&\\!\\@\\#\\$\\%\\^\\*\\(\\)\\,\\.\\:\\_a-zA-Z0-9]{6,128}$)(?:(?!^0[xs]).).*$/`.
* `resource_type` - (String) The resource type.
  * Constraints: Allowable values are: `vpn_gateway_connection`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `routing_protocol` - (String) Routing protocols are disabled for this VPN gateway connection.
  * Constraints: Allowable values are: `none`.
* `status` - (String) The status of a VPN gateway connection.
  * Constraints: Allowable values are: `down`, `up`. The maximum length is `128` characters. The minimum length is `1` character. The value must match regular expression `/^[a-z][a-z0-9]*(_[a-z0-9]+)*$/`.
* `status_reasons` - (List) The reasons for the current VPN gateway connection status (if any).
  * Constraints: The minimum length is `0` items.
Nested schema for **status_reasons**:
	* `code` - (String) A snake case string succinctly identifying the status reason.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
	  * Constraints: Allowable values are: `cannot_authenticate_connection`, `ike_policy_mismatch`, `ike_v1_id_local_remote_cidr_mismatch`, `ike_v2_local_remote_cidr_mismatch`, `internal_error`, `ipsec_policy_mismatch`, `peer_not_responding`. The value must match regular expression `/^[a-z]+(_[a-z]+)*$/`.
	* `message` - (String) An explanation of the reason for this VPN gateway connection's status.
	* `more_info` - (String) Link to documentation about this status reason.
	  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.
* `tunnels` - (List) The VPN tunnel configuration for this VPN gateway connection (in static route mode).
Nested schema for **tunnels**:
	* `public_ip` - (List) The IP address of the VPN gateway member in which the tunnel resides.
	Nested schema for **public_ip**:
		* `address` - (String) The IP address.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.
		  * Constraints: The maximum length is `15` characters. The minimum length is `7` characters. The value must match regular expression `/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/`.
	* `status` - (String) The status of the VPN Tunnel.
	  * Constraints: Allowable values are: `down`, `up`.
	* `status_reasons` - (List) The reasons for the current status (if any).
	  * Constraints: The minimum length is `0` items.
	Nested schema for **status_reasons**:
		* `code` - (String) A reason code for this status:- `cannot_authenticate_connection`: Failed to authenticate a connection because of  mismatched IKE ID and PSK (check IKE ID and PSK in peer VPN configuration)- `internal_error`: Internal error (contact IBM support)- `ike_policy_mismatch`: None of the proposed IKE crypto suites was acceptable (check   the IKE policies on both sides of the VPN)- `ike_v1_id_local_remote_cidr_mismatch`: Invalid IKE ID or mismatched local CIDRs and  remote CIDRs in IKE V1 (check the IKE ID or the local CIDRs and remote CIDRs in IKE  V1 configuration)- `ike_v2_local_remote_cidr_mismatch`: Mismatched local CIDRs and remote CIDRs in IKE  V2 (check the local CIDRs and remote CIDRs in IKE V2 configuration)- `ipsec_policy_mismatch`: None of the proposed IPsec crypto suites was acceptable  (check the IPsec policies on both sides of the VPN)- `peer_not_responding`: No response from peer (check network ACL configuration, peer  availability, and on-premise firewall configuration)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.
		  * Constraints: Allowable values are: `cannot_authenticate_connection`, `ike_policy_mismatch`, `ike_v1_id_local_remote_cidr_mismatch`, `ike_v2_local_remote_cidr_mismatch`, `internal_error`, `ipsec_policy_mismatch`, `peer_not_responding`. The value must match regular expression `/^[a-z]+(_[a-z]+)*$/`.
		* `message` - (String) An explanation of the reason for this VPN gateway connection tunnel's status.
		* `more_info` - (String) Link to documentation about this status reason.
		  * Constraints: The maximum length is `8000` characters. The minimum length is `10` characters. The value must match regular expression `/^http(s)?:\/\/([^\/?#]*)([^?#]*)(\\?([^#]*))?(#(.*))?$/`.

