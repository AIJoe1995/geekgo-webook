package tencent

import (
	"context"
	"fmt"
)

//https://cloud.tencent.com/document/product/382/43199

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = common.StringPtr(tpl)

	req.PhoneNumberSet = common.StringPtrs(numbers)
	req.TemplateParamSet = common.StringPtrs(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet { // SendStatusSet []*SendStatus `json:"SendStatusSet,omitnil" name:"SendStatusSet"`
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil

	//// 处理异常
	//if _, ok := err.(*errors.TencentCloudSDKError); ok {
	//	fmt.Printf("An API error has returned: %s", err)
	//	return
	//}
	//// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	//if err != nil {
	//	panic(err)
	//}
	//b, _ := json.Marshal(response.Response)
	//// 打印返回的json字符串
	//fmt.Printf("%s", b)
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    common.StringPtr(appId),
		signName: common.StringPtr(signName),
	}
}

//func main() {
//	/* 必要步骤：
//	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
//	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
//	 * 您也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
//	 * 以免泄露密钥对危及您的财产安全。
//	 * SecretId、SecretKey 查询: https://console.cloud.tencent.com/cam/capi */
//	credential := common.NewCredential(
//		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
//		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
//		"SecretId",
//		"SecretKey",
//	)
//}
