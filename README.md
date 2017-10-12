# go-fastly-purge

[![godoc](https://godoc.org/github.com/zencoder/go-fastly-purge?status.svg)](http://godoc.org/github.com/zencoder/go-fastly-purge)
[![Circle CI](https://circleci.com/gh/zencoder/go-fastly-purge.svg?style=svg)](https://circleci.com/gh/zencoder/go-fastly-purge)
[![Coverage Status](https://coveralls.io/repos/zencoder/go-fastly-purge/badge.svg?branch=master)](https://coveralls.io/r/zencoder/go-fastly-purge?branch=master)

A Go library for making Purge requests against the Fastly CDN.

Fastly documentation at https://docs.fastly.com/api/purge

Install
-------

This project uses [Glide](https://github.com/Masterminds/glide) to manage it's dependencies. Please refer to the glide docs to see how to install glide.

```bash
mkdir -p $GOPATH/github.com/zencoder
cd $GOPATH/github.com/zencoder
git clone https://github.com/zencoder/go-fastly-purge
cd go-fastly-purge
glide install
go install ./...
```
