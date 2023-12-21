// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func DataSourceIbmIsStorageOntapInstanceStorageVirtualMachines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesRead,

		Schema: map[string]*schema.Schema{
			"storage_ontap_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage ontap instance identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
			},
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"storage_virtual_machines": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of storage virtual machines.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this storage virtual machine.",
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
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listStorageOntapInstanceStorageVirtualMachinesOptions := &ontapv1.ListStorageOntapInstanceStorageVirtualMachinesOptions{}

	listStorageOntapInstanceStorageVirtualMachinesOptions.SetStorageOntapInstanceID(d.Get("storage_ontap_instance_id").(string))
	if _, ok := d.GetOk("name"); ok {
		listStorageOntapInstanceStorageVirtualMachinesOptions.SetName(d.Get("name").(string))
	}

	var pager *ontapv1.StorageOntapInstanceStorageVirtualMachinesPager
	pager, err = ontapClient.NewStorageOntapInstanceStorageVirtualMachinesPager(listStorageOntapInstanceStorageVirtualMachinesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] StorageOntapInstanceStorageVirtualMachinesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("StorageOntapInstanceStorageVirtualMachinesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("storage_virtual_machines", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_virtual_machines %s", err))
	}

	return nil
}

// dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesID returns a reasonable ID for the list.
func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineCollectionFirstToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineCollectionNextToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachine) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ActiveDirectory != nil {
		activeDirectoryMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineActiveDirectoryToMap(model.ActiveDirectory)
		if err != nil {
			return modelMap, err
		}
		modelMap["active_directory"] = []map[string]interface{}{activeDirectoryMap}
	}
	if model.AdminCredentials != nil {
		adminCredentialsMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineAdminCredentialsToMap(model.AdminCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["admin_credentials"] = []map[string]interface{}{adminCredentialsMap}
	}
	modelMap["created_at"] = model.CreatedAt.String()
	endpointsMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineEndpointsToMap(model.Endpoints)
	if err != nil {
		return modelMap, err
	}
	modelMap["endpoints"] = []map[string]interface{}{endpointsMap}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineActiveDirectoryToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineActiveDirectory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["administrators_group"] = model.AdministratorsGroup
	modelMap["dns_ips"] = model.DnsIps
	modelMap["domain_name"] = model.DomainName
	if model.DomainPasswordCredential != nil {
		domainPasswordCredentialMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesCredentialReferenceToMap(model.DomainPasswordCredential)
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

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesCredentialReferenceToMap(model *ontapv1.CredentialReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineAdminCredentialsToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineAdminCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Http != nil {
		httpMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesCredentialReferenceToMap(model.Http)
		if err != nil {
			return modelMap, err
		}
		modelMap["http"] = []map[string]interface{}{httpMap}
	}
	if model.Password != nil {
		passwordMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesCredentialReferenceToMap(model.Password)
		if err != nil {
			return modelMap, err
		}
		modelMap["password"] = []map[string]interface{}{passwordMap}
	}
	if model.Ssh != nil {
		sshMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesCredentialReferenceToMap(model.Ssh)
		if err != nil {
			return modelMap, err
		}
		modelMap["ssh"] = []map[string]interface{}{sshMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceStorageVirtualMachineEndpointsToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineEndpoints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	managementMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceEndpointToMap(model.Management)
	if err != nil {
		return modelMap, err
	}
	modelMap["management"] = []map[string]interface{}{managementMap}
	if model.Nfs != nil {
		nfsMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceEndpointToMap(model.Nfs)
		if err != nil {
			return modelMap, err
		}
		modelMap["nfs"] = []map[string]interface{}{nfsMap}
	}
	if model.Smb != nil {
		smbMap, err := dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceEndpointToMap(model.Smb)
		if err != nil {
			return modelMap, err
		}
		modelMap["smb"] = []map[string]interface{}{smbMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageVirtualMachinesStorageOntapInstanceEndpointToMap(model *ontapv1.StorageOntapInstanceEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["ipv4_address"] = model.Ipv4Address
	return modelMap, nil
}
