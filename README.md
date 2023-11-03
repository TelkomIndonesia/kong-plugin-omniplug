# Kong Plugin Omniplug

Omniplug is a [Kong](https://docs.konghq.com/gateway/latest/) plugin to inject proxied HTTP request with data obtained from various source.

Currently supported data source includes:

- [static json/yaml](./testdata/config.yaml#L2-L6) strings.
- [Kong context](./testdata/config.yaml#L7-L20) data.
- [HTTP response](./testdata/config.yaml#L21-L34) data.

Currently supported target of injection includes:

- Kong [context](./testdata/config.yaml#L37-L43).
- Upstream HTTP request [Headers](./testdata/config.yaml#L44-L49).
- Upstream HTTP request [Query Parameters](./testdata/config.yaml#L50-L55).

See [example](./examples/docker-compose.yml) on how to run with docker-compose.
