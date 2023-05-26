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
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/data"
	"speech-tts/internal/pkg/pointer"
	"unsafe"
)

/**
 * 语音合成开始（即首包数据已准备好）
 */

//export goOnStart
func goOnStart(pUserData unsafe.Pointer, ttsText *C.char, facialExpressionConfig *C.FacialExpressionConfig, bodyMovementConfig *C.BodyMovementConfig) {

	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnStart")

	var facialExpressionFrameDim, bodyMovemenFrameDim uint64 = 0, 0
	var facialExpressionFrameDurMs, bodyMovemenFrameDurMs float32 = 0, 0
	if facialExpressionConfig != nil {
		facialExpressionFrameDim = uint64(facialExpressionConfig.frameDim)
		facialExpressionFrameDurMs = float32(facialExpressionConfig.frameDurMs)
	}
	if bodyMovementConfig != nil {
		bodyMovemenFrameDim = uint64(bodyMovementConfig.frameDim)
		bodyMovemenFrameDurMs = float32(bodyMovementConfig.frameDurMs)
	}
	paramSetting := make(map[string]interface{})
	paramSetting["FacialExpressionFrameDim"] = int32(facialExpressionFrameDim)
	paramSetting["BodyMovemenFrameDim"] = int32(bodyMovemenFrameDim)
	object.ParamMap = paramSetting
	response := v2.TtsRes{
		Status: 1,
		ResultOneof: &v2.TtsRes_ConfigText{
			ConfigText: &v2.ConfigAndText{
				Text: C.GoString(ttsText),
				FacialExpressionConfig: &v2.FacialExpressionConfig{
					FrameDim:   int32(facialExpressionFrameDim),
					FrameDurMs: facialExpressionFrameDurMs,
				},
				BodyMovementConfig: &v2.BodyMovementConfig{
					FrameDim:   int32(bodyMovemenFrameDim),
					FrameDurMs: bodyMovemenFrameDurMs,
				},
			},
		},
	}
	sendResp(object, response)
	object.Log.Info("end to goOnStart")
}

/**
 * 语音合成结束
 */

//export goOnEnd
func goOnEnd(pUserData unsafe.Pointer, flag C.int) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnEnd")

	response := v2.TtsRes{
		ErrorCode: int32(flag),
		Status:    3,
	}
	sendResp(object, response)
	close(object.BackChan)
	object.Log.Info("end to goOnEnd")
}

//export goOnDebug
func goOnDebug(pUserData unsafe.Pointer, debugtype *C.char, info *C.char) {

	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnDebug")

	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_DebugInfo{
			DebugInfo: &v2.DebugInfo{
				DebugType: C.GoString(debugtype),
				Info:      C.GoString(info),
			},
		},
	}

	sendResp(object, response)
	object.Log.Info("end to goOnDebug")
}

//export goOnTimedMouthShape
func goOnTimedMouthShape(pUserData unsafe.Pointer, mouth *C.TimedMouthShape, size C.int, startTimeMs C.float) {

	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnTimedMouthShape")

	var mouths = make([]*v2.TimedMouthShape, int32(size))
	for i := 0; i < int(size); i++ {
		m := *(*C.TimedMouthShape)(unsafe.Pointer(uintptr(unsafe.Pointer(mouth)) + uintptr(C.sizeof_TimedMouthShape*C.int(i))))
		mouths[i] = &v2.TimedMouthShape{
			DurationUs: uint64(m.durationUs),
			Mouth:      int32(m.mouth),
		}
	}
	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_TimeMouthShapes{
			TimeMouthShapes: &v2.TimedMouthShapes{
				Mouths:      mouths,
				StartTimeMs: float32(startTimeMs),
			},
		},
	}

	sendResp(object, response)
	object.Log.Info("end to goOnTimedMouthShape")
}

//export goOnActionElement
func goOnActionElement(pUserData unsafe.Pointer, ctype C.int, url *C.char, operation_type C.int, coordinate unsafe.Pointer, render_duration C.int) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnActionElement")

	coordinateC := *(*C.Coordinate)(coordinate)
	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_ActionElement{
			ActionElement: &v2.ActionElement{
				ActionType:    int32(ctype),
				Url:           C.GoString(url),
				OperationType: int32(operation_type),
				Coordinate: &v2.Coordinate{
					Off:   int32(coordinateC.off_utf8),
					Len:   int32(coordinateC.len_utf8),
					Order: int32(coordinateC.order),
				},
				RenderDuration: int32(render_duration),
			},
		},
	}
	sendResp(object, response)
	object.Log.Info("end to goOnActionElement")
}

//export goOnSynthesizedData
func goOnSynthesizedData(pUserData unsafe.Pointer, audioData *C.SynthesizedAudio, coordinate *C.Coordinate) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnSynthesizedData")

	length := (C.int)(audioData.audio_size)
	wav := C.GoBytes(unsafe.Pointer(audioData.audio_data), 2*length)

	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_SynthesizedAudio{
			SynthesizedAudio: &v2.SynthesizedAudio{
				Pcm: wav,
				Coordinate: &v2.Coordinate{
					Off:   int32(coordinate.off_utf8),
					Len:   int32(coordinate.len_utf8),
					Order: int32(coordinate.order),
				},
				IsPunctuation: int32(audioData.flags),
			},
		},
	}
	sendResp(object, response)
	object.Log.Info("end to goOnSynthesizedData")
}

//export goOnFacialExpression
func goOnFacialExpression(pUserData unsafe.Pointer, facialExpressionData *C.FacialExpressionSegment) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnFacialExpression")

	var framedDim uint64
	if fefd, ok := (object.ParamMap["FacialExpressionFrameDim"]).(int32); ok {
		framedDim = uint64(fefd)
	}
	var expression *v2.Expression
	if facialExpressionData != nil {
		frameSize := uint64(facialExpressionData.frameSize)
		startTimeMs := float32(facialExpressionData.startTimeMs)
		size := int32(frameSize * framedDim)
		var Expressiondata = make([]float32, size)
		for i := 0; i < int(size); i++ {
			data := *(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(facialExpressionData.data)) + uintptr(C.sizeof_float*C.int(i))))
			Expressiondata[i] = float32(data)
		}
		expression = &v2.Expression{
			Data:        Expressiondata,
			FrameSize:   int32(frameSize),
			StartTimeMs: startTimeMs,
		}
	}
	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_Expression{
			Expression: expression,
		},
	}
	sendResp(object, response)
	object.Log.Info("end to goOnFacialExpression")
}

//export goOnBodyMovement
func goOnBodyMovement(pUserData unsafe.Pointer, bodyMovementData unsafe.Pointer) {
	handlerObject := pointer.Load(pUserData)
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		panic("irregularity type")
	}
	object.Log.Info("start to goOnBodyMovement")

	var bodyMovement *v2.BodyMovement
	bodyMovementDataC := (*C.BodyMovementSegment)(bodyMovementData)
	// _, bodyMovementDim := getTmpDecoder(*id)
	var bodyMovementFrameSizeFrameDim uint64
	if bodyMovementDim, ok := (object.ParamMap["BodyMovemenFrameDim"]).(int32); ok {
		bodyMovementFrameSizeFrameDim = uint64(bodyMovementDim)
	}
	if bodyMovementDataC != nil {
		bodyMovementFrameSize := uint64(bodyMovementDataC.frameSize)
		startTimeMs := float32(bodyMovementDataC.startTimeMs)
		bodyMovementDataSize := int32(bodyMovementFrameSize * bodyMovementFrameSizeFrameDim)
		var bodyMovementDataList = make([]float32, bodyMovementDataSize)
		for i := 0; i < int(bodyMovementDataSize); i++ {
			data := *(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(bodyMovementDataC.data)) + uintptr(C.sizeof_float*C.int(i))))
			bodyMovementDataList[i] = float32(data)
		}
		bodyMovement = &v2.BodyMovement{
			Data:        bodyMovementDataList,
			FrameSize:   int32(bodyMovementFrameSize),
			StartTimeMs: startTimeMs,
		}
	}
	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_BodyMovement{
			BodyMovement: bodyMovement,
		},
	}
	sendResp(object, response)
	object.Log.Info("end to goOnBodyMovement")
}

func sendResp(object *data.HandlerObjectV2, response v2.TtsRes) {
	if object != nil && object.BackChan != nil {
		object.BackChan <- response
	}
}
