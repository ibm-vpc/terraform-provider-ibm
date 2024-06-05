// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsInstanceClusterNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceClusterNetworkAttachmentRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual server instance identifier.",
			},
			"is_instance_cluster_network_attachment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance cluster network attachment identifier.",
			},
			"before": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The instance cluster network attachment that is immediately before. If absent, this is thelast instance cluster network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance cluster network attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance cluster network attachment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"cluster_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cluster network interface for this instance cluster network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this cluster network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this cluster network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
						},
						"primary_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP for this cluster network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this cluster network subnet reserved IP.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this cluster network subnet reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this cluster network subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this cluster network subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance cluster network attachment.",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the instance cluster network attachment.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIBMIsInstanceClusterNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getInstanceClusterNetworkAttachmentOptions := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{}

	getInstanceClusterNetworkAttachmentOptions.SetInstanceID(d.Get("instance_id").(string))
	getInstanceClusterNetworkAttachmentOptions.SetID(d.Get("is_instance_cluster_network_attachment_id").(string))

	instanceClusterNetworkAttachment, _, err := vpcClient.GetInstanceClusterNetworkAttachmentWithContext(context, getInstanceClusterNetworkAttachmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceClusterNetworkAttachmentWithContext failed: %s", err.Error()), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*instanceClusterNetworkAttachment.ID)

	before := []map[string]interface{}{}
	if instanceClusterNetworkAttachment.Before != nil {
		modelMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(instanceClusterNetworkAttachment.Before)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachment", "read")
			return tfErr.GetDiag()
		}
		before = append(before, modelMap)
	}
	if err = d.Set("before", before); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting before: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	clusterNetworkInterface := []map[string]interface{}{}
	if instanceClusterNetworkAttachment.ClusterNetworkInterface != nil {
		modelMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(instanceClusterNetworkAttachment.ClusterNetworkInterface)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachment", "read")
			return tfErr.GetDiag()
		}
		clusterNetworkInterface = append(clusterNetworkInterface, modelMap)
	}
	if err = d.Set("cluster_network_interface", clusterNetworkInterface); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cluster_network_interface: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", instanceClusterNetworkAttachment.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	lifecycleReasons := []map[string]interface{}{}
	if instanceClusterNetworkAttachment.LifecycleReasons != nil {
		for _, modelItem := range instanceClusterNetworkAttachment.LifecycleReasons {
			modelMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachment", "read")
				return tfErr.GetDiag()
			}
			lifecycleReasons = append(lifecycleReasons, modelMap)
		}
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_reasons: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("lifecycle_state", instanceClusterNetworkAttachment.LifecycleState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", instanceClusterNetworkAttachment.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", instanceClusterNetworkAttachment.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_instance_cluster_network_attachment", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(model *vpcv1.InstanceClusterNetworkAttachmentBefore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(model *vpcv1.ClusterNetworkInterfaceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	primaryIPMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = *model.ResourceType
	subnetMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceDeletedToMap(model *vpcv1.ClusterNetworkInterfaceReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceDeletedToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceDeletedToMap(model *vpcv1.ClusterNetworkSubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(model *vpcv1.InstanceClusterNetworkAttachmentLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
