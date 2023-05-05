#ifndef ANIMATION_H
#define ANIMATION_H

#include <AnimationDef.h>

#ifdef __cplusplus
extern "C"
{
#endif

/**
* 数据的特性集
*/
enum {
    /**
    * 数据开始
    */
    DATA_START = 1<<0,
    /**
    * 数据结束
    */
    DATA_END = 1<<1,
};
/**
* 合成类型
*/
enum SynthetiseType{
    SynthetiseType_Min = 0,

    /**
    * 不合成
    */
    SynthetiseType_Disable = 0,
    /**
    * 基于发音（拼音或音标）合成
    */
    SynthetiseType_Use_Pron = 1,
    /**
    * 基于语音合成
    */
    SynthetiseType_Use_Audio = 2,

    /**
    新定义的类型添加于此
    */

    SynthetiseType_Max,
};
/**
 * 输入数据
 */
typedef struct {
    /**
    * 16kHz 16bit mono PCM数据
    */
    const short* data;
    /**
    * 数据short数
    */
    unsigned int dataShortSize;
    /**
    * 仅USE_PRON时使用（通常是TTS使用），其他情况应置NULL
    */
    void* pron;
    /**
    * 数据特性
    */
    unsigned int dataFlags;

} AnimationInput;
/**
 * 输出数据
 */
typedef struct {
    /**
    * 表情数据
    */
    FacialExpressionSegment* face;
    unsigned int faceSize;
    /**
    * 动作数据
    */
    BodyMovementSegment* body;
    unsigned int bodySize;
} AnimationOutput;

/**
 * 设置资源路径并初始化
 * @param res_path 资源路径
 */
void Animation_InitRes(const char* res_path);
/**
* 获取所支持的角色
* @param supportedValues 用于输出角色的取值数组，不要修改其内容！！！
* @param element_num 用于输出取值数组的大小
* @return 0--成功,&lt;0--失败
*/
int Animation_GetSupportedRoles(const char*** supportedValues, unsigned int *element_num);
/**
* 获取角色的配置
* @param role_id 指定角色
* @param face 用于输出表情配置
* @param body 用于输出动作配置
* @return 0--成功,&lt;0--失败
*/
int Animation_GetRoleConfig(const char* role_id, FacialExpressionConfig* face, BodyMovementConfig* body);

/**
* 反初始化
*/
void Animation_DeinitRes();


/**
* 创建预测器
*/
void* Animation_CreatePredictor();
/**
* 配置预测器,首次预测前或者需要改变配置时（改变配置只应在完整句子的边界处）调用
* @param predictor 待配置的预测器
* @param role_id 角色
* @return 0--成功,&lt;0--失败
*/
int Animation_ConfigurePredictor(void *predictor, const char* role_id);

/**
* 释放预测器
* @param predictor 预测器
*/
void Animation_ReleasePredictor(void *predictor);

/**
* 重置预测器状态，请确保在一次完整请求的最开始（或结束）时调用一次
* @param predictor 预测器
*/
void Animation_ResetPredictor(void *predictor);
/**
* 预测器预测
* @param predictor 预测器
* @param input 输入数据
* @param output 输出数据，不用释放
* @return 0--成功,&lt;0--失败
*/
int Animation_Predict(void* predictor, AnimationInput* input, AnimationOutput* output);
/**
* 获取版本信息
* @return 版本信息的字符串，不应修改，无需释放
*/
const char* Animation_GetVersion();

#ifdef TTS //以下接口目前仅供TTS出于兼容性目的使用
/**
* 配置预测器,首次预测前或者需要改变配置时（改变配置只应在完整句子的边界处）调用
* @param predictor 待配置的预测器
* @param role_id 角色
* @param faceType 表情的合成类型,参见SynthetiseType
* @param bodyType 动作的合成类型,参见SynthetiseType
* @return 0--成功,&lt;0--失败
*/
int Animation_ConfigurePredictorForTts(void *predictor, const char* role_id, int faceType, int bodyType);
/**
* 预测器预测v1接口
* @param predictor 预测器
* @param input 输入数据
* @param face 表情数据，不用释放
* @return 0--成功,&lt;0--失败
*/
int Animation_Predict_V1(void* p, AnimationInput* input, FacialExpression** face, int* faceSize);
#endif //TTS

#ifdef __cplusplus
}
#endif

#endif //ANIMATION_H
