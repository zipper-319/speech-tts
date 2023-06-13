#ifndef ANIMATION_H
#define ANIMATION_H

#include <AnimationDef.h>

#ifdef __cplusplus
extern "C"
{
#endif

/**
 * 设置资源路径并初始化
 * @param res_path 资源路径
 */
void Animation_InitRes(const char* res_path);

/**
* 反初始化
*/
void Animation_DeinitRes();

/**
* 获取所支持的角色
* @param animationType 用于控制获取面部或身体角色
* @param supportedValues 用于输出角色的取值数组，不要修改其内容！！！
* @param element_num 用于输出取值数组的大小
* @return 0--成功,&lt;0--失败
*/
int Animation_GetSupportedRoles(enum ComponentType animationType, const char*** supportedValues, unsigned int *element_num);

/**
* 获取角色的配置
* @param animationType 指定获取面部或身体配置
* @param face role 用于指定具体角色，支持列表由Animation_GetSupportedRoles获取
* @param metadata 用于输出配置信息
* @return 0--成功,&lt;0--失败
*/
int Animation_GetRolesMetadata(enum ComponentType animationType, const char* role, AccompanyDataConfig* metadata);

/**
* 创建预测器
*/
void* Animation_CreatePredictor();

/**
* 释放预测器
* @param predictor 预测器
*/
void Animation_ReleasePredictor(void *predictor);

/**
* 配置预测器，请确保在一次完整请求的最开始时调用一次
* @param predictor 待配置的预测器
* @param face 面部设置参数
* @param body 动作设置参数
* @return Ready_Face--表情配置成功;Ready_Body--动作配置成功;Ready_Face|Ready_Body--表情、动作配置成功;Ready_Not--配置失败, 
*/
int Animation_ConfigurePredictor(void *predictor, ComponentConfig *face, ComponentConfig *body);

/**
* 预测器预测
* @param predictor 预测器
* @param input 输入数据
* @param output 输出数据，不用释放
* @return Ready_Face--表情预测成功;Ready_Body--动作预测成功;Ready_Face|Ready_Body--表情、动作预测成功;Ready_Not--预测失败, 
*/
int Animation_Predict(void* predictor, AnimationInput* input, AnimationOutput* output);
/**
* 获取版本信息
* @return 版本信息的字符串，不应修改，无需释放
*/
const char* Animation_GetVersion();

#ifdef TTS //以下接口目前仅供TTS出于兼容性目的使用
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
