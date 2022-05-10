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

package postgresql

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "postgresql", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
