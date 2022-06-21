module github.com/jenkins-x-plugins/jx-updatebot

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/cpuguy83/go-md2man v1.0.10
	github.com/google/go-cmp v0.5.6
	github.com/jenkins-x-plugins/jx-gitops v0.3.13
	github.com/jenkins-x-plugins/jx-pipeline v0.0.152
	github.com/jenkins-x-plugins/jx-promote v0.0.277
	github.com/jenkins-x/go-scm v1.11.5
	github.com/jenkins-x/jx-api/v4 v4.3.6
	github.com/jenkins-x/jx-helpers/v3 v3.2.16
	github.com/jenkins-x/jx-logging/v3 v3.0.7
	github.com/jenkins-x/lighthouse-client v0.0.295
	github.com/pkg/errors v0.9.1
	github.com/roboll/helmfile v0.139.0
	github.com/shurcooL/githubv4 v0.0.0-20191102174205-af46314aec7b
	github.com/spf13/cobra v1.2.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/yargevad/filepathx v0.0.0-20161019152617-907099cb5a62
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	k8s.io/apimachinery v0.21.0
	sigs.k8s.io/kustomize/kyaml v0.10.17
)

replace (
	// helm dependencies
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	// override the go-scm from tekton
	github.com/jenkins-x/go-scm => github.com/jenkins-x/go-scm v1.10.10
	// fix yaml comment parsing issue
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776

	k8s.io/api => k8s.io/api v0.20.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.3
	k8s.io/client-go => k8s.io/client-go v0.20.3

	// fix yaml comment parsing issue
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.6.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0
)

go 1.15
