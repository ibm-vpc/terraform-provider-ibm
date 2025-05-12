// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isNetworkACLRuleHref = "href"
)

func DataSourceIBMISNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISNetworkACLRuleRead,

		Schema: map[string]*schema.Schema{
			isNwACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network ACL id",
			},
			isNwACLRuleId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL rule id",
			},
			isNwACLRuleBefore: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule that this rule is immediately before. If absent, this is the last rule.",
			},
			isNetworkACLRuleName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleName),
				Description:  "The user-defined name for this rule",
			},
			isNetworkACLRuleProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the network protocol",
			},
			isNetworkACLRuleHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network ACL rule",
			},
			isNetworkACLRuleAction: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether to allow or deny matching traffic.",
			},
			isNetworkACLRuleIPVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for this rule.",
			},
			isNetworkACLRuleSource: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source IP address or CIDR block.",
			},
			isNetworkACLRuleDestination: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination IP address or CIDR block.",
			},
			isNetworkACLRuleDirection: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the traffic to be matched is inbound or outbound.",
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
	}
}

func dataSourceIBMISNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	nwACLID := d.Get(isNwACLID).(string)
	name := d.Get(isNetworkACLRuleName).(string)
	err := nawaclRuleDataGet(d, meta, name, nwACLID)
	if err != nil {
		return err
	}

	return nil
}

func nawaclRuleDataGet(d *schema.ResourceData, meta interface{}, name, nwACLID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.NetworkACLRuleItemIntf{}
	for {
		listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
			NetworkACLID: &nwACLID,
		}
		if start != "" {
			listNetworkACLRulesOptions.Start = &start
		}

		ruleList, response, err := sess.ListNetworkACLRules(listNetworkACLRulesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching network acl ruless %s\n%s", err, response)
		}
		start = flex.GetNext(ruleList.Next)

		allrecs = append(allrecs, ruleList.Rules...)
		if start == "" {
			break
		}
	}

	for _, rule := range allrecs {
		switch rulex := rule.(type) {
		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
			if *rulex.Name == name {
				setCommonACLRuleFields(d, nwACLID, rulex.ID, rulex.Href, rulex.Protocol,
					rulex.Name, rulex.Action, rulex.IPVersion, rulex.Source, rulex.Destination, rulex.Direction, rulex.Before)
				d.Set(isNetworkACLRuleTCP, []map[string]int{})
				d.Set(isNetworkACLRuleUDP, []map[string]int{})
				icmp := []map[string]int{}
				if rulex.Code != nil && rulex.Type != nil {
					icmp = append(icmp, map[string]int{
						isNetworkACLRuleICMPCode: int(*rulex.Code),
						isNetworkACLRuleICMPType: int(*rulex.Type),
					})
				}
				d.Set(isNetworkACLRuleICMP, icmp)
				break
			}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
			if *rulex.Name == name {
				setCommonACLRuleFields(d, nwACLID, rulex.ID, rulex.Href, rulex.Protocol,
					rulex.Name, rulex.Action, rulex.IPVersion, rulex.Source, rulex.Destination, rulex.Direction, rulex.Before)
				d.Set(isNetworkACLRuleICMP, []map[string]int{})
				if *rulex.Protocol == "tcp" {
					d.Set(isNetworkACLRuleUDP, []map[string]int{})
					tcp := []map[string]int{{
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						isNetworkACLRulePortMax:       checkNetworkACLNil(rulex.DestinationPortMax),
						isNetworkACLRulePortMin:       checkNetworkACLNil(rulex.DestinationPortMin),
					}}
					d.Set(isNetworkACLRuleTCP, tcp)
				} else if *rulex.Protocol == "udp" {
					d.Set(isNetworkACLRuleTCP, []map[string]int{})
					udp := []map[string]int{{
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						isNetworkACLRulePortMax:       checkNetworkACLNil(rulex.DestinationPortMax),
						isNetworkACLRulePortMin:       checkNetworkACLNil(rulex.DestinationPortMin),
					}}
					d.Set(isNetworkACLRuleUDP, udp)
					break
				}
			}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny:
			if *rulex.Name == name {
				setCommonACLRuleFields(d, nwACLID, rulex.ID, rulex.Href, rulex.Protocol,
					rulex.Name, rulex.Action, rulex.IPVersion, rulex.Source, rulex.Destination, rulex.Direction, rulex.Before)
				d.Set(isNetworkACLRuleICMP, []map[string]int{})
				d.Set(isNetworkACLRuleTCP, []map[string]int{})
				d.Set(isNetworkACLRuleUDP, []map[string]int{})
				break
			}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp:
			if *rulex.Name == name {
				setCommonACLRuleFields(d, nwACLID, rulex.ID, rulex.Href, rulex.Protocol,
					rulex.Name, rulex.Action, rulex.IPVersion, rulex.Source, rulex.Destination, rulex.Direction, rulex.Before)
				d.Set(isNetworkACLRuleICMP, []map[string]int{})
				d.Set(isNetworkACLRuleTCP, []map[string]int{})
				d.Set(isNetworkACLRuleUDP, []map[string]int{})
				break
			}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual:
			if *rulex.Name == name {
				setCommonACLRuleFields(d, nwACLID, rulex.ID, rulex.Href, rulex.Protocol,
					rulex.Name, rulex.Action, rulex.IPVersion, rulex.Source, rulex.Destination, rulex.Direction, rulex.Before)
				d.Set(isNetworkACLRuleICMP, []map[string]int{})
				d.Set(isNetworkACLRuleTCP, []map[string]int{})
				d.Set(isNetworkACLRuleUDP, []map[string]int{})
				break
			}
		}

	}
	return nil
}

func setCommonACLRuleFields(d *schema.ResourceData, nwACLID string, id, href, protocol, name, action, ipVersion, source, destination, direction *string, before *vpcv1.NetworkACLRuleReference) {
	d.SetId(makeTerraformACLRuleID(nwACLID, *id))
	d.Set(isNwACLRuleId, *id)
	if before != nil {
		d.Set(isNwACLRuleBefore, *before)
	}
	d.Set(isNetworkACLRuleHref, *href)
	d.Set(isNetworkACLRuleProtocol, *protocol)
	d.Set(isNetworkACLRuleName, *name)
	d.Set(isNetworkACLRuleAction, *action)
	d.Set(isNetworkACLRuleIPVersion, *ipVersion)
	d.Set(isNetworkACLRuleSource, *source)
	d.Set(isNetworkACLRuleDestination, *destination)
	d.Set(isNetworkACLRuleDirection, *direction)
}
