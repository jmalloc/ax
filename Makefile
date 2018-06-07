REQ += $(shell find src -name "*.proto")
REQ += src/internal/endpointtest/sinkmock.go
REQ += src/internal/endpointtest/transportmock.go
REQ += src/internal/endpointtest/pipelinemock.go
REQ += src/internal/endpointtest/validatormock.go
REQ += src/internal/routingtest/handlermock.go
REQ += src/internal/persistencetest/datastoremock.go
REQ += src/internal/persistencetest/transactionmock.go
REQ += src/internal/observabilitytest/observermock.go

-include artifacts/make/go/Makefile

.PHONY: banking
banking:
	protoc --go_out=. examples/banking/messages/*.proto
	protoc --go_out=. examples/banking/domain/*.proto
	AX_RMQ_DSN="amqp://localhost" \
	AX_MYSQL_DSN="ax:ax@tcp(127.0.0.1:3306)/ax" \
		go run examples/banking/main.go $(RUN_ARGS)

%.pb.go: %.proto
	protoc --go_out=. $(@D)/*.proto

MOQ := $(GOPATH)/bin/moq
$(MOQ): | vendor # ensure dependencies are installed before trying to build mocks
	go get -u github.com/matryer/moq

src/internal/endpointtest/sinkmock.go: src/ax/endpoint/sink.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "endpointtest" src/ax/endpoint MessageSink

src/internal/endpointtest/transportmock.go: src/ax/endpoint/transport.go src/ax/endpoint/sink.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "endpointtest" src/ax/endpoint Transport

src/internal/endpointtest/pipelinemock.go: src/ax/endpoint/pipeline.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "endpointtest" src/ax/endpoint InboundPipeline OutboundPipeline

src/internal/endpointtest/validatormock.go: src/ax/endpoint/validator.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "endpointtest" src/ax/endpoint Validator SelfValidatingMessage

src/internal/routingtest/handlermock.go: src/ax/routing/handler.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "routingtest" src/ax/routing MessageHandler

src/internal/persistencetest/datastoremock.go: src/ax/persistence/datastore.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence DataStore

src/internal/persistencetest/transactionmock.go: src/ax/persistence/transaction.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence Tx Committer

src/internal/observabilitytest/observermock.go: src/ax/observability/observer.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "observabilitytest" src/ax/observability BeforeInboundObserver AfterInboundObserver BeforeOutboundObserver AfterOutboundObserver


artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
