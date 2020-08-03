# koffer-go container entrypoint binary
```
root@mordor:~$ podman run -it --rm docker.io/codesparta/koffer bundle -h
Koffer Engine Bundle:
  Bundle is intended to run against koffer collector plugin
  repos to build artifact bundles capable of transporting all
  dependencies required for build or operations time engagement.

  Koffer bundles are designed to be deployed with the Konductor 
  engine and artifacts served via the CloudCtl delivery pod.

Usage:
  koffer bundle [flags]

Flags:
  -b, --branch string      Git Branch (default "master")
  -d, --dir string         Clone Path (default "/root/koffer")
  -h, --help               koffer bundle help
  -r, --repo stringArray   Plugin Repo Name
  -s, --service string     Git Server (default "github.com")
  -u, --user string        Repo {User,Organization}/path (default "CodeSparta")

Global Flags:
      --config string   config file (default is $HOME/.koffer/config.yml)
```
