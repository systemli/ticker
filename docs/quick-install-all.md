# Full Installation

This was tested with an Ubuntu 20.04 LTS server.

Results may differ on other releases or distributions!

## Requirements

- `git`
- `go`
- `nodejs`
- `yarn`
- `nginx`
- certificate (`certbot` & `python3-certbot-nginx` to use free Let's Encrypt Certs)
- `git`
- Public IPv4
- Public IPv6 (Please!)

### Getting Go

_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_

Instead use:
 [golang.org install guide](https://golang.org/doc/install)

Please be also aware, that it's best practice to build your version of "ticker" not on the production machine. In order to keep the hurdle as low as possible, we will build the app on the system we're going to run it.
To enhance security maybe you want to remove `go` afterwards.

### Getting NodeJS

_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_

Instead use:
 [nodesource/distributions](https://github.com/nodesource/distributions/blob/master/README.md#debinstall)

## Install ticker

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

`vim /etc/nginx/sites-available/ticker-api`

```nginx.conf
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

__If you don't want to use certbot for your installation, please keep in mind, that nontheless you'll need an TLS cert for running this in production and need to configure the nginx accordingly by yourself!__

It should generate a certificate after answering a few questions like a email address, etc.

__done. This domain is now serving a ticker API. :)__

## Install ticker-admin

1. `cd /var/www/`
The directory where we gonna install all the things
2. `git clone https://github.com/systemli/ticker-admin`
Clone the repository onto your disk
2. `cd ticker-admin`
Go into the just cloned repository
1. `yarn`
Install the dependencies
5. `vim .env`
Fill your .env file with the following content:

```.env
REACT_APP_API_URL=https://api.domain.tld/v1
```

_Change `api.domain.tld` to the URL you chose at ticker API server creation_

7. `yarn build`
Build the application
6. `chown www-data:www-data dist/ -R`
Sets the owner for the freshly created dist repository to your nginx user

### Exposing

`vim /etc/nginx/sites-available/ticker-admin`

```nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name admin.domain.tld;
    root /var/www/ticker-admin/dist;
    index index.html;
    location / {
        try_files $uri $uri/ =404;
    }
}

```

Create a symlink to enable this config:
`ln -s /etc/nginx/sites-available/ticker-admin /etc/nginx/sites-enabled/`

Now run `nginx -t` to check if the config is correct.

If your output looks like this:

```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

then you can proceed. Otherwise: look for the error or ask someone to help.

Run `certbot --nginx --redirect -d admin.domain.tld` to get a free SSL certificate. _Please keep in mind, that you need to point the `A` & `AAAA` Records to your machine!_

__If you don't want to use certbot for your installation, please keep in mind, that nontheless you'll need an TLS cert for running this in production and need to configure the nginx accordingly by yourself!__

It should generate a certificate after answering a few questions like a email address, etc.

done. This domain is now serving a ticker frontend. :)

**You need to create the ticker in ticker-admin in order to see something!**

## Install ticker-frontend

1. `cd /var/www/`
The directory where we gonna install all the things
2. `git clone https://github.com/systemli/ticker-frontend`
Clone the repository onto your disk
2. `cd ticker-frontend`
Go into the just cloned repository
3. `git checkout d03982f3059d6335a9e9ec0abcb71813ccbafef7`
Checkout this branch. _Hopefully, this won't be necessary in the future, but right now, this seems to be the the last working commit (for me)_
4. `yarn`
Install the dependencies
5. `vim .env`
Fill your .env file with the following content:

```.env
REACT_APP_API_URL=https://api.domain.tld/v1
```

_Change `api.domain.tld` to the URL you chose at ticker API server creation_

7. `yarn build`
Build the application
6. `chown www-data:www-data dist/ -R`
Sets the owner for the freshly created dist repository to your nginx user

### Exposing

`vim /etc/nginx/sites-available/ticker-frontend`

_The following config is for a single domain only! For wildcard configs, just replace the `sub.domain.tld` with a `*.domain.tld`. You need to validate your domain via DNS challenge or use another provider then Let's Encrypt!_

```nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name sub.domain.tld;
    root /var/www/ticker-frontend/dist;
    index index.html;
    location / {
        try_files $uri $uri/ =404;
    }
}

```

Create a symlink to enable this config:
`ln -s /etc/nginx/sites-available/ticker-frontend /etc/nginx/sites-enabled/`

Now run `nginx -t` to check if the config is correct.

If your output looks like this:

```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

then you can proceed. Otherwise: look for the error or ask someone to help.

Run `certbot --nginx --redirect -d sub.domain.tld` to get a free SSL certificate. _Please keep in mind, that you need to point the `A` & `AAAA` Records to your machine!_

__If you don't want to use certbot for your installation, please keep in mind, that nontheless you'll need an TLS cert for running this in production and need to configure the nginx accordingly by yourself!__

It should generate a certificate after answering a few questions like a email address, etc.

done. This domain is now serving a ticker frontend. :)

**You need to create the ticker in ticker-admin in order to see something!**

## First touch

- Go to `https://admin.domain.tld` and log in with the credentials provided at your first start of the ticker api.
- Change the provided credentials (and use a password manager)
- Create a new ticker for with the domain `sub.domain.tld`
- Create test content
- Try to open `https://sub.domain.tld`
- ...
- profit?
