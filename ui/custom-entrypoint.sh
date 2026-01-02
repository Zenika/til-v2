#!/bin/sh

export RESOLVER="$(cat /etc/resolv.conf|grep  '^nameserver' |cut -d ' ' -f2 | tr '\n' ' ')"

envsubst '$PUBLIC_TIL_SERVER_URL,$RESOLVER' < "/etc/nginx/conf.d/default.conf" | sponge "/etc/nginx/conf.d/default.conf"

/docker-entrypoint.sh

nginx -g "daemon off;"