#ifndef TTS_TEXT_NORMALIZE_H
#define TTS_TEXT_NORMALIZE_H

#include <list>
#include <vector>
#include <string>
#include <TtsDataType.h>

/**
* TTS Engine类型
*/
enum {

    /**
    * 未知TTS
    */
    TYPE_UNKNOWN = 0,
    /**
    * 自研TTS
    */
    TYPE_CM = 1,
    /**
    * 谷歌TTS
    */
    TYPE_GOOGLE = 2,
};

/**
* TTS标志位
*/
enum {
    /**
    * 离线TTS的标志位
    */
    TYPE_OFFLINE = 1 << 31,
};

/**
* TTS文本规整初始化，设置TTS文本规整的资源路径.
* @param dir res目录的路径
*/
void TtsTextNormalizeInit(const char* dir);

/**
* TTS文本规整.
* 转换TTS文本中的数字、单位、特殊符号，并按需将其切分为不多于64个字符（一个汉字认为是一个字符）的字符串构成的数组。
* @param text 待处理的TTS文本（utf8字符编码）
* @param ttsType TtsEngine的类型（以TYPE_*的组合来表示）
* <p>(TYPE_CM)--在线自研，(TYPE_GOOGLE)--在线谷歌，
* <p>(TYPE_OFFLINE | TYPE_CM)--离线自研，(TYPE_OFFLINE | TYPE_GOOGLE)--离线谷歌，
* @param language 当前TtsEngine所处理的语言
* @param offBase 偏移的基础量
* @param cutType 分句模式，0-普通模式（句号类分句+逗号类分句+顿号双空格），1-VITS模式（句号类分句+逗号类双空格+顿号双空格），2-停顿分句模式（句号类分句+逗号类分句+双空格分句）
* @param cutLenCh 中文句子分句长度阈值
* @param cutLenEn 英文句子分句长度阈值
* @param ttsTextInfo TTS文本规整信息
* @return int 0-正常
*/
int TtsTextNormalize(const char* text, unsigned int ttsType, const char* language, int offBase, int cutType,
        int cutLenCh, int cutLenEn, TtsTextInfo& ttsTextInfo);

/**
 * 不需要TTS文本规整的文本，填充相应的TTS文本规整信息
 *
 * @param text 待处理的TTS文本（utf8字符编码）
 * @param ttsTextInfo TTS文本规整信息
 * @return int 0-正常
 */
int dummyTextNormalize(const char* text, TtsTextInfo& ttsTextInfo);

void printTtsTextNormalizeOutput(TtsTextInfo& ttsTextInfo);


void getShortSentenceBoundaries(const char* text, vector<SentenceBoundary> & boundaries);


#endif //TTS_TEXT_NORMALIZE_H

