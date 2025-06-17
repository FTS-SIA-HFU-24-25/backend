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
	lib.Print(lib.UDP_SERVICE, "Data type: %v", dataType)
	switch dataType {
	case types.UDP_EKG_SENSOR:
		data := lib.Float64FromBytes(data[1:])
		result := new(types.EKG_SENSOR).Init(data)
		return dataType, result

	case types.UDP_TEMPERATURE_SENSOR:
		data := lib.Float32FromBytes(data[1:5])
		result := new(types.TEMPERATURE_SENSOR).Init(data)
		return dataType, result

	case types.UDP_GYRO_SENSOR:
		x := lib.Float32FromBytes(data[1:5])
		y := lib.Float32FromBytes(data[5:9])
		z := lib.Float32FromBytes(data[9:13])
		result := new(types.GYRO_SENSOR).Init(x, y, z)
		return dataType, result

	case types.UDP_ACCEL_SENSOR:
		x := lib.Float32FromBytes(data[1:5])
		y := lib.Float32FromBytes(data[5:9])
		z := lib.Float32FromBytes(data[9:13])
		result := new(types.ACCEL_SENSOR).Init(x, y, z)
		return dataType, result

	case types.END_CONNECTION:
		result := types.END_REQUEST{
			SENSOR_ID: int(data[1]),
		}
		return dataType, &result
		
	default:
		lib.Print(lib.UDP_SERVICE, "Unkown data")
		return -1, nil
	}
}
