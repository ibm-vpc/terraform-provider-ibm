// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBPoolMembersPort               = "port"
	isLBPoolMembersTargetAddress      = "target_address"
	isLBPoolMembersTargetID           = "target_id"
	isLBPoolMembersWeight             = "weight"
	isLBPoolMembersProvisioningStatus = "provisioning_status"
	isLBPoolMembersHealth             = "health"
	isLBPoolMembersHref               = "href"
	isLBPoolMembersDeletePending      = "delete_pending"
	isLBPoolMembersDeleted            = "done"
	isLBPoolMembersActive             = "active"
	isLBPoolMembers                   = "members"
	isLBPoolMembersTarget             = "target"
)

func ResourceIBMISLBPoolMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMISLBPoolMembersCreate,
		Read:   resourceIBMISLBPoolMembersRead,
		Update: resourceIBMISLBPoolMembersUpdate,
		Delete: resourceIBMISLBPoolMembersDelete,
		//Exists:   resourceIBMISLBPoolMembersExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isLBPoolID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					// if state file entry and tf file entry matches
					if strings.Compare(n, o) == 0 {
						return true
					}

					if strings.Contains(n, "/") {
						new := strings.Split(n, "/")
						if strings.Compare(new[1], o) == 0 {
							return true
						}
					}

					return false
				},
				Description: "Loadblancer Poold ID",
			},

			isLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Load balancer ID",
			},
			isLBPoolMembers: {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceIBMIsLBPoolMembersHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBPoolMembersPort: {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool_members", isLBPoolMembersPort),
							Description:  "Load Balancer Pool port",
						},
						isLBPoolMembersTarget: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Load balancer pool member target, Either target id or ip address",
						},
						isLBPoolMembersWeight: {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      50,
							ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool_members", isLBPoolMembersWeight),
							Description:  "Load balcner pool member weight",
						},
						isLBPoolMembersProvisioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancer Pool member provisioning status",
						},

						isLBPoolMembersHealth: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "LB Pool member health",
						},

						isLBPoolMembersHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "LB pool member Href value",
						},
						"member": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "LB pool member ID",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this member was created",
						},
					},
				},
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func ResourceIBMISLBPoolMembersValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolMembersWeight,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "100"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolMembersPort,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "65535"})

	ibmISLBResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_pool_members", Schema: validateSchema}
	return &ibmISLBResourceValidator
}

func resourceIBMISLBPoolMembersCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] LB Pool create")
	lbPoolID, err := getPoolId(d.Get(isLBPoolID).(string))
	if err != nil {
		return err
	}

	lbID := d.Get(isLBID).(string)
	replaceLBPoolMembersOptions := &vpcv1.ReplaceLoadBalancerPoolMembersOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
	}

	membersSet := d.Get(isLBPoolMembers)
	membersList := membersSet.(*schema.Set).List()
	members := []vpcv1.LoadBalancerPoolMemberPrototype{}
	for _, a := range membersList {
		memberItem := a.(map[string]interface{})
		var member = vpcv1.LoadBalancerPoolMemberPrototype{}
		port := int64(memberItem["port"].(int))
		target := memberItem["target"].(string)
		weight := int64(memberItem["weight"].(int))
		member.Port = &port
		member.Weight = &weight
		if net.ParseIP(target) == nil {
			member.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeInstanceIdentity{
				ID: &target,
			}
		} else {
			member.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeIP{
				Address: &target,
			}
		}
		members = append(members, member)
	}
	replaceLBPoolMembersOptions.Members = members

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	lbPoolMembers, response, err := sess.ReplaceLoadBalancerPoolMembers(replaceLBPoolMembersOptions)
	if err != nil {
		return fmt.Errorf("[DEBUG] lbpool members replace err: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lbID, lbPoolID, *lbPoolMembers.Members[0].ID))
	//log.Printf("[INFO] lbpool member : %s", *lbPoolMembers.ID)

	membersIntfList := make([]interface{}, 0, len(lbPoolMembers.Members))
	for _, member := range lbPoolMembers.Members {
		log.Println("calling refresh for pool member id ", *member.ID)
		_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, *member.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		log.Println("calling get for pool member id ", *member.ID)
		memberMap, err := lbpmembersGet(d, meta, lbID, lbPoolID, *member.ID)
		if err != nil {
			return err
		}
		log.Println("after calling get for pool member id ", memberMap["member"].(string))
		membersIntfList = append(membersIntfList, memberMap)
	}
	memberslistSet := schema.NewSet(resourceIBMIsLBPoolMembersHash,
		membersIntfList).List()
	for _, a := range memberslistSet {
		memberItem := a.(map[string]interface{})

		log.Println("set data member: ", memberItem["member"].(string))
		log.Println("set data target: ", memberItem["target"].(string))
	}
	if err := d.Set(isLBPoolMembers, schema.NewSet(resourceIBMIsLBPoolMembersHash,
		membersIntfList)); err != nil {
		return fmt.Errorf("[ERROR] Error setting members: %s", err)
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	return resourceIBMISLBPoolMembersRead(d, meta)
}

func isWaitForLBPoolMembersAvailable(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer pool member(%s) to be available.", lbPoolMemID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBPoolMembersActive, ""},
		Refresh:    isLBPoolMembersRefreshFunc(lbc, lbID, lbPoolID, lbPoolMemID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBPoolMembersRefreshFunc(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Println("pool member id refresh", lbPoolMemID)
		getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &lbPoolID,
			ID:             &lbPoolMemID,
		}
		lbPoolMem, response, err := lbc.GetLoadBalancerPoolMember(getlbpmoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Members: %s\n%s", err, response)
		}

		if *lbPoolMem.ProvisioningStatus == isLBPoolMembersActive {
			return lbPoolMem, *lbPoolMem.ProvisioningStatus, nil
		}

		return lbPoolMem, *lbPoolMem.ProvisioningStatus, nil
	}
}

func resourceIBMISLBPoolMembersRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return nil
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	if len(parts) < 3 {
		return fmt.Errorf(
			"The id should contain loadbalancer Id, loadbalancer pool Id and loadbalancer poolmemebr Id")
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	_, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Load Balancer : %s\n%s", err, response)
	}
	getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}

	_, response, err = sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Pool : %s\n%s", err, response)
	}
	//lbPoolMemID := parts[2]
	membersSet := d.Get(isLBPoolMembers)
	membersList := membersSet.(*schema.Set).List()
	//members := []vpcv1.LoadBalancerPoolMemberPrototype{}
	membersIntfList := make([]interface{}, 0, len(membersList))
	for _, a := range membersList {
		memberItem := a.(map[string]interface{})
		//var member = &vpcv1.LoadBalancerPoolMemberPrototype{}
		memberID := memberItem["member"].(string)
		member, err := lbpmembersGet(d, meta, lbID, lbPoolID, memberID)
		if err != nil {
			return err
		}
		membersIntfList = append(membersIntfList, member)
	}
	d.Set(isLBPoolMembers, schema.NewSet(resourceIBMIsLBPoolMembersHash,
		membersIntfList))
	return nil
}

func lbpmembersGet(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) (map[string]interface{}, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return nil, err
	}
	log.Println("pool member id", lbPoolMemID)
	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	lbPoolMem, response, err := sess.GetLoadBalancerPoolMember(getlbpmoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil, err
		}
		return nil, fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Members: %s\n%s", err, response)
	}

	// getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
	// 	ID: &lbID,
	// }
	// lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	// if err != nil {
	// 	return nil, fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response)
	// }
	// d.Set(flex.RelatedCRN, *lb.CRN)
	return memberToMap(lbPoolMem), nil
}

func memberToMap(lbPoolMem *vpcv1.LoadBalancerPoolMember) map[string]interface{} {

	mem := make(map[string]interface{})
	mem["member"] = *lbPoolMem.ID
	mem["port"] = *lbPoolMem.Port
	mem["weight"] = *lbPoolMem.Weight
	target := lbPoolMem.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
	if target.Address != nil {
		mem["target"] = *target.Address
	}
	if target.ID != nil {
		mem["target"] = *target.ID
	}

	mem["created_at"] = flex.DateTimeToString(lbPoolMem.CreatedAt)
	mem["href"] = *lbPoolMem.Href
	mem["health"] = *lbPoolMem.Health
	mem["provisioning_status"] = *lbPoolMem.ProvisioningStatus
	return mem
}

func resourceIBMISLBPoolMembersUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	err = lbpmembersUpdate(d, meta, lbID, lbPoolID)
	if err != nil {
		return err
	}

	return resourceIBMISLBPoolMembersRead(d, meta)
}

func lbpmembersUpdate(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	membersSet := d.Get(isLBPoolMembers)
	membersList := membersSet.(*schema.Set).List()
	//members := []vpcv1.LoadBalancerPoolMemberPrototype{}
	membersIntfList := make([]interface{}, 0, len(membersList))

	if d.HasChange(isLBPoolMembers) {
		old, new := d.GetChange(isLBPoolMembers)
		oldset := old.(*schema.Set)
		newset := new.(*schema.Set)
		if oldset.Difference(newset).Len() > 0 {
			log.Println("Update: members len oldset differend ", oldset.Difference(newset).Len())
		} else {
			log.Println("Update: members len oldset differend false")
		}
		diffset := oldset.Difference(newset)

		for _, difflistitem := range diffset.List() {
			difflistmap := difflistitem.(map[string]interface{})
			log.Println("Update:difflistmap members port: ", difflistmap["port"].(int))
			log.Println("Update:difflistmap members ID: ", difflistmap["member"].(string))
			log.Println("Update:difflistmap members weight: ", difflistmap["weight"].(int))
			log.Println("Update:difflistmap members target: ", difflistmap["target"].(string))
		}

		diffset1 := newset.Difference(oldset)
		for _, difflistitem := range diffset1.List() {
			difflistmap1 := difflistitem.(map[string]interface{})
			log.Println("Update:difflistmap1 members port: ", difflistmap1["port"].(int))
			log.Println("Update:difflistmap1 members ID: ", difflistmap1["member"].(string))
			log.Println("Update:difflistmap1 members weight: ", difflistmap1["weight"].(int))
			log.Println("Update:difflistmap1 members target: ", difflistmap1["target"].(string))
		}
		if diffset.Len() > diffset1.Len() {
			for i := 0; i < diffset.Len(); i++ {
				member := diffset.List()[i].(map[string]interface{})
				memberId := member["member"].(string)
				err = lbpmemberDelete(d, meta, lbID, lbPoolID, memberId)
				if err != nil {
					return err
				}
			}
		} else if diffset1.Len() == diffset.Len() {
			for i := 0; i < diffset.Len(); i++ {

				oldMember := diffset.List()[i].(map[string]interface{})
				//newMember := diffset1.List()[i].(map[string]interface{})

				statePort := int64(oldMember["port"].(int))
				stateWeight := int64(oldMember["weight"].(int))
				stateTarget := oldMember["target"].(string)
				memberId := oldMember["member"].(string)
				for i := 0; i < diffset1.Len(); i++ {
					newMember := diffset1.List()[i].(map[string]interface{})
					if statePort == int64(newMember["port"].(int)) || stateWeight == int64(newMember["weight"].(int)) || stateTarget == newMember["target"].(string) {

						_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf(
								"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
						}

						_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, memberId, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}

						_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf(
								"Error checking for load balancer (%s) is active: %s", lbID, err)
						}
						updatelbpmoptions := &vpcv1.UpdateLoadBalancerPoolMemberOptions{
							LoadBalancerID: &lbID,
							PoolID:         &lbPoolID,
							ID:             &memberId,
						}
						loadBalancerPoolMembersPatchModel := &vpcv1.LoadBalancerPoolMemberPatch{}
						if newMember["port"].(int) != oldMember["port"].(int) {
							port := int64(newMember["port"].(int))
							loadBalancerPoolMembersPatchModel.Port = &port
						}
						if newMember["weight"].(int) != oldMember["weight"].(int) {
							weight := int64(newMember["weight"].(int))
							loadBalancerPoolMembersPatchModel.Weight = &weight
						}

						if newMember["target"].(string) != oldMember["target"].(string) {

							target := newMember["target"].(string)

							if net.ParseIP(target) == nil {
								loadBalancerPoolMembersPatchModel.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeInstanceIdentity{
									ID: &target,
								}
							} else {
								loadBalancerPoolMembersPatchModel.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeIP{
									Address: &target,
								}
							}
						}
						loadBalancerPoolMembersPatch, err := loadBalancerPoolMembersPatchModel.AsPatch()
						if err != nil {
							return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPoolMembersPatch: %s", err)
						}
						updatelbpmoptions.LoadBalancerPoolMemberPatch = loadBalancerPoolMembersPatch

						member, response, err := sess.UpdateLoadBalancerPoolMember(updatelbpmoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error Updating Load Balancer Pool Members: %s\n%s", err, response)
						}
						_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, memberId, d.Timeout(schema.TimeoutCreate))
						if err != nil {
							return err
						}

						_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf(
								"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
						}

						_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf(
								"Error checking for load balancer (%s) is active: %s", lbID, err)
						}
						membersIntfList = append(membersIntfList, memberToMap(member))

						diffset1.Remove(newMember)
						break
					}
				}
			}
		}
		for _, a := range membersList {
			memberItem := a.(map[string]interface{})
			//var member = &vpcv1.LoadBalancerPoolMemberPrototype{}
			memberID := memberItem["member"].(string)
			if memberID != "" {
				member, err := lbpmembersGet(d, meta, lbID, lbPoolID, memberID)
				if err != nil {
					return err
				}
				membersIntfList = append(membersIntfList, member)
			}
		}
		d.Set(isLBPoolMembers, schema.NewSet(resourceIBMIsLBPoolMembersHash,
			membersIntfList))
		// } else if diffset1.Len() > diffset.Len() {
		// 	for i := 0; i < diffset1.Len(); i++ {
		// 		member := diffset1.List()[i].(map[string]interface{})
		// 		memberId := member["member"].(string)
		// 		err = lbpmemberDelete(d, meta, lbID, lbPoolID, memberId)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }

		// membersSet := d.Get(isLBPoolMembers)
		// membersList := membersSet.(*schema.Set).List()
		// log.Println("Update: members len", len(membersList))
		// for _, memberIntf := range membersList {
		// 	member := memberIntf.(map[string]interface{})
		// 	log.Println("Update: members port: ", member["port"].(int))
		// 	log.Println("Update: members weight: ", member["weight"].(int))
		// 	log.Println("Update: members ID: ", member["member"].(string))
		// 	log.Println("Update: members target: ", member["target"].(string))
		// }
		// //count := len(members.([]interface{}))
		// isLBKey := "load_balancer_key_" + lbID
		// conns.IbmMutexKV.Lock(isLBKey)
		// defer conns.IbmMutexKV.Unlock(isLBKey)

		// for i := 0; i < len(membersList); i++ {
		// 	if d.HasChange("members." + strconv.Itoa(i)) {
		// 		log.Println("Update: has change true")
		// 		old, _ := d.GetChange("members." + strconv.Itoa(i))
		// 		member := old.(map[string]interface{})
		// 		id := member["member"].(string)
		// 		log.Println("Update: has change true, id: ", id)
		// 	} else {
		// 		log.Println("Update: has change false")
		// 	}
		// }
		// for i, memberIntf := range membersList {
		// 	member := memberIntf.(map[string]interface{})
		// 	hasChange := false
		// 	memberId := member["member"].(string)
		// 	port := member["port"].(int)
		// 	log.Println("Update: member id", memberId)
		// 	log.Println("Update: port", port)
		// 	if d.HasChange(isLBPoolMembers + "." + strconv.Itoa(i) + "member") {
		// 		old, new := d.GetChange(isLBPoolMembers + "." + strconv.Itoa(i) + "member")
		// 		memberId = old.(string)
		// 		newID := new.(string)
		// 		log.Println("Update: before wait old memberID", memberId)
		// 		log.Println("Update: before wait new memberID", newID)
		// 		// hasChange = true
		// 	}
		// 	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		// 	if err != nil {
		// 		return fmt.Errorf(
		// 			"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		// 	}

		// 	_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, memberId, d.Timeout(schema.TimeoutUpdate))
		// 	if err != nil {
		// 		return err
		// 	}

		// 	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		// 	if err != nil {
		// 		return fmt.Errorf(
		// 			"Error checking for load balancer (%s) is active: %s", lbID, err)
		// 	}
		// 	updatelbpmoptions := &vpcv1.UpdateLoadBalancerPoolMemberOptions{
		// 		LoadBalancerID: &lbID,
		// 		PoolID:         &lbPoolID,
		// 		ID:             &memberId,
		// 	}

		// 	loadBalancerPoolMembersPatchModel := &vpcv1.LoadBalancerPoolMemberPatch{}
		// 	if d.HasChange(isLBPoolMembers + "." + strconv.Itoa(i) + "target") {
		// 		target := member["target"].(string)
		// 		if net.ParseIP(target) == nil {
		// 			loadBalancerPoolMembersPatchModel.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeInstanceIdentity{
		// 				ID: &target,
		// 			}
		// 		} else {
		// 			loadBalancerPoolMembersPatchModel.Target = &vpcv1.LoadBalancerPoolMemberTargetPrototypeIP{
		// 				Address: &target,
		// 			}
		// 		}
		// 		hasChange = true
		// 	}
		// 	if d.HasChange(isLBPoolMembers + "." + strconv.Itoa(i) + "weight") {
		// 		weight := int64(member["weight"].(int))
		// 		loadBalancerPoolMembersPatchModel.Weight = &weight
		// 		hasChange = true
		// 	}
		// 	if d.HasChange(isLBPoolMembers + "." + strconv.Itoa(i) + "port") {
		// 		port := int64(member["port"].(int))
		// 		loadBalancerPoolMembersPatchModel.Port = &port
		// 		hasChange = true
		// 	}
		// 	if d.HasChange(isLBPoolMembers + "." + strconv.Itoa(i) + "member") {
		// 		old, new := d.GetChange(isLBPoolMembers + "." + strconv.Itoa(i) + "member")
		// 		memberID := old.(string)
		// 		newID := new.(string)
		// 		log.Println("Update: old memberID", memberID)
		// 		log.Println("Update: new memberID", newID)
		// 		updatelbpmoptions.ID = &memberID

		// 		hasChange = true
		// 	}
		// 	if hasChange {
		// 		loadBalancerPoolMembersPatch, err := loadBalancerPoolMembersPatchModel.AsPatch()
		// 		if err != nil {
		// 			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPoolMembersPatch: %s", err)
		// 		}
		// 		updatelbpmoptions.LoadBalancerPoolMemberPatch = loadBalancerPoolMembersPatch

		// 		_, response, err := sess.UpdateLoadBalancerPoolMember(updatelbpmoptions)
		// 		if err != nil {
		// 			return fmt.Errorf("[ERROR] Error Updating Load Balancer Pool Members: %s\n%s", err, response)
		// 		}
		// 		_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, memberId, d.Timeout(schema.TimeoutCreate))
		// 		if err != nil {
		// 			return err
		// 		}

		// 		_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		// 		if err != nil {
		// 			return fmt.Errorf(
		// 				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		// 		}

		// 		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		// 		if err != nil {
		// 			return fmt.Errorf(
		// 				"Error checking for load balancer (%s) is active: %s", lbID, err)
		// 		}
		// 	}
		// }

	}
	// 	port := int64(d.Get(isLBPoolMembersPort).(int))
	// 	weight := int64(d.Get(isLBPoolMembersWeight).(int))

	// 	isLBKey := "load_balancer_key_" + lbID
	// 	conns.IbmMutexKV.Lock(isLBKey)
	// 	defer conns.IbmMutexKV.Unlock(isLBKey)

	// 	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf(
	// 			"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	// 	}

	// 	_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutCreate))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf(
	// 			"Error checking for load balancer (%s) is active: %s", lbID, err)
	// 	}

	// 	updatelbpmoptions := &vpcv1.UpdateLoadBalancerPoolMembersOptions{
	// 		LoadBalancerID: &lbID,
	// 		PoolID:         &lbPoolID,
	// 		ID:             &lbPoolMemID,
	// 	}

	// 	loadBalancerPoolMembersPatchModel := &vpcv1.LoadBalancerPoolMembersPatch{
	// 		Port:   &port,
	// 		Weight: &weight,
	// 	}

	// 	if _, ok := d.GetOk(isLBPoolMembersTargetAddress); ok {
	// 		targetAddress := d.Get(isLBPoolMembersTargetAddress).(string)
	// 		target := &vpcv1.LoadBalancerPoolMembersTargetPrototype{
	// 			Address: &targetAddress,
	// 		}
	// 		loadBalancerPoolMembersPatchModel.Target = target
	// 	} else {
	// 		targetID := d.Get(isLBPoolMembersTargetID).(string)
	// 		target := &vpcv1.LoadBalancerPoolMembersTargetPrototype{
	// 			ID: &targetID,
	// 		}
	// 		loadBalancerPoolMembersPatchModel.Target = target
	// 	}

	// 	loadBalancerPoolMembersPatch, err := loadBalancerPoolMembersPatchModel.AsPatch()
	// 	if err != nil {
	// 		return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPoolMembersPatch: %s", err)
	// 	}
	// 	updatelbpmoptions.LoadBalancerPoolMembersPatch = loadBalancerPoolMembersPatch

	// 	_, response, err := sess.UpdateLoadBalancerPoolMembers(updatelbpmoptions)
	// 	if err != nil {
	// 		return fmt.Errorf("[ERROR] Error Updating Load Balancer Pool Members: %s\n%s", err, response)
	// 	}
	// 	_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutCreate))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf(
	// 			"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	// 	}

	// 	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf(
	// 			"Error checking for load balancer (%s) is active: %s", lbID, err)
	// 	}
	// }
	return nil
}

func resourceIBMISLBPoolMembersDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]
	//lbPoolMemID := parts[2]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)
	members := d.Get(isLBPoolMembers)

	for _, memberIntf := range members.(*schema.Set).List() {
		member := memberIntf.(map[string]interface{})
		memberId := member["member"].(string)
		err = lbpmemberDelete(d, meta, lbID, lbPoolID, memberId)
		if err != nil {
			return err
		}
	}
	return nil
}

func lbpmembersDelete(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	_, response, err := sess.GetLoadBalancerPoolMember(getlbpmoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Members: %s\n%s", err, response)
	}
	_, err = isWaitForLBPoolMembersAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	dellbpmoptions := &vpcv1.DeleteLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	response, err = sess.DeleteLoadBalancerPoolMember(dellbpmoptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Load Balancer Pool Members: %s\n%s", err, response)
	}

	_, err = isWaitForLBPoolMembersDeleted(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	d.SetId("")
	return nil
}

func isWaitForLBPoolMembersDeleted(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbPoolMemID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolMembersDeletePending},
		Target:     []string{isLBPoolMembersDeleted, ""},
		Refresh:    isDeleteLBPoolMembersRefreshFunc(lbc, lbID, lbPoolID, lbPoolMemID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isDeleteLBPoolMembersRefreshFunc(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &lbPoolID,
			ID:             &lbPoolMemID,
		}
		lbPoolMem, response, err := lbc.GetLoadBalancerPoolMember(getlbpmoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbPoolMem, isLBPoolMembersDeleted, nil
			}
			return nil, "", fmt.Errorf("[ERROR] Error Deleting Load balancer pool member: %s\n%s", err, response)
		}
		return lbPoolMem, isLBPoolMembersDeletePending, nil
	}
}

// func resourceIBMISLBPoolMembersExists(d *schema.ResourceData, meta interface{}) (bool, error) {
// 	parts, err := flex.IdParts(d.Id())
// 	if err != nil {
// 		return false, err
// 	}
// 	if len(parts) != 3 {
// 		return false, fmt.Errorf(
// 			"The id should contain loadbalancer Id, loadbalancer pool Id and loadbalancer poolmemebr Id")
// 	}

// 	lbID := parts[0]
// 	lbPoolID := parts[1]
// 	lbPoolMemID := parts[2]

// 	exists, err := lbpmemberExists(d, meta, lbID, lbPoolID, lbPoolMemID)
// 	return exists, err

// }

// func lbpmemberExists(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) (bool, error) {
// 	sess, err := vpcClient(meta)
// 	if err != nil {
// 		return false, err
// 	}

// 	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMembersOptions{
// 		LoadBalancerID: &lbID,
// 		PoolID:         &lbPoolID,
// 		ID:             &lbPoolMemID,
// 	}
// 	_, response, err := sess.GetLoadBalancerPoolMembers(getlbpmoptions)
// 	if err != nil {
// 		if response != nil && response.StatusCode == 404 {
// 			return false, nil
// 		}
// 		return false, fmt.Errorf("[ERROR] Error getting Load balancer pool member: %s\n%s", err, response)
// 	}
// 	return true, nil
// }

// func getPoolId(id string) (string, error) {
// 	if strings.Contains(id, "/") {
// 		parts, err := flex.IdParts(id)
// 		if err != nil {
// 			return "", err
// 		}

// 		return parts[1], nil
// 	} else {
// 		return id, nil
// 	}
// }

func resourceIBMIsLBPoolMembersHash(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	switch v := a["port"].(type) {
	case int:
		buf.WriteString(fmt.Sprintf("%d-", v))
	case int64:
		buf.WriteString(fmt.Sprintf("%d-", v))
	default:
		buf.WriteString(fmt.Sprintf("%d-", 8888))
	}
	switch v := a["weight"].(type) {
	case int:
		buf.WriteString(fmt.Sprintf("%d-", v))
	case int64:
		buf.WriteString(fmt.Sprintf("%d-", v))
	default:
		buf.WriteString(fmt.Sprintf("%d-", 8888))
	}
	buf.WriteString(fmt.Sprintf("%s-", a["target"].(string)))
	log.Println("hash value string:", buf.String())
	return conns.String(buf.String())
}
