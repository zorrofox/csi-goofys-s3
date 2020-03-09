package s3

import (
	"fmt"
	"context"
    "github.com/aws/aws-sdk-go/aws/endpoints"

	goofysApi "github.com/kahing/goofys/api"
	common "github.com/kahing/goofys/api/common"
)

const (
	defaultRegion = "us-east-1"
)

// Implements Mounter
type goofysMounter struct {
	bucket *bucket
	region string
	readonly bool
}

func newGoofysMounter(b *bucket, cfg *Config) (Mounter, error) {
	region := cfg.Region
	// if endpoint is set we need a default region
	if region == "" {
		region = defaultRegion
	}
	return &goofysMounter{
		bucket: b,
		region: region,
		readonly: cfg.ReadOnly,
	}, nil
}

func (goofys *goofysMounter) Stage(stageTarget string) error {
	return nil
}

func (goofys *goofysMounter) Unstage(stageTarget string) error {
	return nil
}

func (goofys *goofysMounter) Mount(source string, target string) error {
	var mountOptions map[string]string
	resolver := endpoints.DefaultResolver()
	endpoint, err := resolver.EndpointFor(endpoints.S3ServiceID, goofys.region)
	if err != nil {
	    return err
	}
	mountOptions = map[string]string{
		"allow_other": "",
	}
	if goofys.readonly {
        mountOptions["ro"] = "" 
	}
	goofysCfg := common.FlagStorage{
		MountPoint: target,
		Endpoint:   endpoint.URL,
		DirMode:    0755,
		FileMode:   0644,
		MountOptions: mountOptions,
	}

	fullPath := fmt.Sprintf("%s:%s", goofys.bucket.Name, goofys.bucket.FSPath)

	_, _, errApi := goofysApi.Mount(context.Background(), fullPath, &goofysCfg)

	if errApi != nil {
		return fmt.Errorf("Error mounting via goofys: %s", errApi)
	}
	return nil
}
