package simssl

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
		The total length of all simssl packets
	*/
	Length uint16 //Required
	/*
		SHA256 hash of client identification
	*/
	ClientID [32]uint8 //Required
	/*
		SHA256 hash of server identification
	*/
	ServerID [32]uint8 //optional
	/*
		The initial string for check the validation of EncrptKey
	*/
	RandomInit [32]uint8 //optional
	/*
		Method 0x01 AES128
	*/
	Method uint8
	Mode   uint8
	/*Timestamp + 60*60*24*/
	ExpirationTime int64 //optional
	EncryptKey     [16]uint8
	CheckSum       uint16
}
