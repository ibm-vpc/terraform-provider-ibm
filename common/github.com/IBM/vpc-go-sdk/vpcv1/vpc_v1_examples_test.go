//go:build examples

/**
 * (C) Copyright IBM Corp. 2026.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vpcv1_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This file provides an example of how to use the vpc service.
//
// The following configuration properties are assumed to be defined:
// VPC_URL=<service base url>
// VPC_AUTH_TYPE=iam
// VPC_APIKEY=<IAM apikey>
// VPC_AUTH_URL=<IAM token service base URL - omit this if using the production environment>
//
// These configuration properties can be exported as environment variables, or stored
// in a configuration file and then:
// export IBM_CREDENTIALS_FILE=<name of configuration file>
var _ = Describe(`VpcV1 Examples Tests`, func() {

	const externalConfigFile = "../vpc_v1.env"

	var (
		vpcService *vpcv1.VpcV1
		config     map[string]string
	)

	var shouldSkipTest = func() {
		Skip("External configuration is not available, skipping examples...")
	}

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			var err error
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping examples: " + err.Error())
			}

			os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(vpcv1.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping examples: " + err.Error())
			} else if len(config) == 0 {
				Skip("Unable to load service properties, skipping examples")
			}

			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			var err error

			// begin-common

			vpcServiceOptions := &vpcv1.VpcV1Options{
				Version: core.StringPtr("2024-06-23"),
			}

			vpcService, err = vpcv1.NewVpcV1UsingExternalConfig(vpcServiceOptions)

			if err != nil {
				panic(err)
			}

			// end-common

			Expect(vpcService).ToNot(BeNil())
		})
	})

	Describe(`VpcV1 request examples`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBackupPolicies request example`, func() {
			fmt.Println("\nListBackupPolicies() result:")
			// begin-list_backup_policies
			listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Tag:             core.StringPtr("testString"),
			}

			pager, err := vpcService.NewBackupPoliciesPager(listBackupPoliciesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BackupPolicyIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_backup_policies
		})
		It(`CreateBackupPolicy request example`, func() {
			fmt.Println("\nCreateBackupPolicy() result:")
			// begin-create_backup_policy

			backupPolicyPrototypeModel := &vpcv1.BackupPolicyPrototypeBackupPolicyMatchResourceTypeInstancePrototype{
				MatchUserTags:     []string{"my-daily-backup-policy"},
				MatchResourceType: core.StringPtr("instance"),
			}

			createBackupPolicyOptions := vpcService.NewCreateBackupPolicyOptions(
				backupPolicyPrototypeModel,
			)

			backupPolicy, response, err := vpcService.CreateBackupPolicy(createBackupPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_backup_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(backupPolicy).ToNot(BeNil())
		})
		It(`ListBackupPolicyJobs request example`, func() {
			fmt.Println("\nListBackupPolicyJobs() result:")
			// begin-list_backup_policy_jobs
			listBackupPolicyJobsOptions := &vpcv1.ListBackupPolicyJobsOptions{
				BackupPolicyID:     core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Sort:               core.StringPtr("name"),
				SourceID:           core.StringPtr("testString"),
				Status:             core.StringPtr("failed"),
				TargetSnapshotsCRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot:r006-f6bfa329-0e36-433f-a3bb-0df632e79263"),
				TargetSnapshotsID:  core.StringPtr("testString"),
			}

			pager, err := vpcService.NewBackupPolicyJobsPager(listBackupPolicyJobsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BackupPolicyJob
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_backup_policy_jobs
		})
		It(`GetBackupPolicyJob request example`, func() {
			fmt.Println("\nGetBackupPolicyJob() result:")
			// begin-get_backup_policy_job

			getBackupPolicyJobOptions := vpcService.NewGetBackupPolicyJobOptions(
				"testString",
				"testString",
			)

			backupPolicyJob, response, err := vpcService.GetBackupPolicyJob(getBackupPolicyJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyJob, "", "  ")
			fmt.Println(string(b))

			// end-get_backup_policy_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyJob).ToNot(BeNil())
		})
		It(`ListBackupPolicyPlans request example`, func() {
			fmt.Println("\nListBackupPolicyPlans() result:")
			// begin-list_backup_policy_plans

			listBackupPolicyPlansOptions := vpcService.NewListBackupPolicyPlansOptions(
				"testString",
			)
			listBackupPolicyPlansOptions.SetName("my-name")

			backupPolicyPlanCollection, response, err := vpcService.ListBackupPolicyPlans(listBackupPolicyPlansOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyPlanCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_backup_policy_plans

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlanCollection).ToNot(BeNil())
		})
		It(`CreateBackupPolicyPlan request example`, func() {
			fmt.Println("\nCreateBackupPolicyPlan() result:")
			// begin-create_backup_policy_plan

			createBackupPolicyPlanOptions := vpcService.NewCreateBackupPolicyPlanOptions(
				"testString",
				"30 */2 * * 1-5",
			)

			backupPolicyPlan, response, err := vpcService.CreateBackupPolicyPlan(createBackupPolicyPlanOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyPlan, "", "  ")
			fmt.Println(string(b))

			// end-create_backup_policy_plan

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
		It(`GetBackupPolicyPlan request example`, func() {
			fmt.Println("\nGetBackupPolicyPlan() result:")
			// begin-get_backup_policy_plan

			getBackupPolicyPlanOptions := vpcService.NewGetBackupPolicyPlanOptions(
				"testString",
				"testString",
			)

			backupPolicyPlan, response, err := vpcService.GetBackupPolicyPlan(getBackupPolicyPlanOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyPlan, "", "  ")
			fmt.Println(string(b))

			// end-get_backup_policy_plan

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
		It(`UpdateBackupPolicyPlan request example`, func() {
			fmt.Println("\nUpdateBackupPolicyPlan() result:")
			// begin-update_backup_policy_plan

			backupPolicyPlanPatchModel := &vpcv1.BackupPolicyPlanPatch{
				CopyUserTags: core.BoolPtr(false),
			}
			backupPolicyPlanPatchModelAsPatch, asPatchErr := backupPolicyPlanPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBackupPolicyPlanOptions := vpcService.NewUpdateBackupPolicyPlanOptions(
				"testString",
				"testString",
				backupPolicyPlanPatchModelAsPatch,
			)
			updateBackupPolicyPlanOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			backupPolicyPlan, response, err := vpcService.UpdateBackupPolicyPlan(updateBackupPolicyPlanOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyPlan, "", "  ")
			fmt.Println(string(b))

			// end-update_backup_policy_plan

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
		It(`GetBackupPolicy request example`, func() {
			fmt.Println("\nGetBackupPolicy() result:")
			// begin-get_backup_policy

			getBackupPolicyOptions := vpcService.NewGetBackupPolicyOptions(
				"testString",
			)

			backupPolicy, response, err := vpcService.GetBackupPolicy(getBackupPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_backup_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicy).ToNot(BeNil())
		})
		It(`UpdateBackupPolicy request example`, func() {
			fmt.Println("\nUpdateBackupPolicy() result:")
			// begin-update_backup_policy

			backupPolicyPatchModel := &vpcv1.BackupPolicyPatch{
				MatchUserTags: []string{"my-tag-1", "my-tag-2", "my-tag-3"},
			}
			backupPolicyPatchModelAsPatch, asPatchErr := backupPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBackupPolicyOptions := vpcService.NewUpdateBackupPolicyOptions(
				"testString",
				backupPolicyPatchModelAsPatch,
			)
			updateBackupPolicyOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			backupPolicy, response, err := vpcService.UpdateBackupPolicy(updateBackupPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_backup_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicy).ToNot(BeNil())
		})
		It(`ListBareMetalServerCapacities request example`, func() {
			fmt.Println("\nListBareMetalServerCapacities() result:")
			// begin-list_bare_metal_server_capacities
			listBareMetalServerCapacitiesOptions := &vpcv1.ListBareMetalServerCapacitiesOptions{
				Limit:       core.Int64Ptr(int64(10)),
				ProfileName: core.StringPtr("bx2d-metal-192x768"),
				ZoneName:    core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewBareMetalServerCapacitiesPager(listBareMetalServerCapacitiesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServerCapacity
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_bare_metal_server_capacities
		})
		It(`ListBareMetalServerProfiles request example`, func() {
			fmt.Println("\nListBareMetalServerProfiles() result:")
			// begin-list_bare_metal_server_profiles
			listBareMetalServerProfilesOptions := &vpcv1.ListBareMetalServerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewBareMetalServerProfilesPager(listBareMetalServerProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServerProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_bare_metal_server_profiles
		})
		It(`GetBareMetalServerProfile request example`, func() {
			fmt.Println("\nGetBareMetalServerProfile() result:")
			// begin-get_bare_metal_server_profile

			getBareMetalServerProfileOptions := vpcService.NewGetBareMetalServerProfileOptions(
				"bx2-metal-192x768",
			)

			bareMetalServerProfile, response, err := vpcService.GetBareMetalServerProfile(getBareMetalServerProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerProfile).ToNot(BeNil())
		})
		It(`ListBareMetalServers request example`, func() {
			fmt.Println("\nListBareMetalServers() result:")
			// begin-list_bare_metal_servers
			listBareMetalServersOptions := &vpcv1.ListBareMetalServersOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ReservationCRN:  core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
				ReservationID:   core.StringPtr("testString"),
				ReservationName: core.StringPtr("my-reservation"),
				ResourceGroupID: core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewBareMetalServersPager(listBareMetalServersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServer
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_bare_metal_servers
		})
		It(`CreateBareMetalServer request example`, func() {
			fmt.Println("\nCreateBareMetalServer() result:")
			// begin-create_bare_metal_server

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			bareMetalServerInitializationByCatalogOfferingCatalogOfferingModel := &vpcv1.BareMetalServerInitializationByCatalogOfferingCatalogOfferingBareMetalServerCatalogOfferingPrototypeCatalogOfferingByOffering{
				Offering: catalogOfferingIdentityModel,
			}

			bareMetalServerInitializationPrototypeModel := &vpcv1.BareMetalServerInitializationPrototypeBareMetalServerInitializationByCatalogOffering{
				Keys:            []vpcv1.KeyIdentityIntf{keyIdentityModel},
				CatalogOffering: bareMetalServerInitializationByCatalogOfferingCatalogOfferingModel,
			}

			bareMetalServerPrototypeProfileModel := &vpcv1.BareMetalServerPrototypeProfileBareMetalServerProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/bare_metal_server/profiles/bx2-metal-192x768"),
			}

			bareMetalServerPrototypeZoneModel := &vpcv1.BareMetalServerPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			bareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.BareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentBareMetalServerPrimaryNetworkAttachmentPrototypeBareMetalServerPrimaryNetworkAttachmentByPciPrototype{
				VirtualNetworkInterface: bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			bareMetalServerPrototypeModel := &vpcv1.BareMetalServerPrototypeBareMetalServerByNetworkAttachment{
				Initialization:           bareMetalServerInitializationPrototypeModel,
				Profile:                  bareMetalServerPrototypeProfileModel,
				Zone:                     bareMetalServerPrototypeZoneModel,
				PrimaryNetworkAttachment: bareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createBareMetalServerOptions := vpcService.NewCreateBareMetalServerOptions(
				bareMetalServerPrototypeModel,
			)

			bareMetalServer, response, err := vpcService.CreateBareMetalServer(createBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServer, "", "  ")
			fmt.Println(string(b))

			// end-create_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServer).ToNot(BeNil())
		})
		It(`CreateBareMetalServerConsoleAccessToken request example`, func() {
			fmt.Println("\nCreateBareMetalServerConsoleAccessToken() result:")
			// begin-create_bare_metal_server_console_access_token

			createBareMetalServerConsoleAccessTokenOptions := vpcService.NewCreateBareMetalServerConsoleAccessTokenOptions(
				"testString",
				"serial",
			)

			bareMetalServerConsoleAccessToken, response, err := vpcService.CreateBareMetalServerConsoleAccessToken(createBareMetalServerConsoleAccessTokenOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerConsoleAccessToken, "", "  ")
			fmt.Println(string(b))

			// end-create_bare_metal_server_console_access_token

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerConsoleAccessToken).ToNot(BeNil())
		})
		It(`ListBareMetalServerDisks request example`, func() {
			fmt.Println("\nListBareMetalServerDisks() result:")
			// begin-list_bare_metal_server_disks

			listBareMetalServerDisksOptions := vpcService.NewListBareMetalServerDisksOptions(
				"testString",
			)

			bareMetalServerDiskCollection, response, err := vpcService.ListBareMetalServerDisks(listBareMetalServerDisksOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerDiskCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_bare_metal_server_disks

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDiskCollection).ToNot(BeNil())
		})
		It(`GetBareMetalServerDisk request example`, func() {
			fmt.Println("\nGetBareMetalServerDisk() result:")
			// begin-get_bare_metal_server_disk

			getBareMetalServerDiskOptions := vpcService.NewGetBareMetalServerDiskOptions(
				"testString",
				"testString",
			)

			bareMetalServerDisk, response, err := vpcService.GetBareMetalServerDisk(getBareMetalServerDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerDisk, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDisk).ToNot(BeNil())
		})
		It(`UpdateBareMetalServerDisk request example`, func() {
			fmt.Println("\nUpdateBareMetalServerDisk() result:")
			// begin-update_bare_metal_server_disk

			bareMetalServerDiskPatchModel := &vpcv1.BareMetalServerDiskPatch{}
			bareMetalServerDiskPatchModelAsPatch, asPatchErr := bareMetalServerDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerDiskOptions := vpcService.NewUpdateBareMetalServerDiskOptions(
				"testString",
				"testString",
				bareMetalServerDiskPatchModelAsPatch,
			)

			bareMetalServerDisk, response, err := vpcService.UpdateBareMetalServerDisk(updateBareMetalServerDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerDisk, "", "  ")
			fmt.Println(string(b))

			// end-update_bare_metal_server_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDisk).ToNot(BeNil())
		})
		It(`ListBareMetalServerNetworkAttachments request example`, func() {
			fmt.Println("\nListBareMetalServerNetworkAttachments() result:")
			// begin-list_bare_metal_server_network_attachments
			listBareMetalServerNetworkAttachmentsOptions := &vpcv1.ListBareMetalServerNetworkAttachmentsOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewBareMetalServerNetworkAttachmentsPager(listBareMetalServerNetworkAttachmentsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServerNetworkAttachmentIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_bare_metal_server_network_attachments
		})
		It(`CreateBareMetalServerNetworkAttachment request example`, func() {
			fmt.Println("\nCreateBareMetalServerNetworkAttachment() result:")
			// begin-create_bare_metal_server_network_attachment

			bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			bareMetalServerNetworkAttachmentPrototypeModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{
				VirtualNetworkInterface: bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
				InterfaceType:           core.StringPtr("pci"),
			}

			createBareMetalServerNetworkAttachmentOptions := vpcService.NewCreateBareMetalServerNetworkAttachmentOptions(
				"testString",
				bareMetalServerNetworkAttachmentPrototypeModel,
			)

			bareMetalServerNetworkAttachment, response, err := vpcService.CreateBareMetalServerNetworkAttachment(createBareMetalServerNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-create_bare_metal_server_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
		It(`GetBareMetalServerNetworkAttachment request example`, func() {
			fmt.Println("\nGetBareMetalServerNetworkAttachment() result:")
			// begin-get_bare_metal_server_network_attachment

			getBareMetalServerNetworkAttachmentOptions := vpcService.NewGetBareMetalServerNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			bareMetalServerNetworkAttachment, response, err := vpcService.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
		It(`UpdateBareMetalServerNetworkAttachment request example`, func() {
			fmt.Println("\nUpdateBareMetalServerNetworkAttachment() result:")
			// begin-update_bare_metal_server_network_attachment

			bareMetalServerNetworkAttachmentPatchModel := &vpcv1.BareMetalServerNetworkAttachmentPatch{
				Name: core.StringPtr("my-bare-metal-server-network-attachment-updated"),
			}
			bareMetalServerNetworkAttachmentPatchModelAsPatch, asPatchErr := bareMetalServerNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerNetworkAttachmentOptions := vpcService.NewUpdateBareMetalServerNetworkAttachmentOptions(
				"testString",
				"testString",
				bareMetalServerNetworkAttachmentPatchModelAsPatch,
			)

			bareMetalServerNetworkAttachment, response, err := vpcService.UpdateBareMetalServerNetworkAttachment(updateBareMetalServerNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_bare_metal_server_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
		It(`ListBareMetalServerNetworkInterfaces request example`, func() {
			fmt.Println("\nListBareMetalServerNetworkInterfaces() result:")
			// begin-list_bare_metal_server_network_interfaces
			listBareMetalServerNetworkInterfacesOptions := &vpcv1.ListBareMetalServerNetworkInterfacesOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewBareMetalServerNetworkInterfacesPager(listBareMetalServerNetworkInterfacesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServerNetworkInterfaceIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_bare_metal_server_network_interfaces
		})
		It(`CreateBareMetalServerNetworkInterface request example`, func() {
			fmt.Println("\nCreateBareMetalServerNetworkInterface() result:")
			// begin-create_bare_metal_server_network_interface

			bareMetalServerNetworkInterfacePrototypeSubnetModel := &vpcv1.BareMetalServerNetworkInterfacePrototypeSubnetSubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			bareMetalServerNetworkInterfacePrototypeModel := &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByHiperSocketPrototype{
				Subnet:        bareMetalServerNetworkInterfacePrototypeSubnetModel,
				InterfaceType: core.StringPtr("hipersocket"),
			}

			createBareMetalServerNetworkInterfaceOptions := vpcService.NewCreateBareMetalServerNetworkInterfaceOptions(
				"testString",
				bareMetalServerNetworkInterfacePrototypeModel,
			)

			bareMetalServerNetworkInterface, response, err := vpcService.CreateBareMetalServerNetworkInterface(createBareMetalServerNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-create_bare_metal_server_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
		It(`GetBareMetalServerNetworkInterface request example`, func() {
			fmt.Println("\nGetBareMetalServerNetworkInterface() result:")
			// begin-get_bare_metal_server_network_interface

			getBareMetalServerNetworkInterfaceOptions := vpcService.NewGetBareMetalServerNetworkInterfaceOptions(
				"testString",
				"testString",
			)

			bareMetalServerNetworkInterface, response, err := vpcService.GetBareMetalServerNetworkInterface(getBareMetalServerNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
		It(`UpdateBareMetalServerNetworkInterface request example`, func() {
			fmt.Println("\nUpdateBareMetalServerNetworkInterface() result:")
			// begin-update_bare_metal_server_network_interface

			bareMetalServerNetworkInterfacePatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{}
			bareMetalServerNetworkInterfacePatchModelAsPatch, asPatchErr := bareMetalServerNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerNetworkInterfaceOptions := vpcService.NewUpdateBareMetalServerNetworkInterfaceOptions(
				"testString",
				"testString",
				bareMetalServerNetworkInterfacePatchModelAsPatch,
			)

			bareMetalServerNetworkInterface, response, err := vpcService.UpdateBareMetalServerNetworkInterface(updateBareMetalServerNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-update_bare_metal_server_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
		It(`ListBareMetalServerNetworkInterfaceFloatingIps request example`, func() {
			fmt.Println("\nListBareMetalServerNetworkInterfaceFloatingIps() result:")
			// begin-list_bare_metal_server_network_interface_floating_ips

			listBareMetalServerNetworkInterfaceFloatingIpsOptions := vpcService.NewListBareMetalServerNetworkInterfaceFloatingIpsOptions(
				"testString",
				"testString",
			)

			floatingIPUnpaginatedCollection, response, err := vpcService.ListBareMetalServerNetworkInterfaceFloatingIps(listBareMetalServerNetworkInterfaceFloatingIpsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIPUnpaginatedCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_bare_metal_server_network_interface_floating_ips

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPUnpaginatedCollection).ToNot(BeNil())
		})
		It(`GetBareMetalServerNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nGetBareMetalServerNetworkInterfaceFloatingIP() result:")
			// begin-get_bare_metal_server_network_interface_floating_ip

			getBareMetalServerNetworkInterfaceFloatingIPOptions := vpcService.NewGetBareMetalServerNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			floatingIP, response, err := vpcService.GetBareMetalServerNetworkInterfaceFloatingIP(getBareMetalServerNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`AddBareMetalServerNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nAddBareMetalServerNetworkInterfaceFloatingIP() result:")
			// begin-add_bare_metal_server_network_interface_floating_ip

			addBareMetalServerNetworkInterfaceFloatingIPOptions := vpcService.NewAddBareMetalServerNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			floatingIP, response, err := vpcService.AddBareMetalServerNetworkInterfaceFloatingIP(addBareMetalServerNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-add_bare_metal_server_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`ListBareMetalServerNetworkInterfaceIps request example`, func() {
			fmt.Println("\nListBareMetalServerNetworkInterfaceIps() result:")
			// begin-list_bare_metal_server_network_interface_ips

			listBareMetalServerNetworkInterfaceIpsOptions := vpcService.NewListBareMetalServerNetworkInterfaceIpsOptions(
				"testString",
				"testString",
			)

			reservedIPCollectionBareMetalServerNetworkInterfaceContext, response, err := vpcService.ListBareMetalServerNetworkInterfaceIps(listBareMetalServerNetworkInterfaceIpsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIPCollectionBareMetalServerNetworkInterfaceContext, "", "  ")
			fmt.Println(string(b))

			// end-list_bare_metal_server_network_interface_ips

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPCollectionBareMetalServerNetworkInterfaceContext).ToNot(BeNil())
		})
		It(`GetBareMetalServerNetworkInterfaceIP request example`, func() {
			fmt.Println("\nGetBareMetalServerNetworkInterfaceIP() result:")
			// begin-get_bare_metal_server_network_interface_ip

			getBareMetalServerNetworkInterfaceIPOptions := vpcService.NewGetBareMetalServerNetworkInterfaceIPOptions(
				"testString",
				"testString",
				"testString",
			)

			reservedIP, response, err := vpcService.GetBareMetalServerNetworkInterfaceIP(getBareMetalServerNetworkInterfaceIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_network_interface_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`GetBareMetalServer request example`, func() {
			fmt.Println("\nGetBareMetalServer() result:")
			// begin-get_bare_metal_server

			getBareMetalServerOptions := vpcService.NewGetBareMetalServerOptions(
				"testString",
			)

			bareMetalServer, response, err := vpcService.GetBareMetalServer(getBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServer, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServer).ToNot(BeNil())
		})
		It(`UpdateBareMetalServer request example`, func() {
			fmt.Println("\nUpdateBareMetalServer() result:")
			// begin-update_bare_metal_server

			bareMetalServerPatchModel := &vpcv1.BareMetalServerPatch{}
			bareMetalServerPatchModelAsPatch, asPatchErr := bareMetalServerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerOptions := vpcService.NewUpdateBareMetalServerOptions(
				"testString",
				bareMetalServerPatchModelAsPatch,
			)

			bareMetalServer, response, err := vpcService.UpdateBareMetalServer(updateBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServer, "", "  ")
			fmt.Println(string(b))

			// end-update_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServer).ToNot(BeNil())
		})
		It(`UpdateFirmwareForBareMetalServer request example`, func() {
			// begin-update_firmware_for_bare_metal_server

			updateFirmwareForBareMetalServerOptions := vpcService.NewUpdateFirmwareForBareMetalServerOptions(
				"testString",
			)

			response, err := vpcService.UpdateFirmwareForBareMetalServer(updateFirmwareForBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from UpdateFirmwareForBareMetalServer(): %d\n", response.StatusCode)
			}

			// end-update_firmware_for_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`GetBareMetalServerInitialization request example`, func() {
			fmt.Println("\nGetBareMetalServerInitialization() result:")
			// begin-get_bare_metal_server_initialization

			getBareMetalServerInitializationOptions := vpcService.NewGetBareMetalServerInitializationOptions(
				"testString",
			)

			bareMetalServerInitialization, response, err := vpcService.GetBareMetalServerInitialization(getBareMetalServerInitializationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerInitialization, "", "  ")
			fmt.Println(string(b))

			// end-get_bare_metal_server_initialization

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerInitialization).ToNot(BeNil())
		})
		It(`ReplaceBareMetalServerInitialization request example`, func() {
			fmt.Println("\nReplaceBareMetalServerInitialization() result:")
			// begin-replace_bare_metal_server_initialization

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			bareMetalServerInitializationByImageImageModel := &vpcv1.BareMetalServerInitializationByImageImageImageIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"),
			}

			bareMetalServerInitializationReplacePrototypeModel := &vpcv1.BareMetalServerInitializationReplacePrototypeBareMetalServerInitializationByImage{
				Keys:  []vpcv1.KeyIdentityIntf{keyIdentityModel},
				Image: bareMetalServerInitializationByImageImageModel,
			}

			replaceBareMetalServerInitializationOptions := vpcService.NewReplaceBareMetalServerInitializationOptions(
				"testString",
				bareMetalServerInitializationReplacePrototypeModel,
			)

			bareMetalServerInitialization, response, err := vpcService.ReplaceBareMetalServerInitialization(replaceBareMetalServerInitializationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(bareMetalServerInitialization, "", "  ")
			fmt.Println(string(b))

			// end-replace_bare_metal_server_initialization

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerInitialization).ToNot(BeNil())
		})
		It(`RestartBareMetalServer request example`, func() {
			// begin-restart_bare_metal_server

			restartBareMetalServerOptions := vpcService.NewRestartBareMetalServerOptions(
				"testString",
			)

			response, err := vpcService.RestartBareMetalServer(restartBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RestartBareMetalServer(): %d\n", response.StatusCode)
			}

			// end-restart_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`StartBareMetalServer request example`, func() {
			// begin-start_bare_metal_server

			startBareMetalServerOptions := vpcService.NewStartBareMetalServerOptions(
				"testString",
			)

			response, err := vpcService.StartBareMetalServer(startBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from StartBareMetalServer(): %d\n", response.StatusCode)
			}

			// end-start_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`StopBareMetalServer request example`, func() {
			// begin-stop_bare_metal_server

			stopBareMetalServerOptions := vpcService.NewStopBareMetalServerOptions(
				"testString",
				"hard",
			)

			response, err := vpcService.StopBareMetalServer(stopBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from StopBareMetalServer(): %d\n", response.StatusCode)
			}

			// end-stop_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`ListClusterNetworkProfiles request example`, func() {
			fmt.Println("\nListClusterNetworkProfiles() result:")
			// begin-list_cluster_network_profiles
			listClusterNetworkProfilesOptions := &vpcv1.ListClusterNetworkProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewClusterNetworkProfilesPager(listClusterNetworkProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ClusterNetworkProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_cluster_network_profiles
		})
		It(`GetClusterNetworkProfile request example`, func() {
			fmt.Println("\nGetClusterNetworkProfile() result:")
			// begin-get_cluster_network_profile

			getClusterNetworkProfileOptions := vpcService.NewGetClusterNetworkProfileOptions(
				"h100",
			)

			clusterNetworkProfile, response, err := vpcService.GetClusterNetworkProfile(getClusterNetworkProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_cluster_network_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkProfile).ToNot(BeNil())
		})
		It(`ListClusterNetworks request example`, func() {
			fmt.Println("\nListClusterNetworks() result:")
			// begin-list_cluster_networks
			listClusterNetworksOptions := &vpcv1.ListClusterNetworksOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewClusterNetworksPager(listClusterNetworksOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ClusterNetwork
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_cluster_networks
		})
		It(`CreateClusterNetwork request example`, func() {
			fmt.Println("\nCreateClusterNetwork() result:")
			// begin-create_cluster_network

			clusterNetworkPrototypeProfileModel := &vpcv1.ClusterNetworkPrototypeProfileClusterNetworkProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"),
			}

			clusterNetworkPrototypeVPCModel := &vpcv1.ClusterNetworkPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			clusterNetworkPrototypeZoneModel := &vpcv1.ClusterNetworkPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createClusterNetworkOptions := vpcService.NewCreateClusterNetworkOptions(
				clusterNetworkPrototypeProfileModel,
				clusterNetworkPrototypeVPCModel,
				clusterNetworkPrototypeZoneModel,
			)

			clusterNetwork, response, err := vpcService.CreateClusterNetwork(createClusterNetworkOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetwork, "", "  ")
			fmt.Println(string(b))

			// end-create_cluster_network

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetwork).ToNot(BeNil())
		})
		It(`ListClusterNetworkInterfaces request example`, func() {
			fmt.Println("\nListClusterNetworkInterfaces() result:")
			// begin-list_cluster_network_interfaces
			listClusterNetworkInterfacesOptions := &vpcv1.ListClusterNetworkInterfacesOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
			}

			pager, err := vpcService.NewClusterNetworkInterfacesPager(listClusterNetworkInterfacesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ClusterNetworkInterface
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_cluster_network_interfaces
		})
		It(`CreateClusterNetworkInterface request example`, func() {
			fmt.Println("\nCreateClusterNetworkInterface() result:")
			// begin-create_cluster_network_interface

			createClusterNetworkInterfaceOptions := vpcService.NewCreateClusterNetworkInterfaceOptions(
				"testString",
			)

			clusterNetworkInterface, response, err := vpcService.CreateClusterNetworkInterface(createClusterNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-create_cluster_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
		It(`GetClusterNetworkInterface request example`, func() {
			fmt.Println("\nGetClusterNetworkInterface() result:")
			// begin-get_cluster_network_interface

			getClusterNetworkInterfaceOptions := vpcService.NewGetClusterNetworkInterfaceOptions(
				"testString",
				"testString",
			)

			clusterNetworkInterface, response, err := vpcService.GetClusterNetworkInterface(getClusterNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-get_cluster_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
		It(`UpdateClusterNetworkInterface request example`, func() {
			fmt.Println("\nUpdateClusterNetworkInterface() result:")
			// begin-update_cluster_network_interface

			clusterNetworkInterfacePatchModel := &vpcv1.ClusterNetworkInterfacePatch{}
			clusterNetworkInterfacePatchModelAsPatch, asPatchErr := clusterNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkInterfaceOptions := vpcService.NewUpdateClusterNetworkInterfaceOptions(
				"testString",
				"testString",
				clusterNetworkInterfacePatchModelAsPatch,
			)
			updateClusterNetworkInterfaceOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkInterface, response, err := vpcService.UpdateClusterNetworkInterface(updateClusterNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-update_cluster_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
		It(`ListClusterNetworkSubnets request example`, func() {
			fmt.Println("\nListClusterNetworkSubnets() result:")
			// begin-list_cluster_network_subnets
			listClusterNetworkSubnetsOptions := &vpcv1.ListClusterNetworkSubnetsOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
			}

			pager, err := vpcService.NewClusterNetworkSubnetsPager(listClusterNetworkSubnetsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ClusterNetworkSubnet
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_cluster_network_subnets
		})
		It(`CreateClusterNetworkSubnet request example`, func() {
			fmt.Println("\nCreateClusterNetworkSubnet() result:")
			// begin-create_cluster_network_subnet

			clusterNetworkSubnetPrototypeModel := &vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByIPv4CIDRBlockPrototype{
				Ipv4CIDRBlock: core.StringPtr("10.0.0.0/24"),
			}

			createClusterNetworkSubnetOptions := vpcService.NewCreateClusterNetworkSubnetOptions(
				"testString",
				clusterNetworkSubnetPrototypeModel,
			)

			clusterNetworkSubnet, response, err := vpcService.CreateClusterNetworkSubnet(createClusterNetworkSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnet, "", "  ")
			fmt.Println(string(b))

			// end-create_cluster_network_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
		It(`ListClusterNetworkSubnetReservedIps request example`, func() {
			fmt.Println("\nListClusterNetworkSubnetReservedIps() result:")
			// begin-list_cluster_network_subnet_reserved_ips
			listClusterNetworkSubnetReservedIpsOptions := &vpcv1.ListClusterNetworkSubnetReservedIpsOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
				Name:                   core.StringPtr("my-name"),
				Sort:                   core.StringPtr("name"),
			}

			pager, err := vpcService.NewClusterNetworkSubnetReservedIpsPager(listClusterNetworkSubnetReservedIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ClusterNetworkSubnetReservedIP
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_cluster_network_subnet_reserved_ips
		})
		It(`CreateClusterNetworkSubnetReservedIP request example`, func() {
			fmt.Println("\nCreateClusterNetworkSubnetReservedIP() result:")
			// begin-create_cluster_network_subnet_reserved_ip

			createClusterNetworkSubnetReservedIPOptions := vpcService.NewCreateClusterNetworkSubnetReservedIPOptions(
				"testString",
				"testString",
			)

			clusterNetworkSubnetReservedIP, response, err := vpcService.CreateClusterNetworkSubnetReservedIP(createClusterNetworkSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnetReservedIP, "", "  ")
			fmt.Println(string(b))

			// end-create_cluster_network_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
		It(`GetClusterNetworkSubnetReservedIP request example`, func() {
			fmt.Println("\nGetClusterNetworkSubnetReservedIP() result:")
			// begin-get_cluster_network_subnet_reserved_ip

			getClusterNetworkSubnetReservedIPOptions := vpcService.NewGetClusterNetworkSubnetReservedIPOptions(
				"testString",
				"testString",
				"testString",
			)

			clusterNetworkSubnetReservedIP, response, err := vpcService.GetClusterNetworkSubnetReservedIP(getClusterNetworkSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnetReservedIP, "", "  ")
			fmt.Println(string(b))

			// end-get_cluster_network_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
		It(`UpdateClusterNetworkSubnetReservedIP request example`, func() {
			fmt.Println("\nUpdateClusterNetworkSubnetReservedIP() result:")
			// begin-update_cluster_network_subnet_reserved_ip

			clusterNetworkSubnetReservedIPPatchModel := &vpcv1.ClusterNetworkSubnetReservedIPPatch{}
			clusterNetworkSubnetReservedIPPatchModelAsPatch, asPatchErr := clusterNetworkSubnetReservedIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkSubnetReservedIPOptions := vpcService.NewUpdateClusterNetworkSubnetReservedIPOptions(
				"testString",
				"testString",
				"testString",
				clusterNetworkSubnetReservedIPPatchModelAsPatch,
			)
			updateClusterNetworkSubnetReservedIPOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkSubnetReservedIP, response, err := vpcService.UpdateClusterNetworkSubnetReservedIP(updateClusterNetworkSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnetReservedIP, "", "  ")
			fmt.Println(string(b))

			// end-update_cluster_network_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
		It(`GetClusterNetworkSubnet request example`, func() {
			fmt.Println("\nGetClusterNetworkSubnet() result:")
			// begin-get_cluster_network_subnet

			getClusterNetworkSubnetOptions := vpcService.NewGetClusterNetworkSubnetOptions(
				"testString",
				"testString",
			)

			clusterNetworkSubnet, response, err := vpcService.GetClusterNetworkSubnet(getClusterNetworkSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnet, "", "  ")
			fmt.Println(string(b))

			// end-get_cluster_network_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
		It(`UpdateClusterNetworkSubnet request example`, func() {
			fmt.Println("\nUpdateClusterNetworkSubnet() result:")
			// begin-update_cluster_network_subnet

			clusterNetworkSubnetPatchModel := &vpcv1.ClusterNetworkSubnetPatch{}
			clusterNetworkSubnetPatchModelAsPatch, asPatchErr := clusterNetworkSubnetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkSubnetOptions := vpcService.NewUpdateClusterNetworkSubnetOptions(
				"testString",
				"testString",
				clusterNetworkSubnetPatchModelAsPatch,
			)
			updateClusterNetworkSubnetOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkSubnet, response, err := vpcService.UpdateClusterNetworkSubnet(updateClusterNetworkSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnet, "", "  ")
			fmt.Println(string(b))

			// end-update_cluster_network_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
		It(`GetClusterNetwork request example`, func() {
			fmt.Println("\nGetClusterNetwork() result:")
			// begin-get_cluster_network

			getClusterNetworkOptions := vpcService.NewGetClusterNetworkOptions(
				"testString",
			)

			clusterNetwork, response, err := vpcService.GetClusterNetwork(getClusterNetworkOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetwork, "", "  ")
			fmt.Println(string(b))

			// end-get_cluster_network

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetwork).ToNot(BeNil())
		})
		It(`UpdateClusterNetwork request example`, func() {
			fmt.Println("\nUpdateClusterNetwork() result:")
			// begin-update_cluster_network

			clusterNetworkPatchModel := &vpcv1.ClusterNetworkPatch{}
			clusterNetworkPatchModelAsPatch, asPatchErr := clusterNetworkPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkOptions := vpcService.NewUpdateClusterNetworkOptions(
				"testString",
				clusterNetworkPatchModelAsPatch,
			)
			updateClusterNetworkOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetwork, response, err := vpcService.UpdateClusterNetwork(updateClusterNetworkOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetwork, "", "  ")
			fmt.Println(string(b))

			// end-update_cluster_network

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetwork).ToNot(BeNil())
		})
		It(`ListDedicatedHostGroups request example`, func() {
			fmt.Println("\nListDedicatedHostGroups() result:")
			// begin-list_dedicated_host_groups
			listDedicatedHostGroupsOptions := &vpcv1.ListDedicatedHostGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				ZoneName:        core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewDedicatedHostGroupsPager(listDedicatedHostGroupsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.DedicatedHostGroup
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_dedicated_host_groups
		})
		It(`CreateDedicatedHostGroup request example`, func() {
			fmt.Println("\nCreateDedicatedHostGroup() result:")
			// begin-create_dedicated_host_group

			dedicatedHostGroupPrototypeZoneModel := &vpcv1.DedicatedHostGroupPrototypeZoneZoneIdentityByName{
				Name: core.StringPtr("us-south-1"),
			}

			createDedicatedHostGroupOptions := vpcService.NewCreateDedicatedHostGroupOptions(
				"mx2",
				"balanced",
				dedicatedHostGroupPrototypeZoneModel,
			)

			dedicatedHostGroup, response, err := vpcService.CreateDedicatedHostGroup(createDedicatedHostGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostGroup, "", "  ")
			fmt.Println(string(b))

			// end-create_dedicated_host_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
		It(`GetDedicatedHostGroup request example`, func() {
			fmt.Println("\nGetDedicatedHostGroup() result:")
			// begin-get_dedicated_host_group

			getDedicatedHostGroupOptions := vpcService.NewGetDedicatedHostGroupOptions(
				"testString",
			)

			dedicatedHostGroup, response, err := vpcService.GetDedicatedHostGroup(getDedicatedHostGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_dedicated_host_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
		It(`UpdateDedicatedHostGroup request example`, func() {
			fmt.Println("\nUpdateDedicatedHostGroup() result:")
			// begin-update_dedicated_host_group

			dedicatedHostGroupPatchModel := &vpcv1.DedicatedHostGroupPatch{
				Name: core.StringPtr("my-dedicated-host-group-updated"),
			}
			dedicatedHostGroupPatchModelAsPatch, asPatchErr := dedicatedHostGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostGroupOptions := vpcService.NewUpdateDedicatedHostGroupOptions(
				"testString",
				dedicatedHostGroupPatchModelAsPatch,
			)

			dedicatedHostGroup, response, err := vpcService.UpdateDedicatedHostGroup(updateDedicatedHostGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostGroup, "", "  ")
			fmt.Println(string(b))

			// end-update_dedicated_host_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
		It(`ListDedicatedHostProfiles request example`, func() {
			fmt.Println("\nListDedicatedHostProfiles() result:")
			// begin-list_dedicated_host_profiles
			listDedicatedHostProfilesOptions := &vpcv1.ListDedicatedHostProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewDedicatedHostProfilesPager(listDedicatedHostProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.DedicatedHostProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_dedicated_host_profiles
		})
		It(`GetDedicatedHostProfile request example`, func() {
			fmt.Println("\nGetDedicatedHostProfile() result:")
			// begin-get_dedicated_host_profile

			getDedicatedHostProfileOptions := vpcService.NewGetDedicatedHostProfileOptions(
				"mx2-host-152x1216",
			)

			dedicatedHostProfile, response, err := vpcService.GetDedicatedHostProfile(getDedicatedHostProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_dedicated_host_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostProfile).ToNot(BeNil())
		})
		It(`ListDedicatedHosts request example`, func() {
			fmt.Println("\nListDedicatedHosts() result:")
			// begin-list_dedicated_hosts
			listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{
				DedicatedHostGroupID: core.StringPtr("testString"),
				Limit:                core.Int64Ptr(int64(10)),
				Name:                 core.StringPtr("my-name"),
				ResourceGroupID:      core.StringPtr("testString"),
				ZoneName:             core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewDedicatedHostsPager(listDedicatedHostsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.DedicatedHost
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_dedicated_hosts
		})
		It(`CreateDedicatedHost request example`, func() {
			fmt.Println("\nCreateDedicatedHost() result:")
			// begin-create_dedicated_host

			dedicatedHostPrototypeProfileModel := &vpcv1.DedicatedHostPrototypeProfileDedicatedHostProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/dedicated_host/profiles/mx2-host-152x1216"),
			}

			dedicatedHostByGroupGroupModel := &vpcv1.DedicatedHostByGroupGroupDedicatedHostGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host-group:0717-bcc5b834-1258-4b9c-c3b4-43bc9cf5cde0"),
			}

			dedicatedHostPrototypeModel := &vpcv1.DedicatedHostPrototypeDedicatedHostByGroup{
				Profile: dedicatedHostPrototypeProfileModel,
				Group:   dedicatedHostByGroupGroupModel,
			}

			createDedicatedHostOptions := vpcService.NewCreateDedicatedHostOptions(
				dedicatedHostPrototypeModel,
			)

			dedicatedHost, response, err := vpcService.CreateDedicatedHost(createDedicatedHostOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHost, "", "  ")
			fmt.Println(string(b))

			// end-create_dedicated_host

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(dedicatedHost).ToNot(BeNil())
		})
		It(`ListDedicatedHostDisks request example`, func() {
			fmt.Println("\nListDedicatedHostDisks() result:")
			// begin-list_dedicated_host_disks

			listDedicatedHostDisksOptions := vpcService.NewListDedicatedHostDisksOptions(
				"testString",
			)

			dedicatedHostDiskCollection, response, err := vpcService.ListDedicatedHostDisks(listDedicatedHostDisksOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostDiskCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_dedicated_host_disks

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDiskCollection).ToNot(BeNil())
		})
		It(`GetDedicatedHostDisk request example`, func() {
			fmt.Println("\nGetDedicatedHostDisk() result:")
			// begin-get_dedicated_host_disk

			getDedicatedHostDiskOptions := vpcService.NewGetDedicatedHostDiskOptions(
				"testString",
				"testString",
			)

			dedicatedHostDisk, response, err := vpcService.GetDedicatedHostDisk(getDedicatedHostDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostDisk, "", "  ")
			fmt.Println(string(b))

			// end-get_dedicated_host_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDisk).ToNot(BeNil())
		})
		It(`UpdateDedicatedHostDisk request example`, func() {
			fmt.Println("\nUpdateDedicatedHostDisk() result:")
			// begin-update_dedicated_host_disk

			dedicatedHostDiskPatchModel := &vpcv1.DedicatedHostDiskPatch{
				Name: core.StringPtr("my-dedicated-host-disk-updated"),
			}
			dedicatedHostDiskPatchModelAsPatch, asPatchErr := dedicatedHostDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostDiskOptions := vpcService.NewUpdateDedicatedHostDiskOptions(
				"testString",
				"testString",
				dedicatedHostDiskPatchModelAsPatch,
			)

			dedicatedHostDisk, response, err := vpcService.UpdateDedicatedHostDisk(updateDedicatedHostDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHostDisk, "", "  ")
			fmt.Println(string(b))

			// end-update_dedicated_host_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDisk).ToNot(BeNil())
		})
		It(`GetDedicatedHost request example`, func() {
			fmt.Println("\nGetDedicatedHost() result:")
			// begin-get_dedicated_host

			getDedicatedHostOptions := vpcService.NewGetDedicatedHostOptions(
				"testString",
			)

			dedicatedHost, response, err := vpcService.GetDedicatedHost(getDedicatedHostOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHost, "", "  ")
			fmt.Println(string(b))

			// end-get_dedicated_host

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHost).ToNot(BeNil())
		})
		It(`UpdateDedicatedHost request example`, func() {
			fmt.Println("\nUpdateDedicatedHost() result:")
			// begin-update_dedicated_host

			dedicatedHostPatchModel := &vpcv1.DedicatedHostPatch{
				InstancePlacementEnabled: core.BoolPtr(false),
			}
			dedicatedHostPatchModelAsPatch, asPatchErr := dedicatedHostPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostOptions := vpcService.NewUpdateDedicatedHostOptions(
				"testString",
				dedicatedHostPatchModelAsPatch,
			)

			dedicatedHost, response, err := vpcService.UpdateDedicatedHost(updateDedicatedHostOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(dedicatedHost, "", "  ")
			fmt.Println(string(b))

			// end-update_dedicated_host

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHost).ToNot(BeNil())
		})
		It(`ListEndpointGateways request example`, func() {
			fmt.Println("\nListEndpointGateways() result:")
			// begin-list_endpoint_gateways
			listEndpointGatewaysOptions := &vpcv1.ListEndpointGatewaysOptions{
				DnsResolutionBindingMode: []string{"disabled"},
				LifecycleState:           []string{"stable"},
				Limit:                    core.Int64Ptr(int64(10)),
				Name:                     core.StringPtr("my-name"),
				ResourceGroupID:          core.StringPtr("testString"),
				VPCCRN:                   core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:                    core.StringPtr("testString"),
				VPCName:                  core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewEndpointGatewaysPager(listEndpointGatewaysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.EndpointGateway
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_endpoint_gateways
		})
		It(`CreateEndpointGateway request example`, func() {
			fmt.Println("\nCreateEndpointGateway() result:")
			// begin-create_endpoint_gateway

			endpointGatewayTargetPrototypeModel := &vpcv1.EndpointGatewayTargetPrototypeEndpointGatewayTargetResourceTypePrivatePathServiceGatewayPrototype{
				CRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::private-path-service-gateway:r006-7268d425-59b7-48fd-9735-81a7271657d5"),
				ResourceType: core.StringPtr("private_path_service_gateway"),
			}

			endpointGatewayPrototypeVPCModel := &vpcv1.EndpointGatewayPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			createEndpointGatewayOptions := vpcService.NewCreateEndpointGatewayOptions(
				endpointGatewayTargetPrototypeModel,
				endpointGatewayPrototypeVPCModel,
			)

			endpointGateway, response, err := vpcService.CreateEndpointGateway(createEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGateway, "", "  ")
			fmt.Println(string(b))

			// end-create_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(endpointGateway).ToNot(BeNil())
		})
		It(`ListEndpointGatewayIps request example`, func() {
			fmt.Println("\nListEndpointGatewayIps() result:")
			// begin-list_endpoint_gateway_ips
			listEndpointGatewayIpsOptions := &vpcv1.ListEndpointGatewayIpsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Sort:              core.StringPtr("name"),
			}

			pager, err := vpcService.NewEndpointGatewayIpsPager(listEndpointGatewayIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_endpoint_gateway_ips
		})
		It(`GetEndpointGatewayIP request example`, func() {
			fmt.Println("\nGetEndpointGatewayIP() result:")
			// begin-get_endpoint_gateway_ip

			getEndpointGatewayIPOptions := vpcService.NewGetEndpointGatewayIPOptions(
				"testString",
				"testString",
			)

			reservedIP, response, err := vpcService.GetEndpointGatewayIP(getEndpointGatewayIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-get_endpoint_gateway_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`AddEndpointGatewayIP request example`, func() {
			fmt.Println("\nAddEndpointGatewayIP() result:")
			// begin-add_endpoint_gateway_ip

			addEndpointGatewayIPOptions := vpcService.NewAddEndpointGatewayIPOptions(
				"testString",
				"testString",
			)

			reservedIP, response, err := vpcService.AddEndpointGatewayIP(addEndpointGatewayIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-add_endpoint_gateway_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`ListEndpointGatewayResourceBindings request example`, func() {
			fmt.Println("\nListEndpointGatewayResourceBindings() result:")
			// begin-list_endpoint_gateway_resource_bindings
			listEndpointGatewayResourceBindingsOptions := &vpcv1.ListEndpointGatewayResourceBindingsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewEndpointGatewayResourceBindingsPager(listEndpointGatewayResourceBindingsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.EndpointGatewayResourceBinding
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_endpoint_gateway_resource_bindings
		})
		It(`CreateEndpointGatewayResourceBinding request example`, func() {
			fmt.Println("\nCreateEndpointGatewayResourceBinding() result:")
			// begin-create_endpoint_gateway_resource_binding

			endpointGatewayResourceBindingTargetPrototypeModel := &vpcv1.EndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"),
			}

			createEndpointGatewayResourceBindingOptions := vpcService.NewCreateEndpointGatewayResourceBindingOptions(
				"testString",
				endpointGatewayResourceBindingTargetPrototypeModel,
			)

			endpointGatewayResourceBinding, response, err := vpcService.CreateEndpointGatewayResourceBinding(createEndpointGatewayResourceBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGatewayResourceBinding, "", "  ")
			fmt.Println(string(b))

			// end-create_endpoint_gateway_resource_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
		It(`GetEndpointGatewayResourceBinding request example`, func() {
			fmt.Println("\nGetEndpointGatewayResourceBinding() result:")
			// begin-get_endpoint_gateway_resource_binding

			getEndpointGatewayResourceBindingOptions := vpcService.NewGetEndpointGatewayResourceBindingOptions(
				"testString",
				"testString",
			)

			endpointGatewayResourceBinding, response, err := vpcService.GetEndpointGatewayResourceBinding(getEndpointGatewayResourceBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGatewayResourceBinding, "", "  ")
			fmt.Println(string(b))

			// end-get_endpoint_gateway_resource_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
		It(`UpdateEndpointGatewayResourceBinding request example`, func() {
			fmt.Println("\nUpdateEndpointGatewayResourceBinding() result:")
			// begin-update_endpoint_gateway_resource_binding

			endpointGatewayResourceBindingPatchModel := &vpcv1.EndpointGatewayResourceBindingPatch{}
			endpointGatewayResourceBindingPatchModelAsPatch, asPatchErr := endpointGatewayResourceBindingPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateEndpointGatewayResourceBindingOptions := vpcService.NewUpdateEndpointGatewayResourceBindingOptions(
				"testString",
				"testString",
				endpointGatewayResourceBindingPatchModelAsPatch,
			)

			endpointGatewayResourceBinding, response, err := vpcService.UpdateEndpointGatewayResourceBinding(updateEndpointGatewayResourceBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGatewayResourceBinding, "", "  ")
			fmt.Println(string(b))

			// end-update_endpoint_gateway_resource_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
		It(`GetEndpointGateway request example`, func() {
			fmt.Println("\nGetEndpointGateway() result:")
			// begin-get_endpoint_gateway

			getEndpointGatewayOptions := vpcService.NewGetEndpointGatewayOptions(
				"testString",
			)

			endpointGateway, response, err := vpcService.GetEndpointGateway(getEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGateway, "", "  ")
			fmt.Println(string(b))

			// end-get_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGateway).ToNot(BeNil())
		})
		It(`UpdateEndpointGateway request example`, func() {
			fmt.Println("\nUpdateEndpointGateway() result:")
			// begin-update_endpoint_gateway

			endpointGatewayPatchModel := &vpcv1.EndpointGatewayPatch{}
			endpointGatewayPatchModelAsPatch, asPatchErr := endpointGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateEndpointGatewayOptions := vpcService.NewUpdateEndpointGatewayOptions(
				"testString",
				endpointGatewayPatchModelAsPatch,
			)

			endpointGateway, response, err := vpcService.UpdateEndpointGateway(updateEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(endpointGateway, "", "  ")
			fmt.Println(string(b))

			// end-update_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGateway).ToNot(BeNil())
		})
		It(`ListFloatingIps request example`, func() {
			fmt.Println("\nListFloatingIps() result:")
			// begin-list_floating_ips
			listFloatingIpsOptions := &vpcv1.ListFloatingIpsOptions{
				Limit:              core.Int64Ptr(int64(10)),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewFloatingIpsPager(listFloatingIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.FloatingIP
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_floating_ips
		})
		It(`CreateFloatingIP request example`, func() {
			fmt.Println("\nCreateFloatingIP() result:")
			// begin-create_floating_ip

			floatingIPTargetPrototypeModel := &vpcv1.FloatingIPTargetPrototypeBareMetalServerNetworkInterfaceIdentityBareMetalServerNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/bare_metal_servers/0717-aad10fd7-8a02-4f3e-97f3-b18bd82cf304/network_interfaces/0717-da8c43ec-b6ca-4bd2-871e-72e288c66ee6"),
			}

			floatingIPPrototypeModel := &vpcv1.FloatingIPPrototypeFloatingIPByTarget{
				Target: floatingIPTargetPrototypeModel,
			}

			createFloatingIPOptions := vpcService.NewCreateFloatingIPOptions(
				floatingIPPrototypeModel,
			)

			floatingIP, response, err := vpcService.CreateFloatingIP(createFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-create_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`GetFloatingIP request example`, func() {
			fmt.Println("\nGetFloatingIP() result:")
			// begin-get_floating_ip

			getFloatingIPOptions := vpcService.NewGetFloatingIPOptions(
				"testString",
			)

			floatingIP, response, err := vpcService.GetFloatingIP(getFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-get_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`UpdateFloatingIP request example`, func() {
			fmt.Println("\nUpdateFloatingIP() result:")
			// begin-update_floating_ip

			floatingIPPatchModel := &vpcv1.FloatingIPPatch{}
			floatingIPPatchModelAsPatch, asPatchErr := floatingIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateFloatingIPOptions := vpcService.NewUpdateFloatingIPOptions(
				"testString",
				floatingIPPatchModelAsPatch,
			)

			floatingIP, response, err := vpcService.UpdateFloatingIP(updateFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-update_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`ListFlowLogCollectors request example`, func() {
			fmt.Println("\nListFlowLogCollectors() result:")
			// begin-list_flow_log_collectors
			listFlowLogCollectorsOptions := &vpcv1.ListFlowLogCollectorsOptions{
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				ResourceGroupID:    core.StringPtr("testString"),
				TargetID:           core.StringPtr("testString"),
				TargetResourceType: core.StringPtr("testString"),
				VPCCRN:             core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:              core.StringPtr("testString"),
				VPCName:            core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewFlowLogCollectorsPager(listFlowLogCollectorsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.FlowLogCollector
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_flow_log_collectors
		})
		It(`CreateFlowLogCollector request example`, func() {
			fmt.Println("\nCreateFlowLogCollector() result:")
			// begin-create_flow_log_collector

			flowLogCollectorPrototypeStorageBucketModel := &vpcv1.FlowLogCollectorPrototypeStorageBucketLegacyCloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName{
				Name: core.StringPtr("bucket-27200-lwx4cfvcue"),
			}

			flowLogCollectorTargetPrototypeModel := &vpcv1.FlowLogCollectorTargetPrototypeInstanceIdentityInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			createFlowLogCollectorOptions := vpcService.NewCreateFlowLogCollectorOptions(
				flowLogCollectorPrototypeStorageBucketModel,
				flowLogCollectorTargetPrototypeModel,
			)

			flowLogCollector, response, err := vpcService.CreateFlowLogCollector(createFlowLogCollectorOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(flowLogCollector, "", "  ")
			fmt.Println(string(b))

			// end-create_flow_log_collector

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(flowLogCollector).ToNot(BeNil())
		})
		It(`GetFlowLogCollector request example`, func() {
			fmt.Println("\nGetFlowLogCollector() result:")
			// begin-get_flow_log_collector

			getFlowLogCollectorOptions := vpcService.NewGetFlowLogCollectorOptions(
				"testString",
			)

			flowLogCollector, response, err := vpcService.GetFlowLogCollector(getFlowLogCollectorOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(flowLogCollector, "", "  ")
			fmt.Println(string(b))

			// end-get_flow_log_collector

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(flowLogCollector).ToNot(BeNil())
		})
		It(`UpdateFlowLogCollector request example`, func() {
			fmt.Println("\nUpdateFlowLogCollector() result:")
			// begin-update_flow_log_collector

			flowLogCollectorPatchModel := &vpcv1.FlowLogCollectorPatch{
				Active: core.BoolPtr(false),
			}
			flowLogCollectorPatchModelAsPatch, asPatchErr := flowLogCollectorPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateFlowLogCollectorOptions := vpcService.NewUpdateFlowLogCollectorOptions(
				"testString",
				flowLogCollectorPatchModelAsPatch,
			)

			flowLogCollector, response, err := vpcService.UpdateFlowLogCollector(updateFlowLogCollectorOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(flowLogCollector, "", "  ")
			fmt.Println(string(b))

			// end-update_flow_log_collector

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(flowLogCollector).ToNot(BeNil())
		})
		It(`ListRegions request example`, func() {
			fmt.Println("\nListRegions() result:")
			// begin-list_regions

			listRegionsOptions := vpcService.NewListRegionsOptions()

			regionCollection, response, err := vpcService.ListRegions(listRegionsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(regionCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_regions

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(regionCollection).ToNot(BeNil())
		})
		It(`GetRegion request example`, func() {
			fmt.Println("\nGetRegion() result:")
			// begin-get_region

			getRegionOptions := vpcService.NewGetRegionOptions(
				"us-south",
			)

			region, response, err := vpcService.GetRegion(getRegionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(region, "", "  ")
			fmt.Println(string(b))

			// end-get_region

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(region).ToNot(BeNil())
		})
		It(`ListRegionZones request example`, func() {
			fmt.Println("\nListRegionZones() result:")
			// begin-list_region_zones

			listRegionZonesOptions := vpcService.NewListRegionZonesOptions(
				"us-south",
			)

			zoneCollection, response, err := vpcService.ListRegionZones(listRegionZonesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(zoneCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_region_zones

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(zoneCollection).ToNot(BeNil())
		})
		It(`GetRegionZone request example`, func() {
			fmt.Println("\nGetRegionZone() result:")
			// begin-get_region_zone

			getRegionZoneOptions := vpcService.NewGetRegionZoneOptions(
				"us-south-1",
				"us-south",
			)

			zone, response, err := vpcService.GetRegionZone(getRegionZoneOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(zone, "", "  ")
			fmt.Println(string(b))

			// end-get_region_zone

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(zone).ToNot(BeNil())
		})
		It(`ListImages request example`, func() {
			fmt.Println("\nListImages() result:")
			// begin-list_images
			listImagesOptions := &vpcv1.ListImagesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				RemoteAccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				ResourceGroupID: core.StringPtr("testString"),
				Status:          []string{"available"},
				UserDataFormat:  []string{"cloud_init"},
				Visibility:      core.StringPtr("private"),
			}

			pager, err := vpcService.NewImagesPager(listImagesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Image
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_images
		})
		It(`CreateImage request example`, func() {
			fmt.Println("\nCreateImage() result:")
			// begin-create_image

			imageByFileFileModel := &vpcv1.ImageByFileFile{
				Href: core.StringPtr("cos://us-south/custom-image-vpc-bucket/customImage-0.vhd"),
			}

			imageByFileOperatingSystemModel := &vpcv1.ImageByFileOperatingSystemOperatingSystemIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/operating_systems/ubuntu-24-04-amd64"),
			}

			imagePrototypeModel := &vpcv1.ImagePrototypeImageByFile{
				File:            imageByFileFileModel,
				OperatingSystem: imageByFileOperatingSystemModel,
			}

			createImageOptions := vpcService.NewCreateImageOptions(
				imagePrototypeModel,
			)

			image, response, err := vpcService.CreateImage(createImageOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(b))

			// end-create_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(image).ToNot(BeNil())
		})
		It(`GetImage request example`, func() {
			fmt.Println("\nGetImage() result:")
			// begin-get_image

			getImageOptions := vpcService.NewGetImageOptions(
				"testString",
			)

			image, response, err := vpcService.GetImage(getImageOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(b))

			// end-get_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(image).ToNot(BeNil())
		})
		It(`UpdateImage request example`, func() {
			fmt.Println("\nUpdateImage() result:")
			// begin-update_image

			imagePatchModel := &vpcv1.ImagePatch{}
			imagePatchModelAsPatch, asPatchErr := imagePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateImageOptions := vpcService.NewUpdateImageOptions(
				"testString",
				imagePatchModelAsPatch,
			)

			image, response, err := vpcService.UpdateImage(updateImageOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(b))

			// end-update_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(image).ToNot(BeNil())
		})
		It(`ListImageBareMetalServerProfiles request example`, func() {
			fmt.Println("\nListImageBareMetalServerProfiles() result:")
			// begin-list_image_bare_metal_server_profiles
			listImageBareMetalServerProfilesOptions := &vpcv1.ListImageBareMetalServerProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewImageBareMetalServerProfilesPager(listImageBareMetalServerProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.BareMetalServerProfileReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_image_bare_metal_server_profiles
		})
		It(`DeprecateImage request example`, func() {
			// begin-deprecate_image

			deprecateImageOptions := vpcService.NewDeprecateImageOptions(
				"testString",
			)

			response, err := vpcService.DeprecateImage(deprecateImageOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeprecateImage(): %d\n", response.StatusCode)
			}

			// end-deprecate_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`ListImageInstanceProfiles request example`, func() {
			fmt.Println("\nListImageInstanceProfiles() result:")
			// begin-list_image_instance_profiles
			listImageInstanceProfilesOptions := &vpcv1.ListImageInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewImageInstanceProfilesPager(listImageInstanceProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_image_instance_profiles
		})
		It(`ObsoleteImage request example`, func() {
			// begin-obsolete_image

			obsoleteImageOptions := vpcService.NewObsoleteImageOptions(
				"testString",
			)

			response, err := vpcService.ObsoleteImage(obsoleteImageOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from ObsoleteImage(): %d\n", response.StatusCode)
			}

			// end-obsolete_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`ListImageExportJobs request example`, func() {
			fmt.Println("\nListImageExportJobs() result:")
			// begin-list_image_export_jobs

			listImageExportJobsOptions := vpcService.NewListImageExportJobsOptions(
				"testString",
			)
			listImageExportJobsOptions.SetName("my-name")

			imageExportJobUnpaginatedCollection, response, err := vpcService.ListImageExportJobs(listImageExportJobsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(imageExportJobUnpaginatedCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_image_export_jobs

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJobUnpaginatedCollection).ToNot(BeNil())
		})
		It(`CreateImageExportJob request example`, func() {
			fmt.Println("\nCreateImageExportJob() result:")
			// begin-create_image_export_job

			cloudObjectStorageBucketIdentityModel := &vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName{
				Name: core.StringPtr("bucket-27200-lwx4cfvcue"),
			}

			createImageExportJobOptions := vpcService.NewCreateImageExportJobOptions(
				"testString",
				cloudObjectStorageBucketIdentityModel,
			)
			createImageExportJobOptions.SetName("my-image-export")

			imageExportJob, response, err := vpcService.CreateImageExportJob(createImageExportJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(imageExportJob, "", "  ")
			fmt.Println(string(b))

			// end-create_image_export_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(imageExportJob).ToNot(BeNil())
		})
		It(`GetImageExportJob request example`, func() {
			fmt.Println("\nGetImageExportJob() result:")
			// begin-get_image_export_job

			getImageExportJobOptions := vpcService.NewGetImageExportJobOptions(
				"testString",
				"testString",
			)

			imageExportJob, response, err := vpcService.GetImageExportJob(getImageExportJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(imageExportJob, "", "  ")
			fmt.Println(string(b))

			// end-get_image_export_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJob).ToNot(BeNil())
		})
		It(`UpdateImageExportJob request example`, func() {
			fmt.Println("\nUpdateImageExportJob() result:")
			// begin-update_image_export_job

			imageExportJobPatchModel := &vpcv1.ImageExportJobPatch{
				Name: core.StringPtr("my-image-export-job-updated"),
			}
			imageExportJobPatchModelAsPatch, asPatchErr := imageExportJobPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateImageExportJobOptions := vpcService.NewUpdateImageExportJobOptions(
				"testString",
				"testString",
				imageExportJobPatchModelAsPatch,
			)

			imageExportJob, response, err := vpcService.UpdateImageExportJob(updateImageExportJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(imageExportJob, "", "  ")
			fmt.Println(string(b))

			// end-update_image_export_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJob).ToNot(BeNil())
		})
		It(`ListOperatingSystems request example`, func() {
			fmt.Println("\nListOperatingSystems() result:")
			// begin-list_operating_systems
			listOperatingSystemsOptions := &vpcv1.ListOperatingSystemsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewOperatingSystemsPager(listOperatingSystemsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.OperatingSystem
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_operating_systems
		})
		It(`GetOperatingSystem request example`, func() {
			fmt.Println("\nGetOperatingSystem() result:")
			// begin-get_operating_system

			getOperatingSystemOptions := vpcService.NewGetOperatingSystemOptions(
				"red-7-amd64",
			)

			operatingSystem, response, err := vpcService.GetOperatingSystem(getOperatingSystemOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(operatingSystem, "", "  ")
			fmt.Println(string(b))

			// end-get_operating_system

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(operatingSystem).ToNot(BeNil())
		})
		It(`ListInstanceGroups request example`, func() {
			fmt.Println("\nListInstanceGroups() result:")
			// begin-list_instance_groups
			listInstanceGroupsOptions := &vpcv1.ListInstanceGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceGroupsPager(listInstanceGroupsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceGroup
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_groups
		})
		It(`CreateInstanceGroup request example`, func() {
			fmt.Println("\nCreateInstanceGroup() result:")
			// begin-create_instance_group

			instanceTemplateIdentityModel := &vpcv1.InstanceTemplateIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance-template:0717-e6c2c7d8-ad57-4f38-a21c-a86265b6aeb2"),
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			createInstanceGroupOptions := vpcService.NewCreateInstanceGroupOptions(
				instanceTemplateIdentityModel,
				[]vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			)

			instanceGroup, response, err := vpcService.CreateInstanceGroup(createInstanceGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroup, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroup).ToNot(BeNil())
		})
		It(`GetInstanceGroup request example`, func() {
			fmt.Println("\nGetInstanceGroup() result:")
			// begin-get_instance_group

			getInstanceGroupOptions := vpcService.NewGetInstanceGroupOptions(
				"testString",
			)

			instanceGroup, response, err := vpcService.GetInstanceGroup(getInstanceGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroup).ToNot(BeNil())
		})
		It(`UpdateInstanceGroup request example`, func() {
			fmt.Println("\nUpdateInstanceGroup() result:")
			// begin-update_instance_group

			instanceGroupPatchModel := &vpcv1.InstanceGroupPatch{
				MembershipCount: core.Int64Ptr(int64(20)),
			}
			instanceGroupPatchModelAsPatch, asPatchErr := instanceGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupOptions := vpcService.NewUpdateInstanceGroupOptions(
				"testString",
				instanceGroupPatchModelAsPatch,
			)

			instanceGroup, response, err := vpcService.UpdateInstanceGroup(updateInstanceGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroup, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroup).ToNot(BeNil())
		})
		It(`ListInstanceGroupManagers request example`, func() {
			fmt.Println("\nListInstanceGroupManagers() result:")
			// begin-list_instance_group_managers
			listInstanceGroupManagersOptions := &vpcv1.ListInstanceGroupManagersOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceGroupManagersPager(listInstanceGroupManagersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceGroupManagerIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_group_managers
		})
		It(`CreateInstanceGroupManager request example`, func() {
			fmt.Println("\nCreateInstanceGroupManager() result:")
			// begin-create_instance_group_manager

			instanceGroupManagerPrototypeModel := &vpcv1.InstanceGroupManagerPrototypeInstanceGroupManagerAutoScalePrototype{
				ManagerType:        core.StringPtr("autoscale"),
				MaxMembershipCount: core.Int64Ptr(int64(10)),
			}

			createInstanceGroupManagerOptions := vpcService.NewCreateInstanceGroupManagerOptions(
				"testString",
				instanceGroupManagerPrototypeModel,
			)

			instanceGroupManager, response, err := vpcService.CreateInstanceGroupManager(createInstanceGroupManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManager, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_group_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
		It(`GetInstanceGroupManager request example`, func() {
			fmt.Println("\nGetInstanceGroupManager() result:")
			// begin-get_instance_group_manager

			getInstanceGroupManagerOptions := vpcService.NewGetInstanceGroupManagerOptions(
				"testString",
				"testString",
			)

			instanceGroupManager, response, err := vpcService.GetInstanceGroupManager(getInstanceGroupManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManager, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_group_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
		It(`UpdateInstanceGroupManager request example`, func() {
			fmt.Println("\nUpdateInstanceGroupManager() result:")
			// begin-update_instance_group_manager

			instanceGroupManagerPatchModel := &vpcv1.InstanceGroupManagerPatch{
				ManagementEnabled: core.BoolPtr(false),
			}
			instanceGroupManagerPatchModelAsPatch, asPatchErr := instanceGroupManagerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerOptions := vpcService.NewUpdateInstanceGroupManagerOptions(
				"testString",
				"testString",
				instanceGroupManagerPatchModelAsPatch,
			)

			instanceGroupManager, response, err := vpcService.UpdateInstanceGroupManager(updateInstanceGroupManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManager, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_group_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
		It(`ListInstanceGroupManagerActions request example`, func() {
			fmt.Println("\nListInstanceGroupManagerActions() result:")
			// begin-list_instance_group_manager_actions
			listInstanceGroupManagerActionsOptions := &vpcv1.ListInstanceGroupManagerActionsOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceGroupManagerActionsPager(listInstanceGroupManagerActionsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceGroupManagerActionIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_group_manager_actions
		})
		It(`CreateInstanceGroupManagerAction request example`, func() {
			fmt.Println("\nCreateInstanceGroupManagerAction() result:")
			// begin-create_instance_group_manager_action

			instanceGroupManagerScheduledActionGroupPrototypeModel := &vpcv1.InstanceGroupManagerScheduledActionGroupPrototype{
				MembershipCount: core.Int64Ptr(int64(10)),
			}

			instanceGroupManagerActionPrototypeModel := &vpcv1.InstanceGroupManagerActionPrototypeScheduledActionPrototypeByCronSpecWithGroup{
				CronSpec: core.StringPtr("30 */2 * * 1-5"),
				Group:    instanceGroupManagerScheduledActionGroupPrototypeModel,
			}

			createInstanceGroupManagerActionOptions := vpcService.NewCreateInstanceGroupManagerActionOptions(
				"testString",
				"testString",
				instanceGroupManagerActionPrototypeModel,
			)

			instanceGroupManagerAction, response, err := vpcService.CreateInstanceGroupManagerAction(createInstanceGroupManagerActionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerAction, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_group_manager_action

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
		It(`GetInstanceGroupManagerAction request example`, func() {
			fmt.Println("\nGetInstanceGroupManagerAction() result:")
			// begin-get_instance_group_manager_action

			getInstanceGroupManagerActionOptions := vpcService.NewGetInstanceGroupManagerActionOptions(
				"testString",
				"testString",
				"testString",
			)

			instanceGroupManagerAction, response, err := vpcService.GetInstanceGroupManagerAction(getInstanceGroupManagerActionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerAction, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_group_manager_action

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
		It(`UpdateInstanceGroupManagerAction request example`, func() {
			fmt.Println("\nUpdateInstanceGroupManagerAction() result:")
			// begin-update_instance_group_manager_action

			instanceGroupManagerActionPatchModel := &vpcv1.InstanceGroupManagerActionPatch{
				AutoDelete: core.BoolPtr(false),
			}
			instanceGroupManagerActionPatchModelAsPatch, asPatchErr := instanceGroupManagerActionPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerActionOptions := vpcService.NewUpdateInstanceGroupManagerActionOptions(
				"testString",
				"testString",
				"testString",
				instanceGroupManagerActionPatchModelAsPatch,
			)

			instanceGroupManagerAction, response, err := vpcService.UpdateInstanceGroupManagerAction(updateInstanceGroupManagerActionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerAction, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_group_manager_action

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
		It(`ListInstanceGroupManagerPolicies request example`, func() {
			fmt.Println("\nListInstanceGroupManagerPolicies() result:")
			// begin-list_instance_group_manager_policies
			listInstanceGroupManagerPoliciesOptions := &vpcv1.ListInstanceGroupManagerPoliciesOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceGroupManagerPoliciesPager(listInstanceGroupManagerPoliciesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceGroupManagerPolicyIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_group_manager_policies
		})
		It(`CreateInstanceGroupManagerPolicy request example`, func() {
			fmt.Println("\nCreateInstanceGroupManagerPolicy() result:")
			// begin-create_instance_group_manager_policy

			instanceGroupManagerPolicyPrototypeModel := &vpcv1.InstanceGroupManagerPolicyPrototypeInstanceGroupManagerTargetPolicyPrototype{
				MetricType:  core.StringPtr("cpu"),
				MetricValue: core.Int64Ptr(int64(0)),
				PolicyType:  core.StringPtr("target"),
			}

			createInstanceGroupManagerPolicyOptions := vpcService.NewCreateInstanceGroupManagerPolicyOptions(
				"testString",
				"testString",
				instanceGroupManagerPolicyPrototypeModel,
			)

			instanceGroupManagerPolicy, response, err := vpcService.CreateInstanceGroupManagerPolicy(createInstanceGroupManagerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_group_manager_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
		It(`GetInstanceGroupManagerPolicy request example`, func() {
			fmt.Println("\nGetInstanceGroupManagerPolicy() result:")
			// begin-get_instance_group_manager_policy

			getInstanceGroupManagerPolicyOptions := vpcService.NewGetInstanceGroupManagerPolicyOptions(
				"testString",
				"testString",
				"testString",
			)

			instanceGroupManagerPolicy, response, err := vpcService.GetInstanceGroupManagerPolicy(getInstanceGroupManagerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_group_manager_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
		It(`UpdateInstanceGroupManagerPolicy request example`, func() {
			fmt.Println("\nUpdateInstanceGroupManagerPolicy() result:")
			// begin-update_instance_group_manager_policy

			instanceGroupManagerPolicyPatchModel := &vpcv1.InstanceGroupManagerPolicyPatch{
				MetricValue: core.Int64Ptr(int64(75)),
			}
			instanceGroupManagerPolicyPatchModelAsPatch, asPatchErr := instanceGroupManagerPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerPolicyOptions := vpcService.NewUpdateInstanceGroupManagerPolicyOptions(
				"testString",
				"testString",
				"testString",
				instanceGroupManagerPolicyPatchModelAsPatch,
			)

			instanceGroupManagerPolicy, response, err := vpcService.UpdateInstanceGroupManagerPolicy(updateInstanceGroupManagerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupManagerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_group_manager_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
		It(`ListInstanceGroupMemberships request example`, func() {
			fmt.Println("\nListInstanceGroupMemberships() result:")
			// begin-list_instance_group_memberships
			listInstanceGroupMembershipsOptions := &vpcv1.ListInstanceGroupMembershipsOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceGroupMembershipsPager(listInstanceGroupMembershipsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceGroupMembership
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_group_memberships
		})
		It(`GetInstanceGroupMembership request example`, func() {
			fmt.Println("\nGetInstanceGroupMembership() result:")
			// begin-get_instance_group_membership

			getInstanceGroupMembershipOptions := vpcService.NewGetInstanceGroupMembershipOptions(
				"testString",
				"testString",
			)

			instanceGroupMembership, response, err := vpcService.GetInstanceGroupMembership(getInstanceGroupMembershipOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupMembership, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_group_membership

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupMembership).ToNot(BeNil())
		})
		It(`UpdateInstanceGroupMembership request example`, func() {
			fmt.Println("\nUpdateInstanceGroupMembership() result:")
			// begin-update_instance_group_membership

			instanceGroupMembershipPatchModel := &vpcv1.InstanceGroupMembershipPatch{
				DeleteInstanceOnMembershipDelete: core.BoolPtr(false),
			}
			instanceGroupMembershipPatchModelAsPatch, asPatchErr := instanceGroupMembershipPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupMembershipOptions := vpcService.NewUpdateInstanceGroupMembershipOptions(
				"testString",
				"testString",
				instanceGroupMembershipPatchModelAsPatch,
			)

			instanceGroupMembership, response, err := vpcService.UpdateInstanceGroupMembership(updateInstanceGroupMembershipOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceGroupMembership, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_group_membership

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupMembership).ToNot(BeNil())
		})
		It(`ListInstanceTemplates request example`, func() {
			fmt.Println("\nListInstanceTemplates() result:")
			// begin-list_instance_templates

			listInstanceTemplatesOptions := vpcService.NewListInstanceTemplatesOptions()

			instanceTemplateCollection, response, err := vpcService.ListInstanceTemplates(listInstanceTemplatesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceTemplateCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_templates

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplateCollection).ToNot(BeNil())
		})
		It(`CreateInstanceTemplate request example`, func() {
			fmt.Println("\nCreateInstanceTemplate() result:")
			// begin-create_instance_template

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			instanceCatalogOfferingPrototypeModel := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
				Offering: catalogOfferingIdentityModel,
			}

			instanceTemplateByCatalogOfferingZoneModel := &vpcv1.InstanceTemplateByCatalogOfferingZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			instanceByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.InstanceByNetworkAttachmentPrimaryNetworkAttachment{
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instanceTemplatePrototypeModel := &vpcv1.InstanceTemplatePrototypeInstanceTemplateByCatalogOfferingInstanceTemplateByCatalogOfferingInstanceByNetworkAttachment{
				CatalogOffering:          instanceCatalogOfferingPrototypeModel,
				Zone:                     instanceTemplateByCatalogOfferingZoneModel,
				PrimaryNetworkAttachment: instanceByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createInstanceTemplateOptions := vpcService.NewCreateInstanceTemplateOptions(
				instanceTemplatePrototypeModel,
			)

			instanceTemplate, response, err := vpcService.CreateInstanceTemplate(createInstanceTemplateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceTemplate, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_template

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceTemplate).ToNot(BeNil())
		})
		It(`GetInstanceTemplate request example`, func() {
			fmt.Println("\nGetInstanceTemplate() result:")
			// begin-get_instance_template

			getInstanceTemplateOptions := vpcService.NewGetInstanceTemplateOptions(
				"testString",
			)

			instanceTemplate, response, err := vpcService.GetInstanceTemplate(getInstanceTemplateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceTemplate, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_template

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplate).ToNot(BeNil())
		})
		It(`UpdateInstanceTemplate request example`, func() {
			fmt.Println("\nUpdateInstanceTemplate() result:")
			// begin-update_instance_template

			instanceTemplatePatchModel := &vpcv1.InstanceTemplatePatch{
				Name: core.StringPtr("my-instance-template-updated"),
			}
			instanceTemplatePatchModelAsPatch, asPatchErr := instanceTemplatePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceTemplateOptions := vpcService.NewUpdateInstanceTemplateOptions(
				"testString",
				instanceTemplatePatchModelAsPatch,
			)

			instanceTemplate, response, err := vpcService.UpdateInstanceTemplate(updateInstanceTemplateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceTemplate, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_template

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplate).ToNot(BeNil())
		})
		It(`ListInstanceProfiles request example`, func() {
			fmt.Println("\nListInstanceProfiles() result:")
			// begin-list_instance_profiles
			listInstanceProfilesOptions := &vpcv1.ListInstanceProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceProfilesPager(listInstanceProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_profiles
		})
		It(`GetInstanceProfile request example`, func() {
			fmt.Println("\nGetInstanceProfile() result:")
			// begin-get_instance_profile

			getInstanceProfileOptions := vpcService.NewGetInstanceProfileOptions(
				"mx2-host-152x1216",
			)

			instanceProfile, response, err := vpcService.GetInstanceProfile(getInstanceProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceProfile).ToNot(BeNil())
		})
		It(`ListInstances request example`, func() {
			fmt.Println("\nListInstances() result:")
			// begin-list_instances
			listInstancesOptions := &vpcv1.ListInstancesOptions{
				AvailabilityClass:                       core.StringPtr("standard"),
				ClusterNetworkCRN:                       core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::cluster-network:0717-da0df18c-7598-4633-a648-fdaac28a5573"),
				ClusterNetworkID:                        core.StringPtr("testString"),
				ClusterNetworkName:                      core.StringPtr("my-cluster-network"),
				DedicatedHostCRN:                        core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host:0717-1e09281b-f177-46fb-baf1-bc152b2e391a"),
				DedicatedHostID:                         core.StringPtr("testString"),
				DedicatedHostName:                       core.StringPtr("my-dedicated-host"),
				InstanceGroupMembershipInstanceGroupCRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::instance-group:r006-7b3ac170-01f3-43d6-87ec-f0ed11ed3f60"),
				InstanceGroupMembershipInstanceGroupID:  core.StringPtr("r006-7b3ac170-01f3-43d6-87ec-f0ed11ed3f60"),
				Limit:                                   core.Int64Ptr(int64(10)),
				Name:                                    core.StringPtr("my-name"),
				PlacementGroupCRN:                       core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::placement-group:r006-418fe842-a3e9-47b9-a938-1aa5bd632871"),
				PlacementGroupID:                        core.StringPtr("testString"),
				PlacementGroupName:                      core.StringPtr("my-placement-group"),
				ReservationAffinityPolicy:               core.StringPtr("automatic"),
				ReservationCRN:                          core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
				ReservationID:                           core.StringPtr("testString"),
				ReservationName:                         core.StringPtr("my-reservation"),
				ResourceGroupID:                         core.StringPtr("testString"),
				VPCCRN:                                  core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:                                   core.StringPtr("testString"),
				VPCName:                                 core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewInstancesPager(listInstancesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Instance
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instances
		})
		It(`CreateInstance request example`, func() {
			fmt.Println("\nCreateInstance() result:")
			// begin-create_instance

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			instanceCatalogOfferingPrototypeModel := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
				Offering: catalogOfferingIdentityModel,
			}

			instanceByCatalogOfferingZoneModel := &vpcv1.InstanceByCatalogOfferingZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			instanceByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.InstanceByNetworkAttachmentPrimaryNetworkAttachment{
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instancePrototypeModel := &vpcv1.InstancePrototypeInstanceByCatalogOfferingInstanceByCatalogOfferingInstanceByNetworkAttachment{
				CatalogOffering:          instanceCatalogOfferingPrototypeModel,
				Zone:                     instanceByCatalogOfferingZoneModel,
				PrimaryNetworkAttachment: instanceByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createInstanceOptions := vpcService.NewCreateInstanceOptions(
				instancePrototypeModel,
			)

			instance, response, err := vpcService.CreateInstance(createInstanceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instance, "", "  ")
			fmt.Println(string(b))

			// end-create_instance

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instance).ToNot(BeNil())
		})
		It(`GetInstance request example`, func() {
			fmt.Println("\nGetInstance() result:")
			// begin-get_instance

			getInstanceOptions := vpcService.NewGetInstanceOptions(
				"testString",
			)

			instance, response, err := vpcService.GetInstance(getInstanceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instance, "", "  ")
			fmt.Println(string(b))

			// end-get_instance

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instance).ToNot(BeNil())
		})
		It(`UpdateInstance request example`, func() {
			fmt.Println("\nUpdateInstance() result:")
			// begin-update_instance

			instancePatchModel := &vpcv1.InstancePatch{}
			instancePatchModelAsPatch, asPatchErr := instancePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceOptions := vpcService.NewUpdateInstanceOptions(
				"testString",
				instancePatchModelAsPatch,
			)
			updateInstanceOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			instance, response, err := vpcService.UpdateInstance(updateInstanceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instance, "", "  ")
			fmt.Println(string(b))

			// end-update_instance

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instance).ToNot(BeNil())
		})
		It(`GetInstanceInitialization request example`, func() {
			fmt.Println("\nGetInstanceInitialization() result:")
			// begin-get_instance_initialization

			getInstanceInitializationOptions := vpcService.NewGetInstanceInitializationOptions(
				"testString",
			)

			instanceInitialization, response, err := vpcService.GetInstanceInitialization(getInstanceInitializationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceInitialization, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_initialization

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceInitialization).ToNot(BeNil())
		})
		It(`CreateInstanceReinitialization request example`, func() {
			// begin-create_instance_reinitialization

			instanceReinitializeByImageImageModel := &vpcv1.InstanceReinitializeByImageImageImageIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"),
			}

			instanceReinitializePrototypeModel := &vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage{
				Image: instanceReinitializeByImageImageModel,
			}

			createInstanceReinitializationOptions := vpcService.NewCreateInstanceReinitializationOptions(
				"testString",
				instanceReinitializePrototypeModel,
			)

			response, err := vpcService.CreateInstanceReinitialization(createInstanceReinitializationOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from CreateInstanceReinitialization(): %d\n", response.StatusCode)
			}

			// end-create_instance_reinitialization

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`CreateInstanceAction request example`, func() {
			fmt.Println("\nCreateInstanceAction() result:")
			// begin-create_instance_action

			createInstanceActionOptions := vpcService.NewCreateInstanceActionOptions(
				"testString",
				"reboot",
			)

			instanceAction, response, err := vpcService.CreateInstanceAction(createInstanceActionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceAction, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_action

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceAction).ToNot(BeNil())
		})
		It(`ListInstanceClusterNetworkAttachments request example`, func() {
			fmt.Println("\nListInstanceClusterNetworkAttachments() result:")
			// begin-list_instance_cluster_network_attachments
			listInstanceClusterNetworkAttachmentsOptions := &vpcv1.ListInstanceClusterNetworkAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
				Limit:      core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceClusterNetworkAttachmentsPager(listInstanceClusterNetworkAttachmentsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceClusterNetworkAttachment
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_cluster_network_attachments
		})
		It(`CreateClusterNetworkAttachment request example`, func() {
			fmt.Println("\nCreateClusterNetworkAttachment() result:")
			// begin-create_cluster_network_attachment

			instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"),
			}

			createClusterNetworkAttachmentOptions := vpcService.NewCreateClusterNetworkAttachmentOptions(
				"testString",
				instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel,
			)

			instanceClusterNetworkAttachment, response, err := vpcService.CreateClusterNetworkAttachment(createClusterNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceClusterNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-create_cluster_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
		It(`GetInstanceClusterNetworkAttachment request example`, func() {
			fmt.Println("\nGetInstanceClusterNetworkAttachment() result:")
			// begin-get_instance_cluster_network_attachment

			getInstanceClusterNetworkAttachmentOptions := vpcService.NewGetInstanceClusterNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			instanceClusterNetworkAttachment, response, err := vpcService.GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceClusterNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_cluster_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
		It(`UpdateInstanceClusterNetworkAttachment request example`, func() {
			fmt.Println("\nUpdateInstanceClusterNetworkAttachment() result:")
			// begin-update_instance_cluster_network_attachment

			instanceClusterNetworkAttachmentPatchModel := &vpcv1.InstanceClusterNetworkAttachmentPatch{
				Name: core.StringPtr("my-instance-network-attachment-updated"),
			}
			instanceClusterNetworkAttachmentPatchModelAsPatch, asPatchErr := instanceClusterNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceClusterNetworkAttachmentOptions := vpcService.NewUpdateInstanceClusterNetworkAttachmentOptions(
				"testString",
				"testString",
				instanceClusterNetworkAttachmentPatchModelAsPatch,
			)

			instanceClusterNetworkAttachment, response, err := vpcService.UpdateInstanceClusterNetworkAttachment(updateInstanceClusterNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceClusterNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_cluster_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
		It(`CreateInstanceConsoleAccessToken request example`, func() {
			fmt.Println("\nCreateInstanceConsoleAccessToken() result:")
			// begin-create_instance_console_access_token

			createInstanceConsoleAccessTokenOptions := vpcService.NewCreateInstanceConsoleAccessTokenOptions(
				"testString",
				"serial",
			)

			instanceConsoleAccessToken, response, err := vpcService.CreateInstanceConsoleAccessToken(createInstanceConsoleAccessTokenOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceConsoleAccessToken, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_console_access_token

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceConsoleAccessToken).ToNot(BeNil())
		})
		It(`ListInstanceDisks request example`, func() {
			fmt.Println("\nListInstanceDisks() result:")
			// begin-list_instance_disks

			listInstanceDisksOptions := vpcService.NewListInstanceDisksOptions(
				"testString",
			)

			instanceDiskCollection, response, err := vpcService.ListInstanceDisks(listInstanceDisksOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceDiskCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_disks

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDiskCollection).ToNot(BeNil())
		})
		It(`GetInstanceDisk request example`, func() {
			fmt.Println("\nGetInstanceDisk() result:")
			// begin-get_instance_disk

			getInstanceDiskOptions := vpcService.NewGetInstanceDiskOptions(
				"testString",
				"testString",
			)

			instanceDisk, response, err := vpcService.GetInstanceDisk(getInstanceDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceDisk, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDisk).ToNot(BeNil())
		})
		It(`UpdateInstanceDisk request example`, func() {
			fmt.Println("\nUpdateInstanceDisk() result:")
			// begin-update_instance_disk

			instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{}
			instanceDiskPatchModelAsPatch, asPatchErr := instanceDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceDiskOptions := vpcService.NewUpdateInstanceDiskOptions(
				"testString",
				"testString",
				instanceDiskPatchModelAsPatch,
			)

			instanceDisk, response, err := vpcService.UpdateInstanceDisk(updateInstanceDiskOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceDisk, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_disk

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDisk).ToNot(BeNil())
		})
		It(`ListInstanceNetworkAttachments request example`, func() {
			fmt.Println("\nListInstanceNetworkAttachments() result:")
			// begin-list_instance_network_attachments

			listInstanceNetworkAttachmentsOptions := vpcService.NewListInstanceNetworkAttachmentsOptions(
				"testString",
			)

			instanceNetworkAttachmentCollection, response, err := vpcService.ListInstanceNetworkAttachments(listInstanceNetworkAttachmentsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceNetworkAttachmentCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_network_attachments

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachmentCollection).ToNot(BeNil())
		})
		It(`CreateInstanceNetworkAttachment request example`, func() {
			fmt.Println("\nCreateInstanceNetworkAttachment() result:")
			// begin-create_instance_network_attachment

			virtualNetworkInterfacePrimaryIPPrototypeModel := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext{
				Address:    core.StringPtr("10.0.0.7"),
				AutoDelete: core.BoolPtr(false),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext{
				PrimaryIP: virtualNetworkInterfacePrimaryIPPrototypeModel,
			}

			createInstanceNetworkAttachmentOptions := vpcService.NewCreateInstanceNetworkAttachmentOptions(
				"testString",
				instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			)

			instanceNetworkAttachment, response, err := vpcService.CreateInstanceNetworkAttachment(createInstanceNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
		It(`GetInstanceNetworkAttachment request example`, func() {
			fmt.Println("\nGetInstanceNetworkAttachment() result:")
			// begin-get_instance_network_attachment

			getInstanceNetworkAttachmentOptions := vpcService.NewGetInstanceNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			instanceNetworkAttachment, response, err := vpcService.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
		It(`UpdateInstanceNetworkAttachment request example`, func() {
			fmt.Println("\nUpdateInstanceNetworkAttachment() result:")
			// begin-update_instance_network_attachment

			instanceNetworkAttachmentPatchModel := &vpcv1.InstanceNetworkAttachmentPatch{
				Name: core.StringPtr("my-instance-network-attachment-updated"),
			}
			instanceNetworkAttachmentPatchModelAsPatch, asPatchErr := instanceNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceNetworkAttachmentOptions := vpcService.NewUpdateInstanceNetworkAttachmentOptions(
				"testString",
				"testString",
				instanceNetworkAttachmentPatchModelAsPatch,
			)

			instanceNetworkAttachment, response, err := vpcService.UpdateInstanceNetworkAttachment(updateInstanceNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
		It(`ListInstanceNetworkInterfaces request example`, func() {
			fmt.Println("\nListInstanceNetworkInterfaces() result:")
			// begin-list_instance_network_interfaces

			listInstanceNetworkInterfacesOptions := vpcService.NewListInstanceNetworkInterfacesOptions(
				"testString",
			)

			networkInterfaceUnpaginatedCollection, response, err := vpcService.ListInstanceNetworkInterfaces(listInstanceNetworkInterfacesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkInterfaceUnpaginatedCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_network_interfaces

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterfaceUnpaginatedCollection).ToNot(BeNil())
		})
		It(`CreateInstanceNetworkInterface request example`, func() {
			fmt.Println("\nCreateInstanceNetworkInterface() result:")
			// begin-create_instance_network_interface

			networkInterfacePrototypeSubnetModel := &vpcv1.NetworkInterfacePrototypeSubnetSubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			createInstanceNetworkInterfaceOptions := vpcService.NewCreateInstanceNetworkInterfaceOptions(
				"testString",
				networkInterfacePrototypeSubnetModel,
			)

			networkInterface, response, err := vpcService.CreateInstanceNetworkInterface(createInstanceNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkInterface, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkInterface).ToNot(BeNil())
		})
		It(`GetInstanceNetworkInterface request example`, func() {
			fmt.Println("\nGetInstanceNetworkInterface() result:")
			// begin-get_instance_network_interface

			getInstanceNetworkInterfaceOptions := vpcService.NewGetInstanceNetworkInterfaceOptions(
				"testString",
				"testString",
			)

			networkInterface, response, err := vpcService.GetInstanceNetworkInterface(getInstanceNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkInterface, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterface).ToNot(BeNil())
		})
		It(`UpdateInstanceNetworkInterface request example`, func() {
			fmt.Println("\nUpdateInstanceNetworkInterface() result:")
			// begin-update_instance_network_interface

			networkInterfacePatchModel := &vpcv1.NetworkInterfacePatch{
				Name: core.StringPtr("my-network-interface-updated"),
			}
			networkInterfacePatchModelAsPatch, asPatchErr := networkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceNetworkInterfaceOptions := vpcService.NewUpdateInstanceNetworkInterfaceOptions(
				"testString",
				"testString",
				networkInterfacePatchModelAsPatch,
			)

			networkInterface, response, err := vpcService.UpdateInstanceNetworkInterface(updateInstanceNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkInterface, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterface).ToNot(BeNil())
		})
		It(`ListInstanceNetworkInterfaceFloatingIps request example`, func() {
			fmt.Println("\nListInstanceNetworkInterfaceFloatingIps() result:")
			// begin-list_instance_network_interface_floating_ips

			listInstanceNetworkInterfaceFloatingIpsOptions := vpcService.NewListInstanceNetworkInterfaceFloatingIpsOptions(
				"testString",
				"testString",
			)

			floatingIPUnpaginatedCollection, response, err := vpcService.ListInstanceNetworkInterfaceFloatingIps(listInstanceNetworkInterfaceFloatingIpsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIPUnpaginatedCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_network_interface_floating_ips

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPUnpaginatedCollection).ToNot(BeNil())
		})
		It(`GetInstanceNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nGetInstanceNetworkInterfaceFloatingIP() result:")
			// begin-get_instance_network_interface_floating_ip

			getInstanceNetworkInterfaceFloatingIPOptions := vpcService.NewGetInstanceNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			floatingIP, response, err := vpcService.GetInstanceNetworkInterfaceFloatingIP(getInstanceNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`AddInstanceNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nAddInstanceNetworkInterfaceFloatingIP() result:")
			// begin-add_instance_network_interface_floating_ip

			addInstanceNetworkInterfaceFloatingIPOptions := vpcService.NewAddInstanceNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			floatingIP, response, err := vpcService.AddInstanceNetworkInterfaceFloatingIP(addInstanceNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIP, "", "  ")
			fmt.Println(string(b))

			// end-add_instance_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
		It(`ListInstanceNetworkInterfaceIps request example`, func() {
			fmt.Println("\nListInstanceNetworkInterfaceIps() result:")
			// begin-list_instance_network_interface_ips
			listInstanceNetworkInterfaceIpsOptions := &vpcv1.ListInstanceNetworkInterfaceIpsOptions{
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewInstanceNetworkInterfaceIpsPager(listInstanceNetworkInterfaceIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_instance_network_interface_ips
		})
		It(`GetInstanceNetworkInterfaceIP request example`, func() {
			fmt.Println("\nGetInstanceNetworkInterfaceIP() result:")
			// begin-get_instance_network_interface_ip

			getInstanceNetworkInterfaceIPOptions := vpcService.NewGetInstanceNetworkInterfaceIPOptions(
				"testString",
				"testString",
				"testString",
			)

			reservedIP, response, err := vpcService.GetInstanceNetworkInterfaceIP(getInstanceNetworkInterfaceIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_network_interface_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`ListInstanceSoftwareAttachments request example`, func() {
			fmt.Println("\nListInstanceSoftwareAttachments() result:")
			// begin-list_instance_software_attachments

			listInstanceSoftwareAttachmentsOptions := vpcService.NewListInstanceSoftwareAttachmentsOptions(
				"testString",
			)

			instanceSoftwareAttachmentCollection, response, err := vpcService.ListInstanceSoftwareAttachments(listInstanceSoftwareAttachmentsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceSoftwareAttachmentCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_software_attachments

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachmentCollection).ToNot(BeNil())
		})
		It(`GetInstanceSoftwareAttachment request example`, func() {
			fmt.Println("\nGetInstanceSoftwareAttachment() result:")
			// begin-get_instance_software_attachment

			getInstanceSoftwareAttachmentOptions := vpcService.NewGetInstanceSoftwareAttachmentOptions(
				"testString",
				"testString",
			)

			instanceSoftwareAttachment, response, err := vpcService.GetInstanceSoftwareAttachment(getInstanceSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachment).ToNot(BeNil())
		})
		It(`UpdateInstanceSoftwareAttachment request example`, func() {
			fmt.Println("\nUpdateInstanceSoftwareAttachment() result:")
			// begin-update_instance_software_attachment

			instanceSoftwareAttachmentPatchModel := &vpcv1.InstanceSoftwareAttachmentPatch{}
			instanceSoftwareAttachmentPatchModelAsPatch, asPatchErr := instanceSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceSoftwareAttachmentOptions := vpcService.NewUpdateInstanceSoftwareAttachmentOptions(
				"testString",
				"testString",
				instanceSoftwareAttachmentPatchModelAsPatch,
			)

			instanceSoftwareAttachment, response, err := vpcService.UpdateInstanceSoftwareAttachment(updateInstanceSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachment).ToNot(BeNil())
		})
		It(`ListInstanceVolumeAttachments request example`, func() {
			fmt.Println("\nListInstanceVolumeAttachments() result:")
			// begin-list_instance_volume_attachments

			listInstanceVolumeAttachmentsOptions := vpcService.NewListInstanceVolumeAttachmentsOptions(
				"testString",
			)

			volumeAttachmentCollection, response, err := vpcService.ListInstanceVolumeAttachments(listInstanceVolumeAttachmentsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeAttachmentCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_instance_volume_attachments

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachmentCollection).ToNot(BeNil())
		})
		It(`CreateInstanceVolumeAttachment request example`, func() {
			fmt.Println("\nCreateInstanceVolumeAttachment() result:")
			// begin-create_instance_volume_attachment

			volumeAttachmentPrototypeVolumeModel := &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentityVolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			createInstanceVolumeAttachmentOptions := vpcService.NewCreateInstanceVolumeAttachmentOptions(
				"testString",
				volumeAttachmentPrototypeVolumeModel,
			)

			volumeAttachment, response, err := vpcService.CreateInstanceVolumeAttachment(createInstanceVolumeAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeAttachment, "", "  ")
			fmt.Println(string(b))

			// end-create_instance_volume_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volumeAttachment).ToNot(BeNil())
		})
		It(`GetInstanceVolumeAttachment request example`, func() {
			fmt.Println("\nGetInstanceVolumeAttachment() result:")
			// begin-get_instance_volume_attachment

			getInstanceVolumeAttachmentOptions := vpcService.NewGetInstanceVolumeAttachmentOptions(
				"testString",
				"testString",
			)

			volumeAttachment, response, err := vpcService.GetInstanceVolumeAttachment(getInstanceVolumeAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_instance_volume_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachment).ToNot(BeNil())
		})
		It(`UpdateInstanceVolumeAttachment request example`, func() {
			fmt.Println("\nUpdateInstanceVolumeAttachment() result:")
			// begin-update_instance_volume_attachment

			volumeAttachmentPatchModel := &vpcv1.VolumeAttachmentPatch{}
			volumeAttachmentPatchModelAsPatch, asPatchErr := volumeAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceVolumeAttachmentOptions := vpcService.NewUpdateInstanceVolumeAttachmentOptions(
				"testString",
				"testString",
				volumeAttachmentPatchModelAsPatch,
			)

			volumeAttachment, response, err := vpcService.UpdateInstanceVolumeAttachment(updateInstanceVolumeAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_instance_volume_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachment).ToNot(BeNil())
		})
		It(`ListKeys request example`, func() {
			fmt.Println("\nListKeys() result:")
			// begin-list_keys
			listKeysOptions := &vpcv1.ListKeysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewKeysPager(listKeysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Key
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_keys
		})
		It(`CreateKey request example`, func() {
			fmt.Println("\nCreateKey() result:")
			// begin-create_key

			createKeyOptions := vpcService.NewCreateKeyOptions(
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDDGe50Bxa5T5NDddrrtbx2Y4/VGbiCgXqnBsYToIUKoFSHTQl5IX3PasGnneKanhcLwWz5M5MoCRvhxTp66NKzIfAz7r+FX9rxgR+ZgcM253YAqOVeIpOU408simDZKriTlN8kYsXL7P34tsWuAJf4MgZtJAQxous/2byetpdCv8ddnT4X3ltOg9w+LqSCPYfNivqH00Eh7S1Ldz7I8aw5WOp5a+sQFP/RbwfpwHp+ny7DfeIOokcuI42tJkoBn7UsLTVpCSmXr2EDRlSWe/1M/iHNRBzaT3CK0+SwZWd2AEjePxSnWKNGIEUJDlUYp7hKhiQcgT5ZAnWU121oc5En",
			)

			key, response, err := vpcService.CreateKey(createKeyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(key, "", "  ")
			fmt.Println(string(b))

			// end-create_key

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(key).ToNot(BeNil())
		})
		It(`GetKey request example`, func() {
			fmt.Println("\nGetKey() result:")
			// begin-get_key

			getKeyOptions := vpcService.NewGetKeyOptions(
				"testString",
			)

			key, response, err := vpcService.GetKey(getKeyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(key, "", "  ")
			fmt.Println(string(b))

			// end-get_key

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(key).ToNot(BeNil())
		})
		It(`UpdateKey request example`, func() {
			fmt.Println("\nUpdateKey() result:")
			// begin-update_key

			keyPatchModel := &vpcv1.KeyPatch{}
			keyPatchModelAsPatch, asPatchErr := keyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateKeyOptions := vpcService.NewUpdateKeyOptions(
				"testString",
				keyPatchModelAsPatch,
			)

			key, response, err := vpcService.UpdateKey(updateKeyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(key, "", "  ")
			fmt.Println(string(b))

			// end-update_key

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(key).ToNot(BeNil())
		})
		It(`ListLoadBalancerProfiles request example`, func() {
			fmt.Println("\nListLoadBalancerProfiles() result:")
			// begin-list_load_balancer_profiles
			listLoadBalancerProfilesOptions := &vpcv1.ListLoadBalancerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewLoadBalancerProfilesPager(listLoadBalancerProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.LoadBalancerProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_load_balancer_profiles
		})
		It(`GetLoadBalancerProfile request example`, func() {
			fmt.Println("\nGetLoadBalancerProfile() result:")
			// begin-get_load_balancer_profile

			getLoadBalancerProfileOptions := vpcService.NewGetLoadBalancerProfileOptions(
				"network-fixed",
			)

			loadBalancerProfile, response, err := vpcService.GetLoadBalancerProfile(getLoadBalancerProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerProfile).ToNot(BeNil())
		})
		It(`ListLoadBalancers request example`, func() {
			fmt.Println("\nListLoadBalancers() result:")
			// begin-list_load_balancers
			listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewLoadBalancersPager(listLoadBalancersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.LoadBalancer
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_load_balancers
		})
		It(`CreateLoadBalancer request example`, func() {
			fmt.Println("\nCreateLoadBalancer() result:")
			// begin-create_load_balancer

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			createLoadBalancerOptions := vpcService.NewCreateLoadBalancerOptions(
				true,
				[]vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			)

			loadBalancer, response, err := vpcService.CreateLoadBalancer(createLoadBalancerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancer, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancer).ToNot(BeNil())
		})
		It(`GetLoadBalancer request example`, func() {
			fmt.Println("\nGetLoadBalancer() result:")
			// begin-get_load_balancer

			getLoadBalancerOptions := vpcService.NewGetLoadBalancerOptions(
				"testString",
			)

			loadBalancer, response, err := vpcService.GetLoadBalancer(getLoadBalancerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancer, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancer).ToNot(BeNil())
		})
		It(`UpdateLoadBalancer request example`, func() {
			fmt.Println("\nUpdateLoadBalancer() result:")
			// begin-update_load_balancer

			loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{}
			loadBalancerPatchModelAsPatch, asPatchErr := loadBalancerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerOptions := vpcService.NewUpdateLoadBalancerOptions(
				"testString",
				loadBalancerPatchModelAsPatch,
			)
			updateLoadBalancerOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			loadBalancer, response, err := vpcService.UpdateLoadBalancer(updateLoadBalancerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancer, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancer).ToNot(BeNil())
		})
		It(`GetLoadBalancerStatistics request example`, func() {
			fmt.Println("\nGetLoadBalancerStatistics() result:")
			// begin-get_load_balancer_statistics

			getLoadBalancerStatisticsOptions := vpcService.NewGetLoadBalancerStatisticsOptions(
				"testString",
			)

			loadBalancerStatistics, response, err := vpcService.GetLoadBalancerStatistics(getLoadBalancerStatisticsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerStatistics, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_statistics

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerStatistics).ToNot(BeNil())
		})
		It(`ListLoadBalancerListeners request example`, func() {
			fmt.Println("\nListLoadBalancerListeners() result:")
			// begin-list_load_balancer_listeners

			listLoadBalancerListenersOptions := vpcService.NewListLoadBalancerListenersOptions(
				"testString",
			)

			loadBalancerListenerCollection, response, err := vpcService.ListLoadBalancerListeners(listLoadBalancerListenersOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_load_balancer_listeners

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerCollection).ToNot(BeNil())
		})
		It(`CreateLoadBalancerListener request example`, func() {
			fmt.Println("\nCreateLoadBalancerListener() result:")
			// begin-create_load_balancer_listener

			createLoadBalancerListenerOptions := vpcService.NewCreateLoadBalancerListenerOptions(
				"testString",
				"http",
			)

			loadBalancerListener, response, err := vpcService.CreateLoadBalancerListener(createLoadBalancerListenerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListener, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer_listener

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
		It(`GetLoadBalancerListener request example`, func() {
			fmt.Println("\nGetLoadBalancerListener() result:")
			// begin-get_load_balancer_listener

			getLoadBalancerListenerOptions := vpcService.NewGetLoadBalancerListenerOptions(
				"testString",
				"testString",
			)

			loadBalancerListener, response, err := vpcService.GetLoadBalancerListener(getLoadBalancerListenerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListener, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_listener

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
		It(`UpdateLoadBalancerListener request example`, func() {
			fmt.Println("\nUpdateLoadBalancerListener() result:")
			// begin-update_load_balancer_listener

			loadBalancerListenerPatchModel := &vpcv1.LoadBalancerListenerPatch{}
			loadBalancerListenerPatchModelAsPatch, asPatchErr := loadBalancerListenerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerOptions := vpcService.NewUpdateLoadBalancerListenerOptions(
				"testString",
				"testString",
				loadBalancerListenerPatchModelAsPatch,
			)

			loadBalancerListener, response, err := vpcService.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListener, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer_listener

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
		It(`ListLoadBalancerListenerPolicies request example`, func() {
			fmt.Println("\nListLoadBalancerListenerPolicies() result:")
			// begin-list_load_balancer_listener_policies

			listLoadBalancerListenerPoliciesOptions := vpcService.NewListLoadBalancerListenerPoliciesOptions(
				"testString",
				"testString",
			)

			loadBalancerListenerPolicyCollection, response, err := vpcService.ListLoadBalancerListenerPolicies(listLoadBalancerListenerPoliciesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicyCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_load_balancer_listener_policies

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyCollection).ToNot(BeNil())
		})
		It(`CreateLoadBalancerListenerPolicy request example`, func() {
			fmt.Println("\nCreateLoadBalancerListenerPolicy() result:")
			// begin-create_load_balancer_listener_policy

			createLoadBalancerListenerPolicyOptions := vpcService.NewCreateLoadBalancerListenerPolicyOptions(
				"testString",
				"testString",
				"forward_to_listener",
				int64(5),
			)

			loadBalancerListenerPolicy, response, err := vpcService.CreateLoadBalancerListenerPolicy(createLoadBalancerListenerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer_listener_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
		It(`GetLoadBalancerListenerPolicy request example`, func() {
			fmt.Println("\nGetLoadBalancerListenerPolicy() result:")
			// begin-get_load_balancer_listener_policy

			getLoadBalancerListenerPolicyOptions := vpcService.NewGetLoadBalancerListenerPolicyOptions(
				"testString",
				"testString",
				"testString",
			)

			loadBalancerListenerPolicy, response, err := vpcService.GetLoadBalancerListenerPolicy(getLoadBalancerListenerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_listener_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
		It(`UpdateLoadBalancerListenerPolicy request example`, func() {
			fmt.Println("\nUpdateLoadBalancerListenerPolicy() result:")
			// begin-update_load_balancer_listener_policy

			loadBalancerListenerPolicyPatchModel := &vpcv1.LoadBalancerListenerPolicyPatch{
				Name: core.StringPtr("my-load-balancer-listener-policy-updated"),
			}
			loadBalancerListenerPolicyPatchModelAsPatch, asPatchErr := loadBalancerListenerPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerPolicyOptions := vpcService.NewUpdateLoadBalancerListenerPolicyOptions(
				"testString",
				"testString",
				"testString",
				loadBalancerListenerPolicyPatchModelAsPatch,
			)

			loadBalancerListenerPolicy, response, err := vpcService.UpdateLoadBalancerListenerPolicy(updateLoadBalancerListenerPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer_listener_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
		It(`ListLoadBalancerListenerPolicyRules request example`, func() {
			fmt.Println("\nListLoadBalancerListenerPolicyRules() result:")
			// begin-list_load_balancer_listener_policy_rules

			listLoadBalancerListenerPolicyRulesOptions := vpcService.NewListLoadBalancerListenerPolicyRulesOptions(
				"testString",
				"testString",
				"testString",
			)

			loadBalancerListenerPolicyRuleCollection, response, err := vpcService.ListLoadBalancerListenerPolicyRules(listLoadBalancerListenerPolicyRulesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicyRuleCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_load_balancer_listener_policy_rules

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRuleCollection).ToNot(BeNil())
		})
		It(`CreateLoadBalancerListenerPolicyRule request example`, func() {
			fmt.Println("\nCreateLoadBalancerListenerPolicyRule() result:")
			// begin-create_load_balancer_listener_policy_rule

			createLoadBalancerListenerPolicyRuleOptions := vpcService.NewCreateLoadBalancerListenerPolicyRuleOptions(
				"testString",
				"testString",
				"testString",
				"contains",
				"body",
				"my-value",
			)

			loadBalancerListenerPolicyRule, response, err := vpcService.CreateLoadBalancerListenerPolicyRule(createLoadBalancerListenerPolicyRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicyRule, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer_listener_policy_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
		It(`GetLoadBalancerListenerPolicyRule request example`, func() {
			fmt.Println("\nGetLoadBalancerListenerPolicyRule() result:")
			// begin-get_load_balancer_listener_policy_rule

			getLoadBalancerListenerPolicyRuleOptions := vpcService.NewGetLoadBalancerListenerPolicyRuleOptions(
				"testString",
				"testString",
				"testString",
				"testString",
			)

			loadBalancerListenerPolicyRule, response, err := vpcService.GetLoadBalancerListenerPolicyRule(getLoadBalancerListenerPolicyRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicyRule, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_listener_policy_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
		It(`UpdateLoadBalancerListenerPolicyRule request example`, func() {
			fmt.Println("\nUpdateLoadBalancerListenerPolicyRule() result:")
			// begin-update_load_balancer_listener_policy_rule

			loadBalancerListenerPolicyRulePatchModel := &vpcv1.LoadBalancerListenerPolicyRulePatch{
				Value: core.StringPtr("my-example-hostname-updated"),
			}
			loadBalancerListenerPolicyRulePatchModelAsPatch, asPatchErr := loadBalancerListenerPolicyRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerPolicyRuleOptions := vpcService.NewUpdateLoadBalancerListenerPolicyRuleOptions(
				"testString",
				"testString",
				"testString",
				"testString",
				loadBalancerListenerPolicyRulePatchModelAsPatch,
			)

			loadBalancerListenerPolicyRule, response, err := vpcService.UpdateLoadBalancerListenerPolicyRule(updateLoadBalancerListenerPolicyRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerListenerPolicyRule, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer_listener_policy_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
		It(`ListLoadBalancerPools request example`, func() {
			fmt.Println("\nListLoadBalancerPools() result:")
			// begin-list_load_balancer_pools

			listLoadBalancerPoolsOptions := vpcService.NewListLoadBalancerPoolsOptions(
				"testString",
			)

			loadBalancerPoolCollection, response, err := vpcService.ListLoadBalancerPools(listLoadBalancerPoolsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_load_balancer_pools

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolCollection).ToNot(BeNil())
		})
		It(`CreateLoadBalancerPool request example`, func() {
			fmt.Println("\nCreateLoadBalancerPool() result:")
			// begin-create_load_balancer_pool

			loadBalancerPoolPrototypeHealthMonitorModel := &vpcv1.LoadBalancerPoolPrototypeHealthMonitorLoadBalancerPoolHealthMonitorPrototypeLoadBalancerPoolHealthMonitorTypeHttphttpsPrototype{
				Delay:      core.Int64Ptr(int64(5)),
				MaxRetries: core.Int64Ptr(int64(2)),
				Timeout:    core.Int64Ptr(int64(2)),
				Type:       core.StringPtr("http"),
			}

			createLoadBalancerPoolOptions := vpcService.NewCreateLoadBalancerPoolOptions(
				"testString",
				"least_connections",
				loadBalancerPoolPrototypeHealthMonitorModel,
				"http",
			)

			loadBalancerPool, response, err := vpcService.CreateLoadBalancerPool(createLoadBalancerPoolOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPool, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer_pool

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
		It(`GetLoadBalancerPool request example`, func() {
			fmt.Println("\nGetLoadBalancerPool() result:")
			// begin-get_load_balancer_pool

			getLoadBalancerPoolOptions := vpcService.NewGetLoadBalancerPoolOptions(
				"testString",
				"testString",
			)

			loadBalancerPool, response, err := vpcService.GetLoadBalancerPool(getLoadBalancerPoolOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPool, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_pool

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
		It(`UpdateLoadBalancerPool request example`, func() {
			fmt.Println("\nUpdateLoadBalancerPool() result:")
			// begin-update_load_balancer_pool

			loadBalancerPoolPatchModel := &vpcv1.LoadBalancerPoolPatch{
				Algorithm: core.StringPtr("weighted_round_robin"),
			}
			loadBalancerPoolPatchModelAsPatch, asPatchErr := loadBalancerPoolPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerPoolOptions := vpcService.NewUpdateLoadBalancerPoolOptions(
				"testString",
				"testString",
				loadBalancerPoolPatchModelAsPatch,
			)
			updateLoadBalancerPoolOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			loadBalancerPool, response, err := vpcService.UpdateLoadBalancerPool(updateLoadBalancerPoolOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPool, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer_pool

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
		It(`ListLoadBalancerPoolMembers request example`, func() {
			fmt.Println("\nListLoadBalancerPoolMembers() result:")
			// begin-list_load_balancer_pool_members

			listLoadBalancerPoolMembersOptions := vpcService.NewListLoadBalancerPoolMembersOptions(
				"testString",
				"testString",
			)

			loadBalancerPoolMemberCollection, response, err := vpcService.ListLoadBalancerPoolMembers(listLoadBalancerPoolMembersOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolMemberCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_load_balancer_pool_members

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMemberCollection).ToNot(BeNil())
		})
		It(`CreateLoadBalancerPoolMember request example`, func() {
			fmt.Println("\nCreateLoadBalancerPoolMember() result:")
			// begin-create_load_balancer_pool_member

			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			createLoadBalancerPoolMemberOptions := vpcService.NewCreateLoadBalancerPoolMemberOptions(
				"testString",
				"testString",
				int64(80),
				loadBalancerPoolMemberTargetPrototypeModel,
			)

			loadBalancerPoolMember, response, err := vpcService.CreateLoadBalancerPoolMember(createLoadBalancerPoolMemberOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolMember, "", "  ")
			fmt.Println(string(b))

			// end-create_load_balancer_pool_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
		It(`ReplaceLoadBalancerPoolMembers request example`, func() {
			fmt.Println("\nReplaceLoadBalancerPoolMembers() result:")
			// begin-replace_load_balancer_pool_members

			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			loadBalancerPoolMemberPrototypeModel := &vpcv1.LoadBalancerPoolMemberPrototype{
				Port:   core.Int64Ptr(int64(80)),
				Target: loadBalancerPoolMemberTargetPrototypeModel,
			}

			replaceLoadBalancerPoolMembersOptions := vpcService.NewReplaceLoadBalancerPoolMembersOptions(
				"testString",
				"testString",
				[]vpcv1.LoadBalancerPoolMemberPrototype{*loadBalancerPoolMemberPrototypeModel},
			)

			loadBalancerPoolMemberCollection, response, err := vpcService.ReplaceLoadBalancerPoolMembers(replaceLoadBalancerPoolMembersOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolMemberCollection, "", "  ")
			fmt.Println(string(b))

			// end-replace_load_balancer_pool_members

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMemberCollection).ToNot(BeNil())
		})
		It(`GetLoadBalancerPoolMember request example`, func() {
			fmt.Println("\nGetLoadBalancerPoolMember() result:")
			// begin-get_load_balancer_pool_member

			getLoadBalancerPoolMemberOptions := vpcService.NewGetLoadBalancerPoolMemberOptions(
				"testString",
				"testString",
				"testString",
			)

			loadBalancerPoolMember, response, err := vpcService.GetLoadBalancerPoolMember(getLoadBalancerPoolMemberOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolMember, "", "  ")
			fmt.Println(string(b))

			// end-get_load_balancer_pool_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
		It(`UpdateLoadBalancerPoolMember request example`, func() {
			fmt.Println("\nUpdateLoadBalancerPoolMember() result:")
			// begin-update_load_balancer_pool_member

			loadBalancerPoolMemberPatchModel := &vpcv1.LoadBalancerPoolMemberPatch{
				Weight: core.Int64Ptr(int64(99)),
			}
			loadBalancerPoolMemberPatchModelAsPatch, asPatchErr := loadBalancerPoolMemberPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerPoolMemberOptions := vpcService.NewUpdateLoadBalancerPoolMemberOptions(
				"testString",
				"testString",
				"testString",
				loadBalancerPoolMemberPatchModelAsPatch,
			)

			loadBalancerPoolMember, response, err := vpcService.UpdateLoadBalancerPoolMember(updateLoadBalancerPoolMemberOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(loadBalancerPoolMember, "", "  ")
			fmt.Println(string(b))

			// end-update_load_balancer_pool_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
		It(`ListNetworkAcls request example`, func() {
			fmt.Println("\nListNetworkAcls() result:")
			// begin-list_network_acls
			listNetworkAclsOptions := &vpcv1.ListNetworkAclsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewNetworkAclsPager(listNetworkAclsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.NetworkACL
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_network_acls
		})
		It(`CreateNetworkACL request example`, func() {
			fmt.Println("\nCreateNetworkACL() result:")
			// begin-create_network_acl

			networkACLPrototypeVPCModel := &vpcv1.NetworkACLPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			networkACLPrototypeModel := &vpcv1.NetworkACLPrototypeNetworkACLByRules{
				VPC: networkACLPrototypeVPCModel,
			}

			createNetworkACLOptions := vpcService.NewCreateNetworkACLOptions(
				networkACLPrototypeModel,
			)

			networkACL, response, err := vpcService.CreateNetworkACL(createNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACL, "", "  ")
			fmt.Println(string(b))

			// end-create_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACL).ToNot(BeNil())
		})
		It(`GetNetworkACL request example`, func() {
			fmt.Println("\nGetNetworkACL() result:")
			// begin-get_network_acl

			getNetworkACLOptions := vpcService.NewGetNetworkACLOptions(
				"testString",
			)

			networkACL, response, err := vpcService.GetNetworkACL(getNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACL, "", "  ")
			fmt.Println(string(b))

			// end-get_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
		It(`UpdateNetworkACL request example`, func() {
			fmt.Println("\nUpdateNetworkACL() result:")
			// begin-update_network_acl

			networkACLPatchModel := &vpcv1.NetworkACLPatch{
				Name: core.StringPtr("my-network-acl-updated"),
			}
			networkACLPatchModelAsPatch, asPatchErr := networkACLPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateNetworkACLOptions := vpcService.NewUpdateNetworkACLOptions(
				"testString",
				networkACLPatchModelAsPatch,
			)

			networkACL, response, err := vpcService.UpdateNetworkACL(updateNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACL, "", "  ")
			fmt.Println(string(b))

			// end-update_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
		It(`ListNetworkACLRules request example`, func() {
			fmt.Println("\nListNetworkACLRules() result:")
			// begin-list_network_acl_rules
			listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
				NetworkACLID: core.StringPtr("testString"),
				Direction:    core.StringPtr("inbound"),
				Limit:        core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewNetworkACLRulesPager(listNetworkACLRulesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.NetworkACLRuleItemIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_network_acl_rules
		})
		It(`CreateNetworkACLRule request example`, func() {
			fmt.Println("\nCreateNetworkACLRule() result:")
			// begin-create_network_acl_rule

			networkACLRulePrototypeModel := &vpcv1.NetworkACLRulePrototypeNetworkACLRuleProtocolAnyPrototype{
				Action:      core.StringPtr("allow"),
				Destination: core.StringPtr("192.168.3.2/32"),
				Direction:   core.StringPtr("inbound"),
				Source:      core.StringPtr("192.168.3.2/32"),
				Protocol:    core.StringPtr("any"),
			}

			createNetworkACLRuleOptions := vpcService.NewCreateNetworkACLRuleOptions(
				"testString",
				networkACLRulePrototypeModel,
			)

			networkACLRule, response, err := vpcService.CreateNetworkACLRule(createNetworkACLRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACLRule, "", "  ")
			fmt.Println(string(b))

			// end-create_network_acl_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACLRule).ToNot(BeNil())
		})
		It(`GetNetworkACLRule request example`, func() {
			fmt.Println("\nGetNetworkACLRule() result:")
			// begin-get_network_acl_rule

			getNetworkACLRuleOptions := vpcService.NewGetNetworkACLRuleOptions(
				"testString",
				"testString",
			)

			networkACLRule, response, err := vpcService.GetNetworkACLRule(getNetworkACLRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACLRule, "", "  ")
			fmt.Println(string(b))

			// end-get_network_acl_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACLRule).ToNot(BeNil())
		})
		It(`UpdateNetworkACLRule request example`, func() {
			fmt.Println("\nUpdateNetworkACLRule() result:")
			// begin-update_network_acl_rule

			networkACLRulePatchModel := &vpcv1.NetworkACLRulePatch{
				Source: core.StringPtr("10.0.0.0/0"),
			}
			networkACLRulePatchModelAsPatch, asPatchErr := networkACLRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateNetworkACLRuleOptions := vpcService.NewUpdateNetworkACLRuleOptions(
				"testString",
				"testString",
				networkACLRulePatchModelAsPatch,
			)

			networkACLRule, response, err := vpcService.UpdateNetworkACLRule(updateNetworkACLRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACLRule, "", "  ")
			fmt.Println(string(b))

			// end-update_network_acl_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACLRule).ToNot(BeNil())
		})
		It(`ListPlacementGroups request example`, func() {
			fmt.Println("\nListPlacementGroups() result:")
			// begin-list_placement_groups
			listPlacementGroupsOptions := &vpcv1.ListPlacementGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewPlacementGroupsPager(listPlacementGroupsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PlacementGroup
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_placement_groups
		})
		It(`CreatePlacementGroup request example`, func() {
			fmt.Println("\nCreatePlacementGroup() result:")
			// begin-create_placement_group

			createPlacementGroupOptions := vpcService.NewCreatePlacementGroupOptions(
				"host_spread",
			)

			placementGroup, response, err := vpcService.CreatePlacementGroup(createPlacementGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(placementGroup, "", "  ")
			fmt.Println(string(b))

			// end-create_placement_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(placementGroup).ToNot(BeNil())
		})
		It(`GetPlacementGroup request example`, func() {
			fmt.Println("\nGetPlacementGroup() result:")
			// begin-get_placement_group

			getPlacementGroupOptions := vpcService.NewGetPlacementGroupOptions(
				"testString",
			)

			placementGroup, response, err := vpcService.GetPlacementGroup(getPlacementGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(placementGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_placement_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(placementGroup).ToNot(BeNil())
		})
		It(`UpdatePlacementGroup request example`, func() {
			fmt.Println("\nUpdatePlacementGroup() result:")
			// begin-update_placement_group

			placementGroupPatchModel := &vpcv1.PlacementGroupPatch{
				Name: core.StringPtr("my-placement-group-updated"),
			}
			placementGroupPatchModelAsPatch, asPatchErr := placementGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePlacementGroupOptions := vpcService.NewUpdatePlacementGroupOptions(
				"testString",
				placementGroupPatchModelAsPatch,
			)

			placementGroup, response, err := vpcService.UpdatePlacementGroup(updatePlacementGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(placementGroup, "", "  ")
			fmt.Println(string(b))

			// end-update_placement_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(placementGroup).ToNot(BeNil())
		})
		It(`ListPrivatePathServiceGateways request example`, func() {
			fmt.Println("\nListPrivatePathServiceGateways() result:")
			// begin-list_private_path_service_gateways
			listPrivatePathServiceGatewaysOptions := &vpcv1.ListPrivatePathServiceGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewPrivatePathServiceGatewaysPager(listPrivatePathServiceGatewaysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PrivatePathServiceGateway
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_private_path_service_gateways
		})
		It(`CreatePrivatePathServiceGateway request example`, func() {
			fmt.Println("\nCreatePrivatePathServiceGateway() result:")
			// begin-create_private_path_service_gateway

			privatePathServiceGatewayPrototypeLoadBalancerModel := &vpcv1.PrivatePathServiceGatewayPrototypeLoadBalancerLoadBalancerIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
			}

			createPrivatePathServiceGatewayOptions := vpcService.NewCreatePrivatePathServiceGatewayOptions(
				privatePathServiceGatewayPrototypeLoadBalancerModel,
				[]string{"*.example.com"},
			)

			privatePathServiceGateway, response, err := vpcService.CreatePrivatePathServiceGateway(createPrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGateway, "", "  ")
			fmt.Println(string(b))

			// end-create_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
		It(`GetPrivatePathServiceGateway request example`, func() {
			fmt.Println("\nGetPrivatePathServiceGateway() result:")
			// begin-get_private_path_service_gateway

			getPrivatePathServiceGatewayOptions := vpcService.NewGetPrivatePathServiceGatewayOptions(
				"testString",
			)

			privatePathServiceGateway, response, err := vpcService.GetPrivatePathServiceGateway(getPrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGateway, "", "  ")
			fmt.Println(string(b))

			// end-get_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
		It(`UpdatePrivatePathServiceGateway request example`, func() {
			fmt.Println("\nUpdatePrivatePathServiceGateway() result:")
			// begin-update_private_path_service_gateway

			privatePathServiceGatewayPatchModel := &vpcv1.PrivatePathServiceGatewayPatch{
				ZonalAffinity: core.BoolPtr(true),
			}
			privatePathServiceGatewayPatchModelAsPatch, asPatchErr := privatePathServiceGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePrivatePathServiceGatewayOptions := vpcService.NewUpdatePrivatePathServiceGatewayOptions(
				"testString",
				privatePathServiceGatewayPatchModelAsPatch,
			)

			privatePathServiceGateway, response, err := vpcService.UpdatePrivatePathServiceGateway(updatePrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGateway, "", "  ")
			fmt.Println(string(b))

			// end-update_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
		It(`ListPrivatePathServiceGatewayAccountPolicies request example`, func() {
			fmt.Println("\nListPrivatePathServiceGatewayAccountPolicies() result:")
			// begin-list_private_path_service_gateway_account_policies
			listPrivatePathServiceGatewayAccountPoliciesOptions := &vpcv1.ListPrivatePathServiceGatewayAccountPoliciesOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewPrivatePathServiceGatewayAccountPoliciesPager(listPrivatePathServiceGatewayAccountPoliciesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PrivatePathServiceGatewayAccountPolicy
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_private_path_service_gateway_account_policies
		})
		It(`CreatePrivatePathServiceGatewayAccountPolicy request example`, func() {
			fmt.Println("\nCreatePrivatePathServiceGatewayAccountPolicy() result:")
			// begin-create_private_path_service_gateway_account_policy

			privatePathServiceGatewayAccountPolicyPrototypeAccountModel := &vpcv1.PrivatePathServiceGatewayAccountPolicyPrototypeAccountAccountIdentityByID{
				ID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
			}

			createPrivatePathServiceGatewayAccountPolicyOptions := vpcService.NewCreatePrivatePathServiceGatewayAccountPolicyOptions(
				"testString",
				"deny",
				privatePathServiceGatewayAccountPolicyPrototypeAccountModel,
			)

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.CreatePrivatePathServiceGatewayAccountPolicy(createPrivatePathServiceGatewayAccountPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGatewayAccountPolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_private_path_service_gateway_account_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
		It(`GetPrivatePathServiceGatewayAccountPolicy request example`, func() {
			fmt.Println("\nGetPrivatePathServiceGatewayAccountPolicy() result:")
			// begin-get_private_path_service_gateway_account_policy

			getPrivatePathServiceGatewayAccountPolicyOptions := vpcService.NewGetPrivatePathServiceGatewayAccountPolicyOptions(
				"testString",
				"testString",
			)

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.GetPrivatePathServiceGatewayAccountPolicy(getPrivatePathServiceGatewayAccountPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGatewayAccountPolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_private_path_service_gateway_account_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
		It(`UpdatePrivatePathServiceGatewayAccountPolicy request example`, func() {
			fmt.Println("\nUpdatePrivatePathServiceGatewayAccountPolicy() result:")
			// begin-update_private_path_service_gateway_account_policy

			privatePathServiceGatewayAccountPolicyPatchModel := &vpcv1.PrivatePathServiceGatewayAccountPolicyPatch{
				AccessPolicy: core.StringPtr("review"),
			}
			privatePathServiceGatewayAccountPolicyPatchModelAsPatch, asPatchErr := privatePathServiceGatewayAccountPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePrivatePathServiceGatewayAccountPolicyOptions := vpcService.NewUpdatePrivatePathServiceGatewayAccountPolicyOptions(
				"testString",
				"testString",
				privatePathServiceGatewayAccountPolicyPatchModelAsPatch,
			)

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.UpdatePrivatePathServiceGatewayAccountPolicy(updatePrivatePathServiceGatewayAccountPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGatewayAccountPolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_private_path_service_gateway_account_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
		It(`ListPrivatePathServiceGatewayEndpointGatewayBindings request example`, func() {
			fmt.Println("\nListPrivatePathServiceGatewayEndpointGatewayBindings() result:")
			// begin-list_private_path_service_gateway_endpoint_gateway_bindings
			listPrivatePathServiceGatewayEndpointGatewayBindingsOptions := &vpcv1.ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
				Status:                      core.StringPtr("abandoned"),
			}

			pager, err := vpcService.NewPrivatePathServiceGatewayEndpointGatewayBindingsPager(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PrivatePathServiceGatewayEndpointGatewayBinding
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_private_path_service_gateway_endpoint_gateway_bindings
		})
		It(`GetPrivatePathServiceGatewayEndpointGatewayBinding request example`, func() {
			fmt.Println("\nGetPrivatePathServiceGatewayEndpointGatewayBinding() result:")
			// begin-get_private_path_service_gateway_endpoint_gateway_binding

			getPrivatePathServiceGatewayEndpointGatewayBindingOptions := vpcService.NewGetPrivatePathServiceGatewayEndpointGatewayBindingOptions(
				"testString",
				"testString",
			)

			privatePathServiceGatewayEndpointGatewayBinding, response, err := vpcService.GetPrivatePathServiceGatewayEndpointGatewayBinding(getPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(privatePathServiceGatewayEndpointGatewayBinding, "", "  ")
			fmt.Println(string(b))

			// end-get_private_path_service_gateway_endpoint_gateway_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayEndpointGatewayBinding).ToNot(BeNil())
		})
		It(`DenyPrivatePathServiceGatewayEndpointGatewayBinding request example`, func() {
			// begin-deny_private_path_service_gateway_endpoint_gateway_binding

			denyPrivatePathServiceGatewayEndpointGatewayBindingOptions := vpcService.NewDenyPrivatePathServiceGatewayEndpointGatewayBindingOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DenyPrivatePathServiceGatewayEndpointGatewayBinding(denyPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DenyPrivatePathServiceGatewayEndpointGatewayBinding(): %d\n", response.StatusCode)
			}

			// end-deny_private_path_service_gateway_endpoint_gateway_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`PermitPrivatePathServiceGatewayEndpointGatewayBinding request example`, func() {
			// begin-permit_private_path_service_gateway_endpoint_gateway_binding

			permitPrivatePathServiceGatewayEndpointGatewayBindingOptions := vpcService.NewPermitPrivatePathServiceGatewayEndpointGatewayBindingOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.PermitPrivatePathServiceGatewayEndpointGatewayBinding(permitPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from PermitPrivatePathServiceGatewayEndpointGatewayBinding(): %d\n", response.StatusCode)
			}

			// end-permit_private_path_service_gateway_endpoint_gateway_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`PublishPrivatePathServiceGateway request example`, func() {
			// begin-publish_private_path_service_gateway

			publishPrivatePathServiceGatewayOptions := vpcService.NewPublishPrivatePathServiceGatewayOptions(
				"testString",
			)

			response, err := vpcService.PublishPrivatePathServiceGateway(publishPrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from PublishPrivatePathServiceGateway(): %d\n", response.StatusCode)
			}

			// end-publish_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RevokeAccountForPrivatePathServiceGateway request example`, func() {
			// begin-revoke_account_for_private_path_service_gateway

			privatePathServiceGatewayRevokeAccountOptionsAccountModel := &vpcv1.PrivatePathServiceGatewayRevokeAccountOptionsAccountAccountIdentityByID{
				ID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
			}

			revokeAccountForPrivatePathServiceGatewayOptions := vpcService.NewRevokeAccountForPrivatePathServiceGatewayOptions(
				"testString",
				privatePathServiceGatewayRevokeAccountOptionsAccountModel,
			)

			response, err := vpcService.RevokeAccountForPrivatePathServiceGateway(revokeAccountForPrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RevokeAccountForPrivatePathServiceGateway(): %d\n", response.StatusCode)
			}

			// end-revoke_account_for_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`UnpublishPrivatePathServiceGateway request example`, func() {
			// begin-unpublish_private_path_service_gateway

			unpublishPrivatePathServiceGatewayOptions := vpcService.NewUnpublishPrivatePathServiceGatewayOptions(
				"testString",
			)

			response, err := vpcService.UnpublishPrivatePathServiceGateway(unpublishPrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from UnpublishPrivatePathServiceGateway(): %d\n", response.StatusCode)
			}

			// end-unpublish_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`ListPublicAddressRanges request example`, func() {
			fmt.Println("\nListPublicAddressRanges() result:")
			// begin-list_public_address_ranges
			listPublicAddressRangesOptions := &vpcv1.ListPublicAddressRangesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewPublicAddressRangesPager(listPublicAddressRangesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PublicAddressRange
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_public_address_ranges
		})
		It(`CreatePublicAddressRange request example`, func() {
			fmt.Println("\nCreatePublicAddressRange() result:")
			// begin-create_public_address_range

			createPublicAddressRangeOptions := vpcService.NewCreatePublicAddressRangeOptions(
				int64(8),
			)

			publicAddressRange, response, err := vpcService.CreatePublicAddressRange(createPublicAddressRangeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicAddressRange, "", "  ")
			fmt.Println(string(b))

			// end-create_public_address_range

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicAddressRange).ToNot(BeNil())
		})
		It(`GetPublicAddressRange request example`, func() {
			fmt.Println("\nGetPublicAddressRange() result:")
			// begin-get_public_address_range

			getPublicAddressRangeOptions := vpcService.NewGetPublicAddressRangeOptions(
				"testString",
			)

			publicAddressRange, response, err := vpcService.GetPublicAddressRange(getPublicAddressRangeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicAddressRange, "", "  ")
			fmt.Println(string(b))

			// end-get_public_address_range

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicAddressRange).ToNot(BeNil())
		})
		It(`UpdatePublicAddressRange request example`, func() {
			fmt.Println("\nUpdatePublicAddressRange() result:")
			// begin-update_public_address_range

			publicAddressRangePatchModel := &vpcv1.PublicAddressRangePatch{
				Name: core.StringPtr("my-public-address-range-updated"),
			}
			publicAddressRangePatchModelAsPatch, asPatchErr := publicAddressRangePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePublicAddressRangeOptions := vpcService.NewUpdatePublicAddressRangeOptions(
				"testString",
				publicAddressRangePatchModelAsPatch,
			)

			publicAddressRange, response, err := vpcService.UpdatePublicAddressRange(updatePublicAddressRangeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicAddressRange, "", "  ")
			fmt.Println(string(b))

			// end-update_public_address_range

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicAddressRange).ToNot(BeNil())
		})
		It(`ListPublicGateways request example`, func() {
			fmt.Println("\nListPublicGateways() result:")
			// begin-list_public_gateways
			listPublicGatewaysOptions := &vpcv1.ListPublicGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewPublicGatewaysPager(listPublicGatewaysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.PublicGateway
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_public_gateways
		})
		It(`CreatePublicGateway request example`, func() {
			fmt.Println("\nCreatePublicGateway() result:")
			// begin-create_public_gateway

			publicGatewayPrototypeVPCModel := &vpcv1.PublicGatewayPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			publicGatewayPrototypeZoneModel := &vpcv1.PublicGatewayPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createPublicGatewayOptions := vpcService.NewCreatePublicGatewayOptions(
				publicGatewayPrototypeVPCModel,
				publicGatewayPrototypeZoneModel,
			)

			publicGateway, response, err := vpcService.CreatePublicGateway(createPublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicGateway, "", "  ")
			fmt.Println(string(b))

			// end-create_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicGateway).ToNot(BeNil())
		})
		It(`GetPublicGateway request example`, func() {
			fmt.Println("\nGetPublicGateway() result:")
			// begin-get_public_gateway

			getPublicGatewayOptions := vpcService.NewGetPublicGatewayOptions(
				"testString",
			)

			publicGateway, response, err := vpcService.GetPublicGateway(getPublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicGateway, "", "  ")
			fmt.Println(string(b))

			// end-get_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
		It(`UpdatePublicGateway request example`, func() {
			fmt.Println("\nUpdatePublicGateway() result:")
			// begin-update_public_gateway

			publicGatewayPatchModel := &vpcv1.PublicGatewayPatch{
				Name: core.StringPtr("my-public-gateway-updated"),
			}
			publicGatewayPatchModelAsPatch, asPatchErr := publicGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePublicGatewayOptions := vpcService.NewUpdatePublicGatewayOptions(
				"testString",
				publicGatewayPatchModelAsPatch,
			)

			publicGateway, response, err := vpcService.UpdatePublicGateway(updatePublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicGateway, "", "  ")
			fmt.Println(string(b))

			// end-update_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
		It(`ListReservations request example`, func() {
			fmt.Println("\nListReservations() result:")
			// begin-list_reservations
			listReservationsOptions := &vpcv1.ListReservationsOptions{
				AffinityPolicy:      core.StringPtr("automatic"),
				Limit:               core.Int64Ptr(int64(10)),
				Name:                core.StringPtr("my-name"),
				ProfileResourceType: core.StringPtr("testString"),
				ResourceGroupID:     core.StringPtr("testString"),
				ZoneName:            core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewReservationsPager(listReservationsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Reservation
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_reservations
		})
		It(`CreateReservation request example`, func() {
			fmt.Println("\nCreateReservation() result:")
			// begin-create_reservation

			reservationCapacityPrototypeModel := &vpcv1.ReservationCapacityPrototype{
				Total: core.Int64Ptr(int64(10)),
			}

			reservationPrototypeCommittedUseModel := &vpcv1.ReservationPrototypeCommittedUse{
				Term: core.StringPtr("one_year"),
			}

			reservationProfilePrototypeModel := &vpcv1.ReservationProfilePrototype{
				Name:         core.StringPtr("bx2-4x16"),
				ResourceType: core.StringPtr("bare_metal_server_profile"),
			}

			reservationPrototypeZoneModel := &vpcv1.ReservationPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createReservationOptions := vpcService.NewCreateReservationOptions(
				reservationCapacityPrototypeModel,
				reservationPrototypeCommittedUseModel,
				reservationProfilePrototypeModel,
				reservationPrototypeZoneModel,
			)

			reservation, response, err := vpcService.CreateReservation(createReservationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservation, "", "  ")
			fmt.Println(string(b))

			// end-create_reservation

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservation).ToNot(BeNil())
		})
		It(`GetReservation request example`, func() {
			fmt.Println("\nGetReservation() result:")
			// begin-get_reservation

			getReservationOptions := vpcService.NewGetReservationOptions(
				"testString",
			)

			reservation, response, err := vpcService.GetReservation(getReservationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservation, "", "  ")
			fmt.Println(string(b))

			// end-get_reservation

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservation).ToNot(BeNil())
		})
		It(`UpdateReservation request example`, func() {
			fmt.Println("\nUpdateReservation() result:")
			// begin-update_reservation

			reservationPatchModel := &vpcv1.ReservationPatch{}
			reservationPatchModelAsPatch, asPatchErr := reservationPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateReservationOptions := vpcService.NewUpdateReservationOptions(
				"testString",
				reservationPatchModelAsPatch,
			)

			reservation, response, err := vpcService.UpdateReservation(updateReservationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservation, "", "  ")
			fmt.Println(string(b))

			// end-update_reservation

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservation).ToNot(BeNil())
		})
		It(`ActivateReservation request example`, func() {
			// begin-activate_reservation

			activateReservationOptions := vpcService.NewActivateReservationOptions(
				"testString",
			)

			response, err := vpcService.ActivateReservation(activateReservationOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from ActivateReservation(): %d\n", response.StatusCode)
			}

			// end-activate_reservation

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`ListSecurityGroups request example`, func() {
			fmt.Println("\nListSecurityGroups() result:")
			// begin-list_security_groups
			listSecurityGroupsOptions := &vpcv1.ListSecurityGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewSecurityGroupsPager(listSecurityGroupsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.SecurityGroup
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_security_groups
		})
		It(`CreateSecurityGroup request example`, func() {
			fmt.Println("\nCreateSecurityGroup() result:")
			// begin-create_security_group

			securityGroupPrototypeVPCModel := &vpcv1.SecurityGroupPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			createSecurityGroupOptions := vpcService.NewCreateSecurityGroupOptions(
				securityGroupPrototypeVPCModel,
			)

			securityGroup, response, err := vpcService.CreateSecurityGroup(createSecurityGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroup, "", "  ")
			fmt.Println(string(b))

			// end-create_security_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(securityGroup).ToNot(BeNil())
		})
		It(`GetSecurityGroup request example`, func() {
			fmt.Println("\nGetSecurityGroup() result:")
			// begin-get_security_group

			getSecurityGroupOptions := vpcService.NewGetSecurityGroupOptions(
				"testString",
			)

			securityGroup, response, err := vpcService.GetSecurityGroup(getSecurityGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_security_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroup).ToNot(BeNil())
		})
		It(`UpdateSecurityGroup request example`, func() {
			fmt.Println("\nUpdateSecurityGroup() result:")
			// begin-update_security_group

			securityGroupPatchModel := &vpcv1.SecurityGroupPatch{}
			securityGroupPatchModelAsPatch, asPatchErr := securityGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSecurityGroupOptions := vpcService.NewUpdateSecurityGroupOptions(
				"testString",
				securityGroupPatchModelAsPatch,
			)

			securityGroup, response, err := vpcService.UpdateSecurityGroup(updateSecurityGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroup, "", "  ")
			fmt.Println(string(b))

			// end-update_security_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroup).ToNot(BeNil())
		})
		It(`ListSecurityGroupRules request example`, func() {
			fmt.Println("\nListSecurityGroupRules() result:")
			// begin-list_security_group_rules

			listSecurityGroupRulesOptions := vpcService.NewListSecurityGroupRulesOptions(
				"testString",
			)

			securityGroupRuleCollection, response, err := vpcService.ListSecurityGroupRules(listSecurityGroupRulesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupRuleCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_security_group_rules

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRuleCollection).ToNot(BeNil())
		})
		It(`CreateSecurityGroupRule request example`, func() {
			fmt.Println("\nCreateSecurityGroupRule() result:")
			// begin-create_security_group_rule

			securityGroupRulePrototypeModel := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAnyPrototype{
				Direction: core.StringPtr("inbound"),
				Protocol:  core.StringPtr("any"),
			}

			createSecurityGroupRuleOptions := vpcService.NewCreateSecurityGroupRuleOptions(
				"testString",
				securityGroupRulePrototypeModel,
			)

			securityGroupRule, response, err := vpcService.CreateSecurityGroupRule(createSecurityGroupRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupRule, "", "  ")
			fmt.Println(string(b))

			// end-create_security_group_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(securityGroupRule).ToNot(BeNil())
		})
		It(`GetSecurityGroupRule request example`, func() {
			fmt.Println("\nGetSecurityGroupRule() result:")
			// begin-get_security_group_rule

			getSecurityGroupRuleOptions := vpcService.NewGetSecurityGroupRuleOptions(
				"testString",
				"testString",
			)

			securityGroupRule, response, err := vpcService.GetSecurityGroupRule(getSecurityGroupRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupRule, "", "  ")
			fmt.Println(string(b))

			// end-get_security_group_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRule).ToNot(BeNil())
		})
		It(`UpdateSecurityGroupRule request example`, func() {
			fmt.Println("\nUpdateSecurityGroupRule() result:")
			// begin-update_security_group_rule

			securityGroupRuleRemotePatchModel := &vpcv1.SecurityGroupRuleRemotePatchSecurityGroupRuleCIDRPrototype{
				CIDRBlock: core.StringPtr("10.0.0.0/8"),
			}

			securityGroupRulePatchModel := &vpcv1.SecurityGroupRulePatch{
				Remote: securityGroupRuleRemotePatchModel,
			}
			securityGroupRulePatchModelAsPatch, asPatchErr := securityGroupRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSecurityGroupRuleOptions := vpcService.NewUpdateSecurityGroupRuleOptions(
				"testString",
				"testString",
				securityGroupRulePatchModelAsPatch,
			)

			securityGroupRule, response, err := vpcService.UpdateSecurityGroupRule(updateSecurityGroupRuleOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupRule, "", "  ")
			fmt.Println(string(b))

			// end-update_security_group_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRule).ToNot(BeNil())
		})
		It(`ListSecurityGroupTargets request example`, func() {
			fmt.Println("\nListSecurityGroupTargets() result:")
			// begin-list_security_group_targets
			listSecurityGroupTargetsOptions := &vpcv1.ListSecurityGroupTargetsOptions{
				SecurityGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewSecurityGroupTargetsPager(listSecurityGroupTargetsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.SecurityGroupTargetReferenceIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_security_group_targets
		})
		It(`GetSecurityGroupTarget request example`, func() {
			fmt.Println("\nGetSecurityGroupTarget() result:")
			// begin-get_security_group_target

			getSecurityGroupTargetOptions := vpcService.NewGetSecurityGroupTargetOptions(
				"testString",
				"testString",
			)

			securityGroupTargetReference, response, err := vpcService.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupTargetReference, "", "  ")
			fmt.Println(string(b))

			// end-get_security_group_target

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupTargetReference).ToNot(BeNil())
		})
		It(`CreateSecurityGroupTargetBinding request example`, func() {
			fmt.Println("\nCreateSecurityGroupTargetBinding() result:")
			// begin-create_security_group_target_binding

			createSecurityGroupTargetBindingOptions := vpcService.NewCreateSecurityGroupTargetBindingOptions(
				"testString",
				"testString",
			)

			securityGroupTargetReference, response, err := vpcService.CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(securityGroupTargetReference, "", "  ")
			fmt.Println(string(b))

			// end-create_security_group_target_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupTargetReference).ToNot(BeNil())
		})
		It(`ListShareProfiles request example`, func() {
			fmt.Println("\nListShareProfiles() result:")
			// begin-list_share_profiles
			listShareProfilesOptions := &vpcv1.ListShareProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Sort:  core.StringPtr("name"),
			}

			pager, err := vpcService.NewShareProfilesPager(listShareProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ShareProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_share_profiles
		})
		It(`GetShareProfile request example`, func() {
			fmt.Println("\nGetShareProfile() result:")
			// begin-get_share_profile

			getShareProfileOptions := vpcService.NewGetShareProfileOptions(
				"dp2",
			)

			shareProfile, response, err := vpcService.GetShareProfile(getShareProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_share_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareProfile).ToNot(BeNil())
		})
		It(`ListShares request example`, func() {
			fmt.Println("\nListShares() result:")
			// begin-list_shares
			listSharesOptions := &vpcv1.ListSharesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ReplicationRole: core.StringPtr("none"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			pager, err := vpcService.NewSharesPager(listSharesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Share
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_shares
		})
		It(`CreateShare request example`, func() {
			fmt.Println("\nCreateShare() result:")
			// begin-create_share

			shareByOriginShareOriginShareModel := &vpcv1.ShareByOriginShareOriginShareShareIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::share:r006-0fe9e5d8-0a4d-4818-96ec-e99708644a58"),
			}

			sharePrototypeModel := &vpcv1.SharePrototypeShareByOriginShare{
				OriginShare: shareByOriginShareOriginShareModel,
			}

			createShareOptions := vpcService.NewCreateShareOptions(
				sharePrototypeModel,
			)

			share, response, err := vpcService.CreateShare(createShareOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(share, "", "  ")
			fmt.Println(string(b))

			// end-create_share

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(share).ToNot(BeNil())
		})
		It(`GetShare request example`, func() {
			fmt.Println("\nGetShare() result:")
			// begin-get_share

			getShareOptions := vpcService.NewGetShareOptions(
				"testString",
			)

			share, response, err := vpcService.GetShare(getShareOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(share, "", "  ")
			fmt.Println(string(b))

			// end-get_share

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(share).ToNot(BeNil())
		})
		It(`UpdateShare request example`, func() {
			fmt.Println("\nUpdateShare() result:")
			// begin-update_share

			sharePatchModel := &vpcv1.SharePatch{
				Size: core.Int64Ptr(int64(8000)),
			}
			sharePatchModelAsPatch, asPatchErr := sharePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareOptions := vpcService.NewUpdateShareOptions(
				"testString",
				sharePatchModelAsPatch,
			)
			updateShareOptions.SetIfMatch("W/96d225c4-56bd-43d9-98fc-d7148e5c5028")

			share, response, err := vpcService.UpdateShare(updateShareOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(share, "", "  ")
			fmt.Println(string(b))

			// end-update_share

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(share).ToNot(BeNil())
		})
		It(`ListShareAccessorBindings request example`, func() {
			fmt.Println("\nListShareAccessorBindings() result:")
			// begin-list_share_accessor_bindings
			listShareAccessorBindingsOptions := &vpcv1.ListShareAccessorBindingsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewShareAccessorBindingsPager(listShareAccessorBindingsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ShareAccessorBinding
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_share_accessor_bindings
		})
		It(`GetShareAccessorBinding request example`, func() {
			fmt.Println("\nGetShareAccessorBinding() result:")
			// begin-get_share_accessor_binding

			getShareAccessorBindingOptions := vpcService.NewGetShareAccessorBindingOptions(
				"testString",
				"testString",
			)

			shareAccessorBinding, response, err := vpcService.GetShareAccessorBinding(getShareAccessorBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareAccessorBinding, "", "  ")
			fmt.Println(string(b))

			// end-get_share_accessor_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareAccessorBinding).ToNot(BeNil())
		})
		It(`FailoverShare request example`, func() {
			// begin-failover_share

			failoverShareOptions := vpcService.NewFailoverShareOptions(
				"testString",
			)

			response, err := vpcService.FailoverShare(failoverShareOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from FailoverShare(): %d\n", response.StatusCode)
			}

			// end-failover_share

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`ListShareMountTargets request example`, func() {
			fmt.Println("\nListShareMountTargets() result:")
			// begin-list_share_mount_targets
			listShareMountTargetsOptions := &vpcv1.ListShareMountTargetsOptions{
				ShareID: core.StringPtr("testString"),
				Limit:   core.Int64Ptr(int64(10)),
				Name:    core.StringPtr("my-name"),
			}

			pager, err := vpcService.NewShareMountTargetsPager(listShareMountTargetsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ShareMountTarget
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_share_mount_targets
		})
		It(`CreateShareMountTarget request example`, func() {
			fmt.Println("\nCreateShareMountTarget() result:")
			// begin-create_share_mount_target

			shareMountTargetVirtualNetworkInterfacePrototypeModel := &vpcv1.ShareMountTargetVirtualNetworkInterfacePrototypeVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			shareMountTargetPrototypeModel := &vpcv1.ShareMountTargetPrototypeShareMountTargetByAccessControlModeSecurityGroup{
				AccessProtocol:          core.StringPtr("nfs4"),
				TransitEncryption:       core.StringPtr("none"),
				VirtualNetworkInterface: shareMountTargetVirtualNetworkInterfacePrototypeModel,
			}

			createShareMountTargetOptions := vpcService.NewCreateShareMountTargetOptions(
				"testString",
				shareMountTargetPrototypeModel,
			)

			shareMountTarget, response, err := vpcService.CreateShareMountTarget(createShareMountTargetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareMountTarget, "", "  ")
			fmt.Println(string(b))

			// end-create_share_mount_target

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(shareMountTarget).ToNot(BeNil())
		})
		It(`GetShareMountTarget request example`, func() {
			fmt.Println("\nGetShareMountTarget() result:")
			// begin-get_share_mount_target

			getShareMountTargetOptions := vpcService.NewGetShareMountTargetOptions(
				"testString",
				"testString",
			)

			shareMountTarget, response, err := vpcService.GetShareMountTarget(getShareMountTargetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareMountTarget, "", "  ")
			fmt.Println(string(b))

			// end-get_share_mount_target

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareMountTarget).ToNot(BeNil())
		})
		It(`UpdateShareMountTarget request example`, func() {
			fmt.Println("\nUpdateShareMountTarget() result:")
			// begin-update_share_mount_target

			shareMountTargetPatchModel := &vpcv1.ShareMountTargetPatch{}
			shareMountTargetPatchModelAsPatch, asPatchErr := shareMountTargetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareMountTargetOptions := vpcService.NewUpdateShareMountTargetOptions(
				"testString",
				"testString",
				shareMountTargetPatchModelAsPatch,
			)

			shareMountTarget, response, err := vpcService.UpdateShareMountTarget(updateShareMountTargetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareMountTarget, "", "  ")
			fmt.Println(string(b))

			// end-update_share_mount_target

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareMountTarget).ToNot(BeNil())
		})
		It(`ListShareSnapshots request example`, func() {
			fmt.Println("\nListShareSnapshots() result:")
			// begin-list_share_snapshots
			listShareSnapshotsOptions := &vpcv1.ListShareSnapshotsOptions{
				ShareID:            core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				Sort:               core.StringPtr("name"),
			}

			pager, err := vpcService.NewShareSnapshotsPager(listShareSnapshotsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ShareSnapshot
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_share_snapshots
		})
		It(`CreateShareSnapshot request example`, func() {
			fmt.Println("\nCreateShareSnapshot() result:")
			// begin-create_share_snapshot

			createShareSnapshotOptions := vpcService.NewCreateShareSnapshotOptions(
				"testString",
			)

			shareSnapshot, response, err := vpcService.CreateShareSnapshot(createShareSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareSnapshot, "", "  ")
			fmt.Println(string(b))

			// end-create_share_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(shareSnapshot).ToNot(BeNil())
		})
		It(`GetShareSnapshot request example`, func() {
			fmt.Println("\nGetShareSnapshot() result:")
			// begin-get_share_snapshot

			getShareSnapshotOptions := vpcService.NewGetShareSnapshotOptions(
				"testString",
				"testString",
			)

			shareSnapshot, response, err := vpcService.GetShareSnapshot(getShareSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareSnapshot, "", "  ")
			fmt.Println(string(b))

			// end-get_share_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareSnapshot).ToNot(BeNil())
		})
		It(`UpdateShareSnapshot request example`, func() {
			fmt.Println("\nUpdateShareSnapshot() result:")
			// begin-update_share_snapshot

			shareSnapshotPatchModel := &vpcv1.ShareSnapshotPatch{}
			shareSnapshotPatchModelAsPatch, asPatchErr := shareSnapshotPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareSnapshotOptions := vpcService.NewUpdateShareSnapshotOptions(
				"testString",
				"testString",
				shareSnapshotPatchModelAsPatch,
			)
			updateShareSnapshotOptions.SetIfMatch("W/96d225c4-56bd-43d9-98fc-d7148e5c5028")

			shareSnapshot, response, err := vpcService.UpdateShareSnapshot(updateShareSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareSnapshot, "", "  ")
			fmt.Println(string(b))

			// end-update_share_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareSnapshot).ToNot(BeNil())
		})
		It(`GetShareSource request example`, func() {
			fmt.Println("\nGetShareSource() result:")
			// begin-get_share_source

			getShareSourceOptions := vpcService.NewGetShareSourceOptions(
				"testString",
			)

			shareReference, response, err := vpcService.GetShareSource(getShareSourceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareReference, "", "  ")
			fmt.Println(string(b))

			// end-get_share_source

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareReference).ToNot(BeNil())
		})
		It(`ListSnapshotConsistencyGroups request example`, func() {
			fmt.Println("\nListSnapshotConsistencyGroups() result:")
			// begin-list_snapshot_consistency_groups
			listSnapshotConsistencyGroupsOptions := &vpcv1.ListSnapshotConsistencyGroupsOptions{
				BackupPolicyJobID:  core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
			}

			pager, err := vpcService.NewSnapshotConsistencyGroupsPager(listSnapshotConsistencyGroupsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.SnapshotConsistencyGroup
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_snapshot_consistency_groups
		})
		It(`CreateSnapshotConsistencyGroup request example`, func() {
			fmt.Println("\nCreateSnapshotConsistencyGroup() result:")
			// begin-create_snapshot_consistency_group

			volumeIdentityModel := &vpcv1.VolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			snapshotPrototypeSnapshotConsistencyGroupContextModel := &vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext{
				SourceVolume: volumeIdentityModel,
			}

			snapshotConsistencyGroupPrototypeModel := &vpcv1.SnapshotConsistencyGroupPrototypeSnapshotConsistencyGroupBySnapshots{
				Snapshots: []vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext{*snapshotPrototypeSnapshotConsistencyGroupContextModel},
			}

			createSnapshotConsistencyGroupOptions := vpcService.NewCreateSnapshotConsistencyGroupOptions(
				snapshotConsistencyGroupPrototypeModel,
			)

			snapshotConsistencyGroup, response, err := vpcService.CreateSnapshotConsistencyGroup(createSnapshotConsistencyGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotConsistencyGroup, "", "  ")
			fmt.Println(string(b))

			// end-create_snapshot_consistency_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
		It(`GetSnapshotConsistencyGroup request example`, func() {
			fmt.Println("\nGetSnapshotConsistencyGroup() result:")
			// begin-get_snapshot_consistency_group

			getSnapshotConsistencyGroupOptions := vpcService.NewGetSnapshotConsistencyGroupOptions(
				"testString",
			)

			snapshotConsistencyGroup, response, err := vpcService.GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotConsistencyGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_snapshot_consistency_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
		It(`UpdateSnapshotConsistencyGroup request example`, func() {
			fmt.Println("\nUpdateSnapshotConsistencyGroup() result:")
			// begin-update_snapshot_consistency_group

			snapshotConsistencyGroupPatchModel := &vpcv1.SnapshotConsistencyGroupPatch{
				DeleteSnapshotsOnDelete: core.BoolPtr(false),
			}
			snapshotConsistencyGroupPatchModelAsPatch, asPatchErr := snapshotConsistencyGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotConsistencyGroupOptions := vpcService.NewUpdateSnapshotConsistencyGroupOptions(
				"testString",
				snapshotConsistencyGroupPatchModelAsPatch,
			)
			updateSnapshotConsistencyGroupOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			snapshotConsistencyGroup, response, err := vpcService.UpdateSnapshotConsistencyGroup(updateSnapshotConsistencyGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotConsistencyGroup, "", "  ")
			fmt.Println(string(b))

			// end-update_snapshot_consistency_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
		It(`ListSnapshots request example`, func() {
			fmt.Println("\nListSnapshots() result:")
			// begin-list_snapshots
			listSnapshotsOptions := &vpcv1.ListSnapshotsOptions{
				BackupPolicyPlanID:             core.StringPtr("testString"),
				ClonesZoneName:                 core.StringPtr("us-south-1"),
				CopiesCRN:                      core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot:r006-f6bfa329-0e36-433f-a3bb-0df632e79263"),
				CopiesID:                       core.StringPtr("testString"),
				CopiesName:                     core.StringPtr("my-snapshot-copy"),
				CopiesRemoteRegionName:         core.StringPtr("us-south"),
				Limit:                          core.Int64Ptr(int64(10)),
				Name:                           core.StringPtr("my-name"),
				ResourceGroupID:                core.StringPtr("testString"),
				SnapshotConsistencyGroupCRN:    core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot-consistency-group:r006-fa329f6b-0e36-433f-a3bb-0df632e79263"),
				SnapshotConsistencyGroupID:     core.StringPtr("testString"),
				Sort:                           core.StringPtr("name"),
				SourceImageCRN:                 core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"),
				SourceImageID:                  core.StringPtr("testString"),
				SourceImageRemoteRegionName:    core.StringPtr("us-south"),
				SourceSnapshotID:               core.StringPtr("testString"),
				SourceSnapshotRemoteRegionName: core.StringPtr("us-south"),
				SourceVolumeCRN:                core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
				SourceVolumeID:                 core.StringPtr("testString"),
				SourceVolumeRemoteRegionName:   core.StringPtr("us-east"),
				Tag:                            core.StringPtr("testString"),
			}

			pager, err := vpcService.NewSnapshotsPager(listSnapshotsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Snapshot
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_snapshots
		})
		It(`CreateSnapshot request example`, func() {
			fmt.Println("\nCreateSnapshot() result:")
			// begin-create_snapshot

			snapshotIdentityByCRNModel := &vpcv1.SnapshotIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot:r006-f6bfa329-0e36-433f-a3bb-0df632e79263"),
			}

			snapshotPrototypeModel := &vpcv1.SnapshotPrototypeSnapshotBySourceSnapshot{
				SourceSnapshot: snapshotIdentityByCRNModel,
			}

			createSnapshotOptions := vpcService.NewCreateSnapshotOptions(
				snapshotPrototypeModel,
			)

			snapshot, response, err := vpcService.CreateSnapshot(createSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshot, "", "  ")
			fmt.Println(string(b))

			// end-create_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(snapshot).ToNot(BeNil())
		})
		It(`GetSnapshot request example`, func() {
			fmt.Println("\nGetSnapshot() result:")
			// begin-get_snapshot

			getSnapshotOptions := vpcService.NewGetSnapshotOptions(
				"testString",
			)

			snapshot, response, err := vpcService.GetSnapshot(getSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshot, "", "  ")
			fmt.Println(string(b))

			// end-get_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshot).ToNot(BeNil())
		})
		It(`UpdateSnapshot request example`, func() {
			fmt.Println("\nUpdateSnapshot() result:")
			// begin-update_snapshot

			snapshotPatchModel := &vpcv1.SnapshotPatch{
				Name: core.StringPtr("my-snapshot-updated"),
			}
			snapshotPatchModelAsPatch, asPatchErr := snapshotPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotOptions := vpcService.NewUpdateSnapshotOptions(
				"testString",
				snapshotPatchModelAsPatch,
			)
			updateSnapshotOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			snapshot, response, err := vpcService.UpdateSnapshot(updateSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshot, "", "  ")
			fmt.Println(string(b))

			// end-update_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshot).ToNot(BeNil())
		})
		It(`ListSnapshotClones request example`, func() {
			fmt.Println("\nListSnapshotClones() result:")
			// begin-list_snapshot_clones

			listSnapshotClonesOptions := vpcService.NewListSnapshotClonesOptions(
				"testString",
			)

			snapshotCloneCollection, response, err := vpcService.ListSnapshotClones(listSnapshotClonesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotCloneCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_snapshot_clones

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotCloneCollection).ToNot(BeNil())
		})
		It(`GetSnapshotClone request example`, func() {
			fmt.Println("\nGetSnapshotClone() result:")
			// begin-get_snapshot_clone

			getSnapshotCloneOptions := vpcService.NewGetSnapshotCloneOptions(
				"testString",
				"us-south-1",
			)

			snapshotClone, response, err := vpcService.GetSnapshotClone(getSnapshotCloneOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotClone, "", "  ")
			fmt.Println(string(b))

			// end-get_snapshot_clone

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotClone).ToNot(BeNil())
		})
		It(`CreateSnapshotClone request example`, func() {
			fmt.Println("\nCreateSnapshotClone() result:")
			// begin-create_snapshot_clone

			createSnapshotCloneOptions := vpcService.NewCreateSnapshotCloneOptions(
				"testString",
				"us-south-1",
			)

			snapshotClone, response, err := vpcService.CreateSnapshotClone(createSnapshotCloneOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotClone, "", "  ")
			fmt.Println(string(b))

			// end-create_snapshot_clone

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotClone).ToNot(BeNil())
		})
		It(`ListSnapshotInstanceProfiles request example`, func() {
			fmt.Println("\nListSnapshotInstanceProfiles() result:")
			// begin-list_snapshot_instance_profiles
			listSnapshotInstanceProfilesOptions := &vpcv1.ListSnapshotInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewSnapshotInstanceProfilesPager(listSnapshotInstanceProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_snapshot_instance_profiles
		})
		It(`ListSnapshotSoftwareAttachments request example`, func() {
			fmt.Println("\nListSnapshotSoftwareAttachments() result:")
			// begin-list_snapshot_software_attachments

			listSnapshotSoftwareAttachmentsOptions := vpcService.NewListSnapshotSoftwareAttachmentsOptions(
				"testString",
			)

			snapshotSoftwareAttachmentCollection, response, err := vpcService.ListSnapshotSoftwareAttachments(listSnapshotSoftwareAttachmentsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotSoftwareAttachmentCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_snapshot_software_attachments

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachmentCollection).ToNot(BeNil())
		})
		It(`GetSnapshotSoftwareAttachment request example`, func() {
			fmt.Println("\nGetSnapshotSoftwareAttachment() result:")
			// begin-get_snapshot_software_attachment

			getSnapshotSoftwareAttachmentOptions := vpcService.NewGetSnapshotSoftwareAttachmentOptions(
				"testString",
				"testString",
			)

			snapshotSoftwareAttachment, response, err := vpcService.GetSnapshotSoftwareAttachment(getSnapshotSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_snapshot_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachment).ToNot(BeNil())
		})
		It(`UpdateSnapshotSoftwareAttachment request example`, func() {
			fmt.Println("\nUpdateSnapshotSoftwareAttachment() result:")
			// begin-update_snapshot_software_attachment

			snapshotSoftwareAttachmentPatchModel := &vpcv1.SnapshotSoftwareAttachmentPatch{}
			snapshotSoftwareAttachmentPatchModelAsPatch, asPatchErr := snapshotSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotSoftwareAttachmentOptions := vpcService.NewUpdateSnapshotSoftwareAttachmentOptions(
				"testString",
				"testString",
				snapshotSoftwareAttachmentPatchModelAsPatch,
			)

			snapshotSoftwareAttachment, response, err := vpcService.UpdateSnapshotSoftwareAttachment(updateSnapshotSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_snapshot_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachment).ToNot(BeNil())
		})
		It(`ListSubnets request example`, func() {
			fmt.Println("\nListSubnets() result:")
			// begin-list_subnets
			listSubnetsOptions := &vpcv1.ListSubnetsOptions{
				Limit:            core.Int64Ptr(int64(10)),
				ResourceGroupID:  core.StringPtr("testString"),
				RoutingTableID:   core.StringPtr("testString"),
				RoutingTableName: core.StringPtr("my-routing-table"),
				VPCCRN:           core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:            core.StringPtr("testString"),
				VPCName:          core.StringPtr("my-vpc"),
				ZoneName:         core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewSubnetsPager(listSubnetsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Subnet
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_subnets
		})
		It(`CreateSubnet request example`, func() {
			fmt.Println("\nCreateSubnet() result:")
			// begin-create_subnet

			subnetPrototypeVPCModel := &vpcv1.SubnetPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			subnetPrototypeModel := &vpcv1.SubnetPrototypeSubnetByCIDR{
				VPC:           subnetPrototypeVPCModel,
				Ipv4CIDRBlock: core.StringPtr("10.0.0.0/24"),
			}

			createSubnetOptions := vpcService.NewCreateSubnetOptions(
				subnetPrototypeModel,
			)

			subnet, response, err := vpcService.CreateSubnet(createSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(subnet, "", "  ")
			fmt.Println(string(b))

			// end-create_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(subnet).ToNot(BeNil())
		})
		It(`GetSubnet request example`, func() {
			fmt.Println("\nGetSubnet() result:")
			// begin-get_subnet

			getSubnetOptions := vpcService.NewGetSubnetOptions(
				"testString",
			)

			subnet, response, err := vpcService.GetSubnet(getSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(subnet, "", "  ")
			fmt.Println(string(b))

			// end-get_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(subnet).ToNot(BeNil())
		})
		It(`UpdateSubnet request example`, func() {
			fmt.Println("\nUpdateSubnet() result:")
			// begin-update_subnet

			subnetPatchModel := &vpcv1.SubnetPatch{}
			subnetPatchModelAsPatch, asPatchErr := subnetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSubnetOptions := vpcService.NewUpdateSubnetOptions(
				"testString",
				subnetPatchModelAsPatch,
			)

			subnet, response, err := vpcService.UpdateSubnet(updateSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(subnet, "", "  ")
			fmt.Println(string(b))

			// end-update_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(subnet).ToNot(BeNil())
		})
		It(`GetSubnetNetworkACL request example`, func() {
			fmt.Println("\nGetSubnetNetworkACL() result:")
			// begin-get_subnet_network_acl

			getSubnetNetworkACLOptions := vpcService.NewGetSubnetNetworkACLOptions(
				"testString",
			)

			networkACL, response, err := vpcService.GetSubnetNetworkACL(getSubnetNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACL, "", "  ")
			fmt.Println(string(b))

			// end-get_subnet_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
		It(`ReplaceSubnetNetworkACL request example`, func() {
			fmt.Println("\nReplaceSubnetNetworkACL() result:")
			// begin-replace_subnet_network_acl

			networkACLIdentityModel := &vpcv1.NetworkACLIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::network-acl:r006-a4e28308-8ee7-46ab-8108-9f881f22bdbf"),
			}

			replaceSubnetNetworkACLOptions := vpcService.NewReplaceSubnetNetworkACLOptions(
				"testString",
				networkACLIdentityModel,
			)

			networkACL, response, err := vpcService.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(networkACL, "", "  ")
			fmt.Println(string(b))

			// end-replace_subnet_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACL).ToNot(BeNil())
		})
		It(`GetSubnetPublicGateway request example`, func() {
			fmt.Println("\nGetSubnetPublicGateway() result:")
			// begin-get_subnet_public_gateway

			getSubnetPublicGatewayOptions := vpcService.NewGetSubnetPublicGatewayOptions(
				"testString",
			)

			publicGateway, response, err := vpcService.GetSubnetPublicGateway(getSubnetPublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicGateway, "", "  ")
			fmt.Println(string(b))

			// end-get_subnet_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
		It(`SetSubnetPublicGateway request example`, func() {
			fmt.Println("\nSetSubnetPublicGateway() result:")
			// begin-set_subnet_public_gateway

			publicGatewayIdentityModel := &vpcv1.PublicGatewayIdentityPublicGatewayIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::public-gateway:r006-dc5431ef-1fc6-4861-adc9-a59d077d1241"),
			}

			setSubnetPublicGatewayOptions := vpcService.NewSetSubnetPublicGatewayOptions(
				"testString",
				publicGatewayIdentityModel,
			)

			publicGateway, response, err := vpcService.SetSubnetPublicGateway(setSubnetPublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicGateway, "", "  ")
			fmt.Println(string(b))

			// end-set_subnet_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicGateway).ToNot(BeNil())
		})
		It(`GetSubnetRoutingTable request example`, func() {
			fmt.Println("\nGetSubnetRoutingTable() result:")
			// begin-get_subnet_routing_table

			getSubnetRoutingTableOptions := vpcService.NewGetSubnetRoutingTableOptions(
				"testString",
			)

			routingTable, response, err := vpcService.GetSubnetRoutingTable(getSubnetRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(routingTable, "", "  ")
			fmt.Println(string(b))

			// end-get_subnet_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
		It(`ReplaceSubnetRoutingTable request example`, func() {
			fmt.Println("\nReplaceSubnetRoutingTable() result:")
			// begin-replace_subnet_routing_table

			routingTableIdentityModel := &vpcv1.RoutingTableIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc-routing-table:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b/r006-6885e83f-03b2-4603-8a86-db2a0f55c840"),
			}

			replaceSubnetRoutingTableOptions := vpcService.NewReplaceSubnetRoutingTableOptions(
				"testString",
				routingTableIdentityModel,
			)

			routingTable, response, err := vpcService.ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(routingTable, "", "  ")
			fmt.Println(string(b))

			// end-replace_subnet_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(routingTable).ToNot(BeNil())
		})
		It(`ListSubnetReservedIps request example`, func() {
			fmt.Println("\nListSubnetReservedIps() result:")
			// begin-list_subnet_reserved_ips
			listSubnetReservedIpsOptions := &vpcv1.ListSubnetReservedIpsOptions{
				SubnetID:           core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Sort:               core.StringPtr("name"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewSubnetReservedIpsPager(listSubnetReservedIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_subnet_reserved_ips
		})
		It(`CreateSubnetReservedIP request example`, func() {
			fmt.Println("\nCreateSubnetReservedIP() result:")
			// begin-create_subnet_reserved_ip

			createSubnetReservedIPOptions := vpcService.NewCreateSubnetReservedIPOptions(
				"testString",
			)

			reservedIP, response, err := vpcService.CreateSubnetReservedIP(createSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-create_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`GetSubnetReservedIP request example`, func() {
			fmt.Println("\nGetSubnetReservedIP() result:")
			// begin-get_subnet_reserved_ip

			getSubnetReservedIPOptions := vpcService.NewGetSubnetReservedIPOptions(
				"testString",
				"testString",
			)

			reservedIP, response, err := vpcService.GetSubnetReservedIP(getSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-get_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`UpdateSubnetReservedIP request example`, func() {
			fmt.Println("\nUpdateSubnetReservedIP() result:")
			// begin-update_subnet_reserved_ip

			reservedIPPatchModel := &vpcv1.ReservedIPPatch{
				Name: core.StringPtr("my-reserved-ip-updated"),
			}
			reservedIPPatchModelAsPatch, asPatchErr := reservedIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSubnetReservedIPOptions := vpcService.NewUpdateSubnetReservedIPOptions(
				"testString",
				"testString",
				reservedIPPatchModelAsPatch,
			)

			reservedIP, response, err := vpcService.UpdateSubnetReservedIP(updateSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIP, "", "  ")
			fmt.Println(string(b))

			// end-update_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
		It(`ListVirtualNetworkInterfaces request example`, func() {
			fmt.Println("\nListVirtualNetworkInterfaces() result:")
			// begin-list_virtual_network_interfaces
			listVirtualNetworkInterfacesOptions := &vpcv1.ListVirtualNetworkInterfacesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewVirtualNetworkInterfacesPager(listVirtualNetworkInterfacesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VirtualNetworkInterface
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_virtual_network_interfaces
		})
		It(`CreateVirtualNetworkInterface request example`, func() {
			fmt.Println("\nCreateVirtualNetworkInterface() result:")
			// begin-create_virtual_network_interface

			virtualNetworkInterfacePrototypeSubnetModel := &vpcv1.VirtualNetworkInterfacePrototypeSubnetSubnetIdentityByID{
				ID: core.StringPtr("69e55145-cc7d-4d8e-9e1f-cc3fb60b1793"),
			}

			createVirtualNetworkInterfaceOptions := vpcService.NewCreateVirtualNetworkInterfaceOptions()
			createVirtualNetworkInterfaceOptions.SetName("my-virtual-network-interface")
			createVirtualNetworkInterfaceOptions.SetSubnet(virtualNetworkInterfacePrototypeSubnetModel)

			virtualNetworkInterface, response, err := vpcService.CreateVirtualNetworkInterface(createVirtualNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-create_virtual_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
		It(`GetVirtualNetworkInterface request example`, func() {
			fmt.Println("\nGetVirtualNetworkInterface() result:")
			// begin-get_virtual_network_interface

			getVirtualNetworkInterfaceOptions := vpcService.NewGetVirtualNetworkInterfaceOptions(
				"testString",
			)

			virtualNetworkInterface, response, err := vpcService.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-get_virtual_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
		It(`UpdateVirtualNetworkInterface request example`, func() {
			fmt.Println("\nUpdateVirtualNetworkInterface() result:")
			// begin-update_virtual_network_interface

			virtualNetworkInterfacePatchModel := &vpcv1.VirtualNetworkInterfacePatch{}
			virtualNetworkInterfacePatchModelAsPatch, asPatchErr := virtualNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVirtualNetworkInterfaceOptions := vpcService.NewUpdateVirtualNetworkInterfaceOptions(
				"testString",
				virtualNetworkInterfacePatchModelAsPatch,
			)
			updateVirtualNetworkInterfaceOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			virtualNetworkInterface, response, err := vpcService.UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-update_virtual_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
		It(`ListNetworkInterfaceFloatingIps request example`, func() {
			fmt.Println("\nListNetworkInterfaceFloatingIps() result:")
			// begin-list_network_interface_floating_ips
			listNetworkInterfaceFloatingIpsOptions := &vpcv1.ListNetworkInterfaceFloatingIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
			}

			pager, err := vpcService.NewNetworkInterfaceFloatingIpsPager(listNetworkInterfaceFloatingIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.FloatingIPReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_network_interface_floating_ips
		})
		It(`GetNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nGetNetworkInterfaceFloatingIP() result:")
			// begin-get_network_interface_floating_ip

			getNetworkInterfaceFloatingIPOptions := vpcService.NewGetNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
			)

			floatingIPReference, response, err := vpcService.GetNetworkInterfaceFloatingIP(getNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIPReference, "", "  ")
			fmt.Println(string(b))

			// end-get_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPReference).ToNot(BeNil())
		})
		It(`AddNetworkInterfaceFloatingIP request example`, func() {
			fmt.Println("\nAddNetworkInterfaceFloatingIP() result:")
			// begin-add_network_interface_floating_ip

			addNetworkInterfaceFloatingIPOptions := vpcService.NewAddNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
			)

			floatingIPReference, response, err := vpcService.AddNetworkInterfaceFloatingIP(addNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(floatingIPReference, "", "  ")
			fmt.Println(string(b))

			// end-add_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPReference).ToNot(BeNil())
		})
		It(`ListVirtualNetworkInterfaceIps request example`, func() {
			fmt.Println("\nListVirtualNetworkInterfaceIps() result:")
			// begin-list_virtual_network_interface_ips
			listVirtualNetworkInterfaceIpsOptions := &vpcv1.ListVirtualNetworkInterfaceIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
			}

			pager, err := vpcService.NewVirtualNetworkInterfaceIpsPager(listVirtualNetworkInterfaceIpsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.ReservedIPReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_virtual_network_interface_ips
		})
		It(`GetVirtualNetworkInterfaceIP request example`, func() {
			fmt.Println("\nGetVirtualNetworkInterfaceIP() result:")
			// begin-get_virtual_network_interface_ip

			getVirtualNetworkInterfaceIPOptions := vpcService.NewGetVirtualNetworkInterfaceIPOptions(
				"testString",
				"testString",
			)

			reservedIPReference, response, err := vpcService.GetVirtualNetworkInterfaceIP(getVirtualNetworkInterfaceIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIPReference, "", "  ")
			fmt.Println(string(b))

			// end-get_virtual_network_interface_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPReference).ToNot(BeNil())
		})
		It(`AddVirtualNetworkInterfaceIP request example`, func() {
			fmt.Println("\nAddVirtualNetworkInterfaceIP() result:")
			// begin-add_virtual_network_interface_ip

			addVirtualNetworkInterfaceIPOptions := vpcService.NewAddVirtualNetworkInterfaceIPOptions(
				"testString",
				"testString",
			)

			reservedIPReference, response, err := vpcService.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservedIPReference, "", "  ")
			fmt.Println(string(b))

			// end-add_virtual_network_interface_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPReference).ToNot(BeNil())
		})
		It(`ListVolumeProfiles request example`, func() {
			fmt.Println("\nListVolumeProfiles() result:")
			// begin-list_volume_profiles
			listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVolumeProfilesPager(listVolumeProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VolumeProfile
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_volume_profiles
		})
		It(`GetVolumeProfile request example`, func() {
			fmt.Println("\nGetVolumeProfile() result:")
			// begin-get_volume_profile

			getVolumeProfileOptions := vpcService.NewGetVolumeProfileOptions(
				"10iops-tier",
			)

			volumeProfile, response, err := vpcService.GetVolumeProfile(getVolumeProfileOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeProfile, "", "  ")
			fmt.Println(string(b))

			// end-get_volume_profile

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeProfile).ToNot(BeNil())
		})
		It(`ListVolumes request example`, func() {
			fmt.Println("\nListVolumes() result:")
			// begin-list_volumes
			listVolumesOptions := &vpcv1.ListVolumesOptions{
				AttachmentState:             core.StringPtr("attached"),
				Encryption:                  core.StringPtr("provider_managed"),
				Limit:                       core.Int64Ptr(int64(10)),
				Name:                        core.StringPtr("my-name"),
				OperatingSystemArchitecture: core.StringPtr("amd64"),
				OperatingSystemFamily:       core.StringPtr("Ubuntu Server"),
				StorageGeneration:           core.Int64Ptr(int64(1)),
				Tag:                         core.StringPtr("testString"),
				ZoneName:                    core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewVolumesPager(listVolumesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Volume
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_volumes
		})
		It(`CreateVolume request example`, func() {
			fmt.Println("\nCreateVolume() result:")
			// begin-create_volume

			volumePrototypeProfileModel := &vpcv1.VolumePrototypeProfileVolumeProfileIdentityByName{
				Name: core.StringPtr("custom"),
			}

			volumePrototypeZoneModel := &vpcv1.VolumePrototypeZoneZoneIdentityByName{
				Name: core.StringPtr("us-south-1"),
			}

			volumePrototypeModel := &vpcv1.VolumePrototypeVolumeByCapacity{
				Name:     core.StringPtr("my-volume"),
				Profile:  volumePrototypeProfileModel,
				Zone:     volumePrototypeZoneModel,
				Capacity: core.Int64Ptr(int64(100)),
			}

			createVolumeOptions := vpcService.NewCreateVolumeOptions(
				volumePrototypeModel,
			)

			volume, response, err := vpcService.CreateVolume(createVolumeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(b))

			// end-create_volume

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volume).ToNot(BeNil())
		})
		It(`GetVolume request example`, func() {
			fmt.Println("\nGetVolume() result:")
			// begin-get_volume

			getVolumeOptions := vpcService.NewGetVolumeOptions(
				"testString",
			)

			volume, response, err := vpcService.GetVolume(getVolumeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(b))

			// end-get_volume

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volume).ToNot(BeNil())
		})
		It(`UpdateVolume request example`, func() {
			fmt.Println("\nUpdateVolume() result:")
			// begin-update_volume

			volumePatchModel := &vpcv1.VolumePatch{}
			volumePatchModelAsPatch, asPatchErr := volumePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeOptions := vpcService.NewUpdateVolumeOptions(
				"testString",
				volumePatchModelAsPatch,
			)
			updateVolumeOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			volume, response, err := vpcService.UpdateVolume(updateVolumeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(b))

			// end-update_volume

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volume).ToNot(BeNil())
		})
		It(`ListVolumeInstanceProfiles request example`, func() {
			fmt.Println("\nListVolumeInstanceProfiles() result:")
			// begin-list_volume_instance_profiles
			listVolumeInstanceProfilesOptions := &vpcv1.ListVolumeInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVolumeInstanceProfilesPager(listVolumeInstanceProfilesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_volume_instance_profiles
		})
		It(`ListVolumeJobs request example`, func() {
			fmt.Println("\nListVolumeJobs() result:")
			// begin-list_volume_jobs
			listVolumeJobsOptions := &vpcv1.ListVolumeJobsOptions{
				VolumeID: core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVolumeJobsPager(listVolumeJobsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VolumeJobIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_volume_jobs
		})
		It(`CreateVolumeJob request example`, func() {
			fmt.Println("\nCreateVolumeJob() result:")
			// begin-create_volume_job

			volumeProfileIdentityModel := &vpcv1.VolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumeJobTypeMigrateParametersModel := &vpcv1.VolumeJobTypeMigrateParameters{
				Profile: volumeProfileIdentityModel,
			}

			volumeJobPrototypeModel := &vpcv1.VolumeJobPrototypeVolumeJobTypeMigratePrototype{
				JobType:    core.StringPtr("migrate"),
				Parameters: volumeJobTypeMigrateParametersModel,
			}

			createVolumeJobOptions := vpcService.NewCreateVolumeJobOptions(
				"testString",
				volumeJobPrototypeModel,
			)

			volumeJob, response, err := vpcService.CreateVolumeJob(createVolumeJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeJob, "", "  ")
			fmt.Println(string(b))

			// end-create_volume_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volumeJob).ToNot(BeNil())
		})
		It(`GetVolumeJob request example`, func() {
			fmt.Println("\nGetVolumeJob() result:")
			// begin-get_volume_job

			getVolumeJobOptions := vpcService.NewGetVolumeJobOptions(
				"testString",
				"testString",
			)

			volumeJob, response, err := vpcService.GetVolumeJob(getVolumeJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeJob, "", "  ")
			fmt.Println(string(b))

			// end-get_volume_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeJob).ToNot(BeNil())
		})
		It(`UpdateVolumeJob request example`, func() {
			fmt.Println("\nUpdateVolumeJob() result:")
			// begin-update_volume_job

			volumeJobPatchModel := &vpcv1.VolumeJobPatch{
				Name: core.StringPtr("my-updated-volume-job-name"),
			}
			volumeJobPatchModelAsPatch, asPatchErr := volumeJobPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeJobOptions := vpcService.NewUpdateVolumeJobOptions(
				"testString",
				"testString",
				volumeJobPatchModelAsPatch,
			)

			volumeJob, response, err := vpcService.UpdateVolumeJob(updateVolumeJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeJob, "", "  ")
			fmt.Println(string(b))

			// end-update_volume_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeJob).ToNot(BeNil())
		})
		It(`CancelVolumeJob request example`, func() {
			fmt.Println("\nCancelVolumeJob() result:")
			// begin-cancel_volume_job

			cancelVolumeJobOptions := vpcService.NewCancelVolumeJobOptions(
				"testString",
				"testString",
			)

			volumeJob, response, err := vpcService.CancelVolumeJob(cancelVolumeJobOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeJob, "", "  ")
			fmt.Println(string(b))

			// end-cancel_volume_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(volumeJob).ToNot(BeNil())
		})
		It(`ListVolumeSoftwareAttachments request example`, func() {
			fmt.Println("\nListVolumeSoftwareAttachments() result:")
			// begin-list_volume_software_attachments

			listVolumeSoftwareAttachmentsOptions := vpcService.NewListVolumeSoftwareAttachmentsOptions(
				"testString",
			)

			volumeSoftwareAttachmentCollection, response, err := vpcService.ListVolumeSoftwareAttachments(listVolumeSoftwareAttachmentsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeSoftwareAttachmentCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_volume_software_attachments

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachmentCollection).ToNot(BeNil())
		})
		It(`GetVolumeSoftwareAttachment request example`, func() {
			fmt.Println("\nGetVolumeSoftwareAttachment() result:")
			// begin-get_volume_software_attachment

			getVolumeSoftwareAttachmentOptions := vpcService.NewGetVolumeSoftwareAttachmentOptions(
				"testString",
				"testString",
			)

			volumeSoftwareAttachment, response, err := vpcService.GetVolumeSoftwareAttachment(getVolumeSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-get_volume_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachment).ToNot(BeNil())
		})
		It(`UpdateVolumeSoftwareAttachment request example`, func() {
			fmt.Println("\nUpdateVolumeSoftwareAttachment() result:")
			// begin-update_volume_software_attachment

			volumeSoftwareAttachmentPatchModel := &vpcv1.VolumeSoftwareAttachmentPatch{}
			volumeSoftwareAttachmentPatchModelAsPatch, asPatchErr := volumeSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeSoftwareAttachmentOptions := vpcService.NewUpdateVolumeSoftwareAttachmentOptions(
				"testString",
				"testString",
				volumeSoftwareAttachmentPatchModelAsPatch,
			)

			volumeSoftwareAttachment, response, err := vpcService.UpdateVolumeSoftwareAttachment(updateVolumeSoftwareAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(volumeSoftwareAttachment, "", "  ")
			fmt.Println(string(b))

			// end-update_volume_software_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachment).ToNot(BeNil())
		})
		It(`ListVpcs request example`, func() {
			fmt.Println("\nListVpcs() result:")
			// begin-list_vpcs
			listVpcsOptions := &vpcv1.ListVpcsOptions{
				ClassicAccess:   core.BoolPtr(true),
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			pager, err := vpcService.NewVpcsPager(listVpcsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPC
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpcs
		})
		It(`CreateVPC request example`, func() {
			fmt.Println("\nCreateVPC() result:")
			// begin-create_vpc

			createVPCOptions := vpcService.NewCreateVPCOptions()

			vpc, response, err := vpcService.CreateVPC(createVPCOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpc, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpc).ToNot(BeNil())
		})
		It(`GetVPC request example`, func() {
			fmt.Println("\nGetVPC() result:")
			// begin-get_vpc

			getVPCOptions := vpcService.NewGetVPCOptions(
				"testString",
			)

			vpc, response, err := vpcService.GetVPC(getVPCOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpc, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpc).ToNot(BeNil())
		})
		It(`UpdateVPC request example`, func() {
			fmt.Println("\nUpdateVPC() result:")
			// begin-update_vpc

			vpcPatchModel := &vpcv1.VPCPatch{}
			vpcPatchModelAsPatch, asPatchErr := vpcPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCOptions := vpcService.NewUpdateVPCOptions(
				"testString",
				vpcPatchModelAsPatch,
			)
			updateVPCOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			vpc, response, err := vpcService.UpdateVPC(updateVPCOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpc, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpc).ToNot(BeNil())
		})
		It(`GetVPCDefaultNetworkACL request example`, func() {
			fmt.Println("\nGetVPCDefaultNetworkACL() result:")
			// begin-get_vpc_default_network_acl

			getVPCDefaultNetworkACLOptions := vpcService.NewGetVPCDefaultNetworkACLOptions(
				"testString",
			)

			defaultNetworkACL, response, err := vpcService.GetVPCDefaultNetworkACL(getVPCDefaultNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(defaultNetworkACL, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_default_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultNetworkACL).ToNot(BeNil())
		})
		It(`GetVPCDefaultRoutingTable request example`, func() {
			fmt.Println("\nGetVPCDefaultRoutingTable() result:")
			// begin-get_vpc_default_routing_table

			getVPCDefaultRoutingTableOptions := vpcService.NewGetVPCDefaultRoutingTableOptions(
				"testString",
			)

			defaultRoutingTable, response, err := vpcService.GetVPCDefaultRoutingTable(getVPCDefaultRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(defaultRoutingTable, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_default_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultRoutingTable).ToNot(BeNil())
		})
		It(`GetVPCDefaultSecurityGroup request example`, func() {
			fmt.Println("\nGetVPCDefaultSecurityGroup() result:")
			// begin-get_vpc_default_security_group

			getVPCDefaultSecurityGroupOptions := vpcService.NewGetVPCDefaultSecurityGroupOptions(
				"testString",
			)

			defaultSecurityGroup, response, err := vpcService.GetVPCDefaultSecurityGroup(getVPCDefaultSecurityGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(defaultSecurityGroup, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_default_security_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultSecurityGroup).ToNot(BeNil())
		})
		It(`ListVPCAddressPrefixes request example`, func() {
			fmt.Println("\nListVPCAddressPrefixes() result:")
			// begin-list_vpc_address_prefixes
			listVPCAddressPrefixesOptions := &vpcv1.ListVPCAddressPrefixesOptions{
				VPCID: core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVPCAddressPrefixesPager(listVPCAddressPrefixesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.AddressPrefix
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpc_address_prefixes
		})
		It(`CreateVPCAddressPrefix request example`, func() {
			fmt.Println("\nCreateVPCAddressPrefix() result:")
			// begin-create_vpc_address_prefix

			addressPrefixPrototypeZoneModel := &vpcv1.AddressPrefixPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createVPCAddressPrefixOptions := vpcService.NewCreateVPCAddressPrefixOptions(
				"testString",
				"10.240.0.0/18",
				addressPrefixPrototypeZoneModel,
			)

			addressPrefix, response, err := vpcService.CreateVPCAddressPrefix(createVPCAddressPrefixOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(addressPrefix, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc_address_prefix

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(addressPrefix).ToNot(BeNil())
		})
		It(`GetVPCAddressPrefix request example`, func() {
			fmt.Println("\nGetVPCAddressPrefix() result:")
			// begin-get_vpc_address_prefix

			getVPCAddressPrefixOptions := vpcService.NewGetVPCAddressPrefixOptions(
				"testString",
				"testString",
			)

			addressPrefix, response, err := vpcService.GetVPCAddressPrefix(getVPCAddressPrefixOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(addressPrefix, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_address_prefix

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(addressPrefix).ToNot(BeNil())
		})
		It(`UpdateVPCAddressPrefix request example`, func() {
			fmt.Println("\nUpdateVPCAddressPrefix() result:")
			// begin-update_vpc_address_prefix

			addressPrefixPatchModel := &vpcv1.AddressPrefixPatch{
				Name: core.StringPtr("my-vpc-address-prefix-updated"),
			}
			addressPrefixPatchModelAsPatch, asPatchErr := addressPrefixPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCAddressPrefixOptions := vpcService.NewUpdateVPCAddressPrefixOptions(
				"testString",
				"testString",
				addressPrefixPatchModelAsPatch,
			)

			addressPrefix, response, err := vpcService.UpdateVPCAddressPrefix(updateVPCAddressPrefixOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(addressPrefix, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc_address_prefix

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(addressPrefix).ToNot(BeNil())
		})
		It(`ListVPCDnsResolutionBindings request example`, func() {
			fmt.Println("\nListVPCDnsResolutionBindings() result:")
			// begin-list_vpc_dns_resolution_bindings
			listVPCDnsResolutionBindingsOptions := &vpcv1.ListVPCDnsResolutionBindingsOptions{
				VPCID:     core.StringPtr("testString"),
				AccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:     core.Int64Ptr(int64(10)),
				Name:      core.StringPtr("my-name"),
				Sort:      core.StringPtr("name"),
				VPCCRN:    core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCName:   core.StringPtr("my-vpc"),
			}

			pager, err := vpcService.NewVPCDnsResolutionBindingsPager(listVPCDnsResolutionBindingsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VpcdnsResolutionBinding
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpc_dns_resolution_bindings
		})
		It(`CreateVPCDnsResolutionBinding request example`, func() {
			fmt.Println("\nCreateVPCDnsResolutionBinding() result:")
			// begin-create_vpc_dns_resolution_binding

			vpcIdentityModel := &vpcv1.VPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			createVPCDnsResolutionBindingOptions := vpcService.NewCreateVPCDnsResolutionBindingOptions(
				"testString",
				vpcIdentityModel,
			)

			vpcdnsResolutionBinding, response, err := vpcService.CreateVPCDnsResolutionBinding(createVPCDnsResolutionBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpcdnsResolutionBinding, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc_dns_resolution_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
		It(`GetVPCDnsResolutionBinding request example`, func() {
			fmt.Println("\nGetVPCDnsResolutionBinding() result:")
			// begin-get_vpc_dns_resolution_binding

			getVPCDnsResolutionBindingOptions := vpcService.NewGetVPCDnsResolutionBindingOptions(
				"testString",
				"testString",
			)

			vpcdnsResolutionBinding, response, err := vpcService.GetVPCDnsResolutionBinding(getVPCDnsResolutionBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpcdnsResolutionBinding, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_dns_resolution_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
		It(`UpdateVPCDnsResolutionBinding request example`, func() {
			fmt.Println("\nUpdateVPCDnsResolutionBinding() result:")
			// begin-update_vpc_dns_resolution_binding

			vpcdnsResolutionBindingPatchModel := &vpcv1.VpcdnsResolutionBindingPatch{}
			vpcdnsResolutionBindingPatchModelAsPatch, asPatchErr := vpcdnsResolutionBindingPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCDnsResolutionBindingOptions := vpcService.NewUpdateVPCDnsResolutionBindingOptions(
				"testString",
				"testString",
				vpcdnsResolutionBindingPatchModelAsPatch,
			)

			vpcdnsResolutionBinding, response, err := vpcService.UpdateVPCDnsResolutionBinding(updateVPCDnsResolutionBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpcdnsResolutionBinding, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc_dns_resolution_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
		It(`ListVPCRoutes request example`, func() {
			fmt.Println("\nListVPCRoutes() result:")
			// begin-list_vpc_routes
			listVPCRoutesOptions := &vpcv1.ListVPCRoutesOptions{
				VPCID:    core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
				ZoneName: core.StringPtr("us-south-1"),
			}

			pager, err := vpcService.NewVPCRoutesPager(listVPCRoutesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Route
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpc_routes
		})
		It(`CreateVPCRoute request example`, func() {
			fmt.Println("\nCreateVPCRoute() result:")
			// begin-create_vpc_route

			routePrototypeZoneModel := &vpcv1.RoutePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createVPCRouteOptions := vpcService.NewCreateVPCRouteOptions(
				"testString",
				"192.168.3.0/24",
				routePrototypeZoneModel,
			)

			createVPCRouteResponse, response, err := vpcService.CreateVPCRoute(createVPCRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(createVPCRouteResponse, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(createVPCRouteResponse).ToNot(BeNil())
		})
		It(`GetVPCRoute request example`, func() {
			fmt.Println("\nGetVPCRoute() result:")
			// begin-get_vpc_route

			getVPCRouteOptions := vpcService.NewGetVPCRouteOptions(
				"testString",
				"testString",
			)

			getVPCRouteResponse, response, err := vpcService.GetVPCRoute(getVPCRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(getVPCRouteResponse, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(getVPCRouteResponse).ToNot(BeNil())
		})
		It(`UpdateVPCRoute request example`, func() {
			fmt.Println("\nUpdateVPCRoute() result:")
			// begin-update_vpc_route

			routeNextHopPatchModel := &vpcv1.RouteNextHopPatchRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("10.240.0.12"),
			}

			routePatchModel := &vpcv1.RoutePatch{
				NextHop: routeNextHopPatchModel,
			}
			routePatchModelAsPatch, asPatchErr := routePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRouteOptions := vpcService.NewUpdateVPCRouteOptions(
				"testString",
				"testString",
				routePatchModelAsPatch,
			)

			updateVPCRouteResponse, response, err := vpcService.UpdateVPCRoute(updateVPCRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(updateVPCRouteResponse, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(updateVPCRouteResponse).ToNot(BeNil())
		})
		It(`ListVPCRoutingTables request example`, func() {
			fmt.Println("\nListVPCRoutingTables() result:")
			// begin-list_vpc_routing_tables
			listVPCRoutingTablesOptions := &vpcv1.ListVPCRoutingTablesOptions{
				VPCID:     core.StringPtr("testString"),
				IsDefault: core.BoolPtr(true),
				Limit:     core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVPCRoutingTablesPager(listVPCRoutingTablesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.RoutingTable
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpc_routing_tables
		})
		It(`CreateVPCRoutingTable request example`, func() {
			fmt.Println("\nCreateVPCRoutingTable() result:")
			// begin-create_vpc_routing_table

			createVPCRoutingTableOptions := vpcService.NewCreateVPCRoutingTableOptions(
				"testString",
			)

			routingTable, response, err := vpcService.CreateVPCRoutingTable(createVPCRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(routingTable, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(routingTable).ToNot(BeNil())
		})
		It(`GetVPCRoutingTable request example`, func() {
			fmt.Println("\nGetVPCRoutingTable() result:")
			// begin-get_vpc_routing_table

			getVPCRoutingTableOptions := vpcService.NewGetVPCRoutingTableOptions(
				"testString",
				"testString",
			)

			routingTable, response, err := vpcService.GetVPCRoutingTable(getVPCRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(routingTable, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
		It(`UpdateVPCRoutingTable request example`, func() {
			fmt.Println("\nUpdateVPCRoutingTable() result:")
			// begin-update_vpc_routing_table

			routingTablePatchModel := &vpcv1.RoutingTablePatch{}
			routingTablePatchModelAsPatch, asPatchErr := routingTablePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRoutingTableOptions := vpcService.NewUpdateVPCRoutingTableOptions(
				"testString",
				"testString",
				routingTablePatchModelAsPatch,
			)
			updateVPCRoutingTableOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			routingTable, response, err := vpcService.UpdateVPCRoutingTable(updateVPCRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(routingTable, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
		It(`ListVPCRoutingTableRoutes request example`, func() {
			fmt.Println("\nListVPCRoutingTableRoutes() result:")
			// begin-list_vpc_routing_table_routes
			listVPCRoutingTableRoutesOptions := &vpcv1.ListVPCRoutingTableRoutesOptions{
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
				Limit:          core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVPCRoutingTableRoutesPager(listVPCRoutingTableRoutesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.Route
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpc_routing_table_routes
		})
		It(`CreateVPCRoutingTableRoute request example`, func() {
			fmt.Println("\nCreateVPCRoutingTableRoute() result:")
			// begin-create_vpc_routing_table_route

			routePrototypeZoneModel := &vpcv1.RoutePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createVPCRoutingTableRouteOptions := vpcService.NewCreateVPCRoutingTableRouteOptions(
				"testString",
				"testString",
				"192.168.3.0/24",
				routePrototypeZoneModel,
			)

			route, response, err := vpcService.CreateVPCRoutingTableRoute(createVPCRoutingTableRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(route, "", "  ")
			fmt.Println(string(b))

			// end-create_vpc_routing_table_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(route).ToNot(BeNil())
		})
		It(`GetVPCRoutingTableRoute request example`, func() {
			fmt.Println("\nGetVPCRoutingTableRoute() result:")
			// begin-get_vpc_routing_table_route

			getVPCRoutingTableRouteOptions := vpcService.NewGetVPCRoutingTableRouteOptions(
				"testString",
				"testString",
				"testString",
			)

			route, response, err := vpcService.GetVPCRoutingTableRoute(getVPCRoutingTableRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(route, "", "  ")
			fmt.Println(string(b))

			// end-get_vpc_routing_table_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
		It(`UpdateVPCRoutingTableRoute request example`, func() {
			fmt.Println("\nUpdateVPCRoutingTableRoute() result:")
			// begin-update_vpc_routing_table_route

			routePatchModel := &vpcv1.RoutePatch{
				Advertise: core.BoolPtr(false),
			}
			routePatchModelAsPatch, asPatchErr := routePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRoutingTableRouteOptions := vpcService.NewUpdateVPCRoutingTableRouteOptions(
				"testString",
				"testString",
				"testString",
				routePatchModelAsPatch,
			)

			route, response, err := vpcService.UpdateVPCRoutingTableRoute(updateVPCRoutingTableRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(route, "", "  ")
			fmt.Println(string(b))

			// end-update_vpc_routing_table_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
		It(`ListIkePolicies request example`, func() {
			fmt.Println("\nListIkePolicies() result:")
			// begin-list_ike_policies
			listIkePoliciesOptions := &vpcv1.ListIkePoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewIkePoliciesPager(listIkePoliciesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.IkePolicy
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_ike_policies
		})
		It(`CreateIkePolicy request example`, func() {
			fmt.Println("\nCreateIkePolicy() result:")
			// begin-create_ike_policy

			createIkePolicyOptions := vpcService.NewCreateIkePolicyOptions(
				int64(1),
			)

			ikePolicy, response, err := vpcService.CreateIkePolicy(createIkePolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(ikePolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_ike_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(ikePolicy).ToNot(BeNil())
		})
		It(`GetIkePolicy request example`, func() {
			fmt.Println("\nGetIkePolicy() result:")
			// begin-get_ike_policy

			getIkePolicyOptions := vpcService.NewGetIkePolicyOptions(
				"testString",
			)

			ikePolicy, response, err := vpcService.GetIkePolicy(getIkePolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(ikePolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_ike_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(ikePolicy).ToNot(BeNil())
		})
		It(`UpdateIkePolicy request example`, func() {
			fmt.Println("\nUpdateIkePolicy() result:")
			// begin-update_ike_policy

			ikePolicyPatchModel := &vpcv1.IkePolicyPatch{
				KeyLifetime: core.Int64Ptr(int64(86400)),
			}
			ikePolicyPatchModelAsPatch, asPatchErr := ikePolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateIkePolicyOptions := vpcService.NewUpdateIkePolicyOptions(
				"testString",
				ikePolicyPatchModelAsPatch,
			)

			ikePolicy, response, err := vpcService.UpdateIkePolicy(updateIkePolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(ikePolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_ike_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(ikePolicy).ToNot(BeNil())
		})
		It(`ListIkePolicyConnections request example`, func() {
			fmt.Println("\nListIkePolicyConnections() result:")
			// begin-list_ike_policy_connections
			listIkePolicyConnectionsOptions := &vpcv1.ListIkePolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewIkePolicyConnectionsPager(listIkePolicyConnectionsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_ike_policy_connections
		})
		It(`ListIpsecPolicies request example`, func() {
			fmt.Println("\nListIpsecPolicies() result:")
			// begin-list_ipsec_policies
			listIpsecPoliciesOptions := &vpcv1.ListIpsecPoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewIpsecPoliciesPager(listIpsecPoliciesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.IPsecPolicy
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_ipsec_policies
		})
		It(`CreateIpsecPolicy request example`, func() {
			fmt.Println("\nCreateIpsecPolicy() result:")
			// begin-create_ipsec_policy

			createIpsecPolicyOptions := vpcService.NewCreateIpsecPolicyOptions()

			iPsecPolicy, response, err := vpcService.CreateIpsecPolicy(createIpsecPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(iPsecPolicy, "", "  ")
			fmt.Println(string(b))

			// end-create_ipsec_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
		It(`GetIpsecPolicy request example`, func() {
			fmt.Println("\nGetIpsecPolicy() result:")
			// begin-get_ipsec_policy

			getIpsecPolicyOptions := vpcService.NewGetIpsecPolicyOptions(
				"testString",
			)

			iPsecPolicy, response, err := vpcService.GetIpsecPolicy(getIpsecPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(iPsecPolicy, "", "  ")
			fmt.Println(string(b))

			// end-get_ipsec_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
		It(`UpdateIpsecPolicy request example`, func() {
			fmt.Println("\nUpdateIpsecPolicy() result:")
			// begin-update_ipsec_policy

			iPsecPolicyPatchModel := &vpcv1.IPsecPolicyPatch{
				Pfs: core.StringPtr("disabled"),
			}
			iPsecPolicyPatchModelAsPatch, asPatchErr := iPsecPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateIpsecPolicyOptions := vpcService.NewUpdateIpsecPolicyOptions(
				"testString",
				iPsecPolicyPatchModelAsPatch,
			)

			iPsecPolicy, response, err := vpcService.UpdateIpsecPolicy(updateIpsecPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(iPsecPolicy, "", "  ")
			fmt.Println(string(b))

			// end-update_ipsec_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
		It(`ListIpsecPolicyConnections request example`, func() {
			fmt.Println("\nListIpsecPolicyConnections() result:")
			// begin-list_ipsec_policy_connections
			listIpsecPolicyConnectionsOptions := &vpcv1.ListIpsecPolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewIpsecPolicyConnectionsPager(listIpsecPolicyConnectionsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_ipsec_policy_connections
		})
		It(`ListVPNGateways request example`, func() {
			fmt.Println("\nListVPNGateways() result:")
			// begin-list_vpn_gateways
			listVPNGatewaysOptions := &vpcv1.ListVPNGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Mode:            core.StringPtr("policy"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			pager, err := vpcService.NewVPNGatewaysPager(listVPNGatewaysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNGatewayIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_gateways
		})
		It(`CreateVPNGateway request example`, func() {
			fmt.Println("\nCreateVPNGateway() result:")
			// begin-create_vpn_gateway

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			vpnGatewayMemberPrivateIPPrototypeModel := &vpcv1.VPNGatewayMemberPrivateIPPrototypeReservedIPPrototypeVPNGatewayContext{
				Subnet: subnetIdentityModel,
			}

			vpnGatewayMemberPrototypeModel := &vpcv1.VPNGatewayMemberPrototype{
				PrivateIP: vpnGatewayMemberPrivateIPPrototypeModel,
			}

			vpnGatewayPrototypeModel := &vpcv1.VPNGatewayPrototypeVPNGatewayRegionalAvailabilityModeByMemberVPNGatewayRegionalAvailabilityModeByMemberVPNGatewayRouteModePrototype{
				AvailabilityMode: core.StringPtr("regional"),
				Members:          []vpcv1.VPNGatewayMemberPrototype{*vpnGatewayMemberPrototypeModel},
			}

			createVPNGatewayOptions := vpcService.NewCreateVPNGatewayOptions(
				vpnGatewayPrototypeModel,
			)

			vpnGateway, response, err := vpcService.CreateVPNGateway(createVPNGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGateway, "", "  ")
			fmt.Println(string(b))

			// end-create_vpn_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnGateway).ToNot(BeNil())
		})
		It(`GetVPNGateway request example`, func() {
			fmt.Println("\nGetVPNGateway() result:")
			// begin-get_vpn_gateway

			getVPNGatewayOptions := vpcService.NewGetVPNGatewayOptions(
				"testString",
			)

			vpnGateway, response, err := vpcService.GetVPNGateway(getVPNGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGateway, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGateway).ToNot(BeNil())
		})
		It(`UpdateVPNGateway request example`, func() {
			fmt.Println("\nUpdateVPNGateway() result:")
			// begin-update_vpn_gateway

			vpnGatewayPatchModel := &vpcv1.VPNGatewayPatch{
				Name: core.StringPtr("my-vpn-gateway-updated"),
			}
			vpnGatewayPatchModelAsPatch, asPatchErr := vpnGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNGatewayOptions := vpcService.NewUpdateVPNGatewayOptions(
				"testString",
				vpnGatewayPatchModelAsPatch,
			)

			vpnGateway, response, err := vpcService.UpdateVPNGateway(updateVPNGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGateway, "", "  ")
			fmt.Println(string(b))

			// end-update_vpn_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGateway).ToNot(BeNil())
		})
		It(`ListVPNGatewayAdvertisedCidrs request example`, func() {
			fmt.Println("\nListVPNGatewayAdvertisedCidrs() result:")
			// begin-list_vpn_gateway_advertised_cidrs

			listVPNGatewayAdvertisedCidrsOptions := vpcService.NewListVPNGatewayAdvertisedCidrsOptions(
				"testString",
			)

			vpnGatewayAdvertisedCIDRCollection, response, err := vpcService.ListVPNGatewayAdvertisedCidrs(listVPNGatewayAdvertisedCidrsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayAdvertisedCIDRCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_vpn_gateway_advertised_cidrs

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayAdvertisedCIDRCollection).ToNot(BeNil())
		})
		It(`CheckVPNGatewayAdvertisedCIDR request example`, func() {
			// begin-check_vpn_gateway_advertised_cidr

			checkVPNGatewayAdvertisedCIDROptions := vpcService.NewCheckVPNGatewayAdvertisedCIDROptions(
				"192.168.1.0/24",
				"testString",
			)

			response, err := vpcService.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from CheckVPNGatewayAdvertisedCIDR(): %d\n", response.StatusCode)
			}

			// end-check_vpn_gateway_advertised_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`AddVPNGatewayAdvertisedCIDR request example`, func() {
			// begin-add_vpn_gateway_advertised_cidr

			addVPNGatewayAdvertisedCIDROptions := vpcService.NewAddVPNGatewayAdvertisedCIDROptions(
				"192.168.1.0/24",
				"testString",
			)

			response, err := vpcService.AddVPNGatewayAdvertisedCIDR(addVPNGatewayAdvertisedCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from AddVPNGatewayAdvertisedCIDR(): %d\n", response.StatusCode)
			}

			// end-add_vpn_gateway_advertised_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`ListVPNGatewayConnections request example`, func() {
			fmt.Println("\nListVPNGatewayConnections() result:")
			// begin-list_vpn_gateway_connections
			listVPNGatewayConnectionsOptions := &vpcv1.ListVPNGatewayConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
				Status:       core.StringPtr("down"),
			}

			pager, err := vpcService.NewVPNGatewayConnectionsPager(listVPNGatewayConnectionsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_gateway_connections
		})
		It(`CreateVPNGatewayConnection request example`, func() {
			fmt.Println("\nCreateVPNGatewayConnection() result:")
			// begin-create_vpn_gateway_connection

			vpnGatewayConnectionDynamicRouteModePeerPrototypeModel := &vpcv1.VPNGatewayConnectionDynamicRouteModePeerPrototypeVPNGatewayConnectionPeerByAddress{
				Asn:     core.Int64Ptr(int64(38)),
				Address: core.StringPtr("192.0.2.5"),
			}

			ipModel := &vpcv1.IP{
				Address: core.StringPtr("169.254.0.2"),
			}

			vpnGatewayConnectionTunnelPrototypeModel := &vpcv1.VPNGatewayConnectionTunnelPrototype{
				NeighborIP:        ipModel,
				TunnelInterfaceIP: ipModel,
			}

			vpnGatewayConnectionPrototypeModel := &vpcv1.VPNGatewayConnectionPrototypeVPNGatewayConnectionDynamicRouteModePrototype{
				Psk:             core.StringPtr("lkj14b1oi0alcniejkso"),
				Peer:            vpnGatewayConnectionDynamicRouteModePeerPrototypeModel,
				RoutingProtocol: core.StringPtr("bgp"),
				Tunnels:         []vpcv1.VPNGatewayConnectionTunnelPrototype{*vpnGatewayConnectionTunnelPrototypeModel},
			}

			createVPNGatewayConnectionOptions := vpcService.NewCreateVPNGatewayConnectionOptions(
				"testString",
				vpnGatewayConnectionPrototypeModel,
			)

			vpnGatewayConnection, response, err := vpcService.CreateVPNGatewayConnection(createVPNGatewayConnectionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayConnection, "", "  ")
			fmt.Println(string(b))

			// end-create_vpn_gateway_connection

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
		It(`GetVPNGatewayConnection request example`, func() {
			fmt.Println("\nGetVPNGatewayConnection() result:")
			// begin-get_vpn_gateway_connection

			getVPNGatewayConnectionOptions := vpcService.NewGetVPNGatewayConnectionOptions(
				"testString",
				"testString",
			)

			vpnGatewayConnection, response, err := vpcService.GetVPNGatewayConnection(getVPNGatewayConnectionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayConnection, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_gateway_connection

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
		It(`UpdateVPNGatewayConnection request example`, func() {
			fmt.Println("\nUpdateVPNGatewayConnection() result:")
			// begin-update_vpn_gateway_connection

			vpnGatewayConnectionPatchModel := &vpcv1.VPNGatewayConnectionPatch{}
			vpnGatewayConnectionPatchModelAsPatch, asPatchErr := vpnGatewayConnectionPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNGatewayConnectionOptions := vpcService.NewUpdateVPNGatewayConnectionOptions(
				"testString",
				"testString",
				vpnGatewayConnectionPatchModelAsPatch,
			)
			updateVPNGatewayConnectionOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			vpnGatewayConnection, response, err := vpcService.UpdateVPNGatewayConnection(updateVPNGatewayConnectionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayConnection, "", "  ")
			fmt.Println(string(b))

			// end-update_vpn_gateway_connection

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
		It(`ListVPNGatewayConnectionsLocalCidrs request example`, func() {
			fmt.Println("\nListVPNGatewayConnectionsLocalCidrs() result:")
			// begin-list_vpn_gateway_connections_local_cidrs

			listVPNGatewayConnectionsLocalCidrsOptions := vpcService.NewListVPNGatewayConnectionsLocalCidrsOptions(
				"testString",
				"testString",
			)

			vpnGatewayConnectionCidRs, response, err := vpcService.ListVPNGatewayConnectionsLocalCidrs(listVPNGatewayConnectionsLocalCidrsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayConnectionCidRs, "", "  ")
			fmt.Println(string(b))

			// end-list_vpn_gateway_connections_local_cidrs

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnectionCidRs).ToNot(BeNil())
		})
		It(`CheckVPNGatewayConnectionsLocalCIDR request example`, func() {
			// begin-check_vpn_gateway_connections_local_cidr

			checkVPNGatewayConnectionsLocalCIDROptions := vpcService.NewCheckVPNGatewayConnectionsLocalCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.CheckVPNGatewayConnectionsLocalCIDR(checkVPNGatewayConnectionsLocalCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from CheckVPNGatewayConnectionsLocalCIDR(): %d\n", response.StatusCode)
			}

			// end-check_vpn_gateway_connections_local_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`AddVPNGatewayConnectionsLocalCIDR request example`, func() {
			// begin-add_vpn_gateway_connections_local_cidr

			addVPNGatewayConnectionsLocalCIDROptions := vpcService.NewAddVPNGatewayConnectionsLocalCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.AddVPNGatewayConnectionsLocalCIDR(addVPNGatewayConnectionsLocalCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 201 {
				fmt.Printf("\nUnexpected response status code received from AddVPNGatewayConnectionsLocalCIDR(): %d\n", response.StatusCode)
			}

			// end-add_vpn_gateway_connections_local_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
		})
		It(`ListVPNGatewayConnectionsPeerCidrs request example`, func() {
			fmt.Println("\nListVPNGatewayConnectionsPeerCidrs() result:")
			// begin-list_vpn_gateway_connections_peer_cidrs

			listVPNGatewayConnectionsPeerCidrsOptions := vpcService.NewListVPNGatewayConnectionsPeerCidrsOptions(
				"testString",
				"testString",
			)

			vpnGatewayConnectionCidRs, response, err := vpcService.ListVPNGatewayConnectionsPeerCidrs(listVPNGatewayConnectionsPeerCidrsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayConnectionCidRs, "", "  ")
			fmt.Println(string(b))

			// end-list_vpn_gateway_connections_peer_cidrs

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnectionCidRs).ToNot(BeNil())
		})
		It(`CheckVPNGatewayConnectionsPeerCIDR request example`, func() {
			// begin-check_vpn_gateway_connections_peer_cidr

			checkVPNGatewayConnectionsPeerCIDROptions := vpcService.NewCheckVPNGatewayConnectionsPeerCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.CheckVPNGatewayConnectionsPeerCIDR(checkVPNGatewayConnectionsPeerCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from CheckVPNGatewayConnectionsPeerCIDR(): %d\n", response.StatusCode)
			}

			// end-check_vpn_gateway_connections_peer_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`AddVPNGatewayConnectionsPeerCIDR request example`, func() {
			// begin-add_vpn_gateway_connections_peer_cidr

			addVPNGatewayConnectionsPeerCIDROptions := vpcService.NewAddVPNGatewayConnectionsPeerCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.AddVPNGatewayConnectionsPeerCIDR(addVPNGatewayConnectionsPeerCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 201 {
				fmt.Printf("\nUnexpected response status code received from AddVPNGatewayConnectionsPeerCIDR(): %d\n", response.StatusCode)
			}

			// end-add_vpn_gateway_connections_peer_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
		})
		It(`ListVPNGatewayMembers request example`, func() {
			fmt.Println("\nListVPNGatewayMembers() result:")
			// begin-list_vpn_gateway_members

			listVPNGatewayMembersOptions := vpcService.NewListVPNGatewayMembersOptions(
				"testString",
			)

			vpnGatewayMemberCollection, response, err := vpcService.ListVPNGatewayMembers(listVPNGatewayMembersOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayMemberCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_vpn_gateway_members

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayMemberCollection).ToNot(BeNil())
		})
		It(`GetVPNGatewayMember request example`, func() {
			fmt.Println("\nGetVPNGatewayMember() result:")
			// begin-get_vpn_gateway_member

			getVPNGatewayMemberOptions := vpcService.NewGetVPNGatewayMemberOptions(
				"testString",
				"testString",
			)

			vpnGatewayMember, response, err := vpcService.GetVPNGatewayMember(getVPNGatewayMemberOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayMember, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_gateway_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayMember).ToNot(BeNil())
		})
		It(`ReplaceVPNGatewayMember request example`, func() {
			fmt.Println("\nReplaceVPNGatewayMember() result:")
			// begin-replace_vpn_gateway_member

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			vpnGatewayMemberPrivateIPPrototypeModel := &vpcv1.VPNGatewayMemberPrivateIPPrototypeReservedIPPrototypeVPNGatewayContext{
				Subnet: subnetIdentityModel,
			}

			replaceVPNGatewayMemberOptions := vpcService.NewReplaceVPNGatewayMemberOptions(
				"testString",
				"testString",
				vpnGatewayMemberPrivateIPPrototypeModel,
			)

			vpnGatewayMember, response, err := vpcService.ReplaceVPNGatewayMember(replaceVPNGatewayMemberOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayMember, "", "  ")
			fmt.Println(string(b))

			// end-replace_vpn_gateway_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(vpnGatewayMember).ToNot(BeNil())
		})
		It(`ListVPNGatewayServiceConnections request example`, func() {
			fmt.Println("\nListVPNGatewayServiceConnections() result:")
			// begin-list_vpn_gateway_service_connections
			listVPNGatewayServiceConnectionsOptions := &vpcv1.ListVPNGatewayServiceConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
			}

			pager, err := vpcService.NewVPNGatewayServiceConnectionsPager(listVPNGatewayServiceConnectionsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNGatewayServiceConnection
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_gateway_service_connections
		})
		It(`GetVPNGatewayServiceConnection request example`, func() {
			fmt.Println("\nGetVPNGatewayServiceConnection() result:")
			// begin-get_vpn_gateway_service_connection

			getVPNGatewayServiceConnectionOptions := vpcService.NewGetVPNGatewayServiceConnectionOptions(
				"testString",
				"testString",
			)

			vpnGatewayServiceConnection, response, err := vpcService.GetVPNGatewayServiceConnection(getVPNGatewayServiceConnectionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnGatewayServiceConnection, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_gateway_service_connection

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayServiceConnection).ToNot(BeNil())
		})
		It(`ListVPNServers request example`, func() {
			fmt.Println("\nListVPNServers() result:")
			// begin-list_vpn_servers
			listVPNServersOptions := &vpcv1.ListVPNServersOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			pager, err := vpcService.NewVPNServersPager(listVPNServersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNServer
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_servers
		})
		It(`CreateVPNServer request example`, func() {
			fmt.Println("\nCreateVPNServer() result:")
			// begin-create_vpn_server

			vpnServerPrototypeCertificateModel := &vpcv1.VPNServerPrototypeCertificateCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationByCertificatePrototypeClientCaModel := &vpcv1.VPNServerAuthenticationByCertificatePrototypeClientCaCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationPrototypeModel := &vpcv1.VPNServerAuthenticationPrototypeVPNServerAuthenticationByCertificatePrototype{
				Method:   core.StringPtr("certificate"),
				ClientCa: vpnServerAuthenticationByCertificatePrototypeClientCaModel,
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			createVPNServerOptions := vpcService.NewCreateVPNServerOptions(
				vpnServerPrototypeCertificateModel,
				[]vpcv1.VPNServerAuthenticationPrototypeIntf{vpnServerAuthenticationPrototypeModel},
				"172.16.0.0/16",
				[]vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			)

			vpnServer, response, err := vpcService.CreateVPNServer(createVPNServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServer, "", "  ")
			fmt.Println(string(b))

			// end-create_vpn_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnServer).ToNot(BeNil())
		})
		It(`GetVPNServer request example`, func() {
			fmt.Println("\nGetVPNServer() result:")
			// begin-get_vpn_server

			getVPNServerOptions := vpcService.NewGetVPNServerOptions(
				"testString",
			)

			vpnServer, response, err := vpcService.GetVPNServer(getVPNServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServer, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServer).ToNot(BeNil())
		})
		It(`UpdateVPNServer request example`, func() {
			fmt.Println("\nUpdateVPNServer() result:")
			// begin-update_vpn_server

			vpnServerPatchModel := &vpcv1.VPNServerPatch{
				Name: core.StringPtr("my-vpn-server-updated"),
			}
			vpnServerPatchModelAsPatch, asPatchErr := vpnServerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNServerOptions := vpcService.NewUpdateVPNServerOptions(
				"testString",
				vpnServerPatchModelAsPatch,
			)
			updateVPNServerOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			vpnServer, response, err := vpcService.UpdateVPNServer(updateVPNServerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServer, "", "  ")
			fmt.Println(string(b))

			// end-update_vpn_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServer).ToNot(BeNil())
		})
		It(`GetVPNServerClientConfiguration request example`, func() {
			fmt.Println("\nGetVPNServerClientConfiguration() result:")
			// begin-get_vpn_server_client_configuration

			getVPNServerClientConfigurationOptions := vpcService.NewGetVPNServerClientConfigurationOptions(
				"testString",
			)

			vpnServerClientConfiguration, response, err := vpcService.GetVPNServerClientConfiguration(getVPNServerClientConfigurationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServerClientConfiguration, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_server_client_configuration

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerClientConfiguration).ToNot(BeNil())
		})
		It(`ListVPNServerClients request example`, func() {
			fmt.Println("\nListVPNServerClients() result:")
			// begin-list_vpn_server_clients
			listVPNServerClientsOptions := &vpcv1.ListVPNServerClientsOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("created_at"),
			}

			pager, err := vpcService.NewVPNServerClientsPager(listVPNServerClientsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNServerClient
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_server_clients
		})
		It(`GetVPNServerClient request example`, func() {
			fmt.Println("\nGetVPNServerClient() result:")
			// begin-get_vpn_server_client

			getVPNServerClientOptions := vpcService.NewGetVPNServerClientOptions(
				"testString",
				"testString",
			)

			vpnServerClient, response, err := vpcService.GetVPNServerClient(getVPNServerClientOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServerClient, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_server_client

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerClient).ToNot(BeNil())
		})
		It(`DisconnectVPNClient request example`, func() {
			// begin-disconnect_vpn_client

			disconnectVPNClientOptions := vpcService.NewDisconnectVPNClientOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DisconnectVPNClient(disconnectVPNClientOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DisconnectVPNClient(): %d\n", response.StatusCode)
			}

			// end-disconnect_vpn_client

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`ListVPNServerRoutes request example`, func() {
			fmt.Println("\nListVPNServerRoutes() result:")
			// begin-list_vpn_server_routes
			listVPNServerRoutesOptions := &vpcv1.ListVPNServerRoutesOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("name"),
			}

			pager, err := vpcService.NewVPNServerRoutesPager(listVPNServerRoutesOptions)
			if err != nil {
				panic(err)
			}

			var allResults []vpcv1.VPNServerRoute
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_vpn_server_routes
		})
		It(`CreateVPNServerRoute request example`, func() {
			fmt.Println("\nCreateVPNServerRoute() result:")
			// begin-create_vpn_server_route

			createVPNServerRouteOptions := vpcService.NewCreateVPNServerRouteOptions(
				"testString",
				"172.16.0.0/16",
			)

			vpnServerRoute, response, err := vpcService.CreateVPNServerRoute(createVPNServerRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServerRoute, "", "  ")
			fmt.Println(string(b))

			// end-create_vpn_server_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
		It(`GetVPNServerRoute request example`, func() {
			fmt.Println("\nGetVPNServerRoute() result:")
			// begin-get_vpn_server_route

			getVPNServerRouteOptions := vpcService.NewGetVPNServerRouteOptions(
				"testString",
				"testString",
			)

			vpnServerRoute, response, err := vpcService.GetVPNServerRoute(getVPNServerRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServerRoute, "", "  ")
			fmt.Println(string(b))

			// end-get_vpn_server_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
		It(`UpdateVPNServerRoute request example`, func() {
			fmt.Println("\nUpdateVPNServerRoute() result:")
			// begin-update_vpn_server_route

			vpnServerRoutePatchModel := &vpcv1.VPNServerRoutePatch{
				Name: core.StringPtr("my-vpn-server-route-updated"),
			}
			vpnServerRoutePatchModelAsPatch, asPatchErr := vpnServerRoutePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNServerRouteOptions := vpcService.NewUpdateVPNServerRouteOptions(
				"testString",
				"testString",
				vpnServerRoutePatchModelAsPatch,
			)

			vpnServerRoute, response, err := vpcService.UpdateVPNServerRoute(updateVPNServerRouteOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpnServerRoute, "", "  ")
			fmt.Println(string(b))

			// end-update_vpn_server_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
		It(`DeleteBackupPolicyPlan request example`, func() {
			fmt.Println("\nDeleteBackupPolicyPlan() result:")
			// begin-delete_backup_policy_plan

			deleteBackupPolicyPlanOptions := vpcService.NewDeleteBackupPolicyPlanOptions(
				"testString",
				"testString",
			)
			deleteBackupPolicyPlanOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			backupPolicyPlan, response, err := vpcService.DeleteBackupPolicyPlan(deleteBackupPolicyPlanOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicyPlan, "", "  ")
			fmt.Println(string(b))

			// end-delete_backup_policy_plan

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
		It(`DeleteBackupPolicy request example`, func() {
			fmt.Println("\nDeleteBackupPolicy() result:")
			// begin-delete_backup_policy

			deleteBackupPolicyOptions := vpcService.NewDeleteBackupPolicyOptions(
				"testString",
			)
			deleteBackupPolicyOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			backupPolicy, response, err := vpcService.DeleteBackupPolicy(deleteBackupPolicyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(backupPolicy, "", "  ")
			fmt.Println(string(b))

			// end-delete_backup_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(backupPolicy).ToNot(BeNil())
		})
		It(`DeleteBareMetalServerNetworkAttachment request example`, func() {
			// begin-delete_bare_metal_server_network_attachment

			deleteBareMetalServerNetworkAttachmentOptions := vpcService.NewDeleteBareMetalServerNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteBareMetalServerNetworkAttachment(deleteBareMetalServerNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteBareMetalServerNetworkAttachment(): %d\n", response.StatusCode)
			}

			// end-delete_bare_metal_server_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteBareMetalServerNetworkInterface request example`, func() {
			// begin-delete_bare_metal_server_network_interface

			deleteBareMetalServerNetworkInterfaceOptions := vpcService.NewDeleteBareMetalServerNetworkInterfaceOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteBareMetalServerNetworkInterface(deleteBareMetalServerNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteBareMetalServerNetworkInterface(): %d\n", response.StatusCode)
			}

			// end-delete_bare_metal_server_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RemoveBareMetalServerNetworkInterfaceFloatingIP request example`, func() {
			// begin-remove_bare_metal_server_network_interface_floating_ip

			removeBareMetalServerNetworkInterfaceFloatingIPOptions := vpcService.NewRemoveBareMetalServerNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveBareMetalServerNetworkInterfaceFloatingIP(removeBareMetalServerNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveBareMetalServerNetworkInterfaceFloatingIP(): %d\n", response.StatusCode)
			}

			// end-remove_bare_metal_server_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteBareMetalServer request example`, func() {
			// begin-delete_bare_metal_server

			deleteBareMetalServerOptions := vpcService.NewDeleteBareMetalServerOptions(
				"testString",
			)

			response, err := vpcService.DeleteBareMetalServer(deleteBareMetalServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteBareMetalServer(): %d\n", response.StatusCode)
			}

			// end-delete_bare_metal_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteClusterNetworkInterface request example`, func() {
			fmt.Println("\nDeleteClusterNetworkInterface() result:")
			// begin-delete_cluster_network_interface

			deleteClusterNetworkInterfaceOptions := vpcService.NewDeleteClusterNetworkInterfaceOptions(
				"testString",
				"testString",
			)
			deleteClusterNetworkInterfaceOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkInterface, response, err := vpcService.DeleteClusterNetworkInterface(deleteClusterNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-delete_cluster_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
		It(`DeleteClusterNetworkSubnetReservedIP request example`, func() {
			fmt.Println("\nDeleteClusterNetworkSubnetReservedIP() result:")
			// begin-delete_cluster_network_subnet_reserved_ip

			deleteClusterNetworkSubnetReservedIPOptions := vpcService.NewDeleteClusterNetworkSubnetReservedIPOptions(
				"testString",
				"testString",
				"testString",
			)
			deleteClusterNetworkSubnetReservedIPOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkSubnetReservedIP, response, err := vpcService.DeleteClusterNetworkSubnetReservedIP(deleteClusterNetworkSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnetReservedIP, "", "  ")
			fmt.Println(string(b))

			// end-delete_cluster_network_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
		It(`DeleteClusterNetworkSubnet request example`, func() {
			fmt.Println("\nDeleteClusterNetworkSubnet() result:")
			// begin-delete_cluster_network_subnet

			deleteClusterNetworkSubnetOptions := vpcService.NewDeleteClusterNetworkSubnetOptions(
				"testString",
				"testString",
			)
			deleteClusterNetworkSubnetOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetworkSubnet, response, err := vpcService.DeleteClusterNetworkSubnet(deleteClusterNetworkSubnetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetworkSubnet, "", "  ")
			fmt.Println(string(b))

			// end-delete_cluster_network_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
		It(`DeleteClusterNetwork request example`, func() {
			fmt.Println("\nDeleteClusterNetwork() result:")
			// begin-delete_cluster_network

			deleteClusterNetworkOptions := vpcService.NewDeleteClusterNetworkOptions(
				"testString",
			)
			deleteClusterNetworkOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			clusterNetwork, response, err := vpcService.DeleteClusterNetwork(deleteClusterNetworkOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(clusterNetwork, "", "  ")
			fmt.Println(string(b))

			// end-delete_cluster_network

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetwork).ToNot(BeNil())
		})
		It(`DeleteDedicatedHostGroup request example`, func() {
			// begin-delete_dedicated_host_group

			deleteDedicatedHostGroupOptions := vpcService.NewDeleteDedicatedHostGroupOptions(
				"testString",
			)

			response, err := vpcService.DeleteDedicatedHostGroup(deleteDedicatedHostGroupOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteDedicatedHostGroup(): %d\n", response.StatusCode)
			}

			// end-delete_dedicated_host_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteDedicatedHost request example`, func() {
			// begin-delete_dedicated_host

			deleteDedicatedHostOptions := vpcService.NewDeleteDedicatedHostOptions(
				"testString",
			)

			response, err := vpcService.DeleteDedicatedHost(deleteDedicatedHostOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteDedicatedHost(): %d\n", response.StatusCode)
			}

			// end-delete_dedicated_host

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RemoveEndpointGatewayIP request example`, func() {
			// begin-remove_endpoint_gateway_ip

			removeEndpointGatewayIPOptions := vpcService.NewRemoveEndpointGatewayIPOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveEndpointGatewayIP(removeEndpointGatewayIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveEndpointGatewayIP(): %d\n", response.StatusCode)
			}

			// end-remove_endpoint_gateway_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteEndpointGatewayResourceBinding request example`, func() {
			// begin-delete_endpoint_gateway_resource_binding

			deleteEndpointGatewayResourceBindingOptions := vpcService.NewDeleteEndpointGatewayResourceBindingOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteEndpointGatewayResourceBinding(deleteEndpointGatewayResourceBindingOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteEndpointGatewayResourceBinding(): %d\n", response.StatusCode)
			}

			// end-delete_endpoint_gateway_resource_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteEndpointGateway request example`, func() {
			// begin-delete_endpoint_gateway

			deleteEndpointGatewayOptions := vpcService.NewDeleteEndpointGatewayOptions(
				"testString",
			)

			response, err := vpcService.DeleteEndpointGateway(deleteEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteEndpointGateway(): %d\n", response.StatusCode)
			}

			// end-delete_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteFloatingIP request example`, func() {
			// begin-delete_floating_ip

			deleteFloatingIPOptions := vpcService.NewDeleteFloatingIPOptions(
				"testString",
			)

			response, err := vpcService.DeleteFloatingIP(deleteFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteFloatingIP(): %d\n", response.StatusCode)
			}

			// end-delete_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteFlowLogCollector request example`, func() {
			// begin-delete_flow_log_collector

			deleteFlowLogCollectorOptions := vpcService.NewDeleteFlowLogCollectorOptions(
				"testString",
			)

			response, err := vpcService.DeleteFlowLogCollector(deleteFlowLogCollectorOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteFlowLogCollector(): %d\n", response.StatusCode)
			}

			// end-delete_flow_log_collector

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteImage request example`, func() {
			// begin-delete_image

			deleteImageOptions := vpcService.NewDeleteImageOptions(
				"testString",
			)

			response, err := vpcService.DeleteImage(deleteImageOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteImage(): %d\n", response.StatusCode)
			}

			// end-delete_image

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteImageExportJob request example`, func() {
			// begin-delete_image_export_job

			deleteImageExportJobOptions := vpcService.NewDeleteImageExportJobOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteImageExportJob(deleteImageExportJobOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteImageExportJob(): %d\n", response.StatusCode)
			}

			// end-delete_image_export_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteInstanceGroup request example`, func() {
			// begin-delete_instance_group

			deleteInstanceGroupOptions := vpcService.NewDeleteInstanceGroupOptions(
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroup(deleteInstanceGroupOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroup(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupLoadBalancer request example`, func() {
			// begin-delete_instance_group_load_balancer

			deleteInstanceGroupLoadBalancerOptions := vpcService.NewDeleteInstanceGroupLoadBalancerOptions(
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupLoadBalancer(deleteInstanceGroupLoadBalancerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupLoadBalancer(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_load_balancer

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupManager request example`, func() {
			// begin-delete_instance_group_manager

			deleteInstanceGroupManagerOptions := vpcService.NewDeleteInstanceGroupManagerOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupManager(deleteInstanceGroupManagerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupManager(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupManagerAction request example`, func() {
			// begin-delete_instance_group_manager_action

			deleteInstanceGroupManagerActionOptions := vpcService.NewDeleteInstanceGroupManagerActionOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupManagerAction(deleteInstanceGroupManagerActionOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupManagerAction(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_manager_action

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupManagerPolicy request example`, func() {
			// begin-delete_instance_group_manager_policy

			deleteInstanceGroupManagerPolicyOptions := vpcService.NewDeleteInstanceGroupManagerPolicyOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupManagerPolicy(deleteInstanceGroupManagerPolicyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupManagerPolicy(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_manager_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupMemberships request example`, func() {
			// begin-delete_instance_group_memberships

			deleteInstanceGroupMembershipsOptions := vpcService.NewDeleteInstanceGroupMembershipsOptions(
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupMemberships(deleteInstanceGroupMembershipsOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupMemberships(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_memberships

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceGroupMembership request example`, func() {
			// begin-delete_instance_group_membership

			deleteInstanceGroupMembershipOptions := vpcService.NewDeleteInstanceGroupMembershipOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceGroupMembership(deleteInstanceGroupMembershipOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceGroupMembership(): %d\n", response.StatusCode)
			}

			// end-delete_instance_group_membership

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceTemplate request example`, func() {
			// begin-delete_instance_template

			deleteInstanceTemplateOptions := vpcService.NewDeleteInstanceTemplateOptions(
				"testString",
			)

			response, err := vpcService.DeleteInstanceTemplate(deleteInstanceTemplateOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceTemplate(): %d\n", response.StatusCode)
			}

			// end-delete_instance_template

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstance request example`, func() {
			// begin-delete_instance

			deleteInstanceOptions := vpcService.NewDeleteInstanceOptions(
				"testString",
			)
			deleteInstanceOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteInstance(deleteInstanceOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstance(): %d\n", response.StatusCode)
			}

			// end-delete_instance

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceClusterNetworkAttachment request example`, func() {
			fmt.Println("\nDeleteInstanceClusterNetworkAttachment() result:")
			// begin-delete_instance_cluster_network_attachment

			deleteInstanceClusterNetworkAttachmentOptions := vpcService.NewDeleteInstanceClusterNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			instanceClusterNetworkAttachment, response, err := vpcService.DeleteInstanceClusterNetworkAttachment(deleteInstanceClusterNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(instanceClusterNetworkAttachment, "", "  ")
			fmt.Println(string(b))

			// end-delete_instance_cluster_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
		It(`DeleteInstanceNetworkAttachment request example`, func() {
			// begin-delete_instance_network_attachment

			deleteInstanceNetworkAttachmentOptions := vpcService.NewDeleteInstanceNetworkAttachmentOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceNetworkAttachment(deleteInstanceNetworkAttachmentOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceNetworkAttachment(): %d\n", response.StatusCode)
			}

			// end-delete_instance_network_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteInstanceNetworkInterface request example`, func() {
			// begin-delete_instance_network_interface

			deleteInstanceNetworkInterfaceOptions := vpcService.NewDeleteInstanceNetworkInterfaceOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceNetworkInterface(deleteInstanceNetworkInterfaceOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceNetworkInterface(): %d\n", response.StatusCode)
			}

			// end-delete_instance_network_interface

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RemoveInstanceNetworkInterfaceFloatingIP request example`, func() {
			// begin-remove_instance_network_interface_floating_ip

			removeInstanceNetworkInterfaceFloatingIPOptions := vpcService.NewRemoveInstanceNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveInstanceNetworkInterfaceFloatingIP(removeInstanceNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveInstanceNetworkInterfaceFloatingIP(): %d\n", response.StatusCode)
			}

			// end-remove_instance_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteInstanceVolumeAttachment request example`, func() {
			// begin-delete_instance_volume_attachment

			deleteInstanceVolumeAttachmentOptions := vpcService.NewDeleteInstanceVolumeAttachmentOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteInstanceVolumeAttachment(deleteInstanceVolumeAttachmentOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteInstanceVolumeAttachment(): %d\n", response.StatusCode)
			}

			// end-delete_instance_volume_attachment

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteKey request example`, func() {
			// begin-delete_key

			deleteKeyOptions := vpcService.NewDeleteKeyOptions(
				"testString",
			)

			response, err := vpcService.DeleteKey(deleteKeyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteKey(): %d\n", response.StatusCode)
			}

			// end-delete_key

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteLoadBalancer request example`, func() {
			// begin-delete_load_balancer

			deleteLoadBalancerOptions := vpcService.NewDeleteLoadBalancerOptions(
				"testString",
			)
			deleteLoadBalancerOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteLoadBalancer(deleteLoadBalancerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancer(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteLoadBalancerListener request example`, func() {
			// begin-delete_load_balancer_listener

			deleteLoadBalancerListenerOptions := vpcService.NewDeleteLoadBalancerListenerOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancerListener(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer_listener

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteLoadBalancerListenerPolicy request example`, func() {
			// begin-delete_load_balancer_listener_policy

			deleteLoadBalancerListenerPolicyOptions := vpcService.NewDeleteLoadBalancerListenerPolicyOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteLoadBalancerListenerPolicy(deleteLoadBalancerListenerPolicyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancerListenerPolicy(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer_listener_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteLoadBalancerListenerPolicyRule request example`, func() {
			// begin-delete_load_balancer_listener_policy_rule

			deleteLoadBalancerListenerPolicyRuleOptions := vpcService.NewDeleteLoadBalancerListenerPolicyRuleOptions(
				"testString",
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteLoadBalancerListenerPolicyRule(deleteLoadBalancerListenerPolicyRuleOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancerListenerPolicyRule(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer_listener_policy_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteLoadBalancerPool request example`, func() {
			// begin-delete_load_balancer_pool

			deleteLoadBalancerPoolOptions := vpcService.NewDeleteLoadBalancerPoolOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancerPool(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer_pool

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteLoadBalancerPoolMember request example`, func() {
			// begin-delete_load_balancer_pool_member

			deleteLoadBalancerPoolMemberOptions := vpcService.NewDeleteLoadBalancerPoolMemberOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteLoadBalancerPoolMember(deleteLoadBalancerPoolMemberOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteLoadBalancerPoolMember(): %d\n", response.StatusCode)
			}

			// end-delete_load_balancer_pool_member

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteNetworkACL request example`, func() {
			// begin-delete_network_acl

			deleteNetworkACLOptions := vpcService.NewDeleteNetworkACLOptions(
				"testString",
			)

			response, err := vpcService.DeleteNetworkACL(deleteNetworkACLOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteNetworkACL(): %d\n", response.StatusCode)
			}

			// end-delete_network_acl

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteNetworkACLRule request example`, func() {
			// begin-delete_network_acl_rule

			deleteNetworkACLRuleOptions := vpcService.NewDeleteNetworkACLRuleOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteNetworkACLRule(deleteNetworkACLRuleOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteNetworkACLRule(): %d\n", response.StatusCode)
			}

			// end-delete_network_acl_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeletePlacementGroup request example`, func() {
			// begin-delete_placement_group

			deletePlacementGroupOptions := vpcService.NewDeletePlacementGroupOptions(
				"testString",
			)

			response, err := vpcService.DeletePlacementGroup(deletePlacementGroupOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeletePlacementGroup(): %d\n", response.StatusCode)
			}

			// end-delete_placement_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeletePrivatePathServiceGateway request example`, func() {
			// begin-delete_private_path_service_gateway

			deletePrivatePathServiceGatewayOptions := vpcService.NewDeletePrivatePathServiceGatewayOptions(
				"testString",
			)

			response, err := vpcService.DeletePrivatePathServiceGateway(deletePrivatePathServiceGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeletePrivatePathServiceGateway(): %d\n", response.StatusCode)
			}

			// end-delete_private_path_service_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeletePrivatePathServiceGatewayAccountPolicy request example`, func() {
			// begin-delete_private_path_service_gateway_account_policy

			deletePrivatePathServiceGatewayAccountPolicyOptions := vpcService.NewDeletePrivatePathServiceGatewayAccountPolicyOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeletePrivatePathServiceGatewayAccountPolicy(deletePrivatePathServiceGatewayAccountPolicyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeletePrivatePathServiceGatewayAccountPolicy(): %d\n", response.StatusCode)
			}

			// end-delete_private_path_service_gateway_account_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeletePublicAddressRange request example`, func() {
			fmt.Println("\nDeletePublicAddressRange() result:")
			// begin-delete_public_address_range

			deletePublicAddressRangeOptions := vpcService.NewDeletePublicAddressRangeOptions(
				"testString",
			)

			publicAddressRange, response, err := vpcService.DeletePublicAddressRange(deletePublicAddressRangeOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(publicAddressRange, "", "  ")
			fmt.Println(string(b))

			// end-delete_public_address_range

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(publicAddressRange).ToNot(BeNil())
		})
		It(`DeletePublicGateway request example`, func() {
			// begin-delete_public_gateway

			deletePublicGatewayOptions := vpcService.NewDeletePublicGatewayOptions(
				"testString",
			)

			response, err := vpcService.DeletePublicGateway(deletePublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeletePublicGateway(): %d\n", response.StatusCode)
			}

			// end-delete_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteReservation request example`, func() {
			fmt.Println("\nDeleteReservation() result:")
			// begin-delete_reservation

			deleteReservationOptions := vpcService.NewDeleteReservationOptions(
				"testString",
			)

			reservation, response, err := vpcService.DeleteReservation(deleteReservationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(reservation, "", "  ")
			fmt.Println(string(b))

			// end-delete_reservation

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(reservation).ToNot(BeNil())
		})
		It(`DeleteSecurityGroup request example`, func() {
			// begin-delete_security_group

			deleteSecurityGroupOptions := vpcService.NewDeleteSecurityGroupOptions(
				"testString",
			)

			response, err := vpcService.DeleteSecurityGroup(deleteSecurityGroupOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteSecurityGroup(): %d\n", response.StatusCode)
			}

			// end-delete_security_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteSecurityGroupRule request example`, func() {
			// begin-delete_security_group_rule

			deleteSecurityGroupRuleOptions := vpcService.NewDeleteSecurityGroupRuleOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteSecurityGroupRule(deleteSecurityGroupRuleOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteSecurityGroupRule(): %d\n", response.StatusCode)
			}

			// end-delete_security_group_rule

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteSecurityGroupTargetBinding request example`, func() {
			// begin-delete_security_group_target_binding

			deleteSecurityGroupTargetBindingOptions := vpcService.NewDeleteSecurityGroupTargetBindingOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteSecurityGroupTargetBinding(): %d\n", response.StatusCode)
			}

			// end-delete_security_group_target_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteShare request example`, func() {
			fmt.Println("\nDeleteShare() result:")
			// begin-delete_share

			deleteShareOptions := vpcService.NewDeleteShareOptions(
				"testString",
			)
			deleteShareOptions.SetIfMatch("W/96d225c4-56bd-43d9-98fc-d7148e5c5028")

			share, response, err := vpcService.DeleteShare(deleteShareOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(share, "", "  ")
			fmt.Println(string(b))

			// end-delete_share

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(share).ToNot(BeNil())
		})
		It(`DeleteShareAccessorBinding request example`, func() {
			// begin-delete_share_accessor_binding

			deleteShareAccessorBindingOptions := vpcService.NewDeleteShareAccessorBindingOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteShareAccessorBinding(deleteShareAccessorBindingOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteShareAccessorBinding(): %d\n", response.StatusCode)
			}

			// end-delete_share_accessor_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteShareMountTarget request example`, func() {
			fmt.Println("\nDeleteShareMountTarget() result:")
			// begin-delete_share_mount_target

			deleteShareMountTargetOptions := vpcService.NewDeleteShareMountTargetOptions(
				"testString",
				"testString",
			)

			shareMountTarget, response, err := vpcService.DeleteShareMountTarget(deleteShareMountTargetOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareMountTarget, "", "  ")
			fmt.Println(string(b))

			// end-delete_share_mount_target

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(shareMountTarget).ToNot(BeNil())
		})
		It(`DeleteShareSnapshot request example`, func() {
			fmt.Println("\nDeleteShareSnapshot() result:")
			// begin-delete_share_snapshot

			deleteShareSnapshotOptions := vpcService.NewDeleteShareSnapshotOptions(
				"testString",
				"testString",
			)

			shareSnapshot, response, err := vpcService.DeleteShareSnapshot(deleteShareSnapshotOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(shareSnapshot, "", "  ")
			fmt.Println(string(b))

			// end-delete_share_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(shareSnapshot).ToNot(BeNil())
		})
		It(`DeleteShareSource request example`, func() {
			// begin-delete_share_source

			deleteShareSourceOptions := vpcService.NewDeleteShareSourceOptions(
				"testString",
			)

			response, err := vpcService.DeleteShareSource(deleteShareSourceOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteShareSource(): %d\n", response.StatusCode)
			}

			// end-delete_share_source

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteSnapshotConsistencyGroup request example`, func() {
			fmt.Println("\nDeleteSnapshotConsistencyGroup() result:")
			// begin-delete_snapshot_consistency_group

			deleteSnapshotConsistencyGroupOptions := vpcService.NewDeleteSnapshotConsistencyGroupOptions(
				"testString",
			)

			snapshotConsistencyGroup, response, err := vpcService.DeleteSnapshotConsistencyGroup(deleteSnapshotConsistencyGroupOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(snapshotConsistencyGroup, "", "  ")
			fmt.Println(string(b))

			// end-delete_snapshot_consistency_group

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
		It(`DeleteSnapshots request example`, func() {
			// begin-delete_snapshots

			deleteSnapshotsOptions := vpcService.NewDeleteSnapshotsOptions(
				"testString",
			)

			response, err := vpcService.DeleteSnapshots(deleteSnapshotsOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteSnapshots(): %d\n", response.StatusCode)
			}

			// end-delete_snapshots

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteSnapshot request example`, func() {
			// begin-delete_snapshot

			deleteSnapshotOptions := vpcService.NewDeleteSnapshotOptions(
				"testString",
			)
			deleteSnapshotOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteSnapshot(deleteSnapshotOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteSnapshot(): %d\n", response.StatusCode)
			}

			// end-delete_snapshot

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteSnapshotClone request example`, func() {
			// begin-delete_snapshot_clone

			deleteSnapshotCloneOptions := vpcService.NewDeleteSnapshotCloneOptions(
				"testString",
				"us-south-1",
			)

			response, err := vpcService.DeleteSnapshotClone(deleteSnapshotCloneOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteSnapshotClone(): %d\n", response.StatusCode)
			}

			// end-delete_snapshot_clone

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteSubnet request example`, func() {
			// begin-delete_subnet

			deleteSubnetOptions := vpcService.NewDeleteSubnetOptions(
				"testString",
			)

			response, err := vpcService.DeleteSubnet(deleteSubnetOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteSubnet(): %d\n", response.StatusCode)
			}

			// end-delete_subnet

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`UnsetSubnetPublicGateway request example`, func() {
			// begin-unset_subnet_public_gateway

			unsetSubnetPublicGatewayOptions := vpcService.NewUnsetSubnetPublicGatewayOptions(
				"testString",
			)

			response, err := vpcService.UnsetSubnetPublicGateway(unsetSubnetPublicGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from UnsetSubnetPublicGateway(): %d\n", response.StatusCode)
			}

			// end-unset_subnet_public_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteSubnetReservedIP request example`, func() {
			// begin-delete_subnet_reserved_ip

			deleteSubnetReservedIPOptions := vpcService.NewDeleteSubnetReservedIPOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteSubnetReservedIP(deleteSubnetReservedIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteSubnetReservedIP(): %d\n", response.StatusCode)
			}

			// end-delete_subnet_reserved_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVirtualNetworkInterfaces request example`, func() {
			fmt.Println("\nDeleteVirtualNetworkInterfaces() result:")
			// begin-delete_virtual_network_interfaces

			deleteVirtualNetworkInterfacesOptions := vpcService.NewDeleteVirtualNetworkInterfacesOptions(
				"testString",
			)
			deleteVirtualNetworkInterfacesOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			virtualNetworkInterface, response, err := vpcService.DeleteVirtualNetworkInterfaces(deleteVirtualNetworkInterfacesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualNetworkInterface, "", "  ")
			fmt.Println(string(b))

			// end-delete_virtual_network_interfaces

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
		It(`RemoveNetworkInterfaceFloatingIP request example`, func() {
			// begin-remove_network_interface_floating_ip

			removeNetworkInterfaceFloatingIPOptions := vpcService.NewRemoveNetworkInterfaceFloatingIPOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveNetworkInterfaceFloatingIP(removeNetworkInterfaceFloatingIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveNetworkInterfaceFloatingIP(): %d\n", response.StatusCode)
			}

			// end-remove_network_interface_floating_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RemoveVirtualNetworkInterfaceIP request example`, func() {
			// begin-remove_virtual_network_interface_ip

			removeVirtualNetworkInterfaceIPOptions := vpcService.NewRemoveVirtualNetworkInterfaceIPOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveVirtualNetworkInterfaceIP(): %d\n", response.StatusCode)
			}

			// end-remove_virtual_network_interface_ip

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVolume request example`, func() {
			// begin-delete_volume

			deleteVolumeOptions := vpcService.NewDeleteVolumeOptions(
				"testString",
			)
			deleteVolumeOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteVolume(deleteVolumeOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVolume(): %d\n", response.StatusCode)
			}

			// end-delete_volume

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVolumeJob request example`, func() {
			// begin-delete_volume_job

			deleteVolumeJobOptions := vpcService.NewDeleteVolumeJobOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVolumeJob(deleteVolumeJobOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVolumeJob(): %d\n", response.StatusCode)
			}

			// end-delete_volume_job

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteVPC request example`, func() {
			// begin-delete_vpc

			deleteVPCOptions := vpcService.NewDeleteVPCOptions(
				"testString",
			)
			deleteVPCOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteVPC(deleteVPCOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPC(): %d\n", response.StatusCode)
			}

			// end-delete_vpc

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPCAddressPrefix request example`, func() {
			// begin-delete_vpc_address_prefix

			deleteVPCAddressPrefixOptions := vpcService.NewDeleteVPCAddressPrefixOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVPCAddressPrefix(deleteVPCAddressPrefixOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPCAddressPrefix(): %d\n", response.StatusCode)
			}

			// end-delete_vpc_address_prefix

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPCDnsResolutionBinding request example`, func() {
			fmt.Println("\nDeleteVPCDnsResolutionBinding() result:")
			// begin-delete_vpc_dns_resolution_binding

			deleteVPCDnsResolutionBindingOptions := vpcService.NewDeleteVPCDnsResolutionBindingOptions(
				"testString",
				"testString",
			)

			vpcdnsResolutionBinding, response, err := vpcService.DeleteVPCDnsResolutionBinding(deleteVPCDnsResolutionBindingOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(vpcdnsResolutionBinding, "", "  ")
			fmt.Println(string(b))

			// end-delete_vpc_dns_resolution_binding

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
		It(`DeleteVPCRoute request example`, func() {
			// begin-delete_vpc_route

			deleteVPCRouteOptions := vpcService.NewDeleteVPCRouteOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVPCRoute(deleteVPCRouteOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPCRoute(): %d\n", response.StatusCode)
			}

			// end-delete_vpc_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPCRoutingTable request example`, func() {
			// begin-delete_vpc_routing_table

			deleteVPCRoutingTableOptions := vpcService.NewDeleteVPCRoutingTableOptions(
				"testString",
				"testString",
			)
			deleteVPCRoutingTableOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteVPCRoutingTable(deleteVPCRoutingTableOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPCRoutingTable(): %d\n", response.StatusCode)
			}

			// end-delete_vpc_routing_table

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPCRoutingTableRoute request example`, func() {
			// begin-delete_vpc_routing_table_route

			deleteVPCRoutingTableRouteOptions := vpcService.NewDeleteVPCRoutingTableRouteOptions(
				"testString",
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVPCRoutingTableRoute(deleteVPCRoutingTableRouteOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPCRoutingTableRoute(): %d\n", response.StatusCode)
			}

			// end-delete_vpc_routing_table_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteIkePolicy request example`, func() {
			// begin-delete_ike_policy

			deleteIkePolicyOptions := vpcService.NewDeleteIkePolicyOptions(
				"testString",
			)

			response, err := vpcService.DeleteIkePolicy(deleteIkePolicyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteIkePolicy(): %d\n", response.StatusCode)
			}

			// end-delete_ike_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteIpsecPolicy request example`, func() {
			// begin-delete_ipsec_policy

			deleteIpsecPolicyOptions := vpcService.NewDeleteIpsecPolicyOptions(
				"testString",
			)

			response, err := vpcService.DeleteIpsecPolicy(deleteIpsecPolicyOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteIpsecPolicy(): %d\n", response.StatusCode)
			}

			// end-delete_ipsec_policy

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPNGateway request example`, func() {
			// begin-delete_vpn_gateway

			deleteVPNGatewayOptions := vpcService.NewDeleteVPNGatewayOptions(
				"testString",
			)

			response, err := vpcService.DeleteVPNGateway(deleteVPNGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPNGateway(): %d\n", response.StatusCode)
			}

			// end-delete_vpn_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`RemoveVPNGatewayAdvertisedCIDR request example`, func() {
			// begin-remove_vpn_gateway_advertised_cidr

			removeVPNGatewayAdvertisedCIDROptions := vpcService.NewRemoveVPNGatewayAdvertisedCIDROptions(
				"192.168.1.0/24",
				"testString",
			)

			response, err := vpcService.RemoveVPNGatewayAdvertisedCIDR(removeVPNGatewayAdvertisedCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveVPNGatewayAdvertisedCIDR(): %d\n", response.StatusCode)
			}

			// end-remove_vpn_gateway_advertised_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPNGatewayConnection request example`, func() {
			// begin-delete_vpn_gateway_connection

			deleteVPNGatewayConnectionOptions := vpcService.NewDeleteVPNGatewayConnectionOptions(
				"testString",
				"testString",
			)
			deleteVPNGatewayConnectionOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteVPNGatewayConnection(deleteVPNGatewayConnectionOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPNGatewayConnection(): %d\n", response.StatusCode)
			}

			// end-delete_vpn_gateway_connection

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`RemoveVPNGatewayConnectionsLocalCIDR request example`, func() {
			// begin-remove_vpn_gateway_connections_local_cidr

			removeVPNGatewayConnectionsLocalCIDROptions := vpcService.NewRemoveVPNGatewayConnectionsLocalCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveVPNGatewayConnectionsLocalCIDR(removeVPNGatewayConnectionsLocalCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveVPNGatewayConnectionsLocalCIDR(): %d\n", response.StatusCode)
			}

			// end-remove_vpn_gateway_connections_local_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`RemoveVPNGatewayConnectionsPeerCIDR request example`, func() {
			// begin-remove_vpn_gateway_connections_peer_cidr

			removeVPNGatewayConnectionsPeerCIDROptions := vpcService.NewRemoveVPNGatewayConnectionsPeerCIDROptions(
				"192.168.1.0/24",
				"testString",
				"testString",
			)

			response, err := vpcService.RemoveVPNGatewayConnectionsPeerCIDR(removeVPNGatewayConnectionsPeerCIDROptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from RemoveVPNGatewayConnectionsPeerCIDR(): %d\n", response.StatusCode)
			}

			// end-remove_vpn_gateway_connections_peer_cidr

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteVPNServer request example`, func() {
			// begin-delete_vpn_server

			deleteVPNServerOptions := vpcService.NewDeleteVPNServerOptions(
				"testString",
			)
			deleteVPNServerOptions.SetIfMatch("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\"")

			response, err := vpcService.DeleteVPNServer(deleteVPNServerOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPNServer(): %d\n", response.StatusCode)
			}

			// end-delete_vpn_server

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteVPNServerClient request example`, func() {
			// begin-delete_vpn_server_client

			deleteVPNServerClientOptions := vpcService.NewDeleteVPNServerClientOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVPNServerClient(deleteVPNServerClientOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPNServerClient(): %d\n", response.StatusCode)
			}

			// end-delete_vpn_server_client

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
		It(`DeleteVPNServerRoute request example`, func() {
			// begin-delete_vpn_server_route

			deleteVPNServerRouteOptions := vpcService.NewDeleteVPNServerRouteOptions(
				"testString",
				"testString",
			)

			response, err := vpcService.DeleteVPNServerRoute(deleteVPNServerRouteOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 202 {
				fmt.Printf("\nUnexpected response status code received from DeleteVPNServerRoute(): %d\n", response.StatusCode)
			}

			// end-delete_vpn_server_route

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})
})
