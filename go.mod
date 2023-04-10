module github.com/evenh/az-blob-hashdeep

go 1.19

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0
	github.com/openlyinc/pointy v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.6.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0 => github.com/evenh/azure-sdk-for-go/sdk/storage/azblob v0.2.1-0.20220128100502-5d716a1d24c2
