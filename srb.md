---
template_version: 2025-07-11 # DO NOT CHANGE -- this comes from the copied template.md

## These _MUST_ be updated by the mentor as the proposal goes through the SRB process.
state: open # open or closed
status_label: review-follow-up (inception completed) # See https://github.ibm.com/cloudlab/srb#status-labels
last_update: 2026-09-17 # Set to today's date whenever the proposal is updated
---

<!-- cSpell:ignore Aneesha Gullena discoverability Sprintf rocm cudnn miopen scipy SPDXID -->
<!-- cSpell:ignore WCAG VHDs DKMS oneapi dpkg dkms queryformat -->
<!-- cSpell:ignore jjfall Sarith Vijayan RIAS ibmcloud QCOW VPCIMG BYOK KYOK ROKS bootable qcow -->
<!-- cSpell:ignore publicization Sysdig diagnosability NACL -->
<!-- cSpell:ignore psirt PAYG nccl hclwrite bluemix agullen sthornt genctl IFV RIOS -->
<!-- cSpell:ignore cyclonedx -->
<!-- cSpell:ignore CUDA cuda Gaudi gaudi Velero FIPS SSAD lifecycles textwidth -->
<!-- cSpell:ignore Tekton toolchain CBR MZone MZones etcd toposort IPOPS VALCMD -->
<!-- cSpell:ignore cmds buildable templating inspectable hashable unpassable customizers -->
<!-- cSpell:ignore cxf sbom onepipeline Flex CISO -->
<!-- cSpell:ignore HashiCorp Artifactory CycloneDX NGDC pollable DATACLASS -->
<!-- cSpell:ignore unconfigured requiredness vcpu Astha Deepak -->
<!-- cSpell:ignore Sneha Sequeira BCDR misattribution -->

# 1. Introduction

## 1.1 Project Name

Image Builder Service - Phase 1

## 1.2 Submitter and Contributors

- Aneesha Gullena, <agullen@us.ibm.com>,
  [@aneesha](https://ibm.enterprise.slack.com/team/W4EDAKWR3)

## 1.3 Mentor(s)

- Jim Fall, <jjfall@us.ibm.com>,
  [@jjfall](https://ibm.enterprise.slack.com/team/W45DRF7RP))

## 1.4 Offering Manager(s)

- Sarith Vijayan, <sarith.vijayan@ibm.com>,
  [@Sarith.Vijayan](https://ibm.enterprise.slack.com/team/U04E9RLCC0N))

## 1.5 API focal(s)

- Sean Thornton, <sthornt@us.ibm.com>,
  [@sean.thornton](https://ibm.enterprise.slack.com/team/W3Q7CJKCH)
- Sneha Sequeira, <sneha.sequeira@ibm.com>,
  [@Sneha.Sequeira](https://ibm.enterprise.slack.com/team/W015PP1N0LC)

## 1.6 Review Flow (determined by mentor)

[`Full`](README.md#full-flow)

## 1.7 Production Target Date

4Q'26

# 2. Project Overview

## 2.1 Overview

This proposal introduces a new Image Builder capability for VPC, delivered as a feature
extension of the existing image service. It adds two new top-level regional resources,
`image-recipe` (CRN service name `is.image-recipe`) and `image-component` (CRN service
name `is.image-component`), and lets customers create customized virtual machine images,
pre-configured with GPU drivers and specialized software, without having to provision or
manage a virtual server instance themselves. Right now, IBM Cloud customers have no
repeatable, observable way to produce such images, and this closes that gap at no service
cost. Customers pay only for the ephemeral build VSI's runtime, roughly $0.09 and 30 minutes
per build.

The design rests on two new resource types. A recipe describes the image to be built: a base
operating system image plus a chosen set of components, the "ingredients." Components are
IBM-curated software packages. In Phase 1 these are the GPU driver stacks for NVIDIA CUDA,
AMD, and Intel Gaudi, and each one carries a declared operating system and architecture
compatibility, checked
against the recipe's base image. Customers
assemble recipes from the component catalog, or start from an IBM-provided stock recipe, then
trigger builds that produce ready to use custom images in their own account. Each image
carries a recorded source recipe reference and a downloadable software inventory (SBOM).

Phase 1 lays the foundation. New IBM control-plane components, an Image Builder API
controller and a Status Watch controller detailed in 3.1.1.1, orchestrate the pipeline above
the cloud, while the ephemeral
build VSI runs in the customer's account. This is a hybrid model, consistent with how both
AWS EC2
Image Builder and Azure VM Image Builder work. On top of that we're building an automated
pipeline with in build package validation and lifecycle governance for everything in the
catalog, along with the operational groundwork a service like this needs: support channels, a
trained SRE team, and so on. None of this starts from scratch. The design leans heavily on
HashiCorp Packer and our own packer plugin ibmcloud, the IBM Cloud Continuous Delivery Tekton
service as the pipeline runtime, and the VPC image service's existing image archive, its
"create" API extension point, and its stock content catalog management model.

## 2.2 Consumer Interaction Model

Image Builder is customer facing, discoverable through several entry points in the IBM Cloud
ecosystem: a new Image Builder tile in the IBM Cloud Catalog, or the option to create an image
from a recipe right inside the existing image service GUI flows. Customers can use whichever
interface they already prefer, whether that's the VPC GUI, REST API, CLI, or Terraform, and get the same
functionality from all of them. There are a few different paths into the pipeline itself:

- Starting in image builder:
  - **Path 1, Stock recipe**: Pick an unmodified IBM provided recipe and trigger a build
    from it, same as any other recipe. Since every build produces a new image, the customer
    ends up with an image of their own, a real, customer account owned image, not the
    IBM maintained stock image itself. Picking a stock recipe doesn't drop an existing
    IBM maintained stock image into the customer's account on its own; it's just the starting
    content for a real build, and that build always produces a genuinely new image. No
    pre built or shared images exist for any recipe, stock or custom.
  - **Path 2, Start from a stock recipe**: Copy a stock recipe's content into a new custom
    recipe. The console pre populates the create form with the stock recipe's
    current values; the customer edits as they like and submits it as an entirely ordinary
    new recipe with its own name and initial content hash version. It's really just a
    client side
    convenience the VPC GUI gives you, built on the plain `GET v1/image_recipes`/
    `POST v1/image_recipes` calls. There's no dedicated mechanism behind it, no
    fork/derivation API, no lineage field. It's copy and paste, not a fork.
  - **Path 3, Create from scratch**: Build a new custom recipe without starting from an
    existing one. Pick a base VPC stock image, then choose components à la carte from the
    IBM curated catalog (`GET /v1/image_components`). That recipe then produces a custom
    image in the customer's account.

**Every `POST /v1/image/jobs` produces a new image.** Rebuilding from the same recipe, even at
the identical `(id, version)`, even a stock recipe, triggers a real build and produces a new,
distinct image, the same as rebuilding after an edit or against an older version. Recipes are
meant to be reused, but reuse here means "build again," not "just hand me back what's already
there." We did consider a fast complete behavior instead (resolve to an existing image rather
than building again); see 5.2 for that
alternative and why it was not adopted.

> **Note on terminology**: "stock" is used throughout this Phase 1 proposal in the sense
> established by existing VPC images (IBM owned, publicly visible content). This is not a
> naming decision for future vendor or third party owned recipes. What those get called
> (e.g. "vendor recipes," "catalog recipes") is explicitly out of scope for Phase 1 and
> undecided.

| ![Image Builder Consumption Flow](materials/image-builder-consumption-flow.png?raw=true) |
| :---: |
| _Diagram 2.2-1. Image Builder Consumption Flow_ |

The image creation workflow for custom recipes follows a two-step process designed for
flexibility and reusability.

In Step 1, customers create a recipe by picking a base VPC stock image and choosing
components from the IBM-curated catalog. That catalog
(`GET /v1/image_components`) has 10
GPU driver components across three vendor families (NVIDIA, AMD, and Intel), and all of them
use the `precompiled_package` install method in Phase 1, so the build VSI doesn't need GPU
hardware at all. In build validation stays at the package level for Phase 1, since there's no
GPU
hardware present
to validate against anyway (see 3.1.1.3.6a and 2.8.2.9).

In Step 2, customers trigger a build by creating an image job
(`POST /v1/image/jobs` naming the recipe; image identity is
service-derived in Phase 1, more on that below). That kicks off a Tekton pipeline using
`packer-plugin-ibmcloud`, which spins up
a temporary Flex profile (`cxf-4x8`) build VSI
in the customer's account, installs the GPU driver packages from
Artifactory via `dnf`
or `apt-get`, runs package verification and component validation, captures the build
VSI's boot volume and registers it directly as a custom image ("image from volume"), pulls
together an
inventory from the Packer plugin's scan (only the extracted name/version
list sticks around; see below), and
cleans up every ephemeral resource (VSI, SSH key, floating IP, security group) no matter how
the build ends, including failures. The
customer only pays for the build VSI's runtime, roughly $0.09 and 30 minutes per
build. Since every build produces its own new image, this two-step approach lets customers
reuse the same recipe over and over for consistent images.

Customers keep track of the build through the image job, a first-class
regional resource (`GET /v1/image/jobs/{id}`) whose coarse `status` moves through
`queued → running → succeeded | failed | deleting`. There's no pollable
`cancelling`/`cancelled` state (see 3.1.1.9.3a). `status_reasons` doesn't
carry live stage-progress while the job is `running`. Internally the pipeline
does move through provisioning, installing, validating, and capturing, but none of
that shows up to the customer via the API. `status_reasons` only gets populated on
`failed`, from a set of 6 terminal failure
codes (see 3.1.1.9.3a), and stays empty the whole time the job is running.

A user-configurable `timeout` (default 120 minutes,
30–480 minutes) exists as a circuit breaker, if the image job hangs or gets stuck during
internal processing, the timeout's expiration cancels the job, ensuring the temporary build
VSI doesn't sit running (and billing) indefinitely in the customer's account if something goes
wrong. It's scoped to the customer paying window: the timer starts
when the job hits `running`, and the job counts as complete, stopping the timeout, once the
image exists, even if the image's own status is still `pending`. Capture and upload both have
to finish for the job to get there, since that's when the image actually has content; time
spent queued doesn't count against the clock. The job also references
the build worker and the recipe `(id, version)` it's building from, and a customer can pin a
specific historical recipe version for the build
(`ImageRecipeIdentity.version`; the latest version is used if this is left out).

The job's request body can optionally carry an `image_template`, distinct from
`image_template.name`/`image_template.resource_group`. Leave it out and the image name
gets generated from a recipe-name prefix plus a random suffix (the same convention instance
groups use), with the resource group inherited from the recipe's. Things like `allowed_use`
and lifecycle dates only become settable through an ordinary image PATCH after creation. The
prototype doesn't reach that far.

The image is the job's output. The plugin
calls `ImageBySourceVolume` → `POST /v1/images` itself during Task 1's Capturing stage. There's
no placeholder created earlier. The
image shows up with its real content already there, in one step. The plugin's
capture step then polls the image's own status and blocks until it hits `available`
before the Tekton task finishes, so `job.status` only reaches `succeeded` once zone
distribution is fully done, not as soon as the image is merely created. Zone distribution
across every zone happens automatically, inherited for free from the plugin's use of the
standard
`ImageBySourceVolume` call rather than anything Image Builder does itself. The image moves
`pending` →
`partially_available` → `available`. A build-phase failure produces a failed job with no
surviving image; the whole failure lives on the job, and its `status_reasons` names the failed
stage and component, with actionable messages and a support reference ID for internal errors.
The image
object
carries the `source_image_recipe`
reference and an inventory descriptor pointing to the downloadable, extracted
name/version list (`GET /v1/images/{id}/inventory`). Cancelling a
build is just `DELETE /v1/image/jobs/{id}`. The pipeline run gets cancelled and
cleanup always runs. Deleting a terminal job just removes the historical record
without touching any produced image. The
deletion-scenarios table in 3.1.1.9.3a walks through how deletion interacts with each build
stage.

The table below traces both resources together through every path. The image doesn't
exist until the plugin's own `ImageBySourceVolume` call succeeds, during Task 1's Capturing
stage. The stage names in
parentheses describe the pipeline's internal progression, not customer-visible
`status_reasons` content. `status_reasons` stays empty
while the job is `running`. These labels are here purely to show what's actually happening
under the hood at each point:

| Phase | Image job `status` (internal pipeline stage) | Image | Customer billing |
| --- | --- | --- | --- |
| Job accepted (`POST /v1/image/jobs`) | `queued` | does not exist yet | none |
| Worker VSI provisioning | `running` (internal: provisioning) | still does not exist | build VSI starts billing |
| Components installing | `running` (internal: installing) | still does not exist | VSI |
| Validation | `running` (internal: validating) | still does not exist | VSI |
| Boot volume capture: plugin calls `ImageBySourceVolume` (Task 1, Capturing stage) | `running` (internal: capturing) | **created here**, with real content immediately | VSI, still billing |
| Same stage, plugin polls and blocks on the image reaching `available` | `running` (internal: capturing, the same stage, not a separate one) | `pending` → `partially_available` (a real intermediate state) → `available` | VSI, still billing (the task has not completed yet) |
| Task 1 completes once `available` observed; `Cleanup()` runs | `running` | `available` | **VSI billing stops** |
| Task 2: inventory extraction, AT event | `running` | `available` (already, from the prior stage) | — |
| **Success** | `succeeded` (terminal; `completed_at` set) | `available`, ready for VSI provisioning; carries `source_image_recipe` + inventory reference | image storage |
| **Build-phase failure** (before the plugin's `ImageBySourceVolume` call) | `failed`, one of 6 category codes naming stage/component (see 3.1.1.9.3a) | none, the image cannot exist yet at this point | bounded to VSI runtime; cleanup ran |
| **Cancellation** (`DELETE /v1/image/jobs/{id}`) | `deleting` → record removed (non-terminal) | Before the plugin's call, nothing exists to delete; after it (and before `succeeded`, which only fires once `available`), the existing image is deleted | stops at cleanup |
| Build timeout (default 120 min, 30–480 range) | `failed` (`failed_timeout`) | — | bounded by `timeout`, covering the whole of Task 1, including the plugin's blocking wait for `available` |

The detailed derivation (which upstream signal drives each transition) is in 3.1.1.9.3a.

## 2.3 Requirements and Use Cases

### 2.3.1 Overview

The IBM Cloud VPC Image Builder capability addresses the need to streamline the creation,
customization, and deployment of secure machine images, improving efficiency, reducing costs,
and ensuring consistent environments for scalable application management.

Aha Reference:
[VPCIMG-277: Image Builder Tool - Phase 1](https://bigblue.aha.io/features/VPCIMG-277)

### 2.3.2 Core Requirements

#### 2.3.2.1 Service Agility and Responsiveness

Priority: High

The VPC Image Builder capability must be architected with agility and efficiency as foundational
principles to enable the VPC program to be highly responsive in delivering new image recipes,
stock images, and software components to customers. Specifically:

- **Decoupled from CD Pipeline**: New stock images, software
  components, and recipes need to reach VPC customers without waiting on the VPC continuous
  delivery
  (CD) pipeline. The **dynamic catalog property** guarantees this structurally:
  a component's publish payload carries everything the build derives from (`software`,
  `dependencies`); the HCL
  generator is a pure
  schema-to-template transform with no per-component code; and Artifactory holds binaries only —
  no authored install or validate scripts exist anywhere. A new component is buildable the
  moment it's published, with zero code deploys.
- **Rapid Component Updates**: Support expeditious delivery of new image components
  when any of the following become available:
  - New base (stock) system images
  - New software components
  - New software component versions

- **Independent Release Cadence**: The image builder capability must operate on its own release
  schedule, independent of broader VPC platform releases, to maximize responsiveness to customer
  needs and market demands.

#### 2.3.2.2 Competitive Parity

The service must provide feature parity with image builder services available from competitive
cloud service providers, including but not limited to:

- **AWS EC2 Image Builder**: Pipeline-based image creation with automated testing and
  distribution
- **Azure VM Image Builder**: Template-based customization with integration into Azure DevOps
- **Red Hat Image Builder**: Hybrid/multi-cloud image creation with enterprise Linux focus

Key parity features include:

- Component-based image customization
- Automated build pipelines
- Image validation and testing
- Version management and tracking
- Recipe reusability
- Integration with existing cloud services

No major competitor ships stock recipes — first-class, buildable, lifecycle-managed
IBM-curated recipes are a genuine differentiator here. AWS ships managed components
and parent images, but recipes there are always customer-owned; Azure doesn't even have a
component concept. For
NVIDIA and GPU workloads especially, just getting a working environment set up today is a
real source of customer friction. Recipes help beyond what a flat stock image can offer:
customers can see exactly which components went into a base image, and extend that baseline
with their own additions on top of it.

**Note on recipe reuse**: both AWS EC2 Image Builder and Azure VM
Image Builder are built around repeatedly rebuilding from the same
recipe or template over time, including scheduled rebuilds that produce updated, patched
images — that's core to their value proposition, not a side case. VPC Image Builder's recipes
work the same way: every build against a recipe, whether it's the same version or a new one,
produces its own new image, which matches this reuse pattern rather than falling short of it.

#### 2.3.2.3 Image Customization Capabilities

The service must enable users to:

1. **Select Starting Point**: Choose an existing stock recipe or an IBM Cloud stock image
2. **Progressive Customization**: Add customization components incrementally to build complex
   image configurations
3. **Pipeline Execution**: Invoke automated pipeline processes that include:
   - Build stage: Apply customizations to base image
   - Validation stage: package verification plus per-component validation commands

4. **Recipe Management**:
   - View rendered recipes showing all included components and customizations
   - Save recipes for future reuse
   - Automatic content-hash versioning with a complete, addressable version history

5. **Image Retrieval**: Obtain the resulting custom image in the user's cloud account

#### 2.3.2.4 Software Component Governance

The service must implement robust governance for software components:

- **Curated Components**: Guarantee that all software components the system makes available
  are curated and validated
- **Metadata Tracking**: Maintain comprehensive image metadata including:
  - Complete record of all included software components (as-installed versions, via the
    extracted inventory)
  - Component versions
  - Build date and configuration
  - Source recipe `(id, version)` coordinate

- **Lifecycle Management**: Track End-of-Life (EOL) and End-of-Service (EOS) status for all
  software components via explicit `deprecation_at` / `obsolescence_at` dates from which
  component status is derived
- **Vulnerability Monitoring Foundation**: Explicit lifecycle dates and a downloadable software
  inventory per image give a basis for identifying which images contain outdated components
- **No In-Place Updates**: Images are immutable; updates require building new images with
  updated components

### 2.3.3 Use Cases

#### 2.3.3.1 GPU-Enabled Images

Create custom images containing GPU drivers and supporting software for use with VPC GPU profiles:

- **NVIDIA CUDA**: Images with NVIDIA Compute Unified Device Architecture (CUDA) drivers and
  toolkit
- **AMD GPU**: Images with AMD GPU drivers and ROCm platform
- **Intel Gaudi**: Images with Intel Gaudi accelerator drivers and software stack

**Business Value**: Today, producing a working GPU image is a manual, error-prone effort.
Customers must select mutually compatible driver, toolkit, and library versions, match them to
the OS kernel and the target GPU profile, and hand-script the installation, with mistakes
typically surfacing only when the instance boots and the driver fails to load. Image Builder
replaces this with an IBM-governed automated process that produces images purpose-built for
accelerators, with package-level installation verified
in-build (see 3.1.1.3.6a). Hardware/GPU-manufacturer compatibility
across components is not automatically validated (see 3.1.1.3.6a and 2.8.2.9 for that
accepted gap). This
eliminates manual driver installation and reduces time-to-productivity for GPU workloads.

#### 2.3.3.2 Customer Discovery and Self-Service

Address current gaps in the Packer script approach available on github.ibm.com:

- **Discoverability**: Provide UI/API for users to discover available stock recipes, components,
  and custom recipes
- **User Experience**: Offer intuitive interface for image customization
- **Observability**: Track usage patterns and component adoption
- **Self-Service**: Enable customers to build custom images without provisioning or managing
  any build infrastructure themselves (the ephemeral build VSI is provisioned, managed, and
  cleaned up by the service)

**Business Value**: Improves customer experience, increases adoption, provides usage insights
for service improvement, expands the automation toolbox for customers looking to migrate workloads
from IaaS Classic to VPC.

### 2.3.4 Non-Functional Requirements

#### 2.3.4.1 Service Maturity

Phase 1 implementation must establish:

- Dedicated support channel for the service
- Trained Site Reliability Engineering (SRE) team
- Operational runbooks and procedures
- Service monitoring and alerting

#### 2.3.4.2 Technology Foundation

- **Tekton Pipeline**: The IBM-managed Continuous Delivery Tekton pipeline forms the basis
  for the image builder runtime build pipeline. The pipeline consists of two named tasks plus
  a cleanup task. **Task 1 — Packer Execute** performs the build itself: it fetches and
  verifies the generated template, provisions the ephemeral build VSI in the customer's
  account, installs and validates the components, and captures the disk image. **Task 2 —
  Post-process** publishes the results: it assembles the inventory, uploads
  artifacts to Object Storage, and writes the build metadata to the image service. The
  **Cleanup task** deletes all ephemeral resources and runs on every exit path — success,
  failure, and cancellation. HCL generation runs in the Image Builder
  control plane (RIAS) before the pipeline is invoked; the generated template is staged
  per-build, integrity-hashed (`hcl_sha256`), and verified by Task 1 before use. Pipeline
  dispatch uses the Continuous Delivery **REST API** (`CreatePipelineRun`) — not gRPC (see
  3.1.1.3.2 and 3.1.3). See OI-PIPE for the evaluation, still in progress with the CD/Tekton
  team, of whether this dispatches through raw Tekton or through OnePipeline, IBM's shared
  compliance/DevSecOps pipeline platform — both ultimately run on the same underlying CD
  Tekton service, so this is a definitions/dispatch-model decision, not a different execution
  engine; this document is being kept describing the raw-Tekton model through inception
  review, with any change reflected at commitment review once the evaluation concludes.
- **Packer Integration**: HashiCorp Packer (via packer-plugin-ibmcloud, Apache 2.0) is
  used inside Tekton Task 1 to provision the build VSI, install components, capture the image,
  and clean up ephemeral resources. The Packer task container image is **built and owned by
  IBM** rather than consuming the upstream hashicorp image (details in 2.5.3). See 3.1.1.3.5
  for the full Packer authentication/authorization flow.
- **Automation**: Build lifecycle managed by packer-plugin-ibmcloud — no manual cleanup
  required.

#### 2.3.4.3 Phased Delivery

**Phase 1 (2026)**:

- Support for 1-2 stock images (RHEL initially, then Ubuntu, Windows)
- Limited set of customization components focused on GPU drivers
- Foundation for future expansion
- No licensing or billing requirements for initial components

See section 2.4.2 "Phased Delivery" for the projected phase outlook.

### 2.3.5 Scope Clarifications

This project (Phase 1) specifically covers:

- Initial service infrastructure and pipeline
- GPU driver customization components (NVIDIA, AMD, Intel Gaudi)
- Any VPC stock image usable as a base image, no API-level restriction (see 3.1.1.6)
- Recipe management with content-hash versioning and history, and image catalog integration
- Foundation for governance and lifecycle management

### 2.3.6 User Stories

#### 2.3.6.1 Epic 1: Create a Custom Recipe from Stock Recipe

##### 2.3.6.1.1 User Story 1.1: Create Custom Recipe from IBM-Provided Stock Recipe

**Title**: Create Custom Recipe Using IBM Stock Recipe

**As a** Platform Engineer or DevOps Engineer supporting AI/ML workloads,
**I need to** create a new custom image recipe using an IBM-provided stock recipe,
**So that** I can quickly start with a pre-configured, best-practice recipe and optionally customize
it for my specific infrastructure needs.

**Acceptance Criteria**:

1. User can view all available IBM-provided stock recipes in the interface
2. Stock recipes display relevant metadata (base image, included components)
3. User can select a single stock recipe to create a new custom recipe
4. Selected stock recipe provides the baseline for a new editable recipe
5. New recipe inherits all configurations from the stock recipe; inherited components and
   configuration can be edited during creation and changed after creation (content changes
   mint a new version)
6. User can modify or delete the recipe after creation without affecting the original stock
   recipe
7. Recipe creation includes a creation timestamp
8. User receives confirmation when recipe is successfully created from stock recipe
9. Created custom recipe is retrievable and usable by the creating account but not by other
   accounts

**Story Category**: Functional

**Considerations**:

- Handle deprecated or obsolete stock recipes gracefully (console pre-population from a
  deprecated stock recipe carries the deprecation warning; obsolete stock recipes are not
  offered as starting points)
- Stock recipes are readable by any authenticated user without
  recipe-specific IAM grants; IAM permissions gate creation and management of the customer's
  own recipes
- Handle stock recipe content updates without breaking existing recipes (the customer's
  recipe is an independent resource — copied content, no live reference)
- Handle naming conflicts when creating recipes from stock recipes (`recipe_name_conflict`;
  the `ibm-` prefix is reserved — customer recipes cannot use it)
- Ensure stock recipe metadata is preserved in created recipe

---

#### 2.3.6.2 Epic 2: Create a Custom Recipe

##### 2.3.6.2.1 User Story 2.1: Create Custom Recipe with Base Image and GPU Driver Components

**Title**: Create Custom Recipe with Stock Base Image and GPU Driver Selection

**As a** Platform Engineer or DevOps Engineer supporting AI/ML workloads,
**I need to** select a VPC stock base image and optionally add GPU drivers (Nvidia, AMD, Intel
Gaudi) to create a custom image recipe,
**So that** I can build standardized, GPU-enabled images from a trusted foundation for my
infrastructure deployments.

**Acceptance Criteria**:

1. User can view the full VPC stock image catalog as base image candidates (any stock image
   is eligible; deprecated stock images excluded or flagged)
2. Stock images display relevant metadata (OS type, version, architecture)
3. User can select a single stock image as the base for their custom image recipe
4. There is no separate base-image OS validation step — compatibility is determined entirely
   by whether the selected components' `allowed_use.image` expressions match the chosen base
   image (criterion 6)
5. User can optionally select from available GPU driver options (Nvidia, AMD, Intel Gaudi)
6. System validates GPU driver compatibility with selected base image
7. GPU drivers are automatically installed and configured during build process when selected
8. Resulting recipe metadata indicates base image and which GPU drivers are included
9. System provides clear error messages if incompatible selections are made

**Story Category**: Functional

**Considerations**:

- Handle scenarios where the base image is deprecated or deleted from the VPC stock image
  list while a recipe references it or a build is in progress
- Ensure proper IAM permissions are checked before allowing image selection
- Handle version conflicts between OS and GPU driver versions
- Handle network failures during driver download

---

#### 2.3.6.3 Epic 3: Automated Image Build Pipeline

##### 2.3.6.3.1 User Story 3.1: Automate Image Creation Pipeline

**Title**: Automate End-to-End Image Build Process

**As an** Infrastructure Architect,
**I need to** automate the complete image build pipeline from base selection to validation,
**So that** I can eliminate manual errors and ensure consistent, repeatable image creation.

**Acceptance Criteria**:

1. User can initiate an automated image build pipeline through the IBM Cloud portal
2. Pipeline automatically progresses through stages: selection, configuration, build, validation
3. Pipeline includes automated validation checks before proceeding to next stage
4. Failed builds provide detailed logs and error information (image job `status_reasons` names
   the failed stage and component)
5. Successful builds automatically generate versioned image artifacts
6. Pipeline execution time is tracked and reported (`image_job.started_at` / `completed_at`)
7. Build process includes verification that components (GPU drivers) are properly installed
   (package-manager-level verification — see 2.8.2.9 for the scope and limits of this in Phase 1)
8. Build fails gracefully with clear error messages if component installation fails

**Story Category**: Functional

**Considerations**:

- Handle pipeline interruptions due to infrastructure failures
- Implement retry logic for transient failures
- Ensure pipeline state is preserved across system restarts
- Handle concurrent build requests from same user
- Validate resource quotas before starting build
- Implement timeout mechanisms for long-running builds
- Deleting the image mid-build cancels the build with guaranteed cleanup

---

#### 2.3.6.4 Epic 4: Version Control and Traceability

##### 2.3.6.4.1 User Story 4.1: Automated Recipe Versioning

**Title**: Automatic Version Management for Recipes

**As a** Platform Engineer,
**I need to** have my recipes automatically versioned on every content change,
**So that** I can track changes, roll back if needed, and maintain clear history of recipe
evolution.

**Acceptance Criteria**:

1. System automatically assigns a content-hash version to each distinct recipe content state
2. The recipe `id` is immutable; the `(id, version)` pair permanently identifies exact
   content (name is mutable identity metadata)
3. User can view version history for each recipe (`GET /v1/image_recipes/{id}/history`), with
   every version — current and superseded — addressable and `superseded_at` timestamps
   establishing lifetime intervals
4. Each historical version is a complete frozen snapshot (comparison = diff of snapshots)
5. User can roll back by restoring a historical version's content (a rollback mints a version)
6. Version conflicts are impossible by construction (identical content vs. the current version
   → `409` duplicate; distinct content → distinct hash)
7. User can pin a build to a specific historical version (the API
   spec) — the resulting build uses exactly that version's frozen content

**Story Category**: Functional

**Considerations**:

- Ensure version history is preserved when recipes are deprecated
- History remains readable in `deprecated` and `obsolete` states (reads never gate on status;
  only builds do)
- History is removed only when the recipe itself is deleted

---

##### 2.3.6.4.2 User Story 4.2: Build Audit Trail and Traceability

**Title**: Comprehensive Audit Trail for Recipe Actions and Image Builds

**As a** Compliance Auditor,
**I need to** access complete audit trails showing CRUD operations on recipes and images,
**So that** I can verify compliance, troubleshoot issues, and maintain accountability.

**Acceptance Criteria**:

1. Activity Tracker (AT) audit logs include user identity, actions taken, and results
2. AT logs capture the resource (recipe, image) and actions taken
3. Image metadata includes the `source_image_recipe` `(id, version)` reference and an inventory
   for all installed packages
4. Audit logs are immutable and tamper-evident
5. User can search and filter audit logs
6. Logs are retained according to compliance requirements
7. System provides audit log export functionality

**Story Category**: Non-functional (Compliance)

**Considerations**:

- Ensure logs don't contain sensitive information
- Validate log integrity periodically

---

#### 2.3.6.5 Epic 5: User Experience and Accessibility

##### 2.3.6.5.1 User Story 5.1: Intuitive Build Configuration Interface

**Title**: User-Friendly Image Build Configuration

**As a** Platform Engineer new to IBM Cloud,
**I need to** easily configure and initiate image builds through an intuitive interface,
**So that** I can create custom images without extensive training or documentation.

**Acceptance Criteria**:

1. Build configuration interface follows IBM Cloud design standards
2. Configuration options are clearly labeled and explained
3. Interface provides contextual help and tooltips
4. User can save and reuse build configurations
5. Interface validates inputs in real-time
6. User can preview configuration before starting build
7. Interface is accessible and follows WCAG guidelines
8. Configuration wizard guides users through complex setups

**Story Category**: Functional

**Considerations**:

- Handle browser compatibility issues
- Ensure interface works on different screen sizes
- Validate accessibility for users with disabilities
- Handle slow network connections gracefully
- Provide keyboard navigation support
- Handle session timeouts during configuration

---

#### 2.3.6.6 User Story Summary

##### 2.3.6.6.1 Hill 1 (VPCIMG-277)

**Description**: Users will be able to select stock recipes or a base stock image and the
necessary GPU drivers available from the list of Nvidia, AMD, and Intel Gaudi GPU drivers.

**Size**: Large

**Dependencies**:

- VPC Stock Images catalog integration
- Component artifact repository (Artifactory) — exact contents (packages vs. install assets)
  pending legal review of GPU software redistribution (OI-DIST)
- HashiCorp Packer and packer-plugin-ibmcloud
- HCL generation (`hclwrite` Go library)
- IBM Cloud Continuous Delivery Tekton pipelines
- IBM Cloud Object Storage (template staging, build artifacts, inventory)
- IAM integration for access control (S2S authorization, new IAM actions)

---

##### 2.3.6.6.2 Target Customer Segments

1. **Regulated Enterprises** - Financial institutions, insurance, healthcare requiring compliance
2. **Enterprise IT** - Large organizations needing standardized infrastructure
3. **DevOps & SRE Teams** - Teams requiring rapid iteration and deployment
4. **Managed Service Providers** - MSPs supporting multiple clients with varied needs

---

##### 2.3.6.6.3 Success Metrics

- **User Adoption Rate**: Track number of active users and images built
- **Revenue**: Monitor consumption-based revenue from underlying IaaS services (image/storage)
- **NPS (Net Promoter Score)**: Measure customer satisfaction
- **Build Success Rate**: Percentage of successful builds vs. failures
- **Time to Deploy**: Average time from recipe creation to image build

---

##### 2.3.6.6.4 Competitive Positioning

Image Builder aims to achieve market parity with:

- **AWS EC2 Image Builder** - Mature console experience, deep AWS integration
- **Azure VM Image Builder** - Strong Azure-native integration, CLI/IaC focused

**Key Differentiators**:

- No cost for using the service (only underlying IaaS consumption)
- IBM-curated **stock recipes** — no competitor ships first-class, buildable, lifecycle-managed
  recipes
- Component lifecycle semantics (explicit deprecation/obsolescence dates) — ahead of AWS's
  version pile-up, on par with GCP's family/deprecation model
- Integration with IBM Cloud PaaS services
- Enterprise-grade compliance and security features

**Note**: see 2.3.2.2 for the recipe-reuse comparison with AWS and Azure.

## 2.4 Out-of-Scope

This Phase 1 implementation focuses on establishing the foundational Image Builder
capability with a limited set of features. This section lists the features and capabilities that
are explicitly out-of-scope for this phase and may be implemented in future phases based on
business need and requirements.

Future phases will address:

- **Licensed and Billable software components** — the motivating use case is Windows
  Microsoft SQL Server
  images for customers migrating from Classic infrastructure to VPC (replicating the Classic
  offering); requires the licensing and billing integration deferred from Phase 1
- **Customer components** — a declarative approach first, where customers author a schema
  (`package_details`) that follows
  the same generator path as IBM's own components, with a frozen `custom_script` install
  method as the
  bounded exception for software with no package form (script taken at publish, content-addressed,
  never live-fetched). A public write endpoint and the `install_method` extension
  point are designed for this; `package_details`/`install_method` promote to the public schema
  at that point
- **Image scan stage and post-capture boot test** — in-pipeline security scanning of the built
  image, and booting an ephemeral instance from it to re-run component validation before
  availability (the genuine capability gap vs. AWS's `test` phase). A platform-independent
  alternative exists for the scanning half regardless of
  timing — feeding the inventory to SCC Workload Protection (Sysdig) as one additional
  call, alongside `cosign`-style image signing — see the open item tracked as OI-SCAN
- **Automated re-build notifications** — "your image contains a deprecated component"
  (with the Phase 2 vulnerability scanning / image state indicator work)
- **IAM service-group authorization** — reduce the 4-policy S2S onboarding friction (Aha idea
  filed)
- **Image templates** — a reusable, referenced-by-ID image-creation resource (the
  instance-template pattern; effectively an image service feature with Image Builder
  interoperability). This is deferred consciously, since Phase 1 derives image
  identity instead. See the alternatives-considered record in 5.2 for the full rationale,
  the Phase 1 derived-identity behaviors, and the design constraints recorded for the future
  implementation
- Customer-created stock recipes or cross-account shared custom recipes
- Third-party marketplace (software components with billing requirements)
- Advanced vulnerability scanning and remediation
- **Automated image refresh capabilities** — this covers two distinct things, both out of
  scope for Phase 1: (1) auto-triggering a rebuild (e.g. on a schedule, or when a component
  updates), and (2) auto-rebasing a recipe to the latest stock image within an OS family (e.g.
  automatically picking up a new Ubuntu 22.04 point release). Today, a manual rebuild uses the
  exact base image ID stored in the recipe at build time — getting the latest stock image
  refresh requires the customer to `PATCH` the recipe's `base_image` first (minting a new
  version), then trigger a build against that version. There is no automatic rebase and no
  scheduled/triggered rebuild in Phase 1.
- **A dedicated `build_configuration` CRUD resource** does not exist in Phase 1 — this is a
  rejected
  design alternative, not a deferred capability — see the decision record in 5.2.1 for the
  full rationale.

- **Phase 2 (tentatively 2027) - Out of Scope:**

  This list is subject to change but the current phase 2 project includes:
  - Restraints at the appropriate integration points: the desired approach is to obtain some
    degree of service adoption in order to understand customer consumption and usage before adding
    constraints to the new recipe and component architecture elements
  - Licensed customization components with associated costs
  - Automated vulnerability scanning of images created by the image builder service
  - Image state indicators for out-of-date or vulnerable images
  - Regular scanning of built images to detect vulnerabilities or end-of-life/end-of-service
    software
  - Integration with industry scanning tools such as Nessus for vulnerability detection
  - Explicit build retry (Phase 1 requires creating a new job; see 5.2 for the rejected
    artifact-recovery alternative)
  - Recurring/scheduled rebuild support, addressing the recipe-reuse gap noted in 2.3.2.2

- **Phase 3 (Future) - Out of Scope:**

  This list is subject to change but the current phase 3 project includes:
  - Customer-created customization components (initially experimental)
  - Third-party marketplace integration for custom components
  - Cost attribution and billing for third-party components
  - Support for third-party vendors to charge for their custom components
  - Expanded component ecosystem beyond IBM-curated offerings

- **Additional Out-of-Scope Items:**
  - Stock recipes provided by 3rd parties
  - Automatic software component updates within existing images (image builder does not update
    software components; building a new image with updated components is required)
  - Support for BYOK (Bring Your Own Key) or KYOK (Keep Your Own Key) encrypted images in
    vulnerability scanning workflows
  - Image builder capabilities for IBM Cloud Classic infrastructure (this service is VPC-specific)

The rationale for this phased approach is to establish service maturity, operational readiness,
and governance processes before expanding the scope of customization options and introducing
licensed or third-party components. Phase 1 will validate the service architecture, pipeline
infrastructure, and governance model with a focused set of use cases before scaling to more
complex scenarios.

### 2.4.1 Out-of-Scope Use Case

#### 2.4.1.1 Database Software Images

Provide pre-configured images with database software:

- **Microsoft SQL Server**: Replicate functionality available in IBM Cloud Classic Infrastructure.
  Time-permitting this use case could be accelerated to phase 1 or a phase 1 fast follow on.
- **Other Database Platforms**: Support for popular database systems as prioritized by business
  needs

**Business Value**: Simplifies database deployment, ensures proper configuration, provides
migration path from Classic Infrastructure.

#### 2.4.1.2 IBM Cloud Platform Services Images

Enable IBM Cloud platform services to build standardized images containing:

- Platform agents and monitoring tools
- Security and compliance software
- Service-specific utilities and configurations

**Business Value**: Provides repeatable mechanism for platform services to create and maintain
their required images, improves consistency across platform deployments.

#### 2.4.1.3 Popular Application Images

Support images for widely-used applications:

- **cPanel**: Web hosting control panel (migration from Classic Infrastructure)
- **Other Popular Applications**: As prioritized by customer demand and business value

**Business Value**: Attracts customers migrating from other platforms, reduces barrier to entry
for common workloads.

#### 2.4.1.4 Security-Hardened Images

Provide repeatable mechanism for security customizations:

- Apply security baselines and hardening configurations
- Include security monitoring and compliance tools
- Implement organization-specific security policies

**Business Value**: Ensures consistent security posture, reduces security configuration errors,
accelerates compliance certification.

## 2.5 Dependencies, Dependents, and Incompatibilities

### 2.5.1 Other IaaS Projects Under Development

There are no other IaaS projects under development upon which this proposal is dependent and there
are no projects dependent on this proposal.

Although these projects are not directly dependent on this proposal, these services or projects may
benefit from this proposal and thus are likely interested projects:

- Existing and new GPU accelerator projects
- IaaS projects designed to promote IaaS classic to VPC migration (e.g. Sandbox)
- PaaS projects with a need to create custom images for their service VSIs (e.g. ROKS/IKS)

### 2.5.2 Other IBM Projects Under Development

None.

### 2.5.3 Third-party

This proposal introduces a dependency on **HashiCorp Packer**, an open-source tool for creating
machine images from source configuration. Packer is required to execute the generated build
templates that transform recipes into VPC custom images.

#### HashiCorp Packer Dependency

**Purpose**: Packer is used to automate the creation of VPC custom images by:

- Executing the service-generated HCL template (install and validate provisioner blocks derived
  from component schema)
- Creating custom images directly via "image from volume" (see 3.1.1.3.3 for the full
  mechanism)
- Validating image integrity before deployment

**Deployment Model**:

- The Packer binary and the IBM Cloud plugin are **pre-baked into an IBM-built task container
  image** (UBI base, pinned versions, checksum-verified at image build, published to ICR with
  Vulnerability Advisor scanning and CVE-triggered rebuilds) — the pipeline does not pull the
  upstream hashicorp image and performs no network fetch of build logic at run time
- Packer build templates are generated per build in HCL by
  the control plane, integrity-hashed (`hcl_sha256`), staged per-build, and verified by Task 1
  before use
- No runtime dependency: Packer is only required during image creation, not for image usage

**License and Legal Considerations**:

- Packer is licensed under the Mozilla Public License 2.0 (MPL 2.0)
- MPL 2.0 is a weak copyleft license that allows commercial use
- No source code disclosure requirements for IBM Cloud VPC code
- No licensing fees required
- Legal approval status: **TBD - pending IBM Legal review**

**Version and Support**:

- Minimum supported version: **Packer 1.10.x or later**, version pinned in the IBM-built task
  image Dockerfile — build behavior changes only when IBM changes it
- IBM will maintain compatibility with current stable Packer releases
- Packer updates are managed through standard dependency management processes

**Security Considerations**:

- Packer binaries obtained from official HashiCorp releases, verified against HashiCorp's
  published checksums at task-image build time
- Packer execution occurs in controlled build environments with appropriate access controls
- Generated templates are integrity-verified (`hcl_sha256`) before execution; the hash is
  recorded in the image's build metadata
- No customer credentials or secrets are embedded in Packer templates (credentials injected at
  run time from Tekton environment properties)

**Operational Impact**:

- Build environment must have network access to download software component packages from
  Artifactory (binaries only — no scripts)
- Build environment requires IBM Cloud API access for packer-plugin-ibmcloud operations
- Packer build failures are isolated to individual image creation operations
- Failed builds do not impact existing images or running instances
- Build logs are captured for troubleshooting and audit purposes

**Alternatives Considered**:

- Custom image builder implementation: rejected due to development complexity and
  maintenance burden
- Cloud-init only approach: insufficient for complex software installation and configuration
- Packer selected for its maturity, community support, VPC and compatibility

### 2.5.4 External and/or Circular

None.

### 2.5.5 Brittle or Unusual

None.

### 2.5.6 Incompatibilities

None.

### 2.5.7 Databases

This project does not depend on deploying new database instances or introduce new dependencies on
database instances.

## 2.6 Integration, Deployment, and Operations Considerations

### 2.6.1 Integration Considerations

**VPC Service Integration**:

The image builder capability integrates with existing VPC infrastructure and services:

- **RIAS Regional Control Plane**: Recipes (stock and custom) and components are exposed
  through the existing RIAS control plane (customer-facing API, IAM, AT — standard VPC
  patterns); the underlying storage uses Kubernetes CRDs rather than RIAS's direct etcd
  pattern, since the Image Builder controllers already run on Kubernetes and CRDs let the
  same reconciliation framework manage both controller logic and resource lifecycle
- **Image Service**: Image Builder is delivered as a **feature extension of the image service**
  (two new top-level resources under it). Built images are registered with the VPC image service
  and follow existing image lifecycle management
- **IAM Integration**: Authentication and authorization use standard VPC IAM patterns with resource
  groups and access policies
- **Activity Tracker**: All API operations emit Activity Tracker events following VPC event
  standards
- **Metering and Billing**: Underlying build resources are metered through existing VPC billing
  in the customer's account

**Build Pipeline Integration**:

- **Packer Integration**: HashiCorp Packer is integrated as the core image build engine via an
  IBM-built task container image
- **Component Catalog**: Component records provide the declarative build knowledge
  (`software`, `dependencies`) and Artifactory
  download locations for IBM-approved software components
- **Build Orchestration**: Image jobs are dispatched to the IBM Continuous Delivery Tekton
  service via its REST API (`CreatePipelineRun`); a Status Watch controller monitors pipeline
  runs and pushes image job status to the image service
- **Image Validation**: Builds run package-manager integrity checks
  in-build (package-level only in Phase 1 — see 3.1.1.3.6a). Per-component validation commands
  are not part of Phase 1, and post-capture boot test does not run — see 2.4 for both

**Storage Integration**:

- **Object Storage**: Build artifacts, generated templates, logs, and intermediate files are
  stored in IBM Cloud Object Storage with appropriate lifecycle policies
- **Image Storage**: Completed images leverage existing VPC image storage infrastructure
- **Recipe Metadata**: Recipe and component metadata (including recipe version history) are
  modeled as Kubernetes Custom Resources (CRDs) in the regional cluster, reconciled by a
  dedicated controller — consistent with how the Image Builder controllers already run on
  Kubernetes

**Security Integration**:

- **Secrets Management**: Service credentials and pipeline secrets are managed through the
  VPC Vault/Secrets Manager model (controller API keys, Tekton pipeline key, TLS certificates)
- **Encryption**: All data at rest is encrypted using IBM Cloud encryption services; all data in
  transit uses TLS 1.2+

### 2.6.2 Deployment Considerations

**Service Components**:

The image builder capability consists of the following deployable components:

1. **Image Builder API Controller** (Regional):
   - Customer-facing REST endpoints for recipes, components, and image jobs
   - Pre-dispatch validation, HCL generation, pipeline dispatch (CD REST)
   - Stateless microservice in the RIAS control plane; rolling updates with zero downtime

2. **Status Watch Controller** (Regional):
   - Polls Tekton PipelineRun status via the CD REST API
   - Writes the image job record (coarse status; `status_reasons` only on `failed` — internal
     stage tracking is not surfaced to customers, see 3.1.1.9.3a) over the internal surface;
     tracks the image's own status directly, since there is no separate import phase
   - Stateless; rolling updates with zero downtime

3. **Component & Recipe Catalog** (Regional storage):
   - Component and recipe resources (with content-hash version history) modeled as Kubernetes
     CRDs (`ImageComponent`, `ImageRecipe` — named to avoid ambiguity for CORE/SRE
     on the shared cluster), reconciled by the Catalog Controller
   - IBM catalog entries (components, stock recipes) are **owned by a dedicated IBM catalog
     account** (staging/production pair, documented ownership succession) and onboarded through
     internal/operator tooling as the owning account —
     no public write endpoint exists for catalog authoring in Phase 1;
     publicization (visibility flip to `public` + `ibm-` prefix) is a privileged operator action

4. **Tekton Pipelines** (Regional, in a dedicated IBM pay-as-you-go account):
   - One pipeline per supported region in production; a shared Tekton test environment serves
     staging/integration; **off by default in development MZones** (enabled per-MZone via a
     global, following the Share precedent)
   - Tasks run in Tekton-managed pods using the IBM-built Packer task image

**Deployment Pipeline**:

- **CI/CD Integration**: Service components use existing VPC deployment pipelines with standard
  stages (build, test, stage, production)
- **Canary Deployments**: New versions are deployed using canary strategy (5% → 25% → 100%) with
  automated rollback on errors
- **Feature Flags**: New features are controlled by feature flags to enable gradual rollout and
  quick rollback if needed
- **Packer Version Management**: Packer and plugin versions are pinned in the IBM-built task
  image and tested in staging before production deployment; CVE-triggered rebuilds are automated
- **Catalog publishing**: A dedicated publish tool (schema-validated payloads, CI approval
  gates) is the sole holder of the catalog account credential — no human use

**Infrastructure Requirements**:

- **Compute**: Tekton-managed pods (no dedicated IBM zone-resident build nodes); ephemeral build
  VSIs in customer VPCs
- **Storage**: Object storage buckets for build artifacts and logs
- **Network**: Outbound connectivity from build VSIs to Artifactory (public gateway on the
  customer's chosen subnet)
- **Monitoring**: Integration with existing osprey and IBM Cloud Logs (ICL) infrastructure

**No Customer Downtime**:

Deploying this project will **not** cause any customer downtime. The capability is entirely new
and does not modify existing VPC functionality. Customers can continue using existing custom
image workflows without interruption.

### 2.6.3 Operations Considerations

**Operational Responsibilities**:

The VPC SRE team will be responsible for:

- Monitoring service health and performance metrics
- Responding to alerts and incidents
- Managing capacity and scaling
- Troubleshooting build failures and customer issues

**New Operational Risks**:

1. **Build Environment Capacity**: High concurrent build demand could exhaust Tekton capacity
   - **Mitigation**: Build queued, capacity monitoring and alerting for staging/integration load

2. **Component Availability**: Component package download failures could cause build failures
   - **Mitigation**: Retry logic, Artifactory availability monitoring

3. **Packer Dependency**: Packer bugs or breaking changes could impact builds
   - **Mitigation**: Version pinning in the IBM-built task image, thorough testing, rollback
     procedures

4. **Security Vulnerabilities**: Discovered vulnerabilities in components require rapid response
   - **Mitigation**: Component lifecycle dates enable immediate retirement — in an emergency,
     set
     `deprecation_at` and `obsolescence_at` to now; notification system; rapid publish of
     replacement versions

**Recurring Operational Tasks**:

- **Component Lifecycle Management**: Regular review and update of component versions (monthly)
- **Stock Recipe Updates**: Periodic updates to stock recipes with new OS versions and
  components (quarterly)
- **Capacity Planning**: Regular review of build pipeline capacity (monthly)
- **Performance Tuning**: Optimization of build times and resource utilization (quarterly)
- **Security Patching**: Task-image CVE rebuilds (automated) and service component patching (as
  needed)

**Monitoring and Alerting**:

Key metrics to monitor:

- Build success/failure rates
- Build duration (p50, p95, p99)
- Build queue depth and wait times
- Component package download success rates
- Validation pass/fail rates (per component)
- Service API latency and error rates
- Pipeline resource utilization

Critical alerts:

- Build failure rate exceeds threshold (>5%)
- Build queue depth exceeds capacity
- Component download failures
- Service API errors or high latency
- Orphaned build resources detected (cleanup task failures)
- Security vulnerability detected in components

**Support and Escalation**:

- **Support Level**: The image builder capability will enter production at the same support level
  as other VPC services (24x7 support)
- **Incident and Case Routing**: Incidents and customer cases will be routed to CORE IPOPS as
  the primary support team, following the existing CORE IPOPS → Compute SRE support model, with
  escalation to the image builder development team for complex issues
- **Documentation**: Comprehensive operational runbooks will be provided covering common failure
  scenarios and troubleshooting steps
- **Training**: VPC SRE team will receive training on image builder architecture, common issues, and
  troubleshooting procedures before GA

SRE support-model alignment will be completed through the service framework process. Note that
compute SRE does not currently support the image service; their role in the escalation tier of
the CORE IPOPS → Compute SRE model for this capability is to be agreed with them as part of
that process (tracked as OI-SRE).

**Service-Level Objectives**: The SLO referenced here is the service-specific SLO owned by the
Image Builder team, not a team-internal performance measurement. Per-severity (Sev-1/2/3/4)
response and resolution targets will be defined following IBM Cloud's standard resiliency/SLO
framework (<https://cloud.ibm.com/docs/resiliency?topic=resiliency-slo>) and finalized as part
of VPC SRE sign-off ahead of GA. No severity-specific numbers are committed in this proposal;
placeholder values are deliberately omitted here so this section is not read as authoritative
ahead of that sign-off.

**SRE Team Sign-off**:

- **VPC SRE Team**: [TBD — pending sign-off, including confirmation of Sev-1/2/3/4 SLO targets
  per the framework above] — Responsible for service monitoring, incident response, and
  operational maintenance

**Upgrade Strategy**:

All service components can be upgraded non-disruptively:

- **Control Plane Services**: Rolling updates with health checks ensure zero downtime
- **Build Pipelines**: New builds use updated task images; in-progress builds complete on
  existing images
- **Packer Version Updates**: New task image tested in staging, gradual rollout
- **Component Catalog Updates**: Versioned components allow gradual migration

**Disaster Recovery**:

- **Recipe Metadata**: Backed up via the regional Kubernetes cluster's standard backup/DR
  tooling (Velero snapshots) — the same mechanism used for other platform CRDs
- **Build Artifacts**: Stored in COS with versioning and lifecycle policies
- **Recovery Time Objective (RTO)**: 4 hours for full service restoration
- **Recovery Point Objective (RPO)**: 1 hour for metadata; no data loss for completed images

### 2.6.4 Cryptography and Export Considerations

This project does not introduce high-strength cryptography beyond what is already used in VPC
services. All encryption uses standard IBM Cloud encryption services (AES-256) and TLS 1.2+ for data
in transit. There are no export restrictions or legal considerations that would prevent deployment
in any IBM Cloud region.

### 2.6.5 Customer Communication

**Feature Announcement**:

- Customers will be notified of the new image builder capability through standard IBM Cloud
  channels (blog posts, documentation, release notes)
- Beta program (WDC, DAL, FRA) will allow early adopters to test the service and provide feedback
- GA announcement will include comprehensive documentation, tutorials, and examples

**Service Limitations**:

- Phase 1 limitations (IBM-only stock recipes and components) will be clearly documented
- Build time expectations and quotas will be communicated in service documentation
- Component lifecycle and update policies (deprecation/obsolescence dates) will be published and
  maintained

**Support Channels**:

- Standard VPC support channels (tickets, chat, phone)
- Community forum for questions and discussions
- Documentation and tutorials on IBM Cloud Docs
- API reference documentation

## 2.7 Big Rules

### 2.7.1 Consistent Customer Experience

The image builder supports custom image builds for both **virtual server instance (VSI)** and **bare
metal** compute offerings — via the plugin's "image from volume" mechanism. While both use
cases are supported by the architecture, the VSI use case
is the priority for this proposal and will be delivered first. Bare metal image support will follow
in subsequent phases.

### 2.7.2 VPC Big Rules

This project conforms to all [VPC Big Rules](ng-big-rules.md):

- [x] **All components must be able to be upgraded live (without service degradation/outage) and
      without manual intervention** - Yes.

- [x] **All components must present usage metrics** - Yes.

- [x] **All components must use the VPC/NG monitoring framework** - Yes.

- [x] **All components must use existing VPC/NG control planes -- no new control planes** - Yes —
      uses the existing RIAS regional control plane.

- [x] **All components must be controlled via well-defined APIs** - Yes.

- [x] **All components must provide functional parity across all modalities (UI, CLI, API,
      Terraform)** - Yes.

- [x] **All components must provide functional parity across all platform architectures (x86 and
      Z)** - Yes.

- [x] **All customer-facing GA features must interoperate with all other customer-facing GA
      features** - Yes.

- [x] **All customer-facing UIs and CLIs must use the customer-facing VPC API to operate on VPC/NG
      resources** - Yes.

- [x] **All customer-facing APIs must adhere to the API handbook guidelines** - Yes.

- [x] **All regional APIs must have the same API contract in all regions** - Yes.

- [x] **All components must perform authorization using the VPC/NG IAM services** - Yes.

- [x] **All components must be internally stateless and recreate their runtime state from external
      sources** - Yes.

- [x] **Any caching or shadowing of state must never result in an API consumer observing
      inconsistent or incorrect behavior** - Yes.

- [x] **All hardware components must use the VPC/NG diagnosability framework and implement the
      on-demand diagnostic API** - Yes — note this is primarily a control-plane service.

- [x] **All services must present an API for get/set/watch of state** - Yes — standard CRUD plus
      polling on the image/jobs resource and image status.

- [x] **All connections between components must be mutually authenticated, and encrypted using
      AES128 or higher** - Yes.

- [x] **All components must integrate into the existing operations tooling for deployment and
      maintenance** - Yes.

### 2.7.3 Cloud Security Baseline Requirements

This project meets all Cloud Security Baseline Requirements:

- [x] **[Authentication & Access Control](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-authentication-access-control)**
      - Yes.

- [x] **[Compute Isolation](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-compute-isolation)**
      - Yes.

- [x] **[Encryption](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-data-encryption)**
      - Yes, including encryption at rest with secure deletion, encryption in transit (TLS 1.2+,
      mutual TLS between services), and FIPS 140-2 compliant algorithms (AES-256 at rest).

- [x] **[Injection prevention](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-injection-prevention)**
      - Yes — generated Packer templates are derived from validated schema fields only (no
      arbitrary customer script execution in Phase 1) and are validated (`packer validate`) plus
      integrity-hashed before execution.

- [x] **[Least privilege operation](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-least-privilege-operation)**
      - Yes — the customer's S2S grant is scoped to exactly four resource types.

- [x] **[Secret Management](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-secret-management)**
      - Yes.

- [x] **[Secure Development Practices](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-711-software-development-lifecycle)**
      - Yes.

- [x] **[Update and Secure Configuration](https://test.cloud.ibm.com/docs-internal/ibmcloud-security?topic=ibmcloud-security-update-secure-configuration)**
      - Yes.

## 2.8 Assumptions, Risks, and Constraints

### 2.8.1 Assumptions

#### 2.8.1.1 Stock Recipe Management

- IBM will maintain and update stock recipes centrally, ensuring consistency across all regions
- Stock recipes carry the same content-hash versioning and history as customer recipes, allowing
  reproducible pinned builds of specific stock recipe versions
- Stock recipe updates will not automatically affect existing customer recipes or images created
  from those recipes (customer recipes are independent resources with copied content)
- Users will be notified of stock recipe updates through standard IBM Cloud notification channels

#### 2.8.1.2 Component Availability

- Phase 1 components (GPU drivers, accelerators) will be sourced exclusively from IBM-approved
  vendors
- Component versions are distinct catalog entries with lifecycle dates (`deprecation_at`,
  `obsolescence_at`) from which status (`pending`/`available`/`deprecated`/`obsolete`/`deleting`) is derived
- Component binaries will be hosted in IBM-controlled repositories (Artifactory — binaries only,
  no scripts) with appropriate access controls
- Component updates will follow standard VPC change management processes

#### 2.8.1.3 Packer Build Environment

- Build VSIs will have reliable network connectivity to download software component packages
  from Artifactory (public gateway on the customer's configured subnet)
- Packer execution occurs per-build in the customer's own VPC with ephemeral resources —
  no cross-build contamination
- Build times will vary based on image size and component complexity (estimated 15-45 minutes per
  build; for Phase 1 precompiled_package installs on the Flex build profile, best case is
  closer to ~25 minutes — see 3.2.2)
- Failed builds will **not** be automatically retried — re-running is a customer action (a
  new job); transient-failure resilience lives inside the pipeline (a bounded upload window
  within the build timeout, and the plugin's own internal polling retries)

#### 2.8.1.4 Image Distribution

- Images built from recipes will be distributed to availability zones using existing VPC image
  replication mechanisms
- Every build, including one from a stock recipe, produces a new image on demand — there is no
  pre-built stock image cache (see 2.2 for the "every build produces a new image" decision)
- Image storage will leverage existing VPC image storage infrastructure and quotas

#### 2.8.1.5 Integration with Existing VPC Services

- Images created through the image builder will be fully compatible with existing VPC instance
  provisioning
- Recipe reference and inventory metadata will be stored alongside image metadata for
  governance and tracking
- IAM policies for image builder resources will follow existing VPC resource authorization patterns

### 2.8.2 Risks

#### 2.8.2.1 Third-Party Dependency Risk (HashiCorp Packer)

- **Risk**: Packer updates or breaking changes could impact image build reliability
- **Mitigation**: Pin Packer and plugin versions in the IBM-built task image; test new versions
  in staging before production deployment; maintain compatibility matrix
- **Impact**: Medium: could cause build failures but does not affect existing images or running
  instances

#### 2.8.2.2 Component Lifecycle Management Risk

- **Risk**: Component versions may reach end-of-life or have security vulnerabilities discovered
- **Mitigation**: Explicit lifecycle dates with derived statuses; recipes referencing obsoleted
  components surface `degraded` status with actionable warnings; emergency retirement via
  immediate dates; notifications to users with affected recipes/images
- **Impact**: Medium: requires ongoing maintenance and user communication

#### 2.8.2.3 Build Environment Scalability Risk

- **Risk**: High concurrent build demand could exhaust Tekton pipeline capacity
- **Mitigation**: Implement build queued with priority levels; per-region pipelines; capacity
  confirmation with the CD team
- **Impact**: Low: affects build latency but not service availability

#### 2.8.2.4 Hybrid Execution Boundary Risk

- **Risk**: The hybrid model's boundary between IBM control plane (above the cloud) and
  customer VPC (zone-resident) creates a dependency on the customer's VPC configuration.
  Customer NACL rules, subnet configuration, or IAM S2S grant revocation can silently
  break builds after the build VSI is already provisioned.
- **Mitigation**: Pre-dispatch validation catches the most common configuration
  issues before any VSI is provisioned. The cleanup task runs on all exit paths to prevent
  orphaned resources. Monitoring alerts on SSH connectivity failures and iam_s2s_missing rate.
- **Impact**: Medium — customer-side failures are outside IBM's direct control but are
  detectable via build log analysis and monitoring.

#### 2.8.2.5 Stock Recipe Content Update Risk

- **Risk**: A stock recipe's content changing (e.g. a new component version) has no effect on
  images already built from an earlier version, since every build produces its own independent
  image; there is no shared cache to go stale
- **Mitigation**: Content-hash versioning makes every stock recipe version permanently
  addressable, so a customer can always rebuild against the exact version they used before, or
  pin a build to a specific historical version
- **Impact**: Low: each build is independent, so a content update cannot retroactively affect an
  existing image

#### 2.8.2.6 Cross-Region Consistency Risk

- **Risk**: Stock recipes or components could become inconsistent across regions during updates
- **Mitigation**: Automated publishing (single publish tool, schema-validated payloads) across
  regions; validate consistency before publicizing
- **Impact**: Medium: could cause user confusion or build failures

#### 2.8.2.7 Security Vulnerability in Components Risk

- **Risk**: Discovered vulnerabilities in included components could affect many customer images
- **Mitigation**: Implement continuous vulnerability scanning; provide automated notifications;
  enable rapid component retirement (immediate lifecycle dates) and replacement publication;
  maintain security advisory process through psirt
- **Impact**: High: could require coordinated customer response for critical vulnerabilities

#### 2.8.2.8 Catalog Account Credential Risk

- **Risk**: The IBM catalog account's credential controls what GPU software IBM publishes to
  all customers — a supply-chain-relevant credential
- **Mitigation**: Key held in Secrets Manager with automated rotation; the publish tooling is
  the sole consumer (no human use); minimal policy scope (catalog resource types only);
  Activity Tracker alerting on any catalog write not originating from the tooling identity
- **Impact**: High if compromised; controls above reduce likelihood and detection time

#### 2.8.2.9 No Functional Validation or Hardware/Vendor Compatibility Enforcement Risk

Phase 1 ships without per-component
`validation_commands`, and without any hardware or GPU-manufacturer compatibility mechanism.
OS and architecture compatibility are covered separately, via the `allowed_use.image` CEL expression
on the component (see 3.1.1.3.6a) — the risk here is narrower than it once was, scoped to what
that expression doesn't cover.

- **Risk — silent hardware/GPU-manufacturer mismatches**: nothing at the API layer stops a
  customer from combining components in a recipe that are mutually incompatible with each
  other or with the profile the resulting image will eventually run on (e.g. an NVIDIA
  driver alongside an AMD toolkit, or a component built for one GPU model used against
  another). The failure only surfaces later — either the build itself fails in a way that
  doesn't clearly name the real cause, or the build succeeds and the resulting image later
  fails when a customer actually boots an instance from it, at which point the connection back
  to an incompatible recipe is much harder to trace.
- **Risk — no functional verification of installed software**: in-build validation is limited
  to package-manager integrity checks, authored entirely by the HCL generator itself rather
  than any user- or component-supplied validation logic (confirming files installed and
  checksums match; the specific commands are being finalized before commitment review — see
  OI-VALCMD). Nothing confirms the software actually works — a driver that installs
  cleanly but fails to load, or a toolkit binary that runs but reports a broken version, would
  both pass Phase 1's validation and produce an `available` image. This gap is only partially
  inherent to the platform — no GPU hardware exists on the build VSI, so a driver load can
  never be confirmed at build time regardless; a `custom_script`/`validation_commands`-style mechanism
  would still let package-manager-independent checks that don't require GPU hardware run (for
  example, confirming a toolkit binary executes and reports its version correctly), which
  Phase 1 currently has no way to express at all.
- **Risk — no vendor exclusivity enforcement**: with no `vendor` field and no
  hardware compatibility mechanism, nothing prevents a recipe from combining components that
  were never intended to coexist (e.g. an NVIDIA driver alongside an AMD toolkit). Note that
  OS and architecture mismatches specifically are no longer part of this gap — those are
  caught by the `allowed_use.image` expression.
- **Who is actually exposed**: stock recipes are unaffected in practice, since IBM curates and
  tests those combinations directly before publishing them. The exposure is concentrated
  entirely in **Path 3 — fully custom recipes assembled à la carte from the component
  catalog** (2.2). This is a real, not hypothetical, gap for that path, and worth being
  explicit that "most customers use stock recipes" is a mitigation by usage pattern, not a
  mitigation built into the system.
- **Mitigation for Phase 1**: OS/architecture mismatches are now caught at recipe validation
  time. For everything else, the mitigations remain reactive — clear, specific
  `failed_component_installation`/
  `failed_component_validation` status reasons naming the failing component (3.1.1.9.3a) so a
  customer can at least diagnose a failed build quickly, and documentation steering customers
  toward stock recipes for GPU workloads where hardware compatibility matters most.
- **Impact**: Medium for Phase 1 overall, given the stock-recipe usage pattern most customers
  are expected to follow, and now narrower than before given OS/architecture coverage — but
  still High specifically for any customer who does assemble a custom
  recipe from individual components mixing hardware vendors, since they get no compatibility
  or functional-validation
  safety net for that. This should be weighed directly against the time saved by not building
  this validation and compatibility tooling for Phase 1, not treated as a free simplification.

### 2.8.3 Constraints

#### 2.8.3.1 Phase 1 Scope Limitations

- Stock recipes are IBM-provided only; customer-created stock recipes are not supported in Phase 1
- Components are limited to IBM-approved software (GPU drivers, accelerators); no third-party
  marketplace integration
- No support for licensed software components requiring customer-provided licenses
- Build customization limited to component selection; no arbitrary script execution (component
  validation commands are IBM-authored in Phase 1)
- Maximum of 20 components per recipe (`recipe_component_limit_exceeded`), raisable later
  without API change

#### 2.8.3.2 Platform Architecture Constraints

- The plugin
  performs "image from volume" via `ImageBySourceVolume`, during Task 1's Capturing stage. See
  3.1.1.3.3 and 3.1.1.3.4 for the full mechanism.
- Build environments are Linux-based; no Phase 1 component targets Windows, so a Windows
  base image cannot practically be built from today, even though nothing at the API level
  blocks selecting one — see 3.1.1.6
- Images must pass in-build validation (package verification + component validation commands)
  before being made available; an image is never `available` before its build metadata is written

#### 2.8.3.3 Resource and Performance Constraints

- Build operations are asynchronous; no synchronous image creation API
- Maximum image size is constrained by the VPC platform image size limit. Phase 1 images
  consist of a VPC stock base image with IBM-curated GPU driver components added on top.
- Build timeout limits will be enforced (default 120 minutes, 30–480 minute customer-configurable
  range)
- Concurrent builds per account may be limited to prevent resource exhaustion

#### 2.8.3.4 Operational Constraints

- IBM catalog entries (components, stock recipes) cannot be deleted while customer recipes or
  in-flight builds reference them; deletion is gated behind obsolescence plus a cross-account
  safe-to-delete check (operator-side)
- Cancelling a build is `DELETE /v1/image/jobs/{id}`: the pipeline run is cancelled and the
  cleanup task always runs — cancellation never strands build resources. **Import-failure
  handling has been removed as a distinct case** — at the source-code level, the
  plugin's capture step blocks until the image reaches `available` before the task completes,
  so there's no "build succeeded, import failed independently afterward" scenario for a
  recovery model to address. The DELETE-vs-import race scenario is removed for the same
  reason — see 3.1.1.9.3a and the retired decision record in 5.2.1.
- Build logs retained for limited time period (30-90 days) for troubleshooting

#### 2.8.3.5 Integration Constraints

- Images created through image builder follow same lifecycle as manually created custom images
- No automatic image updates; users must rebuild images to incorporate component updates
- Build provenance metadata on an image is written once and immutable; it cannot be
  retroactively modified and persists after recipe deletion
- Cross-account image sharing follows existing VPC image sharing constraints

#### 2.8.3.6 Evolution and Maintenance Constraints

- Recipe schema changes must maintain backward compatibility
- Component schema must support future expansion (licenses, marketplace, customer
  components — see 2.4 for the planned extension points)
- Build pipeline must be extensible to support future phases (Windows, licensed software,
  frozen customer scripts)
- API versioning must accommodate future feature additions without breaking existing integrations

#### 2.8.3.7 Dependency on External Services

- Requires reliable access to stock OS image repositories (Ubuntu, Red Hat, etc.) from the image
  service
- Depends on component vendor repositories for driver downloads (when offering new components)
- Relies on HashiCorp Packer continued support and compatibility
- Assumes stable VPC infrastructure for image storage and distribution

# 3. Architecture, Interfaces, and Impact

## 3.1 Architecture and Interfaces

### 3.1.1 Architecture and Technical Design

VPC Image Builder is a new VPC capability that enables customers to create custom VM images
through declarative recipes. Customers define their own recipes specifying a base OS image and a
set of IBM-curated GPU driver components, then trigger builds that produce ready-to-use custom
images.

**Architectural model — hybrid above-the-cloud:** Builds execute in the
**customer's own VPC** using the customer's quota and billing. The system has three aspects:
(1) the regional Image Builder controllers running in RIAS, above the cloud; (2)
pipeline jobs running in an IBM-owned image builder service account, also above the cloud; and
(3) temporary build VSIs running in the customer's account (zone-resident) — the sole
zone-resident piece, IBM-managed via IAM S2S but not IBM-located, with zero IBM zone-resident
service nodes otherwise. This hybrid boundary
is the defining characteristic
of the architecture:

| Layer | Where it runs | Who manages it |
| --- | --- | --- |
| Image Builder API controller · Status Watch Controller | IBM Regional Control Plane | IBM |
| HCL generation · profile default resolution · pre-dispatch validation | Image Builder API controller (RIAS) | IBM |
| Tekton pipeline · Task 1 · Task 2 · Cleanup task | IBM Continuous Delivery Tekton (dedicated IBM PAYG account) | IBM |
| Build VSI (Flex profile, `cxf-4x8`) | Customer VPC (zone-resident) | IBM via IAM S2S — customer billed |
| packer.hcl · inventory · build artifacts | IBM COS | IBM |

The build VSI executes inside the customer's own VPC. This **hybrid above-the-cloud model**
means IBM orchestrates the pipeline from above while compute runs in the customer's zone,
enabling VPC-native build isolation and customer-side billing without requiring IBM to manage
zone-resident build infrastructure. Build orchestration runs in IBM's Continuous Delivery
Tekton pipelines in an image builder service account; `packer-plugin-ibmcloud` (Apache 2.0)
provisions and manages the build VSI via IAM service-to-service authorization.

**Resource and identity model:** Image Builder is a **feature extension of the image service**.
Two new top-level regional resources are registered with the resource controller:
**`image-recipe`** (service name
`is.image-recipe`) and **`image-component`** (service name `is.image-component`). Two service
IDs are created in the VPC Compute
account under `is.image-recipe`: **Operator** (under which both controllers run) and
**Administrator**. `image-component` requires no service IDs or API keys.

**API path naming:** the customer-facing paths nest under an `/image/` prefix rather than
sitting at the top level or matching the RC service name verbatim — `/v1/image_recipes`,
`/v1/image_components`, `/v1/image/jobs`. The resource name
following the prefix does not repeat "image" (e.g. not `/v1/image/image-recipe`), since the
prefix already establishes the namespace and a second "image" would be redundant. This reflects
Image Builder's status as a feature extension of the image service (2.6.1, 3.1.1) rather than a
sibling top-level service — the path structure is namespaced under `/image/` the same way the
resources themselves are namespaced under `is.image` at the RC/CRN level, even though the
individual RC service names (`is.image-recipe`, `is.image-component`) are still their own
distinct identifiers.

**Phase 1 build profile:** All 10 IBM-curated components install via precompiled packages.
The build VSI uses the **Flex profile (`cxf-4x8`)**, since NGDC does not have Gen2 capacity.
No GPU hardware
or quota is required for Phase 1 builds either way. GPU
quota is only needed to boot the resulting image.

**Subnet and profile configuration:** `subnet` and `worker.profile` are specified directly on the
image job at creation time (`POST /v1/image/jobs`). `subnet` is required on every job and
tells the service which subnet, and thereby which VPC, to create the temporary build VSI in
within the customer's account. `worker.profile` is optional and defaults to `cxf-4x8` when omitted.
There is **no subnet auto-detection**: IBM
Cloud VPC has no default-VPC/default-subnet primitive, and the subnet determines the security
boundary of the build VSI, so the service will not guess — the customer must supply one on
every job. See
"Subnet and Profile on the Job, Not a Settings API" under Key Architecture Design Points
(3.1.1.2) for the design rationale and history.

---

#### 3.1.1.1 High-Level Architecture

##### Service Architecture Overview

VPC Image Builder uses a hybrid above the cloud approach. The control plane, the
pipeline orchestration, and the build artifacts all run above the cloud, in infrastructure
that IBM manages. The one exception, and the reason this is a hybrid approach rather than a
purely above the cloud one, is the ephemeral build VSI. IBM creates it, controls it,
and deletes it again through the plugin, using the customer's IAM S2S grant, but it actually
runs inside the customer's own VPC. So IBM manages it fully, even though it doesn't live on
IBM's own infrastructure. The
capability manages three primary responsibilities:

1. **Component Catalog** — 10 IBM-curated GPU driver components (customer read-only in Phase 1,
   aside from the deprecate/obsolete actions;
   no public write endpoint for content — authoring happens via internal/operator tooling only)
2. **Recipe Management** — Customer-owned recipes with content-hash versioning and history
3. **Build Orchestration** — Automated image creation via Tekton and packer-plugin-ibmcloud.
   IBM orchestrates above the cloud; the build VSI executes in the customer's VPC.

```text
+-----------------------------------------------------------------------+
|         VPC Image Builder  (is.image-recipe / is.image-component)     |
|         IBM Regional Control Plane -- Above the Cloud                 |
+-----------------------------------------------------------------------+
|                                                                       |
|  +-------------------------------------------------------------------+
|  | Image Builder API controller                                     | |
|  | ----------------------------                                     | |
|  | - Customer-facing REST  (RIAS-Ingress)                           | |
|  | - IAM auth . Pre-dispatch validation                             | |
|  | - POST/GET/PATCH/DELETE /v1/image_recipes (+ /history)           | |
|  | - GET /v1/image_components . POST/GET/DELETE /v1/image/jobs      | |
|  |   (subnet required, profile optional on job creation -- 3.1.1.2) | |
|  | - HCL generation -> packer validate -> hcl_sha256 -> COS stage   | |
|  | - Pipeline dispatch: CreatePipelineRun (CD REST API)             | |
|  +-------------------------------------------------------------------+
|                                                                       |
|  +-------------------------------------------------------------------+
|  | Tekton Pipelines  (IBM Continuous Delivery . IBM PAYG account)   | |
|  | -----------------------------------------------                  | |
|  | - Task 1  Packer Execute -- IBM-built task image                 | |
|  |           Fetch packer.hcl from COS . VERIFY hcl_sha256          | |
|  |           Provision build VSI in customer VPC                    | |
|  |           Install components (generated inline provisioners)     | |
|  |           Validate (package-manager integrity checks only)       | |
|  |           Capture image                                          | |
|  | - Task 2  Post-process -- inventory . COS upload                 | |
|  | - Cleanup task -- ALWAYS runs (success, failure, cancel)         | |
|  +-------------------------------------------------------------------+
|                                                                       |
|  +-------------------------------------------------------------------+
|  | Status Watch Controller                                          | |
|  | ------------------------                                         | |
|  | - Polls Tekton PipelineRun status (CD REST)                      | |
|  | - Writes the image job record (internal surface):                | |
|  |   queued -> running -> succeeded | failed | deleting .           | |
|  |   status_reasons: empty while running, 6 failure codes on failed | |
|  | - Never patches image.status -- the image service derives it     | |
|  |   (image created by the plugin's own ImageBySourceVolume call,   | |
|  |    Task 1 Capturing stage, not by this controller -- zone        | |
|  |    distribution then automatic, task blocks until available)     | |
|  +-------------------------------------------------------------------+
|                                                                       |
|  +-------------------------------------------------------------------+
|  | Storage Layer  (Kubernetes CRDs, regional cluster)               | |
|  | --------------------------------                                 | |
|  | - Kind: ImageComponent (image.vpc.cloud.ibm.com/v1)              | |
|  | - Kind: ImageRecipe     -- labeled by account  (+ history)       | |
|  | - Kind: ImageBuildJob                                            | |
|  +-------------------------------------------------------------------+
|                                                                       |
+------------------------+------------------------------------------------+
                         |
                         |  VPC API . IAM S2S . packer-plugin-ibmcloud
                         v
+-----------------------------------------------------------------------+
|                Customer VPC  (zone-resident)                          |
+-----------------------------------------------------------------------+
|  - Build VSI (Flex `cxf-4x8`) -- ephemeral . IBM managed              |
|    Provisions in customer subnet . Customer billed for runtime        |
|    Plugin deletes VSI + SSH key + SG + FIP on all exit paths          |
|                                                                       |
|  - IS.Image (output)      -- status: available after Task 2           |
|    inventory . image_job sub-resource                                 |
|                                                                       |
|  - VPC Subnet             -- customer owned                           |
|    Public gateway required (for Artifactory downloads)                |
|    Specified on POST /v1/image/jobs (required)                        |
|    (no auto-detection -- see 3.1.1.2)                                 |
+-----------------------------------------------------------------------+
```

| ![VPC Image Builder Service Architecture — Above the Cloud](materials/vpc-ib-hla.png?raw=true) |

**Service characteristics:**

- **Resources**: `image-recipe` + `image-component`, feature extension of the image service
- **Model**: Hybrid above-the-cloud — IBM control plane above; build VSI zone-resident in
  customer VPC
- **Regional**: Each region has independent Image Builder controllers and a dedicated Tekton
  pipeline (Beta: WDC, DAL, FRA)
- **Stateless controllers**: Standard VPC controller pattern
- **Storage**: Kubernetes CRDs (regional cluster) for service metadata; IBM COS for build
  artifacts
- **Execution**: Customer VPC via packer-plugin-ibmcloud and IAM S2S authorization
- **Build profile**: Flex profile (`cxf-4x8`) for all Phase 1 builds — no GPU hardware or quota
  required

---

#### 3.1.1.2 Key Architecture Design Points

##### Customer Account Model

**Decision**: Build VSI runs in the customer's own VPC. Customer is billed for VSI runtime.
IBM's Tekton pipeline runs above the cloud in IBM's infrastructure.

**Rationale**:

- Consistent with AWS EC2 Image Builder and Azure VM Image Builder — customers pay for the
  build instance in all major CSP implementations
- Eliminates the need for IBM zone-resident service nodes and NBD device pool management
- Simplifies the security model — build isolation is VPC-native
- Customer creates IAM S2S authorization once; IBM uses it to provision build resources

**Flow**:

```text
Customer: POST /v1/image/jobs { image_recipe: { id [, version] },
                                 worker: { subnet: { id } [, profile: { name }] }
                                 [, image_template: { name, resource_group }]
                                 [, name] [, timeout] }
                                 (worker.subnet required; worker.profile optional, defaults to cxf-4x8;
                                  image_template optional -- identity generated if omitted;
                                  no separate top-level resource_group on the job itself --
                                  only image_template.resource_group, for the resulting image)
  v IBM pre-dispatch validation
    -> validate subnet  ->  default profile if omitted  ->  validate S2S  ->  create job (queued)
  v Generate packer.hcl . packer validate . hcl_sha256 . stage to COS
  v CreatePipelineRun (CD REST) -> Tekton
    -> Task 1: fetch + verify template . provision Flex (cxf-4x8) VSI in customer VPC
              install GPU drivers via dnf / apt-get (Artifactory)
              validate (package-manager integrity checks only) . capture boot volume
    -> Task 2: inventory extracted . COS upload . plugin calls POST /v1/images
              itself ("image from volume") -> image CREATED
              here, with real content, in one step . build metadata written
    -> Cleanup task: delete SSH key, SG, FIP, VSI (always)
  v image created in Task 1 Capturing stage . job.image populated . task blocks until available
    -> image available in all zones . job: succeeded -- exact
       timing of succeeded relative to zone distribution completing is OPEN,
       see 3.1.1.3.3
  v Customer: GET /v1/image/jobs/{id} -> succeeded . image ref -> GET /v1/images/{id} -> available
```

##### Subnet and Profile on the Job

**Decision**: `subnet` and `profile` are both specified on `POST /v1/image/jobs`, nested
together under `worker`. `worker.subnet`
is required; `worker.profile` is optional and defaults to `cxf-4x8` when omitted. There is no
account-level settings API, and no subnet auto-detection. See 5.2 for the design history
behind this.

**Customer journeys**:

```text
Journey 1 -- One-time setup, then build

  [one-time] Create IAM S2S authorization grant in IBM Cloud IAM console
  POST /v1/image_recipes  { base_image: {...}, components: [...] }
  POST /v1/image/jobs { image_recipe: { id }, worker: { subnet: { id: "0717-abc123" } } }
  GET  /v1/images/{id}  ->  poll until available

Journey 2 -- Repeated builds against the same subnet

  POST /v1/image/jobs { image_recipe: { id: "..." }, worker: { subnet: { id: "0717-abc123" } } }
  POST /v1/image/jobs { image_recipe: { id: "..." }, worker: { subnet: { id: "0717-abc123" } } }
             (customer or their own tooling supplies the same subnet each time --
              no service-side memory of a "preferred" subnet across jobs)
```

**No account-level `default_build_profile` field is needed either.** The default lives
directly in the job-creation endpoint's own field description, self-contained. The
resolution model is simply: `worker.profile` on the job, defaulting to `cxf-4x8` when unset, full
stop. The
resolved profile is recorded in the image's build metadata, so the resolution is always
auditable per build. There's deliberately no service-invented default beyond `cxf-4x8` —
worth noting the asymmetry
with `base_image`, which has no default at any layer at all, since a base default would mean
the service choosing the customer's OS for them.

Phase 1: every component has `requires_gpu_for_build: false`, so the Flex profile (`cxf-4x8`)
is the documented recommended profile.

##### precompiled_package for All Phase 1 Components

**Decision**: All 10 Phase 1 components install via precompiled packages. GPU drivers
are installed as pre-compiled RPM/deb packages via dnf/apt-get from Artifactory. No DKMS kernel
module compilation. No GPU hardware on the build VSI.

**Rationale**:

- GPU quota is not required for Phase 1 builds — customers only need GPU quota to boot images
- The Flex profile (`cxf-4x8`, ~$0.09/build) is sufficient for all Phase 1 package installations
- Eliminates the largest adoption barrier (GPU quota prerequisite)

See 2.4 for the `install_method` extension point for customer
components.

##### Plugin-Managed Ephemeral Resources

**Decision**: packer-plugin-ibmcloud auto-manages SSH key, security group, and floating IP per
build. All resources deleted on all exit paths including failures and cancellation.

**Rationale**:

- Customer never configures keys, SGs, or FIPs — zero networking setup required
- Azure VM Image Builder requires customers to configure NSG rules manually; VPC Image Builder
  eliminates this entirely via the plugin's ephemeral SG
- Guaranteed cleanup prevents resource leaks in the customer account

##### Component Catalog — Owner-Managed, Generated Builds

**Decision**: The 10-component Phase 1 catalog is managed under the **hybrid
ownership model**: catalog entries (components and stock recipes) are owned by a dedicated IBM
catalog account and onboarded/lifecycle-managed **through internal/operator tooling as the
owning account**,
without a public-facing write endpoint in Phase 1. **Publicization**
(visibility flip to `public`, applying the `ibm-` prefix) is a privileged
operator action, via an internal `PATCH` on the resource. `visibility` is a
**read-only**
`public | private` property that also doubles as a list filter (`?visibility=`) on both
`GET /v1/image_recipes` and `GET /v1/image_components`, never
customer-settable. Operator APIs additionally cover
cross-account safety checks (safe-to-delete).

Customers see a read-only catalog in Phase 1, aside from the ability to deprecate or
obsolete a component directly via its dedicated action endpoints — the one narrow exception
to an otherwise read-only surface. Component
authoring and content changes happen entirely through internal/operator tooling, not a customer-facing endpoint.
Components carry **declarative build knowledge**: `software.name`/`software.version`
identify what to install, and `dependencies` (see 3.1.1.3.6a
for hardware/OS compatibility) order the install sequence, from which the HCL generator derives
the install provisioner blocks. **No authored install or validate scripts
exist**; Artifactory holds binaries only.

**Phase 1 component catalog** — all `precompiled_package`, all `requires_gpu_for_build: false`.
"Vendor family" below is this document's own organizational grouping for readability, not a
schema field — see 3.1.1.3.6a and 5.2:

| Vendor family | Components | Target GPU profiles |
| --- | --- | --- |
| **NVIDIA** | nvidia-data-center-driver, nvidia-cuda-toolkit, nvidia-cudnn, nvidia-nccl | gx2, gx3 |
| **AMD** | amd-gpu-driver, amd-rocm, amd-miopen | gx2-amd, gx3-amd |
| **Intel** | intel-gaudi-driver, intel-synapse-ai, intel-oneapi | gx3d-intel |

Specific component versions are not pinned in this document — which versions are actually
needed, and how many distinct stock recipe/component combinations are warranted per GPU
generation, is being worked out through testing and is expected to be settled before
commitment review (see OI-18).

##### Free Service

**Decision**: Phase 1 is completely free with no billing, Resource Controller integration, or
usage metering. Customers are billed only for the underlying VPC resources (build VSI runtime,
image storage) through standard VPC billing. No cost-estimation field exists on the API —
consistent with AWS and Azure (pricing documentation before, billing tools after, nothing on
the resource), and the build bills to the customer's own account where their cost tooling
already applies.

---

#### 3.1.1.3 Component Design

##### 3.1.1.3.1 Image Builder API Controller

**Responsibilities**:

- Customer-facing REST API registered via RIAS-Ingress
- IAM authentication and authorization enforcement
- Recipe content management: content-hash version minting, history maintenance,
  whole-array `image_components` replacement on `PATCH` (no element-level merge
  in the generated API spec; there are no separate add/remove sub-operation endpoints),
  duplicate-content detection (`409` — scoped to the _current_
  version only, so restoring a historical version's content via PATCH remains legal)
- Pre-dispatch validation (before any VSI is provisioned)
- profile resolution (job request value, defaulting to `cxf-4x8`)
- HCL generation using the **`hclwrite`** library (Go, programmatic HCL construction — no
  string templating): read recipe frozen content and component records from their CRDs · generate
  packer.hcl (install/validate provisioners derived from component schema, dependency-ordered)
  · `packer validate` · compute `hcl_sha256` · stage to IBM COS at `builds/{build_id}/packer.hcl`
- Pipeline dispatch: `CreatePipelineRun` via the Continuous Delivery **REST API** with
  parameters `{build_id, image_id, hcl_location, hcl_sha256, region}`

**Pre-dispatch validation sequence** (all before any VSI is provisioned; **"build" here refers
specifically to this pre-dispatch-through-capture/upload window — see the note on "build" vs.
"import" terminology at the end of 3.1.1.3.4**):

```text
+---------------------------------------------------------------------+
| Pre-dispatch Validation  (POST /v1/image/jobs)                      |
+---------------------------------------------------------------------+
|  1. Validate recipe exists; status gate (deprecated -> warn,        |
|     obsolete -> reject; degraded -> reject, component_obsolete)     |
|  2. Resolve pinned version (request version | current content)      |
|  3. Resolve profile  (request value, defaulting to cxf-4x8)         |
|  4. Validate subnet  (required on the request; 404 if it doesn't    |
|     exist in the account -- see 3.1.1.2)                            |
|  5. GET /v1/subnets/{id} -> validate public gateway                 |
|  6. Validate profile available in subnet zone                       |
|  7. Verify IAM S2S grant exists (customer-created prerequisite)     |
|  8. Create the image job record (status: queued) -- no image        |
|     exists yet; the image is the job's eventual output              |
|  9. Generate HCL -> packer validate -> hcl_sha256 -> COS ->         |
|     CreatePipelineRun (CD REST)                                     |
+---------------------------------------------------------------------+
```

##### 3.1.1.3.2 Pipeline Dispatch (REST) and the gRPC Exceptions

Dispatch to Tekton uses the IBM Continuous Delivery **REST API** (Go SDK) — not gRPC. The
internal image job record write is REST as well. **Two** gRPC calls exist in this design,
in opposite directions: (1) image service → Image Builder, notifying Image Builder of image
state changes (including deletion), so Image Builder can clean up its own job record; and (2)
Image Builder → image service, writing the `source_image_recipe` and `image_job` fields onto the image once it exists,
since the plugin's generic `ImageBySourceVolume` call has no way to populate
Image-Builder-specific data itself. Every other internal
and customer-facing surface in this design is REST. Image creation happens entirely via the
packer
plugin's own existing call — `ImageBySourceVolume` → `POST /v1/images` — during Task 1's
Capturing
stage; there is no separate `ImageByRecipe` creation mechanism — the plugin's
existing volume-based call already covers image creation. See 3.1.3 and 5.2.

**On when call (2) actually happens**: in practice, writing `source_image_recipe`/`image_job` onto the image and the
image reaching `available` happen at effectively the same time today — but this is a
consequence of how the Packer plugin currently works, not a structural requirement of the
design. The plugin's own capture step polls and blocks on `available` internally, and there is
currently no way for the image ID it obtains from its `CreateImage` call to reach the Image
Builder control plane any earlier than that — the plugin doesn't surface it as an accessible
interim result. The Image Builder control plane itself has no inherent reason it couldn't act
sooner if that information became available to it earlier. This is a named candidate for a
future Packer/plugin enhancement (e.g. surfacing the image ID via a Tekton task result once
known, before the plugin's own blocking wait completes), not a permanent limitation.

**`source_volume` remains on the image.** Considered for removal (since it references a
temporary build volume with no lasting customer relevance), but kept — see 5.2 for the
rationale.

**CreatePipelineRun parameters**:

```text
job_id         the image job (created at pre-dispatch step 8, status: queued)
hcl_location   builds/{job_id}/packer.hcl  (IBM COS)
hcl_sha256     integrity hash -- Task 1 verifies before use; recorded in build metadata
region
account_id
```

> No image identifier is passed to `CreatePipelineRun` — the pipeline doesn't need one; Task 1
> works from the recipe content already embedded in packer.hcl. No image record exists yet at
> dispatch time; the image comes into existence later, via the plugin's own call in Task 1's
> Capturing stage — see 3.1.1.3.3 and 3.1.1.9.3a.
> The base image for the build VSI is resolved by the control plane and embedded in packer.hcl
> (`vsi_base_image_id`, by ID — IDs pin).

##### 3.1.1.3.3 Status Watch Controller

**Responsibilities**:

- Polls Tekton PipelineRun status (CD REST API) during the build phases, and answers
  `GET /v1/image/jobs/{id}` by deriving the job's coarse status live from that same
  `PipelineRun` state at query time: `queued → running → succeeded | failed | deleting`.
  There is no separate internal write endpoint on the IS.Image side for this — status is
  read-derived from Tekton on demand, not pushed to and persisted in a separate record.
  `status_reasons` does
  **not** carry live stage-progress detail while `running` — it stays empty until the job
  either succeeds (stays empty) or fails (populated with one of 6 failure codes, derived from
  the `PipelineRun`'s own failure detail at that same read time). See
  3.1.1.9.3a for the full status model.
- The plugin's `step_capture_image.go` builds an
  `imagePrototype` and calls the VPC Go SDK's `CreateImage`, using `ImageBySourceVolume` →
  `POST /v1/images`, during Task 1's Capturing stage.
  - The plugin's capture step polls the image's own status and blocks until it reaches
    `available`, only calling `Cleanup()` (tearing down the build VSI and boot volume) once
    that happens. The Tekton task does not complete until zone distribution has fully
    finished — so Status Watch Controller, polling the PipelineRun, does not see the pipeline
    complete until `available` is reached. `job.status` therefore stays `running` through zone
    distribution, reaching `succeeded` only once it's fully complete.
  - Two gRPC calls exist in this design, the exceptions to its otherwise REST-only
    pattern: RIOS notifies Image Builder when a customer deletes an image through the ordinary
    image APIs, so Image Builder can clean up its own related build job records; and Image
    Builder notifies the image service to write the `source_image_recipe` and `image_job`
    fields onto the image once it exists, since the plugin's generic call has no way to
    populate Image-Builder-specific data itself — see OI-RIOS for both. The extracted
    inventory is copied directly into the COS archive bucket by the Tekton task itself (Task
    2, see 3.1.1.3.4), so no separate call is needed to write it onto the image.
- **Image Builder never patches `image.status` directly.** The image service derives
  `image.status` itself, through `ImageBySourceVolume`'s standard zone-distribution behavior.
  The image moves `pending` → `partially_available` → `available` — a real, progressive
  intermediate state, on the standard VPC image-creation path, with zero role for Image
  Builder to design or trigger. `partially_available` is a pre-existing status on the base
  `Image` resource, not something Image Builder introduces. Cancelling/deleting an image moves it to `deleting`, a
  transitional status, before the resource is removed entirely — there is no separate
  terminal `deleted` status. **Image creation is entirely via `ImageBySourceVolume`**.
  Full detail in 3.1.3
  and 3.1.1.3.4.

##### 3.1.1.3.4 Tekton Pipeline

**Task 1 — Packer Execute** (IBM-built task image — pinned Packer + plugin, from ICR):

> **On the stage markers below**: the pipeline genuinely progresses through these stages
> internally, and the `[job: running . build_*]` labels below describe that real internal
> progression — none of this is surfaced to the
> customer via `status_reasons` (see 2.2).

```text
+---------------------------------------------------------------------+
| Task 1 -- Packer Execute                                            |
|          Tekton-managed pod -- no dedicated IBM IKS worker nodes    |
+---------------------------------------------------------------------+
|  - Read env properties: is.image-recipe API key (secret),           |
|    desired-account CRN                                              |
|  - Fetch packer.hcl from IBM COS -> VERIFY hcl_sha256               |
|  - packer init (plugin pre-baked -- no network pull)                |
|  - packer build:                                                    |
|                                                                     |
|  Provisioning  [job: running . build_provisioning_worker]           |
|    Auto-generate RSA 4096 key pair in memory                        |
|    POST /v1/keys . POST /v1/instances (Flex cxf-4x8, customer       |
|    subnet)                                                          |
|    POST /v1/floating_ips . ephemeral SG (TCP/22 in, all out)        |
|                                                                     |
|  Installing  [job: running . build_installing_components]           |
|    SSH into build VSI via floating IP                               |
|    Generated inline install provisioner: dnf/apt-get install        |
|    with exact pinned versions, dependency (toposort) order          |
|                                                                     |
|  Validating  [job: running . build_validating]                      |
|    Generated inline validate provisioner:                           |
|    package-manager integrity check per package -- exact commands    |
|    TBD before commitment review, see OI-VALCMD (package-level only, |
|    see 3.1.1.3.6a and 2.8.2.9; no GPU hardware present, and no      |
|    per-component validation commands exist to run beyond this)      |
|    + rpm -qa / dpkg-query -W inventory (feeds the software          |
|    inventory)                                                       |
|                                                                     |
|  Capturing  [job: running . build_capturing]                        |
|    Per step_capture_image.go: plugin builds an imagePrototype,      |
|    calls VPC SDK CreateImage using                                  |
|    ImageBySourceVolume -> POST /v1/images. IMAGE IS CREATED HERE,   |
|    with real content, in Task 1. Plugin then                        |
|    POLLS the image's own status and BLOCKS until it reaches         |
|    available (pending -> partially_available -> available)          |
|    before calling Cleanup() -- this task does not                   |
|    complete until zone distribution is fully done                   |
+---------------------------------------------------------------------+
```

**Task 2 — Post-process** (Tekton-managed pod · above the cloud):

```text
+---------------------------------------------------------------------+
| Task 2 -- Post-process                                              |
+---------------------------------------------------------------------+
|  - Retrieve the plugin-generated scan output (Packer scanning;      |
|    wrapper-injected recipe provenance) and extract the software     |
|    inventory (name/version list) from it                            |
|  - Image creation happens in Task 1's Capturing stage, per          |
|    step_capture_image.go, via the plugin's                          |
|    ImageBySourceVolume call, which blocks until the image reaches   |
|    available before Task 1 completes -- see 3.1.1.3.3. By the time  |
|    Task 2 runs, the image already exists and is already available   |
|  - Upload the extracted inventory file -> IBM COS:                  |
|    builds/{job_id}/inventory.json, following                        |
|    the Images team's existing archive-bucket convention --          |
|    see 3.1.1.9.3                                                    |
|  - Assemble build metadata: recipe (id, name, version), resolved    |
|    base image ID, components with as-installed versions,            |
|    hcl_sha256 -- assembled locally, not written to IS.Image         |
|    separately. Image Builder calls the image service via gRPC to    |
|    write the recipe and image_job fields onto the image             |
|    after Task 1's ImageBySourceVolume call creates it -- the        |
|    plugin's generic call has no inherent way to populate either     |
|    field itself (see 3.1.4, OI-RIOS)                                |
|  - Zone distribution is complete by this point already --           |
|    Task 1's Capturing stage does not hand off to Task 2 until the   |
|    plugin's own polling loop observes available, per                |
|    step_capture_image.go. job.status therefore stays running        |
|    through zone distribution and reaches succeeded only once it's   |
|    fully complete (see 3.1.1.3.3)                                   |
|  - Emit AT event (success)                                          |
+---------------------------------------------------------------------+
```

**Cleanup task — ALWAYS runs** (success, failure, cancel, image deletion):
deletes build VSI, SSH key, floating IP, security group, and per-build COS scratch.
Customer billing is bounded to the VSI's running duration.

**"Build" and "import" are not separate stages.** Zone distribution happens
inside Task 1's Capturing stage — the plugin's own capture step blocks on the image reaching
`available` before the task completes. `timeout` (2.2, 3.1.1.9.3a) bounds the whole of Task 1,
including that blocking wait. Any zone-distribution
failure surfaces as an ordinary Task 1 failure under the existing failure codes (3.1.1.9.3a).

| ![VPC Image Builder — Tekton Pipeline and Dispatch](materials/vpc-ib-tekton-dispatch.png?raw=true) |

##### 3.1.1.3.5 IAM S2S Authorization and Credential Lifecycle — including Packer's Auth Flow

The customer creates the IAM S2S authorization grant once as a prerequisite. The customer must
create **4 individual authorization policies** before their first build:

| Source | Target | Role |
| --- | --- | --- |
| is.image-recipe | is.instance | Operator |
| is.image-recipe | is.key | Operator |
| is.image-recipe | is.floating-ip | Operator |
| is.image-recipe | is.security-group | Operator |

Notes: VPC grants support resource-type scoping only (no per-instance scoping); a broader
`is.*` scope is not used. See 2.4 for IAM service-group authorization as a future
consideration.

**How Packer specifically authenticates (per-build flow)**: because builds can exceed the
1-hour IAM token lifetime, an **API key (not a token)** is provided to the Packer plugin —
this is the credential Packer itself consumes, distinct from the customer's one-time S2S grant
above:

1. The regional `is.image-recipe` **Operator** API key (created and rotated by IBM — see the
   resource-lifecycle table below) is stored as a Tekton pipeline environment property
   (secret).
2. **Task 1 passes this key, plus the desired-account CRN, to `packer-plugin-ibmcloud`.**
3. The plugin performs a **two-token exchange**: the API key is exchanged for a short-lived
   IAM token; that token, combined with the target account CRN
   (`crn:...:image-recipe:...`), is exchanged again for a **customer-scoped token**.
4. The plugin uses the customer-scoped token to create and later delete the build VSI, SSH
   key, floating IP, and security group **in the customer's account** — authorized by the
   4 policies the customer created, billed to the customer.

No customer credential is ever handled by Packer or the plugin — the customer's only
involvement is the one-time policy grants above; everything else in this flow is IBM-owned
and IBM-rotated.

**How the plugin produces the custom image**: it captures the build VSI's boot
volume directly and calls the standard `POST /v1/images` itself ("image from volume"). Images
built through Image Builder are
available across all zones in the region, the same as any other VPC custom image — not
scoped to wherever the build VSI happened to run. Zone
distribution is a genuine,
progressive, automatic process — image status moves `pending` → `partially_available`
→ `available`, with `partially_available` specifically reflecting some-but-not-all zones
completed — inherited for free from the standard `POST /v1/images` call, with zero role for
Image Builder to design or trigger. The plugin's separate image-export-job capability is an
optional post-processor, only runs if explicitly included, calls its own export-image-job
API — not part of the creation path at all. See 3.1.1.3.3 for the full detail and 3.1.1.9.3
for the inventory discussion.

```text
+---------------------------------------------------------------------+
| Credential flow (per build)                                         |
+---------------------------------------------------------------------+
|  is.image-recipe Operator API key (regional; Secrets Manager /      |
|  Vault) -> stored as a Tekton pipeline environment property         |
|  (secret) -- rotation AUTOMATED: a controller (Vault Agent          |
|  sidecar) pushes the rotated key into the Tekton property via       |
|  the Tekton properties API                                          |
|                                                                      |
|  Task 1 passes the key + desired-account CRN to the plugin.         |
|  Plugin two-token flow: API key -> IAM token; token + target        |
|  account CRN (crn:...:image-recipe:...) -> customer-scoped token -> |
|  create/delete build resources in the customer's account            |
|  (grants above authorize it; resources bill to the customer)        |
+---------------------------------------------------------------------+
```

**Complete resource lifecycle** — customer manages nothing beyond the one-time grant:

| Resource | Who creates | Who deletes | Lifecycle |
| --- | --- | --- | --- |
| IAM S2S authorization | Customer (one-time prerequisite) | Customer | Persistent per account |
| is.image-recipe Operator API keys (controller, watcher, Tekton property — regional) | IBM | IBM rotates (automated) | IBM-managed |
| Customer-scoped token (plugin-derived) | Plugin per build | Expires naturally | ≤1 hour, renewable from key |
| SSH key pair | Plugin | Cleanup task (all exits) | Build duration only |
| Security group | Plugin | Cleanup task (all exits) | Build duration only |
| Floating IP | Plugin | Cleanup task (all exits) | Build duration only |
| Build VSI | Plugin | Cleanup task (all exits) | Build duration only |

##### 3.1.1.3.6 Recipe Validation

**Validation at recipe create/PATCH** (every content change is fully re-validated):

- OS/architecture compatibility: each component's `allowed_use.image` CEL expression is
  evaluated against the recipe's `base_image` (`component_os_incompatible`, 400) — see
  3.1.1.3.6a
- Dependency resolution: all component dependencies must be present in the recipe —
  a flat array of required component references, with no `all_of`/`any_of` operator wrapper
- Exclusivity: no vendor-mutual-exclusivity check exists — see 5.2 for that history
- Hardware compatibility: fully dropped from Phase 1, with no replacement built — see
  3.1.1.3.6a for the accepted tradeoff
- Component status gate: `deprecated`/`degraded` → warning; `obsolete` → rejected
- Component count: 1–20 (`recipe_component_limit_exceeded`)
- Base image change: full re-validation of every component against the new base;
  cross-family change carries the `base_image_family_changed` warning

**Validation at `POST /v1/image/jobs`** pre-dispatch (before any VSI provisioned):

| Constraint | Error code | HTTP | Action |
| --- | --- | --- | --- |
| Recipe obsolete | `recipe_obsolete` | 409 | Reject |
| Recipe degraded (component obsoleted) | `component_obsolete` | 409 | Reject with target component |
| Recipe degraded (base image obsoleted) | `base_image_obsolete` | 409 | Reject — adopted by symmetry with component handling, see 3.1.1.9.4 |
| `subnet` missing from the request | Standard missing-required-field `400` | 400 | Reject — `subnet` is a required request field on `POST /v1/image/jobs`, not a separate resolution step |
| Subnet doesn't exist in the account | Standard `404` | 404 | Reject |
| Subnet has no public gateway | `subnet_no_public_gateway` | 409 | Reject |
| Profile not available in subnet zone | `profile_not_available_in_zone` | 409 | Reject |
| IAM S2S grant not found | `iam_s2s_missing` | 409 | Reject — customer must create it before first build |

**Five of these are `409`; two (missing or nonexistent `subnet`) are standard `400`/`404`
rejections, not part of the `409` set.** The
codes named above come from earlier, separate confirmation and aren't independently
re-verifiable from a single spec pull — treat the trigger conditions as the team's best
current understanding, not a byte-for-byte transcription of a discriminated error schema.

###### 3.1.1.3.6a Hardware and OS Compatibility

**OS and architecture compatibility now has a real mechanism: `allowed_use.image`, a CEL
expression on the component.** A component declares a
[Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md)
string using `operating_system.family`, `operating_system.version`, and
`operating_system.architecture` as variables; recipe validation evaluates this expression
against the recipe's
`base_image` at recipe creation/PATCH time. This directly solves the problem that surfaced
with the NVIDIA driver example —
Ubuntu and RHEL builds of the same conceptual driver have different `software.name`/`version`
values with nothing distinguishing which OS each one targets, so a customer combining the
wrong one with their base image previously only found out via a failed install at build time.
The component's `allowed_use` is a bespoke schema with only `api_version` and `image` — no
`instance` or `bare_metal_server`, unlike the image's own pre-existing `allowed_use`. See 5.2
for the fuller history of this decision, including the earlier, fully-dropped
`allowed_use` design and the plain-object alternative considered before landing here.

**Hardware and GPU-manufacturer compatibility still has no mechanism.** There's no
`bare_metal_server`/`instance` expression coverage on the component, and no `vendor` field — nothing prevents mixing
components from different vendors (e.g. an NVIDIA driver alongside an AMD toolkit) in the same
recipe. See 2.8.2.9 for the risk this still leaves open, narrowed now to hardware/vendor only
rather than covering OS/architecture as well.

**`dependencies` is a flat array** (`ImageComponentReference[]`, `maxItems: 19` — one less
than the 20-component-per-recipe cap, since a component's own dependencies are counted against
that same limit)
— no `all_of`/`any_of` wrapper structure exists in the real schema.

**`validation_commands` is planned for Phase 2, not Phase 1 — components carry no declared
validation logic today.** Phase 1's package-level validation commands are authored entirely
by the HCL generator itself, not sourced from any user- or customer-provided validation
field — there is no way for a customer, or even the component's own catalog entry, to supply
validation logic in Phase 1. In-build validation is limited to whatever the generator emits
independent of the component
itself: package-manager-level integrity checks (confirming
files installed correctly, checksums match; the specific commands the generator will emit are
being finalized before
commitment review — see OI-VALCMD). There is no per-component check confirming the
installed software actually functions (e.g. a driver loading, a toolkit reporting its version
correctly) at build time — see 2.8.2.9 for the risk this creates.

---

#### 3.1.1.4 Data Models

> **Note on examples**: reference objects (`base_image`,
> `resource_group`) and
> `href` values are abbreviated for readability. In the API specification, references follow
> the standard VPC reference schemas (`ImageReference`, `ResourceGroupReference`, etc.) with
> full `href` URLs and `crn`/`resource_type`/`deleted` where applicable. **`image_components[]`
> entries are the one exception** — `{id, href}` only, not an abbreviated form of
> a fuller reference shape.

##### 3.1.1.4.1 Component (public representation)

```json
{
  "id": "my-image-component",
  "crn": "crn:v1:bluemix:public:is:us-south:a/{ibm-catalog-account}::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8",
  "href": "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8",
  "name": "my-image-component",
  "software": {
    "name": "cuda-toolkit",
    "version": "12.4.0"
  },
  "allowed_use": {
    "api_version": "2024-06-23",
    "image": "operating_system.family == 'ubuntu' && operating_system.version == '22.04' && operating_system.architecture == 'amd64'"
  },
  "dependencies": [],
  "resource_group": { "id": "fee82deba12e4c0fb69c3b09d1f12345", "href": "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345", "name": "default" },
  "resource_type": "image_component",
  "status": "available",
  "created_at": "2026-01-02T03:04:05.006Z",
  "deprecation_at": null,
  "obsolescence_at": null
}
```

Field notes:

- **The nested identity object is named `software`, not `component`.** The resource's own
  catalog-level `id`/`name` is distinct from the actual software
  identity, which lives in a nested `software: {name, version}` object (e.g.
  `software.name: "cuda-toolkit"`, `software.version: "12.4.0"`). There is no outer-level
  `version` field — it's fully replaced by `software.version`. `software` is designed as a
  `oneOf`, even though Phase 1 has exactly one branch (a plain package identity). This is
  forward-looking: a future phase could add a different branch — for instance, a catalog
  offering reference — without disturbing how Phase 1 components are described. It's the same
  extensibility pattern used for the image job's own `type` field (3.1.1.4.3), applied here to
  what a component's underlying software actually is.
- **`allowed_use` is back on the component, and it's a CEL expression, not a plain object.**
  It carries `{api_version, image}`, where `image` is a
  [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md)
  string checked against the recipe's `base_image` at recipe creation/PATCH time. The
  variables available inside that expression are `operating_system.family`,
  `operating_system.version`, and `operating_system.architecture` — all three live inside the
  one string, e.g. `"operating_system.family == 'ubuntu' && operating_system.architecture ==
  'amd64'"`. This resolves the NVIDIA per-OS naming problem directly: since Ubuntu and RHEL
  builds of the same driver have different `software.name`/`software.version` values, this
  expression is what actually lets recipe validation catch a mismatch between a component and
  its base image, rather than relying on the component's name alone. `allowed_use` has no
  `instance` or `bare_metal_server` properties on the component — its own dedicated schema
  leaves those out, unlike the base `AllowedUse` type
  the image itself uses (see 5.2 for why). This is narrower than the image's own `allowed_use` (no
  `bare_metal_server`/`instance` expressions, no hardware/GPU-manufacturer coverage) — see
  2.8.2.9 for what's still not covered.
- **There is no standalone `architecture` field.** Architecture is expressed as
  `operating_system.architecture` inside the `allowed_use.image` expression, not as a separate
  top-level property.
- **There is no `display_name` field on the component.** Only `name` exists.
- **`resource_group` and `resource_type` are required fields** on the component response.
- **`visibility` is a real, public field** (`private | public`) — readable by customers as part
  of the component response. Its write path (publicization, the stock/private distinction) is
  still operator-only; customers cannot set it themselves.
- **`validation_commands` is planned for Phase 2, not Phase 1.** Phase 1 components carry no
  validation commands; the HCL generator does not derive a validate provisioner from anything
  on the component. See 3.1.1.3.6a and the risk discussion in 2.8.2.9 for what in-build
  validation is scoped to without it in Phase 1.
- **`package_details` and `install_method` are not part of the Phase 1 component.** Phase 1's
  install approach is precompiled packages throughout — worth carrying `install_method` as a
  concept into a future phase's design once components support more than one install
  approach, but neither field exists on Phase 1 components today.
- `dependencies` is deliberately generic, not named something like `component_dependencies` —
  today it represents component-to-component dependencies, but the field is meant to extend
  to other dependency types in the future. It's capped at **19 entries**, not 20 — since a
  recipe allows at most 20 components total (3.1.1.7), and a component's own dependencies are
  themselves counted against that same limit, the most a single component could ever
  legitimately require is every other slot in the recipe.
- **`software` and `dependencies` are both immutable once the component is created.** A new
  package version, or a change to what a component depends on, means publishing a new
  component entry — never updating an existing one. See Component Versioning in 3.1.1.9.1 for
  the full rationale.
- `created_at` is present on this resource.

The Phase 1 field list is: `allowed_use` (object: `api_version`, `image` — a CEL expression),
`created_at`,
`crn`, `dependencies`, `deprecation_at`, `href`, `id`, `name`,
`obsolescence_at`,
`resource_group`, `resource_type`, `software` (object: `name`, `version`), `status`.

- `resource_type`'s enum value (`"image_component"` above) is illustrative.
- OS and architecture compatibility are both expressed entirely through `allowed_use.image`.
  `GET /v1/image_components` filters on `software.name`, `status`, and `visibility` — `status`
  takes a single
  value, unlike `GET /v1/images`, which supports an array. See 5.2 for the history of
  alternative filter/field designs considered for this.
- `allowed_use` is the only top-level field of its kind on the component — see 5.2 for
  the history of how it got here.
- `dependencies` is a flat array (`ImageComponentReference[]`, `maxItems: 19`).

Notes: `install_method` and `package_details` (the mechanism a future phase would use for
OS-keyed package lists) are not part of Phase 1 at all — Phase 1's install approach is
precompiled packages throughout, with no discriminator field needed. Both would be introduced
as part of a documented Phase 2 event when customers author components with more than one
install approach. Each component version is a distinct, immutable catalog entry; lifecycle status is **derived**
from `deprecation_at`/`obsolescence_at`. In the API specification, most references follow the
standard VPC reference schemas (`ImageReference` for `base_image`,
`ResourceGroupReference`, `InstanceProfileReference` for the instance profile) with full-URL
`href` values and the standard optional
`deleted` sub-object; abbreviated forms in this document are illustrative. `image_components[]`
entries are the exception — `{id, href}` only.

**Storage**: Kubernetes Custom Resource (`Kind: ImageComponent`, `apiVersion:
image.vpc.cloud.ibm.com/v1`) in the regional cluster — reconciled by the Catalog Controller,
which is already running on Kubernetes alongside the other Image Builder controllers.
Components are not IBM Cloud Catalog entries; adopting the CRD pattern here reuses the same
reconciliation framework the controllers already depend on, rather than introducing a
second, hand-rolled etcd access pattern alongside it.

##### 3.1.1.4.2 Recipe (customer-owned; stock recipes are the same resource with `visibility: public`)

```json
{
  "id": "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8",
  "crn": "crn:v1:bluemix:public:is:us-south-1:a/{account}::image-recipe:my-rhel9-cuda12",
  "href": "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8",
  "resource_type": "image_recipe",
  "name": "my-recipe-cb8",
  "version": "a3f8c2d1",
  "base_image": { "id": "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8", "crn": "...", "href": "...", "name": "ibm-redhat-9-4-minimal-amd64-3", "resource_type": "image" },
  "image_components": [
    { "id": "r006-drv560", "href": "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-drv560" }
  ],
  "status": "available",
  "status_reasons": [],
  "superseded_at": null,
  "visibility": "private",
  "deprecation_at": null,
  "obsolescence_at": null,
  "resource_group": {
    "id": "fee82deba12e4c0fb69c3b09d1f12345",
    "href": "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345",
    "name": "default"
  },
  "created_at": "2026-05-11T09:00:00Z",
  "updated_at": "2026-07-10T14:22:00Z"
}
```

**Field notes:**

- **There is no `images[]` field on the recipe.** Recipe→images discovery is
  entirely via `GET /v1/images?source_image_recipe.id={id}[&source_image_recipe.version={hash}]` instead
  (see 3.1.1.4.3) — one mechanism for both the current version and any historical version,
  rather than a field on the recipe that only ever covered the current version anyway.
  `source_image_recipe.id` additionally supports `null`/`not:null` values, filtering for images with
  no recipe or any recipe respectively.
- **`superseded_at` is present directly on the live `ImageRecipe` resource**,
  included in the `POST`, `GET`, and `PATCH` recipe responses — not only via
  `/recipes/{id}/history`. Since history entries are themselves recipe objects that already
  carry `superseded_at`, the standard recipe responses carry it too, for consistency.
  `superseded_at` is set once a newer version exists, `null`/absent on the current version.
- **There is no `display_name` field on the recipe.** Only `name` exists.
- **There is no `profile` field on the recipe.** `profile` moved to the image job
  (`POST /v1/image/jobs`), defaulting to `cxf-4x8` — see 3.1.1.2.
- **`base_image` is a reference to an image object** — a full `ImageReference` (`crn`,
  `deleted`, `href`, `id`,
  `name`, `remote`, `resource_type`).
- **`visibility` is read-only.**
  `public` for stock recipes, `private` for customer-owned ones — visible to any authenticated
  user who can read the recipe, not gated behind an operator-only view. This resolves
  the stock-vs-custom distinction the design narrative depends on throughout this
  document.

The `ImageRecipe` field list: `base_image`, `created_at`, `crn`, `deprecation_at`,
`href`, `id`, `image_components`, `name`, `obsolescence_at`,
`resource_group`, `resource_type`, `status`, `status_reasons`, `superseded_at`, `updated_at`,
`visibility`,
`version`. There is no `images` field — see the field notes above. There is no `profile`
field — that moved to the image job (see 3.1.1.2). There is
no `warnings[]` field — how advisory information (e.g. a deprecated
referenced component) surfaces is addressed below.

Key properties: `id` remains immutable and unique per
account+region (`recipe_name_conflict` 409 on duplicate name).
`version` is a content hash minted on every content change; the `(id, version)` pair
permanently identifies exact content. `image_components[]` entries are plain references
(`id` + `href`, the standard VPC
reference shape) — `name` and the component's own `version` are not
duplicated onto the recipe; resolving a component's descriptive details is a live
`GET /v1/image_components/{id}` lookup. A direct consequence
worth noting: if a component were ever referenced by a recipe
after being deleted from the catalog (which the safe-to-delete check is meant to prevent, but
whose exact scope — current recipe content only, or also frozen history entries — isn't fully
specified), there would be nothing on the recipe itself
to fall back on for that component's name or version — the reference would resolve to nothing.
Every
version, current and superseded, is a frozen snapshot retrievable via
`GET /v1/image_recipes/{id}/history/{version}` — a path segment. `status` is derived (see
3.1.1.9.2) — including `degraded` when a referenced component is obsoleted, with
`status_reasons` carrying one standard `{code, message, more_info, target}` entry per
obsolete component: `target` is the machine-readable component reference, `message` names
it. Structured `warnings[]`
(`{code, message, more_info, target}`) carry advisories such as `component_deprecated` and
`base_image_family_changed` — though `warnings[]` is not actually in the
schema, so this mechanism needs reconciling with what does exist. Recipe→images
discovery is entirely via `GET /v1/images?source_image_recipe.id={id}[&source_image_recipe.version={hash}]`
(see 3.1.1.4.3) — filtering by `id` alone returns images built from any version of the
recipe; adding `&source_image_recipe.version={hash}` narrows to one specific version.

**Storage**: Kubernetes Custom Resource (`Kind: ImageRecipe`, `apiVersion:
image.vpc.cloud.ibm.com/v1`) in the regional cluster,
labeled by owning account; history entries as sub-resources reconciled by the same Catalog
Controller

**On `DELETE /image_recipes/{id}`**: deleting a recipe
removes its history. Provenance on already-built images is preserved via the image's own
`source_image_recipe` reference field, which persists independently of the recipe's own
lifecycle. The spec's earlier "recipe manifest" phrasing — an undefined term that never
appeared elsewhere in the spec — has since been removed rather than clarified.

##### Is `ImageRecipeHistory` a distinct schema at all?

`ImageRecipeHistory`'s field list is identical to `ImageRecipe`'s, field for
field (same fields listed above, including `images[]`, `superseded_at`, `status`,
`status_reasons`). This raises the question of whether a separate schema is
warranted at all, versus history entries simply being `ImageRecipe` objects at a past content
state — worth resolving rather than maintaining two identically-shaped schemas.

##### Recipe History Entry (`ImageRecipeHistory`)

```json
{
  "version": "9a11f0c2",
  "name": "ml-training",
  "base_image": { "id": "r006-72b2...", "name": "ibm-redhat-9-4-minimal-amd64-3" },
  "image_components": [
    { "id": "r006-drv560", "href": "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-drv560" }
  ],
  "created_at": "2026-05-11T09:00:00Z",
  "superseded_at": "2026-07-10T14:22:00Z"
}
```

Field list: `base_image`, `created_at`,
`image_components`, `name`, `superseded_at`, `version`. There is no `images` field here
either — discovering which images were built from this specific version is via
`GET /v1/images?source_image_recipe.id={id}&source_image_recipe.version={hash}` (see 3.1.1.4.3), the same
mechanism used for the current version.
`created_at` here is the version's own creation instant
(when that content became current), distinct from `superseded_at` (when it stopped being
current) — consistent with the version-lifetime-interval framing used elsewhere (3.1.1.5).

##### 3.1.1.4.3 IS.Image Extensions and the Image Job Resource

**IS.Image extensions** — two nullable fields plus a back-reference (all additive — null
for non-recipe images); no new image status values — `partially_available` is a pre-existing
status Image Builder simply makes use of (see 3.1.1.9.3a):

```json
{
  "id": "r006-img-abc123",
  "status": "available",
  "source_image_recipe": { "id": "r134-aaa111", "name": "ml-training", "version": "c7d2a9e3", "href": "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r134-aaa111/history/c7d2a9e3" },
  "inventory": {
    "created_at": "2026-05-11T09:28:00Z",
    "href": "https://us-south.iaas.cloud.ibm.com/v1/images/r006-img-abc123/inventory",
    "packages": [ { "name": "nvidia-data-center-driver", "version": "x.y.z" } ],
    "resource_type": "inventory"
  },
  "image_job": { "id": "r006-job-7f2c1a...", "href": "https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-job-7f2c1a..." }
}
```

**`source_image_recipe` is present whenever the image status is not `pending`** — not just
`available`. This covers every status past the initial creation window: `partially_available`,
`available` itself, and any later state the image reaches (`deprecated`, `obsolete`,
`deleting`), since the field is written once, early, and never removed afterward.

The inventory is an extracted name/version list, derived from the Packer plugin's
scan output during the build, and is served at `GET /v1/images/{id}/inventory`, with response
fields `created_at`, `href`, `packages`, and `resource_type`. Per current guidance from the
ITS/CISO teams, the
full original scan document is not retained or exposed anywhere, including internally —
only the extracted list. `image_job` carries both `id` and `href` (matching the shape already
used by `source_image_recipe` on this same object) so the job identifier is directly readable without
parsing it out of the href. See 5.2 for alternatives considered.

**The image job — a first-class top-level regional resource** (`POST/GET/PATCH/DELETE
/v1/image/jobs`). The open question about future job types is now answered: `ImageJob` and
`ImageJobPrototype` carry a required `type` field (currently a single-value enum, `build`,
explicitly flagged as expandable), and the schema is built around a `oneOf`: `ImageJob` merges
common fields (`id`, `href`, `status`, `created_at`, etc.) with `type` and a `oneOf: [ImageJobBuild]`,
where `ImageJobBuild` holds everything specific to building an image —
`image_recipe`, `image`, `image_template`, `timeout`, `worker`,
`delete_image_job_on_image_delete`. The same split applies to the request
(`ImageJobPrototype`/`ImageJobBuildPrototype`) and the patch body
(`ImageJobPatch`/`ImageJobBuildPatch`). A future job type would define its own
`ImageJob<Type>` branch with its own required fields — it would not need `image_recipe`,
`worker`, or `timeout` at all — while `type`, `id`, `status`, and the other common fields stay
shared across every job type. `ImageJobWorkerProfile` uses the same pattern in miniature: a
`oneOf: [InstanceProfileReference]` wrapper, rather than a plain reference, leaving room for a
worker profile to be expressed some other way in the future too.

```json
{
  "id": "r006-job-7f2c1a...",
  "href": "https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-job-7f2c1a...",
  "resource_type": "image_job",
  "type": "build",
  "name": "ml-training-build-1",
  "status": "running",
  "status_reasons": [],
  "image": { "id": "r006-72b27b5c...", "name": "ml-image",
             "href": "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c...",
             "crn": "crn:v1:bluemix:public:is:us-south:a/{account}::image:r006-72b27b5c..." },
  "image_recipe": { "id": "r134-aaa111", "name": "ml-training", "version": "c7d2a9e3", "href": "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r134-aaa111/history/c7d2a9e3" },
  "worker": { "instance": { "id": "0717-vsi-9f3a...", "name": "packer-vpc-vsi-build-7c21",
              "href": "https://us-south.iaas.cloud.ibm.com/v1/instances/0717-vsi-9f3a..." },
              "profile": { "name": "cxf-4x8" } },
  "timeout": 120,
  "delete_image_job_on_image_delete": true,
  "created_at": "2026-05-11T08:58:40Z",
  "started_at": "2026-05-11T09:01:12Z",
  "completed_at": null
}
```

No `crn` field exists on the
image job resource. `resource_type` exists (exact enum value illustrative). There is no
separate top-level `resource_group` on the job itself — only `image_template.resource_group`,
for the resulting image.

**`delete_image_job_on_image_delete`** is an optional, customer-configurable boolean (default
`true`)
controlling cascade behavior: when `true`, deleting the produced image also deletes the job
record that created it. Set at job creation and echoed back in the response.

Key contracts: `status` is the coarse, monotonic enum
`queued → running → succeeded | failed | deleting`. `running` persists through
zone distribution — the job reaches `succeeded` only once distribution fully completes, since
the plugin's own capture step blocks on this before the Tekton task finishes (see 3.1.1.3.3).
`image` is present once the image has been created — an
`ImageReference` — which happens during Task 1's Capturing stage, via the plugin's own
`ImageBySourceVolume` call (see 3.1.1.3.3, 3.1.1.3.4, 3.1.1.9.3a); `worker` is an object
wrapping three things: `instance` (a
reference to
the build instance — standard optional `deleted` sub-object after cleanup, never null),
`subnet` (the subnet the worker runs in, required at job creation), and
`profile` (the resolved build profile, defaulting to `cxf-4x8`);
`timeout` (minutes, 30–480 range, default 120) bounds the whole of Task 1, including the
plugin's blocking wait; `name` is the only PATCH-able field; `DELETE` on a non-terminal job
(`queued`/`running`) cancels the job and deletes the image being built, where one exists;
before the plugin's call has run, there is no image yet to delete. On a
terminal job (`succeeded`/`failed`), `DELETE` only removes the historical record and an
`available` image is explicitly unaffected. Strictly one job per image and one image
per job; no automated retries — re-running a build is a new job.

**Optional image template**: the request body may include an `image_template` object
specifying the resulting image's `name` and `resource_group`. If omitted, the name is
generated (recipe-name prefix + random suffix) and the resource group is inherited from the
recipe's. Image Builder passes this prototype data through to the plugin's own
`ImageBySourceVolume` call as part of the generated build configuration.

**Job retention**: the image job record is retained for as long as its image exists, and is
cleaned up when the image is deleted, via a single gRPC call, image service → Image Builder,
on deletion through the ordinary image APIs — Image Builder cleans up its own job record
itself, rather than the image service reaching in directly. Image
deletion is intended to cascade to clean up the job record and the inventory file together,
not as
separate, independently-timed cleanups (see OI-RIOS).

#### 3.1.1.5 Versioning Strategy

##### The model: one mutable resource per immutable id, content-hash versions, full history

A recipe is a single resource whose **id never changes**, and whose **content may**. Every
content change — `image_components` (whole-array replacement, no element-level merge),
`base_image` change, and `name` change — is fully re-validated
and mints a new
**content-hash version** on the same resource. Only
`deprecation_at` or `obsolescence_at` changes are exempt from minting a new
version — `name` is not on that exemption list. The superseded content becomes a frozen,
addressable snapshot in `recipe_history`, stamped with `superseded_at`.

| Property | Rule |
| --- | --- |
| Identity | `id` immutable; `name` mutable, unique per account+region; `(id, version)` = permanent content coordinate |
| Version | Content hash — minted per distinct content state; identical content vs. the _current_ version → `409` duplicate (an idempotent PATCH — identical to the recipe's current state — is a silent no-op, not an error) |
| History | `GET /v1/image_recipes/{id}/history` (list every version, current included, most-recent first) or `.../history/{version}` (one specific version — current or superseded — a plain path parameter) — every version retrievable forever (until recipe deletion) |
| Rollback | Legal — restoring a historical version's content mints a version (possibly re-reaching a prior hash) |
| Pinned builds | `POST /v1/image/jobs` may pin `image_recipe.version` — builds exact historical content, gated only by the recipe's status |
| Provenance | An image's recorded `source_image_recipe` `(id, version)` resolves via history for the life of the recipe |

Lifecycle dates (`deprecation_at` / `obsolescence_at`) are metadata, not content — changing
them never mints a version.

##### Example

```text
Recipe "ml-training"   version 9a11... -> 3c07... -> c7d2...    (content changes on ONE resource;
                        history: [9a11..., 3c07...]              prior versions frozen in history)

A new, separate recipe (e.g. "ml-training-cudnn") is simply POST /v1/image_recipes with the desired
content -- its own name, its own initial hash. No derivation API or lineage field.
```

##### Version formats

| Entity | Format | Mutable |
| --- | --- | --- |
| Recipe version | Content hash (e.g. `c7d2a9e3`) | Versions immutable; resource content mutable |
| Component | Vendor-native `version`; each version = distinct immutable catalog entry | No |
| Per-build record | Recorded per image at build time | Never |

---

#### 3.1.1.6 Build Execution Flow

##### End-to-End Build Flow

1. **Customer browses component catalog** — `GET /v1/image_components` →
   nvidia-data-center-driver,
   nvidia-cuda-toolkit, amd-rocm, etc.
2. **Customer creates recipe** — `POST /v1/image_recipes` with `name`, `base_image`,
   `image_components[]`
   → full validation (dependencies, component status, component count) → `201`, `version`
   = content hash (e.g. `9a11f0c2`).
3. **Customer triggers build** — `POST /v1/image/jobs { image_recipe: {id},
   timeout }` (image identity service-derived, no prototype) → pre-dispatch validation
   → image job created, `status: queued` (no
   image exists — the image is the job's output) → returns the job immediately (`201`).
4. **Control plane — HCL generation + dispatch (above the cloud)** — read recipe frozen
   content + component records (CRDs) → generate packer.hcl (source block with
   `vsi_base_image_id`, resolved profile, subnet; install provisioner with exact pinned
   packages, dependency order; validate provisioner from schema, no scripts) → `packer
   validate` → `hcl_sha256` → stage to COS → `CreatePipelineRun` (CD REST) → image_job:
   `queued` (PipelineRun created).
5. **Task 1 — Packer Execute** — fetch packer.hcl → verify `hcl_sha256` → provisioning →
   installing → validating → capturing the boot volume. During Capturing, **the plugin
   calls `POST /v1/images` itself** ("image from volume") — this is where the image is
   actually created, with real content, in
   one step. The plugin then polls and blocks until the image reaches `available` before
   the task completes (`job: running`
   throughout; internal stage
   tracking only — `status_reasons` stays empty to the customer, see 3.1.1.9.3a). See
   `OI-PACKER-EXIT` for a narrow edge case in this polling logic.
6. **Task 2 — Post-process · Cleanup task (always)** — inventory extracted; AT event;
   cleanup deletes VSI/key/FIP/SG
   (billing stops). Zone distribution across all zones has already completed automatically
   by this point, inherited from Task 1's blocking wait -> `available`, `job: succeeded`.
7. **Customer polls for completion** — `GET /v1/image/jobs/{id}` → `succeeded` · image ref ·
   `started_at`/`completed_at` · `timeout`; `GET /v1/images/{image_id}` → `available` ·
   inventory reference.

##### OS Support

Any VPC stock image can be used as a recipe's `base_image` — there is no API-level
restriction on which one, beyond it being a stock image. Customers select from the full VPC
stock image catalog (`GET /v1/images?visibility=public`). In practice, `allowed_use.image`
(3.1.1.3.6a) is what makes this workable: since a recipe requires at least 1 component
(3.1.1.7), and every Phase 1 component's `allowed_use.image` expression only matches RHEL 8/9
or Ubuntu 20.04/22.04, a recipe built against a base image outside that set can never satisfy
the minimum-component requirement — every component addition fails
`component_os_incompatible`. This is a side effect of today's component catalog, not a
designed allowlist: if a component targeting a different OS were added to the catalog later,
this protection would need to be revisited. See 5.2 for the fuller history of this decision.

The OS family determines the package manager and inventory method used by the pipeline:

| OS family | Package manager | Install | Inventory (software list) |
| --- | --- | --- | --- |
| RHEL 8 | `dnf` / `yum` | `dnf install` | `rpm -qa --queryformat` |
| RHEL 9 | `dnf` | `dnf install` | `rpm -qa --queryformat` |
| Ubuntu 20.04 | `apt-get` | `apt-get install` | `dpkg-query -W` |
| Ubuntu 22.04 | `apt-get` | `apt-get install` | `dpkg-query -W` |

---

#### 3.1.1.7 Component Catalog

##### Phase 1 Component Catalog — 10 Components

**Tracked as `OI-18`**: no documented rationale exists
for why this specific 10-component NVIDIA/AMD/Intel set was chosen. Related to the OS support
list (2.8.3.2/3.1.1.2, tracked separately as `OI-OS-VERSIONS`).

**All Phase 1 components install via precompiled packages**, with no GPU hardware required
on the build VSI.

| Vendor family | Component names |
| --- | --- |
| **NVIDIA** | nvidia-data-center-driver, nvidia-cuda-toolkit, nvidia-cudnn, nvidia-nccl |
| **AMD** | amd-gpu-driver, amd-rocm, amd-miopen |
| **Intel** | intel-gaudi-driver, intel-synapse-ai, intel-oneapi |

**Component characteristics**:

- Free/open-source (no licensing)
- OS support: RHEL 8/9, Ubuntu 20.04/22.04 — see 3.1.1.3.6a for the hardware/OS
  compatibility discussion
- Dependency declarations: `dependencies` is a flat array — there
  is no `.all_of` sub-property, and it is deliberately generic rather than a more specific
  name — today it represents component-to-component dependencies, but the field is meant to
  extend to other dependency types in the future.
  `dependencies`
  is itself the array of component references (max 19, since a recipe caps at 20 components
  total), all of which are required
  in the same recipe; install order is derived from the graph by toposort. There is no
  alternative-prerequisites operator (`any_of` or similar) in Phase 1.
- **Declarative build knowledge**: `software` and `allowed_use` are what the HCL generator
  actually reads to produce install and validate provisioners; the full public `ImageComponent`
  schema, including the `software` `oneOf` and the `allowed_use.image` CEL mechanism, is
  described in 3.1.1.4.1.
- Binaries stored in Artifactory: `vpc-image-builder-components` repository (binaries only)
- Metadata modeled as an `ImageComponent` CRD (`image.vpc.cloud.ibm.com/v1`)
- Each component version = a distinct immutable catalog entry with derived lifecycle status

##### Validation Rules

1. **Vendor exclusivity**: no mixed-vendor-recipe check (NVIDIA / AMD / Intel) exists — see 5.2
   for that history
2. **Dependency satisfaction**: e.g. cuda-toolkit requires nvidia-data-center-driver
   (a flat, top-level `dependencies` array — all listed components required, no
   `all_of`/`any_of` operator; install order derived by toposort)
3. **Hardware compatibility**: dropped from Phase 1, no automated check, no replacement built —
   see 3.1.1.3.6a for the accepted tradeoff
4. **OS compatibility**: each component's `allowed_use.image` CEL expression is
   validated against the recipe's `base_image` — see 3.1.1.3.6a
5. **Component count**: 1–20 per recipe

---

#### 3.1.1.8 IAM Integration

##### Resource Hierarchy

```text
Account
  |-- Resource Group
        |-- Recipe (image-recipe resource)
        |-- Image job (image-job resource -- build trigger + record; subnet and
        |   profile specified per-job, no account-level settings resource)
        |-- Image (is.image resource -- the job's output)
```

##### IAM Actions (Phase 1)

```text
is.image-component.image-component.read           (public -- all authenticated users)
is.image-component.image-component.list           (public -- all authenticated users)
is.image-component.image-component.deprecate      (authorizes the deprecate action only,
                                                     distinct from .update)
is.image-component.image-component.obsolete       (authorizes the obsolete action only,
                                                     distinct from .update)
is.image-recipe.image-recipe.create
is.image-recipe.image-recipe.read
is.image-recipe.image-recipe.update
is.image-recipe.image-recipe.delete
is.image-recipe.image-recipe.list
is.image-recipe.image-recipe.deprecate   (authorizes the deprecate action only,
                                           distinct from .update)
is.image-recipe.image-recipe.obsolete   (authorizes the obsolete action only,
                                          distinct from .update)
is.image-recipe.image-recipe.operate   (required for image job creation, since the
                                         job operates against the named recipe)
is.image.image-job.create   (creates an image job; is.image.image.create is invoked
                              internally by the service at job creation time, using the
                              caller's optional image template if provided; also requires
                              is.image-recipe.image-recipe.operate, since the job operates
                              against the named recipe)
is.image.image-job.read / list / update / delete
is.image.image.read     (returns source_image_recipe, inventory + image_job back-reference)
```

##### Authorization Matrix

| API | Action | Scope | Notes |
| --- | --- | --- | --- |
| `GET /v1/image_recipes` | `is.image-recipe.image-recipe.read` **and** `.list` (list endpoints require both actions, not `.list` alone) | Account | List recipes; sorted `created_at` desc, ties by ascending `name` |
| `POST /v1/image_recipes` | `is.image-recipe.image-recipe.create` | Resource group | Base image + components |
| `PATCH /v1/image_recipes/{id}` | `is.image-recipe.image-recipe.update` | Recipe | Content changes (mint version) + lifecycle dates; name immutable |
| `POST /v1/image_recipes/{id}/deprecate` | `is.image-recipe.image-recipe.deprecate` | Recipe | Dedicated action; not a reuse of `.update` |
| `POST /v1/image_recipes/{id}/obsolete` | `is.image-recipe.image-recipe.obsolete` | Recipe | Dedicated action; not a reuse of `.update` |
| `GET /v1/image_recipes/{id}/history` | `is.image-recipe.image-recipe.list` | Recipe | List every version, current included — uses `.list`, not `.read` |
| `GET /v1/image_recipes/{id}/history/{version}` | `is.image-recipe.image-recipe.read` | Recipe | Retrieve one specific version — uses `.read`, unlike the list endpoint above |
| `GET /v1/image_components` | `is.image-component.image-component.read` **and** `.list` (both required) | Account | List components |
| `POST /v1/image_components/{id}/deprecate` | `is.image-component.image-component.deprecate` | Component | Dedicated action; not a reuse of `.update` |
| `POST /v1/image_components/{id}/obsolete` | `is.image-component.image-component.obsolete` | Component | Dedicated action; not a reuse of `.update` |
| `GET /v1/images/{id}` | `is.image.image.read` | Image | Returns `source_image_recipe`, inventory reference |
| `GET /v1/images/{id}/image_components` | `is.image.image.read` | Image | Lists this image's components, as `ImageComponentReference[]` — lightweight references, not full component objects |
| `GET /v1/image/jobs` | `is.image.image-job.read` **and** `.list` (both required) | Account | List jobs; sorted `created_at` desc, ties by ascending `name` |
| `POST /v1/image/jobs` | `is.image.image-job.create`, `is.image-recipe.image-recipe.operate` (service internally invokes `is.image.image.create` at job creation time, using the caller's optional image template if supplied) | Resource group | One-call build trigger; `worker.subnet` required, `worker.profile` optional |
| `GET /v1/image/jobs/{id}` · `PATCH` · `DELETE` | `is.image.image-job.read` / `.update` / `.delete` | Job | Status / rename / cancel (`DELETE` returns `202`, not `200`/`204`) |
| `GET /v1/images/{id}/inventory` | `is.image.image.read` | Image | Extracted inventory (name/version pairs) |

---

#### 3.1.1.9 Recipe, Component, and Image Lifecycle

##### Overview

Image Builder introduces resource types with independent lifecycles: **components**
(IBM-managed), **recipes** (customer-managed; stock recipes IBM-owned), and **images**
(IS.Image, customer-owned output). All lifecycle statuses are **derived from explicit date
fields** (`deprecation_at`, `obsolescence_at`) rather than stored enum writes — the dates are
the inputs; status is computed at read time.

![Component and Recipe Versioning & Lifecycle](materials/ib-component-recipe-lifecycle.png?raw=true)

---

##### 3.1.1.9.1 Component Lifecycle

Components are IBM-curated catalog entries. Customers cannot create, modify content, or
delete components in Phase 1 — no general-purpose write endpoint exists. **Deprecate and
obsolete are the one narrow exception**: both are public, customer-facing action endpoints,
letting a customer transition a component's lifecycle state directly, even though they still
cannot author or edit a component's actual content. Everything else about component lifecycle
(publicization, cross-account safety checks) remains internal/operator-side.

Components have dedicated action endpoints, mirroring the pattern recipes use:
`POST /image_components/{id}/deprecate` and
`POST /image_components/{id}/obsolete`. These two actions are the only public, customer-facing
write surface on a component — `deprecation_at`/`obsolescence_at` remain directly PATCHable
too, but only internally, by the operator; there is no public PATCH endpoint for components at
all, unlike recipes where both the dedicated actions and direct date PATCHes are public.

**The precondition and response-code pattern matches recipes exactly**, now that
components have their own 6-state model
(`pending`/`available`/`deprecated`/`degraded`/`obsolete`/`deleting`), the same count as the
recipe's:

- `.../deprecate` requires the component to have status `available` and no `deprecation_at`
  already set. Success returns `204`; a component that doesn't exist returns `404`; a
  component not meeting these preconditions (already deprecated or obsolete) returns `409`.
- `.../obsolete` requires status `available` or `deprecated`, `deprecation_at` not set in the
  future (i.e. not still pending a scheduled deprecation), and no `obsolescence_at` already
  set. Same `204`/`404`/`409` response pattern.

**This requires new IAM actions, the same as recipes.** Both use their own dedicated
`.deprecate`/`.obsolete` actions rather than reusing `.update` — components previously had only
read-oriented actions (`is.image-component.image-component.read`/`.list`)
— no write-style action existed to reuse, since customers
had no write capability of any kind before this. Two new actions,
`is.image-component.image-component.deprecate` and `.obsolete`, are introduced specifically to
authorize these two endpoints, matching the recipe pattern exactly.

This replaces relying solely on
raw date PATCHes for what is functionally a state transition, giving operator tooling and any
future admin UI a single, explicit action to call rather than having to construct the correct
date payload by hand each time.

##### Component Statuses (derived)

| Status | Derivation | Behavior |
| --- | --- | --- |
| `pending` | Component is being created, not yet ready | Not usable in any recipe until it transitions to `available` |
| `available` | No dates set, or `deprecation_at` in the future | Usable in new recipes; builds proceed |
| `deprecated` | `deprecation_at` passed, `obsolescence_at` not | Existing recipes build with a warning (`component_deprecated`); cannot be added to new recipes |
| `degraded` | A component listed in this component's own `dependencies` became `obsolete` | Cannot be added to new recipes; status is derived live, so a recipe already referencing it sees the same `degraded` status the next time it's read or validated |
| `obsolete` | `obsolescence_at` passed | Cannot be added; recipes referencing it derive `degraded`; builds referencing it fail pre-dispatch (`component_obsolete`, 409) — no VSI, no billing |
| `deleting` | Deletion in progress | Terminal, transitional; component is being removed from the catalog |

`degraded` mirrors the recipe's own `degraded` status exactly, one level down: a recipe
degrades when a component it references goes obsolete, and now a component degrades when a
component it depends on goes obsolete. The two exits are the same shape too: publish a new
version of the component with an updated dependency (a fresh catalog entry, since
`dependencies` is immutable per 3.1.1.9.1), or let the affected component itself move toward
deprecation.

Emergency retirement sets both dates to now (immediate obsolescence) — either directly via the
dedicated `.../obsolete` action, or via a date PATCH. Restoration is possible
while the entry exists (dates PATCHed by the owner) — deletion is gated behind obsolescence
plus an operator-side cross-account safe-to-delete check (no customer recipe or in-flight build
references the entry).

##### Component Versioning

Each catalog entry represents a specific vendor-native `software.version`, alongside
`software.name` identifying the actual software. Both `software` and `dependencies` are
immutable once the component is
created; getting a new package version, or changing what a component depends on, means
publishing an entirely new component entry, not
updating an existing one. This matters specifically for `dependencies`: if it were mutable,
changing a component's required dependencies after recipes already reference it could silently
invalidate those recipes without the recipes themselves ever being touched. IBM does not auto-update component versions
within an existing recipe. When IBM releases a new version, the prior entry may remain
`available` or receive deprecation dates with customer notification.

##### End-of-Life and End-of-Service Tracking

```json
{
  "name": "nvidia-data-center-driver",
  "version": "12.4.0",
  "status": "deprecated",
  "deprecation_at": "2026-04-01T00:00:00Z",
  "obsolescence_at": "2026-10-01T00:00:00Z"
}
```

- `deprecation_at` — instant from which status derives `deprecated`
- `obsolescence_at` — instant from which status derives `obsolete`; builds referencing the
  component fail from this point
- Replacement guidance is carried in the deprecation warning's `more_info`

---

##### 3.1.1.9.2 Recipe Lifecycle

Recipes are customer-owned, regionally-scoped resources (stock recipes: same resource,
IBM-owned, `visibility: public`). Recipe `status` is **derived** — from the resource's own
lifecycle dates and from the statuses of its referenced components:

###### Recipe Statuses (derived; 6 states)

| Status | Meaning |
| --- | --- |
| `pending` | Content validation in progress |
| `available` | Valid; builds proceed. May carry an advisory when a referenced component is `deprecated` — see 3.1.1.4.2 |
| `degraded` | A referenced component became `obsolete` — builds fail pre-dispatch; `status_reasons` names the component (`component_obsolete` + target). Two exits: swap the component (content change → new version, un-degrades) or set the recipe's own lifecycle dates |
| `deprecated` | `deprecation_at` passed — builds proceed with `recipe_deprecated` warning (current AND pinned historical versions) |
| `obsolete` | `obsolescence_at` passed — ALL builds rejected 409 (current and pinned); deletion now permitted |
| `deleting` | Deletion in progress |

6 values, no `failed` state. A recipe that fails content
validation never gets created at all (rejected at `POST`/`PATCH` time with a `400`),
rather than existing in a `failed` state.

The recipe status enum (6 values, above) is intentionally distinct from the build
job status enum (5 values — see 3.1.1.9.3a); these are separate resources with
separate lifecycles and should not be conflated.

**Why `pending` and `deleting` are retained** (recipe validation is fully
synchronous today, so it's fair to ask why an in-progress state is needed at all): the
governing consideration is not whether _this specific_ operation happens to be fast, but
whether the API follows VPC's general convention for create/delete operations, which is
asynchronous by default. Diverging from that platform-wide pattern for one resource needs a
specific reason that doesn't apply here — recipe creation and deletion may in practice
complete quickly enough that `pending`/`deleting` are rarely or never observed, but the states
exist for consistency with how VPC resources generally behave, not because either operation is
expected to take meaningful time.

###### Lifecycle Dates (customer-PATCHable on their own recipes)

`deprecation_at` and `obsolescence_at` are ordinary PATCHable fields — schedule, accelerate,
extend, cancel (null), and resurrect-after-obsolete are all legal maneuvers, bounded by:
`obsolescence_at` requires `deprecation_at`; ordering (`obsolescence_at >= deprecation_at`);
nulling `deprecation_at` requires nulling both in the same PATCH. Dates never mint a content
version. Stock recipes reject customer date PATCHes (403) — their lifecycle moves only through
IBM's publishing flow.

Recipes have dedicated, public
`POST /image_recipes/{id}/deprecate` and `POST /image_recipes/{id}/obsolete`
action endpoints, matching the same pattern images use
(`POST /images/{id}/deprecate`/`.../obsolete`), in Phase 1. `deprecation_at`/`obsolescence_at`
remain PATCHable directly as
described above (schedule, accelerate, extend, cancel), alongside the dedicated action endpoints.

**Preconditions and responses**: `.../deprecate` requires the recipe to have status
`available` and no `deprecation_at` already set; success returns `204`. `.../obsolete`
requires status `available` or `deprecated`, `deprecation_at` not set in the future, and no
`obsolescence_at` already set; success also returns `204`. Both return `404` if the recipe
doesn't exist, and `409` if the recipe doesn't meet the preconditions. Each has its own
dedicated IAM action — `is.image-recipe.image-recipe.deprecate` and
`.obsolete` respectively — not a reuse of `.update`. Both fire the same AT event,
`is.image-recipe.image-recipe-status.update`, a distinct status-transition event separate from
the plain content-update event ordinary `PATCH` calls emit.

**Why this mechanism is needed, rather than simply deleting a recipe that's no longer wanted**:
hard-deleting a recipe would leave every image already built from it with a dangling
reference — the recipe that produced a given image would no longer be resolvable, which is a
real usability loss for anything auditing or reviewing an existing image, and applies
regardless of whether the recipe is stock or customer-owned. Deprecation and obsolescence let
a recipe be marked unusable for _new_ builds while remaining fully readable by anything that
already references it — the provenance chain survives the recipe's own retirement.

###### Recipe Deletion Constraints

- A build in progress using the recipe → `recipe_in_use` (409)
- Deletion is otherwise permitted; images built from the recipe are unaffected — their
  provenance records persist. Deleting the recipe is the only event that
  removes its version history
- `recipe_in_use` (409) also applies to `PATCH`, specifically to blocking a deprecation
  attempt while a build is in progress using the recipe.

###### Recipe Content Changes and History

All content changes happen on the same resource via `PATCH` — `components` is replaced
as a whole array (no element-level merge, no separate add/remove sub-operation endpoints);
every change is fully re-validated and mints a
content-hash version, with the
prior content frozen into `recipe_history` (`superseded_at`-stamped). History reads never gate
on status (a `deprecated`/`obsolete` recipe's history remains readable — provenance outlives
retirement); only builds gate on status. See 3.1.1.5.

---

##### 3.1.1.9.3 Image Lifecycle

Images produced by Image Builder are standard VPC IS.Image objects. Image Builder adds two
nullable fields (`source_image_recipe`, inventory reference) and one sub-resource (`image_job`), all
read only extensions that do not affect the existing image lifecycle or status model.
`partially_available` is a pre existing status on the base `Image` resource — not something
Image Builder introduces — reflecting the normal multi-zone distribution every VPC image goes
through. There is no `deleted` status on either the image or the image job, `deleting` is
transitional on both, and the resource simply ceases to exist (404 on further access) once
deletion completes, rather than settling into a separate terminal status.

The base `Image` resource's status enum has two further values Image Builder's flow does not
reach: `failed` and `unusable`. Both appear tied to the image **import** path — uploading an
external disk image file, tracked via the image's own `file` property (explicitly scoped to
"the imported image file") and its own `status_reasons`. `ImageBySourceVolume`, which Image
Builder always uses, has no equivalent asynchronous processing step to fail during: the call
either succeeds atomically, producing an image with real content in one step, or fails
outright with no image record created at all (see 3.1.1.3.3). There is no intermediate state
where an image record exists but subsequently fails to finish, the way an uploaded file can
fail parsing after the record is already created. On that basis, `failed`/`unusable` should
not be reachable via Image Builder's own path — though this is inferred from the two calls'
different semantics, not confirmed directly in the spec — tracked as `OI-IMG-STATUS`.

See 3.1.1.9.3a for the image job status model, terminal failure codes, and the
deletion-scenarios table.

###### Image Availability After Build

Once `image.status = available`: the inventory
descriptor (fields: `created_at`, `href`, `packages`, `resource_type` — there is
no `format` field) references the downloadable, extracted
name/version list; the image can be used immediately for VSI provisioning.

###### Image Retention and Deletion

Deleting an image does **not** delete the recipe used to build it. The extracted software
inventory in IBM COS is retained for IBM-internal audit purposes but is no longer
customer-accessible after image deletion. Only the extracted (name/version) information is
retained — the full original scan document is not kept, per current ITS/CISO guidance. This
storage lives in the Images team's existing archive bucket, and its deletion follows that
team's existing deletion strategy rather than a separately invented Image Builder mechanism —
in practice this means the inventory file (like the image job record) is cleaned up
when the image itself is deleted.

###### Build Metadata Write Timing

Internal build metadata is assembled by Task 2, after the plugin's own
`POST /v1/images` call has already created the image in Task 1's Capturing stage (see
3.1.1.3.3, 3.1.1.3.4). The image
comes
into existence with its content already present, via that one call; there is no
window where an image record exists without content — the
metadata and the image's real content arrive together. The
customer-visible provenance surface (exact per-build record) is defined in the API
specification.

###### 3.1.1.9.3a Image Job Status Model

Image job status enum — 5 values, no pollable `cancelling`/`cancelled` states:

| Value | Terminal? | Meaning |
| --- | --- | --- |
| `queued` | No | Job accepted, waiting to be scheduled |
| `running` | No | In progress — persists through the entire build and the subsequent image zone distribution |
| `succeeded` | Yes | Completed successfully |
| `failed` | Yes | Could not complete; customer may retry via a new job |
| `deleting` | No (transitional) | Being deleted — the job record is then removed entirely, not moved to a separate terminal status |

`status_reasons` is capped at 20 entries (`maxItems: 20, minItems: 0`) — sized generously
relative to what's actually populated today,
since only the 6 failure codes below ever appear, and realistically at most one at a time per
job.

**`status_reasons` while `running` stays empty.** None of the pipeline's internal
stage-progress (provisioning → installing → validating → capturing) is surfaced to
the customer while the job is `running`. A customer polling
`GET /v1/image/jobs/{id}` while `running` should expect `status_reasons: []`.

**`status_reasons` on `failed`** — 6 codes, exactly
one category reason per
terminal failure (not a stack of lower-level detail), and the only codes that ever actually
populate `status_reasons`:

| Code | Fires when | `target` |
| --- | --- | --- |
| `failed_worker_provisioning` | Worker instantiation failed (quota, subnet, NACL/SSH) | subnet/profile where applicable |
| `failed_component_installation` | Package install failed | the component |
| `failed_component_validation` | Package-level validation failed on an installed component (exact commands TBD before commitment review, see OI-VALCMD) | the component |
| `failed_image_capture` | Boot volume capture, or the plugin's `ImageBySourceVolume` call (`POST /v1/images`) failed, or the plugin's own polling loop never observed the image reach `available` | — |
| `failed_timeout` | Exceeded `timeout` — bounds the whole of Task 1, including the plugin's blocking wait for `available`, since the task doesn't complete until that point (see 3.1.1.3.3) | last-active-stage component where applicable |
| `internal_error` | Unattributable to customer input | — (carries a support reference ID) |

###### Image and Job Status During Build

| image job `status` | Image | Meaning |
| --- | --- | --- |
| `queued` | Does not exist yet | Job created; PipelineRun created, not yet executing |
| `running` (internal: provisioning → installing → validating — not customer-visible; `status_reasons` stays empty) | Still does not exist | Build executing on the worker in the customer's VPC |
| `running` (internal: capturing — plugin calls `ImageBySourceVolume`) | Created here, with real content, immediately — no placeholder-then-fill sequence | Still Task 1's Capturing stage; the plugin's `POST /v1/images` call has just succeeded |
| `running` (still — internal: capturing, same stage) | `pending` → `partially_available` (a real intermediate state) → `available` | The plugin is polling and blocking, waiting for `available`, before the task can complete — this is why the job stays `running` through this whole window |
| `succeeded` | `available` | Only reached once the plugin's poll observes `available` and the task completes |
| `failed` (one of the 6 codes above) | No image, if failure occurred before the `ImageBySourceVolume` call, since the task cannot complete (and therefore cannot reach a terminal status) while the plugin is still blocked waiting on distribution | The failure lives entirely on the job; `status_reasons` names the stage/component |
| `deleting` | If it exists yet, deleted along with the job — no separate `deleted` status on the image either; it's simply removed | `DELETE /v1/image/jobs/{id}` on a non-terminal job |

###### Deletion scenarios (job x image)

The API spec's stated behavior: _"If the image job has a status of `running`, the job will be
canceled and the image being created will be deleted. If the image job has a status of
`failed`, `queued`, or `succeeded`, only the job is deleted."_ An `available` image is
separately confirmed unaffected. This groups `queued` with the terminal states rather than
with `running` — consistent with the fact that a `queued` job has no image yet to delete in
the first place, since the image
doesn't exist until the plugin's `ImageBySourceVolume` call succeeds, partway through Task 1's
Capturing stage. During that earlier window, `DELETE` simply cancels
the pipeline with nothing to delete.

There is no window where the job is cancellable and the image is independently already
usable — since the
plugin's capture step blocks until `available` before the task completes, `job.status`
cannot reach `succeeded` while the image is still short of `available`; by the
time `succeeded` is observable at all, the image already is too.

| Scenario | Job | Image | Cleanup |
| --- | --- | --- | --- |
| DELETE job while `queued` or early `running` (before the plugin's `ImageBySourceVolume` call, in Task 1's Capturing stage) | → `deleting` → record removed (PipelineRun cancelled) | No image exists yet — nothing to delete | Always runs |
| DELETE job while `running`, after the plugin's call has created the image but before it reaches `available` | → `deleting` → record removed | The image exists but the job hasn't reached `succeeded` yet (plugin still polling) — deletion applies | Always runs |
| DELETE job racing an already-succeeded job | Cannot occur — `succeeded` only fires once the image is already `available`, so there's no window for this race | N/A | N/A |
| DELETE a **terminal** job (`succeeded`/`failed`) | Historical record removed | An `available` image is explicitly unaffected. A `failed` job has no image at all, since failure before the plugin's call means nothing was created, and the plugin's own polling loop blocks the task from completing at all while distribution is still in progress | Already clean |
| DELETE the image after `available` | Job record and inventory file cleaned up together via the deletion-notification gRPC call (image service → Image Builder, Image Builder cleans up its own records — see 3.1.1.3.3, 3.1.1.4.3) | Normal deletion | Cascading cleanup |
| Build-metadata write fails post-capture | `failed` (`internal_error`) after retry window | No image — the write failure occurs as part of the same Task 1/Task 2 sequence that produces it | Runs |

The cleanup task runs on ALL exit paths — success, failure, cancellation — so customer
billing is bounded and no ephemeral resources are stranded.

###### 3.1.1.9.4 Cross-Resource Lifecycle Interactions

| Event | Effect on Recipe | Effect on Image | Customer action required |
| --- | --- | --- | --- |
| Component deprecated | Status stays `available`; advisory surfaces via `status_reasons` | None on existing images | Optional: swap component version (content change) |
| Component obsoleted | Status derives `degraded`; builds blocked | None on existing images | Swap the component (new version, un-degrades) or retire the recipe |
| Base image deprecated | Status stays `available`; same advisory mechanism as component-deprecated | None on existing images | Optional: swap base image (content change) |
| Base image obsoleted | Status derives `degraded`; builds blocked, code `base_image_obsolete` (409) | None on existing images | Swap the base image (new version, un-degrades) or retire the recipe |
| Recipe deleted | — (history removed with it) | Image provenance preserved | None — existing images continue to work |
| Recipe content change | New content-hash version; prior → history | None on existing images | Trigger new build to get updated image |
| Recipe dates set | Status derives `deprecated`/`obsolete` per dates | None on existing images | Builds warn, then stop at obsolescence |
| Image deleted | No effect on recipe | — (mid-build: cancels the build, cleanup runs); image job record and inventory file are also cleaned up | None |
| Build fails | No effect on recipe | `failed`; image job names the stage | Review error, fix, trigger new build |

Base image handling follows the same pattern as component handling by design decision: a
base image is the same kind of build input as a component, so it derives the recipe's
status the same way. `degraded` is computed at read/validate time against the base image's
own (VPC-standard) status, the same way component status is checked. No new proactive
notification mechanism is required. In-flight builds are unaffected either way, since content is
pinned at dispatch time (`packer.hcl` references fixed IDs), so nothing that happens to a
component or base image afterward touches a build already running.

---

###### 3.1.1.9.5 Deprecation Notices and Customer Notifications

IBM communicates component and recipe deprecations through:

1. **API response fields** — derived `status`, `deprecation_at` / `obsolescence_at`, structured
   `warnings[]` (`{code, message, more_info, target}`) on recipes, `status_reasons` on
   degraded recipes naming the offending component
2. **Build responses and logs** — builds using deprecated content carry the warning
3. **IBM Cloud notifications** — standard channels (email, in-console alerts) for component
   EOL events
4. **Image Builder API** — `GET /v1/image_components` includes lifecycle dates so customers
   can audit
   their recipes programmatically

Phase 1 does not include automated image re-build notifications (e.g. "your image contains a
deprecated component") — see 2.4.

### 3.1.2 Customer Experience

#### 3.1.2.1 User Interface (UI)

**UI (GUI)**: Recipe browsing/CRUD (incl. version history), component catalog browsing, image
creation from recipe (including subnet/profile selection at build time), image job
monitoring. **Decision**:
images built from a stock recipe are shown under their own class of
stock image — a new "stock recipe images" tab, alongside the existing stock images and custom images
tabs (see 2.2). Coordination:
VPC UI, UX, Content teams.

#### 3.1.2.2 Command Line Interface (CLI)

```bash
ibmcloud is components / component <ID>
ibmcloud is recipe-create / recipes / recipe <ID> / recipe-update / recipe-delete
ibmcloud is recipe-history <ID>
ibmcloud is image-build-job-create --recipe <ID> [--recipe-version <hash>] --subnet <ID> [--profile <name>]
ibmcloud is image-build-job <ID>
```

#### 3.1.2.3 Terraform and Packer

new resources `ibm_is_recipe`,
`ibm_is_image_build_job`; new data sources
`ibm_is_components`, `ibm_is_component`, `ibm_is_recipes`, `ibm_is_recipe`,
`ibm_is_recipe_history`.

#### 3.1.2.4 Customer-facing Documentation

This proposal introduces customer-visible concepts that do not exist anywhere in VPC's
documentation today: recipes, components, image jobs, and the recipe-lifecycle model
(deprecation/obsolescence, content-hash versioning). New customer-facing documentation pages
are required, not just updates to existing pages, since there is no existing "custom image"
documentation to extend that already covers this concept.

- **Content team engagement**: not yet done as of this proposal. Will be initiated via
  `#docs-iaas-platform` once this design is settled, to scope the new pages and confirm
  ownership and timeline ahead of Beta.
- **Release notes**: an update is needed at GA (and at Beta program launch) announcing the new
  capability, consistent with how other new VPC capabilities are announced (2.6.5).
- **Migration guide**: no update is needed. This proposal introduces new resources; it does
  not change, deprecate, or require migration away from any existing customer-facing API or
  behavior.
- **API change log**: an update is needed, since this proposal adds new customer-facing API
  operations (3.1.3) that the change log's existing convention covers for every VPC API
  addition.
- **Task flows**: new task flows are needed for the core use cases this proposal enables —
  creating a recipe from the component catalog, creating a recipe from a stock recipe (2.2,
  Path 2), triggering and monitoring a build, and managing a recipe's lifecycle
  (deprecating/obsoleting, reading version history). None of these exist today, since the
  underlying capability doesn't exist today. Scoping the exact task-flow set is part of the
  content-team engagement above, not yet finalized as of this proposal.

#### 3.1.2.5 Customer-facing Metrics

None of the three new resources this proposal introduces (`image-recipe`, `image-component`,
the image job) currently carry a customer-facing metric through Metrics Router.

The two quantitative limits a customer can actually hit are handled as ordinary API responses,
not as Metrics-Router-reported quota metrics:

- **20 components per recipe** (per-recipe, not per-account; 2.8.3.1) — surfaces as an
  ordinary API validation error (`recipe_component_limit_exceeded`) at the point a customer
  tries to exceed it, not as a preemptively reported quota metric a customer would need to
  poll or alert on.
- **5 concurrent builds per account per region** (3.2.2) — enforced at pre-dispatch
  (3.1.1.3.1); a customer exceeding it receives a rejection on the `POST /v1/image/jobs` call
  itself, rather than needing to track consumption against the limit via a separate metric.

**Rationale for not introducing a customer-facing metric for either limit**: both are
naturally self-evident to a customer through ordinary use — a recipe's component count is
visible directly on the recipe object, and concurrent build count is directly the number of
non-terminal image jobs the customer's own `GET /v1/image/jobs` calls already show them.
Neither limit is the kind of slowly-accumulating, hard-to-observe consumption (like storage or
outbound data volume) that the Metrics Router pattern is meant to address per the [best
practices
guidance](https://github.ibm.com/genctl/genesis-platform-metrics/wiki/Developer-Handbook#best-practice-to-introduce-new-platform-metrics).

If a genuine customer-facing metric need is identified before commitment review, for instance
if usage patterns during Beta show customers want to track concurrent-build consumption
proactively rather than reactively hitting the limit, it will be added here with the required
detail (name, type, frequency, cardinality, labels) at that time, following the [Resource
Quota Producer Metric
Event](https://github.ibm.com/cloudlab/srb/blob/master/proposals/2334/materials/service_team_requirements.md)
pattern — not assumed or pre-committed to in this proposal.

### 3.1.3 Customer-Facing APIs

Full endpoint, field, and error-code detail lives in the API specification itself
(<https://github.ibm.com/vpc/vpc-spec-artifacts/pull/5844>); this section stays at the level of
what the abstractions are, how they fit together, and which design points are worth a
reviewer's attention rather than re-deriving what swagger already shows.

#### The shape of the API: four resources, one flow

**Component** → **Recipe** → **Image job** → **Image**. A component is a single piece of
curated software (a GPU driver, say) — the smallest building block. A recipe assembles
components onto a base image, declaring what an image _should_ contain. An image job executes
a recipe once, producing a real, customer-owned image as its output. The image itself then
carries a record of the recipe (and inventory) it came from, closing the loop back to where it
started.

Two of these are catalog resources customers browse and reference (`image_components`,
`image_recipes`); one is a build trigger with no independent existence beyond the build it
executes (`image/jobs`); the last is the ordinary VPC image resource, lightly extended.

#### New and changed operations

For the full property-level detail on any of these, see the swagger view linked above.

**Component catalog**: `GET /v1/image_components`, `GET /v1/image_components/{id}`,
`POST /v1/image_components/{id}/deprecate`, `POST /v1/image_components/{id}/obsolete`

**Recipes**: `POST /v1/image_recipes`, `GET /v1/image_recipes`, `GET /v1/image_recipes/{id}`,
`PATCH /v1/image_recipes/{id}`, `DELETE /v1/image_recipes/{id}`,
`GET /v1/image_recipes/{id}/history`, `GET /v1/image_recipes/{id}/history/{version}`,
`POST /v1/image_recipes/{id}/deprecate`, `POST /v1/image_recipes/{id}/obsolete`

**Image jobs**: `POST /v1/image/jobs`, `GET /v1/image/jobs`, `GET /v1/image/jobs/{id}`,
`PATCH /v1/image/jobs/{id}`, `DELETE /v1/image/jobs/{id}`

**Existing image endpoints, extended**: `GET /v1/images` (new `source_image_recipe.id`/
`.version` filters), `GET /v1/images/{id}` (new `source_image_recipe`, inventory reference,
`image_job` back-reference), plus two new sub-resources,
`GET /v1/images/{id}/image_components` and `GET /v1/images/{id}/inventory`

#### Notable design points

- **`software` on a component is a `oneOf`, even though Phase 1 has exactly one branch.** Today
  it's always a plain package identity (`{name, version}`) — but the schema leaves room for a
  future branch describing software some other way, such as a reference to a catalog offering,
  without touching how Phase 1 components are described. See 3.1.1.4.1.
- **`allowed_use.image` is a CEL expression, not a static compatibility list.** A component
  declares which base images it's compatible with as a
  [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md)
  string (e.g. `operating_system.family == 'ubuntu'`), evaluated against the recipe's
  `base_image` at recipe creation. This is what actually catches an incompatible
  component-and-base-image pairing before a build ever starts. See 3.1.1.3.6a.
- **A recipe is one mutable resource with an immutable `id`, not a family of versioned
  resources.** Every content change mints a content-hash `version`; nothing is ever deleted
  until the recipe itself is. `(id, version)` is the permanent coordinate an image's provenance
  resolves through, for the recipe's whole lifetime. See 3.1.1.5.
- **An image job's `image_recipe` reference is a `oneOf`** — a customer can name the recipe to
  build from by `id`, CRN, or `href`, not just `id`. See 3.1.1.4.3.
- **The image job resource itself is `type`-discriminated via a `oneOf`,** currently with one
  branch (`build`). Fields specific to building — `image_recipe`, `worker`, `timeout` — live on
  that branch, not on the job's common fields, so a future job type could exist without
  inheriting any of them. See 3.1.1.4.3.
- **Deprecate and obsolete are dedicated actions on both recipes and components**, not a
  reused generic update — `is.image-recipe.image-recipe.deprecate`/`.obsolete` and the
  component equivalents, each firing a distinct `-status.update` event separate from ordinary
  content changes. This gives operator tooling and any future admin UI an explicit lifecycle
  verb to call, rather than requiring a raw date `PATCH` for what is functionally a state
  transition. See 3.1.1.9.1 and 3.1.1.9.2.
- **The image carries its own provenance, not the other way around.** There's no `images[]`
  field on a recipe — recipe→images discovery is a filtered `GET /v1/images` query instead
  (`source_image_recipe.id`), since a recipe can produce many images over its lifetime and the
  image is the more natural place to ask "where did this come from." See 3.1.1.4.2.

**IS.Image schema changes**:

- `source_image_recipe`, inventory reference fields and a `image_job` back-reference href added to
  the Image response (nullable, backward compatible); no new image status values — Image
  Builder makes use of the pre-existing `partially_available` status
  (see 3.1.1.9.3a).

---

#### 3.1.3.1 Activity Tracker Events

| Resource | Action | Description |
| --- | --- | --- |
| component | `is.image-component.image-component.list` | Components listed |
| component | `is.image-component.image-component.read` | Component retrieved |
| component | `is.image-component.image-component-status.update` | Component deprecated or obsoleted, via either dedicated action endpoint (both fire this same event, distinct from ordinary content updates) |
| recipe | `is.image-recipe.image-recipe.create` | Recipe created |
| recipe | `is.image-recipe.image-recipe.read` | Recipe retrieved |
| recipe | `is.image-recipe.image-recipe.update` | Recipe updated (content change, or lifecycle dates set directly via `PATCH`) |
| recipe | `is.image-recipe.image-recipe-status.update` | Recipe deprecated or obsoleted, via either dedicated action endpoint (both fire this same event, distinct from ordinary content updates) |
| recipe | `is.image-recipe.image-recipe.delete` | Recipe deleted |
| recipe | `is.image-recipe.image-recipe.list` | Recipes listed, including the history list endpoint (`GET .../history`) — uses `.list`, not `.read`; retrieving one specific historical version uses `.read` instead, see 3.1.1.4.2 |
| image-job | `is.image.image-job.create` | Image job created (user-attributed build action) |
| image-job | `is.image.image-job.read` | Job retrieved |
| image-job | `is.image.image-job.list` | Jobs listed |
| image-job | `is.image.image-job.update` | Job renamed |
| image-job | `is.image.image-job.delete` | Job cancelled / record removed (user-attributed) |
| image | `is.image.image.create` | Image created by the service as the job's output (internal POST /images) — attributed to the job's initiator |

Total new AT events: 15. `/deprecate` and `/obsolete` share a single event between them —
`is.image-recipe.image-recipe-status.update` — distinct from the plain content-update event,
even though each endpoint has its own dedicated IAM action.

---

#### 3.1.3.2 Global Search and Tagging (GhoST)

The two new resource types this proposal introduces, `image-recipe` and `image-component`,
along with the image job, are new candidates for GhoST/ReX indexing, not deprecations of
anything already indexed. Since these are net-new resources rather than changes to an
existing one's properties, the relevant engagement with the GhoST team (via
`#global-search-tagging`) and with Offering Management is to decide which properties on each
resource should be visible to GhoST and ReX, not to plan a transition away from anything.
Likely candidates for visibility: a recipe's `name`, `status`, and `visibility`; a
component's `name`, `software`, and `status`; the image job's `status` and `image_recipe`
reference — the same properties a customer would already search or filter on via the API
itself. This engagement has not yet happened and is tracked to occur ahead of commitment
review.

#### 3.1.3.3 IBM Cloud Compliance and Security Center (IBM SCC)

Phase 1 does not give customers any new compliance or security-configuration surface on the
resources it introduces. Recipes and components are declarative build inputs (what software
goes into an image), not security controls in themselves; the compliance-relevant properties
that already exist on VPC resources, encryption at rest, TLS in transit, IAM authorization,
apply unchanged to the images this service produces, since those images are ordinary VPC
custom images once built (see 3.1.1.9.3). No update to the [VPC Security and Compliance
document](https://github.ibm.com/cloudlab/srb/blob/master/architecture/platformIntegration/ibmSecurityConfigCompliance.md)
is anticipated for Phase 1 on that basis; this will be confirmed, and the document updated if
that assessment turns out to be wrong, ahead of commitment review.

#### 3.1.3.4 Customer Logging

Image Builder does not send log messages to the customer beyond what the API responses
themselves already carry. Everything a customer needs to troubleshoot a failed build is
surfaced structurally, through `status_reasons` (3.1.1.9.3a) naming the failing stage and
component, not through a separate customer-facing log stream. No new customer-facing log
message is planned for Phase 1; if one is identified as genuinely necessary during
implementation, it will go through the review and approval process in [Customer-Facing
Logs](https://github.ibm.com/cloudlab/srb/blob/master/architecture/telemetry/logging/customerFacingLogs.md)
before being added here.

### 3.1.4 Internal APIs and Schemas

**Internal APIs**: the public APIs in 3.1.3 are also reachable internally, additionally the
following internal only operations exist:

**Components**:

- `POST /image_components` — creates a component
- `DELETE /image_components/{id}` — deletes a component by ID
- `PATCH /image_components/{id}` — updates a component by ID, including flipping
  `visibility` between `private` and `public`
- `POST /image_components/{id}/deprecate` — deprecates a component by ID
- `POST /image_components/{id}/obsolete` — obsoletes a component by ID

**Recipes**:

- `PATCH /image_recipes/{id}` — updates a recipe by ID

### 3.1.5 Deployment Architecture

**Feature flags**:

- Primary: `is-image-builder-allowlist` (per VPC feature-flag naming standard)
- Granular: Component-level flags for progressive enablement
- Runtime toggle — no restart required
- Dev MZones: controllers **off by default** (replicas 0); enabled per-MZone via a global
  (Share precedent)
- One exception: `GET /v1/images/{id}/image_components` sits outside this flag entirely, since
  it's a read-only extension of the pre-existing image resource and is visible in both public
  and internal exposures from the start

**Backward compatibility**:

- Existing image creation unaffected
- `source_image_recipe`, inventory reference fields are nullable and `image_job` returns 404 for
  non-recipe images — null/404 for all existing flows
- API clients continue working without changes

**Deployment order**: No hard constraints — independent components, soft dependencies
validated at runtime

**Rollback**: Instant via feature flag disable

**Validation:** Automated test suite (including the synthetic-component CI test — a fixture
component that never existed must generate correct HCL purely from its schema, enforcing the
no-per-component-code property)

**Phased Rollout:** Region-by-region with monitoring (Beta: WDC, DAL, FRA). A CD
service discontinuation notice for `us-east` (12-Feb-2027) surfaced during design review —
tracked separately as OI-CD, since it affects the underlying Tekton service regardless of
region priority here and needs its own resolution with the CD team.

**Feature API PR:** <https://github.ibm.com/vpc/vpc-spec-artifacts/pull/5844>

---

## 3.2 Architectural Considerations and Impact

### 3.2.1 Security

**Build Isolation**: Each build provisions an independent Flex (`cxf-4x8`) VSI in the
customer's VPC zone. The plugin creates an ephemeral security group: inbound TCP/22 from the
Tekton cluster egress CIDR only (not `0.0.0.0/0`); outbound all traffic (required for
Artifactory downloads on port 443). The SG is deleted unconditionally after the build — the
customer never configures SG rules, unlike Azure VM Image Builder's manual NSG requirement.

**Network Isolation**: Tekton pipelines run in IBM Continuous Delivery, no customer access to
the orchestration layer. Build VSI is in the customer's VPC subnet, visible in their
dashboard. Outbound from build VSI: port 443 to Artifactory only. No inbound except Tekton via
floating IP during build.

**Data Isolation**: Build artifacts (extracted inventory, packer.hcl; the custom image itself
is created directly via the plugin's
own `POST /v1/images` call, not staged through COS) are stored in
IBM-owned COS (the Images team's existing archive bucket for the inventory
file — see 3.1.1.9.3) — customer cannot access directly. Build metadata is
account-scoped in IS.Image. Customer data never exposed across accounts. All ephemeral
resources deleted after build.

**Data Classification**: Per-artifact, for Phase 1:

- `packer.hcl` — fully IBM-generated from component schema (`software.name`/`software.version`,
  `dependencies`),
  the recipe's `base_image` reference, and the job's `profile`; no
  customer-authored content, no secrets (injected at runtime, never templated in). Phase 1 has
  no customer script execution, so this artifact contains no customer-originated data.
- extracted inventory — a name/version package list derived from the Packer plugin's
  scan of installed IBM-curated packages, plus recipe provenance (`image_recipe.id`, `image_recipe.version`,
  and the customer-chosen `image_recipe.name` string). The recipe name is the only customer-authored
  value present.
- **The custom image itself** (created via the plugin's `POST /v1/images` call) — base stock
  OS image plus
  installed GPU driver packages. The build VSI is freshly provisioned per build with no
  customer volumes attached, so the captured content carries no pre-existing customer data.
  This item's classification determination needs re-scoping given it's not a COS-resident
  artifact under Image Builder's own storage — see below.

An explicit classification determination for each of these (packer.hcl, the inventory
file, and the custom image itself) — to confirm the above against IBM's formal data classification
taxonomy — is
outstanding and needs security/compliance sign-off before commitment review; tracked as
**OI-DATACLASS**. This determination is Phase 1-scoped and is not expected to survive unchanged
into Phase 2: customer-authored components (declarative `package_details`)
and the bounded `custom_script` install method (5.2.1) introduce customer-originated content into
the pipeline and its artifacts, which will need its own classification review at that time
rather than assuming Phase 1's determination still holds.

**Template Integrity**: The generated packer.hcl is validated (`packer validate`),
integrity-hashed (`hcl_sha256`), staged per-build, and verified by Task 1 before use — the
executed template is provably byte-identical to what the control plane generated. No secrets
in the template (injected at runtime from Tekton secret properties); no build logic fetched
over the network at build time.

#### 3.2.1.1 FS Cloud Considerations

No components operate with elevated privileges or on sensitive customer data. Communications
encrypted with TLS 1.2+; mutual authentication via PKI Vault certificates; all persisted data
at rest encrypted.

#### 3.2.1.2 Secrets

- `is.image-recipe` Operator API keys (regional; API controller, Status Watch controller, and
  the Tekton pipeline environment property) stored in the Vault/Secrets Manager model — never
  exposed to customers; secrets renamed from earlier "image builder" naming to "image recipe"
- Tekton pipeline key rotation **automated**: a controller pushes the rotated key into the
  Tekton environment property via the Tekton properties API
- TLS certificates for both controllers
- Dispatch uses a direct `CreatePipelineRun` call rather than a webhook
  trigger, which eliminates the webhook-secret-rotation problem entirely (a webhook needs a
  rotating shared secret; a direct API call means IBM owns and rotates its own credential on
  its own terms)
- IBM catalog account credential: Secrets Manager, automated rotation, publish-tooling sole
  consumer, minimal policy scope, AT alerting on non-tooling writes
- No customer credentials used in the build pipeline — customer manages nothing
- No secrets stored in recipes, CRD spec/status fields, or generated templates

#### 3.2.1.3 VPC Control Plane Hardening

The Control Plane Hardening Checklist still needs to be completed for each new
microservice, before commitment review.

#### 3.2.1.4 Significant Change Considerations

The existing IS.Image SSAD must
follow the
major change workflow and be updated to reflect the `source_image_recipe`, inventory field
additions, the image/jobs resource, the image-creation flow (via the plugin's existing
`ImageBySourceVolume` call, with `source_image_recipe` and `image_job` written on afterward via
gRPC — no import kickoff
involved), and its use of the pre-existing `partially_available` status.

#### 3.2.1.5 Identity and Access Management

All public APIs require IAM authentication. See Section 3.1.1.8 for the complete new IAM
action set.

### 3.2.2 Performance, Scalability, and Resource Consumption

Image Builder is designed for asynchronous batch processing, off-path from customer
workloads. **Concurrent builds**: 5 per account per region (quota enforced at pre-dispatch).
**Build duration**: ~25–35 minutes for Phase 1 precompiled_package installs (Flex `cxf-4x8`).
**Control plane impact**: Kubernetes API/CRD storage negligible (<0.1%); IS.Image negligible
(<0.01% traffic); RIAS standard routing overhead.

### 3.2.3 Operational Monitoring

Completed during Service Framework ARCH001b.7 review.

**Key fragilities and alert thresholds**: pre-dispatch validation failure rate >5%;
CreatePipelineRun success rate <95%; Tekton task failure rate >5%; build-VSI SSH success rate
<95%; orphaned resources >0; COS upload success <95%;
`iam_s2s_missing` rate >2%; Artifactory download success
<95%. **`failed_image_import` removed** — no longer a distinct case to alert on (see 3.1.1.9.3a);
any zone-distribution issue now surfaces as an ordinary Task 1 failure, already covered by the
Tekton task failure rate threshold above.

**Planned dashboards**: Build Orchestration, Build VSI Health, Recipe Management, Integration
Health, Customer Experience.

#### 3.2.3.1 Fragility Analysis

Scenarios drawn directly from the architecture in 3.1 and the Packer dependency in 2.5.3:

- **Cross-boundary fragility (hybrid above-the-cloud model, 3.1.1)**: the build VSI is the one
  piece of this architecture that runs outside IBM's own infrastructure, in the customer's
  VPC. A customer editing their own NACL, subnet, or security group configuration, or revoking
  the IAM S2S grant, can silently break a build already in flight, with no IBM-side signal
  that anything changed until the build actually fails (see 2.8.2.4).
- **External dependency fragility (2.5.3)**: HashiCorp Packer version or behavior changes,
  even though pinned in the IBM-built task image, are a real upstream risk. Artifactory
  outages or slow downloads directly stall or fail every in-flight build, since every Phase 1
  component installs from there.
- **Cross-service consistency risk**: the two gRPC calls this design depends on (`OI-RIOS`),
  image service to Image Builder on deletion, and Image Builder to image service to write
  `source_image_recipe`/`image_job`, are the only non-REST surfaces in the design. A failure or
  delay in either leaves the image and the job record briefly, or in a failure case
  persistently, out of sync.
- **Resource exhaustion**: high concurrent build demand can exhaust Tekton pipeline capacity
  (2.8.2.3), and the CD service's `us-east` discontinuation (`OI-CD`) removes capacity from
  one of the three Beta regions on a known date.
- **Cleanup-failure fragility**: the cleanup task is designed to run on every exit path, but if
  it fails to run at all (a genuine bug, a Tekton-level failure), the customer is left with
  orphaned billable resources (VSI, floating IP, security group, SSH key) with no automatic
  second attempt.
- **Silent misattribution**: because the build VSI runs in the customer's account, a failure
  caused by the customer's own environment (a NACL blocking egress, a revoked grant) and a
  failure caused by Image Builder's own pipeline can look identical from the job's
  `status_reasons` alone without careful log correlation.

#### 3.2.3.2 Observability

- **Metrics**: build success/failure rates, build duration percentiles, queue depth, Tekton
  task failure rate, component download success rate, and the other rates listed in 3.2.3.3.
- **Events**: the Activity Tracker events already defined in 3.1.3.1 for every customer action,
  plus three internal operator events (3.1.4): `component.catalog.updated`,
  `build.pipeline_run.created`, `build.status.changed`.
- **Logs**: Tekton task logs (Task 1, Task 2, Cleanup) and the two controllers' own logs
  (Image Builder API controller, Status Watch Controller), routed to IBM Cloud Logs.
- **Agents and forwarders**: standard IBM Cloud Logs/osprey agents run on the Tekton-managed
  pods and the two controllers, all of which are IBM-owned and above the cloud. The one place
  this doesn't apply is the ephemeral build VSI itself: it exists for roughly 25–35 minutes, is
  zone-resident in the customer's own account, and is deleted by the cleanup task, so a
  persistent monitoring agent isn't a fit. Observability from the build VSI instead comes from
  Packer's own capture logs, uploaded to COS before the VSI is torn down.
- **Where this is hardest**: the build VSI window is the genuine gap. Anything that happens on
  that VSI, a package failing to install, a validation check failing, has to be captured into
  the pipeline's own logs and the job's `status_reasons` before cleanup runs; there is no way
  to inspect the VSI itself after the fact, by design.

#### 3.2.3.3 Metrics and Golden Signals

Thresholds already established in this proposal, restated here as the golden signals: build
success rate, `CreatePipelineRun` success rate <95%, Tekton task failure rate >5%, build-VSI
SSH success rate <95%, orphaned resources >0, COS upload success
<95%, `iam_s2s_missing` rate >2%, Artifactory download success <95%, pre-dispatch validation
failure rate >5%. Build duration itself (p50/p95/p99, expected around 25–35 minutes for
Phase 1 precompiled installs per 3.2.2) is a goal-side signal, not just a failure threshold.

As an `is.image-recipe`/`is.image-component` service, provisioning and consumption metrics
will be emitted to OSS EDB. Cataloguing those metrics through the Responsibility Matrix's own
process is tracked for completion ahead of the ARCH001b.7 review, not yet done as of this
proposal. Whether any metric will be produced without a corresponding alert or dashboard, for
instance a finer-grained per-component download metric kept for capacity planning rather than
live alerting, is still to be finalized alongside that same cataloguing work.

#### 3.2.3.4 Events

The Activity Tracker events in 3.1.3.1 are the customer-facing surface: every create, read,
update, delete, list, and the deprecate/obsolete actions on both recipes and components, each
carrying a rich key/value payload usable for troubleshooting and landing automatically in the
observability data lake. The three internal operator events in 3.1.4 cover the catalog and
pipeline side.

Phase 1 has no dedicated billing lifecycle events, since the service itself is free (3.1.1.2,
"Free Service"); the underlying build VSI runtime and image storage are billed and metered
through ordinary, pre-existing VPC billing mechanisms, not something Image Builder needs to
separately emit events for.

#### 3.2.3.5 Logging

Log volume follows the standard IBM Cloud Logs [Logging
Standards](https://github.ibm.com/cloudlab/srb/blob/master/architecture/telemetry/logging/loggingStandards.md):
error level and above only in production. Given how much of this design already routes
through AT events (3.1.4) and the metrics in 3.2.3.3, most alerting is expected to be
event-based or metric-based rather than log-based; whether any log-based alert turns out to be
genuinely necessary, versus fully replaceable by an event or a metric threshold, is expected to
be settled during the ARCH001b.7 review itself rather than pre-decided here.

#### 3.2.3.6 Alerting

Mapped to the fragility list in 3.2.3.1: build failure rate and Tekton task failure rate
(pipeline-level fragility), `iam_s2s_missing` rate (cross-boundary fragility), orphaned
resources greater than zero (cleanup-failure fragility), and Artifactory/COS availability
(external-dependency fragility). Alert routing to TIP versus other channels (e.g. Slack) has
not yet been finalized; this proposal does not currently plan any alert that deliberately stays
outside TIP, but that is subject to the SRE team's own review. Alerts are expected to be caught
and handled by the same CORE IPOPS → Compute SRE model already described for support
escalation in 2.6.3, with SRE sign-off on the concrete alert thresholds and severities still
pending (2.6.3). `failed_image_import` is not a distinct case worth alerting on separately:
since the plugin's
capture step blocks until `available` before the task completes (3.1.1.9.3a), any
zone-distribution issue simply surfaces as an ordinary Task 1 failure, already covered by the
Tekton task failure rate threshold above.

#### 3.2.3.7 Dashboards

The five dashboards already planned (2.6.3, 3.2.3 above) are Build Orchestration, Build VSI
Health, Recipe Management, Integration Health, and Customer Experience, intended to correspond
one-to-one with the alert groupings in 3.2.3.6 so that an alert firing has a matching
dashboard to diagnose from. New resource types (`image-recipe`, `image-component`, the image
job) will need to be registered in Ops Dashboard, since these are new IAM/RC resource types
with no existing representation there. Whether any of the five dashboards additionally needs
to be manually monitored, rather than purely alert-driven, for outage or incident prevention is
expected to be decided during SRE sign-off, alongside the SLO framework already referenced in
2.6.3.

#### 3.2.3.8 Responsibility Matrix

Development leads are aware of the [IaaS Responsibility
Matrix](https://pages.github.ibm.com/cloud-iaas-observability/docs/guidance/observability/matrix/#analysis)
and the duties it documents for engineers on this service. No responsibility in that document
is currently expected to be fulfilled in an alternative fashion; if the ARCH001b.7 review
surfaces one, it will be recorded here at that point rather than assumed in advance.

#### 3.2.3.9 SSAD and Approval Checklist

Two SSADs are impacted: a new SSAD for Image Builder itself (`is.image-recipe`/
`is.image-component`), and the existing IS.Image SSAD, which needs the major-change-workflow
update already described in 3.2.1.4 to reflect the `source_image_recipe`/`image_job`
extensions and the pre-existing `partially_available` status Image Builder makes use of. The
[SSAD
template](https://github.ibm.com/cloud-iaas-observability/docs/blob/main/docs/guidance/observability/ssad.md)'s
`6.2 Monitoring` section will be populated directly from 3.2.3.1 through 3.2.3.7 above. This
service is expected to successfully complete the [ARCH001b.7 Review
Checklist](https://w3.ibm.com/w3publisher/ibm-observability/arch001b7-checklist) before Beta,
consistent with the SRE sign-off and SLO commitments already noted throughout 2.6.3 and this
section.

#### 3.2.3.10 Additional Observability Strategy and Troubleshooting

The main troubleshooting surface for a failed build is already built into the design rather
than added on: `status_reasons` names the failing stage and component (3.1.1.9.3a),
`internal_error` carries its own support-reference ID, `hcl_sha256` lets a support engineer
confirm the exact template that executed, and the extracted inventory records exactly what
was actually installed before any failure. No additional troubleshooting tooling beyond these
is currently planned for Phase 1; this is expected to be revisited if the ARCH001b.7 review or
early Beta usage surfaces a genuine gap.

### 3.2.4 Lifecycle, Metering, and Billing

Phase 1 is completely free — no service charges. Customers billed via standard VPC billing
for build VSI runtime (Flex `cxf-4x8`, ~$0.09/build at 30 min) and custom image storage.
IBM-side: Tekton pipelines run on the CD pay-as-you-go professional plan (~$40/month). No
billing integration in Phase 1 (no Resource Controller metering, no BSS). Image Builder's
storage overhead is
close to a typical image, not meaningfully more — the inventory
file in COS is small text, not a second copy of the image content. Whether any pricing
adjustment is warranted is an open question to work through with the pricing/finance
team, not yet resolved.

### 3.2.5 Compliance

The project does not add new hardware inventory and does not add procedural changes to
operations and maintenance procedures.

### 3.2.6 Faults and Disasters

#### 3.2.6.1 Fault Tolerance

**How component faults are handled, by component:**

- **Image Builder API controller**: stateless, standard VPC controller pattern (3.1.1.1). A
  pod crash or restart loses no state, since all durable state lives in the Kubernetes CRDs
  (recipes, components, history), not in controller memory. Rolling restart recovers in under
  2 minutes.
- **Status Watch Controller**: also stateless. If it crashes mid-poll, it resumes polling
  Tekton `PipelineRun` status on restart with no data loss — the `PipelineRun` itself is the
  source of truth for build progress, not the controller's own memory.
- **Tekton pipeline (Task 1/Task 2/Cleanup)**: Tekton's own pod-level self-healing handles a
  task pod crash. A failure here surfaces as an ordinary job failure to the customer
  (3.1.1.9.3a's failure-code taxonomy), not a silent hang, since the Status Watch Controller's
  polling loop will observe the `PipelineRun`'s own failed state.
- **Build VSI**: ephemeral and disposable by design. If the VSI itself fails mid-build (for
  instance, a VPC-side provisioning issue), the failure surfaces as
  `failed_worker_provisioning` (3.1.1.9.3a) and the cleanup task tears down whatever partial
  resources exist. No retry is automatic; the customer creates a new job.

**Open blast-radius question — post-capture failure**: whether an image already created by
the plugin's `POST /v1/images` call (3.1.1.3.3) is left orphaned if a later step in the same
Task 1/Task 2 sequence fails, or whether failure can only occur before that call succeeds, is
not yet confirmed. The inventory file is preserved in COS for retry either way, and Task 2
makes no separate IS.Image write of its own to retry, so this gap is narrow regardless of how
it resolves — but it is not yet a settled fact, and is called out here rather than assumed.

**Single points of failure:**

- **The two gRPC calls this design depends on** (`OI-RIOS`: image service → Image Builder on
  deletion, and Image Builder → image service to write `source_image_recipe`/`image_job`) are
  the only non-REST surfaces in the design, and each is a real dependency edge: a failure or
  delay in either leaves the image and the job record briefly (or, in a genuine failure case,
  persistently) out of sync. Neither call is on the customer-facing build path itself, so a
  failure here does not block a build from completing successfully — it affects only the
  provenance metadata written onto the image afterward and the job-record cleanup on image
  deletion.
- **No single point of failure exists on the customer-facing build path itself** beyond the
  gRPC dependency above: the API controller, Status Watch Controller, and Tekton pipeline are each
  independently horizontally scaled and stateless, and the ephemeral build VSI is
  per-build and disposable by construction.

#### 3.2.6.2 Business Continuity Disaster Recovery (BCDR)

**Service shutdown and restart**: no new procedure is required beyond the standard VPC
controller pattern this service already follows. Both controllers (API controller, Status
Watch Controller) are stateless and can be shut down and restarted with zero data loss, since
all durable state lives in the regional Kubernetes cluster's CRDs, not in controller memory —
the same property that gives the <2 minute recovery time in 3.2.6.1. Tekton pipelines
self-heal at the pod level through Tekton's own existing restart behavior; no
Image-Builder-specific shutdown/restart procedure is being introduced beyond what Tekton and
the standard controller pattern already provide.

**Point-in-time backup and restore**: recipe and component metadata (the new persistent data
store this proposal introduces) is backed up via the regional Kubernetes cluster's standard
backup/DR tooling — Velero snapshots, the same mechanism already used for other platform CRDs
on that cluster. No separately invented backup mechanism is introduced. Build artifacts
(extracted inventory, `packer.hcl`) are stored in IBM COS and share the Images team's existing
archive-bucket lifecycle and retention policy, rather than a separately invented Image Builder
retention window, since that bucket already has its own tested backup/retention procedures.

**New persistent data store — recovery after outage or failure**: the new persistent store
this proposal introduces is the set of Kubernetes CRDs (`ImageComponent`, `ImageRecipe`, plus
recipe history entries) in the regional cluster.

- **Zonal outage**: the regional Kubernetes cluster's existing multi-zone resilience applies
  unchanged; this proposal introduces no new zonal dependency beyond what the cluster already
  provides for every other CRD-based service running on it.
- **Regional outage**: recipe and component metadata is regenerated from the most recent
  Velero snapshot for that region. There is no cross-region replication of this metadata in
  Phase 1 (consistent with 2.6.3's stated Phase 1 scope), so a regional outage's recovery
  point is bounded by the snapshot cadence, not by continuous replication. Build artifacts in
  COS follow the Images team's own existing regional-outage recovery procedure for that
  archive bucket, unchanged by this proposal.
- **In-flight builds during an outage**: a build in progress at the moment of an outage is not
  specially preserved — it either completes if the outage doesn't affect its specific
  dependencies, or fails and is retried as a new job once service is restored, consistent with
  this proposal's existing no-automatic-retry model (2.8.1.3).

**Recovery objectives**: Recovery Time Objective (RTO) of 4 hours for full service
restoration; Recovery Point Objective (RPO) of 1 hour for recipe/component metadata. No data
loss applies to already-completed images, since a completed image is an ordinary IS.Image
resource governed by IS.Image's own existing BCDR procedures, not a separate
Image-Builder-specific data store.

**Sufficiency of existing procedures**: the existing regional-cluster Velero-backup procedure
and the Images team's existing COS archive-bucket procedure are both judged sufficient for
Phase 1's scope, since this proposal introduces no data store with different durability or
consistency requirements than what those existing procedures were already designed to
protect. No new BCDR procedure specific to Image Builder is being introduced.

### 3.2.7 Hardware Platform Support

This project does not require changes to hardware platform components.

### 3.2.8 Component Level Development

**VPC Image Builder**: API controller, Status Watch Controller, CRD storage. Build pipeline:
Task 1/2/Cleanup. Customer VPC: Flex (`cxf-4x8`) build
VSI. **Coordination**: IS.Image team (the two gRPC calls, `OI-RIOS`), plugin team,
IAM team (CRN registration — critical path), CD/Tekton team (pipeline capacity, and the
`us-east` discontinuation notice — OI-CD), Regional Compute team (Flex profile zone
availability — OI-11), RIAS team.

**IS.Image Service**: add `source_image_recipe`/inventory-reference fields + `image_job`
back-reference (nullable, backward compatible); write the `source_image_recipe`/`image_job` fields onto the
image once it exists, via the gRPC call described in 3.1.4 (not at job creation time — the image
itself does not exist until the plugin's own call creates it during Task 1's Capturing stage);
no new image status values — Image Builder makes use of the pre-existing `partially_available`
status (see 3.1.1.9.3a). No internal job-record write endpoint is needed: job status is
read-derived live from Tekton at query time (3.1.1.3.3), and the extracted inventory is copied
directly into the COS archive bucket by the Tekton task itself (3.1.1.3.4), rather than
written via an IS.Image API call. Coordination: RIOS.

**SDK** (external team, IB-508; Go/Python/Node.js/Java): `CreateRecipe`,
`ListRecipeHistory` (filterable by `recipe_version`), `CreateImageBuildJob` (with recipe/version,
subnet, and optional profile), `GetImageBuildJob`, `UpdateImageBuildJob`.

**RIAS**: new routes `/v1/image_components`, `/v1/image_recipes` (+history), `/v1/image/jobs`,
`/v1/images/{id}/inventory`. Removed:
`/v1/build_configurations` and any account-level settings route (both considered and
superseded — see 3.1.1.2).

**Observability**: new dashboards as above.

**Internal/Operator Surfaces**: publicization (operator), cross-account safe-to-delete
(operator), build status read/cancel (internal support).

**Support Structure**: L1/L2/L3 training on recipe concepts, profile resolution,
cleanup guarantees, IAM S2S grant, provenance and inventory
interpretation. Escalation: L1→L2 recipe/profile/subnet config; L2→L3 build failures, orphan
resources, `iam_s2s_missing`; L3→Engineering IS.Image write failures, Tekton issues.

**Release Engineering**: `vpc-image-builder-components` Artifactory repo (10 GPU driver
components, binaries only, GPG-signed). IBM-built Packer task image on ICR with CVE-rebuild
automation. Catalog publish tooling (schema-validated, approval gates).

**ITS**: recipe builds across all 3 vendor stacks, component compatibility matrix, recipe
versioning/history/rollback, degraded-status derivation, subnet/profile validation, IAM S2S
verification, cleanup on all exit paths,
build metadata write, image job status transitions, end-to-end build (RHEL +
Ubuntu), synthetic-component test.

### Summary of Component Impacts

| Component | Impact type | Risk level | Coordination |
| --- | --- | --- | --- |
| VPC Image Builder | New capability | Low (isolated) | IS.Image, IAM, RIAS, CD/Tekton |
| IS.Image | Schema + sub-resource + gRPC calls | Low | RIOS, IFV team |
| UI (GUI) | New components | Low | VPC UI, UX, Content |
| CLI | New commands | Low | External CLI team |
| Terraform | New resources | Low | External Terraform team |
| SDK | New methods | Low | External SDK team |
| Observability | New dashboards | Low | Observability, SRE |
| Internal/operator surfaces | New endpoints | Low | Operations, SRE |
| RIAS | New routes | Low | RIAS, API team |
| Support | New training | Medium | Support, Content |
| Release Engineering | New Artifactory repo + task image | Low | Release Eng, Security |
| ITS | New test pipeline | Medium | ITS, CI/CD |

### 3.2.9 Network Topology and Segmentation Impact

One new network flow: Tekton Task 1 pods → customer VPC build VSI (floating IP, TCP/22). All
other flows covered by existing VPC network architecture. Public gateway required on the
customer's chosen subnet for Artifactory access. HCL generation and profile defaulting
run in RIAS before Tekton is invoked. No dedicated IBM IKS worker nodes.

---

# 4. Testing

Test automation pipelines must be in place such that when ITS delivers a new stock base
image, it's immediately tested for compatibility with existing recipe definitions and
component installations. Testing for roles/actions must cover all authorization flows in
3.2.1.5, including: secondary user with fine-grained access; unauthorized user rejection;
access-tags-only user; CBR-gated administrator access; missing-S2S-grant rejection
(`iam_s2s_missing`, 409); and a full end-to-end build with valid grant + a valid `subnet` on
the request,
confirming the job transitions `queued → running → succeeded` and the image reaches
`available`.

**Key test scenarios for Phase 1:**

| Scenario | Expected outcome |
| --- | --- |
| Build with NVIDIA components on RHEL 9 | `available`; inventory contains nvidia-data-center-driver + cuda-toolkit packages; package-manager integrity checks passed |
| Build with AMD components on Ubuntu 22.04 | `available`; inventory contains amd-gpu-driver + rocm packages |
| Build with Intel Gaudi components | `available`; inventory contains gaudi-driver + synapse packages |
| Recipe content change (component version) | New content-hash version; prior content in history with `superseded_at` |
| Recipe history timestamp semantics across multiple versions | Each history entry's own `created_at` matches the prior version's `superseded_at`; the current version has `superseded_at: null` and its own `created_at` unaffected by later changes |
| `components` array replaced via PATCH | New version minted; full re-validation ran (there is no add/remove sub-operation — the whole array is replaced) |
| Rollback to historical version content | Legal; version minted; builds reproduce |
| Pinned build (`recipe.version` = historical hash) | Builds exact historical content |
| Duplicate name create | `recipe_name_conflict` 409 |
| Idempotent PATCH (identical to current content) | No-op, no version minted, not an error |
| PATCH name on recipe | 400 (name immutable) |
| Component obsoleted while referenced | Recipe derives `degraded`; build → `component_obsolete` 409; no VSI |
| Build trigger without IAM S2S grant | `iam_s2s_missing` 409, no VSI provisioned |
| Build trigger with `subnet` omitted from the request | Standard `400` missing-required-field rejection |
| Build trigger with a `subnet` that doesn't exist in the account | Standard `404` |
| Build trigger with subnet missing public gateway | `subnet_no_public_gateway` 409, no VSI provisioned |
| Build trigger for a recipe (id, version) with an existing available image | Job builds normally and produces a new, distinct image — no fast-complete or lookup behavior |
| Component with `allowed_use.image` expression evaluating false against the recipe's `base_image` | `component_os_incompatible` 400 at recipe create/PATCH time, before any build is triggered |
| Mixed-vendor recipe (e.g. NVIDIA + AMD components) | Not rejected — no `vendor` field and no vendor-exclusivity mechanism exist, so nothing at the API layer stops this; see 2.8.2.9 for the accepted risk |
| Validation stage failure (validation command non-zero) | job `failed` with `failed_component_validation` naming the component; cleanup runs; **no image exists** — validation runs before capturing, and the image is only created once capturing (incl. the plugin's blocking wait for `available`) succeeds |
| DELETE job mid-build | job `deleting` → record removed; PipelineRun cancelled; cleanup runs; no stranded resources |
| DELETE job racing an already-created image | **Not possible under this design**: `job.status` doesn't reach `succeeded` until the plugin's blocking wait observes `available`, so there's no window where the job is still cancellable while the image is already usable — see 3.1.1.9.3a and the retired decision record in 5.2.1 |
| DELETE recipe in use by pending build | `recipe_in_use` 409 |
| DELETE terminal job | Historical record removed; produced image unaffected |
| Build trigger with a different `subnet` on each job | Each build uses the subnet specified on that job — no cross-job memory or fallback |
| Synthetic component (fixture, never published) | Correct HCL generated purely from schema |

## 4.1 CI Pipelines — Automation

All necessary CI pipelines will be updated to support VPC Image Builder, following standard
VPC CI/CD practices (unit, functional, integration testing, static analysis, build/packaging
automation). The Tekton pipeline definition and the IBM-built Packer task image are themselves
under CI (including CVE-triggered rebuilds). The synthetic-component test enforces the
no-per-component-code property on every change.

---

# 5. Solution Space

## 5.1 Lessons Learned

- Builds must execute in the customer's own VPC — consistent with AWS EC2 Image Builder and
  Azure VM Image Builder, avoiding IBM zone-resident service nodes and NBD device pools
- packer-plugin-ibmcloud handles VSI lifecycle management with guaranteed cleanup on all exit
  paths — eliminates the main operational risk of orphaned resources
- `precompiled_package` for all Phase 1 GPU drivers means the build VSI does not need GPU
  hardware — the Flex profile is sufficient and eliminates the GPU quota prerequisite
- IAM S2S authorization must be created by the customer before first build — a prerequisite,
  not something IBM can auto-create
- Recipe templates add complexity without clear Phase 1 value — customers create recipes
  directly from the component catalog
- Separate versioned-recipe resources (the EC2 model) were rejected in favor of one mutable
  resource per immutable name with content-hash versions and history
- Dedicated hardware-compatibility fields (`hardware_class`/
  `hardware_generation`) do not survive contact with Flex profiles, which have no fixed
  hardware mapping. The mechanism went through several iterations before landing on a CEL
  expression (`allowed_use.image`) that covers OS and architecture compatibility; hardware and
  GPU-manufacturer compatibility remain uncovered, with no mechanism built for that part (see
  3.1.1.3.6a)
- Fields that live on different resources for no reason other than sequencing (when each was
  designed) are worth reconciling before they ship — `subnet` and `profile` were split across
  a settings resource and the recipe for that reason alone, and moving both onto the job where
  they're actually used resolved the inconsistency directly (see 3.1.1.2)

## 5.2 Alternatives Considered

### 5.2.1 Phase 1 Architectural Decisions

---

#### `source_volume` on the Image — Removal Considered, Kept

**Decision**: `source_volume` remains on the image after `source_image_recipe` is added, rather than being
removed or replaced.

**Rationale for considering removal**: the temporary build volume it references is not
customer-created, not customer-usable, and disappears once cleanup runs — its ongoing presence
on the image could read as confusing to a customer who has no other reason to know it existed.

**Rationale for keeping it**: if cleanup ever fails to run for some reason, `source_volume`
remains a useful diagnostic link back to the (possibly still-existing) temporary resource. The
field already exists on every VPC image created via `ImageBySourceVolume`, recipe-built or not,
so removing it only for recipe-built images would be an inconsistency rather than a
simplification.

---

#### New `is.image-builder` Service vs Feature Extension of the Image Service

**Decision**: Deliver Image Builder as a feature extension of the existing image service
(`is.image`), registering two new top-level regional resources — `is.image-recipe` and
`is.image-component` — rather than creating a new `is.image-builder` service.

**Rationale**: IAM's inadequate composite-service support for multiple resources under a
child service. The output artifact is an image — the capability extends the image resource
itself — so a separate service would split ownership of the image write path. Directly reuses
the image service's decade-validated stock-content ownership model.

**Rejected alternative — new `is.image-builder` service**: IAM composite-service
limitations; not a registrable S2S source identity for the resource types the build plane
needs; full new-service registration without a corresponding isolation benefit.

---

#### Customer Account Model — Hybrid Above-the-Cloud

**Decision**: IBM orchestration runs above the cloud. The build VSI is the one exception —
IBM-managed via the plugin and the customer's IAM S2S grant, but it runs zone-resident in the
customer's VPC, billed to the customer.

**Rationale**: Consistent with AWS EC2 Image Builder and Azure VM Image Builder. Eliminates
IBM zone-resident service nodes, NBD device pool management, qemu-nbd dependencies, and
chroot isolation complexity.

**Rejected alternative — IBM service nodes (qemu-nbd)**: requires IBM zone-resident
infrastructure; qemu-nbd cannot install kernel modules or GPU drivers requiring hardware
presence; inconsistent with how all major CSP image builders work.

**Rejected alternative — Containers**: cannot install kernel modules or GPU drivers requiring
hardware presence.

---

#### Recipe Versioning — Mutable Resource with Content-Hash History vs Versioned Resources

**Decision**: One mutable recipe resource per immutable id; content changes mint content-hash
versions; superseded content preserved in an addressable history; `(id, version)` as the
permanent provenance coordinate.

**Rejected alternative — semantic-versioned recipes with structural-change forking**:
multiplies resources for routine edits; semantic version numbers imply human-judged
compatibility semantics the service cannot verify; content hashes make duplicate detection
automatic. The EC2 model was assessed and rejected on the same grounds.

---

#### packer-plugin-ibmcloud vs. Custom Orchestration

**Decision**: Use packer-plugin-ibmcloud (Apache 2.0) inside Tekton Task 1, in an IBM-built
task container image.

**Rejected alternative — Custom Go executor**: would replicate plugin functionality without
community benefit.

**Rejected alternative — upstream hashicorp/packer container image**: would still be layered
on (plugin, CA certs, tooling); Docker Hub pulls carry rate-limit/egress/trust concerns vs
ICR; CVE response would wait on HashiCorp's cadence rather than our own automated rebuilds.

---

#### Generated Provisioners vs. Authored Install/Validate Scripts

**Decision**: The HCL generator (`hclwrite`, programmatic AST construction) emits install and
validate provisioner blocks inline from component schema. Artifactory holds binaries only; no
authored scripts exist anywhere.

**Rationale**: Preserves the dynamic-catalog property, keeps a single source of truth, makes
builds inspectable and hashable, enforced by the synthetic-component CI test.

**Rejected alternative — per-component install.sh/validate.sh in Artifactory**: would break
the dynamic catalog, allow schema/script drift, and fail the synthetic-component test. (Phase
2+ `custom_script` is the sole bounded exception, for software with no package form.)

---

#### precompiled_package as the Install Method for All Phase 1 Components

**Decision**: All 10 Phase 1 GPU driver components install via precompiled packages.
No DKMS kernel module compilation in Phase 1.

**Rationale**: Eliminates the GPU quota prerequisite for builds; the Flex profile is
sufficient for pre-compiled package installs; GPU quota is only needed to boot the resulting
image.

**Rejected alternative — dkms_compilation (Phase 2+)**: requires GPU hardware on the build
VSI and customer GPU quota available before building; deferred to Phase 2+.

---

#### Customer-Owned Recipes (No Recipe Templates)

**Decision**: No recipe templates in Phase 1. Customers create recipes directly, optionally
starting from a stock recipe's content (console pre-population).

**Rationale**: Recipe templates add an extra API resource and onboarding step with limited
value given the Phase 1 catalog's small size.

**Image templates — deferred, consciously**: customer-specified image identity properly
belongs in a future **image template** resource (instance-template pattern). Phase 1 defers
it — no requirement today, and 4Q deliverability wins. Recorded as conscious technical debt,
not omission. **Phase 1 consequence — service-derived image identity**: generated image name
(recipe-prefix + random suffix), resource group inherited from the recipe's, `allowed_use`
and lifecycle dates settable only via image PATCH after creation (acknowledged
forget-to-patch weakness — a named driver for the future template).

---

#### Inventory attached to Image (not Recipe)

**Decision**: The inventory is generated per build and referenced from IS.Image
(fields: `created_at`, `href`, `packages`, `resource_type` — no `format`
field). Endpoint:
`GET /v1/images/{id}/inventory`.

**Rationale**: The inventory describes what was _actually installed_ in a specific image (the
build artifact), not what a recipe _intends_ to install. Different builds of the same recipe
may produce different inventories if package patch versions resolve differently. 1:1 image →
inventory relationship is the correct model.

**History of this decision — it changed twice**. The original design served the full,
plugin-generated CycloneDX scan document as a downloadable file (`GET
/v1/images/{id}/sbom`), with a typed metadata descriptor on the image, and recipe provenance
injected into the file's `metadata.properties`. Per guidance surfaced through discussions with
the ITS and CISO teams, a full software bill of materials cannot be exposed to customers via
API without per-image approval, which isn't practical at Image Builder's scale. The design was
revised: the endpoint now serves only the **extracted name/version list**, not the original
scan document, which is not retained anywhere — including internally — beyond the extracted
list. An earlier draft of this document assumed a `format` field would exist as a
forward-looking parameter for a possible richer format later — confirmed **not** part of the
real schema; whatever future format flexibility is needed will have to be designed
separately, not via this field.

**Rejected alternative — Inventory on Recipe**:

- A recipe can produce multiple images; one inventory per recipe loses per-build precision
- Requested vs actually-installed version divergence makes recipe-level inventory misleading
- Recipe version history was also considered as a substitute for a dedicated endpoint — this
  was rejected because recipe history only records the components a customer _selected_, not
  the full as-installed package inventory (base OS packages, transitive dependencies) that the
  inventory is meant to cover; the two answer different questions

**Rejected — typed IBM document**: an IBM-defined schema would require maintaining a
transform from scanner output and duplicate identity the job/image metadata already carries —
this concern is now largely superseded by the extracted-list decision above, which is itself
a form of typed, IBM-defined shape, adopted for the CISO-compliance reason stated rather than
the reasons originally used to reject it.

---

#### Image Job Creation Model — Top-Level Jobs Resource

**Decision**: builds are triggered by `POST /v1/image/jobs` — a first-class
top-level regional resource; strictly 1:1 job↔image; DELETE = cancel; the image is created by
the plugin's own call during Task 1's Capturing stage, not at job creation time (see 2.2 and
3.1.1.3.3) — a reversal from earlier
design discussion, which had assumed import kickoff. **Note**: earlier drafts of this document
framed the resource as `kind`-discriminated to
generalize to future job types; that framing didn't correspond to the real schema at the time.
The real mechanism has since been confirmed: a required `type` field (currently the
single-value enum `build`, explicitly expandable) plus a `oneOf` split between common job
fields and a build-specific `ImageJobBuild` branch — see 3.1.1.4.3.

**The arc**: the original design tacked build state onto the image (a `waiting` status + a
`image_job` sub-resource) — API review established the job abstraction as required, since
build state on the image is conceptually wrong and prevents retry semantics. Five options were
weighed; the adopted top-level-jobs-resource option dissolves the `waiting` question entirely,
makes jobs listable as build history, and gives cancel a natural verb (DELETE) instead of
overloading image deletion.

**Why the job carries a `image_job` back-reference on the image, and whether it's needed at
all**: raised more than once in review — is this reference actually necessary, given the job
already carries the authoritative reference the other way? Two points settled it. First, it is
not a parent-child relationship: the job and image are related 1-to-1, not hierarchically, and
their lifecycles are independent once the job completes — so keeping an explicit reference,
rather than relying on any implied hierarchy, is the more correct model, not just a
convenience. Second, retention: whether this reference stays meaningful over time depends on
how long the job record itself is kept, which is why job retention is now explicitly tied to
the image's own lifecycle (3.1.1.4.3) — a job that could be deleted independently of its image
would eventually leave this reference dangling, which the chosen retention model avoids.

---

#### Generalized Fast-Complete — Considered, Superseded by "New Image on Rebuild"

**Current design**: every `POST /v1/image/jobs` triggers a real build and produces a new
image — rebuilding from the same recipe, even at the identical `(id, version)`, even a stock
recipe, does not resolve to an existing image. This directly reverses an intermediate
decision (below), following direct API-focal feedback that a customer should always get a
new image from a build request.

**Fast-complete — considered, then superseded**: if an available image already existed for
the requested `(recipe id, version)` — stock or custom — `POST /v1/image/jobs` would
fast-complete: `succeeded` immediately, `image`
pointing at the existing image, no worker provisioned, no billing. Rationale at the time: a
`409` communicates failure for what is, from
the customer's perspective, actually a success — fast-complete does the lookup for them and
returns the answer in-band, and avoids Image Builder adding a new proliferation vector
through its own request path.

**Why fast-complete was not carried forward**: direct feedback from the API focal was that
build requests should behave uniformly — a customer triggering a build expects a build to
happen, not a lookup disguised as one. The "odd, ownership-mismatched side effect" concern
(a "create" verb sometimes behaving as a query) that fast-complete was partly designed to
avoid ended up being the same concern that argued against it.

**A more general alternative was also raised during this discussion**: a "search"-style read
(matching an image by base OS family/version and a components
list, independent of any specific recipe) was proposed as a way to give customers an explicit
query path, separate from the build-trigger endpoint. Not adopted for Phase 1 — additive
scope beyond what Phase 1's catalog
size and use cases currently justify — worth reconsidering if the catalog grows substantially
in a later phase.

**Relationship to competitor recipe-reuse models**: AWS and Azure both support producing many
images from the same recipe/template over time (chiefly scheduled rebuilds). With every
Image Builder rebuild now producing a new image, this is closer to parity with that model
than fast-complete was — see 2.3.2.2.

**Rejected alternative — `409` Conflict**: considered at the same time as fast-complete, for
the same duplicate-build scenario; communicates failure for what could be read as a
customer-perspective success, and forces an extra round-trip the service could do itself.
Superseded along with fast-complete by the "always build" decision above, for the same
uniformity reasoning.

---

#### Artifact-Recovery Machinery for Failed Builds — Rejected

**Decision**: a failed build produces a failed job with no surviving image — the customer
retries by creating a new job. No distinct import-failure category exists; any
zone-distribution issue surfaces as an ordinary Task 1 failure under the existing failure-code
taxonomy (3.1.1.9.3a).

**Rejected alternative — artifact-recovery machinery**: preserving the boot
volume, or exposing a captured artifact for customer-driven re-import, was considered and
rejected because service-side retries would be unbounded and implicit in the customer's
account, and the frequency data doesn't justify the machinery.

---

#### `status` List Filter — Array, Then Single Value

**Decision, current state**: `GET /v1/image_components` and `GET /v1/image_recipes` both
filter on `status` as a single value, not an array. This matches `GET /v1/image/jobs`, which
was always single-value.

**How this decision came about**: raised directly in spec review — recipes and components had
been given an array-based `status` filter while jobs had not, with no documented reason for
the difference. The concern: status filtering isn't common practice across VPC list APIs
generally, since most status values are derived rather than cached, and filtering on a derived
property carries real performance implications; the few existing precedents (e.g.
`privatePathServiceGatewayEndpointGatewayBindingStatus`, `backupPolicyJobStatus`) are
single-value only. Rather than treat the array-based filter as a standing precedent to justify
keeping, or expand it further, both were brought in line with the existing single-value
pattern.

**Rationale for single-value over array**: the array case (matching multiple statuses in one
call) was assessed as a nice-to-have rather than a demonstrated customer need — the real use
case (find things in a given status) is served just as well by a single value. Single-value
is also backward-compatible with a future array if real demand appears later, since a
single-value filter is equivalent to a one-element array from the caller's perspective — the
reverse migration (array to single-value) would have been a breaking change.

---

#### Hardware/OS Compatibility — CEL, Dropped, Reintroduced as a Plain Object, Then CEL Again

**Decision, current state**: `allowed_use` exists on the component, scoped to a single `image`
CEL expression (`{api_version, image}`), using `operating_system.family`,
`operating_system.version`, and `operating_system.architecture` as variables — no standalone
`architecture` field, and no `operating_system` object. Hardware/GPU-manufacturer compatibility
and vendor exclusivity remain dropped,
with no mechanism and no replacement built — see 2.8.2.9 for that residual gap.

**The full history**: `hardware_class`/`hardware_generation` →
`profile_allowed_to_use.instance` → a full CEL-based `allowed_use` (covering hardware, vendor,
and OS together) → dropped entirely, with no replacement built → reintroduced as a plain
`operating_system: {family, version}` object plus a separate `architecture` field, once the
NVIDIA per-OS naming problem showed the OS gap specifically was not a narrow edge case →
rebuilt again as a CEL expression (`allowed_use.image`) once implementing the plain-object
version against the base `AllowedUse` schema surfaced a real conflict: that schema requires
`instance` and `bare_metal_server` as inherited properties, which would have forced those back
onto the component against the team's own decision to keep them off. Extending the base type
to make those optional was judged too large a refactor to do now, so a bespoke
`ImageComponentAllowedUse` schema was built instead — CEL again, but scoped to `image` only,
with no `instance`/`bare_metal_server` inherited or present. Vendor exclusivity (a
`vendor` field plus a `vendor_conflict` validation error) was part of an earlier design and
still does not exist in the real schema. A broader `constraints` wrapper object, meant to
group compatibility-related fields together, and a `type` field (`driver`/`toolkit`/`library`)
for classifying components, were both part of the same earlier design and were dropped along
with it — components have no such wrapper, and `allowed_use` is a single top-level field
instead. A `supported_os`/`supported_operating_systems` list-based approach
was also considered along the way and rejected as redundant with the CEL expression.

**Rationale for having moved to CEL in the first place (historical)**: a
dedicated
`hardware_class` field would duplicate the `family` field VPC
instance profiles already carry. More substantially, a fixed `hardware_generation` field
cannot be correctly populated for Flex profiles, which have no single fixed hardware mapping
across provisioning attempts — a structural blocker, not a minor gap. A CEL expression
evaluated at compatibility-check time sidesteps this.

**Rationale for dropping the whole mechanism, then narrowing rather than fully rebuilding
(intermediate step, since superseded)**:
Phase 1 relies on
IBM-curated stock recipes for NVIDIA (and possibly AMD), where IBM manages hardware
compatibility directly when assembling those recipes — automated enforcement at the API layer
isn't needed for that path. The accepted gap is for customers building a fully custom recipe
from individual components. Once it became clear that OS mismatches specifically are the
normal case for any component with real per-OS packaging differences (not a rare edge case),
a plain `operating_system`/`architecture` pair was reintroduced to close that specific gap
without rebuilding the full CEL mechanism or its hardware/vendor scope — later replaced by the
CEL-based `allowed_use.image` described above, once the plain-object version proved
incompatible with the base schema.

**Rejected alternative — dedicated `hardware_class`/`hardware_generation` fields**: has no
correct answer for Flex profiles specifically, which are the confirmed Phase 1 default.
Retaining the fields would mean either leaving them perpetually null for the default profile
or blocking Flex-profile builds entirely.

**Rejected alternative — extending the base `AllowedUse` schema to make `instance`/
`bare_metal_server` optional**: would have solved the inheritance conflict cleanly and kept
the component's schema consistent with every other schema built on `AllowedUse`, but was
judged too large a refactor to complete on the current timeline. Accepted as a real tradeoff —
the bespoke `ImageComponentAllowedUse` schema risks future divergence from the base type as
the API evolves — rather than a cost-free choice.

---

#### Account-Level Settings Singleton — Considered, Then Superseded by Job-Level Fields

**Decision, current state**: there is no account-level settings resource at all. `subnet`
(required) and `worker.profile` (optional, defaulting to `cxf-4x8`) are specified directly on
`POST /v1/image/jobs`. See 3.1.1.2 for where this lives in the current design.

**Intermediate design, now superseded**: an account-level settings singleton
(`GET`/`PATCH /v1/image/build/settings`, no `POST`, no `DELETE`) held `subnet` per
account+region, following the Code Engine/IKS pattern. `build_configuration` — a dedicated
CRUD resource considered as an alternative to the singleton — was rejected at that stage too:
over-engineered for one optional field, would need customers to create/manage/reference an ID
on every build, and Azure has no equivalent resource.

**How this decision came about**: the original design split subnet and profile across two
different places — `subnet` lived in the account-level settings singleton, while `profile`
lived on the recipe. In review, this split turned out to be an artifact of sequencing rather
than a deliberate choice: the settings API was designed before the image job resource existed
at all, and `profile` was added to the recipe separately, later, without anyone going back to
reconcile the two. Once that became clear, moving both onto the job removed the inconsistency
directly, and matched customers' actual expectation that creating a build job means agreeing
to a worker instance being provisioned in their account — at which point asking for the subnet
that instance runs in isn't an unreasonable thing to ask for at that same moment.

**The one real tradeoff, worked through directly**: an account-level settings API lets a
platform admin configure a subnet once and let developers build without ever thinking about
it — a real convenience, especially in the GUI, where every build would otherwise mean
picking a subnet from a dropdown with no reliable account default to pre-select. The
counter-argument that settled it: customers already do exactly this every time they create an
ordinary VPC instance, with no `/instance/settings` API to fall back on, and instance creation
happens far more often than image builds. Anyone building frequently enough for this to matter
is already automating against the API, where specifying a subnet is a non-issue. And if a
settings-style convenience is ever needed later, it can be reintroduced without breaking
anyone: make the job's `subnet` optional, falling back to settings when omitted, with an
explicit job-level `subnet` always overriding it.

---

#### Catalog Management — Hybrid Ownership vs Full Operator APIs

**Decision**: Catalog CRUD through internal/operator tooling as the owning IBM catalog
account; publicization via an internal `PATCH` flipping the undocumented `visibility` field;
operator-side cross-account
safety checks. No public write endpoint in Phase 1.

**Rationale**: Phase 1's catalog is small and entirely IBM-curated, so a full operator API
would be scope without a demonstrated need. Phase 2 customer components will require adding a
public write endpoint at that point —
this is a deliberate deferral, not a pre-built-but-gated placeholder; keeping Phase 1's
surface strictly internal avoids documenting a customer-facing endpoint that no customer can
actually call.

**Rejected alternative — full operator API management**: new operator-resource registration
scope without a named benefit; a parallel lifecycle implementation duplicating the customer
deprecation code path.

---

### 5.2.2 Alternatives Not Considered for Phase 1

**Customer-Managed Packer Scripts** (ibmcloud-image-builder open-source tool): requires Packer
expertise, no governance or SBOM tracking, customer manages their own build VSI provisioning
and credentials. Managed service with curated components is the core value proposition of VPC
Image Builder.

---

## 5.3 Competitive Offerings

### 5.3.1 Amazon Web Services — EC2 Image Builder

Fully managed service, no additional cost beyond underlying AWS resource usage (build EC2
instance runs in customer account, customer-billed — same model as VPC Image Builder Phase
1). Automated pipeline (build → test → distribute), Systems Manager integration, built-in
security scanning, recipe-based YAML/JSON configuration. Amazon-managed components owned by
the pseudo-account `aws`; no managed-content lifecycle (old versions accumulate). Recipes
always customer-owned — AWS ships no stock recipes. Notably, AWS pipelines are designed to
run repeatedly against the same recipe — including on a schedule — to produce updated AMIs
over time; VPC Image Builder Phase 1 matches this in spirit, since every build against a
recipe produces a new image (see 2.3.2.2).

**Documentation**: <https://docs.aws.amazon.com/imagebuilder/>

### 5.3.2 Microsoft Azure — Azure VM Image Builder

Fully managed, built on HashiCorp Packer. Build VM in customer subscription, customer-billed.
Customer must configure NSG rules for Packer SSH access — VPC Image Builder eliminates this
via the plugin's ephemeral SG. No component concept: customizers are inline scripts or
live-fetched scriptUri references (the drift model our frozen, generated approach avoids).
Also documents rerunning the same image template repeatedly (e.g. monthly, as source
images update) — same reuse-model gap noted for AWS above.

**Documentation**: <https://learn.microsoft.com/en-us/azure/virtual-machines/image-builder-overview>

### 5.3.3 Red Hat — Image Builder (Hybrid/Multi-Cloud)

Open-source (formerly Composer), blueprint-based (TOML) configuration for hybrid/multi-cloud
RHEL images. On-premises and edge suitable.

### 5.3.4 IBM Cloud — ibmcloud-image-builder (Open Source)

Community-maintained CLI using HashiCorp Packer. Requires Packer expertise, customer-managed
VSI provisioning and credential handling. VPC Image Builder provides a managed alternative
with inventory tracking, component governance, and automated cleanup.

**Repository**: <https://github.com/IBM-Cloud/ibmcloud-image-builder>

### 5.3.5 Competitive Analysis Summary

| Capability | AWS EC2 IB | Azure VM IB | VPC Image Builder Phase 1 |
| --- | --- | --- | --- |
| Build runs in customer account | Yes | Yes | Yes |
| Customer billed for build instance | Yes | Yes | Yes |
| Managed ephemeral resource cleanup | Yes | Partial (NSG manual) | Yes (fully automated) |
| Inventory generation | Yes | No | Yes (extracted name/version list, component-level) |
| GPU driver component catalog | No | No | Yes (10 components) |
| Stock recipes (first-class, buildable, lifecycle-managed) | No | No | Yes (unique) |
| Managed-content lifecycle (deprecation/obsolescence dates) | No (version pile-up) | Marketplace pipeline | Yes |
| No GPU quota needed to build | No | No | Yes (precompiled packages, no GPU hardware needed) |
| Recurring/scheduled rebuild from same recipe | Yes | Yes | No (Phase 2 candidate — see 2.3.2.2) |
| IAM S2S (customer creates) | IAM Role | Managed Identity | IAM S2S grant |
| Free service | Yes | Yes | Yes (Phase 1) |

---

# 6. Open Issues

| # | Issue | Impact | Status |
| --- | --- | --- | --- |
| OI-RIOS | Image creation happens in Task 1's Capturing stage (`step_capture_image.go`) — the plugin builds an `imagePrototype` and calls the VPC SDK's `CreateImage` using `ImageBySourceVolume` → `POST /v1/images`. No import-ready gRPC signal is needed, and no separate `ImageByRecipe` creation mechanism exists — the plugin creates the image with real content directly; an internal API for creating images from a recipe was considered and removed. The plugin's capture step polls the image's own status and blocks until `available` before calling `Cleanup()`, so the Tekton task doesn't complete until zone distribution fully finishes; `job.status` stays `running` through zone distribution and reaches `succeeded` only once it's complete. **Two gRPC calls confirmed needed, in opposite directions**: (1) image service → Image Builder, on deletion, so Image Builder can clean up its own job record, since the extracted inventory file itself is written directly to the COS archive bucket by the Tekton task and needs no separate cleanup call; (2) Image Builder → image service, writing the `source_image_recipe` and `image_job` fields onto the image once it exists, since the plugin's generic call has no way to populate either itself. Call (2) currently coincides with `available` in practice — not a structural requirement, but a consequence of the plugin not surfacing the image ID as an accessible interim result before its own blocking wait completes; a candidate for a future Packer/plugin enhancement, not a permanent limitation. `source_volume` remains on the image (considered for removal, kept — see 5.2). | Both calls' mechanism and direction settled; exact payload/trigger conditions for call (2) still to be finalized | **Substantially resolved** — payload details remain |
| OI-PACKER-EXIT | Raised directly in review: the plugin's blocking wait polls specifically for `available` (`step_capture_image.go`, lines 180–197). A narrow but real edge case exists — if `deprecation_at`/`obsolescence_at` is set on the image while it's still `pending` (confirmed possible even mid-build), the image can transition straight to `deprecated` or `obsolete` without the plugin ever observing `available`, leaving the plugin waiting indefinitely for a status it will never see. Confirmed fix direction: update the plugin to treat `partially_available`, `deprecated`, and `obsolete` as valid build-complete signals alongside `available`, not just `available` alone. `partially_available` specifically is safe to exit on for our integration, since Task 1's `Cleanup()` never provisions from the freshly-built image immediately after — that constraint only matters for a Packer user who does. | Fix direction confirmed; implementation and Packer plugin PR still pending | Open — Deepak's team implementing; sync meeting planned with Astha |
| OI-CRN | `image-recipe` (and `image-component`) CRN registration in the service framework | **CRITICAL PATH** — blocks all IAM S2S grant setup and POC work | Open — Image Builder + IAM/service framework |
| OI-11 | Flex profile (`cxf-4x8`) zone availability — whether it's available in every zone across the planned Beta regions. General suitability for Phase 1 precompiled installs is well-established throughout this proposal (see 3.1.1.2, 3.1.1.7) and is not in question; zone availability specifically has not been confirmed | Open — Regional Compute team, coordination pending | Open — Regional Compute team |
| OI-16 | Evaluate COS Data Lake vs. standard COS buckets for build artifact storage | Storage architecture decision required before SA | Open — Image Builder team |
| OI-18 | Stock recipe list not yet confirmed — Phase 1 includes IBM-provided stock recipes covering NVIDIA, AMD ROCm, and Intel Gaudi. Specific recipe combinations, naming (ibm-* prefix at publicization), and component version pinning must be confirmed with the GPU architecture team and product before the commitment review. Related: which specific NVIDIA GPU generations (e.g. V100, A100/H100, B300) are covered by one stock recipe vs. requiring separate recipes per generation is still under discussion with the GPU architecture team, weighed against not wanting an excessive number of stock recipes. | Open — GPU architect + product decision before commitment review | Image Builder team |
| OI-OS-VERSIONS | Raised directly in review: the component catalog's OS targeting (RHEL 8/9, Ubuntu 20.04/22.04 — see 3.1.1.7) has no documented rationale for its specific cutoffs and is not currently flagged as provisional anywhere. Whether newer releases (Ubuntu 24.04/26.04, RHEL 10) are deliberately excluded (e.g. not yet validated) or simply weren't added when this section was drafted is unconfirmed. This is now purely a component-catalog scoping question — since any stock image is a valid `base_image` (3.1.1.6), there's no separate base-image OS cutoff to explain. Same open question applies to the 10-component NVIDIA/AMD/Intel scoping itself (3.1.1.7). | Open — needs confirmation with the team before commitment review | Image Builder team |
| OI-IMG-STATUS | The base `Image` resource's status enum includes `failed`/`unusable`, alongside `available`/`deleting`/`deprecated`/`obsolete`/`partially_available`/`pending` (8 values total). This proposal documents `failed`/`unusable` as unreachable via Image Builder's `ImageBySourceVolume` path (see 3.1.1.9.3), based on the two calls' differing atomicity — not a direct confirmation from the API team | Needed before commitment review, to confirm Image Builder's build-phase failure handling (3.1.1.9.3a) is complete and doesn't need to also account for these two statuses | Open — Image Builder team, needs API team confirmation |
| OI-VALCMD | With `validation_commands` removed, package-level in-build validation (3.1.1.3.6a) is generator-authored rather than a plugin default. The exact commands to run per package manager have not been decided — determining the right commands depends on testing and understanding what's actually suitable for this purpose. Will be added before commitment review. | Does not block inception; must be resolved before commitment review | Open — Image Builder team |
| OI-TPL | HCL template staging: per-build COS object vs git-commit trigger (template history + webhook-secret elimination) | Pipeline design refinement | Open — Image Builder team |
| OI-CAP | Tekton test-environment capacity for staging/integration load | Staging topology | Open — CD/Tekton team |
| OI-DIST | GPU software distribution rights — whether IBM may host vendor driver packages in Artifactory (per-vendor legal review); determines Artifactory contents and possibly the install path | Could force vendor-repo downloads at build time; affects reliability and inventory story | Open — Legal + Image Builder team |
| OI-SRE | Compute SRE does not currently support the image service — support model for the Image Builder capability to be agreed via service framework | Support readiness for GA | Open — SRE + service framework |
| OI-CD | IBM Continuous Delivery discontinuation notice for `us-east` (12-Feb-2027) — `us-east` (WDC) is one of the three Beta regions | Affects Beta regional plan, since the underlying CD Tekton service is what's affected | Open — CD team, needs its own resolution track |
| OI-SCAN | Image signing / deeper scan evidence — not built, not committed for Phase 1, but flagged as a plausible near-future ask. A platform-independent path exists (SCC Workload Protection consuming the inventory, `cosign`-style signing) if pursued | Low urgency, Phase 2+ candidate | Open — Image Builder team, no current owner |
| OI-DATACLASS | Explicit data classification determination for each COS-stored build artifact (`packer.hcl`, the extracted inventory file, the custom image itself) against IBM's formal data classification taxonomy (see 3.2.1). Phase 1 content is assessed as IBM-generated plus a customer-chosen recipe name only; Phase 2 customer components and `custom_script` will need a separate re-review since they introduce customer-originated content | Needed before commitment review; downstream Phase 2 re-review also required | Open — Security/Compliance sign-off pending |
| OI-PIPE | Whether to migrate the build pipeline (Task 1/Task 2/Cleanup) from raw Tekton to OnePipeline, IBM's shared compliance/DevSecOps pipeline platform. OnePipeline offers additional capabilities out of the box — e.g. built-in SBOM generation, secret scanning, and potentially image signing/scan evidence — that would otherwise need to be built or integrated separately. Both platforms ultimately run on the same underlying CD Tekton service, so this is a definitions/dispatch-model decision, not a different execution engine. Open questions include: whether OnePipeline's capabilities (typically scoped to VM image artifacts generally, not any specific format) actually apply to our own artifact type (vs. being scoped to containers/source repos); whether OnePipeline supports a parameterized, per-invocation trigger suited to our per-customer image job model, or only git-triggered flows; and the real migration/rework cost for the existing Task 1/Task 2/Cleanup definitions. This document is intentionally being kept describing the raw-Tekton model through inception review; any adoption of OnePipeline will be reflected here at commitment review instead, once the evaluation concludes. | Could reduce custom-build work for SBOM/scanning/signing if applicable to our artifact type; otherwise migration cost may not be justified. Not required for Phase 1 GA | Open — Image Builder team, evaluation in progress with the CD/Tekton team |
| OI-PUBLISH-EP | The internal/operator API surface for managing stock components is now confirmed from the internal swagger: `POST`/`DELETE`/`PATCH /image_components/{id}`, plus dedicated `.../deprecate` and `.../obsolete` actions — see 3.1.4. `PATCH` is also how `visibility` flips between `private` and `public`. One detail remains unconfirmed with the team: the exact request/response schema for the `POST /image_components` create call | Narrowed from "no mechanism described" to a single schema detail needing confirmation before commitment review | Open — Image Builder team, narrowed scope |

## 6.1 HashiCorp Packer Legal Approval

**Status**: Pending IBM Legal review

- HashiCorp Packer is licensed under Mozilla Public License 2.0 (MPL 2.0) — weak copyleft, no
  source disclosure requirements for IBM Cloud VPC code, no licensing fees
- Packer binaries obtained from official HashiCorp releases, checksum-verified, baked into the
  IBM-built task image in ICR; integrity verified using HashiCorp's published checksums/GPG
  signatures

**Impact**: Must be resolved before GA approval — critical dependency for the Tekton build
pipeline.

**Action items**: submit formal legal review request; provide complete dependency analysis and
usage documentation; obtain written approval for production use.

---

# 7. Patents

None.

---

# 8. Meta

None.

---

# 9. Glossary

| Term | Description |
| :--- | :--- |
| **Image job** | A first-class top-level regional resource (`POST/GET/PATCH/DELETE /v1/image/jobs`) that executes a recipe build and produces an image as its output. A required `type` field (currently a single-value enum, `build`, explicitly expandable) discriminates job types via a `oneOf`: the common fields (`id`, `status`, `created_at`, `type`, etc.) sit on the base `ImageJob` schema, while everything build-specific (`image_recipe`, `image`, `image_template`, `timeout`, `worker`) lives on a separate `ImageJobBuild` branch — the mechanism a future job type would extend without needing any of the build-specific fields. Coarse `status` — the 5-value enum `queued → running → succeeded \| failed \| deleting` (`running` persists through zone distribution; no `pending`/`cancelling`/`cancelled` states); `status_reasons` (capped at 20 entries) is populated only on `failed` (a 6-code set) — it stays empty while `running`, since live stage-progress isn't surfaced to customers, even though the pipeline internally still tracks it (see 3.1.1.9.3a) — a customer-configurable `timeout` (default 120 min, 30–480 range), a `worker` object (wrapping `instance`, `subnet`, and `profile` together), and an `image` reference. `image_recipe`, naming the recipe to build from, accepts a `oneOf`: by `id`, by CRN, or by `href`. There is no separate top-level `resource_group` on the job itself, only on `image_template`, for the resulting image. `name` is the only PATCH-able field. `DELETE` = cancel (non-terminal — cancels the job and deletes the image being built, where one exists — moves through `deleting`) or record removal (terminal, does not affect the produced image). Strictly 1:1 with its image, and retained for as long as its image exists; no automated retries — re-running a build is a new job. |
| **Profile** | The VPC VSI profile used for the build VSI — specified as `worker.profile` on `POST /v1/image/jobs`, defaulting to `cxf-4x8` if unspecified, since NGDC does not have Gen2 capacity. Not a recipe field or a settings field — see 3.1.1.2 for why it lives on the job. |
| **CEL (Common Expression Language)** | Used in Phase 1 by the component's `allowed_use.image` expression, checking `operating_system.family`/`version`/`architecture` against the recipe's `base_image`. Also the expression language for the image's own pre-existing `allowed_use.instance`/`allowed_use.bare_metal_server`. The component's `allowed_use` briefly existed as a plain `{operating_system: {family, version}}` object before being rebuilt as a CEL expression — see 3.1.1.3.6a and 5.2. |
| **Component** | An IBM-curated software package or GPU driver that can be included in a recipe. Each version is a distinct immutable catalog entry. Phase 1 public response schema: `allowed_use` (object: `api_version`, `image` — a CEL expression using `operating_system.family`/`version`/`architecture` as variables), `software` (a `oneOf`, currently one branch: `{name, version}` — the actual software identity), `created_at`, `crn`, `dependencies` (flat array, no `.all_of` sub-property, max 19), `deprecation_at`, `href`, `id`, `name` (the catalog resource's own identifier, distinct from `software.name`), `obsolescence_at`, `resource_group`, `resource_type`, `status` (6 values: `pending`, `available`, `deprecated`, `degraded`, `obsolete`, `deleting`, with `degraded` deriving when this component's own dependency becomes obsolete), `visibility` (`private` or `public` — a real, public field; only its write path is operator-only). `resource_group` and `resource_type` are required fields. Customers cannot create, edit, or delete component content — the one exception is the public `.../deprecate` and `.../obsolete` action endpoints, each with its own dedicated IAM action, which let a customer transition a component's lifecycle state directly (see 3.1.1.9.1). All Phase 1 components install via precompiled packages, with in-build validation scoped to package-level integrity checks accordingly (see 2.8.2.9). |
| **Fast-complete (superseded, not current design)** | An earlier design considered for `POST /v1/image/jobs`: if a recipe `(id, version)` already had an available image, the job would return `succeeded` immediately rather than building again. This was superseded by the current design — every build produces a new image, regardless of whether identical content was already built (see 2.2 and the "Generalized Fast-Complete" decision record in 5.2). Retained here only so the term isn't confusing if encountered elsewhere. |
| **HCL (HashiCorp Configuration Language)** | The configuration language in which Packer build templates are expressed. Image Builder generates HCL programmatically (via the `hclwrite` Go library) from component schema — no string templating, no authored scripts. |
| **IAM S2S authorization** | An IBM Cloud IAM service-to-service authorization grant created by the customer as a one-time prerequisite. Grants `is.image-recipe` permission on exactly four target types (instances, keys, floating IPs, security groups) in the customer's account. |
| **ImageComponent / ImageRecipe / ImageBuildJob** | The Kubernetes CRD Kind names for component, recipe, and image job metadata (`image.vpc.cloud.ibm.com/v1`). |
| **OnePipeline** | IBM's shared, centrally-maintained DevSecOps/compliance pipeline platform, still running on the underlying CD Tekton service. Under evaluation as a possible replacement for the raw Tekton dispatch model — see OI-PIPE. |
| **packer-plugin-ibmcloud** | The IBM-maintained HashiCorp Packer plugin (Apache 2.0) executing inside Task 1's IBM-built task image. Handles the full build VSI lifecycle: provisioning, SSH, package installation, image capture; the cleanup task guarantees ephemeral-resource deletion on all exit paths. |
| **Recipe** | A regionally-scoped resource that defines how to build a custom VPC image: a base OS image and a set of components. `worker.profile` and `subnet` are specified on the job, not the recipe — see 3.1.1.2. One mutable resource per **immutable id** (name is mutable identity metadata); content changes mint content-hash versions with superseded content preserved in an addressable history. Stock recipes are the same resource, IBM-owned, `visibility: public`, `ibm-` name prefix. |
| **Recipe history** | The recipe's version store: every version, current and superseded, as a frozen, permanently addressable snapshot, ordered by `created_at` descending (most-recent first) — not `superseded_at`, since the current version has none. List all via `GET /v1/image_recipes/{id}/history` (paginated, `limit` default 50/min 1/max 100; `404` if the recipe doesn't exist; uses `is.image-recipe.image-recipe.list`); look up one specific version — current or superseded — via `.../history/{version}` (a plain path parameter, not a literal segment; uses `is.image-recipe.image-recipe.read`, not `.list`). Each entry's own `created_at` is that version's own creation instant, matching the prior version's `superseded_at`; the current version is stamped `superseded_at: null`. Resolves the provenance chain directly via the image's `source_image_recipe` reference field for the life of the recipe. |
| **Inventory** | The per-image record of installed packages: an extracted name/version list, derived from the Packer plugin's scan output during the build, and served at `GET /v1/images/{id}/inventory`. Fields: `created_at`, `href`, `packages` (the name/version list itself), `resource_type: "inventory"` — no `format` field. The full original scan document is not retained or exposed anywhere, per current ITS/CISO guidance. |
| **Subnet** | The VPC subnet the temporary build VSI is provisioned in — required as `worker.subnet` on `POST /v1/image/jobs`, with no account-level default or auto-detection (VPC has no default-subnet primitive, and the subnet determines the build VSI's security boundary). There is no settings API for this; an earlier design considered one and superseded it — see 3.1.1.2. |
| **Task 1 — Packer Execute** | The Tekton pipeline task that performs the build: fetch + verify the generated HCL template, provision the ephemeral build VSI, install and validate components, capture the image. Runs the IBM-built Packer task image. |
| **Task 2 — Post-process** | The Tekton pipeline task that publishes build results: extract the inventory, upload artifacts to COS, and (via the gRPC call described in 3.1.4) write the `source_image_recipe` and `image_job` fields onto the image once it exists. |
| **Tekton pipeline** | The IBM Continuous Delivery Tekton pipeline (regional; dedicated IBM PAYG account) executing Task 1 (Packer Execute), Task 2 (Post-process), and a Cleanup task that always runs. Dispatched via the CD REST API (`CreatePipelineRun`); HCL is generated, validated, and hashed in the control plane before dispatch. See OI-PIPE for the OnePipeline migration consideration. |
| **(id, version)** | The permanent recipe-content coordinate: immutable `id` plus content-hash `version`. Recorded in every build's provenance; resolvable via recipe history. (`name` is mutable identity metadata — unique per account+region but not the provenance anchor.) |

---

# 10. References

1. [Aha Feature: VPCIMG-277 — VPC Image Builder Phase 1](https://bigblue.aha.io/features/VPCIMG-277)
2. [IBM Cloud ibmcloud-image-builder (Open Source)](https://github.com/IBM-Cloud/ibmcloud-image-builder)
3. [HashiCorp Packer Documentation](https://www.packer.io/docs)
4. [HashiCorp Packer Plugin for IBM Cloud](https://github.com/IBM/packer-plugin-ibmcloud)
5. [AWS EC2 Image Builder](https://docs.aws.amazon.com/imagebuilder/latest/userguide/what-is-image-builder.html)
6. [Azure VM Image Builder](https://learn.microsoft.com/en-us/azure/virtual-machines/image-builder-overview)
7. [Red Hat Image Builder](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/composing_a_customized_rhel_system_image/index)
8. [IBM Cloud IAM Service-to-Service Authorization](https://cloud.ibm.com/docs/account?topic=account-serviceauth)
9. [IBM Cloud VPC API Reference](https://cloud.ibm.com/apidocs/vpc)
10. [Tekton Pipelines Documentation](https://tekton.dev/docs/pipelines/)

<!-- Emacs formatting setup -->
<!-- Local Variables: -->
<!-- fill-column: 100 -->
<!-- End: -->

<!-- Vim formatting setup. Must be in first or last 5 lines by default -->
<!-- vim: set textwidth=100 : -->