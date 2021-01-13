module github.com/jenkins-x-plugins/jx-updatebot

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/cpuguy83/go-md2man v1.0.10
	github.com/jenkins-x/go-scm v1.5.208
	github.com/jenkins-x/jx-gitops v0.0.519
	github.com/jenkins-x/jx-helpers/v3 v3.0.55
	github.com/jenkins-x/jx-logging/v3 v3.0.2
	github.com/jenkins-x/jx-promote v0.0.165
	github.com/jenkins-x/lighthouse v0.0.906 // indirect
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/githubv4 v0.0.0-20191102174205-af46314aec7b
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/yargevad/filepathx v0.0.0-20161019152617-907099cb5a62
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	k8s.io/apimachinery v0.19.4
)

replace (
	github.com/jenkins-x/lighthouse => github.com/rawlingsj/lighthouse v0.0.0-20201005083317-4d21277f7992
	k8s.io/client-go => k8s.io/client-go v0.19.2
)

go 1.15
