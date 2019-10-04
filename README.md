# reporting-client

## About `reporting-client`

`reporting` is an internal service at
Solo.io that we use to check version information, broadcast security
bulletins, etc.

We understand that software making remote calls over the internet
for any reason can be undesirable. Because of this, `reporting-client` can be
disabled in all of our software that includes it. You can view the source
of this client to see that we're not sending any private information.

If you would like to disable this functionality, you can set the following
environment variables:
* `USAGE_REPORTING_DISABLE` -> `"true"`: This disables the reporting of usage
statistics like total number of requests handled by Envoy, number of active
Envoy instances, etc. For `Gloo`, we start up this client in the `gloo` pod,
so set this var on that pod.

**Note:** This repository is probably useless outside of internal Solo.io
use. It is open source for disclosure and because our open source projects
must be able to link to it.

## Local testing

`main.go` starts up a client that reports every two seconds to an instance of
`reporting-server` (the closed-source server that receives these reports) running
on `localhost:3000`. Do `go run main.go`
