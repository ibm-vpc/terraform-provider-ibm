---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : ibm_is_virtual_endpoint_gateway"
description: |-
  Manages IBM Virtual endpoint gateway
---

# ibm_is_virtual_endpoint_gateway

Provides a Virtual endpoint gateway resource. This allows Virtual endpoint gateway to be created, updated, and cancelled.

## Example Usage

In the following example, you can create a VPN gateway:

```terraform
resource "ibm_is_virtual_endpoint_gateway" "endpoint_gateway1" {

  name = "my-endpoint-gateway-1"
  target {
	name          = "ibm-ntp-server"
    resource_type = "provider_infrastructure_service"
  }
  vpc = ibm_is_vpc.testacc_vpc.id
  resource_group = data.ibm_resource_group.test_acc.id
}

resource "ibm_is_virtual_endpoint_gateway" "endpoint_gateway2" {
	name = "my-endpoint-gateway-1"
	target {
	  name          = "ibm-ntp-server"
	  resource_type = "provider_infrastructure_service"
	}
	vpc = ibm_is_vpc.testacc_vpc.id
	ips {
		subnet   = ibm_is_subnet.testacc_subnet.id
		name        = "test-reserved-ip1"
	}
	resource_group = data.ibm_resource_group.test_acc.id
}

resource "ibm_is_virtual_endpoint_gateway" "endpoint_gateway3" {
	name = "my-endpoint-gateway-1"
	target {
	  name          = "ibm-ntp-server"
	  resource_type = "provider_infrastructure_service"
	}
	vpc = ibm_is_vpc.testacc_vpc.id
	ips {
		id   = "0737-5ab3c18e-6f6c-4a69-8f48-20e3456647b5"
	}
	resource_group = data.ibm_resource_group.test_acc.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string,ForceNew) Endpoint gateway name
- `target` - (Required, stringList) Endpoint gateway target
  - `name` - (Optional, string, ForceNew) Endpoint gateway target name. If it is not specified, `crn` must be specified.
  - `crn` - (Optional, string, ForceNew) Endpoint gateway target crn. If it is not specified, `name` must be specified.
  - `resource_type`- (Required, string) Endpoint gateway target resource type
- `vpc` - (Required, string, ForceNew) The VPC id
- `ips` - (Optional, stringList)Endpoint gateway resource group
  - `id` - (Optional, string)Endpoint gateway resource group IPs id
  - `name` - (Optional, string)Endpoint gateway resource group IPs name
  - `subnet` - (Optional, string)Endpoint gateway resource group Subnet id
  - `resource_type` - (Computed, string)Endpoint gateway resource group VPC Resource Type
- `resource_group` - (Optional, string,ForceNew)The resource group id
- `tags` - (Optional, array of strings) Tags associated with the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the endpoint gateway connection.
- `resource_type` - Endpoint gateway resource type
- `created_at` - Endpoint gateway created date and time
- `health_state` - Endpoint gateway health state
- `lifecycle_state` - Endpoint gateway lifecycle state

## Import

ibm_is_virtual_endpoint_gateway can be imported using virtual endpoint gateway ID, eg

```
$ terraform import ibm_is_virtual_endpoint_gateway.example d7bec597-4726-451f-8a63-e62e6f19c32c

```
