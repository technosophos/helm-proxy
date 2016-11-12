# Helm Proxy: A transcoding proxy for JSON-to-gRPC

Tiller, Helm's server-side component, serves gRPC natively. This
provides a simple transcoding proxy that converts JSON to gRPC, and
exposes a number of endpoints. It also adds an authentication layer.
