#!/bin/bash 

# Copyright 2015 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
	   echo "This script must be run as root" 1>&2
	      exit 1
fi



WWW=/home/jeffmc/clientproject/src/github.com/crunchydata/crunchy-postgres-container-client/www

echo "restarting crunchy-pg-client container..."
docker stop crunchy-pg-client
docker rm crunchy-pg-client
chcon -Rt svirt_sandbox_file_t /home/jeffmc/clientproject/src/github.com/crunchydata/crunchy-postgres-container-client/www
docker run --name=crunchy-pg-client -d \
	-e PG_HOST=192.168.0.106 \
	-e PG_USER=user1 \
	-e PG_DATABASE=testdb \
	-e PG_PASSWORD=password \
	-v /home/jeffmc/clientproject/src/github.com/crunchydata/crunchy-postgres-container-client/www:/www crunchydata/crunchy-pg-client:latest

