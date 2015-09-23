#!/bin/bash

# start up the REST API
/var/cpm/bin/backendapi &

# start up the web server
/usr/sbin/nginx -c /var/cpm/conf/nginx.conf > /tmp/nginx.log 2> /tmp/nginx.err
