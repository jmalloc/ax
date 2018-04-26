# REQ += src/ax/internal/messagetest/messages.pb.go
# REQ += src/ax/eventsourcing/messages.pb.go

-include artifacts/make/go/Makefile

.PHONY: banking
banking:
	protoc --go_out=. examples/banking/messages/*.proto
	go run examples/banking/main.go $(RUN_ARGS)

%/messages.pb.go:
	protoc --go_out=. $(@D)/*.proto

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
