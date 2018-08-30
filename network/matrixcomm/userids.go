package matrixcomm

import (
	"bytes"
	"fmt"
	"encoding/hex"
	"strings"
)

const lowerhex = "0123456789abcdef"

/*type uidRegexp struct {
	*regexp.Regexp
}

var uidExp = uidRegexp{regexp.MustCompile(`(?P<addr_local>\d+)\.(?P<addr_random>\d+)`)}

func (r *uidRegexp) FindStringSubmatchMap(s string) map[string]string {
	_match:=uidExp.MatchString(s)
	if _match!=true{
		return nil
	}
	captures := make(map[string]string)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}
	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}*/

//从userID中提取Localpart,e.g:"@myname:smartraiden.org"->"myname"
func ExtractUserLocalpart(userID string) (string, error) {
	if len(userID) == 0 || userID[0] != '@' {
		return "", fmt.Errorf("%s is not a valid user id", userID)
	}
	return strings.TrimPrefix(
		strings.SplitN(userID, ":", 2)[0],
		"@", // remove "@" prefix
	), nil
}


func encode(buf *bytes.Buffer, b byte) {
	buf.WriteByte('=')
	buf.WriteByte(lowerhex[b>>4])
	buf.WriteByte(lowerhex[b&0x0f])
}

func escape(buf *bytes.Buffer, b byte) {
	buf.WriteByte('_')
	if b == '_' {
		buf.WriteByte('_')
	} else {
		buf.WriteByte(b + 0x20)
	}
}

func shouldEncode(b byte) bool {
	return b != '-' && b != '.' && b != '_' && !(b >= '0' && b <= '9') && !(b >= 'a' && b <= 'z') && !(b >= 'A' && b <= 'Z')
}

func shouldEscape(b byte) bool {
	return (b >= 'A' && b <= 'Z') || b == '_'
}

func isValidByte(b byte) bool {
	return isValidEscapedChar(b) || (b >= '0' && b <= '9') || b == '.' || b == '=' || b == '-'
}

func isValidEscapedChar(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z')
}

//将给定的字符串编码为符合matrix的userID格式
func EncodeUserLocalpart(str string) string {
	strBytes := []byte(str)
	var outputBuffer bytes.Buffer
	for _, b := range strBytes {
		if shouldEncode(b) {
			encode(&outputBuffer, b)
		} else if shouldEscape(b) {
			escape(&outputBuffer, b)
		} else {
			outputBuffer.WriteByte(b)
		}
	}
	return outputBuffer.String()
}

//将给定字符串解码回原始输入字符串。_alph=40_bet=5f50up  =>  Alph@Bet_50up
func DecodeUserLocalpart(str string) (string, error) {
	strBytes := []byte(str)
	var outputBuffer bytes.Buffer
	for i := 0; i < len(strBytes); i++ {
		b := strBytes[i]
		if !isValidByte(b) {
			return "", fmt.Errorf("Byte pos %d: Invalid byte", i)
		}

		if b == '_' {
			if i+1 >= len(strBytes) {
				return "", fmt.Errorf("Byte pos %d: expected _[a-z_] encoding but ran out of string", i)
			}
			if !isValidEscapedChar(strBytes[i+1]) {
				return "", fmt.Errorf("Byte pos %d: expected _[a-z_] encoding", i)
			}
			if strBytes[i+1] == '_' {
				outputBuffer.WriteByte('_')
			} else {
				outputBuffer.WriteByte(strBytes[i+1] - 0x20)
			}
			i++
		} else if b == '=' {
			if i+2 >= len(strBytes) {
				return "", fmt.Errorf("Byte pos: %d: expected quote-printable encoding but ran out of string", i)
			}
			dst := make([]byte, 1)
			_, err := hex.Decode(dst, strBytes[i+1:i+3])
			if err != nil {
				return "", err
			}
			outputBuffer.WriteByte(dst[0])
			i += 2
		} else {
			outputBuffer.WriteByte(b)
		}
	}
	return outputBuffer.String(), nil
}

