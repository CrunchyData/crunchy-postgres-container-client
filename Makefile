OSFLAVOR=centos7

gendeps:
	godep save \
	github.com/crunchydata/crunchy-postgres-container-client

build:
	godep go install cmd/backendapi.go

image:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        sudo docker build -t crunchy-pg-client .
		sudo docker tag -f crunchy-pg-client:latest crunchydata/crunchy-pg-client
