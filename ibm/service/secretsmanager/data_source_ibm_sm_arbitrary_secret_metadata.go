// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmArbitrarySecretMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmArbitrarySecretMetadataRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the secret.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A UUID identifier, or `default` secret group.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
			},
			"state_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A text representation of the secret state.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"versions_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of versions of the secret.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
		},
	}
}

func dataSourceIbmSmArbitrarySecretMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, endpointsFile, err := getSecretsManagerSession(meta.(conns.ClientSession))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d), endpointsFile)

	getSecretMetadataOptions := &secretsmanagerv2.GetSecretMetadataOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretMetadataOptions.SetID(secretId)

	arbitrarySecretMetadataIntf, response, err := secretsManagerClient.GetSecretMetadataWithContext(context, getSecretMetadataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretMetadataWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	arbitrarySecretMetadata := arbitrarySecretMetadataIntf.(*secretsmanagerv2.ArbitrarySecretMetadata)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("created_by", arbitrarySecretMetadata.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(arbitrarySecretMetadata.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", arbitrarySecretMetadata.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if arbitrarySecretMetadata.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(arbitrarySecretMetadata.CustomMetadata))
		for k, v := range arbitrarySecretMetadata.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", arbitrarySecretMetadata.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", arbitrarySecretMetadata.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(arbitrarySecretMetadata.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", arbitrarySecretMetadata.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", arbitrarySecretMetadata.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", arbitrarySecretMetadata.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(arbitrarySecretMetadata.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", arbitrarySecretMetadata.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(arbitrarySecretMetadata.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(arbitrarySecretMetadata.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(arbitrarySecretMetadata.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s_metadata", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
