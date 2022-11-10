// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO completely account for custom renaming, with phase separation.
//
// Currently schems reuse tfgen which calls PulumiToTerraformName, and
// among other things plurlizes names of list properties. This code
// accounts for this for now to pass unit tests. But all the cases
// need to be covered.

package tfbridge

import (
	"github.com/gedex/inflector"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func toPropertyKey(name string, typ tftypes.Type) resource.PropertyKey {
	if willPluralize(name, typ) {
		return resource.PropertyKey(inflector.Pluralize(name))
	}
	return resource.PropertyKey(name)
}

func willPluralize(name string, typ tftypes.Type) bool {
	if typ.Is(tftypes.List{}) {
		plu := inflector.Pluralize(name)
		distinct := plu != name
		valid := inflector.Singularize(plu) == name
		return valid && distinct
	}
	return false
}
