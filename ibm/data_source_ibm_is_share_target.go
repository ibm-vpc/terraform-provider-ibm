// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsShareTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareTargetRead,

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The share target identifier.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the share target was created.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share target.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the mount target.",
			},
			"mount_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this share target.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's CRN.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this security group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
						},
					},
				},
			},
			"subnet": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The subnet associated with this file share target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this subnet.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this subnet.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this subnet.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this subnet.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC to which this share target is allowing to mount the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this VPC.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsShareTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareTargetOptions := &vpcv1.GetShareTargetOptions{}

	share_id := d.Get("share").(string)
	getShareTargetOptions.SetShareID(share_id)
	getShareTargetOptions.SetID(d.Get("target").(string))

	shareTarget, response, err := vpcClient.GetShareTargetWithContext(context, getShareTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] GetShareTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", share_id, *shareTarget.ID))
	if err = d.Set("created_at", shareTarget.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", shareTarget.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", shareTarget.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("mount_path", shareTarget.MountPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_path: %s", err))
	}
	if err = d.Set("name", shareTarget.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("resource_type", shareTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if shareTarget.SecurityGroups != nil {
		err = d.Set("security_groups", dataSourceShareTargetFlattenSecurityGroups(shareTarget.SecurityGroups))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting security_groups %s", err))
		}
	}

	if shareTarget.Subnet != nil {
		err = d.Set("subnet", dataSourceShareTargetFlattenSubnet(*shareTarget.Subnet))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subnet %s", err))
		}
	}

	if shareTarget.VPC != nil {
		err = d.Set("vpc", dataSourceShareTargetFlattenVpc(*shareTarget.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	return nil
}

func dataSourceShareTargetFlattenSecurityGroups(result []vpcv1.SecurityGroupReference) (securityGroups []map[string]interface{}) {
	for _, securityGroupsItem := range result {
		securityGroups = append(securityGroups, dataSourceShareTargetSecurityGroupsToMap(securityGroupsItem))
	}

	return securityGroups
}

func dataSourceShareTargetSecurityGroupsToMap(securityGroupsItem vpcv1.SecurityGroupReference) (securityGroupsMap map[string]interface{}) {
	securityGroupsMap = map[string]interface{}{}

	if securityGroupsItem.CRN != nil {
		securityGroupsMap["crn"] = securityGroupsItem.CRN
	}
	if securityGroupsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetSecurityGroupsDeletedToMap(*securityGroupsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		securityGroupsMap["deleted"] = deletedList
	}
	if securityGroupsItem.Href != nil {
		securityGroupsMap["href"] = securityGroupsItem.Href
	}
	if securityGroupsItem.ID != nil {
		securityGroupsMap["id"] = securityGroupsItem.ID
	}
	if securityGroupsItem.Name != nil {
		securityGroupsMap["name"] = securityGroupsItem.Name
	}

	return securityGroupsMap
}

func dataSourceShareTargetSecurityGroupsDeletedToMap(deletedItem vpcv1.SecurityGroupReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareTargetFlattenSubnet(result vpcv1.SubnetReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareTargetSubnetToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareTargetSubnetToMap(subnetItem vpcv1.SubnetReference) (subnetMap map[string]interface{}) {
	subnetMap = map[string]interface{}{}

	if subnetItem.CRN != nil {
		subnetMap["crn"] = subnetItem.CRN
	}
	if subnetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetSubnetDeletedToMap(*subnetItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetMap["deleted"] = deletedList
	}
	if subnetItem.Href != nil {
		subnetMap["href"] = subnetItem.Href
	}
	if subnetItem.ID != nil {
		subnetMap["id"] = subnetItem.ID
	}
	if subnetItem.Name != nil {
		subnetMap["name"] = subnetItem.Name
	}
	if subnetItem.ResourceType != nil {
		subnetMap["resource_type"] = subnetItem.ResourceType
	}

	return subnetMap
}

func dataSourceShareTargetSubnetDeletedToMap(deletedItem vpcv1.SubnetReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareTargetFlattenVpc(result vpcv1.VPCReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareTargetVpcToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareTargetVpcToMap(vpcItem vpcv1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetVpcDeletedToMap(*vpcItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcMap["deleted"] = deletedList
	}
	if vpcItem.Href != nil {
		vpcMap["href"] = vpcItem.Href
	}
	if vpcItem.ID != nil {
		vpcMap["id"] = vpcItem.ID
	}
	if vpcItem.Name != nil {
		vpcMap["name"] = vpcItem.Name
	}
	if vpcItem.ResourceType != nil {
		vpcMap["resource_type"] = vpcItem.ResourceType
	}

	return vpcMap
}

func dataSourceShareTargetVpcDeletedToMap(deletedItem vpcv1.VPCReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
