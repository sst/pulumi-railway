module github.com/pulumi/pulumi-terraform-bridge/v3

go 1.18

require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/go-test/deep v1.0.3
	github.com/golang/glog v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-getter v1.6.1
	github.com/hashicorp/go-hclog v1.2.1
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.4.4
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl/v2 v2.14.0
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.19.0
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi-java/pkg v0.6.0
	github.com/pulumi/pulumi-yaml v0.5.7
	github.com/pulumi/pulumi/pkg/v3 v3.40.2-0.20220924091327-a18d7fb421dc
	github.com/pulumi/pulumi/sdk/v3 v3.40.2-0.20220924091327-a18d7fb421dc
	github.com/pulumi/terraform-diff-reader v0.0.0-20201211191010-ad4715e9285e
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/zclconf/go-cty v1.10.0
	golang.org/x/crypto v0.0.0-20220824171710-5757bc0c5503
	golang.org/x/mod v0.5.0
	golang.org/x/net v0.0.0-20220325170049-de3da57026de
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.6.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/kms v1.1.0 // indirect
	cloud.google.com/go/logging v1.0.0 // indirect
	cloud.google.com/go/storage v1.22.1 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v57.0.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.14.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.20 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.15 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.3 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/aws/aws-sdk-go v1.40.34 // indirect
	github.com/aws/aws-sdk-go-v2 v1.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.7.0 // indirect
	github.com/aws/smithy-go v1.8.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cheggaaa/pb v1.0.29 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/djherbis/times v1.5.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/ettle/strcase v0.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/gax-go/v2 v2.3.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.12.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.6.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.0.0-20220623143253-7d51757b572c // indirect
	github.com/hashicorp/vault/api v1.1.1 // indirect
	github.com/hashicorp/vault/sdk v0.2.1 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/ijc/Gotty v0.0.0-20170406111628-a8b993ba6abd // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pgavlin/goldmark v1.1.33-0.20200616210433-b5eb04559386 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/term v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/posener/complete v1.2.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.3.5 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/tweekmonster/luser v0.0.0-20161003172636-3fa38070dbd7 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/xanzy/ssh-agent v0.3.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/zclconf/go-cty-yaml v1.0.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	gocloud.dev v0.24.0 // indirect
	gocloud.dev/secrets/hashivault v0.24.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220823224334-20c2bfdbfe24 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/api v0.74.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220518221133-4f43b3371335 // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.9-0.20200217094205-6773bdf39b7f // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/frand v1.4.2 // indirect
	sourcegraph.com/sourcegraph/appdash v0.0.0-20211028080628-e2786a622600 // indirect
)

replace github.com/hashicorp/terraform-plugin-sdk/v2 => github.com/pulumi/terraform-plugin-sdk/v2 v2.0.0-20220725190814-23001ad6ec03
