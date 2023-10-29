# Installation

This is tested under Ubuntu 20.04 LTS

We're assuming, that the ticker api will be available under the `api.domain.tld` domain. Please change accordingly.

__This should be considered a QUICK INSTALL GUIDE! Some best practices may differ.__

## Requirements

- `nginx`
- certificate (`certbot` & `python3-certbot-nginx` to use free Let's Encrypt Certs)
- `git`
- `go`
- Public IPv4
- Public IPv6 (Please!)

### Getting Go

_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_

Instead use:
 [golang.org install guide](https://golang.org/doc/install)

Please be also aware, that it's best practice to build your version of "ticker" not on the production machine. In order to keep the hurdle as low as possible, we will build the app on the system we're going to run it.
To enhance security maybe you want to remove `go` afterwards.

## Installation

### Build from source

_As mentioned above, this isn't best practice._
__You can also build it from source on your dedicated build server, your own pc at home, etc. Then just scp it over to the production Server afterwards.__

1. `cd /var/www/`
The directory where we gonna install all the things
2. `git clone https://github.com/systemli/ticker`
Clone the repository onto your disk
3. `cd ticker`
Go into the just cloned repository
4. `go build -o build/ticker`
Build the application
5. Go to "Configuration, Service and Stuff"

### Downloading a release from GitHub

1. Go to <https://github.com/systemli/ticker/releases>
4. Pick the latest release and download it via `wget https://github.com/systemli/ticker/releases/download/<version>/ticker-<version>-<architecture>`
5. `mv ticker-<version>-<architecture> /var/www/ticker/ticker`
6. `chmod +x /var/www/ticker/ticker`
7. Go to "Configuration, Service and Stuff"

### Configuration, Service and Stuff

1. `vim config.yml`
Fill your config file with the following content:

```yaml
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# configuration for the database
database:
    type: "sqlite" # postgres, mysql, sqlite
    dsn: "ticker.db" # postgres: "host=localhost port=5432 user=ticker dbname=ticker password=ticker sslmode=disable"
# secret used for JSON Web Tokens
secret: "<your special little secret> (make it LOOOONG!)"
# listen port for prometheus metrics exporter
metrics_listen: ":8181"
upload:
    # path where to store the uploaded files
    path: "uploads"
    # base url for uploaded assets
    url: "https://api.domain.tld"
```

2. Create a systemd Task (see [docs/ticker-api.service](assets/ticker-api.service) for reference)
2. `systemctl enable ticker-api.service`
3. `systemctl start ticker-api.service`
4. If you enter `systemctl status ticker-api.service` you'll see the generated admin password. __Please change it immediately!__
5. __Done. \o/__ You now have a fully functional ticker API.

## Exposing

In order to expose your ticker API to the users and not only yourself on the server, you'll need some sort of reverse proxy.
The following config expects you to use nginx, but apache2, caddy, etc. works just fine too.

`vim /etc/nginx/sites-available/ticker-api`

__This config is only for use with `cerbot`! Please create a secure SSL config if you won't let certbot do the job!__

```ticker-api.nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name api.domain.tld;
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://127.0.0.1:8080;
    }
}

```

_This is an example config for using TLS/SSL without certbot:_

```ticker-api.secure.nginx.conf
server {
    server_name api.domain.tld;
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://127.0.0.1:8080;
    }

    listen [::]:443 ssl ipv6only=on;
    listen 443 ssl;
    ssl_certificate /etc/ssl/api.domain.tld-fullchain.pem;
    ssl_certificate_key /etc/ssl/api.domain.tld-privkey.pem;
    ssl_session_cache shared:le_nginx_SSL:10m;
    ssl_session_timeout 1440m;
    ssl_session_tickets off;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;

    ssl_ciphers "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA";

}


server {
    return 301 https://$host$request_uri;
   
    listen 80;
    listen [::]:80;
    server_name api.domain.tld;
   
}

```

Create a symlink to enable this config:
`ln -s /etc/nginx/sites-available/ticker-api /etc/nginx/sites-enabled/`

Now run `nginx -t` to check if the config is correct.

If your output looks like this:

```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

then you can proceed. Otherwise: look for the error or ask someone to help.

Run `certbot --nginx --redirect -d api.domain.tld` to get a free SSL certificate. _Please keep in mind, that you need to point the `A` & `AAAA` Records to your machine!_

It should generate a certificate after answering a few questions like a email address, etc.

__done. This domain is now serving a ticker API. :)__
