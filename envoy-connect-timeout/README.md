This is a simple sandbox with envoy as front proxy and [httpbin](https://eu.httpbin.org/) as upstream service
that can be used to illustrate scenario or bug.

We are testing here having a TCP connect timeout large (10s) on the cluster service, how envoy will react if the service (httpbin) port is not bonded (by using a different port).

A RST is sent immediately and envoy upstream is reset immediately.
