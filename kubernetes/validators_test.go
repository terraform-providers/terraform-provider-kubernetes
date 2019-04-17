package kubernetes

import (
	"testing"
)

func TestValidateModeBits(t *testing.T) {
	validCases := []string{
		"0", "0001", "0644", "0777",
	}
	for _, mode := range validCases {
		_, es := validateModeBits(mode, "mode")
		if len(es) > 0 {
			t.Fatalf("Expected %s to be valid: %#v", mode, es)
		}
	}

	invalidCases := []string{
		"-5", "-1", "512", "777",
	}
	for _, mode := range invalidCases {
		_, es := validateModeBits(mode, "mode")
		if len(es) == 0 {
			t.Fatalf("Expected %s to be invalid", mode)
		}
	}
}

func TestValidateBase64Encoded(t *testing.T) {
	validCases := []interface{}{
		"",         // the plain empty string
		"Cg==",     // the encoded empty string
		"blah",     // "nV"
		"VGVzdAo=", // "Test"
		"f0VMRgIBAQAAAAAAAAAAAAMAPgABAAAAMGEAAAAAAABAAAAAAAAAAKATAgAAAAAAAAAAAEAAOAALAEAAHAAbAAYAAAAEAAAAQAAAAAAAAABAAAAAAAAAAEAAAAAAAAAAaAIAAA==", // `head /bin/ls -c 100 | base64`
	}
	for _, data := range validCases {
		_, es := validateBase64Encoded(data, "binary_data")
		if len(es) > 0 {
			t.Fatalf("Expected %#o to be valid: %#v", data, es)
		}
	}

	invalidCases := []interface{}{
		nil,
		"bl ah",
		"blahd",
		"C=",
	}
	for _, data := range invalidCases {
		_, es := validateBase64Encoded(data, "binary_data")
		if len(es) == 0 {
			t.Fatalf("Expected %#o to be invalid", data)
		}
	}
}
