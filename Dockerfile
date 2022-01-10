# syntax=docker/dockerfile:1
FROM golang:1.17.5-alpine as builder

RUN mkdir /opt/crowsnest

WORKDIR /opt/crowsnest

COPY ./ ./

RUN go build -o crowsnest .

FROM golang:1.17.5-alpine

WORKDIR /root/

COPY --from=builder /opt/crowsnest/crowsnest ./

CMD [ "./crowsnest" ]
