# datadog-github-pr
Send a count of Pull-Request as custom metrics to Datadog

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBBW0dpdEh1Yl1cbiAgICBCW2RhdGFkb2ctZ2l0aHViLXByXVxuICAgIENbRGF0YWRvZ11cbiAgICBCIC0tPnxnZXQgUHVsbCBSZXF1ZXN0c3wgQVxuICAgIEIgLS0-fHNuZWQgYSBjb3VudCBvZiBQdWxsIFJlcXVlc3RzIGFzIGN1c3RvbSBtZXRyaWNzfCBDIiwibWVybWFpZCI6e30sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBBW0dpdEh1Yl1cbiAgICBCW2RhdGFkb2ctZ2l0aHViLXByXVxuICAgIENbRGF0YWRvZ11cbiAgICBCIC0tPnxnZXQgUHVsbCBSZXF1ZXN0c3wgQVxuICAgIEIgLS0-fHNuZWQgYSBjb3VudCBvZiBQdWxsIFJlcXVlc3RzIGFzIGN1c3RvbSBtZXRyaWNzfCBDIiwibWVybWFpZCI6e30sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

## Preparation

Copy .envrc and lod it.

```
$ cp .envrc.sample .envrc
$ # edit .envrc
$ # source .envrc
```

The target repositories are specified by GITHUB_REPOSITORIES environment varibales, that should be written in org/reponame, separated by commas.

>export GITHUB_REPOSITORIES="chaspy/datadog-github-pr,chaspy/favsearch"

## How to run

### Local

```
$ go run main.go
```

### Binary

Get the binary file from [Releases](https://github.com/chaspy/datadog-github-pr/releases) and run it.

### Docker

```
$ docker run -e DATADOG_API_KEY="${DATADOG_API_KEY}" -e DATADOG_APP_KEY="${DATADOG_APP_KEY}" -e GITHUB_TOKEN="${GITHUB_TOKEN}" -e GITHUB_REPOSITORIES="${GITHUB_REPOSITORIES}" chaspy/datadog-github-pr:v0.1.1
```
