# WordOfWisdom

## Linux

The code has been tested only on Ubuntu 18.04

## Requirements

Client mode `local` allows you to type commands to request quotes from the server, where `docker` mode sends one quote and that's it
export CLIENT_MODE=local
OR
export CLIENT_MODE=docker

## Run in docker

1. make build-docker
2. make run-on-docker

## Testing

1. Run all tests:

    make test-with-component
2. Run only unit tests:

    make test

## Other Makefile commands

1. make start-server
2. make start-client

## PoW

I choose to work with CPU-bound function - hashcash.
I could find most documentation about this pow scheme and
it was perfect for the task requirements.
