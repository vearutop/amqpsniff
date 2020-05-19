# AMQP Sniffer

Collect messages with AMQP protocol.

![Code lines](https://sloc.xyz/github/vearutop/amqpsniff/?category=code)
![Comments](https://sloc.xyz/github/vearutop/amqpsniff/?category=comments)

## Install

```
go get -u github.com/vearutop/amqpsniff
```

## Usage

```
usage: amqpsniff --dsn=DSN --queue=QUEUE --bindings=BINDINGS [<flags>]

AMQP sniffer creates queue and dumps incoming messages.

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --dsn=DSN                  Connection credentials, env var AMQP_DSN, example: amqp://guest:guest@rabbit:5672
  --queue=QUEUE              Queue name, example: my-queue
  --bindings=BINDINGS ...    Bindings, colon-separated exchange and routing key, example: users:*.user.created
  --output="messages.jsonl"  Path to output JSONL file
  --limit=100                Number messages to collect before exiting
```

## Example

```bash
AMQP_DSN=amqp://guest:guest@rabbit:5672 amqpsniff --queue test-test --bindings myexchange:*.users.created --limit 12
```
