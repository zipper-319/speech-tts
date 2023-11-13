package resource

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	ttsData "speech-tts/export/service/proto"
	"speech-tts/internal/pkg/util"
	"time"
)

const (
	AddrDefault        = "10.12.32.198"
	HttpPortDefault    = "8001"
	GrpcPortDefault    = "9001"
	ResPath            = "./res/read_and_speak/speak"
	ResourceSplit      = "@@"
	PronounceSplit     = ":"
	TmpPath            = "./tmp"
	ResVersionFileName = "./res/res_version.json"
)

var ResVersionMap map[string]string

type CallbackFn func(ttsData.ResType, ttsData.LanguageType, string)

var (
	HttpUrl    string
	GrpcUrl    string
	IsOpenGrpc bool
)

func init() {
	addr := os.Getenv("dataServiceAddr")
	if addr == "" {
		addr = AddrDefault
	}
	httpPort := os.Getenv("dataServiceHttpPort")
	if httpPort == "" {
		httpPort = HttpPortDefault
	}
	grpcPort := os.Getenv("dataServiceGrpcPort")
	if grpcPort == "" {
		grpcPort = GrpcPortDefault
	}
	if os.Getenv("IsOpenGrpc") != "" {
		IsOpenGrpc = true
	}
	ResVersionMap = make(map[string]string, len(ttsData.ResType_name)*2)
	_, err := os.Stat(ResVersionFileName)
	if err == nil {
		content, err := os.ReadFile(ResVersionFileName)
		if err != nil {
			panic(err)
		}
		if len(content) > 0 {
			json.Unmarshal(content, &ResVersionMap)
		}
	}
	log.Infof("resource version map: %v", ResVersionMap)

	HttpUrl = fmt.Sprintf("%s:%s", addr, httpPort)
	GrpcUrl = fmt.Sprintf("%s:%s", addr, grpcPort)
	log.Info("dataServiceAddr:", addr, " dataServiceHttpPort:", httpPort, " dataServiceGrpcPort:", grpcPort)
}

func RegisterResService(ctx context.Context, serviceName, callbackUrl string) error {
	log.Infof("RegisterResService: serviceName: %s, callbackUrl: %s; IsOpenGrpc:%t", serviceName, callbackUrl, IsOpenGrpc)
	if IsOpenGrpc {
		return RegisterResServiceByGrpc(ctx, serviceName, callbackUrl)
	} else {
		return RegisterResServiceByHttp(serviceName, callbackUrl)
	}
}

func UnRegisterResService(ctx context.Context, serviceName, callbackUrl string) error {
	if IsOpenGrpc {
		return UnRegisterResServiceByGrpc(ctx, serviceName, callbackUrl)
	} else {
		return UnregisterResServiceByHttp(serviceName, callbackUrl)
	}
}

func GetTTSResAndSave(ctx context.Context, resType ttsData.ResType, languageType ttsData.LanguageType) (string, error) {
	var (
		dataList []*ttsData.GetTtsDataResponse_TTSData
		err      error
	)
	dataList, err = GetTTSResByGrpc(ctx, resType, languageType)
	if err != nil {
		return "", err
	}
	return SaveResource(dataList, resType, languageType)
}

func InitTTSResource(ctx context.Context, fn CallbackFn) error {

	for v, _ := range ttsData.ResType_name {
		resType := ttsData.ResType(v)
		if int(v) >= int(ttsData.ResType_Model) {

			if err := GetSpeakerModel(ctx, fn); err != nil {
				log.Errorf("GetSpeakerModel error:%v", err)
				continue
			}
			log.Info("GetSpeakerModel success")

		} else {
			if int(v) >= int(ttsData.ResType_Rhythm) {
				fileName, err := GetTTSResAndSave(ctx, resType, ttsData.LanguageType_Chinese)
				if err != nil {
					log.Errorf("Save tts resource; resourceType:%s,language:%s, error:%v", resType, ttsData.LanguageType_Chinese, err)
					continue
				}
				fn(resType, ttsData.LanguageType_Chinese, fileName)
			} else {
				for lang, _ := range ttsData.LanguageType_name {
					languageType := ttsData.LanguageType(lang)
					fileName, err := GetTTSResAndSave(ctx, resType, languageType)
					if err != nil {
						log.Errorf("Save tts resource; resourceType:%s,language:%s, error:%v", resType, languageType, err)
						continue
					}
					fn(resType, languageType, fileName)
				}
			}

		}

	}
	return nil
}

func SaveResource(dataList []*ttsData.GetTtsDataResponse_TTSData, resType ttsData.ResType, languageType ttsData.LanguageType) (string, error) {
	fileName := fmt.Sprintf("%s/%s_%s.txt", ResPath, resType.String(), languageType.String())
	dataListByte, _ := json.Marshal(dataList)
	d := md5.Sum(dataListByte)
	version := hex.EncodeToString(d[:])
	if ResVersionMap[fileName] == version {
		return fileName, nil
	}
	ResVersionMap[fileName] = version

	os.Rename(fileName, fileName+".bak"+time.Now().Format("20060102150405"))
	split := ResourceSplit
	if resType == ttsData.ResType_Pronounce {
		split = PronounceSplit
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var n int
	for _, v := range dataList {
		i, _ := f.WriteString(fmt.Sprintf("%s%s%s\n", v.Key, split, v.Value))
		n += i
	}
	log.Infof("write file:%s, length:%d", fileName, n)
	return fileName, nil
}

func GetSpeakerModel(ctx context.Context, fn CallbackFn) error {
	var (
		err       error
		modelList []*ttsData.GetSpeakerModelResult_SpeakerModel
	)
	if IsOpenGrpc {
		modelList, err = GetSpeakerModelByGrpc(ctx)
	} else {
		modelList, err = GetSpeakerModelByHttp()
	}
	if err != nil {
		return err
	}

	for _, speakerModel := range modelList {

		dstPath, err := SaveSpeakerModel(speakerModel.ModelUrl, speakerModel.SpeakerOwner, speakerModel.SpeakerName)
		if err != nil {
			return err
		}
		fn(ttsData.ResType_Model, ttsData.LanguageType_Chinese, dstPath)
	}
	return nil
}

func TransForm(dataMap map[string]string) []*ttsData.GetTtsDataResponse_TTSData {
	result := make([]*ttsData.GetTtsDataResponse_TTSData, 0, len(dataMap))
	for k, v := range dataMap {
		result = append(result, &ttsData.GetTtsDataResponse_TTSData{
			Key:   k,
			Value: v,
		})
	}
	return result
}

func SaveSpeakerModel(modelUrl, speakerOwner, speakerName string) (string, error) {
	if err := os.MkdirAll(TmpPath, 0777); err != nil {
		return "", err
	}
	tmpFile := fmt.Sprintf("%s/%s_%s.zip", TmpPath, speakerOwner, speakerName)
	if err := util.DownloadFile(modelUrl, tmpFile); err != nil {
		return "", err
	}
	dstPath := fmt.Sprintf("%s/%s/%s", ResPath, speakerOwner, speakerName)
	if err := util.DeCompressToPath(tmpFile, dstPath); err != nil {
		return "", err
	}
	return dstPath, nil
}

func SaveResVersion() error {
	v, _ := json.Marshal(ResVersionMap)
	f, err := os.OpenFile(ResVersionFileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(v))
	return err
}
