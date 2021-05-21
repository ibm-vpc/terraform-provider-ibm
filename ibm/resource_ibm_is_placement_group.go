/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func resourceIbmIsPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsPlacementGroupCreate,
		ReadContext:   resourceIbmIsPlacementGroupRead,
		UpdateContext: resourceIbmIsPlacementGroupUpdate,
		DeleteContext: resourceIbmIsPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_placement_group", "strategy"),
				Description:  "The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_placement_group", "name"),
				Description:  "The unique user-defined name for this placement group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the placement group was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this placement group.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this placement group.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the placement group.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func resourceIbmIsPlacementGroupValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "strategy",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "host_spread, power_spread",
		},
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_placement_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsPlacementGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createPlacementGroupOptions := &vpcv1.CreatePlacementGroupOptions{}

	createPlacementGroupOptions.SetStrategy(d.Get("strategy").(string))
	if pgnameIntf, ok := d.GetOk("name"); ok {
		createPlacementGroupOptions.SetName(pgnameIntf.(string))
	}
	if resourceGroupIntf, ok := d.GetOk("resource_group"); ok {
		resourceGroup := resourceGroupIntf.(string)
		resourceGroupIdentity := &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroup,
		}
		createPlacementGroupOptions.SetResourceGroup(resourceGroupIdentity)
	}

	placementGroup, response, err := vpcClient.CreatePlacementGroupWithContext(context, createPlacementGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*placementGroup.ID)

	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{}

	getPlacementGroupOptions.SetID(d.Id())

	placementGroup, response, err := vpcClient.GetPlacementGroupWithContext(context, getPlacementGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("strategy", placementGroup.Strategy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy: %s", err))
	}
	if err = d.Set("name", placementGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if placementGroup.ResourceGroup != nil {
		if err = d.Set("resource_group", *placementGroup.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("created_at", placementGroup.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", placementGroup.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", placementGroup.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", placementGroup.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", placementGroup.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIbmIsPlacementGroupResourceGroupIdentityToMap(resourceGroupIdentity vpcv1.ResourceGroupIdentity) map[string]interface{} {
	resourceGroupIdentityMap := map[string]interface{}{}

	resourceGroupIdentityMap["id"] = resourceGroupIdentity.ID

	return resourceGroupIdentityMap
}

func resourceIbmIsPlacementGroupResourceGroupIdentityByIDToMap(resourceGroupIdentityByID vpcv1.ResourceGroupIdentityByID) map[string]interface{} {
	resourceGroupIdentityByIDMap := map[string]interface{}{}

	resourceGroupIdentityByIDMap["id"] = resourceGroupIdentityByID.ID

	return resourceGroupIdentityByIDMap
}

func resourceIbmIsPlacementGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updatePlacementGroupOptions := &vpcv1.UpdatePlacementGroupOptions{}

	updatePlacementGroupOptions.SetID(d.Id())

	hasChange := false

	placementGroupPatchModel := &vpcv1.PlacementGroupPatch{}
	if d.HasChange("name") {
		plName := d.Get("name").(string)
		placementGroupPatchModel.Name = &plName
		hasChange = true
	}
	if hasChange {
		placementGroupPatch, err := placementGroupPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling AsPatch for PlacementGroupPatch %s", err)
			return diag.FromErr(err)
		}
		updatePlacementGroupOptions.SetPlacementGroupPatch(placementGroupPatch)
		_, response, err := vpcClient.UpdatePlacementGroupWithContext(context, updatePlacementGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdatePlacementGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{}

	deletePlacementGroupOptions.SetID(d.Id())

	response, err := vpcClient.DeletePlacementGroupWithContext(context, deletePlacementGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] DeletePlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
