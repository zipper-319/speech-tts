package service

/*
#cgo CFLAGS: -I../include -DCOMPATIBLE_V1=1
#cgo LDFLAGS: -L../libs -lCmTts
#include <stdio.h>
#include <stdlib.h>
#include "ActionSynthesizer.h"
#include "VitsDef.h"
#include "TtsSetting.h"
#include "AnimationDef.h"


extern void goOnStart(void*, const char*, FacialExpressionConfig*, BodyMovementConfig* movementConfig);
extern void goOnSynthesizedData(void*, SynthesizedAudio*, Coordinate*);
extern void goOnTimedMouthShape(void*, TimedMouthShape*, int, float);
extern void goOnFacialExpression(void*, FacialExpression*);
extern void goOnBodyMovement(void*, BodyMovementSegment*);
extern void goOnEnd(void*, int flags);
extern void goOnDebug(void*, const char *, const char *);
extern void goOnActionElement(void*, int, const char*, int, Coordinate*, int);

typedef void (*typOnStart)(void* pUserData, const char* ttsText, FacialExpressionConfig* expressionConfig, BodyMovementConfig* movementConfig);
typedef void (*typOnSynthesizedData)(void* pUserData, SynthesizedAudio* data, Coordinate* coordinate);
typedef void (*typOnDebug)(void* pUserData, const char* type, const char *info);
typedef void (*typOnTimedMouthShape)(void* pUserData, TimedMouthShape* mouth, int size, float startTimeMs);
typedef void (*typOnFacialExpression)(void* pUserData, FacialExpression* expression);
typedef void (*typOnBodyMovement)(void* pUserData, BodyMovementSegment* movement);
typedef void (*typOnActionElement)(void* pUserData, int type, const char* url, int operation_type, Coordinate* coordinate, int render_duration);
typedef void (*typOnEnd)(void* pUserData, int flags);

extern void goOnStartV1(void*);
extern void goOnAudioV1(void*, char*,int);
extern void goOnEndV1(void*);
extern void goOnDebugV1(void*, const char *);
extern void goOnTimedMouthShapeV1(void*, TimedMouthShape*,int, const char*);
extern void goOnFacialExpressionV1(void*, FacialExpression*);
extern void goOnCurTextSegmentV1(void*, const char*, const char*);

typedef void (*typOnStartV1)(void* pUserData);
typedef void (*typOnAudioV1)(void* pUserData, char *data, int size);
typedef void (*typOnEndV1)(void* pUserData);
typedef void (*typOnDebugV1)(void* pUserData, const char *info);
typedef void (*typOnTimedMouthShapeV1)(void* pUserData, TimedMouthShape* mouth, int size, const char* subText);
typedef void (*typOnCurTextSegmentV1)(void* pUserData, const char* normalizedText, const char* originalText);
typedef void (*typOnFacialExpressionV1)(void* pUserData, FacialExpression* expression);

char** speakerList = NULL;
*/
import "C"
import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/conf"
	"speech-tts/internal/data"
	"speech-tts/internal/utils"
	"strings"
	"time"
	"unsafe"
)

var ProviderSet = wire.NewSet(NewTTSService, NewNotify)

var (
	actionCallback = C.ActionCallback{}
	ttsCallback    = C.TTS_Callback{}
)

func init() {

	actionCallback.onStart = C.typOnStart(C.goOnStart)
	actionCallback.onSynthesizedData = C.typOnSynthesizedData(C.goOnSynthesizedData)
	actionCallback.onEnd = C.typOnEnd(C.goOnEnd)
	actionCallback.onDebug = C.typOnDebug(C.goOnDebug)
	actionCallback.onTimedMouthShape = C.typOnTimedMouthShape(C.goOnTimedMouthShape)
	actionCallback.onFacialExpression = C.typOnFacialExpression(C.goOnFacialExpression)
	actionCallback.onBodyMovement = C.typOnBodyMovement(C.goOnBodyMovement)
	actionCallback.onActionElement = C.typOnActionElement(C.goOnActionElement)

	ttsCallback.onStart = C.typOnStartV1(C.goOnStartV1)
	ttsCallback.onAudio = C.typOnAudioV1(C.goOnAudioV1)
	ttsCallback.onEnd = C.typOnEndV1(C.goOnEndV1)
	ttsCallback.onDebug = C.typOnDebugV1(C.goOnDebugV1)
	ttsCallback.onTimedMouthShape = C.typOnTimedMouthShapeV1(C.goOnTimedMouthShapeV1)
	ttsCallback.onCurTextSegment = C.typOnCurTextSegmentV1(C.goOnCurTextSegmentV1)
	ttsCallback.onFacialExpression = C.typOnFacialExpressionV1(C.goOnFacialExpressionV1)

}

type TTSService struct {
	ChanLen           int
	ResPath           string
	version           string
	resServiceVersion string
	Speakers          []*data.SpeakerInfo
	*data.SpeakerSetting
	log.Logger
}

func NewNotify(logger log.Logger) chan interface{} {
	notify := make(chan interface{})
	go func() {
		select {
		case <-notify:
			log.NewHelper(logger).Info("success to init tts service")
		case <-time.After(5 * time.Minute):
			log.NewHelper(logger).Error("fail_to_init_tts_service after 5 minute")
		}
	}()
	return notify
}

func NewTTSService(s *conf.Server, speakerSetting *data.SpeakerSetting, logger log.Logger, notify chan interface{}) *TTSService {
	cResPath := C.CString(s.App.Path)
	defer C.free(unsafe.Pointer(cResPath))
	C.ActionSynthesizer_Init(cResPath)
	version := C.ActionSynthesizer_GetVersion()
	resServiceVersion := C.ActionSynthesizer_GetResServiceVersion()
	// 发音人初始化
	speakers := make([]*data.SpeakerInfo, len(speakerSetting.SupportedSpeaker))

	for i, supportedSpeaker := range speakerSetting.SupportedSpeaker {
		cname := C.CString(supportedSpeaker.Name)
		userSpace := C.CString(utils.DefaultUser)
		m := C.GetSpeakerDescriptor(userSpace, cname)
		isEmotion := m.flags&C.SUPPORT_EMOTION != 0
		isMixedVoice := m.flags&C.SUPPORT_MIXED_VOICE != 0
		log.NewHelper(logger).Infof("----ChineseName:%s------cname:%s------flags:%d, isEmotion:%t, isMixedVoice:%t, language:%s",
			supportedSpeaker.ChineseName, supportedSpeaker.Name, int(m.flags), isEmotion, isMixedVoice, C.GoString(m.language))
		speakers[i] = &data.SpeakerInfo{
			SpeakerName:          supportedSpeaker.ChineseName,
			ParameterSpeakerName: supportedSpeaker.Name,
			IsSupportEmotion:     isEmotion,
			IsSupportMixedVoice:  isMixedVoice,
		}
		C.free(unsafe.Pointer(cname))
		C.free(unsafe.Pointer(userSpace))
	}
	close(notify)
	return &TTSService{
		ChanLen:           20,
		ResPath:           s.App.Path,
		version:           C.GoString(version),
		resServiceVersion: C.GoString(resServiceVersion),
		Speakers:          speakers,
		SpeakerSetting:    speakerSetting,
		Logger:            logger,
	}
}

// GetSDKVersion 获取sdk的版本
func (t *TTSService) GetSDKVersion() string {
	return t.version
}

// GetResServiceVersion 获取sdk的热修复版本
func (t *TTSService) GetResServiceVersion() string {
	version := C.ActionSynthesizer_GetResServiceVersion()
	return C.GoString(version)
}

func (t *TTSService) GetSpeakers() []*data.SpeakerInfo {
	return t.Speakers
}

func (t *TTSService) GetUserSpeakers(userspace string) ([]string, error) {
	cUserspace := C.CString(userspace)
	defer C.free(unsafe.Pointer(cUserspace))
	var speakerNum C.uint

	if ok := C.GetUserSpaceSupportedSpeaker(cUserspace, (***C.char)(unsafe.Pointer(&C.speakerList)), &speakerNum); ok != 0 {
		return nil, errors.New("fail to get user space supported speaker")
	}
	speakers := make([]string, speakerNum)
	if C.speakerList != nil && speakerNum > 0 {
		tmpSlice := (*[1 << 30]*C.char)(unsafe.Pointer(C.speakerList))[:int(speakerNum):int(speakerNum)]
		for i, v := range tmpSlice {
			speakers[i] = C.GoString(v)
		}
	}
	return speakers, nil

}

func (t *TTSService) GetSpeakerSetting() *data.SpeakerSetting {
	return t.SpeakerSetting
}

func (t *TTSService) GetSupportedPitch() []*v2.MessagePitch {
	pitchList := make([]*v2.MessagePitch, 0, len(t.SpeakerSetting.SupportedPitch))
	for _, pitch := range t.SpeakerSetting.SupportedPitch {
		pitchList = append(pitchList, &v2.MessagePitch{
			Name:        pitch.Name,
			ChineseName: pitch.ChineseName,
		})
	}
	return pitchList
}

func (t *TTSService) GetSupportedEmotion() []*v2.MessageEmotion {
	emotionList := make([]*v2.MessageEmotion, 0, len(t.SpeakerSetting.SupportedEmotion))
	for _, emotion := range t.SpeakerSetting.SupportedEmotion {
		emotionList = append(emotionList, &v2.MessageEmotion{
			Name:        emotion.Name,
			ChineseName: emotion.ChineseName,
		})
	}
	return emotionList
}

func (t *TTSService) GetSupportedMovement() []*v2.MessageMovement {
	movementList := make([]*v2.MessageMovement, 0, len(t.SpeakerSetting.SupportedEmotion))
	for _, emotion := range t.SpeakerSetting.SupportedMovementDescriptor {
		movementList = append(movementList, &v2.MessageMovement{
			Name:        emotion.Name,
			ChineseName: emotion.ChineseName,
		})
	}
	return movementList
}

func (t *TTSService) GetSupportedExpression() []*v2.MessageExpression {
	ExpressionList := make([]*v2.MessageExpression, 0, len(t.SpeakerSetting.SupportedExpressionDescriptor))
	for _, express := range t.SpeakerSetting.SupportedExpressionDescriptor {
		ExpressionList = append(ExpressionList, &v2.MessageExpression{
			Name:        express.Name,
			ChineseName: express.ChineseName,
		})
	}
	return ExpressionList
}

//func (t *TTSService) GetSupportedDigitalPerson() []*v2.MessageDigitalPerson {
//	digitalPersonList := make([]*v2.MessageDigitalPerson, 0, len(t.SpeakerSetting.SupportedDigitalPerson))
//	for _, digitalPerson := range t.SpeakerSetting.SupportedDigitalPerson {
//		digitalPersonList = append(digitalPersonList, &v2.MessageDigitalPerson{
//			Id:          int32(digitalPerson.Id),
//			Name:        digitalPerson.Name,
//			ChineseName: digitalPerson.ChineseName,
//		})
//	}
//	return digitalPersonList
//}

func (t *TTSService) CallTTSServiceV2(req *data.Speaker, pUserData int32, traceId string) (int, error) {
	var sdkSettings = C.TtsSetting{}

	sdkSettings.speaker = C.CString(req.ParameterSpeakerName)
	defer C.free(unsafe.Pointer(sdkSettings.speaker))
	sdkSettings.speed = C.CString(req.Speed)
	defer C.free(unsafe.Pointer(sdkSettings.speed))
	sdkSettings.volume = C.CString(req.Volume)
	defer C.free(unsafe.Pointer(sdkSettings.volume))
	sdkSettings.pitch = C.CString(req.Pitch)
	defer C.free(unsafe.Pointer(sdkSettings.pitch))
	sdkSettings.speakingStyle = C.CString(req.Emotions)
	defer C.free(unsafe.Pointer(sdkSettings.speakingStyle))
	sdkSettings.featureSet = C.uint(paramFormatter(req.ParameterFlag))
	sdkSettings.expressionDescriptor = C.CString(req.Expression)
	defer C.free(unsafe.Pointer(sdkSettings.expressionDescriptor))
	sdkSettings.movementDescriptor = C.CString(req.Movement)
	defer C.free(unsafe.Pointer(sdkSettings.movementDescriptor))
	text := C.CString(req.Text)
	defer C.free(unsafe.Pointer(text))
	traceIdC := C.CString(traceId)
	defer C.free(unsafe.Pointer(traceIdC))
	sdkSettings.languageTip = C.CString(req.Language)
	defer C.free(unsafe.Pointer(sdkSettings.languageTip))
	sdkSettings.userSpace = C.CString(req.Userspace)
	defer C.free(unsafe.Pointer(sdkSettings.userSpace))
	sdkSettings.audioEncoding = C.int(req.AudioEncoding)
	idC := C.ActionSynthesizer_SynthesizeAction(
		text,
		&sdkSettings,
		&actionCallback,
		unsafe.Pointer(uintptr(pUserData)),
		traceIdC,
	)
	id := int(idC)
	log.NewHelper(t.Logger).Infof("ActionSynthesizer_SynthesizeAction return id:%d", int(id))

	if id < 0 {
		switch id {
		case -1:
			return id, errors.New("tts req param is invalid")
		default:
			return id, fmt.Errorf("fail to call the api of the sdk;id(%d)", id)
		}
	}
	return id, nil
}

func (t *TTSService) CallTTSServiceV1(req *v1.TtsReq, pUserData int32) (int, error) {
	var setting = C.TtsSetting{}

	setting.speaker = C.CString(req.ParameterSpeakerName)
	defer C.free(unsafe.Pointer(setting.speaker))
	setting.speed = C.CString(req.Speed)
	defer C.free(unsafe.Pointer(setting.speed))
	setting.volume = C.CString(req.Volume)
	defer C.free(unsafe.Pointer(setting.volume))
	setting.pitch = C.CString(req.Pitch)
	defer C.free(unsafe.Pointer(setting.pitch))
	setting.speakingStyle = C.CString(req.Emotions)
	defer C.free(unsafe.Pointer(setting.languageTip))
	setting.languageTip = C.CString(req.Language)
	defer C.free(unsafe.Pointer(setting.speakingStyle))
	setting.featureSet = C.uint(3)

	text := C.CString(req.Text)
	defer C.free(unsafe.Pointer(text))
	traceId := C.CString(req.RootTraceId + "_" + req.TraceId)
	defer C.free(unsafe.Pointer(traceId))
	idC := C.ActionSynthesizer_SynthesizeAction_V1(
		text,
		&setting,
		&ttsCallback,
		unsafe.Pointer(uintptr(pUserData)),
		traceId,
	)
	id := int(idC)
	log.NewHelper(t.Logger).Infof("ActionSynthesizer_SynthesizeAction_V1 return id:%d", int(id))
	if id < 0 {
		switch id {
		case -1:
			return id, errors.New("tts req param is invalid")
		default:
			return id, fmt.Errorf("fail to call the api of the sdk;id(%d)", id)
		}
	}
	return id, nil
}

func (t *TTSService) GeneHandlerObjectV2(ctx context.Context, speaker string, logger *log.Helper) *data.HandlerObjectV2 {
	backChan := make(chan v2.TtsRes, t.ChanLen)
	paramMap := make(map[string]interface{})
	return &data.HandlerObjectV2{
		HandlerObject: data.HandlerObject{
			Context: ctx,
			SpeakerInfo: data.SpeakerInfo{
				ParameterSpeakerName: speaker,
			},
			ParamMap: paramMap,
			Log:      logger,
		},
		BackChan: backChan,
	}
}

func (t *TTSService) GeneHandlerObjectV1(ctx context.Context, speaker string, logger *log.Helper) *data.HandlerObjectV1 {
	backChan := make(chan v1.TtsRes, 10)
	paramMap := make(map[string]interface{})
	return &data.HandlerObjectV1{
		HandlerObject: data.HandlerObject{
			Context: ctx,
			SpeakerInfo: data.SpeakerInfo{
				ParameterSpeakerName: speaker,
			},
			ParamMap: paramMap,
			Log:      logger,
		},
		BackChan: backChan,
	}
}

func paramFormatter(flag map[string]string) uint {
	if flag == nil {
		return 0
	}
	// 口型的key:mouth、动作的key:movement、表情的key:expression，c中有一定顺序表情为0，口型1，动作2
	var flagList = [3]string{"expression", "mouth", "movement"}
	var res uint = 0
	for i, s := range flagList {
		v, exists := flag[s]
		if exists && strings.ToLower(v) == "true" {
			res += 1 << i
		}
	}
	return res
}

func (t *TTSService) CancelTTSService(id int) {
	C.ActionSynthesizer_CancelAction(C.int(id))
}
