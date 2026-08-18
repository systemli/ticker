# API Specification

!!! note

    This specification currently covers the **public** endpoints only. The admin API under
    `/v1/admin/**` is not described here yet; the
    [route definitions](https://github.com/systemli/ticker/blob/main/internal/api/api.go) are the
    authoritative reference for it.

All endpoints are served under the `/v1` prefix, including `GET /v1/media/{file}` for uploaded
attachments. The only exception is `GET /healthz`, which lives at the root.

Attachment URLs in responses are **relative to the site that served them**, of the form
`/api/media/{file}` — that is the path the admin and frontend expose the API under. When you talk to
the API directly, replace `/api` with `/v1`.

## Identifying a ticker

The public endpoints (`/init`, `/timeline`, `/feed`, `/manifest.json`, `/ws`) do not take a ticker
ID. Instead the API resolves the ticker from the request's `Origin` header, matched against the
origins registered on that ticker.

Clients that do not send an `Origin` header — RSS readers, scripts — can pass it explicitly as a
query parameter, which takes precedence:

```shell
curl 'https://ticker.example.org/api/feed?origin=https://ticker.example.org'
```

Requests that match no ticker return HTTP 200 with a `ticker not found` error body, or, for `/init`,
a null ticker together with the instance's inactive-page settings.

## Specification

!!swagger swagger.yaml!!
