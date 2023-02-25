module github.com/evenh/az-blob-hashdeep

go 1.19

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0
	github.com/openlyinc/pointy v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.6.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0 => github.com/evenh/azure-sdk-for-go/sdk/storage/azblob v0.2.1-0.20220128100502-5d716a1d24c2
