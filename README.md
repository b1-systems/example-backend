# example-backend

## Overview

Demonstration of passing an OpenID Connect ID Token to a web application via Authorization header.

*Important Notes:*

1. This example does not endorse the practise of passing the ID token as an Authorization header but only demonstrates that it works.
2. This example assumes that "example-frontend" communicates HTTP requests to "example-backend" in the same security domain, i.e. "example-backend" is trusted to
   - validate the ID token against the public keys of the IdP and
   - adhere to the audience claim in the ID token.

## Requirements

"example-frontend" should be up and running.
* See also <https://tk-sls.de/gitlab/golang-oidc/example-frontend/>.
* *Note:* Keycloak client "example-frontend" should exist; this will become relevant in section "Configuration", step 2 (see below).

## Installation

```bash
git clone https://tk-sls.de/gitlab/golang-oidc/example-backend.git
cd example-backend
go mod tidy
go build
sudo mkdir /usr/local/example-backend
sudo cp example-backend /usr/local/example-backend
sudo cp example-backend.service /etc/systemd/system
sudo systemctl daemon-reload
```

## Configuration

1. In Keycloak, create a client "example-backend".
   * Set Access Type to "Bearer only".

2. In Keycloak, extend the audience of the ID token of client "example-frontend".
   * Go to Clients -> "example-frontend" -> tab Mappers
   * Create a mapper
     - Name (for example): `aud-add-example-backend`
     - Mapper Type: Audience
     - Included Custom Audience: "example-backend"
     - Add to ID token: ON

3. Create configuration file "example-backend.ini"

```bash
sudo cp example-backend.ini.sample /usr/local/example-backend/example-backend.ini
sudo vi /usr/local/example-backend/example-backend.ini
```

Example `example-backend.ini`:

```
[example-backend]
# Client ID as set in Keycloak:
clientID = example-backend

# This URL will be used for endpoint discovery of your IdP:
providerUrl = https://your_idp_server/realms/golang-oidc

# Plain HTTP service address of this "example-frontend" server:
listenAddress = 0.0.0.0:8080
```

# Start

```bash
systemctl start example-backend.service
journalctl -xefu  example-backend.service
```

