package common

import (
	"Lumino/common/http_error_code"
	"Lumino/common/logger"
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/spf13/viper"
	"mime/multipart"
	"time"
)

// OssClient -
type OssClient struct {
	Client *oss.Client
}

// NewOssClient 新建oss对象
func NewOssClient() *OssClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			viper.GetString("oss.accessKey"), viper.GetString("oss.secretKey"))).
		WithRegion(viper.GetString("oss.region"))

	// 创建OSS客户端
	client := oss.NewClient(cfg)
	return &OssClient{Client: client}
}

// UploadFile 上传oss文件
func (o *OssClient) UploadFile(name string, file multipart.File) error {
	putRequest := &oss.PutObjectRequest{
		Bucket: oss.Ptr(viper.GetString("oss.bucket")), // 存储空间名称
		Key:    oss.Ptr(name),                          // 对象名称
		Body:   file,
	}
	result, err := o.Client.PutObject(context.TODO(), putRequest)
	logger.Info(result)
	if err != nil {
		return http_error_code.Internal("oss上传异常",
			http_error_code.WithInternal(err))
	}
	return nil
}

// DownloadFile 下载oss文件
func (o *OssClient) DownloadFile(name string) (string, error) {
	// 生成GetObject的预签名URL
	result, err := o.Client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr(viper.GetString("oss.bucket")), // 存储空间名称
		Key:    oss.Ptr(name),                          // 对象名称
	},
		oss.PresignExpires(10*time.Minute),
	)
	if err != nil {
		return "", http_error_code.Internal("oss下载异常",
			http_error_code.WithInternal(err))
	}
	return result.URL, err
}
