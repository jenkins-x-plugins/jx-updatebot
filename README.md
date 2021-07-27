# jx updatebot

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-updatebot?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-updatebot)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-updatebot)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-updatebot)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x-plugins/jx-updatebot.svg)](https://github.com/jenkins-x-plugins/jx-updatebot/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x-plugins/jx-updatebot.svg)](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-updatebot` is a small command line tool for generating downstream Pull Requests


## Getting Started

Download the [jx-updatebot binary](https://github.com/jenkins-x-plugins/jx-updatebot/releases) for your operating system and add it to your `$PATH`.

Or you can use `jx updatebot` directly in the [Jenkins X 3.x CLI](https://github.com/jenkins-x/jx)


## Configuration

The [jx updatebot pr](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/cmd/jx-updatebot_pr.md) command looks in for the `.jx/updatebot.yaml` file to find the repositories to modify along with the list of change rules to make.

You can see the [configuration documentation here](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/config.md#updatebot.jenkins-x.io/v1alpha1.UpdateConfig) for how to format your `.jx/updatebot.yaml` file.
         
## Examples

Here are some example updatebot configurations:

* [jenkins-x/go-scm](https://github.com/jenkins-x/go-scm) which is a go library uses this [.jx/updatebot.yaml](https://github.com/jenkins-x/go-scm/blob/main/.jx/updatebot.yaml) to create a downstream pull request on [jenkins-xlighthouse](https://github.com/jenkins-x/lighthouse) whenever its released. 
* [jenkins-xlighthouse](https://github.com/jenkins-x/lighthouse) releases a chart which then uses this [.jx/updatebot.yaml](https://github.com/jenkins-x/lighthouse/blob/main/.jx/updatebot.yaml) to create a downstream pull request on the Jenkins X [Version Stream](https://jenkins-x.io/blog/2021/01/26/jx3-walkthroughs/#version-streams)
* [jenkins-x/jx](https://github.com/jenkins-x/jx) is the Jenkins X command line which releases binaries and charts and uses this [.jx/updatebot.yaml](https://github.com/jenkins-x/jx/blob/main/.jx/updatebot.yaml) to create a downstream pull request on the Jenkins X [Version Stream](https://jenkins-x.io/blog/2021/01/26/jx3-walkthroughs/#version-streams)
          
          
## Using updatebot in Tekton

You can see an [example of invoking the Updatebot step](https://github.com/jenkins-x/go-scm/blob/main/.lighthouse/jenkins-x/release.yaml#L22-L23) inside a Tekton pipeline. This reuses the [Jenkins X Pipeline as Code GitOps approach](https://jenkins-x.io/blog/2021/02/25/gitops-pipelines/)

e.g. check out the final step in this Task:

``` 
        steps:
        - image: uses:jenkins-x/jx3-pipeline-catalog/tasks/git-clone/git-clone.yaml@versionStream
        - name: next-version
        - name: jx-variables
        - name: build-make-build
        - name: promote-changelog
        - image: uses:jenkins-x/jx3-pipeline-catalog/tasks/updatebot/release.yaml@versionStream
```

If you are not using Jenkins X but are using Tekton you could [copy/paste this Task step](https://github.com/jenkins-x/jx3-pipeline-catalog/blob/master/tasks/updatebot/release.yaml#L14-L18)

## Commands

See the [jx-updatebot command reference](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/cmd/jx-updatebot.md)
