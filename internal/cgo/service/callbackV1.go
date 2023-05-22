package service

/*
#include <stdio.h>
#include <stdlib.h>
#include "ActionSynthesizer.h"
#include "AnimationDef.h"
#include "MouthShape.h"
*/
import "C"
import (
	v1 "speech-tts/api/tts/v1"
	"speech-tts/internal/data"
	"speech-tts/internal/pkg/pointer"
	"speech-tts/internal/pkg/trace"
	"unsafe"
)

//export goOnStartV1
func goOnStartV1(pUserData unsafe.Pointer) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnStartV1", nil)
	defer span.End()
	object.Log.Infof("enter to OnAudioV1")
	return
}

/**
 * 合成的音频数据
 */

//export goOnAudioV1
func goOnAudioV1(pUserData unsafe.Pointer, dataAudio *C.char, len C.int) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnAudioV1", nil)
	defer span.End()
	object.Log.Infof("start to OnAudioV1")
	response := v1.TtsRes{
		Pcm:    C.GoBytes(unsafe.Pointer(dataAudio), len),
		Status: v1.PcmStatus_STATUS_MID,
		Error:  v1.TtsErr_TTS_ERR_OK,
	}
	if object.ParamMap != nil {
		if mouths, ok := object.ParamMap["mouths"].([]*v1.TimedMouthShape); ok {
			response.Mouths = mouths
		}
		if normalizedText, ok := object.ParamMap["normalizedText"].(string); ok {
			response.NormalizedText = normalizedText
		}
		if originalText, ok := object.ParamMap["originalText"].(string); ok {
			response.OriginalText = originalText
		}
		if expression, ok := object.ParamMap["expression"].(v1.Expression); ok {
			response.Expression = &expression
		}
		object.ParamMap = make(map[string]interface{})
	}
	sendRespV1(object, response)
	object.Log.Infof("end to OnAudioV1")
}

/**
 * 语音合成结束
 */

//export goOnEndV1
func goOnEndV1(pUserData unsafe.Pointer, flag C.int) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnEndV1", nil)
	defer span.End()
	object.Log.Infof("start to OnEndV1")

	var err v1.TtsErr

	if flag == 0 {
		err = v1.TtsErr_TTS_ERR_OK
	} else if flag == 1 {
		err = v1.TtsErr_TTS_ERR_SYN_CANCELLED
	} else if flag == 2 {
		err = v1.TtsErr_TTS_ERR_SYN_FAILURE
	}
	response := v1.TtsRes{
		Status: v1.PcmStatus_STATUS_END,
		Error:  err,
	}

	sendRespV1(object, response)
	close(object.BackChan)
	object.Log.Infof("end to OnEndV1")
}

//export goOnDebugV1
func goOnDebugV1(pUserData unsafe.Pointer, debugtype *C.char, info *C.char) {
	return
}

//export goOnTimedMouthShapeV1
func goOnTimedMouthShapeV1(pUserData unsafe.Pointer, mouth *C.TimedMouthShape, size C.int, text *C.char) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnTimedMouthShapeV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnTimedMouthShapeV1")

	var mouthShapes = make([]*v1.TimedMouthShape, int32(size))
	for i := 0; i < int(size); i++ {
		m := *(*C.TimedMouthShape)(unsafe.Pointer(uintptr(unsafe.Pointer(mouth)) + uintptr(C.sizeof_TimedMouthShape*C.int(i))))
		mouthShapes[i] = &v1.TimedMouthShape{
			DurationUs: uint64(m.durationUs),
			Mouth:      int32(m.mouth),
		}
	}
	temp := object.ParamMap["mouths"]
	if mouths, ok := temp.([]*v1.TimedMouthShape); ok {
		object.ParamMap["mouths"] = append(mouths, mouthShapes...)
	} else {
		object.ParamMap["mouths"] = mouthShapes
	}
	object.Log.Infof("end to goOnTimedMouthShapeV1")
}

//export goOnCurTextSegmentV1
func goOnCurTextSegmentV1(pUserData unsafe.Pointer, normalizedText *C.char, originalText *C.char) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnCurTextSegmentV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnCurTextSegmentV1")

	if object.ParamMap != nil {
		object.ParamMap["normalizedText"] = C.GoString(normalizedText)
		object.ParamMap["originalText"] = C.GoString(originalText)
	}
	object.Log.Infof("end to goOnCurTextSegmentV1")
}

//export goOnFacialExpressionV1
func goOnFacialExpressionV1(pUserData unsafe.Pointer, expression *C.FacialExpression) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		panic("irregularity type")
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnFacialExpressionV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnFacialExpressionV1")

	frameSize := uint64(expression.frame_size)
	frameDim := uint64(expression.frame_dim)
	size := int32(frameSize * frameDim)
	var expressionData = make([]float32, size)
	for i := 0; i < int(size); i++ {
		data := *(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(expression.data)) + uintptr(C.sizeof_float*C.int(i))))
		expressionData[i] = float32(data)
	}

	if object.ParamMap != nil {
		object.ParamMap["expression"] = v1.Expression{
			Data:      expressionData,
			FrameSize: int32(frameSize),
			FrameDim:  int32(frameDim),
			FrameTime: float32(expression.frame_time),
		}
	}
	object.Log.Infof("end to goOnFacialExpressionV1")
}

func sendRespV1(object *data.HandlerObjectV1, response v1.TtsRes) {
	if object != nil && object.BackChan != nil && !object.IsInterrupted {
		object.BackChan <- response
	}
}
