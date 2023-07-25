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
	"fmt"
	"log"
	v1 "speech-tts/api/tts/v1"
	"speech-tts/internal/data"
	"speech-tts/internal/pkg/pointer"
	"speech-tts/internal/pkg/trace"
	"unsafe"
)

//export goOnStartV1
func goOnStartV1(pUserData unsafe.Pointer) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnStartV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnStartV1", nil)
	defer span.End()
	object.Log.Info("end to goOnStartV1;pUserData:", pUserData)
	return
}

/**
 * 合成的音频数据
 */

//export goOnAudioV1
func goOnAudioV1(pUserData unsafe.Pointer, dataAudio *C.char, len C.int) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnAudioV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnAudioV1", nil)
	defer span.End()
	object.Log.Info("start to OnAudioV1; pUserData:", pUserData)
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
		if debugInfo, ok := object.ParamMap["debugInfo"].(string); ok {
			response.DebugInfo = debugInfo
		}
		object.ParamMap = make(map[string]interface{})
	}
	sendRespV1(object, response)
	object.Log.Info("end to OnAudioV1 pUserData:%d", pUserData)
}

/**
 * 语音合成结束
 */

//export goOnEndV1
func goOnEndV1(pUserData unsafe.Pointer, flag C.int) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnEndV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnEndV1", nil)
	defer span.End()
	object.Log.Info("start to OnEndV1;pUserData:", pUserData)

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
	object.Log.Infof("end to OnEndV1 pUserData:%d", pUserData)
}

//export goOnDebugV1
func goOnDebugV1(pUserData unsafe.Pointer, info *C.char) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnDebugV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnDebugV1", nil)
	defer span.End()
	debugInfo := string(C.GoString(info))
	object.Log.Infof("start to goOnDebugV1;debugInfo:%s; pUserData:%d", debugInfo, pUserData)

	temp := object.ParamMap["debugInfo"]
	if i, exist := temp.(string); exist {
		object.ParamMap["debugInfo"] = fmt.Sprintf("%s\n%s", i, debugInfo)
	} else {
		object.ParamMap["debugInfo"] = debugInfo
	}
	object.Log.Infof("end to goOnDebugV1; pUserData:%d", pUserData)
	return
}

//export goOnTimedMouthShapeV1
func goOnTimedMouthShapeV1(pUserData unsafe.Pointer, mouth *C.TimedMouthShape, size C.int, text *C.char) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnTimedMouthShapeV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnTimedMouthShapeV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnTimedMouthShapeV1; pUserData:%d", pUserData)

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
	object.Log.Infof("end to goOnTimedMouthShapeV1 pUserData:%d", pUserData)
}

//export goOnCurTextSegmentV1
func goOnCurTextSegmentV1(pUserData unsafe.Pointer, normalizedText *C.char, originalText *C.char) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnCurTextSegmentV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnCurTextSegmentV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnCurTextSegmentV1; pUserData:%d", pUserData)

	if object.ParamMap != nil {
		object.ParamMap["normalizedText"] = C.GoString(normalizedText)
		object.ParamMap["originalText"] = C.GoString(originalText)
	}
	object.Log.Infof("end to goOnCurTextSegmentV1 pUserData:%d", pUserData)
}

//export goOnFacialExpressionV1
func goOnFacialExpressionV1(pUserData unsafe.Pointer, expression *C.FacialExpression) {
	object := getHandlerObjectV1(pUserData)
	if object == nil {
		log.Println("goOnFacialExpressionV1; irregularity type")
		return
	}
	_, span := trace.NewTraceSpan(object.Context, "goOnFacialExpressionV1", nil)
	defer span.End()
	object.Log.Infof("start to goOnFacialExpressionV1; pUserData:%d", pUserData)

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
	object.Log.Infof("end to goOnFacialExpressionV1 pUserData:%d", pUserData)
}

func sendRespV1(object *data.HandlerObjectV1, response v1.TtsRes) {
	if object != nil {
		return
	}
	object.Lock()
	defer object.Unlock()
	if object.IsInterrupted {
		log.Println("HandlerObjectV1;interrupted by cancel")
		return
	}
	if object.BackChan != nil {
		object.BackChan <- response
	}
}

func getHandlerObjectV1(pUserData unsafe.Pointer) *data.HandlerObjectV1 {
	handlerObject := pointer.Load(pUserData)
	if handlerObject == nil {
		log.Println("don't find to handler object; pUserData", pUserData)
		return nil
	}
	object, ok := handlerObject.(*data.HandlerObjectV1)
	if !ok {
		log.Println(" irregularity handler object;pUserData: ", pUserData)
		return nil
	}

	return object
}
