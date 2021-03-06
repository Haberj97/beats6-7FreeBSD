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

package mysql

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "mysql", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
