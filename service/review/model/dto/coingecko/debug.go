package dto_coingecko

import (
	"review-service/pkg/cache"
	"review-service/service/constant"
	"time"
)

type Debug struct {
}

// #################################################################

type ProductCategoryDebug struct {
	//Key
	CategoryName string
	//Value
	Url       string
	PageIndex uint8
	IsSuccess bool
}

func (dao *Debug) AddProductCategory(productCategoryDebug ProductCategoryDebug) {
	data, foundCacheRevainProductInfo := cache.LocalCache.Get(constant.KEY_CACHE_COINGECKO_PRODUCT_CATEGORY_INFO)
	if !foundCacheRevainProductInfo {
		data = make(map[string][]ProductCategoryDebug)
	}
	productInfoDebugList := data.(map[string][]ProductCategoryDebug)[productCategoryDebug.CategoryName]
	productInfoDebugList = append(productInfoDebugList, productCategoryDebug)
	data.(map[string][]ProductCategoryDebug)[productCategoryDebug.CategoryName] = productInfoDebugList
	time10days := 240 * time.Hour
	cache.LocalCache.SetByKey(constant.KEY_CACHE_COINGECKO_PRODUCT_CATEGORY_INFO, data, time10days)
}

func (dao *Debug) GetProductCategory() map[string][]ProductCategoryDebug {
	data, foundCacheRevainProductInfo := cache.LocalCache.Get(constant.KEY_CACHE_COINGECKO_PRODUCT_CATEGORY_INFO)
	if foundCacheRevainProductInfo {
		return data.(map[string][]ProductCategoryDebug)
	}
	return nil
}

// #################################################################
