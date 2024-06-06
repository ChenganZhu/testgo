FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /myprogram

FROM gcr.io/distroless/static-debian12:debug AS build-release-stage

WORKDIR /

COPY --from=build-stage /myprogram /myprogram

USER nonroot:nonroot

ENTRYPOINT ["/myprogram"]
