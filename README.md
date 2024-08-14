# Golang Web API Example

This is my first attempt at building a Golang json API from scratch.  I choose to make an attempt in Golang so you can have an idea of how I might instinctually approach the problem.

I did not use Copilot or any AI tools, and stuct to looking mostly at Gin's API documentation, Golangs official docs, and a few blog articals for project structure inspiration.


## Quick Start

```shell
docker compose up
```

Wait for init container to run and exit

API is exposed on `localhost:8080`

## Follow along

The main entrypoint for the application is in `main.go`.

The endpoints can be found in the controllers directory.  Here I have split the controllers into three seperate files, though once complete, found it would have probably made more sense to put them all under `user_controller.go` since all the functionality hangs off of the `/user` endpoint.

Endpoints follow a pattern where user input is validated, a service is called, and the service response is turned into an appropriate respose for a json API.

All of the buisness logic is grouped into services which should make creating tooling, such as CLI interfaces, easier.  A CLI would just need to validate user input from the CLI, call the same service, and return a response suitable for a CLI.

The services are quite thin since there is very little application logic, or complex queries.  Each service call is grouped into a transaction (even queries).  This is to express the idea that every service call should be considered something that has a transaction boundry, and not simply a wrapper around a single repository call (which many service calls ended up being).

The usefulness of the service pattern here ended up being minimal, but I committed to the idea since it is something I feel a larger project should evolve towards.

I ended up with one service, since I feel the User object should be the aggragate through which all data is accessed and modified.  I didn't end up taking this concept as far as I would have liked, and kept having the nagging feeling I would have been much happier using a document datastore as a backend since thing about the table and queries ended up flavoring how I wanted my apps data model to work.

It has been a while since I have used raw SQL instead of an ORM, but I tought this the perfect excuse to brush up on organizing code around raw sql queries.

I abstracted all querying into repositories and overall am pleased with how it turned out, though I had hoped to end up with a more epressive user model, what I ended up with fits the scope of the project.

```go
// What I was thinking of achiving... is pretty much an ORM 😅
//   Probably defeats the purpose of hand rolling all your own queries
//   And should just use an ORM to begin with
func unfriendPeopleBetterThanYou(userRepo IUserRepo) {
    user := userRepo.Get("id")
    for _, friend := range user.Friends() {
        if friend.Score() > user.Score() {
            user.Unfriend(friend)
        }
    }

    // ...
}
```

I will likely continue with this project to try and achive the above out of personal interest.

https://github.com/LucasRxTx/golang-rest-api-example

## Improvements

> Your friend asks for any suggestions you might have to improve the API in future versions of the game...

### Authentication and Authorization

The API is currently insecure, and anyone can view and change data of any user.

I would suggest JWT's and role based auth.

Admin features where hinted at in the "List all users" endpoint specification.  There are other considerations to take into account for admin activities, for example logging all change events, whom made them, on whos account.

### Error Handling

The API specification made no note about what an error response should contain.  Having uniform error responses makes writing error handling in client code easier.

For example https://www.rfc-editor.org/rfc/rfc7807

### Cacheability

Adding E-Tag headers allows for response cacheability client side, minimizing data transfered over the network, and improving overall speed.

Instead of returning a response body with the creation of a resource, consider redirecting to the newly created resource, or returnig a location header so the client can decide if and when to fetch the new resource with a GET request.

GET requests with E-Tag and Cache-Control headers can often be handled automaticlly by browser based clients, and proxies such as Cloudflare, or NGINX.

The user data is highly cacheable and does not change often, as well as the friends list.  Caching this data at the edge would greately improve response times and server load.

The origin server can maintain it's own cache, and return empty 304 Not Modified responses which would cut down on data egress costs, and improve response times.

### Verbs

The POST for /user/:id/relationship should be a PUT.

The behavior expected is for whatever is currently on the server to be replaced with the given details.

The HTTP verb that most closely matches that behavior is PUT.

A PUT allows for upserts, and replacements of current resorces with the provided data.

### Response envaloping

The friends list is unnessisarilly envaloped in an object with a friends property.  Simply returning a list of friend id's is sufficient, and is what would typically be expected.

### Publishing events

Publishing events from the system will allow other systems to build on events of the API.

For example

- Publishing a UserCreated event
    - Triggers a welcome email
    - Gets logged in a Buisness Analytics system
        - Did a new campaign drive new signups?
- Publishing a GameStateUpdated event
    - Did the user hit a high score?
        - Ask the user to brag on social
    - Did the user start playing poorly?
        - Send a one time offer for a boost
- Publishing a UserUpdatedFriends event
    - Notifiy the other user they are now friends
    - Give a promotional offer

### Internationalization / Language Support

Does the game have a broad international audiance?  If so, will the frontend be responsible for all text presented to the user, mapping backend messages, to messages suitable for the user, and in thier language, or should the API provided responses that are translated and user ready?

This includes any string representations of numbers and dates.

Language choice can be handled through the Language header which is automatically set by most browsers and is a list of languages ordered by the users preference.  In non-browser settings, it can manually be included in the Language header.

If the company has a translation server/API we can integrate with that, or, we can use classic i19n gettext options, or other libraries.

## Personal notes about my implementation

### UUIDs

Should probably use a ULID instead of a UUID.  ULIDs are like a UUID but they contain a timestamp prefix that makes them sortable which can improve performance for certain queries.

### Timestamps

Should likely consider storing created at, and last modified timestamps.

### Overwriting state

Should very much consider keeping a log of GameState for data analytics purposes, and potential game features later on.  Allows the system to answer questions such as

- Is the player getting better or worse over time?
- When does the player log in to look around, instead of log in to play a game?
- What time of day is the player the most active?

### User input validation

The validation library that comes with Gin left alot to be desired

https://github.com/go-playground/validator

The error messages returned reference the struct field names in go, which are not helpful when returned to the end user who will want to see the json field names.  The library gives no out of the box solution and expects you to implament your own handler.  I decided this was outside of the scope for a code interview due to time restraints.

### Repeated db connection handling code

There was repeated code for transaction handling that I tried to refactor into something reusable.  It kind of worked, but the Generics make the code noisy.  It is still a win over accidently missing a Rollback, but there is room for improvement.

### Testing

Concreat implamentations of sql.Db and sql.Transaction are used everywhere making it hard to test.  A wrapper should likely be used to allow a mock to be injected during testing.

The API as it is now would not make much sense to test without a database since most of the logic is about storing and retreiving from a database.  Testing the code would result in mostly testing mocks and fakes.

If I where to write tests for the system, I would write end to end tests (tests requiring the full system to be running, database included), minimal unit tests, and contract tests (Is the API returning data in a format that is expected given the spec).
