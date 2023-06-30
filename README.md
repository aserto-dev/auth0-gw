# auth0-gw

Auth0 Gateway

config.yaml
```
---
gateway:
  port: 8383
  path: /events

directory:
  host: directory.prod.aserto.com:8443
  api_key: <read-write directory service API key>
  tenant_id: <tenant-id uuid>
  insecure: false

auth0:
  domain: <auth0 domain name>
  client_id: <auth0 client id>
  client_secret: <auth0 client secret>

loader:
  bin_path: /app
  template: /tmpl/transform.tmpl

scheduler:
  interval: 15m
```

```
docker run -ti \
--platform=linux/amd64 \
--name auth0-gw \
--rm \
-p 8383:8383 \
-v $PWD:/cfg \
-v $PWD:/tmpl \
ghcr.io/aserto-dev/auth0-gw:v0.0.0-amd64 run --config=/cfg/config.yaml --template=/tmpl/transform.tmpl
```
