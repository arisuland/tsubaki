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

/*
  # Configures using S3 to host your projects, once the bucket is gone,
  # Arisu will attempt to create the bucket but your data will be lost.
  #
  # Type: S3StorageConfig
  s3:
    # Returns the provider to use when authenticating to S3. Arisu supports
    # Amazon S3, Wasabi, or using Minio. By default, it will attempt to use
    # Amazon S3.
    #
    # Type: S3Provider?
    # Variable: TSUBAKI_STORAGE_S3_PROVIDER
    # Default: S3Provider.AMAZON
    provider: S3Provider

    # Returns the bucket to use when storing files. If this bucket
    # doesn't exist, Arisu will attempt to create the bucket.
    # By default, Arisu will use `arisu` as the default bucket name
    # if this is not set.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_BUCKET
    # Default: "arisu"
    bucket: String

    # Returns the access key for authenticating to S3. If this isn't provided,
    # it will attempt to look for your credentials stored in `~/.aws`. This is a
    # recommended variable to set if using the S3 provider.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_ACCESS_KEY
    # Default: "access_key" key in ~/.aws/tsubaki_config
    access_key: String

    # Returns the secret key for authenticating to S3. If this isn't provided,
    # it will attempt to look for your credentials stored in `~/.aws`. This is a
    # recommended variable to set if using the S3 provider.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_SECRET_KEY
    # Default: "access_key" key in ~/.aws/tsubaki_config
    secret_key: String

    # Returns the region to host your bucket, this is dependant on if you
    # created the bucket without running Tsubaki. This is required to set to
    # so no errors will occur while authenticating to S3.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_REGION
    # Default: "us-east1"
    region: String
*/
