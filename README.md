# TiLv2

TIL is a simple software that aim to enable article and reflexion-sharing across all Zenika agencies in the world. It
acts as an aggregator (as daily.dev can do), but on-premises and with filters on content to avoid being spammed with
news that don't interest us.

## Structure

* `/server` and `/ui`: the backend and frontend of the application. Each directory have its own README file for more
  information.
* `/extensions`: Firefox & Chrome-based browser extensions.
* `/chart`: Helm chart for Openshift.

## How to use it?

You can either:

- Use it as two separate images (frontend and backend). The Dockerfiles and respective configuration details are in
  `/server` and `/ui` folders.
- Use it as a "All in one" image, with frontend bundled in server (for ease of use). In this case, just run `build.sh`
  and then run `til-aio:latest`. In this case, only `server` configuration file applies (
  check [README.md](./server/README.md)) for more details.

