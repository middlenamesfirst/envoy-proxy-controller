# envoy-proxy-controller

A Kuberentes controller for the Envoy proxy (https://www.envoyproxy.io/) data plane.

## Adding Dependencies

Confluence documentation: [Vendoring Using go mod and Maintaining Vendored Code](https://confluence.internal.digitalocean.com/display/DEVTOOLS/Vendoring+Using+go+mod+and+Maintaining+Vendored+Code)

To add a new dependency, include `GO111MODULE=on` and the `-m` flag and run your `go get` command normally.

For example if you want to install the dependency github.com/envoyproxy/go-control-plane, use the following command.

```sh
GO111MODULE=on go get -m github.com/envoyproxy/go-control-plane
```

Once you've downloaded and used a dependency. Run the `make version` command to vendor your code.

```sh
make version
```
