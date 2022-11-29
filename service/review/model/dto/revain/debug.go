package dto_revain

import (
	"review-service/pkg/cache"
	"review-service/service/constant"
	"time"
)

type Debug struct {
}

// #################################################################
type ProductInfoDebug struct {
	//Key
	EndpointProduct string
	//Value
	Url       string
	PageIndex uint8
	IsSuccess bool
}

func (dao *Debug) AddProductInfo(productInfoDebug ProductInfoDebug) {
	data, foundCacheRevainProductInfo := cache.LocalCache.Get(constant.KEY_CACHE_REVAIN_PRODUCT_INFO)
	if !foundCacheRevainProductInfo {
		data = make(map[string][]ProductInfoDebug)
	}
	productInfoDebugList := data.(map[string][]ProductInfoDebug)[productInfoDebug.EndpointProduct]
	productInfoDebugList = append(productInfoDebugList, productInfoDebug)
	data.(map[string][]ProductInfoDebug)[productInfoDebug.EndpointProduct] = productInfoDebugList
	time10days := 240 * time.Hour
	cache.LocalCache.SetByKey(constant.KEY_CACHE_REVAIN_PRODUCT_INFO, data, time10days)
}

func (dao *Debug) GetProductInfo() map[string][]ProductInfoDebug {
	data, foundCacheRevainProductInfo := cache.LocalCache.Get(constant.KEY_CACHE_REVAIN_PRODUCT_INFO)
	if foundCacheRevainProductInfo {
		return data.(map[string][]ProductInfoDebug)
	}
	return nil
}

// #################################################################

// #################################################################
type ProductDetailDebug struct {
}

// #################################################################
