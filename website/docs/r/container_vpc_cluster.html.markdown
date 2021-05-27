---

subcategory: "Kubernetes Service"
layout: "ibm"
page_title: "IBM: container_vpc_cluster"
description: |-
  Manages IBM VPC container cluster.
---

# ibm\_container_vpc_cluster

Create or delete a Kubernetes VPC cluster.

In the following example, you can create a Gen-2 VPC cluster with a default worker pool with one worker:
```terraform
resource "ibm_container_vpc_cluster" "cluster" {
  name              = "my_vpc_cluster"
  vpc_id            = "r006-abb7c7ea-aadf-41bd-94c5-b8521736fadf"
  kube_version      = "1.17.5"
  flavor            = "bx2.2x8"
  worker_count      = "1"
  resource_group_id = data.ibm_resource_group.resource_group.id
  zones {
      subnet_id = "0717-0c0899ce-48ac-4eb6-892d-4e2e1ff8c9478"
      name      = "us-south-1"
    }
}
```

Create the Openshift Cluster with default worker Pool entitlement with one worker node:
```terraform
resource "ibm_resource_instance" "cos_instance" {
  name     = "my_cos_instance"
  service  = "cloud-object-storage"
  plan     = "standard"
  location = "global"
}

resource "ibm_container_vpc_cluster" "cluster" {
  name              = "my_vpc_cluster"
  vpc_id            = "r006-abb7c7ea-aadf-41bd-94c5-b8521736fadf"
  kube_version      = "4.3_openshift"
  flavor            = "bx2.16x64"
  worker_count      = "2"
  entitlement       = "cloud_pak"
  cos_instance_crn  = ibm_resource_instance.cos_instance.id
  resource_group_id = data.ibm_resource_group.resource_group.id
  zones {
      subnet_id = "0717-0c0899ce-48ac-4eb6-892d-4e2e1ff8c9478"
      name      = "us-south-1"
    }
}
```

Create a Kms Enabled Kubernetes cluster:

```terraform
resource "ibm_container_vpc_cluster" "cluster" {
  name              = "cluster2"
  vpc_id            = ibm_is_vpc.vpc1.id
  flavor            = "bx2.2x8"
  worker_count      = "1"
  wait_till         = "OneWorkerNodeReady"
  resource_group_id = data.ibm_resource_group.resource_group.id
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }

  kms_config {
      instance_id = "12043812-757f-4e1e-8436-6af3245e6a69"
      crk_id = "0792853c-b9f9-4b35-9d9e-ffceab51d3c1"
      private_endpoint = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `flavor` - (Required, Forces new resource, string) The flavor of the VPC worker node.
* `name` - (Required, Forces new resource, string) The name of the cluster.
* `vpc_id` - (Required, Forces new resource, string) The ID of the VPC in which to create the worker nodes. To list available IDs, run 'ibmcloud ks vpcs'.
* `zones` - (Required, set) A nested block describing the zones of this VPC cluster. Nested zones blocks have the following structure:
  * `subnet-id` - (Required, string) The VPC subnet to assign the cluster. 
  * `name` - (Required, string) Name of the zone.
* `disable_public_service_endpoint` - (Optional,Bool) Disable the public service endpoint to prevent public access to the master. Default Value 'false'.
* `kube_version` - (Optional,String) Specify the Kubernetes version, including at least the major.minor version. If you do not include this flag, the default version is used. To see available versions, run 'ibmcloud ks versions'.
* `update_all_workers` - (Optional, bool)  Set to `true` if you want to update workers kube version.
* `wait_for_worker_update` - (Optional, bool) Set to `true` to wait for kube version of woker nodes to update during the wokrer node kube version update.
  **NOTE**: setting `wait_for_worker_update` to `false` is not recommended. This results in upgradign all the worker nodes in the cluster at the same time causing the cluster downtime
* `pod_subnet` - (Optional, Forces new resource,String) Specify a custom subnet CIDR to provide private IP addresses for pods. The subnet must be at least '/23' or larger. For more info, refer [here](https://cloud.ibm.com/docs/containers?topic=containers-cli-plugin-kubernetes-service-cli#pod-subnet).
* `service_subnet` - (Optional, Forces new resource,String) Specify a custom subnet CIDR to provide private IP addresses for services. The subnet must be at least '/24' or larger. For more info, refer [here](https://cloud.ibm.com/docs/containers?topic=containers-cli-plugin-kubernetes-service-cli#service-subnet).
* `worker_count` - (Optional, Int) The number of worker nodes per zone in the default worker pool. Default value '1'.
* `worker_labels` - (Optional, map) Labels on all the workers in the default worker pool.
* `resource_group_id` - (Optional, Forces new resource, string) The ID of the resource group. You can retrieve the value from data source `ibm_resource_group`. If not provided defaults to default resource group.
* `tags` - (Optional, array of strings) Tags associated with the container cluster instance.
* `kms_config` -  (Optional, list) Used to attach a key protect instance to a cluster. Nested `kms_config` block has the following structure:
	* `instance_id` - The guid of the key protect instance.
	* `crk_id` - Id of the customer root key (CRK).
	* `private_endpoint` - Set this to true to configure the KMS private service endpoint. Default is false.
* `entitlement` - (Optional, String) The openshift cluster entitlement avoids the OCP licence charges incurred. Use cloud paks with OCP Licence entitlement to create the Openshift cluster.
  **NOTE**:
  1. It is set only for the first time creation of the cluster, modification in the further runs will not have any impacts.
  2. Set this argument to 'cloud_pak' only if you use this cluster with a Cloud Pak that has an OpenShift entitlement
* `cos_instance_crn` - (Optional, String) Required for OpenShift clusters only. The standard cloud object storage instance CRN to back up the internal registry in your OpenShift on VPC Gen 2 cluster.
* `wait_till` - (Optional, String) The cluster creation happens in multi-stages. To avoid the longer wait times for resource execution, this field is introduced.
Resource will wait for only the specified stage and complete execution. The supported stages are
  - *MasterNodeReady*: resource will wait till the master node is ready
  - *OneWorkerNodeReady*: resource will wait till atleast one worker node becomes to ready state
  - *IngressReady*: resource will wait till the ingress-host and ingress-secret are available.

  Default value: IngressReady
* `force_delete_storage` - (Optional, bool) If set to true, force the removal of persistent storage associated with the cluster during cluster deletion. Default: false
    **NOTE**: Before doing terraform destroy if force_delete_storage param is introduced after provisioning the cluster, a terraform apply must be done before terraform destroy for force_delete_storage param to take effect.
* `patch_version` - (Optional, string) Set this to update the worker nodes with the required patch version. 
   The patch_version should be in the format - `patch_version_fixpack_version`. Learn more about the Kuberentes version [here](https://cloud.ibm.com/docs/containers?topic=containers-cs_versions).
    **NOTE**: To update the patch/fixpack versions of the worker nodes, Run the command `ibmcloud ks workers -c <cluster_name_or_id> --output json`, fetch the required patch & fixpack versions from `kubeVersion.target` and set the patch_version parameter.
* `retry_patch_version` - (Optional, int) This argument helps to retry the update of patch_version if the previous update fails. Increment the value to retry the update of patch_version on worker nodes.

**NOTE**:
1. For users on account to add tags to a resource, they must be assigned the appropriate access. Learn more about tags permission [here](https://cloud.ibm.com/docs/resources?topic=resources-access)
2. `wait_till` is set only for the first time creation of the resource, modification in the further runs will not any impacts.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Id of the cluster
* `crn` - CRN of the cluster.
* `ingress_hostname` - The Ingress hostname.
* `ingress_secret` - The Ingress secret.
* `master_status` - Status of kubernetes master.
* `master_url` - The Master server URL.
* `private_service_endpoint_url` - Private service endpoint url.
* `public_service_endpoint_url` - Public service endpoint url.
* `state` - State.
* `albs` - Application load balancer (ALB)'s attached to the cluster
  * `id` - The application load balancer (ALB) id.
  * `name` - The name of the application load balancer (ALB).
  * `alb_type` - The application load balancer (ALB) type public or private.
  * `enable` -  Enable (true) or disable(false) application load balancer (ALB).
  * `state` - The status of the application load balancer (ALB)(enabled or disabled).
  * `resize` - Indicate whether resizing should be done.
  * `disable_deployment` - Indicate whether to disable deployment only on disable application load balancer (ALB).
  * `load_balancer_hostname` - The host name of the application load balancer (ALB).


## Import

`ibm_container_vpc_cluster` can be imported using clusterID, eg ibm_container_vpc_cluster.cluster

```
$ terraform import ibm_container_vpc_cluster.cluster bmonvocd0i8m2v5dmb6g
```
