// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsClusterNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsClusterNetworkCreate,
		ReadContext:   resourceIBMIsClusterNetworkRead,
		UpdateContext: resourceIBMIsClusterNetworkUpdate,
		DeleteContext: resourceIBMIsClusterNetworkDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_cluster_network", "name"),
				Description:  "The name for this cluster network. The name must not be used by another cluster network in the region.",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The profile for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL for this cluster network profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The globally unique name for this cluster network profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The resource group for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			"subnet_prefixes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The IP address ranges available for subnets for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_policy": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The allocation policy for this subnet prefix:- `auto`: Subnets created by total count in this cluster network can use this prefix.",
						},
						"cidr": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CIDR block for this prefix.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The VPC this cluster network resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
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
							Required:    true,
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The zone this cluster network resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL for this zone.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this cluster network.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network.",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the cluster network.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsClusterNetworkValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_cluster_network", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsClusterNetworkCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClusterNetworkOptions := &vpcv1.CreateClusterNetworkOptions{}

	profileModel, err := ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentity(d.Get("profile.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createClusterNetworkOptions.SetProfile(profileModel)
	vpcModel, err := ResourceIBMIsClusterNetworkMapToVPCIdentity(d.Get("vpc.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createClusterNetworkOptions.SetVPC(vpcModel)
	zoneModel, err := ResourceIBMIsClusterNetworkMapToZoneIdentity(d.Get("zone.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createClusterNetworkOptions.SetZone(zoneModel)
	if _, ok := d.GetOk("name"); ok {
		createClusterNetworkOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroupModel, err := ResourceIBMIsClusterNetworkMapToResourceGroupIdentity(d.Get("resource_group.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createClusterNetworkOptions.SetResourceGroup(resourceGroupModel)
	}
	if _, ok := d.GetOk("subnet_prefixes"); ok {
		var subnetPrefixes []vpcv1.ClusterNetworkSubnetPrefixPrototype
		for _, v := range d.Get("subnet_prefixes").([]interface{}) {
			value := v.(map[string]interface{})
			subnetPrefixesItem, err := ResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			subnetPrefixes = append(subnetPrefixes, *subnetPrefixesItem)
		}
		createClusterNetworkOptions.SetSubnetPrefixes(subnetPrefixes)
	}

	clusterNetwork, _, err := vpcClient.CreateClusterNetworkWithContext(context, createClusterNetworkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*clusterNetwork.ID)

	return resourceIBMIsClusterNetworkRead(context, d, meta)
}

func resourceIBMIsClusterNetworkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

	getClusterNetworkOptions.SetID(d.Id())

	clusterNetwork, response, err := vpcClient.GetClusterNetworkWithContext(context, getClusterNetworkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(clusterNetwork.Name) {
		if err = d.Set("name", clusterNetwork.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	profileMap, err := ResourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(clusterNetwork.Profile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("profile", []map[string]interface{}{profileMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile: %s", err))
	}
	if !core.IsNil(clusterNetwork.ResourceGroup) {
		resourceGroupMap, err := ResourceIBMIsClusterNetworkResourceGroupReferenceToMap(clusterNetwork.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if !core.IsNil(clusterNetwork.SubnetPrefixes) {
		subnetPrefixes := []map[string]interface{}{}
		for _, subnetPrefixesItem := range clusterNetwork.SubnetPrefixes {
			subnetPrefixesItemMap, err := ResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(&subnetPrefixesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			subnetPrefixes = append(subnetPrefixes, subnetPrefixesItemMap)
		}
		if err = d.Set("subnet_prefixes", subnetPrefixes); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subnet_prefixes: %s", err))
		}
	}
	vpcMap, err := ResourceIBMIsClusterNetworkVPCReferenceToMap(clusterNetwork.VPC)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vpc", []map[string]interface{}{vpcMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc: %s", err))
	}
	zoneMap, err := ResourceIBMIsClusterNetworkZoneReferenceToMap(clusterNetwork.Zone)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("zone", []map[string]interface{}{zoneMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(clusterNetwork.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", clusterNetwork.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", clusterNetwork.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetwork.LifecycleReasons {
		lifecycleReasonsItemMap, err := ResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(&lifecycleReasonsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_reasons: %s", err))
	}
	if err = d.Set("lifecycle_state", clusterNetwork.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", clusterNetwork.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_is_cluster_network", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIBMIsClusterNetworkUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClusterNetworkOptions := &vpcv1.UpdateClusterNetworkOptions{}

	updateClusterNetworkOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vpcv1.ClusterNetworkPatch{}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	updateClusterNetworkOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateClusterNetworkOptions.ClusterNetworkPatch, _ = patchVals.AsPatch()
		_, _, err = vpcClient.UpdateClusterNetworkWithContext(context, updateClusterNetworkOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsClusterNetworkRead(context, d, meta)
}

func resourceIBMIsClusterNetworkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_cluster_network", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClusterNetworkOptions := &vpcv1.DeleteClusterNetworkOptions{}

	deleteClusterNetworkOptions.SetID(d.Id())

	_, _, err = vpcClient.DeleteClusterNetworkWithContext(context, deleteClusterNetworkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClusterNetworkWithContext failed: %s", err.Error()), "ibm_is_cluster_network", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentity(modelMap map[string]interface{}) (vpcv1.ClusterNetworkProfileIdentityIntf, error) {
	model := &vpcv1.ClusterNetworkProfileIdentity{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByName(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkProfileIdentityByName, error) {
	model := &vpcv1.ClusterNetworkProfileIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkProfileIdentityByHref, error) {
	model := &vpcv1.ClusterNetworkProfileIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentity(modelMap map[string]interface{}) (vpcv1.VPCIdentityIntf, error) {
	model := &vpcv1.VPCIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByID(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByID, error) {
	model := &vpcv1.VPCIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByCRN, error) {
	model := &vpcv1.VPCIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToVPCIdentityByHref(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByHref, error) {
	model := &vpcv1.VPCIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToZoneIdentity(modelMap map[string]interface{}) (vpcv1.ZoneIdentityIntf, error) {
	model := &vpcv1.ZoneIdentity{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToZoneIdentityByName(modelMap map[string]interface{}) (*vpcv1.ZoneIdentityByName, error) {
	model := &vpcv1.ZoneIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToZoneIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ZoneIdentityByHref, error) {
	model := &vpcv1.ZoneIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetPrefixPrototype, error) {
	model := &vpcv1.ClusterNetworkSubnetPrefixPrototype{}
	if modelMap["cidr"] != nil && modelMap["cidr"].(string) != "" {
		model.CIDR = core.StringPtr(modelMap["cidr"].(string))
	}
	return model, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(model *vpcv1.ClusterNetworkProfileReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(model *vpcv1.ClusterNetworkSubnetPrefix) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["allocation_policy"] = *model.AllocationPolicy
	modelMap["cidr"] = *model.CIDR
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsClusterNetworkVPCReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkVPCReferenceDeletedToMap(model *vpcv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(model *vpcv1.ClusterNetworkLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
