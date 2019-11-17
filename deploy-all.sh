#!/bin/bash

docker build --tag peterzandbergen/myipaddress:scratch . && \
docker push peterzandbergen/myipaddress:scratch && \
cf push # --docker-image peterzandbergen/myipaddress:scratch
