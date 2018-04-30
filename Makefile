REQ += $(shell find src -name "*.proto")
REQ += src/ax/internal/bustest/sender.go
REQ += src/ax/internal/bustest/transport.go
REQ += src/ax/internal/bustest/pipeline.go

-include artifacts/make/go/Makefile

%.pb.go: %.proto
	protoc --go_out=. $(@D)/*.proto

MOQ := $(GOPATH)/bin/moq
$(MOQ):
	go get -u github.com/matryer/moq

src/ax/internal/bustest/sender.go: src/ax/bus/sender.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus MessageSender

src/ax/internal/bustest/transport.go: src/ax/bus/transport.go src/ax/bus/sender.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus Transport

src/ax/internal/bustest/pipeline.go: src/ax/bus/pipeline.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus InboundPipeline OutboundPipeline

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
