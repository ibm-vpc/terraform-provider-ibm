// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isPublicAddressRangeAuthorizedCIDRDeleting  = "deleting"
	isPublicAddressRangeAuthorizedCIDRDeleted   = "deleted"
	isPublicAddressRangeAuthorizedCIDRAvailable = "stable"
	isPublicAddressRangeAuthorizedCIDRFailed    = "failed"
	isPublicAddressRangeAuthorizedCIDRPending   = "pending"
	isPublicAddressRangeAuthorizedCIDRSuspended = "suspended"
	isPublicAddressRangeAuthorizedCIDRUpdating  = "updating"
	isPublicAddressRangeAuthorizedCIDRWaiting   = "waiting"
)

func ResourceIBMIsPublicAddressRangeAuthorizedCIDR() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsPublicAddressRangeAuthorizedCIDRCreate,
		ReadContext:   resourceIBMIsPublicAddressRangeAuthorizedCIDRRead,
		UpdateContext: resourceIBMIsPublicAddressRangeAuthorizedCIDRUpdate,
		DeleteContext: resourceIBMIsPublicAddressRangeAuthorizedCIDRDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ip_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_public_address_range_authorized_cidr", "ip_version"),
				Description:  "The IP version for this public address range authorized CIDR. Currently only `ipv6` is supported.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name for this public address range authorized CIDR. The name must not be used by another public address range authorized CIDR in the region. Names beginning with `ibm-` are reserved for provider-managed resources, and are not allowed.",
			},
			"availability_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_public_address_range_authorized_cidr", "availability_mode"),
				Description:  "The availability mode of the public address range authorized CIDR. Currently only `zonal` is supported.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name of the zone this public address range authorized CIDR will reside in.",
			},
			"network_prefix_length": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_public_address_range_authorized_cidr", "network_prefix_length"),
				Description:  "The network prefix length for this public address range authorized CIDR.",
			},
			// Computed fields
			"allocation": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The allocation for this public address range authorized CIDR.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of resources allocated from this public address range authorized CIDR.",
						},
						"profile_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile family for resources allocated from this public address range authorized CIDR.",
						},
					},
				},
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP address block for the public address range authorized CIDR, expressed in CIDR format.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this public address range authorized CIDR.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this public address range authorized CIDR.",
			},
			"lifecycle_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the public address range authorized CIDR.",
			},
			"resource_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this public address range authorized CIDR.",
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
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func ResourceIBMIsPublicAddressRangeAuthorizedCIDRValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "ip_version",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "ipv6",
		},
		validate.ValidateSchema{
			Identifier:                 "availability_mode",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "zonal",
		},
		validate.ValidateSchema{
			Identifier:                 "network_prefix_length",
			ValidateFunctionIdentifier: validate.ValidateAllowedIntValue,
			Type:                       validate.TypeInt,
			Required:                   true,
			AllowedValues:              "64",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_public_address_range_authorized_cidr", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	prototype := &vpcv1.PublicAddressRangeAuthorizedCIDRPrototype{}
	prototype.IPVersion = core.StringPtr(d.Get("ip_version").(string))

	if v, ok := d.GetOk("name"); ok {
		prototype.Name = core.StringPtr(v.(string))
	}
	prototype.AvailabilityMode = core.StringPtr(d.Get("availability_mode").(string))
	prototype.Zone = &vpcv1.ZoneIdentity{Name: core.StringPtr(d.Get("zone").(string))}
	prototype.NetworkPrefixLength = core.Int64Ptr(int64(d.Get("network_prefix_length").(int)))

	createOptions := &vpcv1.CreatePublicAddressRangeAuthorizedCIDROptions{
		PublicAddressRangeAuthorizedCIDRPrototype: prototype,
	}

	authorizedCIDR, _, err := vpcClient.CreatePublicAddressRangeAuthorizedCIDRWithContext(context, createOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreatePublicAddressRangeAuthorizedCIDRWithContext failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*authorizedCIDR.ID)
	log.Printf("[INFO] PublicAddressRangeAuthorizedCIDR : %s", *authorizedCIDR.ID)

	_, err = isWaitForPublicAddressRangeAuthorizedCIDRAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForPublicAddressRangeAuthorizedCIDRAvailable failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMIsPublicAddressRangeAuthorizedCIDRRead(context, d, meta)
}

func isPublicAddressRangeAuthorizedCIDRRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{
			ID: &id,
		}
		authorizedCIDR, response, err := sess.GetPublicAddressRangeAuthorizedCIDR(getOptions)
		if err != nil {
			return nil, isPublicAddressRangeAuthorizedCIDRFailed, fmt.Errorf("[ERROR] Error getting PublicAddressRangeAuthorizedCIDR : %s\n%s", err, response)
		}

		if *authorizedCIDR.LifecycleState == isPublicAddressRangeAuthorizedCIDRAvailable {
			return authorizedCIDR, *authorizedCIDR.LifecycleState, nil
		} else if *authorizedCIDR.LifecycleState == isPublicAddressRangeAuthorizedCIDRFailed {
			return authorizedCIDR, *authorizedCIDR.LifecycleState, fmt.Errorf("PublicAddressRangeAuthorizedCIDR (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted PublicAddressRangeAuthorizedCIDR and attempt to create the PublicAddressRangeAuthorizedCIDR again replacing the previous configuration", *authorizedCIDR.ID)
		}

		return authorizedCIDR, isPublicAddressRangeAuthorizedCIDRPending, nil
	}
}

func isWaitForPublicAddressRangeAuthorizedCIDRAvailable(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for public address range authorized CIDR (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isPublicAddressRangeAuthorizedCIDRPending, isPublicAddressRangeAuthorizedCIDRWaiting, isPublicAddressRangeAuthorizedCIDRSuspended},
		Target:     []string{isPublicAddressRangeAuthorizedCIDRAvailable, isPublicAddressRangeAuthorizedCIDRFailed},
		Refresh:    isPublicAddressRangeAuthorizedCIDRRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{}
	getOptions.SetID(d.Id())

	authorizedCIDR, response, err := vpcClient.GetPublicAddressRangeAuthorizedCIDRWithContext(context, getOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetPublicAddressRangeAuthorizedCIDRWithContext failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("ip_version", authorizedCIDR.IPVersion); err != nil {
		err = fmt.Errorf("Error setting ip_version: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-ip_version").GetDiag()
	}
	if !core.IsNil(authorizedCIDR.Name) {
		if err = d.Set("name", authorizedCIDR.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(authorizedCIDR.ResourceGroup) {
		rgMap, err := resourceIBMIsPublicAddressRangeAuthorizedCIDRResourceGroupReferenceToMap(authorizedCIDR.ResourceGroup)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "resource_group-to-map").GetDiag()
		}
		if err = d.Set("resource_group", []map[string]interface{}{rgMap}); err != nil {
			err = fmt.Errorf("Error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-resource_group").GetDiag()
		}
	}
	if err = d.Set("availability_mode", authorizedCIDR.AvailabilityMode); err != nil {
		err = fmt.Errorf("Error setting availability_mode: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-availability_mode").GetDiag()
	}
	if !core.IsNil(authorizedCIDR.Zone) {
		if err = d.Set("zone", authorizedCIDR.Zone.Name); err != nil {
			err = fmt.Errorf("Error setting zone: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-zone").GetDiag()
		}
	}
	if err = d.Set("network_prefix_length", flex.IntValue(authorizedCIDR.NetworkPrefixLength)); err != nil {
		err = fmt.Errorf("Error setting network_prefix_length: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-network_prefix_length").GetDiag()
	}
	if !core.IsNil(authorizedCIDR.Allocation) {
		allocationMap, err := resourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationToMap(authorizedCIDR.Allocation)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "allocation-to-map").GetDiag()
		}
		if err = d.Set("allocation", []map[string]interface{}{allocationMap}); err != nil {
			err = fmt.Errorf("Error setting allocation: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-allocation").GetDiag()
		}
	}
	if err = d.Set("cidr", authorizedCIDR.CIDR); err != nil {
		err = fmt.Errorf("Error setting cidr: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-cidr").GetDiag()
	}
	if err = d.Set("crn", authorizedCIDR.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-crn").GetDiag()
	}
	if err = d.Set("href", authorizedCIDR.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-href").GetDiag()
	}
	if authorizedCIDR.LifecycleReasons != nil {
		lifecycleReasonsMap, err := resourceIBMIsPublicAddressRangeAuthorizedCIDRLifecycleReasonsToMap(authorizedCIDR.LifecycleReasons)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		if err = d.Set("lifecycle_reasons", lifecycleReasonsMap); err != nil {
			err = fmt.Errorf("Error setting lifecycle_reasons: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-lifecycle_reasons").GetDiag()
		}
	}
	if err = d.Set("lifecycle_state", authorizedCIDR.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("resource_type", authorizedCIDR.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "read", "set-resource_type").GetDiag()
	}

	return nil
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if d.HasChange("name") {
		updateOptions := &vpcv1.UpdatePublicAddressRangeAuthorizedCIDROptions{}
		updateOptions.SetID(d.Id())

		patchVals := &vpcv1.PublicAddressRangeAuthorizedCIDRPatch{}
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		updateOptions.PublicAddressRangeAuthorizedCIDRPatch, _ = patchVals.AsPatch()

		_, _, err = vpcClient.UpdatePublicAddressRangeAuthorizedCIDRWithContext(context, updateOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdatePublicAddressRangeAuthorizedCIDRWithContext failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		_, err = isWaitForPublicAddressRangeAuthorizedCIDRUpdate(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForPublicAddressRangeAuthorizedCIDRUpdate failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsPublicAddressRangeAuthorizedCIDRRead(context, d, meta)
}

func isWaitForPublicAddressRangeAuthorizedCIDRUpdate(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PublicAddressRangeAuthorizedCIDR (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isPublicAddressRangeAuthorizedCIDRUpdating, isPublicAddressRangeAuthorizedCIDRWaiting, isPublicAddressRangeAuthorizedCIDRSuspended},
		Target:     []string{isPublicAddressRangeAuthorizedCIDRAvailable, isPublicAddressRangeAuthorizedCIDRFailed},
		Refresh:    isPublicAddressRangeAuthorizedCIDRUpdateRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isPublicAddressRangeAuthorizedCIDRUpdateRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{
			ID: &id,
		}
		authorizedCIDR, response, err := sess.GetPublicAddressRangeAuthorizedCIDR(getOptions)
		if err != nil {
			return nil, isPublicAddressRangeAuthorizedCIDRFailed, fmt.Errorf("[ERROR] Error getting PublicAddressRangeAuthorizedCIDR : %s\n%s", err, response)
		}

		if *authorizedCIDR.LifecycleState == isPublicAddressRangeAuthorizedCIDRAvailable || *authorizedCIDR.LifecycleState == isPublicAddressRangeAuthorizedCIDRFailed {
			return authorizedCIDR, *authorizedCIDR.LifecycleState, nil
		}

		return authorizedCIDR, isPublicAddressRangeAuthorizedCIDRUpdating, nil
	}
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_public_address_range_authorized_cidr", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteOptions := &vpcv1.DeletePublicAddressRangeAuthorizedCIDROptions{}
	deleteOptions.SetID(d.Id())

	_, err = vpcClient.DeletePublicAddressRangeAuthorizedCIDRWithContext(context, deleteOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeletePublicAddressRangeAuthorizedCIDRWithContext failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	_, err = isWaitForPublicAddressRangeAuthorizedCIDRDeleted(vpcClient, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForPublicAddressRangeAuthorizedCIDRDeleted failed: %s", err.Error()), "ibm_is_public_address_range_authorized_cidr", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func isWaitForPublicAddressRangeAuthorizedCIDRDeleted(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for PublicAddressRangeAuthorizedCIDR (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isPublicAddressRangeAuthorizedCIDRDeleting},
		Target:     []string{isPublicAddressRangeAuthorizedCIDRDeleted, isPublicAddressRangeAuthorizedCIDRFailed},
		Refresh:    isPublicAddressRangeAuthorizedCIDRDeleteRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPublicAddressRangeAuthorizedCIDRDeleteRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh function for PublicAddressRangeAuthorizedCIDR delete.")
		getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{
			ID: &id,
		}
		authorizedCIDR, response, err := sess.GetPublicAddressRangeAuthorizedCIDR(getOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return authorizedCIDR, isPublicAddressRangeAuthorizedCIDRDeleted, nil
			}
			return nil, isPublicAddressRangeAuthorizedCIDRFailed, fmt.Errorf("[ERROR] The PublicAddressRangeAuthorizedCIDR %s failed to delete: %s\n%s", id, err, response)
		}
		return authorizedCIDR, *authorizedCIDR.LifecycleState, nil
	}
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDRAllocation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["count"] = flex.IntValue(model.Count)
	modelMap["profile_family"] = *model.ProfileFamily
	return modelMap, nil
}

func resourceIBMIsPublicAddressRangeAuthorizedCIDRLifecycleReasonsToMap(reasons []vpcv1.PublicAddressRangeAuthorizedCIDRLifecycleReason) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(reasons))
	for _, reason := range reasons {
		reasonMap := make(map[string]interface{})
		reasonMap["code"] = *reason.Code
		reasonMap["message"] = *reason.Message
		if !core.IsNil(reason.MoreInfo) {
			reasonMap["more_info"] = *reason.MoreInfo
		}
		result = append(result, reasonMap)
	}
	return result, nil
}
