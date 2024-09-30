# Dogebox reflector service

This service is a basic in-memory key-pair cache. It is used by the Dogebox (by default) to persist its local, internal IP address somewhere that the users client can find it once the Dogebox has switched networks.

### API

#### GET /:token

Fetch the IP submitted via token. Is a one-shot fetch, will be removed after you have retrieved it.

Return `200` `{ "ip": "1.2.3.4" }` if found.

Returns `404` if not found.

#### POST /

Required body: `{ "token": "abc", "ip": "1.2.3.4" }`

Returns `201` on create.

Rate limited to 1 per minute via IP.
