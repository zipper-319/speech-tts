package main

/*
#include <stdlib.h>

// -resType 0--pron_res_cn 1--pron_res_en 2--norm_res
// +resType 0--pron_res_cn 1--pron_res_en 2--string_res_cn 3--string_res_en 4--norm_res_cn 5--norm_res_en
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
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"speech-tts/export/service"
	ttsData "speech-tts/export/service/proto"
	"speech-tts/internal/pkg/util"
	"unsafe"
)

const port = ":8080"

var serviceName = "speech-tts"

//export ResService_Init
func ResService_Init(cb *C.ResService_Callback, pUserData unsafe.Pointer) C.int {
	ctx := context.Background()
	fn := Callback(cb, pUserData)
	go func() {
		r := gin.Default()
		r.POST("/resource/update", ReLoadTTSResource(fn))
		if err := r.Run(port); err != nil {
			log.Error(err)
		}
	}()

	if err := service.InitTTSResource(ctx, fn); err != nil {
		log.Error(err)
	}
	// 注册
	if err := service.RegisterResService(ctx, serviceName, util.GetHostIp()+port); err != nil {
		log.Error(err)
	}

	return C.int(0)
}

//export EndInit
func EndInit() {
	if err := service.UnRegisterResService(context.Background(), serviceName, util.GetHostIp()+port); err != nil {
		log.Error(err)
	}
}

func main() {

}

type UpdateResourceReq struct {
	ResType  ttsData.ResType
	Language ttsData.LanguageType
}

func ReLoadTTSResource(callback service.CallbackFn) gin.HandlerFunc {
	return func(g *gin.Context) {
		var req UpdateResourceReq
		if err := g.BindJSON(&req); err != nil {
			log.Error(err)
			return
		}
		fileName, err := service.GetTTSResAndSave(context.Background(), req.ResType, req.Language)
		if err != nil {
			log.Error(err)
		}
		if int(req.ResType) < 3 {
			callback(req.ResType, req.Language, fileName)
		}
		return
	}
}

func Callback(cb *C.ResService_Callback, pUserData unsafe.Pointer) service.CallbackFn {
	return func(resType ttsData.ResType, languageType ttsData.LanguageType, fileName string) {
		fileNameC := C.CString(fileName)
		defer C.free(unsafe.Pointer(fileNameC))
		C.bridge_event_cb(cb, C.int(2*int(resType)+int(languageType)), fileNameC, pUserData)
	}
}
