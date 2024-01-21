# go-switchbot-metrics

Expose metrics page scrapable by know agent (like Telegraf) for SwitchBot API

## Supported devices

- Meter
- Meter Plus
- Outdoor Meter

## Docker usage

```sh
docker run -p 3000:3000 -d ghcr.io/albinos/go-switchbot-metrics serve --bind_ip=0.0.0.0 --switchbot_openapi_token=<SWITCHBOT_API_TOKEN> --switchbot_secret_key=<SWITCHBOT_SECRET_KEY>
```
