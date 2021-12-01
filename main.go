package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config := GetConfig()

	r := gin.Default()

	for _, router := range config.Routers {
		r.GET(router.Route, func(c *gin.Context) {
			for _, hook := range router.Hooks {
				HandleHook(parseConfig(c, hook))
			}
		})
	}
	if err := r.Run(config.Listen); err != nil {
		log.Fatalln(err.Error())
	}
}

func parseConfig(c *gin.Context, config HookConfig) HookConfig {
	parseString := makeParseStringFunc(c)

	config.Command = parseString(config.Command)
	config.URL = parseString(config.URL)
	config.Freq = parseString(config.Freq)
	return config
}

var placeholder = regexp.MustCompile(`#{\w+}`)

func makeParseStringFunc(c *gin.Context) func(src string) string {
	return func(src string) string {
		return placeholder.ReplaceAllStringFunc(src, func(s string) string {
			key := s[2 : len(s)-1]
			if v := strings.TrimLeft(c.Param(key), "/"); v != "" {
				return v
			}
			if v := c.Query(key); v != "" {
				return v
			}
			return ""
		})
	}
}

func HandleHook(hook HookConfig) {
	switch {
	case hook.URL != "":
		if !shouldLimit(hook.URL, hook.Freq) {
			go http.Get(hook.URL)
		}
	case hook.Command != "":
		if !shouldLimit(hook.Command, hook.Freq) {
			cmd := exec.Command("sh", "-c", hook.Command)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Start()
		}
	}
}

func shouldLimit(key, freqStr string) bool {
	if freqStr == "" {
		return false
	}
	freq, err := time.ParseDuration(freqStr)
	if err != nil {
		return false
	}
	encodedKey := "noti_" + base64.StdEncoding.EncodeToString([]byte(key))
	res, err := rdb.SetNX(context.Background(), encodedKey, 1, freq).Result()
	fmt.Println(encodedKey, freq, res)
	if err != nil {
		return false
	}
	return !res
}
