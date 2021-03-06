module github.com/systemli/ticker

require (
	github.com/DataDog/zstd v1.4.0 // indirect
	github.com/Sereal/Sereal v0.0.0-20190606082811-cf1bab6c7a3a // indirect
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/appleboy/gofight/v2 v2.1.2
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.11.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/dghubble/go-twitter v0.0.0-20200725221434-4bc8ad7ad1b4
	github.com/dghubble/oauth1 v0.6.0
	github.com/disintegration/imaging v1.6.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/size v0.0.0-20200815104238-dc717522c4e2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/pat v0.0.0-20180118222023-199c85a7f6d1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/labstack/gommon v0.2.9 // indirect
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/paulmach/go.geojson v1.4.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.13.0 // indirect
	github.com/sethvargo/go-password v0.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.3.4
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1
	github.com/toorop/gin-logrus v0.0.0-20190701131413-6c374ad36b67
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/image v0.0.0-20200801110659-972c09e46d76 // indirect
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/ini.v1 v1.60.2 // indirect
)

go 1.15

replace github.com/dghubble/go-twitter => ./forks/go-twitter
