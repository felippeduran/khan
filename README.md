# Khan

[![Build Status](https://travis-ci.org/topfreegames/khan.svg?branch=master)](https://travis-ci.org/topfreegames/khan)
[![Coverage Status](https://coveralls.io/repos/github/topfreegames/khan/badge.svg?branch=master)](https://coveralls.io/github/topfreegames/khan?branch=master)
[![Code Climate](https://codeclimate.com/github/topfreegames/khan/badges/gpa.svg)](https://codeclimate.com/github/topfreegames/khan)
[![Go Report Card](https://goreportcard.com/badge/github.com/topfreegames/khan)](https://goreportcard.com/report/github.com/topfreegames/khan)

Khan will drive all your enemies to the sea (and also take care of your game's clans)!

## Setup

Make sure you have go installed on your machine.
If you use homebrew you can install it with `brew install go`.

Run `make setup`.

## Running the application

Create the development database with `make migrate` (first time only).

Run the api with `make run`.

## Running with docker

Provided you have docker installed, to build Khan's image run:

    $ make build-docker

To run a new khan instance, run:

    $ make run-docker

## Docker Image

You can get a docker image from our dockerhub page at https://hub.docker.com/r/tfgco/khan/.

## Tests

Running tests can be done with `make test`, while creating the test database can be accomplished with `make drop-test` and `make db-test`.

## Coverage

Getting coverage data can be achieved with `make coverage`, while reading the actual results can be done with `make coverage-html`.

## Static Analysis

Khan goes through some static analysis tools for go. To run them just use `make static`.

Right now, gocyclo can't process the vendor folder, so we just ignore the exit code for it, while maintaining the output for anything not in the vendor folder.
