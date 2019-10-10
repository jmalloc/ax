# Mock generation is disabled until this PR is merged:
# https://github.com/matryer/moq/pull/105
#
# Until then, we will just leave the existing mocks committed without
# regenerating them.
#
# GENERATED_FILES += axtest/mocks/endpoint.go
# GENERATED_FILES += axtest/mocks/routing.go
# GENERATED_FILES += axtest/mocks/persistence.go
# GENERATED_FILES += axtest/mocks/observability.go

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v1/Makefile
-include .makefiles/pkg/go/v1/Makefile

.PHONY: banking
banking: $(GENERATED_FILES)
	AX_RMQ_DSN="amqp://localhost" \
	AX_MYSQL_DSN="banking:banking@tcp(127.0.0.1:3306)/banking" \
	JAEGER_SERVICE_NAME="ax.examples.banking" \
	JAEGER_SAMPLER_TYPE="const" \
	JAEGER_SAMPLER_PARAM="1" \
	JAEGER_REPORTER_LOG_SPANS=true \
		go run examples/banking/main.go $(RUN_ARGS)

MOQ := $(GOPATH)/bin/moq
$(MOQ):
	go get -u github.com/matryer/moq

axtest/mocks/endpoint.go: $(wildcard endpoint/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" endpoint \
		InboundPipeline \
		MessageSink \
		OutboundPipeline \
		SelfValidatingMessage \
		InboundTransport \
		OutboundTransport \
		Validator

axtest/mocks/routing.go: $(wildcard routing/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" routing \
		MessageHandler

axtest/mocks/persistence.go: $(wildcard persistence/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" persistence \
		Committer \
		DataStore \
		Tx

axtest/mocks/observability.go: $(wildcard observability/*.go) | $(MOQ)
	$(MOQ) -out "$@" -pkg "mocks" observability \
		Observer

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
