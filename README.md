# seeder

## seeder main features:

1. Receive queries from clients and returns a list of available nodes.

2. Receive “hello” message from other clients which want to introduce themselves to the Seeder. The seeder will add such clients to the peer list.

3. Frequently ping existing nodes in the database to update their information (e.g. check whether they are alive, which clients they are using).

4. Support queries like: which nodes are alive during the last 1 hours, 2 hours, 1 day.

## Architecture

Uncle Bob Clean Code architecture is used for seeder apps

## How do we run?

### To run tests
`make test`

### To run seeder server
`make run-seeder`
OR
`make local-seeder`

### To run seeder client
`make run-client PORT=27001`
OR
`make local-client PORT=27001`

### To list nodes
`curl --location --request GET 'http://localhost:5000/v1/nodes'`

### To list nodes alive for 30 seconds
`curl --location --request GET 'http://localhost:5000/v1/nodes?alive=true&age=30s'`

### To send hello message
`curl -X POST 'http://localhost:5000/v1/nodes' -H 'Content-Type: application/json' -d '{"ip":"192.168.0.12:27002", "name":"testClientName2", "client":"testClient", "version":"v1.0.1"}`
