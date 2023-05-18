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

const SpeakerDescriptor* supportedSpeakers = NULL;
unsigned int element_num = 0;

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


static void createUser(void **pUser) {
	if(pUser) *pUser = malloc(sizeof(SpeakerDescriptor));
}
const char* getCArrValue(const char** src, int i) {
	return src[i];
}
*/
import "C"
import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "speech-tts/api/tts/v1"
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/data"
	"strings"
	"unsafe"
)

var ProviderSet = wire.NewSet(NewTTSService)


// GetSDKVersion 获取sdk的版本
func GetSDKVersion() string {
	return C.GoString(C.ActionSynthesizer_GetVersion())
}

// GetResServiceVersion 获取sdk的热修复版本
func GetResServiceVersion() string {
	return C.GoString(C.ActionSynthesizer_GetResServiceVersion())
}

func getCallbackV2() *C.ActionCallback {
	callback := C.ActionCallback{}
	callback.onStart = C.typOnStart(C.goOnStart)
	callback.onSynthesizedData = C.typOnSynthesizedData(C.goOnSynthesizedData)
	callback.onEnd = C.typOnEnd(C.goOnEnd)
	callback.onDebug = C.typOnDebug(C.goOnDebug)
	callback.onTimedMouthShape = C.typOnTimedMouthShape(C.goOnTimedMouthShape)
	callback.onFacialExpression = C.typOnFacialExpression(C.goOnFacialExpression)
	callback.onBodyMovement = C.typOnBodyMovement(C.goOnBodyMovement)
	callback.onActionElement = C.typOnActionElement(C.goOnActionElement)
	return &callback
}

func getCallbackV1() *C.TTS_Callback{
	var ttsCallback = C.TTS_Callback{}

	ttsCallback.onStart = C.typOnStartV1(C.goOnStartV1)
	ttsCallback.onAudio = C.typOnAudioV1(C.goOnAudioV1)
	ttsCallback.onEnd = C.typOnEndV1(C.goOnEndV1)
	ttsCallback.onDebug = C.typOnDebugV1(C.goOnDebugV1)
	ttsCallback.onTimedMouthShape = C.typOnTimedMouthShapeV1(C.goOnTimedMouthShapeV1)
	ttsCallback.onCurTextSegment = C.typOnCurTextSegmentV1(C.goOnCurTextSegmentV1)
	ttsCallback.onFacialExpression = C.typOnFacialExpressionV1(C.goOnFacialExpressionV1)

	return &ttsCallback
}

type TTSService struct {
	ResPath           string
	Version           string
	ResServiceVersion string
	Speakers          []*data.SpeakerInfo
	*data.SpeakerSetting
}

func NewTTSService(resPath string, speakerSetting *data.SpeakerSetting) *TTSService {
	cResPath := C.CString(resPath)
	defer C.free(unsafe.Pointer(cResPath))
	C.ActionSynthesizer_Init(cResPath)
	version := C.ActionSynthesizer_GetVersion()
	resServiceVersion := C.ActionSynthesizer_GetResServiceVersion()
	// 发音人初始化
	speakers := make([]*data.SpeakerInfo, len(speakerSetting.SupportedSpeaker))
	for i, supportedSpeaker := range speakerSetting.SupportedSpeaker {
		cname := C.CString(supportedSpeaker.Name)
		m := C.GetSpeakerDescriptor(cname)
		speakers[i] = &data.SpeakerInfo{
			SpeakerId:            supportedSpeaker.Id,
			SpeakerName:          supportedSpeaker.ChineseName,
			ParameterSpeakerName: supportedSpeaker.Name,
			IsSupportEmotion:     m.flags>>1 == 1,
		}
		C.free(unsafe.Pointer(cname))
	}
	return &TTSService{
		ResPath:           resPath,
		Version:           C.GoString(version),
		ResServiceVersion: C.GoString(resServiceVersion),
		SpeakerSetting:    speakerSetting,
		Speakers:          speakers,
	}
}

func (t *TTSService) GetSpeakers() []*data.SpeakerInfo {
	return t.Speakers
}

func (t *TTSService) GetSpeakerSetting() *data.SpeakerSetting {
	return t.SpeakerSetting
}

func (t *TTSService) CallTTSServiceV2(req *v2.TtsReq, object *data.HandlerObjectV2) error {
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
	sdkSettings.digitalPerson = C.CString(req.ParameterDigitalPerson)
	defer C.free(unsafe.Pointer(sdkSettings.digitalPerson))
	id := C.ActionSynthesizer_SynthesizeAction(
		C.CString(req.Text),
		&sdkSettings,
		getCallbackV2(),
		unsafe.Pointer(object),
		C.CString(req.RootTraceId+"_"+req.TraceId))

	if id < 0 {
		return errors.New("fail to call api of the sdk")
	}
	return nil
}

func (t *TTSService) CallTTSServiceV1(req *v1.TtsReq, pUserData unsafe.Pointer) error {
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
	defer C.free(unsafe.Pointer(setting.speakingStyle))
	setting.featureSet = C.uint(3)

	if len(t.SupportedDigitalPerson) > 0 {
		setting.digitalPerson = C.CString(t.SupportedDigitalPerson[0].Name)
	} else {
		setting.digitalPerson = C.CString("SweetGirl")
	}
	defer C.free(unsafe.Pointer(setting.digitalPerson))

	id := C.ActionSynthesizer_SynthesizeAction_V1(
		C.CString(req.Text),
		&setting,
		getCallbackV1(),
		pUserData,
		C.CString(req.RootTraceId+"_"+req.TraceId),
	)
	if id < 0 {
		return errors.New("fail to call api of the sdk")
	}
	return nil
}

func (t *TTSService) GeneHandlerObjectV2(ctx context.Context, speaker string) *data.HandlerObjectV2 {
	backChan := make(chan v2.TtsRes, 10)
	paramMap := make(map[string]interface{})
	return &data.HandlerObjectV2{
		HandlerObject: data.HandlerObject{
			Context: ctx,
			SpeakerInfo: data.SpeakerInfo{
				ParameterSpeakerName: speaker,
			},
			ParamMap: paramMap,
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
		} else {
			res += 0 << i
		}
	}
	return res
}
