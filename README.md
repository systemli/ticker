# ticker

## Setup 

  * Clone Repo to your $GOPATH (e.g. `/home/$USER/go/src/git.codecoop.org/systemli/ticker`) 
  * install go (e.g. `sudo apt install golang-go`)
  * [optional] adjust config.yml.dist
  * switch to ticker directory (e.g. `cd /home/$USER/go/src/git.codecoop.org/systemli/ticker`)
  * run `go run main.go -config config.yml.dist`

  If everything works correct, you should see the following output:
```
user@laptop:ticker: $ go run main.go -config config.yml.dist                      
INFO[0000] admin user created (change password now!)     email=admin@systemli.org password="ApasswordString"
INFO[0000] starting ticker at localhost:8080            
```

## Development

```
go run main.go -config config.yml.dist
```

## Configuration

  * Example config.yml.dist

```
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# initiator is the email for the first admin user (see password in logs)
initiator: "admin@systemli.org"
# database is the path to the bolt file
database: "ticker.db"
# secret used for JSON Web Tokens
secret: "slorp-panfil-becall-dorp-hashab-incus-biter-lyra-pelage-sarraf-drunk"
```

## Testing

```
go test ./... -cover
```
