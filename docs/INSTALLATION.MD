# Installation
This is tested under Ubuntu 20.04 LTS

We're assuming, that the ticker api will be available under the `api.domain.tld` domain. Please change accordingly.

## Requirements:
- `nginx`
- certificate (`certbot` & `python3-certbot-nginx` to use free Let's Encrypt Certs)
- `git`
- Public IPv4
- Public IPv6 (Please!)
- `make`


### Getting Go:
_Don't use the shipped version of your system, if you're working on a Debian based OS (Ubuntu, etc)_ 

Instead use:
 [golang.org install guide](https://golang.org/doc/install)

## Installation

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

__Done__

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

It should generate a certificate after answering a few questions like a email address, etc.

__done. This domain is now serving a ticker API. :)__
