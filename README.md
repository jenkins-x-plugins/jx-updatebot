# jx-updatebot

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-updatebot?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-updatebot)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-updatebot)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-updatebot)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x/jx-updatebot.svg)](https://github.com/jenkins-x-plugins/jx-updatebot/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x/jx-updatebot.svg)](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-updatebot` is a small command line tool for generating downstream Pull Requests


## Getting Started

Download the [jx-updatebot binary](https://github.com/jenkins-x-plugins/jx-updatebot/releases) for your operating system and add it to your `$PATH`.

Or you can use `jx updatebot` directly in the [Jenkins X 3.x CLI](https://github.com/jenkins-x/jx-cli)


## Configuration

By default the [jx updatebot pr](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/cmd/jx-updatebot_pr.md) command looks in for the `.jx/updatebot.yaml` file to find the repositories to modify along with the list of change rules to make.

You can see the [configuration documentation here](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/config.md#updatebot.jenkins-x.io/v1alpha1.UpdateConfig) for how to format your `.jx/updatebot.yaml` file.

Here's an example: [.jx/updatebot.yaml](https://github.com/jenkins-x/jx-cli/blob/master/.jx/updatebot.yaml)
 
## Commands

See the [jx-updatebot command reference](https://github.com/jenkins-x-plugins/jx-updatebot/blob/master/docs/cmd/jx-updatebot.md)