package data

import (
	"context"
	"github.com/spf13/viper"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
)

type SpeakerSetting struct {
	SupportedSpeaker []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ChineseName string `mapstructure:"chinese_name"`
	} `json:"SupportedSpeaker"`
	SupportedSpeed  []string `json:"SupportedSpeed"`
	SupportedVolume []string `json:"SupportedVolume"`
	SupportedPitch  []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ChineseName string `mapstructure:"chinese_name"`
	} `json:"SupportedPitch"`
	SupportedEmotion []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ChineseName string `mapstructure:"chinese_name"`
	} `json:"SupportedEmotion"`
	SupportedDigitalPerson []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ChineseName string `mapstructure:"chinese_name"`
	} `json:"SupportedDigitalPerson"`
}

type SpeakerInfo struct {
	SpeakerId            int
	SpeakerName          string
	ParameterSpeakerName string
	IsSupportEmotion     bool
}

type HandlerObject struct {
	context.Context
	SpeakerInfo
	ParamMap      map[string]interface{}
	IsInterrupted bool
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
	viper.SetConfigName("speaker")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.Unmarshal(&speakerSetting)
	return &speakerSetting, nil
}
