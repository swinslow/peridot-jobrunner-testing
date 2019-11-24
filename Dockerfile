# SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

FROM golang:1.13

RUN mkdir -p /peridot-jobrunner-testing
WORKDIR /peridot-jobrunner-testing

ADD . /peridot-jobrunner-testing

RUN go get -v ./...
RUN go build
