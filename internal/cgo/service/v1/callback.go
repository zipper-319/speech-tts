package TtsSdkv1

/*
#cgo CFLAGS: -I../include
#cgo LDFLAGS: -L../lib_interface -lCmTts
#include <stdio.h>
#include <stdlib.h>
#include "ActionSynthesizer.h"
#include "AnimationDef.h"
#include "MouthShape.h"
// #include "ExpressionDef.h"
*/
import "C"
import (
	v1 "speech-tts/api/tts/v1"
	"speech-tts/internal/data"
	"unsafe"
)

//export goOnStartV1
func goOnStartV1(pUserData unsafe.Pointer) {
	// todo: 这里start我先把它屏蔽了
	return
}

/**
 * 合成的音频数据
 */

//export goOnAudioV1
func goOnAudioV1(pUserData unsafe.Pointer, dataAudio *C.char, len C.int) {

	object := (*data.HandlerObjectV1)(pUserData)
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
	sendResp(object, response)

}

/**
 * 语音合成结束
 */

//export goOnEndV1
func goOnEndV1(pUserData unsafe.Pointer, flag C.int) {
	object := (*data.HandlerObjectV1)(pUserData)
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

	sendResp(object, response)
	close(object.BackChan)

}

//export goOnDebugV1
func goOnDebugV1(pUserData unsafe.Pointer, debugtype *C.char, info *C.char) {
	return
}

//export goOnTimedMouthShapeV1
func goOnTimedMouthShapeV1(pUserData unsafe.Pointer, mouth *C.TimedMouthShape, size C.int, text *C.char) {
	object := (*data.HandlerObjectV1)(pUserData)

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
}

//export goOnCurTextSegmentV1
func goOnCurTextSegmentV1(pUserData unsafe.Pointer, normalizedText *C.char, originalText *C.char) {
	object := (*data.HandlerObjectV1)(pUserData)

	if object.ParamMap != nil {
		object.ParamMap["normalizedText"] = C.GoString(normalizedText)
		object.ParamMap["originalText"] = C.GoString(originalText)
	}
}

//export goOnFacialExpressionV1
func goOnFacialExpressionV1(pUserData unsafe.Pointer, expression *C.FacialExpression) {
	object := (*data.HandlerObjectV1)(pUserData)

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
}

func sendResp(object *data.HandlerObjectV1, response v1.TtsRes) {
	if object != nil && object.BackChan != nil {
		object.BackChan <- response
	}
}
