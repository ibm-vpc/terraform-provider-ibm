//go:build integration

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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/**
 * This file contains an integration test for the vpcv1 package.
 *
 * Notes:
 *
 * The integration test will automatically skip tests if the required config file is not available.
 */

var _ = Describe(`VpcV1 Integration Tests`, func() {
	const externalConfigFile = "../vpc_v1.env"

	var (
		err        error
		vpcService *vpcv1.VpcV1
		serviceURL string
		config     map[string]string
	)

	var shouldSkipTest = func() {
		Skip("External configuration is not available, skipping tests...")
	}

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping tests: " + err.Error())
			}

			os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(vpcv1.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping tests: " + err.Error())
			}
			serviceURL = config["URL"]
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}

			fmt.Fprintf(GinkgoWriter, "Service URL: %v\n", serviceURL)
			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			vpcServiceOptions := &vpcv1.VpcV1Options{
				Version: core.StringPtr("2024-06-23"),
			}

			vpcService, err = vpcv1.NewVpcV1UsingExternalConfig(vpcServiceOptions)
			Expect(err).To(BeNil())
			Expect(vpcService).ToNot(BeNil())
			Expect(vpcService.Service.Options.URL).To(Equal(serviceURL))

			core.SetLogger(core.NewLogger(core.LevelDebug, log.New(GinkgoWriter, "", log.LstdFlags), log.New(GinkgoWriter, "", log.LstdFlags)))
			vpcService.EnableRetries(4, 30*time.Second)
		})
	})

	Describe(`ListBackupPolicies - List backup policies`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBackupPolicies(listBackupPoliciesOptions *ListBackupPoliciesOptions) with pagination`, func() {
			listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
				Tag:             core.StringPtr("testString"),
			}

			listBackupPoliciesOptions.Start = nil
			listBackupPoliciesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BackupPolicyIntf
			for {
				backupPolicyCollection, response, err := vpcService.ListBackupPolicies(listBackupPoliciesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(backupPolicyCollection).ToNot(BeNil())
				allResults = append(allResults, backupPolicyCollection.BackupPolicies...)

				nextToken, nextErr := backupPolicyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBackupPoliciesOptions.Start = nextToken

				if listBackupPoliciesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBackupPolicies(listBackupPoliciesOptions *ListBackupPoliciesOptions) using BackupPoliciesPager`, func() {
			listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Tag:             core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewBackupPoliciesPager(listBackupPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BackupPolicyIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBackupPoliciesPager(listBackupPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBackupPolicies() returned a total of %d item(s) using BackupPoliciesPager.\n", len(allResults))
		})
	})

	Describe(`CreateBackupPolicy - Create a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBackupPolicy(createBackupPolicyOptions *CreateBackupPolicyOptions)`, func() {
			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			backupPolicyPlanClonePolicyPrototypeModel := &vpcv1.BackupPolicyPlanClonePolicyPrototype{
				MaxSnapshots: core.Int64Ptr(int64(5)),
				Zones:        []vpcv1.ZoneIdentityIntf{zoneIdentityModel},
			}

			backupPolicyPlanDeletionTriggerPrototypeModel := &vpcv1.BackupPolicyPlanDeletionTriggerPrototype{
				DeleteAfter:     core.Int64Ptr(int64(20)),
				DeleteOverCount: core.Int64Ptr(int64(20)),
			}

			encryptionKeyIdentityModel := &vpcv1.EncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			regionIdentityModel := &vpcv1.RegionIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"),
			}

			backupPolicyPlanRemoteRegionPolicyPrototypeModel := &vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{
				DeleteOverCount: core.Int64Ptr(int64(5)),
				EncryptionKey:   encryptionKeyIdentityModel,
				Region:          regionIdentityModel,
			}

			backupPolicyPlanPrototypeModel := &vpcv1.BackupPolicyPlanPrototype{
				Active:               core.BoolPtr(true),
				AttachUserTags:       []string{"my-daily-backup-plan"},
				ClonePolicy:          backupPolicyPlanClonePolicyPrototypeModel,
				CopyUserTags:         core.BoolPtr(true),
				CronSpec:             core.StringPtr("30 */2 * * 1-5"),
				DeletionTrigger:      backupPolicyPlanDeletionTriggerPrototypeModel,
				Name:                 core.StringPtr("my-policy-plan"),
				RemoteRegionPolicies: []vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{*backupPolicyPlanRemoteRegionPolicyPrototypeModel},
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			backupPolicyScopePrototypeModel := &vpcv1.BackupPolicyScopePrototypeEnterpriseIdentityEnterpriseIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:enterprise::a/aa2432b1fa4d4ace891e9b80fc104e34::enterprise:ebc2b430240943458b9e91e1432cfcce"),
			}

			backupPolicyPrototypeModel := &vpcv1.BackupPolicyPrototypeBackupPolicyMatchResourceTypeInstancePrototype{
				MatchUserTags:     []string{"my-daily-backup-policy"},
				Name:              core.StringPtr("my-backup-policy"),
				Plans:             []vpcv1.BackupPolicyPlanPrototype{*backupPolicyPlanPrototypeModel},
				ResourceGroup:     resourceGroupIdentityModel,
				Scope:             backupPolicyScopePrototypeModel,
				IncludedContent:   []string{"boot_volume", "data_volumes"},
				MatchResourceType: core.StringPtr("instance"),
			}

			createBackupPolicyOptions := &vpcv1.CreateBackupPolicyOptions{
				BackupPolicyPrototype: backupPolicyPrototypeModel,
			}

			backupPolicy, response, err := vpcService.CreateBackupPolicy(createBackupPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(backupPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListBackupPolicyJobs - List jobs for a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBackupPolicyJobs(listBackupPolicyJobsOptions *ListBackupPolicyJobsOptions) with pagination`, func() {
			listBackupPolicyJobsOptions := &vpcv1.ListBackupPolicyJobsOptions{
				BackupPolicyID:     core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Sort:               core.StringPtr("name"),
				SourceID:           core.StringPtr("testString"),
				Start:              core.StringPtr("testString"),
				Status:             core.StringPtr("failed"),
				TargetSnapshotsCRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot:r006-f6bfa329-0e36-433f-a3bb-0df632e79263"),
				TargetSnapshotsID:  core.StringPtr("testString"),
			}

			listBackupPolicyJobsOptions.Start = nil
			listBackupPolicyJobsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BackupPolicyJob
			for {
				backupPolicyJobCollection, response, err := vpcService.ListBackupPolicyJobs(listBackupPolicyJobsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(backupPolicyJobCollection).ToNot(BeNil())
				allResults = append(allResults, backupPolicyJobCollection.Jobs...)

				nextToken, nextErr := backupPolicyJobCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBackupPolicyJobsOptions.Start = nextToken

				if listBackupPolicyJobsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBackupPolicyJobs(listBackupPolicyJobsOptions *ListBackupPolicyJobsOptions) using BackupPolicyJobsPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewBackupPolicyJobsPager(listBackupPolicyJobsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BackupPolicyJob
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBackupPolicyJobsPager(listBackupPolicyJobsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBackupPolicyJobs() returned a total of %d item(s) using BackupPolicyJobsPager.\n", len(allResults))
		})
	})

	Describe(`GetBackupPolicyJob - Retrieve a backup policy job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBackupPolicyJob(getBackupPolicyJobOptions *GetBackupPolicyJobOptions)`, func() {
			getBackupPolicyJobOptions := &vpcv1.GetBackupPolicyJobOptions{
				BackupPolicyID: core.StringPtr("testString"),
				ID:             core.StringPtr("testString"),
			}

			backupPolicyJob, response, err := vpcService.GetBackupPolicyJob(getBackupPolicyJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyJob).ToNot(BeNil())
		})
	})

	Describe(`ListBackupPolicyPlans - List plans for a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBackupPolicyPlans(listBackupPolicyPlansOptions *ListBackupPolicyPlansOptions)`, func() {
			listBackupPolicyPlansOptions := &vpcv1.ListBackupPolicyPlansOptions{
				BackupPolicyID: core.StringPtr("testString"),
				Name:           core.StringPtr("my-name"),
			}

			backupPolicyPlanCollection, response, err := vpcService.ListBackupPolicyPlans(listBackupPolicyPlansOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlanCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateBackupPolicyPlan - Create a plan for a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBackupPolicyPlan(createBackupPolicyPlanOptions *CreateBackupPolicyPlanOptions)`, func() {
			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			backupPolicyPlanClonePolicyPrototypeModel := &vpcv1.BackupPolicyPlanClonePolicyPrototype{
				MaxSnapshots: core.Int64Ptr(int64(5)),
				Zones:        []vpcv1.ZoneIdentityIntf{zoneIdentityModel},
			}

			backupPolicyPlanDeletionTriggerPrototypeModel := &vpcv1.BackupPolicyPlanDeletionTriggerPrototype{
				DeleteAfter:     core.Int64Ptr(int64(20)),
				DeleteOverCount: core.Int64Ptr(int64(20)),
			}

			encryptionKeyIdentityModel := &vpcv1.EncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			regionIdentityModel := &vpcv1.RegionIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"),
			}

			backupPolicyPlanRemoteRegionPolicyPrototypeModel := &vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{
				DeleteOverCount: core.Int64Ptr(int64(5)),
				EncryptionKey:   encryptionKeyIdentityModel,
				Region:          regionIdentityModel,
			}

			createBackupPolicyPlanOptions := &vpcv1.CreateBackupPolicyPlanOptions{
				BackupPolicyID:       core.StringPtr("testString"),
				CronSpec:             core.StringPtr("30 */2 * * 1-5"),
				Active:               core.BoolPtr(true),
				AttachUserTags:       []string{"my-daily-backup-plan"},
				ClonePolicy:          backupPolicyPlanClonePolicyPrototypeModel,
				CopyUserTags:         core.BoolPtr(true),
				DeletionTrigger:      backupPolicyPlanDeletionTriggerPrototypeModel,
				Name:                 core.StringPtr("my-policy-plan"),
				RemoteRegionPolicies: []vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{*backupPolicyPlanRemoteRegionPolicyPrototypeModel},
			}

			backupPolicyPlan, response, err := vpcService.CreateBackupPolicyPlan(createBackupPolicyPlanOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
	})

	Describe(`GetBackupPolicyPlan - Retrieve a backup policy plan`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBackupPolicyPlan(getBackupPolicyPlanOptions *GetBackupPolicyPlanOptions)`, func() {
			getBackupPolicyPlanOptions := &vpcv1.GetBackupPolicyPlanOptions{
				BackupPolicyID: core.StringPtr("testString"),
				ID:             core.StringPtr("testString"),
			}

			backupPolicyPlan, response, err := vpcService.GetBackupPolicyPlan(getBackupPolicyPlanOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
	})

	Describe(`UpdateBackupPolicyPlan - Update a backup policy plan`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBackupPolicyPlan(updateBackupPolicyPlanOptions *UpdateBackupPolicyPlanOptions)`, func() {
			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			backupPolicyPlanClonePolicyPatchModel := &vpcv1.BackupPolicyPlanClonePolicyPatch{
				MaxSnapshots: core.Int64Ptr(int64(1)),
				Zones:        []vpcv1.ZoneIdentityIntf{zoneIdentityModel},
			}

			backupPolicyPlanDeletionTriggerPatchModel := &vpcv1.BackupPolicyPlanDeletionTriggerPatch{
				DeleteAfter:     core.Int64Ptr(int64(20)),
				DeleteOverCount: core.Int64Ptr(int64(1)),
			}

			encryptionKeyIdentityModel := &vpcv1.EncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			regionIdentityModel := &vpcv1.RegionIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"),
			}

			backupPolicyPlanRemoteRegionPolicyPrototypeModel := &vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{
				DeleteOverCount: core.Int64Ptr(int64(5)),
				EncryptionKey:   encryptionKeyIdentityModel,
				Region:          regionIdentityModel,
			}

			backupPolicyPlanPatchModel := &vpcv1.BackupPolicyPlanPatch{
				Active:               core.BoolPtr(true),
				AttachUserTags:       []string{"my-daily-backup-plan"},
				ClonePolicy:          backupPolicyPlanClonePolicyPatchModel,
				CopyUserTags:         core.BoolPtr(false),
				CronSpec:             core.StringPtr("30 */2 * * 1-5"),
				DeletionTrigger:      backupPolicyPlanDeletionTriggerPatchModel,
				Name:                 core.StringPtr("my-policy-plan"),
				RemoteRegionPolicies: []vpcv1.BackupPolicyPlanRemoteRegionPolicyPrototype{*backupPolicyPlanRemoteRegionPolicyPrototypeModel},
			}
			backupPolicyPlanPatchModelAsPatch, asPatchErr := backupPolicyPlanPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBackupPolicyPlanOptions := &vpcv1.UpdateBackupPolicyPlanOptions{
				BackupPolicyID:        core.StringPtr("testString"),
				ID:                    core.StringPtr("testString"),
				BackupPolicyPlanPatch: backupPolicyPlanPatchModelAsPatch,
				IfMatch:               core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			backupPolicyPlan, response, err := vpcService.UpdateBackupPolicyPlan(updateBackupPolicyPlanOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
	})

	Describe(`GetBackupPolicy - Retrieve a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBackupPolicy(getBackupPolicyOptions *GetBackupPolicyOptions)`, func() {
			getBackupPolicyOptions := &vpcv1.GetBackupPolicyOptions{
				ID: core.StringPtr("testString"),
			}

			backupPolicy, response, err := vpcService.GetBackupPolicy(getBackupPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdateBackupPolicy - Update a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBackupPolicy(updateBackupPolicyOptions *UpdateBackupPolicyOptions)`, func() {
			backupPolicyPatchModel := &vpcv1.BackupPolicyPatch{
				IncludedContent: []string{"data_volumes"},
				MatchUserTags:   []string{"my-tag-1", "my-tag-2", "my-tag-3"},
				Name:            core.StringPtr("my-backup-policy"),
			}
			backupPolicyPatchModelAsPatch, asPatchErr := backupPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBackupPolicyOptions := &vpcv1.UpdateBackupPolicyOptions{
				ID:                core.StringPtr("testString"),
				BackupPolicyPatch: backupPolicyPatchModelAsPatch,
				IfMatch:           core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			backupPolicy, response, err := vpcService.UpdateBackupPolicy(updateBackupPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(backupPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerCapacities - List capacities for bare metal servers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerCapacities(listBareMetalServerCapacitiesOptions *ListBareMetalServerCapacitiesOptions) with pagination`, func() {
			listBareMetalServerCapacitiesOptions := &vpcv1.ListBareMetalServerCapacitiesOptions{
				Limit:       core.Int64Ptr(int64(10)),
				ProfileName: core.StringPtr("bx2d-metal-192x768"),
				Start:       core.StringPtr("testString"),
				ZoneName:    core.StringPtr("us-south-1"),
			}

			listBareMetalServerCapacitiesOptions.Start = nil
			listBareMetalServerCapacitiesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServerCapacity
			for {
				bareMetalServerCapacityCollection, response, err := vpcService.ListBareMetalServerCapacities(listBareMetalServerCapacitiesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(bareMetalServerCapacityCollection).ToNot(BeNil())
				allResults = append(allResults, bareMetalServerCapacityCollection.Capacities...)

				nextToken, nextErr := bareMetalServerCapacityCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBareMetalServerCapacitiesOptions.Start = nextToken

				if listBareMetalServerCapacitiesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBareMetalServerCapacities(listBareMetalServerCapacitiesOptions *ListBareMetalServerCapacitiesOptions) using BareMetalServerCapacitiesPager`, func() {
			listBareMetalServerCapacitiesOptions := &vpcv1.ListBareMetalServerCapacitiesOptions{
				Limit:       core.Int64Ptr(int64(10)),
				ProfileName: core.StringPtr("bx2d-metal-192x768"),
				ZoneName:    core.StringPtr("us-south-1"),
			}

			// Test GetNext().
			pager, err := vpcService.NewBareMetalServerCapacitiesPager(listBareMetalServerCapacitiesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServerCapacity
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBareMetalServerCapacitiesPager(listBareMetalServerCapacitiesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBareMetalServerCapacities() returned a total of %d item(s) using BareMetalServerCapacitiesPager.\n", len(allResults))
		})
	})

	Describe(`ListBareMetalServerProfiles - List bare metal server profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerProfiles(listBareMetalServerProfilesOptions *ListBareMetalServerProfilesOptions) with pagination`, func() {
			listBareMetalServerProfilesOptions := &vpcv1.ListBareMetalServerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listBareMetalServerProfilesOptions.Start = nil
			listBareMetalServerProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServerProfile
			for {
				bareMetalServerProfileCollection, response, err := vpcService.ListBareMetalServerProfiles(listBareMetalServerProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(bareMetalServerProfileCollection).ToNot(BeNil())
				allResults = append(allResults, bareMetalServerProfileCollection.Profiles...)

				nextToken, nextErr := bareMetalServerProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBareMetalServerProfilesOptions.Start = nextToken

				if listBareMetalServerProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBareMetalServerProfiles(listBareMetalServerProfilesOptions *ListBareMetalServerProfilesOptions) using BareMetalServerProfilesPager`, func() {
			listBareMetalServerProfilesOptions := &vpcv1.ListBareMetalServerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewBareMetalServerProfilesPager(listBareMetalServerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServerProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBareMetalServerProfilesPager(listBareMetalServerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBareMetalServerProfiles() returned a total of %d item(s) using BareMetalServerProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetBareMetalServerProfile - Retrieve a bare metal server profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerProfile(getBareMetalServerProfileOptions *GetBareMetalServerProfileOptions)`, func() {
			getBareMetalServerProfileOptions := &vpcv1.GetBareMetalServerProfileOptions{
				Name: core.StringPtr("bx2-metal-192x768"),
			}

			bareMetalServerProfile, response, err := vpcService.GetBareMetalServerProfile(getBareMetalServerProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerProfile).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServers - List bare metal servers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServers(listBareMetalServersOptions *ListBareMetalServersOptions) with pagination`, func() {
			listBareMetalServersOptions := &vpcv1.ListBareMetalServersOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ReservationCRN:  core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
				ReservationID:   core.StringPtr("testString"),
				ReservationName: core.StringPtr("my-reservation"),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			listBareMetalServersOptions.Start = nil
			listBareMetalServersOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServer
			for {
				bareMetalServerCollection, response, err := vpcService.ListBareMetalServers(listBareMetalServersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(bareMetalServerCollection).ToNot(BeNil())
				allResults = append(allResults, bareMetalServerCollection.BareMetalServers...)

				nextToken, nextErr := bareMetalServerCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBareMetalServersOptions.Start = nextToken

				if listBareMetalServersOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBareMetalServers(listBareMetalServersOptions *ListBareMetalServersOptions) using BareMetalServersPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewBareMetalServersPager(listBareMetalServersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServer
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBareMetalServersPager(listBareMetalServersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBareMetalServers() returned a total of %d item(s) using BareMetalServersPager.\n", len(allResults))
		})
	})

	Describe(`CreateBareMetalServer - Create a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBareMetalServer(createBareMetalServerOptions *CreateBareMetalServerOptions)`, func() {
			bareMetalServerInitializationDefaultTrustedProfilePrototypeTargetModel := &vpcv1.BareMetalServerInitializationDefaultTrustedProfilePrototypeTargetTrustedProfileIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:iam-identity::a/aa2432b1fa4d4ace891e9b80fc104e34::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"),
			}

			bareMetalServerInitializationPrototypeDefaultTrustedProfileModel := &vpcv1.BareMetalServerInitializationPrototypeDefaultTrustedProfile{
				AutoLink: core.BoolPtr(false),
				Target:   bareMetalServerInitializationDefaultTrustedProfilePrototypeTargetModel,
			}

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			catalogOfferingVersionPlanIdentityModel := &vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:51c9e0db-2911-45a6-adb0-ac5332d27cf2:plan:sw.51c9e0db-2911-45a6-adb0-ac5332d27cf2.772c0dbe-aa62-482e-adbe-a3fc20101e0e"),
			}

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			bareMetalServerInitializationByCatalogOfferingCatalogOfferingModel := &vpcv1.BareMetalServerInitializationByCatalogOfferingCatalogOfferingBareMetalServerCatalogOfferingPrototypeCatalogOfferingByOffering{
				Plan:     catalogOfferingVersionPlanIdentityModel,
				Offering: catalogOfferingIdentityModel,
			}

			bareMetalServerInitializationPrototypeModel := &vpcv1.BareMetalServerInitializationPrototypeBareMetalServerInitializationByCatalogOffering{
				DefaultTrustedProfile: bareMetalServerInitializationPrototypeDefaultTrustedProfileModel,
				Keys:                  []vpcv1.KeyIdentityIntf{keyIdentityModel},
				UserData:              core.StringPtr("[...]"),
				CatalogOffering:       bareMetalServerInitializationByCatalogOfferingCatalogOfferingModel,
			}

			bareMetalServerMetadataServicePrototypeModel := &vpcv1.BareMetalServerMetadataServicePrototype{
				Enabled:  core.BoolPtr(true),
				Protocol: core.StringPtr("https"),
			}

			bareMetalServerPrototypeProfileModel := &vpcv1.BareMetalServerPrototypeProfileBareMetalServerProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/bare_metal_server/profiles/bx2-metal-192x768"),
			}

			reservationIdentityModel := &vpcv1.ReservationIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
			}

			bareMetalServerReservationAffinityPrototypeModel := &vpcv1.BareMetalServerReservationAffinityPrototype{
				Policy: core.StringPtr("automatic"),
				Pool:   []vpcv1.ReservationIdentityIntf{reservationIdentityModel},
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			bareMetalServerTrustedPlatformModulePrototypeModel := &vpcv1.BareMetalServerTrustedPlatformModulePrototype{
				Mode: core.StringPtr("disabled"),
			}

			bareMetalServerPrototypeVPCModel := &vpcv1.BareMetalServerPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			bareMetalServerPrototypeZoneModel := &vpcv1.BareMetalServerPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			bareMetalServerNetworkAttachmentPrototypeModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{
				Name:                    core.StringPtr("my-bare-metal-server-network-attachment"),
				VirtualNetworkInterface: bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
				AllowedVlans:            []int64{int64(4)},
				InterfaceType:           core.StringPtr("pci"),
			}

			bareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.BareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentBareMetalServerPrimaryNetworkAttachmentPrototypeBareMetalServerPrimaryNetworkAttachmentByPciPrototype{
				Name:                    core.StringPtr("my-bare-metal-server-network-attachment"),
				VirtualNetworkInterface: bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
				AllowedVlans:            []int64{int64(4)},
				InterfaceType:           core.StringPtr("pci"),
			}

			bareMetalServerPrototypeModel := &vpcv1.BareMetalServerPrototypeBareMetalServerByNetworkAttachment{
				Bandwidth:                core.Int64Ptr(int64(20000)),
				EnableSecureBoot:         core.BoolPtr(false),
				Initialization:           bareMetalServerInitializationPrototypeModel,
				MetadataService:          bareMetalServerMetadataServicePrototypeModel,
				Name:                     core.StringPtr("my-bare-metal-server"),
				Profile:                  bareMetalServerPrototypeProfileModel,
				ReservationAffinity:      bareMetalServerReservationAffinityPrototypeModel,
				ResourceGroup:            resourceGroupIdentityModel,
				TrustedPlatformModule:    bareMetalServerTrustedPlatformModulePrototypeModel,
				VPC:                      bareMetalServerPrototypeVPCModel,
				Zone:                     bareMetalServerPrototypeZoneModel,
				NetworkAttachments:       []vpcv1.BareMetalServerNetworkAttachmentPrototypeIntf{bareMetalServerNetworkAttachmentPrototypeModel},
				PrimaryNetworkAttachment: bareMetalServerByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createBareMetalServerOptions := &vpcv1.CreateBareMetalServerOptions{
				BareMetalServerPrototype: bareMetalServerPrototypeModel,
			}

			bareMetalServer, response, err := vpcService.CreateBareMetalServer(createBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServer).ToNot(BeNil())
		})
	})

	Describe(`CreateBareMetalServerConsoleAccessToken - Create a console access token for a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBareMetalServerConsoleAccessToken(createBareMetalServerConsoleAccessTokenOptions *CreateBareMetalServerConsoleAccessTokenOptions)`, func() {
			createBareMetalServerConsoleAccessTokenOptions := &vpcv1.CreateBareMetalServerConsoleAccessTokenOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ConsoleType:       core.StringPtr("serial"),
				Force:             core.BoolPtr(false),
			}

			bareMetalServerConsoleAccessToken, response, err := vpcService.CreateBareMetalServerConsoleAccessToken(createBareMetalServerConsoleAccessTokenOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerConsoleAccessToken).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerDisks - List disks on a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerDisks(listBareMetalServerDisksOptions *ListBareMetalServerDisksOptions)`, func() {
			listBareMetalServerDisksOptions := &vpcv1.ListBareMetalServerDisksOptions{
				BareMetalServerID: core.StringPtr("testString"),
			}

			bareMetalServerDiskCollection, response, err := vpcService.ListBareMetalServerDisks(listBareMetalServerDisksOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDiskCollection).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServerDisk - Retrieve a bare metal server disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerDisk(getBareMetalServerDiskOptions *GetBareMetalServerDiskOptions)`, func() {
			getBareMetalServerDiskOptions := &vpcv1.GetBareMetalServerDiskOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			bareMetalServerDisk, response, err := vpcService.GetBareMetalServerDisk(getBareMetalServerDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDisk).ToNot(BeNil())
		})
	})

	Describe(`UpdateBareMetalServerDisk - Update a bare metal server disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBareMetalServerDisk(updateBareMetalServerDiskOptions *UpdateBareMetalServerDiskOptions)`, func() {
			bareMetalServerDiskPatchModel := &vpcv1.BareMetalServerDiskPatch{
				Name: core.StringPtr("my-bare-metal-server-disk-updated"),
			}
			bareMetalServerDiskPatchModelAsPatch, asPatchErr := bareMetalServerDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerDiskOptions := &vpcv1.UpdateBareMetalServerDiskOptions{
				BareMetalServerID:        core.StringPtr("testString"),
				ID:                       core.StringPtr("testString"),
				BareMetalServerDiskPatch: bareMetalServerDiskPatchModelAsPatch,
			}

			bareMetalServerDisk, response, err := vpcService.UpdateBareMetalServerDisk(updateBareMetalServerDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerDisk).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerNetworkAttachments - List network attachments on a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerNetworkAttachments(listBareMetalServerNetworkAttachmentsOptions *ListBareMetalServerNetworkAttachmentsOptions) with pagination`, func() {
			listBareMetalServerNetworkAttachmentsOptions := &vpcv1.ListBareMetalServerNetworkAttachmentsOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Start:             core.StringPtr("testString"),
			}

			listBareMetalServerNetworkAttachmentsOptions.Start = nil
			listBareMetalServerNetworkAttachmentsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServerNetworkAttachmentIntf
			for {
				bareMetalServerNetworkAttachmentCollection, response, err := vpcService.ListBareMetalServerNetworkAttachments(listBareMetalServerNetworkAttachmentsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(bareMetalServerNetworkAttachmentCollection).ToNot(BeNil())
				allResults = append(allResults, bareMetalServerNetworkAttachmentCollection.NetworkAttachments...)

				nextToken, nextErr := bareMetalServerNetworkAttachmentCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBareMetalServerNetworkAttachmentsOptions.Start = nextToken

				if listBareMetalServerNetworkAttachmentsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBareMetalServerNetworkAttachments(listBareMetalServerNetworkAttachmentsOptions *ListBareMetalServerNetworkAttachmentsOptions) using BareMetalServerNetworkAttachmentsPager`, func() {
			listBareMetalServerNetworkAttachmentsOptions := &vpcv1.ListBareMetalServerNetworkAttachmentsOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewBareMetalServerNetworkAttachmentsPager(listBareMetalServerNetworkAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServerNetworkAttachmentIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBareMetalServerNetworkAttachmentsPager(listBareMetalServerNetworkAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBareMetalServerNetworkAttachments() returned a total of %d item(s) using BareMetalServerNetworkAttachmentsPager.\n", len(allResults))
		})
	})

	Describe(`CreateBareMetalServerNetworkAttachment - Create a network attachment on a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBareMetalServerNetworkAttachment(createBareMetalServerNetworkAttachmentOptions *CreateBareMetalServerNetworkAttachmentOptions)`, func() {
			bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			bareMetalServerNetworkAttachmentPrototypeModel := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{
				Name:                    core.StringPtr("my-bare-metal-server-network-attachment"),
				VirtualNetworkInterface: bareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
				AllowedVlans:            []int64{int64(4)},
				InterfaceType:           core.StringPtr("pci"),
			}

			createBareMetalServerNetworkAttachmentOptions := &vpcv1.CreateBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID:                         core.StringPtr("testString"),
				BareMetalServerNetworkAttachmentPrototype: bareMetalServerNetworkAttachmentPrototypeModel,
			}

			bareMetalServerNetworkAttachment, response, err := vpcService.CreateBareMetalServerNetworkAttachment(createBareMetalServerNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServerNetworkAttachment - Retrieve a bare metal server network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachmentOptions *GetBareMetalServerNetworkAttachmentOptions)`, func() {
			getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			bareMetalServerNetworkAttachment, response, err := vpcService.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateBareMetalServerNetworkAttachment - Update a bare metal server network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBareMetalServerNetworkAttachment(updateBareMetalServerNetworkAttachmentOptions *UpdateBareMetalServerNetworkAttachmentOptions)`, func() {
			bareMetalServerNetworkAttachmentPatchModel := &vpcv1.BareMetalServerNetworkAttachmentPatch{
				AllowedVlans: []int64{int64(4)},
				Name:         core.StringPtr("my-bare-metal-server-network-attachment-updated"),
			}
			bareMetalServerNetworkAttachmentPatchModelAsPatch, asPatchErr := bareMetalServerNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerNetworkAttachmentOptions := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID:                     core.StringPtr("testString"),
				ID:                                    core.StringPtr("testString"),
				BareMetalServerNetworkAttachmentPatch: bareMetalServerNetworkAttachmentPatchModelAsPatch,
			}

			bareMetalServerNetworkAttachment, response, err := vpcService.UpdateBareMetalServerNetworkAttachment(updateBareMetalServerNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerNetworkInterfaces - List network interfaces on a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerNetworkInterfaces(listBareMetalServerNetworkInterfacesOptions *ListBareMetalServerNetworkInterfacesOptions) with pagination`, func() {
			listBareMetalServerNetworkInterfacesOptions := &vpcv1.ListBareMetalServerNetworkInterfacesOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Start:             core.StringPtr("testString"),
			}

			listBareMetalServerNetworkInterfacesOptions.Start = nil
			listBareMetalServerNetworkInterfacesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServerNetworkInterfaceIntf
			for {
				bareMetalServerNetworkInterfaceCollection, response, err := vpcService.ListBareMetalServerNetworkInterfaces(listBareMetalServerNetworkInterfacesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(bareMetalServerNetworkInterfaceCollection).ToNot(BeNil())
				allResults = append(allResults, bareMetalServerNetworkInterfaceCollection.NetworkInterfaces...)

				nextToken, nextErr := bareMetalServerNetworkInterfaceCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listBareMetalServerNetworkInterfacesOptions.Start = nextToken

				if listBareMetalServerNetworkInterfacesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListBareMetalServerNetworkInterfaces(listBareMetalServerNetworkInterfacesOptions *ListBareMetalServerNetworkInterfacesOptions) using BareMetalServerNetworkInterfacesPager`, func() {
			listBareMetalServerNetworkInterfacesOptions := &vpcv1.ListBareMetalServerNetworkInterfacesOptions{
				BareMetalServerID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewBareMetalServerNetworkInterfacesPager(listBareMetalServerNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServerNetworkInterfaceIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewBareMetalServerNetworkInterfacesPager(listBareMetalServerNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListBareMetalServerNetworkInterfaces() returned a total of %d item(s) using BareMetalServerNetworkInterfacesPager.\n", len(allResults))
		})
	})

	Describe(`CreateBareMetalServerNetworkInterface - Create a network interface on a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateBareMetalServerNetworkInterface(createBareMetalServerNetworkInterfaceOptions *CreateBareMetalServerNetworkInterfaceOptions)`, func() {
			networkInterfaceIPPrototypeModel := &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			bareMetalServerNetworkInterfacePrototypeSubnetModel := &vpcv1.BareMetalServerNetworkInterfacePrototypeSubnetSubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			bareMetalServerNetworkInterfacePrototypeModel := &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByHiperSocketPrototype{
				AllowIPSpoofing:         core.BoolPtr(true),
				EnableInfrastructureNat: core.BoolPtr(true),
				Name:                    core.StringPtr("my-bare-metal-server-network-interface"),
				PrimaryIP:               networkInterfaceIPPrototypeModel,
				SecurityGroups:          []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
				Subnet:                  bareMetalServerNetworkInterfacePrototypeSubnetModel,
				InterfaceType:           core.StringPtr("hipersocket"),
			}

			createBareMetalServerNetworkInterfaceOptions := &vpcv1.CreateBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID:                        core.StringPtr("testString"),
				BareMetalServerNetworkInterfacePrototype: bareMetalServerNetworkInterfacePrototypeModel,
			}

			bareMetalServerNetworkInterface, response, err := vpcService.CreateBareMetalServerNetworkInterface(createBareMetalServerNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServerNetworkInterface - Retrieve a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerNetworkInterface(getBareMetalServerNetworkInterfaceOptions *GetBareMetalServerNetworkInterfaceOptions)`, func() {
			getBareMetalServerNetworkInterfaceOptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			bareMetalServerNetworkInterface, response, err := vpcService.GetBareMetalServerNetworkInterface(getBareMetalServerNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`UpdateBareMetalServerNetworkInterface - Update a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBareMetalServerNetworkInterface(updateBareMetalServerNetworkInterfaceOptions *UpdateBareMetalServerNetworkInterfaceOptions)`, func() {
			bareMetalServerNetworkInterfacePatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{
				AllowIPSpoofing:         core.BoolPtr(true),
				AllowedVlans:            []int64{int64(4)},
				EnableInfrastructureNat: core.BoolPtr(true),
				Name:                    core.StringPtr("my-bare-metal-server-network-interface"),
			}
			bareMetalServerNetworkInterfacePatchModelAsPatch, asPatchErr := bareMetalServerNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerNetworkInterfaceOptions := &vpcv1.UpdateBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID:                    core.StringPtr("testString"),
				ID:                                   core.StringPtr("testString"),
				BareMetalServerNetworkInterfacePatch: bareMetalServerNetworkInterfacePatchModelAsPatch,
			}

			bareMetalServerNetworkInterface, response, err := vpcService.UpdateBareMetalServerNetworkInterface(updateBareMetalServerNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerNetworkInterfaceFloatingIps - List floating IPs associated with a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerNetworkInterfaceFloatingIps(listBareMetalServerNetworkInterfaceFloatingIpsOptions *ListBareMetalServerNetworkInterfaceFloatingIpsOptions)`, func() {
			listBareMetalServerNetworkInterfaceFloatingIpsOptions := &vpcv1.ListBareMetalServerNetworkInterfaceFloatingIpsOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIPUnpaginatedCollection, response, err := vpcService.ListBareMetalServerNetworkInterfaceFloatingIps(listBareMetalServerNetworkInterfaceFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPUnpaginatedCollection).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServerNetworkInterfaceFloatingIP - Retrieve associated floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerNetworkInterfaceFloatingIP(getBareMetalServerNetworkInterfaceFloatingIPOptions *GetBareMetalServerNetworkInterfaceFloatingIPOptions)`, func() {
			getBareMetalServerNetworkInterfaceFloatingIPOptions := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				ID:                 core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIP, response, err := vpcService.GetBareMetalServerNetworkInterfaceFloatingIP(getBareMetalServerNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`AddBareMetalServerNetworkInterfaceFloatingIP - Associate a floating IP with a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddBareMetalServerNetworkInterfaceFloatingIP(addBareMetalServerNetworkInterfaceFloatingIPOptions *AddBareMetalServerNetworkInterfaceFloatingIPOptions)`, func() {
			addBareMetalServerNetworkInterfaceFloatingIPOptions := &vpcv1.AddBareMetalServerNetworkInterfaceFloatingIPOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				ID:                 core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIP, response, err := vpcService.AddBareMetalServerNetworkInterfaceFloatingIP(addBareMetalServerNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`ListBareMetalServerNetworkInterfaceIps - List the primary reserved IP for a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListBareMetalServerNetworkInterfaceIps(listBareMetalServerNetworkInterfaceIpsOptions *ListBareMetalServerNetworkInterfaceIpsOptions)`, func() {
			listBareMetalServerNetworkInterfaceIpsOptions := &vpcv1.ListBareMetalServerNetworkInterfaceIpsOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			reservedIPCollectionBareMetalServerNetworkInterfaceContext, response, err := vpcService.ListBareMetalServerNetworkInterfaceIps(listBareMetalServerNetworkInterfaceIpsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPCollectionBareMetalServerNetworkInterfaceContext).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServerNetworkInterfaceIP - Retrieve the primary reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerNetworkInterfaceIP(getBareMetalServerNetworkInterfaceIPOptions *GetBareMetalServerNetworkInterfaceIPOptions)`, func() {
			getBareMetalServerNetworkInterfaceIPOptions := &vpcv1.GetBareMetalServerNetworkInterfaceIPOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				ID:                 core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			reservedIP, response, err := vpcService.GetBareMetalServerNetworkInterfaceIP(getBareMetalServerNetworkInterfaceIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`GetBareMetalServer - Retrieve a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServer(getBareMetalServerOptions *GetBareMetalServerOptions)`, func() {
			getBareMetalServerOptions := &vpcv1.GetBareMetalServerOptions{
				ID: core.StringPtr("testString"),
			}

			bareMetalServer, response, err := vpcService.GetBareMetalServer(getBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServer).ToNot(BeNil())
		})
	})

	Describe(`UpdateBareMetalServer - Update a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateBareMetalServer(updateBareMetalServerOptions *UpdateBareMetalServerOptions)`, func() {
			bareMetalServerMetadataServicePatchModel := &vpcv1.BareMetalServerMetadataServicePatch{
				Enabled:  core.BoolPtr(true),
				Protocol: core.StringPtr("http"),
			}

			reservationIdentityModel := &vpcv1.ReservationIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
			}

			bareMetalServerReservationAffinityPatchModel := &vpcv1.BareMetalServerReservationAffinityPatch{
				Policy: core.StringPtr("automatic"),
				Pool:   []vpcv1.ReservationIdentityIntf{reservationIdentityModel},
			}

			bareMetalServerTrustedPlatformModulePatchModel := &vpcv1.BareMetalServerTrustedPlatformModulePatch{
				Mode: core.StringPtr("disabled"),
			}

			bareMetalServerPatchModel := &vpcv1.BareMetalServerPatch{
				Bandwidth:             core.Int64Ptr(int64(20000)),
				EnableSecureBoot:      core.BoolPtr(false),
				MetadataService:       bareMetalServerMetadataServicePatchModel,
				Name:                  core.StringPtr("my-bare-metal-server"),
				ReservationAffinity:   bareMetalServerReservationAffinityPatchModel,
				TrustedPlatformModule: bareMetalServerTrustedPlatformModulePatchModel,
			}
			bareMetalServerPatchModelAsPatch, asPatchErr := bareMetalServerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateBareMetalServerOptions := &vpcv1.UpdateBareMetalServerOptions{
				ID:                   core.StringPtr("testString"),
				BareMetalServerPatch: bareMetalServerPatchModelAsPatch,
			}

			bareMetalServer, response, err := vpcService.UpdateBareMetalServer(updateBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServer).ToNot(BeNil())
		})
	})

	Describe(`UpdateFirmwareForBareMetalServer - Update firmware for a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateFirmwareForBareMetalServer(updateFirmwareForBareMetalServerOptions *UpdateFirmwareForBareMetalServerOptions)`, func() {
			updateFirmwareForBareMetalServerOptions := &vpcv1.UpdateFirmwareForBareMetalServerOptions{
				ID:        core.StringPtr("testString"),
				AutoStart: core.BoolPtr(true),
			}

			response, err := vpcService.UpdateFirmwareForBareMetalServer(updateFirmwareForBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`GetBareMetalServerInitialization - Retrieve initialization configuration for a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetBareMetalServerInitialization(getBareMetalServerInitializationOptions *GetBareMetalServerInitializationOptions)`, func() {
			getBareMetalServerInitializationOptions := &vpcv1.GetBareMetalServerInitializationOptions{
				ID: core.StringPtr("testString"),
			}

			bareMetalServerInitialization, response, err := vpcService.GetBareMetalServerInitialization(getBareMetalServerInitializationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerInitialization).ToNot(BeNil())
		})
	})

	Describe(`ReplaceBareMetalServerInitialization - Reinitialize a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceBareMetalServerInitialization(replaceBareMetalServerInitializationOptions *ReplaceBareMetalServerInitializationOptions)`, func() {
			bareMetalServerInitializationDefaultTrustedProfilePrototypeTargetModel := &vpcv1.BareMetalServerInitializationDefaultTrustedProfilePrototypeTargetTrustedProfileIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:iam-identity::a/aa2432b1fa4d4ace891e9b80fc104e34::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"),
			}

			bareMetalServerInitializationReplacePrototypeDefaultTrustedProfileModel := &vpcv1.BareMetalServerInitializationReplacePrototypeDefaultTrustedProfile{
				AutoLink: core.BoolPtr(false),
				Target:   bareMetalServerInitializationDefaultTrustedProfilePrototypeTargetModel,
			}

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			bareMetalServerInitializationByImageImageModel := &vpcv1.BareMetalServerInitializationByImageImageImageIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"),
			}

			bareMetalServerInitializationReplacePrototypeModel := &vpcv1.BareMetalServerInitializationReplacePrototypeBareMetalServerInitializationByImage{
				DefaultTrustedProfile: bareMetalServerInitializationReplacePrototypeDefaultTrustedProfileModel,
				Keys:                  []vpcv1.KeyIdentityIntf{keyIdentityModel},
				UserData:              core.StringPtr("[...]"),
				Image:                 bareMetalServerInitializationByImageImageModel,
			}

			replaceBareMetalServerInitializationOptions := &vpcv1.ReplaceBareMetalServerInitializationOptions{
				ID: core.StringPtr("testString"),
				BareMetalServerInitializationReplacePrototype: bareMetalServerInitializationReplacePrototypeModel,
			}

			bareMetalServerInitialization, response, err := vpcService.ReplaceBareMetalServerInitialization(replaceBareMetalServerInitializationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(bareMetalServerInitialization).ToNot(BeNil())
		})
	})

	Describe(`RestartBareMetalServer - Restart a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RestartBareMetalServer(restartBareMetalServerOptions *RestartBareMetalServerOptions)`, func() {
			restartBareMetalServerOptions := &vpcv1.RestartBareMetalServerOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.RestartBareMetalServer(restartBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`StartBareMetalServer - Start a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`StartBareMetalServer(startBareMetalServerOptions *StartBareMetalServerOptions)`, func() {
			startBareMetalServerOptions := &vpcv1.StartBareMetalServerOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.StartBareMetalServer(startBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`StopBareMetalServer - Stop a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`StopBareMetalServer(stopBareMetalServerOptions *StopBareMetalServerOptions)`, func() {
			stopBareMetalServerOptions := &vpcv1.StopBareMetalServerOptions{
				ID:   core.StringPtr("testString"),
				Type: core.StringPtr("hard"),
			}

			response, err := vpcService.StopBareMetalServer(stopBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`ListClusterNetworkProfiles - List cluster network profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListClusterNetworkProfiles(listClusterNetworkProfilesOptions *ListClusterNetworkProfilesOptions) with pagination`, func() {
			listClusterNetworkProfilesOptions := &vpcv1.ListClusterNetworkProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listClusterNetworkProfilesOptions.Start = nil
			listClusterNetworkProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ClusterNetworkProfile
			for {
				clusterNetworkProfileCollection, response, err := vpcService.ListClusterNetworkProfiles(listClusterNetworkProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(clusterNetworkProfileCollection).ToNot(BeNil())
				allResults = append(allResults, clusterNetworkProfileCollection.Profiles...)

				nextToken, nextErr := clusterNetworkProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listClusterNetworkProfilesOptions.Start = nextToken

				if listClusterNetworkProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListClusterNetworkProfiles(listClusterNetworkProfilesOptions *ListClusterNetworkProfilesOptions) using ClusterNetworkProfilesPager`, func() {
			listClusterNetworkProfilesOptions := &vpcv1.ListClusterNetworkProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewClusterNetworkProfilesPager(listClusterNetworkProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ClusterNetworkProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewClusterNetworkProfilesPager(listClusterNetworkProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListClusterNetworkProfiles() returned a total of %d item(s) using ClusterNetworkProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetClusterNetworkProfile - Retrieve a cluster network profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetClusterNetworkProfile(getClusterNetworkProfileOptions *GetClusterNetworkProfileOptions)`, func() {
			getClusterNetworkProfileOptions := &vpcv1.GetClusterNetworkProfileOptions{
				Name: core.StringPtr("h100"),
			}

			clusterNetworkProfile, response, err := vpcService.GetClusterNetworkProfile(getClusterNetworkProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkProfile).ToNot(BeNil())
		})
	})

	Describe(`ListClusterNetworks - List cluster networks`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListClusterNetworks(listClusterNetworksOptions *ListClusterNetworksOptions) with pagination`, func() {
			listClusterNetworksOptions := &vpcv1.ListClusterNetworksOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				Start:           core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			listClusterNetworksOptions.Start = nil
			listClusterNetworksOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ClusterNetwork
			for {
				clusterNetworkCollection, response, err := vpcService.ListClusterNetworks(listClusterNetworksOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(clusterNetworkCollection).ToNot(BeNil())
				allResults = append(allResults, clusterNetworkCollection.ClusterNetworks...)

				nextToken, nextErr := clusterNetworkCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listClusterNetworksOptions.Start = nextToken

				if listClusterNetworksOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListClusterNetworks(listClusterNetworksOptions *ListClusterNetworksOptions) using ClusterNetworksPager`, func() {
			listClusterNetworksOptions := &vpcv1.ListClusterNetworksOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			// Test GetNext().
			pager, err := vpcService.NewClusterNetworksPager(listClusterNetworksOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ClusterNetwork
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewClusterNetworksPager(listClusterNetworksOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListClusterNetworks() returned a total of %d item(s) using ClusterNetworksPager.\n", len(allResults))
		})
	})

	Describe(`CreateClusterNetwork - Create a cluster network`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateClusterNetwork(createClusterNetworkOptions *CreateClusterNetworkOptions)`, func() {
			clusterNetworkPrototypeProfileModel := &vpcv1.ClusterNetworkPrototypeProfileClusterNetworkProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"),
			}

			clusterNetworkPrototypeVPCModel := &vpcv1.ClusterNetworkPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			clusterNetworkPrototypeZoneModel := &vpcv1.ClusterNetworkPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			clusterNetworkSubnetPrefixPrototypeModel := &vpcv1.ClusterNetworkSubnetPrefixPrototype{
				CIDR: core.StringPtr("10.0.0.0/24"),
			}

			createClusterNetworkOptions := &vpcv1.CreateClusterNetworkOptions{
				Profile:        clusterNetworkPrototypeProfileModel,
				VPC:            clusterNetworkPrototypeVPCModel,
				Zone:           clusterNetworkPrototypeZoneModel,
				Name:           core.StringPtr("my-cluster-network"),
				ResourceGroup:  resourceGroupIdentityModel,
				SubnetPrefixes: []vpcv1.ClusterNetworkSubnetPrefixPrototype{*clusterNetworkSubnetPrefixPrototypeModel},
			}

			clusterNetwork, response, err := vpcService.CreateClusterNetwork(createClusterNetworkOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetwork).ToNot(BeNil())
		})
	})

	Describe(`ListClusterNetworkInterfaces - List cluster network interfaces`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListClusterNetworkInterfaces(listClusterNetworkInterfacesOptions *ListClusterNetworkInterfacesOptions) with pagination`, func() {
			listClusterNetworkInterfacesOptions := &vpcv1.ListClusterNetworkInterfacesOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
				Start:            core.StringPtr("testString"),
			}

			listClusterNetworkInterfacesOptions.Start = nil
			listClusterNetworkInterfacesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ClusterNetworkInterface
			for {
				clusterNetworkInterfaceCollection, response, err := vpcService.ListClusterNetworkInterfaces(listClusterNetworkInterfacesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(clusterNetworkInterfaceCollection).ToNot(BeNil())
				allResults = append(allResults, clusterNetworkInterfaceCollection.Interfaces...)

				nextToken, nextErr := clusterNetworkInterfaceCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listClusterNetworkInterfacesOptions.Start = nextToken

				if listClusterNetworkInterfacesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListClusterNetworkInterfaces(listClusterNetworkInterfacesOptions *ListClusterNetworkInterfacesOptions) using ClusterNetworkInterfacesPager`, func() {
			listClusterNetworkInterfacesOptions := &vpcv1.ListClusterNetworkInterfacesOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewClusterNetworkInterfacesPager(listClusterNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ClusterNetworkInterface
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewClusterNetworkInterfacesPager(listClusterNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListClusterNetworkInterfaces() returned a total of %d item(s) using ClusterNetworkInterfacesPager.\n", len(allResults))
		})
	})

	Describe(`CreateClusterNetworkInterface - Create a cluster network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateClusterNetworkInterface(createClusterNetworkInterfaceOptions *CreateClusterNetworkInterfaceOptions)`, func() {
			clusterNetworkInterfacePrimaryIPPrototypeModel := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/0717-d4d6489a-3bf5-4104-a33a-3572faf2d117"),
			}

			clusterNetworkInterfaceSubnetPrototypeModel := &vpcv1.ClusterNetworkInterfaceSubnetPrototypeClusterNetworkSubnetIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"),
			}

			createClusterNetworkInterfaceOptions := &vpcv1.CreateClusterNetworkInterfaceOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Name:             core.StringPtr("my-cluster-network-interface"),
				PrimaryIP:        clusterNetworkInterfacePrimaryIPPrototypeModel,
				Subnet:           clusterNetworkInterfaceSubnetPrototypeModel,
			}

			clusterNetworkInterface, response, err := vpcService.CreateClusterNetworkInterface(createClusterNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`GetClusterNetworkInterface - Retrieve a cluster network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetClusterNetworkInterface(getClusterNetworkInterfaceOptions *GetClusterNetworkInterfaceOptions)`, func() {
			getClusterNetworkInterfaceOptions := &vpcv1.GetClusterNetworkInterfaceOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				ID:               core.StringPtr("testString"),
			}

			clusterNetworkInterface, response, err := vpcService.GetClusterNetworkInterface(getClusterNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`UpdateClusterNetworkInterface - Update a cluster network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateClusterNetworkInterface(updateClusterNetworkInterfaceOptions *UpdateClusterNetworkInterfaceOptions)`, func() {
			clusterNetworkInterfacePatchModel := &vpcv1.ClusterNetworkInterfacePatch{
				AutoDelete: core.BoolPtr(false),
				Name:       core.StringPtr("my-cluster-network-interface"),
			}
			clusterNetworkInterfacePatchModelAsPatch, asPatchErr := clusterNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkInterfaceOptions := &vpcv1.UpdateClusterNetworkInterfaceOptions{
				ClusterNetworkID:             core.StringPtr("testString"),
				ID:                           core.StringPtr("testString"),
				ClusterNetworkInterfacePatch: clusterNetworkInterfacePatchModelAsPatch,
				IfMatch:                      core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkInterface, response, err := vpcService.UpdateClusterNetworkInterface(updateClusterNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`ListClusterNetworkSubnets - List cluster network subnets`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListClusterNetworkSubnets(listClusterNetworkSubnetsOptions *ListClusterNetworkSubnetsOptions) with pagination`, func() {
			listClusterNetworkSubnetsOptions := &vpcv1.ListClusterNetworkSubnetsOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
				Start:            core.StringPtr("testString"),
			}

			listClusterNetworkSubnetsOptions.Start = nil
			listClusterNetworkSubnetsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ClusterNetworkSubnet
			for {
				clusterNetworkSubnetCollection, response, err := vpcService.ListClusterNetworkSubnets(listClusterNetworkSubnetsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(clusterNetworkSubnetCollection).ToNot(BeNil())
				allResults = append(allResults, clusterNetworkSubnetCollection.Subnets...)

				nextToken, nextErr := clusterNetworkSubnetCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listClusterNetworkSubnetsOptions.Start = nextToken

				if listClusterNetworkSubnetsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListClusterNetworkSubnets(listClusterNetworkSubnetsOptions *ListClusterNetworkSubnetsOptions) using ClusterNetworkSubnetsPager`, func() {
			listClusterNetworkSubnetsOptions := &vpcv1.ListClusterNetworkSubnetsOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				Limit:            core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-name"),
				Sort:             core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewClusterNetworkSubnetsPager(listClusterNetworkSubnetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ClusterNetworkSubnet
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewClusterNetworkSubnetsPager(listClusterNetworkSubnetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListClusterNetworkSubnets() returned a total of %d item(s) using ClusterNetworkSubnetsPager.\n", len(allResults))
		})
	})

	Describe(`CreateClusterNetworkSubnet - Create a cluster network subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateClusterNetworkSubnet(createClusterNetworkSubnetOptions *CreateClusterNetworkSubnetOptions)`, func() {
			clusterNetworkSubnetPrototypeModel := &vpcv1.ClusterNetworkSubnetPrototypeClusterNetworkSubnetByIPv4CIDRBlockPrototype{
				IPVersion:     core.StringPtr("ipv4"),
				Name:          core.StringPtr("my-cluster-network-subnet"),
				Ipv4CIDRBlock: core.StringPtr("10.0.0.0/24"),
			}

			createClusterNetworkSubnetOptions := &vpcv1.CreateClusterNetworkSubnetOptions{
				ClusterNetworkID:              core.StringPtr("testString"),
				ClusterNetworkSubnetPrototype: clusterNetworkSubnetPrototypeModel,
			}

			clusterNetworkSubnet, response, err := vpcService.CreateClusterNetworkSubnet(createClusterNetworkSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
	})

	Describe(`ListClusterNetworkSubnetReservedIps - List cluster network subnet reserved IPs`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListClusterNetworkSubnetReservedIps(listClusterNetworkSubnetReservedIpsOptions *ListClusterNetworkSubnetReservedIpsOptions) with pagination`, func() {
			listClusterNetworkSubnetReservedIpsOptions := &vpcv1.ListClusterNetworkSubnetReservedIpsOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
				Name:                   core.StringPtr("my-name"),
				Sort:                   core.StringPtr("name"),
				Start:                  core.StringPtr("testString"),
			}

			listClusterNetworkSubnetReservedIpsOptions.Start = nil
			listClusterNetworkSubnetReservedIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ClusterNetworkSubnetReservedIP
			for {
				clusterNetworkSubnetReservedIPCollection, response, err := vpcService.ListClusterNetworkSubnetReservedIps(listClusterNetworkSubnetReservedIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(clusterNetworkSubnetReservedIPCollection).ToNot(BeNil())
				allResults = append(allResults, clusterNetworkSubnetReservedIPCollection.ReservedIps...)

				nextToken, nextErr := clusterNetworkSubnetReservedIPCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listClusterNetworkSubnetReservedIpsOptions.Start = nextToken

				if listClusterNetworkSubnetReservedIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListClusterNetworkSubnetReservedIps(listClusterNetworkSubnetReservedIpsOptions *ListClusterNetworkSubnetReservedIpsOptions) using ClusterNetworkSubnetReservedIpsPager`, func() {
			listClusterNetworkSubnetReservedIpsOptions := &vpcv1.ListClusterNetworkSubnetReservedIpsOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
				Name:                   core.StringPtr("my-name"),
				Sort:                   core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewClusterNetworkSubnetReservedIpsPager(listClusterNetworkSubnetReservedIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ClusterNetworkSubnetReservedIP
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewClusterNetworkSubnetReservedIpsPager(listClusterNetworkSubnetReservedIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListClusterNetworkSubnetReservedIps() returned a total of %d item(s) using ClusterNetworkSubnetReservedIpsPager.\n", len(allResults))
		})
	})

	Describe(`CreateClusterNetworkSubnetReservedIP - Create a cluster network subnet reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateClusterNetworkSubnetReservedIP(createClusterNetworkSubnetReservedIPOptions *CreateClusterNetworkSubnetReservedIPOptions)`, func() {
			createClusterNetworkSubnetReservedIPOptions := &vpcv1.CreateClusterNetworkSubnetReservedIPOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				Address:                core.StringPtr("192.168.3.4"),
				Name:                   core.StringPtr("my-cluster-network-subnet-reserved-ip"),
			}

			clusterNetworkSubnetReservedIP, response, err := vpcService.CreateClusterNetworkSubnetReservedIP(createClusterNetworkSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
	})

	Describe(`GetClusterNetworkSubnetReservedIP - Retrieve a cluster network subnet reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetClusterNetworkSubnetReservedIP(getClusterNetworkSubnetReservedIPOptions *GetClusterNetworkSubnetReservedIPOptions)`, func() {
			getClusterNetworkSubnetReservedIPOptions := &vpcv1.GetClusterNetworkSubnetReservedIPOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				ID:                     core.StringPtr("testString"),
			}

			clusterNetworkSubnetReservedIP, response, err := vpcService.GetClusterNetworkSubnetReservedIP(getClusterNetworkSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
	})

	Describe(`UpdateClusterNetworkSubnetReservedIP - Update a cluster network subnet reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateClusterNetworkSubnetReservedIP(updateClusterNetworkSubnetReservedIPOptions *UpdateClusterNetworkSubnetReservedIPOptions)`, func() {
			clusterNetworkSubnetReservedIPPatchModel := &vpcv1.ClusterNetworkSubnetReservedIPPatch{
				AutoDelete: core.BoolPtr(false),
				Name:       core.StringPtr("my-cluster-network-subnet-reserved-ip"),
			}
			clusterNetworkSubnetReservedIPPatchModelAsPatch, asPatchErr := clusterNetworkSubnetReservedIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkSubnetReservedIPOptions := &vpcv1.UpdateClusterNetworkSubnetReservedIPOptions{
				ClusterNetworkID:                    core.StringPtr("testString"),
				ClusterNetworkSubnetID:              core.StringPtr("testString"),
				ID:                                  core.StringPtr("testString"),
				ClusterNetworkSubnetReservedIPPatch: clusterNetworkSubnetReservedIPPatchModelAsPatch,
				IfMatch:                             core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkSubnetReservedIP, response, err := vpcService.UpdateClusterNetworkSubnetReservedIP(updateClusterNetworkSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
	})

	Describe(`GetClusterNetworkSubnet - Retrieve a cluster network subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetClusterNetworkSubnet(getClusterNetworkSubnetOptions *GetClusterNetworkSubnetOptions)`, func() {
			getClusterNetworkSubnetOptions := &vpcv1.GetClusterNetworkSubnetOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				ID:               core.StringPtr("testString"),
			}

			clusterNetworkSubnet, response, err := vpcService.GetClusterNetworkSubnet(getClusterNetworkSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
	})

	Describe(`UpdateClusterNetworkSubnet - Update a cluster network subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateClusterNetworkSubnet(updateClusterNetworkSubnetOptions *UpdateClusterNetworkSubnetOptions)`, func() {
			clusterNetworkSubnetPatchModel := &vpcv1.ClusterNetworkSubnetPatch{
				Name: core.StringPtr("my-cluster-network-subnet"),
			}
			clusterNetworkSubnetPatchModelAsPatch, asPatchErr := clusterNetworkSubnetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkSubnetOptions := &vpcv1.UpdateClusterNetworkSubnetOptions{
				ClusterNetworkID:          core.StringPtr("testString"),
				ID:                        core.StringPtr("testString"),
				ClusterNetworkSubnetPatch: clusterNetworkSubnetPatchModelAsPatch,
				IfMatch:                   core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkSubnet, response, err := vpcService.UpdateClusterNetworkSubnet(updateClusterNetworkSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
	})

	Describe(`GetClusterNetwork - Retrieve a cluster network`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetClusterNetwork(getClusterNetworkOptions *GetClusterNetworkOptions)`, func() {
			getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{
				ID: core.StringPtr("testString"),
			}

			clusterNetwork, response, err := vpcService.GetClusterNetwork(getClusterNetworkOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetwork).ToNot(BeNil())
		})
	})

	Describe(`UpdateClusterNetwork - Update a cluster`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateClusterNetwork(updateClusterNetworkOptions *UpdateClusterNetworkOptions)`, func() {
			clusterNetworkPatchModel := &vpcv1.ClusterNetworkPatch{
				Name: core.StringPtr("my-cluster-network"),
			}
			clusterNetworkPatchModelAsPatch, asPatchErr := clusterNetworkPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateClusterNetworkOptions := &vpcv1.UpdateClusterNetworkOptions{
				ID:                  core.StringPtr("testString"),
				ClusterNetworkPatch: clusterNetworkPatchModelAsPatch,
				IfMatch:             core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetwork, response, err := vpcService.UpdateClusterNetwork(updateClusterNetworkOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(clusterNetwork).ToNot(BeNil())
		})
	})

	Describe(`ListDedicatedHostGroups - List dedicated host groups`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListDedicatedHostGroups(listDedicatedHostGroupsOptions *ListDedicatedHostGroupsOptions) with pagination`, func() {
			listDedicatedHostGroupsOptions := &vpcv1.ListDedicatedHostGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
				ZoneName:        core.StringPtr("us-south-1"),
			}

			listDedicatedHostGroupsOptions.Start = nil
			listDedicatedHostGroupsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.DedicatedHostGroup
			for {
				dedicatedHostGroupCollection, response, err := vpcService.ListDedicatedHostGroups(listDedicatedHostGroupsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(dedicatedHostGroupCollection).ToNot(BeNil())
				allResults = append(allResults, dedicatedHostGroupCollection.Groups...)

				nextToken, nextErr := dedicatedHostGroupCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listDedicatedHostGroupsOptions.Start = nextToken

				if listDedicatedHostGroupsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListDedicatedHostGroups(listDedicatedHostGroupsOptions *ListDedicatedHostGroupsOptions) using DedicatedHostGroupsPager`, func() {
			listDedicatedHostGroupsOptions := &vpcv1.ListDedicatedHostGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				ZoneName:        core.StringPtr("us-south-1"),
			}

			// Test GetNext().
			pager, err := vpcService.NewDedicatedHostGroupsPager(listDedicatedHostGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.DedicatedHostGroup
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewDedicatedHostGroupsPager(listDedicatedHostGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListDedicatedHostGroups() returned a total of %d item(s) using DedicatedHostGroupsPager.\n", len(allResults))
		})
	})

	Describe(`CreateDedicatedHostGroup - Create a dedicated host group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateDedicatedHostGroup(createDedicatedHostGroupOptions *CreateDedicatedHostGroupOptions)`, func() {
			dedicatedHostGroupPrototypeZoneModel := &vpcv1.DedicatedHostGroupPrototypeZoneZoneIdentityByName{
				Name: core.StringPtr("us-south-1"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createDedicatedHostGroupOptions := &vpcv1.CreateDedicatedHostGroupOptions{
				Class:         core.StringPtr("mx2"),
				Family:        core.StringPtr("balanced"),
				Zone:          dedicatedHostGroupPrototypeZoneModel,
				Name:          core.StringPtr("testString"),
				ResourceGroup: resourceGroupIdentityModel,
			}

			dedicatedHostGroup, response, err := vpcService.CreateDedicatedHostGroup(createDedicatedHostGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
	})

	Describe(`GetDedicatedHostGroup - Retrieve a dedicated host group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetDedicatedHostGroup(getDedicatedHostGroupOptions *GetDedicatedHostGroupOptions)`, func() {
			getDedicatedHostGroupOptions := &vpcv1.GetDedicatedHostGroupOptions{
				ID: core.StringPtr("testString"),
			}

			dedicatedHostGroup, response, err := vpcService.GetDedicatedHostGroup(getDedicatedHostGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
	})

	Describe(`UpdateDedicatedHostGroup - Update a dedicated host group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateDedicatedHostGroup(updateDedicatedHostGroupOptions *UpdateDedicatedHostGroupOptions)`, func() {
			dedicatedHostGroupPatchModel := &vpcv1.DedicatedHostGroupPatch{
				Name: core.StringPtr("my-dedicated-host-group-updated"),
			}
			dedicatedHostGroupPatchModelAsPatch, asPatchErr := dedicatedHostGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostGroupOptions := &vpcv1.UpdateDedicatedHostGroupOptions{
				ID:                      core.StringPtr("testString"),
				DedicatedHostGroupPatch: dedicatedHostGroupPatchModelAsPatch,
			}

			dedicatedHostGroup, response, err := vpcService.UpdateDedicatedHostGroup(updateDedicatedHostGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostGroup).ToNot(BeNil())
		})
	})

	Describe(`ListDedicatedHostProfiles - List dedicated host profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListDedicatedHostProfiles(listDedicatedHostProfilesOptions *ListDedicatedHostProfilesOptions) with pagination`, func() {
			listDedicatedHostProfilesOptions := &vpcv1.ListDedicatedHostProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listDedicatedHostProfilesOptions.Start = nil
			listDedicatedHostProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.DedicatedHostProfile
			for {
				dedicatedHostProfileCollection, response, err := vpcService.ListDedicatedHostProfiles(listDedicatedHostProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(dedicatedHostProfileCollection).ToNot(BeNil())
				allResults = append(allResults, dedicatedHostProfileCollection.Profiles...)

				nextToken, nextErr := dedicatedHostProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listDedicatedHostProfilesOptions.Start = nextToken

				if listDedicatedHostProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListDedicatedHostProfiles(listDedicatedHostProfilesOptions *ListDedicatedHostProfilesOptions) using DedicatedHostProfilesPager`, func() {
			listDedicatedHostProfilesOptions := &vpcv1.ListDedicatedHostProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewDedicatedHostProfilesPager(listDedicatedHostProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.DedicatedHostProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewDedicatedHostProfilesPager(listDedicatedHostProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListDedicatedHostProfiles() returned a total of %d item(s) using DedicatedHostProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetDedicatedHostProfile - Retrieve a dedicated host profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetDedicatedHostProfile(getDedicatedHostProfileOptions *GetDedicatedHostProfileOptions)`, func() {
			getDedicatedHostProfileOptions := &vpcv1.GetDedicatedHostProfileOptions{
				Name: core.StringPtr("mx2-host-152x1216"),
			}

			dedicatedHostProfile, response, err := vpcService.GetDedicatedHostProfile(getDedicatedHostProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostProfile).ToNot(BeNil())
		})
	})

	Describe(`ListDedicatedHosts - List dedicated hosts`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListDedicatedHosts(listDedicatedHostsOptions *ListDedicatedHostsOptions) with pagination`, func() {
			listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{
				DedicatedHostGroupID: core.StringPtr("testString"),
				Limit:                core.Int64Ptr(int64(10)),
				Name:                 core.StringPtr("my-name"),
				ResourceGroupID:      core.StringPtr("testString"),
				Start:                core.StringPtr("testString"),
				ZoneName:             core.StringPtr("us-south-1"),
			}

			listDedicatedHostsOptions.Start = nil
			listDedicatedHostsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.DedicatedHost
			for {
				dedicatedHostCollection, response, err := vpcService.ListDedicatedHosts(listDedicatedHostsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(dedicatedHostCollection).ToNot(BeNil())
				allResults = append(allResults, dedicatedHostCollection.DedicatedHosts...)

				nextToken, nextErr := dedicatedHostCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listDedicatedHostsOptions.Start = nextToken

				if listDedicatedHostsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListDedicatedHosts(listDedicatedHostsOptions *ListDedicatedHostsOptions) using DedicatedHostsPager`, func() {
			listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{
				DedicatedHostGroupID: core.StringPtr("testString"),
				Limit:                core.Int64Ptr(int64(10)),
				Name:                 core.StringPtr("my-name"),
				ResourceGroupID:      core.StringPtr("testString"),
				ZoneName:             core.StringPtr("us-south-1"),
			}

			// Test GetNext().
			pager, err := vpcService.NewDedicatedHostsPager(listDedicatedHostsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.DedicatedHost
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewDedicatedHostsPager(listDedicatedHostsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListDedicatedHosts() returned a total of %d item(s) using DedicatedHostsPager.\n", len(allResults))
		})
	})

	Describe(`CreateDedicatedHost - Create a dedicated host`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateDedicatedHost(createDedicatedHostOptions *CreateDedicatedHostOptions)`, func() {
			dedicatedHostPrototypeProfileModel := &vpcv1.DedicatedHostPrototypeProfileDedicatedHostProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/dedicated_host/profiles/mx2-host-152x1216"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			dedicatedHostByGroupGroupModel := &vpcv1.DedicatedHostByGroupGroupDedicatedHostGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host-group:0717-bcc5b834-1258-4b9c-c3b4-43bc9cf5cde0"),
			}

			dedicatedHostPrototypeModel := &vpcv1.DedicatedHostPrototypeDedicatedHostByGroup{
				InstancePlacementEnabled: core.BoolPtr(true),
				Name:                     core.StringPtr("my-host"),
				Profile:                  dedicatedHostPrototypeProfileModel,
				ResourceGroup:            resourceGroupIdentityModel,
				Group:                    dedicatedHostByGroupGroupModel,
			}

			createDedicatedHostOptions := &vpcv1.CreateDedicatedHostOptions{
				DedicatedHostPrototype: dedicatedHostPrototypeModel,
			}

			dedicatedHost, response, err := vpcService.CreateDedicatedHost(createDedicatedHostOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(dedicatedHost).ToNot(BeNil())
		})
	})

	Describe(`ListDedicatedHostDisks - List disks on a dedicated host`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListDedicatedHostDisks(listDedicatedHostDisksOptions *ListDedicatedHostDisksOptions)`, func() {
			listDedicatedHostDisksOptions := &vpcv1.ListDedicatedHostDisksOptions{
				DedicatedHostID: core.StringPtr("testString"),
			}

			dedicatedHostDiskCollection, response, err := vpcService.ListDedicatedHostDisks(listDedicatedHostDisksOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDiskCollection).ToNot(BeNil())
		})
	})

	Describe(`GetDedicatedHostDisk - Retrieve a dedicated host disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetDedicatedHostDisk(getDedicatedHostDiskOptions *GetDedicatedHostDiskOptions)`, func() {
			getDedicatedHostDiskOptions := &vpcv1.GetDedicatedHostDiskOptions{
				DedicatedHostID: core.StringPtr("testString"),
				ID:              core.StringPtr("testString"),
			}

			dedicatedHostDisk, response, err := vpcService.GetDedicatedHostDisk(getDedicatedHostDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDisk).ToNot(BeNil())
		})
	})

	Describe(`UpdateDedicatedHostDisk - Update a dedicated host disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateDedicatedHostDisk(updateDedicatedHostDiskOptions *UpdateDedicatedHostDiskOptions)`, func() {
			dedicatedHostDiskPatchModel := &vpcv1.DedicatedHostDiskPatch{
				Name: core.StringPtr("my-dedicated-host-disk-updated"),
			}
			dedicatedHostDiskPatchModelAsPatch, asPatchErr := dedicatedHostDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostDiskOptions := &vpcv1.UpdateDedicatedHostDiskOptions{
				DedicatedHostID:        core.StringPtr("testString"),
				ID:                     core.StringPtr("testString"),
				DedicatedHostDiskPatch: dedicatedHostDiskPatchModelAsPatch,
			}

			dedicatedHostDisk, response, err := vpcService.UpdateDedicatedHostDisk(updateDedicatedHostDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHostDisk).ToNot(BeNil())
		})
	})

	Describe(`GetDedicatedHost - Retrieve a dedicated host`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetDedicatedHost(getDedicatedHostOptions *GetDedicatedHostOptions)`, func() {
			getDedicatedHostOptions := &vpcv1.GetDedicatedHostOptions{
				ID: core.StringPtr("testString"),
			}

			dedicatedHost, response, err := vpcService.GetDedicatedHost(getDedicatedHostOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHost).ToNot(BeNil())
		})
	})

	Describe(`UpdateDedicatedHost - Update a dedicated host`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateDedicatedHost(updateDedicatedHostOptions *UpdateDedicatedHostOptions)`, func() {
			dedicatedHostPatchModel := &vpcv1.DedicatedHostPatch{
				InstancePlacementEnabled: core.BoolPtr(false),
				Name:                     core.StringPtr("my-host"),
			}
			dedicatedHostPatchModelAsPatch, asPatchErr := dedicatedHostPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateDedicatedHostOptions := &vpcv1.UpdateDedicatedHostOptions{
				ID:                 core.StringPtr("testString"),
				DedicatedHostPatch: dedicatedHostPatchModelAsPatch,
			}

			dedicatedHost, response, err := vpcService.UpdateDedicatedHost(updateDedicatedHostOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(dedicatedHost).ToNot(BeNil())
		})
	})

	Describe(`ListEndpointGateways - List endpoint gateways`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListEndpointGateways(listEndpointGatewaysOptions *ListEndpointGatewaysOptions) with pagination`, func() {
			listEndpointGatewaysOptions := &vpcv1.ListEndpointGatewaysOptions{
				DnsResolutionBindingMode: []string{"disabled"},
				LifecycleState:           []string{"stable"},
				Limit:                    core.Int64Ptr(int64(10)),
				Name:                     core.StringPtr("my-name"),
				ResourceGroupID:          core.StringPtr("testString"),
				Start:                    core.StringPtr("testString"),
				VPCCRN:                   core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:                    core.StringPtr("testString"),
				VPCName:                  core.StringPtr("my-vpc"),
			}

			listEndpointGatewaysOptions.Start = nil
			listEndpointGatewaysOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.EndpointGateway
			for {
				endpointGatewayCollection, response, err := vpcService.ListEndpointGateways(listEndpointGatewaysOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(endpointGatewayCollection).ToNot(BeNil())
				allResults = append(allResults, endpointGatewayCollection.EndpointGateways...)

				nextToken, nextErr := endpointGatewayCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listEndpointGatewaysOptions.Start = nextToken

				if listEndpointGatewaysOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListEndpointGateways(listEndpointGatewaysOptions *ListEndpointGatewaysOptions) using EndpointGatewaysPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewEndpointGatewaysPager(listEndpointGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.EndpointGateway
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewEndpointGatewaysPager(listEndpointGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListEndpointGateways() returned a total of %d item(s) using EndpointGatewaysPager.\n", len(allResults))
		})
	})

	Describe(`CreateEndpointGateway - Create an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateEndpointGateway(createEndpointGatewayOptions *CreateEndpointGatewayOptions)`, func() {
			endpointGatewayTargetPrototypeModel := &vpcv1.EndpointGatewayTargetPrototypeEndpointGatewayTargetResourceTypePrivatePathServiceGatewayPrototype{
				CRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::private-path-service-gateway:r006-7268d425-59b7-48fd-9735-81a7271657d5"),
				ResourceType: core.StringPtr("private_path_service_gateway"),
			}

			endpointGatewayPrototypeVPCModel := &vpcv1.EndpointGatewayPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			endpointGatewayReservedIPModel := &vpcv1.EndpointGatewayReservedIPReservedIPIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			createEndpointGatewayOptions := &vpcv1.CreateEndpointGatewayOptions{
				Target:                   endpointGatewayTargetPrototypeModel,
				VPC:                      endpointGatewayPrototypeVPCModel,
				DnsResolutionBindingMode: core.StringPtr("primary"),
				Ips:                      []vpcv1.EndpointGatewayReservedIPIntf{endpointGatewayReservedIPModel},
				Name:                     core.StringPtr("my-endpoint-gateway"),
				ResourceGroup:            resourceGroupIdentityModel,
				SecurityGroups:           []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
			}

			endpointGateway, response, err := vpcService.CreateEndpointGateway(createEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(endpointGateway).ToNot(BeNil())
		})
	})

	Describe(`ListEndpointGatewayIps - List reserved IPs bound to an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListEndpointGatewayIps(listEndpointGatewayIpsOptions *ListEndpointGatewayIpsOptions) with pagination`, func() {
			listEndpointGatewayIpsOptions := &vpcv1.ListEndpointGatewayIpsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Sort:              core.StringPtr("name"),
				Start:             core.StringPtr("testString"),
			}

			listEndpointGatewayIpsOptions.Start = nil
			listEndpointGatewayIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ReservedIP
			for {
				reservedIPCollectionEndpointGatewayContext, response, err := vpcService.ListEndpointGatewayIps(listEndpointGatewayIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(reservedIPCollectionEndpointGatewayContext).ToNot(BeNil())
				allResults = append(allResults, reservedIPCollectionEndpointGatewayContext.Ips...)

				nextToken, nextErr := reservedIPCollectionEndpointGatewayContext.GetNextStart()
				Expect(nextErr).To(BeNil())
				listEndpointGatewayIpsOptions.Start = nextToken

				if listEndpointGatewayIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListEndpointGatewayIps(listEndpointGatewayIpsOptions *ListEndpointGatewayIpsOptions) using EndpointGatewayIpsPager`, func() {
			listEndpointGatewayIpsOptions := &vpcv1.ListEndpointGatewayIpsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Sort:              core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewEndpointGatewayIpsPager(listEndpointGatewayIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewEndpointGatewayIpsPager(listEndpointGatewayIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListEndpointGatewayIps() returned a total of %d item(s) using EndpointGatewayIpsPager.\n", len(allResults))
		})
	})

	Describe(`GetEndpointGatewayIP - Retrieve a reserved IP bound to an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetEndpointGatewayIP(getEndpointGatewayIPOptions *GetEndpointGatewayIPOptions)`, func() {
			getEndpointGatewayIPOptions := &vpcv1.GetEndpointGatewayIPOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			reservedIP, response, err := vpcService.GetEndpointGatewayIP(getEndpointGatewayIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`AddEndpointGatewayIP - Bind a reserved IP to an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddEndpointGatewayIP(addEndpointGatewayIPOptions *AddEndpointGatewayIPOptions)`, func() {
			addEndpointGatewayIPOptions := &vpcv1.AddEndpointGatewayIPOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			reservedIP, response, err := vpcService.AddEndpointGatewayIP(addEndpointGatewayIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`ListEndpointGatewayResourceBindings - List resource bindings for an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListEndpointGatewayResourceBindings(listEndpointGatewayResourceBindingsOptions *ListEndpointGatewayResourceBindingsOptions) with pagination`, func() {
			listEndpointGatewayResourceBindingsOptions := &vpcv1.ListEndpointGatewayResourceBindingsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
				Start:             core.StringPtr("testString"),
			}

			listEndpointGatewayResourceBindingsOptions.Start = nil
			listEndpointGatewayResourceBindingsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.EndpointGatewayResourceBinding
			for {
				endpointGatewayResourceBindingCollection, response, err := vpcService.ListEndpointGatewayResourceBindings(listEndpointGatewayResourceBindingsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(endpointGatewayResourceBindingCollection).ToNot(BeNil())
				allResults = append(allResults, endpointGatewayResourceBindingCollection.ResourceBindings...)

				nextToken, nextErr := endpointGatewayResourceBindingCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listEndpointGatewayResourceBindingsOptions.Start = nextToken

				if listEndpointGatewayResourceBindingsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListEndpointGatewayResourceBindings(listEndpointGatewayResourceBindingsOptions *ListEndpointGatewayResourceBindingsOptions) using EndpointGatewayResourceBindingsPager`, func() {
			listEndpointGatewayResourceBindingsOptions := &vpcv1.ListEndpointGatewayResourceBindingsOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Limit:             core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewEndpointGatewayResourceBindingsPager(listEndpointGatewayResourceBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.EndpointGatewayResourceBinding
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewEndpointGatewayResourceBindingsPager(listEndpointGatewayResourceBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListEndpointGatewayResourceBindings() returned a total of %d item(s) using EndpointGatewayResourceBindingsPager.\n", len(allResults))
		})
	})

	Describe(`CreateEndpointGatewayResourceBinding - Create a resource binding for an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateEndpointGatewayResourceBinding(createEndpointGatewayResourceBindingOptions *CreateEndpointGatewayResourceBindingOptions)`, func() {
			endpointGatewayResourceBindingTargetPrototypeModel := &vpcv1.EndpointGatewayResourceBindingTargetPrototypeEndpointGatewayResourceBindingTargetByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:cloud-object-storage:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1a0ec336-f391-4091-a6fb-5e084a4c56f4:bucket:bucket-27200-lwx4cfvcue"),
			}

			createEndpointGatewayResourceBindingOptions := &vpcv1.CreateEndpointGatewayResourceBindingOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				Target:            endpointGatewayResourceBindingTargetPrototypeModel,
				Name:              core.StringPtr("my-resource-binding"),
			}

			endpointGatewayResourceBinding, response, err := vpcService.CreateEndpointGatewayResourceBinding(createEndpointGatewayResourceBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
	})

	Describe(`GetEndpointGatewayResourceBinding - Retrieve an endpoint gateway resource binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetEndpointGatewayResourceBinding(getEndpointGatewayResourceBindingOptions *GetEndpointGatewayResourceBindingOptions)`, func() {
			getEndpointGatewayResourceBindingOptions := &vpcv1.GetEndpointGatewayResourceBindingOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			endpointGatewayResourceBinding, response, err := vpcService.GetEndpointGatewayResourceBinding(getEndpointGatewayResourceBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
	})

	Describe(`UpdateEndpointGatewayResourceBinding - Update an endpoint gateway resource binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateEndpointGatewayResourceBinding(updateEndpointGatewayResourceBindingOptions *UpdateEndpointGatewayResourceBindingOptions)`, func() {
			endpointGatewayResourceBindingPatchModel := &vpcv1.EndpointGatewayResourceBindingPatch{
				Name: core.StringPtr("my-resource-binding-updated"),
			}
			endpointGatewayResourceBindingPatchModelAsPatch, asPatchErr := endpointGatewayResourceBindingPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateEndpointGatewayResourceBindingOptions := &vpcv1.UpdateEndpointGatewayResourceBindingOptions{
				EndpointGatewayID:                   core.StringPtr("testString"),
				ID:                                  core.StringPtr("testString"),
				EndpointGatewayResourceBindingPatch: endpointGatewayResourceBindingPatchModelAsPatch,
			}

			endpointGatewayResourceBinding, response, err := vpcService.UpdateEndpointGatewayResourceBinding(updateEndpointGatewayResourceBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGatewayResourceBinding).ToNot(BeNil())
		})
	})

	Describe(`GetEndpointGateway - Retrieve an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetEndpointGateway(getEndpointGatewayOptions *GetEndpointGatewayOptions)`, func() {
			getEndpointGatewayOptions := &vpcv1.GetEndpointGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			endpointGateway, response, err := vpcService.GetEndpointGateway(getEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGateway).ToNot(BeNil())
		})
	})

	Describe(`UpdateEndpointGateway - Update an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateEndpointGateway(updateEndpointGatewayOptions *UpdateEndpointGatewayOptions)`, func() {
			endpointGatewayPatchModel := &vpcv1.EndpointGatewayPatch{
				DnsResolutionBindingMode: core.StringPtr("primary"),
				Name:                     core.StringPtr("my-endpoint-gateway"),
			}
			endpointGatewayPatchModelAsPatch, asPatchErr := endpointGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateEndpointGatewayOptions := &vpcv1.UpdateEndpointGatewayOptions{
				ID:                   core.StringPtr("testString"),
				EndpointGatewayPatch: endpointGatewayPatchModelAsPatch,
			}

			endpointGateway, response, err := vpcService.UpdateEndpointGateway(updateEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(endpointGateway).ToNot(BeNil())
		})
	})

	Describe(`ListFloatingIps - List floating IPs`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListFloatingIps(listFloatingIpsOptions *ListFloatingIpsOptions) with pagination`, func() {
			listFloatingIpsOptions := &vpcv1.ListFloatingIpsOptions{
				Limit:              core.Int64Ptr(int64(10)),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
				Start:              core.StringPtr("testString"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			listFloatingIpsOptions.Start = nil
			listFloatingIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.FloatingIP
			for {
				floatingIPCollection, response, err := vpcService.ListFloatingIps(listFloatingIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(floatingIPCollection).ToNot(BeNil())
				allResults = append(allResults, floatingIPCollection.FloatingIps...)

				nextToken, nextErr := floatingIPCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listFloatingIpsOptions.Start = nextToken

				if listFloatingIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListFloatingIps(listFloatingIpsOptions *ListFloatingIpsOptions) using FloatingIpsPager`, func() {
			listFloatingIpsOptions := &vpcv1.ListFloatingIpsOptions{
				Limit:              core.Int64Ptr(int64(10)),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewFloatingIpsPager(listFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.FloatingIP
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewFloatingIpsPager(listFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListFloatingIps() returned a total of %d item(s) using FloatingIpsPager.\n", len(allResults))
		})
	})

	Describe(`CreateFloatingIP - Reserve a floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateFloatingIP(createFloatingIPOptions *CreateFloatingIPOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			floatingIPTargetPrototypeModel := &vpcv1.FloatingIPTargetPrototypeBareMetalServerNetworkInterfaceIdentityBareMetalServerNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/bare_metal_servers/0717-aad10fd7-8a02-4f3e-97f3-b18bd82cf304/network_interfaces/0717-da8c43ec-b6ca-4bd2-871e-72e288c66ee6"),
			}

			floatingIPPrototypeModel := &vpcv1.FloatingIPPrototypeFloatingIPByTarget{
				Name:          core.StringPtr("my-floating-ip"),
				ResourceGroup: resourceGroupIdentityModel,
				Target:        floatingIPTargetPrototypeModel,
			}

			createFloatingIPOptions := &vpcv1.CreateFloatingIPOptions{
				FloatingIPPrototype: floatingIPPrototypeModel,
			}

			floatingIP, response, err := vpcService.CreateFloatingIP(createFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`GetFloatingIP - Retrieve a floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetFloatingIP(getFloatingIPOptions *GetFloatingIPOptions)`, func() {
			getFloatingIPOptions := &vpcv1.GetFloatingIPOptions{
				ID: core.StringPtr("testString"),
			}

			floatingIP, response, err := vpcService.GetFloatingIP(getFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`UpdateFloatingIP - Update a floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateFloatingIP(updateFloatingIPOptions *UpdateFloatingIPOptions)`, func() {
			floatingIPTargetPatchModel := &vpcv1.FloatingIPTargetPatchBareMetalServerNetworkInterfaceIdentityBareMetalServerNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/bare_metal_servers/0717-aad10fd7-8a02-4f3e-97f3-b18bd82cf304/network_interfaces/0717-da8c43ec-b6ca-4bd2-871e-72e288c66ee6"),
			}

			floatingIPPatchModel := &vpcv1.FloatingIPPatch{
				Name:   core.StringPtr("my-floating-ip"),
				Target: floatingIPTargetPatchModel,
			}
			floatingIPPatchModelAsPatch, asPatchErr := floatingIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateFloatingIPOptions := &vpcv1.UpdateFloatingIPOptions{
				ID:              core.StringPtr("testString"),
				FloatingIPPatch: floatingIPPatchModelAsPatch,
			}

			floatingIP, response, err := vpcService.UpdateFloatingIP(updateFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`ListFlowLogCollectors - List flow log collectors`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListFlowLogCollectors(listFlowLogCollectorsOptions *ListFlowLogCollectorsOptions) with pagination`, func() {
			listFlowLogCollectorsOptions := &vpcv1.ListFlowLogCollectorsOptions{
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				ResourceGroupID:    core.StringPtr("testString"),
				Start:              core.StringPtr("testString"),
				TargetID:           core.StringPtr("testString"),
				TargetResourceType: core.StringPtr("testString"),
				VPCCRN:             core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:              core.StringPtr("testString"),
				VPCName:            core.StringPtr("my-vpc"),
			}

			listFlowLogCollectorsOptions.Start = nil
			listFlowLogCollectorsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.FlowLogCollector
			for {
				flowLogCollectorCollection, response, err := vpcService.ListFlowLogCollectors(listFlowLogCollectorsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(flowLogCollectorCollection).ToNot(BeNil())
				allResults = append(allResults, flowLogCollectorCollection.FlowLogCollectors...)

				nextToken, nextErr := flowLogCollectorCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listFlowLogCollectorsOptions.Start = nextToken

				if listFlowLogCollectorsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListFlowLogCollectors(listFlowLogCollectorsOptions *ListFlowLogCollectorsOptions) using FlowLogCollectorsPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewFlowLogCollectorsPager(listFlowLogCollectorsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.FlowLogCollector
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewFlowLogCollectorsPager(listFlowLogCollectorsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListFlowLogCollectors() returned a total of %d item(s) using FlowLogCollectorsPager.\n", len(allResults))
		})
	})

	Describe(`CreateFlowLogCollector - Create a flow log collector`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateFlowLogCollector(createFlowLogCollectorOptions *CreateFlowLogCollectorOptions)`, func() {
			flowLogCollectorPrototypeStorageBucketModel := &vpcv1.FlowLogCollectorPrototypeStorageBucketLegacyCloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName{
				Name: core.StringPtr("bucket-27200-lwx4cfvcue"),
			}

			flowLogCollectorTargetPrototypeModel := &vpcv1.FlowLogCollectorTargetPrototypeInstanceIdentityInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createFlowLogCollectorOptions := &vpcv1.CreateFlowLogCollectorOptions{
				StorageBucket: flowLogCollectorPrototypeStorageBucketModel,
				Target:        flowLogCollectorTargetPrototypeModel,
				Active:        core.BoolPtr(false),
				Name:          core.StringPtr("my-flow-log-collector"),
				ResourceGroup: resourceGroupIdentityModel,
			}

			flowLogCollector, response, err := vpcService.CreateFlowLogCollector(createFlowLogCollectorOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(flowLogCollector).ToNot(BeNil())
		})
	})

	Describe(`GetFlowLogCollector - Retrieve a flow log collector`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetFlowLogCollector(getFlowLogCollectorOptions *GetFlowLogCollectorOptions)`, func() {
			getFlowLogCollectorOptions := &vpcv1.GetFlowLogCollectorOptions{
				ID: core.StringPtr("testString"),
			}

			flowLogCollector, response, err := vpcService.GetFlowLogCollector(getFlowLogCollectorOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(flowLogCollector).ToNot(BeNil())
		})
	})

	Describe(`UpdateFlowLogCollector - Update a flow log collector`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateFlowLogCollector(updateFlowLogCollectorOptions *UpdateFlowLogCollectorOptions)`, func() {
			flowLogCollectorPatchModel := &vpcv1.FlowLogCollectorPatch{
				Active: core.BoolPtr(false),
				Name:   core.StringPtr("my-flow-log-collector"),
			}
			flowLogCollectorPatchModelAsPatch, asPatchErr := flowLogCollectorPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateFlowLogCollectorOptions := &vpcv1.UpdateFlowLogCollectorOptions{
				ID:                    core.StringPtr("testString"),
				FlowLogCollectorPatch: flowLogCollectorPatchModelAsPatch,
			}

			flowLogCollector, response, err := vpcService.UpdateFlowLogCollector(updateFlowLogCollectorOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(flowLogCollector).ToNot(BeNil())
		})
	})

	Describe(`ListRegions - List regions`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListRegions(listRegionsOptions *ListRegionsOptions)`, func() {
			listRegionsOptions := &vpcv1.ListRegionsOptions{}

			regionCollection, response, err := vpcService.ListRegions(listRegionsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(regionCollection).ToNot(BeNil())
		})
	})

	Describe(`GetRegion - Retrieve a region`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetRegion(getRegionOptions *GetRegionOptions)`, func() {
			getRegionOptions := &vpcv1.GetRegionOptions{
				Name: core.StringPtr("us-south"),
			}

			region, response, err := vpcService.GetRegion(getRegionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(region).ToNot(BeNil())
		})
	})

	Describe(`ListRegionZones - List zones in a region`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListRegionZones(listRegionZonesOptions *ListRegionZonesOptions)`, func() {
			listRegionZonesOptions := &vpcv1.ListRegionZonesOptions{
				RegionName: core.StringPtr("us-south"),
			}

			zoneCollection, response, err := vpcService.ListRegionZones(listRegionZonesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(zoneCollection).ToNot(BeNil())
		})
	})

	Describe(`GetRegionZone - Retrieve a zone`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetRegionZone(getRegionZoneOptions *GetRegionZoneOptions)`, func() {
			getRegionZoneOptions := &vpcv1.GetRegionZoneOptions{
				Name:       core.StringPtr("us-south-1"),
				RegionName: core.StringPtr("us-south"),
			}

			zone, response, err := vpcService.GetRegionZone(getRegionZoneOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(zone).ToNot(BeNil())
		})
	})

	Describe(`ListImages - List images`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListImages(listImagesOptions *ListImagesOptions) with pagination`, func() {
			listImagesOptions := &vpcv1.ListImagesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				RemoteAccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
				Status:          []string{"available"},
				UserDataFormat:  []string{"cloud_init"},
				Visibility:      core.StringPtr("private"),
			}

			listImagesOptions.Start = nil
			listImagesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Image
			for {
				imageCollection, response, err := vpcService.ListImages(listImagesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(imageCollection).ToNot(BeNil())
				allResults = append(allResults, imageCollection.Images...)

				nextToken, nextErr := imageCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listImagesOptions.Start = nextToken

				if listImagesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListImages(listImagesOptions *ListImagesOptions) using ImagesPager`, func() {
			listImagesOptions := &vpcv1.ListImagesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				RemoteAccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				ResourceGroupID: core.StringPtr("testString"),
				Status:          []string{"available"},
				UserDataFormat:  []string{"cloud_init"},
				Visibility:      core.StringPtr("private"),
			}

			// Test GetNext().
			pager, err := vpcService.NewImagesPager(listImagesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Image
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewImagesPager(listImagesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListImages() returned a total of %d item(s) using ImagesPager.\n", len(allResults))
		})
	})

	Describe(`CreateImage - Create an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateImage(createImageOptions *CreateImageOptions)`, func() {
			imageAllowedUsePrototypeModel := &vpcv1.ImageAllowedUsePrototype{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			encryptionKeyIdentityModel := &vpcv1.EncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			imageByFileFileModel := &vpcv1.ImageByFileFile{
				Href: core.StringPtr("cos://us-south/custom-image-vpc-bucket/customImage-0.vhd"),
			}

			imageByFileOperatingSystemModel := &vpcv1.ImageByFileOperatingSystemOperatingSystemIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/operating_systems/ubuntu-24-04-amd64"),
			}

			imagePrototypeModel := &vpcv1.ImagePrototypeImageByFile{
				AllowedUse:       imageAllowedUsePrototypeModel,
				DeprecationAt:    CreateMockDateTime("2026-01-02T03:04:05.006Z"),
				Name:             core.StringPtr("my-image"),
				ObsolescenceAt:   CreateMockDateTime("2026-01-02T03:04:05.006Z"),
				ResourceGroup:    resourceGroupIdentityModel,
				EncryptedDataKey: core.StringPtr("testString"),
				EncryptionKey:    encryptionKeyIdentityModel,
				File:             imageByFileFileModel,
				OperatingSystem:  imageByFileOperatingSystemModel,
			}

			createImageOptions := &vpcv1.CreateImageOptions{
				ImagePrototype: imagePrototypeModel,
			}

			image, response, err := vpcService.CreateImage(createImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(image).ToNot(BeNil())
		})
	})

	Describe(`GetImage - Retrieve an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetImage(getImageOptions *GetImageOptions)`, func() {
			getImageOptions := &vpcv1.GetImageOptions{
				ID: core.StringPtr("testString"),
			}

			image, response, err := vpcService.GetImage(getImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(image).ToNot(BeNil())
		})
	})

	Describe(`UpdateImage - Update an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateImage(updateImageOptions *UpdateImageOptions)`, func() {
			imageAllowedUsePatchModel := &vpcv1.ImageAllowedUsePatch{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			imagePatchModel := &vpcv1.ImagePatch{
				AllowedUse:     imageAllowedUsePatchModel,
				DeprecationAt:  CreateMockDateTime("2026-01-02T03:04:05.006Z"),
				Name:           core.StringPtr("my-image"),
				ObsolescenceAt: CreateMockDateTime("2026-01-02T03:04:05.006Z"),
			}
			imagePatchModelAsPatch, asPatchErr := imagePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateImageOptions := &vpcv1.UpdateImageOptions{
				ID:         core.StringPtr("testString"),
				ImagePatch: imagePatchModelAsPatch,
			}

			image, response, err := vpcService.UpdateImage(updateImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(image).ToNot(BeNil())
		})
	})

	Describe(`ListImageBareMetalServerProfiles - List bare metal server profiles compatible with an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListImageBareMetalServerProfiles(listImageBareMetalServerProfilesOptions *ListImageBareMetalServerProfilesOptions) with pagination`, func() {
			listImageBareMetalServerProfilesOptions := &vpcv1.ListImageBareMetalServerProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listImageBareMetalServerProfilesOptions.Start = nil
			listImageBareMetalServerProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.BareMetalServerProfileReference
			for {
				imageBareMetalServerProfileCollection, response, err := vpcService.ListImageBareMetalServerProfiles(listImageBareMetalServerProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(imageBareMetalServerProfileCollection).ToNot(BeNil())
				allResults = append(allResults, imageBareMetalServerProfileCollection.BareMetalServerProfiles...)

				nextToken, nextErr := imageBareMetalServerProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listImageBareMetalServerProfilesOptions.Start = nextToken

				if listImageBareMetalServerProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListImageBareMetalServerProfiles(listImageBareMetalServerProfilesOptions *ListImageBareMetalServerProfilesOptions) using ImageBareMetalServerProfilesPager`, func() {
			listImageBareMetalServerProfilesOptions := &vpcv1.ListImageBareMetalServerProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewImageBareMetalServerProfilesPager(listImageBareMetalServerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.BareMetalServerProfileReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewImageBareMetalServerProfilesPager(listImageBareMetalServerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListImageBareMetalServerProfiles() returned a total of %d item(s) using ImageBareMetalServerProfilesPager.\n", len(allResults))
		})
	})

	Describe(`DeprecateImage - Deprecate an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeprecateImage(deprecateImageOptions *DeprecateImageOptions)`, func() {
			deprecateImageOptions := &vpcv1.DeprecateImageOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeprecateImage(deprecateImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`ListImageInstanceProfiles - List instance profiles compatible with an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListImageInstanceProfiles(listImageInstanceProfilesOptions *ListImageInstanceProfilesOptions) with pagination`, func() {
			listImageInstanceProfilesOptions := &vpcv1.ListImageInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listImageInstanceProfilesOptions.Start = nil
			listImageInstanceProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceProfileReference
			for {
				imageInstanceProfileCollection, response, err := vpcService.ListImageInstanceProfiles(listImageInstanceProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(imageInstanceProfileCollection).ToNot(BeNil())
				allResults = append(allResults, imageInstanceProfileCollection.InstanceProfiles...)

				nextToken, nextErr := imageInstanceProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listImageInstanceProfilesOptions.Start = nextToken

				if listImageInstanceProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListImageInstanceProfiles(listImageInstanceProfilesOptions *ListImageInstanceProfilesOptions) using ImageInstanceProfilesPager`, func() {
			listImageInstanceProfilesOptions := &vpcv1.ListImageInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewImageInstanceProfilesPager(listImageInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewImageInstanceProfilesPager(listImageInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListImageInstanceProfiles() returned a total of %d item(s) using ImageInstanceProfilesPager.\n", len(allResults))
		})
	})

	Describe(`ObsoleteImage - Obsolete an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ObsoleteImage(obsoleteImageOptions *ObsoleteImageOptions)`, func() {
			obsoleteImageOptions := &vpcv1.ObsoleteImageOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.ObsoleteImage(obsoleteImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`ListImageExportJobs - List export jobs for an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListImageExportJobs(listImageExportJobsOptions *ListImageExportJobsOptions)`, func() {
			listImageExportJobsOptions := &vpcv1.ListImageExportJobsOptions{
				ImageID: core.StringPtr("testString"),
				Name:    core.StringPtr("my-name"),
			}

			imageExportJobUnpaginatedCollection, response, err := vpcService.ListImageExportJobs(listImageExportJobsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJobUnpaginatedCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateImageExportJob - Create an export job for an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateImageExportJob(createImageExportJobOptions *CreateImageExportJobOptions)`, func() {
			cloudObjectStorageBucketIdentityModel := &vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName{
				Name: core.StringPtr("bucket-27200-lwx4cfvcue"),
			}

			createImageExportJobOptions := &vpcv1.CreateImageExportJobOptions{
				ImageID:       core.StringPtr("testString"),
				StorageBucket: cloudObjectStorageBucketIdentityModel,
				Format:        core.StringPtr("qcow2"),
				Name:          core.StringPtr("my-image-export"),
			}

			imageExportJob, response, err := vpcService.CreateImageExportJob(createImageExportJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(imageExportJob).ToNot(BeNil())
		})
	})

	Describe(`GetImageExportJob - Retrieve an image export job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetImageExportJob(getImageExportJobOptions *GetImageExportJobOptions)`, func() {
			getImageExportJobOptions := &vpcv1.GetImageExportJobOptions{
				ID:      core.StringPtr("testString"),
				ImageID: core.StringPtr("testString"),
			}

			imageExportJob, response, err := vpcService.GetImageExportJob(getImageExportJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJob).ToNot(BeNil())
		})
	})

	Describe(`UpdateImageExportJob - Update an image export job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateImageExportJob(updateImageExportJobOptions *UpdateImageExportJobOptions)`, func() {
			imageExportJobPatchModel := &vpcv1.ImageExportJobPatch{
				Name: core.StringPtr("my-image-export-job-updated"),
			}
			imageExportJobPatchModelAsPatch, asPatchErr := imageExportJobPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateImageExportJobOptions := &vpcv1.UpdateImageExportJobOptions{
				ID:                  core.StringPtr("testString"),
				ImageID:             core.StringPtr("testString"),
				ImageExportJobPatch: imageExportJobPatchModelAsPatch,
			}

			imageExportJob, response, err := vpcService.UpdateImageExportJob(updateImageExportJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(imageExportJob).ToNot(BeNil())
		})
	})

	Describe(`ListOperatingSystems - List operating systems`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListOperatingSystems(listOperatingSystemsOptions *ListOperatingSystemsOptions) with pagination`, func() {
			listOperatingSystemsOptions := &vpcv1.ListOperatingSystemsOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listOperatingSystemsOptions.Start = nil
			listOperatingSystemsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.OperatingSystem
			for {
				operatingSystemCollection, response, err := vpcService.ListOperatingSystems(listOperatingSystemsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(operatingSystemCollection).ToNot(BeNil())
				allResults = append(allResults, operatingSystemCollection.OperatingSystems...)

				nextToken, nextErr := operatingSystemCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listOperatingSystemsOptions.Start = nextToken

				if listOperatingSystemsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListOperatingSystems(listOperatingSystemsOptions *ListOperatingSystemsOptions) using OperatingSystemsPager`, func() {
			listOperatingSystemsOptions := &vpcv1.ListOperatingSystemsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewOperatingSystemsPager(listOperatingSystemsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.OperatingSystem
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewOperatingSystemsPager(listOperatingSystemsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListOperatingSystems() returned a total of %d item(s) using OperatingSystemsPager.\n", len(allResults))
		})
	})

	Describe(`GetOperatingSystem - Retrieve an operating system`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetOperatingSystem(getOperatingSystemOptions *GetOperatingSystemOptions)`, func() {
			getOperatingSystemOptions := &vpcv1.GetOperatingSystemOptions{
				Name: core.StringPtr("red-7-amd64"),
			}

			operatingSystem, response, err := vpcService.GetOperatingSystem(getOperatingSystemOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(operatingSystem).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceGroups - List instance groups`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceGroups(listInstanceGroupsOptions *ListInstanceGroupsOptions) with pagination`, func() {
			listInstanceGroupsOptions := &vpcv1.ListInstanceGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listInstanceGroupsOptions.Start = nil
			listInstanceGroupsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceGroup
			for {
				instanceGroupCollection, response, err := vpcService.ListInstanceGroups(listInstanceGroupsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceGroupCollection).ToNot(BeNil())
				allResults = append(allResults, instanceGroupCollection.InstanceGroups...)

				nextToken, nextErr := instanceGroupCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceGroupsOptions.Start = nextToken

				if listInstanceGroupsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceGroups(listInstanceGroupsOptions *ListInstanceGroupsOptions) using InstanceGroupsPager`, func() {
			listInstanceGroupsOptions := &vpcv1.ListInstanceGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceGroupsPager(listInstanceGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceGroup
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceGroupsPager(listInstanceGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceGroups() returned a total of %d item(s) using InstanceGroupsPager.\n", len(allResults))
		})
	})

	Describe(`CreateInstanceGroup - Create an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceGroup(createInstanceGroupOptions *CreateInstanceGroupOptions)`, func() {
			instanceTemplateIdentityModel := &vpcv1.InstanceTemplateIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance-template:0717-e6c2c7d8-ad57-4f38-a21c-a86265b6aeb2"),
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			loadBalancerIdentityModel := &vpcv1.LoadBalancerIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
			}

			loadBalancerPoolIdentityModel := &vpcv1.LoadBalancerPoolIdentityLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createInstanceGroupOptions := &vpcv1.CreateInstanceGroupOptions{
				InstanceTemplate: instanceTemplateIdentityModel,
				Subnets:          []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
				ApplicationPort:  core.Int64Ptr(int64(22)),
				LoadBalancer:     loadBalancerIdentityModel,
				LoadBalancerPool: loadBalancerPoolIdentityModel,
				MembershipCount:  core.Int64Ptr(int64(10)),
				Name:             core.StringPtr("my-instance-group"),
				ResourceGroup:    resourceGroupIdentityModel,
			}

			instanceGroup, response, err := vpcService.CreateInstanceGroup(createInstanceGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroup).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceGroup - Retrieve an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceGroup(getInstanceGroupOptions *GetInstanceGroupOptions)`, func() {
			getInstanceGroupOptions := &vpcv1.GetInstanceGroupOptions{
				ID: core.StringPtr("testString"),
			}

			instanceGroup, response, err := vpcService.GetInstanceGroup(getInstanceGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroup).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceGroup - Update an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceGroup(updateInstanceGroupOptions *UpdateInstanceGroupOptions)`, func() {
			instanceTemplateIdentityModel := &vpcv1.InstanceTemplateIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance-template:0717-e6c2c7d8-ad57-4f38-a21c-a86265b6aeb2"),
			}

			loadBalancerIdentityModel := &vpcv1.LoadBalancerIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
			}

			loadBalancerPoolIdentityModel := &vpcv1.LoadBalancerPoolIdentityLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			instanceGroupPatchModel := &vpcv1.InstanceGroupPatch{
				ApplicationPort:  core.Int64Ptr(int64(22)),
				InstanceTemplate: instanceTemplateIdentityModel,
				LoadBalancer:     loadBalancerIdentityModel,
				LoadBalancerPool: loadBalancerPoolIdentityModel,
				MembershipCount:  core.Int64Ptr(int64(20)),
				Name:             core.StringPtr("my-instance-group"),
				Subnets:          []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			}
			instanceGroupPatchModelAsPatch, asPatchErr := instanceGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupOptions := &vpcv1.UpdateInstanceGroupOptions{
				ID:                 core.StringPtr("testString"),
				InstanceGroupPatch: instanceGroupPatchModelAsPatch,
			}

			instanceGroup, response, err := vpcService.UpdateInstanceGroup(updateInstanceGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroup).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceGroupManagers - List managers for an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceGroupManagers(listInstanceGroupManagersOptions *ListInstanceGroupManagersOptions) with pagination`, func() {
			listInstanceGroupManagersOptions := &vpcv1.ListInstanceGroupManagersOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
				Start:           core.StringPtr("testString"),
			}

			listInstanceGroupManagersOptions.Start = nil
			listInstanceGroupManagersOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceGroupManagerIntf
			for {
				instanceGroupManagerCollection, response, err := vpcService.ListInstanceGroupManagers(listInstanceGroupManagersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceGroupManagerCollection).ToNot(BeNil())
				allResults = append(allResults, instanceGroupManagerCollection.Managers...)

				nextToken, nextErr := instanceGroupManagerCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceGroupManagersOptions.Start = nextToken

				if listInstanceGroupManagersOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceGroupManagers(listInstanceGroupManagersOptions *ListInstanceGroupManagersOptions) using InstanceGroupManagersPager`, func() {
			listInstanceGroupManagersOptions := &vpcv1.ListInstanceGroupManagersOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceGroupManagersPager(listInstanceGroupManagersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceGroupManagerIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceGroupManagersPager(listInstanceGroupManagersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceGroupManagers() returned a total of %d item(s) using InstanceGroupManagersPager.\n", len(allResults))
		})
	})

	Describe(`CreateInstanceGroupManager - Create a manager for an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceGroupManager(createInstanceGroupManagerOptions *CreateInstanceGroupManagerOptions)`, func() {
			instanceGroupManagerPrototypeModel := &vpcv1.InstanceGroupManagerPrototypeInstanceGroupManagerAutoScalePrototype{
				ManagementEnabled:  core.BoolPtr(true),
				Name:               core.StringPtr("my-instance-group-manager"),
				AggregationWindow:  core.Int64Ptr(int64(120)),
				Cooldown:           core.Int64Ptr(int64(210)),
				ManagerType:        core.StringPtr("autoscale"),
				MaxMembershipCount: core.Int64Ptr(int64(10)),
				MinMembershipCount: core.Int64Ptr(int64(10)),
			}

			createInstanceGroupManagerOptions := &vpcv1.CreateInstanceGroupManagerOptions{
				InstanceGroupID:               core.StringPtr("testString"),
				InstanceGroupManagerPrototype: instanceGroupManagerPrototypeModel,
			}

			instanceGroupManager, response, err := vpcService.CreateInstanceGroupManager(createInstanceGroupManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceGroupManager - Retrieve an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceGroupManager(getInstanceGroupManagerOptions *GetInstanceGroupManagerOptions)`, func() {
			getInstanceGroupManagerOptions := &vpcv1.GetInstanceGroupManagerOptions{
				ID:              core.StringPtr("testString"),
				InstanceGroupID: core.StringPtr("testString"),
			}

			instanceGroupManager, response, err := vpcService.GetInstanceGroupManager(getInstanceGroupManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceGroupManager - Update an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceGroupManager(updateInstanceGroupManagerOptions *UpdateInstanceGroupManagerOptions)`, func() {
			instanceGroupManagerPatchModel := &vpcv1.InstanceGroupManagerPatch{
				AggregationWindow:  core.Int64Ptr(int64(120)),
				Cooldown:           core.Int64Ptr(int64(210)),
				ManagementEnabled:  core.BoolPtr(false),
				MaxMembershipCount: core.Int64Ptr(int64(10)),
				MinMembershipCount: core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-instance-group-manager"),
			}
			instanceGroupManagerPatchModelAsPatch, asPatchErr := instanceGroupManagerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerOptions := &vpcv1.UpdateInstanceGroupManagerOptions{
				ID:                        core.StringPtr("testString"),
				InstanceGroupID:           core.StringPtr("testString"),
				InstanceGroupManagerPatch: instanceGroupManagerPatchModelAsPatch,
			}

			instanceGroupManager, response, err := vpcService.UpdateInstanceGroupManager(updateInstanceGroupManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManager).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceGroupManagerActions - List actions for an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceGroupManagerActions(listInstanceGroupManagerActionsOptions *ListInstanceGroupManagerActionsOptions) with pagination`, func() {
			listInstanceGroupManagerActionsOptions := &vpcv1.ListInstanceGroupManagerActionsOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
				Start:                  core.StringPtr("testString"),
			}

			listInstanceGroupManagerActionsOptions.Start = nil
			listInstanceGroupManagerActionsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceGroupManagerActionIntf
			for {
				instanceGroupManagerActionsCollection, response, err := vpcService.ListInstanceGroupManagerActions(listInstanceGroupManagerActionsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceGroupManagerActionsCollection).ToNot(BeNil())
				allResults = append(allResults, instanceGroupManagerActionsCollection.Actions...)

				nextToken, nextErr := instanceGroupManagerActionsCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceGroupManagerActionsOptions.Start = nextToken

				if listInstanceGroupManagerActionsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceGroupManagerActions(listInstanceGroupManagerActionsOptions *ListInstanceGroupManagerActionsOptions) using InstanceGroupManagerActionsPager`, func() {
			listInstanceGroupManagerActionsOptions := &vpcv1.ListInstanceGroupManagerActionsOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceGroupManagerActionsPager(listInstanceGroupManagerActionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceGroupManagerActionIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceGroupManagerActionsPager(listInstanceGroupManagerActionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceGroupManagerActions() returned a total of %d item(s) using InstanceGroupManagerActionsPager.\n", len(allResults))
		})
	})

	Describe(`CreateInstanceGroupManagerAction - Create an instance group manager action`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceGroupManagerAction(createInstanceGroupManagerActionOptions *CreateInstanceGroupManagerActionOptions)`, func() {
			instanceGroupManagerScheduledActionGroupPrototypeModel := &vpcv1.InstanceGroupManagerScheduledActionGroupPrototype{
				MembershipCount: core.Int64Ptr(int64(10)),
			}

			instanceGroupManagerActionPrototypeModel := &vpcv1.InstanceGroupManagerActionPrototypeScheduledActionPrototypeByCronSpecWithGroup{
				Name:     core.StringPtr("my-instance-group-manager-action"),
				CronSpec: core.StringPtr("30 */2 * * 1-5"),
				Group:    instanceGroupManagerScheduledActionGroupPrototypeModel,
			}

			createInstanceGroupManagerActionOptions := &vpcv1.CreateInstanceGroupManagerActionOptions{
				InstanceGroupID:                     core.StringPtr("testString"),
				InstanceGroupManagerID:              core.StringPtr("testString"),
				InstanceGroupManagerActionPrototype: instanceGroupManagerActionPrototypeModel,
			}

			instanceGroupManagerAction, response, err := vpcService.CreateInstanceGroupManagerAction(createInstanceGroupManagerActionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceGroupManagerAction - Retrieve specified instance group manager action`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceGroupManagerAction(getInstanceGroupManagerActionOptions *GetInstanceGroupManagerActionOptions)`, func() {
			getInstanceGroupManagerActionOptions := &vpcv1.GetInstanceGroupManagerActionOptions{
				ID:                     core.StringPtr("testString"),
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
			}

			instanceGroupManagerAction, response, err := vpcService.GetInstanceGroupManagerAction(getInstanceGroupManagerActionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceGroupManagerAction - Update specified instance group manager action`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceGroupManagerAction(updateInstanceGroupManagerActionOptions *UpdateInstanceGroupManagerActionOptions)`, func() {
			instanceGroupManagerActionGroupPatchModel := &vpcv1.InstanceGroupManagerActionGroupPatch{
				MembershipCount: core.Int64Ptr(int64(10)),
			}

			instanceGroupManagerActionManagerPatchModel := &vpcv1.InstanceGroupManagerActionManagerPatch{
				MaxMembershipCount: core.Int64Ptr(int64(10)),
				MinMembershipCount: core.Int64Ptr(int64(10)),
			}

			instanceGroupManagerActionPatchModel := &vpcv1.InstanceGroupManagerActionPatch{
				AutoDelete:        core.BoolPtr(false),
				AutoDeleteTimeout: core.Int64Ptr(int64(24)),
				CronSpec:          core.StringPtr("30 */2 * * 1-5"),
				Group:             instanceGroupManagerActionGroupPatchModel,
				Manager:           instanceGroupManagerActionManagerPatchModel,
				Name:              core.StringPtr("my-instance-group-manager-action"),
				RunAt:             CreateMockDateTime("2026-01-02T03:04:05.006Z"),
			}
			instanceGroupManagerActionPatchModelAsPatch, asPatchErr := instanceGroupManagerActionPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerActionOptions := &vpcv1.UpdateInstanceGroupManagerActionOptions{
				ID:                              core.StringPtr("testString"),
				InstanceGroupID:                 core.StringPtr("testString"),
				InstanceGroupManagerID:          core.StringPtr("testString"),
				InstanceGroupManagerActionPatch: instanceGroupManagerActionPatchModelAsPatch,
			}

			instanceGroupManagerAction, response, err := vpcService.UpdateInstanceGroupManagerAction(updateInstanceGroupManagerActionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerAction).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceGroupManagerPolicies - List policies for an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceGroupManagerPolicies(listInstanceGroupManagerPoliciesOptions *ListInstanceGroupManagerPoliciesOptions) with pagination`, func() {
			listInstanceGroupManagerPoliciesOptions := &vpcv1.ListInstanceGroupManagerPoliciesOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
				Start:                  core.StringPtr("testString"),
			}

			listInstanceGroupManagerPoliciesOptions.Start = nil
			listInstanceGroupManagerPoliciesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceGroupManagerPolicyIntf
			for {
				instanceGroupManagerPolicyCollection, response, err := vpcService.ListInstanceGroupManagerPolicies(listInstanceGroupManagerPoliciesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceGroupManagerPolicyCollection).ToNot(BeNil())
				allResults = append(allResults, instanceGroupManagerPolicyCollection.Policies...)

				nextToken, nextErr := instanceGroupManagerPolicyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceGroupManagerPoliciesOptions.Start = nextToken

				if listInstanceGroupManagerPoliciesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceGroupManagerPolicies(listInstanceGroupManagerPoliciesOptions *ListInstanceGroupManagerPoliciesOptions) using InstanceGroupManagerPoliciesPager`, func() {
			listInstanceGroupManagerPoliciesOptions := &vpcv1.ListInstanceGroupManagerPoliciesOptions{
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
				Limit:                  core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceGroupManagerPoliciesPager(listInstanceGroupManagerPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceGroupManagerPolicyIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceGroupManagerPoliciesPager(listInstanceGroupManagerPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceGroupManagerPolicies() returned a total of %d item(s) using InstanceGroupManagerPoliciesPager.\n", len(allResults))
		})
	})

	Describe(`CreateInstanceGroupManagerPolicy - Create a policy for an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceGroupManagerPolicy(createInstanceGroupManagerPolicyOptions *CreateInstanceGroupManagerPolicyOptions)`, func() {
			instanceGroupManagerPolicyPrototypeModel := &vpcv1.InstanceGroupManagerPolicyPrototypeInstanceGroupManagerTargetPolicyPrototype{
				Name:        core.StringPtr("my-instance-group-manager-policy"),
				MetricType:  core.StringPtr("cpu"),
				MetricValue: core.Int64Ptr(int64(0)),
				PolicyType:  core.StringPtr("target"),
			}

			createInstanceGroupManagerPolicyOptions := &vpcv1.CreateInstanceGroupManagerPolicyOptions{
				InstanceGroupID:                     core.StringPtr("testString"),
				InstanceGroupManagerID:              core.StringPtr("testString"),
				InstanceGroupManagerPolicyPrototype: instanceGroupManagerPolicyPrototypeModel,
			}

			instanceGroupManagerPolicy, response, err := vpcService.CreateInstanceGroupManagerPolicy(createInstanceGroupManagerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceGroupManagerPolicy - Retrieve an instance group manager policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceGroupManagerPolicy(getInstanceGroupManagerPolicyOptions *GetInstanceGroupManagerPolicyOptions)`, func() {
			getInstanceGroupManagerPolicyOptions := &vpcv1.GetInstanceGroupManagerPolicyOptions{
				ID:                     core.StringPtr("testString"),
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
			}

			instanceGroupManagerPolicy, response, err := vpcService.GetInstanceGroupManagerPolicy(getInstanceGroupManagerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceGroupManagerPolicy - Update an instance group manager policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceGroupManagerPolicy(updateInstanceGroupManagerPolicyOptions *UpdateInstanceGroupManagerPolicyOptions)`, func() {
			instanceGroupManagerPolicyPatchModel := &vpcv1.InstanceGroupManagerPolicyPatch{
				MetricType:  core.StringPtr("cpu"),
				MetricValue: core.Int64Ptr(int64(75)),
				Name:        core.StringPtr("my-instance-group-manager-policy"),
			}
			instanceGroupManagerPolicyPatchModelAsPatch, asPatchErr := instanceGroupManagerPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupManagerPolicyOptions := &vpcv1.UpdateInstanceGroupManagerPolicyOptions{
				ID:                              core.StringPtr("testString"),
				InstanceGroupID:                 core.StringPtr("testString"),
				InstanceGroupManagerID:          core.StringPtr("testString"),
				InstanceGroupManagerPolicyPatch: instanceGroupManagerPolicyPatchModelAsPatch,
			}

			instanceGroupManagerPolicy, response, err := vpcService.UpdateInstanceGroupManagerPolicy(updateInstanceGroupManagerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupManagerPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceGroupMemberships - List memberships for an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceGroupMemberships(listInstanceGroupMembershipsOptions *ListInstanceGroupMembershipsOptions) with pagination`, func() {
			listInstanceGroupMembershipsOptions := &vpcv1.ListInstanceGroupMembershipsOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
				Start:           core.StringPtr("testString"),
			}

			listInstanceGroupMembershipsOptions.Start = nil
			listInstanceGroupMembershipsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceGroupMembership
			for {
				instanceGroupMembershipCollection, response, err := vpcService.ListInstanceGroupMemberships(listInstanceGroupMembershipsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceGroupMembershipCollection).ToNot(BeNil())
				allResults = append(allResults, instanceGroupMembershipCollection.Memberships...)

				nextToken, nextErr := instanceGroupMembershipCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceGroupMembershipsOptions.Start = nextToken

				if listInstanceGroupMembershipsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceGroupMemberships(listInstanceGroupMembershipsOptions *ListInstanceGroupMembershipsOptions) using InstanceGroupMembershipsPager`, func() {
			listInstanceGroupMembershipsOptions := &vpcv1.ListInstanceGroupMembershipsOptions{
				InstanceGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceGroupMembershipsPager(listInstanceGroupMembershipsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceGroupMembership
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceGroupMembershipsPager(listInstanceGroupMembershipsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceGroupMemberships() returned a total of %d item(s) using InstanceGroupMembershipsPager.\n", len(allResults))
		})
	})

	Describe(`GetInstanceGroupMembership - Retrieve an instance group membership`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceGroupMembership(getInstanceGroupMembershipOptions *GetInstanceGroupMembershipOptions)`, func() {
			getInstanceGroupMembershipOptions := &vpcv1.GetInstanceGroupMembershipOptions{
				ID:              core.StringPtr("testString"),
				InstanceGroupID: core.StringPtr("testString"),
			}

			instanceGroupMembership, response, err := vpcService.GetInstanceGroupMembership(getInstanceGroupMembershipOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupMembership).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceGroupMembership - Update an instance group membership`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceGroupMembership(updateInstanceGroupMembershipOptions *UpdateInstanceGroupMembershipOptions)`, func() {
			instanceGroupMembershipPatchModel := &vpcv1.InstanceGroupMembershipPatch{
				DeleteInstanceOnMembershipDelete: core.BoolPtr(false),
				Name:                             core.StringPtr("my-instance-group-membership"),
			}
			instanceGroupMembershipPatchModelAsPatch, asPatchErr := instanceGroupMembershipPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceGroupMembershipOptions := &vpcv1.UpdateInstanceGroupMembershipOptions{
				ID:                           core.StringPtr("testString"),
				InstanceGroupID:              core.StringPtr("testString"),
				InstanceGroupMembershipPatch: instanceGroupMembershipPatchModelAsPatch,
			}

			instanceGroupMembership, response, err := vpcService.UpdateInstanceGroupMembership(updateInstanceGroupMembershipOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceGroupMembership).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceTemplates - List instance templates`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceTemplates(listInstanceTemplatesOptions *ListInstanceTemplatesOptions)`, func() {
			listInstanceTemplatesOptions := &vpcv1.ListInstanceTemplatesOptions{}

			instanceTemplateCollection, response, err := vpcService.ListInstanceTemplates(listInstanceTemplatesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplateCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceTemplate - Create an instance template`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceTemplate(createInstanceTemplateOptions *CreateInstanceTemplateOptions)`, func() {
			instanceAvailabilityPrototypeModel := &vpcv1.InstanceAvailabilityPrototype{
				Class: core.StringPtr("standard"),
			}

			instanceAvailabilityPolicyPrototypeModel := &vpcv1.InstanceAvailabilityPolicyPrototype{
				HostFailure: core.StringPtr("restart"),
				Preemption:  core.StringPtr("stop"),
			}

			instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"),
			}

			instanceClusterNetworkAttachmentPrototypeInstanceContextModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{
				ClusterNetworkInterface: instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel,
				Name:                    core.StringPtr("my-instance-network-attachment"),
			}

			instanceDefaultTrustedProfilePrototypeTargetModel := &vpcv1.InstanceDefaultTrustedProfilePrototypeTargetTrustedProfileIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:iam-identity::a/aa2432b1fa4d4ace891e9b80fc104e34::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"),
			}

			instanceDefaultTrustedProfilePrototypeModel := &vpcv1.InstanceDefaultTrustedProfilePrototype{
				AutoLink: core.BoolPtr(false),
				Target:   instanceDefaultTrustedProfilePrototypeTargetModel,
			}

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			instanceMetadataServicePrototypeModel := &vpcv1.InstanceMetadataServicePrototype{
				Enabled:          core.BoolPtr(false),
				Protocol:         core.StringPtr("https"),
				ResponseHopLimit: core.Int64Ptr(int64(2)),
			}

			instancePlacementTargetPrototypeModel := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentityDedicatedHostGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host-group:0717-bcc5b834-1258-4b9c-c3b4-43bc9cf5cde0"),
			}

			instanceTemplatePrototypeProfileModel := &vpcv1.InstanceTemplatePrototypeProfileInstanceProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"),
			}

			reservationIdentityModel := &vpcv1.ReservationIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
			}

			instanceReservationAffinityPrototypeModel := &vpcv1.InstanceReservationAffinityPrototype{
				Policy: core.StringPtr("automatic"),
				Pool:   []vpcv1.ReservationIdentityIntf{reservationIdentityModel},
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			instanceVcpuPrototypeModel := &vpcv1.InstanceVcpuPrototype{
				Percentage: core.Int64Ptr(int64(100)),
			}

			volumeAttachmentPrototypeVolumeModel := &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentityVolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			volumeAttachmentPrototypeModel := &vpcv1.VolumeAttachmentPrototype{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(false),
				Name:                         core.StringPtr("my-volume-attachment"),
				Volume:                       volumeAttachmentPrototypeVolumeModel,
			}

			instanceTemplatePrototypeVPCModel := &vpcv1.InstanceTemplatePrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			volumeAllowedUsePrototypeModel := &vpcv1.VolumeAllowedUsePrototype{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			volumePrototypeInstanceByImageContextEncryptionKeyModel := &vpcv1.VolumePrototypeInstanceByImageContextEncryptionKeyEncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			volumePrototypeInstanceByImageContextProfileModel := &vpcv1.VolumePrototypeInstanceByImageContextProfileVolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumePrototypeInstanceByImageContextResourceGroupModel := &vpcv1.VolumePrototypeInstanceByImageContextResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			volumeAttachmentPrototypeInstanceByImageContextVolumeModel := &vpcv1.VolumeAttachmentPrototypeInstanceByImageContextVolume{
				AllowedUse:    volumeAllowedUsePrototypeModel,
				Bandwidth:     core.Int64Ptr(int64(1000)),
				Capacity:      core.Int64Ptr(int64(100)),
				EncryptionKey: volumePrototypeInstanceByImageContextEncryptionKeyModel,
				Iops:          core.Int64Ptr(int64(10000)),
				Name:          core.StringPtr("my-volume"),
				Profile:       volumePrototypeInstanceByImageContextProfileModel,
				ResourceGroup: volumePrototypeInstanceByImageContextResourceGroupModel,
				UserTags:      []string{"testString"},
			}

			instanceTemplateByCatalogOfferingBootVolumeAttachmentModel := &vpcv1.InstanceTemplateByCatalogOfferingBootVolumeAttachment{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(true),
				Name:                         core.StringPtr("my-volume-attachment"),
				Volume:                       volumeAttachmentPrototypeInstanceByImageContextVolumeModel,
			}

			catalogOfferingVersionPlanIdentityModel := &vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:51c9e0db-2911-45a6-adb0-ac5332d27cf2:plan:sw.51c9e0db-2911-45a6-adb0-ac5332d27cf2.772c0dbe-aa62-482e-adbe-a3fc20101e0e"),
			}

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			instanceCatalogOfferingPrototypeModel := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
				Plan:     catalogOfferingVersionPlanIdentityModel,
				Offering: catalogOfferingIdentityModel,
			}

			instanceTemplateByCatalogOfferingZoneModel := &vpcv1.InstanceTemplateByCatalogOfferingZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			instanceNetworkAttachmentPrototypeModel := &vpcv1.InstanceNetworkAttachmentPrototype{
				Name:                    core.StringPtr("my-instance-network-attachment"),
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instanceByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.InstanceByNetworkAttachmentPrimaryNetworkAttachment{
				Name:                    core.StringPtr("my-instance-network-attachment"),
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instanceTemplatePrototypeModel := &vpcv1.InstanceTemplatePrototypeInstanceTemplateByCatalogOfferingInstanceTemplateByCatalogOfferingInstanceByNetworkAttachment{
				Availability:              instanceAvailabilityPrototypeModel,
				AvailabilityPolicy:        instanceAvailabilityPolicyPrototypeModel,
				ClusterNetworkAttachments: []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{*instanceClusterNetworkAttachmentPrototypeInstanceContextModel},
				ConfidentialComputeMode:   core.StringPtr("disabled"),
				DefaultTrustedProfile:     instanceDefaultTrustedProfilePrototypeModel,
				EnableSecureBoot:          core.BoolPtr(true),
				Keys:                      []vpcv1.KeyIdentityIntf{keyIdentityModel},
				MetadataService:           instanceMetadataServicePrototypeModel,
				Name:                      core.StringPtr("my-instance"),
				PlacementTarget:           instancePlacementTargetPrototypeModel,
				Profile:                   instanceTemplatePrototypeProfileModel,
				ReservationAffinity:       instanceReservationAffinityPrototypeModel,
				ResourceGroup:             resourceGroupIdentityModel,
				ThreadsPerCore:            core.Int64Ptr(int64(1)),
				TotalVolumeBandwidth:      core.Int64Ptr(int64(500)),
				UserData:                  core.StringPtr("[...]"),
				Vcpu:                      instanceVcpuPrototypeModel,
				VolumeAttachments:         []vpcv1.VolumeAttachmentPrototype{*volumeAttachmentPrototypeModel},
				VolumeBandwidthQosMode:    core.StringPtr("pooled"),
				VPC:                       instanceTemplatePrototypeVPCModel,
				BootVolumeAttachment:      instanceTemplateByCatalogOfferingBootVolumeAttachmentModel,
				CatalogOffering:           instanceCatalogOfferingPrototypeModel,
				Zone:                      instanceTemplateByCatalogOfferingZoneModel,
				NetworkAttachments:        []vpcv1.InstanceNetworkAttachmentPrototype{*instanceNetworkAttachmentPrototypeModel},
				PrimaryNetworkAttachment:  instanceByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createInstanceTemplateOptions := &vpcv1.CreateInstanceTemplateOptions{
				InstanceTemplatePrototype: instanceTemplatePrototypeModel,
			}

			instanceTemplate, response, err := vpcService.CreateInstanceTemplate(createInstanceTemplateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceTemplate).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceTemplate - Retrieve an instance template`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceTemplate(getInstanceTemplateOptions *GetInstanceTemplateOptions)`, func() {
			getInstanceTemplateOptions := &vpcv1.GetInstanceTemplateOptions{
				ID: core.StringPtr("testString"),
			}

			instanceTemplate, response, err := vpcService.GetInstanceTemplate(getInstanceTemplateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplate).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceTemplate - Update an instance template`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceTemplate(updateInstanceTemplateOptions *UpdateInstanceTemplateOptions)`, func() {
			instanceTemplatePatchModel := &vpcv1.InstanceTemplatePatch{
				Name: core.StringPtr("my-instance-template-updated"),
			}
			instanceTemplatePatchModelAsPatch, asPatchErr := instanceTemplatePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceTemplateOptions := &vpcv1.UpdateInstanceTemplateOptions{
				ID:                    core.StringPtr("testString"),
				InstanceTemplatePatch: instanceTemplatePatchModelAsPatch,
			}

			instanceTemplate, response, err := vpcService.UpdateInstanceTemplate(updateInstanceTemplateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceTemplate).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceProfiles - List instance profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceProfiles(listInstanceProfilesOptions *ListInstanceProfilesOptions) with pagination`, func() {
			listInstanceProfilesOptions := &vpcv1.ListInstanceProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listInstanceProfilesOptions.Start = nil
			listInstanceProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceProfile
			for {
				instanceProfileCollection, response, err := vpcService.ListInstanceProfiles(listInstanceProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceProfileCollection).ToNot(BeNil())
				allResults = append(allResults, instanceProfileCollection.Profiles...)

				nextToken, nextErr := instanceProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceProfilesOptions.Start = nextToken

				if listInstanceProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceProfiles(listInstanceProfilesOptions *ListInstanceProfilesOptions) using InstanceProfilesPager`, func() {
			listInstanceProfilesOptions := &vpcv1.ListInstanceProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceProfilesPager(listInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceProfilesPager(listInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceProfiles() returned a total of %d item(s) using InstanceProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetInstanceProfile - Retrieve an instance profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceProfile(getInstanceProfileOptions *GetInstanceProfileOptions)`, func() {
			getInstanceProfileOptions := &vpcv1.GetInstanceProfileOptions{
				Name: core.StringPtr("mx2-host-152x1216"),
			}

			instanceProfile, response, err := vpcService.GetInstanceProfile(getInstanceProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceProfile).ToNot(BeNil())
		})
	})

	Describe(`ListInstances - List instances`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstances(listInstancesOptions *ListInstancesOptions) with pagination`, func() {
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
				Start:                                   core.StringPtr("testString"),
				VPCCRN:                                  core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:                                   core.StringPtr("testString"),
				VPCName:                                 core.StringPtr("my-vpc"),
			}

			listInstancesOptions.Start = nil
			listInstancesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Instance
			for {
				instanceCollection, response, err := vpcService.ListInstances(listInstancesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceCollection).ToNot(BeNil())
				allResults = append(allResults, instanceCollection.Instances...)

				nextToken, nextErr := instanceCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstancesOptions.Start = nextToken

				if listInstancesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstances(listInstancesOptions *ListInstancesOptions) using InstancesPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewInstancesPager(listInstancesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Instance
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstancesPager(listInstancesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstances() returned a total of %d item(s) using InstancesPager.\n", len(allResults))
		})
	})

	Describe(`CreateInstance - Create an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstance(createInstanceOptions *CreateInstanceOptions)`, func() {
			instanceAvailabilityPrototypeModel := &vpcv1.InstanceAvailabilityPrototype{
				Class: core.StringPtr("standard"),
			}

			instanceAvailabilityPolicyPrototypeModel := &vpcv1.InstanceAvailabilityPolicyPrototype{
				HostFailure: core.StringPtr("restart"),
				Preemption:  core.StringPtr("stop"),
			}

			instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"),
			}

			instanceClusterNetworkAttachmentPrototypeInstanceContextModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{
				ClusterNetworkInterface: instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel,
				Name:                    core.StringPtr("my-instance-network-attachment"),
			}

			instanceDefaultTrustedProfilePrototypeTargetModel := &vpcv1.InstanceDefaultTrustedProfilePrototypeTargetTrustedProfileIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:iam-identity::a/aa2432b1fa4d4ace891e9b80fc104e34::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"),
			}

			instanceDefaultTrustedProfilePrototypeModel := &vpcv1.InstanceDefaultTrustedProfilePrototype{
				AutoLink: core.BoolPtr(false),
				Target:   instanceDefaultTrustedProfilePrototypeTargetModel,
			}

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			instanceMetadataServicePrototypeModel := &vpcv1.InstanceMetadataServicePrototype{
				Enabled:          core.BoolPtr(false),
				Protocol:         core.StringPtr("https"),
				ResponseHopLimit: core.Int64Ptr(int64(2)),
			}

			instancePlacementTargetPrototypeModel := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentityDedicatedHostGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host-group:0717-bcc5b834-1258-4b9c-c3b4-43bc9cf5cde0"),
			}

			instancePrototypeProfileModel := &vpcv1.InstancePrototypeProfileInstanceProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"),
			}

			reservationIdentityModel := &vpcv1.ReservationIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
			}

			instanceReservationAffinityPrototypeModel := &vpcv1.InstanceReservationAffinityPrototype{
				Policy: core.StringPtr("automatic"),
				Pool:   []vpcv1.ReservationIdentityIntf{reservationIdentityModel},
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			instanceVcpuPrototypeModel := &vpcv1.InstanceVcpuPrototype{
				Percentage: core.Int64Ptr(int64(100)),
			}

			volumeAttachmentPrototypeVolumeModel := &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentityVolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			volumeAttachmentPrototypeModel := &vpcv1.VolumeAttachmentPrototype{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(false),
				Name:                         core.StringPtr("my-volume-attachment"),
				Volume:                       volumeAttachmentPrototypeVolumeModel,
			}

			instancePrototypeVPCModel := &vpcv1.InstancePrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			volumeAllowedUsePrototypeModel := &vpcv1.VolumeAllowedUsePrototype{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			volumePrototypeInstanceByImageContextEncryptionKeyModel := &vpcv1.VolumePrototypeInstanceByImageContextEncryptionKeyEncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			volumePrototypeInstanceByImageContextProfileModel := &vpcv1.VolumePrototypeInstanceByImageContextProfileVolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumePrototypeInstanceByImageContextResourceGroupModel := &vpcv1.VolumePrototypeInstanceByImageContextResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			volumeAttachmentPrototypeInstanceByImageContextVolumeModel := &vpcv1.VolumeAttachmentPrototypeInstanceByImageContextVolume{
				AllowedUse:    volumeAllowedUsePrototypeModel,
				Bandwidth:     core.Int64Ptr(int64(1000)),
				Capacity:      core.Int64Ptr(int64(100)),
				EncryptionKey: volumePrototypeInstanceByImageContextEncryptionKeyModel,
				Iops:          core.Int64Ptr(int64(10000)),
				Name:          core.StringPtr("my-volume"),
				Profile:       volumePrototypeInstanceByImageContextProfileModel,
				ResourceGroup: volumePrototypeInstanceByImageContextResourceGroupModel,
				UserTags:      []string{"testString"},
			}

			instanceByCatalogOfferingBootVolumeAttachmentModel := &vpcv1.InstanceByCatalogOfferingBootVolumeAttachment{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(true),
				Name:                         core.StringPtr("my-volume-attachment"),
				Volume:                       volumeAttachmentPrototypeInstanceByImageContextVolumeModel,
			}

			catalogOfferingVersionPlanIdentityModel := &vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:51c9e0db-2911-45a6-adb0-ac5332d27cf2:plan:sw.51c9e0db-2911-45a6-adb0-ac5332d27cf2.772c0dbe-aa62-482e-adbe-a3fc20101e0e"),
			}

			catalogOfferingIdentityModel := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:globalcatalog-collection:global:a/aa2432b1fa4d4ace891e9b80fc104e34:1082e7d2-5e2f-0a11-a3bc-f88a8e1931fc:offering:00111601-0ec5-41ac-b142-96d1e64e6442"),
			}

			instanceCatalogOfferingPrototypeModel := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
				Plan:     catalogOfferingVersionPlanIdentityModel,
				Offering: catalogOfferingIdentityModel,
			}

			instanceByCatalogOfferingZoneModel := &vpcv1.InstanceByCatalogOfferingZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			instanceNetworkAttachmentPrototypeModel := &vpcv1.InstanceNetworkAttachmentPrototype{
				Name:                    core.StringPtr("my-instance-network-attachment"),
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instanceByNetworkAttachmentPrimaryNetworkAttachmentModel := &vpcv1.InstanceByNetworkAttachmentPrimaryNetworkAttachment{
				Name:                    core.StringPtr("my-instance-network-attachment"),
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
			}

			instancePrototypeModel := &vpcv1.InstancePrototypeInstanceByCatalogOfferingInstanceByCatalogOfferingInstanceByNetworkAttachment{
				Availability:              instanceAvailabilityPrototypeModel,
				AvailabilityPolicy:        instanceAvailabilityPolicyPrototypeModel,
				ClusterNetworkAttachments: []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{*instanceClusterNetworkAttachmentPrototypeInstanceContextModel},
				ConfidentialComputeMode:   core.StringPtr("disabled"),
				DefaultTrustedProfile:     instanceDefaultTrustedProfilePrototypeModel,
				EnableSecureBoot:          core.BoolPtr(true),
				Keys:                      []vpcv1.KeyIdentityIntf{keyIdentityModel},
				MetadataService:           instanceMetadataServicePrototypeModel,
				Name:                      core.StringPtr("my-instance"),
				PlacementTarget:           instancePlacementTargetPrototypeModel,
				Profile:                   instancePrototypeProfileModel,
				ReservationAffinity:       instanceReservationAffinityPrototypeModel,
				ResourceGroup:             resourceGroupIdentityModel,
				ThreadsPerCore:            core.Int64Ptr(int64(1)),
				TotalVolumeBandwidth:      core.Int64Ptr(int64(500)),
				UserData:                  core.StringPtr("[...]"),
				Vcpu:                      instanceVcpuPrototypeModel,
				VolumeAttachments:         []vpcv1.VolumeAttachmentPrototype{*volumeAttachmentPrototypeModel},
				VolumeBandwidthQosMode:    core.StringPtr("pooled"),
				VPC:                       instancePrototypeVPCModel,
				BootVolumeAttachment:      instanceByCatalogOfferingBootVolumeAttachmentModel,
				CatalogOffering:           instanceCatalogOfferingPrototypeModel,
				Zone:                      instanceByCatalogOfferingZoneModel,
				NetworkAttachments:        []vpcv1.InstanceNetworkAttachmentPrototype{*instanceNetworkAttachmentPrototypeModel},
				PrimaryNetworkAttachment:  instanceByNetworkAttachmentPrimaryNetworkAttachmentModel,
			}

			createInstanceOptions := &vpcv1.CreateInstanceOptions{
				InstancePrototype: instancePrototypeModel,
			}

			instance, response, err := vpcService.CreateInstance(createInstanceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instance).ToNot(BeNil())
		})
	})

	Describe(`GetInstance - Retrieve an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstance(getInstanceOptions *GetInstanceOptions)`, func() {
			getInstanceOptions := &vpcv1.GetInstanceOptions{
				ID: core.StringPtr("testString"),
			}

			instance, response, err := vpcService.GetInstance(getInstanceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instance).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstance - Update an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstance(updateInstanceOptions *UpdateInstanceOptions)`, func() {
			instanceAvailabilityPatchModel := &vpcv1.InstanceAvailabilityPatch{
				Class: core.StringPtr("standard"),
			}

			instanceAvailabilityPolicyPatchModel := &vpcv1.InstanceAvailabilityPolicyPatch{
				HostFailure: core.StringPtr("restart"),
				Preemption:  core.StringPtr("delete"),
			}

			instanceMetadataServicePatchModel := &vpcv1.InstanceMetadataServicePatch{
				Enabled:          core.BoolPtr(true),
				Protocol:         core.StringPtr("http"),
				ResponseHopLimit: core.Int64Ptr(int64(1)),
			}

			instancePlacementTargetPatchModel := &vpcv1.InstancePlacementTargetPatchDedicatedHostGroupIdentityDedicatedHostGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::dedicated-host-group:0717-bcc5b834-1258-4b9c-c3b4-43bc9cf5cde0"),
			}

			instancePatchProfileModel := &vpcv1.InstancePatchProfileInstanceProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"),
			}

			reservationIdentityModel := &vpcv1.ReservationIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::reservation:0717-ba49df72-37b8-43ac-98da-f8e029de0e63"),
			}

			instanceReservationAffinityPatchModel := &vpcv1.InstanceReservationAffinityPatch{
				Policy: core.StringPtr("automatic"),
				Pool:   []vpcv1.ReservationIdentityIntf{reservationIdentityModel},
			}

			instanceVcpuPatchModel := &vpcv1.InstanceVcpuPatch{
				Percentage: core.Int64Ptr(int64(100)),
			}

			instancePatchModel := &vpcv1.InstancePatch{
				Availability:            instanceAvailabilityPatchModel,
				AvailabilityPolicy:      instanceAvailabilityPolicyPatchModel,
				ConfidentialComputeMode: core.StringPtr("disabled"),
				EnableSecureBoot:        core.BoolPtr(true),
				MetadataService:         instanceMetadataServicePatchModel,
				Name:                    core.StringPtr("my-instance"),
				PlacementTarget:         instancePlacementTargetPatchModel,
				Profile:                 instancePatchProfileModel,
				ReservationAffinity:     instanceReservationAffinityPatchModel,
				ThreadsPerCore:          core.Int64Ptr(int64(1)),
				TotalVolumeBandwidth:    core.Int64Ptr(int64(500)),
				Vcpu:                    instanceVcpuPatchModel,
				VolumeBandwidthQosMode:  core.StringPtr("pooled"),
			}
			instancePatchModelAsPatch, asPatchErr := instancePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceOptions := &vpcv1.UpdateInstanceOptions{
				ID:            core.StringPtr("testString"),
				InstancePatch: instancePatchModelAsPatch,
				IfMatch:       core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			instance, response, err := vpcService.UpdateInstance(updateInstanceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instance).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceInitialization - Retrieve initialization configuration for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceInitialization(getInstanceInitializationOptions *GetInstanceInitializationOptions)`, func() {
			getInstanceInitializationOptions := &vpcv1.GetInstanceInitializationOptions{
				ID: core.StringPtr("testString"),
			}

			instanceInitialization, response, err := vpcService.GetInstanceInitialization(getInstanceInitializationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceInitialization).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceReinitialization - Reinitialize an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceReinitialization(createInstanceReinitializationOptions *CreateInstanceReinitializationOptions)`, func() {
			instanceDefaultTrustedProfilePrototypeTargetModel := &vpcv1.InstanceDefaultTrustedProfilePrototypeTargetTrustedProfileIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:iam-identity::a/aa2432b1fa4d4ace891e9b80fc104e34::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"),
			}

			instanceReinitializePrototypeDefaultTrustedProfileModel := &vpcv1.InstanceReinitializePrototypeDefaultTrustedProfile{
				AutoLink: core.BoolPtr(false),
				Target:   instanceDefaultTrustedProfilePrototypeTargetModel,
			}

			keyIdentityModel := &vpcv1.KeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"),
			}

			volumeAllowedUsePrototypeModel := &vpcv1.VolumeAllowedUsePrototype{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			volumePrototypeInstanceByImageContextEncryptionKeyModel := &vpcv1.VolumePrototypeInstanceByImageContextEncryptionKeyEncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			volumePrototypeInstanceByImageContextProfileModel := &vpcv1.VolumePrototypeInstanceByImageContextProfileVolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumePrototypeInstanceByImageContextResourceGroupModel := &vpcv1.VolumePrototypeInstanceByImageContextResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			volumeAttachmentPrototypeInstanceByImageContextVolumeModel := &vpcv1.VolumeAttachmentPrototypeInstanceByImageContextVolume{
				AllowedUse:    volumeAllowedUsePrototypeModel,
				Bandwidth:     core.Int64Ptr(int64(1000)),
				Capacity:      core.Int64Ptr(int64(100)),
				EncryptionKey: volumePrototypeInstanceByImageContextEncryptionKeyModel,
				Iops:          core.Int64Ptr(int64(10000)),
				Name:          core.StringPtr("my-volume"),
				Profile:       volumePrototypeInstanceByImageContextProfileModel,
				ResourceGroup: volumePrototypeInstanceByImageContextResourceGroupModel,
				UserTags:      []string{"testString"},
			}

			instanceReinitializeByImageBootVolumeAttachmentModel := &vpcv1.InstanceReinitializeByImageBootVolumeAttachment{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(true),
				Name:                         core.StringPtr("my-volume-attachment"),
				Volume:                       volumeAttachmentPrototypeInstanceByImageContextVolumeModel,
			}

			instanceReinitializeByImageImageModel := &vpcv1.InstanceReinitializeByImageImageImageIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"),
			}

			instanceReinitializePrototypeModel := &vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage{
				DefaultTrustedProfile: instanceReinitializePrototypeDefaultTrustedProfileModel,
				Keys:                  []vpcv1.KeyIdentityIntf{keyIdentityModel},
				UserData:              core.StringPtr("[...]"),
				BootVolumeAttachment:  instanceReinitializeByImageBootVolumeAttachmentModel,
				Image:                 instanceReinitializeByImageImageModel,
			}

			createInstanceReinitializationOptions := &vpcv1.CreateInstanceReinitializationOptions{
				ID:                            core.StringPtr("testString"),
				InstanceReinitializePrototype: instanceReinitializePrototypeModel,
			}

			response, err := vpcService.CreateInstanceReinitialization(createInstanceReinitializationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`CreateInstanceAction - Create an instance action`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceAction(createInstanceActionOptions *CreateInstanceActionOptions)`, func() {
			createInstanceActionOptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: core.StringPtr("testString"),
				Type:       core.StringPtr("reboot"),
				Force:      core.BoolPtr(true),
			}

			instanceAction, response, err := vpcService.CreateInstanceAction(createInstanceActionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceAction).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceClusterNetworkAttachments - List cluster network attachments on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceClusterNetworkAttachments(listInstanceClusterNetworkAttachmentsOptions *ListInstanceClusterNetworkAttachmentsOptions) with pagination`, func() {
			listInstanceClusterNetworkAttachmentsOptions := &vpcv1.ListInstanceClusterNetworkAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
				Limit:      core.Int64Ptr(int64(10)),
				Start:      core.StringPtr("testString"),
			}

			listInstanceClusterNetworkAttachmentsOptions.Start = nil
			listInstanceClusterNetworkAttachmentsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceClusterNetworkAttachment
			for {
				instanceClusterNetworkAttachmentCollection, response, err := vpcService.ListInstanceClusterNetworkAttachments(listInstanceClusterNetworkAttachmentsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(instanceClusterNetworkAttachmentCollection).ToNot(BeNil())
				allResults = append(allResults, instanceClusterNetworkAttachmentCollection.ClusterNetworkAttachments...)

				nextToken, nextErr := instanceClusterNetworkAttachmentCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceClusterNetworkAttachmentsOptions.Start = nextToken

				if listInstanceClusterNetworkAttachmentsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceClusterNetworkAttachments(listInstanceClusterNetworkAttachmentsOptions *ListInstanceClusterNetworkAttachmentsOptions) using InstanceClusterNetworkAttachmentsPager`, func() {
			listInstanceClusterNetworkAttachmentsOptions := &vpcv1.ListInstanceClusterNetworkAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
				Limit:      core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceClusterNetworkAttachmentsPager(listInstanceClusterNetworkAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceClusterNetworkAttachment
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceClusterNetworkAttachmentsPager(listInstanceClusterNetworkAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceClusterNetworkAttachments() returned a total of %d item(s) using InstanceClusterNetworkAttachmentsPager.\n", len(allResults))
		})
	})

	Describe(`CreateClusterNetworkAttachment - Create a cluster network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateClusterNetworkAttachment(createClusterNetworkAttachmentOptions *CreateClusterNetworkAttachmentOptions)`, func() {
			instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"),
			}

			instanceClusterNetworkAttachmentBeforePrototypeModel := &vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"),
			}

			createClusterNetworkAttachmentOptions := &vpcv1.CreateClusterNetworkAttachmentOptions{
				InstanceID:              core.StringPtr("testString"),
				ClusterNetworkInterface: instanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceModel,
				Before:                  instanceClusterNetworkAttachmentBeforePrototypeModel,
				Name:                    core.StringPtr("my-instance-network-attachment"),
			}

			instanceClusterNetworkAttachment, response, err := vpcService.CreateClusterNetworkAttachment(createClusterNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceClusterNetworkAttachment - Retrieve an instance cluster network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachmentOptions *GetInstanceClusterNetworkAttachmentOptions)`, func() {
			getInstanceClusterNetworkAttachmentOptions := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			instanceClusterNetworkAttachment, response, err := vpcService.GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceClusterNetworkAttachment - Update an instance cluster network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceClusterNetworkAttachment(updateInstanceClusterNetworkAttachmentOptions *UpdateInstanceClusterNetworkAttachmentOptions)`, func() {
			instanceClusterNetworkAttachmentPatchModel := &vpcv1.InstanceClusterNetworkAttachmentPatch{
				Name: core.StringPtr("my-instance-network-attachment-updated"),
			}
			instanceClusterNetworkAttachmentPatchModelAsPatch, asPatchErr := instanceClusterNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceClusterNetworkAttachmentOptions := &vpcv1.UpdateInstanceClusterNetworkAttachmentOptions{
				ID:                                    core.StringPtr("testString"),
				InstanceID:                            core.StringPtr("testString"),
				InstanceClusterNetworkAttachmentPatch: instanceClusterNetworkAttachmentPatchModelAsPatch,
			}

			instanceClusterNetworkAttachment, response, err := vpcService.UpdateInstanceClusterNetworkAttachment(updateInstanceClusterNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceConsoleAccessToken - Create a console access token for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceConsoleAccessToken(createInstanceConsoleAccessTokenOptions *CreateInstanceConsoleAccessTokenOptions)`, func() {
			createInstanceConsoleAccessTokenOptions := &vpcv1.CreateInstanceConsoleAccessTokenOptions{
				InstanceID:  core.StringPtr("testString"),
				ConsoleType: core.StringPtr("serial"),
				Force:       core.BoolPtr(false),
			}

			instanceConsoleAccessToken, response, err := vpcService.CreateInstanceConsoleAccessToken(createInstanceConsoleAccessTokenOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceConsoleAccessToken).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceDisks - List disks on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceDisks(listInstanceDisksOptions *ListInstanceDisksOptions)`, func() {
			listInstanceDisksOptions := &vpcv1.ListInstanceDisksOptions{
				InstanceID: core.StringPtr("testString"),
			}

			instanceDiskCollection, response, err := vpcService.ListInstanceDisks(listInstanceDisksOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDiskCollection).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceDisk - Retrieve an instance disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceDisk(getInstanceDiskOptions *GetInstanceDiskOptions)`, func() {
			getInstanceDiskOptions := &vpcv1.GetInstanceDiskOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			instanceDisk, response, err := vpcService.GetInstanceDisk(getInstanceDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDisk).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceDisk - Update an instance disk`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceDisk(updateInstanceDiskOptions *UpdateInstanceDiskOptions)`, func() {
			instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{
				Name: core.StringPtr("my-instance-disk-updated"),
			}
			instanceDiskPatchModelAsPatch, asPatchErr := instanceDiskPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceDiskOptions := &vpcv1.UpdateInstanceDiskOptions{
				ID:                core.StringPtr("testString"),
				InstanceID:        core.StringPtr("testString"),
				InstanceDiskPatch: instanceDiskPatchModelAsPatch,
			}

			instanceDisk, response, err := vpcService.UpdateInstanceDisk(updateInstanceDiskOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceDisk).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceNetworkAttachments - List network attachments on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceNetworkAttachments(listInstanceNetworkAttachmentsOptions *ListInstanceNetworkAttachmentsOptions)`, func() {
			listInstanceNetworkAttachmentsOptions := &vpcv1.ListInstanceNetworkAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
			}

			instanceNetworkAttachmentCollection, response, err := vpcService.ListInstanceNetworkAttachments(listInstanceNetworkAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachmentCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceNetworkAttachment - Create a network attachment on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceNetworkAttachment(createInstanceNetworkAttachmentOptions *CreateInstanceNetworkAttachmentOptions)`, func() {
			virtualNetworkInterfaceIPPrototypeModel := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			virtualNetworkInterfacePrimaryIPPrototypeModel := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext{
				Address:    core.StringPtr("10.0.0.7"),
				AutoDelete: core.BoolPtr(false),
				Name:       core.StringPtr("my-reserved-ip"),
			}

			virtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextResourceGroupModel := &vpcv1.VirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			virtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextSubnetModel := &vpcv1.VirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextSubnetSubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext{
				AllowIPSpoofing:            core.BoolPtr(true),
				AutoDelete:                 core.BoolPtr(false),
				EnableInfrastructureNat:    core.BoolPtr(true),
				Ips:                        []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{virtualNetworkInterfaceIPPrototypeModel},
				Name:                       core.StringPtr("my-virtual-network-interface"),
				PrimaryIP:                  virtualNetworkInterfacePrimaryIPPrototypeModel,
				ProtocolStateFilteringMode: core.StringPtr("auto"),
				ResourceGroup:              virtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextResourceGroupModel,
				SecurityGroups:             []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
				Subnet:                     virtualNetworkInterfacePrototypeInstanceNetworkAttachmentContextSubnetModel,
			}

			createInstanceNetworkAttachmentOptions := &vpcv1.CreateInstanceNetworkAttachmentOptions{
				InstanceID:              core.StringPtr("testString"),
				VirtualNetworkInterface: instanceNetworkAttachmentPrototypeVirtualNetworkInterfaceModel,
				Name:                    core.StringPtr("testString"),
			}

			instanceNetworkAttachment, response, err := vpcService.CreateInstanceNetworkAttachment(createInstanceNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceNetworkAttachment - Retrieve an instance network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions *GetInstanceNetworkAttachmentOptions)`, func() {
			getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			instanceNetworkAttachment, response, err := vpcService.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceNetworkAttachment - Update an instance network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceNetworkAttachment(updateInstanceNetworkAttachmentOptions *UpdateInstanceNetworkAttachmentOptions)`, func() {
			instanceNetworkAttachmentPatchModel := &vpcv1.InstanceNetworkAttachmentPatch{
				Name: core.StringPtr("my-instance-network-attachment-updated"),
			}
			instanceNetworkAttachmentPatchModelAsPatch, asPatchErr := instanceNetworkAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceNetworkAttachmentOptions := &vpcv1.UpdateInstanceNetworkAttachmentOptions{
				ID:                             core.StringPtr("testString"),
				InstanceID:                     core.StringPtr("testString"),
				InstanceNetworkAttachmentPatch: instanceNetworkAttachmentPatchModelAsPatch,
			}

			instanceNetworkAttachment, response, err := vpcService.UpdateInstanceNetworkAttachment(updateInstanceNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceNetworkInterfaces - List network interfaces on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceNetworkInterfaces(listInstanceNetworkInterfacesOptions *ListInstanceNetworkInterfacesOptions)`, func() {
			listInstanceNetworkInterfacesOptions := &vpcv1.ListInstanceNetworkInterfacesOptions{
				InstanceID: core.StringPtr("testString"),
			}

			networkInterfaceUnpaginatedCollection, response, err := vpcService.ListInstanceNetworkInterfaces(listInstanceNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterfaceUnpaginatedCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceNetworkInterface - Create a network interface on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceNetworkInterface(createInstanceNetworkInterfaceOptions *CreateInstanceNetworkInterfaceOptions)`, func() {
			networkInterfacePrototypeSubnetModel := &vpcv1.NetworkInterfacePrototypeSubnetSubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			networkInterfaceIPPrototypeModel := &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			createInstanceNetworkInterfaceOptions := &vpcv1.CreateInstanceNetworkInterfaceOptions{
				InstanceID:      core.StringPtr("testString"),
				Subnet:          networkInterfacePrototypeSubnetModel,
				AllowIPSpoofing: core.BoolPtr(true),
				Name:            core.StringPtr("my-instance-network-interface"),
				PrimaryIP:       networkInterfaceIPPrototypeModel,
				SecurityGroups:  []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
			}

			networkInterface, response, err := vpcService.CreateInstanceNetworkInterface(createInstanceNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkInterface).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceNetworkInterface - Retrieve an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceNetworkInterface(getInstanceNetworkInterfaceOptions *GetInstanceNetworkInterfaceOptions)`, func() {
			getInstanceNetworkInterfaceOptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			networkInterface, response, err := vpcService.GetInstanceNetworkInterface(getInstanceNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterface).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceNetworkInterface - Update an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceNetworkInterface(updateInstanceNetworkInterfaceOptions *UpdateInstanceNetworkInterfaceOptions)`, func() {
			networkInterfacePatchModel := &vpcv1.NetworkInterfacePatch{
				AllowIPSpoofing: core.BoolPtr(true),
				Name:            core.StringPtr("my-network-interface-updated"),
			}
			networkInterfacePatchModelAsPatch, asPatchErr := networkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceNetworkInterfaceOptions := &vpcv1.UpdateInstanceNetworkInterfaceOptions{
				ID:                    core.StringPtr("testString"),
				InstanceID:            core.StringPtr("testString"),
				NetworkInterfacePatch: networkInterfacePatchModelAsPatch,
			}

			networkInterface, response, err := vpcService.UpdateInstanceNetworkInterface(updateInstanceNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkInterface).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceNetworkInterfaceFloatingIps - List floating IPs associated with an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceNetworkInterfaceFloatingIps(listInstanceNetworkInterfaceFloatingIpsOptions *ListInstanceNetworkInterfaceFloatingIpsOptions)`, func() {
			listInstanceNetworkInterfaceFloatingIpsOptions := &vpcv1.ListInstanceNetworkInterfaceFloatingIpsOptions{
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIPUnpaginatedCollection, response, err := vpcService.ListInstanceNetworkInterfaceFloatingIps(listInstanceNetworkInterfaceFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPUnpaginatedCollection).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceNetworkInterfaceFloatingIP - Retrieve associated floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceNetworkInterfaceFloatingIP(getInstanceNetworkInterfaceFloatingIPOptions *GetInstanceNetworkInterfaceFloatingIPOptions)`, func() {
			getInstanceNetworkInterfaceFloatingIPOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
				ID:                 core.StringPtr("testString"),
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIP, response, err := vpcService.GetInstanceNetworkInterfaceFloatingIP(getInstanceNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`AddInstanceNetworkInterfaceFloatingIP - Associate a floating IP with an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddInstanceNetworkInterfaceFloatingIP(addInstanceNetworkInterfaceFloatingIPOptions *AddInstanceNetworkInterfaceFloatingIPOptions)`, func() {
			addInstanceNetworkInterfaceFloatingIPOptions := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
				ID:                 core.StringPtr("testString"),
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIP, response, err := vpcService.AddInstanceNetworkInterfaceFloatingIP(addInstanceNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(floatingIP).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceNetworkInterfaceIps - List the primary reserved IP for an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceNetworkInterfaceIps(listInstanceNetworkInterfaceIpsOptions *ListInstanceNetworkInterfaceIpsOptions) with pagination`, func() {
			listInstanceNetworkInterfaceIpsOptions := &vpcv1.ListInstanceNetworkInterfaceIpsOptions{
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Start:              core.StringPtr("testString"),
			}

			listInstanceNetworkInterfaceIpsOptions.Start = nil
			listInstanceNetworkInterfaceIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ReservedIP
			for {
				reservedIPCollectionInstanceNetworkInterfaceContext, response, err := vpcService.ListInstanceNetworkInterfaceIps(listInstanceNetworkInterfaceIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(reservedIPCollectionInstanceNetworkInterfaceContext).ToNot(BeNil())
				allResults = append(allResults, reservedIPCollectionInstanceNetworkInterfaceContext.Ips...)

				nextToken, nextErr := reservedIPCollectionInstanceNetworkInterfaceContext.GetNextStart()
				Expect(nextErr).To(BeNil())
				listInstanceNetworkInterfaceIpsOptions.Start = nextToken

				if listInstanceNetworkInterfaceIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListInstanceNetworkInterfaceIps(listInstanceNetworkInterfaceIpsOptions *ListInstanceNetworkInterfaceIpsOptions) using InstanceNetworkInterfaceIpsPager`, func() {
			listInstanceNetworkInterfaceIpsOptions := &vpcv1.ListInstanceNetworkInterfaceIpsOptions{
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewInstanceNetworkInterfaceIpsPager(listInstanceNetworkInterfaceIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewInstanceNetworkInterfaceIpsPager(listInstanceNetworkInterfaceIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListInstanceNetworkInterfaceIps() returned a total of %d item(s) using InstanceNetworkInterfaceIpsPager.\n", len(allResults))
		})
	})

	Describe(`GetInstanceNetworkInterfaceIP - Retrieve the primary reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceNetworkInterfaceIP(getInstanceNetworkInterfaceIPOptions *GetInstanceNetworkInterfaceIPOptions)`, func() {
			getInstanceNetworkInterfaceIPOptions := &vpcv1.GetInstanceNetworkInterfaceIPOptions{
				ID:                 core.StringPtr("testString"),
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			reservedIP, response, err := vpcService.GetInstanceNetworkInterfaceIP(getInstanceNetworkInterfaceIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceSoftwareAttachments - List instance software attachments associated with an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceSoftwareAttachments(listInstanceSoftwareAttachmentsOptions *ListInstanceSoftwareAttachmentsOptions)`, func() {
			listInstanceSoftwareAttachmentsOptions := &vpcv1.ListInstanceSoftwareAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
			}

			instanceSoftwareAttachmentCollection, response, err := vpcService.ListInstanceSoftwareAttachments(listInstanceSoftwareAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachmentCollection).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceSoftwareAttachment - Retrieve an instance software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceSoftwareAttachment(getInstanceSoftwareAttachmentOptions *GetInstanceSoftwareAttachmentOptions)`, func() {
			getInstanceSoftwareAttachmentOptions := &vpcv1.GetInstanceSoftwareAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			instanceSoftwareAttachment, response, err := vpcService.GetInstanceSoftwareAttachment(getInstanceSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceSoftwareAttachment - Update an instance software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceSoftwareAttachment(updateInstanceSoftwareAttachmentOptions *UpdateInstanceSoftwareAttachmentOptions)`, func() {
			instanceSoftwareAttachmentPatchModel := &vpcv1.InstanceSoftwareAttachmentPatch{
				Name: core.StringPtr("my-software-attachment-patch"),
			}
			instanceSoftwareAttachmentPatchModelAsPatch, asPatchErr := instanceSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceSoftwareAttachmentOptions := &vpcv1.UpdateInstanceSoftwareAttachmentOptions{
				ID:                              core.StringPtr("testString"),
				InstanceID:                      core.StringPtr("testString"),
				InstanceSoftwareAttachmentPatch: instanceSoftwareAttachmentPatchModelAsPatch,
			}

			instanceSoftwareAttachment, response, err := vpcService.UpdateInstanceSoftwareAttachment(updateInstanceSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(instanceSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListInstanceVolumeAttachments - List volumes attachments on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListInstanceVolumeAttachments(listInstanceVolumeAttachmentsOptions *ListInstanceVolumeAttachmentsOptions)`, func() {
			listInstanceVolumeAttachmentsOptions := &vpcv1.ListInstanceVolumeAttachmentsOptions{
				InstanceID: core.StringPtr("testString"),
			}

			volumeAttachmentCollection, response, err := vpcService.ListInstanceVolumeAttachments(listInstanceVolumeAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachmentCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateInstanceVolumeAttachment - Create a volume attachment on an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateInstanceVolumeAttachment(createInstanceVolumeAttachmentOptions *CreateInstanceVolumeAttachmentOptions)`, func() {
			volumeAttachmentPrototypeVolumeModel := &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentityVolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			createInstanceVolumeAttachmentOptions := &vpcv1.CreateInstanceVolumeAttachmentOptions{
				InstanceID:                   core.StringPtr("testString"),
				Volume:                       volumeAttachmentPrototypeVolumeModel,
				DeleteVolumeOnInstanceDelete: core.BoolPtr(false),
				Name:                         core.StringPtr("my-volume-attachment"),
			}

			volumeAttachment, response, err := vpcService.CreateInstanceVolumeAttachment(createInstanceVolumeAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volumeAttachment).ToNot(BeNil())
		})
	})

	Describe(`GetInstanceVolumeAttachment - Retrieve a volume attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetInstanceVolumeAttachment(getInstanceVolumeAttachmentOptions *GetInstanceVolumeAttachmentOptions)`, func() {
			getInstanceVolumeAttachmentOptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			volumeAttachment, response, err := vpcService.GetInstanceVolumeAttachment(getInstanceVolumeAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateInstanceVolumeAttachment - Update a volume attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateInstanceVolumeAttachment(updateInstanceVolumeAttachmentOptions *UpdateInstanceVolumeAttachmentOptions)`, func() {
			volumeAttachmentPatchModel := &vpcv1.VolumeAttachmentPatch{
				DeleteVolumeOnInstanceDelete: core.BoolPtr(true),
				Name:                         core.StringPtr("my-volume-attachment"),
			}
			volumeAttachmentPatchModelAsPatch, asPatchErr := volumeAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateInstanceVolumeAttachmentOptions := &vpcv1.UpdateInstanceVolumeAttachmentOptions{
				ID:                    core.StringPtr("testString"),
				InstanceID:            core.StringPtr("testString"),
				VolumeAttachmentPatch: volumeAttachmentPatchModelAsPatch,
			}

			volumeAttachment, response, err := vpcService.UpdateInstanceVolumeAttachment(updateInstanceVolumeAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListKeys - List keys`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListKeys(listKeysOptions *ListKeysOptions) with pagination`, func() {
			listKeysOptions := &vpcv1.ListKeysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listKeysOptions.Start = nil
			listKeysOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Key
			for {
				keyCollection, response, err := vpcService.ListKeys(listKeysOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(keyCollection).ToNot(BeNil())
				allResults = append(allResults, keyCollection.Keys...)

				nextToken, nextErr := keyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listKeysOptions.Start = nextToken

				if listKeysOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListKeys(listKeysOptions *ListKeysOptions) using KeysPager`, func() {
			listKeysOptions := &vpcv1.ListKeysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewKeysPager(listKeysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Key
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewKeysPager(listKeysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListKeys() returned a total of %d item(s) using KeysPager.\n", len(allResults))
		})
	})

	Describe(`CreateKey - Create a key`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateKey(createKeyOptions *CreateKeyOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createKeyOptions := &vpcv1.CreateKeyOptions{
				PublicKey:     core.StringPtr("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDDGe50Bxa5T5NDddrrtbx2Y4/VGbiCgXqnBsYToIUKoFSHTQl5IX3PasGnneKanhcLwWz5M5MoCRvhxTp66NKzIfAz7r+FX9rxgR+ZgcM253YAqOVeIpOU408simDZKriTlN8kYsXL7P34tsWuAJf4MgZtJAQxous/2byetpdCv8ddnT4X3ltOg9w+LqSCPYfNivqH00Eh7S1Ldz7I8aw5WOp5a+sQFP/RbwfpwHp+ny7DfeIOokcuI42tJkoBn7UsLTVpCSmXr2EDRlSWe/1M/iHNRBzaT3CK0+SwZWd2AEjePxSnWKNGIEUJDlUYp7hKhiQcgT5ZAnWU121oc5En"),
				Name:          core.StringPtr("my-key-1"),
				ResourceGroup: resourceGroupIdentityModel,
				Type:          core.StringPtr("rsa"),
			}

			key, response, err := vpcService.CreateKey(createKeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(key).ToNot(BeNil())
		})
	})

	Describe(`GetKey - Retrieve a key`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetKey(getKeyOptions *GetKeyOptions)`, func() {
			getKeyOptions := &vpcv1.GetKeyOptions{
				ID: core.StringPtr("testString"),
			}

			key, response, err := vpcService.GetKey(getKeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(key).ToNot(BeNil())
		})
	})

	Describe(`UpdateKey - Update a key`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateKey(updateKeyOptions *UpdateKeyOptions)`, func() {
			keyPatchModel := &vpcv1.KeyPatch{
				Name: core.StringPtr("my-key-1"),
			}
			keyPatchModelAsPatch, asPatchErr := keyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateKeyOptions := &vpcv1.UpdateKeyOptions{
				ID:       core.StringPtr("testString"),
				KeyPatch: keyPatchModelAsPatch,
			}

			key, response, err := vpcService.UpdateKey(updateKeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(key).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerProfiles - List load balancer profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerProfiles(listLoadBalancerProfilesOptions *ListLoadBalancerProfilesOptions) with pagination`, func() {
			listLoadBalancerProfilesOptions := &vpcv1.ListLoadBalancerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listLoadBalancerProfilesOptions.Start = nil
			listLoadBalancerProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.LoadBalancerProfile
			for {
				loadBalancerProfileCollection, response, err := vpcService.ListLoadBalancerProfiles(listLoadBalancerProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(loadBalancerProfileCollection).ToNot(BeNil())
				allResults = append(allResults, loadBalancerProfileCollection.Profiles...)

				nextToken, nextErr := loadBalancerProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listLoadBalancerProfilesOptions.Start = nextToken

				if listLoadBalancerProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListLoadBalancerProfiles(listLoadBalancerProfilesOptions *ListLoadBalancerProfilesOptions) using LoadBalancerProfilesPager`, func() {
			listLoadBalancerProfilesOptions := &vpcv1.ListLoadBalancerProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewLoadBalancerProfilesPager(listLoadBalancerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.LoadBalancerProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewLoadBalancerProfilesPager(listLoadBalancerProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListLoadBalancerProfiles() returned a total of %d item(s) using LoadBalancerProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetLoadBalancerProfile - Retrieve a load balancer profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerProfile(getLoadBalancerProfileOptions *GetLoadBalancerProfileOptions)`, func() {
			getLoadBalancerProfileOptions := &vpcv1.GetLoadBalancerProfileOptions{
				Name: core.StringPtr("network-fixed"),
			}

			loadBalancerProfile, response, err := vpcService.GetLoadBalancerProfile(getLoadBalancerProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerProfile).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancers - List load balancers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancers(listLoadBalancersOptions *ListLoadBalancersOptions) with pagination`, func() {
			listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listLoadBalancersOptions.Start = nil
			listLoadBalancersOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.LoadBalancer
			for {
				loadBalancerCollection, response, err := vpcService.ListLoadBalancers(listLoadBalancersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(loadBalancerCollection).ToNot(BeNil())
				allResults = append(allResults, loadBalancerCollection.LoadBalancers...)

				nextToken, nextErr := loadBalancerCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listLoadBalancersOptions.Start = nextToken

				if listLoadBalancersOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListLoadBalancers(listLoadBalancersOptions *ListLoadBalancersOptions) using LoadBalancersPager`, func() {
			listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewLoadBalancersPager(listLoadBalancersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.LoadBalancer
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewLoadBalancersPager(listLoadBalancersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListLoadBalancers() returned a total of %d item(s) using LoadBalancersPager.\n", len(allResults))
		})
	})

	Describe(`CreateLoadBalancer - Create a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancer(createLoadBalancerOptions *CreateLoadBalancerOptions)`, func() {
			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			loadBalancerDnsPrototypeInstanceModel := &vpcv1.LoadBalancerDnsPrototypeInstanceDnsInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:dns-svcs:global:a/bb1b52262f7441a586f49068482f1e60:f761b566-030a-4696-8649-cc9d09889e88::"),
			}

			loadBalancerDnsPrototypeZoneModel := &vpcv1.LoadBalancerDnsPrototypeZoneDnsZoneIdentityByID{
				ID: core.StringPtr("d66662cc-aa23-4fe1-9987-858487a61f45"),
			}

			loadBalancerDnsPrototypeModel := &vpcv1.LoadBalancerDnsPrototype{
				Instance: loadBalancerDnsPrototypeInstanceModel,
				Zone:     loadBalancerDnsPrototypeZoneModel,
			}

			loadBalancerListenerPrototypeLoadBalancerContextCertificateInstanceModel := &vpcv1.LoadBalancerListenerPrototypeLoadBalancerContextCertificateInstanceCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			certificateInstanceIdentityModel := &vpcv1.CertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerListenerClientAuthenticationPrototypeModel := &vpcv1.LoadBalancerListenerClientAuthenticationPrototype{
				CertificateAuthority:      certificateInstanceIdentityModel,
				CertificateRevocationList: core.StringPtr("-----BEGIN X509 CRL-----\nMIICvTCBpgIBATANBgkqhkiG9w0BAQsFADBMMQswCQYDVQQGEwJVUzEOMAwGA1UE\nCAwFZGVsYXMxDDAKBgNVBAoMA0lCTTENMAsGA1UECwwEcm9vdDEQMA4GA1UEAwwH\ncm9vdC1jYRcNMjUwOTA4MDUwMjQwWhcNMjUxMDA4MDUwMjQwWjAVMBMCAhAAFw0y\nNTA5MDgwNTAxNTlaoA8wDTALBgNVHRQEBAICEAAwDQYJKoZIhvcNAQELBQADggIB\nACeEcj7ompUepc5qTvTrNA5PoK5bN71gNI7Rbhq/Bxf1YPMp2iU3qMSj7YpVP7aw\nGNrxFoIZcQ4X7PYyHMfDk6Z83PSTVMnSOVk09fZW49tyVTWmzBVLz3R1bPasnWTZ\n0hRIv9j9n7Lemin+0ubIR/2zmsfBs1JFAFEbbRcgwg+qotsfZNLkX6bjHDpsRQzE\nmXUEu4/AqAsWPbFzG2uMKZ9pKOK+Nn3bt/NEK+AFlnSmgjEqzQ+0zhsrCExIReJV\nc2oiLBkLG6rBwxlGDog+PqwjP+1wGNIL1J3c2lMW1IGMNcts/aDBO5LtPVIY1LsQ\nFoeaTfm3U3GKC/pTczoDk/pKN756f8O05nTWUHgktcNsPvgqDKnpvEkI3VPf9Y4a\nfMOzKgVTgY1dSgjzHO8+4ZfcVGpBePsjOe0/RCUwkgtgOyGtcmBPTMJa0elJzjaM\njD9myqIXkB359sqbuEmcrjgo5uUUvubFYpmT/W0YxOi/py/bDK+7uUs38nUElNkZ\n+YFRpNWjLF9JtAghX5MhA5BwhTTuATvWYuDdK769ifi9qcYvE4u+VNxYfOpPY6sv\nx4FnkZ9+A7s2hk11d+DEq29Efa0xak8rO1LzT5hCSFT0P3KfZEZMpbuXpzVGiZoM\ng5cWHgYcNnzhUatKodvzZizAOVGRR7UFg42O4ylhxDVe\n-----END X509 CRL-----\n"),
			}

			loadBalancerListenerPrototypeLoadBalancerContextDefaultPoolModel := &vpcv1.LoadBalancerListenerPrototypeLoadBalancerContextDefaultPool{
				Name: core.StringPtr("my-load-balancer-pool"),
			}

			loadBalancerListenerHTTPSRedirectPrototypeListenerModel := &vpcv1.LoadBalancerListenerHTTPSRedirectPrototypeListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			loadBalancerListenerPrototypeLoadBalancerContextHTTPSRedirectModel := &vpcv1.LoadBalancerListenerPrototypeLoadBalancerContextHTTPSRedirect{
				HTTPStatusCode: core.Int64Ptr(int64(301)),
				Listener:       loadBalancerListenerHTTPSRedirectPrototypeListenerModel,
				URI:            core.StringPtr("/example?doc=get"),
			}

			loadBalancerListenerPrototypeLoadBalancerContextModel := &vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
				AcceptProxyProtocol:   core.BoolPtr(true),
				CertificateInstance:   loadBalancerListenerPrototypeLoadBalancerContextCertificateInstanceModel,
				ClientAuthentication:  loadBalancerListenerClientAuthenticationPrototypeModel,
				ConnectionLimit:       core.Int64Ptr(int64(2000)),
				DefaultPool:           loadBalancerListenerPrototypeLoadBalancerContextDefaultPoolModel,
				HTTPSRedirect:         loadBalancerListenerPrototypeLoadBalancerContextHTTPSRedirectModel,
				IdleConnectionTimeout: core.Int64Ptr(int64(100)),
				Port:                  core.Int64Ptr(int64(443)),
				PortMax:               core.Int64Ptr(int64(499)),
				PortMin:               core.Int64Ptr(int64(443)),
				Protocol:              core.StringPtr("http"),
			}

			loadBalancerLoggingDatapathPrototypeModel := &vpcv1.LoadBalancerLoggingDatapathPrototype{
				Active: core.BoolPtr(true),
			}

			loadBalancerLoggingPrototypeModel := &vpcv1.LoadBalancerLoggingPrototype{
				Datapath: loadBalancerLoggingDatapathPrototypeModel,
			}

			loadBalancerPoolClientAuthenticationPrototypeModel := &vpcv1.LoadBalancerPoolClientAuthenticationPrototype{
				CertificateInstance: certificateInstanceIdentityModel,
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{
				Field: core.StringPtr("Content-Type"),
				Value: core.StringPtr("text/plain"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeLoadBalancerPoolHealthMonitorTypeHttphttpsRequestGetPrototype{
				HeadersVar: []vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{*loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel},
				Method:     core.StringPtr("get"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototype{
				BodyRegex: core.StringPtr("ACTIVE"),
				Codes:     []string{"200"},
			}

			loadBalancerPoolPrototypeLoadBalancerContextHealthMonitorModel := &vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextHealthMonitorLoadBalancerPoolHealthMonitorPrototypeLoadBalancerPoolHealthMonitorTypeHttphttpsPrototype{
				Delay:      core.Int64Ptr(int64(5)),
				MaxRetries: core.Int64Ptr(int64(2)),
				Port:       core.Int64Ptr(int64(22)),
				Timeout:    core.Int64Ptr(int64(2)),
				Request:    loadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeModel,
				Response:   loadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototypeModel,
				Type:       core.StringPtr("http"),
				URLPath:    core.StringPtr("/"),
			}

			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			loadBalancerPoolMemberPrototypeModel := &vpcv1.LoadBalancerPoolMemberPrototype{
				Port:   core.Int64Ptr(int64(80)),
				Target: loadBalancerPoolMemberTargetPrototypeModel,
				Weight: core.Int64Ptr(int64(50)),
			}

			loadBalancerPoolServerAuthenticationPrototypeModel := &vpcv1.LoadBalancerPoolServerAuthenticationPrototype{
				CertificateAuthority: certificateInstanceIdentityModel,
				VerifyCertificate:    core.BoolPtr(false),
			}

			loadBalancerPoolPrototypeLoadBalancerContextSessionPersistenceModel := &vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextSessionPersistence{
				CookieName: core.StringPtr("my-cookie-name"),
				Type:       core.StringPtr("app_cookie"),
			}

			loadBalancerPoolPrototypeLoadBalancerContextModel := &vpcv1.LoadBalancerPoolPrototypeLoadBalancerContext{
				Algorithm:            core.StringPtr("least_connections"),
				ClientAuthentication: loadBalancerPoolClientAuthenticationPrototypeModel,
				HealthMonitor:        loadBalancerPoolPrototypeLoadBalancerContextHealthMonitorModel,
				Members:              []vpcv1.LoadBalancerPoolMemberPrototype{*loadBalancerPoolMemberPrototypeModel},
				Name:                 core.StringPtr("my-load-balancer-pool"),
				Protocol:             core.StringPtr("http"),
				ProxyProtocol:        core.StringPtr("disabled"),
				ServerAuthentication: loadBalancerPoolServerAuthenticationPrototypeModel,
				SessionPersistence:   loadBalancerPoolPrototypeLoadBalancerContextSessionPersistenceModel,
			}

			loadBalancerPrototypeProfileModel := &vpcv1.LoadBalancerPrototypeProfileLoadBalancerProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancer/profiles/network-fixed"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			createLoadBalancerOptions := &vpcv1.CreateLoadBalancerOptions{
				IsPublic:       core.BoolPtr(true),
				Subnets:        []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
				Dns:            loadBalancerDnsPrototypeModel,
				IsPrivatePath:  core.BoolPtr(true),
				Listeners:      []vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{*loadBalancerListenerPrototypeLoadBalancerContextModel},
				Logging:        loadBalancerLoggingPrototypeModel,
				Name:           core.StringPtr("my-load-balancer"),
				Pools:          []vpcv1.LoadBalancerPoolPrototypeLoadBalancerContext{*loadBalancerPoolPrototypeLoadBalancerContextModel},
				Profile:        loadBalancerPrototypeProfileModel,
				ResourceGroup:  resourceGroupIdentityModel,
				RouteMode:      core.BoolPtr(true),
				SecurityGroups: []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
			}

			loadBalancer, response, err := vpcService.CreateLoadBalancer(createLoadBalancerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancer).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancer - Retrieve a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancer(getLoadBalancerOptions *GetLoadBalancerOptions)`, func() {
			getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
				ID: core.StringPtr("testString"),
			}

			loadBalancer, response, err := vpcService.GetLoadBalancer(getLoadBalancerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancer).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancer - Update a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancer(updateLoadBalancerOptions *UpdateLoadBalancerOptions)`, func() {
			loadBalancerDnsPatchInstanceModel := &vpcv1.LoadBalancerDnsPatchInstanceDnsInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:dns-svcs:global:a/bb1b52262f7441a586f49068482f1e60:f761b566-030a-4696-8649-cc9d09889e88::"),
			}

			loadBalancerDnsPatchZoneModel := &vpcv1.LoadBalancerDnsPatchZoneDnsZoneIdentityByID{
				ID: core.StringPtr("d66662cc-aa23-4fe1-9987-858487a61f45"),
			}

			loadBalancerDnsPatchModel := &vpcv1.LoadBalancerDnsPatch{
				Instance: loadBalancerDnsPatchInstanceModel,
				Zone:     loadBalancerDnsPatchZoneModel,
			}

			loadBalancerLoggingDatapathPatchModel := &vpcv1.LoadBalancerLoggingDatapathPatch{
				Active: core.BoolPtr(true),
			}

			loadBalancerLoggingPatchModel := &vpcv1.LoadBalancerLoggingPatch{
				Datapath: loadBalancerLoggingDatapathPatchModel,
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{
				Dns:     loadBalancerDnsPatchModel,
				Logging: loadBalancerLoggingPatchModel,
				Name:    core.StringPtr("my-load-balancer"),
				Subnets: []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			}
			loadBalancerPatchModelAsPatch, asPatchErr := loadBalancerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerOptions := &vpcv1.UpdateLoadBalancerOptions{
				ID:                core.StringPtr("testString"),
				LoadBalancerPatch: loadBalancerPatchModelAsPatch,
				IfMatch:           core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			loadBalancer, response, err := vpcService.UpdateLoadBalancer(updateLoadBalancerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancer).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerStatistics - List statistics of a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerStatistics(getLoadBalancerStatisticsOptions *GetLoadBalancerStatisticsOptions)`, func() {
			getLoadBalancerStatisticsOptions := &vpcv1.GetLoadBalancerStatisticsOptions{
				ID: core.StringPtr("testString"),
			}

			loadBalancerStatistics, response, err := vpcService.GetLoadBalancerStatistics(getLoadBalancerStatisticsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerStatistics).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerListeners - List listeners for a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerListeners(listLoadBalancerListenersOptions *ListLoadBalancerListenersOptions)`, func() {
			listLoadBalancerListenersOptions := &vpcv1.ListLoadBalancerListenersOptions{
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerListenerCollection, response, err := vpcService.ListLoadBalancerListeners(listLoadBalancerListenersOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateLoadBalancerListener - Create a listener for a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancerListener(createLoadBalancerListenerOptions *CreateLoadBalancerListenerOptions)`, func() {
			loadBalancerListenerPrototypeCertificateInstanceModel := &vpcv1.LoadBalancerListenerPrototypeCertificateInstanceCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			certificateInstanceIdentityModel := &vpcv1.CertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerListenerClientAuthenticationPrototypeModel := &vpcv1.LoadBalancerListenerClientAuthenticationPrototype{
				CertificateAuthority:      certificateInstanceIdentityModel,
				CertificateRevocationList: core.StringPtr("-----BEGIN X509 CRL-----\nMIICvTCBpgIBATANBgkqhkiG9w0BAQsFADBMMQswCQYDVQQGEwJVUzEOMAwGA1UE\nCAwFZGVsYXMxDDAKBgNVBAoMA0lCTTENMAsGA1UECwwEcm9vdDEQMA4GA1UEAwwH\ncm9vdC1jYRcNMjUwOTA4MDUwMjQwWhcNMjUxMDA4MDUwMjQwWjAVMBMCAhAAFw0y\nNTA5MDgwNTAxNTlaoA8wDTALBgNVHRQEBAICEAAwDQYJKoZIhvcNAQELBQADggIB\nACeEcj7ompUepc5qTvTrNA5PoK5bN71gNI7Rbhq/Bxf1YPMp2iU3qMSj7YpVP7aw\nGNrxFoIZcQ4X7PYyHMfDk6Z83PSTVMnSOVk09fZW49tyVTWmzBVLz3R1bPasnWTZ\n0hRIv9j9n7Lemin+0ubIR/2zmsfBs1JFAFEbbRcgwg+qotsfZNLkX6bjHDpsRQzE\nmXUEu4/AqAsWPbFzG2uMKZ9pKOK+Nn3bt/NEK+AFlnSmgjEqzQ+0zhsrCExIReJV\nc2oiLBkLG6rBwxlGDog+PqwjP+1wGNIL1J3c2lMW1IGMNcts/aDBO5LtPVIY1LsQ\nFoeaTfm3U3GKC/pTczoDk/pKN756f8O05nTWUHgktcNsPvgqDKnpvEkI3VPf9Y4a\nfMOzKgVTgY1dSgjzHO8+4ZfcVGpBePsjOe0/RCUwkgtgOyGtcmBPTMJa0elJzjaM\njD9myqIXkB359sqbuEmcrjgo5uUUvubFYpmT/W0YxOi/py/bDK+7uUs38nUElNkZ\n+YFRpNWjLF9JtAghX5MhA5BwhTTuATvWYuDdK769ifi9qcYvE4u+VNxYfOpPY6sv\nx4FnkZ9+A7s2hk11d+DEq29Efa0xak8rO1LzT5hCSFT0P3KfZEZMpbuXpzVGiZoM\ng5cWHgYcNnzhUatKodvzZizAOVGRR7UFg42O4ylhxDVe\n-----END X509 CRL-----\n"),
			}

			loadBalancerListenerPrototypeDefaultPoolModel := &vpcv1.LoadBalancerListenerPrototypeDefaultPoolLoadBalancerPoolIdentityLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			loadBalancerListenerHTTPSRedirectPrototypeListenerModel := &vpcv1.LoadBalancerListenerHTTPSRedirectPrototypeListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			loadBalancerListenerPrototypeHTTPSRedirectModel := &vpcv1.LoadBalancerListenerPrototypeHTTPSRedirect{
				HTTPStatusCode: core.Int64Ptr(int64(301)),
				Listener:       loadBalancerListenerHTTPSRedirectPrototypeListenerModel,
				URI:            core.StringPtr("/example?doc=get"),
			}

			loadBalancerListenerPolicyRulePrototypeModel := &vpcv1.LoadBalancerListenerPolicyRulePrototype{
				Condition: core.StringPtr("contains"),
				Field:     core.StringPtr("MY-APP-HEADER"),
				Type:      core.StringPtr("body"),
				Value:     core.StringPtr("my-value"),
			}

			loadBalancerListenerPolicyTargetPrototypeModel := &vpcv1.LoadBalancerListenerPolicyTargetPrototypeByListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			loadBalancerListenerPolicyPrototypeModel := &vpcv1.LoadBalancerListenerPolicyPrototype{
				Action:   core.StringPtr("forward_to_listener"),
				Name:     core.StringPtr("my-load-balancer-listener-policy"),
				Priority: core.Int64Ptr(int64(5)),
				Rules:    []vpcv1.LoadBalancerListenerPolicyRulePrototype{*loadBalancerListenerPolicyRulePrototypeModel},
				Target:   loadBalancerListenerPolicyTargetPrototypeModel,
			}

			createLoadBalancerListenerOptions := &vpcv1.CreateLoadBalancerListenerOptions{
				LoadBalancerID:        core.StringPtr("testString"),
				Protocol:              core.StringPtr("http"),
				AcceptProxyProtocol:   core.BoolPtr(true),
				CertificateInstance:   loadBalancerListenerPrototypeCertificateInstanceModel,
				ClientAuthentication:  loadBalancerListenerClientAuthenticationPrototypeModel,
				ConnectionLimit:       core.Int64Ptr(int64(2000)),
				DefaultPool:           loadBalancerListenerPrototypeDefaultPoolModel,
				HTTPSRedirect:         loadBalancerListenerPrototypeHTTPSRedirectModel,
				IdleConnectionTimeout: core.Int64Ptr(int64(100)),
				Policies:              []vpcv1.LoadBalancerListenerPolicyPrototype{*loadBalancerListenerPolicyPrototypeModel},
				Port:                  core.Int64Ptr(int64(443)),
				PortMax:               core.Int64Ptr(int64(499)),
				PortMin:               core.Int64Ptr(int64(443)),
			}

			loadBalancerListener, response, err := vpcService.CreateLoadBalancerListener(createLoadBalancerListenerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerListener - Retrieve a load balancer listener`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerListener(getLoadBalancerListenerOptions *GetLoadBalancerListenerOptions)`, func() {
			getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerListener, response, err := vpcService.GetLoadBalancerListener(getLoadBalancerListenerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancerListener - Update a load balancer listener`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancerListener(updateLoadBalancerListenerOptions *UpdateLoadBalancerListenerOptions)`, func() {
			loadBalancerListenerPatchCertificateInstanceModel := &vpcv1.LoadBalancerListenerPatchCertificateInstanceCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerListenerClientAuthenticationCertificateAuthorityPatchModel := &vpcv1.LoadBalancerListenerClientAuthenticationCertificateAuthorityPatchCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerListenerClientAuthenticationPatchModel := &vpcv1.LoadBalancerListenerClientAuthenticationPatch{
				CertificateAuthority:      loadBalancerListenerClientAuthenticationCertificateAuthorityPatchModel,
				CertificateRevocationList: core.StringPtr("-----BEGIN X509 CRL-----\nMIICvTCBpgIBATANBgkqhkiG9w0BAQsFADBMMQswCQYDVQQGEwJVUzEOMAwGA1UE\nCAwFZGVsYXMxDDAKBgNVBAoMA0lCTTENMAsGA1UECwwEcm9vdDEQMA4GA1UEAwwH\ncm9vdC1jYRcNMjUwOTA4MDUwMjQwWhcNMjUxMDA4MDUwMjQwWjAVMBMCAhAAFw0y\nNTA5MDgwNTAxNTlaoA8wDTALBgNVHRQEBAICEAAwDQYJKoZIhvcNAQELBQADggIB\nACeEcj7ompUepc5qTvTrNA5PoK5bN71gNI7Rbhq/Bxf1YPMp2iU3qMSj7YpVP7aw\nGNrxFoIZcQ4X7PYyHMfDk6Z83PSTVMnSOVk09fZW49tyVTWmzBVLz3R1bPasnWTZ\n0hRIv9j9n7Lemin+0ubIR/2zmsfBs1JFAFEbbRcgwg+qotsfZNLkX6bjHDpsRQzE\nmXUEu4/AqAsWPbFzG2uMKZ9pKOK+Nn3bt/NEK+AFlnSmgjEqzQ+0zhsrCExIReJV\nc2oiLBkLG6rBwxlGDog+PqwjP+1wGNIL1J3c2lMW1IGMNcts/aDBO5LtPVIY1LsQ\nFoeaTfm3U3GKC/pTczoDk/pKN756f8O05nTWUHgktcNsPvgqDKnpvEkI3VPf9Y4a\nfMOzKgVTgY1dSgjzHO8+4ZfcVGpBePsjOe0/RCUwkgtgOyGtcmBPTMJa0elJzjaM\njD9myqIXkB359sqbuEmcrjgo5uUUvubFYpmT/W0YxOi/py/bDK+7uUs38nUElNkZ\n+YFRpNWjLF9JtAghX5MhA5BwhTTuATvWYuDdK769ifi9qcYvE4u+VNxYfOpPY6sv\nx4FnkZ9+A7s2hk11d+DEq29Efa0xak8rO1LzT5hCSFT0P3KfZEZMpbuXpzVGiZoM\ng5cWHgYcNnzhUatKodvzZizAOVGRR7UFg42O4ylhxDVe\n-----END X509 CRL-----\n"),
			}

			loadBalancerListenerDefaultPoolPatchModel := &vpcv1.LoadBalancerListenerDefaultPoolPatchLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			loadBalancerListenerHTTPSRedirectPatchListenerModel := &vpcv1.LoadBalancerListenerHTTPSRedirectPatchListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			loadBalancerListenerPatchHTTPSRedirectModel := &vpcv1.LoadBalancerListenerPatchHTTPSRedirect{
				HTTPStatusCode: core.Int64Ptr(int64(301)),
				Listener:       loadBalancerListenerHTTPSRedirectPatchListenerModel,
				URI:            core.StringPtr("/example?doc=get"),
			}

			loadBalancerListenerPatchModel := &vpcv1.LoadBalancerListenerPatch{
				AcceptProxyProtocol:   core.BoolPtr(true),
				CertificateInstance:   loadBalancerListenerPatchCertificateInstanceModel,
				ClientAuthentication:  loadBalancerListenerClientAuthenticationPatchModel,
				ConnectionLimit:       core.Int64Ptr(int64(2000)),
				DefaultPool:           loadBalancerListenerDefaultPoolPatchModel,
				HTTPSRedirect:         loadBalancerListenerPatchHTTPSRedirectModel,
				IdleConnectionTimeout: core.Int64Ptr(int64(100)),
				Port:                  core.Int64Ptr(int64(443)),
				PortMax:               core.Int64Ptr(int64(499)),
				PortMin:               core.Int64Ptr(int64(443)),
				Protocol:              core.StringPtr("http"),
			}
			loadBalancerListenerPatchModelAsPatch, asPatchErr := loadBalancerListenerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerOptions := &vpcv1.UpdateLoadBalancerListenerOptions{
				ID:                        core.StringPtr("testString"),
				LoadBalancerID:            core.StringPtr("testString"),
				LoadBalancerListenerPatch: loadBalancerListenerPatchModelAsPatch,
			}

			loadBalancerListener, response, err := vpcService.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListener).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerListenerPolicies - List policies for a load balancer listener`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerListenerPolicies(listLoadBalancerListenerPoliciesOptions *ListLoadBalancerListenerPoliciesOptions)`, func() {
			listLoadBalancerListenerPoliciesOptions := &vpcv1.ListLoadBalancerListenerPoliciesOptions{
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerListenerPolicyCollection, response, err := vpcService.ListLoadBalancerListenerPolicies(listLoadBalancerListenerPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateLoadBalancerListenerPolicy - Create a policy for a load balancer listener`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancerListenerPolicy(createLoadBalancerListenerPolicyOptions *CreateLoadBalancerListenerPolicyOptions)`, func() {
			loadBalancerListenerPolicyRulePrototypeModel := &vpcv1.LoadBalancerListenerPolicyRulePrototype{
				Condition: core.StringPtr("contains"),
				Field:     core.StringPtr("MY-APP-HEADER"),
				Type:      core.StringPtr("body"),
				Value:     core.StringPtr("my-value"),
			}

			loadBalancerListenerPolicyTargetPrototypeModel := &vpcv1.LoadBalancerListenerPolicyTargetPrototypeByListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			createLoadBalancerListenerPolicyOptions := &vpcv1.CreateLoadBalancerListenerPolicyOptions{
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				Action:         core.StringPtr("forward_to_listener"),
				Priority:       core.Int64Ptr(int64(5)),
				Name:           core.StringPtr("my-load-balancer-listener-policy"),
				Rules:          []vpcv1.LoadBalancerListenerPolicyRulePrototype{*loadBalancerListenerPolicyRulePrototypeModel},
				Target:         loadBalancerListenerPolicyTargetPrototypeModel,
			}

			loadBalancerListenerPolicy, response, err := vpcService.CreateLoadBalancerListenerPolicy(createLoadBalancerListenerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerListenerPolicy - Retrieve a load balancer listener policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerListenerPolicy(getLoadBalancerListenerPolicyOptions *GetLoadBalancerListenerPolicyOptions)`, func() {
			getLoadBalancerListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
				ID:             core.StringPtr("testString"),
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerListenerPolicy, response, err := vpcService.GetLoadBalancerListenerPolicy(getLoadBalancerListenerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancerListenerPolicy - Update a load balancer listener policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancerListenerPolicy(updateLoadBalancerListenerPolicyOptions *UpdateLoadBalancerListenerPolicyOptions)`, func() {
			loadBalancerListenerPolicyTargetPatchModel := &vpcv1.LoadBalancerListenerPolicyTargetPatchByListenerLoadBalancerListenerIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/listeners/r006-97c05c26-4748-4db8-ac3f-d8a9c79a7091"),
			}

			loadBalancerListenerPolicyPatchModel := &vpcv1.LoadBalancerListenerPolicyPatch{
				Name:     core.StringPtr("my-load-balancer-listener-policy-updated"),
				Priority: core.Int64Ptr(int64(5)),
				Target:   loadBalancerListenerPolicyTargetPatchModel,
			}
			loadBalancerListenerPolicyPatchModelAsPatch, asPatchErr := loadBalancerListenerPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerPolicyOptions := &vpcv1.UpdateLoadBalancerListenerPolicyOptions{
				ID:                              core.StringPtr("testString"),
				ListenerID:                      core.StringPtr("testString"),
				LoadBalancerID:                  core.StringPtr("testString"),
				LoadBalancerListenerPolicyPatch: loadBalancerListenerPolicyPatchModelAsPatch,
			}

			loadBalancerListenerPolicy, response, err := vpcService.UpdateLoadBalancerListenerPolicy(updateLoadBalancerListenerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerListenerPolicyRules - List rules of a load balancer listener policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerListenerPolicyRules(listLoadBalancerListenerPolicyRulesOptions *ListLoadBalancerListenerPolicyRulesOptions)`, func() {
			listLoadBalancerListenerPolicyRulesOptions := &vpcv1.ListLoadBalancerListenerPolicyRulesOptions{
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PolicyID:       core.StringPtr("testString"),
			}

			loadBalancerListenerPolicyRuleCollection, response, err := vpcService.ListLoadBalancerListenerPolicyRules(listLoadBalancerListenerPolicyRulesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRuleCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateLoadBalancerListenerPolicyRule - Create a rule for a load balancer listener policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancerListenerPolicyRule(createLoadBalancerListenerPolicyRuleOptions *CreateLoadBalancerListenerPolicyRuleOptions)`, func() {
			createLoadBalancerListenerPolicyRuleOptions := &vpcv1.CreateLoadBalancerListenerPolicyRuleOptions{
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PolicyID:       core.StringPtr("testString"),
				Condition:      core.StringPtr("contains"),
				Type:           core.StringPtr("body"),
				Value:          core.StringPtr("my-value"),
				Field:          core.StringPtr("MY-APP-HEADER"),
			}

			loadBalancerListenerPolicyRule, response, err := vpcService.CreateLoadBalancerListenerPolicyRule(createLoadBalancerListenerPolicyRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerListenerPolicyRule - Retrieve a load balancer listener policy rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerListenerPolicyRule(getLoadBalancerListenerPolicyRuleOptions *GetLoadBalancerListenerPolicyRuleOptions)`, func() {
			getLoadBalancerListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
				ID:             core.StringPtr("testString"),
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PolicyID:       core.StringPtr("testString"),
			}

			loadBalancerListenerPolicyRule, response, err := vpcService.GetLoadBalancerListenerPolicyRule(getLoadBalancerListenerPolicyRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancerListenerPolicyRule - Update a load balancer listener policy rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancerListenerPolicyRule(updateLoadBalancerListenerPolicyRuleOptions *UpdateLoadBalancerListenerPolicyRuleOptions)`, func() {
			loadBalancerListenerPolicyRulePatchModel := &vpcv1.LoadBalancerListenerPolicyRulePatch{
				Condition: core.StringPtr("contains"),
				Field:     core.StringPtr("MY-APP-HEADER"),
				Type:      core.StringPtr("body"),
				Value:     core.StringPtr("my-example-hostname-updated"),
			}
			loadBalancerListenerPolicyRulePatchModelAsPatch, asPatchErr := loadBalancerListenerPolicyRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerListenerPolicyRuleOptions := &vpcv1.UpdateLoadBalancerListenerPolicyRuleOptions{
				ID:                                  core.StringPtr("testString"),
				ListenerID:                          core.StringPtr("testString"),
				LoadBalancerID:                      core.StringPtr("testString"),
				PolicyID:                            core.StringPtr("testString"),
				LoadBalancerListenerPolicyRulePatch: loadBalancerListenerPolicyRulePatchModelAsPatch,
			}

			loadBalancerListenerPolicyRule, response, err := vpcService.UpdateLoadBalancerListenerPolicyRule(updateLoadBalancerListenerPolicyRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerListenerPolicyRule).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerPools - List pools of a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerPools(listLoadBalancerPoolsOptions *ListLoadBalancerPoolsOptions)`, func() {
			listLoadBalancerPoolsOptions := &vpcv1.ListLoadBalancerPoolsOptions{
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerPoolCollection, response, err := vpcService.ListLoadBalancerPools(listLoadBalancerPoolsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateLoadBalancerPool - Create a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancerPool(createLoadBalancerPoolOptions *CreateLoadBalancerPoolOptions)`, func() {
			loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{
				Field: core.StringPtr("Content-Type"),
				Value: core.StringPtr("text/plain"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeLoadBalancerPoolHealthMonitorTypeHttphttpsRequestGetPrototype{
				HeadersVar: []vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{*loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel},
				Method:     core.StringPtr("get"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototype{
				BodyRegex: core.StringPtr("ACTIVE"),
				Codes:     []string{"200"},
			}

			loadBalancerPoolPrototypeHealthMonitorModel := &vpcv1.LoadBalancerPoolPrototypeHealthMonitorLoadBalancerPoolHealthMonitorPrototypeLoadBalancerPoolHealthMonitorTypeHttphttpsPrototype{
				Delay:      core.Int64Ptr(int64(5)),
				MaxRetries: core.Int64Ptr(int64(2)),
				Port:       core.Int64Ptr(int64(22)),
				Timeout:    core.Int64Ptr(int64(2)),
				Request:    loadBalancerPoolHealthMonitorTypeHttphttpsRequestPrototypeModel,
				Response:   loadBalancerPoolHealthMonitorTypeHttphttpsResponsePrototypeModel,
				Type:       core.StringPtr("http"),
				URLPath:    core.StringPtr("/"),
			}

			certificateInstanceIdentityModel := &vpcv1.CertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerPoolClientAuthenticationPrototypeModel := &vpcv1.LoadBalancerPoolClientAuthenticationPrototype{
				CertificateInstance: certificateInstanceIdentityModel,
			}

			loadBalancerPoolFailsafePolicyPrototypeTargetModel := &vpcv1.LoadBalancerPoolFailsafePolicyPrototypeTargetLoadBalancerPoolIdentityLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			loadBalancerPoolFailsafePolicyPrototypeModel := &vpcv1.LoadBalancerPoolFailsafePolicyPrototype{
				Action: core.StringPtr("forward"),
				Target: loadBalancerPoolFailsafePolicyPrototypeTargetModel,
			}

			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			loadBalancerPoolMemberPrototypeModel := &vpcv1.LoadBalancerPoolMemberPrototype{
				Port:   core.Int64Ptr(int64(80)),
				Target: loadBalancerPoolMemberTargetPrototypeModel,
				Weight: core.Int64Ptr(int64(50)),
			}

			loadBalancerPoolServerAuthenticationPrototypeModel := &vpcv1.LoadBalancerPoolServerAuthenticationPrototype{
				CertificateAuthority: certificateInstanceIdentityModel,
				VerifyCertificate:    core.BoolPtr(false),
			}

			loadBalancerPoolPrototypeSessionPersistenceModel := &vpcv1.LoadBalancerPoolPrototypeSessionPersistence{
				CookieName: core.StringPtr("my-cookie-name"),
				Type:       core.StringPtr("app_cookie"),
			}

			createLoadBalancerPoolOptions := &vpcv1.CreateLoadBalancerPoolOptions{
				LoadBalancerID:       core.StringPtr("testString"),
				Algorithm:            core.StringPtr("least_connections"),
				HealthMonitor:        loadBalancerPoolPrototypeHealthMonitorModel,
				Protocol:             core.StringPtr("http"),
				ClientAuthentication: loadBalancerPoolClientAuthenticationPrototypeModel,
				FailsafePolicy:       loadBalancerPoolFailsafePolicyPrototypeModel,
				Members:              []vpcv1.LoadBalancerPoolMemberPrototype{*loadBalancerPoolMemberPrototypeModel},
				Name:                 core.StringPtr("my-load-balancer-pool"),
				ProxyProtocol:        core.StringPtr("disabled"),
				ServerAuthentication: loadBalancerPoolServerAuthenticationPrototypeModel,
				SessionPersistence:   loadBalancerPoolPrototypeSessionPersistenceModel,
			}

			loadBalancerPool, response, err := vpcService.CreateLoadBalancerPool(createLoadBalancerPoolOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerPool - Retrieve a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerPool(getLoadBalancerPoolOptions *GetLoadBalancerPoolOptions)`, func() {
			getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			loadBalancerPool, response, err := vpcService.GetLoadBalancerPool(getLoadBalancerPoolOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancerPool - Update a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancerPool(updateLoadBalancerPoolOptions *UpdateLoadBalancerPoolOptions)`, func() {
			certificateInstanceIdentityModel := &vpcv1.CertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			loadBalancerPoolClientAuthenticationPatchModel := &vpcv1.LoadBalancerPoolClientAuthenticationPatch{
				CertificateInstance: certificateInstanceIdentityModel,
			}

			loadBalancerPoolFailsafePolicyTargetPatchModel := &vpcv1.LoadBalancerPoolFailsafePolicyTargetPatchLoadBalancerPoolIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/load_balancers/r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727/pools/r006-70294e14-4e61-11e8-bcf4-0242ac110004"),
			}

			loadBalancerPoolFailsafePolicyPatchModel := &vpcv1.LoadBalancerPoolFailsafePolicyPatch{
				Action: core.StringPtr("forward"),
				Target: loadBalancerPoolFailsafePolicyTargetPatchModel,
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{
				Field: core.StringPtr("Content-Type"),
				Value: core.StringPtr("text/plain"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsRequestPatchModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestPatch{
				Body:       core.StringPtr("ACTIVE"),
				HeadersVar: []vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototype{*loadBalancerPoolHealthMonitorTypeHttphttpsRequestHeaderPrototypeModel},
				Method:     core.StringPtr("get"),
			}

			loadBalancerPoolHealthMonitorTypeHttphttpsResponsePatchModel := &vpcv1.LoadBalancerPoolHealthMonitorTypeHttphttpsResponsePatch{
				BodyRegex: core.StringPtr("ACTIVE"),
				Codes:     []string{"2XX"},
			}

			loadBalancerPoolPatchHealthMonitorModel := &vpcv1.LoadBalancerPoolPatchHealthMonitor{
				Delay:      core.Int64Ptr(int64(5)),
				MaxRetries: core.Int64Ptr(int64(2)),
				Port:       core.Int64Ptr(int64(22)),
				Request:    loadBalancerPoolHealthMonitorTypeHttphttpsRequestPatchModel,
				Response:   loadBalancerPoolHealthMonitorTypeHttphttpsResponsePatchModel,
				Timeout:    core.Int64Ptr(int64(2)),
				Type:       core.StringPtr("http"),
				URLPath:    core.StringPtr("/"),
			}

			loadBalancerPoolServerAuthenticationPatchModel := &vpcv1.LoadBalancerPoolServerAuthenticationPatch{
				CertificateAuthority: certificateInstanceIdentityModel,
				VerifyCertificate:    core.BoolPtr(false),
			}

			loadBalancerPoolPatchSessionPersistenceModel := &vpcv1.LoadBalancerPoolPatchSessionPersistence{
				CookieName: core.StringPtr("my-cookie-name"),
				Type:       core.StringPtr("app_cookie"),
			}

			loadBalancerPoolPatchModel := &vpcv1.LoadBalancerPoolPatch{
				Algorithm:            core.StringPtr("weighted_round_robin"),
				ClientAuthentication: loadBalancerPoolClientAuthenticationPatchModel,
				FailsafePolicy:       loadBalancerPoolFailsafePolicyPatchModel,
				HealthMonitor:        loadBalancerPoolPatchHealthMonitorModel,
				Name:                 core.StringPtr("my-load-balancer-pool"),
				Protocol:             core.StringPtr("http"),
				ProxyProtocol:        core.StringPtr("disabled"),
				ServerAuthentication: loadBalancerPoolServerAuthenticationPatchModel,
				SessionPersistence:   loadBalancerPoolPatchSessionPersistenceModel,
			}
			loadBalancerPoolPatchModelAsPatch, asPatchErr := loadBalancerPoolPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerPoolOptions := &vpcv1.UpdateLoadBalancerPoolOptions{
				ID:                    core.StringPtr("testString"),
				LoadBalancerID:        core.StringPtr("testString"),
				LoadBalancerPoolPatch: loadBalancerPoolPatchModelAsPatch,
				IfMatch:               core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			loadBalancerPool, response, err := vpcService.UpdateLoadBalancerPool(updateLoadBalancerPoolOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPool).ToNot(BeNil())
		})
	})

	Describe(`ListLoadBalancerPoolMembers - List members of a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListLoadBalancerPoolMembers(listLoadBalancerPoolMembersOptions *ListLoadBalancerPoolMembersOptions)`, func() {
			listLoadBalancerPoolMembersOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{
				LoadBalancerID: core.StringPtr("testString"),
				PoolID:         core.StringPtr("testString"),
			}

			loadBalancerPoolMemberCollection, response, err := vpcService.ListLoadBalancerPoolMembers(listLoadBalancerPoolMembersOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMemberCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateLoadBalancerPoolMember - Create a member in a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateLoadBalancerPoolMember(createLoadBalancerPoolMemberOptions *CreateLoadBalancerPoolMemberOptions)`, func() {
			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			createLoadBalancerPoolMemberOptions := &vpcv1.CreateLoadBalancerPoolMemberOptions{
				LoadBalancerID: core.StringPtr("testString"),
				PoolID:         core.StringPtr("testString"),
				Port:           core.Int64Ptr(int64(80)),
				Target:         loadBalancerPoolMemberTargetPrototypeModel,
				Weight:         core.Int64Ptr(int64(50)),
			}

			loadBalancerPoolMember, response, err := vpcService.CreateLoadBalancerPoolMember(createLoadBalancerPoolMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
	})

	Describe(`ReplaceLoadBalancerPoolMembers - Replace load balancer pool members`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceLoadBalancerPoolMembers(replaceLoadBalancerPoolMembersOptions *ReplaceLoadBalancerPoolMembersOptions)`, func() {
			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			loadBalancerPoolMemberPrototypeModel := &vpcv1.LoadBalancerPoolMemberPrototype{
				Port:   core.Int64Ptr(int64(80)),
				Target: loadBalancerPoolMemberTargetPrototypeModel,
				Weight: core.Int64Ptr(int64(50)),
			}

			replaceLoadBalancerPoolMembersOptions := &vpcv1.ReplaceLoadBalancerPoolMembersOptions{
				LoadBalancerID: core.StringPtr("testString"),
				PoolID:         core.StringPtr("testString"),
				Members:        []vpcv1.LoadBalancerPoolMemberPrototype{*loadBalancerPoolMemberPrototypeModel},
			}

			loadBalancerPoolMemberCollection, response, err := vpcService.ReplaceLoadBalancerPoolMembers(replaceLoadBalancerPoolMembersOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMemberCollection).ToNot(BeNil())
		})
	})

	Describe(`GetLoadBalancerPoolMember - Retrieve a load balancer pool member`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetLoadBalancerPoolMember(getLoadBalancerPoolMemberOptions *GetLoadBalancerPoolMemberOptions)`, func() {
			getLoadBalancerPoolMemberOptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PoolID:         core.StringPtr("testString"),
			}

			loadBalancerPoolMember, response, err := vpcService.GetLoadBalancerPoolMember(getLoadBalancerPoolMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
	})

	Describe(`UpdateLoadBalancerPoolMember - Update a load balancer pool member`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateLoadBalancerPoolMember(updateLoadBalancerPoolMemberOptions *UpdateLoadBalancerPoolMemberOptions)`, func() {
			loadBalancerPoolMemberTargetPrototypeModel := &vpcv1.LoadBalancerPoolMemberTargetPrototypeByInstanceInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"),
			}

			loadBalancerPoolMemberPatchModel := &vpcv1.LoadBalancerPoolMemberPatch{
				Port:   core.Int64Ptr(int64(80)),
				Target: loadBalancerPoolMemberTargetPrototypeModel,
				Weight: core.Int64Ptr(int64(99)),
			}
			loadBalancerPoolMemberPatchModelAsPatch, asPatchErr := loadBalancerPoolMemberPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateLoadBalancerPoolMemberOptions := &vpcv1.UpdateLoadBalancerPoolMemberOptions{
				ID:                          core.StringPtr("testString"),
				LoadBalancerID:              core.StringPtr("testString"),
				PoolID:                      core.StringPtr("testString"),
				LoadBalancerPoolMemberPatch: loadBalancerPoolMemberPatchModelAsPatch,
			}

			loadBalancerPoolMember, response, err := vpcService.UpdateLoadBalancerPoolMember(updateLoadBalancerPoolMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(loadBalancerPoolMember).ToNot(BeNil())
		})
	})

	Describe(`ListNetworkAcls - List network ACLs`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListNetworkAcls(listNetworkAclsOptions *ListNetworkAclsOptions) with pagination`, func() {
			listNetworkAclsOptions := &vpcv1.ListNetworkAclsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listNetworkAclsOptions.Start = nil
			listNetworkAclsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.NetworkACL
			for {
				networkACLCollection, response, err := vpcService.ListNetworkAcls(listNetworkAclsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(networkACLCollection).ToNot(BeNil())
				allResults = append(allResults, networkACLCollection.NetworkAcls...)

				nextToken, nextErr := networkACLCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listNetworkAclsOptions.Start = nextToken

				if listNetworkAclsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListNetworkAcls(listNetworkAclsOptions *ListNetworkAclsOptions) using NetworkAclsPager`, func() {
			listNetworkAclsOptions := &vpcv1.ListNetworkAclsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewNetworkAclsPager(listNetworkAclsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.NetworkACL
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewNetworkAclsPager(listNetworkAclsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListNetworkAcls() returned a total of %d item(s) using NetworkAclsPager.\n", len(allResults))
		})
	})

	Describe(`CreateNetworkACL - Create a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateNetworkACL(createNetworkACLOptions *CreateNetworkACLOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			networkACLPrototypeVPCModel := &vpcv1.NetworkACLPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			networkACLRulePrototypeNetworkACLContextModel := &vpcv1.NetworkACLRulePrototypeNetworkACLContextNetworkACLRuleProtocolAnyPrototype{
				Action:      core.StringPtr("allow"),
				Destination: core.StringPtr("192.168.3.2/32"),
				Direction:   core.StringPtr("inbound"),
				IPVersion:   core.StringPtr("ipv4"),
				Name:        core.StringPtr("my-network-acl-rule"),
				Source:      core.StringPtr("192.168.3.2/32"),
				Protocol:    core.StringPtr("any"),
			}

			networkACLPrototypeModel := &vpcv1.NetworkACLPrototypeNetworkACLByRules{
				Name:          core.StringPtr("my-network-acl"),
				ResourceGroup: resourceGroupIdentityModel,
				VPC:           networkACLPrototypeVPCModel,
				Rules:         []vpcv1.NetworkACLRulePrototypeNetworkACLContextIntf{networkACLRulePrototypeNetworkACLContextModel},
			}

			createNetworkACLOptions := &vpcv1.CreateNetworkACLOptions{
				NetworkACLPrototype: networkACLPrototypeModel,
			}

			networkACL, response, err := vpcService.CreateNetworkACL(createNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACL).ToNot(BeNil())
		})
	})

	Describe(`GetNetworkACL - Retrieve a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetNetworkACL(getNetworkACLOptions *GetNetworkACLOptions)`, func() {
			getNetworkACLOptions := &vpcv1.GetNetworkACLOptions{
				ID: core.StringPtr("testString"),
			}

			networkACL, response, err := vpcService.GetNetworkACL(getNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
	})

	Describe(`UpdateNetworkACL - Update a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateNetworkACL(updateNetworkACLOptions *UpdateNetworkACLOptions)`, func() {
			networkACLPatchModel := &vpcv1.NetworkACLPatch{
				Name: core.StringPtr("my-network-acl-updated"),
			}
			networkACLPatchModelAsPatch, asPatchErr := networkACLPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateNetworkACLOptions := &vpcv1.UpdateNetworkACLOptions{
				ID:              core.StringPtr("testString"),
				NetworkACLPatch: networkACLPatchModelAsPatch,
			}

			networkACL, response, err := vpcService.UpdateNetworkACL(updateNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
	})

	Describe(`ListNetworkACLRules - List rules for a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListNetworkACLRules(listNetworkACLRulesOptions *ListNetworkACLRulesOptions) with pagination`, func() {
			listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
				NetworkACLID: core.StringPtr("testString"),
				Direction:    core.StringPtr("inbound"),
				Limit:        core.Int64Ptr(int64(10)),
				Start:        core.StringPtr("testString"),
			}

			listNetworkACLRulesOptions.Start = nil
			listNetworkACLRulesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.NetworkACLRuleItemIntf
			for {
				networkACLRuleCollection, response, err := vpcService.ListNetworkACLRules(listNetworkACLRulesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(networkACLRuleCollection).ToNot(BeNil())
				allResults = append(allResults, networkACLRuleCollection.Rules...)

				nextToken, nextErr := networkACLRuleCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listNetworkACLRulesOptions.Start = nextToken

				if listNetworkACLRulesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListNetworkACLRules(listNetworkACLRulesOptions *ListNetworkACLRulesOptions) using NetworkACLRulesPager`, func() {
			listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
				NetworkACLID: core.StringPtr("testString"),
				Direction:    core.StringPtr("inbound"),
				Limit:        core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewNetworkACLRulesPager(listNetworkACLRulesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.NetworkACLRuleItemIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewNetworkACLRulesPager(listNetworkACLRulesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListNetworkACLRules() returned a total of %d item(s) using NetworkACLRulesPager.\n", len(allResults))
		})
	})

	Describe(`CreateNetworkACLRule - Create a rule for a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateNetworkACLRule(createNetworkACLRuleOptions *CreateNetworkACLRuleOptions)`, func() {
			networkACLRuleBeforePrototypeModel := &vpcv1.NetworkACLRuleBeforePrototypeNetworkACLRuleIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/network_acls/r006-3217cb8b-5368-452a-9399-a84f14fb539d/rules/r006-8daca77a-4980-4d33-8f3e-7038797be8f9"),
			}

			networkACLRulePrototypeModel := &vpcv1.NetworkACLRulePrototypeNetworkACLRuleProtocolAnyPrototype{
				Action:      core.StringPtr("allow"),
				Before:      networkACLRuleBeforePrototypeModel,
				Destination: core.StringPtr("192.168.3.2/32"),
				Direction:   core.StringPtr("inbound"),
				IPVersion:   core.StringPtr("ipv4"),
				Name:        core.StringPtr("my-network-acl-rule"),
				Source:      core.StringPtr("192.168.3.2/32"),
				Protocol:    core.StringPtr("any"),
			}

			createNetworkACLRuleOptions := &vpcv1.CreateNetworkACLRuleOptions{
				NetworkACLID:            core.StringPtr("testString"),
				NetworkACLRulePrototype: networkACLRulePrototypeModel,
			}

			networkACLRule, response, err := vpcService.CreateNetworkACLRule(createNetworkACLRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACLRule).ToNot(BeNil())
		})
	})

	Describe(`GetNetworkACLRule - Retrieve a network ACL rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetNetworkACLRule(getNetworkACLRuleOptions *GetNetworkACLRuleOptions)`, func() {
			getNetworkACLRuleOptions := &vpcv1.GetNetworkACLRuleOptions{
				ID:           core.StringPtr("testString"),
				NetworkACLID: core.StringPtr("testString"),
			}

			networkACLRule, response, err := vpcService.GetNetworkACLRule(getNetworkACLRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACLRule).ToNot(BeNil())
		})
	})

	Describe(`UpdateNetworkACLRule - Update a network ACL rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateNetworkACLRule(updateNetworkACLRuleOptions *UpdateNetworkACLRuleOptions)`, func() {
			networkACLRuleBeforePatchModel := &vpcv1.NetworkACLRuleBeforePatchNetworkACLRuleIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/network_acls/r006-3217cb8b-5368-452a-9399-a84f14fb539d/rules/r006-8daca77a-4980-4d33-8f3e-7038797be8f9"),
			}

			networkACLRulePatchModel := &vpcv1.NetworkACLRulePatch{
				Action:             core.StringPtr("allow"),
				Before:             networkACLRuleBeforePatchModel,
				Code:               core.Int64Ptr(int64(0)),
				Destination:        core.StringPtr("192.168.3.2/32"),
				DestinationPortMax: core.Int64Ptr(int64(22)),
				DestinationPortMin: core.Int64Ptr(int64(22)),
				Direction:          core.StringPtr("inbound"),
				Name:               core.StringPtr("my-network-acl-rule"),
				Source:             core.StringPtr("10.0.0.0/0"),
				SourcePortMax:      core.Int64Ptr(int64(65535)),
				SourcePortMin:      core.Int64Ptr(int64(49152)),
				Type:               core.Int64Ptr(int64(8)),
			}
			networkACLRulePatchModelAsPatch, asPatchErr := networkACLRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateNetworkACLRuleOptions := &vpcv1.UpdateNetworkACLRuleOptions{
				ID:                  core.StringPtr("testString"),
				NetworkACLID:        core.StringPtr("testString"),
				NetworkACLRulePatch: networkACLRulePatchModelAsPatch,
			}

			networkACLRule, response, err := vpcService.UpdateNetworkACLRule(updateNetworkACLRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACLRule).ToNot(BeNil())
		})
	})

	Describe(`ListPlacementGroups - List placement groups`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPlacementGroups(listPlacementGroupsOptions *ListPlacementGroupsOptions) with pagination`, func() {
			listPlacementGroupsOptions := &vpcv1.ListPlacementGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listPlacementGroupsOptions.Start = nil
			listPlacementGroupsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PlacementGroup
			for {
				placementGroupCollection, response, err := vpcService.ListPlacementGroups(listPlacementGroupsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(placementGroupCollection).ToNot(BeNil())
				allResults = append(allResults, placementGroupCollection.PlacementGroups...)

				nextToken, nextErr := placementGroupCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPlacementGroupsOptions.Start = nextToken

				if listPlacementGroupsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPlacementGroups(listPlacementGroupsOptions *ListPlacementGroupsOptions) using PlacementGroupsPager`, func() {
			listPlacementGroupsOptions := &vpcv1.ListPlacementGroupsOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewPlacementGroupsPager(listPlacementGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PlacementGroup
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPlacementGroupsPager(listPlacementGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPlacementGroups() returned a total of %d item(s) using PlacementGroupsPager.\n", len(allResults))
		})
	})

	Describe(`CreatePlacementGroup - Create a placement group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreatePlacementGroup(createPlacementGroupOptions *CreatePlacementGroupOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createPlacementGroupOptions := &vpcv1.CreatePlacementGroupOptions{
				Strategy:      core.StringPtr("host_spread"),
				Name:          core.StringPtr("my-placement-group"),
				ResourceGroup: resourceGroupIdentityModel,
			}

			placementGroup, response, err := vpcService.CreatePlacementGroup(createPlacementGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(placementGroup).ToNot(BeNil())
		})
	})

	Describe(`GetPlacementGroup - Retrieve a placement group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPlacementGroup(getPlacementGroupOptions *GetPlacementGroupOptions)`, func() {
			getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{
				ID: core.StringPtr("testString"),
			}

			placementGroup, response, err := vpcService.GetPlacementGroup(getPlacementGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(placementGroup).ToNot(BeNil())
		})
	})

	Describe(`UpdatePlacementGroup - Update a placement group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdatePlacementGroup(updatePlacementGroupOptions *UpdatePlacementGroupOptions)`, func() {
			placementGroupPatchModel := &vpcv1.PlacementGroupPatch{
				Name: core.StringPtr("my-placement-group-updated"),
			}
			placementGroupPatchModelAsPatch, asPatchErr := placementGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePlacementGroupOptions := &vpcv1.UpdatePlacementGroupOptions{
				ID:                  core.StringPtr("testString"),
				PlacementGroupPatch: placementGroupPatchModelAsPatch,
			}

			placementGroup, response, err := vpcService.UpdatePlacementGroup(updatePlacementGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(placementGroup).ToNot(BeNil())
		})
	})

	Describe(`ListPrivatePathServiceGateways - List private path service gateways`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPrivatePathServiceGateways(listPrivatePathServiceGatewaysOptions *ListPrivatePathServiceGatewaysOptions) with pagination`, func() {
			listPrivatePathServiceGatewaysOptions := &vpcv1.ListPrivatePathServiceGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listPrivatePathServiceGatewaysOptions.Start = nil
			listPrivatePathServiceGatewaysOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PrivatePathServiceGateway
			for {
				privatePathServiceGatewayCollection, response, err := vpcService.ListPrivatePathServiceGateways(listPrivatePathServiceGatewaysOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(privatePathServiceGatewayCollection).ToNot(BeNil())
				allResults = append(allResults, privatePathServiceGatewayCollection.PrivatePathServiceGateways...)

				nextToken, nextErr := privatePathServiceGatewayCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPrivatePathServiceGatewaysOptions.Start = nextToken

				if listPrivatePathServiceGatewaysOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPrivatePathServiceGateways(listPrivatePathServiceGatewaysOptions *ListPrivatePathServiceGatewaysOptions) using PrivatePathServiceGatewaysPager`, func() {
			listPrivatePathServiceGatewaysOptions := &vpcv1.ListPrivatePathServiceGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewPrivatePathServiceGatewaysPager(listPrivatePathServiceGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PrivatePathServiceGateway
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPrivatePathServiceGatewaysPager(listPrivatePathServiceGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPrivatePathServiceGateways() returned a total of %d item(s) using PrivatePathServiceGatewaysPager.\n", len(allResults))
		})
	})

	Describe(`CreatePrivatePathServiceGateway - Create a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreatePrivatePathServiceGateway(createPrivatePathServiceGatewayOptions *CreatePrivatePathServiceGatewayOptions)`, func() {
			privatePathServiceGatewayPrototypeLoadBalancerModel := &vpcv1.PrivatePathServiceGatewayPrototypeLoadBalancerLoadBalancerIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createPrivatePathServiceGatewayOptions := &vpcv1.CreatePrivatePathServiceGatewayOptions{
				LoadBalancer:        privatePathServiceGatewayPrototypeLoadBalancerModel,
				ServiceEndpoints:    []string{"*.example.com"},
				DefaultAccessPolicy: core.StringPtr("deny"),
				Name:                core.StringPtr("my-private-path-service-gateway"),
				ResourceGroup:       resourceGroupIdentityModel,
				ZonalAffinity:       core.BoolPtr(false),
			}

			privatePathServiceGateway, response, err := vpcService.CreatePrivatePathServiceGateway(createPrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
	})

	Describe(`GetPrivatePathServiceGateway - Retrieve a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPrivatePathServiceGateway(getPrivatePathServiceGatewayOptions *GetPrivatePathServiceGatewayOptions)`, func() {
			getPrivatePathServiceGatewayOptions := &vpcv1.GetPrivatePathServiceGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			privatePathServiceGateway, response, err := vpcService.GetPrivatePathServiceGateway(getPrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
	})

	Describe(`UpdatePrivatePathServiceGateway - Update a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdatePrivatePathServiceGateway(updatePrivatePathServiceGatewayOptions *UpdatePrivatePathServiceGatewayOptions)`, func() {
			privatePathServiceGatewayPatchLoadBalancerModel := &vpcv1.PrivatePathServiceGatewayPatchLoadBalancerLoadBalancerIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
			}

			privatePathServiceGatewayPatchModel := &vpcv1.PrivatePathServiceGatewayPatch{
				DefaultAccessPolicy: core.StringPtr("deny"),
				LoadBalancer:        privatePathServiceGatewayPatchLoadBalancerModel,
				Name:                core.StringPtr("my-private-path-service-gateway"),
				ZonalAffinity:       core.BoolPtr(true),
			}
			privatePathServiceGatewayPatchModelAsPatch, asPatchErr := privatePathServiceGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePrivatePathServiceGatewayOptions := &vpcv1.UpdatePrivatePathServiceGatewayOptions{
				ID:                             core.StringPtr("testString"),
				PrivatePathServiceGatewayPatch: privatePathServiceGatewayPatchModelAsPatch,
			}

			privatePathServiceGateway, response, err := vpcService.UpdatePrivatePathServiceGateway(updatePrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGateway).ToNot(BeNil())
		})
	})

	Describe(`ListPrivatePathServiceGatewayAccountPolicies - List account policies for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPrivatePathServiceGatewayAccountPolicies(listPrivatePathServiceGatewayAccountPoliciesOptions *ListPrivatePathServiceGatewayAccountPoliciesOptions) with pagination`, func() {
			listPrivatePathServiceGatewayAccountPoliciesOptions := &vpcv1.ListPrivatePathServiceGatewayAccountPoliciesOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
				Start:                       core.StringPtr("testString"),
			}

			listPrivatePathServiceGatewayAccountPoliciesOptions.Start = nil
			listPrivatePathServiceGatewayAccountPoliciesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PrivatePathServiceGatewayAccountPolicy
			for {
				privatePathServiceGatewayAccountPolicyCollection, response, err := vpcService.ListPrivatePathServiceGatewayAccountPolicies(listPrivatePathServiceGatewayAccountPoliciesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(privatePathServiceGatewayAccountPolicyCollection).ToNot(BeNil())
				allResults = append(allResults, privatePathServiceGatewayAccountPolicyCollection.AccountPolicies...)

				nextToken, nextErr := privatePathServiceGatewayAccountPolicyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPrivatePathServiceGatewayAccountPoliciesOptions.Start = nextToken

				if listPrivatePathServiceGatewayAccountPoliciesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPrivatePathServiceGatewayAccountPolicies(listPrivatePathServiceGatewayAccountPoliciesOptions *ListPrivatePathServiceGatewayAccountPoliciesOptions) using PrivatePathServiceGatewayAccountPoliciesPager`, func() {
			listPrivatePathServiceGatewayAccountPoliciesOptions := &vpcv1.ListPrivatePathServiceGatewayAccountPoliciesOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewPrivatePathServiceGatewayAccountPoliciesPager(listPrivatePathServiceGatewayAccountPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PrivatePathServiceGatewayAccountPolicy
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPrivatePathServiceGatewayAccountPoliciesPager(listPrivatePathServiceGatewayAccountPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPrivatePathServiceGatewayAccountPolicies() returned a total of %d item(s) using PrivatePathServiceGatewayAccountPoliciesPager.\n", len(allResults))
		})
	})

	Describe(`CreatePrivatePathServiceGatewayAccountPolicy - Create an account policy for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreatePrivatePathServiceGatewayAccountPolicy(createPrivatePathServiceGatewayAccountPolicyOptions *CreatePrivatePathServiceGatewayAccountPolicyOptions)`, func() {
			privatePathServiceGatewayAccountPolicyPrototypeAccountModel := &vpcv1.PrivatePathServiceGatewayAccountPolicyPrototypeAccountAccountIdentityByID{
				ID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
			}

			createPrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.CreatePrivatePathServiceGatewayAccountPolicyOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccessPolicy:                core.StringPtr("deny"),
				Account:                     privatePathServiceGatewayAccountPolicyPrototypeAccountModel,
			}

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.CreatePrivatePathServiceGatewayAccountPolicy(createPrivatePathServiceGatewayAccountPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
	})

	Describe(`GetPrivatePathServiceGatewayAccountPolicy - Retrieve an account policy for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPrivatePathServiceGatewayAccountPolicy(getPrivatePathServiceGatewayAccountPolicyOptions *GetPrivatePathServiceGatewayAccountPolicyOptions)`, func() {
			getPrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.GetPrivatePathServiceGatewayAccountPolicyOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
			}

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.GetPrivatePathServiceGatewayAccountPolicy(getPrivatePathServiceGatewayAccountPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdatePrivatePathServiceGatewayAccountPolicy - Update an account policy for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdatePrivatePathServiceGatewayAccountPolicy(updatePrivatePathServiceGatewayAccountPolicyOptions *UpdatePrivatePathServiceGatewayAccountPolicyOptions)`, func() {
			privatePathServiceGatewayAccountPolicyPatchModel := &vpcv1.PrivatePathServiceGatewayAccountPolicyPatch{
				AccessPolicy: core.StringPtr("review"),
			}
			privatePathServiceGatewayAccountPolicyPatchModelAsPatch, asPatchErr := privatePathServiceGatewayAccountPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.UpdatePrivatePathServiceGatewayAccountPolicyOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				PrivatePathServiceGatewayAccountPolicyPatch: privatePathServiceGatewayAccountPolicyPatchModelAsPatch,
			}

			privatePathServiceGatewayAccountPolicy, response, err := vpcService.UpdatePrivatePathServiceGatewayAccountPolicy(updatePrivatePathServiceGatewayAccountPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayAccountPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListPrivatePathServiceGatewayEndpointGatewayBindings - List endpoint gateway bindings for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPrivatePathServiceGatewayEndpointGatewayBindings(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions *ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions) with pagination`, func() {
			listPrivatePathServiceGatewayEndpointGatewayBindingsOptions := &vpcv1.ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
				Start:                       core.StringPtr("testString"),
				Status:                      core.StringPtr("abandoned"),
			}

			listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.Start = nil
			listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PrivatePathServiceGatewayEndpointGatewayBinding
			for {
				privatePathServiceGatewayEndpointGatewayBindingCollection, response, err := vpcService.ListPrivatePathServiceGatewayEndpointGatewayBindings(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(privatePathServiceGatewayEndpointGatewayBindingCollection).ToNot(BeNil())
				allResults = append(allResults, privatePathServiceGatewayEndpointGatewayBindingCollection.EndpointGatewayBindings...)

				nextToken, nextErr := privatePathServiceGatewayEndpointGatewayBindingCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.Start = nextToken

				if listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPrivatePathServiceGatewayEndpointGatewayBindings(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions *ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions) using PrivatePathServiceGatewayEndpointGatewayBindingsPager`, func() {
			listPrivatePathServiceGatewayEndpointGatewayBindingsOptions := &vpcv1.ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				AccountID:                   core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:                       core.Int64Ptr(int64(10)),
				Status:                      core.StringPtr("abandoned"),
			}

			// Test GetNext().
			pager, err := vpcService.NewPrivatePathServiceGatewayEndpointGatewayBindingsPager(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PrivatePathServiceGatewayEndpointGatewayBinding
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPrivatePathServiceGatewayEndpointGatewayBindingsPager(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPrivatePathServiceGatewayEndpointGatewayBindings() returned a total of %d item(s) using PrivatePathServiceGatewayEndpointGatewayBindingsPager.\n", len(allResults))
		})
	})

	Describe(`GetPrivatePathServiceGatewayEndpointGatewayBinding - Retrieve an endpoint gateway binding for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPrivatePathServiceGatewayEndpointGatewayBinding(getPrivatePathServiceGatewayEndpointGatewayBindingOptions *GetPrivatePathServiceGatewayEndpointGatewayBindingOptions)`, func() {
			getPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.GetPrivatePathServiceGatewayEndpointGatewayBindingOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
			}

			privatePathServiceGatewayEndpointGatewayBinding, response, err := vpcService.GetPrivatePathServiceGatewayEndpointGatewayBinding(getPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(privatePathServiceGatewayEndpointGatewayBinding).ToNot(BeNil())
		})
	})

	Describe(`DenyPrivatePathServiceGatewayEndpointGatewayBinding - Deny an endpoint gateway binding for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DenyPrivatePathServiceGatewayEndpointGatewayBinding(denyPrivatePathServiceGatewayEndpointGatewayBindingOptions *DenyPrivatePathServiceGatewayEndpointGatewayBindingOptions)`, func() {
			denyPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.DenyPrivatePathServiceGatewayEndpointGatewayBindingOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				SetAccountPolicy:            core.BoolPtr(true),
			}

			response, err := vpcService.DenyPrivatePathServiceGatewayEndpointGatewayBinding(denyPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`PermitPrivatePathServiceGatewayEndpointGatewayBinding - Permit an endpoint gateway binding for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`PermitPrivatePathServiceGatewayEndpointGatewayBinding(permitPrivatePathServiceGatewayEndpointGatewayBindingOptions *PermitPrivatePathServiceGatewayEndpointGatewayBindingOptions)`, func() {
			permitPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.PermitPrivatePathServiceGatewayEndpointGatewayBindingOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				SetAccountPolicy:            core.BoolPtr(true),
			}

			response, err := vpcService.PermitPrivatePathServiceGatewayEndpointGatewayBinding(permitPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`PublishPrivatePathServiceGateway - Publish a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`PublishPrivatePathServiceGateway(publishPrivatePathServiceGatewayOptions *PublishPrivatePathServiceGatewayOptions)`, func() {
			publishPrivatePathServiceGatewayOptions := &vpcv1.PublishPrivatePathServiceGatewayOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.PublishPrivatePathServiceGateway(publishPrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RevokeAccountForPrivatePathServiceGateway - Revoke access to a private path service gateway for an account`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RevokeAccountForPrivatePathServiceGateway(revokeAccountForPrivatePathServiceGatewayOptions *RevokeAccountForPrivatePathServiceGatewayOptions)`, func() {
			privatePathServiceGatewayRevokeAccountOptionsAccountModel := &vpcv1.PrivatePathServiceGatewayRevokeAccountOptionsAccountAccountIdentityByID{
				ID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
			}

			revokeAccountForPrivatePathServiceGatewayOptions := &vpcv1.RevokeAccountForPrivatePathServiceGatewayOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
				Account:                     privatePathServiceGatewayRevokeAccountOptionsAccountModel,
			}

			response, err := vpcService.RevokeAccountForPrivatePathServiceGateway(revokeAccountForPrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`UnpublishPrivatePathServiceGateway - Unpublish a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UnpublishPrivatePathServiceGateway(unpublishPrivatePathServiceGatewayOptions *UnpublishPrivatePathServiceGatewayOptions)`, func() {
			unpublishPrivatePathServiceGatewayOptions := &vpcv1.UnpublishPrivatePathServiceGatewayOptions{
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.UnpublishPrivatePathServiceGateway(unpublishPrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`ListPublicAddressRanges - List public address ranges`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPublicAddressRanges(listPublicAddressRangesOptions *ListPublicAddressRangesOptions) with pagination`, func() {
			listPublicAddressRangesOptions := &vpcv1.ListPublicAddressRangesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listPublicAddressRangesOptions.Start = nil
			listPublicAddressRangesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PublicAddressRange
			for {
				publicAddressRangeCollection, response, err := vpcService.ListPublicAddressRanges(listPublicAddressRangesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(publicAddressRangeCollection).ToNot(BeNil())
				allResults = append(allResults, publicAddressRangeCollection.PublicAddressRanges...)

				nextToken, nextErr := publicAddressRangeCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPublicAddressRangesOptions.Start = nextToken

				if listPublicAddressRangesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPublicAddressRanges(listPublicAddressRangesOptions *ListPublicAddressRangesOptions) using PublicAddressRangesPager`, func() {
			listPublicAddressRangesOptions := &vpcv1.ListPublicAddressRangesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewPublicAddressRangesPager(listPublicAddressRangesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PublicAddressRange
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPublicAddressRangesPager(listPublicAddressRangesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPublicAddressRanges() returned a total of %d item(s) using PublicAddressRangesPager.\n", len(allResults))
		})
	})

	Describe(`CreatePublicAddressRange - Create a public address range`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreatePublicAddressRange(createPublicAddressRangeOptions *CreatePublicAddressRangeOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			publicAddressRangeTargetPrototypeVPCModel := &vpcv1.PublicAddressRangeTargetPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			publicAddressRangeTargetPrototypeModel := &vpcv1.PublicAddressRangeTargetPrototype{
				VPC:  publicAddressRangeTargetPrototypeVPCModel,
				Zone: zoneIdentityModel,
			}

			createPublicAddressRangeOptions := &vpcv1.CreatePublicAddressRangeOptions{
				Ipv4AddressCount: core.Int64Ptr(int64(8)),
				Name:             core.StringPtr("my-public-address-range"),
				ResourceGroup:    resourceGroupIdentityModel,
				Target:           publicAddressRangeTargetPrototypeModel,
			}

			publicAddressRange, response, err := vpcService.CreatePublicAddressRange(createPublicAddressRangeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicAddressRange).ToNot(BeNil())
		})
	})

	Describe(`GetPublicAddressRange - Retrieve a public address range`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPublicAddressRange(getPublicAddressRangeOptions *GetPublicAddressRangeOptions)`, func() {
			getPublicAddressRangeOptions := &vpcv1.GetPublicAddressRangeOptions{
				ID: core.StringPtr("testString"),
			}

			publicAddressRange, response, err := vpcService.GetPublicAddressRange(getPublicAddressRangeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicAddressRange).ToNot(BeNil())
		})
	})

	Describe(`UpdatePublicAddressRange - Update a public address range`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdatePublicAddressRange(updatePublicAddressRangeOptions *UpdatePublicAddressRangeOptions)`, func() {
			vpcIdentityModel := &vpcv1.VPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			publicAddressRangeTargetPatchModel := &vpcv1.PublicAddressRangeTargetPatch{
				VPC:  vpcIdentityModel,
				Zone: zoneIdentityModel,
			}

			publicAddressRangePatchModel := &vpcv1.PublicAddressRangePatch{
				Name:   core.StringPtr("my-public-address-range-updated"),
				Target: publicAddressRangeTargetPatchModel,
			}
			publicAddressRangePatchModelAsPatch, asPatchErr := publicAddressRangePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePublicAddressRangeOptions := &vpcv1.UpdatePublicAddressRangeOptions{
				ID:                      core.StringPtr("testString"),
				PublicAddressRangePatch: publicAddressRangePatchModelAsPatch,
			}

			publicAddressRange, response, err := vpcService.UpdatePublicAddressRange(updatePublicAddressRangeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicAddressRange).ToNot(BeNil())
		})
	})

	Describe(`ListPublicGateways - List public gateways`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListPublicGateways(listPublicGatewaysOptions *ListPublicGatewaysOptions) with pagination`, func() {
			listPublicGatewaysOptions := &vpcv1.ListPublicGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listPublicGatewaysOptions.Start = nil
			listPublicGatewaysOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.PublicGateway
			for {
				publicGatewayCollection, response, err := vpcService.ListPublicGateways(listPublicGatewaysOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(publicGatewayCollection).ToNot(BeNil())
				allResults = append(allResults, publicGatewayCollection.PublicGateways...)

				nextToken, nextErr := publicGatewayCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listPublicGatewaysOptions.Start = nextToken

				if listPublicGatewaysOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListPublicGateways(listPublicGatewaysOptions *ListPublicGatewaysOptions) using PublicGatewaysPager`, func() {
			listPublicGatewaysOptions := &vpcv1.ListPublicGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewPublicGatewaysPager(listPublicGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.PublicGateway
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewPublicGatewaysPager(listPublicGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListPublicGateways() returned a total of %d item(s) using PublicGatewaysPager.\n", len(allResults))
		})
	})

	Describe(`CreatePublicGateway - Create a public gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreatePublicGateway(createPublicGatewayOptions *CreatePublicGatewayOptions)`, func() {
			publicGatewayPrototypeVPCModel := &vpcv1.PublicGatewayPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			publicGatewayPrototypeZoneModel := &vpcv1.PublicGatewayPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			publicGatewayFloatingIPPrototypeModel := &vpcv1.PublicGatewayFloatingIPPrototypeFloatingIPIdentityFloatingIPIdentityByAddress{
				Address: core.StringPtr("203.0.113.1"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createPublicGatewayOptions := &vpcv1.CreatePublicGatewayOptions{
				VPC:           publicGatewayPrototypeVPCModel,
				Zone:          publicGatewayPrototypeZoneModel,
				FloatingIP:    publicGatewayFloatingIPPrototypeModel,
				Name:          core.StringPtr("my-public-gateway"),
				ResourceGroup: resourceGroupIdentityModel,
			}

			publicGateway, response, err := vpcService.CreatePublicGateway(createPublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicGateway).ToNot(BeNil())
		})
	})

	Describe(`GetPublicGateway - Retrieve a public gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetPublicGateway(getPublicGatewayOptions *GetPublicGatewayOptions)`, func() {
			getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			publicGateway, response, err := vpcService.GetPublicGateway(getPublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
	})

	Describe(`UpdatePublicGateway - Update a public gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdatePublicGateway(updatePublicGatewayOptions *UpdatePublicGatewayOptions)`, func() {
			publicGatewayPatchModel := &vpcv1.PublicGatewayPatch{
				Name: core.StringPtr("my-public-gateway-updated"),
			}
			publicGatewayPatchModelAsPatch, asPatchErr := publicGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updatePublicGatewayOptions := &vpcv1.UpdatePublicGatewayOptions{
				ID:                 core.StringPtr("testString"),
				PublicGatewayPatch: publicGatewayPatchModelAsPatch,
			}

			publicGateway, response, err := vpcService.UpdatePublicGateway(updatePublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
	})

	Describe(`ListReservations - List reservations`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListReservations(listReservationsOptions *ListReservationsOptions) with pagination`, func() {
			listReservationsOptions := &vpcv1.ListReservationsOptions{
				AffinityPolicy:      core.StringPtr("automatic"),
				Limit:               core.Int64Ptr(int64(10)),
				Name:                core.StringPtr("my-name"),
				ProfileResourceType: core.StringPtr("testString"),
				ResourceGroupID:     core.StringPtr("testString"),
				Start:               core.StringPtr("testString"),
				ZoneName:            core.StringPtr("us-south-1"),
			}

			listReservationsOptions.Start = nil
			listReservationsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Reservation
			for {
				reservationCollection, response, err := vpcService.ListReservations(listReservationsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(reservationCollection).ToNot(BeNil())
				allResults = append(allResults, reservationCollection.Reservations...)

				nextToken, nextErr := reservationCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listReservationsOptions.Start = nextToken

				if listReservationsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListReservations(listReservationsOptions *ListReservationsOptions) using ReservationsPager`, func() {
			listReservationsOptions := &vpcv1.ListReservationsOptions{
				AffinityPolicy:      core.StringPtr("automatic"),
				Limit:               core.Int64Ptr(int64(10)),
				Name:                core.StringPtr("my-name"),
				ProfileResourceType: core.StringPtr("testString"),
				ResourceGroupID:     core.StringPtr("testString"),
				ZoneName:            core.StringPtr("us-south-1"),
			}

			// Test GetNext().
			pager, err := vpcService.NewReservationsPager(listReservationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Reservation
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewReservationsPager(listReservationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListReservations() returned a total of %d item(s) using ReservationsPager.\n", len(allResults))
		})
	})

	Describe(`CreateReservation - Create a reservation`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateReservation(createReservationOptions *CreateReservationOptions)`, func() {
			reservationCapacityPrototypeModel := &vpcv1.ReservationCapacityPrototype{
				Total: core.Int64Ptr(int64(10)),
			}

			reservationPrototypeCommittedUseModel := &vpcv1.ReservationPrototypeCommittedUse{
				ExpirationPolicy: core.StringPtr("release"),
				Term:             core.StringPtr("one_year"),
			}

			reservationProfilePrototypeModel := &vpcv1.ReservationProfilePrototype{
				Name:         core.StringPtr("bx2-4x16"),
				ResourceType: core.StringPtr("bare_metal_server_profile"),
			}

			reservationPrototypeZoneModel := &vpcv1.ReservationPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createReservationOptions := &vpcv1.CreateReservationOptions{
				Capacity:       reservationCapacityPrototypeModel,
				CommittedUse:   reservationPrototypeCommittedUseModel,
				Profile:        reservationProfilePrototypeModel,
				Zone:           reservationPrototypeZoneModel,
				AffinityPolicy: core.StringPtr("automatic"),
				Name:           core.StringPtr("my-reservation"),
				ResourceGroup:  resourceGroupIdentityModel,
			}

			reservation, response, err := vpcService.CreateReservation(createReservationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservation).ToNot(BeNil())
		})
	})

	Describe(`GetReservation - Retrieve a reservation`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetReservation(getReservationOptions *GetReservationOptions)`, func() {
			getReservationOptions := &vpcv1.GetReservationOptions{
				ID: core.StringPtr("testString"),
			}

			reservation, response, err := vpcService.GetReservation(getReservationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservation).ToNot(BeNil())
		})
	})

	Describe(`UpdateReservation - Update a reservation`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateReservation(updateReservationOptions *UpdateReservationOptions)`, func() {
			reservationCapacityPatchModel := &vpcv1.ReservationCapacityPatch{
				Total: core.Int64Ptr(int64(10)),
			}

			reservationPatchCommittedUseModel := &vpcv1.ReservationPatchCommittedUse{
				ExpirationPolicy: core.StringPtr("release"),
				Term:             core.StringPtr("one_year"),
			}

			reservationProfilePatchModel := &vpcv1.ReservationProfilePatch{
				Name:         core.StringPtr("bx2-4x16"),
				ResourceType: core.StringPtr("bare_metal_server_profile"),
			}

			reservationPatchModel := &vpcv1.ReservationPatch{
				AffinityPolicy: core.StringPtr("automatic"),
				Capacity:       reservationCapacityPatchModel,
				CommittedUse:   reservationPatchCommittedUseModel,
				Name:           core.StringPtr("my-reservation"),
				Profile:        reservationProfilePatchModel,
			}
			reservationPatchModelAsPatch, asPatchErr := reservationPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateReservationOptions := &vpcv1.UpdateReservationOptions{
				ID:               core.StringPtr("testString"),
				ReservationPatch: reservationPatchModelAsPatch,
			}

			reservation, response, err := vpcService.UpdateReservation(updateReservationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservation).ToNot(BeNil())
		})
	})

	Describe(`ActivateReservation - Activate a reservation`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ActivateReservation(activateReservationOptions *ActivateReservationOptions)`, func() {
			activateReservationOptions := &vpcv1.ActivateReservationOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.ActivateReservation(activateReservationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`ListSecurityGroups - List security groups`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSecurityGroups(listSecurityGroupsOptions *ListSecurityGroupsOptions) with pagination`, func() {
			listSecurityGroupsOptions := &vpcv1.ListSecurityGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			listSecurityGroupsOptions.Start = nil
			listSecurityGroupsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.SecurityGroup
			for {
				securityGroupCollection, response, err := vpcService.ListSecurityGroups(listSecurityGroupsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(securityGroupCollection).ToNot(BeNil())
				allResults = append(allResults, securityGroupCollection.SecurityGroups...)

				nextToken, nextErr := securityGroupCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSecurityGroupsOptions.Start = nextToken

				if listSecurityGroupsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSecurityGroups(listSecurityGroupsOptions *ListSecurityGroupsOptions) using SecurityGroupsPager`, func() {
			listSecurityGroupsOptions := &vpcv1.ListSecurityGroupsOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				VPCCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:           core.StringPtr("testString"),
				VPCName:         core.StringPtr("my-vpc"),
			}

			// Test GetNext().
			pager, err := vpcService.NewSecurityGroupsPager(listSecurityGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.SecurityGroup
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSecurityGroupsPager(listSecurityGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSecurityGroups() returned a total of %d item(s) using SecurityGroupsPager.\n", len(allResults))
		})
	})

	Describe(`CreateSecurityGroup - Create a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSecurityGroup(createSecurityGroupOptions *CreateSecurityGroupOptions)`, func() {
			securityGroupPrototypeVPCModel := &vpcv1.SecurityGroupPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupRuleLocalPrototypeModel := &vpcv1.SecurityGroupRuleLocalPrototypeSecurityGroupRuleCIDRPrototype{
				CIDRBlock: core.StringPtr("testString"),
			}

			securityGroupRuleRemotePrototypeModel := &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			securityGroupRulePrototypeModel := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAnyPrototype{
				Direction: core.StringPtr("inbound"),
				IPVersion: core.StringPtr("ipv4"),
				Local:     securityGroupRuleLocalPrototypeModel,
				Name:      core.StringPtr("my-rule-1"),
				Remote:    securityGroupRuleRemotePrototypeModel,
				Protocol:  core.StringPtr("any"),
			}

			createSecurityGroupOptions := &vpcv1.CreateSecurityGroupOptions{
				VPC:           securityGroupPrototypeVPCModel,
				Name:          core.StringPtr("my-security-group"),
				ResourceGroup: resourceGroupIdentityModel,
				Rules:         []vpcv1.SecurityGroupRulePrototypeIntf{securityGroupRulePrototypeModel},
			}

			securityGroup, response, err := vpcService.CreateSecurityGroup(createSecurityGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(securityGroup).ToNot(BeNil())
		})
	})

	Describe(`GetSecurityGroup - Retrieve a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSecurityGroup(getSecurityGroupOptions *GetSecurityGroupOptions)`, func() {
			getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
				ID: core.StringPtr("testString"),
			}

			securityGroup, response, err := vpcService.GetSecurityGroup(getSecurityGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroup).ToNot(BeNil())
		})
	})

	Describe(`UpdateSecurityGroup - Update a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSecurityGroup(updateSecurityGroupOptions *UpdateSecurityGroupOptions)`, func() {
			securityGroupPatchModel := &vpcv1.SecurityGroupPatch{
				Name: core.StringPtr("my-security-group"),
			}
			securityGroupPatchModelAsPatch, asPatchErr := securityGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSecurityGroupOptions := &vpcv1.UpdateSecurityGroupOptions{
				ID:                 core.StringPtr("testString"),
				SecurityGroupPatch: securityGroupPatchModelAsPatch,
			}

			securityGroup, response, err := vpcService.UpdateSecurityGroup(updateSecurityGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroup).ToNot(BeNil())
		})
	})

	Describe(`ListSecurityGroupRules - List rules in a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSecurityGroupRules(listSecurityGroupRulesOptions *ListSecurityGroupRulesOptions)`, func() {
			listSecurityGroupRulesOptions := &vpcv1.ListSecurityGroupRulesOptions{
				SecurityGroupID: core.StringPtr("testString"),
			}

			securityGroupRuleCollection, response, err := vpcService.ListSecurityGroupRules(listSecurityGroupRulesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRuleCollection).ToNot(BeNil())
		})
	})

	Describe(`CreateSecurityGroupRule - Create a rule for a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSecurityGroupRule(createSecurityGroupRuleOptions *CreateSecurityGroupRuleOptions)`, func() {
			securityGroupRuleLocalPrototypeModel := &vpcv1.SecurityGroupRuleLocalPrototypeSecurityGroupRuleCIDRPrototype{
				CIDRBlock: core.StringPtr("testString"),
			}

			securityGroupRuleRemotePrototypeModel := &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			securityGroupRulePrototypeModel := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAnyPrototype{
				Direction: core.StringPtr("inbound"),
				IPVersion: core.StringPtr("ipv4"),
				Local:     securityGroupRuleLocalPrototypeModel,
				Name:      core.StringPtr("my-rule-1"),
				Remote:    securityGroupRuleRemotePrototypeModel,
				Protocol:  core.StringPtr("any"),
			}

			createSecurityGroupRuleOptions := &vpcv1.CreateSecurityGroupRuleOptions{
				SecurityGroupID:            core.StringPtr("testString"),
				SecurityGroupRulePrototype: securityGroupRulePrototypeModel,
			}

			securityGroupRule, response, err := vpcService.CreateSecurityGroupRule(createSecurityGroupRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(securityGroupRule).ToNot(BeNil())
		})
	})

	Describe(`GetSecurityGroupRule - Retrieve a security group rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSecurityGroupRule(getSecurityGroupRuleOptions *GetSecurityGroupRuleOptions)`, func() {
			getSecurityGroupRuleOptions := &vpcv1.GetSecurityGroupRuleOptions{
				ID:              core.StringPtr("testString"),
				SecurityGroupID: core.StringPtr("testString"),
			}

			securityGroupRule, response, err := vpcService.GetSecurityGroupRule(getSecurityGroupRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRule).ToNot(BeNil())
		})
	})

	Describe(`UpdateSecurityGroupRule - Update a security group rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSecurityGroupRule(updateSecurityGroupRuleOptions *UpdateSecurityGroupRuleOptions)`, func() {
			securityGroupRuleLocalPatchModel := &vpcv1.SecurityGroupRuleLocalPatchSecurityGroupRuleCIDRPrototype{
				CIDRBlock: core.StringPtr("testString"),
			}

			securityGroupRuleRemotePatchModel := &vpcv1.SecurityGroupRuleRemotePatchSecurityGroupRuleCIDRPrototype{
				CIDRBlock: core.StringPtr("10.0.0.0/8"),
			}

			securityGroupRulePatchModel := &vpcv1.SecurityGroupRulePatch{
				Code:      core.Int64Ptr(int64(0)),
				Direction: core.StringPtr("outbound"),
				IPVersion: core.StringPtr("ipv4"),
				Local:     securityGroupRuleLocalPatchModel,
				Name:      core.StringPtr("my-security-group-rule"),
				PortMax:   core.Int64Ptr(int64(22)),
				PortMin:   core.Int64Ptr(int64(22)),
				Remote:    securityGroupRuleRemotePatchModel,
				Type:      core.Int64Ptr(int64(8)),
			}
			securityGroupRulePatchModelAsPatch, asPatchErr := securityGroupRulePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSecurityGroupRuleOptions := &vpcv1.UpdateSecurityGroupRuleOptions{
				ID:                     core.StringPtr("testString"),
				SecurityGroupID:        core.StringPtr("testString"),
				SecurityGroupRulePatch: securityGroupRulePatchModelAsPatch,
			}

			securityGroupRule, response, err := vpcService.UpdateSecurityGroupRule(updateSecurityGroupRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupRule).ToNot(BeNil())
		})
	})

	Describe(`ListSecurityGroupTargets - List targets associated with a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSecurityGroupTargets(listSecurityGroupTargetsOptions *ListSecurityGroupTargetsOptions) with pagination`, func() {
			listSecurityGroupTargetsOptions := &vpcv1.ListSecurityGroupTargetsOptions{
				SecurityGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
				Start:           core.StringPtr("testString"),
			}

			listSecurityGroupTargetsOptions.Start = nil
			listSecurityGroupTargetsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.SecurityGroupTargetReferenceIntf
			for {
				securityGroupTargetCollection, response, err := vpcService.ListSecurityGroupTargets(listSecurityGroupTargetsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(securityGroupTargetCollection).ToNot(BeNil())
				allResults = append(allResults, securityGroupTargetCollection.Targets...)

				nextToken, nextErr := securityGroupTargetCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSecurityGroupTargetsOptions.Start = nextToken

				if listSecurityGroupTargetsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSecurityGroupTargets(listSecurityGroupTargetsOptions *ListSecurityGroupTargetsOptions) using SecurityGroupTargetsPager`, func() {
			listSecurityGroupTargetsOptions := &vpcv1.ListSecurityGroupTargetsOptions{
				SecurityGroupID: core.StringPtr("testString"),
				Limit:           core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewSecurityGroupTargetsPager(listSecurityGroupTargetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.SecurityGroupTargetReferenceIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSecurityGroupTargetsPager(listSecurityGroupTargetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSecurityGroupTargets() returned a total of %d item(s) using SecurityGroupTargetsPager.\n", len(allResults))
		})
	})

	Describe(`GetSecurityGroupTarget - Retrieve a security group target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSecurityGroupTarget(getSecurityGroupTargetOptions *GetSecurityGroupTargetOptions)`, func() {
			getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
				ID:              core.StringPtr("testString"),
				SecurityGroupID: core.StringPtr("testString"),
			}

			securityGroupTargetReference, response, err := vpcService.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupTargetReference).ToNot(BeNil())
		})
	})

	Describe(`CreateSecurityGroupTargetBinding - Add a target to a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions *CreateSecurityGroupTargetBindingOptions)`, func() {
			createSecurityGroupTargetBindingOptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
				ID:              core.StringPtr("testString"),
				SecurityGroupID: core.StringPtr("testString"),
			}

			securityGroupTargetReference, response, err := vpcService.CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(securityGroupTargetReference).ToNot(BeNil())
		})
	})

	Describe(`ListShareProfiles - List file share profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListShareProfiles(listShareProfilesOptions *ListShareProfilesOptions) with pagination`, func() {
			listShareProfilesOptions := &vpcv1.ListShareProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Sort:  core.StringPtr("name"),
				Start: core.StringPtr("testString"),
			}

			listShareProfilesOptions.Start = nil
			listShareProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ShareProfile
			for {
				shareProfileCollection, response, err := vpcService.ListShareProfiles(listShareProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(shareProfileCollection).ToNot(BeNil())
				allResults = append(allResults, shareProfileCollection.Profiles...)

				nextToken, nextErr := shareProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listShareProfilesOptions.Start = nextToken

				if listShareProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListShareProfiles(listShareProfilesOptions *ListShareProfilesOptions) using ShareProfilesPager`, func() {
			listShareProfilesOptions := &vpcv1.ListShareProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Sort:  core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewShareProfilesPager(listShareProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ShareProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewShareProfilesPager(listShareProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListShareProfiles() returned a total of %d item(s) using ShareProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetShareProfile - Retrieve a file share profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShareProfile(getShareProfileOptions *GetShareProfileOptions)`, func() {
			getShareProfileOptions := &vpcv1.GetShareProfileOptions{
				Name: core.StringPtr("dp2"),
			}

			shareProfile, response, err := vpcService.GetShareProfile(getShareProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareProfile).ToNot(BeNil())
		})
	})

	Describe(`ListShares - List file shares`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListShares(listSharesOptions *ListSharesOptions) with pagination`, func() {
			listSharesOptions := &vpcv1.ListSharesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ReplicationRole: core.StringPtr("none"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				Start:           core.StringPtr("testString"),
			}

			listSharesOptions.Start = nil
			listSharesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Share
			for {
				shareCollection, response, err := vpcService.ListShares(listSharesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(shareCollection).ToNot(BeNil())
				allResults = append(allResults, shareCollection.Shares...)

				nextToken, nextErr := shareCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSharesOptions.Start = nextToken

				if listSharesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListShares(listSharesOptions *ListSharesOptions) using SharesPager`, func() {
			listSharesOptions := &vpcv1.ListSharesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ReplicationRole: core.StringPtr("none"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewSharesPager(listSharesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Share
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSharesPager(listSharesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListShares() returned a total of %d item(s) using SharesPager.\n", len(allResults))
		})
	})

	Describe(`CreateShare - Create a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateShare(createShareOptions *CreateShareOptions)`, func() {
			shareMountTargetVirtualNetworkInterfacePrototypeModel := &vpcv1.ShareMountTargetVirtualNetworkInterfacePrototypeVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			shareMountTargetPrototypeModel := &vpcv1.ShareMountTargetPrototypeShareMountTargetByAccessControlModeSecurityGroup{
				AccessProtocol:          core.StringPtr("nfs4"),
				Name:                    core.StringPtr("my-share-mount-target"),
				TransitEncryption:       core.StringPtr("none"),
				VirtualNetworkInterface: shareMountTargetVirtualNetworkInterfacePrototypeModel,
			}

			shareProfileIdentityModel := &vpcv1.ShareProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/share/profiles/tier-3iops"),
			}

			sharePrototypeShareContextResourceGroupModel := &vpcv1.SharePrototypeShareContextResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			sharePrototypeShareContextModel := &vpcv1.SharePrototypeShareContext{
				AllowedTransitEncryptionModes: []string{"ipsec"},
				Iops:                          core.Int64Ptr(int64(100)),
				MountTargets:                  []vpcv1.ShareMountTargetPrototypeIntf{shareMountTargetPrototypeModel},
				Name:                          core.StringPtr("my-share"),
				Profile:                       shareProfileIdentityModel,
				ReplicationCronSpec:           core.StringPtr("0 */5 * * *"),
				ResourceGroup:                 sharePrototypeShareContextResourceGroupModel,
				UserTags:                      []string{"testString"},
				Zone:                          zoneIdentityModel,
			}

			shareByOriginShareOriginShareModel := &vpcv1.ShareByOriginShareOriginShareShareIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::share:r006-0fe9e5d8-0a4d-4818-96ec-e99708644a58"),
			}

			shareByOriginShareResourceGroupModel := &vpcv1.ShareByOriginShareResourceGroupResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			sharePrototypeModel := &vpcv1.SharePrototypeShareByOriginShare{
				AllowedTransitEncryptionModes: []string{"ipsec"},
				MountTargets:                  []vpcv1.ShareMountTargetPrototypeIntf{shareMountTargetPrototypeModel},
				Name:                          core.StringPtr("my-share"),
				ReplicaShare:                  sharePrototypeShareContextModel,
				UserTags:                      []string{"testString"},
				OriginShare:                   shareByOriginShareOriginShareModel,
				ResourceGroup:                 shareByOriginShareResourceGroupModel,
			}

			createShareOptions := &vpcv1.CreateShareOptions{
				SharePrototype: sharePrototypeModel,
			}

			share, response, err := vpcService.CreateShare(createShareOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(share).ToNot(BeNil())
		})
	})

	Describe(`GetShare - Retrieve a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShare(getShareOptions *GetShareOptions)`, func() {
			getShareOptions := &vpcv1.GetShareOptions{
				ID: core.StringPtr("testString"),
			}

			share, response, err := vpcService.GetShare(getShareOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(share).ToNot(BeNil())
		})
	})

	Describe(`UpdateShare - Update a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateShare(updateShareOptions *UpdateShareOptions)`, func() {
			sharePatchProfileModel := &vpcv1.SharePatchProfileShareProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/share/profiles/tier-3iops"),
			}

			sharePatchModel := &vpcv1.SharePatch{
				AccessControlMode:             core.StringPtr("security_group"),
				AllowedAccessProtocols:        []string{"nfs4"},
				AllowedTransitEncryptionModes: []string{"ipsec"},
				Bandwidth:                     core.Int64Ptr(int64(1)),
				Iops:                          core.Int64Ptr(int64(100)),
				Name:                          core.StringPtr("my-share"),
				Profile:                       sharePatchProfileModel,
				ReplicationCronSpec:           core.StringPtr("0 */5 * * *"),
				Size:                          core.Int64Ptr(int64(8000)),
				UserTags:                      []string{"testString"},
			}
			sharePatchModelAsPatch, asPatchErr := sharePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareOptions := &vpcv1.UpdateShareOptions{
				ID:         core.StringPtr("testString"),
				SharePatch: sharePatchModelAsPatch,
				IfMatch:    core.StringPtr("W/96d225c4-56bd-43d9-98fc-d7148e5c5028"),
			}

			share, response, err := vpcService.UpdateShare(updateShareOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(share).ToNot(BeNil())
		})
	})

	Describe(`ListShareAccessorBindings - List accessor bindings for a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListShareAccessorBindings(listShareAccessorBindingsOptions *ListShareAccessorBindingsOptions) with pagination`, func() {
			listShareAccessorBindingsOptions := &vpcv1.ListShareAccessorBindingsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listShareAccessorBindingsOptions.Start = nil
			listShareAccessorBindingsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ShareAccessorBinding
			for {
				shareAccessorBindingCollection, response, err := vpcService.ListShareAccessorBindings(listShareAccessorBindingsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(shareAccessorBindingCollection).ToNot(BeNil())
				allResults = append(allResults, shareAccessorBindingCollection.AccessorBindings...)

				nextToken, nextErr := shareAccessorBindingCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listShareAccessorBindingsOptions.Start = nextToken

				if listShareAccessorBindingsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListShareAccessorBindings(listShareAccessorBindingsOptions *ListShareAccessorBindingsOptions) using ShareAccessorBindingsPager`, func() {
			listShareAccessorBindingsOptions := &vpcv1.ListShareAccessorBindingsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewShareAccessorBindingsPager(listShareAccessorBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ShareAccessorBinding
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewShareAccessorBindingsPager(listShareAccessorBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListShareAccessorBindings() returned a total of %d item(s) using ShareAccessorBindingsPager.\n", len(allResults))
		})
	})

	Describe(`GetShareAccessorBinding - Retrieve a file share accessor binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShareAccessorBinding(getShareAccessorBindingOptions *GetShareAccessorBindingOptions)`, func() {
			getShareAccessorBindingOptions := &vpcv1.GetShareAccessorBindingOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			shareAccessorBinding, response, err := vpcService.GetShareAccessorBinding(getShareAccessorBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareAccessorBinding).ToNot(BeNil())
		})
	})

	Describe(`FailoverShare - Failover to replica file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`FailoverShare(failoverShareOptions *FailoverShareOptions)`, func() {
			failoverShareOptions := &vpcv1.FailoverShareOptions{
				ShareID:        core.StringPtr("testString"),
				FallbackPolicy: core.StringPtr("fail"),
				Timeout:        core.Int64Ptr(int64(600)),
			}

			response, err := vpcService.FailoverShare(failoverShareOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`ListShareMountTargets - List mount targets for a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListShareMountTargets(listShareMountTargetsOptions *ListShareMountTargetsOptions) with pagination`, func() {
			listShareMountTargetsOptions := &vpcv1.ListShareMountTargetsOptions{
				ShareID: core.StringPtr("testString"),
				Limit:   core.Int64Ptr(int64(10)),
				Name:    core.StringPtr("my-name"),
				Start:   core.StringPtr("testString"),
			}

			listShareMountTargetsOptions.Start = nil
			listShareMountTargetsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ShareMountTarget
			for {
				shareMountTargetCollection, response, err := vpcService.ListShareMountTargets(listShareMountTargetsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(shareMountTargetCollection).ToNot(BeNil())
				allResults = append(allResults, shareMountTargetCollection.MountTargets...)

				nextToken, nextErr := shareMountTargetCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listShareMountTargetsOptions.Start = nextToken

				if listShareMountTargetsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListShareMountTargets(listShareMountTargetsOptions *ListShareMountTargetsOptions) using ShareMountTargetsPager`, func() {
			listShareMountTargetsOptions := &vpcv1.ListShareMountTargetsOptions{
				ShareID: core.StringPtr("testString"),
				Limit:   core.Int64Ptr(int64(10)),
				Name:    core.StringPtr("my-name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewShareMountTargetsPager(listShareMountTargetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ShareMountTarget
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewShareMountTargetsPager(listShareMountTargetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListShareMountTargets() returned a total of %d item(s) using ShareMountTargetsPager.\n", len(allResults))
		})
	})

	Describe(`CreateShareMountTarget - Create a mount target for a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateShareMountTarget(createShareMountTargetOptions *CreateShareMountTargetOptions)`, func() {
			shareMountTargetVirtualNetworkInterfacePrototypeModel := &vpcv1.ShareMountTargetVirtualNetworkInterfacePrototypeVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::virtual-network-interface:0717-54eb57ee-86f2-4796-90bb-d7874e0831ef"),
			}

			shareMountTargetPrototypeModel := &vpcv1.ShareMountTargetPrototypeShareMountTargetByAccessControlModeSecurityGroup{
				AccessProtocol:          core.StringPtr("nfs4"),
				Name:                    core.StringPtr("my-share-mount-target"),
				TransitEncryption:       core.StringPtr("none"),
				VirtualNetworkInterface: shareMountTargetVirtualNetworkInterfacePrototypeModel,
			}

			createShareMountTargetOptions := &vpcv1.CreateShareMountTargetOptions{
				ShareID:                   core.StringPtr("testString"),
				ShareMountTargetPrototype: shareMountTargetPrototypeModel,
			}

			shareMountTarget, response, err := vpcService.CreateShareMountTarget(createShareMountTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(shareMountTarget).ToNot(BeNil())
		})
	})

	Describe(`GetShareMountTarget - Retrieve a file share mount target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShareMountTarget(getShareMountTargetOptions *GetShareMountTargetOptions)`, func() {
			getShareMountTargetOptions := &vpcv1.GetShareMountTargetOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			shareMountTarget, response, err := vpcService.GetShareMountTarget(getShareMountTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareMountTarget).ToNot(BeNil())
		})
	})

	Describe(`UpdateShareMountTarget - Update a file share mount target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateShareMountTarget(updateShareMountTargetOptions *UpdateShareMountTargetOptions)`, func() {
			shareMountTargetPatchModel := &vpcv1.ShareMountTargetPatch{
				Name: core.StringPtr("my-share-mount-target"),
			}
			shareMountTargetPatchModelAsPatch, asPatchErr := shareMountTargetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareMountTargetOptions := &vpcv1.UpdateShareMountTargetOptions{
				ID:                    core.StringPtr("testString"),
				ShareID:               core.StringPtr("testString"),
				ShareMountTargetPatch: shareMountTargetPatchModelAsPatch,
			}

			shareMountTarget, response, err := vpcService.UpdateShareMountTarget(updateShareMountTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareMountTarget).ToNot(BeNil())
		})
	})

	Describe(`ListShareSnapshots - List file share snapshots`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListShareSnapshots(listShareSnapshotsOptions *ListShareSnapshotsOptions) with pagination`, func() {
			listShareSnapshotsOptions := &vpcv1.ListShareSnapshotsOptions{
				ShareID:            core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				Sort:               core.StringPtr("name"),
				Start:              core.StringPtr("testString"),
			}

			listShareSnapshotsOptions.Start = nil
			listShareSnapshotsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ShareSnapshot
			for {
				shareSnapshotCollection, response, err := vpcService.ListShareSnapshots(listShareSnapshotsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(shareSnapshotCollection).ToNot(BeNil())
				allResults = append(allResults, shareSnapshotCollection.Snapshots...)

				nextToken, nextErr := shareSnapshotCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listShareSnapshotsOptions.Start = nextToken

				if listShareSnapshotsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListShareSnapshots(listShareSnapshotsOptions *ListShareSnapshotsOptions) using ShareSnapshotsPager`, func() {
			listShareSnapshotsOptions := &vpcv1.ListShareSnapshotsOptions{
				ShareID:            core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				Sort:               core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewShareSnapshotsPager(listShareSnapshotsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ShareSnapshot
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewShareSnapshotsPager(listShareSnapshotsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListShareSnapshots() returned a total of %d item(s) using ShareSnapshotsPager.\n", len(allResults))
		})
	})

	Describe(`CreateShareSnapshot - Create a snapshot for a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateShareSnapshot(createShareSnapshotOptions *CreateShareSnapshotOptions)`, func() {
			createShareSnapshotOptions := &vpcv1.CreateShareSnapshotOptions{
				ShareID:  core.StringPtr("testString"),
				Name:     core.StringPtr("my-share-snapshot"),
				UserTags: []string{"testString"},
			}

			shareSnapshot, response, err := vpcService.CreateShareSnapshot(createShareSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(shareSnapshot).ToNot(BeNil())
		})
	})

	Describe(`GetShareSnapshot - Retrieve a share snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShareSnapshot(getShareSnapshotOptions *GetShareSnapshotOptions)`, func() {
			getShareSnapshotOptions := &vpcv1.GetShareSnapshotOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			shareSnapshot, response, err := vpcService.GetShareSnapshot(getShareSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareSnapshot).ToNot(BeNil())
		})
	})

	Describe(`UpdateShareSnapshot - Update a share snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateShareSnapshot(updateShareSnapshotOptions *UpdateShareSnapshotOptions)`, func() {
			shareSnapshotPatchModel := &vpcv1.ShareSnapshotPatch{
				Name:     core.StringPtr("my-share-snapshot"),
				UserTags: []string{"testString"},
			}
			shareSnapshotPatchModelAsPatch, asPatchErr := shareSnapshotPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateShareSnapshotOptions := &vpcv1.UpdateShareSnapshotOptions{
				ID:                 core.StringPtr("testString"),
				ShareID:            core.StringPtr("testString"),
				ShareSnapshotPatch: shareSnapshotPatchModelAsPatch,
				IfMatch:            core.StringPtr("W/96d225c4-56bd-43d9-98fc-d7148e5c5028"),
			}

			shareSnapshot, response, err := vpcService.UpdateShareSnapshot(updateShareSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareSnapshot).ToNot(BeNil())
		})
	})

	Describe(`GetShareSource - Retrieve the source file share for a replica file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetShareSource(getShareSourceOptions *GetShareSourceOptions)`, func() {
			getShareSourceOptions := &vpcv1.GetShareSourceOptions{
				ShareID: core.StringPtr("testString"),
			}

			shareReference, response, err := vpcService.GetShareSource(getShareSourceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(shareReference).ToNot(BeNil())
		})
	})

	Describe(`ListSnapshotConsistencyGroups - List snapshot consistency groups`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSnapshotConsistencyGroups(listSnapshotConsistencyGroupsOptions *ListSnapshotConsistencyGroupsOptions) with pagination`, func() {
			listSnapshotConsistencyGroupsOptions := &vpcv1.ListSnapshotConsistencyGroupsOptions{
				BackupPolicyJobID:  core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
				Start:              core.StringPtr("testString"),
			}

			listSnapshotConsistencyGroupsOptions.Start = nil
			listSnapshotConsistencyGroupsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.SnapshotConsistencyGroup
			for {
				snapshotConsistencyGroupCollection, response, err := vpcService.ListSnapshotConsistencyGroups(listSnapshotConsistencyGroupsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(snapshotConsistencyGroupCollection).ToNot(BeNil())
				allResults = append(allResults, snapshotConsistencyGroupCollection.SnapshotConsistencyGroups...)

				nextToken, nextErr := snapshotConsistencyGroupCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSnapshotConsistencyGroupsOptions.Start = nextToken

				if listSnapshotConsistencyGroupsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSnapshotConsistencyGroups(listSnapshotConsistencyGroupsOptions *ListSnapshotConsistencyGroupsOptions) using SnapshotConsistencyGroupsPager`, func() {
			listSnapshotConsistencyGroupsOptions := &vpcv1.ListSnapshotConsistencyGroupsOptions{
				BackupPolicyJobID:  core.StringPtr("testString"),
				BackupPolicyPlanID: core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Name:               core.StringPtr("my-name"),
				ResourceGroupID:    core.StringPtr("testString"),
				Sort:               core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewSnapshotConsistencyGroupsPager(listSnapshotConsistencyGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.SnapshotConsistencyGroup
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSnapshotConsistencyGroupsPager(listSnapshotConsistencyGroupsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSnapshotConsistencyGroups() returned a total of %d item(s) using SnapshotConsistencyGroupsPager.\n", len(allResults))
		})
	})

	Describe(`CreateSnapshotConsistencyGroup - Create a snapshot consistency group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSnapshotConsistencyGroup(createSnapshotConsistencyGroupOptions *CreateSnapshotConsistencyGroupOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			volumeIdentityModel := &vpcv1.VolumeIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"),
			}

			snapshotPrototypeSnapshotConsistencyGroupContextModel := &vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext{
				Name:         core.StringPtr("my-snapshot"),
				SourceVolume: volumeIdentityModel,
				UserTags:     []string{"testString"},
			}

			snapshotConsistencyGroupPrototypeModel := &vpcv1.SnapshotConsistencyGroupPrototypeSnapshotConsistencyGroupBySnapshots{
				DeleteSnapshotsOnDelete: core.BoolPtr(true),
				Name:                    core.StringPtr("my-snapshot-consistency-group"),
				ResourceGroup:           resourceGroupIdentityModel,
				Snapshots:               []vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext{*snapshotPrototypeSnapshotConsistencyGroupContextModel},
			}

			createSnapshotConsistencyGroupOptions := &vpcv1.CreateSnapshotConsistencyGroupOptions{
				SnapshotConsistencyGroupPrototype: snapshotConsistencyGroupPrototypeModel,
			}

			snapshotConsistencyGroup, response, err := vpcService.CreateSnapshotConsistencyGroup(createSnapshotConsistencyGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
	})

	Describe(`GetSnapshotConsistencyGroup - Retrieve a snapshot consistency group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions *GetSnapshotConsistencyGroupOptions)`, func() {
			getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{
				ID: core.StringPtr("testString"),
			}

			snapshotConsistencyGroup, response, err := vpcService.GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
	})

	Describe(`UpdateSnapshotConsistencyGroup - Update a snapshot consistency group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSnapshotConsistencyGroup(updateSnapshotConsistencyGroupOptions *UpdateSnapshotConsistencyGroupOptions)`, func() {
			snapshotConsistencyGroupPatchModel := &vpcv1.SnapshotConsistencyGroupPatch{
				DeleteSnapshotsOnDelete: core.BoolPtr(false),
				Name:                    core.StringPtr("my-snapshot-consistency-group"),
			}
			snapshotConsistencyGroupPatchModelAsPatch, asPatchErr := snapshotConsistencyGroupPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotConsistencyGroupOptions := &vpcv1.UpdateSnapshotConsistencyGroupOptions{
				ID:                            core.StringPtr("testString"),
				SnapshotConsistencyGroupPatch: snapshotConsistencyGroupPatchModelAsPatch,
				IfMatch:                       core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			snapshotConsistencyGroup, response, err := vpcService.UpdateSnapshotConsistencyGroup(updateSnapshotConsistencyGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
	})

	Describe(`ListSnapshots - List snapshots`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSnapshots(listSnapshotsOptions *ListSnapshotsOptions) with pagination`, func() {
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
				Start:                          core.StringPtr("testString"),
				Tag:                            core.StringPtr("testString"),
			}

			listSnapshotsOptions.Start = nil
			listSnapshotsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Snapshot
			for {
				snapshotCollection, response, err := vpcService.ListSnapshots(listSnapshotsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(snapshotCollection).ToNot(BeNil())
				allResults = append(allResults, snapshotCollection.Snapshots...)

				nextToken, nextErr := snapshotCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSnapshotsOptions.Start = nextToken

				if listSnapshotsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSnapshots(listSnapshotsOptions *ListSnapshotsOptions) using SnapshotsPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewSnapshotsPager(listSnapshotsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Snapshot
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSnapshotsPager(listSnapshotsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSnapshots() returned a total of %d item(s) using SnapshotsPager.\n", len(allResults))
		})
	})

	Describe(`CreateSnapshot - Create a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSnapshot(createSnapshotOptions *CreateSnapshotOptions)`, func() {
			snapshotAllowedUsePrototypeModel := &vpcv1.SnapshotAllowedUsePrototype{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			snapshotClonePrototypeZoneModel := &vpcv1.SnapshotClonePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			snapshotClonePrototypeModel := &vpcv1.SnapshotClonePrototype{
				Zone: snapshotClonePrototypeZoneModel,
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			encryptionKeyIdentityModel := &vpcv1.EncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			snapshotIdentityByCRNModel := &vpcv1.SnapshotIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::snapshot:r006-f6bfa329-0e36-433f-a3bb-0df632e79263"),
			}

			snapshotPrototypeModel := &vpcv1.SnapshotPrototypeSnapshotBySourceSnapshot{
				AllowedUse:     snapshotAllowedUsePrototypeModel,
				Clones:         []vpcv1.SnapshotClonePrototype{*snapshotClonePrototypeModel},
				Name:           core.StringPtr("my-snapshot"),
				ResourceGroup:  resourceGroupIdentityModel,
				UserTags:       []string{"testString"},
				EncryptionKey:  encryptionKeyIdentityModel,
				SourceSnapshot: snapshotIdentityByCRNModel,
			}

			createSnapshotOptions := &vpcv1.CreateSnapshotOptions{
				SnapshotPrototype: snapshotPrototypeModel,
			}

			snapshot, response, err := vpcService.CreateSnapshot(createSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(snapshot).ToNot(BeNil())
		})
	})

	Describe(`GetSnapshot - Retrieve a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSnapshot(getSnapshotOptions *GetSnapshotOptions)`, func() {
			getSnapshotOptions := &vpcv1.GetSnapshotOptions{
				ID: core.StringPtr("testString"),
			}

			snapshot, response, err := vpcService.GetSnapshot(getSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshot).ToNot(BeNil())
		})
	})

	Describe(`UpdateSnapshot - Update a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSnapshot(updateSnapshotOptions *UpdateSnapshotOptions)`, func() {
			snapshotAllowedUsePatchModel := &vpcv1.SnapshotAllowedUsePatch{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			snapshotPatchModel := &vpcv1.SnapshotPatch{
				AllowedUse: snapshotAllowedUsePatchModel,
				Name:       core.StringPtr("my-snapshot-updated"),
				UserTags:   []string{"testString"},
			}
			snapshotPatchModelAsPatch, asPatchErr := snapshotPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotOptions := &vpcv1.UpdateSnapshotOptions{
				ID:            core.StringPtr("testString"),
				SnapshotPatch: snapshotPatchModelAsPatch,
				IfMatch:       core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			snapshot, response, err := vpcService.UpdateSnapshot(updateSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshot).ToNot(BeNil())
		})
	})

	Describe(`ListSnapshotClones - List clones for a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSnapshotClones(listSnapshotClonesOptions *ListSnapshotClonesOptions)`, func() {
			listSnapshotClonesOptions := &vpcv1.ListSnapshotClonesOptions{
				ID: core.StringPtr("testString"),
			}

			snapshotCloneCollection, response, err := vpcService.ListSnapshotClones(listSnapshotClonesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotCloneCollection).ToNot(BeNil())
		})
	})

	Describe(`GetSnapshotClone - Retrieve a snapshot clone`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSnapshotClone(getSnapshotCloneOptions *GetSnapshotCloneOptions)`, func() {
			getSnapshotCloneOptions := &vpcv1.GetSnapshotCloneOptions{
				ID:       core.StringPtr("testString"),
				ZoneName: core.StringPtr("us-south-1"),
			}

			snapshotClone, response, err := vpcService.GetSnapshotClone(getSnapshotCloneOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotClone).ToNot(BeNil())
		})
	})

	Describe(`CreateSnapshotClone - Create a clone for a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSnapshotClone(createSnapshotCloneOptions *CreateSnapshotCloneOptions)`, func() {
			createSnapshotCloneOptions := &vpcv1.CreateSnapshotCloneOptions{
				ID:       core.StringPtr("testString"),
				ZoneName: core.StringPtr("us-south-1"),
			}

			snapshotClone, response, err := vpcService.CreateSnapshotClone(createSnapshotCloneOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotClone).ToNot(BeNil())
		})
	})

	Describe(`ListSnapshotInstanceProfiles - List instance profiles compatible with a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSnapshotInstanceProfiles(listSnapshotInstanceProfilesOptions *ListSnapshotInstanceProfilesOptions) with pagination`, func() {
			listSnapshotInstanceProfilesOptions := &vpcv1.ListSnapshotInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listSnapshotInstanceProfilesOptions.Start = nil
			listSnapshotInstanceProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceProfileReference
			for {
				snapshotInstanceProfileCollection, response, err := vpcService.ListSnapshotInstanceProfiles(listSnapshotInstanceProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(snapshotInstanceProfileCollection).ToNot(BeNil())
				allResults = append(allResults, snapshotInstanceProfileCollection.InstanceProfiles...)

				nextToken, nextErr := snapshotInstanceProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSnapshotInstanceProfilesOptions.Start = nextToken

				if listSnapshotInstanceProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSnapshotInstanceProfiles(listSnapshotInstanceProfilesOptions *ListSnapshotInstanceProfilesOptions) using SnapshotInstanceProfilesPager`, func() {
			listSnapshotInstanceProfilesOptions := &vpcv1.ListSnapshotInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewSnapshotInstanceProfilesPager(listSnapshotInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSnapshotInstanceProfilesPager(listSnapshotInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSnapshotInstanceProfiles() returned a total of %d item(s) using SnapshotInstanceProfilesPager.\n", len(allResults))
		})
	})

	Describe(`ListSnapshotSoftwareAttachments - List snapshot software attachments associated with a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSnapshotSoftwareAttachments(listSnapshotSoftwareAttachmentsOptions *ListSnapshotSoftwareAttachmentsOptions)`, func() {
			listSnapshotSoftwareAttachmentsOptions := &vpcv1.ListSnapshotSoftwareAttachmentsOptions{
				SnapshotID: core.StringPtr("testString"),
			}

			snapshotSoftwareAttachmentCollection, response, err := vpcService.ListSnapshotSoftwareAttachments(listSnapshotSoftwareAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachmentCollection).ToNot(BeNil())
		})
	})

	Describe(`GetSnapshotSoftwareAttachment - Retrieve a snapshot software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSnapshotSoftwareAttachment(getSnapshotSoftwareAttachmentOptions *GetSnapshotSoftwareAttachmentOptions)`, func() {
			getSnapshotSoftwareAttachmentOptions := &vpcv1.GetSnapshotSoftwareAttachmentOptions{
				ID:         core.StringPtr("testString"),
				SnapshotID: core.StringPtr("testString"),
			}

			snapshotSoftwareAttachment, response, err := vpcService.GetSnapshotSoftwareAttachment(getSnapshotSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateSnapshotSoftwareAttachment - Update a snapshot software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSnapshotSoftwareAttachment(updateSnapshotSoftwareAttachmentOptions *UpdateSnapshotSoftwareAttachmentOptions)`, func() {
			snapshotSoftwareAttachmentPatchModel := &vpcv1.SnapshotSoftwareAttachmentPatch{
				Name: core.StringPtr("my-software-attachment-patch"),
			}
			snapshotSoftwareAttachmentPatchModelAsPatch, asPatchErr := snapshotSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSnapshotSoftwareAttachmentOptions := &vpcv1.UpdateSnapshotSoftwareAttachmentOptions{
				ID:                              core.StringPtr("testString"),
				SnapshotID:                      core.StringPtr("testString"),
				SnapshotSoftwareAttachmentPatch: snapshotSoftwareAttachmentPatchModelAsPatch,
			}

			snapshotSoftwareAttachment, response, err := vpcService.UpdateSnapshotSoftwareAttachment(updateSnapshotSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(snapshotSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListSubnets - List subnets`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSubnets(listSubnetsOptions *ListSubnetsOptions) with pagination`, func() {
			listSubnetsOptions := &vpcv1.ListSubnetsOptions{
				Limit:            core.Int64Ptr(int64(10)),
				ResourceGroupID:  core.StringPtr("testString"),
				RoutingTableID:   core.StringPtr("testString"),
				RoutingTableName: core.StringPtr("my-routing-table"),
				Start:            core.StringPtr("testString"),
				VPCCRN:           core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCID:            core.StringPtr("testString"),
				VPCName:          core.StringPtr("my-vpc"),
				ZoneName:         core.StringPtr("us-south-1"),
			}

			listSubnetsOptions.Start = nil
			listSubnetsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Subnet
			for {
				subnetCollection, response, err := vpcService.ListSubnets(listSubnetsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(subnetCollection).ToNot(BeNil())
				allResults = append(allResults, subnetCollection.Subnets...)

				nextToken, nextErr := subnetCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSubnetsOptions.Start = nextToken

				if listSubnetsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSubnets(listSubnetsOptions *ListSubnetsOptions) using SubnetsPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewSubnetsPager(listSubnetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Subnet
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSubnetsPager(listSubnetsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSubnets() returned a total of %d item(s) using SubnetsPager.\n", len(allResults))
		})
	})

	Describe(`CreateSubnet - Create a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSubnet(createSubnetOptions *CreateSubnetOptions)`, func() {
			subnetPrototypeNetworkACLModel := &vpcv1.SubnetPrototypeNetworkACLNetworkACLIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::network-acl:r006-a4e28308-8ee7-46ab-8108-9f881f22bdbf"),
			}

			subnetPrototypePublicGatewayModel := &vpcv1.SubnetPrototypePublicGatewayPublicGatewayIdentityPublicGatewayIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::public-gateway:r006-dc5431ef-1fc6-4861-adc9-a59d077d1241"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			routingTableIdentityModel := &vpcv1.RoutingTableIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc-routing-table:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b/r006-6885e83f-03b2-4603-8a86-db2a0f55c840"),
			}

			subnetPrototypeVPCModel := &vpcv1.SubnetPrototypeVPCVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			subnetByCIDRZoneModel := &vpcv1.SubnetByCIDRZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			subnetPrototypeModel := &vpcv1.SubnetPrototypeSubnetByCIDR{
				IPVersion:     core.StringPtr("ipv4"),
				Name:          core.StringPtr("my-subnet"),
				NetworkACL:    subnetPrototypeNetworkACLModel,
				PublicGateway: subnetPrototypePublicGatewayModel,
				ResourceGroup: resourceGroupIdentityModel,
				RoutingTable:  routingTableIdentityModel,
				VPC:           subnetPrototypeVPCModel,
				Ipv4CIDRBlock: core.StringPtr("10.0.0.0/24"),
				Zone:          subnetByCIDRZoneModel,
			}

			createSubnetOptions := &vpcv1.CreateSubnetOptions{
				SubnetPrototype: subnetPrototypeModel,
			}

			subnet, response, err := vpcService.CreateSubnet(createSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(subnet).ToNot(BeNil())
		})
	})

	Describe(`GetSubnet - Retrieve a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSubnet(getSubnetOptions *GetSubnetOptions)`, func() {
			getSubnetOptions := &vpcv1.GetSubnetOptions{
				ID: core.StringPtr("testString"),
			}

			subnet, response, err := vpcService.GetSubnet(getSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(subnet).ToNot(BeNil())
		})
	})

	Describe(`UpdateSubnet - Update a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSubnet(updateSubnetOptions *UpdateSubnetOptions)`, func() {
			subnetPatchNetworkACLModel := &vpcv1.SubnetPatchNetworkACLNetworkACLIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::network-acl:r006-a4e28308-8ee7-46ab-8108-9f881f22bdbf"),
			}

			subnetPublicGatewayPatchModel := &vpcv1.SubnetPublicGatewayPatchPublicGatewayIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::public-gateway:r006-dc5431ef-1fc6-4861-adc9-a59d077d1241"),
			}

			routingTableIdentityModel := &vpcv1.RoutingTableIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc-routing-table:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b/r006-6885e83f-03b2-4603-8a86-db2a0f55c840"),
			}

			subnetPatchModel := &vpcv1.SubnetPatch{
				Name:          core.StringPtr("my-subnet"),
				NetworkACL:    subnetPatchNetworkACLModel,
				PublicGateway: subnetPublicGatewayPatchModel,
				RoutingTable:  routingTableIdentityModel,
			}
			subnetPatchModelAsPatch, asPatchErr := subnetPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSubnetOptions := &vpcv1.UpdateSubnetOptions{
				ID:          core.StringPtr("testString"),
				SubnetPatch: subnetPatchModelAsPatch,
			}

			subnet, response, err := vpcService.UpdateSubnet(updateSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(subnet).ToNot(BeNil())
		})
	})

	Describe(`GetSubnetNetworkACL - Retrieve a subnet's attached network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSubnetNetworkACL(getSubnetNetworkACLOptions *GetSubnetNetworkACLOptions)`, func() {
			getSubnetNetworkACLOptions := &vpcv1.GetSubnetNetworkACLOptions{
				ID: core.StringPtr("testString"),
			}

			networkACL, response, err := vpcService.GetSubnetNetworkACL(getSubnetNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(networkACL).ToNot(BeNil())
		})
	})

	Describe(`ReplaceSubnetNetworkACL - Replace the network ACL for a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptions *ReplaceSubnetNetworkACLOptions)`, func() {
			networkACLIdentityModel := &vpcv1.NetworkACLIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::network-acl:r006-a4e28308-8ee7-46ab-8108-9f881f22bdbf"),
			}

			replaceSubnetNetworkACLOptions := &vpcv1.ReplaceSubnetNetworkACLOptions{
				ID:                 core.StringPtr("testString"),
				NetworkACLIdentity: networkACLIdentityModel,
			}

			networkACL, response, err := vpcService.ReplaceSubnetNetworkACL(replaceSubnetNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(networkACL).ToNot(BeNil())
		})
	})

	Describe(`GetSubnetPublicGateway - Retrieve a subnet's attached public gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSubnetPublicGateway(getSubnetPublicGatewayOptions *GetSubnetPublicGatewayOptions)`, func() {
			getSubnetPublicGatewayOptions := &vpcv1.GetSubnetPublicGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			publicGateway, response, err := vpcService.GetSubnetPublicGateway(getSubnetPublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(publicGateway).ToNot(BeNil())
		})
	})

	Describe(`SetSubnetPublicGateway - Attach a public gateway to a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`SetSubnetPublicGateway(setSubnetPublicGatewayOptions *SetSubnetPublicGatewayOptions)`, func() {
			publicGatewayIdentityModel := &vpcv1.PublicGatewayIdentityPublicGatewayIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::public-gateway:r006-dc5431ef-1fc6-4861-adc9-a59d077d1241"),
			}

			setSubnetPublicGatewayOptions := &vpcv1.SetSubnetPublicGatewayOptions{
				ID:                    core.StringPtr("testString"),
				PublicGatewayIdentity: publicGatewayIdentityModel,
			}

			publicGateway, response, err := vpcService.SetSubnetPublicGateway(setSubnetPublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(publicGateway).ToNot(BeNil())
		})
	})

	Describe(`GetSubnetRoutingTable - Retrieve a subnet's attached routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSubnetRoutingTable(getSubnetRoutingTableOptions *GetSubnetRoutingTableOptions)`, func() {
			getSubnetRoutingTableOptions := &vpcv1.GetSubnetRoutingTableOptions{
				ID: core.StringPtr("testString"),
			}

			routingTable, response, err := vpcService.GetSubnetRoutingTable(getSubnetRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
	})

	Describe(`ReplaceSubnetRoutingTable - Replace the routing table for a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptions *ReplaceSubnetRoutingTableOptions)`, func() {
			routingTableIdentityModel := &vpcv1.RoutingTableIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc-routing-table:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b/r006-6885e83f-03b2-4603-8a86-db2a0f55c840"),
			}

			replaceSubnetRoutingTableOptions := &vpcv1.ReplaceSubnetRoutingTableOptions{
				ID:                   core.StringPtr("testString"),
				RoutingTableIdentity: routingTableIdentityModel,
			}

			routingTable, response, err := vpcService.ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(routingTable).ToNot(BeNil())
		})
	})

	Describe(`ListSubnetReservedIps - List reserved IPs in a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListSubnetReservedIps(listSubnetReservedIpsOptions *ListSubnetReservedIpsOptions) with pagination`, func() {
			listSubnetReservedIpsOptions := &vpcv1.ListSubnetReservedIpsOptions{
				SubnetID:           core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Sort:               core.StringPtr("name"),
				Start:              core.StringPtr("testString"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			listSubnetReservedIpsOptions.Start = nil
			listSubnetReservedIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ReservedIP
			for {
				reservedIPCollection, response, err := vpcService.ListSubnetReservedIps(listSubnetReservedIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(reservedIPCollection).ToNot(BeNil())
				allResults = append(allResults, reservedIPCollection.ReservedIps...)

				nextToken, nextErr := reservedIPCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listSubnetReservedIpsOptions.Start = nextToken

				if listSubnetReservedIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListSubnetReservedIps(listSubnetReservedIpsOptions *ListSubnetReservedIpsOptions) using SubnetReservedIpsPager`, func() {
			listSubnetReservedIpsOptions := &vpcv1.ListSubnetReservedIpsOptions{
				SubnetID:           core.StringPtr("testString"),
				Limit:              core.Int64Ptr(int64(10)),
				Sort:               core.StringPtr("name"),
				TargetCRN:          core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::load-balancer:r006-dd754295-e9e0-4c9d-bf6c-58fbc59e5727"),
				TargetID:           core.StringPtr("testString"),
				TargetName:         core.StringPtr("my-resource"),
				TargetResourceType: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewSubnetReservedIpsPager(listSubnetReservedIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ReservedIP
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewSubnetReservedIpsPager(listSubnetReservedIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListSubnetReservedIps() returned a total of %d item(s) using SubnetReservedIpsPager.\n", len(allResults))
		})
	})

	Describe(`CreateSubnetReservedIP - Reserve an IP in a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateSubnetReservedIP(createSubnetReservedIPOptions *CreateSubnetReservedIPOptions)`, func() {
			reservedIPTargetPrototypeModel := &vpcv1.ReservedIPTargetPrototypeEndpointGatewayIdentityEndpointGatewayIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::endpoint-gateway:r006-7610ebfb-f5dc-4d42-bc18-287d47f7a5b0"),
			}

			createSubnetReservedIPOptions := &vpcv1.CreateSubnetReservedIPOptions{
				SubnetID:   core.StringPtr("testString"),
				Address:    core.StringPtr("192.168.3.4"),
				AutoDelete: core.BoolPtr(false),
				Name:       core.StringPtr("my-reserved-ip"),
				Target:     reservedIPTargetPrototypeModel,
			}

			reservedIP, response, err := vpcService.CreateSubnetReservedIP(createSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`GetSubnetReservedIP - Retrieve a reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSubnetReservedIP(getSubnetReservedIPOptions *GetSubnetReservedIPOptions)`, func() {
			getSubnetReservedIPOptions := &vpcv1.GetSubnetReservedIPOptions{
				ID:       core.StringPtr("testString"),
				SubnetID: core.StringPtr("testString"),
			}

			reservedIP, response, err := vpcService.GetSubnetReservedIP(getSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`UpdateSubnetReservedIP - Update a reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateSubnetReservedIP(updateSubnetReservedIPOptions *UpdateSubnetReservedIPOptions)`, func() {
			reservedIPPatchModel := &vpcv1.ReservedIPPatch{
				AutoDelete: core.BoolPtr(false),
				Name:       core.StringPtr("my-reserved-ip-updated"),
			}
			reservedIPPatchModelAsPatch, asPatchErr := reservedIPPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateSubnetReservedIPOptions := &vpcv1.UpdateSubnetReservedIPOptions{
				ID:              core.StringPtr("testString"),
				SubnetID:        core.StringPtr("testString"),
				ReservedIPPatch: reservedIPPatchModelAsPatch,
			}

			reservedIP, response, err := vpcService.UpdateSubnetReservedIP(updateSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIP).ToNot(BeNil())
		})
	})

	Describe(`ListVirtualNetworkInterfaces - List virtual network interfaces`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVirtualNetworkInterfaces(listVirtualNetworkInterfacesOptions *ListVirtualNetworkInterfacesOptions) with pagination`, func() {
			listVirtualNetworkInterfacesOptions := &vpcv1.ListVirtualNetworkInterfacesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listVirtualNetworkInterfacesOptions.Start = nil
			listVirtualNetworkInterfacesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VirtualNetworkInterface
			for {
				virtualNetworkInterfaceCollection, response, err := vpcService.ListVirtualNetworkInterfaces(listVirtualNetworkInterfacesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(virtualNetworkInterfaceCollection).ToNot(BeNil())
				allResults = append(allResults, virtualNetworkInterfaceCollection.VirtualNetworkInterfaces...)

				nextToken, nextErr := virtualNetworkInterfaceCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVirtualNetworkInterfacesOptions.Start = nextToken

				if listVirtualNetworkInterfacesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVirtualNetworkInterfaces(listVirtualNetworkInterfacesOptions *ListVirtualNetworkInterfacesOptions) using VirtualNetworkInterfacesPager`, func() {
			listVirtualNetworkInterfacesOptions := &vpcv1.ListVirtualNetworkInterfacesOptions{
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVirtualNetworkInterfacesPager(listVirtualNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VirtualNetworkInterface
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVirtualNetworkInterfacesPager(listVirtualNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVirtualNetworkInterfaces() returned a total of %d item(s) using VirtualNetworkInterfacesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVirtualNetworkInterface - Create a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVirtualNetworkInterface(createVirtualNetworkInterfaceOptions *CreateVirtualNetworkInterfaceOptions)`, func() {
			virtualNetworkInterfaceIPPrototypeModel := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			virtualNetworkInterfacePrimaryIPPrototypeModel := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-bea6a632-5e13-42a4-b4b8-31dc877abfe4/reserved_ips/0717-6d353a0f-aeb1-4ae1-832e-1110d10981bb"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			virtualNetworkInterfacePrototypeSubnetModel := &vpcv1.VirtualNetworkInterfacePrototypeSubnetSubnetIdentityByID{
				ID: core.StringPtr("69e55145-cc7d-4d8e-9e1f-cc3fb60b1793"),
			}

			createVirtualNetworkInterfaceOptions := &vpcv1.CreateVirtualNetworkInterfaceOptions{
				AllowIPSpoofing:            core.BoolPtr(false),
				AutoDelete:                 core.BoolPtr(false),
				EnableInfrastructureNat:    core.BoolPtr(true),
				Ips:                        []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{virtualNetworkInterfaceIPPrototypeModel},
				Name:                       core.StringPtr("my-virtual-network-interface"),
				PrimaryIP:                  virtualNetworkInterfacePrimaryIPPrototypeModel,
				ProtocolStateFilteringMode: core.StringPtr("auto"),
				ResourceGroup:              resourceGroupIdentityModel,
				SecurityGroups:             []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
				Subnet:                     virtualNetworkInterfacePrototypeSubnetModel,
			}

			virtualNetworkInterface, response, err := vpcService.CreateVirtualNetworkInterface(createVirtualNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`GetVirtualNetworkInterface - Retrieve a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions *GetVirtualNetworkInterfaceOptions)`, func() {
			getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
				ID: core.StringPtr("testString"),
			}

			virtualNetworkInterface, response, err := vpcService.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`UpdateVirtualNetworkInterface - Update a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions *UpdateVirtualNetworkInterfaceOptions)`, func() {
			virtualNetworkInterfacePatchModel := &vpcv1.VirtualNetworkInterfacePatch{
				AllowIPSpoofing:            core.BoolPtr(true),
				AutoDelete:                 core.BoolPtr(false),
				EnableInfrastructureNat:    core.BoolPtr(true),
				Name:                       core.StringPtr("my-virtual-network-interface"),
				ProtocolStateFilteringMode: core.StringPtr("auto"),
			}
			virtualNetworkInterfacePatchModelAsPatch, asPatchErr := virtualNetworkInterfacePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
				ID:                           core.StringPtr("testString"),
				VirtualNetworkInterfacePatch: virtualNetworkInterfacePatchModelAsPatch,
				IfMatch:                      core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			virtualNetworkInterface, response, err := vpcService.UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`ListNetworkInterfaceFloatingIps - List floating IPs associated with a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListNetworkInterfaceFloatingIps(listNetworkInterfaceFloatingIpsOptions *ListNetworkInterfaceFloatingIpsOptions) with pagination`, func() {
			listNetworkInterfaceFloatingIpsOptions := &vpcv1.ListNetworkInterfaceFloatingIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
				Start:                     core.StringPtr("testString"),
			}

			listNetworkInterfaceFloatingIpsOptions.Start = nil
			listNetworkInterfaceFloatingIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.FloatingIPReference
			for {
				floatingIPCollectionVirtualNetworkInterfaceContext, response, err := vpcService.ListNetworkInterfaceFloatingIps(listNetworkInterfaceFloatingIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(floatingIPCollectionVirtualNetworkInterfaceContext).ToNot(BeNil())
				allResults = append(allResults, floatingIPCollectionVirtualNetworkInterfaceContext.FloatingIps...)

				nextToken, nextErr := floatingIPCollectionVirtualNetworkInterfaceContext.GetNextStart()
				Expect(nextErr).To(BeNil())
				listNetworkInterfaceFloatingIpsOptions.Start = nextToken

				if listNetworkInterfaceFloatingIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListNetworkInterfaceFloatingIps(listNetworkInterfaceFloatingIpsOptions *ListNetworkInterfaceFloatingIpsOptions) using NetworkInterfaceFloatingIpsPager`, func() {
			listNetworkInterfaceFloatingIpsOptions := &vpcv1.ListNetworkInterfaceFloatingIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewNetworkInterfaceFloatingIpsPager(listNetworkInterfaceFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.FloatingIPReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewNetworkInterfaceFloatingIpsPager(listNetworkInterfaceFloatingIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListNetworkInterfaceFloatingIps() returned a total of %d item(s) using NetworkInterfaceFloatingIpsPager.\n", len(allResults))
		})
	})

	Describe(`GetNetworkInterfaceFloatingIP - Retrieve associated floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetNetworkInterfaceFloatingIP(getNetworkInterfaceFloatingIPOptions *GetNetworkInterfaceFloatingIPOptions)`, func() {
			getNetworkInterfaceFloatingIPOptions := &vpcv1.GetNetworkInterfaceFloatingIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIPReference, response, err := vpcService.GetNetworkInterfaceFloatingIP(getNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPReference).ToNot(BeNil())
		})
	})

	Describe(`AddNetworkInterfaceFloatingIP - Add an association between a floating IP and a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddNetworkInterfaceFloatingIP(addNetworkInterfaceFloatingIPOptions *AddNetworkInterfaceFloatingIPOptions)`, func() {
			addNetworkInterfaceFloatingIPOptions := &vpcv1.AddNetworkInterfaceFloatingIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			floatingIPReference, response, err := vpcService.AddNetworkInterfaceFloatingIP(addNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(floatingIPReference).ToNot(BeNil())
		})
	})

	Describe(`ListVirtualNetworkInterfaceIps - List reserved IPs bound to a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVirtualNetworkInterfaceIps(listVirtualNetworkInterfaceIpsOptions *ListVirtualNetworkInterfaceIpsOptions) with pagination`, func() {
			listVirtualNetworkInterfaceIpsOptions := &vpcv1.ListVirtualNetworkInterfaceIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
				Start:                     core.StringPtr("testString"),
			}

			listVirtualNetworkInterfaceIpsOptions.Start = nil
			listVirtualNetworkInterfaceIpsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.ReservedIPReference
			for {
				reservedIPCollectionVirtualNetworkInterfaceContext, response, err := vpcService.ListVirtualNetworkInterfaceIps(listVirtualNetworkInterfaceIpsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(reservedIPCollectionVirtualNetworkInterfaceContext).ToNot(BeNil())
				allResults = append(allResults, reservedIPCollectionVirtualNetworkInterfaceContext.Ips...)

				nextToken, nextErr := reservedIPCollectionVirtualNetworkInterfaceContext.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVirtualNetworkInterfaceIpsOptions.Start = nextToken

				if listVirtualNetworkInterfaceIpsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVirtualNetworkInterfaceIps(listVirtualNetworkInterfaceIpsOptions *ListVirtualNetworkInterfaceIpsOptions) using VirtualNetworkInterfaceIpsPager`, func() {
			listVirtualNetworkInterfaceIpsOptions := &vpcv1.ListVirtualNetworkInterfaceIpsOptions{
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
				Limit:                     core.Int64Ptr(int64(10)),
				Sort:                      core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVirtualNetworkInterfaceIpsPager(listVirtualNetworkInterfaceIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.ReservedIPReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVirtualNetworkInterfaceIpsPager(listVirtualNetworkInterfaceIpsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVirtualNetworkInterfaceIps() returned a total of %d item(s) using VirtualNetworkInterfaceIpsPager.\n", len(allResults))
		})
	})

	Describe(`GetVirtualNetworkInterfaceIP - Retrieve bound reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVirtualNetworkInterfaceIP(getVirtualNetworkInterfaceIPOptions *GetVirtualNetworkInterfaceIPOptions)`, func() {
			getVirtualNetworkInterfaceIPOptions := &vpcv1.GetVirtualNetworkInterfaceIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			reservedIPReference, response, err := vpcService.GetVirtualNetworkInterfaceIP(getVirtualNetworkInterfaceIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPReference).ToNot(BeNil())
		})
	})

	Describe(`AddVirtualNetworkInterfaceIP - Bind a reserved IP to a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions *AddVirtualNetworkInterfaceIPOptions)`, func() {
			addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			reservedIPReference, response, err := vpcService.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(reservedIPReference).ToNot(BeNil())
		})
	})

	Describe(`ListVolumeProfiles - List volume profiles`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVolumeProfiles(listVolumeProfilesOptions *ListVolumeProfilesOptions) with pagination`, func() {
			listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listVolumeProfilesOptions.Start = nil
			listVolumeProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VolumeProfile
			for {
				volumeProfileCollection, response, err := vpcService.ListVolumeProfiles(listVolumeProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(volumeProfileCollection).ToNot(BeNil())
				allResults = append(allResults, volumeProfileCollection.Profiles...)

				nextToken, nextErr := volumeProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVolumeProfilesOptions.Start = nextToken

				if listVolumeProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVolumeProfiles(listVolumeProfilesOptions *ListVolumeProfilesOptions) using VolumeProfilesPager`, func() {
			listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVolumeProfilesPager(listVolumeProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VolumeProfile
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVolumeProfilesPager(listVolumeProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVolumeProfiles() returned a total of %d item(s) using VolumeProfilesPager.\n", len(allResults))
		})
	})

	Describe(`GetVolumeProfile - Retrieve a volume profile`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVolumeProfile(getVolumeProfileOptions *GetVolumeProfileOptions)`, func() {
			getVolumeProfileOptions := &vpcv1.GetVolumeProfileOptions{
				Name: core.StringPtr("10iops-tier"),
			}

			volumeProfile, response, err := vpcService.GetVolumeProfile(getVolumeProfileOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeProfile).ToNot(BeNil())
		})
	})

	Describe(`ListVolumes - List volumes`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVolumes(listVolumesOptions *ListVolumesOptions) with pagination`, func() {
			listVolumesOptions := &vpcv1.ListVolumesOptions{
				AttachmentState:             core.StringPtr("attached"),
				Encryption:                  core.StringPtr("provider_managed"),
				Limit:                       core.Int64Ptr(int64(10)),
				Name:                        core.StringPtr("my-name"),
				OperatingSystemArchitecture: core.StringPtr("amd64"),
				OperatingSystemFamily:       core.StringPtr("Ubuntu Server"),
				Start:                       core.StringPtr("testString"),
				StorageGeneration:           core.Int64Ptr(int64(1)),
				Tag:                         core.StringPtr("testString"),
				ZoneName:                    core.StringPtr("us-south-1"),
			}

			listVolumesOptions.Start = nil
			listVolumesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Volume
			for {
				volumeCollection, response, err := vpcService.ListVolumes(listVolumesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(volumeCollection).ToNot(BeNil())
				allResults = append(allResults, volumeCollection.Volumes...)

				nextToken, nextErr := volumeCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVolumesOptions.Start = nextToken

				if listVolumesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVolumes(listVolumesOptions *ListVolumesOptions) using VolumesPager`, func() {
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

			// Test GetNext().
			pager, err := vpcService.NewVolumesPager(listVolumesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Volume
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVolumesPager(listVolumesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVolumes() returned a total of %d item(s) using VolumesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVolume - Create a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVolume(createVolumeOptions *CreateVolumeOptions)`, func() {
			volumePrototypeProfileModel := &vpcv1.VolumePrototypeProfileVolumeProfileIdentityByName{
				Name: core.StringPtr("custom"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			volumePrototypeZoneModel := &vpcv1.VolumePrototypeZoneZoneIdentityByName{
				Name: core.StringPtr("us-south-1"),
			}

			volumeByCapacityEncryptionKeyModel := &vpcv1.VolumeByCapacityEncryptionKeyEncryptionKeyIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"),
			}

			volumePrototypeModel := &vpcv1.VolumePrototypeVolumeByCapacity{
				Bandwidth:     core.Int64Ptr(int64(1000)),
				Iops:          core.Int64Ptr(int64(10000)),
				Name:          core.StringPtr("my-volume"),
				Profile:       volumePrototypeProfileModel,
				ResourceGroup: resourceGroupIdentityModel,
				UserTags:      []string{"testString"},
				Zone:          volumePrototypeZoneModel,
				Capacity:      core.Int64Ptr(int64(100)),
				EncryptionKey: volumeByCapacityEncryptionKeyModel,
			}

			createVolumeOptions := &vpcv1.CreateVolumeOptions{
				VolumePrototype: volumePrototypeModel,
			}

			volume, response, err := vpcService.CreateVolume(createVolumeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volume).ToNot(BeNil())
		})
	})

	Describe(`GetVolume - Retrieve a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVolume(getVolumeOptions *GetVolumeOptions)`, func() {
			getVolumeOptions := &vpcv1.GetVolumeOptions{
				ID: core.StringPtr("testString"),
			}

			volume, response, err := vpcService.GetVolume(getVolumeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volume).ToNot(BeNil())
		})
	})

	Describe(`UpdateVolume - Update a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVolume(updateVolumeOptions *UpdateVolumeOptions)`, func() {
			volumeAllowedUsePatchModel := &vpcv1.VolumeAllowedUsePatch{
				ApiVersion:      core.StringPtr("2024-06-23"),
				BareMetalServer: core.StringPtr("enable_secure_boot == true"),
				Instance:        core.StringPtr("gpu.count > 0 && enable_secure_boot == true"),
			}

			volumePatchProfileModel := &vpcv1.VolumePatchProfileVolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumePatchModel := &vpcv1.VolumePatch{
				AllowedUse: volumeAllowedUsePatchModel,
				Bandwidth:  core.Int64Ptr(int64(1000)),
				Capacity:   core.Int64Ptr(int64(100)),
				Iops:       core.Int64Ptr(int64(10000)),
				Name:       core.StringPtr("my-volume"),
				Profile:    volumePatchProfileModel,
				UserTags:   []string{"testString"},
			}
			volumePatchModelAsPatch, asPatchErr := volumePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
				ID:          core.StringPtr("testString"),
				VolumePatch: volumePatchModelAsPatch,
				IfMatch:     core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			volume, response, err := vpcService.UpdateVolume(updateVolumeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volume).ToNot(BeNil())
		})
	})

	Describe(`ListVolumeInstanceProfiles - List instance profiles compatible with a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVolumeInstanceProfiles(listVolumeInstanceProfilesOptions *ListVolumeInstanceProfilesOptions) with pagination`, func() {
			listVolumeInstanceProfilesOptions := &vpcv1.ListVolumeInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listVolumeInstanceProfilesOptions.Start = nil
			listVolumeInstanceProfilesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.InstanceProfileReference
			for {
				volumeInstanceProfileCollection, response, err := vpcService.ListVolumeInstanceProfiles(listVolumeInstanceProfilesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(volumeInstanceProfileCollection).ToNot(BeNil())
				allResults = append(allResults, volumeInstanceProfileCollection.InstanceProfiles...)

				nextToken, nextErr := volumeInstanceProfileCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVolumeInstanceProfilesOptions.Start = nextToken

				if listVolumeInstanceProfilesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVolumeInstanceProfiles(listVolumeInstanceProfilesOptions *ListVolumeInstanceProfilesOptions) using VolumeInstanceProfilesPager`, func() {
			listVolumeInstanceProfilesOptions := &vpcv1.ListVolumeInstanceProfilesOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVolumeInstanceProfilesPager(listVolumeInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.InstanceProfileReference
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVolumeInstanceProfilesPager(listVolumeInstanceProfilesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVolumeInstanceProfiles() returned a total of %d item(s) using VolumeInstanceProfilesPager.\n", len(allResults))
		})
	})

	Describe(`ListVolumeJobs - List jobs for a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVolumeJobs(listVolumeJobsOptions *ListVolumeJobsOptions) with pagination`, func() {
			listVolumeJobsOptions := &vpcv1.ListVolumeJobsOptions{
				VolumeID: core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
				Start:    core.StringPtr("testString"),
			}

			listVolumeJobsOptions.Start = nil
			listVolumeJobsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VolumeJobIntf
			for {
				volumeJobCollection, response, err := vpcService.ListVolumeJobs(listVolumeJobsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(volumeJobCollection).ToNot(BeNil())
				allResults = append(allResults, volumeJobCollection.Jobs...)

				nextToken, nextErr := volumeJobCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVolumeJobsOptions.Start = nextToken

				if listVolumeJobsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVolumeJobs(listVolumeJobsOptions *ListVolumeJobsOptions) using VolumeJobsPager`, func() {
			listVolumeJobsOptions := &vpcv1.ListVolumeJobsOptions{
				VolumeID: core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVolumeJobsPager(listVolumeJobsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VolumeJobIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVolumeJobsPager(listVolumeJobsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVolumeJobs() returned a total of %d item(s) using VolumeJobsPager.\n", len(allResults))
		})
	})

	Describe(`CreateVolumeJob - Create a job for a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVolumeJob(createVolumeJobOptions *CreateVolumeJobOptions)`, func() {
			volumeProfileIdentityModel := &vpcv1.VolumeProfileIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"),
			}

			volumeJobTypeMigrateParametersModel := &vpcv1.VolumeJobTypeMigrateParameters{
				Bandwidth: core.Int64Ptr(int64(1000)),
				Iops:      core.Int64Ptr(int64(10000)),
				Profile:   volumeProfileIdentityModel,
			}

			volumeJobPrototypeModel := &vpcv1.VolumeJobPrototypeVolumeJobTypeMigratePrototype{
				Name:       core.StringPtr("my-volume-job"),
				JobType:    core.StringPtr("migrate"),
				Parameters: volumeJobTypeMigrateParametersModel,
			}

			createVolumeJobOptions := &vpcv1.CreateVolumeJobOptions{
				VolumeID:           core.StringPtr("testString"),
				VolumeJobPrototype: volumeJobPrototypeModel,
				Limit:              core.Int64Ptr(int64(50)),
				Start:              core.StringPtr("testString"),
			}

			volumeJob, response, err := vpcService.CreateVolumeJob(createVolumeJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(volumeJob).ToNot(BeNil())
		})
	})

	Describe(`GetVolumeJob - Retrieve a volume job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVolumeJob(getVolumeJobOptions *GetVolumeJobOptions)`, func() {
			getVolumeJobOptions := &vpcv1.GetVolumeJobOptions{
				ID:       core.StringPtr("testString"),
				VolumeID: core.StringPtr("testString"),
			}

			volumeJob, response, err := vpcService.GetVolumeJob(getVolumeJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeJob).ToNot(BeNil())
		})
	})

	Describe(`UpdateVolumeJob - Update a volume job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVolumeJob(updateVolumeJobOptions *UpdateVolumeJobOptions)`, func() {
			volumeJobPatchModel := &vpcv1.VolumeJobPatch{
				Name: core.StringPtr("my-updated-volume-job-name"),
			}
			volumeJobPatchModelAsPatch, asPatchErr := volumeJobPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeJobOptions := &vpcv1.UpdateVolumeJobOptions{
				ID:             core.StringPtr("testString"),
				VolumeID:       core.StringPtr("testString"),
				VolumeJobPatch: volumeJobPatchModelAsPatch,
			}

			volumeJob, response, err := vpcService.UpdateVolumeJob(updateVolumeJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeJob).ToNot(BeNil())
		})
	})

	Describe(`CancelVolumeJob - Cancel a queued or running volume job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CancelVolumeJob(cancelVolumeJobOptions *CancelVolumeJobOptions)`, func() {
			cancelVolumeJobOptions := &vpcv1.CancelVolumeJobOptions{
				ID:       core.StringPtr("testString"),
				VolumeID: core.StringPtr("testString"),
			}

			volumeJob, response, err := vpcService.CancelVolumeJob(cancelVolumeJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(volumeJob).ToNot(BeNil())
		})
	})

	Describe(`ListVolumeSoftwareAttachments - List volume software attachments associated with a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVolumeSoftwareAttachments(listVolumeSoftwareAttachmentsOptions *ListVolumeSoftwareAttachmentsOptions)`, func() {
			listVolumeSoftwareAttachmentsOptions := &vpcv1.ListVolumeSoftwareAttachmentsOptions{
				VolumeID: core.StringPtr("testString"),
			}

			volumeSoftwareAttachmentCollection, response, err := vpcService.ListVolumeSoftwareAttachments(listVolumeSoftwareAttachmentsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachmentCollection).ToNot(BeNil())
		})
	})

	Describe(`GetVolumeSoftwareAttachment - Retrieve a volume software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVolumeSoftwareAttachment(getVolumeSoftwareAttachmentOptions *GetVolumeSoftwareAttachmentOptions)`, func() {
			getVolumeSoftwareAttachmentOptions := &vpcv1.GetVolumeSoftwareAttachmentOptions{
				ID:       core.StringPtr("testString"),
				VolumeID: core.StringPtr("testString"),
			}

			volumeSoftwareAttachment, response, err := vpcService.GetVolumeSoftwareAttachment(getVolumeSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`UpdateVolumeSoftwareAttachment - Update a volume software attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVolumeSoftwareAttachment(updateVolumeSoftwareAttachmentOptions *UpdateVolumeSoftwareAttachmentOptions)`, func() {
			volumeSoftwareAttachmentPatchModel := &vpcv1.VolumeSoftwareAttachmentPatch{
				Name: core.StringPtr("my-software-attachment-patch"),
			}
			volumeSoftwareAttachmentPatchModelAsPatch, asPatchErr := volumeSoftwareAttachmentPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVolumeSoftwareAttachmentOptions := &vpcv1.UpdateVolumeSoftwareAttachmentOptions{
				ID:                            core.StringPtr("testString"),
				VolumeID:                      core.StringPtr("testString"),
				VolumeSoftwareAttachmentPatch: volumeSoftwareAttachmentPatchModelAsPatch,
			}

			volumeSoftwareAttachment, response, err := vpcService.UpdateVolumeSoftwareAttachment(updateVolumeSoftwareAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(volumeSoftwareAttachment).ToNot(BeNil())
		})
	})

	Describe(`ListVpcs - List VPCs`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVpcs(listVpcsOptions *ListVpcsOptions) with pagination`, func() {
			listVpcsOptions := &vpcv1.ListVpcsOptions{
				ClassicAccess:   core.BoolPtr(true),
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
				Start:           core.StringPtr("testString"),
			}

			listVpcsOptions.Start = nil
			listVpcsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPC
			for {
				vpcCollection, response, err := vpcService.ListVpcs(listVpcsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpcCollection).ToNot(BeNil())
				allResults = append(allResults, vpcCollection.Vpcs...)

				nextToken, nextErr := vpcCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVpcsOptions.Start = nextToken

				if listVpcsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVpcs(listVpcsOptions *ListVpcsOptions) using VpcsPager`, func() {
			listVpcsOptions := &vpcv1.ListVpcsOptions{
				ClassicAccess:   core.BoolPtr(true),
				Limit:           core.Int64Ptr(int64(10)),
				ResourceGroupID: core.StringPtr("testString"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVpcsPager(listVpcsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPC
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVpcsPager(listVpcsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVpcs() returned a total of %d item(s) using VpcsPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPC - Create a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPC(createVPCOptions *CreateVPCOptions)`, func() {
			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			dnsServerPrototypeModel := &vpcv1.DnsServerPrototype{
				Address:      core.StringPtr("192.168.3.4"),
				ZoneAffinity: zoneIdentityModel,
			}

			vpcdnsResolverPrototypeModel := &vpcv1.VpcdnsResolverPrototypeVpcdnsResolverTypeManualPrototype{
				ManualServers: []vpcv1.DnsServerPrototype{*dnsServerPrototypeModel},
				Type:          core.StringPtr("manual"),
			}

			vpcdnsPrototypeModel := &vpcv1.VpcdnsPrototype{
				EnableHub: core.BoolPtr(false),
				Resolver:  vpcdnsResolverPrototypeModel,
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createVPCOptions := &vpcv1.CreateVPCOptions{
				AddressPrefixManagement: core.StringPtr("auto"),
				ClassicAccess:           core.BoolPtr(false),
				Dns:                     vpcdnsPrototypeModel,
				Name:                    core.StringPtr("my-vpc"),
				ResourceGroup:           resourceGroupIdentityModel,
			}

			vpc, response, err := vpcService.CreateVPC(createVPCOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpc).ToNot(BeNil())
		})
	})

	Describe(`GetVPC - Retrieve a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPC(getVPCOptions *GetVPCOptions)`, func() {
			getVPCOptions := &vpcv1.GetVPCOptions{
				ID: core.StringPtr("testString"),
			}

			vpc, response, err := vpcService.GetVPC(getVPCOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpc).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPC - Update a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPC(updateVPCOptions *UpdateVPCOptions)`, func() {
			zoneIdentityModel := &vpcv1.ZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			dnsServerPrototypeModel := &vpcv1.DnsServerPrototype{
				Address:      core.StringPtr("192.168.3.4"),
				ZoneAffinity: zoneIdentityModel,
			}

			vpcdnsResolverVPCPatchModel := &vpcv1.VpcdnsResolverVPCPatchVPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			vpcdnsResolverPatchModel := &vpcv1.VpcdnsResolverPatch{
				ManualServers: []vpcv1.DnsServerPrototype{*dnsServerPrototypeModel},
				Type:          core.StringPtr("delegated"),
				VPC:           vpcdnsResolverVPCPatchModel,
			}

			vpcdnsPatchModel := &vpcv1.VpcdnsPatch{
				EnableHub: core.BoolPtr(true),
				Resolver:  vpcdnsResolverPatchModel,
			}

			vpcPatchModel := &vpcv1.VPCPatch{
				Dns:  vpcdnsPatchModel,
				Name: core.StringPtr("my-vpc"),
			}
			vpcPatchModelAsPatch, asPatchErr := vpcPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCOptions := &vpcv1.UpdateVPCOptions{
				ID:       core.StringPtr("testString"),
				VPCPatch: vpcPatchModelAsPatch,
				IfMatch:  core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			vpc, response, err := vpcService.UpdateVPC(updateVPCOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpc).ToNot(BeNil())
		})
	})

	Describe(`GetVPCDefaultNetworkACL - Retrieve a VPC's default network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCDefaultNetworkACL(getVPCDefaultNetworkACLOptions *GetVPCDefaultNetworkACLOptions)`, func() {
			getVPCDefaultNetworkACLOptions := &vpcv1.GetVPCDefaultNetworkACLOptions{
				ID: core.StringPtr("testString"),
			}

			defaultNetworkACL, response, err := vpcService.GetVPCDefaultNetworkACL(getVPCDefaultNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultNetworkACL).ToNot(BeNil())
		})
	})

	Describe(`GetVPCDefaultRoutingTable - Retrieve a VPC's default routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCDefaultRoutingTable(getVPCDefaultRoutingTableOptions *GetVPCDefaultRoutingTableOptions)`, func() {
			getVPCDefaultRoutingTableOptions := &vpcv1.GetVPCDefaultRoutingTableOptions{
				ID: core.StringPtr("testString"),
			}

			defaultRoutingTable, response, err := vpcService.GetVPCDefaultRoutingTable(getVPCDefaultRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultRoutingTable).ToNot(BeNil())
		})
	})

	Describe(`GetVPCDefaultSecurityGroup - Retrieve a VPC's default security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCDefaultSecurityGroup(getVPCDefaultSecurityGroupOptions *GetVPCDefaultSecurityGroupOptions)`, func() {
			getVPCDefaultSecurityGroupOptions := &vpcv1.GetVPCDefaultSecurityGroupOptions{
				ID: core.StringPtr("testString"),
			}

			defaultSecurityGroup, response, err := vpcService.GetVPCDefaultSecurityGroup(getVPCDefaultSecurityGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(defaultSecurityGroup).ToNot(BeNil())
		})
	})

	Describe(`ListVPCAddressPrefixes - List address prefixes for a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPCAddressPrefixes(listVPCAddressPrefixesOptions *ListVPCAddressPrefixesOptions) with pagination`, func() {
			listVPCAddressPrefixesOptions := &vpcv1.ListVPCAddressPrefixesOptions{
				VPCID: core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listVPCAddressPrefixesOptions.Start = nil
			listVPCAddressPrefixesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.AddressPrefix
			for {
				addressPrefixCollection, response, err := vpcService.ListVPCAddressPrefixes(listVPCAddressPrefixesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(addressPrefixCollection).ToNot(BeNil())
				allResults = append(allResults, addressPrefixCollection.AddressPrefixes...)

				nextToken, nextErr := addressPrefixCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPCAddressPrefixesOptions.Start = nextToken

				if listVPCAddressPrefixesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPCAddressPrefixes(listVPCAddressPrefixesOptions *ListVPCAddressPrefixesOptions) using VPCAddressPrefixesPager`, func() {
			listVPCAddressPrefixesOptions := &vpcv1.ListVPCAddressPrefixesOptions{
				VPCID: core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPCAddressPrefixesPager(listVPCAddressPrefixesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.AddressPrefix
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPCAddressPrefixesPager(listVPCAddressPrefixesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPCAddressPrefixes() returned a total of %d item(s) using VPCAddressPrefixesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPCAddressPrefix - Create an address prefix for a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPCAddressPrefix(createVPCAddressPrefixOptions *CreateVPCAddressPrefixOptions)`, func() {
			addressPrefixPrototypeZoneModel := &vpcv1.AddressPrefixPrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			createVPCAddressPrefixOptions := &vpcv1.CreateVPCAddressPrefixOptions{
				VPCID:     core.StringPtr("testString"),
				CIDR:      core.StringPtr("10.240.0.0/18"),
				Zone:      addressPrefixPrototypeZoneModel,
				IsDefault: core.BoolPtr(true),
				Name:      core.StringPtr("my-vpc-address-prefix-zone-1"),
			}

			addressPrefix, response, err := vpcService.CreateVPCAddressPrefix(createVPCAddressPrefixOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(addressPrefix).ToNot(BeNil())
		})
	})

	Describe(`GetVPCAddressPrefix - Retrieve an address prefix`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCAddressPrefix(getVPCAddressPrefixOptions *GetVPCAddressPrefixOptions)`, func() {
			getVPCAddressPrefixOptions := &vpcv1.GetVPCAddressPrefixOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			addressPrefix, response, err := vpcService.GetVPCAddressPrefix(getVPCAddressPrefixOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(addressPrefix).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPCAddressPrefix - Update an address prefix`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPCAddressPrefix(updateVPCAddressPrefixOptions *UpdateVPCAddressPrefixOptions)`, func() {
			addressPrefixPatchModel := &vpcv1.AddressPrefixPatch{
				IsDefault: core.BoolPtr(false),
				Name:      core.StringPtr("my-vpc-address-prefix-updated"),
			}
			addressPrefixPatchModelAsPatch, asPatchErr := addressPrefixPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCAddressPrefixOptions := &vpcv1.UpdateVPCAddressPrefixOptions{
				ID:                 core.StringPtr("testString"),
				VPCID:              core.StringPtr("testString"),
				AddressPrefixPatch: addressPrefixPatchModelAsPatch,
			}

			addressPrefix, response, err := vpcService.UpdateVPCAddressPrefix(updateVPCAddressPrefixOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(addressPrefix).ToNot(BeNil())
		})
	})

	Describe(`ListVPCDnsResolutionBindings - List DNS resolution bindings for a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPCDnsResolutionBindings(listVPCDnsResolutionBindingsOptions *ListVPCDnsResolutionBindingsOptions) with pagination`, func() {
			listVPCDnsResolutionBindingsOptions := &vpcv1.ListVPCDnsResolutionBindingsOptions{
				VPCID:     core.StringPtr("testString"),
				AccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:     core.Int64Ptr(int64(10)),
				Name:      core.StringPtr("my-name"),
				Sort:      core.StringPtr("name"),
				Start:     core.StringPtr("testString"),
				VPCCRN:    core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCName:   core.StringPtr("my-vpc"),
			}

			listVPCDnsResolutionBindingsOptions.Start = nil
			listVPCDnsResolutionBindingsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VpcdnsResolutionBinding
			for {
				vpcdnsResolutionBindingCollection, response, err := vpcService.ListVPCDnsResolutionBindings(listVPCDnsResolutionBindingsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpcdnsResolutionBindingCollection).ToNot(BeNil())
				allResults = append(allResults, vpcdnsResolutionBindingCollection.DnsResolutionBindings...)

				nextToken, nextErr := vpcdnsResolutionBindingCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPCDnsResolutionBindingsOptions.Start = nextToken

				if listVPCDnsResolutionBindingsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPCDnsResolutionBindings(listVPCDnsResolutionBindingsOptions *ListVPCDnsResolutionBindingsOptions) using VPCDnsResolutionBindingsPager`, func() {
			listVPCDnsResolutionBindingsOptions := &vpcv1.ListVPCDnsResolutionBindingsOptions{
				VPCID:     core.StringPtr("testString"),
				AccountID: core.StringPtr("bb1b52262f7441a586f49068482f1e60"),
				Limit:     core.Int64Ptr(int64(10)),
				Name:      core.StringPtr("my-name"),
				Sort:      core.StringPtr("name"),
				VPCCRN:    core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
				VPCName:   core.StringPtr("my-vpc"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPCDnsResolutionBindingsPager(listVPCDnsResolutionBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VpcdnsResolutionBinding
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPCDnsResolutionBindingsPager(listVPCDnsResolutionBindingsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPCDnsResolutionBindings() returned a total of %d item(s) using VPCDnsResolutionBindingsPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPCDnsResolutionBinding - Create a DNS resolution binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPCDnsResolutionBinding(createVPCDnsResolutionBindingOptions *CreateVPCDnsResolutionBindingOptions)`, func() {
			vpcIdentityModel := &vpcv1.VPCIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"),
			}

			createVPCDnsResolutionBindingOptions := &vpcv1.CreateVPCDnsResolutionBindingOptions{
				VPCID: core.StringPtr("testString"),
				VPC:   vpcIdentityModel,
				Name:  core.StringPtr("my-dns-resolution-binding"),
			}

			vpcdnsResolutionBinding, response, err := vpcService.CreateVPCDnsResolutionBinding(createVPCDnsResolutionBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
	})

	Describe(`GetVPCDnsResolutionBinding - Retrieve a DNS resolution binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCDnsResolutionBinding(getVPCDnsResolutionBindingOptions *GetVPCDnsResolutionBindingOptions)`, func() {
			getVPCDnsResolutionBindingOptions := &vpcv1.GetVPCDnsResolutionBindingOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			vpcdnsResolutionBinding, response, err := vpcService.GetVPCDnsResolutionBinding(getVPCDnsResolutionBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPCDnsResolutionBinding - Update a DNS resolution binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPCDnsResolutionBinding(updateVPCDnsResolutionBindingOptions *UpdateVPCDnsResolutionBindingOptions)`, func() {
			vpcdnsResolutionBindingPatchModel := &vpcv1.VpcdnsResolutionBindingPatch{
				Name: core.StringPtr("my-dns-resolution-binding-updated"),
			}
			vpcdnsResolutionBindingPatchModelAsPatch, asPatchErr := vpcdnsResolutionBindingPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCDnsResolutionBindingOptions := &vpcv1.UpdateVPCDnsResolutionBindingOptions{
				ID:                           core.StringPtr("testString"),
				VPCID:                        core.StringPtr("testString"),
				VpcdnsResolutionBindingPatch: vpcdnsResolutionBindingPatchModelAsPatch,
			}

			vpcdnsResolutionBinding, response, err := vpcService.UpdateVPCDnsResolutionBinding(updateVPCDnsResolutionBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
	})

	Describe(`ListVPCRoutes - List routes in a VPC's default routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPCRoutes(listVPCRoutesOptions *ListVPCRoutesOptions) with pagination`, func() {
			listVPCRoutesOptions := &vpcv1.ListVPCRoutesOptions{
				VPCID:    core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
				Start:    core.StringPtr("testString"),
				ZoneName: core.StringPtr("us-south-1"),
			}

			listVPCRoutesOptions.Start = nil
			listVPCRoutesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Route
			for {
				routeCollectionVPCContext, response, err := vpcService.ListVPCRoutes(listVPCRoutesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(routeCollectionVPCContext).ToNot(BeNil())
				allResults = append(allResults, routeCollectionVPCContext.Routes...)

				nextToken, nextErr := routeCollectionVPCContext.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPCRoutesOptions.Start = nextToken

				if listVPCRoutesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPCRoutes(listVPCRoutesOptions *ListVPCRoutesOptions) using VPCRoutesPager`, func() {
			listVPCRoutesOptions := &vpcv1.ListVPCRoutesOptions{
				VPCID:    core.StringPtr("testString"),
				Limit:    core.Int64Ptr(int64(10)),
				ZoneName: core.StringPtr("us-south-1"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPCRoutesPager(listVPCRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Route
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPCRoutesPager(listVPCRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPCRoutes() returned a total of %d item(s) using VPCRoutesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPCRoute - Create a route in a VPC's default routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPCRoute(createVPCRouteOptions *CreateVPCRouteOptions)`, func() {
			routePrototypeZoneModel := &vpcv1.RoutePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			routeNextHopPrototypeModel := &vpcv1.RouteNextHopPrototypeRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("0.0.0.0"),
			}

			createVPCRouteOptions := &vpcv1.CreateVPCRouteOptions{
				VPCID:       core.StringPtr("testString"),
				Destination: core.StringPtr("192.168.3.0/24"),
				Zone:        routePrototypeZoneModel,
				Action:      core.StringPtr("deliver"),
				Advertise:   core.BoolPtr(false),
				Name:        core.StringPtr("my-vpc-routing-table-route"),
				NextHop:     routeNextHopPrototypeModel,
				Priority:    core.Int64Ptr(int64(1)),
			}

			createVPCRouteResponse, response, err := vpcService.CreateVPCRoute(createVPCRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(createVPCRouteResponse).ToNot(BeNil())
		})
	})

	Describe(`GetVPCRoute - Retrieve a VPC route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCRoute(getVPCRouteOptions *GetVPCRouteOptions)`, func() {
			getVPCRouteOptions := &vpcv1.GetVPCRouteOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			getVPCRouteResponse, response, err := vpcService.GetVPCRoute(getVPCRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(getVPCRouteResponse).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPCRoute - Update a VPC route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPCRoute(updateVPCRouteOptions *UpdateVPCRouteOptions)`, func() {
			routeNextHopPatchModel := &vpcv1.RouteNextHopPatchRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("10.240.0.12"),
			}

			routePatchModel := &vpcv1.RoutePatch{
				Advertise: core.BoolPtr(true),
				Name:      core.StringPtr("my-vpc-routing-table-route"),
				NextHop:   routeNextHopPatchModel,
				Priority:  core.Int64Ptr(int64(1)),
			}
			routePatchModelAsPatch, asPatchErr := routePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRouteOptions := &vpcv1.UpdateVPCRouteOptions{
				ID:         core.StringPtr("testString"),
				VPCID:      core.StringPtr("testString"),
				RoutePatch: routePatchModelAsPatch,
			}

			updateVPCRouteResponse, response, err := vpcService.UpdateVPCRoute(updateVPCRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(updateVPCRouteResponse).ToNot(BeNil())
		})
	})

	Describe(`ListVPCRoutingTables - List routing tables for a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPCRoutingTables(listVPCRoutingTablesOptions *ListVPCRoutingTablesOptions) with pagination`, func() {
			listVPCRoutingTablesOptions := &vpcv1.ListVPCRoutingTablesOptions{
				VPCID:     core.StringPtr("testString"),
				IsDefault: core.BoolPtr(true),
				Limit:     core.Int64Ptr(int64(10)),
				Start:     core.StringPtr("testString"),
			}

			listVPCRoutingTablesOptions.Start = nil
			listVPCRoutingTablesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.RoutingTable
			for {
				routingTableCollection, response, err := vpcService.ListVPCRoutingTables(listVPCRoutingTablesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(routingTableCollection).ToNot(BeNil())
				allResults = append(allResults, routingTableCollection.RoutingTables...)

				nextToken, nextErr := routingTableCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPCRoutingTablesOptions.Start = nextToken

				if listVPCRoutingTablesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPCRoutingTables(listVPCRoutingTablesOptions *ListVPCRoutingTablesOptions) using VPCRoutingTablesPager`, func() {
			listVPCRoutingTablesOptions := &vpcv1.ListVPCRoutingTablesOptions{
				VPCID:     core.StringPtr("testString"),
				IsDefault: core.BoolPtr(true),
				Limit:     core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPCRoutingTablesPager(listVPCRoutingTablesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.RoutingTable
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPCRoutingTablesPager(listVPCRoutingTablesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPCRoutingTables() returned a total of %d item(s) using VPCRoutingTablesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPCRoutingTable - Create a routing table for a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPCRoutingTable(createVPCRoutingTableOptions *CreateVPCRoutingTableOptions)`, func() {
			resourceFilterModel := &vpcv1.ResourceFilter{
				ResourceType: core.StringPtr("vpn_server"),
			}

			routeNextHopPrototypeModel := &vpcv1.RouteNextHopPrototypeRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("0.0.0.0"),
			}

			routePrototypeZoneModel := &vpcv1.RoutePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			routePrototypeModel := &vpcv1.RoutePrototype{
				Action:      core.StringPtr("deliver"),
				Advertise:   core.BoolPtr(false),
				Destination: core.StringPtr("192.168.3.0/24"),
				Name:        core.StringPtr("my-vpc-routing-table-route"),
				NextHop:     routeNextHopPrototypeModel,
				Priority:    core.Int64Ptr(int64(1)),
				Zone:        routePrototypeZoneModel,
			}

			createVPCRoutingTableOptions := &vpcv1.CreateVPCRoutingTableOptions{
				VPCID:                      core.StringPtr("testString"),
				AcceptRoutesFrom:           []vpcv1.ResourceFilter{*resourceFilterModel},
				AdvertiseRoutesTo:          []string{"transit_gateway"},
				Name:                       core.StringPtr("my-routing-table-1"),
				RouteDirectLinkIngress:     core.BoolPtr(false),
				RouteInternetIngress:       core.BoolPtr(false),
				RouteTransitGatewayIngress: core.BoolPtr(false),
				RouteVPCZoneIngress:        core.BoolPtr(false),
				Routes:                     []vpcv1.RoutePrototype{*routePrototypeModel},
			}

			routingTable, response, err := vpcService.CreateVPCRoutingTable(createVPCRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(routingTable).ToNot(BeNil())
		})
	})

	Describe(`GetVPCRoutingTable - Retrieve a VPC routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCRoutingTable(getVPCRoutingTableOptions *GetVPCRoutingTableOptions)`, func() {
			getVPCRoutingTableOptions := &vpcv1.GetVPCRoutingTableOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			routingTable, response, err := vpcService.GetVPCRoutingTable(getVPCRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPCRoutingTable - Update a VPC routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPCRoutingTable(updateVPCRoutingTableOptions *UpdateVPCRoutingTableOptions)`, func() {
			resourceFilterModel := &vpcv1.ResourceFilter{
				ResourceType: core.StringPtr("vpn_server"),
			}

			routingTablePatchModel := &vpcv1.RoutingTablePatch{
				AcceptRoutesFrom:           []vpcv1.ResourceFilter{*resourceFilterModel},
				AdvertiseRoutesTo:          []string{"transit_gateway"},
				Name:                       core.StringPtr("my-routing-table-2"),
				RouteDirectLinkIngress:     core.BoolPtr(true),
				RouteInternetIngress:       core.BoolPtr(true),
				RouteTransitGatewayIngress: core.BoolPtr(true),
				RouteVPCZoneIngress:        core.BoolPtr(true),
			}
			routingTablePatchModelAsPatch, asPatchErr := routingTablePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRoutingTableOptions := &vpcv1.UpdateVPCRoutingTableOptions{
				ID:                core.StringPtr("testString"),
				VPCID:             core.StringPtr("testString"),
				RoutingTablePatch: routingTablePatchModelAsPatch,
				IfMatch:           core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			routingTable, response, err := vpcService.UpdateVPCRoutingTable(updateVPCRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routingTable).ToNot(BeNil())
		})
	})

	Describe(`ListVPCRoutingTableRoutes - List routes in a VPC routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPCRoutingTableRoutes(listVPCRoutingTableRoutesOptions *ListVPCRoutingTableRoutesOptions) with pagination`, func() {
			listVPCRoutingTableRoutesOptions := &vpcv1.ListVPCRoutingTableRoutesOptions{
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
				Limit:          core.Int64Ptr(int64(10)),
				Start:          core.StringPtr("testString"),
			}

			listVPCRoutingTableRoutesOptions.Start = nil
			listVPCRoutingTableRoutesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.Route
			for {
				routeCollection, response, err := vpcService.ListVPCRoutingTableRoutes(listVPCRoutingTableRoutesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(routeCollection).ToNot(BeNil())
				allResults = append(allResults, routeCollection.Routes...)

				nextToken, nextErr := routeCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPCRoutingTableRoutesOptions.Start = nextToken

				if listVPCRoutingTableRoutesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPCRoutingTableRoutes(listVPCRoutingTableRoutesOptions *ListVPCRoutingTableRoutesOptions) using VPCRoutingTableRoutesPager`, func() {
			listVPCRoutingTableRoutesOptions := &vpcv1.ListVPCRoutingTableRoutesOptions{
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
				Limit:          core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPCRoutingTableRoutesPager(listVPCRoutingTableRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.Route
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPCRoutingTableRoutesPager(listVPCRoutingTableRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPCRoutingTableRoutes() returned a total of %d item(s) using VPCRoutingTableRoutesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPCRoutingTableRoute - Create a route in a VPC routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPCRoutingTableRoute(createVPCRoutingTableRouteOptions *CreateVPCRoutingTableRouteOptions)`, func() {
			routePrototypeZoneModel := &vpcv1.RoutePrototypeZoneZoneIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"),
			}

			routeNextHopPrototypeModel := &vpcv1.RouteNextHopPrototypeRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("0.0.0.0"),
			}

			createVPCRoutingTableRouteOptions := &vpcv1.CreateVPCRoutingTableRouteOptions{
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
				Destination:    core.StringPtr("192.168.3.0/24"),
				Zone:           routePrototypeZoneModel,
				Action:         core.StringPtr("deliver"),
				Advertise:      core.BoolPtr(false),
				Name:           core.StringPtr("my-vpc-routing-table-route"),
				NextHop:        routeNextHopPrototypeModel,
				Priority:       core.Int64Ptr(int64(1)),
			}

			route, response, err := vpcService.CreateVPCRoutingTableRoute(createVPCRoutingTableRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(route).ToNot(BeNil())
		})
	})

	Describe(`GetVPCRoutingTableRoute - Retrieve a VPC routing table route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPCRoutingTableRoute(getVPCRoutingTableRouteOptions *GetVPCRoutingTableRouteOptions)`, func() {
			getVPCRoutingTableRouteOptions := &vpcv1.GetVPCRoutingTableRouteOptions{
				ID:             core.StringPtr("testString"),
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
			}

			route, response, err := vpcService.GetVPCRoutingTableRoute(getVPCRoutingTableRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPCRoutingTableRoute - Update a VPC routing table route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPCRoutingTableRoute(updateVPCRoutingTableRouteOptions *UpdateVPCRoutingTableRouteOptions)`, func() {
			routeNextHopPatchModel := &vpcv1.RouteNextHopPatchRouteNextHopIPRouteNextHopIPSentinelIP{
				Address: core.StringPtr("0.0.0.0"),
			}

			routePatchModel := &vpcv1.RoutePatch{
				Advertise: core.BoolPtr(false),
				Name:      core.StringPtr("my-vpc-routing-table-route"),
				NextHop:   routeNextHopPatchModel,
				Priority:  core.Int64Ptr(int64(1)),
			}
			routePatchModelAsPatch, asPatchErr := routePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPCRoutingTableRouteOptions := &vpcv1.UpdateVPCRoutingTableRouteOptions{
				ID:             core.StringPtr("testString"),
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
				RoutePatch:     routePatchModelAsPatch,
			}

			route, response, err := vpcService.UpdateVPCRoutingTableRoute(updateVPCRoutingTableRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
	})

	Describe(`ListIkePolicies - List IKE policies`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListIkePolicies(listIkePoliciesOptions *ListIkePoliciesOptions) with pagination`, func() {
			listIkePoliciesOptions := &vpcv1.ListIkePoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listIkePoliciesOptions.Start = nil
			listIkePoliciesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.IkePolicy
			for {
				ikePolicyCollection, response, err := vpcService.ListIkePolicies(listIkePoliciesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(ikePolicyCollection).ToNot(BeNil())
				allResults = append(allResults, ikePolicyCollection.IkePolicies...)

				nextToken, nextErr := ikePolicyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listIkePoliciesOptions.Start = nextToken

				if listIkePoliciesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListIkePolicies(listIkePoliciesOptions *ListIkePoliciesOptions) using IkePoliciesPager`, func() {
			listIkePoliciesOptions := &vpcv1.ListIkePoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewIkePoliciesPager(listIkePoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.IkePolicy
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewIkePoliciesPager(listIkePoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListIkePolicies() returned a total of %d item(s) using IkePoliciesPager.\n", len(allResults))
		})
	})

	Describe(`CreateIkePolicy - Create an IKE policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateIkePolicy(createIkePolicyOptions *CreateIkePolicyOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createIkePolicyOptions := &vpcv1.CreateIkePolicyOptions{
				IkeVersion:               core.Int64Ptr(int64(1)),
				AuthenticationAlgorithm:  core.StringPtr("sha256"),
				AuthenticationAlgorithms: []string{"sha256"},
				DhGroup:                  core.Int64Ptr(int64(14)),
				DhGroups:                 []int64{int64(14)},
				EncryptionAlgorithm:      core.StringPtr("aes128"),
				EncryptionAlgorithms:     []string{"aes128"},
				KeyLifetime:              core.Int64Ptr(int64(28800)),
				Name:                     core.StringPtr("my-ike-policy"),
				ResourceGroup:            resourceGroupIdentityModel,
			}

			ikePolicy, response, err := vpcService.CreateIkePolicy(createIkePolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(ikePolicy).ToNot(BeNil())
		})
	})

	Describe(`GetIkePolicy - Retrieve an IKE policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetIkePolicy(getIkePolicyOptions *GetIkePolicyOptions)`, func() {
			getIkePolicyOptions := &vpcv1.GetIkePolicyOptions{
				ID: core.StringPtr("testString"),
			}

			ikePolicy, response, err := vpcService.GetIkePolicy(getIkePolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(ikePolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdateIkePolicy - Update an IKE policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateIkePolicy(updateIkePolicyOptions *UpdateIkePolicyOptions)`, func() {
			ikePolicyPatchModel := &vpcv1.IkePolicyPatch{
				AuthenticationAlgorithm:  core.StringPtr("sha256"),
				AuthenticationAlgorithms: []string{"sha256"},
				DhGroup:                  core.Int64Ptr(int64(14)),
				DhGroups:                 []int64{int64(14)},
				EncryptionAlgorithm:      core.StringPtr("aes128"),
				EncryptionAlgorithms:     []string{"aes128"},
				IkeVersion:               core.Int64Ptr(int64(1)),
				KeyLifetime:              core.Int64Ptr(int64(86400)),
				Name:                     core.StringPtr("my-ike-policy"),
			}
			ikePolicyPatchModelAsPatch, asPatchErr := ikePolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateIkePolicyOptions := &vpcv1.UpdateIkePolicyOptions{
				ID:             core.StringPtr("testString"),
				IkePolicyPatch: ikePolicyPatchModelAsPatch,
			}

			ikePolicy, response, err := vpcService.UpdateIkePolicy(updateIkePolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(ikePolicy).ToNot(BeNil())
		})
	})

	Describe(`ListIkePolicyConnections - List VPN gateway connections that use a specified IKE policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListIkePolicyConnections(listIkePolicyConnectionsOptions *ListIkePolicyConnectionsOptions) with pagination`, func() {
			listIkePolicyConnectionsOptions := &vpcv1.ListIkePolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listIkePolicyConnectionsOptions.Start = nil
			listIkePolicyConnectionsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for {
				ikePolicyConnectionCollection, response, err := vpcService.ListIkePolicyConnections(listIkePolicyConnectionsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(ikePolicyConnectionCollection).ToNot(BeNil())
				allResults = append(allResults, ikePolicyConnectionCollection.Connections...)

				nextToken, nextErr := ikePolicyConnectionCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listIkePolicyConnectionsOptions.Start = nextToken

				if listIkePolicyConnectionsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListIkePolicyConnections(listIkePolicyConnectionsOptions *ListIkePolicyConnectionsOptions) using IkePolicyConnectionsPager`, func() {
			listIkePolicyConnectionsOptions := &vpcv1.ListIkePolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewIkePolicyConnectionsPager(listIkePolicyConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewIkePolicyConnectionsPager(listIkePolicyConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListIkePolicyConnections() returned a total of %d item(s) using IkePolicyConnectionsPager.\n", len(allResults))
		})
	})

	Describe(`ListIpsecPolicies - List IPsec policies`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListIpsecPolicies(listIpsecPoliciesOptions *ListIpsecPoliciesOptions) with pagination`, func() {
			listIpsecPoliciesOptions := &vpcv1.ListIpsecPoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listIpsecPoliciesOptions.Start = nil
			listIpsecPoliciesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.IPsecPolicy
			for {
				iPsecPolicyCollection, response, err := vpcService.ListIpsecPolicies(listIpsecPoliciesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(iPsecPolicyCollection).ToNot(BeNil())
				allResults = append(allResults, iPsecPolicyCollection.IpsecPolicies...)

				nextToken, nextErr := iPsecPolicyCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listIpsecPoliciesOptions.Start = nextToken

				if listIpsecPoliciesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListIpsecPolicies(listIpsecPoliciesOptions *ListIpsecPoliciesOptions) using IpsecPoliciesPager`, func() {
			listIpsecPoliciesOptions := &vpcv1.ListIpsecPoliciesOptions{
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewIpsecPoliciesPager(listIpsecPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.IPsecPolicy
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewIpsecPoliciesPager(listIpsecPoliciesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListIpsecPolicies() returned a total of %d item(s) using IpsecPoliciesPager.\n", len(allResults))
		})
	})

	Describe(`CreateIpsecPolicy - Create an IPsec policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateIpsecPolicy(createIpsecPolicyOptions *CreateIpsecPolicyOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			createIpsecPolicyOptions := &vpcv1.CreateIpsecPolicyOptions{
				AuthenticationAlgorithm:  core.StringPtr("disabled"),
				AuthenticationAlgorithms: []string{"disabled"},
				EncryptionAlgorithm:      core.StringPtr("aes128"),
				EncryptionAlgorithms:     []string{"aes128"},
				KeyLifetime:              core.Int64Ptr(int64(3600)),
				Name:                     core.StringPtr("my-ipsec-policy"),
				Pfs:                      core.StringPtr("disabled"),
				PfsGroups:                []string{"disabled"},
				ResourceGroup:            resourceGroupIdentityModel,
			}

			iPsecPolicy, response, err := vpcService.CreateIpsecPolicy(createIpsecPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
	})

	Describe(`GetIpsecPolicy - Retrieve an IPsec policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetIpsecPolicy(getIpsecPolicyOptions *GetIpsecPolicyOptions)`, func() {
			getIpsecPolicyOptions := &vpcv1.GetIpsecPolicyOptions{
				ID: core.StringPtr("testString"),
			}

			iPsecPolicy, response, err := vpcService.GetIpsecPolicy(getIpsecPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
	})

	Describe(`UpdateIpsecPolicy - Update an IPsec policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateIpsecPolicy(updateIpsecPolicyOptions *UpdateIpsecPolicyOptions)`, func() {
			iPsecPolicyPatchModel := &vpcv1.IPsecPolicyPatch{
				AuthenticationAlgorithm:  core.StringPtr("disabled"),
				AuthenticationAlgorithms: []string{"disabled"},
				EncryptionAlgorithm:      core.StringPtr("aes128"),
				EncryptionAlgorithms:     []string{"aes128"},
				KeyLifetime:              core.Int64Ptr(int64(3600)),
				Name:                     core.StringPtr("my-ipsec-policy"),
				Pfs:                      core.StringPtr("disabled"),
				PfsGroups:                []string{"disabled"},
			}
			iPsecPolicyPatchModelAsPatch, asPatchErr := iPsecPolicyPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateIpsecPolicyOptions := &vpcv1.UpdateIpsecPolicyOptions{
				ID:               core.StringPtr("testString"),
				IPsecPolicyPatch: iPsecPolicyPatchModelAsPatch,
			}

			iPsecPolicy, response, err := vpcService.UpdateIpsecPolicy(updateIpsecPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(iPsecPolicy).ToNot(BeNil())
		})
	})

	Describe(`ListIpsecPolicyConnections - List VPN gateway connections that use a specified IPsec policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListIpsecPolicyConnections(listIpsecPolicyConnectionsOptions *ListIpsecPolicyConnectionsOptions) with pagination`, func() {
			listIpsecPolicyConnectionsOptions := &vpcv1.ListIpsecPolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
				Start: core.StringPtr("testString"),
			}

			listIpsecPolicyConnectionsOptions.Start = nil
			listIpsecPolicyConnectionsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for {
				iPsecPolicyConnectionCollection, response, err := vpcService.ListIpsecPolicyConnections(listIpsecPolicyConnectionsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(iPsecPolicyConnectionCollection).ToNot(BeNil())
				allResults = append(allResults, iPsecPolicyConnectionCollection.Connections...)

				nextToken, nextErr := iPsecPolicyConnectionCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listIpsecPolicyConnectionsOptions.Start = nextToken

				if listIpsecPolicyConnectionsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListIpsecPolicyConnections(listIpsecPolicyConnectionsOptions *ListIpsecPolicyConnectionsOptions) using IpsecPolicyConnectionsPager`, func() {
			listIpsecPolicyConnectionsOptions := &vpcv1.ListIpsecPolicyConnectionsOptions{
				ID:    core.StringPtr("testString"),
				Limit: core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewIpsecPolicyConnectionsPager(listIpsecPolicyConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewIpsecPolicyConnectionsPager(listIpsecPolicyConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListIpsecPolicyConnections() returned a total of %d item(s) using IpsecPolicyConnectionsPager.\n", len(allResults))
		})
	})

	Describe(`ListVPNGateways - List VPN gateways`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGateways(listVPNGatewaysOptions *ListVPNGatewaysOptions) with pagination`, func() {
			listVPNGatewaysOptions := &vpcv1.ListVPNGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Mode:            core.StringPtr("policy"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				Start:           core.StringPtr("testString"),
			}

			listVPNGatewaysOptions.Start = nil
			listVPNGatewaysOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNGatewayIntf
			for {
				vpnGatewayCollection, response, err := vpcService.ListVPNGateways(listVPNGatewaysOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnGatewayCollection).ToNot(BeNil())
				allResults = append(allResults, vpnGatewayCollection.VPNGateways...)

				nextToken, nextErr := vpnGatewayCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNGatewaysOptions.Start = nextToken

				if listVPNGatewaysOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNGateways(listVPNGatewaysOptions *ListVPNGatewaysOptions) using VPNGatewaysPager`, func() {
			listVPNGatewaysOptions := &vpcv1.ListVPNGatewaysOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Mode:            core.StringPtr("policy"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNGatewaysPager(listVPNGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNGatewayIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNGatewaysPager(listVPNGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNGateways() returned a total of %d item(s) using VPNGatewaysPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPNGateway - Create a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPNGateway(createVPNGatewayOptions *CreateVPNGatewayOptions)`, func() {
			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

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
				Name:             core.StringPtr("my-vpn-gateway"),
				ResourceGroup:    resourceGroupIdentityModel,
				AvailabilityMode: core.StringPtr("regional"),
				Members:          []vpcv1.VPNGatewayMemberPrototype{*vpnGatewayMemberPrototypeModel},
				AdvertisedCidrs:  []string{"192.168.3.0/24"},
				LocalAsn:         core.Int64Ptr(int64(64520)),
				Mode:             core.StringPtr("route"),
			}

			createVPNGatewayOptions := &vpcv1.CreateVPNGatewayOptions{
				VPNGatewayPrototype: vpnGatewayPrototypeModel,
			}

			vpnGateway, response, err := vpcService.CreateVPNGateway(createVPNGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnGateway).ToNot(BeNil())
		})
	})

	Describe(`GetVPNGateway - Retrieve a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNGateway(getVPNGatewayOptions *GetVPNGatewayOptions)`, func() {
			getVPNGatewayOptions := &vpcv1.GetVPNGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			vpnGateway, response, err := vpcService.GetVPNGateway(getVPNGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGateway).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPNGateway - Update a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPNGateway(updateVPNGatewayOptions *UpdateVPNGatewayOptions)`, func() {
			vpnGatewayPatchModel := &vpcv1.VPNGatewayPatch{
				AvailabilityMode: core.StringPtr("regional"),
				LocalAsn:         core.Int64Ptr(int64(64520)),
				Name:             core.StringPtr("my-vpn-gateway-updated"),
			}
			vpnGatewayPatchModelAsPatch, asPatchErr := vpnGatewayPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNGatewayOptions := &vpcv1.UpdateVPNGatewayOptions{
				ID:              core.StringPtr("testString"),
				VPNGatewayPatch: vpnGatewayPatchModelAsPatch,
			}

			vpnGateway, response, err := vpcService.UpdateVPNGateway(updateVPNGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGateway).ToNot(BeNil())
		})
	})

	Describe(`ListVPNGatewayAdvertisedCidrs - List advertised CIDRs for a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayAdvertisedCidrs(listVPNGatewayAdvertisedCidrsOptions *ListVPNGatewayAdvertisedCidrsOptions)`, func() {
			listVPNGatewayAdvertisedCidrsOptions := &vpcv1.ListVPNGatewayAdvertisedCidrsOptions{
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayAdvertisedCIDRCollection, response, err := vpcService.ListVPNGatewayAdvertisedCidrs(listVPNGatewayAdvertisedCidrsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayAdvertisedCIDRCollection).ToNot(BeNil())
		})
	})

	Describe(`CheckVPNGatewayAdvertisedCIDR - Check if the specified advertised CIDR exists on a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions *CheckVPNGatewayAdvertisedCIDROptions)`, func() {
			checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`AddVPNGatewayAdvertisedCIDR - Set an advertised CIDR on a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddVPNGatewayAdvertisedCIDR(addVPNGatewayAdvertisedCIDROptions *AddVPNGatewayAdvertisedCIDROptions)`, func() {
			addVPNGatewayAdvertisedCIDROptions := &vpcv1.AddVPNGatewayAdvertisedCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.AddVPNGatewayAdvertisedCIDR(addVPNGatewayAdvertisedCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`ListVPNGatewayConnections - List connections of a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayConnections(listVPNGatewayConnectionsOptions *ListVPNGatewayConnectionsOptions) with pagination`, func() {
			listVPNGatewayConnectionsOptions := &vpcv1.ListVPNGatewayConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
				Start:        core.StringPtr("testString"),
				Status:       core.StringPtr("down"),
			}

			listVPNGatewayConnectionsOptions.Start = nil
			listVPNGatewayConnectionsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for {
				vpnGatewayConnectionCollection, response, err := vpcService.ListVPNGatewayConnections(listVPNGatewayConnectionsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnGatewayConnectionCollection).ToNot(BeNil())
				allResults = append(allResults, vpnGatewayConnectionCollection.Connections...)

				nextToken, nextErr := vpnGatewayConnectionCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNGatewayConnectionsOptions.Start = nextToken

				if listVPNGatewayConnectionsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNGatewayConnections(listVPNGatewayConnectionsOptions *ListVPNGatewayConnectionsOptions) using VPNGatewayConnectionsPager`, func() {
			listVPNGatewayConnectionsOptions := &vpcv1.ListVPNGatewayConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
				Status:       core.StringPtr("down"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNGatewayConnectionsPager(listVPNGatewayConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNGatewayConnectionIntf
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNGatewayConnectionsPager(listVPNGatewayConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNGatewayConnections() returned a total of %d item(s) using VPNGatewayConnectionsPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPNGatewayConnection - Create a connection for a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPNGatewayConnection(createVPNGatewayConnectionOptions *CreateVPNGatewayConnectionOptions)`, func() {
			vpnGatewayConnectionDpdPrototypeModel := &vpcv1.VPNGatewayConnectionDpdPrototype{
				Action:   core.StringPtr("restart"),
				Interval: core.Int64Ptr(int64(30)),
				Timeout:  core.Int64Ptr(int64(120)),
			}

			vpnGatewayConnectionIkePolicyPrototypeModel := &vpcv1.VPNGatewayConnectionIkePolicyPrototypeIkePolicyIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/ike_policies/r006-e98f46a3-1e4e-4195-b4e5-b8155192689d"),
			}

			vpnGatewayConnectionIPsecPolicyPrototypeModel := &vpcv1.VPNGatewayConnectionIPsecPolicyPrototypeIPsecPolicyIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/ipsec_policies/r006-51eae621-dbbc-4c47-b623-b57a43c19876"),
			}

			vpnGatewayConnectionIkeIdentityPrototypeModel := &vpcv1.VPNGatewayConnectionIkeIdentityPrototypeVPNGatewayConnectionIkeIdentityFqdn{
				Type:  core.StringPtr("fqdn"),
				Value: core.StringPtr("my-service.example.com"),
			}

			vpnGatewayConnectionDynamicRouteModeLocalPrototypeModel := &vpcv1.VPNGatewayConnectionDynamicRouteModeLocalPrototype{
				IkeIdentities: []vpcv1.VPNGatewayConnectionIkeIdentityPrototypeIntf{vpnGatewayConnectionIkeIdentityPrototypeModel},
			}

			vpnGatewayConnectionDynamicRouteModePeerPrototypeModel := &vpcv1.VPNGatewayConnectionDynamicRouteModePeerPrototypeVPNGatewayConnectionPeerByAddress{
				Asn:         core.Int64Ptr(int64(38)),
				IkeIdentity: vpnGatewayConnectionIkeIdentityPrototypeModel,
				Address:     core.StringPtr("192.0.2.5"),
			}

			ipModel := &vpcv1.IP{
				Address: core.StringPtr("169.254.0.2"),
			}

			vpnGatewayConnectionTunnelPrototypeModel := &vpcv1.VPNGatewayConnectionTunnelPrototype{
				NeighborIP:        ipModel,
				TunnelInterfaceIP: ipModel,
			}

			vpnGatewayConnectionPrototypeModel := &vpcv1.VPNGatewayConnectionPrototypeVPNGatewayConnectionDynamicRouteModePrototype{
				AdminStateUp:      core.BoolPtr(true),
				DeadPeerDetection: vpnGatewayConnectionDpdPrototypeModel,
				EstablishMode:     core.StringPtr("bidirectional"),
				IkePolicy:         vpnGatewayConnectionIkePolicyPrototypeModel,
				IpsecPolicy:       vpnGatewayConnectionIPsecPolicyPrototypeModel,
				Name:              core.StringPtr("my-vpn-gateway-connection"),
				Psk:               core.StringPtr("lkj14b1oi0alcniejkso"),
				DistributeTraffic: core.BoolPtr(false),
				Local:             vpnGatewayConnectionDynamicRouteModeLocalPrototypeModel,
				Peer:              vpnGatewayConnectionDynamicRouteModePeerPrototypeModel,
				RoutingProtocol:   core.StringPtr("bgp"),
				Tunnels:           []vpcv1.VPNGatewayConnectionTunnelPrototype{*vpnGatewayConnectionTunnelPrototypeModel},
			}

			createVPNGatewayConnectionOptions := &vpcv1.CreateVPNGatewayConnectionOptions{
				VPNGatewayID:                  core.StringPtr("testString"),
				VPNGatewayConnectionPrototype: vpnGatewayConnectionPrototypeModel,
			}

			vpnGatewayConnection, response, err := vpcService.CreateVPNGatewayConnection(createVPNGatewayConnectionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
	})

	Describe(`GetVPNGatewayConnection - Retrieve a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNGatewayConnection(getVPNGatewayConnectionOptions *GetVPNGatewayConnectionOptions)`, func() {
			getVPNGatewayConnectionOptions := &vpcv1.GetVPNGatewayConnectionOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayConnection, response, err := vpcService.GetVPNGatewayConnection(getVPNGatewayConnectionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPNGatewayConnection - Update a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPNGatewayConnection(updateVPNGatewayConnectionOptions *UpdateVPNGatewayConnectionOptions)`, func() {
			vpnGatewayConnectionDpdPatchModel := &vpcv1.VPNGatewayConnectionDpdPatch{
				Action:   core.StringPtr("restart"),
				Interval: core.Int64Ptr(int64(30)),
				Timeout:  core.Int64Ptr(int64(120)),
			}

			vpnGatewayConnectionIkePolicyPatchModel := &vpcv1.VPNGatewayConnectionIkePolicyPatchIkePolicyIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/ike_policies/r006-e98f46a3-1e4e-4195-b4e5-b8155192689d"),
			}

			vpnGatewayConnectionIPsecPolicyPatchModel := &vpcv1.VPNGatewayConnectionIPsecPolicyPatchIPsecPolicyIdentityByHref{
				Href: core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/ipsec_policies/r006-51eae621-dbbc-4c47-b623-b57a43c19876"),
			}

			vpnGatewayConnectionPeerPatchModel := &vpcv1.VPNGatewayConnectionPeerPatchVPNGatewayConnectionDynamicRouteModePeerPatchVPNGatewayConnectionDynamicRouteModePeerPatchVPNGatewayConnectionPeerAddressPatch{
				Asn:     core.Int64Ptr(int64(38)),
				Address: core.StringPtr("192.0.2.5"),
			}

			ipModel := &vpcv1.IP{
				Address: core.StringPtr("169.254.0.2"),
			}

			vpnGatewayConnectionTunnelModel := &vpcv1.VPNGatewayConnectionTunnel{
				NeighborIP:        ipModel,
				TunnelInterfaceIP: ipModel,
			}

			vpnGatewayConnectionPatchModel := &vpcv1.VPNGatewayConnectionPatch{
				AdminStateUp:      core.BoolPtr(true),
				DeadPeerDetection: vpnGatewayConnectionDpdPatchModel,
				DistributeTraffic: core.BoolPtr(true),
				EstablishMode:     core.StringPtr("bidirectional"),
				IkePolicy:         vpnGatewayConnectionIkePolicyPatchModel,
				IpsecPolicy:       vpnGatewayConnectionIPsecPolicyPatchModel,
				Name:              core.StringPtr("my-vpn-gateway-connection"),
				Peer:              vpnGatewayConnectionPeerPatchModel,
				Psk:               core.StringPtr("lkj14b1oi0alcniejkso"),
				RoutingProtocol:   core.StringPtr("bgp"),
				Tunnels:           []vpcv1.VPNGatewayConnectionTunnel{*vpnGatewayConnectionTunnelModel},
			}
			vpnGatewayConnectionPatchModelAsPatch, asPatchErr := vpnGatewayConnectionPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNGatewayConnectionOptions := &vpcv1.UpdateVPNGatewayConnectionOptions{
				ID:                        core.StringPtr("testString"),
				VPNGatewayID:              core.StringPtr("testString"),
				VPNGatewayConnectionPatch: vpnGatewayConnectionPatchModelAsPatch,
				IfMatch:                   core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			vpnGatewayConnection, response, err := vpcService.UpdateVPNGatewayConnection(updateVPNGatewayConnectionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnection).ToNot(BeNil())
		})
	})

	Describe(`ListVPNGatewayConnectionsLocalCidrs - List local CIDRs for a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayConnectionsLocalCidrs(listVPNGatewayConnectionsLocalCidrsOptions *ListVPNGatewayConnectionsLocalCidrsOptions)`, func() {
			listVPNGatewayConnectionsLocalCidrsOptions := &vpcv1.ListVPNGatewayConnectionsLocalCidrsOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayConnectionCidRs, response, err := vpcService.ListVPNGatewayConnectionsLocalCidrs(listVPNGatewayConnectionsLocalCidrsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnectionCidRs).ToNot(BeNil())
		})
	})

	Describe(`CheckVPNGatewayConnectionsLocalCIDR - Check if the specified local CIDR exists on a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CheckVPNGatewayConnectionsLocalCIDR(checkVPNGatewayConnectionsLocalCIDROptions *CheckVPNGatewayConnectionsLocalCIDROptions)`, func() {
			checkVPNGatewayConnectionsLocalCIDROptions := &vpcv1.CheckVPNGatewayConnectionsLocalCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.CheckVPNGatewayConnectionsLocalCIDR(checkVPNGatewayConnectionsLocalCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`AddVPNGatewayConnectionsLocalCIDR - Set a local CIDR on a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddVPNGatewayConnectionsLocalCIDR(addVPNGatewayConnectionsLocalCIDROptions *AddVPNGatewayConnectionsLocalCIDROptions)`, func() {
			addVPNGatewayConnectionsLocalCIDROptions := &vpcv1.AddVPNGatewayConnectionsLocalCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.AddVPNGatewayConnectionsLocalCIDR(addVPNGatewayConnectionsLocalCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
		})
	})

	Describe(`ListVPNGatewayConnectionsPeerCidrs - List peer CIDRs for a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayConnectionsPeerCidrs(listVPNGatewayConnectionsPeerCidrsOptions *ListVPNGatewayConnectionsPeerCidrsOptions)`, func() {
			listVPNGatewayConnectionsPeerCidrsOptions := &vpcv1.ListVPNGatewayConnectionsPeerCidrsOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayConnectionCidRs, response, err := vpcService.ListVPNGatewayConnectionsPeerCidrs(listVPNGatewayConnectionsPeerCidrsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayConnectionCidRs).ToNot(BeNil())
		})
	})

	Describe(`CheckVPNGatewayConnectionsPeerCIDR - Check if the specified peer CIDR exists on a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CheckVPNGatewayConnectionsPeerCIDR(checkVPNGatewayConnectionsPeerCIDROptions *CheckVPNGatewayConnectionsPeerCIDROptions)`, func() {
			checkVPNGatewayConnectionsPeerCIDROptions := &vpcv1.CheckVPNGatewayConnectionsPeerCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.CheckVPNGatewayConnectionsPeerCIDR(checkVPNGatewayConnectionsPeerCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`AddVPNGatewayConnectionsPeerCIDR - Set a peer CIDR on a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`AddVPNGatewayConnectionsPeerCIDR(addVPNGatewayConnectionsPeerCIDROptions *AddVPNGatewayConnectionsPeerCIDROptions)`, func() {
			addVPNGatewayConnectionsPeerCIDROptions := &vpcv1.AddVPNGatewayConnectionsPeerCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.AddVPNGatewayConnectionsPeerCIDR(addVPNGatewayConnectionsPeerCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
		})
	})

	Describe(`ListVPNGatewayMembers - List members of a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayMembers(listVPNGatewayMembersOptions *ListVPNGatewayMembersOptions)`, func() {
			listVPNGatewayMembersOptions := &vpcv1.ListVPNGatewayMembersOptions{
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayMemberCollection, response, err := vpcService.ListVPNGatewayMembers(listVPNGatewayMembersOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayMemberCollection).ToNot(BeNil())
		})
	})

	Describe(`GetVPNGatewayMember - Retrieve a VPN gateway member`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNGatewayMember(getVPNGatewayMemberOptions *GetVPNGatewayMemberOptions)`, func() {
			getVPNGatewayMemberOptions := &vpcv1.GetVPNGatewayMemberOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayMember, response, err := vpcService.GetVPNGatewayMember(getVPNGatewayMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayMember).ToNot(BeNil())
		})
	})

	Describe(`ReplaceVPNGatewayMember - Replace a VPN gateway member`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceVPNGatewayMember(replaceVPNGatewayMemberOptions *ReplaceVPNGatewayMemberOptions)`, func() {
			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			vpnGatewayMemberPrivateIPPrototypeModel := &vpcv1.VPNGatewayMemberPrivateIPPrototypeReservedIPPrototypeVPNGatewayContext{
				Subnet: subnetIdentityModel,
			}

			replaceVPNGatewayMemberOptions := &vpcv1.ReplaceVPNGatewayMemberOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
				PrivateIP:    vpnGatewayMemberPrivateIPPrototypeModel,
			}

			vpnGatewayMember, response, err := vpcService.ReplaceVPNGatewayMember(replaceVPNGatewayMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(vpnGatewayMember).ToNot(BeNil())
		})
	})

	Describe(`ListVPNGatewayServiceConnections - List VPN gateway service connections`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNGatewayServiceConnections(listVPNGatewayServiceConnectionsOptions *ListVPNGatewayServiceConnectionsOptions) with pagination`, func() {
			listVPNGatewayServiceConnectionsOptions := &vpcv1.ListVPNGatewayServiceConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
				Start:        core.StringPtr("testString"),
			}

			listVPNGatewayServiceConnectionsOptions.Start = nil
			listVPNGatewayServiceConnectionsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNGatewayServiceConnection
			for {
				vpnGatewayServiceConnectionCollection, response, err := vpcService.ListVPNGatewayServiceConnections(listVPNGatewayServiceConnectionsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnGatewayServiceConnectionCollection).ToNot(BeNil())
				allResults = append(allResults, vpnGatewayServiceConnectionCollection.ServiceConnections...)

				nextToken, nextErr := vpnGatewayServiceConnectionCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNGatewayServiceConnectionsOptions.Start = nextToken

				if listVPNGatewayServiceConnectionsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNGatewayServiceConnections(listVPNGatewayServiceConnectionsOptions *ListVPNGatewayServiceConnectionsOptions) using VPNGatewayServiceConnectionsPager`, func() {
			listVPNGatewayServiceConnectionsOptions := &vpcv1.ListVPNGatewayServiceConnectionsOptions{
				VPNGatewayID: core.StringPtr("testString"),
				Limit:        core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNGatewayServiceConnectionsPager(listVPNGatewayServiceConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNGatewayServiceConnection
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNGatewayServiceConnectionsPager(listVPNGatewayServiceConnectionsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNGatewayServiceConnections() returned a total of %d item(s) using VPNGatewayServiceConnectionsPager.\n", len(allResults))
		})
	})

	Describe(`GetVPNGatewayServiceConnection - Retrieve a VPN gateway service connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNGatewayServiceConnection(getVPNGatewayServiceConnectionOptions *GetVPNGatewayServiceConnectionOptions)`, func() {
			getVPNGatewayServiceConnectionOptions := &vpcv1.GetVPNGatewayServiceConnectionOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			vpnGatewayServiceConnection, response, err := vpcService.GetVPNGatewayServiceConnection(getVPNGatewayServiceConnectionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnGatewayServiceConnection).ToNot(BeNil())
		})
	})

	Describe(`ListVPNServers - List VPN servers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNServers(listVPNServersOptions *ListVPNServersOptions) with pagination`, func() {
			listVPNServersOptions := &vpcv1.ListVPNServersOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
				Start:           core.StringPtr("testString"),
			}

			listVPNServersOptions.Start = nil
			listVPNServersOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNServer
			for {
				vpnServerCollection, response, err := vpcService.ListVPNServers(listVPNServersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnServerCollection).ToNot(BeNil())
				allResults = append(allResults, vpnServerCollection.VPNServers...)

				nextToken, nextErr := vpnServerCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNServersOptions.Start = nextToken

				if listVPNServersOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNServers(listVPNServersOptions *ListVPNServersOptions) using VPNServersPager`, func() {
			listVPNServersOptions := &vpcv1.ListVPNServersOptions{
				Limit:           core.Int64Ptr(int64(10)),
				Name:            core.StringPtr("my-name"),
				ResourceGroupID: core.StringPtr("testString"),
				Sort:            core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNServersPager(listVPNServersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNServer
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNServersPager(listVPNServersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNServers() returned a total of %d item(s) using VPNServersPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPNServer - Create a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPNServer(createVPNServerOptions *CreateVPNServerOptions)`, func() {
			vpnServerPrototypeCertificateModel := &vpcv1.VPNServerPrototypeCertificateCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationByCertificatePrototypeClientCaModel := &vpcv1.VPNServerAuthenticationByCertificatePrototypeClientCaCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationPrototypeModel := &vpcv1.VPNServerAuthenticationPrototypeVPNServerAuthenticationByCertificatePrototype{
				Method:   core.StringPtr("certificate"),
				ClientCa: vpnServerAuthenticationByCertificatePrototypeClientCaModel,
				Crl:      core.StringPtr("testString"),
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			ipModel := &vpcv1.IP{
				Address: core.StringPtr("192.168.3.4"),
			}

			resourceGroupIdentityModel := &vpcv1.ResourceGroupIdentityByID{
				ID: core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345"),
			}

			securityGroupIdentityModel := &vpcv1.SecurityGroupIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::security-group:r006-be5df5ca-12a0-494b-907e-aa6ec2bfa271"),
			}

			createVPNServerOptions := &vpcv1.CreateVPNServerOptions{
				Certificate:          vpnServerPrototypeCertificateModel,
				ClientAuthentication: []vpcv1.VPNServerAuthenticationPrototypeIntf{vpnServerAuthenticationPrototypeModel},
				ClientIPPool:         core.StringPtr("172.16.0.0/16"),
				Subnets:              []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
				ClientDnsServerIps:   []vpcv1.IP{*ipModel},
				ClientIdleTimeout:    core.Int64Ptr(int64(600)),
				EnableSplitTunneling: core.BoolPtr(false),
				Name:                 core.StringPtr("my-vpn-server"),
				Port:                 core.Int64Ptr(int64(443)),
				Protocol:             core.StringPtr("udp"),
				ResourceGroup:        resourceGroupIdentityModel,
				SecurityGroups:       []vpcv1.SecurityGroupIdentityIntf{securityGroupIdentityModel},
			}

			vpnServer, response, err := vpcService.CreateVPNServer(createVPNServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnServer).ToNot(BeNil())
		})
	})

	Describe(`GetVPNServer - Retrieve a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNServer(getVPNServerOptions *GetVPNServerOptions)`, func() {
			getVPNServerOptions := &vpcv1.GetVPNServerOptions{
				ID: core.StringPtr("testString"),
			}

			vpnServer, response, err := vpcService.GetVPNServer(getVPNServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServer).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPNServer - Update a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPNServer(updateVPNServerOptions *UpdateVPNServerOptions)`, func() {
			vpnServerPatchCertificateModel := &vpcv1.VPNServerPatchCertificateCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationByCertificatePrototypeClientCaModel := &vpcv1.VPNServerAuthenticationByCertificatePrototypeClientCaCertificateInstanceIdentityCertificateInstanceIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:secrets-manager:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:36fa422d-080d-4d83-8d2d-86851b4001df:secret:2e786aab-42fa-63ed-14f8-d66d552f4dd5"),
			}

			vpnServerAuthenticationPrototypeModel := &vpcv1.VPNServerAuthenticationPrototypeVPNServerAuthenticationByCertificatePrototype{
				Method:   core.StringPtr("certificate"),
				ClientCa: vpnServerAuthenticationByCertificatePrototypeClientCaModel,
				Crl:      core.StringPtr("testString"),
			}

			ipModel := &vpcv1.IP{
				Address: core.StringPtr("192.168.3.4"),
			}

			subnetIdentityModel := &vpcv1.SubnetIdentityByCRN{
				CRN: core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"),
			}

			vpnServerPatchModel := &vpcv1.VPNServerPatch{
				Certificate:          vpnServerPatchCertificateModel,
				ClientAuthentication: []vpcv1.VPNServerAuthenticationPrototypeIntf{vpnServerAuthenticationPrototypeModel},
				ClientDnsServerIps:   []vpcv1.IP{*ipModel},
				ClientIdleTimeout:    core.Int64Ptr(int64(600)),
				ClientIPPool:         core.StringPtr("172.16.0.0/16"),
				EnableSplitTunneling: core.BoolPtr(true),
				Name:                 core.StringPtr("my-vpn-server-updated"),
				Port:                 core.Int64Ptr(int64(443)),
				Protocol:             core.StringPtr("tcp"),
				Subnets:              []vpcv1.SubnetIdentityIntf{subnetIdentityModel},
			}
			vpnServerPatchModelAsPatch, asPatchErr := vpnServerPatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNServerOptions := &vpcv1.UpdateVPNServerOptions{
				ID:             core.StringPtr("testString"),
				VPNServerPatch: vpnServerPatchModelAsPatch,
				IfMatch:        core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			vpnServer, response, err := vpcService.UpdateVPNServer(updateVPNServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServer).ToNot(BeNil())
		})
	})

	Describe(`GetVPNServerClientConfiguration - Retrieve client configuration`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNServerClientConfiguration(getVPNServerClientConfigurationOptions *GetVPNServerClientConfigurationOptions)`, func() {
			getVPNServerClientConfigurationOptions := &vpcv1.GetVPNServerClientConfigurationOptions{
				ID: core.StringPtr("testString"),
			}

			result, response, err := vpcService.GetVPNServerClientConfiguration(getVPNServerClientConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})
	})

	Describe(`ListVPNServerClients - List VPN clients for a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNServerClients(listVPNServerClientsOptions *ListVPNServerClientsOptions) with pagination`, func() {
			listVPNServerClientsOptions := &vpcv1.ListVPNServerClientsOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("created_at"),
				Start:       core.StringPtr("testString"),
			}

			listVPNServerClientsOptions.Start = nil
			listVPNServerClientsOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNServerClient
			for {
				vpnServerClientCollection, response, err := vpcService.ListVPNServerClients(listVPNServerClientsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnServerClientCollection).ToNot(BeNil())
				allResults = append(allResults, vpnServerClientCollection.Clients...)

				nextToken, nextErr := vpnServerClientCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNServerClientsOptions.Start = nextToken

				if listVPNServerClientsOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNServerClients(listVPNServerClientsOptions *ListVPNServerClientsOptions) using VPNServerClientsPager`, func() {
			listVPNServerClientsOptions := &vpcv1.ListVPNServerClientsOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("created_at"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNServerClientsPager(listVPNServerClientsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNServerClient
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNServerClientsPager(listVPNServerClientsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNServerClients() returned a total of %d item(s) using VPNServerClientsPager.\n", len(allResults))
		})
	})

	Describe(`GetVPNServerClient - Retrieve a VPN client`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNServerClient(getVPNServerClientOptions *GetVPNServerClientOptions)`, func() {
			getVPNServerClientOptions := &vpcv1.GetVPNServerClientOptions{
				ID:          core.StringPtr("testString"),
				VPNServerID: core.StringPtr("testString"),
			}

			vpnServerClient, response, err := vpcService.GetVPNServerClient(getVPNServerClientOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerClient).ToNot(BeNil())
		})
	})

	Describe(`DisconnectVPNClient - Disconnect a VPN client`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DisconnectVPNClient(disconnectVPNClientOptions *DisconnectVPNClientOptions)`, func() {
			disconnectVPNClientOptions := &vpcv1.DisconnectVPNClientOptions{
				ID:          core.StringPtr("testString"),
				VPNServerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DisconnectVPNClient(disconnectVPNClientOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`ListVPNServerRoutes - List VPN routes for a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVPNServerRoutes(listVPNServerRoutesOptions *ListVPNServerRoutesOptions) with pagination`, func() {
			listVPNServerRoutesOptions := &vpcv1.ListVPNServerRoutesOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("name"),
				Start:       core.StringPtr("testString"),
			}

			listVPNServerRoutesOptions.Start = nil
			listVPNServerRoutesOptions.Limit = core.Int64Ptr(1)

			var allResults []vpcv1.VPNServerRoute
			for {
				vpnServerRouteCollection, response, err := vpcService.ListVPNServerRoutes(listVPNServerRoutesOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(vpnServerRouteCollection).ToNot(BeNil())
				allResults = append(allResults, vpnServerRouteCollection.Routes...)

				nextToken, nextErr := vpnServerRouteCollection.GetNextStart()
				Expect(nextErr).To(BeNil())
				listVPNServerRoutesOptions.Start = nextToken

				if listVPNServerRoutesOptions.Start == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVPNServerRoutes(listVPNServerRoutesOptions *ListVPNServerRoutesOptions) using VPNServerRoutesPager`, func() {
			listVPNServerRoutesOptions := &vpcv1.ListVPNServerRoutesOptions{
				VPNServerID: core.StringPtr("testString"),
				Limit:       core.Int64Ptr(int64(10)),
				Sort:        core.StringPtr("name"),
			}

			// Test GetNext().
			pager, err := vpcService.NewVPNServerRoutesPager(listVPNServerRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []vpcv1.VPNServerRoute
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = vpcService.NewVPNServerRoutesPager(listVPNServerRoutesOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVPNServerRoutes() returned a total of %d item(s) using VPNServerRoutesPager.\n", len(allResults))
		})
	})

	Describe(`CreateVPNServerRoute - Create a VPN route for a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVPNServerRoute(createVPNServerRouteOptions *CreateVPNServerRouteOptions)`, func() {
			createVPNServerRouteOptions := &vpcv1.CreateVPNServerRouteOptions{
				VPNServerID: core.StringPtr("testString"),
				Destination: core.StringPtr("172.16.0.0/16"),
				Action:      core.StringPtr("deliver"),
				Name:        core.StringPtr("my-vpn-route-2"),
			}

			vpnServerRoute, response, err := vpcService.CreateVPNServerRoute(createVPNServerRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
	})

	Describe(`GetVPNServerRoute - Retrieve a VPN route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVPNServerRoute(getVPNServerRouteOptions *GetVPNServerRouteOptions)`, func() {
			getVPNServerRouteOptions := &vpcv1.GetVPNServerRouteOptions{
				ID:          core.StringPtr("testString"),
				VPNServerID: core.StringPtr("testString"),
			}

			vpnServerRoute, response, err := vpcService.GetVPNServerRoute(getVPNServerRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
	})

	Describe(`UpdateVPNServerRoute - Update a VPN route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UpdateVPNServerRoute(updateVPNServerRouteOptions *UpdateVPNServerRouteOptions)`, func() {
			vpnServerRoutePatchModel := &vpcv1.VPNServerRoutePatch{
				Name: core.StringPtr("my-vpn-server-route-updated"),
			}
			vpnServerRoutePatchModelAsPatch, asPatchErr := vpnServerRoutePatchModel.AsPatch()
			Expect(asPatchErr).To(BeNil())

			updateVPNServerRouteOptions := &vpcv1.UpdateVPNServerRouteOptions{
				ID:                  core.StringPtr("testString"),
				VPNServerID:         core.StringPtr("testString"),
				VPNServerRoutePatch: vpnServerRoutePatchModelAsPatch,
			}

			vpnServerRoute, response, err := vpcService.UpdateVPNServerRoute(updateVPNServerRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(vpnServerRoute).ToNot(BeNil())
		})
	})

	Describe(`DeleteBackupPolicyPlan - Delete a backup policy plan`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteBackupPolicyPlan(deleteBackupPolicyPlanOptions *DeleteBackupPolicyPlanOptions)`, func() {
			deleteBackupPolicyPlanOptions := &vpcv1.DeleteBackupPolicyPlanOptions{
				BackupPolicyID: core.StringPtr("testString"),
				ID:             core.StringPtr("testString"),
				IfMatch:        core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			backupPolicyPlan, response, err := vpcService.DeleteBackupPolicyPlan(deleteBackupPolicyPlanOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(backupPolicyPlan).ToNot(BeNil())
		})
	})

	Describe(`DeleteBackupPolicy - Delete a backup policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteBackupPolicy(deleteBackupPolicyOptions *DeleteBackupPolicyOptions)`, func() {
			deleteBackupPolicyOptions := &vpcv1.DeleteBackupPolicyOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			backupPolicy, response, err := vpcService.DeleteBackupPolicy(deleteBackupPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(backupPolicy).ToNot(BeNil())
		})
	})

	Describe(`DeleteBareMetalServerNetworkAttachment - Delete a bare metal server network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteBareMetalServerNetworkAttachment(deleteBareMetalServerNetworkAttachmentOptions *DeleteBareMetalServerNetworkAttachmentOptions)`, func() {
			deleteBareMetalServerNetworkAttachmentOptions := &vpcv1.DeleteBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteBareMetalServerNetworkAttachment(deleteBareMetalServerNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteBareMetalServerNetworkInterface - Delete a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteBareMetalServerNetworkInterface(deleteBareMetalServerNetworkInterfaceOptions *DeleteBareMetalServerNetworkInterfaceOptions)`, func() {
			deleteBareMetalServerNetworkInterfaceOptions := &vpcv1.DeleteBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteBareMetalServerNetworkInterface(deleteBareMetalServerNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RemoveBareMetalServerNetworkInterfaceFloatingIP - Disassociate a floating IP from a bare metal server network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveBareMetalServerNetworkInterfaceFloatingIP(removeBareMetalServerNetworkInterfaceFloatingIPOptions *RemoveBareMetalServerNetworkInterfaceFloatingIPOptions)`, func() {
			removeBareMetalServerNetworkInterfaceFloatingIPOptions := &vpcv1.RemoveBareMetalServerNetworkInterfaceFloatingIPOptions{
				BareMetalServerID:  core.StringPtr("testString"),
				ID:                 core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveBareMetalServerNetworkInterfaceFloatingIP(removeBareMetalServerNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteBareMetalServer - Delete a bare metal server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteBareMetalServer(deleteBareMetalServerOptions *DeleteBareMetalServerOptions)`, func() {
			deleteBareMetalServerOptions := &vpcv1.DeleteBareMetalServerOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteBareMetalServer(deleteBareMetalServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteClusterNetworkInterface - Delete a cluster network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteClusterNetworkInterface(deleteClusterNetworkInterfaceOptions *DeleteClusterNetworkInterfaceOptions)`, func() {
			deleteClusterNetworkInterfaceOptions := &vpcv1.DeleteClusterNetworkInterfaceOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				ID:               core.StringPtr("testString"),
				IfMatch:          core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkInterface, response, err := vpcService.DeleteClusterNetworkInterface(deleteClusterNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`DeleteClusterNetworkSubnetReservedIP - Delete a cluster network subnet reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteClusterNetworkSubnetReservedIP(deleteClusterNetworkSubnetReservedIPOptions *DeleteClusterNetworkSubnetReservedIPOptions)`, func() {
			deleteClusterNetworkSubnetReservedIPOptions := &vpcv1.DeleteClusterNetworkSubnetReservedIPOptions{
				ClusterNetworkID:       core.StringPtr("testString"),
				ClusterNetworkSubnetID: core.StringPtr("testString"),
				ID:                     core.StringPtr("testString"),
				IfMatch:                core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkSubnetReservedIP, response, err := vpcService.DeleteClusterNetworkSubnetReservedIP(deleteClusterNetworkSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkSubnetReservedIP).ToNot(BeNil())
		})
	})

	Describe(`DeleteClusterNetworkSubnet - Delete a cluster network subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteClusterNetworkSubnet(deleteClusterNetworkSubnetOptions *DeleteClusterNetworkSubnetOptions)`, func() {
			deleteClusterNetworkSubnetOptions := &vpcv1.DeleteClusterNetworkSubnetOptions{
				ClusterNetworkID: core.StringPtr("testString"),
				ID:               core.StringPtr("testString"),
				IfMatch:          core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetworkSubnet, response, err := vpcService.DeleteClusterNetworkSubnet(deleteClusterNetworkSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetworkSubnet).ToNot(BeNil())
		})
	})

	Describe(`DeleteClusterNetwork - Delete a cluster network`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteClusterNetwork(deleteClusterNetworkOptions *DeleteClusterNetworkOptions)`, func() {
			deleteClusterNetworkOptions := &vpcv1.DeleteClusterNetworkOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			clusterNetwork, response, err := vpcService.DeleteClusterNetwork(deleteClusterNetworkOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(clusterNetwork).ToNot(BeNil())
		})
	})

	Describe(`DeleteDedicatedHostGroup - Delete a dedicated host group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteDedicatedHostGroup(deleteDedicatedHostGroupOptions *DeleteDedicatedHostGroupOptions)`, func() {
			deleteDedicatedHostGroupOptions := &vpcv1.DeleteDedicatedHostGroupOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteDedicatedHostGroup(deleteDedicatedHostGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteDedicatedHost - Delete a dedicated host`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteDedicatedHost(deleteDedicatedHostOptions *DeleteDedicatedHostOptions)`, func() {
			deleteDedicatedHostOptions := &vpcv1.DeleteDedicatedHostOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteDedicatedHost(deleteDedicatedHostOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RemoveEndpointGatewayIP - Unbind a reserved IP from an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveEndpointGatewayIP(removeEndpointGatewayIPOptions *RemoveEndpointGatewayIPOptions)`, func() {
			removeEndpointGatewayIPOptions := &vpcv1.RemoveEndpointGatewayIPOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveEndpointGatewayIP(removeEndpointGatewayIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteEndpointGatewayResourceBinding - Delete a resource binding from an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteEndpointGatewayResourceBinding(deleteEndpointGatewayResourceBindingOptions *DeleteEndpointGatewayResourceBindingOptions)`, func() {
			deleteEndpointGatewayResourceBindingOptions := &vpcv1.DeleteEndpointGatewayResourceBindingOptions{
				EndpointGatewayID: core.StringPtr("testString"),
				ID:                core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteEndpointGatewayResourceBinding(deleteEndpointGatewayResourceBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteEndpointGateway - Delete an endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteEndpointGateway(deleteEndpointGatewayOptions *DeleteEndpointGatewayOptions)`, func() {
			deleteEndpointGatewayOptions := &vpcv1.DeleteEndpointGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteEndpointGateway(deleteEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteFloatingIP - Delete a floating IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteFloatingIP(deleteFloatingIPOptions *DeleteFloatingIPOptions)`, func() {
			deleteFloatingIPOptions := &vpcv1.DeleteFloatingIPOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteFloatingIP(deleteFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteFlowLogCollector - Delete a flow log collector`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteFlowLogCollector(deleteFlowLogCollectorOptions *DeleteFlowLogCollectorOptions)`, func() {
			deleteFlowLogCollectorOptions := &vpcv1.DeleteFlowLogCollectorOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteFlowLogCollector(deleteFlowLogCollectorOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteImage - Delete an image`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteImage(deleteImageOptions *DeleteImageOptions)`, func() {
			deleteImageOptions := &vpcv1.DeleteImageOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteImage(deleteImageOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteImageExportJob - Delete an image export job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteImageExportJob(deleteImageExportJobOptions *DeleteImageExportJobOptions)`, func() {
			deleteImageExportJobOptions := &vpcv1.DeleteImageExportJobOptions{
				ID:      core.StringPtr("testString"),
				ImageID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteImageExportJob(deleteImageExportJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteInstanceGroup - Delete an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroup(deleteInstanceGroupOptions *DeleteInstanceGroupOptions)`, func() {
			deleteInstanceGroupOptions := &vpcv1.DeleteInstanceGroupOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroup(deleteInstanceGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupLoadBalancer - Delete an instance group load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupLoadBalancer(deleteInstanceGroupLoadBalancerOptions *DeleteInstanceGroupLoadBalancerOptions)`, func() {
			deleteInstanceGroupLoadBalancerOptions := &vpcv1.DeleteInstanceGroupLoadBalancerOptions{
				InstanceGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupLoadBalancer(deleteInstanceGroupLoadBalancerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupManager - Delete an instance group manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupManager(deleteInstanceGroupManagerOptions *DeleteInstanceGroupManagerOptions)`, func() {
			deleteInstanceGroupManagerOptions := &vpcv1.DeleteInstanceGroupManagerOptions{
				ID:              core.StringPtr("testString"),
				InstanceGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupManager(deleteInstanceGroupManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupManagerAction - Delete specified instance group manager action`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupManagerAction(deleteInstanceGroupManagerActionOptions *DeleteInstanceGroupManagerActionOptions)`, func() {
			deleteInstanceGroupManagerActionOptions := &vpcv1.DeleteInstanceGroupManagerActionOptions{
				ID:                     core.StringPtr("testString"),
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupManagerAction(deleteInstanceGroupManagerActionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupManagerPolicy - Delete an instance group manager policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupManagerPolicy(deleteInstanceGroupManagerPolicyOptions *DeleteInstanceGroupManagerPolicyOptions)`, func() {
			deleteInstanceGroupManagerPolicyOptions := &vpcv1.DeleteInstanceGroupManagerPolicyOptions{
				ID:                     core.StringPtr("testString"),
				InstanceGroupID:        core.StringPtr("testString"),
				InstanceGroupManagerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupManagerPolicy(deleteInstanceGroupManagerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupMemberships - Delete memberships from an instance group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupMemberships(deleteInstanceGroupMembershipsOptions *DeleteInstanceGroupMembershipsOptions)`, func() {
			deleteInstanceGroupMembershipsOptions := &vpcv1.DeleteInstanceGroupMembershipsOptions{
				InstanceGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupMemberships(deleteInstanceGroupMembershipsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceGroupMembership - Delete an instance group membership`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceGroupMembership(deleteInstanceGroupMembershipOptions *DeleteInstanceGroupMembershipOptions)`, func() {
			deleteInstanceGroupMembershipOptions := &vpcv1.DeleteInstanceGroupMembershipOptions{
				ID:              core.StringPtr("testString"),
				InstanceGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceGroupMembership(deleteInstanceGroupMembershipOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceTemplate - Delete an instance template`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceTemplate(deleteInstanceTemplateOptions *DeleteInstanceTemplateOptions)`, func() {
			deleteInstanceTemplateOptions := &vpcv1.DeleteInstanceTemplateOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceTemplate(deleteInstanceTemplateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstance - Delete an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstance(deleteInstanceOptions *DeleteInstanceOptions)`, func() {
			deleteInstanceOptions := &vpcv1.DeleteInstanceOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteInstance(deleteInstanceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceClusterNetworkAttachment - Delete an instance cluster network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceClusterNetworkAttachment(deleteInstanceClusterNetworkAttachmentOptions *DeleteInstanceClusterNetworkAttachmentOptions)`, func() {
			deleteInstanceClusterNetworkAttachmentOptions := &vpcv1.DeleteInstanceClusterNetworkAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			instanceClusterNetworkAttachment, response, err := vpcService.DeleteInstanceClusterNetworkAttachment(deleteInstanceClusterNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(instanceClusterNetworkAttachment).ToNot(BeNil())
		})
	})

	Describe(`DeleteInstanceNetworkAttachment - Delete an instance network attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceNetworkAttachment(deleteInstanceNetworkAttachmentOptions *DeleteInstanceNetworkAttachmentOptions)`, func() {
			deleteInstanceNetworkAttachmentOptions := &vpcv1.DeleteInstanceNetworkAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceNetworkAttachment(deleteInstanceNetworkAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteInstanceNetworkInterface - Delete an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceNetworkInterface(deleteInstanceNetworkInterfaceOptions *DeleteInstanceNetworkInterfaceOptions)`, func() {
			deleteInstanceNetworkInterfaceOptions := &vpcv1.DeleteInstanceNetworkInterfaceOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceNetworkInterface(deleteInstanceNetworkInterfaceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RemoveInstanceNetworkInterfaceFloatingIP - Disassociate a floating IP from an instance network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveInstanceNetworkInterfaceFloatingIP(removeInstanceNetworkInterfaceFloatingIPOptions *RemoveInstanceNetworkInterfaceFloatingIPOptions)`, func() {
			removeInstanceNetworkInterfaceFloatingIPOptions := &vpcv1.RemoveInstanceNetworkInterfaceFloatingIPOptions{
				ID:                 core.StringPtr("testString"),
				InstanceID:         core.StringPtr("testString"),
				NetworkInterfaceID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveInstanceNetworkInterfaceFloatingIP(removeInstanceNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteInstanceVolumeAttachment - Delete a volume attachment`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteInstanceVolumeAttachment(deleteInstanceVolumeAttachmentOptions *DeleteInstanceVolumeAttachmentOptions)`, func() {
			deleteInstanceVolumeAttachmentOptions := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
				ID:         core.StringPtr("testString"),
				InstanceID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteInstanceVolumeAttachment(deleteInstanceVolumeAttachmentOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteKey - Delete a key`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteKey(deleteKeyOptions *DeleteKeyOptions)`, func() {
			deleteKeyOptions := &vpcv1.DeleteKeyOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteKey(deleteKeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteLoadBalancer - Delete a load balancer`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancer(deleteLoadBalancerOptions *DeleteLoadBalancerOptions)`, func() {
			deleteLoadBalancerOptions := &vpcv1.DeleteLoadBalancerOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteLoadBalancer(deleteLoadBalancerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteLoadBalancerListener - Delete a load balancer listener`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions *DeleteLoadBalancerListenerOptions)`, func() {
			deleteLoadBalancerListenerOptions := &vpcv1.DeleteLoadBalancerListenerOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteLoadBalancerListenerPolicy - Delete a load balancer listener policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancerListenerPolicy(deleteLoadBalancerListenerPolicyOptions *DeleteLoadBalancerListenerPolicyOptions)`, func() {
			deleteLoadBalancerListenerPolicyOptions := &vpcv1.DeleteLoadBalancerListenerPolicyOptions{
				ID:             core.StringPtr("testString"),
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteLoadBalancerListenerPolicy(deleteLoadBalancerListenerPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteLoadBalancerListenerPolicyRule - Delete a load balancer listener policy rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancerListenerPolicyRule(deleteLoadBalancerListenerPolicyRuleOptions *DeleteLoadBalancerListenerPolicyRuleOptions)`, func() {
			deleteLoadBalancerListenerPolicyRuleOptions := &vpcv1.DeleteLoadBalancerListenerPolicyRuleOptions{
				ID:             core.StringPtr("testString"),
				ListenerID:     core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PolicyID:       core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteLoadBalancerListenerPolicyRule(deleteLoadBalancerListenerPolicyRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteLoadBalancerPool - Delete a load balancer pool`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions *DeleteLoadBalancerPoolOptions)`, func() {
			deleteLoadBalancerPoolOptions := &vpcv1.DeleteLoadBalancerPoolOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteLoadBalancerPoolMember - Delete a load balancer pool member`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteLoadBalancerPoolMember(deleteLoadBalancerPoolMemberOptions *DeleteLoadBalancerPoolMemberOptions)`, func() {
			deleteLoadBalancerPoolMemberOptions := &vpcv1.DeleteLoadBalancerPoolMemberOptions{
				ID:             core.StringPtr("testString"),
				LoadBalancerID: core.StringPtr("testString"),
				PoolID:         core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteLoadBalancerPoolMember(deleteLoadBalancerPoolMemberOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteNetworkACL - Delete a network ACL`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteNetworkACL(deleteNetworkACLOptions *DeleteNetworkACLOptions)`, func() {
			deleteNetworkACLOptions := &vpcv1.DeleteNetworkACLOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteNetworkACL(deleteNetworkACLOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteNetworkACLRule - Delete a network ACL rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteNetworkACLRule(deleteNetworkACLRuleOptions *DeleteNetworkACLRuleOptions)`, func() {
			deleteNetworkACLRuleOptions := &vpcv1.DeleteNetworkACLRuleOptions{
				ID:           core.StringPtr("testString"),
				NetworkACLID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteNetworkACLRule(deleteNetworkACLRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeletePlacementGroup - Delete a placement group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeletePlacementGroup(deletePlacementGroupOptions *DeletePlacementGroupOptions)`, func() {
			deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeletePlacementGroup(deletePlacementGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeletePrivatePathServiceGateway - Delete a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeletePrivatePathServiceGateway(deletePrivatePathServiceGatewayOptions *DeletePrivatePathServiceGatewayOptions)`, func() {
			deletePrivatePathServiceGatewayOptions := &vpcv1.DeletePrivatePathServiceGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeletePrivatePathServiceGateway(deletePrivatePathServiceGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeletePrivatePathServiceGatewayAccountPolicy - Delete an account policy for a private path service gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeletePrivatePathServiceGatewayAccountPolicy(deletePrivatePathServiceGatewayAccountPolicyOptions *DeletePrivatePathServiceGatewayAccountPolicyOptions)`, func() {
			deletePrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.DeletePrivatePathServiceGatewayAccountPolicyOptions{
				ID:                          core.StringPtr("testString"),
				PrivatePathServiceGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeletePrivatePathServiceGatewayAccountPolicy(deletePrivatePathServiceGatewayAccountPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeletePublicAddressRange - Delete a public address range`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeletePublicAddressRange(deletePublicAddressRangeOptions *DeletePublicAddressRangeOptions)`, func() {
			deletePublicAddressRangeOptions := &vpcv1.DeletePublicAddressRangeOptions{
				ID: core.StringPtr("testString"),
			}

			publicAddressRange, response, err := vpcService.DeletePublicAddressRange(deletePublicAddressRangeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(publicAddressRange).ToNot(BeNil())
		})
	})

	Describe(`DeletePublicGateway - Delete a public gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeletePublicGateway(deletePublicGatewayOptions *DeletePublicGatewayOptions)`, func() {
			deletePublicGatewayOptions := &vpcv1.DeletePublicGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeletePublicGateway(deletePublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteReservation - Delete a reservation`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteReservation(deleteReservationOptions *DeleteReservationOptions)`, func() {
			deleteReservationOptions := &vpcv1.DeleteReservationOptions{
				ID: core.StringPtr("testString"),
			}

			reservation, response, err := vpcService.DeleteReservation(deleteReservationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(reservation).ToNot(BeNil())
		})
	})

	Describe(`DeleteSecurityGroup - Delete a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSecurityGroup(deleteSecurityGroupOptions *DeleteSecurityGroupOptions)`, func() {
			deleteSecurityGroupOptions := &vpcv1.DeleteSecurityGroupOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSecurityGroup(deleteSecurityGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteSecurityGroupRule - Delete a security group rule`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSecurityGroupRule(deleteSecurityGroupRuleOptions *DeleteSecurityGroupRuleOptions)`, func() {
			deleteSecurityGroupRuleOptions := &vpcv1.DeleteSecurityGroupRuleOptions{
				ID:              core.StringPtr("testString"),
				SecurityGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSecurityGroupRule(deleteSecurityGroupRuleOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteSecurityGroupTargetBinding - Remove a target from a security group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions *DeleteSecurityGroupTargetBindingOptions)`, func() {
			deleteSecurityGroupTargetBindingOptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
				ID:              core.StringPtr("testString"),
				SecurityGroupID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteShare - Delete a file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteShare(deleteShareOptions *DeleteShareOptions)`, func() {
			deleteShareOptions := &vpcv1.DeleteShareOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/96d225c4-56bd-43d9-98fc-d7148e5c5028"),
			}

			share, response, err := vpcService.DeleteShare(deleteShareOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(share).ToNot(BeNil())
		})
	})

	Describe(`DeleteShareAccessorBinding - Delete a file share accessor binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteShareAccessorBinding(deleteShareAccessorBindingOptions *DeleteShareAccessorBindingOptions)`, func() {
			deleteShareAccessorBindingOptions := &vpcv1.DeleteShareAccessorBindingOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteShareAccessorBinding(deleteShareAccessorBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteShareMountTarget - Delete a file share mount target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteShareMountTarget(deleteShareMountTargetOptions *DeleteShareMountTargetOptions)`, func() {
			deleteShareMountTargetOptions := &vpcv1.DeleteShareMountTargetOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			shareMountTarget, response, err := vpcService.DeleteShareMountTarget(deleteShareMountTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(shareMountTarget).ToNot(BeNil())
		})
	})

	Describe(`DeleteShareSnapshot - Delete a share snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteShareSnapshot(deleteShareSnapshotOptions *DeleteShareSnapshotOptions)`, func() {
			deleteShareSnapshotOptions := &vpcv1.DeleteShareSnapshotOptions{
				ID:      core.StringPtr("testString"),
				ShareID: core.StringPtr("testString"),
			}

			shareSnapshot, response, err := vpcService.DeleteShareSnapshot(deleteShareSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(shareSnapshot).ToNot(BeNil())
		})
	})

	Describe(`DeleteShareSource - Split the source file share from a replica file share`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteShareSource(deleteShareSourceOptions *DeleteShareSourceOptions)`, func() {
			deleteShareSourceOptions := &vpcv1.DeleteShareSourceOptions{
				ShareID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteShareSource(deleteShareSourceOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteSnapshotConsistencyGroup - Delete a snapshot consistency group`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSnapshotConsistencyGroup(deleteSnapshotConsistencyGroupOptions *DeleteSnapshotConsistencyGroupOptions)`, func() {
			deleteSnapshotConsistencyGroupOptions := &vpcv1.DeleteSnapshotConsistencyGroupOptions{
				ID: core.StringPtr("testString"),
			}

			snapshotConsistencyGroup, response, err := vpcService.DeleteSnapshotConsistencyGroup(deleteSnapshotConsistencyGroupOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(snapshotConsistencyGroup).ToNot(BeNil())
		})
	})

	Describe(`DeleteSnapshots - Delete a filtered collection of snapshots`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSnapshots(deleteSnapshotsOptions *DeleteSnapshotsOptions)`, func() {
			deleteSnapshotsOptions := &vpcv1.DeleteSnapshotsOptions{
				SourceVolumeID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSnapshots(deleteSnapshotsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteSnapshot - Delete a snapshot`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSnapshot(deleteSnapshotOptions *DeleteSnapshotOptions)`, func() {
			deleteSnapshotOptions := &vpcv1.DeleteSnapshotOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteSnapshot(deleteSnapshotOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteSnapshotClone - Delete a snapshot clone`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSnapshotClone(deleteSnapshotCloneOptions *DeleteSnapshotCloneOptions)`, func() {
			deleteSnapshotCloneOptions := &vpcv1.DeleteSnapshotCloneOptions{
				ID:       core.StringPtr("testString"),
				ZoneName: core.StringPtr("us-south-1"),
			}

			response, err := vpcService.DeleteSnapshotClone(deleteSnapshotCloneOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteSubnet - Delete a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSubnet(deleteSubnetOptions *DeleteSubnetOptions)`, func() {
			deleteSubnetOptions := &vpcv1.DeleteSubnetOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSubnet(deleteSubnetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`UnsetSubnetPublicGateway - Detach a public gateway from a subnet`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`UnsetSubnetPublicGateway(unsetSubnetPublicGatewayOptions *UnsetSubnetPublicGatewayOptions)`, func() {
			unsetSubnetPublicGatewayOptions := &vpcv1.UnsetSubnetPublicGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.UnsetSubnetPublicGateway(unsetSubnetPublicGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteSubnetReservedIP - Delete a reserved IP`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteSubnetReservedIP(deleteSubnetReservedIPOptions *DeleteSubnetReservedIPOptions)`, func() {
			deleteSubnetReservedIPOptions := &vpcv1.DeleteSubnetReservedIPOptions{
				ID:       core.StringPtr("testString"),
				SubnetID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteSubnetReservedIP(deleteSubnetReservedIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVirtualNetworkInterfaces - Delete a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVirtualNetworkInterfaces(deleteVirtualNetworkInterfacesOptions *DeleteVirtualNetworkInterfacesOptions)`, func() {
			deleteVirtualNetworkInterfacesOptions := &vpcv1.DeleteVirtualNetworkInterfacesOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			virtualNetworkInterface, response, err := vpcService.DeleteVirtualNetworkInterfaces(deleteVirtualNetworkInterfacesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(virtualNetworkInterface).ToNot(BeNil())
		})
	})

	Describe(`RemoveNetworkInterfaceFloatingIP - Disassociate a floating IP from a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveNetworkInterfaceFloatingIP(removeNetworkInterfaceFloatingIPOptions *RemoveNetworkInterfaceFloatingIPOptions)`, func() {
			removeNetworkInterfaceFloatingIPOptions := &vpcv1.RemoveNetworkInterfaceFloatingIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveNetworkInterfaceFloatingIP(removeNetworkInterfaceFloatingIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RemoveVirtualNetworkInterfaceIP - Unbind a reserved IP from a virtual network interface`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions *RemoveVirtualNetworkInterfaceIPOptions)`, func() {
			removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{
				ID:                        core.StringPtr("testString"),
				VirtualNetworkInterfaceID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVolume - Delete a volume`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVolume(deleteVolumeOptions *DeleteVolumeOptions)`, func() {
			deleteVolumeOptions := &vpcv1.DeleteVolumeOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteVolume(deleteVolumeOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVolumeJob - Delete a volume job`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVolumeJob(deleteVolumeJobOptions *DeleteVolumeJobOptions)`, func() {
			deleteVolumeJobOptions := &vpcv1.DeleteVolumeJobOptions{
				ID:       core.StringPtr("testString"),
				VolumeID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVolumeJob(deleteVolumeJobOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteVPC - Delete a VPC`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPC(deleteVPCOptions *DeleteVPCOptions)`, func() {
			deleteVPCOptions := &vpcv1.DeleteVPCOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteVPC(deleteVPCOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPCAddressPrefix - Delete an address prefix`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPCAddressPrefix(deleteVPCAddressPrefixOptions *DeleteVPCAddressPrefixOptions)`, func() {
			deleteVPCAddressPrefixOptions := &vpcv1.DeleteVPCAddressPrefixOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPCAddressPrefix(deleteVPCAddressPrefixOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPCDnsResolutionBinding - Delete a DNS resolution binding`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPCDnsResolutionBinding(deleteVPCDnsResolutionBindingOptions *DeleteVPCDnsResolutionBindingOptions)`, func() {
			deleteVPCDnsResolutionBindingOptions := &vpcv1.DeleteVPCDnsResolutionBindingOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			vpcdnsResolutionBinding, response, err := vpcService.DeleteVPCDnsResolutionBinding(deleteVPCDnsResolutionBindingOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(vpcdnsResolutionBinding).ToNot(BeNil())
		})
	})

	Describe(`DeleteVPCRoute - Delete a VPC route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPCRoute(deleteVPCRouteOptions *DeleteVPCRouteOptions)`, func() {
			deleteVPCRouteOptions := &vpcv1.DeleteVPCRouteOptions{
				ID:    core.StringPtr("testString"),
				VPCID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPCRoute(deleteVPCRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPCRoutingTable - Delete a VPC routing table`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPCRoutingTable(deleteVPCRoutingTableOptions *DeleteVPCRoutingTableOptions)`, func() {
			deleteVPCRoutingTableOptions := &vpcv1.DeleteVPCRoutingTableOptions{
				ID:      core.StringPtr("testString"),
				VPCID:   core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteVPCRoutingTable(deleteVPCRoutingTableOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPCRoutingTableRoute - Delete a VPC routing table route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPCRoutingTableRoute(deleteVPCRoutingTableRouteOptions *DeleteVPCRoutingTableRouteOptions)`, func() {
			deleteVPCRoutingTableRouteOptions := &vpcv1.DeleteVPCRoutingTableRouteOptions{
				ID:             core.StringPtr("testString"),
				RoutingTableID: core.StringPtr("testString"),
				VPCID:          core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPCRoutingTableRoute(deleteVPCRoutingTableRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteIkePolicy - Delete an IKE policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteIkePolicy(deleteIkePolicyOptions *DeleteIkePolicyOptions)`, func() {
			deleteIkePolicyOptions := &vpcv1.DeleteIkePolicyOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteIkePolicy(deleteIkePolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteIpsecPolicy - Delete an IPsec policy`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteIpsecPolicy(deleteIpsecPolicyOptions *DeleteIpsecPolicyOptions)`, func() {
			deleteIpsecPolicyOptions := &vpcv1.DeleteIpsecPolicyOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteIpsecPolicy(deleteIpsecPolicyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPNGateway - Delete a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPNGateway(deleteVPNGatewayOptions *DeleteVPNGatewayOptions)`, func() {
			deleteVPNGatewayOptions := &vpcv1.DeleteVPNGatewayOptions{
				ID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPNGateway(deleteVPNGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`RemoveVPNGatewayAdvertisedCIDR - Remove an advertised CIDR from a VPN gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveVPNGatewayAdvertisedCIDR(removeVPNGatewayAdvertisedCIDROptions *RemoveVPNGatewayAdvertisedCIDROptions)`, func() {
			removeVPNGatewayAdvertisedCIDROptions := &vpcv1.RemoveVPNGatewayAdvertisedCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveVPNGatewayAdvertisedCIDR(removeVPNGatewayAdvertisedCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPNGatewayConnection - Delete a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPNGatewayConnection(deleteVPNGatewayConnectionOptions *DeleteVPNGatewayConnectionOptions)`, func() {
			deleteVPNGatewayConnectionOptions := &vpcv1.DeleteVPNGatewayConnectionOptions{
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
				IfMatch:      core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteVPNGatewayConnection(deleteVPNGatewayConnectionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`RemoveVPNGatewayConnectionsLocalCIDR - Remove a local CIDR from a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveVPNGatewayConnectionsLocalCIDR(removeVPNGatewayConnectionsLocalCIDROptions *RemoveVPNGatewayConnectionsLocalCIDROptions)`, func() {
			removeVPNGatewayConnectionsLocalCIDROptions := &vpcv1.RemoveVPNGatewayConnectionsLocalCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveVPNGatewayConnectionsLocalCIDR(removeVPNGatewayConnectionsLocalCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`RemoveVPNGatewayConnectionsPeerCIDR - Remove a peer CIDR from a VPN gateway connection`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`RemoveVPNGatewayConnectionsPeerCIDR(removeVPNGatewayConnectionsPeerCIDROptions *RemoveVPNGatewayConnectionsPeerCIDROptions)`, func() {
			removeVPNGatewayConnectionsPeerCIDROptions := &vpcv1.RemoveVPNGatewayConnectionsPeerCIDROptions{
				CIDR:         core.StringPtr("192.168.1.0/24"),
				ID:           core.StringPtr("testString"),
				VPNGatewayID: core.StringPtr("testString"),
			}

			response, err := vpcService.RemoveVPNGatewayConnectionsPeerCIDR(removeVPNGatewayConnectionsPeerCIDROptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteVPNServer - Delete a VPN server`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPNServer(deleteVPNServerOptions *DeleteVPNServerOptions)`, func() {
			deleteVPNServerOptions := &vpcv1.DeleteVPNServerOptions{
				ID:      core.StringPtr("testString"),
				IfMatch: core.StringPtr("W/\"96d225c4-56bd-43d9-98fc-d7148e5c5028\""),
			}

			response, err := vpcService.DeleteVPNServer(deleteVPNServerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteVPNServerClient - Delete a VPN client`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPNServerClient(deleteVPNServerClientOptions *DeleteVPNServerClientOptions)`, func() {
			deleteVPNServerClientOptions := &vpcv1.DeleteVPNServerClientOptions{
				ID:          core.StringPtr("testString"),
				VPNServerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPNServerClient(deleteVPNServerClientOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})

	Describe(`DeleteVPNServerRoute - Delete a VPN route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVPNServerRoute(deleteVPNServerRouteOptions *DeleteVPNServerRouteOptions)`, func() {
			deleteVPNServerRouteOptions := &vpcv1.DeleteVPNServerRouteOptions{
				ID:          core.StringPtr("testString"),
				VPNServerID: core.StringPtr("testString"),
			}

			response, err := vpcService.DeleteVPNServerRoute(deleteVPNServerRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
		})
	})
})

//
// Utility functions are declared in the unit test file
//
