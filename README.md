# xplex-agent

xplex-agent is deployed on edge servers of xplex to control media streaming. It's responsibilities include:

- Handle local nginx process authentication
- Perform stream authentication by looking up data with rig-server
- Spin up and control media worker processes (`ffmpeg`)
- Report usage status

## Install

- [Install Golang 1.9+](https://golang.org/doc/install)
- Install [Glide](https://github.com/Masterminds/glide)

In project root, run:

```sh
$ go build
```
