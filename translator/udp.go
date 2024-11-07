package translator

import (
	"sia/backend/lib"
	"sia/backend/types"
)

/*
Rules for translating UDP server:
1. The first byte, is the data type
- 0: ecg-sensor
- 1: temperature-sensor
- 2: gps-service
2. After the second byte, is the data
*/
func TranslateUDPBinary(data []byte) (int, interface{}) {
	dataType := int(data[0])
	switch dataType {
	case types.UDP_EKG_SENSOR:
		data := lib.Float64FromBytes(data[1:])
		result := new(types.EKG_SENSOR).Init(data)
		return dataType, result

	case types.UDP_TEMPERATURE_SENSOR:
		data := lib.Float64FromBytes(data[1:])
		result := new(types.TEMPERATURE_SENSOR).Init(data)
		return dataType, result

	case types.UDP_GPS_SERVICE:
		latitude := lib.Float64FromBytes(data[1:])
		longitude := lib.Float64FromBytes(data[9:])
		result := new(types.GPS_SERVICE).Init(latitude, longitude)
		return dataType, result

	default:
		return -1, nil
	}
}
