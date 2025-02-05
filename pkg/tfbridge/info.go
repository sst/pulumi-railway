// Copyright 2016-2023, Pulumi Corporation.
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

package tfbridge

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/blang/semver"
	"golang.org/x/net/context"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/schema"
	"github.com/pulumi/pulumi-terraform-bridge/v3/unstable/logging"
)

const (
	MPL20LicenseType      TFProviderLicense = "MPL 2.0"
	MITLicenseType        TFProviderLicense = "MIT"
	Apache20LicenseType   TFProviderLicense = "Apache 2.0"
	UnlicensedLicenseType TFProviderLicense = "UNLICENSED"
)

// ProviderInfo contains information about a Terraform provider plugin that we will use to generate the Pulumi
// metadata.  It primarily contains a pointer to the Terraform schema, but can also contain specific name translations.
//
//nolint:lll
type ProviderInfo struct {
	P              shim.Provider // the TF provider/schema.
	Name           string        // the TF provider name (e.g. terraform-provider-XXXX).
	ResourcePrefix string        // the prefix on resources the provider exposes, if different to `Name`.
	// GitHubOrg is the segment of the upstream provider's Go module path that comes after GitHubHost and before
	// terraform-provider-${Name}. Defaults to `terraform-providers`.
	//
	// Note that this value should match the require directive for the upstream provider, not any replace directives.
	//
	// For example, GitHubOrg should be set to "my-company" given the following go.mod:
	//
	// require github.com/my-company/terraform-repo-example v1.0.0
	// replace github.com/my-company/terraform-repo-example => github.com/some-fork/terraform-repo-example v1.0.0
	GitHubOrg      string                             // the GitHub org of the provider. Defaults to `terraform-providers`.
	GitHubHost     string                             // the GitHub host for the provider. Defaults to `github.com`.
	Description    string                             // an optional descriptive overview of the package (a default supplied).
	Keywords       []string                           // an optional list of keywords to help discovery of this package. e.g. "category/cloud, category/infrastructure"
	License        string                             // the license, if any, the resulting package has (default is none).
	LogoURL        string                             // an optional URL to the logo of the package
	DisplayName    string                             // the human friendly name of the package used in the Pulumi registry
	Publisher      string                             // the name of the person or organization that authored and published the package.
	Homepage       string                             // the URL to the project homepage.
	Repository     string                             // the URL to the project source code repository.
	Version        string                             // the version of the provider package.
	Config         map[string]*SchemaInfo             // a map of TF name to config schema overrides.
	ExtraConfig    map[string]*ConfigInfo             // a list of Pulumi-only configuration variables.
	Resources      map[string]*ResourceInfo           // a map of TF name to Pulumi name; standard mangling occurs if no entry.
	DataSources    map[string]*DataSourceInfo         // a map of TF name to Pulumi resource info.
	ExtraTypes     map[string]pschema.ComplexTypeSpec // a map of Pulumi token to schema type for extra types.
	ExtraResources map[string]pschema.ResourceSpec    // a map of Pulumi token to schema type for extra resources.
	ExtraFunctions map[string]pschema.FunctionSpec    // a map of Pulumi token to schema type for extra functions.

	// ExtraResourceHclExamples is a slice of additional HCL examples attached to resources which are converted to the
	// relevant target language(s)
	ExtraResourceHclExamples []HclExampler
	// ExtraFunctionHclExamples is a slice of additional HCL examples attached to functions which are converted to the
	// relevant target language(s)
	ExtraFunctionHclExamples []HclExampler
	// IgnoreMappings is a list of TF resources and data sources that are known to be unmapped.
	//
	// These resources/data sources do not generate missing mappings errors and will not be automatically
	// mapped.
	//
	// If there is a mapping in Resources or DataSources, it can override IgnoreMappings. This is common
	// when you need to ignore a datasource but not the resource with the same name, or vice versa.
	IgnoreMappings          []string
	PluginDownloadURL       string             // an optional URL to download the provider binary from.
	JavaScript              *JavaScriptInfo    // optional overlay information for augmented JavaScript code-generation.
	Python                  *PythonInfo        // optional overlay information for augmented Python code-generation.
	Golang                  *GolangInfo        // optional overlay information for augmented Golang code-generation.
	CSharp                  *CSharpInfo        // optional overlay information for augmented C# code-generation.
	Java                    *JavaInfo          // optional overlay information for augmented C# code-generation.
	TFProviderVersion       string             // the version of the TF provider on which this was based
	TFProviderLicense       *TFProviderLicense // license that the TF provider is distributed under. Default `MPL 2.0`.
	TFProviderModuleVersion string             // the Go module version of the provider. Default is unversioned e.g. v1

	// a provider-specific callback to invoke prior to TF Configure
	// Any CheckFailureErrors returned from PreConfigureCallback are converted to
	// CheckFailures and returned as failures instead of errors in CheckConfig
	PreConfigureCallback PreConfigureCallback
	// Any CheckFailureErrors returned from PreConfigureCallbackWithLogger are
	// converted to CheckFailures and returned as failures instead of errors in CheckConfig
	PreConfigureCallbackWithLogger PreConfigureCallbackWithLogger

	// Information for the embedded metadata file.
	//
	// See NewProviderMetadata for in-place construction of a *MetadataInfo.
	// If a provider should be mixed in with the Terraform provider with MuxWith (see below)
	// this field must be initialized.
	MetadataInfo *MetadataInfo

	// Rules that control file discovery and edits for any subset of docs in a provider.
	DocRules *DocRuleInfo

	UpstreamRepoPath string // An optional path that overrides upstream location during docs lookup

	// EXPERIMENTAL: the signature may change in minor releases.
	//
	// If set, allows selecting individual examples to skip generating into the PackageSchema (and eventually API
	// docs). The primary use case for this hook is to ignore problematic or flaky examples temporarily until the
	// underlying issues are resolved and the examples can be rendered correctly.
	SkipExamples func(SkipExamplesArgs) bool

	// EXPERIMENTAL: the signature may change in minor releases.
	//
	// Optional function to post-process the generated schema spec after
	// the bridge completed its original version based on the TF schema.
	// A hook to enable custom schema modifications specific to a provider.
	SchemaPostProcessor func(spec *pschema.PackageSpec)

	// The MuxWith array allows the mixin (muxing) of other providers to the wrapped upstream Terraform provider.
	// With a provider mixin it's possible to add or replace resources and/or functions (data sources) in the wrapped
	// Terraform provider without having to change the upstream code itself. If multiple provider mixins are specified
	// the schema generator in pkg/tfgen will call the GetSpec() method of muxer.Provider in sequence. Thus, if more or two
	// of the mixins define the same resource/function, the last definition will end up in the combined schema of the
	// compiled provider.
	MuxWith []MuxProvider

	// Disables validation of provider-level configuration for Plugin Framework based providers.
	// Hybrid providers that utilize a mixture of Plugin Framework and SDKv2 based resources may
	// opt into this to workaround slowdown in PF validators, since their configuration is
	// already being checked by SDKv2 based validators.
	//
	// See also: pulumi/pulumi-terraform-bridge#1448
	SkipValidateProviderConfigForPluginFramework bool

	// Disables using detailed diff to determine diff changes and falls back on the length of TF Diff Attributes.
	//
	// See https://github.com/pulumi/pulumi-terraform-bridge/issues/1501
	XSkipDetailedDiffForChanges bool

	// Enables generation of a trimmed, runtime-only metadata file
	// to help reduce resource plugin start time
	//
	// See also pulumi/pulumi-terraform-bridge#1524
	GenerateRuntimeMetadata bool
}

// Send logs or status logs to the user.
//
// Logged messages are pre-associated with the resource they are called from.
type Logger interface {
	Log

	// Convert to sending ephemeral status logs to the user.
	Status() Log
}

// The set of logs available to show to the user
type Log interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// Get access to the [Logger] associated with this context.
func GetLogger(ctx context.Context) Logger {
	logger := ctx.Value(logging.CtxKey)
	contract.Assertf(logger != nil, "Cannot call GetLogger on a context that is not equipped with a Logger")
	return newLoggerAdapter(logger)
}

func (info *ProviderInfo) GetConfig() map[string]*SchemaInfo {
	if info.Config != nil {
		return info.Config
	}
	return map[string]*SchemaInfo{}
}

// The function used to produce the set of edit rules for a provider.
//
// For example, if you want to skip default edits, you would use the function:
//
//	func([]DocsEdit) []DocsEdit { return nil }
//
// If you wanted to incorporate custom edits, default edits, and then a check that the
// resulting document is valid, you would use the function:
//
//	func(defaults []DocsEdit) []DocsEdit {
//		return append(customEdits, append(defaults, validityCheck)...)
//	}
type MakeEditRules func(defaults []DocsEdit) []DocsEdit

// DocRuleInfo controls file discovery and edits for any subset of docs in a provider.
type DocRuleInfo struct {
	// The function called to get the set of edit rules to use.
	//
	// defaults represents suggested edit rules. If EditRules is `nil`, defaults is
	// used as is.
	EditRules MakeEditRules

	// A function to suggest alternative file names for a TF element.
	//
	// When the bridge loads the documentation for a resource or a datasource, it
	// infers the name of the file that contains the documentation. AlternativeNames
	// allows you to provide a provider specific extension to the override list.
	//
	// For example, when attempting to find the documentation for the resource token
	// aws_waf_instances, the bridge will check the following files (in order):
	//
	//	"waf_instance.html.markdown"
	//	"waf_instance.markdown"
	//	"waf_instance.html.md"
	//	"waf_instance.md"
	//	"aws_waf_instance.html.markdown"
	//	"aws_waf_instance.markdown"
	//	"aws_waf_instance.html.md"
	//	"aws_waf_instance.md"
	//
	// The bridge will check any file names returned by AlternativeNames before
	// checking it's standard list.
	AlternativeNames func(DocsPathInfo) []string
}

// Information for file lookup.
type DocsPathInfo struct {
	TfToken string
}

type DocsEdit struct {
	// The file name at which this rule applies. File names are matched via filepath.Match.
	//
	// To match all files, supply "*".
	//
	// All 4 of these names will match "waf_instances.html.markdown":
	//
	// - "waf_instances.html.markdown"
	// - "waf_instances.*"
	// - "waf*"
	// - "*"
	//
	// Provider resources are sourced directly from the TF schema, and as such have an
	// empty path.
	Path string
	// The function that performs the edit on the file bytes.
	//
	// Must not be nil.
	Edit func(path string, content []byte) ([]byte, error)
}

// TFProviderLicense is a way to be able to pass a license type for the upstream Terraform provider.
type TFProviderLicense string

// GetResourcePrefix returns the resource prefix for the provider: info.ResourcePrefix
// if that is set, or info.Name if not. This is to avoid unexpected behaviour with providers
// which have no need to set ResourcePrefix following its introduction.
func (info ProviderInfo) GetResourcePrefix() string {
	if info.ResourcePrefix == "" {
		return info.Name
	}

	return info.ResourcePrefix
}

func (info ProviderInfo) GetMetadata() ProviderMetadata {
	info.MetadataInfo.assertValid()
	return info.MetadataInfo.Data
}

func (info ProviderInfo) GetGitHubOrg() string {
	if info.GitHubOrg == "" {
		return "terraform-providers"
	}

	return info.GitHubOrg
}

func (info ProviderInfo) GetGitHubHost() string {
	if info.GitHubHost == "" {
		return "github.com"
	}

	return info.GitHubHost
}

func (info ProviderInfo) GetTFProviderLicense() TFProviderLicense {
	if info.TFProviderLicense == nil {
		return MPL20LicenseType
	}

	return *info.TFProviderLicense
}

func (info ProviderInfo) GetProviderModuleVersion() string {
	if info.TFProviderModuleVersion == "" {
		return "" // there is no such thing as a v1 module - there is just a missing version declaration
	}

	return info.TFProviderModuleVersion
}

// AliasInfo is a partial description of prior named used for a resource. It can be processed in the
// context of a resource creation to determine what the full aliased URN would be.
//
// It can be used by Pulumi resource providers to change the aspects of it (i.e. what module it is
// contained in), without causing resources to be recreated for customers who migrate from the
// original resource to the current resource.
type AliasInfo struct {
	Name    *string
	Type    *string
	Project *string
}

// ResourceOrDataSourceInfo is a shared interface to ResourceInfo and DataSourceInfo mappings
type ResourceOrDataSourceInfo interface {
	GetTok() tokens.Token              // a type token to override the default; "" uses the default.
	GetFields() map[string]*SchemaInfo // a map of custom field names; if a type is missing, uses the default.
	GetDocs() *DocInfo                 // overrides for finding and mapping TF docs.
	ReplaceExamplesSection() bool      // whether we are replacing the upstream TF examples generation
}

// ResourceInfo is a top-level type exported by a provider.  This structure can override the type to generate.  It can
// also give custom metadata for fields, using the SchemaInfo structure below.  Finally, a set of composite keys can be
// given; this is used when Terraform needs more than just the ID to uniquely identify and query for a resource.
type ResourceInfo struct {
	Tok    tokens.Type            // a type token to override the default; "" uses the default.
	Fields map[string]*SchemaInfo // a map of custom field names; if a type is missing, uses the default.

	// Deprecated: IDFields is not currently used and will be removed in the next major version of
	// pulumi-terraform-bridge. See [ComputeID].
	IDFields []string

	// list of parameters that we can trust that any change will allow a createBeforeDelete
	UniqueNameFields    []string
	Docs                *DocInfo    // overrides for finding and mapping TF docs.
	DeleteBeforeReplace bool        // if true, Pulumi will delete before creating new replacement resources.
	Aliases             []AliasInfo // aliases for this resources, if any.
	DeprecationMessage  string      // message to use in deprecation warning
	CSharpName          string      // .NET-specific name

	// Optional hook to run before upgrading the state. TODO[pulumi/pulumi-terraform-bridge#864] this is currently
	// only supported for Plugin-Framework based providers.
	PreStateUpgradeHook PreStateUpgradeHook

	// An experimental way to augment the Check function in the Pulumi life cycle.
	PreCheckCallback PreCheckCallback

	// Resource operations such as Create, Read, and Update return the resource outputs to be
	// recored in Pulumi statefile. TransformOutputs provides the last chance to edit these
	// outputs before they are stored. In particular, it can be used as a last resort hook to
	// make corrections in the default translation of the resource state from TF to Pulumi.
	// Should be used sparingly.
	TransformOutputs PropertyTransform

	// Check, Diff, Read, Update and Delete refer to old inputs sourced from the
	// Pulumi statefile. TransformFromState lets providers edit these outputs before they
	// are accessed by other provider functions or by terraform. In particular, it can
	// be used to perform upgrades on old pulumi state.  Should be used sparingly.
	TransformFromState PropertyTransform

	// Customizes inferring resource identity from state.
	//
	// The vast majority of resources define an "id" field that is recognized as the resource
	// identity. This is the default behavior when ComputeID is nil. There are some exceptions,
	// however, such as the RandomBytes resource, that base identity on a different field
	// ("base64" in the case of RandomBytes). ComputeID customization option supports such
	// resources. It is called in Create(preview=false) and Read provider methods.
	//
	// This option is currently only supported for Plugin Framework based resources.
	//
	// To delegate the resource ID to another string field in state, use the helper function
	// [DelegateIDField].
	ComputeID ComputeID
}

type ComputeID = func(ctx context.Context, state resource.PropertyMap) (resource.ID, error)

type PropertyTransform = func(context.Context, resource.PropertyMap) (resource.PropertyMap, error)

type PreCheckCallback = func(
	ctx context.Context, config resource.PropertyMap, meta resource.PropertyMap,
) (resource.PropertyMap, error)

// GetTok returns a resource type token
func (info *ResourceInfo) GetTok() tokens.Token { return tokens.Token(info.Tok) }

// GetFields returns information about the resource's custom fields
func (info *ResourceInfo) GetFields() map[string]*SchemaInfo {
	if info == nil {
		return nil
	}
	return info.Fields
}

// GetDocs returns a resource docs override from the Pulumi provider
func (info *ResourceInfo) GetDocs() *DocInfo { return info.Docs }

// ReplaceExamplesSection returns whether to replace the upstream examples with our own source
func (info *ResourceInfo) ReplaceExamplesSection() bool {
	return info.Docs != nil && info.Docs.ReplaceExamplesSection
}

// DataSourceInfo can be used to override a data source's standard name mangling and argument/return information.
type DataSourceInfo struct {
	Tok                tokens.ModuleMember
	Fields             map[string]*SchemaInfo
	Docs               *DocInfo // overrides for finding and mapping TF docs.
	DeprecationMessage string   // message to use in deprecation warning
}

// GetTok returns a datasource type token
func (info *DataSourceInfo) GetTok() tokens.Token { return tokens.Token(info.Tok) }

// GetFields returns information about the datasource's custom fields
func (info *DataSourceInfo) GetFields() map[string]*SchemaInfo {
	if info == nil {
		return nil
	}
	return info.Fields
}

// GetDocs returns a datasource docs override from the Pulumi provider
func (info *DataSourceInfo) GetDocs() *DocInfo { return info.Docs }

// ReplaceExamplesSection returns whether to replace the upstream examples with our own source
func (info *DataSourceInfo) ReplaceExamplesSection() bool {
	return info.Docs != nil && info.Docs.ReplaceExamplesSection
}

// SchemaInfo contains optional name transformations to apply.
type SchemaInfo struct {
	// a name to override the default; "" uses the default.
	Name string

	// a name to override the default when targeting C#; "" uses the default.
	CSharpName string

	// a type to override the default; "" uses the default.
	Type tokens.Type

	// alternative types that can be used instead of the override.
	AltTypes []tokens.Type

	// a type to override when the property is a nested structure.
	NestedType tokens.Type

	// an optional idemponent transformation, applied before passing to TF.
	Transform Transformer

	// a schema override for elements for arrays, maps, and sets.
	Elem *SchemaInfo

	// a map of custom field names; if a type is missing, the default is used.
	Fields map[string]*SchemaInfo

	// a map of asset translation information, if this is an asset.
	Asset *AssetTranslation

	// an optional default directive to be applied if a value is missing.
	Default *DefaultInfo

	// to override whether a property is stable or not.
	Stable *bool

	// to override whether this property should project as a scalar or array.
	MaxItemsOne *bool

	// to remove empty object array elements
	SuppressEmptyMapElements *bool

	// this will make the parameter as computed and not allow the user to set it
	MarkAsComputedOnly *bool

	// this will make the parameter optional in the schema
	MarkAsOptional *bool

	// the deprecation message for the property
	DeprecationMessage string

	// whether a change in the configuration would force a new resource
	ForceNew *bool

	// Controls whether a change in the provider configuration should trigger a provider
	// replacement. While there is no matching concept in TF, Pulumi supports replacing explicit
	// providers and cascading the replacement to all resources provisioned with the given
	// provider configuration.
	//
	// This property is only relevant for [ProviderInfo.Config] properties.
	ForcesProviderReplace *bool

	// whether or not this property has been removed from the Terraform schema
	Removed bool

	// if set, this property will not be added to the schema and no bindings will be generated for it
	Omit bool

	// whether or not to treat this property as secret
	Secret *bool
}

// ConfigInfo represents a synthetic configuration variable that is Pulumi-only, and not passed to Terraform.
type ConfigInfo struct {
	// Info is the Pulumi schema for this variable.
	Info *SchemaInfo
	// Schema is the Terraform schema for this variable.
	Schema shim.Schema
}

// Transformer is given the option to transform a value in situ before it is processed by the bridge. This
// transformation must be deterministic and idempotent, and any value produced by this transformation must
// be a legal alternative input value. A good example is a resource that accepts either a string or
// JSON-stringable map; a resource provider may opt to store the raw string, but let users pass in maps as
// a convenience mechanism, and have the transformer stringify them on the fly. This is safe to do because
// the raw string is still accepted as a possible input value.
type Transformer func(resource.PropertyValue) (resource.PropertyValue, error)

// DocInfo contains optional overrides for finding and mapping TF docs.
type DocInfo struct {
	Source                         string // an optional override to locate TF docs; "" uses the default.
	Markdown                       []byte // an optional override for the source markdown.
	IncludeAttributesFrom          string // optionally include attributes from another raw resource for docs.
	IncludeArgumentsFrom           string // optionally include arguments from another raw resource for docs.
	IncludeAttributesFromArguments string // optionally include attributes from another raw resource's arguments.
	ImportDetails                  string // Overwrite for import instructions

	// Replace examples with the contents of a specific document
	// this document will satisfy the criteria `docs/pulumiToken.md`
	// The examples need to wrapped in the correct shortcodes
	ReplaceExamplesSection bool

	// Don't error when this doc is missing.
	//
	// This applies when PULUMI_MISSING_DOCS_ERROR="true".
	AllowMissing bool
}

// GetImportDetails returns a string of import instructions defined in the Pulumi provider. Defaults to empty.
func (info *DocInfo) GetImportDetails() string { return info.ImportDetails }

// HasDefault returns true if there is a default value for this property.
func (info SchemaInfo) HasDefault() bool {
	return info.Default != nil
}

// DefaultInfo lets fields get default values at runtime, before they are even passed to Terraform.
type DefaultInfo struct {
	// AutoNamed is true if this default represents an autogenerated name.
	AutoNamed bool
	// Config uses a configuration variable from this package as the default value.
	Config string

	// Deprecated. Use ComputeDefault.
	From func(res *PulumiResource) (interface{}, error)

	// ComputeDefault specifies how to compute a default value for the given property by consulting other properties
	// such as the resource's URN. See [ComputeDefaultOptions] for all available information.
	ComputeDefault func(ctx context.Context, opts ComputeDefaultOptions) (interface{}, error)

	// Value injects a raw literal value as the default.
	// Note that only simple types such as string, int and boolean are currently supported here.
	// Structs, slices and maps are not yet supported.
	Value interface{}
	// EnvVars to use for defaults. If none of these variables have values at runtime, the value of `Value` (if any)
	// will be used as the default.
	EnvVars []string
}

// Configures [DefaultInfo.ComputeDefault].
type ComputeDefaultOptions struct {
	// URN identifying the Resource. Set when computing default properties for a Resource, and unset for functions.
	URN resource.URN

	// Property map before computing the defaults.
	Properties resource.PropertyMap

	// Property map representing prior state, only set for non-Create Resource operations.
	PriorState resource.PropertyMap

	// PriorValue represents the last value of the current property in PriorState. It will have zero value if there
	// is no PriorState or if the property did not have a value in PriorState.
	PriorValue resource.PropertyValue

	// The engine provides a stable seed useful for generating random values consistently. This guarantees, for
	// example, that random values generated across "pulumi preview" and "pulumi up" in the same deployment are
	// consistent. This currently is only available for resource changes.
	Seed []byte
}

// PulumiResource is just a little bundle that carries URN, seed and properties around.
type PulumiResource struct {
	URN        resource.URN
	Properties resource.PropertyMap
	Seed       []byte
}

// OverlayInfo contains optional overlay information.  Each info has a 1:1 correspondence with a module and
// permits extra files to be included from the overlays/ directory when building up packs/.  This allows augmented
// code-generation for convenient things like helper functions, modules, and gradual typing.
type OverlayInfo struct {
	DestFiles []string                // Additional files to include in the index file. Must exist in the destination.
	Modules   map[string]*OverlayInfo // extra modules to inject into the structure.
}

// JavaScriptInfo contains optional overlay information for Python code-generation.
type JavaScriptInfo struct {
	PackageName       string            // Custom name for the NPM package.
	Dependencies      map[string]string // NPM dependencies to add to package.json.
	DevDependencies   map[string]string // NPM dev-dependencies to add to package.json.
	PeerDependencies  map[string]string // NPM peer-dependencies to add to package.json.
	Resolutions       map[string]string // NPM resolutions to add to package.json.
	Overlay           *OverlayInfo      // optional overlay information for augmented code-generation.
	TypeScriptVersion string            // A specific version of TypeScript to include in package.json.
	PluginName        string            // The name of the plugin, which might be
	// different from the package name.  The version of the plugin, which might be
	// different from the version of the package.
	PluginVersion string

	// A map containing overrides for module names to package names.
	ModuleToPackage map[string]string

	// An indicator for whether the package contains enums.
	ContainsEnums bool

	// A map allowing you to map the name of a provider to the name of the module encapsulating the provider.
	ProviderNameToModuleName map[string]string

	// Additional files to include in TypeScript compilation. These paths are added to the `files` section of the
	// generated `tsconfig.json`. A typical use case for this is compiling hand-authored unit test files that check
	// the generated code.
	ExtraTypeScriptFiles []string

	// Determines whether to make single-return-value methods return an output object or the single value.
	LiftSingleValueMethodReturns bool

	// Respect the Pkg.Version field in the schema
	RespectSchemaVersion bool

	// Experimental flag that permits `import type *` style code to be generated to optimize startup time of
	// programs consuming the provider by minimizing the set of Node modules loaded at startup. Turning this on may
	// currently generate non-compiling code for some providers; but if the code compiles it is safe to use. Also,
	// turning this on requires TypeScript 3.8 or higher to compile the generated code.
	UseTypeOnlyReferences bool
}

// PythonInfo contains optional overlay information for Python code-generation.
type PythonInfo struct {
	Requires      map[string]string // Pip install_requires information.
	Overlay       *OverlayInfo      // optional overlay information for augmented code-generation.
	UsesIOClasses bool              // Deprecated: No longer required, all providers use IO classes.
	PackageName   string            // Name of the Python package to generate

	// PythonRequires determines the Python versions that the generated provider supports
	PythonRequires string

	// Optional overrides for Pulumi module names
	//
	//    { "flowcontrol.apiserver.k8s.io/v1alpha1": "flowcontrol/v1alpha1" }
	//
	ModuleNameOverrides map[string]string

	// Determines whether to make single-return-value methods return an output object or the single value.
	LiftSingleValueMethodReturns bool

	// Respect the Pkg.Version field for emitted code.
	RespectSchemaVersion bool

	// If enabled, a pyproject.toml file will be generated.
	PyProject struct {
		Enabled bool
	}
}

// GolangInfo contains optional overlay information for Golang code-generation.
type GolangInfo struct {
	GenerateResourceContainerTypes bool         // Generate container types for resources e.g. arrays, maps, pointers etc.
	ImportBasePath                 string       // Base import path for package.
	Overlay                        *OverlayInfo // optional overlay information for augmented code-generation.

	// Module path for go.mod
	//
	//   go get github.com/pulumi/pulumi-aws-native/sdk/go/aws@v0.16.0
	//          ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ module path
	//                                                  ~~~~~~ package path - can be any number of path parts
	//                                                         ~~~~~~~ version
	ModulePath string

	// Explicit package name, which may be different to the import path.
	RootPackageName string

	// Map from module -> package name
	//
	//    { "flowcontrol.apiserver.k8s.io/v1alpha1": "flowcontrol/v1alpha1" }
	//
	ModuleToPackage map[string]string

	// Map from package name -> package alias
	//
	//    { "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/flowcontrol/v1alpha1": "flowcontrolv1alpha1" }
	//
	PackageImportAliases map[string]string

	// The version of the Pulumi SDK used with this provider, e.g. 3.
	// Used to generate doc links for pulumi builtin types. If omitted, the latest SDK version is used.
	PulumiSDKVersion int

	// Feature flag to disable generating `$fnOutput` invoke
	// function versions to save space.
	DisableFunctionOutputVersions bool

	// Determines whether to make single-return-value methods return an output struct or the value.
	LiftSingleValueMethodReturns bool

	// Feature flag to disable generating input type registration. This is a
	// space saving measure.
	DisableInputTypeRegistrations bool

	// Feature flag to disable generating Pulumi object default functions. This is a
	// space saving measure.
	DisableObjectDefaults bool

	// GenerateExtraInputTypes determines whether or not the code generator generates input (and output) types for
	// all plain types, instead of for only types that are used as input/output types.
	GenerateExtraInputTypes bool

	// omitExtraInputTypes determines whether the code generator generates input (and output) types
	// for all plain types, instead of for only types that are used as input/output types.
	OmitExtraInputTypes bool

	// Respect the Pkg.Version field for emitted code.
	RespectSchemaVersion bool

	// InternalDependencies are blank imports that are emitted in the SDK so that `go mod tidy` does not remove the
	// associated module dependencies from the SDK's go.mod.
	InternalDependencies []string

	// Specifies how to handle generating a variant of the SDK that uses generics.
	// Allowed values are the following:
	// - "none" (default): do not generate a generics variant of the SDK
	// - "side-by-side": generate a side-by-side generics variant of the SDK under the x subdirectory
	// - "only-generics": generate a generics variant of the SDK only
	Generics string
}

// CSharpInfo contains optional overlay information for C# code-generation.
type CSharpInfo struct {
	PackageReferences map[string]string // NuGet package reference information.
	Overlay           *OverlayInfo      // optional overlay information for augmented code-generation.
	Namespaces        map[string]string // Known .NET namespaces with proper capitalization.
	RootNamespace     string            // The root namespace if setting to something other than Pulumi in the package name

	Compatibility          string
	DictionaryConstructors bool
	ProjectReferences      []string

	// Determines whether to make single-return-value methods return an output object or the single value.
	LiftSingleValueMethodReturns bool

	// Allow the Pkg.Version field to filter down to emitted code.
	RespectSchemaVersion bool
}

// See https://github.com/pulumi/pulumi-java/blob/main/pkg/codegen/java/package_info.go#L35C1-L108C1 documenting
// supported options.
type JavaInfo struct {
	BasePackage string // the Base package for the Java SDK

	// If set to "gradle" enables a generation of a basic set of
	// Gradle build files.
	BuildFiles string

	// If non-empty and BuildFiles="gradle", enables the use of a
	// given version of io.github.gradle-nexus.publish-plugin in
	// the generated Gradle build files.
	GradleNexusPublishPluginVersion string

	Packages     map[string]string `json:"packages,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	GradleTest   string            `json:"gradleTest"`
}

// PreConfigureCallback is a function to invoke prior to calling the TF provider Configure
type PreConfigureCallback func(vars resource.PropertyMap, config shim.ResourceConfig) error

// PreConfigureCallbackWithLogger is a function to invoke prior to calling the T
type PreConfigureCallbackWithLogger func(
	ctx context.Context,
	host *provider.HostClient, vars resource.PropertyMap,
	config shim.ResourceConfig,
) error

// The types below are marshallable versions of the schema descriptions associated with a provider. These are used when
// marshalling a provider info as JSON; Note that these types only represent a subset of the information associated
// with a ProviderInfo; thus, a ProviderInfo cannot be round-tripped through JSON.

// MarshallableSchema is the JSON-marshallable form of a Terraform schema.
type MarshallableSchema struct {
	Type               shim.ValueType    `json:"type"`
	Optional           bool              `json:"optional,omitempty"`
	Required           bool              `json:"required,omitempty"`
	Computed           bool              `json:"computed,omitempty"`
	ForceNew           bool              `json:"forceNew,omitempty"`
	Elem               *MarshallableElem `json:"element,omitempty"`
	MaxItems           int               `json:"maxItems,omitempty"`
	MinItems           int               `json:"minItems,omitempty"`
	DeprecationMessage string            `json:"deprecated,omitempty"`
}

// MarshalSchema converts a Terraform schema into a MarshallableSchema.
func MarshalSchema(s shim.Schema) *MarshallableSchema {
	return &MarshallableSchema{
		Type:               s.Type(),
		Optional:           s.Optional(),
		Required:           s.Required(),
		Computed:           s.Computed(),
		ForceNew:           s.ForceNew(),
		Elem:               MarshalElem(s.Elem()),
		MaxItems:           s.MaxItems(),
		MinItems:           s.MinItems(),
		DeprecationMessage: s.Deprecated(),
	}
}

// Unmarshal creates a mostly-initialized Terraform schema from the given MarshallableSchema.
func (m *MarshallableSchema) Unmarshal() shim.Schema {
	return (&schema.Schema{
		Type:       m.Type,
		Optional:   m.Optional,
		Required:   m.Required,
		Computed:   m.Computed,
		ForceNew:   m.ForceNew,
		Elem:       m.Elem.Unmarshal(),
		MaxItems:   m.MaxItems,
		MinItems:   m.MinItems,
		Deprecated: m.DeprecationMessage,
	}).Shim()
}

// MarshallableResource is the JSON-marshallable form of a Terraform resource schema.
type MarshallableResource map[string]*MarshallableSchema

// MarshalResource converts a Terraform resource schema into a MarshallableResource.
func MarshalResource(r shim.Resource) MarshallableResource {
	m := make(MarshallableResource)
	if r.Schema() == nil {
		return m
	}
	r.Schema().Range(func(k string, v shim.Schema) bool {
		m[k] = MarshalSchema(v)
		return true
	})
	return m
}

// Unmarshal creates a mostly-initialized Terraform resource schema from the given MarshallableResource.
func (m MarshallableResource) Unmarshal() shim.Resource {
	s := schema.SchemaMap{}
	for k, v := range m {
		s[k] = v.Unmarshal()
	}
	return (&schema.Resource{Schema: s}).Shim()
}

// MarshallableElem is the JSON-marshallable form of a Terraform schema's element field.
type MarshallableElem struct {
	Schema   *MarshallableSchema  `json:"schema,omitempty"`
	Resource MarshallableResource `json:"resource,omitempty"`
}

// MarshalElem converts a Terraform schema's element field into a MarshallableElem.
func MarshalElem(e interface{}) *MarshallableElem {
	switch v := e.(type) {
	case shim.Schema:
		return &MarshallableElem{Schema: MarshalSchema(v)}
	case shim.Resource:
		return &MarshallableElem{Resource: MarshalResource(v)}
	default:
		contract.Assertf(e == nil, "unexpected schema element of type %T", e)
		return nil
	}
}

// Unmarshal creates a Terraform schema element from a MarshallableElem.
func (m *MarshallableElem) Unmarshal() interface{} {
	switch {
	case m == nil:
		return nil
	case m.Schema != nil:
		return m.Schema.Unmarshal()
	default:
		// m.Resource might be nil in which case it was empty when marshalled. But Unmarshal can be called on
		// nil and returns something sensible.
		return m.Resource.Unmarshal()
	}
}

// MarshallableProvider is the JSON-marshallable form of a Terraform provider schema.
type MarshallableProvider struct {
	Schema      map[string]*MarshallableSchema  `json:"schema,omitempty"`
	Resources   map[string]MarshallableResource `json:"resources,omitempty"`
	DataSources map[string]MarshallableResource `json:"dataSources,omitempty"`
}

// MarshalProvider converts a Terraform provider schema into a MarshallableProvider.
func MarshalProvider(p shim.Provider) *MarshallableProvider {
	if p == nil {
		return nil
	}

	config := make(map[string]*MarshallableSchema)
	p.Schema().Range(func(k string, v shim.Schema) bool {
		config[k] = MarshalSchema(v)
		return true
	})
	resources := make(map[string]MarshallableResource)
	p.ResourcesMap().Range(func(k string, v shim.Resource) bool {
		resources[k] = MarshalResource(v)
		return true
	})
	dataSources := make(map[string]MarshallableResource)
	p.DataSourcesMap().Range(func(k string, v shim.Resource) bool {
		dataSources[k] = MarshalResource(v)
		return true
	})
	return &MarshallableProvider{
		Schema:      config,
		Resources:   resources,
		DataSources: dataSources,
	}
}

// Unmarshal creates a mostly-initialized Terraform provider schema from a MarshallableProvider
func (m *MarshallableProvider) Unmarshal() shim.Provider {
	if m == nil {
		return nil
	}

	config := schema.SchemaMap{}
	for k, v := range m.Schema {
		config[k] = v.Unmarshal()
	}
	resources := schema.ResourceMap{}
	for k, v := range m.Resources {
		resources[k] = v.Unmarshal()
	}
	dataSources := schema.ResourceMap{}
	for k, v := range m.DataSources {
		dataSources[k] = v.Unmarshal()
	}
	return (&schema.Provider{
		Schema:         config,
		ResourcesMap:   resources,
		DataSourcesMap: dataSources,
	}).Shim()
}

// MarshallableSchemaInfo is the JSON-marshallable form of a Pulumi SchemaInfo value.
type MarshallableSchemaInfo struct {
	Name        string                             `json:"name,omitempty"`
	CSharpName  string                             `json:"csharpName,omitempty"`
	Type        tokens.Type                        `json:"typeomitempty"`
	AltTypes    []tokens.Type                      `json:"altTypes,omitempty"`
	Elem        *MarshallableSchemaInfo            `json:"element,omitempty"`
	Fields      map[string]*MarshallableSchemaInfo `json:"fields,omitempty"`
	Asset       *AssetTranslation                  `json:"asset,omitempty"`
	Default     *MarshallableDefaultInfo           `json:"default,omitempty"`
	MaxItemsOne *bool                              `json:"maxItemsOne,omitempty"`
	Deprecated  string                             `json:"deprecated,omitempty"`
	ForceNew    *bool                              `json:"forceNew,omitempty"`
	Secret      *bool                              `json:"secret,omitempty"`
}

// MarshalSchemaInfo converts a Pulumi SchemaInfo value into a MarshallableSchemaInfo value.
func MarshalSchemaInfo(s *SchemaInfo) *MarshallableSchemaInfo {
	if s == nil {
		return nil
	}

	fields := make(map[string]*MarshallableSchemaInfo)
	for k, v := range s.Fields {
		fields[k] = MarshalSchemaInfo(v)
	}
	return &MarshallableSchemaInfo{
		Name:        s.Name,
		CSharpName:  s.CSharpName,
		Type:        s.Type,
		AltTypes:    s.AltTypes,
		Elem:        MarshalSchemaInfo(s.Elem),
		Fields:      fields,
		Asset:       s.Asset,
		Default:     MarshalDefaultInfo(s.Default),
		MaxItemsOne: s.MaxItemsOne,
		Deprecated:  s.DeprecationMessage,
		ForceNew:    s.ForceNew,
		Secret:      s.Secret,
	}
}

// Unmarshal creates a mostly-=initialized Pulumi SchemaInfo value from the given MarshallableSchemaInfo.
func (m *MarshallableSchemaInfo) Unmarshal() *SchemaInfo {
	if m == nil {
		return nil
	}

	fields := make(map[string]*SchemaInfo)
	for k, v := range m.Fields {
		fields[k] = v.Unmarshal()
	}
	return &SchemaInfo{
		Name:               m.Name,
		CSharpName:         m.CSharpName,
		Type:               m.Type,
		AltTypes:           m.AltTypes,
		Elem:               m.Elem.Unmarshal(),
		Fields:             fields,
		Asset:              m.Asset,
		Default:            m.Default.Unmarshal(),
		MaxItemsOne:        m.MaxItemsOne,
		DeprecationMessage: m.Deprecated,
		ForceNew:           m.ForceNew,
		Secret:             m.Secret,
	}
}

// MarshallableDefaultInfo is the JSON-marshallable form of a Pulumi DefaultInfo value.
type MarshallableDefaultInfo struct {
	AutoNamed bool        `json:"autonamed,omitempty"`
	IsFunc    bool        `json:"isFunc,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	EnvVars   []string    `json:"envvars,omitempty"`
}

// MarshalDefaultInfo converts a Pulumi DefaultInfo value into a MarshallableDefaultInfo value.
func MarshalDefaultInfo(d *DefaultInfo) *MarshallableDefaultInfo {
	if d == nil {
		return nil
	}

	return &MarshallableDefaultInfo{
		AutoNamed: d.AutoNamed,
		IsFunc:    d.From != nil || d.ComputeDefault != nil,
		Value:     d.Value,
		EnvVars:   d.EnvVars,
	}
}

// Unmarshal creates a mostly-initialized Pulumi DefaultInfo value from the given MarshallableDefaultInfo.
func (m *MarshallableDefaultInfo) Unmarshal() *DefaultInfo {
	if m == nil {
		return nil
	}

	defInfo := &DefaultInfo{
		AutoNamed: m.AutoNamed,
		Value:     m.Value,
		EnvVars:   m.EnvVars,
	}

	if m.IsFunc {
		defInfo.ComputeDefault = func(context.Context, ComputeDefaultOptions) (interface{}, error) {
			panic("transforms cannot be run on unmarshaled DefaultInfo values")
		}
	}
	return defInfo
}

// MarshallableResourceInfo is the JSON-marshallable form of a Pulumi ResourceInfo value.
type MarshallableResourceInfo struct {
	Tok        tokens.Type                        `json:"tok"`
	CSharpName string                             `json:"csharpName,omitempty"`
	Fields     map[string]*MarshallableSchemaInfo `json:"fields"`

	// Deprecated: IDFields is not currently used and will be deprecated in the next major version of
	// pulumi-terraform-bridge.
	IDFields []string `json:"idFields"`
}

// MarshalResourceInfo converts a Pulumi ResourceInfo value into a MarshallableResourceInfo value.
func MarshalResourceInfo(r *ResourceInfo) *MarshallableResourceInfo {
	fields := make(map[string]*MarshallableSchemaInfo)
	for k, v := range r.Fields {
		fields[k] = MarshalSchemaInfo(v)
	}
	return &MarshallableResourceInfo{
		Tok:        r.Tok,
		CSharpName: r.CSharpName,
		Fields:     fields,
		IDFields:   r.IDFields,
	}
}

// Unmarshal creates a mostly-=initialized Pulumi ResourceInfo value from the given MarshallableResourceInfo.
func (m *MarshallableResourceInfo) Unmarshal() *ResourceInfo {
	fields := make(map[string]*SchemaInfo)
	for k, v := range m.Fields {
		fields[k] = v.Unmarshal()
	}
	return &ResourceInfo{
		Tok:        m.Tok,
		Fields:     fields,
		IDFields:   m.IDFields,
		CSharpName: m.CSharpName,
	}
}

// MarshallableDataSourceInfo is the JSON-marshallable form of a Pulumi DataSourceInfo value.
type MarshallableDataSourceInfo struct {
	Tok    tokens.ModuleMember                `json:"tok"`
	Fields map[string]*MarshallableSchemaInfo `json:"fields"`
}

// MarshalDataSourceInfo converts a Pulumi DataSourceInfo value into a MarshallableDataSourceInfo value.
func MarshalDataSourceInfo(d *DataSourceInfo) *MarshallableDataSourceInfo {
	fields := make(map[string]*MarshallableSchemaInfo)
	for k, v := range d.Fields {
		fields[k] = MarshalSchemaInfo(v)
	}
	return &MarshallableDataSourceInfo{
		Tok:    d.Tok,
		Fields: fields,
	}
}

// Unmarshal creates a mostly-=initialized Pulumi DataSourceInfo value from the given MarshallableDataSourceInfo.
func (m *MarshallableDataSourceInfo) Unmarshal() *DataSourceInfo {
	fields := make(map[string]*SchemaInfo)
	for k, v := range m.Fields {
		fields[k] = v.Unmarshal()
	}
	return &DataSourceInfo{
		Tok:    m.Tok,
		Fields: fields,
	}
}

// MarshallableProviderInfo is the JSON-marshallable form of a Pulumi ProviderInfo value.
type MarshallableProviderInfo struct {
	Provider          *MarshallableProvider                  `json:"provider"`
	Name              string                                 `json:"name"`
	Version           string                                 `json:"version"`
	Config            map[string]*MarshallableSchemaInfo     `json:"config,omitempty"`
	Resources         map[string]*MarshallableResourceInfo   `json:"resources,omitempty"`
	DataSources       map[string]*MarshallableDataSourceInfo `json:"dataSources,omitempty"`
	TFProviderVersion string                                 `json:"tfProviderVersion,omitempty"`
}

// MarshalProviderInfo converts a Pulumi ProviderInfo value into a MarshallableProviderInfo value.
func MarshalProviderInfo(p *ProviderInfo) *MarshallableProviderInfo {
	config := make(map[string]*MarshallableSchemaInfo)
	for k, v := range p.Config {
		config[k] = MarshalSchemaInfo(v)
	}
	resources := make(map[string]*MarshallableResourceInfo)
	for k, v := range p.Resources {
		resources[k] = MarshalResourceInfo(v)
	}
	dataSources := make(map[string]*MarshallableDataSourceInfo)
	for k, v := range p.DataSources {
		dataSources[k] = MarshalDataSourceInfo(v)
	}

	info := MarshallableProviderInfo{
		Provider:          MarshalProvider(p.P),
		Name:              p.Name,
		Version:           p.Version,
		Config:            config,
		Resources:         resources,
		DataSources:       dataSources,
		TFProviderVersion: p.TFProviderVersion,
	}

	return &info
}

// Unmarshal creates a mostly-=initialized Pulumi ProviderInfo value from the given MarshallableProviderInfo.
func (m *MarshallableProviderInfo) Unmarshal() *ProviderInfo {
	config := make(map[string]*SchemaInfo)
	for k, v := range m.Config {
		config[k] = v.Unmarshal()
	}
	resources := make(map[string]*ResourceInfo)
	for k, v := range m.Resources {
		resources[k] = v.Unmarshal()
	}
	dataSources := make(map[string]*DataSourceInfo)
	for k, v := range m.DataSources {
		dataSources[k] = v.Unmarshal()
	}

	info := ProviderInfo{
		P:                 m.Provider.Unmarshal(),
		Name:              m.Name,
		Version:           m.Version,
		Config:            config,
		Resources:         resources,
		DataSources:       dataSources,
		TFProviderVersion: m.TFProviderVersion,
	}

	return &info
}

// Calculates the major version of a go sdk
// go module paths only care about appending a version when the version is
// 2 or greater. github.com/org/my-repo/sdk/v1/go is not a valid
// go module path but github.com/org/my-repo/sdk/v2/go is
func GetModuleMajorVersion(version string) string {
	var majorVersion string
	sver, err := semver.ParseTolerant(version)
	if err != nil {
		panic(err)
	}
	if sver.Major > 1 {
		majorVersion = fmt.Sprintf("v%d", sver.Major)
	}
	return majorVersion
}

// MakeMember manufactures a type token for the package and the given module and type.
//
// Deprecated: Use MakeResource or call into the `tokens` module in
// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens" directly.
func MakeMember(pkg string, mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(pkg + ":" + mod + ":" + mem)
}

// MakeType manufactures a type token for the package and the given module and type.
//
// Deprecated: Use MakeResource or call into the `tokens` module in
// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens" directly.
func MakeType(pkg string, mod string, typ string) tokens.Type {
	return tokens.Type(MakeMember(pkg, mod, typ))
}

// MakeDataSource manufactures a standard Pulumi function token given a package, module, and data source name.  It
// automatically uses the main package and names the file by simply lower casing the data source's
// first character.
//
// Invalid inputs panic.
func MakeDataSource(pkg string, mod string, name string) tokens.ModuleMember {
	contract.Assertf(tokens.IsName(name), "invalid datasource name: '%s'", name)
	modT := makeModule(pkg, mod, name)
	return tokens.NewModuleMemberToken(modT, tokens.ModuleMemberName(name))
}

// makeModule manufactures a standard pulumi module from a (pkg, mod, member) triple.
//
// For example:
//
//	(pkg, mod, Resource) => pkg:mod/resource
//
// Invalid inputs panic.
func makeModule(pkg, mod, member string) tokens.Module {
	mod += "/" + string(unicode.ToLower(rune(member[0]))) + member[1:]
	contract.Assertf(tokens.IsQName(pkg), "invalid pkg name: '%s'", pkg)
	pkgT := tokens.NewPackageToken(tokens.PackageName(pkg))
	contract.Assertf(tokens.IsQName(mod), "invalid module name: '%s'", mod)
	return tokens.NewModuleToken(pkgT, tokens.ModuleName(mod))
}

// MakeResource manufactures a standard resource token given a package, module and resource name.  It
// automatically uses the main package and names the file by simply lower casing the resource's
// first character.
//
// Invalid inputs panic.
func MakeResource(pkg string, mod string, res string) tokens.Type {
	contract.Assertf(tokens.IsName(res), "invalid resource name: '%s'", res)
	modT := makeModule(pkg, mod, res)
	return tokens.NewTypeToken(modT, tokens.TypeName(res))
}

// BoolRef returns a reference to the bool argument.
func BoolRef(b bool) *bool {
	return &b
}

// StringValue gets a string value from a property map if present, else ""
func StringValue(vars resource.PropertyMap, prop resource.PropertyKey) string {
	val, ok := vars[prop]
	if ok && val.IsString() {
		return val.StringValue()
	}
	return ""
}

// ManagedByPulumi is a default used for some managed resources, in the absence of something more meaningful.
var ManagedByPulumi = &DefaultInfo{Value: "Managed by Pulumi"}

// ConfigStringValue gets a string value from a property map, then from environment vars; defaults to empty string ""
func ConfigStringValue(vars resource.PropertyMap, prop resource.PropertyKey, envs []string) string {
	val, ok := vars[prop]
	if ok && val.IsString() {
		return val.StringValue()
	}
	for _, env := range envs {
		val, ok := os.LookupEnv(env)
		if ok {
			return val
		}
	}
	return ""
}

// ConfigArrayValue takes an array value from a property map, then from environment vars; defaults to an empty array
func ConfigArrayValue(vars resource.PropertyMap, prop resource.PropertyKey, envs []string) []string {
	val, ok := vars[prop]
	var vals []string
	if ok && val.IsArray() {
		for _, v := range val.ArrayValue() {
			vals = append(vals, v.StringValue())
		}
		return vals
	}

	for _, env := range envs {
		val, ok := os.LookupEnv(env)
		if ok {
			return strings.Split(val, ";")
		}
	}
	return vals
}

// ConfigBoolValue takes a bool value from a property map, then from environment vars; defaults to false
func ConfigBoolValue(vars resource.PropertyMap, prop resource.PropertyKey, envs []string) bool {
	val, ok := vars[prop]
	if ok && val.IsBool() {
		return val.BoolValue()
	}
	for _, env := range envs {
		val, ok := os.LookupEnv(env)
		if ok && val == "true" {
			return true
		}
	}
	return false
}

// If specified, the hook will run just prior to executing Terraform state upgrades to transform the resource state as
// stored in Pulumi. It can be used to perform idempotent corrections on corrupt state and to compensate for
// Terraform-level state upgrade not working as expected. Returns the corrected resource state and version. To be used
// with care.
//
// See also: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/resource/schema#Schema.Version
type PreStateUpgradeHook = func(PreStateUpgradeHookArgs) (int64, resource.PropertyMap, error)

type PreStateUpgradeHookArgs struct {
	PriorState              resource.PropertyMap
	PriorStateSchemaVersion int64
	ResourceSchemaVersion   int64
}

// EXPERIMENTAL: the signature may change in minor releases.
type SkipExamplesArgs struct {
	// token will be a resource, function, or type token from Pulumi Package Schema. For instance,
	// "aws:acm/certificate:Certificate" would indicate the example pertains to the Certificate resource in the AWS
	// provider.
	Token string

	// examplePath will provide even more information on where the example is found. For instance,
	// "#/resources/aws:acm/certificate:Certificate/arn" would encode that the example pertains to the arn property
	// of the Certificate resource in the AWS provider.
	ExamplePath string
}

func DelegateIDField(field resource.PropertyKey, providerName, repoURL string) ComputeID {
	return func(ctx context.Context, state resource.PropertyMap) (resource.ID, error) {
		err := func(msg string, a ...any) error {
			return delegateIDFieldError{
				msg:          fmt.Sprintf(msg, a...),
				providerName: providerName,
				repoURL:      repoURL,
			}
		}
		fieldValue, ok := state[field]
		if !ok {
			return "", err("Could not find required property '%s' in state", field)
		}

		contract.Assertf(
			!fieldValue.IsComputed() && (!fieldValue.IsOutput() || fieldValue.OutputValue().Known),
			"ComputeID is only called during when preview=false, so we should never need to "+
				"deal with computed properties",
		)

		if fieldValue.IsSecret() || (fieldValue.IsOutput() && fieldValue.OutputValue().Secret) {
			msg := fmt.Sprintf("Setting non-secret resource ID as '%s' (which is secret)", field)
			GetLogger(ctx).Warn(msg)
			if fieldValue.IsSecret() {
				fieldValue = fieldValue.SecretValue().Element
			} else {
				fieldValue = fieldValue.OutputValue().Element
			}
		}

		if !fieldValue.IsString() {
			return "", err("Expected '%s' property to be a string, found %s",
				field, fieldValue.TypeString())
		}

		return resource.ID(fieldValue.StringValue()), nil
	}
}

type delegateIDFieldError struct {
	msg                   string
	providerName, repoURL string
}

func (err delegateIDFieldError) Error() string {
	return fmt.Sprintf("%s. This is an error in %s resource provider, please report at %s",
		err.msg, err.providerName, err.repoURL)
}

func (err delegateIDFieldError) Is(target error) bool {
	target, ok := target.(delegateIDFieldError)
	return ok && err == target
}
