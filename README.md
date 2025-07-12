## Project Wormhole

Wormhole is a lightweight api gateway that acts as a reverse proxy for one or multiple APIs.

The server and the APIs can be configured in a yaml file. An example is priovided in [sample.yaml](sample.yaml)

Server configurations such as execution port and CORS can be changed.

To create a service entry you must provide the name of the service and the upstream path to proxy to.

```yaml
service-one:
    upstream_path: "some-path"
```
Wormhole will match destination of the request by name and proxy to the provided upstream path.



### Future work
- Documentation and set up examples
- Request authentication
- Retry plugin
