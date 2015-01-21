package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"

  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

func home(c web.C, w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

const (
  clientId     = "de761aa8f066479fb7ea069396ae50b5"
  clientSecret = "0dd10b36d467450aaad644ce44e51028"
  redirectURI  = "http://localhost:8000/hello/alex"
)

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

  // http.Redirect(authorizeEndpoint.String())
  http.Redirect(w, r, authorizeEndpoint.String(), http.StatusMovedPermanently)

}

func main() {
  goji.Get("/instagram/auth", auth)
  goji.Get("/home/:name", home)
  goji.Serve()
}





  // GET the endpoint
  // resp, err := http.Get(authorizeEndpoint.String())
  // if err != nil {
  //   fmt.Println(err)
  // }
  // defer resp.Body.Close()
  // body, err := ioutil.ReadAll(resp.Body)
  // if err != nil {
  //   fmt.Println(err)
  // }

  // print the response body
  // fmt.Println(string(body))