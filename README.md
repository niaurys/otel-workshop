# OpenTelemetry Workshop of Engineering Conference 2024

Repository provides environment for OpenTelemetry workshop held at Vinted Engineering Conference 2024.

## Prerequisites

- Docker
- [Docker Compose](https://docs.docker.com/compose/install/#install-compose) v2.0.0+
- 6 GB of RAM for the application
- [Golang](https://go.dev/doc/install) v1.22.0+

## Get and run workshop

1. Clone the workshop repository:

```shell
git clone https://github.com/niaurys/otel-workshop-2024.git
```

2. Change to the workshop folder:

```shell
cd otel-workshop-2024/
```

3. Start the workshop:

```shell
docker compose up --force-recreate --remove-orphans --detach
```

4. Rebuild everything and start again (after changes):

```shell
docker compose up --detach --build 
```

## Instrumentation packages

To get started with instrumentation flowing dependencies can be installed:

```shell
go get "go.opentelemetry.io/otel" \
  "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp" \
  "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" \
  "go.opentelemetry.io/otel/sdk/log" \
  "go.opentelemetry.io/otel/log/global" \
  "go.opentelemetry.io/otel/propagation" \
  "go.opentelemetry.io/otel/sdk/metric" \
  "go.opentelemetry.io/otel/sdk/resource" \
  "go.opentelemetry.io/otel/sdk/trace"
```

## Telemetry infrastructure

Once the images are built and containers are started you can access:

- Grafana: <http://localhost:3000>
- Jaeger UI: <http://localhost:16686>
- Prometheus: <http://localhost:9090>
- Kafka UI: <http://localhost:8081>

## Architecture

**OpenTelemetry Workshop** is composed of microservices written in Golang that talk to each other over gRPC or HTTP and uses clients of popular technologies like Kafka and Redis.

Please check `docker-compose.yaml` and `.env` files for exact configurations parameters.

The **Protocol Buffer Definitions** are in the `/pb/` directory.

```mermaid
graph TD
subgraph Service Diagram
you(*You*)
buyer(Buyer Service)
factory(Factory Service)
Redis[(Redis)]
shop(Shop Service)
warehouse(Warehouse Service)
Kafka[(Kafka)]

you --> |HTTP| buyer
factory --> |TCP| Kafka
warehouse --> |TCP| Kafka
warehouse --> |Redis| Redis
shop --> |Redis| Redis
buyer --> |gRPC| shop
buyer --> |HTTP| factory

end
```

It is possible to trigger Buyer service endpoint `/order` to place "order" for Factory service, which will be picked up by Warehouse service and randomly sold in Shop. For example:

```bash
curl http://localhost:3001/order \
  --data '{ "name": "watch", "color": "purple", "quantity":130}'
```

## Telemetry services architecture

The collector is configured in
[otelcol-config.yml](./config/otelcollector/otelcol-config.yml),
alternative exporters can be configured if needed.

```mermaid
graph TB
subgraph tdf[Telemetry Data Flow]
    subgraph subgraph_padding [ ]
        style subgraph_padding fill:none,stroke:none;
        %% padding to stop the titles clashing
        subgraph od[OpenTelemetry Demo]
        ms(Microservice)
        end

        ms -.->|"OTLP<br/>gRPC"| oc-grpc
        ms -.->|"OTLP<br/>HTTP POST"| oc-http

        subgraph oc[OTel Collector]
            style oc fill:#97aef3,color:black;
            oc-grpc[/"OTLP Receiver<br/>listening on<br/>grpc://localhost:4317"/]
            oc-http[/"OTLP Receiver<br/>listening on <br/>localhost:4318"/]
            oc-proc(Processors)
            oc-prom[/"OTLP HTTP Exporter"/]
            oc-otlp[/"OTLP Exporter"/]
            oc-elastic[/"ElasticSearch"/]

            oc-grpc --> oc-proc
            oc-http --> oc-proc

            oc-proc --> oc-prom
            oc-proc --> oc-otlp
            oc-proc --> oc-elastic
        end

        oc-prom -->|"localhost:9090/api/v1/otlp"| pr-sc
        oc-otlp -->|gRPC| ja-col
        oc-elastic --> |HTTP| es-index-api

        subgraph es[ElasticSearch]
            style es fill:#c27ba0,color:black;
            es-index-api[/"ElasticSearch<br/>Index API on<br/>localhost:9200"/]
            es-index[(*logs* index)]
            es-search-api[/"ElasticSearch<br/>Search API on<br/>localhost:9200"/]

            es-index-api --> es-index
            es-index --> es-search-api
        end

        subgraph pr[Prometheus]
            style pr fill:#e75128,color:black;
            pr-sc[/"Prometheus OTLP Write Receiver"/]
            pr-tsdb[(Prometheus TSDB)]
            pr-http[/"Prometheus HTTP<br/>listening on<br/>localhost:9090"/]

            pr-sc --> pr-tsdb
            pr-tsdb --> pr-http
        end

        pr-b{{"Browser<br/>Prometheus UI"}}
        pr-http ---->|"localhost:9090/graph"| pr-b

        subgraph ja[Jaeger]
            style ja fill:#60d0e4,color:black;
            ja-col[/"Jaeger Collector<br/>listening on<br/>grpc://jaeger:4317"/]
            ja-db[(Jaeger DB)]
            ja-http[/"Jaeger HTTP<br/>listening on<br/>localhost:16686"/]

            ja-col --> ja-db
            ja-db --> ja-http
        end

        subgraph gr[Grafana]
            style gr fill:#f8b91e,color:black;
            gr-srv["Grafana Server"]
            gr-http[/"Grafana HTTP<br/>listening on<br/>localhost:3000"/]

            gr-srv --> gr-http
        end

        pr-http --> |"localhost:9090/api"| gr-srv
        es-search-api --> |"localhost:9200/_search"| gr-srv
        ja-http --> |"localhost:16686/api"| gr-srv

        ja-b{{"Browser<br/>Jaeger UI"}}
        ja-http ---->|"localhost:16686/search"| ja-b

        gr-b{{"Browser<br/>Grafana UI"}}
        gr-http -->|"localhost:3000/dashboard"| gr-b
    end
end
```

## Useful links

- [OpenTelemetry documentation](https://opentelemetry.io/docs/)
- [OpenTelemetry Go documentation](https://opentelemetry.io/docs/languages/go/)
- [OpenTelemetry Go API and SDK](https://github.com/open-telemetry/opentelemetry-go)
- [OpenTelemetry Go API and SDK extensions](https://github.com/open-telemetry/opentelemetry-go-contrib)
- [Metric Supplementary Guidelines](https://opentelemetry.io/docs/specs/otel/metrics/supplementary-guidelines/)

Specification:

- [Configuration Environment Variables](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/)
- [OpenTelemetry Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/)
- [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/specs/otlp/)
- [Tracing](https://opentelemetry.io/docs/specs/otel/trace/)
- [Metrics](https://opentelemetry.io/docs/specs/otel/metrics/)
- [Logs](https://opentelemetry.io/docs/specs/otel/logs/)
- [Baggage](https://opentelemetry.io/docs/specs/otel/baggage/)

OpenTelemetry showcase:

- [OpenTelemetry demo](https://opentelemetry.io/docs/demo/)
