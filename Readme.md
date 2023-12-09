## Unit Testing HTTP Endpoints

THis is a codebase written to accompany a blog post I'm working on.  The code is a modified version of [this example]() in the `go-restful` repo, with changes made to focus on testing.

To run locally, pull this repo and run the following command:
```shell
$ go run .
```

That will start a local server on port `:8080`, which you can then interact with via REST requests:
- GET /users
- GET /users/{id}
- POST /users
- PUT /users
- DELETE /users/{id}

To run the tests, ensure you have it pulled locally, and run the following command:
```shell
$ go test ./... -coverprofile=coverage.out
```

Once the output file has been created, you can visualize line coverage by running:
```shell
$ go tool cover -html coverage.out -o coverage.html
```
This will generate an HTML file that you can view in your browser.
