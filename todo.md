We recently rolled out a new feature called "Chirpy Red". It's a membership program, and members of "Chirpy Red" get pretty incredible features: like the ability to edit chirps after posting them. But that's beside the point...

Chirpy uses "Polka" as its payment provider. They send us webhooks whenever a user subscribes to Chirpy Red. We need to mark users as Chirpy Red members when we receive these webhooks.

[X] - Add a migration to the users table to include a new column called is_chirpy_red. This column should be a boolean, and it should default to false.

[X] - Add a database query that upgrades a user to chirpy red based on their ID.

[ ] - Add a POST /api/polka/webhooks endpoint. It should accept a request of this shape:

{
"event": "user.upgraded",
"data": {
"user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
}
}

[ ] - If the event is anything other than user.upgraded, the endpoint should immediately respond with a 204 status code - we don't care about any other events.
[ ] - If the event is user.upgraded, then it should update the user in the database, and mark that they are a Chirpy Red member.
[ ] - If the user is upgraded successfully, the endpoint should respond with a 204 status code and an empty response body. If the user can't be found, the endpoint should respond with a 404 status code.
Polka uses the response code to know whether or not the webhook was received successfully. If the response code is anything other than 2XX, they'll retry the request.

Update all endpoints that return user resources to include the is_chirpy_red field.

bootdev run 1304e939-bf50-48d3-a351-b35faafc267d

bootdev run 1304e939-bf50-48d3-a351-b35faafc267d -s

Default Base URL: http://localhost:8080
Optionally configure your CLI to override the default base URL by running bootdev config base_url <url>
Run the CLI commands to test your solution.

POST /admin/reset

1. Expecting status code: 200
   POST /api/users
   Request Body:
   {
   "email": "walt@breakingbad.com",
   "password": "123456"
   }

1. Expecting status code: 201
1. Expecting JSON at .email to be equal to walt@breakingbad.com
1. Expecting JSON at .is_chirpy_red to be equal to false
   Parsing userID variable from response body .id
   POST /api/polka/webhooks
   Request Body:
   {
   "data": {
   "user_id": "${userID}"
   },
   "event": "user.payment_failed"
   }

1. Expecting status code: 204
   POST /api/polka/webhooks
   Request Body:
   {
   "data": {
   "user_id": "${userID}"
   },
   "event": "user.upgraded"
   }

1. Expecting status code: 204
   POST /api/login
   Request Body:
   {
   "email": "walt@breakingbad.com",
   "password": "123456"
   }

1. Expecting status code: 200
1. Expecting JSON at .is_chirpy_red to be equal to true
