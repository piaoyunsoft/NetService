package NS

type ISecurity interface {
	Encrypt(in []byte) (out []byte, err error)
	Decrypt(in []byte) (out []byte, err error)
}
