# ChatGPT Proxy

Gets around cloudflare via TLS spoofing

## Notes
There is an IP based rate limit. Set a PUID environment variable to get around it
`export PUID="user-..."`
This requires a ChatGPT Plus account

## Building and running
`go build`
`./ChatGPT-Proxy-V4`

## Limitations
This cannot get around an outright IP ban by OpenAI
