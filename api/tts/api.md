# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [helloworld/v1/error_reason.proto](#helloworld_v1_error_reason-proto)
    - [ErrorReason](#helloworld-v1-ErrorReason)
  
- [helloworld/v1/greeter.proto](#helloworld_v1_greeter-proto)
    - [HelloReply](#helloworld-v1-HelloReply)
    - [HelloRequest](#helloworld-v1-HelloRequest)
  
    - [Greeter](#helloworld-v1-Greeter)
  
- [tts/v2/tts.proto](#tts_v2_tts-proto)
    - [ActionElement](#ttsschema-ActionElement)
    - [BodyMovement](#ttsschema-BodyMovement)
    - [BodyMovementConfig](#ttsschema-BodyMovementConfig)
    - [ConfigAndText](#ttsschema-ConfigAndText)
    - [Coordinate](#ttsschema-Coordinate)
    - [DebugInfo](#ttsschema-DebugInfo)
    - [Expression](#ttsschema-Expression)
    - [FacialExpressionConfig](#ttsschema-FacialExpressionConfig)
    - [MessageDigitalPerson](#ttsschema-MessageDigitalPerson)
    - [MessageEmotion](#ttsschema-MessageEmotion)
    - [MessagePitch](#ttsschema-MessagePitch)
    - [RespGetTtsConfig](#ttsschema-RespGetTtsConfig)
    - [SpeakerList](#ttsschema-SpeakerList)
    - [SpeakerParameter](#ttsschema-SpeakerParameter)
    - [SynthesizedAudio](#ttsschema-SynthesizedAudio)
    - [TimedMouthShape](#ttsschema-TimedMouthShape)
    - [TimedMouthShapes](#ttsschema-TimedMouthShapes)
    - [TtsReq](#ttsschema-TtsReq)
    - [TtsReq.ParameterFlagEntry](#ttsschema-TtsReq-ParameterFlagEntry)
    - [TtsRes](#ttsschema-TtsRes)
    - [VerReq](#ttsschema-VerReq)
    - [VerVersionReq](#ttsschema-VerVersionReq)
    - [VerVersionRsp](#ttsschema-VerVersionRsp)
  
    - [CloudMindsTTS](#ttsschema-CloudMindsTTS)
  
- [tts/v1/schema.proto](#tts_v1_schema-proto)
    - [Expression](#schema-Expression)
    - [MixTtsReq](#schema-MixTtsReq)
    - [SpeakerList](#schema-SpeakerList)
    - [SpeakerParameter](#schema-SpeakerParameter)
    - [TimedMouthShape](#schema-TimedMouthShape)
    - [TtsReq](#schema-TtsReq)
    - [TtsRes](#schema-TtsRes)
    - [VerReq](#schema-VerReq)
    - [VerRsp](#schema-VerRsp)
  
    - [pcmStatus](#schema-pcmStatus)
    - [speakerInfo](#schema-speakerInfo)
    - [ttsErr](#schema-ttsErr)
  
    - [CloudMindsTTS](#schema-CloudMindsTTS)
  
- [Scalar Value Types](#scalar-value-types)



<a name="helloworld_v1_error_reason-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## helloworld/v1/error_reason.proto


 


<a name="helloworld-v1-ErrorReason"></a>

### ErrorReason


| Name | Number | Description |
| ---- | ------ | ----------- |
| GREETER_UNSPECIFIED | 0 |  |
| USER_NOT_FOUND | 1 |  |


 

 

 



<a name="helloworld_v1_greeter-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## helloworld/v1/greeter.proto



<a name="helloworld-v1-HelloReply"></a>

### HelloReply
The response message containing the greetings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |






<a name="helloworld-v1-HelloRequest"></a>

### HelloRequest
The request message containing the user&#39;s name.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |





 

 

 


<a name="helloworld-v1-Greeter"></a>

### Greeter
The greeting service definition.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SayHello | [HelloRequest](#helloworld-v1-HelloRequest) | [HelloReply](#helloworld-v1-HelloReply) | Sends a greeting |

 



<a name="tts_v2_tts-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## tts/v2/tts.proto



<a name="ttsschema-ActionElement"></a>

### ActionElement
Action基元数据


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action_type | [int32](#int32) |  | actionType 基元类型，-100&lt;action_type&lt;100 |
| url | [string](#string) |  | 基元数据的url |
| operation_type | [int32](#int32) |  | action操作类型，-100&lt;operation_type&lt;100 |
| coordinate | [Coordinate](#ttsschema-Coordinate) |  | 该数据的坐标信息 |
| render_duration | [int32](#int32) |  | render_duration 渲染时长（该值不应大于文件时长）单位ms，-1代表持续到指定文件结束，-2代表由coordinate的len_utf8部分决定时长 |






<a name="ttsschema-BodyMovement"></a>

### BodyMovement
动作


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [float](#float) | repeated | 具体动作数据，frame_dim*frame_size个float,frame_dim见BodyMovementConfig |
| frame_size | [int32](#int32) |  | 动作帧数 |
| start_time_ms | [float](#float) |  | 起始时间，单位ms |






<a name="ttsschema-BodyMovementConfig"></a>

### BodyMovementConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| frame_dim | [int32](#int32) |  | 每帧的维度，即一帧由frameDim个float组成 |
| frame_dur_ms | [float](#float) |  | 每帧的持续时长 |






<a name="ttsschema-ConfigAndText"></a>

### ConfigAndText
音频流


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  | 文本信息 |
| facial_expression_config | [FacialExpressionConfig](#ttsschema-FacialExpressionConfig) |  | 表情配置 |
| body_movement_config | [BodyMovementConfig](#ttsschema-BodyMovementConfig) |  | 动作配置 |






<a name="ttsschema-Coordinate"></a>

### Coordinate
坐标信息


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| off | [int32](#int32) |  | 文本起点 |
| len | [int32](#int32) |  | 文本长度 |
| order | [int32](#int32) |  | 当len_utf8 = 0时，order与off_utf8一起描述一个的时间点。具体地，当off_utf8相同时，order相同则认为同时，order不同时，0优先，1次之，依次类推 当len_utf8 &gt; 0时，order无意义 |






<a name="ttsschema-DebugInfo"></a>

### DebugInfo
调试信息


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| debug_type | [string](#string) |  |  |
| info | [string](#string) |  |  |






<a name="ttsschema-Expression"></a>

### Expression
表情


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [float](#float) | repeated | 具体的表情数据, frame_size*frame_dim,frame_dim见FacialExpressionConfig |
| frame_size | [int32](#int32) |  | 表情帧数 |
| start_time_ms | [float](#float) |  | 起始时间，单位ms |






<a name="ttsschema-FacialExpressionConfig"></a>

### FacialExpressionConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| frame_dim | [int32](#int32) |  | 每帧的维度，即一帧由frameDim个float组成 |
| frame_dur_ms | [float](#float) |  | 每帧的持续时长 |






<a name="ttsschema-MessageDigitalPerson"></a>

### MessageDigitalPerson



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| name | [string](#string) |  | 参数 |
| chinese_name | [string](#string) |  | 对应中文 |






<a name="ttsschema-MessageEmotion"></a>

### MessageEmotion



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| name | [string](#string) |  | 参数 |
| chinese_name | [string](#string) |  | 对应中文 |






<a name="ttsschema-MessagePitch"></a>

### MessagePitch



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| name | [string](#string) |  | 参数 |
| chinese_name | [string](#string) |  | 对应中文 |






<a name="ttsschema-RespGetTtsConfig"></a>

### RespGetTtsConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| speaker_list | [SpeakerList](#ttsschema-SpeakerList) |  | 发音人列表 |
| speed_list | [string](#string) | repeated | SupportedSpeed |
| volume_list | [string](#string) | repeated | SupportedVolume |
| pitch_list | [MessagePitch](#ttsschema-MessagePitch) | repeated | SupportedPitch |
| emotion_list | [MessageEmotion](#ttsschema-MessageEmotion) | repeated | SupportedEmotion |
| digital_person_list | [MessageDigitalPerson](#ttsschema-MessageDigitalPerson) | repeated | SupportedDigitalPerson |






<a name="ttsschema-SpeakerList"></a>

### SpeakerList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [SpeakerParameter](#ttsschema-SpeakerParameter) | repeated | 发音人列表 |






<a name="ttsschema-SpeakerParameter"></a>

### SpeakerParameter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| speaker_id | [int32](#int32) |  | 发音人id |
| speaker_name | [string](#string) |  | 发音人名字 |
| parameter_speaker_name | [string](#string) |  | 发音人英文名字 |
| is_support_emotion | [bool](#bool) |  | 是否支持情感 |
| is_support_mixed_voice | [bool](#bool) |  | 是否支持混合发音 |






<a name="ttsschema-SynthesizedAudio"></a>

### SynthesizedAudio
音频流


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pcm | [bytes](#bytes) |  | pcm |
| coordinate | [Coordinate](#ttsschema-Coordinate) |  | 坐标信息 |
| is_punctuation | [int32](#int32) |  | 是否标点1是标点 |






<a name="ttsschema-TimedMouthShape"></a>

### TimedMouthShape



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| duration_us | [uint64](#uint64) |  | 时间 |
| mouth | [int32](#int32) |  | 嘴型 enum MouthShape |






<a name="ttsschema-TimedMouthShapes"></a>

### TimedMouthShapes
口型


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mouths | [TimedMouthShape](#ttsschema-TimedMouthShape) | repeated | 口型数据 |
| start_time_ms | [float](#float) |  | 该段口型的起始时间，单位ms |






<a name="ttsschema-TtsReq"></a>

### TtsReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  | 具体需要合成的文本，支持带json |
| speed | [string](#string) |  | 取值范围请用 GetTtsConfig SpeedList |
| volume | [string](#string) |  | 取值范围请用 GetTtsConfig VolumeList |
| pitch | [string](#string) |  | 取值范围请用 GetTtsConfig PitchList.Name |
| emotions | [string](#string) |  | 如果该发音人支持情感，取值范围请用 GetTtsConfig EmotionList.Name，如果不支持请传&#34;&#34;，否则会报错 |
| parameter_speaker_name | [string](#string) |  | 取值范围请用 GetTtsConfig函数的返回Speakerlist.parameterSpeakerName |
| parameter_digital_person | [string](#string) |  | 数字人形象， |
| parameter_flag | [TtsReq.ParameterFlagEntry](#ttsschema-TtsReq-ParameterFlagEntry) | repeated | 额外信息参数，口型key:mouth,字符串&#34;true&#34;或者&#34;false&#34;、动作key:movement,字符串&#34;true&#34;或者&#34;false&#34;、表情key:expression,字符串&#34;true&#34;或者&#34;false&#34; |
| trace_id | [string](#string) |  |  |
| root_trace_id | [string](#string) |  |  |






<a name="ttsschema-TtsReq-ParameterFlagEntry"></a>

### TtsReq.ParameterFlagEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="ttsschema-TtsRes"></a>

### TtsRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| error_code | [int32](#int32) |  | 错误码，非0为错误，0为正确 |
| status | [int32](#int32) |  | 合成状态，1.开始，2.中间，3，结束 |
| error_msg | [string](#string) |  | 错误信息正确为空，不正确具体字符串 |
| synthesized_audio | [SynthesizedAudio](#ttsschema-SynthesizedAudio) |  | 音频 |
| debug_info | [DebugInfo](#ttsschema-DebugInfo) |  | 调试信息,当有debug_info时需要输出 |
| action_element | [ActionElement](#ttsschema-ActionElement) |  | 基元数据 |
| config_text | [ConfigAndText](#ttsschema-ConfigAndText) |  | 文本配置信息start时会返回 |
| time_mouth_shapes | [TimedMouthShapes](#ttsschema-TimedMouthShapes) |  | 口型数据 |
| expression | [Expression](#ttsschema-Expression) |  | 表情数据 |
| body_movement | [BodyMovement](#ttsschema-BodyMovement) |  | 动作数据 |






<a name="ttsschema-VerReq"></a>

### VerReq







<a name="ttsschema-VerVersionReq"></a>

### VerVersionReq







<a name="ttsschema-VerVersionRsp"></a>

### VerVersionRsp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  |  |





 

 

 


<a name="ttsschema-CloudMindsTTS"></a>

### CloudMindsTTS


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Call | [TtsReq](#ttsschema-TtsReq) | [TtsRes](#ttsschema-TtsRes) stream |  |
| GetVersion | [VerVersionReq](#ttsschema-VerVersionReq) | [VerVersionRsp](#ttsschema-VerVersionRsp) |  |
| GetTtsConfig | [VerReq](#ttsschema-VerReq) | [RespGetTtsConfig](#ttsschema-RespGetTtsConfig) |  |

 



<a name="tts_v1_schema-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## tts/v1/schema.proto



<a name="schema-Expression"></a>

### Expression



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [float](#float) | repeated | 具体的表情数据, frame_size*frame_dim |
| frame_size | [int32](#int32) |  | 表情帧数 |
| frame_dim | [int32](#int32) |  | 一帧的维度，即多少个float |
| frame_time | [float](#float) |  | 一帧的持续时间 |






<a name="schema-MixTtsReq"></a>

### MixTtsReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ttsreq | [TtsReq](#schema-TtsReq) |  |  |
| weight | [float](#float) | repeated | 发音人权重，目前只有27个，超过27个截断，没有补0 |






<a name="schema-SpeakerList"></a>

### SpeakerList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [SpeakerParameter](#schema-SpeakerParameter) | repeated | 发音人列表 |






<a name="schema-SpeakerParameter"></a>

### SpeakerParameter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| speakerId | [int32](#int32) |  | 发音人id |
| speakerName | [string](#string) |  | 发音人名字 |
| parameterSpeakerName | [string](#string) |  | 发音人英文名字 |






<a name="schema-TimedMouthShape"></a>

### TimedMouthShape



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| durationUs | [uint64](#uint64) |  | 时间 |
| mouth | [int32](#int32) |  | 嘴型 enum MouthShape |






<a name="schema-TtsReq"></a>

### TtsReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| location | [string](#string) |  | Cloud 传固定值 |
| type | [string](#string) |  | CloudMinds 传固定值 |
| speaker | [speakerInfo](#schema-speakerInfo) |  | 具体发音人的id 发音人具体参考 |
| speed | [string](#string) |  | 1 2 3 4 5 发音速度 默认3 |
| volume | [string](#string) |  | 1 2 3 4 5 发音音量 默认3 |
| pitch | [string](#string) |  | low medium high 音调 默认 medium |
| streamEnable | [bool](#bool) |  | 是否流式合成 默认 false |
| text | [string](#string) |  | 具体需要合成的文本 |
| textPreHandle | [bool](#bool) |  | 是否需要文本预处理 |
| voiceTuning | [string](#string) |  | voiceTuning开关,传 &#34;on&#34;或者&#34;off&#34; |
| Emotions | [string](#string) |  | 情感发音，取值范围[&#34;Chat&#34;,&#34;Angry&#34;,&#34;Gentle&#34;,&#34;Cheerful&#34;,&#34;Serious&#34;,&#34;General&#34;,&#34;Affectionate&#34;,&#34;Lyrical&#34;,&#34;Newscast&#34;,&#34;CustomerService&#34;] |
| ParameterSpeakerName | [string](#string) |  | parameterSpeakerName 传字符串类型的发音人，兼容speakerInfo，优先选择这个,如果为空，就选择id |
| traceId | [string](#string) |  |  |
| rootTraceId | [string](#string) |  |  |






<a name="schema-TtsRes"></a>

### TtsRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pcm | [bytes](#bytes) |  | 具体的音频流(不会有空帧出现) |
| status | [pcmStatus](#schema-pcmStatus) |  | pcm状态(2中间状态，目前只有2) |
| error | [ttsErr](#schema-ttsErr) |  | 错误码 |
| mouths | [TimedMouthShape](#schema-TimedMouthShape) | repeated | 嘴型，status=2时输出， |
| DebugInfo | [string](#string) |  | 调试信息,当有debufinfo时需要输出 |
| Version | [string](#string) |  | 每次调用都要返回version信息 |
| NormalizedText | [string](#string) |  | 当前正在合成的正则后的文本片段 |
| originalText | [string](#string) |  | 当前正在合成的正则前的文本片段(端侧需要的) |
| expression | [Expression](#schema-Expression) |  | 表情的具体数据 |






<a name="schema-VerReq"></a>

### VerReq







<a name="schema-VerRsp"></a>

### VerRsp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  |  |





 


<a name="schema-pcmStatus"></a>

### pcmStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATUS_DEF | 0 | keep |
| STATUS_START | 1 | 开始 |
| STATUS_MID | 2 | 中间状态 |
| STATUS_END | 3 | 结束 |



<a name="schema-speakerInfo"></a>

### speakerInfo


| Name | Number | Description |
| ---- | ------ | ----------- |
| DaXiaoFang_Chinese_Mandarin_Female_Adult | 0 | 达小芳（温柔女声） |
| DaXiaoQi_Chinese_Mandarin_Female_Adult | 1 | 达小琪（甜美女声） |
| DaDaQiang_Chinese_Mandarin_Male_Adult | 2 | 达大强（磁性男声） |
| DaDaGang_Chinese_Mandarin_Male_Adult | 3 | 达大刚（标准男声） |
| DaTongTong_Chinese_Mandarin_Female_Child | 4 | 达彤彤（儿童女声） |
| DaMingMing_Chinese_Mandarin_Male_Child | 5 | 达明明（儿童男声） |
| DaXiaoChuan_Chinese_Sichuan_Female_Adult | 6 | 达小川（四川女声） |
| DaXiaoYue_Chinese_Yueyu_Female_Adult | 7 | 达小粤（粤语女声） |
| DaXiaoBei_Chinese_Dongbei_Female_Adult | 8 | 达小北（东北女声） |
| DaXiaoTai_Chinese_Taiwan_Female_Adult | 9 | 达小台（台湾女声） |
| DaDaXiang_Chinese_Hunan_Male_Adult | 10 | 达大湘（湖南男声） |
| DaXiaoHu_Chinese_ShangHai_Female_Adult | 11 | 达小沪（上海女声） |
| DaXiaoguan_English_Female_Adult | 12 | 达小罐（英文女声） |
| DaXiaoBo_Emotions_Female_Adult | 13 | 达小波（情感发音人） |
| DaXiaoWu_Chinese_Mandarin_Female_Adult | 14 | 达小舞（电音女声) |



<a name="schema-ttsErr"></a>

### ttsErr


| Name | Number | Description |
| ---- | ------ | ----------- |
| TTS_ERR_INIT | 0 | keep |
| TTS_ERR_OK | 10000 | tts服务正常返回 |
| TTS_ERR_INVALID_RATE | 10001 | rate值错误 |
| TTS_ERR_TEXT_LENGTH_OVERFLOW | 10002 | 请求文本过长 |
| TTS_ERR_INVALID_SPEED | 10003 | 无效的speed[1,5],默认3 |
| TTS_ERR_INVALID_VOLUME | 10004 | 无效的volume[1,5],默认3 |
| TTS_ERR_INVALID_PITCH | 10005 | 无效的pitch[low,medium,high],默认medium |
| TTS_ERR_INVALID_SPEAKER | 10006 | 无效的speed[0,10] |
| TTS_ERR_INVALID_TYPE | 10007 | 无效的type,只支持CloudMinds |
| TTS_ERR_SYN_CANCELLED | 10008 | 音频合成cancelled |
| TTS_ERR_SYN_FAILURE | 10009 | 音频合成异常：cuda oom |
| TTS_ERR_NO_FREE_MOUDLE | 10010 | 没有可用的decoder |
| TTS_ERR_MODULE_INPUT_TEXT | 10011 | 模型输入文本出错 |
| TTS_ERR_TEXT_NORMALIZE | 10012 | 文本序列化出错 |
| TTS_ERR_INVALID_VOICE_TUNING | 10013 | 无效的后置开关值[on,off] |


 

 


<a name="schema-CloudMindsTTS"></a>

### CloudMindsTTS


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Call | [TtsReq](#schema-TtsReq) | [TtsRes](#schema-TtsRes) stream |  |
| GetVersion | [VerReq](#schema-VerReq) | [VerRsp](#schema-VerRsp) |  |
| MixCall | [MixTtsReq](#schema-MixTtsReq) | [TtsRes](#schema-TtsRes) stream |  |
| GetSpeaker | [VerReq](#schema-VerReq) | [SpeakerList](#schema-SpeakerList) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

