module github.com/jenkins-x-plugins/jx-updatebot

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/cpuguy83/go-md2man v1.0.10
	github.com/jenkins-x-plugins/jx-gitops v0.2.58
	github.com/jenkins-x-plugins/jx-pipeline v0.0.129
	github.com/jenkins-x-plugins/jx-promote v0.0.256
	github.com/jenkins-x/go-scm v1.6.17
	github.com/jenkins-x/jx-api/v4 v4.0.28
	github.com/jenkins-x/jx-helpers/v3 v3.0.104
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/jenkins-x/lighthouse-client v0.0.102
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/githubv4 v0.0.0-20191102174205-af46314aec7b
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/yargevad/filepathx v0.0.0-20161019152617-907099cb5a62
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c
	k8s.io/api v0.20.6 // indirect
	k8s.io/apimachinery v0.20.6
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

go 1.15
