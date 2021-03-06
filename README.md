# inmemory-search
In-memory Key Value with prefix / suffix and exact search

- /get/<key> → Return value of the key
- /set → Post call which sets the key/value pair
- /search → Search for keys using the following filters
  - Assume you have keys: abc-1, abc-2, xyz-1, xyz-2
  - /search?prefix=abc return abc-1 and abc-2
  - /search?suffix=-1 return abc-1 and xyz-1

=======================================
## Algorithm / Implementation
- 2 (prefix) trie to store the data
- First prefix trie is storing the key as well as data in the terminal node (word/key end) and then we walk the prefix trie if the "prefix search" happens
- Second prefix trie is storing the "reverse" key but without the "value" / "data" and terminal node is storing nil data. This is done to match the suffix search.
  - When the suffix key pattern comes, we reverse the key pattern
  - walk (reverse) prefix tree
  - reverse thematching keys (so as to get real key)
  - return all possible results

- prefix / suffix search accepts "-1" to return all matching keys, if the numebr is positive, then it limits the search result
- All the cache GET / SET is consurrent safe by adding a sync.mutex lock


## Scope of Improvements
- Could have used adaptive radix tree (compressed trie) to store more efficently
- Could have used suffix array to search on all possible pattern and not use 2 "tries"
- Implement a "delete" to expire the cache
- Implement LRU / autoTTL key expiry
- Introduce singleflight pattern to avoid cache stampeding


## github-action

- On every push to github, it runs linter / test and docker build
- #TODO upload the docker build to GCR (but it requires a paid gcr account)

## Build

- `docker-compose up -d` will run the http path(s) on localhost:8080 and promethus on localhost:9000
- `make build` puts the binary executable in `$root/build` folder.
- it's good to run `make lint test` before the make build command to ensure lint and test passes.
- make `docker-build` generates the docker image in `$root/build` folder.

## Lint
- `make lint` runs the golang-ci lint (installs it if not present) and runs the linter as defined in `root/.golangci.yml`.

## Test
- `make test` runs all the test in main as well as helper packages and generates the coverage report.


## inmemorycache
- this folder contains the inmemroy cache package logic

## cmd
- `main` http service is inside the `cmd/inmemory-search`

## Kubernetes Specs
- k8s folder contains all the specs file as asked

:warning: I am using latest minkube in local so if you use lower version, some yaml config may not work for you :warning:
```
kubectl version
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.2", GitCommit:"092fbfbf53427de67cac1e9fa54aaa09a28371d7", GitTreeState:"clean", BuildDate:"2021-06-16T12:59:11Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"22", GitVersion:"v1.22.1", GitCommit:"632ed300f2c34f6d6d15ca4cef3d3c7073412212", GitTreeState:"clean", BuildDate:"2021-08-19T15:39:34Z", GoVersion:"go1.16.7", Compiler:"gc", Platform:"linux/amd64"}
```

:bulb: Working k8 config in my local :bulb:
```
kubectl apply -f k8s/


namespace/inmemory-production unchanged
deployment.apps/inmemory-http configured
ingress.networking.k8s.io/inmemory-ingress unchanged
service/inmemory-http unchanged
```

## prometheus
- prometheus folder contains promethus yaml and artifacts

![latency sum graph](prometheus/latency_graph.png)


## Postman Code for easy accesibilty

- Set
```
curl -X POST \
  http://localhost:8080/set \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"key":"rohitJ",
	"value":"2929"
}'
````
- Get
```
curl -X GET \
  http://localhost:8080/get/123 \
  -H 'cache-control: no-cache'
```

- Search Suffix
```
curl -X GET \
  'http://localhost:8080/search?suffix=tJ' \
  -H 'cache-control: no-cache'
```
- Search Prefix
```
curl -X GET \
  'http://localhost:8080/search?prefix=ro' \
  -H 'cache-control: no-cache'
```
