Documentation
When you're designing a server-side API, no one is going to know how to interact with it unless you tell them. Are you going to force the front-end developers, mobile developers, or other back-end service teams to sift through your code and reverse engineer your API?

Of course not! You're a good person. You're going to write documentation.

First Be Obvious, Then Document It Anyway
We've talked a lot about how your REST API should follow conventions as much as possible. That said, the conventions are not enough. You still need to document your endpoints. Without documentation, no one will know:

Which resources are available

- What the path to the endpoints are
- Which HTTP methods are supported for each resource
- What the shape of the data is for each resource
- etc.

## Assignment

One type of endpoint that's nearly impossible to interact with without documentation is a plural GET endpoint, that is, an endpoint that returns a list of resources. They often have different sorting, filtering, and pagination features.

[ ] - Update the GET /api/chirps endpoint. It should accept an optional query parameter called author_id.
[ ] - If the author_id query parameter is provided, the endpoint should return only the chirps for that author.
[ ] - If the author_id query parameter is not provided, the endpoint should return all chirps as it did before.

### For example:

GET http://localhost:8080/api/chirps?author_id=1

Continue sorting the chirps by created_at in ascending order.

Be sure to filter by author ID at the database level, not in memory! That will be more efficient on large datasets.

Run and submit the CLI tests.

Tips
The http.Request struct has a way to grab the query parameters from the URL:

s := r.URL.Query().Get("author_id")
// s is a string that contains the value of the author_id query parameter
// if it exists, or an empty string if it doesn't
