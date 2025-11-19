package protocol

import "testing"

func Test_EncryptAndDecrypt(t *testing.T) {
	inputData := []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")

	encryptedData, err := EncryptRSA(inputData)
	if err != nil {
		t.Fatal(err)
	}
	decryptedData, err := DecryptRSA(encryptedData)
	if err != nil {
		t.Fatal(err)
	}
	if string(decryptedData) != string(inputData) {
		t.Fatalf("Decrypted data does not match original. Got '%s', expected '%s'", string(decryptedData), string(inputData))
	}
}
