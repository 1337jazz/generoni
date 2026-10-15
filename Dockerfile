FROM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/server ./cmd/server

RUN mkdir /out/data

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server
COPY --from=build --chown=nonroot:nonroot /out/data /data

WORKDIR /

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/server"]
