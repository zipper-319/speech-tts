#ifndef TTS_SETTING_H
#define TTS_SETTING_H

extern const char* DefaultUserSpace;//value: "cloudminds"

#ifdef __cplusplus
extern "C"
{
#endif

/**
* 正则断句模式
*/
enum {
    /**
    * 默认切分（句号类分句+逗号类分句+顿号双空格）
    */
    Normalize_Cut_Default = 0,
    /**
    * 长切分（句号类双空格+逗号类双空格+顿号双空格）
    */
    Normalize_Cut_Long = 1,
    /**
    * 短切分（句号类分句+逗号类分句+双空格分句）
    */
    Normalize_Cut_Short = 2,
    /**
    * 长切分（句号类分句+逗号类双空格+顿号双空格）
    */
    Normalize_Cut_Long_Period = 3,
    /**
    * 长切分（首句标点切分，非首句等同于Normalize_Cut_Long）
    */
    Normalize_Cut_Long_FirstSentence = 4,
};

/**
* feature标签
*/
enum {
    /**
    * 不支持
    */
    NOT_SUPPORT = 0,
    /**
    * 支持
    */
    SUPPORT = 1,
};

/**
* feature shift
*/
enum {
    LEFT_SHIFT_support_mixed_voice = 4,
    LEFT_SHIFT_support_emotion = 5,
};

/**
* feature helper
*/
enum {
    SUPPORT_MIXED_VOICE = SUPPORT<<LEFT_SHIFT_support_mixed_voice,
    SUPPORT_EMOTION = SUPPORT<<LEFT_SHIFT_support_emotion,
    NORMALIZE_CUT_MASK = 0x0000000F,
};

/*
typedef struct {
    unsigned placeholder: 26;
    unsigned support_emotion: 1;
    unsigned support_mixed_voice: 1;
    unsigned normalize_cut_type: 4;
}FeatureFlag;
*/

/**
* 发音人描述符
*/
typedef struct{
    /**
    * 发音人描述符的id,具有唯一性
    */
    unsigned int id;
    /**
    * 发音人姓名
    */
    const char* name;
    /**
    * 发音人所使用的语言
    */
    const char* language;
    /**
    * 发音人性别
    */
    const char* gender;
    /**
    * 发音人年龄分类
    */
    const char* age_category;

    /**
    * 发音人特性
    * 从高位到低位（31~0）|26bit(undefined)|1bit(support_emotion)|1bit(support_mixed_voice)|4bit(normalize_cut_type)|
    */
    unsigned int flags;

    /**
    * 发音人的扩展特性（未使用）
    */
    const char* extended_features;

    /**
    * 发音人中文最大长度
    */
    int CHINESE_WORDS_LIMIT;

    /**
    * 发音人英文最大长度
    */
    int ENGLISH_WORDS_LIMIT;

}SpeakerDescriptor;

/**
* 合成参数类型
*/
enum {
    /**
    * 语速
    */
    PARAMETER_SPEED = 0,
    /**
    * 音量
    */
    PARAMETER_VOLUME = 1,
    /**
    * 音调
    */
    PARAMETER_PITCH = 2,
    /**
    * @deprecated 输出模式
    */
    PARAMETER_OUTPUT_MODE = 3,
    /**
    * 合成模式
    */
    PARAMETER_SYNTHESIS_MODE = 4,
    /**
    * 说话风格
    */
    PARAMETER_SPEAKING_STYLE = 5,
    /**
    * 动作描述
    */
    PARAMETER_MOVEMENT_DESCRIPTOR = 6,
    /**
    * 表情描述
    */
    PARAMETER_EXPRESSION_DESCRIPTOR = 7,
    /**
    * 音频编码格式
    */
    PARAMETER_AUDIO_ENCODE = 8,
//以下仅供内部debug
    /**
    * @deprecated tuning
    */
    PARAMETER_TUNING = 100,
};

/**
* 功能集
*/
enum {
    /**
    * 使能表情合成
    */
    ENABLE_EXPRESSION = 1<<0,
    /**
    * 使能口型
    */
    ENABLE_MOUTHSHAPE = 1<<1,
    /**
    * 使能肢体动作合成
    */
    ENABLE_MOVEMENT = 1<<2,
};

/**
* 音频编码格式
*/
enum AudioEncoding{
    /**
    * Uncompressed 16-bit signed little-endian samples (Linear PCM)
    */
    LINEAR16 = 0,
    /**
    * opus
    */
    OPUS = 1,
};


/**
* TTS配置
*/
typedef struct{
    /**
    * 发音人
    */
    const char* speaker;
    /**
    * 语速
    */
    const char* speed;
    /**
    * 音量
    */
    const char* volume;
    /**
    * 音调
    */
    const char* pitch;
    /**
    * 说话风格
    */
    const char* speakingStyle;
    /**
    * 功能集设置
    */
    unsigned int featureSet;
    /**
    * 对于v1接口，置为NULL; 对于v2接口, 且表情使能时，指定表情描述符
    */
    const char* expressionDescriptor;
    /**
    * 对于v1接口，置为NULL; 对于v2接口, 且动作使能时，指定动作描述符
    */
    const char* movementDescriptor;
    /**
    * 语种提示，提示合成请求的文本的语种（https://zh.wikipedia.org/zh-hans/ISO_639-1 http://www.lingoes.net/zh/translator/langcode.htm ，当前只支持中文——zh(含变种)，英文——en）
    */
    const char* languageTip;
    /**
    * 指定speaker所属的用户空间，即要访问的该用户空间下的speaker，为空时默认为cloudminds
    */
    const char* userSpace;
    /**
    * 指定音频编码格式，参考enum AudioEncoding，默认为LINEAR16
    */
    int audioEncoding;
}TtsSetting;


/**
* 获取CmTts所支持的发音人数组
* @param supportedSpeakers 用于输出发音人（名字）数组
* @param element_num 用于输出发音人的个数（即发音人数组的大小）
* @return 0--成功,&lt;0--失败参考enum AudioEncoding，默认为LINEAR16
*/
int GetSupportedSpeaker(const char*** supportedSpeakers, unsigned int *element_num);
/**
* 获取指定参数所支持的取值数组
* @param type 在 PARAMETER_* 中取值，用于指定参数的类型，即表明要获取的是哪个参数所支持的取值
* @param supportedValues 用于输出某参数的取值数组，取值用字符串表示
* @param element_num 用于输出取值数组的大小
* @return 0--成功,&lt;0--失败
*/
int GetSupportedParameter(unsigned int type, const char*** supportedValues, unsigned int *element_num);

/**
* 获取特定用户空间下发音人的发音人描述
* @param userSpace 用户空间
* @param speakerName 发音人
* @return SpeakerDescriptor指针
*/
const SpeakerDescriptor* GetSpeakerDescriptor(const char* userSpace, const char* speakerName);


/**
* 获取特定用户空间下的发音人数组
* @param userSpace 用户空间
* @param supportedSpeakers 用于输出发音人（名字）数组
* @param element_num 用于输出发音人的个数（即发音人数组的大小）
* @return 0--成功,&lt;0--失败
*/
int GetUserSpaceSupportedSpeaker(const char* userSpace, const char*** supportedSpeakers, unsigned int *element_num);

#ifdef __cplusplus
}
#endif

#endif //TTS_SETTING_H
