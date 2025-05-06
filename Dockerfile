FROM golang:1.23-alpine as builder

ARG SSH_PRIVATE_KEY
RUN apk update && apk upgrade && apk --no-cache add ca-certificates bash \
    git gcc g++ pkgconfig build-base zlib-dev pkgconf openssh \
    && rm -rf /var/cache/apk/*
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa && chmod 400 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan gitlab.com >> /root/.ssh/known_hosts
RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

WORKDIR /go/src/gitlab.com/clubhub.ai1/organization/backend/payments-api/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -a -tags musl -installsuffix cgo -o app ./cmd

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata chromium
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/
WORKDIR /usr/
COPY --from=builder /go/src/gitlab.com/clubhub.ai1/organization/backend/payments-api/app .
CMD /usr/app
