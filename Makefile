version     ?= latest
log          = 'error'
devimg       = sysdigclidev
GOPATH      ?= $(HOME)/go
packagename  = github.com/NeowayLabs/sysdig-client
workdir      = /go/src/$(packagename)
rundev       = docker run --net=host -e SEVERINO_LOGGER=$(log) --rm -v `pwd`:$(workdir) --workdir $(workdir) $(devimg)
runbuild     = --rm -v `pwd`:$(workdir) -w $(workdir) $(devimg)
gitversion   = $(version)

all: check-integration analyze

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

release: guard-version publish
	git tag -a $(version) -m "Generated release "$(version)
	git push origin $(version)

imagedev:
	docker build -t $(devimg) -f ./hack/Dockerfile.dev .

vendor: imagedev
	$(rundev) ./hack/vendor.sh
	sudo chown -Rf $(USER):$(USER) ./vendor
	sudo chown -Rf $(USER):$(USER) ./Gopkg.toml
	sudo chown -Rf $(USER):$(USER) ./Gopkg.lock

analyze: imagedev
	docker run --rm -v `pwd`:$(workdir) --workdir $(workdir) $(devimg) ./hack/analyze.sh


check: imagedev
	$(rundev) ./hack/check.sh $(pkg) $(test)

check-integration: imagedev
	$(rundev) ./hack/check-integration.sh $(pkg) $(test)

check-all: analyze coverage

coverage: imagedev
	$(rundev) ./hack/coverage.sh

shell:
	docker run -ti --rm -v `pwd`:$(workdir) --workdir $(workdir) $(devimg)
