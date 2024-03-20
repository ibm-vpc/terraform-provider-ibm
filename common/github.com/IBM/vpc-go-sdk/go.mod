module github.com/IBM/vpc-go-sdk

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.14.1
	github.com/go-openapi/strfmt v0.21.5
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.8.2
)

retract (
	v1.0.2
	v1.0.1
	v1.0.0
)
