package http2

import (
	"fmt"
	"testing"
)

var (
	encodedBytes       = []byte{255, 199, 255, 253, 143, 255, 255, 226, 255, 255, 254, 63, 255, 255, 228, 255, 255, 254, 95, 255, 255, 230, 255, 255, 254, 127, 255, 255, 232, 255, 255, 234, 255, 255, 255, 243, 255, 255, 250, 127, 255, 255, 171, 255, 255, 255, 223, 255, 255, 235, 255, 255, 254, 207, 255, 255, 237, 255, 255, 254, 239, 255, 255, 239, 255, 255, 255, 15, 255, 255, 241, 255, 255, 255, 47, 255, 255, 255, 191, 255, 255, 207, 255, 255, 253, 63, 255, 255, 215, 255, 255, 253, 191, 255, 255, 223, 255, 255, 254, 63, 255, 255, 231, 255, 255, 254, 191, 255, 255, 237, 79, 227, 249, 255, 175, 252, 171, 241, 254, 191, 175, 239, 231, 253, 253, 44, 187, 0, 8, 153, 105, 183, 29, 121, 251, 159, 127, 255, 32, 255, 191, 243, 255, 80, 221, 189, 127, 6, 28, 88, 242, 101, 205, 159, 70, 157, 90, 246, 109, 221, 191, 135, 30, 95, 156, 255, 127, 247, 255, 252, 63, 249, 255, 228, 95, 255, 71, 25, 36, 44, 179, 78, 110, 157, 104, 166, 163, 215, 218, 196, 38, 222, 254, 60, 250, 247, 255, 251, 254, 127, 251, 255, 223, 255, 255, 252, 255, 254, 111, 255, 244, 191, 255, 159, 255, 250, 63, 255, 211, 255, 255, 83, 255, 253, 95, 255, 251, 63, 255, 235, 127, 255, 218, 255, 255, 183, 255, 255, 115, 255, 254, 239, 255, 253, 239, 255, 254, 191, 255, 251, 255, 255, 253, 159, 255, 253, 191, 255, 235, 255, 255, 224, 255, 255, 238, 255, 255, 195, 255, 255, 139, 255, 255, 31, 255, 254, 79, 255, 238, 127, 255, 177, 255, 255, 151, 255, 253, 159, 255, 252, 223, 255, 249, 255, 255, 251, 255, 255, 218, 255, 254, 239, 255, 244, 255, 255, 183, 255, 254, 231, 255, 254, 143, 255, 253, 63, 255, 222, 255, 255, 213, 255, 254, 239, 255, 251, 223, 255, 254, 31, 255, 223, 255, 255, 127, 255, 255, 95, 255, 254, 207, 255, 240, 127, 255, 135, 255, 254, 15, 255, 241, 127, 255, 237, 255, 255, 135, 255, 255, 119, 255, 254, 255, 255, 234, 255, 255, 139, 255, 254, 63, 255, 249, 63, 255, 248, 127, 255, 203, 255, 255, 55, 255, 255, 31, 255, 255, 131, 255, 255, 225, 255, 254, 191, 255, 227, 255, 255, 63, 255, 255, 47, 255, 250, 63, 255, 253, 159, 255, 255, 23, 255, 255, 199, 255, 255, 242, 127, 255, 253, 239, 255, 255, 191, 255, 255, 242, 255, 255, 248, 255, 255, 251, 127, 255, 151, 255, 248, 255, 255, 254, 111, 255, 255, 193, 255, 255, 248, 127, 255, 254, 127, 255, 255, 197, 255, 255, 229, 255, 254, 79, 255, 242, 255, 255, 253, 31, 255, 255, 79, 255, 255, 254, 255, 255, 254, 63, 255, 255, 201, 255, 255, 249, 127, 255, 179, 255, 255, 207, 255, 251, 127, 255, 205, 255, 255, 79, 255, 249, 255, 255, 209, 255, 255, 207, 255, 254, 175, 255, 250, 255, 255, 253, 223, 255, 254, 255, 255, 255, 79, 255, 255, 95, 255, 255, 171, 255, 255, 167, 255, 255, 215, 255, 255, 249, 191, 255, 254, 207, 255, 255, 183, 255, 255, 243, 255, 255, 254, 143, 255, 255, 211, 255, 255, 250, 191, 255, 255, 95, 255, 255, 255, 127, 255, 254, 207, 255, 255, 219, 255, 255, 251, 191, 255, 255, 127, 255, 255, 240, 255, 255, 251, 191}
	littleEncodedBytes = []byte{83, 248, 254, 127, 235, 255, 42, 252, 127, 175, 235, 251, 249, 255, 127, 75, 46, 192, 2, 38, 90, 109, 199, 94, 126, 231, 220, 55, 111, 95, 193, 135, 22, 60, 153, 115, 103, 209, 167, 86, 189, 155, 119, 111, 225, 199, 151, 231, 63, 209, 198, 73, 11, 44, 211, 155, 167, 90, 41, 168, 245, 246, 177, 9, 183, 191, 143, 62, 189, 255}

	decodedBytes       = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255}
	littleDecodedBytes = []byte{32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122}
)

func compareBytes(b, bb []byte) error {
	if len(b) != len(bb) {
		return fmt.Errorf("different sizes: %d<>%d\n", len(b), len(bb))
	}
	for i := 0; i < len(b); i++ {
		if b[i] != bb[i] {
			return fmt.Errorf("different bytes at %d: %v<>%v", i, b[i], bb[i])
		}
	}
	return nil
}

func makeCopy(bb []byte) []byte {
	var b = make([]byte, len(bb))
	copy(b, bb)
	return b
}

func decodeHuffman(t *testing.T, b, bb, toCompare []byte) {
	b = HuffmanDecode(b[:0], bb)
	if err := compareBytes(b, toCompare); err != nil {
		t.Fatal(err)
	}
}

func TestHuffmanDecodeHuge(t *testing.T) {
	decodeHuffman(t, nil, encodedBytes, decodedBytes)
}

func TestHuffmanDecode(t *testing.T) {
	decodeHuffman(t, nil, littleEncodedBytes, littleDecodedBytes)
}

/* This algorithm cannot reuse bytes
func TestHuffmanDecodeReusing(t *testing.T) {
	b := makeCopy(encodedBytes)
	decodeHuffman(t, b, b, decodedBytes)
}

func TestHuffmanDecodeReusingLittle(t *testing.T) {
	b := makeCopy(littleEncodedBytes)
	decodeHuffman(t, b, b, littleDecodedBytes)
}
*/

func encodeHuffman(t *testing.T, b, bb, toCompare []byte) {
	b = HuffmanEncode(b[:0], bb)
	if err := compareBytes(b, toCompare); err != nil {
		t.Fatal(err)
	}
}

func TestHuffmanEncodeHuge(t *testing.T) {
	encodeHuffman(t, nil, decodedBytes, encodedBytes)
}

func TestHuffmanEncode(t *testing.T) {
	encodeHuffman(t, nil, littleDecodedBytes, littleEncodedBytes)
}

/* This algorithm cannot reuse bytes
func TestHuffmanEncodeReusing(t *testing.T) {
	b := makeCopy(decodedBytes)
	encodeHuffman(t, b, b, encodedBytes)
}

func TestHuffmanEncodeReusingLittle(t *testing.T) {
	b := makeCopy(littleDecodedBytes)
	encodeHuffman(t, b, b, littleEncodedBytes)
}
*/