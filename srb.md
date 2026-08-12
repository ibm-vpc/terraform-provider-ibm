---
template_version: 2025-07-11 # DO NOT CHANGE -- this comes from the copied template.md

## These *MUST* be updated by the mentor as the proposal goes through the SRB process.
state: open # open or closed
status_label: review-follow-up (re-inception completed) # See https://github.ibm.com/cloudlab/srb#status-labels
last_update: 2026-04-26 # Set to today's date whenever the proposal is updated
---
<!-- cSpell:ignore NGUID nsze nuse SPDK nvmf discsrv cephadm nvmevirt bdev -->
<!-- cSpell:ignore Multipathing Multipath cryptsetup blockresize vola vold volu -->
<!-- cSpell:ignore apidocs elbas bmid volid AEOV cgroups VMFS techsupport Sanjeev -->
<!-- cSpell:ignore Ranjan ranjan Prathamesh Kadam RBOS RBAAS HMAC OMAP nvmeof -->
<!-- cSpell:ignore cnode dhchap CHAPKEY traddr trtype trsvcid treq nnaid vsiid volattid -->
<!-- cspell:ignore unmap portid subnqn eflags adrfam sectype subsys autoconnect nvmexpress -->
<!-- cspell:ignore hostnqn nvmes nvmea nvmeu nvmed nvmeck datavol lbaf lbads naid sasr -->
<!-- cspell:ignore targetsubsystem NSID iface keepalive operstate MITM policers -->

# SRB Proposal Template

# 1. Introduction

## 1.1 Project Name

Block storage for VPC Bare Metal

## 1.2 Submitter and Contributors

- TJ Harris - <theoharr@us.ibm.com> - @theoharr
- Paul Gerver - <pgerver@us.ibm.com> - @pgerver
- Anuj Chandra - <anujchan@in.ibm.com> - @anujchandra
- Sanjeev Ranjan - <sanjeev.ranjan@ibm.com> - @sanjeev.ranjan

## 1.3 Mentor(s)

- TJ Harris - <theoharr@us.ibm.com> - @theoharr
- Amartey Pearson - <apearson@us.ibm.com> - @amartey

## 1.4 Offering Manager(s)

- Prathamesh Kadam - <Prathamesh.Kadam1@ibm.com> - @Prathamesh Kadam

## 1.5 API focal(s)

- Karthik Baskaran - <kbaskar@us.ibm.com> - @kbaskar
- Chris Baker - <cbake@us.ibm.com> - @cbake
- Sean Thornton - <sthornt@us.ibm.com> - @sean.thornton
- Michael Rieth - <riethm@us.ibm.com> - @riethm

## 1.6 Review Flow (determined by mentor)

Full

## 1.7 Production Target Date

* 3Q2026 - Beta
* 4Q2026 - Select Availability
* 2Q2027 - GA

# 2. Project Overview

## 2.1 Overview

Currently, RIAS (block storage) volumes in VPC can only be attached to VSIs. VPC Bare Metal
server instances are not able to attach RIAS volumes due to:

1. the lack of client side libraries (nfs and/or librbd) which are present on compute hypervisors
2. the lack of access to networked block storage systems (clients do not have direct access to the
   underlay network from their bare metal server instances. No block storage overlay network exists)

This project adds VPC Bare Metal server attachment support for RIAS volumes of the `sdp` profile.
Volumes of this profile are backed by an Acadia storage cluster, which runs a version of the Ceph
software defined storage product. The IBM Ceph team has created a Ceph NVMe-over-TCP gateway that
will run on the Acadia storage platform and be used to provide NVMe-over-Fabrics access to Ceph
volumes (RBDs). Bare Metal servers, will be able to access volumes using the NVMeoTCP remote
storage protocol. Access to the `nvme-tcp` block storage network will be made available using
link local addresses (LLAs) backed by a hidden virtual network interface (`VNI`). Using the LLAs,
customers will be able to access block storage volumes using the `nvme-tcp` driver from the guest
operating system.

![Ceph Block Data Path Overview](./materials/acadia-bare-metal-AcadiaRack.png)

## 2.1.1 Project History

Prior to the decision to use Acadia as the backing storage for VPC Bare Metal server attachment,
Anuj Chandra and his team were building a solution based on Netapp as the backing storage.
Appreciation and recognition go to Anuj and team for unpeeling some of the challenges around
bringing NVMe to the cloud and laying the ground work for several designs.

The original architecture for this project leveraged the AMD/Pensando data-processing-unit
(DPU), code named Elba, to perform PCIe NVMe controller virtualization, allowing remote block
storage volumes to appear as locally attached PCIe NVMe devices. Due to challenges closing on a
commitment from the vendor to deliver on the outstanding
requirements, the architecture of this project has shifted to avoid additional dependencies
on the DPU.

## 2.2 Consumer Interaction Model

Consumers of VPC Bare Metal servers will now have the ability to attach or detach certain block
storage volumes to their bare metal servers (specifically data volumes of the `sdp` volume profile).
Additionally, users will be able to specify existing or new block storage volumes for attachment
when creating a new bare metal server--identical to VSI creation.  These interactions will be
available through the standard customer facing tools (API, UI, ibmcloud CLI, terraform).

In this proposal, the concept of a *volume attachment protocol* is being added to the regional API.
Volume, VSI, and Bare Metal server profiles will be updated to include supported volume attachment
protocols (see section [3.1.1.9](#3119-updates-to-global-catalog)). Regional services will enforce
volume attachment compatibility anytime a client requests a block storage volume to be attached to
a VSI or Bare Metal server.

The two volume attachment protocols will be:

* `virtio-blk` - Block storage volumes will be presented to the hosted OS as a virtual block
  device (i.e. `/dev/vd*`). This is what VSIs support and will continue to use until
  a future SRB adds `nvme-tcp` attachment support to VSIs.
* `nvme-tcp` - Block storage volumes will be presented as remote NVMe devices. The hosted OS can connect
  to the remote NVMe devices using the native `nvme-tcp` driver.  This will be the only volume
  attachment protocol Bare Metal server instances will support in this proposal.

In the future, an additional volume attachment protocol may be added to support a DPU enabled
block storage option where volumes are presented as local NVMe devices:

* `nvme-pcie` - Block storage volumes will be presented to the hosted OS as local (PCIe attached) NVMe
  devices (i.e. `/dev/nvme*n*`). Out of scope for this proposal.

Additionally, volume profiles will support an `attachment_mode` property to indicate whether volumes
can be attached to a single compute server (`single`) or multiple compute servers concurrently (`multiple`).
This enables multi-attach capabilities for volumes of the `sdp` volume profile when used with bare
metal servers.

### 2.2.1 Identifying NVMe devices

Most, if not all, Bare Metal server profiles include some amount of local NVMe storage. Clients
need a way to determine which NVMe devices are really local vs network attached and for those that
are network attached, which RIAS volumes they map to.

NVMe devices/volumes/namespaces can have up to three globally unique identifiers:

* Namespace GUID - NGUID
* Extended Unique Identifier 64 - EUI64
* Namespace UUID

```bash
root@r18s20nvmeinit:~# nvme id-ns /dev/nvme0n1
NVME Identify Namespace 1:
nsze    : 0x280000
ncap    : 0x280000
nuse    : 0x280000
...
nguid   : bdefe7b77f2f496f85062cef40bf859f
eui64   : 0000000000000000
```

Documentation and tooling will be provided to allow a client to map
a `/dev/nvme*` device to the RIAS block storage volume backing it. Details can be found in
section [3.1.2.4.1](#31241-identifying-nvme-devices).

## 2.3 Requirements and Use Cases

Beta AHA: [BSVPC-345](https://bigblue.aha.io/features/BSVPC-345)

LA/SA AHA: [BSVPC-252](https://bigblue.aha.io/features/BSVPC-252)

GA AHA: [BSVPC-72](https://bigblue.aha.io/features/BSVPC-72)

- Support for data volume attachment to Bare Metal server instances
- Support attaching volumes of all allowed sizes for the 'sdp' profile (1GB - 32TB)
- Support volume performance ranges for the 'sdp' profile (max 64K IOPS, 1GB/s throughput)
- Support up to 100 volume attachments per Bare Metal server instance
  - For use cases which require many more volume attachments, a feature flag (`is-bare-metal-server-volume-attach-limit-value`)
    can be increased on an account-by-account basis.
  - Note: All volume attachments will not be able to drive max IOPS/throughput concurrently
          as they will be subject to network and Acadia cluster limitations
- Support attaching volumes using customer managed encryption keys (BYOK)
- Support volume expansion
- Support volume snapshots
- Support backup service integration including backup policy support for BM servers
- Support encryption in transit
- Continue to support booting from local NVMe storage (existing offering)
- Over the lifetime of a volume that supports multiple attachment protocols (e.g. `virtio-block`
  and `nvme-tcp`), it can be attached to different VSIs and/or different Bare Metal servers (but not
  at the same time). At this time, the only volume profile that will support multiple attachment
  protocols will be `sdp`.
- Support volumes being attached up to 32 Bare Metal servers concurrently (aka volume multi-attach).

## 2.4 Out-of-Scope

- Support for boot volume attachment to Bare Metal server instances
  - Remote boot volumes are not supported on Classic so this is not viewed as a requirement for this
    project. Boot volume support will be added later in a future SRB proposal.
- NVMe storage attachment of VPC Block storage for VSIs
- NVMe storage attachment of Netapp-backed Volumes for Bare Metal servers
- Support for operator driven volume migration (see [SRB/4732](https://github.ibm.com/cloudlab/srb/tree/master/proposals/4732))
  when volumes are attached to a Bare Metal server.
- Extending Bare Metal server metadata service integration to include remote volume connection
  information (i.e. NVMe host nqn, host CHAP, `nvme-tcp` IPs)
- Volume multi-attach support for VSIs
- Support for a direct-to-guest volume attachment
  - In a BYO-hypervisor or Openshift use case, it could be useful to allow guests to directly
    connect to block storage via nvme-tcp rather than via hypervisor passthrough.
- Client side rate limiting of block storage traffic
  - This feature is dependent on additional function being delivered by AMD/Pensando (flow based policers).
    Target is 4Q26. This can be re-evaluated for inclusion in a future SRB proposal.

Items which are out-of-scope for the Beta/LA commits, but will be included in a future GA commit of
this proposal:

- Volume from snapshot instant availability
  - This feature is dependent on additional function being delivered by Ceph. Target is 4Q26. See
    section [3.1.1.10.2](#311102-snapshot-limitations-for-beta-and-la) for more details.

## 2.5 Dependencies, Dependents, and Incompatibilities

- Availability of Acadia storage hardware in regions/zones targeted for this feature
- Bare Metal server Operating System driver support for `nvme-tcp` including support for NVMe
  in-band authentication (i.e. `DH-HMAC-CHAP`).
  - This will require stock images to include the latest versions of `nvme-tcp` drivers for Bare
    Metal enabled images

### 2.5.1 Other VPC Projects Under Development

None

### 2.5.2 Other IBM Projects Under Development

Additional Ceph NVMeoF Gateway features including:

* BYOK and KMIP support
* Namespace masking
* OS Qualification

These features are being developed by the Ceph team and will be delivered as part of the IBM
Storage Ceph 9.1 and 10.0 releases.

### 2.5.3 Third-party

None

### 2.5.4 External and/or Circular

None

### 2.5.5 Brittle or Unusual

None

### 2.5.6 Incompatibilities

None

### 2.5.7 Databases

None

## 2.6 Integration, Deployment, and Operations Considerations

### 2.6.1 Acadia platform updates

The Acadia platform will host new services required for NVMeoF connectivity, specifically the
the NVMeoTCP Gateways. An Acadia cluster (T2) is comprised of 20 servers, known as `OSD` nodes.
Each `OSD` node will host an instance of the NVMeoTCP Gateway. The Acadia platform deployer will be
enhanced to optionally deploy the new NVMeoF services as resilient containers.

A new feature property, `nvmeof`, will be added to the [ringmaster](https://github.ibm.com/genctl-acadia/ringmaster)
github repo which hold the configuration-as-code for all Acadia storage clusters. When enabled,
the Acadia deployment tooling will deploy the gateway where appropriate.

### 2.6.2 Client enablement

A feature flag `is-bare-metal-server-volumes-phase1-allowlist` along with maturity flag
`is-bare-metal-server-volumes-phase1-maturity` will be used to allow list
accounts so they can be used to provision Bare Metal server instances with network attached block
storage. Another feature flag, `is-bare-metal-server-volume-attach-limit-value` will control how many
volumes a Bare Metal server is able to attach.

## 2.7 Big Rules

### 2.7.1 Consistent Customer Experience

This SRB introduces customer facing changes, so there will be a change/impact to the customer
experience. The goal is to align Bare Metal provisioning and block storage attachment (with `nvme`)
experiences with what we currently offer for VSIs (with `virtio-blk`), without breaking any existing
behaviors.

The volume attachment experience for Bare Metal servers will deviate from the VSI volume attachment
experience.
Access to block storage volumes via `nvme-tcp` from Bare Metal servers will require customer action.
By 'attaching' a volume to a Bare Metal server the cloud orchestration will make that volume available
on the Bare Metal servers link local network. Then the customer must issue an `nvme connect` command
to complete the
attachment, which is more akin to mounting a file share ... or connecting to an iSCSI volume in
the IBM classic infrastructure.

### 2.7.2 VPC Big Rules

- [x] All components must be able to be upgraded live (without service degradation/outage) and without
   manual intervention
  - Control plane components will be upgraded by `razee`
  - Acadia platform upgrades orchestrated by `cephadm`
- [x] All components must present usage metrics
- [x] All components must use the [VPC/NG monitoring framework](https://github.ibm.com/cloudlab/srb/tree/master/architecture/telemetry/monitoring)
- [x] All components must use existing VPC/NG control planes -- no new control planes
- [x] All components must be controlled via well-defined APIs
- [x] All components must provide functional parity across all modalities (UI, CLI, API, Terraform)
- [x] All components must provide functional parity across all targeted platform architectures (x86
   and Z)
  - Note: Z Bare metal does not exist
- [x] All customer-facing GA features must interoperate with all other customer-facing GA features
- [x] All customer-facing UIs and CLIs must use the customer-facing VPC API to operate on VPC/NG resources
- [x] All customer-facing APIs must adhere to the [API handbook guidelines](https://pages.github.ibm.com/CloudEngineering/api_handbook/)
- [x] All regional APIs must be the same in all regions (this implies equivalent functionality in
   all regions)
- [x] All components must perform authorization using the [VPC/NG IAM services](architecture/platformIntegration/genesis-iam)
- [x] All components must be internally stateless and recreate their runtime state from external sources
- [x] Any caching or shadowing of state must never result in an API consumer observing inconsistent or
   incorrect behavior
- [x] All hardware components must use the
   [VPC/NG diagnosability framework](https://github.ibm.com/cloudlab/srb/tree/master/architecture/telemetry/diagnostics)
   and implement the on-demand diagnostic API
- [x] All services must present an API for get/set/watch of state
- [x] All connections between components must be mutually authenticated, and encrypted using AES128
   or higher
- [x] All components must integrate into the existing operations tooling for deployment and
   maintenance

### 2.7.3 Cloud Security Baseline Requirements

- [x] Authentication & access control: Enforce authentication and access control on all public
  interfaces and between service components.
- [x] Secrets management: Store and manage secrets securely.
- [x] Least privilege operation: Always operate at the lowest level of privilege required to execute
  tasks effectively. Don't run processes as root.
- [x] Compute isolation: Isolate critical processes and untrusted code at the compute level.
- [x] Network isolation: Group related systems into network segments that are isolated from
  unrelated systems, adding additional isolation for high risk systems.
- [x] Encryption at rest & secure deletion: Encrypt all sensitive data at rest, all customer data is
  sensitive.
- [x] Encryption in transit: Encrypt all sensitive data "on-the-wire", all customer data is
  sensitive.
- [x] Encrypt to standards: Use AES-256 for symmetric encryption, use recommended TLS cipher suites,
  follow NIST 800-131A for anything else.
- [x] Logging & non-repudiation: Log CRUD operations and forward all log entries off box to prevent
  tampering.
- [x] Logging of credentials: Don't include any type of credentials in logs.
- [x] Injection prevention: Prevent injection vulnerabilities by always validating user input, using
  language-specific safe libraries and APIs, and perform periodic application vulnerabilities scans.
- [x] Secure development practices: Build security into all phases of the software development
  lifecycle.
- [x] Update & secure configuration: Keep all aspects of the service updated and securely
  configured.

## 2.8 Assumptions, Risks, and Constraints

None

# 3. Architecture, Interfaces, and Impact

## 3.1 Architecture and Interfaces

This project incorporates deliverables from two IBM groups:

* IBM Storage
* IBM Cloud IaaS

IBM Storage:

- [Ceph NVMeoTCP gateway](https://github.com/ceph/ceph-nvmeof)

The Ceph NVMeoTCP gateway makes heavy use of code from the Storage Performance Development
Kit (aka [SPDK](https://github.com/spdk/spdk) open source project. All required enhancements
have been contributed upstream by IBM Research and/or the IBM Ceph teams.

IBM Cloud VPC IaaS:

- Updates to Regional and Zonal control planes
- RIAS API updates (primarily to `/volumes` and `/bare_metal_servers`)
- Acadia platform integration of Ceph NVMeoTCP gateway
- Monitoring and observability updates to IaaS tooling to support Ops/SRE
- Local NVMeoF discovery service
- Network underlay support for a new `nvme-tcp` VRF
- Network overlay/SDN support for using endpoints on the `nvme-tcp` VRF in virtual network
  interfaces (VNIs)

### 3.1.1 Architecture and Technical Design

#### 3.1.1.1 High Level Overview

To date, the Acadia team has built a new storage backend based on Ceph. The block storage provided
by Acadia can be consumed by VSIs. For details see SRBs [1793](https://github.ibm.com/cloudlab/srb/tree/master/proposals/1793)
and [3358](https://github.ibm.com/cloudlab/srb/tree/master/proposals/3358). The following diagram
depicts the primary components used to build this solution, namely:

* Acadia Storage Platform
  * Ceph cluster with additional services and automation purpose built for IaaS
  * Monitor (MON) nodes used for primarily for control path
  * Object Storage Daemon (OSD) nodes used primarily for data path
* Compute Nodes
  * Same hypervisor virtualization stack including qemu/libvirt
  * Ceph client library for block storage access (`librbd`)

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-vsi-overview.png alt="Acadia VSI Data Path Overview" width="800">

To allow block storage connectivity to Bare Metal Server instances, a new storage protocol must be
added to Acadia/Ceph.  The NVMe-over-Fabrics / `NVMeoF` (specifically NVMe-over-TCP / `NVMeoTCP`),
storage protocol was chosen for a number of reasons including:

* The DPUs present in Bare Metal Severs aren't powerful enough to run the Acadia/Ceph client
  software stack (`librbd`/`librados`) in a performant way
* The Acadia/Ceph client software stack (`librbd`/`librados`) is not installed by default in most
  cloud enabled OS images
* Industry momentum and adoption behind `NVMe` and `NVMeoTCP` (driver support is standard in most
  cloud enabled OS images)
* Ubiquity of TCP in the cloud infrastructure

The following diagram depicts the primary components used to build this solution, namely

* Ceph NVMeoF Gateway
  * New container running on the OSD nodes
  * Presents Ceph block storage (`RBDs`) as NVMeoF namespaces
  * Leverages the `SPDK` project to provide NVMeoF target functionality

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-bm-overview.png alt="Acadia BM Data Path Overview" width="800">

Shifting from the native Ceph block storage protocol (`rbd`) to NVMeoTCP presents some new
challenges.

* Shifting from a one-to-many storage protocol to a point-to-point presents new performance
  bottlenecks and failure modes that must be addressed by this design.
* The Ceph `msgr2` protocol terminates between the NVMeoTCP gateway and the Acadia cluster, meaning
  data in-flight between the DPU and the gateway is *not* encrypted by `msgr2`. Need to use a
  different encryption in flight technology (see this [section](#3115-authentication-and-authorization)).

A more detailed diagram depicts high level provisioning steps for an example with two Bare Metal servers
and four volumes:

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-nvme-tcp-provisioning.png alt="BM Provisioning flows" width="800">

Subsequent sections will go into details on how the challenges of introducing a new block storage
protocol are being addressed.

#### 3.1.1.2 Network Design

#### 3.1.1.2.1 New VRF

The `nvme-tcp` networking concepts build on top of the VNI and FaaS network proposals delivered
previously with [SRB/1967](https://github.ibm.com/cloudlab/srb/tree/master/proposals/1967) and
[SRB/4734](https://github.ibm.com/cloudlab/srb/tree/master/proposals/4734).

To support `nvme-tcp` traffic, a new VRF will be added. Using a new VRF allows:

* the new block traffic to be segregated from existing file traffic (NFS and upcoming SMB).
* a unique IP address space for `nvme-tcp` so that there are no IP conflicts with file related services

Storage nodes (aka OSD nodes) will advertise their fabric IP (e.g. 10.x.x.x) as a /32 IP addresses
into the `nvme-tcp` VRF.

To support this new VRF, SDN will add a new VXLAN Network ID (aka VNI not to be confused with
a Virtual Network Interface)) of `4073` (as opposed to the the `FaaS` network which uses a VNI
of [`4072`](https://github.ibm.com/cloudlab/srb/tree/master/proposals/4734#3112-network-design)). To
ensure `nvme-tcp` routes are unique, an RD/RT of `64900:4073` will be used for NG regions and an
RD/RT of `64901:4073` will be used for NGDC regions.

Lastly, the EVPN prefix routes advertised from the EENs must reach the storage nodes in the storage
platform in order for them to route traffic back onto the `nvme-tcp` overlay network.
To accomplish this, the STORs will have an export policy for routes advertised with a community tag
of `65201:42` to announce them to the storage nodes. The SDN code on the EEN will attach this community
tag to all `nvme-tcp` routes advertised by them. The routes advertised by SDN are controlled through
a new property in [`platform-inventory`](https://github.ibm.com/cloudlab/platform-inventory) named
`NVME_EVPN_ADDRESS_POOL_CIDR`. The EVPN prefix routes are additionally required to be reflected
into the RIAS globals. EVPN prefix allocations will be stored as part of the
[`acadia-baremetal-storage-workspace`](https://github.ibm.com/genctl/acadia-baremetal-storage-workspace/wiki/Block-for-Barmetal-NVME-CIDR-allocations).

> NOTE:
> During development, the storage team will ensure FaaS and `nvme-tcp` EVPN prefixes do not overlap
> for any environment. Since the FaaS and `nvme-tcp` networks on are different VRFs, there shouldn't
> be any issues with overlapping routes. However, the SDN team needs to add additional code in support
> of this configuration and validate overlapping EVPN prefixes can be supported. This restriction
> will be re-evaluated prior to beta commit.

#### 3.1.1.2.2 Aggregate VNI / Storage VNI Gateway

For this project and the Denali-SMB project SDN will be adding support for a `StorageVNIGateway`,
also known as an aggregate VNI. This type of VNI will be used to provide network connectivity to
block storage from a Bare Metal server using link local addresses. [RFC/3927](https://datatracker.ietf.org/doc/html/rfc3927)
reserves the `169.254.0.0/16` CIDR for link local addresses. A subset of this CIDR will be used for
`nvme-tcp` endpoint mappings, `169.254.172.0/22`. This allows storage and networking to map up to 1024
storage `nvme-tcp` underlay IPs to LLAs without consuming any of the customer's BYOIP space. Customers
can then use the LLAs to connect to their remote block storage.

The `StorageVNIGateway`s used to provide storage connectivity will be not viewable
by customers from the public API (`GET /virtual_network_interfaces`). The lifecycle of `StorageVNIGateway`s
will be managed entirely by the storage and network control planes, and tied to whether a Bare Metal
server has block storage volumes attached.

The aggregate VNI will carry a mapping of all required `nvme-tcp` IPs to customer LLAs.
As volumes are attached/detached, the aggregate VNI's IP mapping
will be updated by the storage control plane. Each volume will have 4 `nvme-tcp` IPs (4 paths). The
Acadia storage cluster and NVMeoF gateway group to which a volume is made available for `nvme-tcp`
connection will influence which storage IPs must be mapped in the aggregate VNI. Customer's have no
visibility into Acadia storage clusters or NVMeoF gateway groups so the mapping must be managed
by the storage control plane.

Aggregate VNIs used to provide `nvme-tcp` connectivity will not span zones. In other words, the Bare
Metal server and any block storage volumes its attaching must be in the same zone.

A /24 EVPN prefix will be configured in the `StorageVNIGateway` so that SDN can allocate an `nvme-tcp`
EVPN IP in order to route traffic back from storage to the bare metal server. The aggregate VNI
will only be used by the Bare Metal server instance which it is bound to, so a /32 would be sufficient.
However, SDN only supports /24 so storage will continue sending /24 EVPN prefixes.

#### 3.1.1.2.2.1 Aggregate VNI / Storage VNI Gateway Security

The `StorageVNIGateway` design provides strong network isolation through the following mechanisms:

* Each `StorageVNIGateway` is bound exclusively to a single Bare Metal server.
* Link-local addresses (169.254.172.0/22) are strictly scoped to the individual Bare Metal server and
  cannot be used for east-west traffic between servers or cross-tenant communication. These addresses
  are only valid within the context of the specific server's storage VNI.
* SDN enforces strict source validation for all storage link-local addresses (LLAs). Any traffic
  originating from an LLA that doesn't match the expected Bare Metal server's aggregate VNI will be
  dropped, preventing spoofing or lateral access attempts.
* Inbound Network ACL rules enforce only NVMe Discovery traffic (port 8009) on the Discovery service
  IP (`169.254.172.1`) within the LLA subnet.
* Inbound Network ACL rules enforce only `nvme-tcp` traffic (port 4420) to other IPs on the LLA subnet.

Following is a diagram summarizing the security controls in place for the `StorageVNIGateway` design:

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-network-controls.png alt="LLA Network Controls" width="800">

#### 3.1.1.2.3 nvme-tcp connectivity to Bare Metal OS

Each Bare Metal server will have a single `nvme-tcp` `StorageVNIGateway` / aggregate VNI to allow for
as-needed mapping of link local addresses to backend storage IPs.  Each volume attachment will
require 4 IPs from the Bare Metal server's LLA space, one for each `nvme-tcp`
data path. As volume attachments are being processed, the storage control plane will map, unmap, or
use existing IPs from the aggregate VNI.

Some volume attachments may be able to share IPs if the volumes exist on the same Acadia
storage cluster and NVMeoF gateway group. Where possible, the storage control plane will reuse IPs
in the aggregate VNI to reduce the overall number of IPs required. Volume attachment API responses
will be extended to include new properties which customer's must use to connect to remote `nvme-tcp`
block volumes and a new status to indicate further action is required by the customer:

| Volume Attachment Property | Example | Purpose |
| --- | --- | --- |
| NVMe Subsystem | `nqn.2016-06.io.spdk:cnode1.group0` | NVMe subsystem which is presenting the volume as an NVMe namespace. |
| `device.id` | `26fc81e15476447f91d8df71465e9884` | NVMe UUID to assist in mapping a volume to an NVMe namespace. |
| List of 4 IPs | `169.254.172.1[0-3]` | Destination IP endpoint, from the customer's BYOIP space, for each `nvme-tcp` data path. |
| `status` | `available` | Distinguish from VSI volume attachments which eventually go to an `attached` status. Customer must now initiate an NVMe connect from their OS. |

Example attachment commands for a linux based OS:

```bash
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.10 -s 4420 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.11 -s 4420 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.12 -s 4420 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.13 -s 4420 --dhchap-secret ${CHAPKEY} --concat
```

The standard `nvme-tcp` port as defined by [IANA](https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=4420)
will be used provide `nvme-tcp` connectivity to Bare Metal servers. The standard port number is
`4420` and will be inferred when omitted from the `nvme connect` command. The following commands
are functionally equivalent to the previous commands which explicitly use port `4420`:

```bash
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.10 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.11 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.12 --dhchap-secret ${CHAPKEY} --concat
nvme connect -t tcp  -n nqn.2016-06.io.spdk:cnode1.group0 -a 169.254.172.13 --dhchap-secret ${CHAPKEY} --concat
```

#### 3.1.1.2.3.1 nvme-tcp connect rate limiting

Establishing `nvme-tcp` connections will be relatively rare; once the `nvme-tcp` connection has been
established the host and target will negotiate a number of long living TCP sessions to create for
processing I/O. Typically, new connections will only be established when
adding a new volume to a Bare Metal server or re-establishing connectivity after a Bare Metal
server reboot. The [Ceph NVMeoF Gateways](#3113-ceph-nvmeof-gateway-overview) do not have any controls
to limit the number of connections or frequency that can be established.
So to prevent abuse, rate limits for new connection establishments will be enforced by
SDN. The rate limit will be two connections per second per LLA.

#### 3.1.1.2.4 NVMe Discovery Service

To streamline the process of connecting to remote block volumes, a bespoke NVMe discovery service
will be added to the DPUs. A well-known/static link local address will be reserved for NVMe discovery
(e.g. `169.254.172.1`), similar to the metadata service. This will allow all Bare Metal clients to
use the same instructions and endpoint to connect to remote block storage.

Clients will then be able to discover any and all NVMe subsystems (which expose volumes) to connect
to. For linux based OSs, the command would be:

```bash
nvme discover -t tcp -a 169.254.172.1 -s 8009
```

And would provide output in this format, where the `traddr`s are IPs from the customer's BYOIP space:

```bash
Discovery Log Number of Records 8, Generation counter 12
=====Discovery Log Entry 0======
trtype:  tcp
adrfam:  ipv4
subtype: nvme subsystem
treq:    not required
portid:  0
trsvcid: 4420
subnqn:  nqn.2016-06.io.spdk:cnode1.group0
traddr:  169.254.172.10
eflags:  none
sectype: none
=====Discovery Log Entry 1======
trtype:  tcp
adrfam:  ipv4
subtype: nvme subsystem
treq:    not required
portid:  1
trsvcid: 4420
subnqn:  nqn.2016-06.io.spdk:cnode1.group0
traddr:  169.254.172.11
eflags:  none
sectype: none
...
```

This allows clients to gather remote block volume connection information from the Bare Metal server
itself, rather than extracting the same information from public RIAS API calls.

In addition, the `nvme-cli` tool set also allows a user to automatically connect to all discovered NVMe
subsystems using a single command:

```bash
nvme connect-all -t tcp -a 169.254.172.1 -s 8009 --dhchap-secret ${CHAPKEY} --concat
```

Clients can then add this configuration, saved in `/etc/nvme/discovery.conf`, in tandem with the
[`nvmf-autoconnect service`](https://github.com/linux-nvme/nvme-cli/blob/master/nvmf-autoconnect/systemd/nvmf-autoconnect.service.in)
to connect all remote block volumes during boot. This can be further extended to cloud-init
vendor data to assist customers in automating volume discovery and connection.

Following is a diagram which depicts the high level interactions between zonal control place, Acadia
storage, and the discovery service:

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-nvme-discovery-service.png alt="NVMe Discovery Service" width="800">

Adding this service to the DPUs will require:

* Storage team to build a new discovery service which can:
  * be dynamically configured with a set of subsystems names and IPs
  * respond to NVMe discovery requests (adhering to the NVMe specification)
* A new snap package to be run on the DPU
* SDN support for an additional link local address and rules to forward requests to the DPU based
  discovery service
* Rate limits enforced by SDN and the discovery service to prevent abuse
  * SDN traffic to the discovery service is rate limited based on connection and packet limits. Only
    one new connection per second will be allowed.
  * Discovery service will only support a single discovery request at a time. Additional requests will
    be queued and processed when the previous request is complete.

#### 3.1.1.2.4 NVMe Discovery Service Security

The NVMe discovery service is designed with security as a core principle:

* The discovery service only returns NVMe subsystems that have already been authorized for the specific
  Bare Metal server making the request. Discovery does not grant new access—it merely reflects the
  current control plane authorization state for that server.
* The discovery service is read-only and informational. It cannot be used to gain access to volumes
  that haven't been explicitly attached through the VPC control plane APIs. All access control decisions
  are made by the control plane before discovery results are returned.
* Multiple layers of protection prevent discovery service abuse:
  * SDN enforces connection and packet rate limits on traffic to the discovery service endpoint. Only
    one new connection per second will be allowed.
  * The discovery service itself processes only one discovery request at a time per Bare Metal server,
    with additional requests queued
  * These limits prevent the discovery service from being used as a reconnaissance tool or for
    denial-of-service attacks

This design ensures that the discovery service improves usability without introducing new security
risks or attack surfaces.

#### 3.1.1.3 Ceph NVMeoF Gateway Overview

The Ceph NVMeoF Gateway is an [open source project](https://github.com/ceph/ceph-nvmeof) with the
goal being to add NVMeoTCP protocol support for `RBDs` (block storage).  Like other projects in Ceph,
all code (fixes, enhancements, etc) will be delivered upstream first and then built into downstream
releases by the Ceph product team.

The gateway is made up of three main components:

![Ceph NVMeoF Gateway overview](./materials/acadia-bare-metal-gw-container.png)

* Configuration Daemon
  * gRPC server to configure NVMeoF objects (subsystems, namespaces, listeners) and map them to
    specific NVMeoF initiators
  * Configuration stored persistently in the Ceph cluster
* Storage Performance Development Kit (`SPDK`)
  * Based on `DPDK`
  * Provides an NVMeoF target implementation
  * Support for multiple backend storage technologies (e.g. Ceph RBD, iSCSI, etc) using a block
    device (aka `bdev`) abstraction
  * Support for volume level QoS limits (both IOPs and throughput) (see section [3.1.1.6](#3116-volume-level-qos))
* `librbd`
  * Ceph block storage client library
  * Provides encryption in flight using the `msgr2` protocol

The gateway is purely software packaged as a container with no hardware dependencies (`NOTE: there
are some cpu instruction set optimizations in SPDK/DPDK but these are not strictly required`).

The gateways will run on the Acadia platform OSD nodes, which will put extra demand on the
OSD nodes in terms of

* CPU (performing NVMeoF-to-RBD translation, encrypting/decrypting packets for encryption-in-flight
  and BYOK)
* Memory (in-memory SPDK structures for NVMeoF, TCP I/O buffers)
* Network (intermediate hop in the data path. doubling bandwidth for NVMeoF traffic)

Measurements and modeling are needed to ensure Acadia OSD nodes have capacity to meet the Bare Metal
volume demand.

By default, the gateways will not expose any block storage as NVMeoF storage. The zonal control
plane will instruct the gateways which volumes/RBDs to map to what Bare Metal server and DPU.

#### 3.1.1.4 NVMeoF Multi Path Overview

The Ceph `RADOS` protocol is one-to-many, meaning that a client connects to *and* performs I/O
directly to the OSD nodes (which persistently store the data). The NVMeoF protocol is point-to-point
so the NVMe initiators must connect through a gateway to access data. If a gateway crashes, the
initiator can no longer access storage via that gateway or path. To provide data availability
during planned or unplanned events which take gateways offline, the Ceph NVMeoF gateways support
'groups'. Each gateway within the same gateway group share the same NVMeoF configuration data
(subsystem, namespaces, allowed hosts). This allows the gateways in a group to present the same
set of volumes to initiator, so that each gateway is seen as a additional path to the target
device/volume/namespace.

When the Acadia platform cluster is configured, the gateway groups will be defined in sets of 4s.
This gives each initiator 4 NVMeoF paths to every device/volume/namespace. NVMe supports an asymmetric
access model or ANA (asymmetric namespace access), which allows some paths to be active/optimized and
others to be non-optimized or inaccessible. Ceph supports a single active client per RBD/volume, so
only one path/gateway will be active per volume; the remaining three paths/gateways will be inaccessible
(until a failover event occurs). The active/optimized paths within a gateway group will rotate between
the attached volumes so that the overall I/O workload is evenly spread across all gateways within
a group.

Ceph is able to notify the Bare Metal host which paths are optimized or inaccessible using
asynchronous event notifications, or AENs. This allows the Bare Metal servers `nvme-tcp` driver to
select the most appropriate data path. Following is an example of what a Bare Metal server with a
linux based OS might see with two data volumes connected:

```bash
# first volume
root@r17s33nvmeinit:~# nvme list-subsys /dev/nvme1n1
nvme-subsys1 - NQN=nqn.2016-06.io.spdk:cnode1.group0
               hostnqn=nqn.2014-08.org.nvmexpress:uuid:73d3b385-e47f-4811-826c-1bef7ed924c2
\
 +- nvme1 tcp traddr=169.254.172.171,trsvcid=4420,src_addr=169.254.172.177 live optimized
 +- nvme2 tcp traddr=169.254.172.173,trsvcid=4420,src_addr=169.254.172.177 live inaccessible
 +- nvme3 tcp traddr=169.254.172.175,trsvcid=4420,src_addr=169.254.172.177 live inaccessible
 +- nvme3 tcp traddr=169.254.172.177,trsvcid=4420,src_addr=169.254.172.177 live inaccessible

# second volume
root@r17s33nvmeinit:~# nvme list-subsys /dev/nvme1n2
nvme-subsys1 - NQN=nqn.2016-06.io.spdk:cnode1.group0
               hostnqn=nqn.2014-08.org.nvmexpress:uuid:73d3b385-e47f-4811-826c-1bef7ed924c2
\
 +- nvme1 tcp traddr=169.254.172.171,trsvcid=4420,src_addr=169.254.172.177 live inaccessible
 +- nvme2 tcp traddr=169.254.172.173,trsvcid=4420,src_addr=169.254.172.177 live optimized
 +- nvme3 tcp traddr=169.254.172.175,trsvcid=4420,src_addr=169.254.172.177 live inaccessible
 +- nvme3 tcp traddr=169.254.172.177,trsvcid=4420,src_addr=192.168.45.177 live inaccessible
```

The NVMeoF data path for Bare Metal server attachment will be:

* Each Bare Metal server will have two physical network connections, via the Elba card,
  to the ToR switches
* NVMeoTCP traffic will traverse the leaves/spine of the underlay network
* Ceph OSD nodes have two dual port CX6 cards which will be connected to the ToR switches
* The Bare Metal server OS will see 4 paths to each device/volume/namespace (via a VNI)
* Only a single path will be active/optimized per device/volume/namespace. Any changes in the NVMe
  path states will be announced to the host via NVMe asynchronous event notifications (AEN)

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-Multipathing.png alt="NVMeoF Multipath" width="800">

Each gateway within a group will perform similarly in response to I/O requests. The data for volumes
(RBDs) is persistently stored by striping across all OSD nodes in the Ceph cluster. This means that
each gateway will likely need to communicate with the other OSD nodes while processing I/Os.

#### 3.1.1.5 Authentication and Authorization

Traditionally, Ceph clients provide a [username and key](https://docs.ceph.com/en/latest/rados/operations/user-management/)
to authenticate with the Ceph cluster. Once the client has been authenticated and authorized, data
is exchanged using the [msgr2](https://docs.ceph.com/en/quincy/rados/configuration/msgr2/) protocol.
The `msgr2` protocol optionally supports a 'secure' mode, in which all data is encrypted. Acadia
platform clusters are configured to require all clients to use 'secure' mode.

Inserting the Ceph NVMeoF gateways into the data path means that the Ceph authentication,
authorization, and encryption-in-flight are only between the gateways and Ceph cluster.

```json
Bare Metal server OS <-> NVMeoTCP <-> Ceph NVMeoTCP Gateway <-> RADOS+msgr2 <-> Ceph Cluster
```

The NVMeoTCP datapath between the Bare Metal server's OS and the Ceph NVMeoTCP gateways need similar
controls in
place so that all data is secured and only visible to Bare Metal servers which should have access.
By default, NVMeoTCP traffic is not encrypted-in-flight and the endpoints (initiator/target) have
very weak authentication (a string based NVMe Qualified Name, NQN, check which could easily be
impersonated).
The [NVMe specification](https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-2.0c-2022.10.04-Ratified.pdf)
does support in-band authentication for stronger authentication and secure channels which offer
authentication, authorization, and encryption. SPDK supports both in-band
authentication and secure channels.

To achieve the required security controls, Nvme in-band authentication will be used to provide:

* encryption-in-flight
* authentication and authorization by only allowing connections from allowed initiators

Each Bare Metal server will use a unique host NQN based on a UUID generated at Bare Metal server
creation (or during a backfill process for previously provisioned Bare Metal servers).
The host NQN will be used to allow-list a Bare Metal server to the set of volumes attached.
Customer's will be able to determine the host NQN value from the RIAS public API (i.e. GET /bare_metal_servers/{id}),
and then use that as an argument to NVMe connect commands. Also, they can choose to store the host NQN
in the guest OS's NVMe configuration (i.e. of the `/etc/nvme/hostnqn` for linux based OSs).

> TODO: See if this can be provided to the guest OS via vendor data

Each Bare Metal server will also be provided with a unique Challenge Handshake Authentication Protocol
(CHAP) key. The gateways will be configured to require that the Bare Metal servers provide the CHAP
key and identity (host NQN) over a TLS 1.3 channel for any `nvme-tcp` connections.

Per section
8.3.5.5.2 of the [NVMe base specification](https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-Revision-2.3-2025.08.01-Ratified.pdf),
three hash functions are supported by the CHAP protocol: `sha256`, `sha384`, and `sha512`. `sha256`
is the only mandatory hash function, so all CHAP keys generated by the control plane will use the
`sha256` hash function.

In the RIAS public API, the CHAP key will be returned using an extensible property name, storage
access secret. The storage access secret is returned encrypted using a customer-provided SSH public
key of type `rsa` (referenced via `public_key`) and encoded in base64 format. The Bare Metal
server API response
includes a `storage_access` object containing the encrypted secret, the encryption key reference,
and `created_at` and `rotated_at` timestamps. If no valid SSH public key was provided during
server creation, the `storage_access` will be omitted from the Bare Metal
server API response. To remedy this, customers can choose to reinitialize the Bare Metal server
with a valid SSH public key or rotate the storage access secret with a valid SSH public key.

For new volume attachments, the service will need the CHAP key to add more subsystems. To store this
securely, the service will use a root key from Key Protect of the service account to store the key
in CRD as a wrapped data encryption key (WDEK).

#### 3.1.1.5.1 Storage Access Secret Rotation

A new API endpoint (`POST /bare_metal_servers/{id}/storage_access/rotate`)
will be introduced to enable
storage access secret rotation for Bare Metal servers. It will optionally accept a new SSH public key
reference to use for encrypting the `storage_access_secret`. To prevent misuse and ensure system stability,
a rate limit of one rotation per hour will be enforced per Bare Metal server. The Bare Metal server
API response will include a `storage_access.rotated_at` field, providing visibility into when
the storage access secret was last rotated for audit and operational purposes.

<!-- markdownlint-disable MD033 -->
<img src=./materials/nvme-chap-rotation-sequence.png alt="NVME Chap Rotation" width="700">

#### 3.1.1.5.2 Authentication and Authorization Security Considerations

***Strength of the CHAP key***

To guess a CHAP key, an attacker has two options:

* Crack the Diffie-Hellman exchange or
* Brute force attack the `sha256` HMAC.

The Diffie-Hellman exchange uses a group size of 2048-bits, which is considered secure against attacks
using current technology. It would take a large GPU cluster with 1000s of nodes millions of years to
crack it.
The `sha256` HMAC is a 256-bit key, which is considered strong and secure
for most use cases. It provides a high level of protection against brute-force attacks and is widely
used in industry-standard authentication protocols. The key is generated using a cryptographically
secure random number generator, ensuring that it is unique and unpredictable. The use of a 256-bit key
ensures that the key is resistant to brute-force attacks, even with the most powerful computing
resources available today. An attacker would need thousands to millions of years with current technology.

***Authentication Credential Lifecycle and Usage***

The following sequence diagram illustrates the complete lifecycle of `nvme-tcp` authentication
credentials and their usage:

```bash
┌─────────────┐         ┌──────────────┐         ┌─────────────┐         ┌──────────────┐
│  Customer   │         │ VPC Control  │         │   NVMeoF    │         │    Ceph      │
│             │         │    Plane     │         │  Gateway    │         │   Cluster    │
└──────┬──────┘         └──────┬───────┘         └──────┬──────┘         └──────┬───────┘
       │                       │                        │                       │
       │ 1. Create BM Server   │                        │                       │
       │──────────────────────>│                        │                       │
       │                       │                        │                       │
       │                  Generate Host NQN             │                       │
       │                  (UUID-based)                  │                       │
       │                       │                        │                       │
       │                  Generate CHAP Key             │                       │
       │                  (sha256, 256-bit)             │                       │
       │                       │                        │                       │
       │<─ Return encrypted ───│                        │                       │
       │   CHAP key via SSH    │                        │                       │
       │                       │                        │                       │
       │ 2. Attach Volume      │                        │                       │
       │──────────────────────>│                        │                       │
       │                       │                        │                       │
       │                       │ Configure Subsystem    │                       │
       │                       │ + Host NQN allowlist   │                       │
       │                       │───────────────────────>│                       │
       │                       │                        │                       │
       │                       │                        │  Create RBD Image     │
       │                       │                        │──────────────────────>│
       │                       │                        │                       │
       │<─ Volume attachment ──│                        │                       │
       │   with LLA IPs        │                        │                       │
       │                       │                        │                       │
       │ 3. NVMe Discovery     │                        │                       │
       │ (nvme discover)       │                        │                       │
       │───────────────────────┼───────────────────────>│                       │
       │                       │                        │                       │
       │<─ Subsystem list ─────┼────────────────────────│                       │
       │   (authorized only)   │                        │                       │
       │                       │                        │                       │
       │ 4. NVMe Connect       │                        │                       │
       │ (with CHAP secret)    │                        │                       │
       │───────────────────────┼───────────────────────>│                       │
       │                       │                        │                       │
       │                       │        TLS 1.3 + CHAP Authentication           │
       │                       │                        │                       │
       │<─ Connection established ──────────────────────│                       │
       │   (encrypted channel) │                        │                       │
       │                       │                        │                       │
       │ 5. Rotate Secret      │                        │                       │
       │──────────────────────>│                        │                       │
       │                       │                        │                       │
       │                  Generate new CHAP             │                       │
       │                       │                        │                       │
       │                       │ Update Gateway Config  │                       │
       │                       │───────────────────────>│                       │
       │                       │                        │                       │
       │                       │      Gateway Restart   │                       │
       │                       │                        │                       │
       │<─ New encrypted ──────│                        │                       │
       │   CHAP key            │                        │                       │
       │                       │                        │                       │
       │ 6. Reconnect with     │                        │                       │
       │    new CHAP           │                        │                       │
       │───────────────────────┼───────────────────────>│                       │
       │                       │                        │                       │
       │<─ Connection re-established ───────────────────│                       │
       │                       │                        │                       │
```

**Key Security Points:**

* Host NQN is generated once at BM creation and remains stable
* CHAP keys are unique per Bare Metal server and encrypted with customer's SSH public key
* All nvme-tcp connections require both Host NQN (identity) and CHAP key (authentication) over TLS 1.3
* Secret rotation is customer-initiated and rate-limited to once per hour
* Gateway configuration updates are atomic and connections persist across rotations

#### 3.1.1.6 Volume level QoS

NVMeoF Gateways support QoS enforcement, both in terms of both IOPs and throughput, on a per-volume
basis. During a volume attachment workflow, IOPs change workflow, or throughput change workload the
cloud control plane will specify the volume's new QoS limit to the gateways. This configuration
will be stored in OMAP on the storage cluster so that the volume QoS limits can be reconfigured in the
event of a gateway crash. Since Ceph only supports a single active gateway/path per volume, there is
no need for a distributed QoS model across the gateways within a gateway group.

#### 3.1.1.7 BYOK

Ceph RBDs support [layering](https://docs.ceph.com/en/latest/dev/rbd-layering/). Acadia block
storage volumes utilize this feature for both boot volumes (based on Images) and snapshots. When
BYOK enabled volumes are attached to a VSI, `librbd` on the hypervisor is used as the encryption
engine (for background see [SRB/1793](https://github.ibm.com/cloudlab/srb/tree/master/proposals/1793#3112-mapping-ceph-to-our-use-case)).
In this proposal, `librbd` runs as part of the the NVMeoTCP Gateway which is a part of the Acadia
storage platform.  At volume attachment time, the NVMeoTCP Gateway will call a KMIP server, which
is a shim layer provided by the zonal storage control plane, to fetch the key. The zonal storage
control plane will receive the KMIP request and use the information in the request to fetch the
WDEK for the volume. The wrapped DEK (WDEK) must be unwrapped by the control plane via a keylore
call and sent securely back to the NVMeoTCP Gateway for use in priming `librbd` with the encryption
key for the volume. The encryption key will never be persisted to storage and will be expunged if
the volume is detached. Any failover type scenario that would cause the key to be lost (e.g. OSD node
down, NVMeoTCP Gateway upgrade, etc) will result in the Ceph manager again calling the KMIP shim to
fetch the key for the volume. This design extends the BYOK support for Denali shares in
[SRB/4734](https://github.ibm.com/cloudlab/srb/tree/master/proposals/4734#31143-customer-managed-encryption-at-rest-byok).

Ceph RBDs use `AES-256-XTS-plain64` for at rest encryption. This will be consistent whether the volume/RBD
is attached to a VSI or Bare Metal server.

<!-- markdownlint-disable MD033 -->
<img src=./materials/acadia-bare-metal-BYOK-with-KMIP.png alt="BYOK with KMIP" width="800">

In summary, the flow will be:

1. Customer imports an encryption key, which is stored in the KMS by the regional control plane and
   assigned a CRN
2. Customer creates a volume and specifies the encryption key CRN as the key for the new volume
   a. The zonal control plane creates a master encryption key and wraps that key with the customer key,
      then creates a volume/RBD in Ceph and crypto formats the volume using the master key. The master
      key stored in the LUKS header and protected by a passphrase.
3. Customer requests the newly created volume be attached to a Bare Metal server
   a. The zonal control plane creates `nvme-tcp` mappings for the volume at the Ceph nvmeof gateway.
   b. The zonal and regional control planes map the underlay/loopback IP of the nvmeof gateway to
      reserved IP in the customer's BYOIP space.
4. nvmeof gateways will fetch the encryption key using the KMIP-shim
   a. Control plane will provide IP addresses for the KMIP server, which the nvmeof gateways will use
      to contact the KMIP server.
   b. nvmeof gateways will send a KMIP `Get` operation with an encryption key identifier (the volume
      ID) provided by the control plane.
5. The KMIP shim calls keylore to unwrap the WDEK for the share and returns the DEK material to the
   nvmeof gateway
6. The customer connects the volume to a Bare Metal server
   a. Customer will issue `nvme connect` commands (or equivalent commands for their OS) using the
      BYOIP provided in the volume attachment.

#### 3.1.1.7.1 BYOK Lifecycle State Changes for Bare Metal Support

Most of the BYOK lifecycle state changes will continue to function as described in
[BYOK Phase 3 – Support CRK Lifecycle States](https://github.ibm.com/cloudlab/srb/blob/master/proposals/1085).

In addition to those, the following enhancements are required to support Bare Metal servers:

- Upon encryption key deletion, Keyreact does not perform any action on the bare metal server.

- Keyreact continues to issue suspend and activate requests solely to update and reconcile the volume
  status in the control plane.

- When an encryption key is deleted, the control plane explicitly disables NVMe access by removing
  namespace visibility for the host and transitions the attachment state to `suspended`

- When a key is restored, the control plane re-enables NVMe access by adding namespace visibility for
  the host and transitions the attachment state to `available`.

#### 3.1.1.8 Volume Expansion

The majority of the existing volume expansion workflow will work as-is for NVMeoTCP attached
block storage volumes. Currently, the `sdp` volume expansion workflow for VSIs is:

1. User requests a volume expansion through RIAS API, CLI, etc.
2. The request is processed by the RSOS service controller and updates regional volume objects in Kubernetes
3. Projection controllers in the zonal control plane mirror updates to respective zonal volume objects
4. Zonal Acadia control plane makes a call to the Acadia platform to resize the RBD
5. Zonal Acadia control plane creates a `compute action` to alert the VSI of the volume size change
6. compute action processing performs a libvirt blockresize

For Bare Metal, there is no need for a `compute action`. Instead, SPDK
supports sending an NVMe Asynchronous Event Notification (AEN) when a `bdev` is resized which will
notify any initiators of the size change. The volume expansion workflow when Bare Metal initiator
is attached will be:

1. Zonal Acadia control plane makes a call to the Acadia platform to resize the RBD
2. Zonal Acadia control plane makes a gRPC call to the gateway to notify SPDK about the resize
3. the gateway (SPDK) sends a resize AEN to any attached NVMeoTCP initiators

#### 3.1.1.9 Updates to Global Catalog

`supported_volume_attachment_protocols` will be added to all Instance profiles. The only element
in the list will be `virtio-blk`:

Example:

```json
      "id": "bx3d-128x640",
      ...
      "kind": "instance.profile",
      "metadata": {
        "other": {
          "profile": {
            ...
            "default_config": {
              "allowed_profile_classes": [
                "bx3d"
              ],
            ...
            }
            "supported_volume_attachment_protocols": [ "virtio-blk" ]
          }
        ...
```

`supported_volume_attachment_protocols` will be added to all Bare Metal server profiles. For the
majority of profiles the only element in the list will be `nvme`. The only exceptions to this will
be the [SAP HANA](https://github.ibm.com/cloudlab/srb/blob/master/proposals/2872/README.md)
optimized profiles which will not have an Elba DPU. For those profiles, the
`supported_volume_attachment_protocols` list will be empty.

Example of a Bare Metal server profile that supports NVMe volume attachment:

```json
      "id": "bx2d-metal-192x768",
      ...
      "kind": "bare-metal-server.profile",
      ...
      "metadata": {
        "other": {
          "profile": {
            "billing_plan": "4e2e0cfc-07fd-46cd-8674-64495feec4fd",
            "default_config": {
              "bandwidth": 100000,
              "cpu": 96,
              "cpu_architecture": "amd64",
              "cpu_sockets": 4,
              ...
            }
            "tpm_supported_modes": [
              "tpm_2"
            ],
            "supported_volume_attachment_protocols": [ "nvme-tcp" ]
          }
        ...
```

`supported_volume_attachment_protocols` will be added to all Volume profiles. The `sdp` profile
will report `nvme-tcp` and `virtio-blk`. For all other profiles, the only element in the list will
be `virtio-blk`.

Example of the `sdp` volume.profile:

```json
      "kind": "volume.profile",
      "metadata": {
        "other": {
          "profile": {
            "adjustable_iops_supported": true,
            "block_size": 16,
            "boot_capacity_config": {
              "max": 32000,
              "min": 1,
              "units": "gb"
            },
            ...
            "supported_volume_attachment_protocols": [ "nvme-tcp", "virtio-blk" ],
            ...
      },
      "name": "sdp",
```

#### 3.1.1.10 Snapshot and Backup Policy

Snapshot and backup policy functionality for bare metal servers is included in this SRB.

**Snapshot Consistency Groups:**

- Snapshot consistency groups now support bare metal servers
- Service tags associated with snapshot consistency groups are prefixed with `is.instance:` or `is.bare-metal-server:`
- Snapshots can be created for volumes attached to bare metal servers

**Backup Policies:**

- Backup policies now support `bare_metal_server` as a match resource type
- For bare metal servers, the `included_content` field supports:
  - `data_volumes`: Include the bare metal server's data volumes

#### 3.1.1.10.1 Snapshot support

The process to take Acadia based snapshots is not dependent how the volume is attached, or whether
the volume is attached. Taking snapshots of volumes which are attached to Bare Metal servers does
not require any additional changes.

For multi-volume snapshots, the same is true. If a volume is multi-attached, each Bare Metal server
could create its own snapshot consistency group. This would result in multiple snapshots being created
for the multi-attached volumes (one for each Bare Metal server which has taken a multi-volume snapshot).

The snapshot restoration process (aka create volume from snapshot) will largely work as-is. The
one exception to that is instant availability which is covered in more detail in the following section.

#### 3.1.1.10.2 Snapshot limitations for Beta and LA

Snapshot restore instant availability for Acadia backed volumes (i.e. volume profile `sdp`) is being
added under [SRB/5324](https://github.ibm.com/cloudlab/srb/issues/5324) for volumes attached to VSIs.
To support instant availability on volumes created from snapshots which are attached to Bare Metal
servers requires additional feature development on the Acadia storage platform, Ceph. The Ceph team
has committed to delivering the additional capabilities required, however their timeline of 4Q26
doesn't align with the Beta/LA timeline for this proposal. For that reason, volumes created from
snapshots will not support instant availability when attached to Bare Metal servers. This means
that the volume hydration must complete before the volume can be attached to a Bare Metal server. While
the volume is in the process of being hydrated, the  `health_state` is `degraded`. Any volume
attachments will remain in the `attaching` status until the hydration completes.

#### 3.1.1.11 Multi-Attach

Allowing a volume to be attached to more than one compute instance concurrently has not been supported
in VPC to date. This is due to the lack of low-level volume locking primitives in the `virtio-blk`
attachment protocol used for VSI volume attachments.

The `nvme-tcp` protocol supports NVMe Reservations
which allows applications to coordinate sharing a block volume (see sections 7.5-7.8 of the
[NVMe specification](https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-Revision-2.3-2025.08.01-Ratified.pdf)).

For customers that require multi-attachment, a new volume property `attachment_mode` will be added to
enable multi-attach. By default, the `attachment_mode` will be `single`. This is done out of an abundance
of caution to prevent multiple Bare Metal servers from sharing a volume unintentionally, which could
lead to data loss. At volume creation, or during a future `PATCH` while the volume is unattached,
the volume's `attachment_mode` can be changed to `multiple`. With an `attachment_mode` of `multiple`,
the control plane will allow up to 32 attachments to a single volume. All attachments *must* use the
`nvme-tcp` attachment protocol; any attachments from a VSI using the attachment protocol `virtio-blk`
will not be allowed.
Each Bare Metal server will create unique volume attachment resources for each shared volume (via
POST `/bare_metal_servers/{id}/volume_attachments`). Each shared volume will report all volume attachments,
including the Bare Metal servers, in the `volume_attachments` property.

When creating a volume attachment to a volume with an `attachment_mode` of `multiple`, the auto-delete
option (i.e. `delete_volume_on_instance_delete` or `delete_volume_on_server_delete`) *must* be false.

Only data volumes will support and `attachment_mode` of `multiple`. Boot volumes cannot be attached
to multiple Bare Metal servers concurrently.

For multi-attached volumes, there will be no QoS fairness or enforcement of IOPS or throughput amongst
the attached Bare Metal servers. In other words, there is no guarantee of:

* an even split of IOPS/throughput between attached Bare Metal servers
* a min IOPS/throughput for each attached Bare Metal server

The split of volume performance between Bare Metal servers is left up to the application/workload.

### 3.1.2 Customer Experience

This SRB introduces customer facing changes, so there will be a change/impact to the customer
experience.  The volume attachment experience for Bare Metal servers will deviate from the VSI
volume attachment experience.  Access to block storage volumes via `nvme-tcp` from Bare Metal
servers will require customer action.  By 'attaching', a volume to a Bare Metal server the cloud
orchestration will make that volume available on the customer's Bare Metal server via link local
addresses.  Then the customer must
issue an `nvme connect` command to complete the attachment, which is more akin to mounting a file
share ... or connecting to an iSCSI volume in the IBM classic infrastructure. UI pages will need
to be updated to provide customers with the details required to connect to `nvme-tcp` volumes, namely

* The Bare Metal server's Host Nvme Qualified Name (NQN)
* `nvme-tcp` endpoints (IP, Subsystem NQN)
  * There will be 4 per volume
* CHAP key (aka storage access secret)
  * used to securely connect to the volume and establish encryption-in-transit

#### 3.1.2.1 User Interface (UI)

Primarily, three experiences in the UI will need to be updated:

* [Provision a BM Server](https://cloud.ibm.com/vpc-ext/provision/bm)
* [View/Edit a BM Server](https://cloud.ibm.com/vpc-ext/compute/bm/{id}/overview)
* [View/Edit a Block storage volume](https://cloud.ibm.com/vpc-ext/storage/storageVolumes/{id}/overview)

Snapshot and backup policy functionality is now supported for bare metal servers. UI flows should support:

- Creating snapshots of volumes attached to bare metal servers
- Creating backup policies with `bare_metal_server` as the match resource type
- Viewing snapshot consistency groups with `is.bare-metal-server:` service tag prefix

This SRB will also cause a paradigm shift that may cause a ripple effect through various
pages/tables:

* Not all block storage volumes can be attached to all compute instance types. More specifically,
*only* volumes of the `sdp` profile can be attached to Bare Metal server instances.

#### 3.1.2.1.1 Provision a BM Server in the UI

Similar to the VSI create page like this:

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-vsi.png alt="VSI UI" width="700">

The Bare Metal server creation will need to have a storage subsection where volumes can be created
or selected:

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-bm.png alt="BM UI" width="700">

In addition, only volumes which match the Bare Metal server's profile volume attachment protocol
should be shown (i.e. volumes of a profile that support `nvme-tcp` attachment).

#### 3.1.2.1.2 View/Edit a BM Server in the UI

Similar to VSI view/edit pages like this:

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-vsi-overview.png alt="VSI OV UI" width="700">

The Bare Metal server view/edit pages will need to have a storage subsection where
volumes can be attached or selected:

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-bm-overview.png alt="BM OV UI" width="700">

#### 3.1.2.1.3 View/Edit a Block storage volume in the UI

When viewing details for a Block storage volume, the attached VSI is shown. Will need to
additionally show the attached Bare Metal server.

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-vol-overview.png alt="VOL UI" width="700">

In the block storage volume table, may want a column for the attached server type (VSI vs BM). Note,
this is fine within the scope of *this* SRB but in the future support may be added to allow a VSI
and Bare Metal server to concurrently attach a volume.

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-vol-columns.png alt="VOL UI COL" width="700">

From the block storage volume table, the customer may choose the `Attach to server` option. In that
scenario, will need to add an option for Bare Metal server attachment in addition to VSI.

<!-- markdownlint-disable MD033 -->
<img src=./materials/ui-vol-attach.png alt="VOL UI ATTACH" width="700">

#### 3.1.2.2 Command Line Interface (CLI)

CLI commands, parameters, and response payload changes are needed to support RIAS volume attachments
for Bare Metal servers.

#### 3.1.2.2.1 New CLI commands

|Command|Short name|Description|
|---|---|---|
|bare-metal-server-volume-attachments|bm-vols|List all volume attachments to a bare metal server|
|bare-metal-server-volume-attachment|bm-vol|View details of a bare metal server volume attachment|
|bare-metal-server-volume-attachment-add|bm-vola|Create a volume attachment, connecting a volume to a bare metal server|
|bare-metal-server-volume-attachment-detach|bm-vold|Delete one or more volume attachments, detaching volume from a bare metal server|
|bare-metal-server-volume-attachment-update|bm-volu|Update a bare metal server volume attachment|
|bare-metal-server-storage-access-secret-rotate|bm-sasr|Rotate the storage access secret for NVMe|

#### 3.1.2.2.2 New parameters

| Command | Flag | Value | Description |
| --- | --- | --- | --- |
| bare-metal-server-create | --volume-attach | VOLUME_ATTACH_JSON <br> @VOLUME_ATTACH_JSON_FILE | volume attachment in JSON or JSON file, list of volumes. For the data schema, see the 'volume_attachments' property in the [API documentation](/apidocs/vpc#create-instance). |
| bare-metal-server-network-attachment-create | --interface-type | `nvme_tcp` | New interface type for network attachment.<br>• `nvme_tcp` - NVMe over TCP protocol for remote block storage access |

#### 3.1.2.2.3 Response payload changes

To support volume attachment protocols instance profiles, Bare Metal server profiles, and volume
profiles need to be extended. This will impact:

* `bare-metal-server-profiles` - New field `supported_volume_attachment_protocols`
  (array of supported protocols: `nvme_tcp`)
* `instance-profiles` - New field `supported_volume_attachment_protocols`
  (array of supported protocols: `virtio_blk`)
* `volume-profiles` - New field `allowed_access_protocols` (array of supported protocols:
  `virtio_blk` and/or `nvme_tcp`)
* `volumes` - New field `allowed_access_protocols` and `attachment_mode` (single/multiple)
* `backup-policies` - Extended to support `match_resource_type: bare_metal_server` with
  `included_content: [data_volumes]`
* `snapshot-consistency-groups` - Service tags now include `is.bare-metal-server:` prefix

Additionally, Bare Metal server response payloads will be extended:

* New field `volume_attachments` to include info about all volume attachments to the Bare Metal server
* New field `nvme_qualified_name` - The NVMe qualified name for the bare metal server
* New field `storage_access_secret` - The encrypted storage access credential with encryption key reference
  and rotation timestamp
* Volume attachment responses include `device.id`, `volume`, `bandwidth`, `status`, `type`, `subsystem_nqn`,
  and `attachment_protocol`
* Volume attachment `type` field for bare metal servers only supports `data` (boot volumes are not supported
  for bare metal servers)

#### 3.1.2.3 Terraform and Packer

The RIAS API changes will necessitate changes in Terraform. The following table has a summary
of changes required:

| New? | Resource | Change Summary |
| --- | --- | --- |
| Existing | [ibm_is_bare_metal_server](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_bare_metal_server) | Supports `volume_attachments`, `nvme_qualified_name` and `storage_access_secret` fields |
| New | ibm_is_bare_metal_server_volume_attachment | Create/Show/Update/Delete Bare Metal server volume attachment with properties: `name`, `volume`, `delete_volume_on_server_delete`, `network_attachment`, `attachment_protocol` |
| New | ibm_is_bare_metal_server_volume_attachments | List Bare Metal server volume attachments |
| New | ibm_is_bare_metal_server_storage_access_secret_rotate | Rotate storage access secret for NVMe volumes |
| Existing | [ibm_is_volume](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_volume) | Show `allowed_access_protocols` array and `attachment_mode` (single/multiple). Include Bare Metal server volume attachments |
| Existing | [ibm_is_volume_profile](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_volume_profile) | Show `allowed_access_protocols` array and `attachment_mode` array |
| Existing | [ibm_is_volume_profiles](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_volume_profiles) | Show `allowed_access_protocols` array and `attachment_mode` array |
| Existing | [ibm_is_bare_metal_server_profile](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_bare_metal_server_profile) | Show `supported_volume_attachment_protocols` array (nvme_tcp) |
| Existing | [ibm_is_instance_profile](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/is_instance_profile) | Show `supported_volume_attachment_protocols` array (virtio_blk) |
| Existing | [ibm_is_backup_policy](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_backup_policy) | Extended to support `match_resource_type: bare_metal_server` with `included_content: [data_volumes]` |

#### 3.1.2.4 Customer-facing Documentation

The Volumes API and documentation will need to be updated to be agnostic to VSI
and bare metal server volume attachments. Currently, the documentation is exclusively written for
VSIs.

The bare metal server volume attachments API changes are influenced by VSI volume attachments, so
documentation for bare metal server volume attachments should be based on its VSI sibling. Some
notable differences, however, include only select volume profiles and volumes will be supported
and adapting instructions for bare metal server usage.

To help customers distinguish what volume profiles and volumes can be attached to a bare metal server,
the concept of "attachment protocols" is needed. Rather than use general "instance" or "bare metal
server" terms for supported attachments, the underlying technology terminology should be surfaced
to the customer. Other cloud providers like
[Amazon](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/nvme-ebs-volumes.html)
and [Microsoft](https://learn.microsoft.com/en-us/azure/virtual-machines/enable-nvme-interface) do
this and if NVMe-attached storage is supported for VSIs in the future there will be less churn.

The changes required to the public VPC documentation to show how to attach and manage external
block storage on Bare Metal servers are summarized in the following table:

| New? | Page Location | Change Summary |
| --- | --- | --- |
| Existing | [Get Started / Storage / VPC Storage service overview](https://cloud.ibm.com/docs/vpc?topic=vpc-storage-overview&interface=ui) | List Bare Metal servers as an attachment option |
| Existing | [Get Started / Storage / Storage overview for bare metal servers / Storage overview for Bare Metal Servers for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-bare-metal-servers-storage&interface=ui) | List Block Storage as an option |
| Existing | [Get Started / Storage / Block storage / About Block Storage for VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-block-storage-about&interface=ui) | List Bare Metal servers as an attachment option |
| Existing | [How To / Bare metal servers / Planning for Bare Metal Servers on VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-planning-for-bare-metal-servers) | Add Auxillary storage to the checklist |
| Existing | [How To / Bare metal servers / Creating Bare Metal Servers on VPC](https://cloud.ibm.com/docs/vpc?topic=vpc-creating-bare-metal-servers&interface=ui) | Show options to create/attach block storage. |
| Existing | [How To / Block storage volumes / Attaching a Block Storage for VPC volume](https://cloud.ibm.com/docs/vpc?topic=vpc-attaching-block-storage&interface=ui) | Show options to attach to Bare Metal servers |
| Existing | [How To / Block storage volumes / Managing Block Storage for VPC volumes](https://cloud.ibm.com/docs/vpc?topic=vpc-managing-block-storage&interface=ui) | Show options for Bare Metal servers |
| New | How To / Block storage volumes / Set up Block Storage for VPC on a Bare Metal server | NVMe driver maturity varies. List any special instructions that may be required for each OS to see storage and adapt to any changes (add volume/remove volume/resize volume) |
| New | How To / Block storage volumes / Mapping NVMe devices to Block Storage volumes | Instructions on how to show the globally unique IDs for an NVMe device and how to associate that to a Block Storage volume. Similar to [this](https://www.ibm.com/blog/using-the-metadata-service-to-identify-disks-in-your-vsi-with-ibm-cloud-vpc/) |
| New | How To / Data encryption / Creating a Bare Metal server with customer managed encryption volumes | Show how to create BYOK block storage volumes and attach to Bare Metal servers. Similar to [Creating instances with BYOK](https://cloud.ibm.com/docs/vpc?topic=vpc-creating-instances-byok&interface=ui) |
| Existing | [How To / Logging and Monitoring / Activity Tracker events](https://cloud.ibm.com/docs/vpc?topic=vpc-at-events&interface=ui) | Update with new AT events for bare metal server volume attachments |

#### 3.1.2.4.1 Identifying NVMe devices

Clients need a method to associate an NVMe device on their Bare Metal server (e.g. `/dev/nvme0n1`)
to a RIAS Volume in their IBM Cloud account. The NVMe namespace UUID and GUID will be set to the
[UUID](https://github.ibm.com/riaas/regional-storage/blob/da188f89fdfeaa5c3f0a70cc5bee144267242f02/pkg/clients/control-plane/common/utils.go#L615-L618)
of the RIAS volume (with no region prefix). This id will be present in volume attachment API responses
(via the `device.id` property) and also visible on the host itself once the volume has been connected.

For example, this volume:

```bash
Listing volumes in resource group Default and region us-south under account IBM as user theoharr@us.ibm.com...
ID                                          Name                     Status      Capacity   IOPS   Profile           Attachment state   Attachment type   Zone         Resource group   Catalog Offering Version   Catalog Offering Plan   Storage Generation
r134-ea248fcd-fb77-4919-bff2-7c4db5f0428c   ve-dev-sdp0-datavol      available   100        3000   sdp               attached           data              us-south-2   Default          -                          -                       2
```

Shows up on the OS as `/dev/nvme2n2`, and the NVMe UUID NVMe GUID can be interrogated via the `nvme`
cli or the `/sys` filesystem:

```bash
root@r11s24nvmeinit:~# nvme id-ns /dev/nvme2n2
NVME Identify Namespace 8:
nsze    : 0x32000
...
nguid   : ea248fcdfb774919bff27c4db5f0428c
eui64   : 0000000000000000
lbaf  0 : ms:0   lbads:9  rp:0 (in use)

root@r11s24nvmeinit:~# cat /sys/class/block/nvme2n2/uuid
ea248fcd-fb77-4919-bff2-7c4db5f0428c

root@r11s24nvmeinit:~# cat /sys/class/block/nvme2n2/nguid
ea248fcd-fb77-4919-bff2-7c4db5f0428c
```

Note, the volume's UUID was chosen for the NVMe UUID/GUID rather than the Bare Metal server's volume
attachment's UUID so that the volume presents a consistent UUID across all attachments. This is
important so in multi-attach configurations, clustered applications can detect
and treat shared block volumes appropriately.

#### 3.1.2.5 Customer-facing Metrics

No new customer-facing metrics.

### 3.1.3 Customer-facing APIs

Proposed changes
can also be viewed in the [api-spec PR](https://github.ibm.com/riaas/api-spec/pull/7344). Note, the
proposed API changes do include a hazardous change. See section [3.1.3.5](#3135-hazardous-api-changes)
for details.

API changes will be implemented following the [Dynamic API Maturity](https://github.ibm.com/cloudlab/srb/tree/master/proposals/2313)
pattern. A feature flag, `is-bare-metal-server-volumes-phase1-allowlist` and maturity flag `is-bare-metal-server-volumes-phase1-maturity`,
will gate access to the new API changes.

Summary of changes:

| Endpoint | Action(s) | Change Summary |
| --- | --- | --- |
| `/bare_metal_server/profiles/{name}` | GET | New field `supported_volume_attachment_protocols` will return an array of supported attachment protocols (only `nvme_tcp`) |
| `/instance/profiles/{name}` | GET | New field `supported_volume_attachment_protocols` will return an array of supported attachment protocols (only `virtio_blk`) |
| `/instance/{instance_id}/volume_attachments/{id}` | GET | New field `attachment_protocol` will return the attachment protocol used (`virtio_blk`) |
| `/volume/profiles/{name}` | GET | <ul><li>New field `supported_attachment_protocols` will return an array of supported attachment protocols (`virtio_blk` and/or `nvme_tcp` depending on the profile)</li><li>New field `attachment_mode` will return an array of supported attachment modes (`single` and/or `multiple`)</li></ul> |
| `/volumes` | GET, POST | <ul><li>New field `attachment_mode` (single/multiple) for multi-attach support</li></ul> |
| `/volumes/{id}` | GET, PATCH | <ul><li>New field `attachment_mode` (single/multiple) for multi-attach support</li><li>Additions to the `volume_attachments` list to support Bare Metal server references</li></ul> |
| `/bare_metal_servers` | POST | <ul><li>Supports `volume_attachments` to allow for data volume attachments</li></ul> |
| `/bare_metal_servers/{bmid}` | GET | <ul><li>New field `volume_attachments` to report all volume attachments</li><li>New field `nvme_qualified_name` - NVMe qualified name for the server</li><li>New field `storage_access_secret` - encrypted storage access credential with encryption key and rotation timestamp</li></ul> |
| `/bare_metal_servers/{bmid}/volume_attachments` | GET, POST | New endpoints to list or create volume attachments. Includes fields: `name`, `volume`, `delete_volume_on_server_delete`, `device`, `bandwidth`, `status`, `type` (only `data` - boot volumes not supported), `attachment_protocol`, `network_attachment`, `subsystem_nqn`, `ips` |
| `/bare_metal_servers/{bmid}/volume_attachments/{id}` | GET, PATCH, DELETE | New endpoints to retrieve, update, or delete volume attachments. PATCH supports updating `name` and `delete_volume_on_server_delete` |
| `/bare_metal_servers/{bmid}/storage_access` | GET | New endpoint to retrieve storage access secret. |
| `/bare_metal_servers/{bmid}/storage_access/rotate` | POST | New endpoint to rotate storage access secret. |
| `/backup_policies` | GET, POST | Extended to support `match_resource_type: bare_metal_server` with `included_content: [data_volumes]`. Creates backups of all volumes with `storage_generation` value of `1` attached to the bare metal server |
| `/backup_policies/{id}` | GET, PATCH | Extended to support bare metal server resource type with `included_content` field for data volumes |

#### 3.1.3.1 Activity Tracker (AT)

New actions will be added to mirror the volume attachment events currently supported for VSIs.

| Resource | Action | Description |
| --- | --- | --- |
| Bare Metal Server | `is.bare-metal-server.volume-attachment.create` | Bare Metal Server volume attachment was created |
| Bare Metal Server | `is.bare-metal-server.volume-attachment.delete` | Bare Metal Server volume attachment was deleted |
| Bare Metal Server | `is.bare-metal-server.volume-attachment.update` | Bare Metal Server volume attachment was updated |
| Bare Metal Server | `is.bare-metal-server.volume-attachment.read` | One or more Bare Metal Server volume attachments was retrieved |
| Bare Metal Server | `is.bare-metal-server.storage-access.read` | Bare Metal Server storage access secret was retrieved |
| Bare Metal Server | `is.bare-metal-server.storage-access.update` | Bare Metal Server storage access secret was rotated |

#### 3.1.3.2 Global Search and Tagging (GhoST)

No changes

#### 3.1.3.3 IBM Cloud Compliance and Security Center (IBM SCC)

No changes. Existing properties for Volumes and Bare Metal servers still apply.

Note: Storage access secrets are encrypted and managed securely through the encryption_key reference.

#### 3.1.3.4 Customer Logging

No additional customer facing logs will be generated. Currently, regional storage does not emit any
customer facing logs for RIAS volumes. The regional baremetal service will continue to generate
quota related logs.

#### 3.1.3.5 Hazardous API Changes

This feature will introduce hazardous changes to the API. The API [hazardous change remediation](
https://test.cloud.ibm.com/docs/api-docs?topic=api-docs-making-api-changes#change-remediation)
process will be followed, covering the items listed here.

In the existing `GET /volumes/{id}` response, a list of `volume_attachments` is returned which represents
each binding of the volume to a compute server. The payload for each volume attachment includes:

- `delete_volume_on_instance_delete`: Boolean to indicate whether deleting the instance will also
  delete the attached volume
- `deleted`: Standard deleted property
- `device`: The device name of the volume attachment.
- `href`: The URL for the volume attachment.
- `id`: ID of the volume attachment
- `instance`: Reference to the Instance which has attached the volume
- `name`: Name of the volume attachment
- `resource_type`: `volume_attachment`
- `type`: `boot` or `data`

Adding the ability for Bare Metal servers to attach to volumes, requires adding new properties to
the `volume_attachments` schema in the form of a `oneOf`.  However, this introduces a hazardous change.
Many properties within the `volume_attachments` schema are common, however two required properties
are not:

- `delete_volume_on_instance_delete` and
- `instance`

When a bare metal server is attached to a volume, the `instance` and `delete_volume_on_instance_delete`
properties will be omitted from the response; instead `server` and `delete_volume_on_server_delete`
properties will be returned.

Existing clients which do not use the Block for Bare Metal feature will
be unaffected.

However, clients that do use the Block for Bare Metal feature must update their client
API versions. Otherwise, older clients will fail to deserialize the GET `/volumes/` responses when
the list of `volume_attachments` contains a reference to a Bare Metal server.

### 3.1.4 Internal APIs and Schemas

The CreateVolumeSpecDryRun API, which validates volume creation requests before provisioning resources,
will accept a new volume_attachment_protocol parameter and optionally perform plan ID and private
catalog visibility validation. The getVolume API will return volume_attachment_protocols in its
response to enable validation of attachment protocol compatibility for existing volumes before
attachment operations.

```golang
enum VolumeAttachmentProtocols {
  virtio_blk = 0;
  nvme_tcp  = 1;
}
message ValidateCreateVolumeRequest {
  // Enable validation of plan ID and private catalog visibility.
  // When true, performs additional checks beyond basic volume creation validation.
  // Default: false
  bool enable_additional_validation
}
message ValidateVolumeInfoRequest {
  VolumeAttachmentProtocols volume_attachment_protocol
}

message GetVolumeResponse{
 repeated VolumeAttachmentProtocols volume_attachment_protocols
}
```

### 3.1.5 Deployment Architecture

The Acadia storage cluster is deployed and managed independently from the VPC production
environment. The new NVMe related services will be optionally deployed based on configuration
options. See section [2.6.2](#261-acadia-platform-updates).

The regional and zonal control plane microservices will be deployed and updated via Razee. A
globals/ConfigMap property (`enable_bm_controllers`) exists to control which environments the
microservices get deployed into.

The following feature flags are currently defined and used:

| Description | FF Name | /internal/v1/features name |
| --- | --- | --- |
| Account enablement | `is-bare-metal-server-volumes-phase1-maturity` and `is-bare-metal-server-volumes-phase1-allowlist` | is-bare-metal-server-volumes-phase1 |
| Per BM server volume attachment limit | `is-bare-metal-server-volume-attach-limit-value` | n/a |

## 3.2 Architectural Considerations and Impact

### 3.2.1 Security

The architecture and design in this proposal strives to achieve the same level of security
as with the existing Acadia based block storage capability with VSIs. Some specific examples are:

* All control plane traffic uses mTLS (including zonal control plane and NVMeoTCP gateways)
* Data plane traffic is encrypted in flight (NVMeoTCP with TLS, RBD/RADOS with msgr2 secure mode)
* Block storage volumes are only exposed to Bare Metal servers which should have access

The following table provides a lightweight threat model addressing key security concerns:

| Threat | Mitigation | Implementation Details |
|--------|-----------|------------------------|
| **Cross-tenant storage access** | `StorageVNIGateway` isolation + CHAP allowlisting | • Each `StorageVNIGateway` is 1:1 bound to a single BM server<br>• Host NQN allowlisting at gateway level ensures only authorized servers can connect<br>• SDN enforces source validation on all storage LLAs<br>• CHAP authentication required for all connections |
| **Man-in-the-Middle (MITM) attacks** | TLS 1.3 encryption | • All nvme-tcp connections use TLS 1.3 with strong cipher suites (AES-128-GCM or AES-256-GCM)<br>• Gateway-to-Ceph traffic uses msgr2 secure mode (AES-128-GCM)<br>• Control plane uses mTLS with certificate-based identity |
| **Compromised Bare Metal OS** | No lateral storage access | • Compromised BM server cannot access volumes attached to other servers due to:<br>&nbsp;&nbsp;- VNI isolation (aggregate VNI is server-specific)<br>&nbsp;&nbsp;- Host NQN allowlisting (only authorized NQN can connect)<br>&nbsp;&nbsp;- CHAP authentication (unique per server)<br>• Discovery service only returns subsystems authorized for that specific server |
| **CHAP secret compromise** | Customer-initiated rotation + rate limiting | • Customers can rotate CHAP secrets on-demand via API<br>• Rate limited to 1 rotation per hour per server<br>• Secrets encrypted with customer's SSH public key for delivery<br>• Stored as wrapped DEK in control plane |
| **Unauthorized volume attachment** | IAM enforcement | • All volume attachment operations require appropriate IAM permissions<br>• Volume `operate` permission required to attach existing volumes<br>• Volume `create` permission required to create new volumes during attachment<br>• See section 3.2.1.5 for complete IAM action matrix |
| **Data at rest exposure** | Encryption at rest | • All volumes encrypted at rest (provider-managed or BYOK)<br>• BYOK keys managed through Key Protect/HPCS<br>• Inherited from existing Acadia block storage security model |
| **Discovery service reconnaissance** | Authorization-based results + rate limiting | • Discovery only returns already-authorized subsystems<br>• Cannot be used to discover other tenants' storage<br>• SDN rate limiting prevents abuse<br>• Single request processing per server prevents DoS |
| **Replay attacks** | TLS 1.3 + CHAP nonce | • TLS 1.3 provides replay protection<br>• CHAP protocol includes challenge-response with nonces<br>• Each connection establishes fresh session keys |

The following table clarifies the security responsibilities between IBM and the customer:

| Security Aspect | IBM Guarantees | Customer Responsibilities |
|----------------|----------------|---------------------------|
| **Network Isolation** | • VNI-based tenant isolation<br>• Link-local address scoping per BM server<br>• SDN source validation enforcement<br>• No cross-tenant storage access | • Monitoring network access patterns |
| **Authentication & Authorization** | • Unique Host NQN generation per BM server<br>• Cryptographically secure CHAP key generation (256-bit)<br>• Host NQN allowlisting at gateway level<br>• TLS 1.3 enforcement for all nvme-tcp connections | • Secure storage of CHAP secrets in guest OS<br>• Proper SSH key management for secret encryption<br>• Initiating secret rotation when needed<br>• Configuring nvme connect with correct credentials |
| **Encryption** | • TLS 1.3 for nvme-tcp data plane (AES-128-GCM or AES-256-GCM)<br>• msgr2 secure mode for gateway-to-Ceph (AES-128-GCM)<br>• Encryption at rest (provider-managed or BYOK)<br>• mTLS for all control plane communications | • BYOK key management (if using customer-managed encryption)<br>• Proper handling of decrypted data in guest OS |
| **Access Control** | • IAM integration for API operations<br>• Volume attachment authorization enforcement<br>• Discovery service returns only authorized subsystems | • Proper IAM policy configuration<br>• Managing user access to BM servers<br>• Guest OS access controls |
| **Secret Management** | • Encrypted CHAP key delivery via SSH public key<br>• Secure storage of CHAP keys in control plane (wrapped DEK in CRD)<br>• Rate limiting on secret rotation (1/hour)<br>• Atomic gateway configuration updates | • Protecting CHAP secrets in guest OS (e.g., file permissions on nvme config)<br>• Secure handling during nvme connect operations<br>• Not logging or exposing CHAP secrets |
| **Manual Operations** | • Providing clear documentation and examples<br>• Returning all necessary connection parameters via API<br>• Discovery service for simplified connection | • Executing nvme connect commands in guest OS<br>• Configuring nvme-cli and autoconnect services<br>• Monitoring connection health<br> |

#### 3.2.1.1 FS Cloud Considerations

All communications are encrypted. All control plane traffic uses mTLS with certificates stored in
vault to provide identity. Data plane traffic is either encrypted with:

* TLS 1.3 - TLS_AES_128_GCM_SHA256 or TLS_AES_256_GCM_SHA384
* msgr2 - AES-128-GCM

To establish secure channels for the NVMeoTCP traffic, DH-HMAC-CHAP keys and the CHAP protocol will
be leveraged to to establish secure channels (i.e. TLS) between the Bare Metal servers and NVMeoTCP
gateways.  Each Bare Metal server will have a unique DH-HMAC-CHAP key, which the customer can
rotate on demand.

Customer data stored on the block storage volumes will always be encrypted at rest, either through
provider managed encryption or BYOK (this is inherited from the Acadia block storage with VSIs
project).

#### 3.2.1.2 Secrets

> If a new secret (key, credential, pass-code) needs to be stored, has a [secrets
> template](https://github.ibm.com/gensec/platform-inventory/blob/main/secops/docs/template.md)
> been created and added to the [vault secret
> documentation](https://github.ibm.com/gensec/platform-inventory/blob/main/secops/README.md)?
> Further, how often will this new secret be rotated?
>
> In the case of a secret for encrypting data at rest:
>
> * Who can access that key and how is that access controlled?
> * Is the key stored in a hardware security module (HSM)?
> * How often is the key rotated?
TODO

#### 3.2.1.3 VPC Control Plane Hardening

> If this project includes new control plane microservices, the [Control Plane Hardening
> Checklist](architecture/security/hardening/README.md) must be completed for each microservice and
> linked to from this proposal, along with rationale for any requested exceptions.
>
> List any changes to HostOS IP packet filtering policies, and any changes to listeners or sources
> for existing HostOS networking traffic.
>
> Also consider the implications of the section "Secure Development Practices" of the
> [Service Framework Design Requirements](https://pages.github.ibm.com/ibmcloud/Service-Framework/12_securitydesignrequirements.html)
TODO

#### 3.2.1.4 Significant Change Considerations

| Yes | No | N/A | |
| :--- | :--- | :--- | :--- |
| | x | | Introduce a new major/minor Linux distribution or base image into your service's compute resources (hypervisor, VSI, container) |
| | x | | Change impacted service's trust zones or change how impacted service's trust zones are mapped to network segments |
| | x | | Change Ports/Protocols/Services exposed to public routes |
| x | | | Change the connections between IBM operated components/control planes that cross the boundaries of your service's trust zones |
| x | | | Change the network encryption used to protect data in flight |
| | x | | Change the database/datastore used to persist service data/metadata |
| | x | | Change the encryption used to protect data at rest |
| | x | | Change the classification of data processed by any portion of your compute resources to a higher level |
| | x | | Introduce new cryptographic libraries or modules, including switching to different versions |
| | x | | Change the list of dependencies that your service leverages |
| | x | | Change any hardware your service leverages |
| | x | | Introduce changes to the software packages used by your service components |
| | x | | Introduce new privileged process (a new system process or container running with sudo/root) |
| | x | | Change how the components of your service authenticate and authorize to each other |
| | x | | Change to operators remote access (including the network encryption used) |
| | x | | Change how impacted services collect audit information, OR what audit information is collected |
| x | | | Change identities and/or access |
| | x | | Change the following IBM Cloud public documentation: regional availability, high availability considerations, BCDR customer responsibilities |

#### 3.2.1.5 Identity and Access Management

API updates will follow existing IAM best practices and patterns, including the IAM action checking
done by instances. The goal of these changes is to maintain consistency with the existing VSI volume
attachment implementation. One area where the Bare Metal volume attachment implementation differs is
when dealing with cascaded volume deletions caused by the deletion of the Bare Metal server. The
inconsistency with cascaded volume deletion behavior between VSI and Bare Metal servers, although
not ideal is the preferred solution. The alternatives are:

* Change the VSI volume attachment implementation to match the Bare Metal server implementation
  * This has a high likelihood of breaking existing clients whose accounts lack the sufficient volume
    permissions.
* Model the Bare Metal server implementation to match the VSI volume attachment implementation
  * This would perpetuate a suboptimal model with regards to IAM permission enforcement, by not
    checking for `is.volume.volume.delete` when a volume is being deleted by way of a Bare Metal server
    deletion.

Following in a summary of additional IAM checks required to support block
storage attachment on Bare Metal servers:

| API Endpoint | Additional IAM Actions | Roles | Workflow | Notes |
| --- | --- | --- | --- | --- |
| `POST /v1/bare_metal_servers/` | When no volume attachments are specified <ul><li>no changes</li></ul> When volume attachment is creating a volume <ul><li>`is.volume.volume.create`</li></ul> When volume attachment volume schema has `attachment_mode` specified <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> When volume attachment is creating a volume from a snapshot in the same account <ul><li>`is.snapshot.snapshot.operate`</li></ul> When volume attachment is creating a volume from a snapshot in a different account <ul><li>`is.volume.volume.allow-remote-account-snapshot-restore`</li></ul> When volume attachment uses an existing volume <ul><li>`is.volume.volume.operate`</li> | AE | Create a BM server. Optionally with data volume attachments | Data volumes attachments can optionally be created as part of BM server creation. When a volume attachment creates a volume, must have volume `create` permission and conditional permissions are checked for snapshots and attachment_mode. When a volume attachment uses an existing volume, must have the volume `operate` permission. |
| `DELETE /v1/bare_metal_servers/<bmid>` | <ul><li>no changes</li></ul> When the BM server has auto-delete volumes attached <ul><li>`is.volume.volume.delete`</li></ul> | AE | Delete a BM server | If the BM server has any volumes attached with the `delete_volume_on_server_delete` property set this causes a cascaded delete of the volumes marked as auto-delete. Conditionally enforce the `delete` volume permission. This is a deviation from the existing VSI behavior. |
| `GET /v1/bare_metal_servers/<bmid>/storage_access` | <ul><li>`is.bare-metal-server.bare-metal-server.read`</li></ul> | AEOV | View storage access secret | Enforce the parent resource `read` permission. |
| `POST /v1/bare_metal_servers/<bmid>/storage_access/rotate` | <ul><li>`is.bare-metal-server.bare-metal-server.update`</li><li>`is.bare-metal-server.bare-metal-server.manage-storage-access-secret`</li><li>`is.key.key.operate`</li></ul> | AE | Rotate storage access secret | Rotating the storage access secret is a higher privileged operation, enforce new `manage-storage-access-secret` action in addition to the parent resource `update` and keys `operate` permissions. |
| `GET /v1/bare_metal_servers/<bmid>/volume_attachments` | <ul><li>`is.bare-metal-server.bare-metal-server.read`</li></ul> | AEOV | List volume attachments | Enforce parent resource `read` permission for volume attachments |
| `GET /v1/bare_metal_servers/<bmid>/volume_attachments/<volattid>` | <ul><li>`is.bare-metal-server.bare-metal-server.read`</li></ul> | AEOV | GET volume attachment | Enforce parent resource `read` permission for volume attachment |
| `POST /v1/bare_metal_servers/<bmid>/volume_attachments/` | <ul><li>`is.bare-metal-server.bare-metal-server.update`</li><li>`is.volume.volume.create`</li></ul> When volume schema has `attachment_mode` specified <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> When volume is created from a snapshot in the same account <ul><li>`is.snapshot.snapshot.operate`</li></ul> When volume is created from a snapshot in a different account <ul><li>`is.volume.volume.allow-remote-account-snapshot-restore`</li></ul> | AE | Create new block storage volumes for attachment | A volume attachment which creates a volume must have the parent resource `update` permission and volume `create` permission. Conditional permissions are checked for snapshots and `attachment_mode` |
| `POST /v1/bare_metal_servers/<bmid>/volume_attachments/` | <ul><li>`is.bare-metal-server.bare-metal-server.update`</li><li>`is.volume.volume.operate`</li></ul> | AEO | Use existing block storage volumes for attachment | Enforce `operate` permission on the parent resource the volume |
| `DELETE /v1/bare_metal_servers/<bmid>/volume_attachments/<volattid>` | <ul><li>`is.bare-metal-server.bare-metal-server.operate`</li><li>`is.volume.volume.operate`</li></ul> | AEO | DELETE volume attachment | Enforce `operate` permission on the parent resource the volume |
| `PATCH /v1/bare_metal_servers/<bmid>/volume_attachments/<volattid>` | <ul><li>`is.bare-metal-server.bare-metal-server.update`</li></ul> When changing the `delete_volume_on_server_delete` property <ul><li>`is.volume.volume.manage-deletion-policy`</li></ul> | AE | Update volume attachment | Enforce `update` permission on the parent. Conditionally enforce `manage-deletion-policy` when the auto-delete property of the volume is being changed (see [link](https://github.ibm.com/riaas/api-spec/issues/1708#issuecomment-58259207) for historical context). |
| `POST /v1/instances` | <ul><li>no changes</li></ul>When volume schema has `attachment_mode` specified <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> | AE | Create VSI and attach new block storage volumes | No changes to the default permission checks. Conditionally, enforce the volume `manage-attachment-mode` permission when the schema has `attachment_mode` specified |
| `POST /v1/instances/<vsiid>/volume_attachments/` | <ul><li>no changes</li></ul>When volume schema has `attachment_mode` specified <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> | AEO | Create new block storage volume for volume attachment | Conditionally, enforce the volume `manage-attachment-mode` permission when the schema has `attachment_mode` specified |
| `PATCH /v1/instances/<vsiid>/volume_attachments/<volattid>` | <ul><li>no changes</li></ul> | AE | Update volume attachment | Existing behavior will remain unchanged to maintain backward compatibility. The volume `manage-deletion-policy` will *not* be enforced when changing the auto-delete property of the volume. This is a deviation from the newly introduced Bare Metal behavior. |
| `DELETE /v1/instances/<vsiid>` | <ul><li>no changes</li></ul> | AE | Delete a VSI | Existing behavior will remain unchanged to maintain backward compatibility. Cascading volume deletes will *not* enforce the volume `delete` permission. This is a deviation from the newly introduced Bare Metal behavior. |
| `POST /v1/volumes` | <ul><li>no changes</li></ul> When the volume specifies the `attachment_mode` property <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> | AE | Create volume | Conditionally enforce the `manage-attachment-mode` permission when the `attachment_mode` property is specified. |
| `PATCH /v1/volumes/<volid>` | <ul><li>no changes</li></ul> When the volume specifies the `attachment_mode` property <ul><li>`is.volume.volume.manage-attachment-mode`</li></ul> | AE | Update volume | Conditionally enforce the `manage-attachment-mode` permission when the `attachment_mode` property is specified. |

Most IAM actions already referenced in this table already exist. The new actions are:

* `is.volume.volume.manage-attachment-mode`
* `is.volume.volume.manage-deletion-policy`

### 3.2.2 Performance, Scalability, and Resource Consumption

Without additional development from AMD/Pensando, there is no opportunity to offload storage encryption
calculation to the programmable P4 h/w pipelines. Instead, CPU cycles on the Bare Metal server and
Acadia storage nodes will be consumed to support encryption-in-transit and encryption-at-rest.

To support EIT over `nvme-tcp`, TLS will be used. The TLS encryption operations will be performed by
the Bare Metal server's CPUs (initiator) and the NVMeoTCP gateway node's CPUs (target). For
encryption-at-rest, the encryption operations will be wholly performed by the NVMeoTCP gateway
node's CPUs.

Adding an NVMeoTCP gateway to the data path presents several performance and scaling considerations:

* Shifting from a one-to-many storage protocol to a point-to-point presents new performance
  bottlenecks and failure modes that must be addressed by this design.
* Extra hop in the data path means the NVMeoTCP network traffic will double at the storage/OSD nodes
* Co-locating the NVMeoTCP on the storage/OSD nodes could create a competition for system resources
  and in the worst case negatively impact the overall Ceph cluster performance.

Additional testing and modeling is needed to see what resource limits and resource isolation
techniques (e.g. cgroups) should be put in place to limit the impact a gateway could have on the
overall Ceph cluster.

The performance plan must include the following assessments in order to establish performance goals
for the gateway service (both in terms of what a single gateway can achieve and what an Acadia
cluster's gateways in aggregate can achieve):

* The performance characteristics of the NVMeoTCP gateway in terms of latency and throughput
* The impact of the NVMeoTCP gateway on the overall Ceph cluster performance
* The impact of multi-attaching a volume to many clients

### 3.2.2.1 Network and Storage bandwidth sharing

Performing reads/writes to remote block storage consumes network bandwidth. In VSIs, there is a customer
control (`total_volume_bandwidth`) to limit how much network bandwidth block storage can consume. Enforcing
this limit is made possible by using Qemu controls for aggregate storage throughput. In this design,
there are no controls we can implement (no DPU support, no limits in the Bare Metal server OS) to provide
a guaranteed split between storage and network throughput. Customer's will be responsible for managing
how the total Bare Metal server network bandwidth is allocated and used.

### 3.2.3 Operational Monitoring

#### 3.2.3.1 Fragility Analysis

| Fragility                                | Microservice/Module  | Impact | Risk/Observe? |
| ---------------------------------------- | -------------------- | ------ | ------------- |
| Pods crashing                            | All controllers      | Failures to provision, update, delete resources | Yes |
| Kubernetes api-extension server crash/unresponsive | All controllers      | Failures to provision, update, delete resources | Yes |
| Unable to retrieve secrets from Vault    | All controllers      | Failures to provision, update, delete resources | Yes |
| gRPC to NVMe GW fail | `acadia-nvme-targetsubsystem-controller` | Failures or slowness to provision, update, delete resources | Yes |
| Keylore service unavailable              | `kmip-controller` | Unable to attach BYOK volumes | Yes |
| DPU discovery service unable to connect to control plane gRPC | DPU discovery service | Customer unable to discover volume attachments | Yes |
| Slow I/O performance | Elba, NVMeoTCP GW, Ceph | Customer application slowness or possible failures | Yes |

#### 3.2.3.2 Observability

The storage platform, Ceph, already support generating and exporting [metrics](https://docs.ceph.com/en/latest/monitoring/)
via a Prometheus exporter. The Acadia platform is sending the available metrics to Osprey currently.
The suite of metrics is being enhanced to include the new NVMe services. The following metrics will
be exported by the Ceph NVMeoTCP Gateway:

| Metric Category | Metric Name | Type | Description |
| --- | --- | --- | --- |
| **Gateway Information** | `ceph_nvmeof_gateway` | Info | Gateway metadata including SPDK version, gateway version, address, port, name, hostname, group, and load balancing group |
| **Block Device (BDEV) Metrics** | `ceph_nvmeof_bdev_metadata` | Gauge | BDEV metadata with labels: bdev_name, pool_name, namespace, rbd_name, block_size |
| | `ceph_nvmeof_bdev_capacity_bytes` | Gauge | BDEV capacity in bytes |
| | `ceph_nvmeof_bdev_reads_completed_total` | Counter | Total read operations completed |
| | `ceph_nvmeof_bdev_writes_completed_total` | Counter | Total write operations completed |
| | `ceph_nvmeof_bdev_read_bytes_total` | Counter | Total bytes read |
| | `ceph_nvmeof_bdev_written_bytes_total` | Counter | Total bytes written |
| | `ceph_nvmeof_bdev_read_seconds_total` | Counter | Total time spent servicing read I/O |
| | `ceph_nvmeof_bdev_write_seconds_total` | Counter | Total time spent servicing write I/O |
| **SPDK Reactor Metrics** | `ceph_nvmeof_reactor_seconds_total` | Counter | Time reactor thread active with I/O (busy/idle modes) |
| **Subsystem Metrics** | `ceph_nvmeof_subsystem_metadata` | Gauge | Subsystem configuration metadata (NQN, serial, model, allow_any_host, HA enabled, group) |
| | `ceph_nvmeof_subsystem_listener_count` | Gauge | Number of listener addresses per subsystem |
| | `ceph_nvmeof_subsystem_host_count` | Gauge | Number of hosts defined per subsystem |
| | `ceph_nvmeof_subsystem_namespace_limit` | Gauge | Maximum namespaces supported |
| | `ceph_nvmeof_subsystem_namespace_count` | Gauge | Number of namespaces per subsystem |
| | `ceph_nvmeof_subsystem_namespace_metadata` | Gauge | Namespace information (NQN, NSID, bdev_name, ANA group ID) |
| **Host Connection Metrics** | `ceph_nvmeof_host_connection_state` | Gauge | Host connection state (0=disconnected, 1=connected) |
| | `ceph_nvmeof_host_keepalive_timeout` | Gauge | Host keepalive timeout indicator (0=no, 1=yes) |
| **Network Interface Metrics** | `ceph_nvmeof_subsystem_listener_iface_info` | Gauge | Interface information (device, operstate, duplex, MAC address) |
| | `ceph_nvmeof_subsystem_listener_iface_speed_bytes` | Gauge | Link speed of listener interface |
| | `ceph_nvmeof_subsystem_listener_iface_nqn_info` | Gauge | Subsystem usage of NIC device |
| **Performance Metrics** | `ceph_nvmeof_rpc_method_seconds` | Gauge | Runtime of RPC method calls |

On top of the standard Kubernetes metrics that will be generated for the new controllers, other
metrics will be added which will be used to detect issues between the central components of this
project (i.e. control plane, Ceph NVMeoTCP Gateway):

| Condition | Metric Type | Alert Condition |
| --- | --- | --- |
| K8s resource processed successfully | Counter | n/a |
| K8s resource processing failed | Counter | >1 % of overall request count |
| K8s resource processing blocked | Counter | >1 % of overall request count |
| K8s resource processing queue for retry | Counter | >1 % of overall request count |
| Nvme namespace attachment time | Histogram/Time | Duration > 100ms |
| Nvme namespace detachment time | Histogram/Time | Duration > 100ms |
| Nvme target subsystem creation time | Histogram/Time | Duration > 100ms |
| Nvme target subsystem deletion time | Histogram/Time | Duration > 100ms |

#### 3.2.3.3 Metrics and Golden Signals

Provisioning and consumption metrics for Volumes and Bare Metal Servers remain unchanged. These
metrics will continue to be sent by the regional control plane.

All other metrics emitted by the control plane, Elba DPU, and Acadia storage platform will flow
into Osprey. Golden signals will continue to be refined/defined as this project progresses:

| Signal | Metric | Anomaly |
| --- | --- | --- |
| K8s resource processed failed | Counter | >1 % of overall request count |
| Control plane controller resources (cpu/mem) consumed | Gauge | 5 min avg > 80%,2G |
| Nvme Gateway resources (cpu/mem) consumed | Gauge | 5 min avg > 80%,2G |
| Latency to attach/detach volumes | Histogram/Time | Duration > TBD |

#### 3.2.3.4 Events

Existing billing related events for Volumes and
Bare Metal Servers will continue to be emitted. See section [3.2.4](#324-lifecycle-metering-and-billing).

The zonal control plane will continue to emit existing events related to Bare Metal servers, volumes,
and virtual network interfaces. New events will be [defined](https://github.ibm.com/genctl/shared-eventing/pull/601/files)
and emitted for new resources including NVMe storage devices, NVMe path groups, and NVMe target subsystems.
Fatal errors which are deemed unlikely to resolve without human intervention will raise alert events.
All events will be generated using the genctl/shared-eventing library.

#### 3.2.3.5 Logging

All new control plane services and data plane services running on the Acadia platform (i.e. NVMeoTCP
gateway) will send all logs to IBM Cloud Logs and metrics to Osprey.

At this time, no log-based alerting is planned. Metrics based alerting is preferred.

#### 3.2.3.6 Alerting

The VPC and Acadia SRE teams will be responsible for handling alerts related to Bare Metal servers
with external storage. Many alerts already exist, and those will continue to be used. For example,
alerts for long running bare metal server [provisions](https://github.ibm.com/cloudlab/ng-sysdig-osstip/blob/master/alerts/is-bare-metal-server/baremetal/long-running-provision.json)
or long running bare metal server [power on](https://github.ibm.com/cloudlab/ng-sysdig-osstip/blob/master/alerts/is-bare-metal-server/baremetal/long-running-poweron.json)

New alerts related to the Elba DPU and Ceph NVMeoTCP gateway are TBD.

#### 3.2.3.7 Dashboards

Both `Bare metal server` and `Volume` objects are already reflected in the Ops Dashboard. A new
resource, `NVMe Target Subsystem` will also be added to the Ops Dashboard.

The Acadia SRE team maintains [dashboards](https://github.ibm.com/acadia-sre-and-operations/dashboards)
to monitor the health of Acadia storage platforms. These dashboard are also available via the Ops
Dashboard via the via the [Common Grafana Dashboard](https://opsdashboard.w3.cloud.ibm.com/graphs/dashboards/f/folder-storage/storage?tag=acadia).
The Acadia dashboards will be enhanced to alert and report on the health of the new NVMe components.

#### 3.2.3.8 Responsibility Matrix

No deviations from responsibility matrix planned at this time.

#### 3.2.3.9 SSAD and Approval Checklist

`is.volume` has an approved [SSAD](https://github.ibm.com/cloud-governance-framework/services-governed-content/blob/main/arch-design/ssad/is.volume/is.volume-ssad.md)
which will be updated to reflect extending support for volume migration, including updates for
SSAD section `6.2 Monitoring`.

#### 3.2.3.10 Additional Observability Strategy

No.

#### 3.2.3.10 Troubleshooting

No further troubleshooting tools planned at this time.

### 3.2.4 Lifecycle, Metering, and Billing

- No new top level resources
- Existing regional quota limits for Bare Metal servers and Volumes still apply
  - Maximum number of volumes attached to a Bare Metal server instance is controlled by feature
    flags `is-bare-metal-server-volume-attach-limit-value`. The default limit will be 16 and enforced
    by regional control plane
- Resource lifecycles for Bare Metal servers and Volumes are not being changed
- Billing for Bare Metal servers and Volumes is not being changed

### 3.2.5 Compliance

No new hardware being added with this SRB.

### 3.2.6 Faults and Disasters

#### 3.2.6.1 Fault Tolerance

Control plane fault tolerance follows the existing design.

Data plane fault toleration follows the existing Acadia platform design
(3 way data replication, etc). NVMeoF gateways will be deployed in sets of 4 so that there is no
single point of failure in the NVMeoTCP data path.

#### 3.2.6.2 Business Continuity Disaster Recovery (BCDR)

No new persistent data store being added. After a zone/region outage, the kubernetes control plane
will reconstruct all the Bare Metal server to block storage volume attachments.

### 3.2.7 Hardware Platform Support

No impact.

### 3.2.8 Component Level Development

| Area | Section(s) |
| --- | --- |
| Bare Metal | All |
| CLI | [3.1.2.2](#3122-command-line-interface-cli) |
| Compute | [2.2](#22-consumer-interaction-model), [3.1.3](#313-customer-facing-apis) |
| Network | All sections under project overview [2](#2-project-overview), [3.1.1.2](#3112-network-design) |
| Observability | [3.2.3](#323-operational-monitoring) |
| Performance | [3.1.1.1](#3111-high-level-overview), [3.1.1.3](#3113-ceph-nvmeof-gateway-overview), [3.2.2](#322-performance-scalability-and-resource-consumption) |
| Platform Integration | [3.1.3.1](#3131-activity-tracker-at) |
| RE/RM | [2.6.1](#261-acadia-platform-updates), [2.6.2](#262-client-enablement), [3.1.5](#315-deployment-architecture) |
| RIAS API | All |
| SDK | [3.1.3](#313-customer-facing-apis) |
| Security | [2.7.3](#273-cloud-security-baseline-requirements), [3.1.1.5](#3115-authentication-and-authorization), all sections under [3.2.1](#321-security) |
| Storage | All |
| Support | All, emphasis on [3.2.3](#323-operational-monitoring) |
| Terraform | [3.1.2.3](#3123-terraform-and-packer) |
| Test | [4](#4-testing) |
| UI | [3.1.2.1](#3121-user-interface-ui) |

### 3.2.9 Network Topology and Segmentation Impact

The control plane flows follow what was established with the initial [Acadia SRB](https://github.ibm.com/cloudlab/srb/tree/master/proposals/1793).
To program the NVMeoTCP gateways, the `genctl` controllers will send gRPC+mTLS requests on
the underlay network. To program the NVMe agent on the Elba DPUs, the `genctl` controllers will send
gRPC+mTLS requests on the VCN.

To allow accessing block storage using an Elba DPU requires some updates to the Acadia data path,
specifically enabling the NVMeoTCP protocol. The underlay network is used for the storage data
plane between Bare Metal servers and Acadia platform. An overview of the data plane follows:

![BM NVMe Data Path Overview](./materials/BM-NVMe-Network-Flows.png)

Note that data is always encrypted in transit using either NVMe secure channels (TLS) or Ceph
`msgr2`.

# 4. Testing

> All projects are expected to be designed to be easily testable. How can one effectively test
> what's proposed?  If this project has a dedicated test plan, please link to it. If a test plan is
> not available, specify what areas of VPC will require testing, the applicable types of testing
> (unit, regression, functional, performance, scalability, stress, system, compatibility, boundary,
> fuzz, etc.), and provide general guidance on how best to exercise and validate the changes.
> See the [test left FAQ](/architecture/testing/test_left_faq.md) for more information.
>
> List any aspects of this project that will pose particular challenges for automated tests, and any
> aspects of the project that will require white-box or other brittle forms of testing. Are there
> aspects of this project that require test results to be manually curated? Are there aspects of
> this project that likely require ongoing test maintenance even if the code itself is not
> undergoing change? How will the tests and code be kept in sync?
>
> Testing must cover any new authorization flows described in [IAM
> Integration](#3215-identity-and-access-management). Examples:
>
> - As a secondary user with fine grained access, I can access the resource
> - As an unauthorized user, I cannot access the resource
> - As a secondary user with *only* access tags, I can access the resource
> - With a CBR policy in place, as an administrator I can access the resource _only_
> if the request originates within the defined network zone
>
> For more information, see the [guidelines for IAM-related
> testing](architecture/platformIntegration/genesis-iam/IAMIntegrationChecklist.md#iam-test#iam-testingchecklist)
>
> If this project does not anticipate requiring any changes to testing, please justify.
TODO

## 4.1 CI Pipelines - Automation

- New `genctl` CI pipelines have been created to support the new [acadia-baremetal-storage-workspace](https://github.ibm.com/genctl/acadia-baremetal-storage-workspace)
- Acadia platform CI pipelines updated to support new Ceph containers and services required for
  NVMeoTCP attachment

# 5. Solution Space

## 5.1 Lessons Learned

This is the first attempt at adding external block storage to the Bare Metal server offering in VPC
so this SRB attempts to borrow much of the design and customer interaction model from block storage
with VSI offering in VPC.

## 5.2 Alternatives Considered

- Running the Ceph client (`librbd`) on the DPU to provide native access to the Acadia storage
  platform
  - IBM Research evaluated this option and found running the `librbd` software stack on the ARM
    cores of the DPU was not performant
- Present volumes as locally attached PCIe NVMe devices
  - Requires additional function in our DPUs. Our DPU vendor has not committed to providing the
    additional function required to make this approach possible.
- Supporting a multi-attach model with finer grained access control.
  - Allowing customers more control over the read/write capabilities for each attachment is desirable.
    Customers could create a read-only attachment for a read-only workload, or a read-write attachment
    for a read-write workload. This design would allow sharing of a volume, but provide guard rails
    to ensure only one attachment is writing to the volume at a time.
  - This level of control is not a requirement for this proposal. Additionally, the storage datapath
    would need enhancements to restrict attachments to a read-only mode.
- Using a deterministic NQN based on the product UUID set in the server's system firmware ((see
  section `7.2.1 System — UUID` of the [SMBIOS specification](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.9.0.pdf)))
  - Although this approach would give customers a consistent NQN value from the RIAS public API and
    standard NVMe tooling, there are challenges backfilling the NQN for existing servers and updating
    the NQN during a a server hardware failure that requires a chassis swap.
- Using sub resource IAM actions for Bare Metal volume attachments
  - Sub resource IAM actions aren't prevalent in the api-spec.
  - Introducing sub resource IAM action enforcement would at a minimum, cause a deviation in the
    Bare Metal server APIs and possibly the Instance APIs (if it was decided to bring the same
    functionality to instances) from the remainder of the api-spec endpoints.
  - The current proposal does not require sub resource IAM actions.
  - Sub resource IAM modeling should be driven across the entire api-spec in a single proposal, instead
    of piecemeal implementations which will lead to inconsistencies in the APIs for an indeterminate
    period of time.
- Using customer BYOIPs to provide `nvme-tcp` connectivity
  - This was deemed too burdensome for customers, especially those that have a limited range of free
    IP addresses within their VPC.

## 5.3 Competitive Offerings

- [AWS Nitro](https://aws.amazon.com/ec2/nitro/) based systems providing [NVMe storage](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/nvme-ebs-volumes.html)
- [Azure Baremetal Infrastructure](https://learn.microsoft.com/en-us/azure/baremetal-infrastructure/concepts-baremetal-infrastructure-overview#storage)
  support NVME over Fibre Channel
- [Oracle iSCSI](https://docs.oracle.com/en-us/iaas/Content/Block/Tasks/connectingtoavolume_topic-Connecting_to_iSCSIAttached_Volumes.htm)
- [AWS Multi-Attach](https://docs.aws.amazon.com/ebs/latest/userguide/working-with-multi-attach.html)
- [Google Cloud Shared Disks](https://docs.cloud.google.com/compute/docs/disks/sharing-disks-between-vms)

# 6. Known Issues

- Measurements and modeling needed to observe how much CPU/mem the NVMeoTCP gateways require across
  various workloads.
- Measurements and modeling needed to ensure Acadia OSD nodes have capacity to meet the Bare Metal
volume demand.
- Need to define performance and scaling targets for the offering prior to beta commit
  review

# 7. Patents

None

# 8. Meta

No

# 9. Glossary

- AEN: Asynchronous Event Notification
- BM: Bare Metal
- BYOK: Bring your own key
- DPU: Data processing unit
- NVMe: NVM Express - An interface that allows host software to communicate with a non-volatile
  memory subsystem
- NQN: NVMe Qualified Name
- NVMe Namespace: The NMVe representation of a block storage volume, aka `RBD` in Ceph terms
- NVMe-oF: NVMe over Fabrics - A common architecture that supports NVMe block storage protocol over
  a storage networking fabric
- NVMeoTCP: NVMe over TCP - An NVMe-oF transport using TCP (transmission control protocol)
- OSD: Object storage daemon - it is responsible for storing objects in an associated device
- RADOS: Reliable Autonomic Distributed Object Store - Ceph software defined
  storage solution is based on RADOS
- RBD: RADOS block device - the block abstraction provided by Ceph
- SPDK: Storage Performance Development kit

# 10. References

- The original source for all diagrams is in [genctl-acadia/team-docs](https://github.ibm.com/genctl-acadia/team-docs/tree/master/control-plane/bare-metal)
- The majority of design documents and reference materials can be found in [lighthouse](https://w3.ibm.com/services/lighthouse/wikis/view/73e97cc0-a119-486e-858d-22bd17be6e3e/6ee820c2-5561-4319-a490-d1bddeaa31d9)

<!-- Emacs formatting setup -->
<!-- Local Variables: -->
<!-- fill-column: 100 -->
<!-- End: -->

<!-- Vim formatting setup. Must be in first or last 5 lines by default -->
<!-- vim: set textwidth=100 : -->