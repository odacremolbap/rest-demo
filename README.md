# rest-demo

This is a REST demo of a TODO list implemented using Go and Postgres

It uses:

- [Restful(https://github.com/emicklei/go-restful) for REST definition and routing
- [Logrus](https://github.com/sirupsen/logrus) wrapped as [Logr](https://github.com/go-logr/logr) for logging
- [pq](github.com/lib/pq) as Postgres database driver
- [pkg/errors](github.com/pkg/errors) for errors stacks and wrapping
- [cobra](github.com/spf13/cobra) for command and flags management

For building and running instructions check [instructions](docs/build.md)

 ## TODO List features

 This is a very simple TODO list with only one entity, `Tasks`

You can CRUD using:

- `GET http://localhost:9101/v1/tasks` for listing all tasks
- `GET http://localhost:9101/v1/tasks/<id>` for retrieving a task
- `POST http://localhost:9101/v1/tasks + <JSON Payload>` to create a task
- `PUT http://localhost:9101/v1/tasks/<id> + <JSON Payload>` to update a task
- `DELETE http://localhost:9101/v1/tasks/<id>` to delete a task

All tasks listing can be filtered by `category`, `name` or `status` adding any of those fields and the exact value at the URL query

- `GET http://localhost:9101/v1/tasks?category=longterm` would return `longterm` category tasks

Task deletion is logical by default. To make it a physical database deletion it must be appended `permanent=true` URL query

- `DELETE http://localhost:9101/v1/tasks/3` would set task 3 status to deleted
- `DELETE http://localhost:9101/v1/tasks/3?permanent=true` would delete task 3 from the database

There is also a watch feature that returns the processed Task:

- `http://localhost:9101/v1/tasks?watch` would block and list all Tasks processed

(Unfortunately it returns a Task, but for now no info about the operation executed, I'll add a wrapper with the operation at some point. Also, when deleting a task, it returns the last state of the task, which by that time no longer exists)

You can find some handy `curl` examples [here](assets/curl)