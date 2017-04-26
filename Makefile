SISYPHUS_GO_EXECUTABLE ?= go
DIST_DIRS := find * -type d -exec
VERSION_GIT != git describe --tags
VERSION_GIT_CLEAN != echo ${VERSION_GIT} | sed -nre 's/^[^0-9]*(([0-9]+\.)*[0-9]+).*/\1/p'
VERSION ?= ${VERSION_GIT_CLEAN}
VERSION_INCHANGELOG != head -n1 CHANGELOG.md | sed -nre 's/^[^0-9]*(([0-9]+\.)*[0-9]+).*/\1/p'

verify-version:
	@if [ "$(VERSION)" = "$(VERSION_INCHANGELOG)" ]; then \
		echo "sisyphus: $(VERSION_INCHANGELOG)"; \
	elif [ "$(VERSION)" = "$(VERSION_INCHANGELOG)-dev" ]; then \
		echo "sisyphus (development): $(VERSION_INCHANGELOG)"; \
	else \
		echo "Version number does not match CHANGELOG.md"; \
		echo "Build Version: $(VERSION)"; \
		echo "CHANGELOG : $(VERSION_INCHANGELOG)"; \
		exit 1; \
	fi

build: verify-version
	${SISYPHUS_GO_EXECUTABLE} build -o sisyphus/sisyphus -ldflags "-X main.version=${VERSION}" sisyphus/sisyphus.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./sisyphus/sisyphus ${DESTDIR}/usr/local/bin/sisyphus

test:
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/onsi/ginkgo/ginkgo
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/onsi/gomega
	${GOPATH}/bin/ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --progress

static-test:
	${SISYPHUS_GO_EXECUTABLE} get -u github.com/alecthomas/gometalinter
	${GOPATH}/bin/gometalinter --install
	${GOPATH}/bin/gometalinter --vendor --deadline=5m --disable-all --enable=gas --enable=goconst --enable=gocyclo --enable=unused --enable=interfacer --enable=lll --enable=misspell --enable=staticcheck --enable=aligncheck --enable=deadcode --enable=goimports --enable=ineffassign --enable=unconvert --enable=unparam --enable=varcheck --enable=dupl --enable=errcheck --enable=golint --enable=structcheck --enable=gosimple --enable=safesql --tests --json . > gometalinter.out

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

.PHONY: build test install clean build-all dist integration-test verify-version

