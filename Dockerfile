ARG GO_VERSION=1.25.5
FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/server .

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN mkdir logs && chmod 755 logs

COPY --from=build /bin/server /app/server
COPY users.json /app/users.json

ENTRYPOINT [ "/app/server" ]