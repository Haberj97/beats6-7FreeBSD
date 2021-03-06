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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "../libbeat/fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
