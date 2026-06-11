# GitHub Issue Checklist: BYOIP (Bring Your Own IP) for VPC

**SRB Number:** 5686  
**Feature:** Bring Your Own IP (IPv4) for VPC  
**Target:** Beta 2Q2026, GA 3Q2026

## Overview

This checklist tracks all implementation tasks for the BYOIP feature. Each major change area has sub-tasks for code implementation, testing, and documentation.

---

## 1. Public Address Range Authorized CIDR

### 1.1 API Implementation
- [ ] **Code**: Implement `GET /public_address_range/authorized_cidrs` endpoint
  - [ ] Add route handler
  - [ ] Implement pagination support
  - [ ] Add filtering by region
  - [ ] Implement IAM authorization check
- [ ] **Code**: Implement `GET /public_address_range/authorized_cidrs/{id}` endpoint
  - [ ] Add route handler
  - [ ] Implement resource lookup
  - [ ] Add error handling for not found
  - [ ] Implement IAM authorization check
- [ ] **Code**: Implement `GET /public_address_range/authorized_cidrs/{id}/allocations` endpoint
  - [ ] Add route handler
  - [ ] Implement allocation tracking
  - [ ] Add pagination support
  - [ ] Return both FIP and PAR allocations
- [ ] **Test**: Unit tests for authorized CIDR endpoints
  - [ ] Test list operation with pagination
  - [ ] Test get operation with valid/invalid IDs
  - [ ] Test allocations list with various scenarios
  - [ ] Test IAM authorization
- [ ] **Test**: Integration tests for authorized CIDR workflows
  - [ ] Test end-to-end authorized CIDR retrieval
  - [ ] Test allocation tracking accuracy
  - [ ] Test cross-account isolation
- [ ] **Docs**: API documentation for authorized CIDR endpoints
  - [ ] Document request/response schemas
  - [ ] Add example requests and responses
  - [ ] Document error codes and messages
  - [ ] Add usage guidelines

---

## 2. Floating IP Profiles

### 2.1 API Implementation
- [ ] **Code**: Implement `GET /floating_ip/profiles` endpoint
  - [ ] Add route handler
  - [ ] Return available profiles
  - [ ] Add pagination support
- [ ] **Code**: Implement `GET /floating_ip/profiles/{name}` endpoint
  - [ ] Add route handler
  - [ ] Implement profile lookup
  - [ ] Add error handling
- [ ] **Test**: Unit tests for profile endpoints
  - [ ] Test list profiles
  - [ ] Test get profile by name
  - [ ] Test invalid profile name
- [ ] **Test**: Integration tests for profiles
  - [ ] Test profile availability across regions
- [ ] **Docs**: Profile API documentation
  - [ ] Document profile schema
  - [ ] List available profiles
  - [ ] Add usage examples

---

## 3. Public Address Range Profiles

### 3.1 API Implementation
- [ ] **Code**: Implement `GET /public_address_range/profiles` endpoint
  - [ ] Add route handler
  - [ ] Return available profiles
  - [ ] Add pagination support
- [ ] **Code**: Implement `GET /public_address_range/profiles/{name}` endpoint
  - [ ] Add route handler
  - [ ] Implement profile lookup
  - [ ] Add error handling
- [ ] **Test**: Unit tests for profile endpoints
  - [ ] Test list profiles
  - [ ] Test get profile by name
  - [ ] Test invalid profile name
- [ ] **Test**: Integration tests for profiles
  - [ ] Test profile availability across regions
- [ ] **Docs**: Profile API documentation
  - [ ] Document profile schema
  - [ ] List available profiles
  - [ ] Add usage examples

---

## 4. Floating IP Updates

### 4.1 API Changes
- [ ] **Code**: Add `address` parameter to `POST /floating_ips`
  - [ ] Implement address validation
  - [ ] Verify address is in authorized CIDR
  - [ ] Verify address is not already allocated
  - [ ] Add conflict detection
- [ ] **Code**: Add `resource_type` property to FloatingIP model
  - [ ] Update schema definition
  - [ ] Populate in all responses
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `profile` property to FloatingIP model
  - [ ] Update schema definition
  - [ ] Populate profile information
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `authorized_cidr` property to FloatingIP model
  - [ ] Update schema definition
  - [ ] Populate when allocated from BYOIP
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `profile.name` filter to `GET /floating_ips`
  - [ ] Implement filter logic
  - [ ] Add to query parameters
  - [ ] Update pagination handling
- [ ] **Test**: Unit tests for floating IP changes
  - [ ] Test create with specific address
  - [ ] Test address validation
  - [ ] Test authorized CIDR verification
  - [ ] Test profile filtering
  - [ ] Test new properties in responses
- [ ] **Test**: Integration tests for floating IP workflows
  - [ ] Test create FIP from BYOIP pool
  - [ ] Test create FIP from IBM pool
  - [ ] Test binding to VNI
  - [ ] Test address conflict scenarios
- [ ] **Docs**: Update floating IP documentation
  - [ ] Document new `address` parameter
  - [ ] Document new properties
  - [ ] Add BYOIP examples
  - [ ] Update all affected endpoints

### 4.2 Related Endpoint Updates
- [ ] **Code**: Update bare metal server network interface endpoints (8 endpoints)
  - [ ] Add `resource_type` to responses
  - [ ] Add `profile` to responses
  - [ ] Add `authorized_cidr` to responses
- [ ] **Code**: Update instance network interface endpoints (8 endpoints)
  - [ ] Add `resource_type` to responses
  - [ ] Add `profile` to responses
  - [ ] Add `authorized_cidr` to responses
- [ ] **Code**: Update public gateway endpoints (4 endpoints)
  - [ ] Add `resource_type` to responses
- [ ] **Code**: Update virtual network interface endpoints (4 endpoints)
  - [ ] Add `resource_type` to responses
  - [ ] Add `profile` to responses
  - [ ] Add `authorized_cidr` to responses
- [ ] **Test**: Unit tests for all updated endpoints
  - [ ] Test each endpoint individually
  - [ ] Verify new properties present
- [ ] **Test**: Integration tests for cross-resource workflows
  - [ ] Test FIP with bare metal
  - [ ] Test FIP with instance
  - [ ] Test FIP with VNI
- [ ] **Docs**: Update documentation for all affected endpoints
  - [ ] Update request/response examples
  - [ ] Add notes about new properties

---

## 5. Public Address Range Updates

### 5.1 API Changes
- [ ] **Code**: Add `cidr` parameter to `POST /public_address_ranges`
  - [ ] Implement CIDR validation
  - [ ] Verify CIDR is in authorized CIDR
  - [ ] Verify CIDR is not already allocated
  - [ ] Add conflict detection
  - [ ] Make `count` and `cidr` mutually exclusive
- [ ] **Code**: Add `ip_version` property to PublicAddressRange model
  - [ ] Update schema definition
  - [ ] Populate in all responses
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `network_prefix_length` property to PublicAddressRange model
  - [ ] Update schema definition
  - [ ] Calculate from CIDR
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `profile` property to PublicAddressRange model
  - [ ] Update schema definition
  - [ ] Populate profile information
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `authorized_cidr` property to PublicAddressRange model
  - [ ] Update schema definition
  - [ ] Populate when allocated from BYOIP
  - [ ] Add to all related endpoints
- [ ] **Code**: Add `profile.name` filter to `GET /public_address_ranges`
  - [ ] Implement filter logic
  - [ ] Add to query parameters
  - [ ] Update pagination handling
- [ ] **Code**: Remove `count` from required fields in `POST /public_address_ranges`
  - [ ] Update schema validation
  - [ ] Ensure either `count` or `cidr` is provided
- [ ] **Test**: Unit tests for public address range changes
  - [ ] Test create with specific CIDR
  - [ ] Test CIDR validation
  - [ ] Test authorized CIDR verification
  - [ ] Test mutual exclusivity of count/cidr
  - [ ] Test profile filtering
  - [ ] Test new properties in responses
- [ ] **Test**: Integration tests for PAR workflows
  - [ ] Test create PAR from BYOIP pool
  - [ ] Test create PAR from IBM pool
  - [ ] Test binding to VPC
  - [ ] Test CIDR conflict scenarios
  - [ ] Test custom routing with BYOIP PAR
- [ ] **Docs**: Update public address range documentation
  - [ ] Document new `cidr` parameter
  - [ ] Document new properties
  - [ ] Add BYOIP examples
  - [ ] Update all affected endpoints
  - [ ] Document count/cidr mutual exclusivity

### 5.2 Related Endpoint Updates
- [ ] **Code**: Update VPC endpoints (4 endpoints)
  - [ ] Add `cidr` to public_address_ranges in responses
  - [ ] Update schema definitions
- [ ] **Test**: Unit tests for VPC endpoint updates
  - [ ] Test VPC list with PARs
  - [ ] Test VPC get with PARs
  - [ ] Verify CIDR property present
- [ ] **Test**: Integration tests for VPC workflows
  - [ ] Test VPC with BYOIP PARs
  - [ ] Test VPC with IBM PARs
- [ ] **Docs**: Update VPC documentation
  - [ ] Update response examples
  - [ ] Add notes about CIDR property

---

## 6. RNOS Implementation

### 6.1 Core Logic
- [ ] **Code**: Implement authorized CIDR resource management
  - [ ] Create PublicAddressRange CRD for BYOIP pools
  - [ ] Implement CRUD operations
  - [ ] Add validation logic
  - [ ] Implement allocation tracking
- [ ] **Code**: Update floating IP creation logic
  - [ ] Support address parameter
  - [ ] Verify address in authorized CIDR
  - [ ] Track allocation in authorized CIDR
  - [ ] Update profile assignment
- [ ] **Code**: Update public address range creation logic
  - [ ] Support CIDR parameter
  - [ ] Verify CIDR in authorized CIDR
  - [ ] Track allocation in authorized CIDR
  - [ ] Update profile assignment
  - [ ] Handle count/cidr mutual exclusivity
- [ ] **Code**: Implement allocation tracking
  - [ ] Track FIP allocations
  - [ ] Track PAR allocations
  - [ ] Implement conflict detection
  - [ ] Support allocation queries
- [ ] **Code**: Implement quota enforcement
  - [ ] Add BYOIP prefix quota (5 per account per region)
  - [ ] Enforce BYOIP prefix size limits
  - [ ] Maintain existing FIP/PAR quotas
- [ ] **Test**: Unit tests for RNOS logic
  - [ ] Test authorized CIDR CRUD
  - [ ] Test FIP creation from BYOIP
  - [ ] Test PAR creation from BYOIP
  - [ ] Test allocation tracking
  - [ ] Test quota enforcement
- [ ] **Test**: Integration tests for RNOS workflows
  - [ ] Test complete BYOIP provisioning
  - [ ] Test resource allocation
  - [ ] Test deprovisioning
  - [ ] Test migration workflows
- [ ] **Docs**: RNOS implementation documentation
  - [ ] Document CRD schema
  - [ ] Document allocation tracking
  - [ ] Document quota rules

### 6.2 Feature Flags
- [ ] **Code**: Implement feature flag support
  - [ ] Add `is-floating-ip-public-address-range-byoip-ipv4` flag
  - [ ] Add `is-floating-ip-public-address-range-byoip-ipv4-allowlist` flag
  - [ ] Implement allowlist checking
  - [ ] Add feature flag validation
- [ ] **Test**: Feature flag tests
  - [ ] Test with flags enabled
  - [ ] Test with flags disabled
  - [ ] Test allowlist enforcement
- [ ] **Docs**: Feature flag documentation
  - [ ] Document flag names and purposes
  - [ ] Document allowlist management

---

## 7. Genctl Implementation

### 7.1 Projection Controller
- [ ] **Code**: Implement BYOIP pool projection
  - [ ] Project authorized CIDR CRDs to zonal side
  - [ ] Handle customer-owned vs IBM-owned logic
  - [ ] Support classic_migration states
- [ ] **Code**: Update FIP projection
  - [ ] Include BYOIP flags in projection
  - [ ] Include migration state
  - [ ] Update grpc calls to SDN
- [ ] **Code**: Update PAR projection
  - [ ] Include BYOIP flags in projection
  - [ ] Include migration state
  - [ ] Update grpc calls to SDN
- [ ] **Test**: Unit tests for projection logic
  - [ ] Test BYOIP pool projection
  - [ ] Test FIP projection with BYOIP
  - [ ] Test PAR projection with BYOIP
- [ ] **Test**: Integration tests for projection
  - [ ] Test end-to-end projection flow
  - [ ] Test state synchronization
- [ ] **Docs**: Projection controller documentation
  - [ ] Document projection logic
  - [ ] Document state handling

### 7.2 Genctl-Network
- [ ] **Code**: Implement BYOIP pool processing
  - [ ] Process customer-owned BYOIP (no action)
  - [ ] Process IBM-owned BYOIP (configure on EENs)
  - [ ] Handle classic_migration state transitions
  - [ ] Make grpc calls to SDN for whole prefix advertisement
- [ ] **Code**: Update FIP processing
  - [ ] Pass BYOIP flags to SDN
  - [ ] Pass migration state to SDN
  - [ ] Update fabcon API calls
- [ ] **Code**: Update PAR processing
  - [ ] Pass BYOIP flags to SDN
  - [ ] Pass migration state to SDN
  - [ ] Update fabcon API calls
- [ ] **Code**: Implement feature flag support
  - [ ] Add `is-floating-ip-public-address-range-byoip-ipv4-genctl-network` flag
  - [ ] Check flag before processing BYOIP
- [ ] **Test**: Unit tests for genctl-network
  - [ ] Test BYOIP pool processing
  - [ ] Test FIP processing with BYOIP
  - [ ] Test PAR processing with BYOIP
  - [ ] Test feature flag handling
- [ ] **Test**: Integration tests for genctl-network
  - [ ] Test complete workflow with SDN
  - [ ] Test migration state transitions
- [ ] **Docs**: Genctl-network documentation
  - [ ] Document BYOIP processing logic
  - [ ] Document SDN integration

---

## 8. SDN Implementation

### 8.1 Core Logic
- [ ] **Code**: Update fabcon API to accept BYOIP flags
  - [ ] Add BYOIP type parameter (customer-owned/IBM-owned)
  - [ ] Add migration state parameter
  - [ ] Update API schema
- [ ] **Code**: Implement BGP community assignment
  - [ ] Add new community strings (65201:32, 65201:33, 65201:37, 65201:39)
  - [ ] Implement community selection logic based on BYOIP type and migration state
  - [ ] Support advertisement scheme table
- [ ] **Code**: Implement FIP advertisement with BYOIP
  - [ ] Advertise with correct communities
  - [ ] Handle migration state
  - [ ] Support customer-owned vs IBM-owned
- [ ] **Code**: Implement PAR advertisement with BYOIP
  - [ ] Advertise with correct communities
  - [ ] Handle migration state
  - [ ] Support customer-owned vs IBM-owned
- [ ] **Code**: Implement whole prefix advertisement
  - [ ] Accept grpc call from genctl-network
  - [ ] Advertise with null community for blackhole
  - [ ] Support migration completion
- [ ] **Test**: Unit tests for SDN logic
  - [ ] Test community assignment
  - [ ] Test FIP advertisement
  - [ ] Test PAR advertisement
  - [ ] Test whole prefix advertisement
- [ ] **Test**: Integration tests with underlay
  - [ ] Test BGP advertisement
  - [ ] Test route propagation
  - [ ] Test traffic routing
- [ ] **Docs**: SDN implementation documentation
  - [ ] Document BGP community usage
  - [ ] Document advertisement logic
  - [ ] Document integration with genctl

---

## 9. Underlay Configuration

### 9.1 BGP Community Configuration
- [ ] **Code**: Add new BGP community strings to routers
  - [ ] Configure BGP-ORIGIN-VPC-INTERNET-CUSTOMER (65201:32)
  - [ ] Configure BGP-ORIGIN-VPC-INTERNET-CUSTOMER-HOST (65201:33)
  - [ ] Configure BGP-ORIGIN-VPC-INTERNET-REMOTE-CUSTOMER (65201:37)
  - [ ] Configure BGP-ORIGIN-VPC-INTERNET-BYOIP (65201:39)
- [ ] **Code**: Update route reflector policies
  - [ ] Match on new communities
  - [ ] Advertise to GER only
- [ ] **Code**: Update GER policies
  - [ ] Match on new communities
  - [ ] Change next-hop to discard for null routes
  - [ ] Advertise northbound with appropriate export policies
- [ ] **Test**: Underlay configuration tests
  - [ ] Test BGP community propagation
  - [ ] Test route filtering
  - [ ] Test advertisement policies
- [ ] **Docs**: Underlay configuration documentation
  - [ ] Document BGP community usage
  - [ ] Document routing policies
  - [ ] Create configuration runbooks

### 9.2 Customer-Owned BYOIP Configuration
- [ ] **Code**: Create change request process for DAR configuration
  - [ ] Define change request template
  - [ ] Implement validation checks
  - [ ] Add rollback procedures
- [ ] **Test**: DAR configuration tests
  - [ ] Test prefix configuration
  - [ ] Test advertisement to BBR
  - [ ] Test internet reachability
- [ ] **Docs**: DAR configuration documentation
  - [ ] Document change request process
  - [ ] Create operator runbook
  - [ ] Add troubleshooting guide

---

## 10. UI Implementation

### 10.1 Authorized CIDR Views
- [ ] **Code**: Implement authorized CIDR list view
  - [ ] Create list component
  - [ ] Add filtering and sorting
  - [ ] Add pagination
  - [ ] Show allocation count
- [ ] **Code**: Implement authorized CIDR details view
  - [ ] Create details component
  - [ ] Show CIDR information
  - [ ] Show allocation list
  - [ ] Add navigation to allocations
- [ ] **Test**: UI tests for authorized CIDR views
  - [ ] Test list rendering
  - [ ] Test details rendering
  - [ ] Test navigation
- [ ] **Docs**: UI documentation for authorized CIDRs
  - [ ] Add user guide
  - [ ] Add screenshots

### 10.2 Floating IP Updates
- [ ] **Code**: Update floating IP creation form
  - [ ] Add address input field
  - [ ] Add authorized CIDR selector
  - [ ] Implement validation
  - [ ] Show profile information
- [ ] **Code**: Update floating IP list view
  - [ ] Show profile badge
  - [ ] Show authorized CIDR reference
  - [ ] Add profile filter
- [ ] **Code**: Update floating IP details view
  - [ ] Show profile information
  - [ ] Show authorized CIDR link
  - [ ] Show resource type
- [ ] **Test**: UI tests for floating IP updates
  - [ ] Test creation with address
  - [ ] Test creation from authorized CIDR
  - [ ] Test list filtering
  - [ ] Test details display
- [ ] **Docs**: UI documentation for floating IPs
  - [ ] Update user guide
  - [ ] Add BYOIP examples
  - [ ] Add screenshots

### 10.3 Public Address Range Updates
- [ ] **Code**: Update PAR creation form
  - [ ] Add CIDR input field
  - [ ] Add authorized CIDR selector
  - [ ] Implement validation
  - [ ] Show profile information
  - [ ] Handle count/CIDR mutual exclusivity
- [ ] **Code**: Update PAR list view
  - [ ] Show profile badge
  - [ ] Show authorized CIDR reference
  - [ ] Add profile filter
  - [ ] Show CIDR value
- [ ] **Code**: Update PAR details view
  - [ ] Show profile information
  - [ ] Show authorized CIDR link
  - [ ] Show CIDR and prefix length
- [ ] **Test**: UI tests for PAR updates
  - [ ] Test creation with CIDR
  - [ ] Test creation from authorized CIDR
  - [ ] Test list filtering
  - [ ] Test details display
- [ ] **Docs**: UI documentation for PARs
  - [ ] Update user guide
  - [ ] Add BYOIP examples
  - [ ] Add screenshots

---

## 11. CLI Implementation

### 11.1 Authorized CIDR Commands
- [ ] **Code**: Implement `ibmcloud is public-address-range-authorized-cidrs` command
  - [ ] Add list functionality
  - [ ] Add output formatting
  - [ ] Add filtering options
- [ ] **Code**: Implement `ibmcloud is public-address-range-authorized-cidr` command
  - [ ] Add get functionality
  - [ ] Add output formatting
- [ ] **Code**: Implement `ibmcloud is public-address-range-authorized-cidr-allocations` command
  - [ ] Add list functionality
  - [ ] Add output formatting
- [ ] **Test**: CLI tests for authorized CIDR commands
  - [ ] Test list command
  - [ ] Test get command
  - [ ] Test allocations command
- [ ] **Docs**: CLI documentation for authorized CIDRs
  - [ ] Add command reference
  - [ ] Add usage examples

### 11.2 Profile Commands
- [ ] **Code**: Implement `ibmcloud is floating-ip-profiles` command
  - [ ] Add list functionality
  - [ ] Add output formatting
- [ ] **Code**: Implement `ibmcloud is floating-ip-profile` command
  - [ ] Add get functionality
  - [ ] Add output formatting
- [ ] **Code**: Implement `ibmcloud is public-address-range-profiles` command
  - [ ] Add list functionality
  - [ ] Add output formatting
- [ ] **Code**: Implement `ibmcloud is public-address-range-profile` command
  - [ ] Add get functionality
  - [ ] Add output formatting
- [ ] **Test**: CLI tests for profile commands
  - [ ] Test all profile commands
- [ ] **Docs**: CLI documentation for profiles
  - [ ] Add command reference
  - [ ] Add usage examples

### 11.3 Floating IP Updates
- [ ] **Code**: Update `ibmcloud is floating-ip-create` command
  - [ ] Add `--address` flag
  - [ ] Add validation
  - [ ] Update output to show new properties
- [ ] **Code**: Update `ibmcloud is floating-ip` command
  - [ ] Show profile information
  - [ ] Show authorized CIDR reference
  - [ ] Show resource type
- [ ] **Code**: Update `ibmcloud is floating-ips` command
  - [ ] Add `--profile` filter
  - [ ] Show profile in list output
- [ ] **Test**: CLI tests for floating IP updates
  - [ ] Test create with address
  - [ ] Test get with new properties
  - [ ] Test list with profile filter
- [ ] **Docs**: CLI documentation for floating IPs
  - [ ] Update command reference
  - [ ] Add BYOIP examples

### 11.4 Public Address Range Updates
- [ ] **Code**: Update `ibmcloud is public-address-range-create` command
  - [ ] Add `--cidr` flag
  - [ ] Make `--count` and `--cidr` mutually exclusive
  - [ ] Add validation
  - [ ] Update output to show new properties
- [ ] **Code**: Update `ibmcloud is public-address-range` command
  - [ ] Show profile information
  - [ ] Show authorized CIDR reference
  - [ ] Show CIDR and prefix length
- [ ] **Code**: Update `ibmcloud is public-address-ranges` command
  - [ ] Add `--profile` filter
  - [ ] Show profile in list output
  - [ ] Show CIDR in list output
- [ ] **Test**: CLI tests for PAR updates
  - [ ] Test create with CIDR
  - [ ] Test get with new properties
  - [ ] Test list with profile filter
- [ ] **Docs**: CLI documentation for PARs
  - [ ] Update command reference
  - [ ] Add BYOIP examples

---

## 12. Terraform Implementation

### 12.1 Data Sources
- [ ] **Code**: Implement `ibm_is_floating_ip_profiles` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_floating_ip_profile` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_public_address_range_profiles` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_public_address_range_profile` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_public_address_range_authorized_cidrs` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_public_address_range_authorized_cidr` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Code**: Implement `ibm_is_public_address_range_authorized_cidr_allocations` data source
  - [ ] Add schema definition
  - [ ] Implement read logic
  - [ ] Add attributes
- [ ] **Test**: Unit tests for data sources
  - [ ] Test each data source
  - [ ] Test error handling
- [ ] **Test**: Acceptance tests for data sources
  - [ ] Test with real API
- [ ] **Docs**: Data source documentation
  - [ ] Add documentation for each data source
  - [ ] Add usage examples

### 12.2 Resource Updates
- [ ] **Code**: Update `ibm_is_floating_ip` resource
  - [ ] Add `address` argument
  - [ ] Add `profile` computed attribute
  - [ ] Add `authorized_cidr` computed attribute
  - [ ] Add `resource_type` computed attribute
  - [ ] Implement validation
- [ ] **Code**: Update `ibm_is_public_address_range` resource
  - [ ] Add `cidr` argument
  - [ ] Make `count` and `cidr` mutually exclusive
  - [ ] Add `ip_version` computed attribute
  - [ ] Add `network_prefix_length` computed attribute
  - [ ] Add `profile` computed attribute
  - [ ] Add `authorized_cidr` computed attribute
  - [ ] Implement validation
- [ ] **Test**: Unit tests for resource updates
  - [ ] Test floating IP with address
  - [ ] Test PAR with CIDR
  - [ ] Test validation rules
- [ ] **Test**: Acceptance tests for resources
  - [ ] Test create FIP from BYOIP
  - [ ] Test create PAR from BYOIP
  - [ ] Test with real API
- [ ] **Docs**: Resource documentation
  - [ ] Update floating IP documentation
  - [ ] Update PAR documentation
  - [ ] Add BYOIP examples

---

## 13. SDK Implementation

### 13.1 Go SDK
- [ ] **Code**: Update FloatingIP model
  - [ ] Add Address field
  - [ ] Add ResourceType field
  - [ ] Add Profile field
  - [ ] Add AuthorizedCIDR field
- [ ] **Code**: Update PublicAddressRange model
  - [ ] Add CIDR field
  - [ ] Add IPVersion field
  - [ ] Add NetworkPrefixLength field
  - [ ] Add Profile field
  - [ ] Add AuthorizedCIDR field
- [ ] **Code**: Add profile service methods
  - [ ] ListFloatingIPProfiles
  - [ ] GetFloatingIPProfile
  - [ ] ListPublicAddressRangeProfiles
  - [ ] GetPublicAddressRangeProfile
- [ ] **Code**: Add authorized CIDR service methods
  - [ ] ListPublicAddressRangeAuthorizedCIDRs
  - [ ] GetPublicAddressRangeAuthorizedCIDR
  - [ ] ListPublicAddressRangeAuthorizedCIDRAllocations
- [ ] **Test**: Unit tests for SDK changes
  - [ ] Test model serialization
  - [ ] Test service methods
- [ ] **Docs**: SDK documentation
  - [ ] Update API reference
  - [ ] Add code examples

### 13.2 Python SDK
- [ ] **Code**: Update FloatingIP model
- [ ] **Code**: Update PublicAddressRange model
- [ ] **Code**: Add profile service methods
- [ ] **Code**: Add authorized CIDR service methods
- [ ] **Test**: Unit tests for SDK changes
- [ ] **Docs**: SDK documentation

### 13.3 Node.js SDK
- [ ] **Code**: Update FloatingIP model
- [ ] **Code**: Update PublicAddressRange model
- [ ] **Code**: Add profile service methods
- [ ] **Code**: Add authorized CIDR service methods
- [ ] **Test**: Unit tests for SDK changes
- [ ] **Docs**: SDK documentation

### 13.4 Java SDK
- [ ] **Code**: Update FloatingIP model
- [ ] **Code**: Update PublicAddressRange model
- [ ] **Code**: Add profile service methods
- [ ] **Code**: Add authorized CIDR service methods
- [ ] **Test**: Unit tests for SDK changes
- [ ] **Docs**: SDK documentation

---

## 14. Documentation

### 14.1 Customer-Facing Documentation
- [ ] **Docs**: Create BYOIP overview guide
  - [ ] Explain BYOIP concept
  - [ ] List use cases
  - [ ] Describe provisioning process
  - [ ] Add prerequisites
- [ ] **Docs**: Create BYOIP getting started guide
  - [ ] Step-by-step provisioning
  - [ ] Creating resources from BYOIP
  - [ ] Best practices
- [ ] **Docs**: Create BYOIP API reference
  - [ ] Document all new endpoints
  - [ ] Add request/response examples
  - [ ] Document error codes
- [ ] **Docs**: Update floating IP documentation
  - [ ] Add BYOIP section
  - [ ] Update examples
  - [ ] Add troubleshooting
- [ ] **Docs**: Update public address range documentation
  - [ ] Add BYOIP section
  - [ ] Update examples
  - [ ] Add troubleshooting
- [ ] **Docs**: Create BYOIP migration guide
  - [ ] Classic to VPC migration
  - [ ] IBM pool to BYOIP migration
  - [ ] Best practices
- [ ] **Docs**: Create BYOIP troubleshooting guide
  - [ ] Common issues
  - [ ] Error messages
  - [ ] Resolution steps

---

## 15. Testing

### 15.1 Unit Tests
- [ ] **Test**: API layer unit tests (all endpoints)
- [ ] **Test**: RNOS unit tests (all logic)
- [ ] **Test**: Genctl unit tests (all logic)
- [ ] **Test**: SDN unit tests (all logic)
- [ ] **Test**: UI component unit tests
- [ ] **Test**: CLI command unit tests
- [ ] **Test**: Terraform provider unit tests
- [ ] **Test**: SDK unit tests (all languages)

### 15.2 Integration Tests
- [ ] **Test**: End-to-end BYOIP provisioning
- [ ] **Test**: End-to-end resource allocation
- [ ] **Test**: End-to-end deprovisioning
- [ ] **Test**: Migration workflows
- [ ] **Test**: Cross-component integration
- [ ] **Test**: UI workflows
- [ ] **Test**: CLI workflows
- [ ] **Test**: Terraform workflows

### 15.3 Performance Tests
- [ ] **Test**: API performance under load
- [ ] **Test**: Large-scale resource allocation
- [ ] **Test**: Concurrent operations
- [ ] **Test**: Database performance
- [ ] **Test**: Network performance

### 15.4 Security Tests
- [ ] **Test**: IAM authorization
- [ ] **Test**: Input validation
- [ ] **Test**: Injection prevention
- [ ] **Test**: Cross-account isolation
- [ ] **Test**: Operator API security

### 15.5 Acceptance Tests
- [ ] **Test**: Customer workflow validation
- [ ] **Test**: Operator workflow validation
- [ ] **Test**: Migration scenario validation
- [ ] **Test**: Error handling validation
- [ ] **Test**: Billing and metering validation

---

## 16. Operations

### 16.1 Monitoring and Alerting
- [ ] **Code**: Implement metrics collection
  - [ ] Authorized CIDR count
  - [ ] Allocation count per CIDR
  - [ ] API request metrics
  - [ ] Error rate metrics
- [ ] **Code**: Create dashboards
  - [ ] BYOIP overview dashboard
  - [ ] Resource allocation dashboard
  - [ ] Error tracking dashboard
- [ ] **Code**: Configure alerts
  - [ ] High error rate alerts
  - [ ] Quota threshold alerts
  - [ ] Provisioning failure alerts
- [ ] **Test**: Monitoring tests
  - [ ] Verify metrics collection
  - [ ] Test alert triggers
- [ ] **Docs**: Monitoring documentation
  - [ ] Document metrics
  - [ ] Document dashboards
  - [ ] Document alerts

### 16.2 Logging
- [ ] **Code**: Implement structured logging
  - [ ] API request logging
  - [ ] Operator action logging
  - [ ] Error logging
- [ ] **Code**: Configure log aggregation
  - [ ] Send logs to LogDNA
  - [ ] Configure retention
- [ ] **Test**: Logging tests
  - [ ] Verify log format
  - [ ] Verify log content
- [ ] **Docs**: Logging documentation
  - [ ] Document log format
  - [ ] Document log queries

### 16.3 Activity Tracker
- [ ] **Code**: Implement AT events
  - [ ] Authorized CIDR events
  - [ ] Floating IP events
  - [ ] Public address range events
- [ ] **Test**: AT tests
  - [ ] Verify event generation
  - [ ] Verify event content
- [ ] **Docs**: AT documentation
  - [ ] Document event types
  - [ ] Document event schema

---

## 17. Compliance and Security

### 17.1 Security Review
- [ ] **Docs**: Complete security questionnaire
- [ ] **Docs**: Document threat model
- [ ] **Docs**: Document security controls
- [ ] **Test**: Security testing
- [ ] **Docs**: Security sign-off

### 17.2 Compliance Review
- [ ] **Docs**: Complete compliance checklist
- [ ] **Docs**: Document data handling
- [ ] **Docs**: Document audit requirements
- [ ] **Test**: Compliance testing
- [ ] **Docs**: Compliance sign-off

---

## 18. Deployment

### 18.1 Staging Deployment
- [ ] **Code**: Deploy to staging environment
  - [ ] Deploy API changes
  - [ ] Deploy RNOS changes
  - [ ] Deploy genctl changes
  - [ ] Deploy SDN changes
- [ ] **Test**: Staging validation
  - [ ] Run smoke tests
  - [ ] Run integration tests
  - [ ] Verify feature flags
- [ ] **Docs**: Staging deployment report

### 18.2 Production Deployment (Beta)
- [ ] **Code**: Deploy to production (beta)
  - [ ] Deploy with feature flags disabled
  - [ ] Enable for allowlisted accounts
  - [ ] Monitor deployment
- [ ] **Test**: Production validation
  - [ ] Run smoke tests
  - [ ] Verify allowlist enforcement
- [ ] **Docs**: Beta deployment report
- [ ] **Docs**: Beta release notes

### 18.3 Production Deployment (GA)
- [ ] **Code**: Deploy to production (GA)
  - [ ] Enable feature flags globally
  - [ ] Remove allowlist restrictions
  - [ ] Monitor deployment
- [ ] **Test**: Production validation
  - [ ] Run smoke tests
  - [ ] Monitor metrics
- [ ] **Docs**: GA deployment report
- [ ] **Docs**: GA release notes

---

## 19. Training and Enablement

### 19.1 Internal Training
- [ ] **Docs**: Create training materials
  - [ ] Feature overview
  - [ ] Technical deep dive
  - [ ] Operator procedures
- [ ] **Docs**: Conduct training sessions
  - [ ] Development team
  - [ ] Operations team
  - [ ] Support team
- [ ] **Docs**: Create knowledge base articles

### 19.2 Customer Enablement
- [ ] **Docs**: Create customer webinar
- [ ] **Docs**: Create demo videos
- [ ] **Docs**: Create blog posts
- [ ] **Docs**: Update product documentation

---

## 20. Post-GA

### 20.1 Monitoring and Support
- [ ] Monitor production metrics
- [ ] Track customer adoption
- [ ] Collect customer feedback
- [ ] Address issues and bugs

### 20.2 Iteration
- [ ] Plan feature enhancements
- [ ] Address technical debt
- [ ] Optimize performance
- [ ] Improve documentation

---

## Summary

**Total Tasks:** ~400+  
**Major Areas:** 20  
**Target Completion:** GA 3Q2026

This checklist provides a comprehensive tracking mechanism for all implementation tasks related to the BYOIP feature. Each task should be marked as complete when finished, with appropriate code reviews, testing, and documentation in place.