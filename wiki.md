
Skip to content
Navigation Menu
CloudEngineering
openapi-sdkgen
Type / to search

    Code
    Pull requests 10
    Actions
    Projects
    Wiki
    Security
    Insights

Terraform Provider Generation
Dustin Popp edited this page on Jun 28, 2022 · 16 revisions

The OpenAPI SDKGen project is capable of generating a Terraform provider for a service, in addition to SDKs. The code produced by the generator must be combined with the code in the IBM Terraform Provider to produce a functioning Terraform provider.

Note that in order to generate a working Terraform provider, a working Go SDK must be generated for the same service. The generated Terraform provider will depend on the generated Go SDK. Depending on the complexity of your service, some hand-modifications may be needed in order to properly build, and run, any generated Terraform provider code.
Preparing the API Definition

There are a few additions to make to the service's API Definition in order to generate a Terraform provider.
Go SDK Package Name

If a Go SDK already exists for this service, skip this step. If not, generating a Go SDK will be necessary. The package name for the generated Go SDK needs to be included in the API Definition, under the key hierarchy info -> x-codegen-config -> go -> apiPackage. This will be used to import the Go SDK in the generated Terraform provider code. Example:

info:
  x-codegen-config:
    go:
      apiPackage: github.com/IBM/platform-services-go-sdk

Terraform Resource

If you do not wish to generate Terraform resource(s) for this service, skip this step.

The Terraform provider generator supports generating Terraform resources from schema objects within the API definition. The annotated schema objects should be schemas that are used to describe an operation's response body. Use the x-resource-name annotation to annotate a schema object as a Terraform resource. The value of the annotation will serve as the name of a generated resource.

Example:

Target:
  description: A Target a ATracker data collection.
  type: object
  x-resource-name: atracker_target
  properties:
    ...

For an API that follows common REST design, we can identify the create, read, update, and delete operations for the resource as follows:

    to find the create operation, the generator looks for a POST operation on a path whose final path segment is not a path parameter and whose response body is the resource schema

    from the create operation, we can find the read, update, and delete operations by looking for the appropriate methods on a path that matches the path of Create but with one final path parameter — which we assume is the resource ID.

If your API doesn’t follow this pattern, use the x-resource-operations annotation to specify the create, read, update, and delete operations.

The arguments of the resource will be the parameters and request body properties of the create operation for the resource.

The attributes of the resource will be the properties of the schema, with nested inner models as appropriate.

Note: The generator makes the assumption that a parameter or request body property on create called "foo" will have the same type as the "foo" property in the resource schema if that exists. This is almost always true but if your API has some good reason for these to be different types, then you'll have some extra work to do to fix up the generated code to make it work.

There are a few other annotations that are important. For example

x-terraform-user-account - on a parameter or property that is used to supply the users account ID to the service

There’s special processing we need to do for the user account so its important that these be flagged appropriately.

There are more annotations — the ones specific to Terraform generation are described here:  https://github.ibm.com/CloudEngineering/openapi-sdkgen/wiki/API-Document-annotations-for-Terraform-generation
Terraform Data Source

If you do not wish to generate Terraform data source(s) for this service, skip this step.

The Terraform provider generator supports generating Terraform data sources from schema objects within the API definition. The annotated schema objects should be the response schema of a GET operation in the API -- this GET operation will be used in the data source "read" method to retrieve the data source content. Use the x-data-source-name annotation to annotate a schema object as a Terraform data source. The value of the annotation will serve as the name of a generated data source.

Example:

Target:
  description: Property values for a target in response. Credentials associated with the target are encrypted and masked as REDACTED in the response.
  type: object
  x-data-source-name: atracker_target
  properties:
    ...

Note: Annotating multiple schema objects as data sources / resources is supported, and a schema object can be annotated as both a resource, and a data source.

The arguments of the data source will be the set of required parameters of the GET operation for the data source.

The attributes of the data source will be the properties of the schema, with nested inner models as appropriate.
Types of Terraform data sources

Terraform data sources can have many forms, but there are two very common types:

resource instance: This type of data source corresponds directly to a Terraform resource and can be used to obtain information about an instance of that Terraform resource. For this type of data source, the x-data-source-name annotation should be on the same schema as x-resource-name, and typically with the same value.

resource collection: This type of data source corresponds to a list or collection of some type of Terraform resource. For this type of data source, the x-data-source-name annotation is typically on the response schema for the "list" operation for the resource.

When the data source is not a resource instance type data source, the generator assumes that the data source is a collection type data source if the data source schema contains a top-level property that is an array whose items are a Terraform resource. If the data source schema contains multiple top-level array properties, an additional annotation, x-data-source-collection, should be specified on the schema to specify the name of the property that is a collection of resources.
Create a project

The generated Terraform code relies on a framework for a Terraform provider that contains the initialization logic and utility methods that are shared across all services. To simplify development and testing of Terraform provider code, we created a Terraform provider template that contains the essential elements of this framework for generated Terraform providers. This allows you to develop your Terraform provider support in a small focused project.

To create your development/test project, follow the instructions in the template project README.md.

After you generated and tested the Terraform provider support in this project, you will merge the appropriate files into the IBM Terraform Provider. Follow the instructions in the provider project README for testing your code in the context of the IBM Terraform provider.
Generating the code
Generate the Go SDK (Optional)

If a Go SDK already exists for this service, skip this step. If not, a Go SDK must be generated. Follow these in-depth instructions on setting up a new Go SDK.
Generate the Terraform Provider

Download the generator and run the generator using the API Definition as input, setting the language to “ibm-terraform”, and setting the output directory to the repository where the IBM Terraform Provider was cloned in an earlier step. For the first generation, the configuration option initialize must be set to true. This can be accomplished with the command line flag --additional-properties initialize=true when generating. Setting this flag will cause a file called instructions.md to be generated, which contains instructions for changes to make when migrating to the IBM Terraform Provider, as well as custom versions of config.go and provider.go to be generated for use in your local project. The code will be generated into the correct locations in the file structure automatically.

Example:

openapi-sdkgen.sh generate -i /my/apidefdir/apidef.json -g ibm-terraform -o /path/to/terraform-provider-ibm --additional-properties initialize=true

Note: Generating a Terraform provider will also generate documentation files and examples for any generated data source(s) / resource(s). Documentation files can be found within the website/docs directory under the output directory containing the generated provider. Generated examples can be found under the examples directory under the output directory containing the generated provider.
Configure the IBM Terraform Provider

Navigate to the repository created from cloning the IBM Terraform Provider, where the code was generated. Follow the instructions located in the generated instructions.md file. As mentioned before, this file is generated as part of generating a Terraform provider when the configuration option initialize is set to true. This file contains instructions to add the generated Terraform provider to the IBM Terraform Provider repo cloned, including changes to provider.go and config.go. These instructions must be completed to properly build and run the generated code.
Building the Terraform Provider

Within the repository created from cloning the IBM Terraform Provider, where the code was generated, build the provider:

make build

This command will install the necessary dependencies and compile the code.
Running the Generated Tests

To run the generated acceptance tests, please see Developing the Provider.

Note: Acceptance tests create real resources, and often cost money to run.
Running the Generated Examples

Navigate to the examples directory within the cloned IBM Terraform provider repo.

To run any example, follow the Usage instructions from the example's generated README.md within the example's own directory. The generated examples have a README.md containing more details on what the example does.

Note: Examples create real resources, and often cost money to run.
Reference
Available API Definition Annotations
Terraform Codgen Config

info -> x-codegen-config -> terraform
apiEndpoint (String)

    The name of an environment variable that may be used to override the API endpoint for the service.

subcategory (String)

    The name of the subcategory for the service that will be inserted into the subcategory header field at the top of each generated resource and data-source doc file.
    If this value is not specified, the service name or title will be used. This list of allowed subcategory values can here. For new subcategories, the website/allowed-subcategories.txt file should be updated to include them.

OpenAPI annotations

These are annotations that can be used within the API definition to control the behavior of the Terraform provider generator. See the API annotations page for more information on these annotations.
x-terraform-user-account (boolean)

    OpenAPI annotation on a parameter or request body property to specify that its value should be set from the userAccount field of UserDetails.

Example:

 /groups:
  post:
  operationId: create_access_group
  summary: Create an Access Group
  parameters:
    - description: IBM Cloud account id under which the group is created
       name: account_id
        in: query
        required: true
        schema:
          type: string
        x-terraform-user-account: true

x-resource-name (String)

    OpenAPI annotation on a schema to specify the name for the Terraform resource to be generated for this schema. Adding this annotation on a schema will result in a resource file to be generated, along with any supporting files for that resource.

x-data-source-name (String)

    OpenAPI annotation on a schema to specify the name for the Terraform data source to be generated for this schema. Adding this annotation on a schema will result in a data source file to be generated, along with any supporting files for that data source.

x-resource-id (String)

    OpenAPI annotation on a schema to specify an expression that will produce the ID for the Terraform resource / data source to be generated for this schema. If not specified, then the generator will attempt to choose a property to serve as the ID for the Terraform resource / data source.

x-terraform-computed (boolean)

    OpenAPI annotation on a property of a resource or data source model/schema to specify that it should defined as a computed attribute. Note that readOnly: true implies a computed field so this annotation is not needed on readOnly properties, and properties within a schema marked as a data source are automatically generated as a computed attribute.

x-terraform-include-tags (boolean)

    OpenAPI annotation on a schema that is also annotated as a Terraform resource to specify that tags should be added to the resource schema.

Example:

Group:
    type: object
    description: An IAM access group.
    title: IAM access group
    x-resource-name: iam_access_group
    x-terraform-include-tags: true
    properties: ...

x-resource-operations (Map<String, String>)

    OpenAPI annotation on a schema to specify the operationIds of the create, read, update, and delete operations associated with a resource. This should only be added if the associated operations of a resource cannot be auto determined by the generator.

Example:

"NamespaceDetail": {
    "description": "Details of a namespace",
    "type": "object",
    "x-resource-name": "cr_namespace",
    "x-resource-operations": {
        "create": "create_namespace",
        "read": "list_namespace_details",
        "delete": "delete_namespace"
    },
}

x-data-source-operations (Map<String, String>)

    OpenAPI annotation on a schema specify the operationId of the read operation associated with a data source. This should only be added if the associated read operation of a data source cannot be auto determined by the generator.

Example:

"NamespaceDetail": {
    "description": "Details of a namespace",
    "type": "object",
    "x-data-source-name": "cr_namespace",
    "x-data-source-operations": {
        "read": "list_namespace_details"
    },
}

x-data-source-collection (String)

    OpenAPI annotation on a schema to specify the property in the data source's model that is a collection of resources. If not specified, then the collection is auto-determined based on the properties of the data source's model. If a data source's read operation supports pagination, or you wish to enable a "filtered" data source, then adding this annotation is recommended.

Example:

GroupsList:
    type: object
    description: The list of access groups returned as part of a response.
    x-data-source-name: iam_access_groups
    x-data-source-collection: groups
    properties:
        groups:
          type: array
          description: An array of access groups.
          items:
            $ref: '#/components/schemas/Group'

x-terraform-data-source-filter (String)

    OpenAPI annotation on a schema to specify the property to use as a filter in the generated data source. Adding this annotation on a data source model will result in a generated data source that is "filterable" by the name of the property specified as the value of the annotation. The filter will be generated as a optional argument to the data source. If x-data-source-collection is also specified, then the generator will attempt to apply the filter to the collection specified.

Example:

GroupsList:
    type: object
    description: The list of access groups returned as part of a response.
    x-data-source-name: iam_access_groups
    x-data-source-collection: groups
    x-terraform-data-source-filter: name
    properties:
        groups:
          type: array
          description: An array of access groups.
          items:
           $ref: '#/components/schemas/Group'
Group:
      type: object
      description: An IAM access group.
      title: IAM access group
      properties:
        id:
          type: string
          description: The group's Access Group ID
        name:
          type: string
          description: The group's name

x-terraform-timeouts (Map<String, String>)

    OpenAPI annotation on a schema to specify the timeouts for the create, read, update, and delete operations associated with a resource or data source. The timeout value is specified in minutes.

Example:

Group:
      type: object
      description: An IAM access group.
      title: IAM access group
      x-terraform-resource-name: iam_access_group
      x-terraform-timeouts:
        create: 60
        update: 20
        delete: 10

x-terraform-exclude (boolean)

    OpenAPI annotation on a parameter or property to specify that it should be excluded from the arguments or attributes for a resource or data source in the generated terraform provider. For properties and parameters that contain a reference object to a schema, this annotation can be added to the referenced schema to exclude that parameter/property from the arguments or attributes.

x-terraform-force-new (boolean)

    OpenAPI annotation on a parameter or property to specify that it should be declared as ForceNew: true in the resource schema. More info: https://www.terraform.io/docs/extend/schemas/schema-behaviors.html#forcenew

x-terraform-resource-crn (String)

    OpenAPI annotation on a schema to specify the name of the property that contains the CRN for the Terraform resource to be generated for this schema.

x-terraform-sensitive (Boolean)

    OpenAPI annotation on a property of a resource or data source model to specify that it should defined as a sensitive attribute.

x-terraform-conflicts-with (List<String>)

    OpenAPI annotation on a property of a resource or data source model to specify that it conflicts with another attribute. This annotation is defined as a list of property names to be added to the conflictsWith{} function repesenting the conflicting properties.

x-terraform-diff-suppress (String)

    OpenAPI annotation on a property of a resource model that specifies the name of a user-defined function that calculates the diff of this field and returns a boolean value, true when the two inputs are equal or deemed "equivalent" by the function, and false otherwise

x-ibm-resource-group-id (boolean)

    OpenAPI annotation on a parameter or request body property to specify that it holds a resource group id.

Example:

"/api/v1/namespaces/{namespace}": {
  "put": {
    "tags": [
      "Namespaces"
    ],
    "summary": "Create namespace",
    "description": "Add a namespace to the targeted IBM Cloud account.",
    "operationId": "create_namespace",
    "parameters": [
      {
        "description": "The ID of the resource group that the namespace will be created within.",
        "name": "X-Auth-Resource-Group",
        "in": "header",
        "schema": {
          "type": "string"
        },
        "x-ibm-resource-group-id": true
      }
    ]
  }
}

x-terraform-polling-operation (String)

    OpenAPI annotation on an operation that specifies an operationID of a get operation to poll for to determine the current state of the operation.

x-terraform-polling-property (String)

    OpenAPI annotation on an operation that specifies the property of a get operation to compare values for to determine the current state of the operation.

x-terraform-polling-success (String)

    OpenAPI annotation on an operation that specifies a state representing the successful completion of the operation that the provider should poll for.

x-terraform-polling-pending (String)

    OpenAPI annotation on an operation that specifies a state representing the pending status of the operation that the provider should poll for.

x-terraform-polling-failure (String)

    OpenAPI annotation on an operation that specifies a state representing the failed completion of the operation that the provider should poll for.

x-cli-flag-name (String)

    OpenAPI annotation on a parameter or property to specify the Terraform argument/attribute name.

x-cli-description (String)

    OpenAPI annotation on a parameter or property to specify the description for a Terraform argument/attribute.

Using Links

Using Open API links within your API definition is encouraged, and the Terraform provider generator makes partial use of them while generating tests and examples of Terraform resources. Enabling links within your API definition may result in better representative mock values in generated tests and examples of Terraform resources.
Pages 85
Home
Getting Started

SDK Generation Overview
Links

    SDK Cores
    SDK Templates
    API Handbook
    SDK Handbook
    OpenAPI Spec
    SDK Video Tutorials

Agile Practices

    Planning
    Retrospectives

Generator Usage

    Usage Instructions
    Config Options
    Compatibility Chart

Clone this wiki locally
Footer
GitHub Enterprise Server 3.19.10
IBM - Office of the CIO
IBM - Office of the CIO avatar IBM - Office of the CIO

    GHE Terms of Service
    IT Support
    Planned Maintenance
    IBM Internal Privacy Statement

    GHE Terms of Service
    IT Support
    Planned Maintenance
    IBM Internal Privacy Statement

© 2026 GitHub, Inc.
Footer navigation

    Help
    Support

