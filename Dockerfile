FROM golang:1.19-alpine as build

RUN apk add --update make git build-base
RUN apk --no-cache add ca-certificates

COPY . /gh_action_listener/
WORKDIR /gh_action_listener/

RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /listener/
COPY --from=build /gh_action_listener/listener .
COPY --from=build /gh_action_listener/listener.yaml .

CMD ["./listener"]
