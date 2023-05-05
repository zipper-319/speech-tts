#ifndef ACTION_DEF_H
#define ACTION_DEF_H

#ifdef __cplusplus
extern "C"
{
#endif

/**
* action基元类型，-100&lt;value&lt;100
*/
enum ActionElement {
    //ActionElement_Nothing = -1,
    ActionElement_TtsText = 0,
    ActionElement_AnnotationText = 1,
    ActionElement_MuteAudio = 2,
    ActionElement_GivenAudio = 3,
    ActionElement_Image = 4,
    ActionElement_Video = 5,
    ActionElement_BodyMovement = 6,
    ActionElement_EmbeddedWeb = 7,
    ActionElement_AudioWithMovement = 8,
    ActionElement_Idle = 9,
};
/**
* action操作类型，-100&lt;value&lt;100
*/
enum ActionOperation {
    ActionOperation_Nothing = 0,
    ActionOperation_TtsSetting = 1,
    ActionOperation_Insert = 2,
    ActionOperation_MixAudio = 3,
    ActionOperation_ReplaceAudio = 4,
    ActionOperation_MuteTts = 5,
    ActionOperation_IndependentRender = 6,
    ActionOperation_BackgroundAudio = 7,
};

#ifdef __cplusplus
}
#endif

#endif //ACTION_DEF_H
