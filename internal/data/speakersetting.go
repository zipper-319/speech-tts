package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/viper"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
)

type SpeakerSetting struct {
	SupportedSpeaker []struct {
		Name        string `json:"name"`
		ChineseName string `json:"chinese_name"`
	} `json:"SupportedSpeaker"`
	SupportedSpeed  []string `json:"SupportedSpeed"`
	SupportedVolume []string `json:"SupportedVolume"`
	SupportedPitch  []struct {
		Name        string `json:"name"`
		ChineseName string `json:"chinese_name"`
	} `json:"SupportedPitch"`
	SupportedEmotion []struct {
		Name        string `json:"name"`
		ChineseName string `json:"chinese_name"`
	} `json:"SupportedEmotion"`
	SupportedMovementDescriptor []struct {
		Name        string `json:"name"`
		ChineseName string `json:"chinese_name"`
	} `json:"SupportedMovementDescriptor"`
	SupportedExpressionDescriptor []struct {
		Name        string `json:"name"`
		ChineseName string `json:"chinese_name"`
	} `json:"SupportedExpressionDescriptor"`
}

type SpeakerInfo struct {
	SpeakerName          string
	ParameterSpeakerName string
	IsSupportEmotion     bool
	IsSupportMixedVoice  bool
}

type HandlerObject struct {
	context.Context
	SpeakerInfo
	ParamMap      map[string]interface{}
	IsInterrupted bool
	Log           *log.Helper
}

type HandlerObjectV2 struct {
	HandlerObject
	BackChan chan v2.TtsRes
}

type HandlerObjectV1 struct {
	HandlerObject
	BackChan chan v1.TtsRes
}

func NewSpeakerSetting(path string) (*SpeakerSetting, error) {
	var speakerSetting SpeakerSetting
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.Unmarshal(&speakerSetting)
	return &speakerSetting, nil
}
