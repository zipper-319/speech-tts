#ifndef ANIMATION_DEF_H
#define ANIMATION_DEF_H

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
}AccompanyDataConfig;

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
