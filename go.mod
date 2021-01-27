module github.com/jenkins-x-plugins/jx-updatebot

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/cpuguy83/go-md2man v1.0.10
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/jenkins-x/go-scm v1.5.216
	github.com/jenkins-x/jx-gitops v0.0.531
	github.com/jenkins-x/jx-helpers/v3 v3.0.72
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/jenkins-x/jx-promote v0.0.179
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/githubv4 v0.0.0-20191102174205-af46314aec7b
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/yargevad/filepathx v0.0.0-20161019152617-907099cb5a62
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	k8s.io/apimachinery v0.20.2
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

go 1.15
