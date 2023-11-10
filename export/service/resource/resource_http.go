package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/idoubi/goz"
	ttsData "speech-tts/export/service/proto"
)

type Response struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Reason   string `json:"reason"`
	Metadata struct {
		DataList []*ttsData.GetTtsDataResponse_TTSData `json:"data_list"`
	} `json:"metadata"`
}

func GetTTSResByHttp(resType ttsData.ResType, languageType ttsData.LanguageType) ([]*ttsData.GetTtsDataResponse_TTSData, error) {

	url := fmt.Sprintf("http://%s/api/ttsData/v1/resource/get", HttpUrl)
	resp, err := goz.Get(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: struct {
			Resource     ttsData.ResType      `json:"resource"`
			LanguageType ttsData.LanguageType `json:"language"`
		}{
			Resource:     resType,
			LanguageType: languageType,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := resp.GetBody()
	if err != nil {
		return nil, err
	}
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Code != 200 {
		return nil, errors.New(result.Message)
	}
	return result.Metadata.DataList, nil
}

func RegisterResServiceByHttp(serviceName, callbackUrl string) error {
	url := fmt.Sprintf("http://%s/api/ttsData/v1/resource/register", HttpUrl)
	_, err := goz.Post(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: struct {
			ServiceName string `json:"service_name"`
			CallbackUrl string `json:"callback_url"`
		}{
			ServiceName: serviceName,
			CallbackUrl: callbackUrl,
		},
	})
	log.Infof("register res service by http, serviceName: %s, callbackUrl: %s", serviceName, callbackUrl)
	return err
}

func UnregisterResServiceByHttp(serviceName, callbackUrl string) error {
	url := fmt.Sprintf("http://%s/api/ttsData/v1/resource/unregister", HttpUrl)
	_, err := goz.Post(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: struct {
			ServiceName string `json:"service_name"`
			CallbackUrl string `json:"callback_url"`
		}{
			ServiceName: serviceName,
			CallbackUrl: callbackUrl,
		},
	})
	log.Infof("register res service by http, serviceName: %s, callbackUrl: %s", serviceName, callbackUrl)
	return err
}

func GetSpeakerModelByHttp() ([]*ttsData.GetSpeakerModelResult_SpeakerModel, error) {
	url := fmt.Sprintf("http://%s/api/ttsData/v1/resource/get-speaker-model", HttpUrl)
	resp, err := goz.Get(url, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}

	var respData struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		Reason   string `json:"reason"`
		Metadata struct {
			SpeakerModels []*ttsData.GetSpeakerModelResult_SpeakerModel `json:"speaker_models"`
		} `json:"metadata"`
	}
	body, _ := resp.GetBody()

	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}
	if respData.Code != 200 {
		return nil, errors.New(respData.Message)
	}
	return respData.Metadata.SpeakerModels, nil
}
