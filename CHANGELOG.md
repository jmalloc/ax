# Changelog

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
