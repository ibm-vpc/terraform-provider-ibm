// Copyright IBM Corp. 2023 All Rights Reserved.
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
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func DataSourceIbmIsStorageOntapInstanceStorageVirtualMachine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineRead,

		Schema: map[string]*schema.Schema{
			"storage_ontap_instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage ontap instance identifier.",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage virtual machine identifier.",
			},
			"active_directory": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Active Directory service this storage virtual machine is joined to.If absent, this storage virtual machine is not joined to an Active Directory service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"administrators_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the domain group whose members have been granted administrative privileges for this storage virtual machine.",
						},
						"dns_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP addresses of the Active Directory DNS servers or domain controllers.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fully qualified domain name of the self-managed Active Directory.",
						},
						"domain_password_credential": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The password credential for the Active Directory domain.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this credential.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"netbios_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Active Directory computer object that will be created for the storage virtual machine.",
						},
						"organizational_unit_distinguished_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The distinguished name of the organizational unit within the self-managed Active Directory.",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username that this storage virtual machine will use when joining the Active Directory domain. This username will be the same as the username credential used in `domain_password_credential`.",
						},
					},
				},
			},
			"admin_credentials": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The credentials used for the administrator to access the storage virtual machine of thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The security certificate credential for ONTAP REST API access for the storage virtualmachine administrator.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this credential.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"password": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The password credential for the storage virtual machine administrator.If present, this password credential is used by the storage virtual machineadministrator for both ONTAP CLI SSH access and ONTAP REST API access.If absent, the storage virtual machine is not accessible through either the ONTAP CLIor ONTAP REST API using password-based authentication.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this credential.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"ssh": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The public key credential for ONTAP CLI based ssh login for the storage virtualmachine administrator.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this credential.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the storage virtual machine was created.",
			},
			"endpoints": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The data and management endpoints for this storage virtual machine.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"management": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The NetApp management endpoint for this storage virtual machine. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique IP address of an endpoint.",
									},
								},
							},
						},
						"nfs": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Network File System (NFS) protocol endpoint for this storage virtual machine.If absent, NFS is not enabled on this storage virtual machine.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique IP address of an endpoint.",
									},
								},
							},
						},
						"smb": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Server Message Block (SMB) protocol endpoint for this storage virtual machine.If absent, SMB is not enabled on this storage virtual machine.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique IP address of an endpoint.",
									},
								},
							},
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this storage virtual machine.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the storage virtual machine.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getStorageOntapInstanceStorageVirtualMachineOptions := &ontapv1.GetStorageOntapInstanceStorageVirtualMachineOptions{}

	getStorageOntapInstanceStorageVirtualMachineOptions.SetStorageOntapInstanceID(d.Get("storage_ontap_instance").(string))
	getStorageOntapInstanceStorageVirtualMachineOptions.SetID(d.Get("identifier").(string))

	storageOntapInstanceStorageVirtualMachine, response, err := ontapClient.GetStorageOntapInstanceStorageVirtualMachineWithContext(context, getStorageOntapInstanceStorageVirtualMachineOptions)
	if err != nil {
		log.Printf("[DEBUG] GetStorageOntapInstanceStorageVirtualMachineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetStorageOntapInstanceStorageVirtualMachineWithContext failed %s\n%s", err, response))
	}

	d.SetId(*storageOntapInstanceStorageVirtualMachine.ID)

	activeDirectory := []map[string]interface{}{}
	if storageOntapInstanceStorageVirtualMachine.ActiveDirectory != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineActiveDirectoryToMap(storageOntapInstanceStorageVirtualMachine.ActiveDirectory)
		if err != nil {
			return diag.FromErr(err)
		}
		activeDirectory = append(activeDirectory, modelMap)
	}
	if err = d.Set("active_directory", activeDirectory); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting active_directory %s", err))
	}

	adminCredentials := []map[string]interface{}{}
	if storageOntapInstanceStorageVirtualMachine.AdminCredentials != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineAdminCredentialsToMap(storageOntapInstanceStorageVirtualMachine.AdminCredentials)
		if err != nil {
			return diag.FromErr(err)
		}
		adminCredentials = append(adminCredentials, modelMap)
	}
	if err = d.Set("admin_credentials", adminCredentials); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting admin_credentials %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(storageOntapInstanceStorageVirtualMachine.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	endpoints := []map[string]interface{}{}
	if storageOntapInstanceStorageVirtualMachine.Endpoints != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineEndpointsToMap(storageOntapInstanceStorageVirtualMachine.Endpoints)
		if err != nil {
			return diag.FromErr(err)
		}
		endpoints = append(endpoints, modelMap)
	}
	if err = d.Set("endpoints", endpoints); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoints %s", err))
	}

	if err = d.Set("href", storageOntapInstanceStorageVirtualMachine.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("lifecycle_state", storageOntapInstanceStorageVirtualMachine.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", storageOntapInstanceStorageVirtualMachine.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("resource_type", storageOntapInstanceStorageVirtualMachine.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineActiveDirectoryToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineActiveDirectory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["administrators_group"] = model.AdministratorsGroup
	modelMap["dns_ips"] = model.DnsIps
	modelMap["domain_name"] = model.DomainName
	if model.DomainPasswordCredential != nil {
		domainPasswordCredentialMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineCredentialReferenceToMap(model.DomainPasswordCredential)
		if err != nil {
			return modelMap, err
		}
		modelMap["domain_password_credential"] = []map[string]interface{}{domainPasswordCredentialMap}
	}
	modelMap["netbios_name"] = model.NetbiosName
	modelMap["organizational_unit_distinguished_name"] = model.OrganizationalUnitDistinguishedName
	modelMap["username"] = model.Username
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineCredentialReferenceToMap(model *ontapv1.CredentialReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineAdminCredentialsToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineAdminCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Http != nil {
		httpMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineCredentialReferenceToMap(model.Http)
		if err != nil {
			return modelMap, err
		}
		modelMap["http"] = []map[string]interface{}{httpMap}
	}
	if model.Password != nil {
		passwordMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineCredentialReferenceToMap(model.Password)
		if err != nil {
			return modelMap, err
		}
		modelMap["password"] = []map[string]interface{}{passwordMap}
	}
	if model.Ssh != nil {
		sshMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineCredentialReferenceToMap(model.Ssh)
		if err != nil {
			return modelMap, err
		}
		modelMap["ssh"] = []map[string]interface{}{sshMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceStorageVirtualMachineEndpointsToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineEndpoints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	managementMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceEndpointToMap(model.Management)
	if err != nil {
		return modelMap, err
	}
	modelMap["management"] = []map[string]interface{}{managementMap}
	if model.Nfs != nil {
		nfsMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceEndpointToMap(model.Nfs)
		if err != nil {
			return modelMap, err
		}
		modelMap["nfs"] = []map[string]interface{}{nfsMap}
	}
	if model.Smb != nil {
		smbMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceEndpointToMap(model.Smb)
		if err != nil {
			return modelMap, err
		}
		modelMap["smb"] = []map[string]interface{}{smbMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachineStorageOntapInstanceEndpointToMap(model *ontapv1.StorageOntapInstanceEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["ipv4_address"] = model.Ipv4Address
	return modelMap, nil
}
