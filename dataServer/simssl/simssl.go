package simssl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

/*
SimSsl packet design of Simple ssl
*/
type SimSsl struct {
	/*Content Type
	0x01 client Hello
	0x02 Server Hello
	0x03 client knows server failed
	0x04 server echo eraser
	*/
	ContentType uint8 //Required
	/* Version
	0x01
	*/
	Version uint8 //Required
	/*
		The total length of all simssl packets (bytes)
	*/
	Length uint16 //Required
	/*
		Method 0x01 AES128
	*/
	Method uint8 //required
	/*
		Mode 0x01 ECB
		Mode 0x02 CBC
		Mode 0x03 CTR
		Mode 0x04 CFB
		Mode oxo5 OFB
	*/
	Mode uint8 //required
	/*
		Checksum
	*/
	CheckSum uint16 //required
	/*
		SHA256 hash of client identification
	*/
	ClientID [32]uint8 //Required
	/*
		SHA256 hash of server identification
		zerospadding when client sends
	*/
	ServerID [32]uint8 //required
	/*
		Timestamp + (60*60*24)
		add one day when key exchange
	*/
	ExpirationTime int64 //required

	/*
		The initial string for check the validation of EncrptKey
	*/
	RandomInit [32]uint8 //optional
	/*
		Encrypt key
	*/
	EncryptKey [16]uint8 //optianl
}

/*
GenerateClientHello Generate a client Hello Packet
*/
func GenerateClientHello(cid []byte) (SimSsl, error) {
	rand.Seed(time.Now().Unix())
	clientHello := SimSsl{
		ContentType:    0x01,
		Version:        0x01,
		Length:         128, //字节
		Method:         0x01,
		Mode:           0x02,
		CheckSum:       0x00,
		ClientID:       sha256.Sum256(cid),
		ServerID:       sha256.Sum256([]byte("unknown")),
		ExpirationTime: time.Now().Unix() + 60*60*24,
		RandomInit:     sha256.Sum256([]byte(strconv.FormatUint(rand.Uint64(), 36))),
		EncryptKey:     md5.Sum([]byte(strconv.FormatUint(rand.Uint64(), 36))),
	}
	encryptMessage, err := AesEncrypt(clientHello.RandomInit[:], clientHello.EncryptKey[:])
	if err != nil {
		return SimSsl{}, err
	}
	clientHello.RandomInit = encryptMessage
	//	fmt.Println(unsafe.Sizeof(clientHello.Mode))
	//	fmt.Println(unsafe.Sizeof(clientHello))
	fmt.Printf("%x\n", clientHello.RandomInit)

	return clientHello, nil
}

/*
AesEncrypt AES decrypt with message and key
*/
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	//origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

/*
AesDecrypt AES decrypt with crypted message and key
*/
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	// origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
