package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	u, err := url.Parse("https://coolcar-1307431695.cos.ap-beijing.myqcloud.com")
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	secID := "AKIDVpseIflbCYeT2KL8gxxg8KFtHPGq9CyB"
	secKey := "0kFSO3SaQ1OUKtBE1bvzmJZaJlxsXv1q"
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secID,
			SecretKey: secKey,
		},
	})
	// 获取预签名URL
	name := "abc.png"
	presignedURL, err := client.Object.GetPresignedURL(context.Background(),
		http.MethodPut, name, secID, secKey, 1*time.Hour, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(presignedURL)
}
