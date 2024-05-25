package testutils

import (
	"fmt"
	"os"

	"github.com/hokkung/go-tumboon/pkg/cipher"
)

func GenerateMockDonationFile(filePath string, payload string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	rotWriter, err := cipher.NewRot128Writer(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := []byte(payload)
	_, err = rotWriter.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func RemoveMockFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}
}
