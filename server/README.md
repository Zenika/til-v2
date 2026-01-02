# Today I Learned (TiL) - v2

TIL is a simple software that aim to enable article and reflexion-sharing across all Zenika agencies in the world. It
acts as an aggregator (as
daily.dev can do), but on-premises and with filters on content to avoid being spammed with news that don't interest us.


## Getting started

### Developer tools

If you would like to develop TiL, the best way is to simply install Go on your computer ; then run
`go run -tags=viper_bind_struct ./cmd`. The server
will start on port 8000.

### Docker

If you simply want to run the application, Docker is the best way to do. Simply run `docker build -t til:latest .`, and
run it! Full configuration
reference below. Please think to mount database file in your container if you want persistent data!

## Configuration

All the configuration is done through environment variables. Viper is in charge of automatically parse them and load
them in our Go application. The
available environment variables are the following:

* **(Mandatory)** `TIL_JWT_SECRET`: The secret key used to sign JWT tokens issued by the app.
* **(Mandatory)** `TIL_GOOGLE_CLIENT_ID`: The Client ID for Google Authentication
* **(Mandatory)** `TIL_GOOGLE_CLIENT_SECRET`: The Client Secret for Google Authentication
* *(Optional)* `TIL_DEBUG` (default: `false`): Set this variable to `true` if you want to enable debug logs on
  application. Please be careful: this mode will log WAY MORE information, even secret tokens! The debug mode will also
  disable token expiry time check, which may lead to a lot of JWT renewal. Do not use this in production.
* *(Optional)* `TIL_DATABASE_FILE_NAME` (default: `til.sqlite3`): Specify database file name (and location).
* *(Optional)* `TIL_SERVER_PORT` (default: `8000`): Server port to listen on
* *(Optional)* `TIL_GOOGLE_TOKEN_ENDPOINT` (default: `https://oauth2.googleapis.com/token`): The Google Authentication
  endpoint to validate user token
* *(Optional)* `TIL_DEFAULT_ADMIN` (default: none): The user to immediately promote as admin. Use his Google ID in this
  variable to make it a superuser of the application. This variable will only take effect if ALL the following conditions are met :
  * The variable is set to a non-void value;
  * There is no other user (admin or not) in the database.
* *(Optional)* `TIL_USE_EMBEDDED_FRONTEND` (default: `false`): If `true`, all content in `/static` will be served on root path.


## Tests

Tests are made using Gherkin (Cucumber). This allows non-technical people to understand and collaborate to tests,
improving security and coverage.

All the tests can be find in `test/` folder. We are using a mock for the Google Authentication Server so you'll never
have to provide a secret to run tests. The `test/` folder includes some secrets (notably an RSA key and a JWT secret
key): these secrets are dedicated to our tests, and have never been used in production. Please do not report them as a
security issue; we're perfectly aware of it. Obviously, do not use these secrets for your own production (otherwise,
magic will probably happen!).

All our tests are containerized to avoid configuration issues. To run them, simply type `go test -v ./test/`, and Golang
will do the rest.

The following sentences are available :

#### Authentication

* `I have a JWT token for my regular|admin user`
    * This call will connect a user with the defined privilege. Every subsequent API calls will be authenticated using
      this JWT token.
* `I clear the current JWT token`
    * Remove the current JWT token (if any). Every subsequent API calls will be unauthenticated.

#### Database

* `I reset the database`
    * Reset the database by restarting the container

#### HTTP

* `I send a "GET|POST|PUT|DELETE" request to "ENDPOINT"`
    * Send a request to the defined endpoint (an endpoint is an API path, for example `/api/posts`).
* `I send a "GET|POST|PUT|DELETE" request to "ENDPOINT" with payload`
    * Same as before, but add a payload to the request. The payload must be placed on the next line, between `"""`. See
      example below or in `test/features/test.feature` for more details.
* `the response code should be XXX`
    * Check that the previous request sent have the defined status code.
* `the response should have XXX items in path "PATH"`
    * **WORKS ONLY ON ARRAYS** Count the number of items returned by the server in the JSON file.
    * The path is the place on the JSON to look, dot-separated. `@` is the main document. For example, if you want to
      access the `items` part of your JSON, the path will be `@.items`.
* `I save the "XXXXX" header as "XXXXXX" for suite`
    * Save the content of a header in a variable for later re-use. For example,
      `I save the "X-Id" header as "postId" for suite` will allow you to write
      `I send a "GET" request to "/api/posts/{{postId}}"` later on.
* `I save the value "XXX" in path "YYY" as "ZZZ" for suite`
  * Save the value XXX located in path YYY in a variable for later reuse. For example, `I save the value "id" in path "@" as "userId" for suite`
* `the response should have the following content in path "XXX"`
    * **WORKS ONLY ON OBJECTS** Check that the specified values for each key are correct.
* `the response should have the following items in path "XXX"`
    * **WORKS ONLY ON ARRAYS** Same as the previous one, but for arrays.
* `the response should have the key "XXX" in path "YYY"`
    * Check if the response contains a specified key in the desired path. For example `the response should have the key "id" in path "@.items.0"`
* `the response should not have the key "XXX" in path "YYY"`
    * The exact opposite as before ;-)


### Example

In this example, we log in as a regular user, then create a post before accessing it.

```gherkin
Feature: Demo

  Scenario: Demo
    When I have a JWT token for my regular user
    And I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Test",
        "link": "https://lamontagne.fr",
        "tags": ["lang:fr"]
      }
      """
    And the response code should be 201
    And I save the "X-Post-Id" header as "postId" for suite
    And I send a "GET" request to "/api/posts/"
    And the response code should be 200
    And the response should have the following content in path "@"
      | total_items | total_pages | current_page | items_per_page |
      | 1           | 1           | 0            | 20             |
    And the response should have the following items in path "@.items"
      | id         | title | link                  |
      | {{postId}} | Test  | https://lamontagne.fr |
    And the response should have 1 items in path "@.items.0.tags"
    And the response should have the following items in path "@.items.0.tags"
      | lang:fr |
```

## Will it scale?

TiL is conceived following the KISS way of life; but it doesn't mean it won't scale. During our tests, we ran the
following Locust file, with 50 concurrent users running 131 RPS:

```py
from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    @task
    def index(self):
        self.client.get("/api/posts", headers={"Authorization": "token"})


    @task
    def publish(self):
        self.client.post("/api/posts", headers={"Authorization": ""}, json={"title":"essai", "link":"https://www.lamontagne.fr/paris-75000/actualites/coupure-de-courant-geante-habitants-dans-la-rue-metros-et-avions-a-l-arret-les-images-du-chaos-en-espagne-et-au-portugal_14679007/", "tags": ["failure", "electricity", "lang:fr"], "content": "severe electricty shortage in Spain and Portugal"})
```

During 10 minutes, the locust file generated roughly 30k entries in database, and the same amount of article queried.
The Locust file ended with 0 failure and a MTRR around 600ms, which is a good result stating that our test is far above
the real frequentation.
