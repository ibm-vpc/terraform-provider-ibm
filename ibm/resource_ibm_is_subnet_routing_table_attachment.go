// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	rtRoutes = "routes"
)

func resourceIBMISSubnetRoutingTableAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISSubnetRoutingTableAttachmentCreate,
		Read:     resourceIBMISSubnetRoutingTableAttachmentRead,
		Update:   resourceIBMISSubnetRoutingTableAttachmentUpdate,
		Delete:   resourceIBMISSubnetRoutingTableAttachmentDelete,
		Exists:   resourceIBMISSubnetRoutingTableAttachmentExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isSubnetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier",
			},

			isRoutingTableID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of routing table",
			},

			rtRouteDirectLinkIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
			},

			rtIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC",
			},
			rtLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "he lifecycle state of the routing table [ deleting, failed, pending, stable, suspended, updating, waiting ]",
			},

			isRoutingTableName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the routing table",
			},
			isRoutingTableResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type",
			},

			rtRouteTransitGatewayIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
			},

			rtRouteVPCZoneIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
			},

			rtSubnets: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID",
						},
					},
				},
			},

			rtRoutes: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "route name",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "route ID",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISSubnetRoutingTableAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnet := d.Get(isSubnetID).(string)
	routingTable := d.Get(isRoutingTableID).(string)

	// Construct an instance of the RoutingTableIdentityByID model
	routingTableIdentityModel := new(vpcv1.RoutingTableIdentityByID)
	routingTableIdentityModel.ID = &routingTable

	// Construct an instance of the ReplaceSubnetRoutingTableOptions model
	replaceSubnetRoutingTableOptionsModel := new(vpcv1.ReplaceSubnetRoutingTableOptions)
	replaceSubnetRoutingTableOptionsModel.ID = &subnet
	replaceSubnetRoutingTableOptionsModel.RoutingTableIdentity = routingTableIdentityModel
	resultACL, response, err := sess.ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptionsModel)

	if err != nil {
		log.Printf("[DEBUG] Error while attaching a routing table to a subnet %s\n%s", err, response)
		return fmt.Errorf("Error while attaching a routing table to a subnet %s\n%s", err, response)
	}
	d.SetId(subnet)
	log.Printf("[INFO] Routing Table : %s", *resultACL.ID)
	log.Printf("[INFO] Subnet ID : %s", subnet)

	return resourceIBMISSubnetRoutingTableAttachmentRead(d, meta)
}

func resourceIBMISSubnetRoutingTableAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSubnetRoutingTableOptionsModel := &vpcv1.GetSubnetRoutingTableOptions{
		ID: &id,
	}
	nwacl, response, err := sess.GetSubnetRoutingTable(getSubnetRoutingTableOptionsModel)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting subnet's (%s) attached routing table: %s\n%s", id, err, response)
	}
	d.Set(isRoutingTableName, *nwacl.Name)
	d.Set(isRoutingTableResourceType, *nwacl.ResourceType)
	d.Set(rtRouteDirectLinkIngress, *nwacl.RouteDirectLinkIngress)
	d.Set(rtIsDefault, *nwacl.IsDefault)
	d.Set(rtLifecycleState, *nwacl.LifecycleState)
	d.Set(isRoutingTableResourceType, *nwacl.ResourceType)
	d.Set(rtRouteTransitGatewayIngress, *nwacl.RouteTransitGatewayIngress)
	d.Set(rtRouteVPCZoneIngress, *nwacl.RouteVPCZoneIngress)
	subnets := make([]map[string]interface{}, 0)

	for _, s := range nwacl.Subnets {
		subnet := make(map[string]interface{})
		subnet[ID] = *s.ID
		subnet["name"] = *s.Name
		subnets = append(subnets, subnet)
	}
	d.Set(rtSubnets, subnets)

	routes := make([]map[string]interface{}, 0)
	for _, s := range nwacl.Routes {
		route := make(map[string]interface{})
		route[ID] = *s.ID
		route["name"] = *s.Name
		routes = append(routes, route)
	}
	d.Set(rtRoutes, routes)
	return nil
}

func resourceIBMISSubnetRoutingTableAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isRoutingTableID) {
		subnet := d.Get(isSubnetID).(string)
		routingTable := d.Get(isRoutingTableID).(string)

		// Construct an instance of the RoutingTableIdentityByID model
		routingTableIdentityModel := new(vpcv1.RoutingTableIdentityByID)
		routingTableIdentityModel.ID = &routingTable

		// Construct an instance of the ReplaceSubnetRoutingTableOptions model
		replaceSubnetRoutingTableOptionsModel := new(vpcv1.ReplaceSubnetRoutingTableOptions)
		replaceSubnetRoutingTableOptionsModel.ID = &subnet
		replaceSubnetRoutingTableOptionsModel.RoutingTableIdentity = routingTableIdentityModel
		resultACL, response, err := sess.ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a routing table to a subnet %s\n%s", err, response)
			return fmt.Errorf("Error while attaching a routing table to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with Routing Table : %s", subnet, *resultACL.ID)

		d.SetId(subnet)
		return resourceIBMISSubnetRoutingTableAttachmentRead(d, meta)
	}

	return resourceIBMISSubnetRoutingTableAttachmentRead(d, meta)
}

func resourceIBMISSubnetRoutingTableAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	// Set the subnet with VPC default routing table
	getSubnetOptions := &vpcv1.GetSubnetOptions{
		ID: &id,
	}
	subnet, response, err := sess.GetSubnet(getSubnetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Subnet (%s): %s\n%s", id, err, response)
	}
	// Fetch VPC
	vpcID := *subnet.VPC.ID

	getvpcOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getvpcOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting VPC : %s\n%s", err, response)
	}

	// Fetch default routing table
	if vpc.DefaultRoutingTable != nil {
		log.Printf("[DEBUG] vpc default routing table is not null :%s", *vpc.DefaultRoutingTable.ID)
		// Construct an instance of the RoutingTableIdentityByID model
		routingTableIdentityModel := new(vpcv1.RoutingTableIdentityByID)
		routingTableIdentityModel.ID = vpc.DefaultRoutingTable.ID

		// Construct an instance of the ReplaceSubnetRoutingTableOptions model
		replaceSubnetRoutingTableOptionsModel := new(vpcv1.ReplaceSubnetRoutingTableOptions)
		replaceSubnetRoutingTableOptionsModel.ID = &id
		replaceSubnetRoutingTableOptionsModel.RoutingTableIdentity = routingTableIdentityModel
		resultACL, response, err := sess.ReplaceSubnetRoutingTable(replaceSubnetRoutingTableOptionsModel)

		if err != nil {
			log.Printf("[DEBUG] Error while attaching a routing table to a subnet %s\n%s", err, response)
			return fmt.Errorf("Error while attaching a routing table to a subnet %s\n%s", err, response)
		}
		log.Printf("[INFO] Updated subnet %s with VPC default Routing Table : %s", id, *resultACL.ID)
	} else {
		log.Printf("[DEBUG] vpc default routing table is  null")
	}

	d.SetId("")
	return nil
}

func resourceIBMISSubnetRoutingTableAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSubnetRoutingTableOptionsModel := &vpcv1.GetSubnetRoutingTableOptions{
		ID: &id,
	}
	_, response, err := sess.GetSubnetRoutingTable(getSubnetRoutingTableOptionsModel)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting subnet's attached routing table: %s\n%s", err, response)
	}
	return true, nil
}
