package entity

type EncryptedField = []byte

type Encryptor struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement"`
	Index int    `gorm:"primary_key auto_increment"`
	Hash  []byte `gorm:"primary_key type:varbinary(512)"`
}
