// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVPNGatewayAdvertisedCidr           = "cidr"
	isVPNGatewayAdvertisedCidrVPNGateway = "vpn_gateway"
	isVPNGatewayAdvertisedCidrDeleting   = "deleting"
	isVPNGatewayAdvertisedCidrDeleted    = "done"
)

func ResourceIBMISVPNGatewayAdvertisedCidr() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPNGatewayAdvertisedCidrCreate,
		Read:     resourceIBMISVPNGatewayAdvertisedCidrRead,
		Delete:   resourceIBMISVPNGatewayAdvertisedCidrDelete,
		Exists:   resourceIBMISVPNGatewayAdvertisedCidrExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			isVPNGatewayAdvertisedCidrVPNGateway: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier",
			},
			isVPNGatewayAdvertisedCidr: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address range in CIDR block notation.",
			},
		},
	}
}

func resourceIBMISVPNGatewayAdvertisedCidrCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Adding advertised cidr to vpn gateway")

	gatewayID := d.Get(isVPNGatewayAdvertisedCidrVPNGateway).(string)
	cidr := d.Get(isVPNGatewayAdvertisedCidr).(string)

	options := &vpcv1.AddVPNGatewayAdvertisedCIDROptions{
		VPNGatewayID: &gatewayID,
		CIDR:         &cidr,
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	response, err := sess.AddVPNGatewayAdvertisedCIDR(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error adding advertised cidr to VPN Gateway err %s\n%s", err, response)
	}

	return nil
}

func resourceIBMISVPNGatewayAdvertisedCidrRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	if len(parts) != 2 {
		return fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of gID/gAdvertisedCidr", d.Id())
	}

	gID := parts[0]
	gAdvertisedCidr := parts[1]

	checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
		VPNGatewayID: &gID,
		CIDR:         &gAdvertisedCidr,
	}
	response, err := sess.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting advertised cidr : %s\n%s", err, response)
	}
	return nil
}

func resourceIBMISVPNGatewayAdvertisedCidrDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gID := parts[0]
	cidr := parts[1]

	err = vpngwAdvertisedCidrDelete(d, meta, gID, cidr)
	if err != nil {
		return err
	}
	return nil
}

func vpngwAdvertisedCidrDelete(d *schema.ResourceData, meta interface{}, gID, gCidr string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
		VPNGatewayID: &gID,
		CIDR:         &gCidr,
	}
	response, err := sess.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Vpn Gateway advertised cidr(%s): %s\n%s", gCidr, err, response)
	}

	removeVPNGatewayAdvertisedCIDROptions := &vpcv1.RemoveVPNGatewayAdvertisedCIDROptions{
		VPNGatewayID: &gID,
		CIDR:         &gCidr,
	}
	response, err = sess.RemoveVPNGatewayAdvertisedCIDR(removeVPNGatewayAdvertisedCIDROptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error removing advertised cidr from Vpn Gateway: %s\n%s", err, response)
	}

	_, err = isWaitForVPNGatewayAdvertisedCIDRDeleted(sess, gID, gCidr, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for Vpn Gateway advertised cidr (%s) is deleted: %s", gCidr, err)
	}

	d.SetId("")
	return nil
}

func isWaitForVPNGatewayAdvertisedCIDRDeleted(vpnGatewayAdverisedCidr *vpcv1.VpcV1, gID, gCidr string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for VPNGatewayAdvertisedCIDR (%s) to be deleted.", gCidr)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayAdvertisedCidrDeleting},
		Target:     []string{"", isVPNGatewayAdvertisedCidrDeleted},
		Refresh:    isVPNGatewayAdvertisedCIDRDeleteRefreshFunc(vpnGatewayAdverisedCidr, gID, gCidr),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVPNGatewayAdvertisedCIDRDeleteRefreshFunc(vpnGatewayAdverisedCidr *vpcv1.VpcV1, gID, gCidr string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
			VPNGatewayID: &gID,
			CIDR:         &gCidr,
		}
		response, err := vpnGatewayAdverisedCidr.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return "", isVPNGatewayConnectionDeleted, nil
			}
			return "", "", fmt.Errorf("[ERROR] The VPNGateway Advertised CIDR %s failed to delete: %s\n%s", gCidr, err, response)
		}
		return nil, isVPNGatewayConnectionDeleting, nil
	}
}

func resourceIBMISVPNGatewayAdvertisedCidrExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of gID/gadvertisedCidr", d.Id())
	}

	gID := parts[0]
	gAdvertisedCidr := parts[1]
	exists, err := vpngwadvertisedCidrExists(d, meta, gID, gAdvertisedCidr)
	return exists, err
}

func vpngwadvertisedCidrExists(d *schema.ResourceData, meta interface{}, gID, gAdvertisedCidr string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
		VPNGatewayID: &gID,
		CIDR:         &gAdvertisedCidr,
	}
	response, err := sess.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error checking Vpn Gateway advertised cidr: %s\n%s", err, response)
	}
	return true, nil
}
