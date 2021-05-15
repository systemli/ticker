Installation of ticker, ticker-admin & ticker-frontend
==
This was tested with an Ubuntu 20.04 LTS server.

Results may differ on other releases or distributions!

# Requirements

- `git`
- `go`
- `make`
- `nodejs`
- `yarn`
- `nginx`
- certificate (`certbot` & `python3-certbot-nginx` to use free Let's Encrypt Certs)
- `git`
- Public IPv4
- Public IPv6 (Please!)

## Getting Go:
_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_ 

Instead use:
 [golang.org install guide](https://golang.org/doc/install)

## Getting NodeJS:
_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_ 

Instead use:
 [nodesource/distributions](https://github.com/nodesource/distributions/blob/master/README.md#debinstall)

# Install ticker

1. `cd /var/www/`
The directory where we gonna install all the things
1. `git clone --recurse-submodules https://github.com/systemli/ticker`
Clone the repository onto your disk
1. `cd ticker`
Go into the just cloned repository
3. `git submodule update --init --recursive`
Initialise all the submodules
1. `vim config.yml.dist`
Fill your config file with the following content:
```yaml
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# initiator is the email for the first admin user (see password in logs)
initiator: "<your e-mail>"
# database is the path to the bolt file
database: "ticker.db"
# secret used for JSON Web Tokens
secret: "<your special little secret> (make it LOOOONG!)"
# twitter configuration
twitter_consumer_key: ""
twitter_consumer_secret: ""
# listen port for prometheus metrics exporter
metrics_listen: ":8181"
# path where to store the uploaded files
upload_path: "uploads"
# base url for uploaded assets
upload_url: "https://api.domain.tld"
```
6. `make run`
Run the application

__You'll receive credentials in the console output. Use them for your first login to ticker-admin and change them immediately!__


## Exposing
`vim /etc/nginx/sites-available/ticker-api`


```nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name api.domain.tld
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

# Install ticker-admin

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

## Exposing
`vim /etc/nginx/sites-available/ticker-admin`


```nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name admin.domain.tld
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


# Install ticker-frontend

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

## Exposing
`vim /etc/nginx/sites-available/ticker-frontend`

_The following config is for a single domain only! For wildcard configs, just replace the `sub.domain.tld` with a `*.domain.tld`. You need to validate your domain via DNS challenge or use another provider then Let's Encrypt!_

```nginx.conf
server {
    listen 80;
    listen [::]:80;
    server_name sub.domain.tld
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

# First touch

- Go to `https://admin.domain.tld` and log in with the credentials provided at your first start of the ticker api.
- Change the provided credentials (and use a password manager)
- Create a new ticker for with the domain `sub.domain.tld`
- Create test content
- Try to open `https://sub.domain.tld`
- ...
- profit?