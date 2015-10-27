package server

import (
	"net/http"

	"os"
	"strconv"

	"appengine"
)

func init() {
	var kohaImageCache ImageCache

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "example", "twitter-favorites-api")
	})

	http.HandleFunc("/v1/api/image", func(w http.ResponseWriter, r *http.Request) {
		// appengine.NewContext(r) を事前に用意できないのでnilだったら生成する感じにしてる
		if kohaImageCache == nil {
			twAPI, err := NewTwitterAPI(NewAppEngineHttpClient(appengine.NewContext(r)), "", "")
			if err != nil {
				renderErrorResponse(w, 400, "twitter auth error "+err.Error())
				return
			}

			twitterID := os.Getenv("TWITTER_ID")
			cacheCount, err := strconv.Atoi(os.Getenv("CACHE_COUNT"))
			if err != nil {
				cacheCount = 10 // TODO: 未設定なら適当に10件くらい
			}
			kohaImageCache = NewImageCache(twAPI, twitterID, cacheCount)
		}

		imageURL := kohaImageCache.GetRandomImageURL()
		renderer.Text(w, http.StatusOK, imageURL)
	})
}
