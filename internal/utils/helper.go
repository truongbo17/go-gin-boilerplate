package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

func FormatStackTrace(stack []byte) string {
	lines := strings.Split(string(stack), "\n")
	return strings.Join(lines, "\n")
}

func GenerateRandomString(isNumber int, length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if isNumber == 1 {
		characters = "0123456789"
	}

	charactersLength := len(characters)
	randomString := make([]byte, length)

	for i := 0; i < length; i++ {
		randomString[i] = characters[r.Intn(charactersLength)]
	}

	return string(randomString)
}

func VNStrFilter(str string) string {
	unicodeMap := map[string]string{
		"a": "á|à|ả|ã|ạ|ă|ắ|ặ|ằ|ẳ|ẵ|â|ấ|ầ|ẩ|ẫ|ậ",
		"d": "đ",
		"e": "é|è|ẻ|ẽ|ẹ|ê|ế|ề|ể|ễ|ệ",
		"i": "í|ì|ỉ|ĩ|ị",
		"o": "ó|ò|ỏ|õ|ọ|ô|ố|ồ|ổ|ỗ|ộ|ơ|ớ|ờ|ở|ỡ|ợ",
		"u": "ú|ù|ủ|ũ|ụ|ư|ứ|ừ|ử|ữ|ự",
		"y": "ý|ỳ|ỷ|ỹ|ỵ",
		"A": "Á|À|Ả|Ã|Ạ|Ă|Ắ|Ặ|Ằ|Ẳ|Ẵ|Â|Ấ|Ầ|Ẩ|Ẫ|Ậ",
		"D": "Đ",
		"E": "É|È|Ẻ|Ẽ|Ẹ|Ê|Ế|Ề|Ể|Ễ|Ệ",
		"I": "Í|Ì|Ỉ|Ĩ|Ị",
		"O": "Ó|Ò|Ỏ|Õ|Ọ|Ô|Ố|Ồ|Ổ|Ỗ|Ộ|Ơ|Ớ|Ờ|Ở|Ỡ|Ợ",
		"U": "Ú|Ù|Ủ|Ũ|Ụ|Ư|Ứ|Ừ|Ử|Ữ|Ự",
		"Y": "Ý|Ỳ|Ỷ|Ỹ|Ỵ",
	}

	for nonUnicode, uni := range unicodeMap {
		re := regexp.MustCompile(uni)
		str = re.ReplaceAllString(str, nonUnicode)
	}

	return str
}

func SortMapKeys(data map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedMap := make(map[string]interface{}, len(data))
	for _, k := range keys {
		if nestedMap, ok := data[k].(map[string]interface{}); ok {
			sortedMap[k] = SortMapKeys(nestedMap)
		} else if nestedSlice, ok := data[k].([]interface{}); ok {
			newSlice := make([]interface{}, len(nestedSlice))
			for i, v := range nestedSlice {
				if nestedItem, ok := v.(map[string]interface{}); ok {
					newSlice[i] = SortMapKeys(nestedItem)
				} else {
					newSlice[i] = v
				}
			}
			sortedMap[k] = newSlice
		} else {
			sortedMap[k] = data[k]
		}
	}
	return sortedMap
}

func CheckLanguage(target string) bool {
	for _, code := range config.LanguageAvailable {
		if code == target {
			return true
		}
	}
	return false
}

func GetMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
