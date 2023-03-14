# OpenTelemetry Spike
## Build and run
Run mod tidy in the root of the project
```shell
go mod tidy
```
### Using docker
Build app docker images
```shell
make buildAllImages
```
Run required services
```shell
make runRequiredServices
```
Run application
```shell
make run
```
## How to use
After you run the application you can send a request to the endpoint
```shell
curl -X POST http://localhost:8080/data -d '{"name": "test"}'
```
### Flow
1. The `main-app` sends a message to kafka
2. The `processor-app` reads a message from kafka and save into the memory

## How to check
### Jaeger
Open Jaeger UI
```shell
open http://localhost:16686/
```
## Example output for file exporter
Without error
```json
{
	"Name": "/data",
	"SpanContext": {
		"TraceID": "e03e88e7f18cf080801dde9481368647",
		"SpanID": "475e934ccdd2cb18",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 2,
	"StartTime": "0001-01-01T00:00:00Z",
	"EndTime": "0001-01-01T00:00:00Z",
	"Attributes": [
		{
			"Key": "http.method",
			"Value": {
				"Type": "STRING",
				"Value": "POST"
			}
		},
		{
			"Key": "http.scheme",
			"Value": {
				"Type": "STRING",
				"Value": "http"
			}
		},
		{
			"Key": "http.flavor",
			"Value": {
				"Type": "STRING",
				"Value": "1.1"
			}
		},
		{
			"Key": "net.host.name",
			"Value": {
				"Type": "STRING",
				"Value": "application"
			}
		},
		{
			"Key": "net.host.port",
			"Value": {
				"Type": "INT64",
				"Value": 8080
			}
		},
		{
			"Key": "net.sock.peer.addr",
			"Value": {
				"Type": "STRING",
				"Value": "127.0.0.1"
			}
		},
		{
			"Key": "net.sock.peer.port",
			"Value": {
				"Type": "INT64",
				"Value": 60996
			}
		},
		{
			"Key": "http.user_agent",
			"Value": {
				"Type": "STRING",
				"Value": "Apache-HttpClient/4.5.14 (Java/17.0.6)"
			}
		},
		{
			"Key": "http.route",
			"Value": {
				"Type": "STRING",
				"Value": "/data"
			}
		},
		{
			"Key": "http.status_code",
			"Value": {
				"Type": "INT64",
				"Value": 200
			}
		}
	],
	"Events": [
		{
			"Name": "Inserting data",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "0001-01-01T00:00:00Z"
		},
		{
			"Name": "log",
			"Attributes": [
				{
					"Key": "log.severity",
					"Value": {
						"Type": "STRING",
						"Value": "INFO"
					}
				},
				{
					"Key": "log.message",
					"Value": {
						"Type": "STRING",
						"Value": "Data inserted"
					}
				},
				{
					"Key": "code.function",
					"Value": {
						"Type": "STRING",
						"Value": "main.insertData.func1"
					}
				},
				{
					"Key": "code.filepath",
					"Value": {
						"Type": "STRING",
						"Value": "/Users/avlasov/git/blueprints/spike-go-opentelemetry-logging/pkg/app.go"
					}
				},
				{
					"Key": "code.lineno",
					"Value": {
						"Type": "INT64",
						"Value": 84
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": "0001-01-01T00:00:00Z"
		}
	],
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": null,
	"InstrumentationLibrary": {
		"Name": "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin",
		"Version": "semver:0.40.0",
		"SchemaURL": ""
	}
}
```
With error
```json
{
	"Name": "/data",
	"SpanContext": {
		"TraceID": "fcf00cff22cfd168020205e657b7eac2",
		"SpanID": "5cc23d4fb4f6c0a5",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 2,
	"StartTime": "0001-01-01T00:00:00Z",
	"EndTime": "0001-01-01T00:00:00Z",
	"Attributes": [
		{
			"Key": "http.method",
			"Value": {
				"Type": "STRING",
				"Value": "POST"
			}
		},
		{
			"Key": "http.scheme",
			"Value": {
				"Type": "STRING",
				"Value": "http"
			}
		},
		{
			"Key": "http.flavor",
			"Value": {
				"Type": "STRING",
				"Value": "1.1"
			}
		},
		{
			"Key": "net.host.name",
			"Value": {
				"Type": "STRING",
				"Value": "application"
			}
		},
		{
			"Key": "net.host.port",
			"Value": {
				"Type": "INT64",
				"Value": 8080
			}
		},
		{
			"Key": "net.sock.peer.addr",
			"Value": {
				"Type": "STRING",
				"Value": "127.0.0.1"
			}
		},
		{
			"Key": "net.sock.peer.port",
			"Value": {
				"Type": "INT64",
				"Value": 61013
			}
		},
		{
			"Key": "http.user_agent",
			"Value": {
				"Type": "STRING",
				"Value": "Apache-HttpClient/4.5.14 (Java/17.0.6)"
			}
		},
		{
			"Key": "http.route",
			"Value": {
				"Type": "STRING",
				"Value": "/data"
			}
		},
		{
			"Key": "http.status_code",
			"Value": {
				"Type": "INT64",
				"Value": 500
			}
		}
	],
	"Events": [
		{
			"Name": "Inserting data",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "0001-01-01T00:00:00Z"
		},
		{
			"Name": "log",
			"Attributes": [
				{
					"Key": "exception.type",
					"Value": {
						"Type": "STRING",
						"Value": "*errors.errorString"
					}
				},
				{
					"Key": "exception.message",
					"Value": {
						"Type": "STRING",
						"Value": "data already exists"
					}
				},
				{
					"Key": "log.severity",
					"Value": {
						"Type": "STRING",
						"Value": "ERROR"
					}
				},
				{
					"Key": "log.message",
					"Value": {
						"Type": "STRING",
						"Value": "Error inserting data"
					}
				},
				{
					"Key": "code.function",
					"Value": {
						"Type": "STRING",
						"Value": "main.insertData.func1"
					}
				},
				{
					"Key": "code.filepath",
					"Value": {
						"Type": "STRING",
						"Value": "/Users/avlasov/git/blueprints/spike-go-opentelemetry-logging/pkg/app.go"
					}
				},
				{
					"Key": "code.lineno",
					"Value": {
						"Type": "INT64",
						"Value": 79
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": "0001-01-01T00:00:00Z"
		}
	],
	"Links": null,
	"Status": {
		"Code": "Error",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": null,
	"InstrumentationLibrary": {
		"Name": "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin",
		"Version": "semver:0.40.0",
		"SchemaURL": ""
	}
}
```