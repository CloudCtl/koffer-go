# koffer-go container entrypoint binary
```
root@mordor:~$ podman run -it --rm docker.io/codesparta/koffer bundle -h
Koffer Engine Bundle:
  Bundle is intended to run against koffer collector plugin
  plugins to build artifact bundles capable of transporting all
  dependencies required for build or operations time engagement.

  Each Koffer plugin should be a reference to a git repository 
  specified by the service, user or organization, and plugin name. 
  The plugin name supports a syntax that allows for overriding the
  default git reference ('master'). The syntax is '--plugin repository@ref'.

  Koffer bundles are designed to be deployed with the Konductor 
  engine and artifacts served via the CloudCtl delivery pod.

  Example - Infra from master:
        koffer bundle --plugin collector-infra

  Example - Infra from tag v1.0
    koffer bundle --plugin collector-infra@v1.0

  Example - All plugins from tag v1.0 by default
    koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps

  Example - Default to tag v1.0 but use master branch on Apps
    koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps@master

Usage:
  koffer bundle [flags]

Flags:
  -d, --dir string           Clone Path (default "/root/koffer")
  -v, --version string       Default git reference to use for all plugin repositories. (default "master")
  -h, --help                 koffer bundle help
  -p, --plugin stringArray   Name of plugin repository to use with optional @gitref.
  -s, --service string       Git Server (default "github.com")
  -a, --silent               Ask for pull secret, if true uses existing value in /root/.docker/config.json
  -u, --user string          Repo {User,Organization}/path (default "CodeSparta")

Global Flags:
      --config string   config file (default is $HOME/.koffer/config.yml)
```
