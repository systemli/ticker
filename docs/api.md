# API Specification

!!! note

    This specification currently covers the **public** endpoints only. The admin API under
    `/v1/admin/**` is not described here yet; the
    [route definitions](https://github.com/systemli/ticker/blob/main/internal/api/api.go) are the
    authoritative reference for it.

All endpoints are served under the `/v1` prefix, with two exceptions that live at the root:

| Endpoint | Purpose |
| --- | --- |
| `GET /media/{file}` | uploaded attachments |
| `GET /healthz` | health check |

## Identifying a ticker

The public endpoints (`/init`, `/timeline`, `/feed`, `/manifest.json`, `/ws`) do not take a ticker
ID. Instead the API resolves the ticker from the request's `Origin` header, matched against the
origins registered on that ticker.

Clients that do not send an `Origin` header — RSS readers, scripts — can pass it explicitly as a
query parameter, which takes precedence:

```shell
curl 'https://api.ticker.example.org/v1/feed?origin=https://ticker.example.org'
```

Requests that match no ticker return HTTP 200 with a `ticker not found` error body, or, for `/init`,
a null ticker together with the instance's inactive-page settings.

## Specification

!!swagger swagger.yaml!!
