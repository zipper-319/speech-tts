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
	"log"
	v2 "speech-tts/api/tts/v2"
	"speech-tts/internal/data"
	"speech-tts/internal/pkg/pointer"
	"unsafe"
)

/**
 * 语音合成开始（即首包数据已准备好）
 */

//export goOnStart
func goOnStart(pUserData unsafe.Pointer, ttsText *C.char, audioConfig *C.AudioConfig, facialExpressionConfig *C.FacialExpressionConfig, bodyMovementConfig *C.BodyMovementConfig) {

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnStart; irregularity type")
		return
	}
	object.Log.Infof("start to OnStart; pUserData: %d", pUserData)

	var facialExpressionFrameDim, bodyMovementFrameDim uint64
	var facialExpressionFrameDurMs, bodyMovementFrameDurMs float32
	var expressControlNameList string
	audioConfigData := new(v2.AudioConfig)
	//var movementControlNameList []string
	if facialExpressionConfig != nil {
		frameDimLen := uint64(facialExpressionConfig.frameDim)
		facialExpressionFrameDim = frameDimLen
		facialExpressionFrameDurMs = float32(facialExpressionConfig.frameDurMs)

		if facialExpressionConfig.meta_data != nil && int(facialExpressionConfig.frameDim) > 0 {
			expressControlNameList = C.GoString(facialExpressionConfig.meta_data)
		}
	}

	if bodyMovementConfig != nil {
		bodyMovementFrameDim = uint64(bodyMovementConfig.frameDim)
		bodyMovementFrameDurMs = float32(bodyMovementConfig.frameDurMs)

		//frameDimLen := bodyMovementFrameDim
		//if bodyMovementConfig.control_name != nil && int(frameDimLen) > 0 {
		//	movementControlNameList = make([]string, frameDimLen)
		//	tmpSlice := (*[1 << 30]*C.char)(unsafe.Pointer(bodyMovementConfig.control_name))[:frameDimLen:frameDimLen]
		//	for i, s := range tmpSlice {
		//		movementControlNameList[i] = C.GoString(s)
		//	}
		//}
	}
	if audioConfig != nil {
		audioConfigData.AudioEncoding = int32(audioConfig.audio_encoding)
		audioConfigData.SamplingRate = int32(audioConfig.sampling_rate)
		audioConfigData.Channels = int32(audioConfig.channels)
	}
	paramSetting := make(map[string]interface{})
	paramSetting["FacialExpressionFrameDim"] = int32(facialExpressionFrameDim)
	paramSetting["BodyMovementFrameDim"] = int32(bodyMovementFrameDim)
	object.ParamMap = paramSetting
	response := v2.TtsRes{
		Status: 1,
		ResultOneof: &v2.TtsRes_ConfigText{
			ConfigText: &v2.ConfigAndText{
				Text: C.GoString(ttsText),
				FacialExpressionConfig: &v2.FacialExpressionConfig{
					FrameDim:   int32(facialExpressionFrameDim),
					FrameDurMs: facialExpressionFrameDurMs,
					MetaData:   expressControlNameList,
				},
				BodyMovementConfig: &v2.BodyMovementConfig{
					FrameDim:   int32(bodyMovementFrameDim),
					FrameDurMs: bodyMovementFrameDurMs,
				},
				AudioConfig: audioConfigData,
			},
		},
	}
	sendResp(object, response)
	object.Log.Infof("end to goOnStart pUserData: %d", pUserData)
}

/**
 * 语音合成结束
 */

//export goOnEnd
func goOnEnd(pUserData unsafe.Pointer, flag C.int) {
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnEnd; irregularity type")
		return
	}
	object.Log.Infof("start to goOnEnd; pUserData: %d", pUserData)

	response := v2.TtsRes{
		ErrorCode: int32(flag),
		Status:    3,
	}
	sendResp(object, response)
	close(object.BackChan)
	object.Log.Infof("end to goOnEnd pUserData: %d", pUserData)
}

//export goOnDebug
func goOnDebug(pUserData unsafe.Pointer, debugtype *C.char, info *C.char) {

	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnDebug; irregularity type")
		return
	}
	object.Log.Infof("start to goOnDebug; pUserData: %d", pUserData)

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
	object.Log.Infof("end to goOnDebug pUserData: %d", pUserData)
}

//export goOnTimedMouthShape
func goOnTimedMouthShape(pUserData unsafe.Pointer, mouth *C.TimedMouthShape, size C.int, startTimeMs C.float) {

	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnTimedMouthShape; irregularity type")
		return
	}
	object.Log.Infof("start to goOnTimedMouthShape;pUserData: %d; size:%d", pUserData, size)

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
	object.Log.Infof("end to goOnTimedMouthShape pUserData: %d", pUserData)
}

//export goOnActionElement
func goOnActionElement(pUserData unsafe.Pointer, actionType C.int, url *C.char, operationType C.int, coordinate unsafe.Pointer, renderDuration C.int) {
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnActionElement; irregularity type")
		return
	}
	object.Log.Infof("start to goOnActionElement;pUserData: %d;actionType:%d, operationType:%d, renderDuration:%d", pUserData, actionType, operationType, renderDuration)

	coordinateC := *(*C.Coordinate)(coordinate)
	response := v2.TtsRes{
		Status: 2,
		ResultOneof: &v2.TtsRes_ActionElement{
			ActionElement: &v2.ActionElement{
				ActionType:    int32(actionType),
				Url:           C.GoString(url),
				OperationType: int32(operationType),
				Coordinate: &v2.Coordinate{
					Off:   int32(coordinateC.off_utf8),
					Len:   int32(coordinateC.len_utf8),
					Order: int32(coordinateC.order),
				},
				RenderDuration: int32(renderDuration),
			},
		},
	}
	sendResp(object, response)
	object.Log.Infof("end to goOnActionElement pUserData: %d", pUserData)
}

//export goOnSynthesizedData
func goOnSynthesizedData(pUserData unsafe.Pointer, audioData *C.SynthesizedAudio, coordinate *C.Coordinate) {
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnSynthesizedData; irregularity type")
		return
	}
	object.Log.Infof("start to goOnSynthesizedData;pUserData:%d", pUserData)

	length := (C.int)(audioData.audio_size)
	wav := C.GoBytes(unsafe.Pointer(audioData.audio_data), length)

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
	object.Log.Infof("end to goOnSynthesizedData pUserData: %d", pUserData)
}

//export goOnFacialExpression
func goOnFacialExpression(pUserData unsafe.Pointer, facialExpressionData *C.FacialExpressionSegment) {
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnFacialExpression; irregularity type")
		return
	}
	object.Log.Infof("start to goOnFacialExpression;pUserData:%d", pUserData)

	var framedDim uint64
	if fefd, ok := (object.ParamMap["FacialExpressionFrameDim"]).(int32); ok {
		framedDim = uint64(fefd)
	}
	var expression *v2.Expression
	if facialExpressionData != nil {
		frameSize := uint64(facialExpressionData.frameSize)
		startTimeMs := float32(facialExpressionData.startTimeMs)
		size := int32(frameSize * framedDim)
		object.Log.Infof("start to goOnFacialExpression;pUserData:%d; size:%d", pUserData, size)
		var expressionData = make([]float32, size)
		for i := 0; i < int(size); i++ {
			data := *(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(facialExpressionData.data)) + uintptr(C.sizeof_float*C.int(i))))
			expressionData[i] = float32(data)
		}
		expression = &v2.Expression{
			Data:        expressionData,
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
	object.Log.Infof("end to goOnFacialExpression pUserData: %d", pUserData)
}

//export goOnBodyMovement
func goOnBodyMovement(pUserData unsafe.Pointer, bodyMovementData unsafe.Pointer) {
	object := getHandlerObject(pUserData)
	if object == nil {
		log.Println("goOnBodyMovement; irregularity type")
		return
	}
	object.Log.Infof("start to goOnBodyMovement;pUserData:%d", pUserData)

	var bodyMovement *v2.BodyMovement
	bodyMovementDataC := (*C.BodyMovementSegment)(bodyMovementData)
	// _, bodyMovementDim := getTmpDecoder(*id)
	var bodyMovementFrameSizeFrameDim uint64
	if bodyMovementDim, ok := (object.ParamMap["BodyMovementFrameDim"]).(int32); ok {
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
	object.Log.Infof("end to goOnBodyMovement pUserData: %d", pUserData)
}

func sendResp(object *data.HandlerObjectV2, response v2.TtsRes) {

	if object != nil && object.BackChan != nil {
		object.BackChan <- response
	}
}

func getHandlerObject(pUserData unsafe.Pointer) *data.HandlerObjectV2 {
	handlerObject := pointer.Load(int32(uintptr(pUserData)))
	if handlerObject == nil {
		log.Println("don't find to handler object;pUserData:", pUserData)
		return nil
	}
	object, ok := handlerObject.(*data.HandlerObjectV2)
	if !ok {
		log.Println("irregularity handler object;pUserData:", pUserData)
		return nil
	}
	return object
}
