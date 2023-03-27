# bc-wallet-common-lib-config

## Description

Library for preparing application configs 
Library can prepare config from different sources to GO-lang structs

Library can prepare config from:
* ENV variables
* JSON files
* Secret management engine which implemented compatible interface

## Usage examples

Examples of create prepare application config

### From ENV variables

```go
package main

import (
	"context"

	commonEnvConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/envconfig"
)

// DO NOT EDIT THIS VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
// The value of these variables will be assigned at the build stage
var (
	Version = "DEVELOPMENT.VESION"
	ReleaseTag = "DEVELOPMENT.RELEASE_TAG"
	CommitID = "DEVELOPMENT.COMMIT_HASH"
	ShortCommitID = "DEVELOPMENT.SHORT_COMMIT_HASH"
	BuildNumber uint64 = 0
	BuildDateTS uint64 = 0
)

type NatsConfig struct {
	NatsAddresses string `envconfig:"NATS_ADDRESSES" default:"nats://ns-1:4223,nats://ns-2:4224,nats://na-3:4225"`
	NatsUser      string `envconfig:"NATS_USER" secret:"true"`
	NatsPassword  string `envconfig:"NATS_PASSWORD" secret:"true""`
}

type DbConfig struct {
	DatabaseDriver              string `envconfig:"DATABASE_DRIVER" required:"true"`
	DatabasePort                uint16 `envconfig:"DATABASE_PORT" default:"54321"`
	DatabaseUser                string `envconfig:"DATABASE_USER" secret:"true"`
	DatabasePassword            string `envconfig:"DATABASE_PASSWORD" secret:"true"`
	TestFieldForSecretOverwrite string `envconfig:"TEST_FIELD_FOR_OVERWRITE_BY_SECRET" secret:"true"`
}

type AppConfig struct {
	*NatsConfig
	*DbConfig
}

func main() {
	ctx := context.Background()

	flagManagerSrv := commonEnvConfig.NewLdFlagsManager(Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)

	baseCfgPreparerSrv := commonEnvConfig.NewConfigManager()
	baseCfg := commonEnvConfig.NewBaseConfig(applicationName)
	err := baseCfgPreparerSrv.PrepareTo(baseCfg).With(flagManagerSrv).Do(ctx)
	if err != nil {
		panic(err)
	}
	
	appCfg := &AppConfig{}
	appCfgPreparerSrv := commonEnvConfig.NewConfigManager()
	err = appCfgPreparerSrv.PrepareTo(appCfg).With(flagManagerSrv, baseCfg).Do(ctx)
	if err != nil {
		panic(err)
	}
}

```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron) 

## Licence

**bc-wallet-common-lib-config** is licensed under the [MIT](./LICENSE) License.