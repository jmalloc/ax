REQ += $(shell find src -name "*.proto")
REQ += src/axtest/mocks/endpoint.go
REQ += src/axtest/mocks/routing.go
REQ += src/axtest/mocks/persistence.go
REQ += src/axtest/mocks/observability.go

-include artifacts/make/go/Makefile

.PHONY: banking
banking:
	protoc --go_out=paths=source_relative:. examples/banking/messages/*.proto
	protoc --go_out=paths=source_relative:. examples/banking/domain/*.proto
	protoc --go_out=paths=source_relative:. examples/banking/workflows/*.proto
	AX_RMQ_DSN="amqp://localhost" \
	AX_MYSQL_DSN="banking:banking@tcp(127.0.0.1:3306)/banking" \
	JAEGER_SERVICE_NAME="ax.examples.banking" \
	JAEGER_SAMPLER_TYPE="const" \
	JAEGER_SAMPLER_PARAM="1" \
	JAEGER_REPORTER_LOG_SPANS=true \
		go run examples/banking/main.go $(RUN_ARGS)

%.pb.go: %.proto
	protoc --go_out=paths=source_relative:. $(@D)/*.proto

MOQ := $(GOPATH)/bin/moq
$(MOQ): | vendor # ensure dependencies are installed before trying to build mocks
	go get -u github.com/matryer/moq

src/axtest/mocks/endpoint.go: $(wildcard src/ax/endpoint/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" src/ax/endpoint \
		InboundPipeline \
		MessageSink \
		OutboundPipeline \
		SelfValidatingMessage \
		InboundTransport \
		OutboundTransport \
		Validator

src/axtest/mocks/routing.go: $(wildcard src/ax/routing/*.go)  | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" src/ax/routing \
		MessageHandler

src/axtest/mocks/persistence.go: $(wildcard src/ax/persistence/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" src/ax/persistence \
		Committer \
		DataStore \
		Tx

src/axtest/mocks/observability.go: $(wildcard src/ax/observability/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" src/ax/observability \
		Observer

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
