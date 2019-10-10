# Changelog

## Next Release

- **[IMPROVED]** Switch to using Go modules instead of Glide.
- **[BC]** Move all source files out of the `src` folder. The `src/ax` package
  has been moved to the root of the repository.

## 0.3.0 (2018-11-20)

- **[BC]** The method signatures and naming conventions supported by `saga.NewAggregate()` have changed, see the documentation for details
- **[BC]** The method signatures and naming conventions supported by `saga.NewWorkflow()` have changed, see the documentation for details
- **[BC]** `saga.Saga.HandleMessage()` and `HandleNotFound()` now accept an `ax.MessageContext` instead of an `Envelope`
- **[BC]** The method signatures supported by `routing.NewMessageHandler()` have changed, see the documentation for details
- **[BC]** `routing.MessageHandler.HandeMessage()` now accepts an `ax.MessageContext` instead of an `Envelope`
- **[BC]** `projection.Projector.ApplyMessage()` now accepts an `ax.MessageContext` instead of an `Envelope`
- **[BC]** The method signatures supported by `axmysql.projection.NewReadModel()` have changed, see the documentation for details
- **[BC]** Renamed `persistence.Injector` to `InboundInjector` and added `OutboundInjector`
- **[BC]** Renamed `endpoint.InboundEnvelope.DeliveryCount` to `AttemptCount`
- **[BC]** `endpoint.WithContext()` and `GetEnvelope()` now operate on an `InboundEnvelope` instead of `ax.Envelope`
- **[BC]** Split `endpoint.Transport` into separate `InboundTransport` and `OutboundTransport` interfaces
- **[BC]** Split `endpoint.Endpoint.Transport` into separate `InboundTransport` and `OutboundTransport` fields
- **[BC]** Renamed `endpoint.Endpoint.In` and `Out` to `InboundPipeline` and `OutboundPipeline`, respectively
- **[BC]** Renamed `endpoint.Sender.Out` to `OutboundPipeline`
- **[BC]** Merged all observer interfaces in `observability` into a single interface
- **[FIX]** Fix issue with `axrmq` whereby a message's attempt count was reported as `1` on both the first and second attempts
- **[NEW]** Added support for OpenTracing
- **[NEW]** Added `endpoint.InboundEnvelope.AttemptID`
- **[NEW]** Added `ax.Envelope.Delay()`, which returns the amount of time the message is delayed by
- **[NEW]** Added logging support to `projection.Consumer`
- **[NEW]** Added logging support to application-defined message handlers via new `ax.MessageContext` type
- **[IMPROVED]** `ax.Delay()` now computes the `SendAt` time based on `CreatedAt`, instead of `time.Now()`
- **[IMPROVED]** Improved message logging output
- **[IMPROVED]** Added guarantee that `saga.UnitOfWork.Save()` will set the new instance revision
- **[IMPROVED]** Prevent extraneous query when detecting the end of a stream in the `axmysql` implementation of `messagestore.Stream`

## 0.2.0 (2018-10-09)

- **[BC]** Workflow handler methods now accept an `ax.CommandExecutor` instead of returning commands
- **[BC]** Add options to `ax.Sender.ExecuteCommand()` and `PublishEvent()`
- **[BC]** Remove `ax.Envelope.Time` in favour of `CreatedAt` and `SendAt` times
- **[BC]** The `axrmq` transport can not communicate with prior versions
- **[FIX]** `saga.MessageHandler` now only catches panics directly related to recording events
- **[NEW]** Add support for sending delayed commands, used to implement business process timeouts, etc
- **[NEW]** Add protocol-buffers based representation of message envelopes
- **[NEW]** Add `ax.Envelope.Equal()` method for checking message envelope equality
- **[NEW]** Add `ax.GenerateMessageID()`, `ParseMessageID()` and `MustParseMessageID()`
- **[NEW]** Add `saga.GenerateInstanceID()`, `ParseInstanceID()` and `MustParseInstanceID()`
- **[NEW]** Add `routing.NewMessageHandler()`
- **[NEW]** Aggregate and workflow handler methods may now optionally accept an `ax.Envelope`
- **[NEW]** Add support for the `bytes` protocol-buffers type to the `axcli` package

## 0.1.0 (2018-06-18)

- Initial release
