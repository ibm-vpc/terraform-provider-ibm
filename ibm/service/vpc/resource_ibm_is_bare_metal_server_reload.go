// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIsBareMetalServerReload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerReloadCreate,
		ReadContext:   resourceIBMISBareMetalServerReloadRead,
		DeleteContext: resourceIBMISBareMetalServerReloadDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bare metal server identifier",
			},
			"user_data": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Bare metal server user data to reload",
			},
		},
	}
}

func resourceIBMISBareMetalServerReloadCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var bareMetalServerId, userdata string
	if bmsId, ok := d.GetOk(isBareMetalServerID); ok {
		bareMetalServerId = bmsId.(string)
	}
	if userdataOk, ok := d.GetOk("user_data"); ok {
		userdata = userdataOk.(string)
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	stopServerIfStartingForInitialization := false
	options := &vpcv1.GetBareMetalServerInitializationOptions{
		ID: &bareMetalServerId,
	}
	stopServerIfStartingForInitialization, err = resourceStopServerIfRunning(bareMetalServerId, "hard", d, context, sess, stopServerIfStartingForInitialization)
	if err != nil {
		return diag.FromErr(err)
	}
	init, response, err := sess.GetBareMetalServerInitializationWithContext(context, options)
	if err != nil || init == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error get bare metal server initialization (%s) err %s\n%s", bareMetalServerId, err, response))
	}
	d.SetId(bareMetalServerId)
	keys := make([]vpcv1.KeyIdentityIntf, 0)
	for _, key := range init.Keys {
		k := &vpcv1.KeyIdentity{
			ID: key.ID,
		}
		keys = append(keys, k)
	}
	reloadoptions := &vpcv1.ReplaceBareMetalServerInitializationOptions{
		ID: &bareMetalServerId,
		Image: &vpcv1.ImageIdentityByID{
			ID: init.Image.ID,
		},
		Keys:     keys,
		UserData: &userdata,
	}
	initReload, response, err := sess.ReplaceBareMetalServerInitializationWithContext(context, reloadoptions)
	if err != nil || initReload == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error reloading bare metal server (%s) err %s\n%s", bareMetalServerId, err, response))
	}
	if stopServerIfStartingForInitialization {
		stopServerIfStartingForInitialization, err = resourceStartServerIfStopped(bareMetalServerId, "hard", d, context, sess, stopServerIfStartingForInitialization)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = BareMetalServerReloadGet(d, sess, bareMetalServerId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func BareMetalServerReloadGet(d *schema.ResourceData, sess *vpcv1.VpcV1, bareMetalServerId string) error {

	options := &vpcv1.GetBareMetalServerInitializationOptions{
		ID: &bareMetalServerId,
	}
	init, response, err := sess.GetBareMetalServerInitialization(options)
	if err != nil || init == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error fetching bare metal server (%s)  initialization err %s\n%s", bareMetalServerId, err, response)
	}

	d.Set(isBareMetalServerID, bareMetalServerId)
	return nil
}

func resourceIBMISBareMetalServerReloadRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var bareMetalServerId string
	if bmsId, ok := d.GetOk(isBareMetalServerID); ok {
		bareMetalServerId = bmsId.(string)
	}
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = BareMetalServerReloadGet(d, sess, bareMetalServerId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
func resourceIBMISBareMetalServerReloadDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
