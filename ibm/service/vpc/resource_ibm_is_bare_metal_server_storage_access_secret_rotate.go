// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBMSStorageAccessBareMetalServer = "bare_metal_server"
	isBMSStorageAccessKey             = "key"
	isBMSStorageAccessEncryptedSecret = "encrypted_secret"
	isBMSStorageAccessCreatedAt       = "created_at"
	isBMSStorageAccessRotatedAt       = "rotated_at"
	isBMSStorageAccessStatus          = "status"
	isBMSStorageAccessPublicKey       = "public_key"
)

func ResourceIBMISBareMetalServerStorageAccessSecretRotate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBMSStorageAccessSecretRotateCreate,
		ReadContext:   resourceIBMISBMSStorageAccessSecretRotateRead,
		UpdateContext: resourceIBMISBMSStorageAccessSecretRotateUpdate,
		DeleteContext: resourceIBMISBMSStorageAccessSecretRotateDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isBMSStorageAccessBareMetalServer: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_storage_access_secret_rotate", isBMSStorageAccessBareMetalServer),
				Description:  "The ID of the bare metal server whose storage access secret will be rotated.",
			},

			isBMSStorageAccessKey: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_storage_access_secret_rotate", isBMSStorageAccessKey),
				Description:  "The ID of an SSH RSA key to use for encrypting the new storage access secret. If not specified, an existing RSA key from the server's initialization is used.",
			},

			// Computed — read back from GET /bare_metal_servers/{id}/storage_access
			isBMSStorageAccessEncryptedSecret: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The storage access secret, encrypted using public_key and returned as a base64-encoded string. This property will only be present when the status of the storage_access secret is active.",
			},
			isBMSStorageAccessCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the storage access secret was created.",
			},
			isBMSStorageAccessRotatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the storage access secret was last rotated.",
			},
			isBMSStorageAccessStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the storage access secret: active or updating.",
			},
			isBMSStorageAccessPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the public SSH key used to encrypt the storage access secret.",
			},
		},
	}
}

func ResourceIBMISBMSStorageAccessSecretRotateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBMSStorageAccessBareMetalServer,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
		},
		validate.ValidateSchema{
			Identifier:                 isBMSStorageAccessKey,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
		},
	)
	validator := validate.ResourceValidator{
		ResourceName: "ibm_is_bare_metal_server_storage_access_secret_rotate",
		Schema:       validateSchema,
	}
	return &validator
}

func bmsStorageAccessRotate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_storage_access_secret_rotate", "rotate", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bmsId := d.Get(isBMSStorageAccessBareMetalServer).(string)
	rotateOpts := &vpcv1.RotateBareMetalServerStorageAccessSecretOptions{
		BareMetalServerID: &bmsId,
	}
	if keyId, ok := d.GetOk(isBMSStorageAccessKey); ok {
		keyIdStr := keyId.(string)
		rotateOpts.Key = &vpcv1.KeyIdentityByID{
			ID: &keyIdStr,
		}
	}

	_, err = sess.RotateBareMetalServerStorageAccessSecretWithContext(context, rotateOpts)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("RotateBareMetalServerStorageAccessSecretWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_storage_access_secret_rotate", "rotate")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	return nil
}

func resourceIBMISBMSStorageAccessSecretRotateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bmsId := d.Get(isBMSStorageAccessBareMetalServer).(string)
	if diagErr := bmsStorageAccessRotate(context, d, meta); diagErr != nil {
		return diagErr
	}
	d.SetId(bmsId)
	return resourceIBMISBMSStorageAccessSecretRotateRead(context, d, meta)
}

func resourceIBMISBMSStorageAccessSecretRotateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_storage_access_secret_rotate", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bmsId := d.Id()
	getOpts := &vpcv1.GetBareMetalServerStorageAccessOptions{
		BareMetalServerID: &bmsId,
	}
	storageAccess, response, err := sess.GetBareMetalServerStorageAccessWithContext(context, getOpts)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBareMetalServerStorageAccessWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_storage_access_secret_rotate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set(isBMSStorageAccessBareMetalServer, bmsId)
	d.Set(isBMSStorageAccessCreatedAt, storageAccess.CreatedAt.String())
	if storageAccess.EncryptedSecret != nil {
		d.Set(isBMSStorageAccessEncryptedSecret, *storageAccess.EncryptedSecret)
	}
	if storageAccess.PublicKey != nil {
		d.Set(isBMSStorageAccessPublicKey, *storageAccess.PublicKey.Fingerprint)
	}
	if storageAccess.RotatedAt != nil {
		d.Set(isBMSStorageAccessRotatedAt, storageAccess.RotatedAt.String())
	}
	d.Set(isBMSStorageAccessStatus, *storageAccess.Status)

	return nil
}

// Update is called when the key changes — re-triggers the rotation.
func resourceIBMISBMSStorageAccessSecretRotateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange(isBMSStorageAccessKey) {
		if diagErr := bmsStorageAccessRotate(context, d, meta); diagErr != nil {
			return diagErr
		}
	}
	return resourceIBMISBMSStorageAccessSecretRotateRead(context, d, meta)
}

func resourceIBMISBMSStorageAccessSecretRotateDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
