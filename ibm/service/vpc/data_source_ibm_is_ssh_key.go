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

func DataSourceIBMISSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group ID",
			},

			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "User Tags for the ssh",
			},

			isKeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ssh key",
			},

			isKeyType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ssh key type",
			},

			isKeyFingerprint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ssh key Fingerprint",
			},

			isKeyPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH Public key data",
			},

			isKeyLength: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ssh key length",
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

			IsKeyCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			isKeyAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func dataSourceIBMISSSHKeyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get(isKeyName).(string)

	err := keyGetByName(context, d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func keyGetByName(context context.Context, d *schema.ResourceData, meta interface{}, name string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_ssh_key", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	listKeysOptions := &vpcv1.ListKeysOptions{}

	start := ""
	allrecs := []vpcv1.Key{}
	for {
		if start != "" {
			listKeysOptions.Start = &start
		}

		keys, response, err := sess.ListKeysWithContext(context, listKeysOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListKeysWithContext failed: %s\n%s", response, err.Error()), "(Data) ibm_is_ssh_key", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		start = flex.GetNext(keys.Next)
		allrecs = append(allrecs, keys.Keys...)
		if start == "" {
			break
		}
	}

	for _, key := range allrecs {
		if *key.Name == name {
			d.SetId(*key.ID)
			if err = d.Set("name", key.Name); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_ssh_key", "read", "set-name").GetDiag()
			}
			if err = d.Set("type", key.Type); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_is_ssh_key", "read", "set-type").GetDiag()
			}
			if err = d.Set("fingerprint", key.Fingerprint); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting fingerprint: %s", err), "(Data) ibm_is_ssh_key", "read", "set-fingerprint").GetDiag()
			}
			if err = d.Set("length", flex.IntValue(key.Length)); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting length: %s", err), "(Data) ibm_is_ssh_key", "read", "set-length").GetDiag()
			}
			controller, err := flex.GetBaseController(meta)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBaseController failed: %s", err.Error()), "(Data) ibm_is_ssh_key", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			d.Set(flex.ResourceControllerURL, controller+"/vpc/compute/sshKeys")
			if err = d.Set("resource_name", key.Name); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_name: %s", err), "(Data) ibm_is_ssh_key", "read", "set-resource_name").GetDiag()
			}
			if err = d.Set("resource_crn", key.CRN); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_crn: %s", err), "(Data) ibm_is_ssh_key", "read", "set-resource_crn").GetDiag()
			}
			if err = d.Set("crn", key.CRN); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_ssh_key", "read", "set-crn").GetDiag()
			}
			if key.ResourceGroup != nil {
				if err = d.Set("resource_group", *key.ResourceGroup.ID); err != nil {
					return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group: %s", err), "(Data) ibm_is_ssh_key", "read", "set-resource_group_name").GetDiag()
				}
				if err = d.Set("resource_group_name", *key.ResourceGroup.ID); err != nil {
					return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group_name: %s", err), "(Data) ibm_is_ssh_key", "read", "set-resource_group_name").GetDiag()
				}
			}
			if err = d.Set("public_key", key.PublicKey); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting public_key: %s", err), "(Data) ibm_is_ssh_key", "read", "set-public_key").GetDiag()
			}
			tags, err := flex.GetGlobalTagsUsingCRN(meta, *key.CRN, "", isUserTagType)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc ssh key (%s) tags: %s", d.Id(), err)
			}
			d.Set("tags", tags)
			accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *key.CRN, "", isKeyAccessTagType)
			if err != nil {
				log.Printf(
					"Error on get of resource SSH Key (%s) access tags: %s", d.Id(), err)
			}
			d.Set(isKeyAccessTags, accesstags)
			return nil
		}
	}
	tfErr := flex.TerraformErrorf(err, fmt.Sprintf("No SSH Key found with name %s", name), "(Data) ibm_is_ssh_key", "read")
	log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
	return tfErr.GetDiag()
}
