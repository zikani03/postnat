postnat
=======

Publish messages to NATS via Postgres' LISTEN/NOTIFY feature.

## Why would I use this?

Let's assume you already have a postgresql backed application and you want to
introduce a Messaging layer but don't have the resources to modify the system 
to introduce more dependencies or you just want to put it off for a bit more.

`postnat` helps you use the facilities of Postgres to publish messages to a NATS
server with minimal changes to your code.

```sql
NOTIFY 'time_us_east', '<binary-or-json-message-here>';
```

What `postnat` does is basically listen to all registered patterns and publishes
the payloads to the configured NATS server.

## Usage

You can run it like so:

```sh
$ postnat --config "postnat.toml" run
```

### CLI 
```text
Usage: postnat <command>

Publish messages to NATS from PostgreSQL LISTEN/NOTIFY messages

Flags:
  -h, --help                     Show context-sensitive help.
      --config="postnat.toml"    Location of configuration file
      --debug                    Enable debug mode
      --version                  Show version and quit

Commands:
  run    Start the postnat daemon

Run "postnat <command> --help" for more information on a command.
```

### Configuration

```toml
[postgres]
host = "localhost"
port = 5432
database = "database"
username = "username"
password = "password"
sslmode = "disable"

[nats]
url = "nats://username:password@localhost:4222"
max_reconnects = 10

[topics]
listen_for = ["users", "users_id"]
# optional topic prefix
prefix = "app."
# This replaces the underscore when publishing to nats, e.g. emails_signup -> emails.signup
replace_underscore_with_dot = true
```

## Building

First clone this repo:

```sh
$ git clone https://github.com/zikani03/postnat
$ cd postnat
```

### Building from source

```sh
$ go build ./cmd/postnat.go
```

### Docker, with docker compose

Create a configuration file named `development.toml` and update as appropriate, then run:

```sh
$ docker compose up
```

---

MIT LICENSE &copy; Zikani Nyirenda Mwase