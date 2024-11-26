# Yikes!

### Compile protobuf

    otel_proto_path=/path/to/opentelemetry-proto
    yikes_path=/path/to/yikes/repo

    cd ~/opentelemetry-proto
    docker run --rm -u 501 \
        -v${otel_proto_path}:/mnt/opentelemetry-proto \
        -v${yikes_path}:/mnt/yikes \
        -w/mnt/opentelemetry-proto \
        otel/build-protobuf:0.9.0 \
        --proto_path=/mnt/opentelemetry-proto \
        --go_out=plugins=grpc:/mnt/yikes/proto \
        opentelemetry/proto/collector/trace/v1/trace_service.proto
