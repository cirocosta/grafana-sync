### grafana-sync - Keeps your Grafana dashboards in sync

At each time that it's run, `grafana-sync` gathers information
about dashboards from a particular source of truth (a grafana
deployment) and then updates the state the filesystem to reflect
such source.

ps.: assumes use of the `folders` feature from grafana 5+

### Usage

1. Create an API key that is capable of visualizing all dashboards

![API Key creation](./assets/create-api-key.png)

2. Run `grafana-sync`

```sh
./grafana-sync \
    --verbose \
    --address http://my-instance.com \
    --access-token=<api_key> \
    --directory=./dashboards
```

