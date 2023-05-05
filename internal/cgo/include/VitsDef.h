#ifndef VITS_DEF_H
#define VITS_DEF_H

/**
* vits发音人个数
*/
#define VITS_SPEAKER_NUM  (38)

#ifdef __cplusplus
extern "C"
{
#endif

/**
* vits混合发音人描述符
*/
typedef struct
{
    /**
    * 主发音人ID
    */
    int mainSpeakerId;
    /**
    * 发音人权重
    */
    float speakerWeight[VITS_SPEAKER_NUM];
}VitsBlendSpeakerDescriptor;

#ifdef __cplusplus
}
#endif

#endif //VITS_DEF_H
