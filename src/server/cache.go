package server

import (
	"math/rand"
	"time"

	"github.com/kyokomi/expcache"
)

type ImageCache interface {
	GetRandomImageURL() string
}

type imageCache struct {
	Expire expcache.ExpireOnMemoryCache
	cache  []string

	twAPI     TwitterAPI
	rd        *rand.Rand
	twitterID string
	cacheCnt  int
}

func NewImageCache(twAPI TwitterAPI, twitterID string, cacheCnt int) ImageCache {
	e := &imageCache{
		cache:     []string{},
		twAPI:     twAPI,
		rd:        rand.New(rand.NewSource(time.Now().UnixNano())),
		twitterID: twitterID,
		cacheCnt:  cacheCnt,
	}
	e.Expire = expcache.NewExpireMemoryCache(e, 24*time.Hour)
	return e
}

func (c *imageCache) GetRandomImageURL() string {
	var result string
	c.Expire.WithRefreshLock(time.Now(), func() {
		idx := c.rd.Int31n(int32(len(c.cache))) - 1
		result = c.cache[idx]
	})
	return result
}

func (c *imageCache) Refresh() error {
	images, err := c.twAPI.GetFavoritesImages(c.cacheCnt, c.twitterID)
	if err != nil {
		return err
	}
	c.cache = images
	return nil
}

var _ expcache.OnMemoryCache = (*imageCache)(nil)
var _ ImageCache = (*imageCache)(nil)
