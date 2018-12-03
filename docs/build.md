# Build

## Requirements

- Go (tested on 1.11.2)
- Postgres (tested using postgres image 11.1)
- [dep](https://github.com/golang/dep)

## Make targets

Main targest at the `Makefile` are

- `make test` will test the code and output the cover file
- `make db` will pull and start a posgres container, then create the empty TODO list schema in it. *WARNING* it will whipe all data, use with caution
- `make run` will execute the application locally
- `make build` will generate the output binary for the current OS and architecure
- `make release` will execute tests and build for all OS and architecures

## Running rest-demo

To successfuly run rest-demo you will need to:

- Create the postgres schema using the `make db` target. If you have a postgres instance, you can execute `assets/deployment/database/schema.sql` on it.
- Copy `run/environment.template` as `run/environment`, and customize as needed. If you are using `make db` no customization for database is needed.
- Execute `make run`

## TODOs for an MVP

This repo haven't had a lot of time to work on, so these are the main issues to work at:
- fmt, vet and lint are not part of Makefile
- Tests coverage are in the very low side
- Watch implementation doesn't fit very well with Restful library, probably adding a subresource for it would sound better. Also, it needs to wrap the response to add the operation being watched (Create, Update, Delete) and allow filtering.
- When logging V(10), database logs look ugly. We need to re-arrange `\t\n` to make it look readable
- We should `test -race`, specially for the watch feature
- Watch feature would need `sync.RWMutex` when managing watchers
- I'm willing to use some of [this](https://github.com/heptio/contour/blob/master/Makefile), will need time to check one by one those tools
- There is no creation date/due date filter for tasks
- There is warning nor status change when a task is due
- There is no category management



