// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/texttospeech/v1beta1/cloud_tts.proto

package texttospeech

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Gender of the voice as described in
// [SSML voice element](https://www.w3.org/TR/speech-synthesis11/#edef_voice).
type SsmlVoiceGender int32

const (
	// An unspecified gender.
	// In VoiceSelectionParams, this means that the client doesn't care which
	// gender the selected voice will have. In the Voice field of
	// ListVoicesResponse, this may mean that the voice doesn't fit any of the
	// other categories in this enum, or that the gender of the voice isn't known.
	SsmlVoiceGender_SSML_VOICE_GENDER_UNSPECIFIED SsmlVoiceGender = 0
	// A male voice.
	SsmlVoiceGender_MALE SsmlVoiceGender = 1
	// A female voice.
	SsmlVoiceGender_FEMALE SsmlVoiceGender = 2
	// A gender-neutral voice.
	SsmlVoiceGender_NEUTRAL SsmlVoiceGender = 3
)

var SsmlVoiceGender_name = map[int32]string{
	0: "SSML_VOICE_GENDER_UNSPECIFIED",
	1: "MALE",
	2: "FEMALE",
	3: "NEUTRAL",
}

var SsmlVoiceGender_value = map[string]int32{
	"SSML_VOICE_GENDER_UNSPECIFIED": 0,
	"MALE":                          1,
	"FEMALE":                        2,
	"NEUTRAL":                       3,
}

func (x SsmlVoiceGender) String() string {
	return proto.EnumName(SsmlVoiceGender_name, int32(x))
}

func (SsmlVoiceGender) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{0}
}

// Configuration to set up audio encoder. The encoding determines the output
// audio format that we'd like.
type AudioEncoding int32

const (
	// Not specified. Will return result [google.rpc.Code.INVALID_ARGUMENT][].
	AudioEncoding_AUDIO_ENCODING_UNSPECIFIED AudioEncoding = 0
	// Uncompressed 16-bit signed little-endian samples (Linear PCM).
	// Audio content returned as LINEAR16 also contains a WAV header.
	AudioEncoding_LINEAR16 AudioEncoding = 1
	// MP3 audio.
	AudioEncoding_MP3 AudioEncoding = 2
	// Opus encoded audio wrapped in an ogg container. The result will be a
	// file which can be played natively on Android, and in browsers (at least
	// Chrome and Firefox). The quality of the encoding is considerably higher
	// than MP3 while using approximately the same bitrate.
	AudioEncoding_OGG_OPUS AudioEncoding = 3
)

var AudioEncoding_name = map[int32]string{
	0: "AUDIO_ENCODING_UNSPECIFIED",
	1: "LINEAR16",
	2: "MP3",
	3: "OGG_OPUS",
}

var AudioEncoding_value = map[string]int32{
	"AUDIO_ENCODING_UNSPECIFIED": 0,
	"LINEAR16":                   1,
	"MP3":                        2,
	"OGG_OPUS":                   3,
}

func (x AudioEncoding) String() string {
	return proto.EnumName(AudioEncoding_name, int32(x))
}

func (AudioEncoding) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{1}
}

// The top-level message sent by the client for the `ListVoices` method.
type ListVoicesRequest struct {
	// Optional (but recommended)
	// [BCP-47](https://www.rfc-editor.org/rfc/bcp/bcp47.txt) language tag. If
	// specified, the ListVoices call will only return voices that can be used to
	// synthesize this language_code. E.g. when specifying "en-NZ", you will get
	// supported "en-*" voices; when specifying "no", you will get supported
	// "no-*" (Norwegian) and "nb-*" (Norwegian Bokmal) voices; specifying "zh"
	// will also get supported "cmn-*" voices; specifying "zh-hk" will also get
	// supported "yue-*" voices.
	LanguageCode         string   `protobuf:"bytes,1,opt,name=language_code,json=languageCode,proto3" json:"language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListVoicesRequest) Reset()         { *m = ListVoicesRequest{} }
func (m *ListVoicesRequest) String() string { return proto.CompactTextString(m) }
func (*ListVoicesRequest) ProtoMessage()    {}
func (*ListVoicesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{0}
}
func (m *ListVoicesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListVoicesRequest.Unmarshal(m, b)
}
func (m *ListVoicesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListVoicesRequest.Marshal(b, m, deterministic)
}
func (m *ListVoicesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListVoicesRequest.Merge(m, src)
}
func (m *ListVoicesRequest) XXX_Size() int {
	return xxx_messageInfo_ListVoicesRequest.Size(m)
}
func (m *ListVoicesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListVoicesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListVoicesRequest proto.InternalMessageInfo

func (m *ListVoicesRequest) GetLanguageCode() string {
	if m != nil {
		return m.LanguageCode
	}
	return ""
}

// The message returned to the client by the `ListVoices` method.
type ListVoicesResponse struct {
	// The list of voices.
	Voices               []*Voice `protobuf:"bytes,1,rep,name=voices,proto3" json:"voices,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListVoicesResponse) Reset()         { *m = ListVoicesResponse{} }
func (m *ListVoicesResponse) String() string { return proto.CompactTextString(m) }
func (*ListVoicesResponse) ProtoMessage()    {}
func (*ListVoicesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{1}
}
func (m *ListVoicesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListVoicesResponse.Unmarshal(m, b)
}
func (m *ListVoicesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListVoicesResponse.Marshal(b, m, deterministic)
}
func (m *ListVoicesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListVoicesResponse.Merge(m, src)
}
func (m *ListVoicesResponse) XXX_Size() int {
	return xxx_messageInfo_ListVoicesResponse.Size(m)
}
func (m *ListVoicesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListVoicesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListVoicesResponse proto.InternalMessageInfo

func (m *ListVoicesResponse) GetVoices() []*Voice {
	if m != nil {
		return m.Voices
	}
	return nil
}

// Description of a voice supported by the TTS service.
type Voice struct {
	// The languages that this voice supports, expressed as
	// [BCP-47](https://www.rfc-editor.org/rfc/bcp/bcp47.txt) language tags (e.g.
	// "en-US", "es-419", "cmn-tw").
	LanguageCodes []string `protobuf:"bytes,1,rep,name=language_codes,json=languageCodes,proto3" json:"language_codes,omitempty"`
	// The name of this voice.  Each distinct voice has a unique name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The gender of this voice.
	SsmlGender SsmlVoiceGender `protobuf:"varint,3,opt,name=ssml_gender,json=ssmlGender,proto3,enum=google.cloud.texttospeech.v1beta1.SsmlVoiceGender" json:"ssml_gender,omitempty"`
	// The natural sample rate (in hertz) for this voice.
	NaturalSampleRateHertz int32    `protobuf:"varint,4,opt,name=natural_sample_rate_hertz,json=naturalSampleRateHertz,proto3" json:"natural_sample_rate_hertz,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *Voice) Reset()         { *m = Voice{} }
func (m *Voice) String() string { return proto.CompactTextString(m) }
func (*Voice) ProtoMessage()    {}
func (*Voice) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{2}
}
func (m *Voice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Voice.Unmarshal(m, b)
}
func (m *Voice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Voice.Marshal(b, m, deterministic)
}
func (m *Voice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Voice.Merge(m, src)
}
func (m *Voice) XXX_Size() int {
	return xxx_messageInfo_Voice.Size(m)
}
func (m *Voice) XXX_DiscardUnknown() {
	xxx_messageInfo_Voice.DiscardUnknown(m)
}

var xxx_messageInfo_Voice proto.InternalMessageInfo

func (m *Voice) GetLanguageCodes() []string {
	if m != nil {
		return m.LanguageCodes
	}
	return nil
}

func (m *Voice) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Voice) GetSsmlGender() SsmlVoiceGender {
	if m != nil {
		return m.SsmlGender
	}
	return SsmlVoiceGender_SSML_VOICE_GENDER_UNSPECIFIED
}

func (m *Voice) GetNaturalSampleRateHertz() int32 {
	if m != nil {
		return m.NaturalSampleRateHertz
	}
	return 0
}

// The top-level message sent by the client for the `SynthesizeSpeech` method.
type SynthesizeSpeechRequest struct {
	// Required. The Synthesizer requires either plain text or SSML as input.
	Input *SynthesisInput `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	// Required. The desired voice of the synthesized audio.
	Voice *VoiceSelectionParams `protobuf:"bytes,2,opt,name=voice,proto3" json:"voice,omitempty"`
	// Required. The configuration of the synthesized audio.
	AudioConfig          *AudioConfig `protobuf:"bytes,3,opt,name=audio_config,json=audioConfig,proto3" json:"audio_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SynthesizeSpeechRequest) Reset()         { *m = SynthesizeSpeechRequest{} }
func (m *SynthesizeSpeechRequest) String() string { return proto.CompactTextString(m) }
func (*SynthesizeSpeechRequest) ProtoMessage()    {}
func (*SynthesizeSpeechRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{3}
}
func (m *SynthesizeSpeechRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynthesizeSpeechRequest.Unmarshal(m, b)
}
func (m *SynthesizeSpeechRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynthesizeSpeechRequest.Marshal(b, m, deterministic)
}
func (m *SynthesizeSpeechRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynthesizeSpeechRequest.Merge(m, src)
}
func (m *SynthesizeSpeechRequest) XXX_Size() int {
	return xxx_messageInfo_SynthesizeSpeechRequest.Size(m)
}
func (m *SynthesizeSpeechRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SynthesizeSpeechRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SynthesizeSpeechRequest proto.InternalMessageInfo

func (m *SynthesizeSpeechRequest) GetInput() *SynthesisInput {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *SynthesizeSpeechRequest) GetVoice() *VoiceSelectionParams {
	if m != nil {
		return m.Voice
	}
	return nil
}

func (m *SynthesizeSpeechRequest) GetAudioConfig() *AudioConfig {
	if m != nil {
		return m.AudioConfig
	}
	return nil
}

// Contains text input to be synthesized. Either `text` or `ssml` must be
// supplied. Supplying both or neither returns
// [google.rpc.Code.INVALID_ARGUMENT][]. The input size is limited to 5000
// characters.
type SynthesisInput struct {
	// The input source, which is either plain text or SSML.
	//
	// Types that are valid to be assigned to InputSource:
	//	*SynthesisInput_Text
	//	*SynthesisInput_Ssml
	InputSource          isSynthesisInput_InputSource `protobuf_oneof:"input_source"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SynthesisInput) Reset()         { *m = SynthesisInput{} }
func (m *SynthesisInput) String() string { return proto.CompactTextString(m) }
func (*SynthesisInput) ProtoMessage()    {}
func (*SynthesisInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{4}
}
func (m *SynthesisInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynthesisInput.Unmarshal(m, b)
}
func (m *SynthesisInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynthesisInput.Marshal(b, m, deterministic)
}
func (m *SynthesisInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynthesisInput.Merge(m, src)
}
func (m *SynthesisInput) XXX_Size() int {
	return xxx_messageInfo_SynthesisInput.Size(m)
}
func (m *SynthesisInput) XXX_DiscardUnknown() {
	xxx_messageInfo_SynthesisInput.DiscardUnknown(m)
}

var xxx_messageInfo_SynthesisInput proto.InternalMessageInfo

type isSynthesisInput_InputSource interface {
	isSynthesisInput_InputSource()
}

type SynthesisInput_Text struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3,oneof"`
}

type SynthesisInput_Ssml struct {
	Ssml string `protobuf:"bytes,2,opt,name=ssml,proto3,oneof"`
}

func (*SynthesisInput_Text) isSynthesisInput_InputSource() {}

func (*SynthesisInput_Ssml) isSynthesisInput_InputSource() {}

func (m *SynthesisInput) GetInputSource() isSynthesisInput_InputSource {
	if m != nil {
		return m.InputSource
	}
	return nil
}

func (m *SynthesisInput) GetText() string {
	if x, ok := m.GetInputSource().(*SynthesisInput_Text); ok {
		return x.Text
	}
	return ""
}

func (m *SynthesisInput) GetSsml() string {
	if x, ok := m.GetInputSource().(*SynthesisInput_Ssml); ok {
		return x.Ssml
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SynthesisInput) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SynthesisInput_OneofMarshaler, _SynthesisInput_OneofUnmarshaler, _SynthesisInput_OneofSizer, []interface{}{
		(*SynthesisInput_Text)(nil),
		(*SynthesisInput_Ssml)(nil),
	}
}

func _SynthesisInput_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SynthesisInput)
	// input_source
	switch x := m.InputSource.(type) {
	case *SynthesisInput_Text:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Text)
	case *SynthesisInput_Ssml:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Ssml)
	case nil:
	default:
		return fmt.Errorf("SynthesisInput.InputSource has unexpected type %T", x)
	}
	return nil
}

func _SynthesisInput_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SynthesisInput)
	switch tag {
	case 1: // input_source.text
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.InputSource = &SynthesisInput_Text{x}
		return true, err
	case 2: // input_source.ssml
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.InputSource = &SynthesisInput_Ssml{x}
		return true, err
	default:
		return false, nil
	}
}

func _SynthesisInput_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SynthesisInput)
	// input_source
	switch x := m.InputSource.(type) {
	case *SynthesisInput_Text:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Text)))
		n += len(x.Text)
	case *SynthesisInput_Ssml:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Ssml)))
		n += len(x.Ssml)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Description of which voice to use for a synthesis request.
type VoiceSelectionParams struct {
	// The language (and optionally also the region) of the voice expressed as a
	// [BCP-47](https://www.rfc-editor.org/rfc/bcp/bcp47.txt) language tag, e.g.
	// "en-US". Required. This should not include a script tag (e.g. use
	// "cmn-cn" rather than "cmn-Hant-cn"), because the script will be inferred
	// from the input provided in the SynthesisInput.  The TTS service
	// will use this parameter to help choose an appropriate voice.  Note that
	// the TTS service may choose a voice with a slightly different language code
	// than the one selected; it may substitute a different region
	// (e.g. using en-US rather than en-CA if there isn't a Canadian voice
	// available), or even a different language, e.g. using "nb" (Norwegian
	// Bokmal) instead of "no" (Norwegian)".
	LanguageCode string `protobuf:"bytes,1,opt,name=language_code,json=languageCode,proto3" json:"language_code,omitempty"`
	// The name of the voice. Optional; if not set, the service will choose a
	// voice based on the other parameters such as language_code and gender.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The preferred gender of the voice. Optional; if not set, the service will
	// choose a voice based on the other parameters such as language_code and
	// name. Note that this is only a preference, not requirement; if a
	// voice of the appropriate gender is not available, the synthesizer should
	// substitute a voice with a different gender rather than failing the request.
	SsmlGender           SsmlVoiceGender `protobuf:"varint,3,opt,name=ssml_gender,json=ssmlGender,proto3,enum=google.cloud.texttospeech.v1beta1.SsmlVoiceGender" json:"ssml_gender,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *VoiceSelectionParams) Reset()         { *m = VoiceSelectionParams{} }
func (m *VoiceSelectionParams) String() string { return proto.CompactTextString(m) }
func (*VoiceSelectionParams) ProtoMessage()    {}
func (*VoiceSelectionParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{5}
}
func (m *VoiceSelectionParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoiceSelectionParams.Unmarshal(m, b)
}
func (m *VoiceSelectionParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoiceSelectionParams.Marshal(b, m, deterministic)
}
func (m *VoiceSelectionParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoiceSelectionParams.Merge(m, src)
}
func (m *VoiceSelectionParams) XXX_Size() int {
	return xxx_messageInfo_VoiceSelectionParams.Size(m)
}
func (m *VoiceSelectionParams) XXX_DiscardUnknown() {
	xxx_messageInfo_VoiceSelectionParams.DiscardUnknown(m)
}

var xxx_messageInfo_VoiceSelectionParams proto.InternalMessageInfo

func (m *VoiceSelectionParams) GetLanguageCode() string {
	if m != nil {
		return m.LanguageCode
	}
	return ""
}

func (m *VoiceSelectionParams) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *VoiceSelectionParams) GetSsmlGender() SsmlVoiceGender {
	if m != nil {
		return m.SsmlGender
	}
	return SsmlVoiceGender_SSML_VOICE_GENDER_UNSPECIFIED
}

// Description of audio data to be synthesized.
type AudioConfig struct {
	// Required. The format of the requested audio byte stream.
	AudioEncoding AudioEncoding `protobuf:"varint,1,opt,name=audio_encoding,json=audioEncoding,proto3,enum=google.cloud.texttospeech.v1beta1.AudioEncoding" json:"audio_encoding,omitempty"`
	// Optional speaking rate/speed, in the range [0.25, 4.0]. 1.0 is the normal
	// native speed supported by the specific voice. 2.0 is twice as fast, and
	// 0.5 is half as fast. If unset(0.0), defaults to the native 1.0 speed. Any
	// other values < 0.25 or > 4.0 will return an error.
	SpeakingRate float64 `protobuf:"fixed64,2,opt,name=speaking_rate,json=speakingRate,proto3" json:"speaking_rate,omitempty"`
	// Optional speaking pitch, in the range [-20.0, 20.0]. 20 means increase 20
	// semitones from the original pitch. -20 means decrease 20 semitones from the
	// original pitch.
	Pitch float64 `protobuf:"fixed64,3,opt,name=pitch,proto3" json:"pitch,omitempty"`
	// Optional volume gain (in dB) of the normal native volume supported by the
	// specific voice, in the range [-96.0, 16.0]. If unset, or set to a value of
	// 0.0 (dB), will play at normal native signal amplitude. A value of -6.0 (dB)
	// will play at approximately half the amplitude of the normal native signal
	// amplitude. A value of +6.0 (dB) will play at approximately twice the
	// amplitude of the normal native signal amplitude. Strongly recommend not to
	// exceed +10 (dB) as there's usually no effective increase in loudness for
	// any value greater than that.
	VolumeGainDb float64 `protobuf:"fixed64,4,opt,name=volume_gain_db,json=volumeGainDb,proto3" json:"volume_gain_db,omitempty"`
	// The synthesis sample rate (in hertz) for this audio. Optional.  If this is
	// different from the voice's natural sample rate, then the synthesizer will
	// honor this request by converting to the desired sample rate (which might
	// result in worse audio quality), unless the specified sample rate is not
	// supported for the encoding chosen, in which case it will fail the request
	// and return [google.rpc.Code.INVALID_ARGUMENT][].
	SampleRateHertz int32 `protobuf:"varint,5,opt,name=sample_rate_hertz,json=sampleRateHertz,proto3" json:"sample_rate_hertz,omitempty"`
	// An identifier which selects 'audio effects' profiles that are applied on
	// (post synthesized) text to speech.
	// Effects are applied on top of each other in the order they are given.
	EffectsProfileId     []string `protobuf:"bytes,6,rep,name=effects_profile_id,json=effectsProfileId,proto3" json:"effects_profile_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AudioConfig) Reset()         { *m = AudioConfig{} }
func (m *AudioConfig) String() string { return proto.CompactTextString(m) }
func (*AudioConfig) ProtoMessage()    {}
func (*AudioConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{6}
}
func (m *AudioConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AudioConfig.Unmarshal(m, b)
}
func (m *AudioConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AudioConfig.Marshal(b, m, deterministic)
}
func (m *AudioConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AudioConfig.Merge(m, src)
}
func (m *AudioConfig) XXX_Size() int {
	return xxx_messageInfo_AudioConfig.Size(m)
}
func (m *AudioConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AudioConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AudioConfig proto.InternalMessageInfo

func (m *AudioConfig) GetAudioEncoding() AudioEncoding {
	if m != nil {
		return m.AudioEncoding
	}
	return AudioEncoding_AUDIO_ENCODING_UNSPECIFIED
}

func (m *AudioConfig) GetSpeakingRate() float64 {
	if m != nil {
		return m.SpeakingRate
	}
	return 0
}

func (m *AudioConfig) GetPitch() float64 {
	if m != nil {
		return m.Pitch
	}
	return 0
}

func (m *AudioConfig) GetVolumeGainDb() float64 {
	if m != nil {
		return m.VolumeGainDb
	}
	return 0
}

func (m *AudioConfig) GetSampleRateHertz() int32 {
	if m != nil {
		return m.SampleRateHertz
	}
	return 0
}

func (m *AudioConfig) GetEffectsProfileId() []string {
	if m != nil {
		return m.EffectsProfileId
	}
	return nil
}

// The message returned to the client by the `SynthesizeSpeech` method.
type SynthesizeSpeechResponse struct {
	// The audio data bytes encoded as specified in the request, including the
	// header (For LINEAR16 audio, we include the WAV header). Note: as
	// with all bytes fields, protobuffers use a pure binary representation,
	// whereas JSON representations use base64.
	AudioContent         []byte   `protobuf:"bytes,1,opt,name=audio_content,json=audioContent,proto3" json:"audio_content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SynthesizeSpeechResponse) Reset()         { *m = SynthesizeSpeechResponse{} }
func (m *SynthesizeSpeechResponse) String() string { return proto.CompactTextString(m) }
func (*SynthesizeSpeechResponse) ProtoMessage()    {}
func (*SynthesizeSpeechResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dfc2f1f47f5721c8, []int{7}
}
func (m *SynthesizeSpeechResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SynthesizeSpeechResponse.Unmarshal(m, b)
}
func (m *SynthesizeSpeechResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SynthesizeSpeechResponse.Marshal(b, m, deterministic)
}
func (m *SynthesizeSpeechResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SynthesizeSpeechResponse.Merge(m, src)
}
func (m *SynthesizeSpeechResponse) XXX_Size() int {
	return xxx_messageInfo_SynthesizeSpeechResponse.Size(m)
}
func (m *SynthesizeSpeechResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SynthesizeSpeechResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SynthesizeSpeechResponse proto.InternalMessageInfo

func (m *SynthesizeSpeechResponse) GetAudioContent() []byte {
	if m != nil {
		return m.AudioContent
	}
	return nil
}

func init() {
	proto.RegisterType((*ListVoicesRequest)(nil), "google.cloud.texttospeech.v1beta1.ListVoicesRequest")
	proto.RegisterType((*ListVoicesResponse)(nil), "google.cloud.texttospeech.v1beta1.ListVoicesResponse")
	proto.RegisterType((*Voice)(nil), "google.cloud.texttospeech.v1beta1.Voice")
	proto.RegisterType((*SynthesizeSpeechRequest)(nil), "google.cloud.texttospeech.v1beta1.SynthesizeSpeechRequest")
	proto.RegisterType((*SynthesisInput)(nil), "google.cloud.texttospeech.v1beta1.SynthesisInput")
	proto.RegisterType((*VoiceSelectionParams)(nil), "google.cloud.texttospeech.v1beta1.VoiceSelectionParams")
	proto.RegisterType((*AudioConfig)(nil), "google.cloud.texttospeech.v1beta1.AudioConfig")
	proto.RegisterType((*SynthesizeSpeechResponse)(nil), "google.cloud.texttospeech.v1beta1.SynthesizeSpeechResponse")
	proto.RegisterEnum("google.cloud.texttospeech.v1beta1.SsmlVoiceGender", SsmlVoiceGender_name, SsmlVoiceGender_value)
	proto.RegisterEnum("google.cloud.texttospeech.v1beta1.AudioEncoding", AudioEncoding_name, AudioEncoding_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TextToSpeechClient is the client API for TextToSpeech service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TextToSpeechClient interface {
	// Returns a list of [Voice][google.cloud.texttospeech.v1beta1.Voice]
	// supported for synthesis.
	ListVoices(ctx context.Context, in *ListVoicesRequest, opts ...grpc.CallOption) (*ListVoicesResponse, error)
	// Synthesizes speech synchronously: receive results after all text input
	// has been processed.
	SynthesizeSpeech(ctx context.Context, in *SynthesizeSpeechRequest, opts ...grpc.CallOption) (*SynthesizeSpeechResponse, error)
}

type textToSpeechClient struct {
	cc *grpc.ClientConn
}

func NewTextToSpeechClient(cc *grpc.ClientConn) TextToSpeechClient {
	return &textToSpeechClient{cc}
}

func (c *textToSpeechClient) ListVoices(ctx context.Context, in *ListVoicesRequest, opts ...grpc.CallOption) (*ListVoicesResponse, error) {
	out := new(ListVoicesResponse)
	err := c.cc.Invoke(ctx, "/google.cloud.texttospeech.v1beta1.TextToSpeech/ListVoices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *textToSpeechClient) SynthesizeSpeech(ctx context.Context, in *SynthesizeSpeechRequest, opts ...grpc.CallOption) (*SynthesizeSpeechResponse, error) {
	out := new(SynthesizeSpeechResponse)
	err := c.cc.Invoke(ctx, "/google.cloud.texttospeech.v1beta1.TextToSpeech/SynthesizeSpeech", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TextToSpeechServer is the server API for TextToSpeech service.
type TextToSpeechServer interface {
	// Returns a list of [Voice][google.cloud.texttospeech.v1beta1.Voice]
	// supported for synthesis.
	ListVoices(context.Context, *ListVoicesRequest) (*ListVoicesResponse, error)
	// Synthesizes speech synchronously: receive results after all text input
	// has been processed.
	SynthesizeSpeech(context.Context, *SynthesizeSpeechRequest) (*SynthesizeSpeechResponse, error)
}

func RegisterTextToSpeechServer(s *grpc.Server, srv TextToSpeechServer) {
	s.RegisterService(&_TextToSpeech_serviceDesc, srv)
}

func _TextToSpeech_ListVoices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVoicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextToSpeechServer).ListVoices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.texttospeech.v1beta1.TextToSpeech/ListVoices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextToSpeechServer).ListVoices(ctx, req.(*ListVoicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TextToSpeech_SynthesizeSpeech_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SynthesizeSpeechRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextToSpeechServer).SynthesizeSpeech(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.texttospeech.v1beta1.TextToSpeech/SynthesizeSpeech",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextToSpeechServer).SynthesizeSpeech(ctx, req.(*SynthesizeSpeechRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TextToSpeech_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.cloud.texttospeech.v1beta1.TextToSpeech",
	HandlerType: (*TextToSpeechServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListVoices",
			Handler:    _TextToSpeech_ListVoices_Handler,
		},
		{
			MethodName: "SynthesizeSpeech",
			Handler:    _TextToSpeech_SynthesizeSpeech_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/cloud/texttospeech/v1beta1/cloud_tts.proto",
}

func init() {
	proto.RegisterFile("google/cloud/texttospeech/v1beta1/cloud_tts.proto", fileDescriptor_dfc2f1f47f5721c8)
}

var fileDescriptor_dfc2f1f47f5721c8 = []byte{
	// 875 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0xee, 0xd8, 0x71, 0xda, 0x1e, 0xaf, 0x1d, 0x67, 0x14, 0x51, 0x13, 0x51, 0x94, 0x6e, 0xa8,
	0x64, 0x45, 0xc8, 0xc6, 0x2e, 0xbf, 0xe9, 0x05, 0x38, 0xf6, 0xd6, 0xb5, 0xe4, 0x3f, 0x66, 0x93,
	0x54, 0xe2, 0x66, 0x35, 0xd9, 0x9d, 0x6c, 0x56, 0xac, 0x67, 0x96, 0x9d, 0x71, 0x54, 0x7a, 0x89,
	0xb8, 0xe6, 0x02, 0x9e, 0x00, 0x89, 0x07, 0xe0, 0x59, 0x40, 0xe2, 0x09, 0x78, 0x08, 0x2e, 0xd1,
	0xce, 0x6c, 0x52, 0x27, 0x41, 0xd4, 0xe1, 0x82, 0xbb, 0x99, 0x73, 0xf6, 0x3b, 0xf3, 0x9d, 0x6f,
	0xbe, 0x3d, 0x03, 0xed, 0x50, 0x88, 0x30, 0x66, 0x2d, 0x3f, 0x16, 0x8b, 0xa0, 0xa5, 0xd8, 0x4b,
	0xa5, 0x84, 0x4c, 0x18, 0xf3, 0xcf, 0x5a, 0xe7, 0xed, 0x13, 0xa6, 0x68, 0xdb, 0xa4, 0x3c, 0xa5,
	0x64, 0x33, 0x49, 0x85, 0x12, 0xf8, 0x91, 0x81, 0x34, 0x75, 0xbc, 0xb9, 0x0c, 0x69, 0xe6, 0x90,
	0xed, 0x77, 0xf2, 0xaa, 0x34, 0x89, 0x5a, 0x94, 0x73, 0xa1, 0xa8, 0x8a, 0x04, 0xcf, 0x0b, 0xd8,
	0x9f, 0xc2, 0xe6, 0x28, 0x92, 0xea, 0x58, 0x44, 0x3e, 0x93, 0x84, 0x7d, 0xb3, 0x60, 0x52, 0xe1,
	0x5d, 0xa8, 0xc4, 0x94, 0x87, 0x0b, 0x1a, 0x32, 0xcf, 0x17, 0x01, 0xab, 0xa3, 0x1d, 0xd4, 0xb8,
	0x4f, 0xac, 0x8b, 0x60, 0x4f, 0x04, 0xcc, 0x3e, 0x06, 0xbc, 0x8c, 0x94, 0x89, 0xe0, 0x92, 0xe1,
	0x2f, 0x60, 0xfd, 0x5c, 0x47, 0xea, 0x68, 0xa7, 0xd8, 0x28, 0x77, 0x1a, 0xcd, 0x37, 0x32, 0x6c,
	0xea, 0x12, 0x24, 0xc7, 0xd9, 0xbf, 0x21, 0x28, 0xe9, 0x08, 0x7e, 0x0c, 0xd5, 0x2b, 0x34, 0x4c,
	0xcd, 0xfb, 0xa4, 0xb2, 0xcc, 0x43, 0x62, 0x0c, 0x6b, 0x9c, 0xce, 0x59, 0xbd, 0xa0, 0x49, 0xea,
	0x35, 0x76, 0xa1, 0x2c, 0xe5, 0x3c, 0xf6, 0x42, 0xc6, 0x03, 0x96, 0xd6, 0x8b, 0x3b, 0xa8, 0x51,
	0xed, 0x74, 0x56, 0xe0, 0xe2, 0xca, 0x79, 0xac, 0x4f, 0x1f, 0x68, 0x24, 0x81, 0xac, 0x8c, 0x59,
	0xe3, 0xcf, 0xe0, 0x6d, 0x4e, 0xd5, 0x22, 0xa5, 0xb1, 0x27, 0xe9, 0x3c, 0x89, 0x99, 0x97, 0x52,
	0xc5, 0xbc, 0x33, 0x96, 0xaa, 0x57, 0xf5, 0xb5, 0x1d, 0xd4, 0x28, 0x91, 0xb7, 0xf2, 0x0f, 0x5c,
	0x9d, 0x27, 0x54, 0xb1, 0xe7, 0x59, 0xd6, 0xfe, 0xbe, 0x00, 0x0f, 0xdc, 0x6f, 0xb9, 0x3a, 0x63,
	0x32, 0x7a, 0xc5, 0x5c, 0x7d, 0xe6, 0x85, 0xda, 0x03, 0x28, 0x45, 0x3c, 0x59, 0x28, 0xad, 0x72,
	0xb9, 0xd3, 0x5e, 0x85, 0x65, 0x5e, 0x4a, 0x0e, 0x33, 0x20, 0x31, 0x78, 0x3c, 0x86, 0x92, 0xd6,
	0x50, 0x2b, 0x51, 0xee, 0x7c, 0xb2, 0xaa, 0xf4, 0x2e, 0x8b, 0x99, 0x9f, 0x99, 0x62, 0x46, 0x53,
	0x3a, 0x97, 0xc4, 0x54, 0xc1, 0x5f, 0x82, 0x45, 0x17, 0x41, 0x24, 0x3c, 0x5f, 0xf0, 0xd3, 0x28,
	0xd4, 0x22, 0x96, 0x3b, 0xcd, 0x15, 0xaa, 0x76, 0x33, 0x58, 0x4f, 0xa3, 0x48, 0x99, 0xbe, 0xde,
	0xd8, 0x23, 0xa8, 0x5e, 0xa5, 0x8e, 0xb7, 0x60, 0x2d, 0x2b, 0x61, 0x1c, 0xf6, 0xfc, 0x0e, 0xd1,
	0xbb, 0x2c, 0x9a, 0xe9, 0x6e, 0xae, 0x34, 0x8b, 0x66, 0xbb, 0x83, 0x2a, 0x58, 0xba, 0x51, 0x4f,
	0x8a, 0x45, 0xea, 0x33, 0xfb, 0x17, 0x04, 0x5b, 0xff, 0xd4, 0xc0, 0x4a, 0xfe, 0xfd, 0xdf, 0x6c,
	0x63, 0xff, 0x5c, 0x80, 0xf2, 0x92, 0x22, 0xf8, 0x05, 0x54, 0x8d, 0xae, 0x8c, 0xfb, 0x22, 0x88,
	0x78, 0xa8, 0xe9, 0x55, 0x3b, 0x1f, 0xac, 0xaa, 0xac, 0x93, 0xe3, 0x48, 0x85, 0x2e, 0x6f, 0xb3,
	0xb6, 0x65, 0xc2, 0xe8, 0xd7, 0x11, 0x0f, 0xb5, 0x33, 0x75, 0x6b, 0x88, 0x58, 0x17, 0xc1, 0xcc,
	0x8e, 0x78, 0x0b, 0x4a, 0x49, 0xa4, 0xfc, 0x33, 0xdd, 0x1c, 0x22, 0x66, 0x83, 0xdf, 0x83, 0xea,
	0xb9, 0x88, 0x17, 0x73, 0xe6, 0x85, 0x34, 0xe2, 0x5e, 0x70, 0xa2, 0xfd, 0x8c, 0x88, 0x65, 0xa2,
	0x03, 0x1a, 0xf1, 0xfe, 0x09, 0xde, 0x83, 0xcd, 0x9b, 0xc6, 0x2f, 0x69, 0xe3, 0x6f, 0xc8, 0xab,
	0x8e, 0xc7, 0xef, 0x03, 0x66, 0xa7, 0xa7, 0xcc, 0x57, 0xd2, 0x4b, 0x52, 0x71, 0x1a, 0xc5, 0xcc,
	0x8b, 0x82, 0xfa, 0xba, 0xfe, 0x81, 0x6b, 0x79, 0x66, 0x66, 0x12, 0xc3, 0xc0, 0xfe, 0x1c, 0xea,
	0x37, 0x7f, 0x8f, 0x7c, 0xa4, 0xec, 0x42, 0xe5, 0xd2, 0x87, 0x8a, 0x71, 0xe3, 0x15, 0x8b, 0x58,
	0x17, 0xc6, 0xca, 0x62, 0x7b, 0x2f, 0x60, 0xe3, 0xda, 0x1d, 0xe0, 0x47, 0xf0, 0xd0, 0x75, 0xc7,
	0x23, 0xef, 0x78, 0x3a, 0xec, 0x39, 0xde, 0xc0, 0x99, 0xf4, 0x1d, 0xe2, 0x1d, 0x4d, 0xdc, 0x99,
	0xd3, 0x1b, 0x3e, 0x1b, 0x3a, 0xfd, 0xda, 0x1d, 0x7c, 0x0f, 0xd6, 0xc6, 0xdd, 0x91, 0x53, 0x43,
	0x18, 0x60, 0xfd, 0x99, 0xa3, 0xd7, 0x05, 0x5c, 0x86, 0xbb, 0x13, 0xe7, 0xe8, 0x90, 0x74, 0x47,
	0xb5, 0xe2, 0xde, 0x21, 0x54, 0xae, 0x88, 0x8e, 0xdf, 0x85, 0xed, 0xee, 0x51, 0x7f, 0x38, 0xf5,
	0x9c, 0x49, 0x6f, 0xda, 0x1f, 0x4e, 0x06, 0xd7, 0x6a, 0x5a, 0x70, 0x6f, 0x34, 0x9c, 0x38, 0x5d,
	0xd2, 0xfe, 0xb8, 0x86, 0xf0, 0x5d, 0x28, 0x8e, 0x67, 0x4f, 0x6a, 0x85, 0x2c, 0x3c, 0x1d, 0x0c,
	0xbc, 0xe9, 0xec, 0xc8, 0xad, 0x15, 0x3b, 0x7f, 0x14, 0xc0, 0x3a, 0x64, 0x2f, 0xd5, 0xa1, 0x30,
	0xcd, 0xe2, 0x1f, 0x11, 0xc0, 0xeb, 0x71, 0x8a, 0x3f, 0x5c, 0xc1, 0x0b, 0x37, 0xe6, 0xf6, 0xf6,
	0x47, 0xb7, 0x44, 0x19, 0x81, 0xed, 0x07, 0xdf, 0xfd, 0xfe, 0xe7, 0x4f, 0x85, 0x4d, 0xbc, 0x71,
	0xf9, 0xcc, 0x98, 0x51, 0x8c, 0x7f, 0x45, 0x50, 0xbb, 0x7e, 0x2d, 0x78, 0xff, 0x16, 0xf3, 0xe9,
	0xda, 0xa8, 0xdb, 0x7e, 0xfa, 0x9f, 0xb0, 0x39, 0xcd, 0x5d, 0x4d, 0xf3, 0xa1, 0x5d, 0xbf, 0xa4,
	0x99, 0xe1, 0xf7, 0xe5, 0xe5, 0xf7, 0xfb, 0x68, 0xef, 0xe0, 0x07, 0x04, 0x8f, 0x7d, 0x31, 0x7f,
	0xf3, 0x39, 0x07, 0x9b, 0xcb, 0xfa, 0xcf, 0xb2, 0xc7, 0x70, 0x86, 0xbe, 0x1a, 0xe7, 0xb8, 0x50,
	0x64, 0xb3, 0xa2, 0x29, 0xd2, 0xb0, 0x15, 0x32, 0xae, 0x9f, 0xca, 0x96, 0x49, 0xd1, 0x24, 0x92,
	0xff, 0xf2, 0x42, 0x3f, 0x5d, 0x0e, 0xfe, 0x85, 0xd0, 0xc9, 0xba, 0x06, 0x3f, 0xf9, 0x3b, 0x00,
	0x00, 0xff, 0xff, 0x26, 0x28, 0xff, 0x42, 0xdd, 0x07, 0x00, 0x00,
}
