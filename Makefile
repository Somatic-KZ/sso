# Build Path
BUILD_PATH=./cmd/apiserver

# This how we want to name the binary output
BINARY=./bin/techno-sso

# These are the values we want to pass for VERSION and BUILD
VERSION=`git describe --abbrev=6 --always --tag`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.version=${VERSION}"
GOFLAGS=-a -tags techno-sso -installsuffix techno-sso -mod=vendor

run:
	go run $(BUILD_PATH)/main.go

bin:
	echo "  >  Building binary \"techno-sso\" $(VERSION)..."
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) $(BUILD_PATH)

bin-linux:
	echo "  >  Building linux-amd64 binary \"techno-sso\" $(VERSION)..."
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY)-linux $(BUILD_PATH)

bin-windows:
	echo "  >  Building windows-amd64 binary \"techno-sso\" $(VERSION)..."
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY).exe $(BUILD_PATH)

bin-cross-platform: bin-linux bin-windows

build:
	$(MAKE) bin

build-cross-platform:
	$(MAKE) bin-cross-platform

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY}-linux ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY}.exe ; fi
	if [ -f coverage.html ] ; then rm coverage.html ; fi
	if [ -d .cover ] ; then rm -rf .cover ; fi
	docker-compose down --rmi all -v 2>/dev/null || true
	docker-compose stop >/dev/null
	docker-compose rm >/dev/null

rebuild:
	docker-compose build techno-sso
	docker-compose build unit

unit:
	docker-compose run --rm unit

coverage:
	docker-compose run --rm unit && [ -f ./coverage.html ] && xdg-open coverage.html

public:
	cp -r doc/* bin/techno-sso-linux bin/techno-sso.exe ~/Dropbox/Public/sso/
	docker-compose build && docker tag techno-sso_techno-sso dbzer0/sso && docker push dbzer0/sso
	slackcli -h "auth" -u "🤖Techno БОТ" -i https://vero.kz/bot-icon.png -m "Доступен новый билд techno-sso: `git describe --abbrev=6 --always --tag`"

swagger:
	swag init -g ./internal/ports/http/resources/swagger/v1/resource.go -o ./api

release: build public

.PHONY: build clean build-cross-platform unit rebuild coverage release swagger
