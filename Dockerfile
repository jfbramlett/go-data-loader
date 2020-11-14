FROM golang:1.13-alpine3.11 as base

FROM base as build
RUN apk add --update make \
    && apk add --update git \
    && apk add --no-cache openssh-client
RUN git config --global url."git@github.com:".insteadOf https://github.com/
RUN mkdir -p /root/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
COPY id_rsa /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa
ENV GO111MODULE=on
ADD . /app
WORKDIR /app
RUN make build


FROM alpine:3.11 as runtime
COPY --from=build /app/bin/server /usr/local/bin
COPY web/ /root/web

ENV MOUNT_ROOT=/root
ENV WBROOT=/root/web

EXPOSE 8085

ENTRYPOINT ["/usr/local/bin/server", "serve"]