default: help

HOST_GOLANG_VERSION     = $(go version | cut -d ' ' -f3 | cut -c 3-)
# this variable is used like a function. First arg is the minimum version, Second arg is the version to be checked.
ALLOWED_GO_VERSION      = $(test '$(/bin/echo -e "$(1)\n$(2)" | sort -V | head -n1)' == '$(1)' && echo 'true')

## Make bin for mesos-metrics.
bin:
	./control build

## Get vet go tools.
vet:
	go get -u golang.org/x/tools/cmd/vet

## `go get github.com/golang/lint/golint`
.golint:
	go get github.com/golang/lint/golint
ifeq ($(call ALLOWED_GO_VERSION,1.5,$(HOST_GOLANG_VERSION)),true)
	golint ./...
endif

## Validate this go project.
validate:
	script/validate-gofmt
#	go vet ./...

## Run test case for this go project.
test:
	go list ./... | grep -v 'vendor' | xargs -L1 go test -v

## Clean everything (including stray volumes).
clean:
	-rm -rf var
	-rm -f mesos-metrics

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb)                                                      \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L)  \
	$(MAKEFILE_LIST)

