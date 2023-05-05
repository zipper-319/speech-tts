#ifndef TTS_SETTING_H
#define TTS_SETTING_H

#ifdef __cplusplus
extern "C"
{
#endif

/**
* 发音人特性集
*/
enum {
    /**
    * 支持混合发音
    */
    SUPPORT_MIXED_VOICE = 1<<0,
    /**
    * 支持情绪
    */
    SUPPORT_EMOTION = 1<<1,
    /**
    * 顿号用于分句
    */
    CUT_SENTENCES_USE_CHINESE_SERIAL_COMMA = 1<<2,
};


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
    */
    unsigned int flags;

    /**
    * 发音人的扩展特性
    */
    const char* extended_features;

    /**
    * 发音人中文最大长度
    */
    const int CHINESE_WORDS_LIMIT;

    /**
    * 发音人英文最大长度
    */
    const int ENGLISH_WORDS_LIMIT;

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
    * 输出模式
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
    * 数字人形象
    */
    PARAMETER_DIGITAL_PERSON = 6,
//以下仅供内部debug
    /**
    * tuning
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
    * 当肢体动作使能时，指定肢体动作对应的数字人
    */
    const char* digitalPerson;
}TtsSetting;


/**
* 获取CmTts所支持的发音人数组
* @param supportedSpeakers 用于输出发音人（名字）数组
* @param element_num 用于输出发音人的个数（即发音人数组的大小）
* @return 0--成功,&lt;0--失败
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

const SpeakerDescriptor* GetSpeakerDescriptor(const char* speakerName);

#ifdef __cplusplus
}
#endif

#endif //TTS_SETTING_H
