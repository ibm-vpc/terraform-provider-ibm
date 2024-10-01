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

func DataSourceIBMIsClusterNetworkInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkInterfaceRead,

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network identifier.",
			},
			"is_cluster_network_interface_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network interface identifier.",
			},
			"allow_ip_spoofing": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether source IP spoofing is allowed on this cluster network interface. If `false`, source IP spoofing is prevented on this cluster network interface. If `true`, source IP spoofing is allowed on this cluster network interface.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network interface was created.",
			},
			"enable_infrastructure_nat": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network interface.",
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
				Description: "The lifecycle state of the cluster network interface.",
			},
			"mac_address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The MAC address of the cluster network interface. May be absent if`lifecycle_state` is `pending`.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
			},
			"primary_ip": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cluster network subnet reserved IP for this cluster network interface.",
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
			"protocol_state_filtering_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol state filtering mode used for this cluster network interface.Protocol state filtering monitors each network connection flowing over this cluster network interface, and drops any packets that are invalid based on the current connection state and protocol. See [About Protocol State Filtering](https://cloud.ibm.com/docs/__TBD__) for more information.",
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
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target of this cluster network interface.If absent, this cluster network interface is not attached to a target.The resources supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this cluster network interface resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
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
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone this cluster network interface resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this zone.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsClusterNetworkInterfaceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkInterfaceOptions := &vpcv1.GetClusterNetworkInterfaceOptions{}

	getClusterNetworkInterfaceOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	getClusterNetworkInterfaceOptions.SetID(d.Get("is_cluster_network_interface_id").(string))

	clusterNetworkInterface, _, err := vpcClient.GetClusterNetworkInterfaceWithContext(context, getClusterNetworkInterfaceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkInterfaceWithContext failed: %s", err.Error()), "(Data) ibm_is_cluster_network_interface", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getClusterNetworkInterfaceOptions.ClusterNetworkID, *getClusterNetworkInterfaceOptions.ID))

	if err = d.Set("allow_ip_spoofing", clusterNetworkInterface.AllowIPSpoofing); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_ip_spoofing: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("auto_delete", clusterNetworkInterface.AutoDelete); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting auto_delete: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkInterface.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("enable_infrastructure_nat", clusterNetworkInterface.EnableInfrastructureNat); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enable_infrastructure_nat: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", clusterNetworkInterface.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	lifecycleReasons := []map[string]interface{}{}
	if clusterNetworkInterface.LifecycleReasons != nil {
		for _, modelItem := range clusterNetworkInterface.LifecycleReasons {
			modelMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
				return tfErr.GetDiag()
			}
			lifecycleReasons = append(lifecycleReasons, modelMap)
		}
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_reasons: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("lifecycle_state", clusterNetworkInterface.LifecycleState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("mac_address", clusterNetworkInterface.MacAddress); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting mac_address: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", clusterNetworkInterface.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	primaryIP := []map[string]interface{}{}
	if clusterNetworkInterface.PrimaryIP != nil {
		modelMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(clusterNetworkInterface.PrimaryIP)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
			return tfErr.GetDiag()
		}
		primaryIP = append(primaryIP, modelMap)
	}
	if err = d.Set("primary_ip", primaryIP); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting primary_ip: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("protocol_state_filtering_mode", clusterNetworkInterface.ProtocolStateFilteringMode); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting protocol_state_filtering_mode: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", clusterNetworkInterface.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	subnet := []map[string]interface{}{}
	if clusterNetworkInterface.Subnet != nil {
		modelMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(clusterNetworkInterface.Subnet)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
			return tfErr.GetDiag()
		}
		subnet = append(subnet, modelMap)
	}
	if err = d.Set("subnet", subnet); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting subnet: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	target := []map[string]interface{}{}
	if clusterNetworkInterface.Target != nil {
		modelMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(clusterNetworkInterface.Target)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
			return tfErr.GetDiag()
		}
		target = append(target, modelMap)
	}
	if err = d.Set("target", target); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting target: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	vpc := []map[string]interface{}{}
	if clusterNetworkInterface.VPC != nil {
		modelMap, err := DataSourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(clusterNetworkInterface.VPC)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
			return tfErr.GetDiag()
		}
		vpc = append(vpc, modelMap)
	}
	if err = d.Set("vpc", vpc); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting vpc: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	zone := []map[string]interface{}{}
	if clusterNetworkInterface.Zone != nil {
		modelMap, err := DataSourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(clusterNetworkInterface.Zone)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_interface", "read")
			return tfErr.GetDiag()
		}
		zone = append(zone, modelMap)
	}
	if err = d.Set("zone", zone); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting zone: %s", err), "(Data) ibm_is_cluster_network_interface", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(model *vpcv1.ClusterNetworkInterfaceLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceDeletedToMap(model.Deleted)
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

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceDeletedToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceDeletedToMap(model.Deleted)
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

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceDeletedToMap(model *vpcv1.ClusterNetworkSubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(model vpcv1.ClusterNetworkInterfaceTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReference); ok {
		return DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceToMap(model.(*vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReference))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfaceTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkInterfaceTarget)
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkInterfaceTargetIntf subtype encountered")
	}
}

func DataSourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceToMap(model *vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworkInterfaceVPCReferenceDeletedToMap(model.Deleted)
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

func DataSourceIBMIsClusterNetworkInterfaceVPCReferenceDeletedToMap(model *vpcv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
