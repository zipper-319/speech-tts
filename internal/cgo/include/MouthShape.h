#ifndef MOUTH_SHAPE_H
#define MOUTH_SHAPE_H

#ifdef __cplusplus
extern "C"
{
#endif

/**
口型值枚举
*/
typedef enum{
    MouthShape_ZeroInitialConsonant  = -1, //zero initial consonant
//
    MouthShape_sil  = 0,
    MouthShape_PP = 1,
    MouthShape_FF = 2,
    MouthShape_TH = 3,
    MouthShape_DD = 4,
    MouthShape_kk = 5,
    MouthShape_CH = 6,
    MouthShape_SS = 7,
    MouthShape_nn = 8,
    MouthShape_RR = 9,
    MouthShape_aa = 10,
    MouthShape_E = 11,
    MouthShape_I = 12,
    MouthShape_O = 13,
    MouthShape_U = 14,
    /**
    新定义的口型添加于此--<
    */
    MouthShape_L = MouthShape_nn,
    MouthShape_JQX = MouthShape_I,
    MouthShape_ZCS = MouthShape_SS,
    MouthShape_ao = MouthShape_aa,
    MouthShape_o = MouthShape_U,
    MouthShape_ou = MouthShape_O,
    MouthShape_e = MouthShape_nn,
    MouthShape_ei = MouthShape_I,
    MouthShape_er = MouthShape_aa,
    MouthShape_v = MouthShape_U,
    /**
    新定义的口型添加于此-->
    */
    MOUTHSHAPE_TOTAL_SIZE = 15,
} MouthShape;

/**
时间化口型
*/
typedef struct
{
    /**
    持续时长
    */
    unsigned long long durationUs;
    /**
    口型值
    */
    MouthShape mouth;
}TimedMouthShape;

/**
加权口型
*/
typedef struct
{
    /**
    持续时长
    */
    unsigned long long durationUs;
    /**
    口型权重
    */
    float weight[MOUTHSHAPE_TOTAL_SIZE];
}MouthShapeWeight;
/**
* 口型平滑。//注意：该函数可能改变入参mouth中元素的值，所以若需保持原值，请传入mouth的一份深拷贝。
* @param mouth TimedMouthShape数组，该内存的释放由调用者负责
* @param size TimedMouthShape数组的元素个数
* @param granularityUs 输出颗粒度（单位us），即每granularityUs输出一个MouthShapeWeight元素。目前未使用
* @param weight 输出的MouthShapeWeight数组，调用者需要释放（free）该内存。
* @param size1 MouthShapeWeight数组的元素个数
* @return 0--成功,&lt;0--失败
*/
int SmoothTimedMouthShape(const TimedMouthShape* mouth, unsigned size, unsigned granularityUs, MouthShapeWeight** weight, unsigned* size1);

/**
* 打印口型，用于调试
* @param mouth TimedMouthShape数组
* @param size TimedMouthShape数组的元素个数
*/
void printMouthShape(const TimedMouthShape* mouth, unsigned size);
#ifdef __cplusplus
}
#endif

#endif //MOUTH_SHAPE_H

