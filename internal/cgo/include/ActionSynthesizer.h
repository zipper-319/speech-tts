#ifndef ACTION_SYNTHESIZER_H
#define ACTION_SYNTHESIZER_H

#include <MouthShape.h>
#include <AnimationDef.h>
#include <TtsSetting.h>
#include <ActionDef.h>

#ifdef __cplusplus
extern "C"
{
#endif

/**
* 文本坐标
*/
typedef struct
{
    /**
    * 文本起点（以utf8字符为单位）
    */
    unsigned int off_utf8;
    /**
    * 文本长度
    */
    unsigned int len_utf8;
    /**
    * 当len_utf8&eq;0时，与off_utf8一起描述一个的时间点。具体地，当off_utf8相同时，order相同则认为同时，order不同时，0优先，1次之，依次类推
    * 当len_utf8&gt;0时，order无意义
    */
    int order;
}Coordinate;

/**
* 时间坐标
*/
typedef struct
{
    /**
    * 时间起点
    */
    unsigned int startMs;
    /**
    * 时长
    */
    unsigned int durMs;
}TimeCoordinate;

/**
* 数据的特性集
*/
enum {
    /**
    * 是标点
    */
    IS_PUNCTUATION = 1<<0,
};

/**
 * 合成的音频
 */
typedef struct
{
    /**
    * 音频数据
    */
    char* audio_data;
    /**
    * 字节长度
    */
    unsigned int audio_size;
    /**
    * 数据的特性
    */
    unsigned int flags;
}SynthesizedAudio;

/**
 * 音频config
 */
typedef struct
{
    /**
    * 采样率，16000Hz
    */
    unsigned int sampling_rate;
    /**
    * 通道数, 1--mono
    */
    unsigned int channels;
    /**
    * 音频编码, 参考enum AudioEncoding
    */
    unsigned int audio_encoding;
}AudioConfig;

/**
 * Action回调
 */
typedef struct
{
    /**
    * 语音合成开始
    * @param pUserData 应用层数据
    * @param ttsText 用于TTS合成的文本
    * @param audioConfig 音频配置
    * @param expressionConfig 表情配置
    * @param movementConfig 动作配置
    */
    void (*onStart)(void* pUserData, const char* ttsText, AudioConfig* audioConfig, FacialExpressionConfig* expressionConfig, BodyMovementConfig* movementConfig);

    /**
    * 合成的音频数据（以文本元素为单位输出）
    * @param pUserData 应用层数据
    * @param data SynthesizedAudio
    * @param coordinate 该数据的坐标信息
    */
    void (*onSynthesizedData)(void* pUserData, SynthesizedAudio* data, Coordinate* coordinate);

    /**
    * 坐标信息（时间坐标与文本坐标的映射）
    * @param pUserData 应用层数据
    * @param coordinate 文本坐标
    * @param timeCoordinate 时间坐标
    */
    void (*onCoordinate)(void* pUserData, Coordinate* coordinate, TimeCoordinate* timeCoordinate);

    /**
    * 编码后的音频数据
    * onEncodedData和onCoordinate一起，构成对onSynthesizedData的一个替代，一次请求中onSynthesizedData和onEncodedData只会有一个被回调
    * @param pUserData 应用层数据
    * @param data SynthesizedAudio
    */
    void (*onEncodedData)(void* pUserData, SynthesizedAudio* data);

    /**
    * 语音合成结束
    * @param pUserData 应用层数据
    * @param flags 0--正常合成完成，1--因cancel而提前终止，2--因异常而提前终止
    */
    void (*onEnd)(void* pUserData, int flags);

    /**
    * 输出调试性信息(for debug)
    * @param pUserData 应用层数据
    * @param type 信息的类型
    * @param info 调试性信息
    */
    void (*onDebug)(void* pUserData, const char* type, const char *info);
    /**
    * 时间化的口型
    * @param pUserData 应用层数据
    * @param mouth TimedMouthShape数组，注意：该对象仅在本函数执行期间有效，使用者无需考虑释放
    * @param size TimedMouthShape数组的大小
    * @param startTimeMs 该段口型的起始时间，单位ms
    */
    void (*onTimedMouthShape)(void* pUserData, TimedMouthShape* mouth, int size, float startTimeMs);

    /**
    * 表情数据
    * @param pUserData 应用层数据
    * @param expression FacialExpressionSegment
    */
    void (*onFacialExpression)(void* pUserData, FacialExpressionSegment* expression);

    /**
    * 动作数据
    * @param pUserData 应用层数据
    * @param expression BodyMovementSegment
    */
    void (*onBodyMovement)(void* pUserData, BodyMovementSegment* movement);

    /**
    * Action基元数据
    * @param pUserData 应用层数据
    * @param type action基元类型，-100&lt;type&lt;100
    * @param url 基元数据的url
    * @param operation_type action操作类型，-100&lt;operation_type&lt;100
    * @param coordinate 该数据的坐标信息
    * @param render_duration 渲染时长（该值不应大于文件时长），单位ms，-1代表持续到指定文件结束，-2代表由coordinate的len_utf8部分决定时长
    */
    void (*onActionElement)(void* pUserData, int type, const char* url, int operation_type, Coordinate* coordinate, int render_duration);
}ActionCallback;



/**
* ActionSynthesizer初始化
* @param 资源目录的路径
*/
void ActionSynthesizer_Init(const char* resDir);
/**
* ActionSynthesizer反初始化
*/
void ActionSynthesizer_Deinit();


/**
* 合成Action
* @param actionDescriptor action描述，json格式
* @param setting TtsSetting
* @param cb 回调
* @param pUserData 应用层数据
* @param traceId 用于跟踪该条请求，由应用层保证其唯一性
* @return &ge;0--当前合成请求的id, &lt;0--失败（-1--TtsSetting非法 -2--actionDescriptor非法 -3--当前服务繁忙）
*/
int ActionSynthesizer_SynthesizeAction(const char* actionDescriptor, TtsSetting* setting, ActionCallback* cb, void* pUserData, const char* traceId);
/**
* 终止指定的Action合成（非阻塞版，函数返回时该请求可能尚未终止）
* @param id 要终止的Action合成请求的id
*/
void ActionSynthesizer_CancelAction(int id);


/**
* 获取ActionSynthesizer的版本信息
* @return 版本信息的字符串，不应修改，无需释放
*/
const char* ActionSynthesizer_GetVersion();
/**
* 获取ActionSynthesizer资源服务的版本信息
* @return 版本信息的字符串，不应修改，无需释放
*/
const char* ActionSynthesizer_GetResServiceVersion();


#if COMPATIBLE_V1

/**
 * TTS回调
 */
typedef struct
{
    /**
    * 语音合成开始（即首包数据已准备好）
    * @param pUserData 应用层数据
    */
    void (*onStart)(void* pUserData);

    /**
    * 合成的音频数据
    * @param pUserData 应用层数据
    * @param data 音频数据（16kHz 16bit mono PCM数据）
    * @param size 该笔音频数据的字节数
    */
    void (*onAudio)(void* pUserData, const char *data, int size);

    /**
    * 语音合成结束
    * @param pUserData 应用层数据
    * @param flags 0--正常合成完成，1--因cancel而提前终止，2--因异常而提前终止
    */
    void (*onEnd)(void* pUserData, int flags);

    /**
    * for debug
    * @param pUserData 应用层数据
    * @param info 调试性信息
    */
    void (*onDebug)(void* pUserData, const char *info);
    /**
    * 时间化的口型
    * @param pUserData 应用层数据
    * @param mouth TimedMouthShape数组，注意：该对象仅在本函数执行期间有效，使用者无需考虑释放
    * @param size TimedMouthShape数组的大小
    * @param subText 口型对应的正则后的文本，它的颗粒度小于onCurTextSegment函数的参数normalizedText
    */
    void (*onTimedMouthShape)(void* pUserData, TimedMouthShape* mouth, int size, const char* subText);

    /**
    * 当前正在合成的文本片段
    * @param pUserData 应用层数据
    * @param normalizedText 当前正在合成的正则后的文本片段，使用者无需考虑释放
    * @param originalText 当前正在合成的正则前的文本片段，使用者无需考虑释放
    */
    void (*onCurTextSegment)(void* pUserData, const char* normalizedText, const char* originalText);

    /**
    * 表情数据
    * @param pUserData 应用层数据
    * @param expression FacialExpression
    */
    void (*onFacialExpression)(void* pUserData, FacialExpression* expression);

}TTS_Callback;


/**
* 合成Action
* @param actionDescriptor action描述，json格式
* @param setting TtsSetting
* @param cb 回调
* @param pUserData 应用层数据
* @param traceId 用于跟踪该条请求，由应用层保证其唯一性
* @return &ge;0--当前合成请求的id, &lt;0--失败（-1--TtsSetting非法 -2--actionDescriptor非法 -3--当前服务繁忙）
*/
int ActionSynthesizer_SynthesizeAction_V1(const char* tts_text, TtsSetting* setting, TTS_Callback* cb, void* pUserData, const char* traceId);


#endif

#ifdef __cplusplus
}
#endif

#endif //ACTION_SYNTHESIZER_H
