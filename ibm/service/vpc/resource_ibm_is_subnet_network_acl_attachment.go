// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSubnetID     = "subnet"
	isNetworkACLID = "network_acl"
)

func ResourceIBMISSubnetNetworkACLAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISSubnetNetworkACLAttachmentCreate,
		Read:     resourceIBMISSubnetNetworkACLAttachmentRead,
		Update:   resourceIBMISSubnetNetworkACLAttachmentUpdate,
		Delete:   resourceIBMISSubnetNetworkACLAttachmentDelete,
		Exists:   resourceIBMISSubnetNetworkACLAttachmentExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isSubnetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier",
			},

			isNetworkACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of network ACL",
			},

			isNetworkACLName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL name",
			},

			isNetworkACLCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn for this Network ACL",
			},

			isNetworkACLVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL VPC",
			},

			isNetworkACLResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group ID for the network ACL",
			},

			isNetworkACLRules: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this Network ACL rule",
						},
						isNetworkACLRuleName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this rule",
						},
						isNetworkACLRuleAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to allow or deny matching traffic",
						},
						isNetworkACLRuleIPVersion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version for this rule",
						},
						isNetworkACLRuleSource: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source CIDR block",
						},
						isNetworkACLRuleDestination: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination CIDR block",
						},
						isNetworkACLRuleDirection: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direction of traffic to enforce, either inbound or outbound",
						},
						isNetworkACLRuleICMP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRuleICMPCode: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic code to allow",
									},
									isNetworkACLRuleICMPType: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic type to allow",
									},
								},
							},
						},

						isNetworkACLRuleTCP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of TCP destination port range",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of TCP destination port range",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of TCP source port range",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of TCP source port range",
									},
								},
							},
						},

						isNetworkACLRuleUDP: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of UDP destination port range",
									},
									isNetworkACLRulePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of UDP destination port range",
									},
									isNetworkACLRuleSourcePortMax: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of UDP source port range",
									},
									isNetworkACLRuleSourcePortMin: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of UDP source port range",
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

func resourceIBMISSubnetNetworkACLAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnet := d.Get(isSubnetID).(string)
	networkACL := d.Get(isNetworkACLID).(string)

	// Construct an instance of the NetworkACLIdentityByID model
	networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
	networkACLIdentityModel.ID = &networkACL

	// Construct an instance of the ReplaceSubnetNetworkACLOptions model
	replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
	replaceSubnetNetworkACLOptionsModel.ID = &subnet
	replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
	resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

	if err != nil {
		log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
		return fmt.Errorf("[ERROR] Error while attaching a network ACL to a subnet %s\n%s", err, response)
	}
	d.SetId(subnet)
	log.Printf("[INFO] Network ACL : %s", *resultACL.ID)
	log.Printf("[INFO] Subnet ID : %s", subnet)

	return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
}

func resourceIBMISSubnetNetworkACLAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSubnetNetworkACLOptionsModel := &vpcv1.GetSubnetNetworkACLOptions{
		ID: &id,
	}
	nwacl, response, err := sess.GetSubnetNetworkACL(getSubnetNetworkACLOptionsModel)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting subnet's (%s) attached network ACL: %s\n%s", id, err, response)
	}
	d.Set(isNetworkACLName, *nwacl.Name)
	d.Set(isNetworkACLCRN, *nwacl.CRN)
	d.Set(isNetworkACLVPC, *nwacl.VPC.ID)
	d.Set(isNetworkACLID, *nwacl.ID)
	if nwacl.ResourceGroup != nil {
		d.Set(isNetworkACLResourceGroup, *nwacl.ResourceGroup.ID)
	}

	rules := make([]interface{}, 0)
	if len(nwacl.Rules) > 0 {
		for _, rulex := range nwacl.Rules {
			log.Println("[DEBUG] Type of the Rule", reflect.TypeOf(rulex))
			r := make(map[string]interface{})
			switch rule := rulex.(type) {
			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
				setNetworkACLRuleCommonFields(r, rule.ID, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction)
				r[isNetworkACLRuleTCP] = []map[string]int{}
				r[isNetworkACLRuleUDP] = []map[string]int{}
				if rule.Code != nil && rule.Type != nil {
					r[isNetworkACLRuleICMP] = []map[string]int{{
						isNetworkACLRuleICMPCode: int(*rule.Code),
						isNetworkACLRuleICMPType: int(*rule.Type),
					}}
				} else {
					r[isNetworkACLRuleICMP] = []map[string]int{}
				}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
				setNetworkACLRuleCommonFields(r, rule.ID, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction)
				r[isNetworkACLRuleICMP] = []map[string]int{}

				if rule.Protocol != nil && *rule.Protocol == "tcp" {
					r[isNetworkACLRuleUDP] = []map[string]int{}
					r[isNetworkACLRuleTCP] = []map[string]int{{
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rule.SourcePortMin),
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rule.SourcePortMax),
						isNetworkACLRulePortMin:       checkNetworkACLNil(rule.DestinationPortMin),
						isNetworkACLRulePortMax:       checkNetworkACLNil(rule.DestinationPortMax),
					}}
				} else if rule.Protocol != nil && *rule.Protocol == "udp" {
					r[isNetworkACLRuleTCP] = []map[string]int{}
					r[isNetworkACLRuleUDP] = []map[string]int{{
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rule.SourcePortMin),
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rule.SourcePortMax),
						isNetworkACLRulePortMin:       checkNetworkACLNil(rule.DestinationPortMin),
						isNetworkACLRulePortMax:       checkNetworkACLNil(rule.DestinationPortMax),
					}}
				}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny:
				setNetworkACLRuleCommonFields(r, rule.ID, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction)
				r[isNetworkACLRuleICMP] = []map[string]int{}
				r[isNetworkACLRuleTCP] = []map[string]int{}
				r[isNetworkACLRuleUDP] = []map[string]int{}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp:
				setNetworkACLRuleCommonFields(r, rule.ID, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction)
				r[isNetworkACLRuleICMP] = []map[string]int{}
				r[isNetworkACLRuleTCP] = []map[string]int{}
				r[isNetworkACLRuleUDP] = []map[string]int{}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual:
				setNetworkACLRuleCommonFields(r, rule.ID, rule.Name, rule.Action, rule.IPVersion, rule.Source, rule.Destination, rule.Direction)
				r[isNetworkACLRuleICMP] = []map[string]int{}
				r[isNetworkACLRuleTCP] = []map[string]int{}
				r[isNetworkACLRuleUDP] = []map[string]int{}

			}

			rules = append(rules, r)
		}
	}
	d.Set(isNetworkACLRules, rules)
	return nil
}

func setNetworkACLRuleCommonFields(r map[string]interface{}, id, name, action, ipVersion, source, destination, direction *string) {
	r[isNetworkACLRuleID] = *id
	r[isNetworkACLRuleName] = *name
	r[isNetworkACLRuleAction] = *action
	r[isNetworkACLRuleIPVersion] = *ipVersion
	r[isNetworkACLRuleSource] = *source
	r[isNetworkACLRuleDestination] = *destination
	r[isNetworkACLRuleDirection] = *direction
}

func resourceIBMISSubnetNetworkACLAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isNetworkACLID) {
		subnet := d.Get(isSubnetID).(string)
		networkACL := d.Get(isNetworkACLID).(string)

		// Construct an instance of the NetworkACLIdentityByID model
		networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
		networkACLIdentityModel.ID = &networkACL

		// Construct an instance of the ReplaceSubnetNetworkACLOptions model
		replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
		replaceSubnetNetworkACLOptionsModel.ID = &subnet
		replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
		resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
			return fmt.Errorf("[ERROR] Error while attaching a network ACL to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with Network ACL : %s", subnet, *resultACL.ID)

		d.SetId(subnet)
		return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
	}

	return resourceIBMISSubnetNetworkACLAttachmentRead(d, meta)
}

func resourceIBMISSubnetNetworkACLAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	// Set the subnet with VPC default network ACL
	getSubnetOptions := &vpcv1.GetSubnetOptions{
		ID: &id,
	}
	subnet, response, err := sess.GetSubnet(getSubnetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Subnet (%s): %s\n%s", id, err, response)
	}
	// Fetch VPC
	vpcID := *subnet.VPC.ID

	getvpcOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getvpcOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting VPC : %s\n%s", err, response)
	}

	// Fetch default network ACL
	if vpc.DefaultNetworkACL != nil {
		log.Printf("[DEBUG] vpc default network acl is not null :%s", *vpc.DefaultNetworkACL.ID)
		// Construct an instance of the NetworkACLIdentityByID model
		networkACLIdentityModel := new(vpcv1.NetworkACLIdentityByID)
		networkACLIdentityModel.ID = vpc.DefaultNetworkACL.ID

		// Construct an instance of the ReplaceSubnetNetworkACLOptions model
		replaceSubnetNetworkACLOptionsModel := new(vpcv1.ReplaceSubnetNetworkACLOptions)
		replaceSubnetNetworkACLOptionsModel.ID = &id
		replaceSubnetNetworkACLOptionsModel.NetworkACLIdentity = networkACLIdentityModel
		resultACL, response, err := sess.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a network ACL to a subnet %s\n%s", err, response)
			return fmt.Errorf("[ERROR] Error while attaching a network ACL to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with VPC default Network ACL : %s", id, *resultACL.ID)
	} else {
		log.Printf("[DEBUG] vpc default network acl is  null")
	}

	d.SetId("")
	return nil
}

func resourceIBMISSubnetNetworkACLAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSubnetNetworkACLOptionsModel := &vpcv1.GetSubnetNetworkACLOptions{
		ID: &id,
	}
	_, response, err := sess.GetSubnetNetworkACL(getSubnetNetworkACLOptionsModel)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting subnet's attached network ACL: %s\n%s", err, response)
	}
	return true, nil
}
