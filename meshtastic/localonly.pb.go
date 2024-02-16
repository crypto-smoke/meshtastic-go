// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: meshtastic/localonly.proto

package generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LocalConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The part of the config that is specific to the Device
	Device *Config_DeviceConfig `protobuf:"bytes,1,opt,name=device,proto3" json:"device,omitempty"`
	// The part of the config that is specific to the GPS Position
	Position *Config_PositionConfig `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	// The part of the config that is specific to the Power settings
	Power *Config_PowerConfig `protobuf:"bytes,3,opt,name=power,proto3" json:"power,omitempty"`
	// The part of the config that is specific to the Wifi Settings
	Network *Config_NetworkConfig `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	// The part of the config that is specific to the Display
	Display *Config_DisplayConfig `protobuf:"bytes,5,opt,name=display,proto3" json:"display,omitempty"`
	// The part of the config that is specific to the Lora Radio
	Lora *Config_LoRaConfig `protobuf:"bytes,6,opt,name=lora,proto3" json:"lora,omitempty"`
	// The part of the config that is specific to the Bluetooth settings
	Bluetooth *Config_BluetoothConfig `protobuf:"bytes,7,opt,name=bluetooth,proto3" json:"bluetooth,omitempty"`
	// A version integer used to invalidate old save files when we make
	// incompatible changes This integer is set at build time and is private to
	// NodeDB.cpp in the device code.
	Version uint32 `protobuf:"varint,8,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *LocalConfig) Reset() {
	*x = LocalConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meshtastic_localonly_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalConfig) ProtoMessage() {}

func (x *LocalConfig) ProtoReflect() protoreflect.Message {
	mi := &file_meshtastic_localonly_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalConfig.ProtoReflect.Descriptor instead.
func (*LocalConfig) Descriptor() ([]byte, []int) {
	return file_meshtastic_localonly_proto_rawDescGZIP(), []int{0}
}

func (x *LocalConfig) GetDevice() *Config_DeviceConfig {
	if x != nil {
		return x.Device
	}
	return nil
}

func (x *LocalConfig) GetPosition() *Config_PositionConfig {
	if x != nil {
		return x.Position
	}
	return nil
}

func (x *LocalConfig) GetPower() *Config_PowerConfig {
	if x != nil {
		return x.Power
	}
	return nil
}

func (x *LocalConfig) GetNetwork() *Config_NetworkConfig {
	if x != nil {
		return x.Network
	}
	return nil
}

func (x *LocalConfig) GetDisplay() *Config_DisplayConfig {
	if x != nil {
		return x.Display
	}
	return nil
}

func (x *LocalConfig) GetLora() *Config_LoRaConfig {
	if x != nil {
		return x.Lora
	}
	return nil
}

func (x *LocalConfig) GetBluetooth() *Config_BluetoothConfig {
	if x != nil {
		return x.Bluetooth
	}
	return nil
}

func (x *LocalConfig) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type LocalModuleConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The part of the config that is specific to the MQTT module
	Mqtt *ModuleConfig_MQTTConfig `protobuf:"bytes,1,opt,name=mqtt,proto3" json:"mqtt,omitempty"`
	// The part of the config that is specific to the Serial module
	Serial *ModuleConfig_SerialConfig `protobuf:"bytes,2,opt,name=serial,proto3" json:"serial,omitempty"`
	// The part of the config that is specific to the ExternalNotification module
	ExternalNotification *ModuleConfig_ExternalNotificationConfig `protobuf:"bytes,3,opt,name=external_notification,json=externalNotification,proto3" json:"external_notification,omitempty"`
	// The part of the config that is specific to the Store & Forward module
	StoreForward *ModuleConfig_StoreForwardConfig `protobuf:"bytes,4,opt,name=store_forward,json=storeForward,proto3" json:"store_forward,omitempty"`
	// The part of the config that is specific to the RangeTest module
	RangeTest *ModuleConfig_RangeTestConfig `protobuf:"bytes,5,opt,name=range_test,json=rangeTest,proto3" json:"range_test,omitempty"`
	// The part of the config that is specific to the Telemetry module
	Telemetry *ModuleConfig_TelemetryConfig `protobuf:"bytes,6,opt,name=telemetry,proto3" json:"telemetry,omitempty"`
	// The part of the config that is specific to the Canned Message module
	CannedMessage *ModuleConfig_CannedMessageConfig `protobuf:"bytes,7,opt,name=canned_message,json=cannedMessage,proto3" json:"canned_message,omitempty"`
	// The part of the config that is specific to the Audio module
	Audio *ModuleConfig_AudioConfig `protobuf:"bytes,9,opt,name=audio,proto3" json:"audio,omitempty"`
	// The part of the config that is specific to the Remote Hardware module
	RemoteHardware *ModuleConfig_RemoteHardwareConfig `protobuf:"bytes,10,opt,name=remote_hardware,json=remoteHardware,proto3" json:"remote_hardware,omitempty"`
	// The part of the config that is specific to the Neighbor Info module
	NeighborInfo *ModuleConfig_NeighborInfoConfig `protobuf:"bytes,11,opt,name=neighbor_info,json=neighborInfo,proto3" json:"neighbor_info,omitempty"`
	// The part of the config that is specific to the Ambient Lighting module
	AmbientLighting *ModuleConfig_AmbientLightingConfig `protobuf:"bytes,12,opt,name=ambient_lighting,json=ambientLighting,proto3" json:"ambient_lighting,omitempty"`
	// The part of the config that is specific to the Detection Sensor module
	DetectionSensor *ModuleConfig_DetectionSensorConfig `protobuf:"bytes,13,opt,name=detection_sensor,json=detectionSensor,proto3" json:"detection_sensor,omitempty"`
	// Paxcounter Config
	Paxcounter *ModuleConfig_PaxcounterConfig `protobuf:"bytes,14,opt,name=paxcounter,proto3" json:"paxcounter,omitempty"`
	// A version integer used to invalidate old save files when we make
	// incompatible changes This integer is set at build time and is private to
	// NodeDB.cpp in the device code.
	Version uint32 `protobuf:"varint,8,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *LocalModuleConfig) Reset() {
	*x = LocalModuleConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_meshtastic_localonly_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalModuleConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalModuleConfig) ProtoMessage() {}

func (x *LocalModuleConfig) ProtoReflect() protoreflect.Message {
	mi := &file_meshtastic_localonly_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalModuleConfig.ProtoReflect.Descriptor instead.
func (*LocalModuleConfig) Descriptor() ([]byte, []int) {
	return file_meshtastic_localonly_proto_rawDescGZIP(), []int{1}
}

func (x *LocalModuleConfig) GetMqtt() *ModuleConfig_MQTTConfig {
	if x != nil {
		return x.Mqtt
	}
	return nil
}

func (x *LocalModuleConfig) GetSerial() *ModuleConfig_SerialConfig {
	if x != nil {
		return x.Serial
	}
	return nil
}

func (x *LocalModuleConfig) GetExternalNotification() *ModuleConfig_ExternalNotificationConfig {
	if x != nil {
		return x.ExternalNotification
	}
	return nil
}

func (x *LocalModuleConfig) GetStoreForward() *ModuleConfig_StoreForwardConfig {
	if x != nil {
		return x.StoreForward
	}
	return nil
}

func (x *LocalModuleConfig) GetRangeTest() *ModuleConfig_RangeTestConfig {
	if x != nil {
		return x.RangeTest
	}
	return nil
}

func (x *LocalModuleConfig) GetTelemetry() *ModuleConfig_TelemetryConfig {
	if x != nil {
		return x.Telemetry
	}
	return nil
}

func (x *LocalModuleConfig) GetCannedMessage() *ModuleConfig_CannedMessageConfig {
	if x != nil {
		return x.CannedMessage
	}
	return nil
}

func (x *LocalModuleConfig) GetAudio() *ModuleConfig_AudioConfig {
	if x != nil {
		return x.Audio
	}
	return nil
}

func (x *LocalModuleConfig) GetRemoteHardware() *ModuleConfig_RemoteHardwareConfig {
	if x != nil {
		return x.RemoteHardware
	}
	return nil
}

func (x *LocalModuleConfig) GetNeighborInfo() *ModuleConfig_NeighborInfoConfig {
	if x != nil {
		return x.NeighborInfo
	}
	return nil
}

func (x *LocalModuleConfig) GetAmbientLighting() *ModuleConfig_AmbientLightingConfig {
	if x != nil {
		return x.AmbientLighting
	}
	return nil
}

func (x *LocalModuleConfig) GetDetectionSensor() *ModuleConfig_DetectionSensorConfig {
	if x != nil {
		return x.DetectionSensor
	}
	return nil
}

func (x *LocalModuleConfig) GetPaxcounter() *ModuleConfig_PaxcounterConfig {
	if x != nil {
		return x.Paxcounter
	}
	return nil
}

func (x *LocalModuleConfig) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_meshtastic_localonly_proto protoreflect.FileDescriptor

var file_meshtastic_localonly_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x6f, 0x6e, 0x6c, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x65,
	0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x1a, 0x17, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61,
	0x73, 0x74, 0x69, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc2, 0x03, 0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x37, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x08, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6d,
	0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x05, 0x70, 0x6f, 0x77,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74,
	0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x6f, 0x77,
	0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x05, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x12,
	0x3a, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x3a, 0x0a, 0x07, 0x64,
	0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6d,
	0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x07,
	0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x12, 0x31, 0x0a, 0x04, 0x6c, 0x6f, 0x72, 0x61, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4c, 0x6f, 0x52, 0x61, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x04, 0x6c, 0x6f, 0x72, 0x61, 0x12, 0x40, 0x0a, 0x09, 0x62, 0x6c,
	0x75, 0x65, 0x74, 0x6f, 0x6f, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x42, 0x6c, 0x75, 0x65, 0x74, 0x6f, 0x6f, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x09, 0x62, 0x6c, 0x75, 0x65, 0x74, 0x6f, 0x6f, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xae, 0x08, 0x0a, 0x11, 0x4c, 0x6f, 0x63, 0x61, 0x6c,
	0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x37, 0x0a, 0x04,
	0x6d, 0x71, 0x74, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x65, 0x73,
	0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4d, 0x51, 0x54, 0x54, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x04, 0x6d, 0x71, 0x74, 0x74, 0x12, 0x3d, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x12, 0x68, 0x0a, 0x15, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63,
	0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x45, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x14, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x50,
	0x0a, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x0c, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64,
	0x12, 0x47, 0x0a, 0x0a, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69,
	0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52,
	0x61, 0x6e, 0x67, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x09, 0x74, 0x65, 0x6c,
	0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6d,
	0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x12, 0x53, 0x0a, 0x0e, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6d, 0x65, 0x73, 0x68,
	0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x43, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0d, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3a, 0x0a, 0x05, 0x61, 0x75, 0x64, 0x69, 0x6f, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x41, 0x75, 0x64, 0x69, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x05, 0x61, 0x75, 0x64,
	0x69, 0x6f, 0x12, 0x56, 0x0a, 0x0f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x68, 0x61, 0x72,
	0x64, 0x77, 0x61, 0x72, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x6d, 0x65,
	0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x48, 0x61, 0x72, 0x64,
	0x77, 0x61, 0x72, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x12, 0x50, 0x0a, 0x0d, 0x6e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4e, 0x65, 0x69, 0x67,
	0x68, 0x62, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0c,
	0x6e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x59, 0x0a, 0x10,
	0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x69, 0x6e, 0x67,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x41, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x69, 0x6e, 0x67,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0f, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x4c,
	0x69, 0x67, 0x68, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x59, 0x0a, 0x10, 0x64, 0x65, 0x74, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2e, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x65, 0x74, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x0f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x6e, 0x73,
	0x6f, 0x72, 0x12, 0x49, 0x0a, 0x0a, 0x70, 0x61, 0x78, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x50, 0x61, 0x78, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x0a, 0x70, 0x61, 0x78, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x64, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67,
	0x65, 0x65, 0x6b, 0x73, 0x76, 0x69, 0x6c, 0x6c, 0x65, 0x2e, 0x6d, 0x65, 0x73, 0x68, 0x42, 0x0f,
	0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x4f, 0x6e, 0x6c, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x5a,
	0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x73, 0x68,
	0x74, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0xaa, 0x02, 0x14, 0x4d, 0x65, 0x73, 0x68, 0x74, 0x61, 0x73, 0x74, 0x69, 0x63,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0xba, 0x02, 0x00, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_meshtastic_localonly_proto_rawDescOnce sync.Once
	file_meshtastic_localonly_proto_rawDescData = file_meshtastic_localonly_proto_rawDesc
)

func file_meshtastic_localonly_proto_rawDescGZIP() []byte {
	file_meshtastic_localonly_proto_rawDescOnce.Do(func() {
		file_meshtastic_localonly_proto_rawDescData = protoimpl.X.CompressGZIP(file_meshtastic_localonly_proto_rawDescData)
	})
	return file_meshtastic_localonly_proto_rawDescData
}

var file_meshtastic_localonly_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_meshtastic_localonly_proto_goTypes = []interface{}{
	(*LocalConfig)(nil),                             // 0: meshtastic.LocalConfig
	(*LocalModuleConfig)(nil),                       // 1: meshtastic.LocalModuleConfig
	(*Config_DeviceConfig)(nil),                     // 2: meshtastic.Config.DeviceConfig
	(*Config_PositionConfig)(nil),                   // 3: meshtastic.Config.PositionConfig
	(*Config_PowerConfig)(nil),                      // 4: meshtastic.Config.PowerConfig
	(*Config_NetworkConfig)(nil),                    // 5: meshtastic.Config.NetworkConfig
	(*Config_DisplayConfig)(nil),                    // 6: meshtastic.Config.DisplayConfig
	(*Config_LoRaConfig)(nil),                       // 7: meshtastic.Config.LoRaConfig
	(*Config_BluetoothConfig)(nil),                  // 8: meshtastic.Config.BluetoothConfig
	(*ModuleConfig_MQTTConfig)(nil),                 // 9: meshtastic.ModuleConfig.MQTTConfig
	(*ModuleConfig_SerialConfig)(nil),               // 10: meshtastic.ModuleConfig.SerialConfig
	(*ModuleConfig_ExternalNotificationConfig)(nil), // 11: meshtastic.ModuleConfig.ExternalNotificationConfig
	(*ModuleConfig_StoreForwardConfig)(nil),         // 12: meshtastic.ModuleConfig.StoreForwardConfig
	(*ModuleConfig_RangeTestConfig)(nil),            // 13: meshtastic.ModuleConfig.RangeTestConfig
	(*ModuleConfig_TelemetryConfig)(nil),            // 14: meshtastic.ModuleConfig.TelemetryConfig
	(*ModuleConfig_CannedMessageConfig)(nil),        // 15: meshtastic.ModuleConfig.CannedMessageConfig
	(*ModuleConfig_AudioConfig)(nil),                // 16: meshtastic.ModuleConfig.AudioConfig
	(*ModuleConfig_RemoteHardwareConfig)(nil),       // 17: meshtastic.ModuleConfig.RemoteHardwareConfig
	(*ModuleConfig_NeighborInfoConfig)(nil),         // 18: meshtastic.ModuleConfig.NeighborInfoConfig
	(*ModuleConfig_AmbientLightingConfig)(nil),      // 19: meshtastic.ModuleConfig.AmbientLightingConfig
	(*ModuleConfig_DetectionSensorConfig)(nil),      // 20: meshtastic.ModuleConfig.DetectionSensorConfig
	(*ModuleConfig_PaxcounterConfig)(nil),           // 21: meshtastic.ModuleConfig.PaxcounterConfig
}
var file_meshtastic_localonly_proto_depIdxs = []int32{
	2,  // 0: meshtastic.LocalConfig.device:type_name -> meshtastic.Config.DeviceConfig
	3,  // 1: meshtastic.LocalConfig.position:type_name -> meshtastic.Config.PositionConfig
	4,  // 2: meshtastic.LocalConfig.power:type_name -> meshtastic.Config.PowerConfig
	5,  // 3: meshtastic.LocalConfig.network:type_name -> meshtastic.Config.NetworkConfig
	6,  // 4: meshtastic.LocalConfig.display:type_name -> meshtastic.Config.DisplayConfig
	7,  // 5: meshtastic.LocalConfig.lora:type_name -> meshtastic.Config.LoRaConfig
	8,  // 6: meshtastic.LocalConfig.bluetooth:type_name -> meshtastic.Config.BluetoothConfig
	9,  // 7: meshtastic.LocalModuleConfig.mqtt:type_name -> meshtastic.ModuleConfig.MQTTConfig
	10, // 8: meshtastic.LocalModuleConfig.serial:type_name -> meshtastic.ModuleConfig.SerialConfig
	11, // 9: meshtastic.LocalModuleConfig.external_notification:type_name -> meshtastic.ModuleConfig.ExternalNotificationConfig
	12, // 10: meshtastic.LocalModuleConfig.store_forward:type_name -> meshtastic.ModuleConfig.StoreForwardConfig
	13, // 11: meshtastic.LocalModuleConfig.range_test:type_name -> meshtastic.ModuleConfig.RangeTestConfig
	14, // 12: meshtastic.LocalModuleConfig.telemetry:type_name -> meshtastic.ModuleConfig.TelemetryConfig
	15, // 13: meshtastic.LocalModuleConfig.canned_message:type_name -> meshtastic.ModuleConfig.CannedMessageConfig
	16, // 14: meshtastic.LocalModuleConfig.audio:type_name -> meshtastic.ModuleConfig.AudioConfig
	17, // 15: meshtastic.LocalModuleConfig.remote_hardware:type_name -> meshtastic.ModuleConfig.RemoteHardwareConfig
	18, // 16: meshtastic.LocalModuleConfig.neighbor_info:type_name -> meshtastic.ModuleConfig.NeighborInfoConfig
	19, // 17: meshtastic.LocalModuleConfig.ambient_lighting:type_name -> meshtastic.ModuleConfig.AmbientLightingConfig
	20, // 18: meshtastic.LocalModuleConfig.detection_sensor:type_name -> meshtastic.ModuleConfig.DetectionSensorConfig
	21, // 19: meshtastic.LocalModuleConfig.paxcounter:type_name -> meshtastic.ModuleConfig.PaxcounterConfig
	20, // [20:20] is the sub-list for method output_type
	20, // [20:20] is the sub-list for method input_type
	20, // [20:20] is the sub-list for extension type_name
	20, // [20:20] is the sub-list for extension extendee
	0,  // [0:20] is the sub-list for field type_name
}

func init() { file_meshtastic_localonly_proto_init() }
func file_meshtastic_localonly_proto_init() {
	if File_meshtastic_localonly_proto != nil {
		return
	}
	file_meshtastic_config_proto_init()
	file_meshtastic_module_config_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_meshtastic_localonly_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_meshtastic_localonly_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalModuleConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_meshtastic_localonly_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_meshtastic_localonly_proto_goTypes,
		DependencyIndexes: file_meshtastic_localonly_proto_depIdxs,
		MessageInfos:      file_meshtastic_localonly_proto_msgTypes,
	}.Build()
	File_meshtastic_localonly_proto = out.File
	file_meshtastic_localonly_proto_rawDesc = nil
	file_meshtastic_localonly_proto_goTypes = nil
	file_meshtastic_localonly_proto_depIdxs = nil
}
