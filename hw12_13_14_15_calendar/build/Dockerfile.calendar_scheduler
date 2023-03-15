FROM golang:1.17 as build

ENV BIN_FILE /opt/calendar/scheduler-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ${CODE_DIR}/cmd/calendar_scheduler/

FROM alpine:3.9

LABEL ORGANIZATION="OTUS"
LABEL SERVICE="calendar_scheduler"
LABEL MAINTAINERS="info@spendmail.ru"

ENV BIN_FILE "/opt/calendar/scheduler-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_FILE_NAME

ENV CONFIG_FILE /etc/calendar_scheduler/${CONFIG_FILE_NAME}.toml
COPY ./configs/${CONFIG_FILE_NAME}.toml ${CONFIG_FILE}

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

CMD /wait && ${BIN_FILE} -config ${CONFIG_FILE}
