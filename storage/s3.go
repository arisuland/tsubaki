// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2021 Noelware
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
	"context"
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
)

// S3StorageProvider is a provider for using S3 for hosting your images.
// This can be any compatible S3 servers like Wasabi or Minio. In the docker-compose
// file, it will have configuration for a Minio instance to use within Arisu.
type S3StorageProvider struct {
	Config *S3StorageConfig
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

func NewS3StorageProvider(config *S3StorageConfig) BaseStorageProvider {
	return S3StorageProvider{
		Config: config,
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

//func get

func (s3 S3StorageProvider) Init() error {
	log.Info(context.Background(), "Creating S3 client...")

	//var cfg aws.Config
	//if c, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s3.Config.Region)); err != nil {
	//	log.Warn(context.Background(), "Credentials is not present in ~/.aws directory, checking from config.yml / environment vars...")
	//	if s3.Config.AccessKey == nil || s3.Config.SecretKey == nil {
	//		return errors.New("missing 'storage.s3.access_key' or 'storage.s3.secret_key' values in config.yml")
	//	}
	//
	//	accessKey, secretKey := *s3.Config.AccessKey, *s3.Config.SecretKey
	//	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	//	if c, err = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(s3.Config.Region)); err != nil {
	//		return err
	//	}
	//}

	log.Info(context.Background(), "Creating S3 session...")
	//session, err := session.NewSession()

	return nil
}

func (s3 S3StorageProvider) Name() string {
	return "s3"
}

func (s3 S3StorageProvider) GetMetadata(id string, project string) *ProjectMetadata {
	return nil
}

func (s3 S3StorageProvider) HandleUpload() {
	// TODO: this
}
