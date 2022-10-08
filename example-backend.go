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
          Expires int64 `json:"exp"`
          Audience []string `json:"aud"`
          Email string `json:"email"`
          Verified bool `json:"email_verified"`
        }

        if err := idToken.Claims(&claims); err != nil {
          log.Printf("Unable to parse claims from ID token: %s", err)
          http.Error(w, "Bad request", http.StatusBadRequest)
        } else {
          w.Write([]byte(fmt.Sprintf(
            "--------%s--------\r\n" +
            "Parsed claims from verified ID token: exp = %s, aud = %s, email = %s, email_verified = %t\r\n",
            clientID,
            time.Unix(claims.Expires, 0).UTC(),
            claims.Audience,
            claims.Email,
            claims.Verified)))
        }
      }
    }
  })

  log.Printf("Listening on http://%s/", listenAddress)
  log.Fatal(http.ListenAndServe(listenAddress, nil))
}
