SHELL          := /bin/bash
PROGRAM        := fstaid
VERSION        := v0.1.7
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
TEST_SRC       := $(wildcard src/*/*_test.go) $(wildcard src/*/test_*.go)
SRC            := main.go $(filter-out $(TEST_SRC),$(wildcard src/*/*.go))

UBUNTU_IMAGE          := docker-go-pkg-build-ubuntu
UBUNTU_CONTAINER_NAME := docker-go-pkg-build-ubuntu-$(shell date +%s)
CENTOS_IMAGE          := docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME := docker-go-pkg-build-centos6-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get github.com/gin-gonic/gin
	go get github.com/BurntSushi/toml
	go get github.com/stretchr/testify
	go get github.com/bouk/monkey
	go get github.com/fvbock/endless
	go get github.com/mattn/go-shellwords

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X fstaid.version=$(VERSION)" -a -tags netgo -installsuffix netgo -o $(PROGRAM)
	[[ "`ldd $(PROGRAM)`" =~ "not a dynamic executable" ]] || exit 1
else
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X fstaid.version=$(VERSION)" -o $(PROGRAM)
endif

.PHONY: test
test: $(TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test -v $(TEST_SRC)

.PHONY: clean
clean: $(TEST_SRC)
	rm -f $(PROGRAM)

.PHONY: package
package: clean test $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm -f $(PROGRAM)

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(UBUNTU_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb
deb:
	docker run --name $(UBUNTU_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) make -C /tmp/src deb/docker
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb/docker
deb/docker: clean go-get
	dpkg-buildpackage -us -uc
	mv ../fstaid_* pkg/

.PHONY: docker/buil/ubuntu
docker/build/ubuntu: etc/Dockerfile.ubuntu
	docker build -f etc/Dockerfile.ubuntu -t $(UBUNTU_IMAGE) .

.PHONY: rpm
rpm:
	docker run --name $(CENTOS_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) make -C /tmp/src rpm/docker
	docker rm $(CENTOS_CONTAINER_NAME)

.PHONY: rpm/docker
rpm/docker: clean go-get
	cd ../ && tar zcf fstaid.tar.gz src
	mv ../fstaid.tar.gz /root/rpmbuild/SOURCES/
	cp fstaid.spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/fstaid.spec
	mv /root/rpmbuild/RPMS/x86_64/fstaid-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/fstaid-*.src.rpm pkg/

.PHONY: docker/build/centos
docker/build/centos:
	docker build -f etc/Dockerfile.centos -t $(CENTOS_IMAGE) .
