// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSnapshots                              = "snapshots"
	isSnapshotId                             = "id"
	isSnapshotCopiesId                       = "snapshot_copies_id"
	isSnapshotCopiesName                     = "snapshot_copies_name"
	isSnapshotCopiesCRN                      = "snapshot_copies_crn"
	isSnapshotCopiesRemoteRegionName         = "snapshot_copies_remote_region_name"
	isSnapshotSourceSnapshotId               = "source_snapshot_id"
	isSnapshotSourceSnapshotRemoteRegionName = "source_snapshot_remote_region_name"
	isSnapshotSourceVolumeRemoteRegionName   = "snapshot_source_volume_remote_region_name"
	isSnapshotConsistencyGroupId             = "snapshot_consistency_group_id"
	isSnapshotConsistencyGroupCrn            = "snapshot_consistency_group_crn"
	isSnapshotConsistencyGroup               = "snapshot_consistency_group"
)

func DataSourceSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISSnapshotsRead,

		Schema: map[string]*schema.Schema{

			isSnapshotResourceGroup: {
				Type:        schema.TypeString,
				Description: "Filters the snapshot collection by resources group id",
				Optional:    true,
			},

			isSnapshotName: {
				Type:        schema.TypeString,
				Description: "Filters the snapshot collection by snapshot name",
				Optional:    true,
			},

			isSnapshotSourceImage: {
				Type:        schema.TypeString,
				Description: "Filters the snapshot collection by source image id",
				Optional:    true,
			},

			isSnapshotSourceVolume: {
				Type:        schema.TypeString,
				Description: "Filters the snapshot collection by source volume id",
				Optional:    true,
			},

			"backup_policy_plan_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to backup policy jobs with the backup plan with the specified identifier",
			},

			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with the exact tag value",
			},

			isSnapshotCopiesId: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with copies with the specified identifier.",
				Optional:    true,
			},

			isSnapshotCopiesName: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with copies with the exact specified name.",
				Optional:    true,
			},

			isSnapshotCopiesCRN: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with copies with the specified CRN.",
				Optional:    true,
			},

			isSnapshotCopiesRemoteRegionName: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with copies with the exact remote region name.",
				Optional:    true,
			},

			isSnapshotSourceSnapshotId: {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with the source snapshot with the specified identifier.",
				Optional:    true,
			},

			isSnapshotSourceSnapshotRemoteRegionName: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with a source snapshot with the exact remote region name.",
				Optional:    true,
			},

			isSnapshotSourceVolumeRemoteRegionName: {
				Type:        schema.TypeString,
				Description: "Filters the collection to snapshots with a source snapshot with the exact remote region name.",
				Optional:    true,
			},

			isSnapshotConsistencyGroupId: {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with a source snapshot with the exact snapshot consistency group id.",
				Optional:    true,
			},

			isSnapshotConsistencyGroupCrn: {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with a source snapshot with the exact snapshot consistency group crn.",
				Optional:    true,
			},

			isSnapshots: {
				Type:        schema.TypeList,
				Description: "List of snapshots",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSnapshotId: {
							Type:     schema.TypeString,
							Computed: true,
						},

						"service_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The [service tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags) prefixed with `is.snapshot:` associated with this snapshot.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						isSnapshotCopies: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The copies of this snapshot in other regions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for the copied snapshot.",
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
										Computed:    true,
										Description: "The URL for the copied snapshot.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the copied snapshot.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for the copied snapshot. The name is unique across all snapshots in the copied snapshot's native region.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},

						isSnapshotConsistencyGroup: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The snapshot consistency group which created this snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of this snapshot consistency group.",
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
										Computed:    true,
										Description: "The URL for the snapshot consistency group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the snapshot consistency group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for the snapshot consistency group. The name is unique across all snapshot consistency groups in the region.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},

						isSnapshotName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot name",
						},

						isSnapshotResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group info",
						},

						isSnapshotSourceVolume: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot source volume",
						},

						isSnapshotSourceSnapshot: {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "If present, the source snapshot this snapshot was created from.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The CRN of the source snapshot.",
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
										Computed:    true,
										Description: "The URL for the source snapshot.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the source snapshot.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for the source snapshot. The name is unique across all snapshots in the source snapshot's native region.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},

						isSnapshotSourceImage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If present, the image id from which the data on this volume was most directly provisioned.",
						},

						isSnapshotOperatingSystem: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for the operating system included in this image",
						},

						isSnapshotLCState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot lifecycle state",
						},
						isSnapshotCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the resource",
						},
						isSnapshotEncryption: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encryption type of the snapshot",
						},
						isSnapshotEncryptionKey: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reference to the root key used to wrap the data encryption key for the source volume.",
						},
						isSnapshotHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for the snapshot",
						},

						isSnapshotBootable: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if a boot volume attachment can be created with a volume created from this snapshot",
						},

						isSnapshotMinCapacity: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum capacity of the snapshot",
						},
						isSnapshotResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type of the snapshot",
						},

						isSnapshotSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the snapshot",
						},

						isSnapshotClones: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Zones for creating the snapshot clone",
						},

						isSnapshotCapturedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this snapshot was created",
						},

						isSnapshotUserTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "User Tags for the snapshot",
						},

						isSnapshotAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},
						isSnapshotCatalogOffering: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The catalog offering inherited from the snapshot's source. If a virtual server instance is provisioned with a source_snapshot specifying this snapshot, the virtual server instance will use this snapshot's catalog offering, including its pricing plan.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isSnapshotCatalogOfferingPlanCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this catalog offering version's billing plan",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									isSnapshotCatalogOfferingVersionCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this version of a catalog offering",
									},
								},
							},
						},

						isSnapshotBackupPolicyPlan: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, the backup policy plan which created this snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
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
										Computed:    true,
										Description: "The URL for this backup policy plan.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this backup policy plan.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this backup policy plan.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced",
									},
								},
							},
						},
						"allowed_use": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The usage constraints to match against the requested instance or bare metal server properties to determine compatibility. Can only be specified for bootable snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bare_metal_server": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expression that must be satisfied by the properties of a bare metal server provisioned using the image data in this snapshot.",
									},
									"instance": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expression that must be satisfied by the properties of a virtual server instance provisioned using this snapshot.",
									},
									"api_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API version with which to evaluate the expressions.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISSnapshotsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := getSnapshots(context, d, meta)
	if err != nil {
		return err
	}
	return nil
}

func getSnapshots(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_snapshots", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	start := ""
	allrecs := []vpcv1.Snapshot{}
	for {
		listSnapshotOptions := &vpcv1.ListSnapshotsOptions{}
		if start != "" {
			listSnapshotOptions.Start = &start
		}
		if rgFilterOk, ok := d.GetOk(isSnapshotResourceGroup); ok {
			rgFilter := rgFilterOk.(string)
			listSnapshotOptions.ResourceGroupID = &rgFilter
		}
		if nameFilterOk, ok := d.GetOk(isSnapshotName); ok {
			nameFilter := nameFilterOk.(string)
			listSnapshotOptions.Name = &nameFilter
		}
		if sourceImageFilterOk, ok := d.GetOk(isSnapshotSourceImage); ok {
			sourceImageFilter := sourceImageFilterOk.(string)
			listSnapshotOptions.SourceImageID = &sourceImageFilter
		}
		if sourceVolumeFilterOk, ok := d.GetOk(isSnapshotSourceVolume); ok {
			sourceVolumeFilter := sourceVolumeFilterOk.(string)
			listSnapshotOptions.SourceVolumeID = &sourceVolumeFilter
		}
		if backupPolicyPlanIdFilterOk, ok := d.GetOk("backup_policy_plan_id"); ok {
			backupPolicyPlanIdFilter := backupPolicyPlanIdFilterOk.(string)
			listSnapshotOptions.BackupPolicyPlanID = &backupPolicyPlanIdFilter
		}
		if tagFilterOk, ok := d.GetOk("tag"); ok {
			tagFilter := tagFilterOk.(string)
			listSnapshotOptions.Tag = &tagFilter
		}
		if copiesId, ok := d.GetOk(isSnapshotCopiesId); ok {
			copiesIdFilter := copiesId.(string)
			listSnapshotOptions.CopiesID = &copiesIdFilter
		}
		if copiesName, ok := d.GetOk(isSnapshotCopiesName); ok {
			copiesNameFilter := copiesName.(string)
			listSnapshotOptions.CopiesName = &copiesNameFilter
		}
		if copiesCRN, ok := d.GetOk(isSnapshotCopiesCRN); ok {
			copiesCRNFilter := copiesCRN.(string)
			listSnapshotOptions.CopiesCRN = &copiesCRNFilter
		}
		if copiesRemoteRegionName, ok := d.GetOk(isSnapshotCopiesRemoteRegionName); ok {
			copiesRemoteRegionNameFilter := copiesRemoteRegionName.(string)
			listSnapshotOptions.CopiesRemoteRegionName = &copiesRemoteRegionNameFilter
		}
		if sourceSnapshotId, ok := d.GetOk(isSnapshotSourceSnapshotId); ok {
			sourceSnapshotIdFilter := sourceSnapshotId.(string)
			listSnapshotOptions.SourceSnapshotID = &sourceSnapshotIdFilter
		}
		if sourceSnapshotRemoteRegionName, ok := d.GetOk(isSnapshotSourceSnapshotRemoteRegionName); ok {
			sourceSnapshotRemoteRegionNameFilter := sourceSnapshotRemoteRegionName.(string)
			listSnapshotOptions.SourceSnapshotRemoteRegionName = &sourceSnapshotRemoteRegionNameFilter
		}
		if sourceVolumeRemoteRegionName, ok := d.GetOk(isSnapshotSourceVolumeRemoteRegionName); ok {
			sourceVolumeRemoteRegionNameFilter := sourceVolumeRemoteRegionName.(string)
			listSnapshotOptions.SourceVolumeRemoteRegionName = &sourceVolumeRemoteRegionNameFilter
		}
		if snapshotConsistencyGroupId, ok := d.GetOk(isSnapshotConsistencyGroupId); ok {
			snapshotConsistencyGroupIdFilter := snapshotConsistencyGroupId.(string)
			listSnapshotOptions.SnapshotConsistencyGroupID = &snapshotConsistencyGroupIdFilter
		}
		if snapshotConsistencyGroupCrn, ok := d.GetOk(isSnapshotConsistencyGroupCrn); ok {
			snapshotConsistencyGroupCrnFilter := snapshotConsistencyGroupCrn.(string)
			listSnapshotOptions.SnapshotConsistencyGroupCRN = &snapshotConsistencyGroupCrnFilter
		}
		snapshots, _, err := sess.ListSnapshotsWithContext(context, listSnapshotOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListSnapshotsWithContext failed %s", err), "(Data) ibm_is_snapshots", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		start = flex.GetNext(snapshots.Next)
		allrecs = append(allrecs, snapshots.Snapshots...)
		if start == "" {
			break
		}
	}

	snapshotsInfo := make([]map[string]interface{}, 0)
	for _, snapshot := range allrecs {
		l := map[string]interface{}{
			isSnapshotId:           *snapshot.ID,
			isSnapshotName:         *snapshot.Name,
			isSnapshotHref:         *snapshot.Href,
			isSnapshotCRN:          *snapshot.CRN,
			isSnapshotMinCapacity:  *snapshot.MinimumCapacity,
			isSnapshotSize:         *snapshot.Size,
			isSnapshotEncryption:   *snapshot.Encryption,
			isSnapshotLCState:      *snapshot.LifecycleState,
			isSnapshotResourceType: *snapshot.ResourceType,
			isSnapshotBootable:     *snapshot.Bootable,
		}
		if snapshot.EncryptionKey != nil {
			l[isSnapshotEncryptionKey] = snapshot.EncryptionKey.CRN
		}
		if snapshot.ServiceTags != nil && len(snapshot.ServiceTags) > 0 {
			l["service_tags"] = snapshot.ServiceTags
		}
		if snapshot.EncryptionKey != nil && snapshot.EncryptionKey.CRN != nil {
			l[isSnapshotEncryptionKey] = *snapshot.EncryptionKey.CRN
		}
		if snapshot.CapturedAt != nil {
			l[isSnapshotCapturedAt] = (*snapshot.CapturedAt).String()
		}
		// source snapshot
		sourceSnapshotList := []map[string]interface{}{}
		if snapshot.SourceSnapshot != nil {
			sourceSnapshot := map[string]interface{}{}
			sourceSnapshot["href"] = snapshot.SourceSnapshot.Href
			sourceSnapshot["crn"] = snapshot.SourceSnapshot.CRN
			if snapshot.SourceSnapshot.Deleted != nil {
				snapshotSourceSnapshotDeletedMap := map[string]interface{}{}
				snapshotSourceSnapshotDeletedMap["more_info"] = snapshot.SourceSnapshot.Deleted.MoreInfo
				sourceSnapshot["deleted"] = []map[string]interface{}{snapshotSourceSnapshotDeletedMap}
			}
			sourceSnapshot["id"] = snapshot.SourceSnapshot.ID
			sourceSnapshot["name"] = snapshot.SourceSnapshot.Name
			sourceSnapshot["resource_type"] = snapshot.SourceSnapshot.ResourceType
			sourceSnapshotList = append(sourceSnapshotList, sourceSnapshot)
		}
		l[isSnapshotSourceSnapshot] = sourceSnapshotList

		// snapshot copies
		snapshotCopies := []map[string]interface{}{}
		if snapshot.Copies != nil {
			for _, copiesItem := range snapshot.Copies {
				copiesMap, err := dataSourceIBMIsSnapshotsSnapshotCopiesItemToMap(&copiesItem)
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_snapshots", "read", "copies-to-map").GetDiag()
				}
				snapshotCopies = append(snapshotCopies, copiesMap)
			}
		}
		l[isSnapshotCopies] = snapshotCopies

		// snapshot consistency group
		snapshotConsistencyGroupList := []map[string]interface{}{}
		if snapshot.SnapshotConsistencyGroup != nil {
			snapshotConsistencyGroup := map[string]interface{}{}
			snapshotConsistencyGroup["href"] = snapshot.SnapshotConsistencyGroup.Href
			snapshotConsistencyGroup["crn"] = snapshot.SnapshotConsistencyGroup.CRN
			if snapshot.SnapshotConsistencyGroup.Deleted != nil {
				snapshotConsistencyGroupDeletedMap := map[string]interface{}{}
				snapshotConsistencyGroupDeletedMap["more_info"] = snapshot.SnapshotConsistencyGroup.Deleted.MoreInfo
				snapshotConsistencyGroup["deleted"] = []map[string]interface{}{snapshotConsistencyGroupDeletedMap}
			}
			snapshotConsistencyGroup["id"] = snapshot.SnapshotConsistencyGroup.ID
			snapshotConsistencyGroup["name"] = snapshot.SnapshotConsistencyGroup.Name
			snapshotConsistencyGroup["resource_type"] = snapshot.SnapshotConsistencyGroup.ResourceType
			snapshotConsistencyGroupList = append(snapshotConsistencyGroupList, snapshotConsistencyGroup)
		}
		l[isSnapshotConsistencyGroup] = snapshotConsistencyGroupList

		if snapshot.UserTags != nil {
			l[isSnapshotUserTags] = snapshot.UserTags
		}
		if snapshot.ResourceGroup != nil && snapshot.ResourceGroup.ID != nil {
			l[isSnapshotResourceGroup] = *snapshot.ResourceGroup.ID
		}
		if snapshot.SourceVolume != nil && snapshot.SourceVolume.ID != nil {
			l[isSnapshotSourceVolume] = *snapshot.SourceVolume.ID
		}
		if snapshot.SourceImage != nil && snapshot.SourceImage.ID != nil {
			l[isSnapshotSourceImage] = *snapshot.SourceImage.ID
		}
		if snapshot.OperatingSystem != nil && snapshot.OperatingSystem.Name != nil {
			l[isSnapshotOperatingSystem] = *snapshot.OperatingSystem.Name
		}
		var clones []string
		clones = make([]string, 0)
		if snapshot.Clones != nil {
			for _, clone := range snapshot.Clones {
				if clone.Zone != nil {
					clones = append(clones, *clone.Zone.Name)
				}
			}
		}
		l[isSnapshotClones] = flex.NewStringSet(schema.HashString, clones)

		// catalog
		catalogList := make([]map[string]interface{}, 0)
		if snapshot.CatalogOffering != nil {
			versionCrn := ""
			if snapshot.CatalogOffering.Version != nil && snapshot.CatalogOffering.Version.CRN != nil {
				versionCrn = *snapshot.CatalogOffering.Version.CRN
			}
			catalogMap := map[string]interface{}{}
			if versionCrn != "" {
				catalogMap[isSnapshotCatalogOfferingVersionCrn] = versionCrn
			}
			if snapshot.CatalogOffering.Plan != nil {
				planCrn := ""
				if snapshot.CatalogOffering.Plan.CRN != nil {
					planCrn = *snapshot.CatalogOffering.Plan.CRN
				}
				if planCrn != "" {
					catalogMap[isSnapshotCatalogOfferingPlanCrn] = planCrn
				}
				if snapshot.CatalogOffering.Plan.Deleted != nil {
					deletedMap := resourceIbmIsSnapshotCatalogOfferingVersionPlanReferenceDeletedToMap(*snapshot.CatalogOffering.Plan.Deleted)
					catalogMap["deleted"] = []map[string]interface{}{deletedMap}
				}
			}
			catalogList = append(catalogList, catalogMap)
		}
		l[isSnapshotCatalogOffering] = catalogList

		backupPolicyPlanList := []map[string]interface{}{}
		if snapshot.BackupPolicyPlan != nil {
			backupPolicyPlan := map[string]interface{}{}
			if snapshot.BackupPolicyPlan.Deleted != nil {
				snapshotBackupPolicyPlanDeletedMap := map[string]interface{}{}
				snapshotBackupPolicyPlanDeletedMap["more_info"] = snapshot.BackupPolicyPlan.Deleted.MoreInfo
				backupPolicyPlan["deleted"] = []map[string]interface{}{snapshotBackupPolicyPlanDeletedMap}
			}
			backupPolicyPlan["href"] = snapshot.BackupPolicyPlan.Href
			backupPolicyPlan["id"] = snapshot.BackupPolicyPlan.ID
			backupPolicyPlan["name"] = snapshot.BackupPolicyPlan.Name
			backupPolicyPlan["resource_type"] = snapshot.BackupPolicyPlan.ResourceType
			backupPolicyPlanList = append(backupPolicyPlanList, backupPolicyPlan)
		}
		l[isSnapshotBackupPolicyPlan] = backupPolicyPlanList
		if snapshot.AllowedUse != nil {
			allowedUseList := []map[string]interface{}{}
			modelMap, err := DataSourceIBMIsSnapshotAllowedUseToMap(snapshot.AllowedUse)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_snapshots", "read")
				log.Println(tfErr.GetDiag())
			}
			allowedUseList = append(allowedUseList, modelMap)
			l["allowed_use"] = allowedUseList
		}
		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshot.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on get of resource snapshot (%s) access tags: %s", d.Id(), err)
		}
		l[isSnapshotAccessTags] = accesstags
		snapshotsInfo = append(snapshotsInfo, l)
	}
	d.SetId(dataSourceIBMISSnapshotsID(d))
	if err = d.Set("snapshots", snapshotsInfo); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting snapshots %s", err), "(Data) ibm_is_snapshots", "read", "snapshots-set").GetDiag()
	}
	return nil
}

// dataSourceIBMISSnapshotsID returns a reasonable ID for the snapshot list.
func dataSourceIBMISSnapshotsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsSnapshotsSnapshotCopiesItemToMap(model *vpcv1.SnapshotCopiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsSnapshotsSnapshotRemoteReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}

	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIsSnapshotsSnapshotRemoteReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
