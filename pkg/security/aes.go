package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func AesEncrypt(data string, key string, mode string) string {
	var encrypt []byte
	switch mode {
	case "CBC":
		encrypt = aesCBCEncrypt([]byte(data), []byte(key))
	}
	return base64.StdEncoding.EncodeToString(encrypt)
}

func AesDecrypt(encrypt string, key string, mode string) string {
	decrypt, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		panic(err)
	}
	switch mode {
	case "CBC":
		decrypt = aesCBCDecrypt([]byte(decrypt), []byte(key))
	}

	return string(decrypt)
}

//对明文进行填充
func padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//对密文删除填充
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	if len(cipherText)-int(end) < 0 {
		return []byte("")
	}
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

func aesCBCEncrypt(plainText []byte, key []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//获取block长度
	blockSize := block.BlockSize()
	//进行填充
	plainText = padding(plainText, blockSize)
	//指定初始向量, 长度和block的块尺寸一致
	cipherText := make([]byte, blockSize+len(plainText))
	iv := cipherText[:blockSize]
	//ReadFull从rand.Reader精确地读取len(b)字节数据填充进b
	//rand.Reader是一个全局、共享的密码用强随机数生成器
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText[blockSize:], plainText)
	//返回密文
	return cipherText
}

func aesCBCDecrypt(cipherText []byte, key []byte) []byte {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		panic("ciphertext too short")
	}

	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	// CBC mode always works in whole blocks.
	if len(cipherText)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(cipherText, cipherText)
	//删除填充
	cipherText = unPadding(cipherText)
	return cipherText
}
