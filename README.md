# Golang Web API Example

This is for a code challenge
## Quick Start

```shell
docker compose up
```

Wait for init container to run and exit

API is exposed on `localhost:8080`

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
