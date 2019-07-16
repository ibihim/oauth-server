// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/metric.proto

/*
Package metric is a generated protocol buffer package.

It is generated from these files:
	google/api/metric.proto

It has these top-level messages:
	MetricDescriptor
	Metric
*/
package metric

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_api "google.golang.org/genproto/googleapis/api/label"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The kind of measurement. It describes how the data is reported.
type MetricDescriptor_MetricKind int32

const (
	// Do not use this default value.
	MetricDescriptor_METRIC_KIND_UNSPECIFIED MetricDescriptor_MetricKind = 0
	// An instantaneous measurement of a value.
	MetricDescriptor_GAUGE MetricDescriptor_MetricKind = 1
	// The change in a value during a time interval.
	MetricDescriptor_DELTA MetricDescriptor_MetricKind = 2
	// A value accumulated over a time interval.  Cumulative
	// measurements in a time series should have the same start time
	// and increasing end times, until an event resets the cumulative
	// value to zero and sets a new start time for the following
	// points.
	MetricDescriptor_CUMULATIVE MetricDescriptor_MetricKind = 3
)

var MetricDescriptor_MetricKind_name = map[int32]string{
	0: "METRIC_KIND_UNSPECIFIED",
	1: "GAUGE",
	2: "DELTA",
	3: "CUMULATIVE",
}
var MetricDescriptor_MetricKind_value = map[string]int32{
	"METRIC_KIND_UNSPECIFIED": 0,
	"GAUGE":                   1,
	"DELTA":                   2,
	"CUMULATIVE":              3,
}

func (x MetricDescriptor_MetricKind) String() string {
	return proto.EnumName(MetricDescriptor_MetricKind_name, int32(x))
}
func (MetricDescriptor_MetricKind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 0}
}

// The value type of a metric.
type MetricDescriptor_ValueType int32

const (
	// Do not use this default value.
	MetricDescriptor_VALUE_TYPE_UNSPECIFIED MetricDescriptor_ValueType = 0
	// The value is a boolean.
	// This value type can be used only if the metric kind is `GAUGE`.
	MetricDescriptor_BOOL MetricDescriptor_ValueType = 1
	// The value is a signed 64-bit integer.
	MetricDescriptor_INT64 MetricDescriptor_ValueType = 2
	// The value is a double precision floating point number.
	MetricDescriptor_DOUBLE MetricDescriptor_ValueType = 3
	// The value is a text string.
	// This value type can be used only if the metric kind is `GAUGE`.
	MetricDescriptor_STRING MetricDescriptor_ValueType = 4
	// The value is a [`Distribution`][google.api.Distribution].
	MetricDescriptor_DISTRIBUTION MetricDescriptor_ValueType = 5
	// The value is money.
	MetricDescriptor_MONEY MetricDescriptor_ValueType = 6
)

var MetricDescriptor_ValueType_name = map[int32]string{
	0: "VALUE_TYPE_UNSPECIFIED",
	1: "BOOL",
	2: "INT64",
	3: "DOUBLE",
	4: "STRING",
	5: "DISTRIBUTION",
	6: "MONEY",
}
var MetricDescriptor_ValueType_value = map[string]int32{
	"VALUE_TYPE_UNSPECIFIED": 0,
	"BOOL":         1,
	"INT64":        2,
	"DOUBLE":       3,
	"STRING":       4,
	"DISTRIBUTION": 5,
	"MONEY":        6,
}

func (x MetricDescriptor_ValueType) String() string {
	return proto.EnumName(MetricDescriptor_ValueType_name, int32(x))
}
func (MetricDescriptor_ValueType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 1}
}

// Defines a metric type and its schema. Once a metric descriptor is created,
// deleting or altering it stops data collection and makes the metric type's
// existing data unusable.
type MetricDescriptor struct {
	// The resource name of the metric descriptor. Depending on the
	// implementation, the name typically includes: (1) the parent resource name
	// that defines the scope of the metric type or of its data; and (2) the
	// metric's URL-encoded type, which also appears in the `type` field of this
	// descriptor. For example, following is the resource name of a custom
	// metric within the GCP project `my-project-id`:
	//
	//     "projects/my-project-id/metricDescriptors/custom.googleapis.com%2Finvoice%2Fpaid%2Famount"
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The metric type, including its DNS name prefix. The type is not
	// URL-encoded.  All user-defined custom metric types have the DNS name
	// `custom.googleapis.com`.  Metric types should use a natural hierarchical
	// grouping. For example:
	//
	//     "custom.googleapis.com/invoice/paid/amount"
	//     "appengine.googleapis.com/http/server/response_latencies"
	Type string `protobuf:"bytes,8,opt,name=type" json:"type,omitempty"`
	// The set of labels that can be used to describe a specific
	// instance of this metric type. For example, the
	// `appengine.googleapis.com/http/server/response_latencies` metric
	// type has a label for the HTTP response code, `response_code`, so
	// you can look at latencies for successful responses or just
	// for responses that failed.
	Labels []*google_api.LabelDescriptor `protobuf:"bytes,2,rep,name=labels" json:"labels,omitempty"`
	// Whether the metric records instantaneous values, changes to a value, etc.
	// Some combinations of `metric_kind` and `value_type` might not be supported.
	MetricKind MetricDescriptor_MetricKind `protobuf:"varint,3,opt,name=metric_kind,json=metricKind,enum=google.api.MetricDescriptor_MetricKind" json:"metric_kind,omitempty"`
	// Whether the measurement is an integer, a floating-point number, etc.
	// Some combinations of `metric_kind` and `value_type` might not be supported.
	ValueType MetricDescriptor_ValueType `protobuf:"varint,4,opt,name=value_type,json=valueType,enum=google.api.MetricDescriptor_ValueType" json:"value_type,omitempty"`
	// The unit in which the metric value is reported. It is only applicable
	// if the `value_type` is `INT64`, `DOUBLE`, or `DISTRIBUTION`. The
	// supported units are a subset of [The Unified Code for Units of
	// Measure](http://unitsofmeasure.org/ucum.html) standard:
	//
	// **Basic units (UNIT)**
	//
	// * `bit`   bit
	// * `By`    byte
	// * `s`     second
	// * `min`   minute
	// * `h`     hour
	// * `d`     day
	//
	// **Prefixes (PREFIX)**
	//
	// * `k`     kilo    (10**3)
	// * `M`     mega    (10**6)
	// * `G`     giga    (10**9)
	// * `T`     tera    (10**12)
	// * `P`     peta    (10**15)
	// * `E`     exa     (10**18)
	// * `Z`     zetta   (10**21)
	// * `Y`     yotta   (10**24)
	// * `m`     milli   (10**-3)
	// * `u`     micro   (10**-6)
	// * `n`     nano    (10**-9)
	// * `p`     pico    (10**-12)
	// * `f`     femto   (10**-15)
	// * `a`     atto    (10**-18)
	// * `z`     zepto   (10**-21)
	// * `y`     yocto   (10**-24)
	// * `Ki`    kibi    (2**10)
	// * `Mi`    mebi    (2**20)
	// * `Gi`    gibi    (2**30)
	// * `Ti`    tebi    (2**40)
	//
	// **Grammar**
	//
	// The grammar includes the dimensionless unit `1`, such as `1/s`.
	//
	// The grammar also includes these connectors:
	//
	// * `/`    division (as an infix operator, e.g. `1/s`).
	// * `.`    multiplication (as an infix operator, e.g. `GBy.d`)
	//
	// The grammar for a unit is as follows:
	//
	//     Expression = Component { "." Component } { "/" Component } ;
	//
	//     Component = [ PREFIX ] UNIT [ Annotation ]
	//               | Annotation
	//               | "1"
	//               ;
	//
	//     Annotation = "{" NAME "}" ;
	//
	// Notes:
	//
	// * `Annotation` is just a comment if it follows a `UNIT` and is
	//    equivalent to `1` if it is used alone. For examples,
	//    `{requests}/s == 1/s`, `By{transmitted}/s == By/s`.
	// * `NAME` is a sequence of non-blank printable ASCII characters not
	//    containing '{' or '}'.
	Unit string `protobuf:"bytes,5,opt,name=unit" json:"unit,omitempty"`
	// A detailed description of the metric, which can be used in documentation.
	Description string `protobuf:"bytes,6,opt,name=description" json:"description,omitempty"`
	// A concise name for the metric, which can be displayed in user interfaces.
	// Use sentence case without an ending period, for example "Request count".
	DisplayName string `protobuf:"bytes,7,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
}

func (m *MetricDescriptor) Reset()                    { *m = MetricDescriptor{} }
func (m *MetricDescriptor) String() string            { return proto.CompactTextString(m) }
func (*MetricDescriptor) ProtoMessage()               {}
func (*MetricDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MetricDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetricDescriptor) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MetricDescriptor) GetLabels() []*google_api.LabelDescriptor {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *MetricDescriptor) GetMetricKind() MetricDescriptor_MetricKind {
	if m != nil {
		return m.MetricKind
	}
	return MetricDescriptor_METRIC_KIND_UNSPECIFIED
}

func (m *MetricDescriptor) GetValueType() MetricDescriptor_ValueType {
	if m != nil {
		return m.ValueType
	}
	return MetricDescriptor_VALUE_TYPE_UNSPECIFIED
}

func (m *MetricDescriptor) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *MetricDescriptor) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *MetricDescriptor) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

// A specific metric, identified by specifying values for all of the
// labels of a [`MetricDescriptor`][google.api.MetricDescriptor].
type Metric struct {
	// An existing metric type, see [google.api.MetricDescriptor][google.api.MetricDescriptor].
	// For example, `custom.googleapis.com/invoice/paid/amount`.
	Type string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	// The set of label values that uniquely identify this metric. All
	// labels listed in the `MetricDescriptor` must be assigned values.
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Metric) Reset()                    { *m = Metric{} }
func (m *Metric) String() string            { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()               {}
func (*Metric) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Metric) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Metric) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func init() {
	proto.RegisterType((*MetricDescriptor)(nil), "google.api.MetricDescriptor")
	proto.RegisterType((*Metric)(nil), "google.api.Metric")
	proto.RegisterEnum("google.api.MetricDescriptor_MetricKind", MetricDescriptor_MetricKind_name, MetricDescriptor_MetricKind_value)
	proto.RegisterEnum("google.api.MetricDescriptor_ValueType", MetricDescriptor_ValueType_name, MetricDescriptor_ValueType_value)
}

func init() { proto.RegisterFile("google/api/metric.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 506 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0x4d, 0x6f, 0xda, 0x40,
	0x10, 0xad, 0x3f, 0x70, 0xc3, 0x10, 0xa1, 0xd5, 0xaa, 0x4a, 0x2c, 0x22, 0x55, 0x94, 0x43, 0xcb,
	0x09, 0xa4, 0xa4, 0x4a, 0xbf, 0x4e, 0x80, 0xb7, 0xd4, 0x8a, 0xb1, 0x91, 0x63, 0x23, 0xa5, 0x17,
	0xcb, 0x81, 0x95, 0x65, 0xc5, 0xd8, 0xae, 0x71, 0x22, 0xf9, 0x57, 0xf4, 0x17, 0xf4, 0xd2, 0x5f,
	0x5a, 0xed, 0xae, 0x03, 0x16, 0x95, 0x72, 0xe2, 0xed, 0x9b, 0x37, 0x6f, 0x67, 0x96, 0x67, 0x38,
	0x8f, 0xb2, 0x2c, 0x4a, 0xe8, 0x38, 0xcc, 0xe3, 0xf1, 0x96, 0x96, 0x45, 0xbc, 0x1e, 0xe5, 0x45,
	0x56, 0x66, 0x18, 0x44, 0x61, 0x14, 0xe6, 0x71, 0xef, 0xac, 0x21, 0x4a, 0xc2, 0x7b, 0x9a, 0x08,
	0xcd, 0xe0, 0x8f, 0x0a, 0x68, 0xc1, 0x9b, 0x0c, 0xba, 0x5b, 0x17, 0x71, 0x5e, 0x66, 0x05, 0xc6,
	0xa0, 0xa6, 0xe1, 0x96, 0xea, 0x52, 0x5f, 0x1a, 0xb6, 0x5d, 0x8e, 0x19, 0x57, 0x56, 0x39, 0xd5,
	0x4f, 0x04, 0xc7, 0x30, 0xbe, 0x02, 0x8d, 0x7b, 0xed, 0x74, 0xb9, 0xaf, 0x0c, 0x3b, 0x97, 0x17,
	0xa3, 0xc3, 0x8d, 0x23, 0x8b, 0x55, 0x0e, 0xa6, 0x6e, 0x2d, 0xc5, 0x3f, 0xa0, 0x23, 0xa6, 0x0c,
	0x1e, 0xe2, 0x74, 0xa3, 0x2b, 0x7d, 0x69, 0xd8, 0xbd, 0xfc, 0xd0, 0xec, 0x3c, 0x9e, 0xa7, 0x26,
	0x6e, 0xe2, 0x74, 0xe3, 0xc2, 0x76, 0x8f, 0x31, 0x01, 0x78, 0x0a, 0x93, 0x47, 0x1a, 0xf0, 0xc1,
	0x54, 0x6e, 0xf4, 0xfe, 0x45, 0xa3, 0x15, 0x93, 0x7b, 0x55, 0x4e, 0xdd, 0xf6, 0xd3, 0x33, 0x64,
	0x9b, 0x3d, 0xa6, 0x71, 0xa9, 0xb7, 0xc4, 0x66, 0x0c, 0xe3, 0x3e, 0x74, 0x36, 0x75, 0x5b, 0x9c,
	0xa5, 0xba, 0xc6, 0x4b, 0x4d, 0x0a, 0xbf, 0x83, 0xd3, 0x4d, 0xbc, 0xcb, 0x93, 0xb0, 0x0a, 0xf8,
	0x5b, 0xbd, 0xae, 0x25, 0x82, 0xb3, 0xc3, 0x2d, 0x1d, 0x38, 0x00, 0x87, 0xc9, 0xf1, 0x05, 0x9c,
	0x2f, 0x88, 0xe7, 0x9a, 0xb3, 0xe0, 0xc6, 0xb4, 0x8d, 0xc0, 0xb7, 0x6f, 0x97, 0x64, 0x66, 0x7e,
	0x37, 0x89, 0x81, 0x5e, 0xe1, 0x36, 0xb4, 0xe6, 0x13, 0x7f, 0x4e, 0x90, 0xc4, 0xa0, 0x41, 0x2c,
	0x6f, 0x82, 0x64, 0xdc, 0x05, 0x98, 0xf9, 0x0b, 0xdf, 0x9a, 0x78, 0xe6, 0x8a, 0x20, 0x65, 0xf0,
	0x0b, 0xda, 0xfb, 0x0d, 0x70, 0x0f, 0xce, 0x56, 0x13, 0xcb, 0x27, 0x81, 0x77, 0xb7, 0x24, 0x47,
	0x76, 0x27, 0xa0, 0x4e, 0x1d, 0xc7, 0x12, 0x6e, 0xa6, 0xed, 0x5d, 0x7f, 0x44, 0x32, 0x06, 0xd0,
	0x0c, 0xc7, 0x9f, 0x5a, 0x04, 0x29, 0x0c, 0xdf, 0x7a, 0xae, 0x69, 0xcf, 0x91, 0x8a, 0x11, 0x9c,
	0x1a, 0x26, 0x3b, 0x4d, 0x7d, 0xcf, 0x74, 0x6c, 0xd4, 0x62, 0x4d, 0x0b, 0xc7, 0x26, 0x77, 0x48,
	0x1b, 0xfc, 0x96, 0x40, 0x13, 0x4b, 0xec, 0x13, 0xa0, 0x34, 0x12, 0x70, 0x7d, 0x94, 0x80, 0xb7,
	0xff, 0x3f, 0xbf, 0x08, 0xc2, 0x8e, 0xa4, 0x65, 0x51, 0x3d, 0x87, 0xa0, 0xf7, 0x05, 0x3a, 0x0d,
	0x1a, 0x23, 0x50, 0x1e, 0x68, 0x55, 0xe7, 0x8d, 0x41, 0xfc, 0x06, 0x5a, 0xfc, 0x1f, 0xd2, 0x65,
	0xce, 0x89, 0xc3, 0x57, 0xf9, 0xb3, 0x34, 0x0d, 0xa0, 0xbb, 0xce, 0xb6, 0x8d, 0x7b, 0xa6, 0x1d,
	0x71, 0xd1, 0x92, 0x05, 0x7a, 0x29, 0xfd, 0xfc, 0x54, 0x97, 0xa2, 0x2c, 0x09, 0xd3, 0x68, 0x94,
	0x15, 0xd1, 0x38, 0xa2, 0x29, 0x8f, 0xfb, 0x58, 0x94, 0xc2, 0x3c, 0xde, 0x35, 0x3e, 0x97, 0x6f,
	0xe2, 0xe7, 0xaf, 0xac, 0xce, 0x27, 0x4b, 0xf3, 0x5e, 0xe3, 0xd2, 0xab, 0x7f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x18, 0x04, 0x05, 0x82, 0x58, 0x03, 0x00, 0x00,
}
