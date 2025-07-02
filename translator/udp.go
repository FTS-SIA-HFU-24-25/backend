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
	if len(data) < 1 {
		lib.Print(lib.UDP_SERVICE, "Invalid data length")
		return -1, nil
	}

	dataType := int(data[0])
	lib.Print(lib.UDP_SERVICE, "Data type: %v", dataType)

	switch dataType {
	case types.UDP_EKG_SENSOR:
		if len(data) < 3 {
			lib.Print(lib.UDP_SERVICE, "Invalid ECG data length")
			return -1, nil
		}

		msb := int(data[1]) & 0x03 // Only lower 2 bits are valid for MSB
		lsb := int(data[2])        // Full 8 bits

		value := (msb << 8) | lsb // Combine MSB and LSB into a 10-bit integer (0–1023)

		result := new(types.EKG_SENSOR).Init(float64(value))
		return dataType, result

	case types.UDP_TEMPERATURE_SENSOR:
		if len(data) < 5 {
			lib.Print(lib.UDP_SERVICE, "Invalid Temperature data length")
			return -1, nil
		}
		temp := lib.Float32FromBytes(data[1:5])
		result := new(types.TEMPERATURE_SENSOR).Init(temp)
		return dataType, result

	case types.UDP_GYRO_SENSOR:
		if len(data) < 13 {
			lib.Print(lib.UDP_SERVICE, "Invalid Gyro data length")
			return -1, nil
		}
		x := lib.Float32FromBytes(data[1:5])
		y := lib.Float32FromBytes(data[5:9])
		z := lib.Float32FromBytes(data[9:13])
		result := new(types.GYRO_SENSOR).Init(x, y, z)
		return dataType, result

	case types.UDP_ACCEL_SENSOR:
		if len(data) < 13 {
			lib.Print(lib.UDP_SERVICE, "Invalid Accel data length")
			return -1, nil
		}
		x := lib.Float32FromBytes(data[1:5])
		y := lib.Float32FromBytes(data[5:9])
		z := lib.Float32FromBytes(data[9:13])
		result := new(types.ACCEL_SENSOR).Init(x, y, z)
		return dataType, result

	case types.END_CONNECTION:
		if len(data) < 2 {
			lib.Print(lib.UDP_SERVICE, "Invalid END_CONNECTION data length")
			return -1, nil
		}
		result := types.END_REQUEST{
			SENSOR_ID: int(data[1]),
		}
		return dataType, &result

	default:
		lib.Print(lib.UDP_SERVICE, "Unknown data type")
		return -1, nil
	}
}
