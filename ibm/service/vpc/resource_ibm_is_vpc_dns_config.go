package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMISVPCDnsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISVPCDnsConfigCreate,
		ReadContext:   resourceIBMISVPCDnsConfigRead,
		UpdateContext: resourceIBMISVPCDnsConfigUpdate,
		DeleteContext: resourceIBMISVPCDnsConfigDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The VPC identifier",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpc_dns_config", "vpc_id"),
			},
			"enable_hub": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this VPC is enabled as a DNS name resolution hub",
			},
			"resolver_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "system",
				Description:  "The type of DNS resolver (system, manual, delegated)",
				ValidateFunc: validation.StringInSlice([]string{"system", "manual", "delegated"}, false),
			},
			"resolver_vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The VPC ID to delegate DNS resolution to (required when resolver_type is delegated)",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpc_dns_config", "resolver_vpc_id"),
			},
			"resolver_vpc_crn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC CRN to delegate DNS resolution to (alternative to resolver_vpc_id)",
			},
			"dns_binding_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name for the DNS resolution binding (auto-generated if not specified)",
			},
			"manual_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Manual DNS servers (required when resolver_type is manual)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the DNS server",
						},
						"zone_affinity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The zone affinity for this DNS server",
						},
					},
				},
			},
			"resolution_binding_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of DNS resolution bindings for this VPC",
			},
			"dns_resolution_binding_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DNS resolution binding created for delegated resolver",
			},
		},
	}
}

func ResourceIBMISVPCDnsConfigValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "vpc_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "resolver_vpc_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
	)
	resourceValidator := validate.ResourceValidator{
		ResourceName: "ibm_is_vpc_dns_config",
		Schema:       validateSchema,
	}
	return &resourceValidator
}

func resourceIBMISVPCDnsConfigCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpcID := d.Get("vpc_id").(string)

	// Validate configuration
	if err := validateDnsConfig(d); err != nil {
		return diag.FromErr(err)
	}

	// Create DNS resolution binding first if resolver type is delegated
	var bindingID string
	var bindingName string

	resolverType := d.Get("resolver_type").(string)
	if resolverType == "delegated" {
		binding, err := createDnsResolutionBinding(sess, d)
		if err != nil {
			return diag.FromErr(err)
		}
		bindingID = *binding.ID
		bindingName = *binding.Name
	}

	// Update VPC DNS configuration
	if err := updateVPCDnsConfiguration(sess, d, bindingID, bindingName); err != nil {
		// Rollback: delete the binding if VPC update fails
		if bindingID != "" {
			deleteBindingOptions := &vpcv1.DeleteVPCDnsResolutionBindingOptions{
				VPCID: &vpcID,
				ID:    &bindingID,
			}
			sess.DeleteVPCDnsResolutionBinding(deleteBindingOptions)
		}
		return diag.FromErr(err)
	}

	// Set resource ID and computed fields
	d.SetId(vpcID)
	if bindingID != "" {
		d.Set("dns_resolution_binding_id", bindingID)
	}
	if bindingName != "" {
		d.Set("dns_binding_name", bindingName)
	}

	return resourceIBMISVPCDnsConfigRead(context, d, meta)
}

func resourceIBMISVPCDnsConfigRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpcID := d.Id()

	getVpcOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}

	vpc, response, err := sess.GetVPCWithContext(context, getVpcOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error getting VPC: %s\n%s", err, response))
	}

	// Read DNS configuration
	if vpc.Dns != nil {
		d.Set("enable_hub", vpc.Dns.EnableHub)

		if vpc.Dns.ResolutionBindingCount != nil {
			d.Set("resolution_binding_count", *vpc.Dns.ResolutionBindingCount)
		}

		if vpc.Dns.Resolver != nil {
			vpcdnsResolver := vpc.Dns.Resolver.(*vpcv1.VpcdnsResolver)
			d.Set("resolver_type", vpcdnsResolver.Type)

			// Handle manual servers
			if vpcdnsResolver.ManualServers != nil && len(vpcdnsResolver.ManualServers) > 0 {
				manualServers := make([]map[string]interface{}, 0)
				for _, server := range vpcdnsResolver.ManualServers {
					serverMap := map[string]interface{}{
						"address": *server.Address,
					}
					if server.ZoneAffinity != nil {
						serverMap["zone_affinity"] = *server.ZoneAffinity.Name
					}
					manualServers = append(manualServers, serverMap)
				}
				d.Set("manual_servers", manualServers)
			}

			// Handle delegated resolver VPC
			if vpcdnsResolver.VPC != nil {
				if vpcdnsResolver.VPC.ID != nil {
					d.Set("resolver_vpc_id", *vpcdnsResolver.VPC.ID)
				}
				if vpcdnsResolver.VPC.CRN != nil {
					d.Set("resolver_vpc_crn", *vpcdnsResolver.VPC.CRN)
				}
			}
		}
	}

	// If delegated, try to find the binding
	if d.Get("resolver_type").(string) == "delegated" {
		bindingID := d.Get("dns_resolution_binding_id").(string)
		if bindingID != "" {
			listBindingsOptions := &vpcv1.ListVPCDnsResolutionBindingsOptions{
				VPCID: &vpcID,
			}
			bindings, _, err := sess.ListVPCDnsResolutionBindingsWithContext(context, listBindingsOptions)
			if err == nil && bindings != nil {
				for _, binding := range bindings.DnsResolutionBindings {
					if *binding.ID == bindingID {
						d.Set("dns_binding_name", *binding.Name)
						break
					}
				}
			}
		}
	}

	return nil
}

func resourceIBMISVPCDnsConfigUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpcID := d.Id()
	hasChange := false

	// Validate new configuration
	if err := validateDnsConfig(d); err != nil {
		return diag.FromErr(err)
	}

	// Check if resolver type changed
	if d.HasChange("resolver_type") {
		oldType, newType := d.GetChange("resolver_type")

		// If changing from delegated, delete old binding
		if oldType.(string) == "delegated" {
			oldBindingID := d.Get("dns_resolution_binding_id").(string)
			if oldBindingID != "" {
				deleteBindingOptions := &vpcv1.DeleteVPCDnsResolutionBindingOptions{
					VPCID: &vpcID,
					ID:    &oldBindingID,
				}
				_, _, err := sess.DeleteVPCDnsResolutionBindingWithContext(context, deleteBindingOptions)
				if err != nil {
					log.Printf("[WARN] Error deleting old DNS resolution binding: %s", err)
				}
				d.Set("dns_resolution_binding_id", "")
				d.Set("dns_binding_name", "")
			}
		}

		// If changing to delegated, create new binding
		if newType.(string) == "delegated" {
			binding, err := createDnsResolutionBinding(sess, d)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set("dns_resolution_binding_id", *binding.ID)
			d.Set("dns_binding_name", *binding.Name)
		}

		hasChange = true
	}

	// Check for other changes
	if d.HasChange("enable_hub") || d.HasChange("manual_servers") ||
		d.HasChange("resolver_vpc_id") || d.HasChange("resolver_vpc_crn") ||
		d.HasChange("dns_binding_name") {
		hasChange = true
	}

	if hasChange {
		bindingID := d.Get("dns_resolution_binding_id").(string)
		bindingName := d.Get("dns_binding_name").(string)

		if err := updateVPCDnsConfiguration(sess, d, bindingID, bindingName); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMISVPCDnsConfigRead(context, d, meta)
}

func resourceIBMISVPCDnsConfigDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpcID := d.Id()

	// Delete DNS resolution binding if exists
	bindingID := d.Get("dns_resolution_binding_id").(string)
	if bindingID != "" {
		deleteBindingOptions := &vpcv1.DeleteVPCDnsResolutionBindingOptions{
			VPCID: &vpcID,
			ID:    &bindingID,
		}
		_, _, err := sess.DeleteVPCDnsResolutionBindingWithContext(context, deleteBindingOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error deleting DNS resolution binding: %s", err))
		}
	}

	// Reset VPC DNS to system defaults
	vpcDnsPatch := &vpcv1.VpcdnsPatch{
		EnableHub: core.BoolPtr(false),
		Resolver: &vpcv1.VpcdnsResolverPatch{
			Type: core.StringPtr("system"),
		},
	}

	vpcPatch := &vpcv1.VPCPatch{
		Dns: vpcDnsPatch,
	}

	vpcPatchAsPatch, err := vpcPatch.AsPatch()
	if err != nil {
		return diag.FromErr(err)
	}

	updateVpcOptions := &vpcv1.UpdateVPCOptions{
		ID:       &vpcID,
		VPCPatch: vpcPatchAsPatch,
	}

	_, _, err = sess.UpdateVPCWithContext(context, updateVpcOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error resetting VPC DNS configuration: %s", err))
	}

	d.SetId("")
	return nil
}

// Helper functions

func validateDnsConfig(d *schema.ResourceData) error {
	resolverType := d.Get("resolver_type").(string)

	switch resolverType {
	case "delegated":
		resolverVpcID := d.Get("resolver_vpc_id").(string)
		resolverVpcCRN := d.Get("resolver_vpc_crn").(string)

		if resolverVpcID == "" && resolverVpcCRN == "" {
			return fmt.Errorf("either resolver_vpc_id or resolver_vpc_crn must be specified when resolver_type is 'delegated'")
		}

		enableHub := d.Get("enable_hub").(bool)
		if enableHub {
			return fmt.Errorf("enable_hub must be false when resolver_type is 'delegated'")
		}

	case "manual":
		manualServers := d.Get("manual_servers").([]interface{})
		if len(manualServers) == 0 {
			return fmt.Errorf("manual_servers must be specified when resolver_type is 'manual'")
		}

	case "system":
		// No additional validation needed
	}

	return nil
}

func createDnsResolutionBinding(sess *vpcv1.VpcV1, d *schema.ResourceData) (*vpcv1.VpcdnsResolutionBinding, error) {
	vpcID := d.Get("vpc_id").(string)

	// Determine which VPC identifier to use
	var vpcIdentity vpcv1.VPCIdentityIntf

	if resolverVpcID := d.Get("resolver_vpc_id").(string); resolverVpcID != "" {
		vpcIdentity = &vpcv1.VPCIdentityByID{
			ID: core.StringPtr(resolverVpcID),
		}
	} else if resolverVpcCRN := d.Get("resolver_vpc_crn").(string); resolverVpcCRN != "" {
		vpcIdentity = &vpcv1.VPCIdentityByCRN{
			CRN: core.StringPtr(resolverVpcCRN),
		}
	}

	createOptions := &vpcv1.CreateVPCDnsResolutionBindingOptions{
		VPCID: &vpcID,
		VPC:   vpcIdentity,
	}

	// Add optional binding name if specified
	if bindingName := d.Get("dns_binding_name").(string); bindingName != "" {
		createOptions.Name = core.StringPtr(bindingName)
	}

	binding, _, err := sess.CreateVPCDnsResolutionBinding(createOptions)
	if err != nil {
		return nil, fmt.Errorf("Error creating DNS resolution binding: %s", err)
	}

	return binding, nil
}

func updateVPCDnsConfiguration(sess *vpcv1.VpcV1, d *schema.ResourceData, bindingID, bindingName string) error {
	vpcID := d.Get("vpc_id").(string)
	resolverType := d.Get("resolver_type").(string)

	// Build DNS resolver patch
	resolverPatch := &vpcv1.VpcdnsResolverPatch{
		Type: core.StringPtr(resolverType),
	}

	// Handle different resolver types
	switch resolverType {
	case "manual":
		manualServersRaw := d.Get("manual_servers").([]interface{})
		manualServers := make([]vpcv1.DnsServerPrototype, 0, len(manualServersRaw))

		for _, serverRaw := range manualServersRaw {
			serverMap := serverRaw.(map[string]interface{})
			server := vpcv1.DnsServerPrototype{
				Address: core.StringPtr(serverMap["address"].(string)),
			}
			if zoneAffinity, ok := serverMap["zone_affinity"].(string); ok && zoneAffinity != "" {
				server.ZoneAffinity = &vpcv1.ZoneIdentityByName{
					Name: core.StringPtr(zoneAffinity),
				}
			}
			manualServers = append(manualServers, server)
		}
		resolverPatch.ManualServers = manualServers

	case "delegated":
		// For delegated, the VPC is referenced through the binding
		// We need to use the VPC from the binding
		if resolverVpcID := d.Get("resolver_vpc_id").(string); resolverVpcID != "" {
			resolverPatch.VPC = &vpcv1.VpcdnsResolverVPCPatchVPCIdentityByID{
				ID: core.StringPtr(resolverVpcID),
			}
		} else if resolverVpcCRN := d.Get("resolver_vpc_crn").(string); resolverVpcCRN != "" {
			resolverPatch.VPC = &vpcv1.VpcdnsResolverVPCPatchVPCIdentityByCRN{
				CRN: core.StringPtr(resolverVpcCRN),
			}
		}

	case "system":
		// System resolver requires no additional configuration
		// Set VPC to null if transitioning from delegated
		resolverPatch.VPC = nil
	}

	// Build VPC DNS patch
	vpcDnsPatch := &vpcv1.VpcdnsPatch{
		EnableHub: core.BoolPtr(d.Get("enable_hub").(bool)),
		Resolver:  resolverPatch,
	}

	// Build VPC patch
	vpcPatch := &vpcv1.VPCPatch{
		Dns: vpcDnsPatch,
	}

	vpcPatchAsPatch, err := vpcPatch.AsPatch()
	if err != nil {
		return err
	}

	updateVpcOptions := &vpcv1.UpdateVPCOptions{
		ID:       &vpcID,
		VPCPatch: vpcPatchAsPatch,
	}

	_, _, err = sess.UpdateVPC(updateVpcOptions)
	if err != nil {
		return fmt.Errorf("Error updating VPC DNS configuration: %s", err)
	}

	return nil
}
