package toolbox

import "github.com/google/uuid"

// crockfordAlphabet is the Base32 alphabet defined by the TypeID spec.
const crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// TypeIDPrefix is the string prefix that identifies the entity type in a TypeID.
type TypeIDPrefix string

const (
	OneTimeTokenPrefix TypeIDPrefix = "ott"
)

// GenerateTypeID returns a TypeID-compatible string of the form "{prefix}_{suffix}",
// where the suffix is a UUIDv7 encoded as 26 Crockford Base32 characters.
// The output is byte-for-byte identical to the go.jetify.com/typeid format it replaces.
func GenerateTypeID(prefix TypeIDPrefix) string {
	id, err := uuid.NewV7()
	if err != nil {
		return string(prefix) + "_" + "00000000000000000000000000"
	}
	return string(prefix) + "_" + encodeBase32(id)
}

// encodeBase32 packs a 128-bit UUID into 26 Crockford Base32 characters.
// 26 × 5 bits = 130 bits; the two most-significant padding bits are always 0,
// which is guaranteed for any UUIDv7 value (top byte never exceeds 0x1F after
// the version/variant fields are applied).
func encodeBase32(id uuid.UUID) string {
	b := [16]byte(id)
	var dst [26]byte
	dst[0] = crockfordAlphabet[(b[0]&0xE0)>>5]
	dst[1] = crockfordAlphabet[b[0]&0x1F]
	dst[2] = crockfordAlphabet[(b[1]&0xF8)>>3]
	dst[3] = crockfordAlphabet[((b[1]&0x07)<<2)|((b[2]&0xC0)>>6)]
	dst[4] = crockfordAlphabet[(b[2]&0x3E)>>1]
	dst[5] = crockfordAlphabet[((b[2]&0x01)<<4)|((b[3]&0xF0)>>4)]
	dst[6] = crockfordAlphabet[((b[3]&0x0F)<<1)|((b[4]&0x80)>>7)]
	dst[7] = crockfordAlphabet[(b[4]&0x7C)>>2]
	dst[8] = crockfordAlphabet[((b[4]&0x03)<<3)|((b[5]&0xE0)>>5)]
	dst[9] = crockfordAlphabet[b[5]&0x1F]
	dst[10] = crockfordAlphabet[(b[6]&0xF8)>>3]
	dst[11] = crockfordAlphabet[((b[6]&0x07)<<2)|((b[7]&0xC0)>>6)]
	dst[12] = crockfordAlphabet[(b[7]&0x3E)>>1]
	dst[13] = crockfordAlphabet[((b[7]&0x01)<<4)|((b[8]&0xF0)>>4)]
	dst[14] = crockfordAlphabet[((b[8]&0x0F)<<1)|((b[9]&0x80)>>7)]
	dst[15] = crockfordAlphabet[(b[9]&0x7C)>>2]
	dst[16] = crockfordAlphabet[((b[9]&0x03)<<3)|((b[10]&0xE0)>>5)]
	dst[17] = crockfordAlphabet[b[10]&0x1F]
	dst[18] = crockfordAlphabet[(b[11]&0xF8)>>3]
	dst[19] = crockfordAlphabet[((b[11]&0x07)<<2)|((b[12]&0xC0)>>6)]
	dst[20] = crockfordAlphabet[(b[12]&0x3E)>>1]
	dst[21] = crockfordAlphabet[((b[12]&0x01)<<4)|((b[13]&0xF0)>>4)]
	dst[22] = crockfordAlphabet[((b[13]&0x0F)<<1)|((b[14]&0x80)>>7)]
	dst[23] = crockfordAlphabet[(b[14]&0x7C)>>2]
	dst[24] = crockfordAlphabet[((b[14]&0x03)<<3)|((b[15]&0xE0)>>5)]
	dst[25] = crockfordAlphabet[b[15]&0x1F]
	return string(dst[:])
}
