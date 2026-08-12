// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsBareMetalServerStorageAccess() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerStorageAccessRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_bare_metal_server_storage_access", isBareMetalServerID),
				Description:  "The bare metal server identifier",
			},
			isBMSStorageAccessCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the storage access secret was created.",
			},
			isBMSStorageAccessEncryptedSecret: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The storage access secret, encrypted using public_key and returned as a base64-encoded string. This property will only be present when the status of the storage_access secret is active.",
			},
			isBMSStorageAccessPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the public SSH key used to encrypt the storage access secret.",
			},
			isBMSStorageAccessRotatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the storage access secret was last rotated.",
			},
			isBMSStorageAccessStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the storage access secret.",
			},
		},
	}
}

func DataSourceIBMIsBareMetalServerStorageAccessValidator() *validate.ResourceValidator {
	validateSchema := []validate.ValidateSchema{
		{
			Identifier:                 isBareMetalServerID,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true,
		},
	}

	validator := validate.ResourceValidator{
		ResourceName: "ibm_is_bare_metal_server_storage_access",
		Schema:       validateSchema,
	}
	return &validator
}

func dataSourceIBMISBareMetalServerStorageAccessRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_storage_access", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &vpcv1.GetBareMetalServerStorageAccessOptions{
		BareMetalServerID: &bareMetalServerID,
	}
	storageAccess, _, err := sess.GetBareMetalServerStorageAccessWithContext(context, options)
	if err != nil || storageAccess == nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBareMetalServerStorageAccessWithContext failed: %s", err.Error()), "(Data) ibm_is_bare_metal_server_storage_access", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(bareMetalServerID)
	d.Set(isBareMetalServerID, bareMetalServerID)
	d.Set(isBMSStorageAccessCreatedAt, storageAccess.CreatedAt.String())

	if storageAccess.EncryptedSecret != nil {
		d.Set(isBMSStorageAccessEncryptedSecret, *storageAccess.EncryptedSecret)
	}

	if storageAccess.PublicKey != nil && storageAccess.PublicKey.Fingerprint != nil {
		d.Set(isBMSStorageAccessPublicKey, *storageAccess.PublicKey.Fingerprint)
	}

	if storageAccess.RotatedAt != nil {
		d.Set(isBMSStorageAccessRotatedAt, storageAccess.RotatedAt.String())
	}

	d.Set(isBMSStorageAccessStatus, *storageAccess.Status)

	return nil
}
