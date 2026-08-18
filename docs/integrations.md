# Integrations

Besides its own public page, a ticker can push every message to Telegram, Mastodon, Bluesky and
Signal groups.

!!! important "Integrations are configured at runtime, not in a config file"

    All credentials live in the database and are managed through the admin interface. There are no
    environment variables and no `config.yml` keys for them. Older documentation described
    `telegram:` and `signal_group:` config blocks and variables such as `TICKER_TELEGRAM_TOKEN` —
    these no longer exist and are silently ignored if present.

There are two levels:

| Integration | Instance-wide setup | Per-ticker setup |
| --- | --- | --- |
| Telegram | required (bot token) | channel |
| Signal | required (signal-cli endpoint) | group |
| Mastodon | none | server and credentials |
| Bluesky | none | handle and app password |

Telegram and Signal need an instance-wide step by a super admin before editors can use them. Until
that is done, the admin interface hides them.

## Telegram

**Instance-wide.** Create a bot with [@BotFather](https://t.me/BotFather) and note the token. As a
super admin, enter it in the admin interface under the Telegram settings.

Stored as the `telegram_settings` record, with fields `token` and `botUsername`.

**Per ticker.** Create a Telegram channel, add the bot as an administrator with permission to post,
then set the channel name on the ticker.

## Signal

**Instance-wide.** Signal support talks to a [signal-cli](https://github.com/AsamK/signal-cli) REST
endpoint, which is **not part of the Ticker stack** — you run it yourself. A common choice is
[signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api).

Add it to the stack on the internal network:

```yaml
  signal-cli:
    image: bbernhard/signal-cli-rest-api:latest
    environment:
      MODE: json-rpc
    volumes:
      - signal-cli-data:/home/.local/share/signal-cli
    networks:
      - internal
```

Register the phone number with signal-cli first, following its own documentation. Then, as a super
admin, configure the Signal settings:

| Field | Meaning |
| --- | --- |
| `apiUrl` | URL of the signal-cli JSON-RPC endpoint, for example `http://signal-cli:8080/api/v1/rpc` |
| `account` | the registered phone number in international format, e.g. `+491234567890` |
| `avatar` | optional path to an avatar image used for created groups |

Stored as the `signal_group_settings` record. Signal is considered enabled once **both** `apiUrl`
and `account` are set.

**Per ticker.** Enabling Signal on a ticker creates a Signal group and returns an invite link that
you can publish. Ticker admins can be promoted to group admins.

## Mastodon

Configured entirely per ticker: the server URL and the application credentials. Register an
application on your Mastodon instance and enter the resulting values on the ticker. It is
considered connected once the token, secret and access token are all present.

## Bluesky

Configured entirely per ticker: the handle and an **app password** — generate one in the Bluesky
settings rather than using your account password.

Optionally restrict who may reply to the posts. Values are `followers`, `following`, `mentioned`
and `nobody`; leave empty to allow anyone. Several can be combined with commas, for example
`followers,mentioned`.

## Behaviour

Dispatch happens when a message is created, and integrations are also told about deletions and
ticker updates. A failing integration is logged and does **not** block the message from appearing
on the ticker itself, so a broken Telegram token will not take your public page down. Check the API
logs if a message did not arrive somewhere:

```shell
docker compose logs ticker | grep bridge_name
```

Attachments are sent along, using the absolute URLs built from `TICKER_UPLOAD_URL`. If that value is
wrong, posts arrive with broken images — see [Configuration](configuration.md#uploads).
