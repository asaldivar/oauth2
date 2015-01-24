#Go OAuth2 for Instagram

###*(in progress)*

Some Go code I put together to help myself start to learn Go as well as help others with their pursuit to integrate Instgram OAuth2 into their applications.

First, you'll need to setup your [instagram client](http://instagram.com/developer/clients/manage/) with the proper `REDIRECT URI` that you'll use as an endpoint on line 134 of the oauth2.go file:

```
goji.Get("/home/:name", home)
```

Then, on lines 17-19 insert your `CLIENT ID`, `CLIENT SECRET` and `REDIRECT URI`.

After this configuration is done you run the command:

`go run oauth2.go --clientId <CLIENT ID> --clientSecret <CLIENT SECRET> --redirectURI <REDIRECT URI`

Follow the steps in the terminal and the User's full name and username will be printed to the console and viewed in the browser, as well as the access token which you can use to hit other endpoints for more data.

*this is currently a work in progress, but hopefully it can still help out. please feel free to grab snippets or use it all.*