package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"

  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

const (
  clientId     = "de761aa8f066479fb7ea069396ae50b5"
  clientSecret = "0dd10b36d467450aaad644ce44e51028"
  redirectURI  = "http://localhost:8000/home/alex"
)

type Response struct {
  Data map[string]interface{}
}

// User authorizes app
func auth(c web.C, w http.ResponseWriter, r *http.Request) {
  // parse Instagram's authorize endpoint
  authorizeEndpoint, _ := url.Parse("https://api.instagram.com/oauth/authorize/")

  // create necessary params for endpoint
  params := url.Values{}
  params.Add("client_id", clientId)
  params.Add("redirect_uri", redirectURI)
  params.Add("response_type", "code")

  // encode values into URL encoded form and append to endpoint
  authorizeEndpoint.RawQuery = params.Encode()

  // redirect to to configured authorize endpoint
  http.Redirect(w, r, authorizeEndpoint.String(), http.StatusMovedPermanently)

}

// request access token from service provider
func home(c web.C, w http.ResponseWriter, r *http.Request) {

  // grab code value from URL
  code := r.URL.Query()["code"][0]

  // ping access_token endpoint with appropriate data to get public user info and access token which can be used for future requests
  resp, err := http.PostForm("https://api.instagram.com/oauth/access_token",
    url.Values{
      "client_id"     : {clientId},
      "client_secret" : {clientSecret},
      "grant_type"    : {"authorization_code"},
      "redirect_uri"  : {redirectURI},
      "code"          : {code},
    },
  )
  if err != nil {
    fmt.Println("There's an error")
    fmt.Println(err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
  }
  // parse json-encoded data
  var v Response
  error := json.Unmarshal(body, &v.Data)
  if error != nil {
    fmt.Println("error:", err)
  }
  // grab access token and user info
  accessToken := v.Data["access_token"]
  userInfo    := v.Data["user"].(map[string]interface{})

  fmt.Println("accessToken:",accessToken)
  fmt.Println("fullName:",userInfo["full_name"])
  fmt.Println("userName:",userInfo["username"])

  fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func main() {
  goji.Get("/instagram/auth", auth)
  goji.Get("/home/:name", home)
  goji.Serve()
}