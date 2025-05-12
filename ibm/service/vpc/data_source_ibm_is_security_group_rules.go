// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIsSecurityGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"security_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group identifier.",
			},
			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direction": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The direction of traffic to enforce, either `inbound` or `outbound`.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this security group rule.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this security group rule.",
						},
						"ip_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version to enforce. The format of `remote.address` or `remote.cidr_block` must match this property, if they are used. Alternatively, if `remote` references a security group, then this rule only applies to IP addresses (network interfaces) in that group matching this IP version.",
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network protocol.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP addresses or security groups from which this rule allows traffic (or to which,for outbound rules). Can be specified as an IP address, a CIDR block, or a securitygroup. A CIDR block of `0.0.0.0/0` allows traffic from any source (or to any source,for outbound rules).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"cidr_block": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CIDR block. This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's CRN.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The security group's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this security group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
									},
								},
							},
						},
						"local": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The local IP address or range of local IP addresses to which this rule will allow inbound traffic (or from which, for outbound traffic). A CIDR block of 0.0.0.0/0 allows traffic to all local IP addresses (or from all local IP addresses, for outbound rules).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"cidr_block": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CIDR block. This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
									},
								},
							},
						},
						"code": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ICMP traffic code to allow.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ICMP traffic type to allow.",
						},
						"port_max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The inclusive upper bound of TCP/UDP port range.",
						},
						"port_min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The inclusive lower bound of TCP/UDP port range.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	secGrpId := d.Get("security_group").(string)
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listSecurityGroupRuleOptions := vpcv1.ListSecurityGroupRulesOptions{
		SecurityGroupID: &secGrpId,
	}

	ruleList, response, err := sess.ListSecurityGroupRules(&listSecurityGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error fetching security group rules %s\n%s", err, response)
	}

	rulesInfo := make([]map[string]interface{}, 0)
	for _, securityGroupRuleIntf := range ruleList.Rules {
		r := map[string]interface{}{}
		switch rule := securityGroupRuleIntf.(type) {
		case *vpcv1.SecurityGroupRuleProtocolAny:
			{
				setCommonSecurityGroupRulesFields(r, rule.ID, rule.Direction, rule.Href, rule.IPVersion, rule.Protocol, rule.Remote, rule.Local)
			}
		case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
			{
				setCommonSecurityGroupRulesFields(r, rule.ID, rule.Direction, rule.Href, rule.IPVersion, rule.Protocol, rule.Remote, rule.Local)
				if rule.Code != nil {
					r["code"] = *rule.Code
				}
				if rule.Type != nil {
					r["type"] = *rule.Type
				}
			}
		case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
			{
				setCommonSecurityGroupRulesFields(r, rule.ID, rule.Direction, rule.Href, rule.IPVersion, rule.Protocol, rule.Remote, rule.Local)
				r["port_max"] = *rule.PortMax
				r["port_min"] = *rule.PortMin
			}
		case *vpcv1.SecurityGroupRuleProtocolIcmptcpudp:
			{
				setCommonSecurityGroupRulesFields(r, rule.ID, rule.Direction, rule.Href, rule.IPVersion, rule.Protocol, rule.Remote, rule.Local)
			}
		case *vpcv1.SecurityGroupRuleProtocolIndividual:
			{
				setCommonSecurityGroupRulesFields(r, rule.ID, rule.Direction, rule.Href, rule.IPVersion, rule.Protocol, rule.Remote, rule.Local)
			}
		}
		rulesInfo = append(rulesInfo, r)
	}
	d.SetId(dataSourceIBMIsSecurityGroupRulesID(d))
	d.Set("rules", rulesInfo)
	return nil
}

func setCommonSecurityGroupRulesFields(r map[string]interface{}, ruleID, direction, href, ipVersion, protocol *string, remoteIntf vpcv1.SecurityGroupRuleRemoteIntf, localIntf vpcv1.SecurityGroupRuleLocalIntf) {
	r["direction"] = direction
	r["href"] = href
	r["id"] = ruleID
	r["ip_version"] = ipVersion
	r["protocol"] = protocol
	// nested map for remote.
	if remoteIntf != nil {
		remoteList := []map[string]interface{}{}
		remoteMap := dataSourceSecurityGroupRuleRemoteToMap(remoteIntf.(*vpcv1.SecurityGroupRuleRemote))
		remoteList = append(remoteList, remoteMap)
		r["remote"] = remoteList
	}
	// nested map for local.
	if localIntf != nil {
		localList := []map[string]interface{}{}
		localMap := dataSourceSecurityGroupRuleLocalToMap(localIntf.(*vpcv1.SecurityGroupRuleLocal))
		localList = append(localList, localMap)
		r["local"] = localList
	}
}

// dataSourceIBMIsSecurityGroupRulesID returns a reasonable ID for the list.
func dataSourceIBMIsSecurityGroupRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceSecurityGroupRuleCollectionRemoteDeletedToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = *deletedItem.MoreInfo
	}

	return resultMap
}
