package domain

import (
	"github.com/gopcua/opcua/ua"
)

func UATypetoDataType(t ua.TypeID) OpcUaDataType {
	switch t {
	case ua.TypeIDFloat, ua.TypeIDDouble:
		return OpcUaDataTypeFloat
	case ua.TypeIDInt16, ua.TypeIDUint16, ua.TypeIDInt32, ua.TypeIDUint32, ua.TypeIDInt64, ua.TypeIDUint64:
		return OpcUaDataTypeInt
	case ua.TypeIDBoolean:
		return OpcUaDataTypeBoolean
	case ua.TypeIDString:
		return OpcUaDataTypeString
	case ua.TypeIDNull, ua.TypeIDSByte, ua.TypeIDDateTime, ua.TypeIDGUID, ua.TypeIDByteString, ua.TypeIDXMLElement, ua.TypeIDNodeID, ua.TypeIDExpandedNodeID, ua.TypeIDStatusCode, ua.TypeIDQualifiedName, ua.TypeIDLocalizedText, ua.TypeIDExtensionObject, ua.TypeIDDataValue, ua.TypeIDVariant, ua.TypeIDDiagnosticInfo:
		return OpcUaDataTypeOther
	}
	return OpcUaDataTypeUnknown
}
