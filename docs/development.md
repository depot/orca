# Development

## MacOS

Get a containerd socket on your Mac:

```shell
brew install lima
limactl create ./hack/lima.yaml --name orca
limactl start orca
export CONTAINERD_ADDRESS=$(limactl list orca --format '{{.Dir}}/containerd.sock')
```
