// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package add_cloud_metadata

import (
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/mapstr"
)

const (
	hetznerMetadataInstanceIDURI       = "/hetzner/v1/metadata/instance-id"
	hetznerMetadataHostnameURI         = "/hetzner/v1/metadata/hostname"
	hetznerMetadataAvailabilityZoneURI = "/hetzner/v1/metadata/availability-zone"
	hetznerMetadataRegionURI           = "/hetzner/v1/metadata/region"
)

// Hetzner Cloud Metadata Service
// Document https://docs.hetzner.cloud/#server-metadata
var hetznerMetadataFetcher = provider{
	Name:           "hetzner-cloud",
	DefaultEnabled: true,
	Create: func(_ string, c *conf.C) (metadataFetcher, error) {
		hetznerSchema := func(m map[string]interface{}) mapstr.M {
			m["service"] = mapstr.M{
				"name": "Cloud",
			}
			return mapstr.M{"cloud": m}
		}

		urls, err := getMetadataURLs(c, metadataHost, []string{
			hetznerMetadataInstanceIDURI,
			hetznerMetadataHostnameURI,
			hetznerMetadataAvailabilityZoneURI,
			hetznerMetadataRegionURI,
		})
		if err != nil {
			return nil, err
		}

		responseHandlers := map[string]responseHandler{
			urls[0]: func(all []byte, result *result) error {
				result.metadata.Put("instance.id", string(all))
				return nil
			},
			urls[1]: func(all []byte, result *result) error {
				result.metadata.Put("instance.name", string(all))
				return nil
			},
			urls[2]: func(all []byte, result *result) error {
				result.metadata["availability_zone"] = string(all)
				return nil
			},
			urls[3]: func(all []byte, result *result) error {
				result.metadata["region"] = string(all)
				return nil
			},
		}
		fetcher := &httpMetadataFetcher{"hetzner", nil, responseHandlers, hetznerSchema}
		return fetcher, nil
	},
}
