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

func DataSourceIbmSmImportedCertificateMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmImportedCertificateMetadataRead,

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
			"signing_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign a certificate.",
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Common Name (AKA CN) represents the server name protected by the SSL certificate.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"intermediate_included": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate was imported with an associated intermediate certificate.",
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
			},
			"key_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm used to generate the public key that is associated with the certificate.",
			},
			"managed_csr": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The data specified to create the CSR and the private key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alt_names": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "With the Subject Alternative Name field, you can specify additional hostnames to be protected by a single SSL certificate.",
						},
						"client_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field indicates whether certificate is flagged for client use.",
						},
						"code_signing_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field indicates whether certificate is flagged for code signing use.",
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Common Name (CN) represents the server name protected by the SSL certificate.",
						},
						"csr": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate signing request.",
						},
						"country": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Country (C) values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"email_protection_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field indicates whether certificate is flagged for email protection use.",
						},
						"exclude_cn_from_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This parameter controls whether the common name is excluded from Subject Alternative Names (SANs).",
						},
						"ext_key_usage": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allowed extended key usage constraint on certificate, in a comma-delimited list.",
						},
						"ext_key_usage_oids": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A comma-delimited list of extended key usage Object Identifiers (OIDs).",
						},
						"ip_sans": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP Subject Alternative Names to define for the certificate, in a comma-delimited list.",
						},
						"key_bits": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of bits to use to generate the private key.",
						},
						"key_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of private key to generate.",
						},
						"key_usage": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allowed key usage constraint to define for certificate, in a comma-delimited list.",
						},
						"locality": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"organization": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"other_sans": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the certificate, in a comma-delimited list.",
						},
						"ou": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"policy_identifiers": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A comma-delimited list of policy Object Identifiers (OIDs).",
						},
						"postal_code": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The postal code values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"province": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"require_cn": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to false, makes the common_name field optional while generating a certificate.",
						},
						"rotate_keys": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field indicates whether the private key will be rotated.",
						},
						"server_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field indicates whether certificate is flagged for server use.",
						},
						"street_address": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The street address values to define in the subject field of the resulting certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"uri_sans": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URI Subject Alternative Names to define for the certificate, in a comma-delimited list.",
						},
						"user_ids": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the list of requested User ID (OID 0.9.2342.19200300.100.1.1) Subject values to be placed on the signed certificate.",
						},
					},
				},
			},
			"private_key_included": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate was imported with an associated private key.",
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
			},
			"validity": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The date and time that the certificate validity period begins and ends.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"not_before": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
						"not_after": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSmImportedCertificateMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, endpointsFile, err := getSecretsManagerSession(meta.(conns.ClientSession))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d), endpointsFile)

	getSecretMetadataOptions := &secretsmanagerv2.GetSecretMetadataOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretMetadataOptions.SetID(secretId)

	importedCertificateMetadataIntf, response, err := secretsManagerClient.GetSecretMetadataWithContext(context, getSecretMetadataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretMetadataWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	importedCertificateMetadata := importedCertificateMetadataIntf.(*secretsmanagerv2.ImportedCertificateMetadata)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", importedCertificateMetadata.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(importedCertificateMetadata.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", importedCertificateMetadata.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if importedCertificateMetadata.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(importedCertificateMetadata.CustomMetadata))
		for k, v := range importedCertificateMetadata.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", importedCertificateMetadata.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", importedCertificateMetadata.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(importedCertificateMetadata.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", importedCertificateMetadata.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", importedCertificateMetadata.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", importedCertificateMetadata.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(importedCertificateMetadata.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", importedCertificateMetadata.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(importedCertificateMetadata.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(importedCertificateMetadata.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("signing_algorithm", importedCertificateMetadata.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("common_name", importedCertificateMetadata.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(importedCertificateMetadata.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("intermediate_included", importedCertificateMetadata.IntermediateIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate_included"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("issuer", importedCertificateMetadata.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_algorithm", importedCertificateMetadata.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("private_key_included", importedCertificateMetadata.PrivateKeyIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key_included"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("serial_number", importedCertificateMetadata.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	validity := []map[string]interface{}{}
	if importedCertificateMetadata.Validity != nil {
		modelMap, err := dataSourceIbmSmImportedCertificateMetadataCertificateValidityToMap(importedCertificateMetadata.Validity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		validity = append(validity, modelMap)
	}
	if err = d.Set("validity", validity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), fmt.Sprintf("(Data) %s_metadata", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if importedCertificateMetadata.ManagedCsr != nil {
		managedCsrMap := managedCsrToMap(importedCertificateMetadata.ManagedCsr)
		if err = d.Set("managed_csr", []map[string]interface{}{managedCsrMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting managed_csr"), ImportedCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func dataSourceIbmSmImportedCertificateMetadataCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}
