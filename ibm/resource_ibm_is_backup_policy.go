// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func resourceIBMIsBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBackupPolicyCreate,
		ReadContext:   resourceIBMIsBackupPolicyRead,
		UpdateContext: resourceIBMIsBackupPolicyUpdate,
		DeleteContext: resourceIBMIsBackupPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"match_resource_types": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Set:         schema.HashString,
				Description: "A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"match_user_tags": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				// ValidateFunc: InvokeValidator("ibm_is_backup_policy", "match_user_tags"),
				Description: "The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.",
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_backup_policy", "match_user_tags")},
				Set:         schema.HashString,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_backup_policy", "name"),
				Description:  "The user-defined name for this backup policy. Names must be unique within the region this backup policy resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"plans": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The prototype objects for backup plans to be created for this backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"attach_user_tags": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							// ValidateFunc: InvokeValidator("ibm_is_backup_policy", "attach_user_tags"),
							Description: "User tags to attach to each resource created by this plan. If unspecified, no user tags will be attached.",
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_backup_policy", "attach_user_tags")},
						},
						"copy_user_tags": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							// Computed: true,
							Default:     true,
							Description: "Indicates whether to copy the source's user tags to the created resource.",
						},
						"cron_spec": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_backup_policy", "cron_spec"),
							Description:  "The cron specification for the backup schedule.",
						},
						"delete_after": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "The number of days to keep the backup.",
						},

						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_backup_policy", "name"),
							Description:  "The user-defined name for this backup policy plan. Names must be unique within the backup policy this plan resides in.",
						},

						"plan_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy plan.",
						},

						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced",
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: InvokeValidator("ibm_is_backup_policy", "id"),
							Description:  "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this backup policy.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy.",
			},
			"last_job_completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the most recent job for this backup policy completed.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the backup policy.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func resourceIBMIsBackupPolicyValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "cron_spec",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^((((\d+,)+\d+|([\d\*]+(\/|-)\d+)|\d+|\*) ?){5,7})$`,
			MinValueLength:             9,
			MaxValueLength:             63,
		},
	)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "id",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[0-9a-f]{32}$`,
			MinValueLength:             9,
			MaxValueLength:             63,
		},
	)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "attach_user_tags",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             0,
			MaxValueLength:             128,
		},
	)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "match_user_tags",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)
	resourceValidator := ResourceValidator{ResourceName: "ibm_is_backup_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBackupPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createBackupPolicyOptions := &vpcv1.CreateBackupPolicyOptions{}

	if _, ok := d.GetOk("match_resource_types"); ok {
		createBackupPolicyOptions.SetMatchResourceTypes(expandStringList((d.Get("match_resource_types").(*schema.Set)).List()))
		// matchResourceTypes := []string{}
		// for _, matchResourceType := range d.Get("match_resource_types").([]interface{}) {
		// 	matchResourceTypes = append(matchResourceTypes, matchResourceType.(string))
		// }
		// createBackupPolicyOptions.SetMatchResourceTypes(matchResourceTypes)
	}

	// matchUserTags := []string{}
	// for _, matchUserTag := range d.Get("match_user_tags").([]interface{}) {
	// 	matchUserTags = append(matchUserTags, matchUserTag.(string))
	// }
	// createBackupPolicyOptions.SetMatchUserTags(matchUserTags)
	createBackupPolicyOptions.SetMatchUserTags(expandStringList((d.Get("match_user_tags").(*schema.Set)).List()))

	if _, ok := d.GetOk("name"); ok {
		createBackupPolicyOptions.SetName(d.Get("name").(string))
	}

	if _, ok := d.GetOk("plans"); ok {
		var plans []vpcv1.BackupPolicyPlanPrototype
		for _, e := range d.Get("plans").([]interface{}) {
			value := e.(map[string]interface{})
			plansItem := resourceIBMIsBackupPolicyMapToBackupPolicyPlanPrototype(value)
			plans = append(plans, plansItem)
		}
		createBackupPolicyOptions.SetPlans(plans)
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroup := resourceIBMIsBackupPolicyMapToResourceGroupIdentity(d.Get("resource_group.0").(map[string]interface{}))
		// createBackupPolicyOptions.SetResourceGroup(&resourceGroup)
		createBackupPolicyOptions.ResourceGroup = resourceGroup
	}

	backupPolicy, response, err := sess.CreateBackupPolicyWithContext(context, createBackupPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateBackupPolicyWithContext failed %s\n%s", err, response))
	}

	d.SetId(*backupPolicy.ID)

	return resourceIBMIsBackupPolicyRead(context, d, meta)
}

func resourceIBMIsBackupPolicyMapToBackupPolicyPlanPrototype(backupPolicyPlanPrototypeMap map[string]interface{}) vpcv1.BackupPolicyPlanPrototype {
	backupPolicyPlanPrototype := vpcv1.BackupPolicyPlanPrototype{}

	if backupPolicyPlanPrototypeMap["attach_user_tags"] != nil {
		attachUserTags := []string{}
		for _, attachUserTagsItem := range backupPolicyPlanPrototypeMap["attach_user_tags"].([]interface{}) {
			attachUserTags = append(attachUserTags, attachUserTagsItem.(string))
		}
		backupPolicyPlanPrototype.AttachUserTags = attachUserTags
	}
	if backupPolicyPlanPrototypeMap["copy_user_tags"] != nil {
		backupPolicyPlanPrototype.CopyUserTags = core.BoolPtr(backupPolicyPlanPrototypeMap["copy_user_tags"].(bool))
	}
	backupPolicyPlanPrototype.CronSpec = core.StringPtr(backupPolicyPlanPrototypeMap["cron_spec"].(string))

	if backupPolicyPlanPrototypeMap["delete_after"] != nil {
		backupPolicyPlanDeletionTriggerPrototype := vpcv1.BackupPolicyPlanDeletionTriggerPrototype{}
		deleteAfter := int64(backupPolicyPlanPrototypeMap["delete_after"].(int))
		backupPolicyPlanDeletionTriggerPrototype.DeleteAfter = &deleteAfter
		backupPolicyPlanPrototype.DeletionTrigger = &backupPolicyPlanDeletionTriggerPrototype
	}
	if backupPolicyPlanPrototypeMap["name"] != nil {
		backupPolicyPlanPrototype.Name = core.StringPtr(backupPolicyPlanPrototypeMap["name"].(string))
	}

	return backupPolicyPlanPrototype
}

func resourceIBMIsBackupPolicyMapToResourceGroupIdentity(resourceGroupIdentityMap map[string]interface{}) vpcv1.ResourceGroupIdentityIntf {
	resourceGroupIdentity := vpcv1.ResourceGroupIdentity{}

	if resourceGroupIdentityMap["id"] != nil {
		resourceGroupIdentity.ID = core.StringPtr(resourceGroupIdentityMap["id"].(string))
	}

	return &resourceGroupIdentity
}

func resourceIBMIsBackupPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyOptions := &vpcv1.GetBackupPolicyOptions{}
	getBackupPolicyOptions.SetID(d.Id())

	backupPolicy, response, err := sess.GetBackupPolicyWithContext(context, getBackupPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBackupPolicyWithContext failed %s\n%s", err, response))
	}
	if backupPolicy.MatchResourceTypes != nil {
		if err = d.Set("match_resource_types", backupPolicy.MatchResourceTypes); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting match_resource_types: %s", err))
		}
	}
	if backupPolicy.MatchUserTags != nil {
		if err = d.Set("match_user_tags", backupPolicy.MatchUserTags); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting match_user_tags: %s", err))
		}
	}
	if err = d.Set("name", backupPolicy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if backupPolicy.Plans != nil {
		plans := []map[string]interface{}{}
		for _, plansItem := range backupPolicy.Plans {
			plansItemMap := resourceIBMIsBackupPolicyBackupPolicyPlanReferenceToMap(d, plansItem)
			plans = append(plans, plansItemMap)
		}
		if err = d.Set("plans", plans); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting plans: %s", err))
		}
	}
	if backupPolicy.ResourceGroup != nil {
		resourceGroupMap := resourceIBMIsBackupPolicyResourceGroupIdentityToMap(*backupPolicy.ResourceGroup)
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("created_at", dateTimeToString(backupPolicy.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	// if err = d.Set("last_job_completed_at", dateTimeToString(backupPolicy.LastJobCompletedAt)); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting last_job_completed_at: %s", err))
	// }
	if err = d.Set("crn", backupPolicy.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", backupPolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", backupPolicy.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", backupPolicy.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIBMIsBackupPolicyBackupPolicyPlanReferenceToMap(d *schema.ResourceData, backupPolicyPlanReference vpcv1.BackupPolicyPlanReference) map[string]interface{} {
	backupPolicyPlanReferenceMap := map[string]interface{}{}

	if backupPolicyPlanReference.ID != nil {
		backupPolicyPlanReferenceMap["plan_id"] = backupPolicyPlanReference.ID
	}
	if backupPolicyPlanReference.Name != nil {
		backupPolicyPlanReferenceMap["name"] = backupPolicyPlanReference.Name
		if _, ok := d.GetOk("plans"); ok {
			for _, e := range d.Get("plans").([]interface{}) {
				backupPolicyPlanPrototypeMap := e.(map[string]interface{})
				if backupPolicyPlanPrototypeMap["name"] != nil && backupPolicyPlanPrototypeMap["name"].(string) == *backupPolicyPlanReference.Name {
					backupPolicyPlanReferenceMap["attach_user_tags"] = backupPolicyPlanPrototypeMap["attach_user_tags"]
					backupPolicyPlanReferenceMap["copy_user_tags"] = backupPolicyPlanPrototypeMap["copy_user_tags"]
					backupPolicyPlanReferenceMap["delete_after"] = backupPolicyPlanPrototypeMap["delete_after"]
					backupPolicyPlanReferenceMap["cron_spec"] = core.StringPtr(backupPolicyPlanPrototypeMap["cron_spec"].(string))
				}
			}
		}
	}
	if backupPolicyPlanReference.ResourceType != nil {
		backupPolicyPlanReferenceMap["resource_type"] = backupPolicyPlanReference.ResourceType
	}

	return backupPolicyPlanReferenceMap
}

func resourceIBMIsBackupPolicyResourceGroupIdentityToMap(resourceGroupReference vpcv1.ResourceGroupReference) map[string]interface{} {
	resourceGroupIdentityMap := map[string]interface{}{}

	if resourceGroupReference.ID != nil {
		resourceGroupIdentityMap["id"] = resourceGroupReference.ID
	}
	if resourceGroupReference.Name != nil {
		resourceGroupIdentityMap["name"] = resourceGroupReference.Name
	}
	return resourceGroupIdentityMap
}

func resourceIBMIsBackupPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBackupPolicyOptions := &vpcv1.UpdateBackupPolicyOptions{}

	updateBackupPolicyOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vpcv1.BackupPolicyPatch{}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		updateBackupPolicyOptions.BackupPolicyPatch, _ = patchVals.AsPatch()
		_, response, err := sess.UpdateBackupPolicyWithContext(context, updateBackupPolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBackupPolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateBackupPolicyWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsBackupPolicyRead(context, d, meta)
}

func resourceIBMIsBackupPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackupPolicyOptions := &vpcv1.DeleteBackupPolicyOptions{}

	deleteBackupPolicyOptions.SetID(d.Id())

	_, response, err := sess.DeleteBackupPolicyWithContext(context, deleteBackupPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBackupPolicyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
