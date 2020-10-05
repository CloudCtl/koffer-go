# Koffer-Go

## Description
Koffer-go is the Golang entrypoint to the Koffer image. Koffer builds offline installation bundles for OpenShift. 

## Command Line Usage
```
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle -h
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
  -S, --silent               Ask for pull secret, if true uses existing value in /root/.docker/config.json
  -u, --user string          Repo {User,Organization}/path (default "containercraft")

Global Flags:
      --config string   config file (default is $HOME/.koffer/config.yml)
```

## Choosing a Version
Koffer-go accepts a command line "--plugin" parameter with an option to define a "version". In the case of koffer-go the
"version" specified is (in order): the tag, the branch, and then the *full* git hash. In the event that a tag and a branch
exist with the same name the tag will be chosen.

Examples:
```bash
# creating a koffer bundle for a tag (refs/tag/v1.0)
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle --plugin collector-infra@v1.0
# creating a koffer bundle for a branch (refs/branches/testing)
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle --plugin collector-infra@testing
# creating a koffer bundle for a full git hash (for a specific commit)
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle --plugin collector-infra@3443e502878b8e3ffc9b405b6648428a208a21b6
```

It is possible to set a version to be used for all plugins:
```bash
# using the v1.0 tag for both plugins
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps
```

Or to mix and match:
```bash
# using the v1.0 tag for infra/operators plugins but the testing branch for collector-apps
root@mordor:~$ podman run -it --rm docker.io/containercraft/koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps@testing --plugin collector-operators
```

## Configuration File
Koffer uses the standard Sparta configuration format from [sparta-libs](https://github.com/containercraft/sparta-libs). The 
`koffer` section of the configuration file configures koffer's behavior.

A sample section might look like:
```
koffer:
  # silent is the same as the command line flag --silent, which will not ask for the quay pull secret
  # and expects it to already be in the docker configuration directory
  silent: true
  # a map of plugin names to the configuration for that plugin
  plugins:
    # specifies github.com/containercraft/collector-infra at version (tag then branch) 1.0.0
    collector-infra:
      version: 1.0.0
      service: github.com
      organization: containercraft
    # specifies github.com/containercraft/collector-operators at version (tag then branch) 1.0.0
    collector-operators:
      version: 1.0.0
      service: github.com
      organization: containercraft
    # specifies github.com/containercraft/collector-apps at branch testing, using "branch" directly does not check tags
    collector-apps:
      service: github.com
      organization: containercraft
      branch: testing
```

You can override the options in the configuration file by setting command line flags. The following examples use
the configuration example (in $HOME/.koffer/config.yml, the default location) above.

```bash
# run the plain configuration from ~/.koffer/config.yml
root@mordor:~$ podman run -it --rm -v ~/config.yml:~/.koffer/config.yml docker.io/containercraft/koffer bundle
# override the version to 1.1.0 for _all_ plugins
root@mordor:~$ podman run -it --rm -v ~/config.yml:~/.koffer/config.yml docker.io/containercraft/koffer bundle --version 1.1.0
# override just collector-apps to use "unstable" branch
root@mordor:~$ podman run -it --rm -v ~/config.yml:~/.koffer/config.yml docker.io/containercraft/koffer bundle --plugin collector-apps@unstable
```



