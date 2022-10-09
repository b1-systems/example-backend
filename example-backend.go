/* Demonstration of passing an OpenID Connect ID Token to a web application via Authorization header. */

package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"
  "time"
  "github.com/coreos/go-oidc/v3/oidc"
  "golang.org/x/net/context"
  "example-backend/ini"
)

var (
  clientName = "example-backend"
  clientID = ""
  clientSecret = ""
  providerUrl = ""
  listenAddress = ""
)

func main() {
  arr := []ini.Ref{
    {"clientID", &clientID},
    {"clientSecret", &clientSecret},
    {"providerUrl", &providerUrl},
    {"listenAddress", &listenAddress}}

  err := ini.ReadIni(clientName, arr)

  if err != nil {
    log.Fatal()
    os.Exit(1)
  }

  ctx := context.Background()
  provider, err := oidc.NewProvider(ctx, providerUrl)

  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  oidcConfig := &oidc.Config{
    ClientID: clientID,
  }

  verifier := provider.Verifier(oidcConfig)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    auth_header := r.Header.Get("Authorization")

    if auth_header == "" {
      log.Printf("Authorization header missing from request")
      http.Error(w, "Bad request", http.StatusBadRequest)
    } else if !strings.HasPrefix(auth_header, "Bearer") {
      log.Printf("Authorization header is not a Bearer token")
      http.Error(w, "Bad request", http.StatusBadRequest)
    } else {
      id_token := strings.TrimPrefix(auth_header, "Bearer ")
      idToken, err := verifier.Verify(ctx, id_token)

      if err != nil {
        log.Printf("Failed to verify ID Token: %s", err.Error())
        http.Error(w, "Internal server error", http.StatusInternalServerError)
      } else {
        var claims struct {
          Exp int64 `json:"exp"`
          Aud []string `json:"aud"`
          Sub string `json:"sub"`
          AuthTime int `json:"auth_time"`
          SessionState string `json:"session_state"`
          Acr string `json:"acr"`
          RealmAccess struct {
            Roles []string `json:"roles"`
          } `json:"realm_access"`
          ResourceAccess struct {
            ExampleFrontend struct {
              Roles []string `json:"roles"`
            } `json:"example-frontend"`
          } `json:"resource_access"`
          Scope string `json:"scope"`
          EmailVerified bool `json:"email_verified"`
          Address struct {
          } `json:"address"`
          Name string `json:"name"`
          PreferredUsername string `json:"preferred_username"`
          GivenName string `json:"given_name"`
          FamilyName string `json:"family_name"`
          Email string `json:"email"`
        }

        if err := idToken.Claims(&claims); err != nil {
          log.Printf("Unable to parse claims from ID token: %s", err)
          http.Error(w, "Bad request", http.StatusBadRequest)
        } else {
          w.Write([]byte(fmt.Sprintf(
            "--------%s--------\r\n" +
            "Parsed claims from verified ID token (excerpt):\r\n" +
            " * exp = %s\r\n" +
            " * aud = %s\r\n" +
            " * email = %s\r\n" +
            " * name = %s\r\n" +
            " * preferred_username = %s\r\n" +
            " * resource_access.example-frontend.roles = %s\r\n",
            clientID,
            time.Unix(claims.Exp, 0).UTC(),
            claims.Aud,
            claims.Email,
            claims.Name,
            claims.PreferredUsername,
            claims.ResourceAccess.ExampleFrontend.Roles)))
        }
      }
    }
  })

  log.Printf("Listening on http://%s/", listenAddress)
  log.Fatal(http.ListenAndServe(listenAddress, nil))
}
