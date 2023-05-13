FROM golang:1.17-buster AS build
WORKDIR /src
COPY . .

# Must build without cgo because libc is unavailable in runtime image
ENV GO111MODULE=on CGO_ENABLED=0
RUN make build

FROM scratch
EXPOSE 8080

COPY --from=build /src/geo-ip-service /opt/bin/geo-ip-service
COPY --from=build /src/data /opt/data
WORKDIR /opt
ENTRYPOINT ["/opt/bin/geo-ip-service"]
