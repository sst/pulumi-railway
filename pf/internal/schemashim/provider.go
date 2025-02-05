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

package schemashim

import (
	"context"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"

	pfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/pulumi/pulumi-terraform-bridge/pf/internal/pfutils"
)

type SchemaOnlyProvider struct {
	ctx context.Context
	tf  pfprovider.Provider
}

func (p *SchemaOnlyProvider) PfProvider() pfprovider.Provider {
	return p.tf
}

var _ shim.Provider = (*SchemaOnlyProvider)(nil)

func (p *SchemaOnlyProvider) Schema() shim.SchemaMap {
	ctx := p.ctx
	schemaResp := &pfprovider.SchemaResponse{}
	p.tf.Schema(ctx, pfprovider.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		panic("Schema() returned error diags")
	}
	return newSchemaMap(pfutils.FromProviderSchema(schemaResp.Schema))
}

func (p *SchemaOnlyProvider) ResourcesMap() shim.ResourceMap {
	resources, err := pfutils.GatherResources(context.TODO(), p.tf)
	if err != nil {
		panic(err)
	}
	return &schemaOnlyResourceMap{resources}
}

func (p *SchemaOnlyProvider) DataSourcesMap() shim.ResourceMap {
	dataSources, err := pfutils.GatherDatasources(context.TODO(), p.tf)
	if err != nil {
		panic(err)
	}
	return &schemaOnlyDataSourceMap{dataSources}
}

func (p *SchemaOnlyProvider) Validate(context.Context, shim.ResourceConfig) ([]string, []error) {
	panic("schemaOnlyProvider does not implement runtime operation Validate")
}

func (p *SchemaOnlyProvider) ValidateResource(
	context.Context, string, shim.ResourceConfig,
) ([]string, []error) {
	panic("schemaOnlyProvider does not implement runtime operation ValidateResource")
}

func (p *SchemaOnlyProvider) ValidateDataSource(
	context.Context, string, shim.ResourceConfig) ([]string, []error) {
	panic("schemaOnlyProvider does not implement runtime operation ValidateDataSource")
}

func (p *SchemaOnlyProvider) Configure(ctx context.Context, c shim.ResourceConfig) error {
	panic("schemaOnlyProvider does not implement runtime operation Configure")
}

func (p *SchemaOnlyProvider) Diff(
	context.Context, string, shim.InstanceState, shim.ResourceConfig, shim.DiffOptions,
) (shim.InstanceDiff, error) {
	panic("schemaOnlyProvider does not implement runtime operation Diff")
}

func (p *SchemaOnlyProvider) Apply(
	context.Context, string, shim.InstanceState, shim.InstanceDiff,
) (shim.InstanceState, error) {
	panic("schemaOnlyProvider does not implement runtime operation Apply")
}

func (p *SchemaOnlyProvider) Refresh(
	context.Context, string, shim.InstanceState, shim.ResourceConfig,
) (shim.InstanceState, error) {
	panic("schemaOnlyProvider does not implement runtime operation Refresh")
}

func (p *SchemaOnlyProvider) ReadDataDiff(
	context.Context, string, shim.ResourceConfig,
) (shim.InstanceDiff, error) {
	panic("schemaOnlyProvider does not implement runtime operation ReadDataDiff")
}

func (p *SchemaOnlyProvider) ReadDataApply(
	context.Context, string, shim.InstanceDiff,
) (shim.InstanceState, error) {
	panic("schemaOnlyProvider does not implement runtime operation ReadDataApply")
}

func (p *SchemaOnlyProvider) Meta(context.Context) interface{} {
	panic("schemaOnlyProvider does not implement runtime operation Meta")
}

func (p *SchemaOnlyProvider) Stop(context.Context) error {
	panic("schemaOnlyProvider does not implement runtime operation Stop")
}

func (p *SchemaOnlyProvider) InitLogging(context.Context) {
	panic("schemaOnlyProvider does not implement runtime operation InitLogging")
}

func (p *SchemaOnlyProvider) NewDestroyDiff(context.Context, string, shim.TimeoutOptions) shim.InstanceDiff {
	panic("schemaOnlyProvider does not implement runtime operation NewDestroyDiff")
}

func (p *SchemaOnlyProvider) NewResourceConfig(context.Context, map[string]interface{}) shim.ResourceConfig {
	panic("schemaOnlyProvider does not implement runtime operation ResourceConfig")
}

func (p *SchemaOnlyProvider) IsSet(context.Context, interface{}) ([]interface{}, bool) {
	panic("schemaOnlyProvider does not implement runtime operation IsSet")
}
