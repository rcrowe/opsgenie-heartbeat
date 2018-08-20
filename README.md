# opsgenie-heartbeat

A super simple client to send ping requests to the [Opsgenie]() [Heartbeat API]().

While the [SDK](https://github.com/opsgenie/opsgenie-go-sdk/tree/master/heartbeat) supports heartbeat requests,
it's more involved to use & I wanted the underlying HTTP GET requests to be more rebust to failure.

## Usage

The opsgenie API key comes from the env variable `OPSGENIE_HEARTBEAT_KEY`.

```golang
hb := heartbeat.New("name-of-the-heartbeat")
hb.Ping()
```