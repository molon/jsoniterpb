module github.com/molon/jsoniterpb

go 1.19

require (
	github.com/dlclark/regexp2 v1.7.0
	github.com/google/go-cmp v0.5.8
	github.com/google/gofuzz v1.2.0
	github.com/json-iterator/go v1.1.12
	github.com/modern-go/reflect2 v1.0.2
	github.com/srikrsna/goprotofuzz v0.0.0-20220606153644-8d0e21b5787a
	github.com/stretchr/testify v1.8.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// go get github.com/molon/jsoniter@jsoniterpb
replace github.com/json-iterator/go => github.com/molon/jsoniter v0.0.0-20230529062209-e42e40bd8588
