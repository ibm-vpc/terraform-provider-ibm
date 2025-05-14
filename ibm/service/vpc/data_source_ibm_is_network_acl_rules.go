// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isNwACLRules = "rules"
)

func DataSourceIBMISNetworkACLRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISNetworkACLRulesRead,

		Schema: map[string]*schema.Schema{
			"direction": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The direction of the rules to filter",
			},
			isNwACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network ACL id",
			},
			isNwACLRules: {
				Type:        schema.TypeList,
				Description: "List of network acl rules",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNwACLRuleId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network acl rule id.",
						},
						isNetworkACLRuleName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this rule",
						},
						isNwACLRuleBefore: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule that this rule is immediately before. If absent, this is the last rule.",
						},
						isNetworkACLRuleHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network ACL rule.",
						},
						isNetworkACLRuleProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network protocol",
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
				},
			},
		},
	}
}

func dataSourceIBMISNetworkACLRulesRead(d *schema.ResourceData, meta interface{}) error {
	nwACLID := d.Get(isNwACLID).(string)
	err := networkACLRulesList(d, meta, nwACLID)
	if err != nil {
		return err
	}
	return nil
}

func networkACLRulesList(d *schema.ResourceData, meta interface{}, nwACLID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.NetworkACLRuleItemIntf{}
	listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
		NetworkACLID: &nwACLID,
	}
	if directionIntf, ok := d.GetOk("direction"); ok {
		direction := directionIntf.(string)
		listNetworkACLRulesOptions.Direction = &direction
	}
	for {

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
	rulesInfo := make([]map[string]interface{}, 0)
	for _, rule := range allrecs {
		l := map[string]interface{}{}
		switch v := rule.(type) {
		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
			setCommonRuleFields(l, v.ID, v.Href, v.Protocol, v.Name, v.Action, v.IPVersion, v.Source, v.Destination, v.Direction, v.Before)
			l[isNetworkACLRuleTCP] = []map[string]int{}
			l[isNetworkACLRuleUDP] = []map[string]int{}
			icmp := []map[string]int{}
			if v.Code != nil && v.Type != nil {
				icmp = append(icmp, map[string]int{
					isNetworkACLRuleICMPCode: int(*v.Code),
					isNetworkACLRuleICMPType: int(*v.Type),
				})
			}
			l[isNetworkACLRuleICMP] = icmp

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
			setCommonRuleFields(l, v.ID, v.Href, v.Protocol, v.Name, v.Action, v.IPVersion, v.Source, v.Destination, v.Direction, v.Before)
			l[isNetworkACLRuleICMP] = []map[string]int{}
			if *v.Protocol == "tcp" {
				l[isNetworkACLRuleUDP] = []map[string]int{}
				tcp := []map[string]int{{
					isNetworkACLRuleSourcePortMax: checkNetworkACLNil(v.SourcePortMax),
					isNetworkACLRuleSourcePortMin: checkNetworkACLNil(v.SourcePortMin),
					isNetworkACLRulePortMax:       checkNetworkACLNil(v.DestinationPortMax),
					isNetworkACLRulePortMin:       checkNetworkACLNil(v.DestinationPortMin),
				}}
				l[isNetworkACLRuleTCP] = tcp
			} else if *v.Protocol == "udp" {
				l[isNetworkACLRuleTCP] = []map[string]int{}
				udp := []map[string]int{{
					isNetworkACLRuleSourcePortMax: checkNetworkACLNil(v.SourcePortMax),
					isNetworkACLRuleSourcePortMin: checkNetworkACLNil(v.SourcePortMin),
					isNetworkACLRulePortMax:       checkNetworkACLNil(v.DestinationPortMax),
					isNetworkACLRulePortMin:       checkNetworkACLNil(v.DestinationPortMin),
				}}
				l[isNetworkACLRuleUDP] = udp
			}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny:
			setCommonRuleFields(l, v.ID, v.Href, v.Protocol, v.Name, v.Action, v.IPVersion, v.Source, v.Destination, v.Direction, v.Before)
			l[isNetworkACLRuleICMP] = []map[string]int{}
			l[isNetworkACLRuleTCP] = []map[string]int{}
			l[isNetworkACLRuleUDP] = []map[string]int{}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp:
			setCommonRuleFields(l, v.ID, v.Href, v.Protocol, v.Name, v.Action, v.IPVersion, v.Source, v.Destination, v.Direction, v.Before)
			l[isNetworkACLRuleICMP] = []map[string]int{}
			l[isNetworkACLRuleTCP] = []map[string]int{}
			l[isNetworkACLRuleUDP] = []map[string]int{}

		case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual:
			setCommonRuleFields(l, v.ID, v.Href, v.Protocol, v.Name, v.Action, v.IPVersion, v.Source, v.Destination, v.Direction, v.Before)
			l[isNetworkACLRuleICMP] = []map[string]int{}
			l[isNetworkACLRuleTCP] = []map[string]int{}
			l[isNetworkACLRuleUDP] = []map[string]int{}
		}
		rulesInfo = append(rulesInfo, l)
	}
	d.SetId(dataSourceIBMISNetworkACLRulesId(d))
	d.Set(isNetworkACLRules, rulesInfo)
	return nil
}

func setCommonRuleFields(l map[string]interface{}, id, href, protocol, name, action, ipVersion, source, destination, direction *string, before *vpcv1.NetworkACLRuleReference) {
	l[isNwACLRuleId] = *id
	l[isNetworkACLRuleHref] = *href
	l[isNetworkACLRuleProtocol] = *protocol
	if before != nil {
		l[isNwACLRuleBefore] = *before.ID
	}
	l[isNetworkACLRuleName] = *name
	l[isNetworkACLRuleAction] = *action
	l[isNetworkACLRuleIPVersion] = *ipVersion
	l[isNetworkACLRuleSource] = *source
	l[isNetworkACLRuleDestination] = *destination
	l[isNetworkACLRuleDirection] = *direction
}

// dataSourceIBMISNetworkACLRulesId returns a reasonable ID for a rule list.
func dataSourceIBMISNetworkACLRulesId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
