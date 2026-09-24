## 12. Terraform Implementation

### New Data Sources

| Data Source | API |
|---|---|
| `ibm_is_floating_ip_profiles` | `GET /floating_ip/profiles` |
| `ibm_is_floating_ip_profile` | `GET /floating_ip/profiles/{name}` |
| `ibm_is_public_address_range_profiles` | `GET /public_address_range/profiles` |
| `ibm_is_public_address_range_profile` | `GET /public_address_range/profiles/{name}` |
| `ibm_is_public_address_range_authorized_cidrs` | `GET /public_address_range/authorized_cidrs` |
| `ibm_is_public_address_range_authorized_cidr` | `GET /public_address_range/authorized_cidrs/{id}` |
| `ibm_is_public_address_range_authorized_cidr_allocations` | `GET …/authorized_cidrs/{id}/allocations` |
| `ibm_is_public_address_range_authorized_cidr_allocation` | `GET …/authorized_cidrs/{id}/allocations/{id}` |

- [ ] **Code** `ibm_is_floating_ip_profiles` — schema, read func, register in provider
- [ ] **Code** `ibm_is_floating_ip_profile` — schema (`name` arg), read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_profiles` — schema, read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_profile` — schema (`name` arg), read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_authorized_cidrs` — schema (filters: `allocation_profile_family`, `availability_mode`), read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_authorized_cidr` — schema (`authorized_cidr_id` arg), read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_authorized_cidr_allocations` — schema (filter: `allocations_resource_type`), read func, register in provider
- [ ] **Code** `ibm_is_public_address_range_authorized_cidr_allocation` — schema (`authorized_cidr_id` + `authorized_cidr_allocation_id` args), read func, register in provider
- [ ] **Test** unit + acceptance tests for all 8 new data sources
- [ ] **Docs** `.html.markdown` page for each new data source

### Updated Resources & Data Sources

**`ibm_is_floating_ip`** — new fields from `POST/GET/PATCH /floating_ips`
- [ ] **Code** add optional `address` arg (ForceNew)
- [ ] **Code** add computed `profile` block, `authorized_cidr` block, `resource_type` attr
- [ ] **Code** update `ibm_is_floating_ip` data source — same computed attrs
- [ ] **Code** update `ibm_is_floating_ips` list data source — add `profile_name` filter
- [ ] **Test** unit: create with `address`; flatten with/without `authorized_cidr`
- [ ] **Test** acceptance: create FIP from BYOIP pool; existing path unbroken; list filter
- [ ] **Docs** update `ibm_is_floating_ip.html.markdown` and `ibm_is_floating_ips.html.markdown`

**`ibm_is_public_address_range`** — new fields from `POST/GET/PATCH/DELETE /public_address_ranges`
- [ ] **Code** add optional `cidr` arg (ForceNew, `ConflictsWith` `ipv4_address_count`)
- [ ] **Code** add computed `ip_version`, `network_prefix_length`, `profile` block, `authorized_cidr` block
- [ ] **Code** flatten new fields in Read and Delete (202) responses
- [ ] **Code** update single + list data sources — same computed attrs; list adds `profile_name` filter
- [ ] **Test** unit: create with `cidr`; `ipv4_address_count` backward compat; conflict validation; Delete 202 flatten
- [ ] **Test** acceptance: create PAR from BYOIP pool; existing path unbroken; list filter
- [ ] **Docs** update `ibm_is_public_address_range.html.markdown` and `ibm_is_public_address_ranges.html.markdown`

**`ibm_is_vpc`** — new `public_address_ranges[*].cidr` from `GET/POST/PATCH /vpcs`
- [ ] **Code** add computed `cidr` to `public_address_ranges` nested schema; populate in Read; update data source
- [ ] **Test** unit + acceptance
- [ ] **Docs** update `ibm_is_vpc.html.markdown`

**`ibm_is_public_gateway`** — new `floating_ip.resource_type` from `GET/POST/PATCH /public_gateways` + `PUT/GET /subnets/{id}/public_gateway`
- [ ] **Code** add computed `resource_type` to `floating_ip` nested schema; populate in Read; update data source
- [ ] **Test** unit + acceptance
- [ ] **Docs** update `ibm_is_public_gateway.html.markdown`

**`ibm_is_virtual_network_interface`** — new `floating_ips[*].resource_type` from `GET/POST/PATCH/DELETE/PUT` VNI + VNI FIP endpoints
- [ ] **Code** add computed `resource_type` to `floating_ips` items schema; populate in Read; update data source
- [ ] **Test** unit + acceptance
- [ ] **Docs** update `ibm_is_virtual_network_interface.html.markdown`

**`ibm_is_bare_metal_server_network_interface`** — new `resource_type` on NI; `profile`, `authorized_cidr`, `resource_type` on FIP sub-endpoints
- [ ] **Code** add computed `resource_type` to NI `floating_ips` items; add `profile` + `authorized_cidr` to FIP sub-resource schema; update data source
- [ ] **Test** unit + acceptance
- [ ] **Docs** update BMS NI resource and data source docs

**Instance network interface** — new `resource_type` on NI; `profile`, `authorized_cidr`, `resource_type` on FIP sub-endpoints
- [ ] **Code** add computed `resource_type` to NI `floating_ips` items; add `profile` + `authorized_cidr` to FIP sub-resource schema; update data source
- [ ] **Test** unit + acceptance
- [ ] **Docs** update instance NI resource and data source docs

### Cross-cutting
- [ ] **Code** register all 8 new data sources in `provider.go`
- [ ] **Code** shared flatten helpers: `flattenFIPProfile`, `flattenFIPAuthorizedCIDR`, `flattenPARProfile`, `flattenPARAuthorizedCIDR`
- [ ] **Code** shared schema helpers for reusable `profile` and `authorized_cidr` nested schemas
- [ ] **Test** unit test each shared flatten helper
- [ ] **Docs** add "Using BYOIP with Terraform" guide; update provider changelog