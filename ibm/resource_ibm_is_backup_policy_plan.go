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

func resourceIBMIsBackupPolicyPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBackupPolicyPlanCreate,
		ReadContext:   resourceIBMIsBackupPolicyPlanRead,
		UpdateContext: resourceIBMIsBackupPolicyPlanUpdate,
		DeleteContext: resourceIBMIsBackupPolicyPlanDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The backup policy identifier.",
			},
			"cron_spec": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_backup_policy_plan", "cron_spec"),
				Description:  "The cron specification for the backup schedule.",
			},
			"active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the plan is active.",
			},
			"attach_user_tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User tags to attach to each resource created by this plan. If unspecified, no user tags will be attached.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"copy_user_tags": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether to copy the source's user tags to the created resource.",
			},
			"deletion_trigger": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_after": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of days to keep the backup.",
						},
						"delete_after_backup_count": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of latest backup to be retained.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_backup_policy_plan", "name"),
				Description:  "The user-defined name for this backup policy plan. Names must be unique within the backup policy this plan resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy plan was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy plan.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this backup policy plan.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func resourceIBMIsBackupPolicyPlanValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "cron_spec",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^((((\\d+,)+\\d+|([\\d\\*]+(\/|-)\\d+)|\\d+|\\*) ?){5,7})$`,
			MinValueLength:             9,
			MaxValueLength:             63,
		},
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

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_backup_policy_plan", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBackupPolicyPlanCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	createBackupPolicyPlanOptions := &vpcv1.CreateBackupPolicyPlanOptions{}

	createBackupPolicyPlanOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))
	createBackupPolicyPlanOptions.SetCronSpec(d.Get("cron_spec").(string))
	if _, ok := d.GetOk("active"); ok {
		createBackupPolicyPlanOptions.SetActive(d.Get("active").(bool))
	}
	if _, ok := d.GetOk("attach_user_tags"); ok {
		createBackupPolicyPlanOptions.SetAttachUserTags(d.Get("attach_user_tags").([]string))
	}
	if _, ok := d.GetOk("copy_user_tags"); ok {
		createBackupPolicyPlanOptions.SetCopyUserTags(d.Get("copy_user_tags").(bool))
	}
	if _, ok := d.GetOk("deletion_trigger"); ok {
		deletionTrigger := resourceIBMIsBackupPolicyPlanMapToBackupPolicyPlanDeletionTrigger(d.Get("deletion_trigger.0").(map[string]interface{}))
		createBackupPolicyPlanOptions.SetDeletionTrigger(deletionTrigger)
	}
	if _, ok := d.GetOk("name"); ok {
		createBackupPolicyPlanOptions.SetName(d.Get("name").(string))
	}

	backupPolicyPlan, response, err := sess.CreateBackupPolicyPlanWithContext(context, createBackupPolicyPlanOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createBackupPolicyPlanOptions.BackupPolicyID, *backupPolicyPlan.ID))

	return resourceIBMIsBackupPolicyPlanRead(context, d, meta)
}

func resourceIBMIsBackupPolicyPlanMapToBackupPolicyPlanDeletionTrigger(backupPolicyPlanDeletionTriggerMap map[string]interface{}) vpcv1.BackupPolicyPlanDeletionTriggerIntf {
	backupPolicyPlanDeletionTrigger := vpcv1.BackupPolicyPlanDeletionTrigger{}

	if backupPolicyPlanDeletionTriggerMap["delete_after"] != nil {
		backupPolicyPlanDeletionTrigger.DeleteAfter = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerMap["delete_after"].(int)))
	}
	if backupPolicyPlanDeletionTriggerMap["delete_after_backup_count"] != nil {
		backupPolicyPlanDeletionTrigger.DeleteAfterBackupCount = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerMap["delete_after_backup_count"].(int)))
	}

	return &backupPolicyPlanDeletionTrigger
}

func resourceIBMIsBackupPolicyPlanMapToBackupPolicyPlanDeletionTriggerByAge(backupPolicyPlanDeletionTriggerByAgeMap map[string]interface{}) vpcv1.BackupPolicyPlanDeletionTriggerByAge {
	backupPolicyPlanDeletionTriggerByAge := vpcv1.BackupPolicyPlanDeletionTriggerByAge{}

	if backupPolicyPlanDeletionTriggerByAgeMap["delete_after"] != nil {
		backupPolicyPlanDeletionTriggerByAge.DeleteAfter = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerByAgeMap["delete_after"].(int)))
	}

	return backupPolicyPlanDeletionTriggerByAge
}

func resourceIBMIsBackupPolicyPlanMapToBackupPolicyPlanDeletionTriggerByCount(backupPolicyPlanDeletionTriggerByCountMap map[string]interface{}) vpcv1.BackupPolicyPlanDeletionTriggerByCount {
	backupPolicyPlanDeletionTriggerByCount := vpcv1.BackupPolicyPlanDeletionTriggerByCount{}

	if backupPolicyPlanDeletionTriggerByCountMap["delete_after_backup_count"] != nil {
		backupPolicyPlanDeletionTriggerByCount.DeleteAfterBackupCount = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerByCountMap["delete_after_backup_count"].(int)))
	}

	return backupPolicyPlanDeletionTriggerByCount
}

func resourceIBMIsBackupPolicyPlanRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	getBackupPolicyPlanOptions := &vpcv1.GetBackupPolicyPlanOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	getBackupPolicyPlanOptions.SetID(parts[1])

	backupPolicyPlan, response, err := sess.GetBackupPolicyPlanWithContext(context, getBackupPolicyPlanOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("backup_policy_id", getBackupPolicyPlanOptions.BackupPolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting backup_policy_id: %s", err))
	}
	if err = d.Set("cron_spec", backupPolicyPlan.CronSpec); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cron_spec: %s", err))
	}
	if err = d.Set("active", backupPolicyPlan.Active); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting active: %s", err))
	}
	if backupPolicyPlan.AttachUserTags != nil {
		if err = d.Set("attach_user_tags", backupPolicyPlan.AttachUserTags); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting attach_user_tags: %s", err))
		}
	}
	if err = d.Set("copy_user_tags", backupPolicyPlan.CopyUserTags); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting copy_user_tags: %s", err))
	}
	if backupPolicyPlan.DeletionTrigger != nil {
		// deletionTriggerMap := resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerToMap(*backupPolicyPlan.DeletionTrigger)
		// if err = d.Set("deletion_trigger", []map[string]interface{}{deletionTriggerMap}); err != nil {
		// 	return diag.FromErr(fmt.Errorf("Error setting deletion_trigger: %s", err))
		// }
	}
	if err = d.Set("name", backupPolicyPlan.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(backupPolicyPlan.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", backupPolicyPlan.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", backupPolicyPlan.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", backupPolicyPlan.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerToMap(backupPolicyPlanDeletionTrigger vpcv1.BackupPolicyPlanDeletionTriggerIntf) map[string]interface{} {
	backupPolicyPlanDeletionTriggerMap := map[string]interface{}{}

	// TODO: Add code here to convert a vpcv1.BackupPolicyPlanDeletionTriggerIntf to map[string]interface{}

	return backupPolicyPlanDeletionTriggerMap
}

func resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerByAgeToMap(backupPolicyPlanDeletionTriggerByAge vpcv1.BackupPolicyPlanDeletionTriggerByAge) map[string]interface{} {
	backupPolicyPlanDeletionTriggerByAgeMap := map[string]interface{}{}

	if backupPolicyPlanDeletionTriggerByAge.DeleteAfter != nil {
		backupPolicyPlanDeletionTriggerByAgeMap["delete_after"] = intValue(backupPolicyPlanDeletionTriggerByAge.DeleteAfter)
	}

	return backupPolicyPlanDeletionTriggerByAgeMap
}

func resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerByCountToMap(backupPolicyPlanDeletionTriggerByCount vpcv1.BackupPolicyPlanDeletionTriggerByCount) map[string]interface{} {
	backupPolicyPlanDeletionTriggerByCountMap := map[string]interface{}{}

	if backupPolicyPlanDeletionTriggerByCount.DeleteAfterBackupCount != nil {
		backupPolicyPlanDeletionTriggerByCountMap["delete_after_backup_count"] = intValue(backupPolicyPlanDeletionTriggerByCount.DeleteAfterBackupCount)
	}

	return backupPolicyPlanDeletionTriggerByCountMap
}

func resourceIBMIsBackupPolicyPlanUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBackupPolicyPlanOptions := &vpcv1.UpdateBackupPolicyPlanOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	updateBackupPolicyPlanOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.BackupPolicyPlanPatch{}
	if d.HasChange("cron_spec") {
		patchVals.CronSpec = core.StringPtr(d.Get("cron_spec").(string))
		hasChange = true
	}
	if d.HasChange("active") {
		patchVals.Active = core.BoolPtr(d.Get("active").(bool))
		hasChange = true
	}
	if d.HasChange("attach_user_tags") {
		// TODO: handle AttachUserTags of type TypeList -- not primitive, not model
		hasChange = true
	}
	if d.HasChange("copy_user_tags") {
		patchVals.CopyUserTags = core.BoolPtr(d.Get("copy_user_tags").(bool))
		hasChange = true
	}
	if d.HasChange("deletion_trigger") {
		deletionTrigger := resourceIBMIsBackupPolicyPlanMapToBackupPolicyPlanDeletionTrigger(d.Get("deletion_trigger.0").(map[string]interface{}))
		patchVals.DeletionTrigger = deletionTrigger
		hasChange = true
	}
	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		updateBackupPolicyPlanOptions.BackupPolicyPlanPatch, _ = patchVals.AsPatch()
		_, response, err := sess.UpdateBackupPolicyPlanWithContext(context, updateBackupPolicyPlanOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBackupPolicyPlanWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateBackupPolicyPlanWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsBackupPolicyPlanRead(context, d, meta)
}

func resourceIBMIsBackupPolicyPlanDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackupPolicyPlanOptions := &vpcv1.DeleteBackupPolicyPlanOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	deleteBackupPolicyPlanOptions.SetID(parts[1])

	_, response, err := sess.DeleteBackupPolicyPlanWithContext(context, deleteBackupPolicyPlanOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
