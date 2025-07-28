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


## Features
There are several features that are provided. They can be toggled depending on the values provided or absence in the configuration file.

### CORS
CORS headers can be added to the response using comma separated lists.

### Rate Limiting
Rate limiting is implemented using the leaky bucket algorithm. The maximum capacity of the bucket can be adjusted using the `max_capacity` field. 

It can be enabled in two modes.

1. **Global bucket** \
A single leaky bucket is used to rate limit all incoming requests. 

2. **One bucket per visitor**\
A leaky bucket is provided to each visitor from a unique domain. The rate limit will apply individually. \
In this mode we clean up visitors that have been inactive following the time period provided in `cleanup_interval_seconds`. 


Choose the specific mode by setting `per_ip_address.enabled`.

### Authentication
A simple jwt authentication has been implemented.\
In order to utilise this you must supply a `JWT_SECRET` secret which is used to create the tokens you want verified.

### IP Whitelisting
Ip addresses in both ipv4 and ipv6 format can be whitelisted to allow for access. They must be provided individually in the `allow` section. \
If the user chooses not to whitelist any addresses, they can leave this section in the configuration file blank.


## Future work
- Retry plugin
- Documentation and set up examples
