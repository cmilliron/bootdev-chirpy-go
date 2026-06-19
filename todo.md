[ ] Add a PUT /api/users endpoint so that users can update their own (but not others') email and password. It requires:
[ ] An access token in the header
[ ] A new password and email in the request body
[ ] Hash the password, then update the hashed password and the email for the authenticated user in the database. Respond with a 200 if everything is successful and the newly updated User resource (omitting the password of course).
[ ] If the access token is malformed or missing, respond with a 401 status code.
