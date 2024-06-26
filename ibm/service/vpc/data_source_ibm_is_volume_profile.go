// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVolumeProfile       = "name"
	isVolumeProfileFamily = "family"
)

func DataSourceIBMISVolumeProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVolumeProfileRead,

		Schema: map[string]*schema.Schema{

			isVolumeProfile: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Volume profile name",
			},
			"boot_capacity": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"capacity": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"family": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this volume profile belongs to.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this volume profile.",
			},
			"iops": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"unattached_capacity_update_supported": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The value for this profile field.",
						},
					},
				},
			},
			"unattached_iops_update_supported": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The value for this profile field.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISVolumeProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	name := d.Get(isVolumeProfile).(string)

	err := volumeProfileGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func volumeProfileGet(d *schema.ResourceData, meta interface{}, name string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getVolumeProfileOptions := &vpcv1.GetVolumeProfileOptions{
		Name: &name,
	}
	volumeProfile, _, err := sess.GetVolumeProfile(getVolumeProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVolumeProfileWithContext failed: %s", err.Error()), "(Data) ibm_is_volume_profile", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	// For lack of anything better, compose our id from profile name.
	d.SetId(*volumeProfile.Name)

	bootCapacity := []map[string]interface{}{}
	if volumeProfile.BootCapacity != nil {
		modelMap, err := DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityToMap(volumeProfile.BootCapacity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
			return tfErr.GetDiag()
		}
		bootCapacity = append(bootCapacity, modelMap)
	}
	if err = d.Set("boot_capacity", bootCapacity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting boot_capacity: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	capacity := []map[string]interface{}{}
	if volumeProfile.Capacity != nil {
		modelMap, err := DataSourceIBMIsVolumeProfileVolumeProfileCapacityToMap(volumeProfile.Capacity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
			return tfErr.GetDiag()
		}
		capacity = append(capacity, modelMap)
	}
	if err = d.Set("capacity", capacity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting capacity: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("family", volumeProfile.Family); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting family: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", volumeProfile.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	iops := []map[string]interface{}{}
	if volumeProfile.Iops != nil {
		modelMap, err := DataSourceIBMIsVolumeProfileVolumeProfileIopsToMap(volumeProfile.Iops)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
			return tfErr.GetDiag()
		}
		iops = append(iops, modelMap)
	}
	if err = d.Set("iops", iops); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting iops: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	unattachedCapacityUpdateSupported := []map[string]interface{}{}
	if volumeProfile.UnattachedCapacityUpdateSupported != nil {
		modelMap, err := DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedToMap(volumeProfile.UnattachedCapacityUpdateSupported)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
			return tfErr.GetDiag()
		}
		unattachedCapacityUpdateSupported = append(unattachedCapacityUpdateSupported, modelMap)
	}
	if err = d.Set("unattached_capacity_update_supported", unattachedCapacityUpdateSupported); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting unattached_capacity_update_supported: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	unattachedIopsUpdateSupported := []map[string]interface{}{}
	if volumeProfile.UnattachedIopsUpdateSupported != nil {
		modelMap, err := DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedToMap(volumeProfile.UnattachedIopsUpdateSupported)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_volume_profile", "read")
			return tfErr.GetDiag()
		}
		unattachedIopsUpdateSupported = append(unattachedIopsUpdateSupported, modelMap)
	}
	if err = d.Set("unattached_iops_update_supported", unattachedIopsUpdateSupported); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting unattached_iops_update_supported: %s", err), "(Data) ibm_is_volume_profile", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityToMap(model vpcv1.VolumeProfileBootCapacityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileBootCapacityFixed); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityFixedToMap(model.(*vpcv1.VolumeProfileBootCapacityFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityRangeToMap(model.(*vpcv1.VolumeProfileBootCapacityRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityEnum); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityEnumToMap(model.(*vpcv1.VolumeProfileBootCapacityEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacityDependentRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityDependentRangeToMap(model.(*vpcv1.VolumeProfileBootCapacityDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileBootCapacity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileBootCapacity)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileBootCapacityIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityFixedToMap(model *vpcv1.VolumeProfileBootCapacityFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityRangeToMap(model *vpcv1.VolumeProfileBootCapacityRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityEnumToMap(model *vpcv1.VolumeProfileBootCapacityEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileBootCapacityDependentRangeToMap(model *vpcv1.VolumeProfileBootCapacityDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileCapacityToMap(model vpcv1.VolumeProfileCapacityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileCapacityFixed); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileCapacityFixedToMap(model.(*vpcv1.VolumeProfileCapacityFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileCapacityRangeToMap(model.(*vpcv1.VolumeProfileCapacityRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityEnum); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileCapacityEnumToMap(model.(*vpcv1.VolumeProfileCapacityEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacityDependentRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileCapacityDependentRangeToMap(model.(*vpcv1.VolumeProfileCapacityDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileCapacity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileCapacity)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileCapacityIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfileVolumeProfileCapacityFixedToMap(model *vpcv1.VolumeProfileCapacityFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileCapacityRangeToMap(model *vpcv1.VolumeProfileCapacityRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileCapacityEnumToMap(model *vpcv1.VolumeProfileCapacityEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileCapacityDependentRangeToMap(model *vpcv1.VolumeProfileCapacityDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileIopsToMap(model vpcv1.VolumeProfileIopsIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileIopsFixed); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileIopsFixedToMap(model.(*vpcv1.VolumeProfileIopsFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileIopsRangeToMap(model.(*vpcv1.VolumeProfileIopsRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsEnum); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileIopsEnumToMap(model.(*vpcv1.VolumeProfileIopsEnum))
	} else if _, ok := model.(*vpcv1.VolumeProfileIopsDependentRange); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileIopsDependentRangeToMap(model.(*vpcv1.VolumeProfileIopsDependentRange))
	} else if _, ok := model.(*vpcv1.VolumeProfileIops); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileIops)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = flex.IntValue(model.Value)
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileIopsIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfileVolumeProfileIopsFixedToMap(model *vpcv1.VolumeProfileIopsFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = flex.IntValue(model.Value)
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileIopsRangeToMap(model *vpcv1.VolumeProfileIopsRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileIopsEnumToMap(model *vpcv1.VolumeProfileIopsEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = flex.IntValue(model.Default)
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileIopsDependentRangeToMap(model *vpcv1.VolumeProfileIopsDependentRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max"] = flex.IntValue(model.Max)
	modelMap["min"] = flex.IntValue(model.Min)
	modelMap["step"] = flex.IntValue(model.Step)
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedToMap(model vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateFixed); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateFixedToMap(model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateDependent); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateDependentToMap(model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateDependent))
	} else if _, ok := model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupported); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileUnattachedCapacityUpdateSupported)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = *model.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateFixedToMap(model *vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateDependentToMap(model *vpcv1.VolumeProfileUnattachedCapacityUpdateSupportedVolumeProfileUnattachedCapacityUpdateDependent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedToMap(model vpcv1.VolumeProfileUnattachedIopsUpdateSupportedIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateFixed); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateFixedToMap(model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateFixed))
	} else if _, ok := model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateDependent); ok {
		return DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateDependentToMap(model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateDependent))
	} else if _, ok := model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupported); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.VolumeProfileUnattachedIopsUpdateSupported)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = *model.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.VolumeProfileUnattachedIopsUpdateSupportedIntf subtype encountered")
	}
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateFixedToMap(model *vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateFixed) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIBMIsVolumeProfileVolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateDependentToMap(model *vpcv1.VolumeProfileUnattachedIopsUpdateSupportedVolumeProfileUnattachedIopsUpdateDependent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	return modelMap, nil
}
