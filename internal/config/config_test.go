package config

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	log.Logger.SetOutput(GinkgoWriter)

	var envs map[string]string = map[string]string{
		"TICKER_LISTEN":         ":7070",
		"TICKER_LOG_LEVEL":      "trace",
		"TICKER_LOG_FORMAT":     "text",
		"TICKER_SECRET":         "secret",
		"TICKER_DATABASE_TYPE":  "mysql",
		"TICKER_DATABASE_DSN":   "user:password@tcp(localhost:3306)/ticker?charset=utf8mb4&parseTime=True&loc=Local",
		"TICKER_METRICS_LISTEN": ":9191",
		"TICKER_UPLOAD_PATH":    "/data/uploads",
		"TICKER_UPLOAD_URL":     "https://example.com",
		"TICKER_TELEGRAM_TOKEN": "token",
	}

	Describe("LoadConfig", func() {
		BeforeEach(func() {
			for key := range envs {
				os.Unsetenv(key)
			}
		})

		Context("when path is empty", func() {
			It("loads config with default values", func() {
				c := LoadConfig("")
				Expect(c.Listen).To(Equal(":8080"))
				Expect(c.LogLevel).To(Equal("debug"))
				Expect(c.LogFormat).To(Equal("json"))
				Expect(c.Secret).ToNot(BeEmpty())
				Expect(c.Database.Type).To(Equal("sqlite"))
				Expect(c.Database.DSN).To(Equal("ticker.db"))
				Expect(c.MetricsListen).To(Equal(":8181"))
				Expect(c.UploadPath).To(Equal("uploads"))
				Expect(c.UploadURL).To(Equal("http://localhost:8080"))
				Expect(c.Telegram.Token).To(BeEmpty())
				Expect(c.Telegram.Enabled()).To(BeFalse())
			})

			It("loads config from env", func() {
				for key, value := range envs {
					err := os.Setenv(key, value)
					Expect(err).ToNot(HaveOccurred())
				}

				c := LoadConfig("")
				Expect(c.Listen).To(Equal(envs["TICKER_LISTEN"]))
				Expect(c.LogLevel).To(Equal(envs["TICKER_LOG_LEVEL"]))
				Expect(c.LogFormat).To(Equal(envs["TICKER_LOG_FORMAT"]))
				Expect(c.Secret).To(Equal(envs["TICKER_SECRET"]))
				Expect(c.Database.Type).To(Equal(envs["TICKER_DATABASE_TYPE"]))
				Expect(c.Database.DSN).To(Equal(envs["TICKER_DATABASE_DSN"]))
				Expect(c.MetricsListen).To(Equal(envs["TICKER_METRICS_LISTEN"]))
				Expect(c.UploadPath).To(Equal(envs["TICKER_UPLOAD_PATH"]))
				Expect(c.UploadURL).To(Equal(envs["TICKER_UPLOAD_URL"]))
				Expect(c.Telegram.Token).To(Equal(envs["TICKER_TELEGRAM_TOKEN"]))
				Expect(c.Telegram.Enabled()).To(BeTrue())
			})
		})

		Context("when path is not empty", func() {
			Context("when path is absolute", func() {
				It("loads config from file", func() {
					path, err := filepath.Abs("../../testdata/config_valid.yml")
					Expect(err).ToNot(HaveOccurred())
					c := LoadConfig(path)
					Expect(c.Listen).To(Equal("127.0.0.1:8888"))
				})
			})

			Context("when path is relative", func() {
				It("loads config from file", func() {
					c := LoadConfig("../../testdata/config_valid.yml")
					Expect(c.Listen).To(Equal("127.0.0.1:8888"))
				})
			})

			Context("when file does not exist", func() {
				It("loads config with default values", func() {
					c := LoadConfig("config_notfound.yml")
					Expect(c.Listen).To(Equal(":8080"))
				})
			})

			Context("when file is invalid", func() {
				It("loads config with default values", func() {
					c := LoadConfig("../../testdata/config_invalid.txt")
					Expect(c.Listen).To(Equal(":8080"))
				})
			})
		})
	})
})
