package domain

import "testing"

func TestCorrectMessage(t *testing.T) {
	rawMessage := []byte("*HQ,1400046168,V1,055600,A,2234.3066,N,11351.6829,E,000.0,000,080813,FFFFFBFF#")

	var gpsMessage GPSMessage
	if err := gpsMessage.UnmarshalText(rawMessage); err != nil {
		t.Error("Should not fail with a valid message")
	}
}

func TestInvalidMessage(t *testing.T) {
	rawMessage := []byte("this message does't even makes sense at all")

	var gpsMessage GPSMessage
	if err := gpsMessage.UnmarshalText(rawMessage); err == nil {
		t.Error("Should fail with a invalid message")
	}
}

func TestLongitudeAndLatitude(t *testing.T) {
	rawMessage := []byte("*HQ,1400046168,V1,055600,A,2234.3066,N,11351.6829,E,000.0,000,080813,FFFFFBFF#")

	var gpsMessage GPSMessage
	if err := gpsMessage.UnmarshalText(rawMessage); err != nil {
		t.Error("Should not fail with a valid message")
	}

	if gpsMessage.Loc.Type != "Point" {
		t.Error("Should have a location of 'Point'")
	}

	expectedLongitude := 113 + (51.6829 / 60)
	if gpsMessage.Loc.Coordinates[0] != expectedLongitude {
		t.Errorf("Unexpected longitude: wanted %f, have %f", expectedLongitude, gpsMessage.Loc.Coordinates[0])
	}

	expectedLatitude := 22 + (34.3066 / 60)
	if gpsMessage.Loc.Coordinates[1] != expectedLatitude {
		t.Errorf("Unexpected latitude: wanted %f, have %f", expectedLatitude, gpsMessage.Loc.Coordinates[1])
	}
}

func TestGPSMessageIsInvalid(t *testing.T) {
	rawMessage := []byte("*HQ,1400046168,V1,055600,V,2234.3066,N,11351.6829,E,000.0,000,080813,FFFFFFFE#")

	var gpsMessage GPSMessage
	if err := gpsMessage.UnmarshalText(rawMessage); err != nil {
		t.Error("Should not fail with a valid message", err)
	}

	if gpsMessage.Valid {
		t.Error("Message should be invalid")
	}
}
