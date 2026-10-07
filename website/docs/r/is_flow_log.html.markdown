---

subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : is_flow_log"
description: |-
  Manages IBM VPC flow log.
---

# ibm_is_flow_log
Create, update, delete and suspend the flow log resource. For more information, about VPC flow log, see [creating a flow log collector](https://cloud.ibm.com/docs/vpc?topic=vpc-ordering-flow-log-collector).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```


## Example usage

```terraform

resource "ibm_is_instance" "example" {
  name    = "example-instance"
  image   = ibm_is_image.example.id
  profile = "bx2-2x8"

  primary_network_interface {
    subnet     = ibm_is_subnet.example.id
  }

  vpc  = ibm_is_vpc.example.id
  zone = "us-south-1"
  keys = [ibm_is_ssh_key.example.id]
}


resource "ibm_resource_group" "example" {
  name = "example-resource-group"
}

resource "ibm_resource_instance" "example" {
  name              = "example-cos-instance"
  resource_group_id = ibm_resource_group.example.id
  service           = "cloud-object-storage"
  plan              = "standard"
  location          = "global"
}

resource "ibm_cos_bucket" "example" {
  bucket_name          = "us-south-bucket-vpc1"
  resource_instance_id = ibm_resource_instance.example.id
  region_location      = var.region
  storage_class        = "standard"
}

resource "ibm_is_flow_log" "example" {
  depends_on     = [ibm_cos_bucket.example]
  name           = "example-instance-flow-log"
  target         = ibm_is_instance.example.id
  active         = true
  storage_bucket = ibm_cos_bucket.example.bucket_name
}

```

### Using the `destination` block (Cloud Object Storage)

```terraform
resource "ibm_is_flow_log" "example_cos_destination" {
  name   = "example-cos-destination-flow-log"
  target = ibm_is_instance.example.id
  active = true
  destination = {
    type           = "cloud_object_storage"
    storage_bucket = ibm_cos_bucket.example.bucket_name
  }
}
```

### Using the `destination` block (IBM Cloud Logs)

```terraform
resource "ibm_is_flow_log" "example_cloud_logs_destination" {
  name   = "example-cloud-logs-flow-log"
  target = ibm_is_subnet.example.id
  active = true
  destination = {
    type = "ibm_cloud_logs"
  }
}
```


## Argument reference
Review the argument references that you can specify for your resource. 

- `access_tags`  - (Optional, List of Strings) A list of access management tags to attach to the flow log.

  ~> **Note:** 
  **&#x2022;** You can attach only those access tags that already exists.</br>
  **&#x2022;** For more information, about creating access tags, see [working with tags](https://cloud.ibm.com/docs/account?topic=account-tag&interface=ui#create-access-console).</br>
  **&#x2022;** You must have the access listed in the [Granting users access to tag resources](https://cloud.ibm.com/docs/account?topic=account-access) for `access_tags`</br>
  **&#x2022;** `access_tags` must be in the format `key:value`.
- `name` - (Required, String) The unique user-defined name for the flow log collector.
- `target` - (Required, Forces new resource, String) The ID of the target to collect flow logs.

  -> **Note:**
  **&#x2022;** If the target is an instance network attachment, flow logs will be collected  for that instance network attachment.</br>
  **&#x2022;** If the target is an instance network interface, flow logs will be collected  for that instance network interface.</br>
  **&#x2022;** If the target is a virtual network interface, flow logs will be collected for the  the virtual network interface's `target` resource if the resource is:  - an instance network attachment.</br>
  **&#x2022;** If the target is a virtual server instance, flow logs will be collected  for all network attachments or network interfaces on that instance.</br>
  **&#x2022;** If the target is a subnet, flow logs will be collected  for all instance network interfaces and virtual network interfaces  attached to that subnet.</br>
  **&#x2022;** If the target is a VPC, flow logs will be collected for all instance network  interfaces and virtual network interfaces  attached to all subnets within that VPC. If the target is an instance, subnet, or VPC, flow logs will not be collectedfor any instance network attachments or instance network interfaces within the targetthat are themselves the target of a more specific flow log collector.</br>
- `storage_bucket` - (Optional, Forces new resource, String) The name of the IBM Cloud Object Storage bucket where the collected flows will be logged. The bucket must exist and an IAM service authorization must grant IBM Cloud flow logs resources of VPC infrastructure services writer access to the bucket. Conflicts with `destination`; use one or the other, not both.
- `destination` - (Optional, Forces new resource, Map) The destination for the collected flow logs. Conflicts with `storage_bucket`; use one or the other, not both. The map supports the following keys:
  - `type` - (Required, String) The destination type. Supported values are `cloud_object_storage` and `ibm_cloud_logs`.
  - `storage_bucket` - (Required when `type` is `cloud_object_storage`, String) The Cloud Object Storage bucket name where the collected flows will be logged. Must not be set when `type` is `ibm_cloud_logs`.
- `active` - (Optional, String) Indicates whether the collector is active. If **false**, this collector is created in inactive mode. Default value is true.
- `resource_group` - (Optional, Forces new resource, String) The resource group ID where the flow log is created.
- `tags` - (Optional, Array of Strings) The tags associated with the flow log.


## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `auto_delete` - (Boolean) Indicates whether this flow log collector will be automatically deleted when target is deleted. At present, this is always true, but may be modifiable in the future.
- `created_at`-  (String) The date and time that the flow log collector created.
- `crn` - (String) The CRN of the flow log collector.
- `destination` - (Map) The destination for the collected flow logs. This is populated from the API even when the resource was created using the legacy `storage_bucket` argument. Map keys:
  - `type` - (String) The destination type. One of `cloud_object_storage` or `ibm_cloud_logs`.
  - `storage_bucket` - (String) The Cloud Object Storage bucket name. Present only when `type` is `cloud_object_storage`.
- `href` - (String) The URL of the flow log collector.
- `id` - (String) The unique identifier of the flow log collector.
- `lifecycle_state` - (String) The lifecycle state of the flow log collector.
- `name`-  (String) The user-defined name of the flow log collector.
- `storage_bucket` - (String) The Cloud Object Storage bucket name where the collected flows are logged.
- `vpc` - (String) The VPC of the flow log collector that is associated.


## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the `ibm_is_flow_log` resource by using `id`.
The `id` property can be formed from `VPC flow log ID`. For example:

```terraform
import {
  to = ibm_is_flow_log.example
  id = "<vpc_flow_log_id>"
}
```

Using `terraform import`. For example:

```console
% terraform import ibm_is_flow_log.example <vpc_flow_log_id>
```