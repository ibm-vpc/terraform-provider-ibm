// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsVPNServerClientConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServerClientConfigurationRead,

		Schema: map[string]*schema.Schema{
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPN server identifier.",
			},
			"file_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The VPN server identifier.",
			},
			"vpn_server_client_configuration": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPN client identifier.",
			},
		},
	}
}

func dataSourceIBMIsVPNServerClientConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("dataSourceIBMIsVPNServerClientConfigurationRead")
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPNServerClientConfigurationOptions := &vpcv1.GetVPNServerClientConfigurationOptions{}
	getVPNServerClientConfigurationOptions.SetID(d.Get("vpn_server").(string))

	result, response, err := sess.GetVPNServerClientConfigurationWithContext(context, getVPNServerClientConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVPNServerClientWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVPNServerClientWithContext failed %s\n%s", err, response))
	}

	d.SetId(d.Get("vpn_server").(string))
	log.Println("result")
	log.Println(*result)
	configStr := *result
	configStr = strings.Trim(configStr, "\n")
	configStr = strings.Trim(configStr, `"`)
	configStr = strings.Replace(configStr, `\n`, "\n", -1)

	if v, ok := d.GetOk("file_path"); ok {
		fileName := v.(string)
		// fileName := cliContext.String(vpn_server.OutputFileFlag)
		f, err := os.Create(fileName)
		if err == nil {
			_, err = f.WriteString(configStr)
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error Saving VPNServerClientConfiguration Result: %s", err))

		}
		log.Printf("OpenVPN client configuration was saved to {{.%s}}", fileName)
	}

	if err = d.Set("vpn_server_client_configuration", *result); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting VPNServerClientConfiguration Result: %s", err))
	}
	return nil
}
