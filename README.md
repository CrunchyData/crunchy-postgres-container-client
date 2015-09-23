# crunchy-postgresql-container-client
This is a test client, web app written in Angular.js that
has a golang REST API.  The REST api makes calls to a
postgresql database running in another pod.  You pass
environment variables to the client to tell it which 
database to use and which user-password to use in its
database connections.  You can also send curl commands
to the REST API to test it or write your own UI client!

##REST API

~~~~~~~~~~~~~~~~~~~
cd test
curl http://crunchy-pg-client:13002/car/list
curl http://crunchy-pg-client:13002/car/1
curl -X POST -d @add-car.json http://crunchy-pg-client:13002/car/add
curl -X POST -d @update-car.json http://crunchy-pg-client:13002/car/update
curl -X POST -d @delete-car.json http://crunchy-pg-client:13002/car/delete
~~~~~~~~~~~~~~~~~~~

##Building the example:

~~~~~~~~~~~~~~~~~~~
mkdir clientproject
mkdir -p clientproject/pkg clientproject/bin clientproject/src
export GOPATH=$HOME/clientproject
export GOBIN=$GOPATH/bin
cd clientproject
go get github.com/tools/godep
go get github.com/crunchydata/crunchy-postgres-container-client
cd src/crunchydata/crunchy-postgres-container-client
godep restore
make build
make image
~~~~~~~~~~~~~~~~~~~

##Setting up the client database
You will need to build the test client database tables using
the following command.  Lookup  the pg-standalone database
credentials using:
~~~~~~~~~~~~~~~~~~~
oc describe pod pg-standalone
~~~~~~~~~~~~~~~~~~~

~~~~~~~~~~~~~~~~~~~
psql -h pg-standalone -U testuser -f setup.sql userdb
~~~~~~~~~~~~~~~~~~~

cd openshift
##Running the example in openshift:

~~~~~~~~~~~~~~~~~~~
cd openshift
vi crunchy-pg-client.json (change env vars to pg-standalone deployed values)
oc create -f crunchy-pg-client.json
oc get pods
~~~~~~~~~~~~~~~~~~~

##Testing the example:

~~~~~~~~~~~~~~~~~~~
curl http://crunchy-pg-client:13002/car/list
~~~~~~~~~~~~~~~~~~~

or if you have external networking set up to your openshift
VM:

browse to http://crunchy-pg-client:13001


