package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger" // new
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	type Sample struct {
		Key  string `json:"key"`
		Text string `json:"text"`
	}

	type SampleEncrypted struct {
		Key string `json:"key"`
		Enc string `json:"enc"`
	}

	// give response when at /
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the endpoint.",
		})
	})

	// give response when at /
	app.Post("*encrypt", func(c *fiber.Ctx) error {

		p := new(Sample)

		if err := c.BodyParser(p); err != nil {
			return err

		}

		if p.Key == "" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Property 'key' or 'text' does not exist in body.",
				"success": false,
			})
		}

		// encrypting the text
		foo1, err := encryptString(p.Text, p.Key)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"enc": foo1,
		})
	})

	// give response when at /
	app.Post("*decrypt", func(c *fiber.Ctx) error {

		p := new(SampleEncrypted)

		if err := c.BodyParser(p); err != nil {
			return err

		}

		if p.Key == "" || p.Enc == "" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Property 'key' or 'enc' does not exist in body.",
				"success": false,
			})
		}

		// encrypting the text
		foo, err := decryptString(p.Enc, p.Key)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"text": string(foo),
		})
	})

	// Listen on server 8000 and catch error if any
	err := app.Listen(":8000")

	// handle error
	if err != nil {
		panic(err)
	}
}

func encryptString(plainText string, keyString string) (cipherTextString string, err error) {

	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES([]byte(key), []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func encryptAES(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func decryptString(cryptoText string, keyString string) (plainTextString string, err error) {

	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES([]byte(hashTo32Bytes(keyString)), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}
func hashTo32Bytes(input string) (output string) {

	if len(input) == 0 {
		return ""
	}

	hasher := sha256.New()
	hasher.Write([]byte(input))

	stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Cut the length down to 32 bytes and return.
	return stringToSHA256[:32]
}
