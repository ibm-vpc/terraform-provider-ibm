// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsNetworkACL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsNetworkACLRead,

		Schema: map[string]*schema.Schema{
			isNetworkACLName: {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vpc_name"},
				ExactlyOneOf: []string{"name", "network_acl"},
				Description:  "The network acl name.",
			},
			"vpc_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{isNetworkACLName},
				Description:  "The name of the vpc the network acl resides in.",
			},
			"network_acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "network_acl"},
				Description:  "The network acl id.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the network ACL was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this network ACL.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network ACL.",
			},
			isNetworkACLResourceGroup: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this network ACL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			isNetworkACLRules: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The ordered rules for this network ACL. If no rules exist, all traffic will be denied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRuleAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to allow or deny matching traffic.",
						},
						"before": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The rule that this rule is immediately before. In a rule collection, this alwaysrefers to the next item in the collection. If absent, this is the last rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this network ACL rule.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network ACL rule.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this network ACL rule.",
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the rule was created.",
						},
						isNetworkACLRuleDestination: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination CIDR block. The CIDR block `0.0.0.0/0` applies to all addresses.",
						},
						isNetworkACLRuleDirection: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the traffic to be matched is `inbound` or `outbound`.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network ACL rule.",
						},
						isNetworkACLRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this network ACL rule.",
						},
						isNetworkACLRuleIPVersion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version for this rule.",
						},
						isNetworkACLRuleName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this rule. Names must be unique within the network ACL the rule resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						isNetworkACLRuleProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network protocol",
						},
						isNetworkACLRuleSource: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source CIDR block. The CIDR block `0.0.0.0/0` applies to all addresses.",
						},
						isNetworkACLRuleICMP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The protocol ICMP",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRuleICMPCode: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic code to allow. Valid values from 0 to 255.",
									},
									isNetworkACLRuleICMPType: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic type to allow. Valid values from 0 to 254.",
									},
								},
							},
						},

						isNetworkACLRuleTCP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TCP protocol",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The highest port in the range of ports to be matched",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The lowest port in the range of ports to be matched",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The highest port in the range of ports to be matched",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The lowest port in the range of ports to be matched",
									},
								},
							},
						},

						isNetworkACLRuleUDP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "UDP protocol",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The highest port in the range of ports to be matched",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The lowest port in the range of ports to be matched",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The highest port in the range of ports to be matched",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The lowest port in the range of ports to be matched",
									},
								},
							},
						},
					},
				},
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The subnets to which this network ACL is attached.",
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
					},
				},
			},
			"vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this network ACL is a part of.",
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
					},
				},
			},
			isNetworkACLAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func dataSourceIBMIsNetworkACLRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpc_name_str := ""
	network_acl_name := ""
	if vpc, ok := d.GetOk("vpc_name"); ok {
		vpc_name_str = vpc.(string)
	}
	networkACL := &vpcv1.NetworkACL{}
	if vpc_name_str != "" {
		network_acl_name = d.Get(isNetworkACLName).(string)

		start := ""
		allrecs := []vpcv1.NetworkACL{}
		for {
			listNetworkAclsOptions := &vpcv1.ListNetworkAclsOptions{}
			if start != "" {
				listNetworkAclsOptions.Start = &start
			}
			networkACLCollection, response, err := vpcClient.ListNetworkAclsWithContext(context, listNetworkAclsOptions)
			if err != nil || networkACLCollection == nil {
				log.Printf("[DEBUG] ListNetworkAclsWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("ListNetworkAclsWithContext failed %s\n%s", err, response))
			}
			start = flex.GetNext(networkACLCollection.Next)
			allrecs = append(allrecs, networkACLCollection.NetworkAcls...)
			if start == "" {
				break
			}
		}
		acl_found := false
		for _, networkAcl := range allrecs {
			if *networkAcl.VPC.Name == vpc_name_str && network_acl_name == *networkAcl.Name {
				networkACL = &networkAcl
				acl_found = true
				break
			}
		}

		if !acl_found {
			log.Printf("[DEBUG] No networkACL found with given VPC %s and ACL name %s", vpc_name_str, network_acl_name)
			return diag.FromErr(fmt.Errorf("[ERROR] No networkACL found with given VPC %s and ACL name %s", vpc_name_str, network_acl_name))
		}
	} else {

		getNetworkACLOptions := &vpcv1.GetNetworkACLOptions{}

		getNetworkACLOptions.SetID(d.Get("network_acl").(string))

		networkACLInst, response, err := vpcClient.GetNetworkACLWithContext(context, getNetworkACLOptions)
		if err != nil || networkACLInst == nil {
			log.Printf("[DEBUG] GetNetworkACLWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetNetworkACLWithContext failed %s\n%s", err, response))
		}
		networkACL = networkACLInst
	}
	d.SetId(*networkACL.ID)
	if err = d.Set("created_at", flex.DateTimeToString(networkACL.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", networkACL.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("href", networkACL.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("name", networkACL.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if networkACL.ResourceGroup != nil {
		err = d.Set(isNetworkACLResourceGroup, dataSourceNetworkACLFlattenResourceGroup(*networkACL.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group %s", err))
		}
	}

	if networkACL.Rules != nil {
		err = d.Set(isNetworkACLRules, dataSourceNetworkACLFlattenRules(networkACL.Rules))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting rules %s", err))
		}
	}

	if networkACL.Subnets != nil {
		err = d.Set("subnets", dataSourceNetworkACLFlattenSubnets(networkACL.Subnets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnets %s", err))
		}
	}

	if networkACL.VPC != nil {
		err = d.Set("vpc", dataSourceNetworkACLFlattenVPC(*networkACL.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting vpc %s", err))
		}
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *networkACL.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource Network ACL (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isNetworkACLAccessTags, accesstags)

	return nil
}

func dataSourceNetworkACLFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceNetworkACLResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceNetworkACLResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceNetworkACLFlattenRules(result []vpcv1.NetworkACLRuleItemIntf) (rules []map[string]interface{}) {
	for _, rulesItem := range result {
		rules = append(rules, dataSourceNetworkACLRulesToMap(rulesItem))
	}

	return rules
}

func dataSourceNetworkACLRulesToMap(rule vpcv1.NetworkACLRuleItemIntf) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	switch rule := rule.(type) {
	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
		setCommonSecurityRuleFields(rulesMap, rule.ID, rule.Href, rule.Protocol, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction, rule.Before)
		rulesMap["created_at"] = flex.DateTimeToString(rule.CreatedAt)
		rulesMap[isNetworkACLRuleTCP] = []map[string]int{}
		rulesMap[isNetworkACLRuleUDP] = []map[string]int{}
		icmp := []map[string]int{}
		if rule.Code != nil && rule.Type != nil {
			icmp = append(icmp, map[string]int{
				isNetworkACLRuleICMPCode: int(*rule.Code),
				isNetworkACLRuleICMPType: int(*rule.Type),
			})
		}
		rulesMap[isNetworkACLRuleICMP] = icmp

	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
		setCommonSecurityRuleFields(rulesMap, rule.ID, rule.Href, rule.Protocol, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction, rule.Before)
		rulesMap["created_at"] = flex.DateTimeToString(rule.CreatedAt)
		rulesMap[isNetworkACLRuleICMP] = []map[string]int{}
		if *rule.Protocol == "tcp" {
			rulesMap[isNetworkACLRuleUDP] = []map[string]int{}
			tcp := []map[string]int{{
				isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rule.SourcePortMax),
				isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rule.SourcePortMin),
				isNetworkACLRulePortMax:       checkNetworkACLNil(rule.DestinationPortMax),
				isNetworkACLRulePortMin:       checkNetworkACLNil(rule.DestinationPortMin),
			}}
			rulesMap[isNetworkACLRuleTCP] = tcp
		} else if *rule.Protocol == "udp" {
			rulesMap[isNetworkACLRuleTCP] = []map[string]int{}
			udp := []map[string]int{{
				isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rule.SourcePortMax),
				isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rule.SourcePortMin),
				isNetworkACLRulePortMax:       checkNetworkACLNil(rule.DestinationPortMax),
				isNetworkACLRulePortMin:       checkNetworkACLNil(rule.DestinationPortMin),
			}}
			rulesMap[isNetworkACLRuleUDP] = udp
		}

	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny:
		setCommonSecurityRuleFields(rulesMap, rule.ID, rule.Href, rule.Protocol, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction, rule.Before)
		rulesMap["created_at"] = flex.DateTimeToString(rule.CreatedAt)
		rulesMap[isNetworkACLRuleICMP] = []map[string]int{}
		rulesMap[isNetworkACLRuleTCP] = []map[string]int{}
		rulesMap[isNetworkACLRuleUDP] = []map[string]int{}

	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp:
		setCommonSecurityRuleFields(rulesMap, rule.ID, rule.Href, rule.Protocol, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction, rule.Before)
		rulesMap["created_at"] = flex.DateTimeToString(rule.CreatedAt)
		rulesMap[isNetworkACLRuleICMP] = []map[string]int{}
		rulesMap[isNetworkACLRuleTCP] = []map[string]int{}
		rulesMap[isNetworkACLRuleUDP] = []map[string]int{}

	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual:
		setCommonSecurityRuleFields(rulesMap, rule.ID, rule.Href, rule.Protocol, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction, rule.Before)
		rulesMap["created_at"] = flex.DateTimeToString(rule.CreatedAt)
		rulesMap[isNetworkACLRuleICMP] = []map[string]int{}
		rulesMap[isNetworkACLRuleTCP] = []map[string]int{}
		rulesMap[isNetworkACLRuleUDP] = []map[string]int{}
	}

	return rulesMap
}

func setCommonSecurityRuleFields(r map[string]interface{}, id, href, protocol, name, action, ipVersion, source, destination, direction *string, before *vpcv1.NetworkACLRuleReference) {
	r["id"] = *id
	r[isNetworkACLRuleHref] = *href
	r[isNetworkACLRuleProtocol] = *protocol
	if before != nil {
		beforeList := []map[string]interface{}{}
		beforeMap := dataSourceNetworkACLRulesBeforeToMap(*before)
		beforeList = append(beforeList, beforeMap)
		r[isNwACLRuleBefore] = beforeList
	}
	r[isNetworkACLRuleName] = *name
	r[isNetworkACLRuleAction] = *action
	r[isNetworkACLRuleIPVersion] = *ipVersion
	r[isNetworkACLRuleSource] = *source
	r[isNetworkACLRuleDestination] = *destination
	r[isNetworkACLRuleDirection] = *direction
}

func dataSourceNetworkACLRulesBeforeToMap(beforeItem vpcv1.NetworkACLRuleReference) (beforeMap map[string]interface{}) {
	beforeMap = map[string]interface{}{}

	if beforeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLBeforeDeletedToMap(*beforeItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		beforeMap["deleted"] = deletedList
	}
	if beforeItem.Href != nil {
		beforeMap["href"] = beforeItem.Href
	}
	if beforeItem.ID != nil {
		beforeMap["id"] = beforeItem.ID
	}
	if beforeItem.Name != nil {
		beforeMap["name"] = beforeItem.Name
	}

	return beforeMap
}

func dataSourceNetworkACLBeforeDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkACLFlattenSubnets(result []vpcv1.SubnetReference) (subnets []map[string]interface{}) {
	for _, subnetsItem := range result {
		subnets = append(subnets, dataSourceNetworkACLSubnetsToMap(subnetsItem))
	}

	return subnets
}

func dataSourceNetworkACLSubnetsToMap(subnetsItem vpcv1.SubnetReference) (subnetsMap map[string]interface{}) {
	subnetsMap = map[string]interface{}{}

	if subnetsItem.CRN != nil {
		subnetsMap["crn"] = subnetsItem.CRN
	}
	if subnetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLSubnetsDeletedToMap(*subnetsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetsMap["deleted"] = deletedList
	}
	if subnetsItem.Href != nil {
		subnetsMap["href"] = subnetsItem.Href
	}
	if subnetsItem.ID != nil {
		subnetsMap["id"] = subnetsItem.ID
	}
	if subnetsItem.Name != nil {
		subnetsMap["name"] = subnetsItem.Name
	}

	return subnetsMap
}

func dataSourceNetworkACLSubnetsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkACLFlattenVPC(result vpcv1.VPCReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceNetworkACLVPCToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceNetworkACLVPCToMap(vpcItem vpcv1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLVPCDeletedToMap(*vpcItem.Deleted)
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

	return vpcMap
}

func dataSourceNetworkACLVPCDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
