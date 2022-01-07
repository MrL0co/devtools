# syntax=docker/dockerfile:1

FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
#COPY go.sum ./

RUN go mod download

COPY . ./

RUN make build


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/pinstall /pinstall

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/pinstall" ]
