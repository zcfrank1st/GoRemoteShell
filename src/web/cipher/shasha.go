package cipher

import (
    "hash"
    "crypto/sha1"
    "encoding/hex"
)

var shasha hash.Hash

func init () {
    shasha = sha1.New()
}

func Encode(url string, timestamp string, id string) string {
    shasha.Write([]byte(url + timestamp + id))
    bs := shasha.Sum(nil)
    return hex.EncodeToString(bs)
}