// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isNetworkACLName              = "name"
	isNetworkACLRules             = "rules"
	isNetworkACLSubnets           = "subnets"
	isNetworkACLRuleID            = "id"
	isNetworkACLRuleName          = "name"
	isNetworkACLRuleAction        = "action"
	isNetworkACLRuleIPVersion     = "ip_version"
	isNetworkACLRuleSource        = "source"
	isNetworkACLRuleDestination   = "destination"
	isNetworkACLRuleDirection     = "direction"
	isNetworkACLRuleProtocol      = "protocol"
	isNetworkACLRuleICMP          = "icmp"
	isNetworkACLRuleICMPCode      = "code"
	isNetworkACLRuleICMPType      = "type"
	isNetworkACLRuleTCP           = "tcp"
	isNetworkACLRuleUDP           = "udp"
	isNetworkACLRulePortMax       = "port_max"
	isNetworkACLRulePortMin       = "port_min"
	isNetworkACLRuleSourcePortMax = "source_port_max"
	isNetworkACLRuleSourcePortMin = "source_port_min"
	isNetworkACLVPC               = "vpc"
	isNetworkACLResourceGroup     = "resource_group"
	isNetworkACLTags              = "tags"
	isNetworkACLAccessTags        = "access_tags"
	isNetworkACLCRN               = "crn"
)

func ResourceIBMISNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISNetworkACLCreate,
		Read:     resourceIBMISNetworkACLRead,
		Update:   resourceIBMISNetworkACLUpdate,
		Delete:   resourceIBMISNetworkACLDelete,
		Exists:   resourceIBMISNetworkACLExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{
			isNetworkACLName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLName),
				Description:  "Network ACL name",
			},
			isNetworkACLVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Network ACL VPC name",
			},
			isNetworkACLResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Resource group ID for the network ACL",
			},
			isNetworkACLTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},

			isNetworkACLAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},

			isNetworkACLCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			isNetworkACLRules: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNetworkACLRuleID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isNetworkACLRuleName: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleName),
						},
						isNetworkACLRuleAction: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleAction),
						},
						isNetworkACLRuleIPVersion: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isNetworkACLRuleSource: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSource),
						},
						isNetworkACLRuleDestination: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleDestination),
						},
						isNetworkACLRuleDirection: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							Description:  "Direction of traffic to enforce, either inbound or outbound",
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleDirection),
						},
						isNetworkACLSubnets: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isNetworkACLRuleProtocol: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The name of the network protocol",
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleProtocol),
						},
						isNetworkACLRuleICMPCode: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleICMPCode),
						},
						isNetworkACLRuleICMPType: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleICMPType),
						},
						isNetworkACLRulePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMax),
						},
						isNetworkACLRulePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMin),
						},
						isNetworkACLRuleSourcePortMax: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMax),
						},
						isNetworkACLRuleSourcePortMin: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMin),
						},
						isNetworkACLRuleICMP: {
							Type:       schema.TypeList,
							MinItems:   0,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							Deprecated: "icmp is deprecated, use 'protocol', 'code', and 'type' instead.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRuleICMPCode: {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleICMPCode),
									},
									isNetworkACLRuleICMPType: {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleICMPType),
									},
								},
							},
						},

						isNetworkACLRuleTCP: {
							Type:       schema.TypeList,
							MinItems:   0,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							Deprecated: "tcp is deprecated, use 'protocol', 'port_min', 'port_max', 'source_port_min', and 'source_port_max' instead.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      65535,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMax),
									},
									isNetworkACLRulePortMin: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMin),
									},
									isNetworkACLRuleSourcePortMax: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      65535,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMax),
									},
									isNetworkACLRuleSourcePortMin: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMin),
									},
								},
							},
						},

						isNetworkACLRuleUDP: {
							Type:       schema.TypeList,
							MinItems:   0,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							Deprecated: "udp is deprecated, use 'protocol', 'port_min', 'port_max', 'source_port_min', and 'source_port_max' instead.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRulePortMax: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      65535,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMax),
									},
									isNetworkACLRulePortMin: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRulePortMin),
									},
									isNetworkACLRuleSourcePortMax: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      65535,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMax),
									},
									isNetworkACLRuleSourcePortMin: {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleSourcePortMin),
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

func ResourceIBMISNetworkACLValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	direction := "inbound, outbound"
	protocol := "tcp, udp, icmp, ah, any, esp, gre, icmp_tcp_udp, ip_in_ip, l2tp, number_0, number_10, number_100, number_101, number_102, number_103, number_104, number_105, number_106, number_107, number_108, number_109, number_11, number_110, number_111, number_113, number_114, number_116, number_117, number_118, number_119, number_12, number_120, number_121, number_122, number_123, number_124, number_125, number_126, number_127, number_128, number_129, number_13, number_130, number_131, number_133, number_134, number_135, number_136, number_137, number_138, number_139, number_14, number_140, number_141, number_142, number_143, number_144, number_145, number_146, number_147, number_148, number_149, number_15, number_150, number_151, number_152, number_153, number_154, number_155, number_156, number_157, number_158, number_159, number_16, number_160, number_161, number_162, number_163, number_164, number_165, number_166, number_167, number_168, number_169, number_170, number_171, number_172, number_173, number_174, number_175, number_176, number_177, number_178, number_179, number_18, number_180, number_181, number_182, number_183, number_184, number_185, number_186, number_187, number_188, number_189, number_19, number_190, number_191, number_192, number_193, number_194, number_195, number_196, number_197, number_198, number_199, number_2, number_20, number_200, number_201, number_202, number_203, number_204, number_205, number_206, number_207, number_208, number_209, number_21, number_210, number_211, number_212, number_213, number_214, number_215, number_216, number_217, number_218, number_219, number_22, number_220, number_221, number_222, number_223, number_224, number_225, number_226, number_227, number_228, number_229, number_23, number_230, number_231, number_232, number_233, number_234, number_235, number_236, number_237, number_238, number_239, number_24, number_240, number_241, number_242, number_243, number_244, number_245, number_246, number_247, number_248, number_249, number_25, number_250, number_251, number_252, number_253, number_254, number_255, number_26, number_27, number_28, number_29, number_3, number_30, number_31, number_32, number_33, number_34, number_35, number_36, number_37, number_38, number_39, number_40, number_41, number_42, number_43, number_44, number_45, number_48, number_49, number_5, number_52, number_53, number_54, number_55, number_56, number_57, number_58, number_59, number_60, number_61, number_62, number_63, number_64, number_65, number_66, number_67, number_68, number_69, number_7, number_70, number_71, number_72, number_73, number_74, number_75, number_76, number_77, number_78, number_79, number_8, number_80, number_81, number_82, number_83, number_84, number_85, number_86, number_87, number_88, number_89, number_9, number_90, number_91, number_92, number_93, number_94, number_95, number_96, number_97, number_98, number_99, rsvp, sctp, vrrp"
	action := "allow, deny"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              action})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleDirection,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              direction})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   false,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleDestination,
			ValidateFunctionIdentifier: validate.ValidateIPorCIDR,
			Type:                       validate.TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSource,
			ValidateFunctionIdentifier: validate.ValidateIPorCIDR,
			Type:                       validate.TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleICMPType,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "0",
			MaxValue:                   "254"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleICMPCode,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "0",
			MaxValue:                   "255"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRulePortMin,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRulePortMax,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSourcePortMin,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isNetworkACLRuleSourcePortMax,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSecurityGroupRuleProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              protocol})

	ibmISNetworkACLResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_network_acl", Schema: validateSchema}
	return &ibmISNetworkACLResourceValidator
}

func resourceIBMISNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isNetworkACLName).(string)
	err := nwaclCreate(d, meta, name)
	if err != nil {
		return err
	}
	return resourceIBMISNetworkACLRead(d, meta)

}

func nwaclCreate(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	var vpc, rg string
	if vpcID, ok := d.GetOk(isNetworkACLVPC); ok {
		vpc = vpcID.(string)
	} else {
		return fmt.Errorf("[ERROR] Required parameter vpc is not set")
	}

	nwaclTemplate := &vpcv1.NetworkACLPrototype{
		VPC: &vpcv1.VPCIdentity{
			ID: &vpc,
		},
	}
	if name != "" {
		nwaclTemplate.Name = &name
	}

	if grp, ok := d.GetOk(isNetworkACLResourceGroup); ok {
		rg = grp.(string)
		nwaclTemplate.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	// validate each rule before attempting to create the ACL
	var rules []interface{}
	if rls, ok := d.GetOk(isNetworkACLRules); ok {
		rules = rls.([]interface{})
	}
	err = validateInlineRules(rules)
	if err != nil {
		return err
	}

	options := &vpcv1.CreateNetworkACLOptions{
		NetworkACLPrototype: nwaclTemplate,
	}

	nwacl, response, err := sess.CreateNetworkACL(options)
	if err != nil {
		return fmt.Errorf("[DEBUG]Error while creating Network ACL err %s\n%s", err, response)
	}
	d.SetId(*nwacl.ID)
	log.Printf("[INFO] Network ACL : %s", *nwacl.ID)
	nwaclid := *nwacl.ID

	//Remove default rules
	err = clearRules(sess, nwaclid)
	if err != nil {
		return err
	}

	err = createInlineRules(d, sess, nwaclid, rules)
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isNetworkACLTags); ok || v != "" {
		oldList, newList := d.GetChange(isNetworkACLTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *nwacl.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource network acl (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isNetworkACLAccessTags); ok {
		oldList, newList := d.GetChange(isNetworkACLAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *nwacl.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource network acl (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMISNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	err := nwaclGet(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func nwaclGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getNetworkAclOptions := &vpcv1.GetNetworkACLOptions{
		ID: &id,
	}
	nwacl, response, err := sess.GetNetworkACL(getNetworkAclOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Network ACL(%s) : %s\n%s", id, err, response)
	}
	d.Set(isNetworkACLName, *nwacl.Name)
	d.Set(isNetworkACLVPC, *nwacl.VPC.ID)
	if nwacl.ResourceGroup != nil {
		d.Set(isNetworkACLResourceGroup, *nwacl.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *nwacl.ResourceGroup.Name)
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *nwacl.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource network acl (%s) tags: %s", d.Id(), err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *nwacl.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource network acl (%s) access tags: %s", d.Id(), err)
	}

	d.Set(isNetworkACLTags, tags)
	d.Set(isNetworkACLAccessTags, accesstags)
	d.Set(isNetworkACLCRN, *nwacl.CRN)
	rules := make([]interface{}, 0)
	if len(nwacl.Rules) > 0 {
		for _, rulex := range nwacl.Rules {
			log.Println("[DEBUG] Type of the Rule", reflect.TypeOf(rulex))
			rule := make(map[string]interface{})
			rule[isNetworkACLSubnets] = len(nwacl.Subnets)
			switch r := rulex.(type) {
			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
				setCommonNetworkACLRuleFields(rule, r.ID, r.Name, r.Action, r.IPVersion, r.Source, r.Destination, r.Direction, r.Protocol)
				rule[isNetworkACLRuleTCP] = []map[string]int{}
				rule[isNetworkACLRuleUDP] = []map[string]int{}
				icmp := []map[string]int{}
				if r.Code != nil && r.Type != nil {
					icmp = append(icmp, map[string]int{
						isNetworkACLRuleICMPCode: int(*r.Code),
						isNetworkACLRuleICMPType: int(*r.Type),
					})
					rule[isNetworkACLRuleICMPCode] = int(*r.Code)
					rule[isNetworkACLRuleICMPType] = int(*r.Type)
				}
				rule[isNetworkACLRuleICMP] = icmp

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
				setCommonNetworkACLRuleFields(rule, r.ID, r.Name, r.Action, r.IPVersion, r.Source, r.Destination, r.Direction, r.Protocol)
				rule[isNetworkACLRuleICMP] = []map[string]int{}
				if r.Protocol != nil {
					switch *r.Protocol {
					case "tcp":
						rule[isNetworkACLRuleICMP] = []map[string]int{}
						rule[isNetworkACLRuleUDP] = []map[string]int{}
						tcp := []map[string]int{{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(r.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(r.SourcePortMin),
							isNetworkACLRulePortMax:       checkNetworkACLNil(r.DestinationPortMax),
							isNetworkACLRulePortMin:       checkNetworkACLNil(r.DestinationPortMin),
						}}
						rule[isNetworkACLRuleTCP] = tcp
						rule[isNetworkACLRuleSourcePortMax] = checkNetworkACLNil(r.SourcePortMax)
						rule[isNetworkACLRuleSourcePortMin] = checkNetworkACLNil(r.SourcePortMin)
						rule[isNetworkACLRulePortMax] = checkNetworkACLNil(r.DestinationPortMax)
						rule[isNetworkACLRulePortMin] = checkNetworkACLNil(r.DestinationPortMin)

					case "udp":
						rule[isNetworkACLRuleICMP] = []map[string]int{}
						rule[isNetworkACLRuleTCP] = []map[string]int{}
						udp := []map[string]int{{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(r.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(r.SourcePortMin),
							isNetworkACLRulePortMax:       checkNetworkACLNil(r.DestinationPortMax),
							isNetworkACLRulePortMin:       checkNetworkACLNil(r.DestinationPortMin),
						}}
						rule[isNetworkACLRuleUDP] = udp
						rule[isNetworkACLRuleSourcePortMax] = checkNetworkACLNil(r.SourcePortMax)
						rule[isNetworkACLRuleSourcePortMin] = checkNetworkACLNil(r.SourcePortMin)
						rule[isNetworkACLRulePortMax] = checkNetworkACLNil(r.DestinationPortMax)
						rule[isNetworkACLRulePortMin] = checkNetworkACLNil(r.DestinationPortMin)
					}
				}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny:
				setCommonNetworkACLRuleFields(rule, r.ID, r.Name, r.Action, r.IPVersion, r.Source, r.Destination, r.Direction, r.Protocol)
				rule[isNetworkACLRuleICMP] = []map[string]int{}
				rule[isNetworkACLRuleTCP] = []map[string]int{}
				rule[isNetworkACLRuleUDP] = []map[string]int{}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp:
				setCommonNetworkACLRuleFields(rule, r.ID, r.Name, r.Action, r.IPVersion, r.Source, r.Destination, r.Direction, r.Protocol)
				rule[isNetworkACLRuleICMP] = []map[string]int{}
				rule[isNetworkACLRuleTCP] = []map[string]int{}
				rule[isNetworkACLRuleUDP] = []map[string]int{}

			case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual:
				setCommonNetworkACLRuleFields(rule, r.ID, r.Name, r.Action, r.IPVersion, r.Source, r.Destination, r.Direction, r.Protocol)
				rule[isNetworkACLRuleICMP] = []map[string]int{}
				rule[isNetworkACLRuleTCP] = []map[string]int{}
				rule[isNetworkACLRuleUDP] = []map[string]int{}
			}

			rules = append(rules, rule)
		}
	}
	d.Set(isNetworkACLRules, rules)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/acl")
	d.Set(flex.ResourceName, *nwacl.Name)
	// d.Set(flex.ResourceCRN, *nwacl.Crn)
	return nil
}

func setCommonNetworkACLRuleFields(rule map[string]interface{}, id, name, action, ipVersion, source, destination, direction, protocol *string) {
	rule[isNetworkACLRuleID] = *id
	rule[isNetworkACLRuleName] = *name
	rule[isNetworkACLRuleAction] = *action
	rule[isNetworkACLRuleIPVersion] = *ipVersion
	rule[isNetworkACLRuleSource] = *source
	rule[isNetworkACLRuleDestination] = *destination
	rule[isNetworkACLRuleDirection] = *direction
	rule[isNetworkACLRuleProtocol] = *protocol
}

func resourceIBMISNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	name := ""
	hasChanged := false

	if d.HasChange(isNetworkACLName) {
		name = d.Get(isNetworkACLName).(string)
		hasChanged = true
	}

	err := nwaclUpdate(d, meta, id, name, hasChanged)
	if err != nil {
		return err
	}
	return resourceIBMISNetworkACLRead(d, meta)
}

func nwaclUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	rules := d.Get(isNetworkACLRules).([]interface{})
	if hasChanged {
		updateNetworkACLOptions := &vpcv1.UpdateNetworkACLOptions{
			ID: &id,
		}
		networkACLPatchModel := &vpcv1.NetworkACLPatch{
			Name: &name,
		}
		networkACLPatch, err := networkACLPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for NetworkACLPatch: %s", err)
		}
		updateNetworkACLOptions.NetworkACLPatch = networkACLPatch
		_, response, err := sess.UpdateNetworkACL(updateNetworkACLOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Network ACL(%s) : %s\n%s", id, err, response)
		}
	}
	if d.HasChange(isNetworkACLTags) {
		oldList, newList := d.GetChange(isNetworkACLTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get(isNetworkACLCRN).(string), "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource network acl (%s) tags: %s", d.Id(), err)
		}
	}
	if d.HasChange(isNetworkACLAccessTags) {
		oldList, newList := d.GetChange(isNetworkACLAccessTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get(isNetworkACLCRN).(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource network acl (%s) access tags: %s", d.Id(), err)
		}
	}
	if d.HasChange(isNetworkACLRules) {
		err := validateInlineRules(rules)
		if err != nil {
			return err
		}
		//Delete all existing rules
		err = clearRules(sess, id)
		if err != nil {
			return err
		}
		//Create the rules as per the def
		err = createInlineRules(d, sess, id, rules)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMISNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	err := nwaclDelete(d, meta, id)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func nwaclDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getNetworkAclOptions := &vpcv1.GetNetworkACLOptions{
		ID: &id,
	}
	_, response, err := sess.GetNetworkACL(getNetworkAclOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Network ACL (%s): %s\n%s", id, err, response)
	}

	deleteNetworkAclOptions := &vpcv1.DeleteNetworkACLOptions{
		ID: &id,
	}
	response, err = sess.DeleteNetworkACL(deleteNetworkAclOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Network ACL : %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISNetworkACLExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	exists, err := nwaclExists(d, meta, id)
	return exists, err
}

func nwaclExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getNetworkAclOptions := &vpcv1.GetNetworkACLOptions{
		ID: &id,
	}
	_, response, err := sess.GetNetworkACL(getNetworkAclOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Network ACL: %s\n%s", err, response)
	}
	return true, nil
}

func checkNetworkACLNil(ptr *int64) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

func clearRules(nwaclC *vpcv1.VpcV1, nwaclid string) error {
	start := ""
	allrecs := []vpcv1.NetworkACLRuleItemIntf{}
	for {
		listNetworkAclRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
			NetworkACLID: &nwaclid,
		}
		if start != "" {
			listNetworkAclRulesOptions.Start = &start
		}
		rawrules, response, err := nwaclC.ListNetworkACLRules(listNetworkAclRulesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Listing network ACL rules : %s\n%s", err, response)
		}
		start = flex.GetNext(rawrules.Next)
		allrecs = append(allrecs, rawrules.Rules...)
		if start == "" {
			break
		}
	}

	for _, rule := range allrecs {
		deleteNetworkAclRuleOptions := &vpcv1.DeleteNetworkACLRuleOptions{
			NetworkACLID: &nwaclid,
		}
		switch reflect.TypeOf(rule).String() {
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp":
			rule := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp)
			deleteNetworkAclRuleOptions.ID = rule.ID
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp":
			rule := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp)
			deleteNetworkAclRuleOptions.ID = rule.ID
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny":
			rule := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAny)
			deleteNetworkAclRuleOptions.ID = rule.ID
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp":
			rule := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmptcpudp)
			deleteNetworkAclRuleOptions.ID = rule.ID
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual":
			rule := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIndividual)
			deleteNetworkAclRuleOptions.ID = rule.ID

		}

		response, err := nwaclC.DeleteNetworkACLRule(deleteNetworkAclRuleOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deleting network ACL rule : %s\n%s", err, response)
		}
	}
	return nil
}

func validateInlineRules(rules []interface{}) error {
	for _, rule := range rules {
		rulex := rule.(map[string]interface{})
		action := rulex[isNetworkACLRuleAction].(string)
		if (action != "allow") && (action != "deny") {
			return fmt.Errorf("[ERROR] Invalid action. valid values are allow|deny")
		}

		direction := rulex[isNetworkACLRuleDirection].(string)
		direction = strings.ToLower(direction)

		icmp := len(rulex[isNetworkACLRuleICMP].([]interface{})) > 0
		tcp := len(rulex[isNetworkACLRuleTCP].([]interface{})) > 0
		udp := len(rulex[isNetworkACLRuleUDP].([]interface{})) > 0

		if (icmp && tcp) || (icmp && udp) || (tcp && udp) {
			return fmt.Errorf("Only one of icmp|tcp|udp can be defined per rule")
		}

	}
	return nil
}

func createInlineRules(d *schema.ResourceData, nwaclC *vpcv1.VpcV1, nwaclid string, rules []interface{}) error {
	before := ""

	for i := 0; i <= len(rules)-1; i++ {
		rulex := rules[i].(map[string]interface{})

		name := rulex[isNetworkACLRuleName].(string)
		source := rulex[isNetworkACLRuleSource].(string)
		destination := rulex[isNetworkACLRuleDestination].(string)
		action := rulex[isNetworkACLRuleAction].(string)
		direction := rulex[isNetworkACLRuleDirection].(string)
		icmp := rulex[isNetworkACLRuleICMP].([]interface{})
		tcp := rulex[isNetworkACLRuleTCP].([]interface{})
		udp := rulex[isNetworkACLRuleUDP].([]interface{})
		icmptype := int64(-1)
		icmpcode := int64(-1)
		minport := int64(-1)
		maxport := int64(-1)
		sourceminport := int64(-1)
		sourcemaxport := int64(-1)
		protocol := "icmp_tcp_udp"
		if action == "deny" {
			protocol = "any"
		}
		fmt.Println("Default protocol ", protocol)
		if protocolVal, ok := rulex[isNetworkACLRuleProtocol]; ok {
			if str, ok := protocolVal.(string); ok && str != "" {
				protocol = str
			}
		}
		fmt.Println("After Default protocol ", protocol)
		ruleTemplate := &vpcv1.NetworkACLRulePrototype{
			Action:      &action,
			Destination: &destination,
			Direction:   &direction,
			Source:      &source,
			Name:        &name,
		}

		if before != "" {
			ruleTemplate.Before = &vpcv1.NetworkACLRuleBeforePrototype{
				ID: &before,
			}
		}
		if protocol != "icmp" {
			if _, ok := d.GetOk("type"); ok {
				return fmt.Errorf("attribute 'type' conflicts with protocol %s; 'type' is only valid for icmp protocol", protocol)
			}
			if _, ok := d.GetOk("code"); ok {
				return fmt.Errorf("attribute 'code' conflicts with protocol %q; 'code' is only valid for icmp protocol", protocol)
			}
		}

		if protocol != "tcp" && protocol != "udp" {
			if _, ok := d.GetOk("port_min"); ok {
				return fmt.Errorf("attribute 'port_min' conflicts with protocol %s; ports apply only to tcp/udp protocol", protocol)
			}
			if _, ok := d.GetOk("port_max"); ok {
				return fmt.Errorf("attribute 'port_max' conflicts with protocol %s; ports apply only to tcp/udp protocol", protocol)
			}
			if _, ok := d.GetOk("source_port_max"); ok {
				return fmt.Errorf("attribute 'source_port_max' conflicts with protocol %s; ports apply only to tcp/udp protocol", protocol)
			}
			if _, ok := d.GetOk("source_port_min"); ok {
				return fmt.Errorf("attribute 'source_port_min' conflicts with protocol %s; ports apply only to tcp/udp protocol", protocol)
			}
		}
		if protocol == "icmp" {
			ruleTemplate.Protocol = &protocol
			if val, ok := rulex["type"]; ok {
				icmptype = int64(val.(int))
				ruleTemplate.Type = &icmptype
			}
			if val, ok := rulex["code"]; ok {
				icmpcode = int64(val.(int))
				ruleTemplate.Code = &icmpcode
			}
		} else if len(icmp) > 0 {
			protocol = "icmp"
			ruleTemplate.Protocol = &protocol
			if !isNil(icmp[0]) {
				icmpval := icmp[0].(map[string]interface{})
				if val, ok := icmpval[isNetworkACLRuleICMPType]; ok {
					icmptype = int64(val.(int))
					ruleTemplate.Type = &icmptype
				}
				if val, ok := icmpval[isNetworkACLRuleICMPCode]; ok {
					icmpcode = int64(val.(int))
					ruleTemplate.Code = &icmpcode
				}
			}
		}

		if protocol == "tcp" {
			ruleTemplate.Protocol = &protocol
			if val, ok := rulex[isNetworkACLRulePortMin]; ok {
				minport = int64(val.(int))
				ruleTemplate.DestinationPortMin = &minport
			}
			if val, ok := rulex[isNetworkACLRulePortMax]; ok {
				maxport = int64(val.(int))
				ruleTemplate.DestinationPortMax = &maxport
			}
			if val, ok := rulex[isNetworkACLRuleSourcePortMin]; ok {
				sourceminport = int64(val.(int))
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if val, ok := rulex[isNetworkACLRuleSourcePortMax]; ok {
				sourcemaxport = int64(val.(int))
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
			// Adding to default values
			if minport == -1 || minport == 0 {
				minPort := int64(1)
				ruleTemplate.DestinationPortMin = &minPort
			}
			if maxport == -1 || minport == 0 {
				maxport := int64(65535)
				ruleTemplate.DestinationPortMax = &maxport
			}
			if sourceminport == -1 || minport == 0 {
				sourceminport := int64(1)
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if sourcemaxport == -1 || minport == 0 {
				sourcemaxport := int64(65535)
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
		} else if len(tcp) > 0 {
			protocol = "tcp"
			ruleTemplate.Protocol = &protocol
			tcpval := tcp[0].(map[string]interface{})
			if val, ok := tcpval[isNetworkACLRulePortMin]; ok {
				minport = int64(val.(int))
				ruleTemplate.DestinationPortMin = &minport
			}
			if val, ok := tcpval[isNetworkACLRulePortMax]; ok {
				maxport = int64(val.(int))
				ruleTemplate.DestinationPortMax = &maxport
			}
			if val, ok := tcpval[isNetworkACLRuleSourcePortMin]; ok {
				sourceminport = int64(val.(int))
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if val, ok := tcpval[isNetworkACLRuleSourcePortMax]; ok {
				sourcemaxport = int64(val.(int))
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
		}

		if protocol == "udp" {
			ruleTemplate.Protocol = &protocol
			if val, ok := rulex[isNetworkACLRulePortMin]; ok {
				minport = int64(val.(int))
				ruleTemplate.DestinationPortMin = &minport
			}
			if val, ok := rulex[isNetworkACLRulePortMax]; ok {
				maxport = int64(val.(int))
				ruleTemplate.DestinationPortMax = &maxport
			}
			if val, ok := rulex[isNetworkACLRuleSourcePortMin]; ok {
				sourceminport = int64(val.(int))
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if val, ok := rulex[isNetworkACLRuleSourcePortMax]; ok {
				sourcemaxport = int64(val.(int))
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
			// Adding to default values
			if minport == -1 || minport == 0 {
				minPort := int64(1)
				ruleTemplate.DestinationPortMin = &minPort
			}
			if maxport == -1 || minport == 0 {
				maxport := int64(65535)
				ruleTemplate.DestinationPortMax = &maxport
			}
			if sourceminport == -1 || minport == 0 {
				sourceminport := int64(1)
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if sourcemaxport == -1 || minport == 0 {
				sourcemaxport := int64(65535)
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
		} else if len(udp) > 0 {
			protocol = "udp"
			ruleTemplate.Protocol = &protocol
			udpval := udp[0].(map[string]interface{})
			if val, ok := udpval[isNetworkACLRulePortMin]; ok {
				minport = int64(val.(int))
				ruleTemplate.DestinationPortMin = &minport
			}
			if val, ok := udpval[isNetworkACLRulePortMax]; ok {
				maxport = int64(val.(int))
				ruleTemplate.DestinationPortMax = &maxport
			}
			if val, ok := udpval[isNetworkACLRuleSourcePortMin]; ok {
				sourceminport = int64(val.(int))
				ruleTemplate.SourcePortMin = &sourceminport
			}
			if val, ok := udpval[isNetworkACLRuleSourcePortMax]; ok {
				sourcemaxport = int64(val.(int))
				ruleTemplate.SourcePortMax = &sourcemaxport
			}
		}

		fmt.Println("Before setting protocol ", protocol)
		ruleTemplate.Protocol = &protocol

		createNetworkAclRuleOptions := &vpcv1.CreateNetworkACLRuleOptions{
			NetworkACLID:            &nwaclid,
			NetworkACLRulePrototype: ruleTemplate,
		}
		_, response, err := nwaclC.CreateNetworkACLRule(createNetworkAclRuleOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating network ACL rule : %s\n%s", err, response)
		}
	}
	return nil
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
