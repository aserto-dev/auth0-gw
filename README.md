# auth0-gw
Auth0 Gateway


```
docker run -ti \
--platform=linux/amd64 \
--name auth0-gw \
--rm \
-p 8080:8383 \
-v $PWD:/cfg \
-v $PWD/.dev:/tmpl \
ghcr.io/aserto-dev/auth0-gw:v0.0.0-amd64 run --config=/cfg/metrikus.yaml --template=/tmpl/metrikus.tmpl