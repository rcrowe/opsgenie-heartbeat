# opsgenie-heartbeat

A super simple client to send ping requests to the [Opsgenie]() [Heartbeat API]().

While the [SDK](https://github.com/opsgenie/opsgenie-go-sdk/tree/master/heartbeat) supports heartbeat requests,
it's more involved to use & I wanted the underlying HTTP GET requests to be more robust to failure.

## Usage

```golang
hb := heartbeat.New("api-key")
hb.Ping(context.Background(), "name-of-the-heartbeat")
```
