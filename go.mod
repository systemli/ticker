module github.com/systemli/ticker

require (
	github.com/appleboy/gin-jwt/v2 v2.9.0
	github.com/asdine/storm v2.1.2+incompatible
	github.com/dghubble/go-twitter v0.0.0-20211115160449-93a8679adecb
	github.com/dghubble/oauth1 v0.7.1
	github.com/disintegration/imaging v1.6.2
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/size v0.0.0-20220102055520-f75bacbc2df3
	github.com/gin-gonic/gin v1.8.1
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/uuid v1.3.0
	github.com/paulmach/go.geojson v1.4.0
	github.com/prometheus/client_golang v1.13.0
	github.com/sethvargo/go-password v0.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/afero v1.9.2
	github.com/spf13/viper v1.13.0
	github.com/stretchr/testify v1.8.0
	github.com/toorop/gin-logrus v0.0.0-20210225092905-2c785434f26f
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
)

require (
	github.com/DataDog/zstd v1.4.0 // indirect
	github.com/Sereal/Sereal v0.0.0-20190606082811-cf1bab6c7a3a // indirect
	github.com/goccy/go-json v0.9.7
	github.com/golang/snappy v0.0.4 // indirect
	github.com/onsi/ginkgo/v2 v2.1.6
	github.com/onsi/gomega v1.20.2
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
)

go 1.16

replace github.com/dghubble/go-twitter => github.com/0x46616c6b/go-twitter v0.0.0-media
