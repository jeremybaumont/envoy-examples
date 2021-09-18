This is a simple sandbox with envoy as front proxy and [httpbin](https://eu.httpbin.org/) as upstream service
that can be used to illustrate scenario or bug.

You can run script against different envoy version like this:
```
for version in v1.14.1 v1.14.2 v1.14.3 v1.14.4 v1.15.0; do docker-compose build --build-arg ENVOY_VERSION=$version; docker-compose up -d --quiet-pull --no-build; curl -sSL -H 'Cache-Control: no-cache'  https://gist.githubusercontent.com/jeremybaumont/f5d7ddc63f6a3a75431a4a6a5016efbe/raw | python - | tee -a /tmp/results; curl -s http://localhost:8001/server_info | grep '"version"' | tee -a /tmp/results ; docker-compose down; done
```
