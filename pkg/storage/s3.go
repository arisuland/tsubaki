// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2022 Noelware
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"time"
)

// Provider is the S3 provider to use the endpoints from
type Provider string

var (
	// Wasabi is a Provider that uses Wasabi as your storage bucket.
	Wasabi Provider = "wasabi"

	// Custom is a custom-set URI to use to connect to Amazon S3
	Custom Provider = "custom"

	// Amazon is the Provider to use Amazon S3, the URI will return nil
	Amazon Provider = "amazon"

	// Empty is a invalid provider.
	Empty Provider = "empty"
)

// String stringifies a Provider value type.
func (p Provider) String() string {
	switch {
	case p == Wasabi:
		return "wasabi"

	case p == Custom:
		return "custom"

	case p == Amazon:
		return "amazon"

	default:
		return "empty"
	}
}

func isValidProvider(provider string) bool {
	allProviders := []string{"wasabi", "custom", "amazon"}
	for _, p := range allProviders {
		if provider == p {
			return true
		}
	}

	return false
}

func FromProvider(provider string) Provider {
	if !isValidProvider(provider) {
		return Empty
	}

	switch provider {
	case "wasabi":
		return Wasabi

	case "custom":
		return Custom

	case "amazon":
	default:
		return Amazon
	}

	// we should never go here
	return Empty
}

// S3StorageConfig is the configuration for a S3StorageProvider instance.
type S3StorageConfig struct {
	// SecretKey is the secret key used to authenticate.
	SecretKey *string `yaml:"secret_key"`

	// AccessKey is the access key used to authenticate.
	AccessKey *string `yaml:"access_key"`

	// Provider is the Provider to use instead of Amazon S3
	Provider Provider `yaml:"provider"`

	// Endpoint is the custom endpoint to use to authenticate.
	Endpoint *string `yaml:"endpoint"`

	// Region is a S3 region to use.
	Region string `yaml:"region"`

	// Bucket is the bucket to use.
	Bucket string `yaml:"bucket"`
}

type S3StorageProvider struct {
	config *S3StorageConfig
	client *s3.S3
}

func NewS3StorageProvider(config *S3StorageConfig) BaseStorageProvider {
	return S3StorageProvider{
		config: config,
		client: nil,
	}
}

func (s S3StorageProvider) Init() error {
	logrus.Info("Now creating S3 client...")

	cfg := aws.NewConfig().WithRegion(s.config.Region)
	if s.config.SecretKey != nil && s.config.AccessKey != nil {
		cfg.WithCredentials(credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     *s.config.AccessKey,
			SecretAccessKey: *s.config.SecretKey,
		}))
	}

	if s.config.Endpoint != nil {
		var endpoint = ""
		switch {
		case s.config.Provider == Wasabi:
			endpoint = "https://s3.wasabisys.com"

		case s.config.Provider == Custom:
			endpoint = *s.config.Endpoint
		}

		if endpoint != "" {
			cfg.WithEndpoint(endpoint)
		}
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return err
	}

	client := s3.New(sess)
	logrus.Info("Created S3 client, checking bucket list...")

	t := time.Now()
	buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	logrus.Infof("Received %d buckets in %s", len(buckets.Buckets), time.Since(t).String())
	exists := false

	for _, b := range buckets.Buckets {
		if b.Name != nil && *b.Name == s.config.Bucket {
			exists = true
		}
	}

	if !exists {
		t := time.Now()
		logrus.Warnf("Bucket %s doesn't exist, now creating...", s.config.Bucket)

		_, err := client.CreateBucket(&s3.CreateBucketInput{
			Bucket: &s.config.Bucket,
		})

		if err != nil {
			return err
		}

		logrus.Infof("Created bucket %s in %s.", s.config.Bucket, time.Since(t).String())
	}

	s.client = client
	return nil
}

func (s S3StorageProvider) Name() string {
	return "s3"
}

func (s S3StorageProvider) GetMetadata(id string, project string) *ProjectMetadata {
	return nil
}

func (s S3StorageProvider) HandleUpload(files []UploadRequest) error {
	return nil
}
