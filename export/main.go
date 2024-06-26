package main

/*
#include <stdlib.h>

// -resType 0--pron_res_cn 1--pron_res_en 2--norm_res
// +resType 0--pron_res_cn 1--pron_res_en 2--string_res_cn 3--string_res_en 4--norm_res_cn 5--norm_res_en  16 Rhythm 18 Participle 30模型文件
typedef struct
{
    void (*reloadRes)(void* pUserData, int resType, char* resPath);
}ResService_Callback;
static void bridge_event_cb(ResService_Callback* cb, int resType, char* resPath,void* pUserData)
{
	cb->reloadRes(pUserData,resType,resPath);
}
*/
import "C"
import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	ttsData "speech-tts/export/service/proto"
	"speech-tts/export/service/resource"
	"speech-tts/internal/pkg/util"
	"unsafe"
)

const port = ":8080"
const serviceName = "speech-tts"

//export ResService_Init
func ResService_Init(cb *C.ResService_Callback, pUserData unsafe.Pointer) C.int {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	ctx := context.Background()
	fn := Callback(cb, pUserData)
	go func() {
		r := gin.Default()
		r.POST("/resource/update", ReLoadTTSResource(fn))
		if err := r.Run(port); err != nil {
			log.Println(err)
		}
	}()

	if err := resource.InitTTSResource(ctx, fn); err != nil {
		log.Println(err)
		return C.int(-1)
	}
	// 注册
	if err := resource.RegisterResService(ctx, serviceName, fmt.Sprintf("http://%s%s/resource/update", util.GetHostIp(), port)); err != nil {
		log.Println(err)
		return C.int(-1)
	}
	resource.SaveResVersion()
	return C.int(0)
}

//export EndInit
func EndInit() C.int {
	if err := resource.UnRegisterResService(context.Background(), serviceName, util.GetHostIp()+port); err != nil {
		log.Println(err)
		return C.int(-1)
	}
	return C.int(0)
}

//export ResService_GetVersion
func ResService_GetVersion() *C.char {
	v, _ := json.Marshal(resource.ResVersionMap)
	d := md5.Sum(v)
	version := hex.EncodeToString(d[:])
	return C.CString(version)
}

func main() {

}

type UpdateResourceReq struct {
	ResType  ttsData.ResType ``
	Language ttsData.LanguageType
	DataMap  map[string]string
}

func ReLoadTTSResource(callback resource.CallbackFn) gin.HandlerFunc {
	return func(g *gin.Context) {
		var req UpdateResourceReq
		if err := g.ShouldBindJSON(&req); err != nil {
			log.Println(err)
			return
		}
		log.Printf("reload resource,DataMap:%#v resType:%d, language:%d", req.DataMap, req.ResType, req.Language)

		if int(req.ResType) < int(ttsData.ResType_Model) {
			fileName, err := resource.SaveResource(req.DataMap, req.ResType, req.Language)
			if err != nil {
				log.Println(err)
			}
			log.Printf("reload resource success, fileName:%s", fileName)
			callback(req.ResType, req.Language, fileName)
		} else if req.ResType == ttsData.ResType_Model {
			speakerName := req.DataMap["speaker_name"]
			speakerOwner := req.DataMap["speaker_owner"]
			modelUrl := req.DataMap["model_url"]
			dstPath, err := resource.SaveSpeakerModel(modelUrl, speakerOwner, speakerName)
			if err != nil {
				log.Printf("fail to SaveSpeakerModel; modelUrl:%s, err:%v",modelUrl, err)
				return
			}
			log.Printf("save speaker model success,speakerName:%s,speakerOwner:%s,modelUrl:%s, path: %s", speakerName, speakerOwner, modelUrl, dstPath)
			callback(ttsData.ResType_Model, ttsData.LanguageType_Chinese, dstPath)
		}
		resource.SaveResVersion()
		g.JSON(200, "success")
	}
}

func Callback(cb *C.ResService_Callback, pUserData unsafe.Pointer) resource.CallbackFn {
	return func(resType ttsData.ResType, languageType ttsData.LanguageType, fileName string) {
		fileNameC := C.CString(fileName)
		defer C.free(unsafe.Pointer(fileNameC))
		C.bridge_event_cb(cb, C.int(2*int(resType)+int(languageType)), fileNameC, pUserData)
	}
}
