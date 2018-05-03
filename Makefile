REQ += $(shell find src -name "*.proto")
REQ += src/internal/bustest/handlermock.go
REQ += src/internal/bustest/sinkmock.go
REQ += src/internal/bustest/transportmock.go
REQ += src/internal/bustest/pipelinemock.go
REQ += src/internal/persistencetest/datastoremock.go
REQ += src/internal/persistencetest/transactionmock.go

-include artifacts/make/go/Makefile

.PHONY: banking
banking:
	protoc --go_out=. examples/banking/messages/*.proto
	go run examples/banking/main.go $(RUN_ARGS)

%.pb.go: %.proto
	protoc --go_out=. $(@D)/*.proto

MOQ := $(GOPATH)/bin/moq
$(MOQ):
	go get -u github.com/matryer/moq

src/internal/bustest/handlermock.go: src/ax/bus/handler.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus MessageHandler

src/internal/bustest/sinkmock.go: src/ax/bus/sink.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus MessageSink

src/internal/bustest/transportmock.go: src/ax/bus/transport.go src/ax/bus/sink.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus Transport

src/internal/bustest/pipelinemock.go: src/ax/bus/pipeline.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus InboundPipeline OutboundPipeline

src/internal/persistencetest/datastoremock.go: src/ax/persistence/datastore.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence DataStore

src/internal/persistencetest/transactionmock.go: src/ax/persistence/transaction.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence Tx Committer

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
