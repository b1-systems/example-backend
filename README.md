# example-backend

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

```bash
sudo cp example-backend.ini.sample /usr/local/example-backend/example-backend.ini
sudo vi /usr/local/example-backend/example-backend.ini
```

Example `example-backend.ini`:

```
[example-backend]
# These values are provided by your IdP for your confidential client
clientID = example-backend
clientSecret = secret

# This URL will be used for endpoint discovery of your IdP
providerUrl = https://example.test/keycloak/realms/golang-oidc

# Plain HTTP service address of this "example-frontend" server:
listenAddress = 0.0.0.0:8080
```

# Start

```bash
systemctl start example-backend.service
journalctl -xefu  example-backend.service
```

