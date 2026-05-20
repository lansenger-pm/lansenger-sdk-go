package lansenger

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type CallbackEvent struct {
	EventType string
	Category  string
	Data      map[string]interface{}
	RawData   map[string]interface{}
}

func ParseCallbackPayload(queryString string) ([]CallbackEvent, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("parsing callback payload: %w", err)
	}

	events := []CallbackEvent{}

	eventTypes := values["eventType"]
	for _, eventType := range eventTypes {
		category := CallbackEventTypes[eventType]

		data := map[string]interface{}{}
		for key, vals := range values {
			if key != "eventType" && key != "signature" && key != "timestamp" && key != "nonce" {
				if len(vals) == 1 {
					data[key] = vals[0]
				} else {
					data[key] = vals
				}
			}
		}

		events = append(events, CallbackEvent{
			EventType: eventType,
			Category:  category,
			Data:      data,
		})
	}

	return events, nil
}

func VerifyCallbackSignature(queryString, appSecret string) bool {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return false
	}

	signature := values.Get("signature")
	if signature == "" {
		return false
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		if key != "signature" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, key := range keys {
		for _, val := range values[key] {
			buf.WriteString(key)
			buf.WriteString("=")
			buf.WriteString(val)
		}
	}
	buf.WriteString(appSecret)

	hash := sha256.Sum256([]byte(buf.String()))
	computedSig := hex.EncodeToString(hash[:])

	return computedSig == signature
}

func GetCallbackEventTypes() map[string]string {
	return CallbackEventTypes
}
