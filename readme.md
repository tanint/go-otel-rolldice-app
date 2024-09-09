# OpenTelemetry Instrumentation PoC

This project demonstrates how to set up OpenTelemetry instrumentation in a Go application.


# Screeshots
### Tracing
![image](https://github.com/user-attachments/assets/d18dd942-6de0-4f6d-bd03-cb704dfb668f)

### Logging
![image](https://github.com/user-attachments/assets/2e954946-cf6b-4d0d-8776-9af3822b7e85)


### Service Map
![image](https://github.com/user-attachments/assets/0de9389d-94bc-4789-a338-2a36246937d1)

### Environment example
| Environment Variable             | Description                        |
|-----------------------------------|------------------------------------|
| `PORT`                           | Port number for the application    |
| `APP_ENV`                        | Application environment (e.g., dev, prod) |
| `SERVICE_NAME`                   | Name of the service                |
| `OTEL_EXPORTER_OTLP_ENDPOINT`     | OpenTelemetry OTLP exporter endpoint |
| `OTEL_HTTP_EXPORTER_AUTH_TOKEN`   | OpenTelemetry HTTP exporter auth token |
| `NOTIFICATION_SERVICE_NAME`       | Name of the notification service   |
| `LINE_BOT_API_AUTH_TOKEN`         | Line bot API authentication token  |
| `LINE_BOT_RECEIVER_ID`            | Receiver ID for Line bot messages  |
| `KAFKA_USERNAME`                  | Username for Kafka authentication  |
| `KAFKA_PASSWORD`                  | Password for Kafka authentication  |
| `KAFKA_BROKERS`                   | Kafka brokers (comma-separated)    |

