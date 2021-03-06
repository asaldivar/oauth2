package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "flag"
  "os"

  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

const (
  // clientId     = /* <CLIENT ID> */ 
  // clientSecret = /* <CLIENT SECRET> */
  // redirectURI  = /* <CLIENT SECRET> */
)

type Client struct {
  clientId     string
  clientSecret string
  redirectURI  string
  code         string
}

func NewClient(clientId, clientSecret, redirectURI string) *Client {
  return &Client{
    clientId:     clientId,
    clientSecret: clientSecret,
    redirectURI:  redirectURI,
  }
}

type Response struct {
  Data map[string]interface{}
}

// User authorizes app
func auth() {
  // Get User Info from command line
  var clientId *string = flag.String(
    "clientId",
    "",
    "CLIENT ID from http://instagram.com/developer/clients/manage/")

  var clientSecret *string = flag.String(
    "clientSecret",
    "",
    "CLIENT SECRET from http://instagram.com/developer/clients/manage/")

  var redirectURI = flag.String(
    "redirectURI",
    "",
    "REDIRECT URI from http://instagram.com/developer/clients/manage/")

  flag.Parse()

  if len(*clientId) == 0 || len(*clientSecret) == 0 || len(*redirectURI) == 0 {
    fmt.Println("Please provide your Instagram CLIENT ID, CLIENT SECRET and REDIRECT URI")
    os.Exit(1)
  }

  // create Client
  client := NewClient(*clientId, *clientSecret, *redirectURI)

  // parse Instagram's authorize endpoint
  authorizeEndpoint, _ := url.Parse("https://api.instagram.com/oauth/authorize/")

  // create necessary params for endpoint
  params := url.Values{}
  params.Add("client_id", client.clientId)
  params.Add("redirect_uri", client.redirectURI)
  params.Add("response_type", "code")

  // encode values into URL encoded form and append to endpoint
  authorizeEndpoint.RawQuery = params.Encode()

  // Give user configured IG Authorization endpoint
  fmt.Printf("\nGo to: %s", authorizeEndpoint.String() + "\n\n")

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

  fmt.Printf("\naccessToken: %s", accessToken)
  fmt.Printf("\nfullName: %s", userInfo["full_name"])
  fmt.Printf("\nuserName: %s", userInfo["username"].(string) + "\n\n")

  fmt.Fprintf(w, "Hello, %s!\n", c.URLParams["name"])
  fmt.Fprintf(w, "Your access token is: %s\n", accessToken)
  fmt.Fprintf(w, "User's username: %s\n", userInfo["username"])
  fmt.Fprintf(w, "User's full name: %s\n", userInfo["full_name"])
}

func main() {
  auth()
  goji.Get("/home/:name", home) 
  goji.Serve()
}