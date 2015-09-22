OSFLAVOR=centos7

gendeps:
	godep save \
	github.com/crunchydata/crunchy-pg-client/backend

build:
	godep go install cmd/backend

image:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        sudo docker build -t crunchy-pg-client .
		sudo docker tag -f crunchy-pg-client:latest crunchydata/crunchy-pg-client
