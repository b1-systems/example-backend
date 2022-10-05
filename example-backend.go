package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "time"
  "github.com/coreos/go-oidc/v3/oidc"
  "golang.org/x/net/context"
  "gopkg.in/ini.v1"
)

var (
  clientName = "example-backend"
  clientID = ""
  clientSecret = ""
  providerUrl = ""
  listenAddress = ""
)

func readIni() {
  ex, err := os.Executable()

  if err != nil {
    panic(err)
  }

  cfg, err := ini.Load(filepath.Join(filepath.Dir(ex), clientName + ".ini"))

  if err != nil {
    panic(err)
  }

  cs := cfg.Section(clientName)

  clientID = cs.Key("clientID").String()

  if clientID == "" {
    log.Fatal(clientName + ".ini does not specify clientID")
    os.Exit(1)
  }

  clientSecret = cs.Key("clientSecret").String()

  if clientSecret == "" {
    log.Fatal(clientName + ".ini does not specify clientSecret")
    os.Exit(1)
  }

  providerUrl = cs.Key("providerUrl").String()

  if providerUrl == "" {
    log.Fatal(clientName + ".ini does not specify providerUrl")
    os.Exit(1)
  }

  listenAddress = cs.Key("listenAddress").String()

  if listenAddress == "" {
    log.Fatal(clientName + ".ini does not specify listenAddress")
    os.Exit(1)
  }

  log.Printf(
    "Read configuration:\n" +
    " clientID = %s\n" +
    " clientSecret = %s\n" +
    " providerUrl = %s\n" +
    " listenAddress = %s\n",
    clientID,
    "*REDACTED*",
    providerUrl,
    listenAddress,
  )
}

func main() {
  readIni()

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
            "Client %s parsed claims from verified ID token: exp = %s, aud = %s, email = %s, email_verified = %t",
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
