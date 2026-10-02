# jwt vs session authentication


# stateful vs stateless
-> server need to store the some state -> stateul
-> suppose in session based
-> user session is stored in database
-> session id is stored in browser 
-> user authenticate by sending session id
-> now it can easily be revoked
-> but not scalable as stateful
-> jwt is also not stateless
-> it can be stateful/ statelesss
-> jwt has 
-> a.b.c
-> a -> header
-> b -> payload
-> c -> signature = HMAC(
    secret,
    base64url(header) + "." + base64url(payload)
)
-> c -> has is created using secret and payload payload changes 
-> then you need to make new c but as you donot have secret 
-> so you cannot do it if you put something wrong then when you will send it to server
-> server will take payload and use secret to create c 
-> if it is same then correct otherwise not
-> it is designed to get forward hash not backward
-> so can use match it forward only not backward
-> the secret will be shared across all services

-> refresh token
-> access token expired
-> get 401
-> automatic handling
-> client will send request to auth service
-> to create new token using refresh token
-> refresh token hash can be stored in db
-> it can be deleted if this token is invlidated
-> if user request then he didnot get any row with this hash
-> it means token doesnot exist/invalidated
-> this refresh token can be refreshed also on each new access token refresh