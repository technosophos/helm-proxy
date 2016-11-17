# Helm Proxy: A transcoding proxy for JSON-to-gRPC

Tiller, Helm's server-side component, serves gRPC natively. This
provides a simple transcoding proxy that converts JSON to gRPC, and
exposes a number of endpoints. It also adds an authentication layer.

## The JSON API:

Status of the server: `GET /`

List all releases: `GET /v1/releases`

Install a Chart: `POST /v1/releases`

Get a release by name: `GET /v1/releases/{RELEASE_NAME}`

Upgrade a Release: `POST /v1/releases/{RELEASE_NAME}`

Get release history: `GET /v1/releases/{RELEASE_NAME}/history`

Rollback to a previous release: `POST /v1/releases/{RELEASE_NAME}/history/{VERSION}`

Delete a release: `DELETE /v1/releases/{RELEASE_NAME}`
