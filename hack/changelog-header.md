### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-updatebot/releases/download/v{{.Version}}/jx-updatebot-linux-amd64.tar.gz | tar xzv 
sudo mv jx-updatebot /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-updatebot/releases/download/v{{.Version}}/jx-updatebot-darwin-amd64.tar.gz | tar xzv
sudo mv jx-updatebot /usr/local/bin
```

