FROM golang:1.21.0 AS build

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY ./ /app/

RUN mkdir -p /app/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/rest-api


FROM scratch

COPY --from=build /app/build/rest-api /usr/bin/rest-api

ENTRYPOINT [ "/usr/bin/rest-api" ]