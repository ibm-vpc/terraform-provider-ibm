// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceNetworkInterfaceFloatingIpAvailable = "available"
	isInstanceNetworkInterfaceFloatingIpDeleting  = "deleting"
	isInstanceNetworkInterfaceFloatingIpPending   = "pending"
	isInstanceNetworkInterfaceFloatingIpDeleted   = "deleted"
	isInstanceNetworkInterfaceFloatingIpFailed    = "failed"
	isInstanceNetworkInterface                    = "network_interface"
	isInstanceNetworkInterfaceFloatingIPID        = "floating_ip"
)

func ResourceIBMIsInstanceNetworkInterfaceFloatingIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISInstanceNetworkInterfaceFloatingIpCreate,
		ReadContext:   resourceIBMISInstanceNetworkInterfaceFloatingIpRead,
		UpdateContext: resourceIBMISInstanceNetworkInterfaceFloatingIpUpdate,
		DeleteContext: resourceIBMISInstanceNetworkInterfaceFloatingIpDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			isInstanceNetworkInterface: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance network interface identifier",
			},
			isInstanceNetworkInterfaceFloatingIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating ip identifier of the network interface associated with the Instance",
			},
			floatingIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the floating IP",
			},

			floatingIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP address",
			},

			floatingIPStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP status",
			},

			floatingIPZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			floatingIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target info",
			},

			floatingIPCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP crn",
			},
			isFloatingIPResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			isFloatingIPAuthorizedCIDR: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The public address range authorized CIDR.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isFloatingIPAuthorizedCIDRID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this public address range.",
						},
						isFloatingIPAuthorizedCIDRCIDR: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR for this public address range.",
						},
						isFloatingIPAuthorizedCIDRHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this public address range.",
						},
						isFloatingIPAuthorizedCIDRName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this public address range.",
						},
						isFloatingIPAuthorizedCIDRResType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			isFloatingIPProfile: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile for this floating IP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isFloatingIPProfileName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this floating IP profile.",
						},
						isFloatingIPProfileHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this floating IP profile.",
						},
						isFloatingIPResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceId := d.Get(isInstanceID).(string)
	instanceNicId := ""
	nicId := d.Get(isInstanceNetworkInterface).(string)
	if strings.Contains(nicId, "/") {
		_, instanceNicId, err = ParseNICTerraformID(nicId)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "create", "sep-id-parts").GetDiag()
		}
	} else {
		instanceNicId = nicId
	}

	instanceNicFipId := d.Get(isInstanceNetworkInterfaceFloatingIPID).(string)

	options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &instanceNicId,
		ID:                 &instanceNicFipId,
	}

	fip, _, err := sess.AddInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("AddInstanceNetworkInterfaceFloatingIPWithContext failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(MakeTerraformNICFipID(instanceId, instanceNicId, *fip.ID))
	diagErr := instanceNICFipGet(context, d, fip, instanceId, instanceNicId)
	if diagErr != nil {
		return diagErr
	}

	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceId, nicID, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "sep-id-parts").GetDiag()
	}

	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	options := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicID,
		ID:                 &fipId,
	}

	fip, response, err := sess.GetInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceNetworkInterfaceFloatingIPWithContext failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	diagErr := instanceNICFipGet(context, d, fip, instanceId, nicID)
	if diagErr != nil {
		return diagErr
	}
	return nil
}

func instanceNICFipGet(context context.Context, d *schema.ResourceData, fip *vpcv1.FloatingIP, instanceId, nicId string) diag.Diagnostics {
	var err error
	d.SetId(MakeTerraformNICFipID(instanceId, nicId, *fip.ID))
	if err = d.Set(floatingIPName, *fip.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-name").GetDiag()
	}
	if err = d.Set(floatingIPAddress, *fip.Address); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-address").GetDiag()
	}

	if err = d.Set(floatingIPStatus, fip.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-status").GetDiag()
	}

	if err = d.Set(floatingIPZone, *fip.Zone.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-zone").GetDiag()
	}

	if err = d.Set(floatingIPCRN, *fip.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-crn").GetDiag()
	}
	target, ok := fip.Target.(*vpcv1.FloatingIPTarget)
	if ok {
		if err = d.Set(floatingIPTarget, target.ID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-target").GetDiag()
		}
	}
	if fip.ResourceType != nil {
		if err = d.Set(isFloatingIPResourceType, *fip.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-resource_type").GetDiag()
		}
	}
	if fip.AuthorizedCIDR != nil {
		authorizedCIDRMap := map[string]interface{}{}
		if fip.AuthorizedCIDR.ID != nil {
			authorizedCIDRMap[isFloatingIPAuthorizedCIDRID] = *fip.AuthorizedCIDR.ID
		}
		if fip.AuthorizedCIDR.CIDR != nil {
			authorizedCIDRMap[isFloatingIPAuthorizedCIDRCIDR] = *fip.AuthorizedCIDR.CIDR
		}
		if fip.AuthorizedCIDR.Href != nil {
			authorizedCIDRMap[isFloatingIPAuthorizedCIDRHref] = *fip.AuthorizedCIDR.Href
		}
		if fip.AuthorizedCIDR.Name != nil {
			authorizedCIDRMap[isFloatingIPAuthorizedCIDRName] = *fip.AuthorizedCIDR.Name
		}
		if fip.AuthorizedCIDR.ResourceType != nil {
			authorizedCIDRMap[isFloatingIPAuthorizedCIDRResType] = *fip.AuthorizedCIDR.ResourceType
		}
		if err = d.Set(isFloatingIPAuthorizedCIDR, []map[string]interface{}{authorizedCIDRMap}); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-authorized_cidr").GetDiag()
		}
	}
	if fip.Profile != nil {
		profileMap := map[string]interface{}{}
		if fip.Profile.Name != nil {
			profileMap[isFloatingIPProfileName] = *fip.Profile.Name
		}
		if fip.Profile.Href != nil {
			profileMap[isFloatingIPProfileHref] = *fip.Profile.Href
		}
		if fip.Profile.ResourceType != nil {
			profileMap[isFloatingIPResourceType] = *fip.Profile.ResourceType
		}
		if err = d.Set(isFloatingIPProfile, []map[string]interface{}{profileMap}); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "read", "set-profile").GetDiag()
		}
	}
	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange(isInstanceNetworkInterfaceFloatingIPID) {
		instanceId, nicId, _, err := ParseNICFipTerraformID(d.Id())
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "update", "sep-id-parts").GetDiag()
		}
		sess, err := vpcClient(meta)
		if err != nil {
			tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "update", "initialize-client")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		floatingIpId := ""
		if fipOk, ok := d.GetOk(isInstanceNetworkInterfaceFloatingIPID); ok {
			floatingIpId = fipOk.(string)
		}
		options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &floatingIpId,
		}

		fip, _, err := sess.AddInstanceNetworkInterfaceFloatingIPWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("AddInstanceNetworkInterfaceFloatingIPWithContext failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(MakeTerraformNICFipID(instanceId, nicId, *fip.ID))
		return instanceNICFipGet(context, d, fip, instanceId, nicId)
	}
	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceId, nicId, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "delete", "sep-id-parts").GetDiag()
	}

	diagErr := instanceNetworkInterfaceFipDelete(context, d, meta, instanceId, nicId, fipId)
	if diagErr != nil {
		return diagErr
	}

	return nil
}

func instanceNetworkInterfaceFipDelete(context context.Context, d *schema.ResourceData, meta interface{}, instanceId, nicId, fipId string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_network_interface_floating_ip", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getBmsNicFipOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	fip, response, err := sess.GetInstanceNetworkInterfaceFloatingIPWithContext(context, getBmsNicFipOptions)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceNetworkInterfaceFloatingIPWithContext failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &vpcv1.RemoveInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	response, err = sess.RemoveInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("RemoveInstanceNetworkInterfaceFloatingIPWithContext failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	_, err = isWaitForInstanceNetworkInterfaceFloatingIpDeleted(sess, instanceId, nicId, fipId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceNetworkInterfaceFloatingIpDeleted failed: %s", err.Error()), "ibm_is_instance_network_interface_floating_ip", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId("")
	return nil
}

func isWaitForInstanceNetworkInterfaceFloatingIpDeleted(instanceC *vpcv1.VpcV1, instanceId, nicId, fipId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for (%s) / (%s) / (%s) to be deleted.", instanceId, nicId, fipId)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceNetworkInterfaceFloatingIpAvailable, isInstanceNetworkInterfaceFloatingIpDeleting, isInstanceNetworkInterfaceFloatingIpPending},
		Target:     []string{isInstanceNetworkInterfaceFloatingIpDeleted, isInstanceNetworkInterfaceFloatingIpFailed, ""},
		Refresh:    isInstanceNetworkInterfaceFloatingIpDeleteRefreshFunc(instanceC, instanceId, nicId, fipId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isInstanceNetworkInterfaceFloatingIpDeleteRefreshFunc(instanceC *vpcv1.VpcV1, instanceId, nicId, fipId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getBmsNicFloatingIpOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := instanceC.GetInstanceNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return fip, isInstanceNetworkInterfaceFloatingIpDeleted, nil
			}
			return fip, isInstanceNetworkInterfaceFloatingIpFailed, fmt.Errorf("[ERROR] Error getting Instance(%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
		}
		return fip, isInstanceNetworkInterfaceFloatingIpDeleting, err
	}
}

func isWaitForInstanceNetworkInterfaceFloatingIpAvailable(client *vpcv1.VpcV1, instanceId, nicId, fipId string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Instance (%s) Network Interface (%s) to be available.", instanceId, nicId)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceNetworkInterfaceFloatingIpPending},
		Target:     []string{isInstanceNetworkInterfaceFloatingIpAvailable, isInstanceNetworkInterfaceFloatingIpFailed},
		Refresh:    isInstanceNetworkInterfaceFloatingIpRefreshFunc(client, instanceId, nicId, fipId, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isInstanceNetworkInterfaceFloatingIpRefreshFunc(client *vpcv1.VpcV1, instanceId, nicId, fipId string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getBmsNicFloatingIpOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := client.GetInstanceNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Instance (%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
		}
		status := ""

		status = *fip.Status
		d.Set(floatingIPStatus, *fip.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if status == "available" || status == "failed" {
			close(communicator)
			return fip, status, nil

		}

		return fip, "pending", nil
	}
}
