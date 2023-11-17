#ifndef ANIMATION_DEF_H
#define ANIMATION_DEF_H

typedef char bool;

#ifdef __cplusplus
extern "C"
{
#endif

/**
 * TTS表情(v1版)
 */
typedef struct
{
    /**
     * 表情数据，frame_dim*frame_size个float
     */
    float* data;
    /**
     * 帧数
     */
    unsigned int frame_size;
    /**
     * 每帧的维度
     */
    unsigned int frame_dim;
    /**
     * 每帧的时长
     */
    float  frame_time;
}FacialExpression;


/**
 * TTS伴生数据
 */
typedef struct
{
    /**
     * 伴生数据，frameSize*frameDim个float
     */
    float* data;
    /**
     * 帧数
     */
    unsigned int frameSize;
    /**
     * 起始时间，单位ms
     */
    float  startTimeMs;
}AccompanyData;

/**
 * 伴生数据的配置信息
 */
typedef struct
{
    /**
     * 每帧的维度，即一帧由frameDim个float组成
     */
    unsigned int frameDim;
    /**
     * 每帧的持续时长
     */
    float frameDurMs;
    /**
     * 伴生数据的其他描述性信息(json)
     */
    const char *meta_data;
}AccompanyDataConfig;

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
* 合成目标类型
*/
enum ComponentType{
    ComponentType_Face = 0,
    ComponentType_Body = 1,
};

/**
* 合成输入数据类型
*/
enum InputType{
    InputType_Pron = 0,
    InputType_Audio = 1,
};

enum ResultStatus{
    Ready_Not = 0,
    Ready_Face = 1<<0,
    Ready_Body = 1<<1,
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
    AccompanyData* face;
    unsigned int faceSize;
    /**
    * 动作数据
    */
    AccompanyData* body;
    unsigned int bodySize;
} AnimationOutput;


/**
 * Configure Predictor 配置入参信息
 */
typedef struct ComponentConfig {
    /**
     * 是否预测animation
     */
    bool enable;
    /**
     * 预测animation输入类型文本或音频
     */
    enum InputType input_type;//Pron   audio
    /**
     * 预测animation角色信息， 支持列表由Animation_GetSupportedRoles获取
     */
    const char* expect_role;
} ComponentConfig;

/**
 * 表情数据段，参考AccompanyData
 */
typedef AccompanyData FacialExpressionSegment;
/**
 * 动作数据段，参考AccompanyData
 */
typedef AccompanyData BodyMovementSegment;
/**
 * 表情配置信息，参考AccompanyDataConfig
 */
typedef AccompanyDataConfig FacialExpressionConfig;
/**
 * 动作配置信息，参考AccompanyDataConfig
 */
typedef AccompanyDataConfig BodyMovementConfig;

#ifdef __cplusplus
}
#endif

#endif //ANIMATION_DEF_H
