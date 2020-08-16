Tweam
===

Tweam provides a Webhook by using `Account Activity API` and `Timeline API`.

Installation
---

### Direct Installation

```
$ go get github.com/flum1025/tweam/...
$ scheduler -config setting.yml
$ worker -config setting.yml
```

### Use Docker

```
$ docker run -v $(pwd)/setting.yml:/setting.yml flum1025/tweam scheduler -config /setting.yml
$ docker run -v $(pwd)/setting.yml:/setting.yml flum1025/tweam worker -config /setting.yml
```

### Configuration File

Configuration file is YAML format file.

```yaml
accounts:
  - id: 123456 # your twitter id
    token:
      consumer_key: consumer_key
      consumer_secret: consumer_secret
      access_token: access_token
      access_token_secret: access_token_secret
    home_timeline_fetch_interval: 60
    mention_timeline_fetch_interval: 60
    webhooks:
      - http://example.com
redis:
  address: localhost:6379
  db: 0
```
