#! /bin/bash

if [ $# -ne 1 ]
then
    echo "Please Provide Resource Name"
    echo " all, vpc_all, subnet_all, vpn_server_all, vpn_gateway_all, volume_all, virtual_endpoint_gateway_all, ssh_key_all, snapshot_all, security_group_all,public_gateway_all supported"
    echo "load_balancer_all, ip_sec_policy_all, instance_all, image_all supportd"
    echo "ike_policy_all, flow_log_all, floating_ip_all, dedicated_host_all, baremetal_server_all, vpc_routing_table, vpc_routing_table_route"
    exit 1
fi

TEST=$1
if [ $TEST == "all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMIS' -timeout 700m
fi


if [ $TEST == "vpc_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVPC' -timeout 700m
fi

if [ $TEST == "subnet_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISSubnet' -timeout 700m
fi

if [ $TEST == "vpn_server_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMIsVPNServer' -timeout 700m
fi

if [ $TEST == "vpn_gateway_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVPNGateway' -timeout 700m
fi

if [ $TEST == "volume_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVolume' -timeout 700m
fi

if [ $TEST == "virtual_endpoint_gateway_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVirtualEndpointGateway' -timeout 700m
fi

if [ $TEST == "ssh_key_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISSSHKey' -timeout 700m
fi

if [ $TEST == "snapshot_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISSnapshot' -timeout 700m
fi

if [ $TEST == "security_group_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISSecurityGroup' -timeout 700m
fi

if [ $TEST == "public_gateway_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISPublicGateway' -timeout 700m
fi

if [ $TEST == "load_balancer_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISLB' -timeout 700m
fi

if [ $TEST == "ip_sec_policy_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISIPSecPolicy' -timeout 700m
fi

if [ $TEST == "instance_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISInstance' -timeout 700m
fi

if [ $TEST == "image_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISImage' -timeout 700m
fi

if [ $TEST == "ike_policy_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISIKEPolicy' -timeout 700m
fi

if [ $TEST == "flow_log_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISFlowLog' -timeout 700m
fi

if [ $TEST == "floating_ip_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISFloatingIP' -timeout 700m
fi

if [ $TEST == "dedicated_host_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIbmIsDedicatedHost' -timeout 700m
fi

if [ $TEST == "" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISBareMetalServer' -timeout 700m
fi

if [ $TEST == "backup_policy_all" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMIsBackupPolicy' -timeout 700m
fi

if [ $TEST == "vpc_routing_table" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVPCRoutingTable_' -timeout 700m
fi

if [ $TEST == "vpc_routing_table_route" ]
then
    TF_ACC=1 go test ./ibm/service/vpc -v '-run=TestAccIBMISVPCRoutingTableRoute_' -timeout 700m
fi

