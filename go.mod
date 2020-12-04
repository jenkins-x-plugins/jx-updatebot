module github.com/jenkins-x-plugins/jx-updatebot

require (
	github.com/cpuguy83/go-md2man v1.0.10
	github.com/jenkins-x/go-scm v1.5.193
	github.com/jenkins-x/jx-api/v4 v4.0.12 // indirect
	github.com/jenkins-x/jx-gitops v0.0.445
	github.com/jenkins-x/jx-helpers/v3 v3.0.31
	github.com/jenkins-x/jx-logging/v3 v3.0.2
	github.com/jenkins-x/jx-promote v0.0.149
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	k8s.io/apimachinery v0.19.3
)

replace (
	github.com/jenkins-x/lighthouse => github.com/rawlingsj/lighthouse v0.0.0-20201005083317-4d21277f7992
	k8s.io/client-go => k8s.io/client-go v0.19.2
)

go 1.15
