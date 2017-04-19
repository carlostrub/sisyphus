SISYPHUS_GO_EXECUTABLE ?= go
DIST_DIRS := find * -type d -exec
VERSION ?= $(shell git describe --tags)
VERSION_INCODE = $(shell perl -ne '/^var version.*"([^"]+)".*$$/ && print "v$$1\n"' sisyphus.go)
VERSION_INCHANGELOG = $(shell perl -ne '/^\# Release (\d+(\.\d+)+) / && print "$$1\n"' CHANGELOG.md | head -n1)

build:
	${SISYPHUS_GO_EXECUTABLE} build -o sisyphus/sisyphus -ldflags "-X main.version=${VERSION}" sisyphus/sisyphus.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./sisyphus/sisyphus ${DESTDIR}/usr/local/bin/sisyphus

test:
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/onsi/ginkgo/ginkgo
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/onsi/gomega
	${GOPATH}/bin/ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --progress

integration-test:
	${SISYPHUS_GO_EXECUTABLE} build
	./sisyphus run
	./sisyphus start
	./sisyphus restart
	./sisyphus status
	./sisyphus stop

clean:
	rm -f ./sisyphus/sisyphus
	rm -rf ./dist

build-all:
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/franciscocpg/gox
	${GOPATH}/bin/gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows freebsd openbsd netbsd" \
	-arch="amd64 386 armv5 armv6 armv7 arm64" \
	-osarch="!darwin/arm64" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" ./sisyphus

dist: build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf sisyphus-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r sisyphus-${VERSION}-{}.zip {} \; && \
	cd ..

verify-version:
	@if [ "$(VERSION_INCODE)" = "v$(VERSION_INCHANGELOG)" ]; then \
		echo "sisyphus: $(VERSION_INCHANGELOG)"; \
	elif [ "$(VERSION_INCODE)" = "v$(VERSION_INCHANGELOG)-dev" ]; then \
		echo "sisyphus (development): $(VERSION_INCHANGELOG)"; \
	else \
		echo "Version number in sisyphus.go does not match CHANGELOG.md"; \
		echo "sisyphus.go: $(VERSION_INCODE)"; \
		echo "CHANGELOG : $(VERSION_INCHANGELOG)"; \
		exit 1; \
	fi

.PHONY: build test install clean build-all dist integration-test verify-version

