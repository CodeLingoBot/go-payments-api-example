# Go API Payments Example

[![Coverage Status](https://coveralls.io/repos/github/ifraixedes/go-payments-api-example/badge.svg?branch=master)](https://coveralls.io/github/ifraixedes/go-payments-api-example/?branch=master)
[![Build Status](https://travis-ci.com/ifraixedes/go-payments-api-example.svg?branch=master)](https://travis-ci.com/ifraixedes/go-payments-api-example)
[![Go Report Card](https://goreportcard.com/badge/github.com/ifraixedes/go-payments-api-example)](https://goreportcard.com/report/github.com/ifraixedes/go-payments-api-example)
[![GoDoc](https://godoc.org/github.com/ifraixedes/go-payments-api-example?status.svg)](https://godoc.org/github.com/ifraixedes/go-payments-api-example)


This repository contains an example of a payments API implemented with Go serving the purpose of an example with the packages, tools, technologies, etc. that are used.

NOTE the packages, tools, technologies, etc., will be listed in the future when the first implementation be done; any other important information to highlight will be added to this document too.

## Goal

Implement, in Go, a small example of a payment RESTful API which allows to:

* Create, read, update and delete a payment.
* Get a list of payments.

## Considerations

The project will be used in the future of a base of experimentation for improving several topics, between some of them, API designs, Software Patterns, DDD (Domain Driven Development), modularity, simplicity.

However all of those changes will be done through the time under my needs and curiosities and without any commitment and iterative in order to have have "releases".

This example isn't intended for being in production because it is already known that it reqiures some improvements which will never been done, because this is just an example; in addition, in the sources you may find a `TODO: won't be implemented` to mark parts which exists because it clafies how some feature/requirement should be done, however, such implementation won't be done.

Following there is a incomplete list of improvements which should be done before using this implementation in a production system.

* Add more input validations


NOTE that some decions have been made without having any specific requirement in terms of business domain, production environment, SLAs, etc., so they may look to you that those are improvements to be made, however, they are not because of the lack of such requirements, that's the reason because they aren't in the list and they will never be, others they aren't because as commented it's an incomplete list, because of the fact that I may forget or may not think about them. If you think that this list should contain some which aren't on it, feel free to send an issue or PR, but be aware that it could be rejected because they may not be improvements, they may be some of those decisions which were made.

## TODO

These will be tracked in the issue tracker ones the branch be sent in a PR.

* Use context in the sqlite implementation using the [`Iterrupt` method of the connection](https://godoc.org/github.com/bvinc/go-sqlite-lite/sqlite3#Conn.Interrupt).


## License

MIT, read [the license file](LICENSE) for more information.
