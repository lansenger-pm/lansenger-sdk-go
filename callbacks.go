package lansenger

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type CallbackEvent struct {
	EventID   string
	EventType string
	Category  string
	Data      map[string]interface{}
	RawData   map[string]interface{}
	AppID     string
	OrgID     string
}

type DecryptedCallbackPayload struct {
	Random string
	OrgID  string
	AppID  string
	Events []CallbackEvent
	Length uint32
}

func ParseCallbackPayload(input string) ([]CallbackEvent, error) {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "{") {
		var wrapper map[string]interface{}
		if err := json.Unmarshal([]byte(input), &wrapper); err != nil {
			return nil, fmt.Errorf("parsing callback JSON payload: %w", err)
		}

		if dataEncrypt, ok := wrapper["dataEncrypt"].(string); ok && dataEncrypt != "" {
			return nil, fmt.Errorf("encrypted callback payload requires encoding_key — use DecryptCallbackPayload instead")
		}

		return parsePlainJSONPayload(wrapper)
	}

	values, err := url.ParseQuery(input)
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
			RawData:   data,
		})
	}

	return events, nil
}

func parsePlainJSONPayload(wrapper map[string]interface{}) ([]CallbackEvent, error) {
	eventList, ok := wrapper["events"].([]interface{})
	if !ok {
		if wrapper["eventType"] != nil || wrapper["type"] != nil {
			eventList = []interface{}{wrapper}
		} else {
			return nil, fmt.Errorf("no events found in callback payload")
		}
	}

	topAppID, _ := wrapper["appId"].(string)
	topOrgID, _ := wrapper["orgId"].(string)

	events := []CallbackEvent{}
	for _, entry := range eventList {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		eventType := ""
		if t, ok := entryMap["eventType"].(string); ok {
			eventType = t
		} else if t, ok := entryMap["type"].(string); ok {
			eventType = t
		}

		category := CallbackEventTypes[eventType]
		rawData, _ := entryMap["data"].(map[string]interface{})
		if rawData == nil && eventType != "" {
			rawData = entryMap
		}

		parsedData := parseEventData(eventType, rawData)
		eventID := ""
		if id, ok := entryMap["eventId"].(string); ok {
			eventID = id
		} else if id, ok := entryMap["id"].(string); ok {
			eventID = id
		} else if id, ok := entryMap["eventId"].(float64); ok {
			eventID = fmt.Sprintf("%v", id)
		} else if id, ok := entryMap["id"].(float64); ok {
			eventID = fmt.Sprintf("%v", id)
		}

		appID, _ := entryMap["appId"].(string)
		orgID, _ := entryMap["orgId"].(string)
		if appID == "" {
			appID = topAppID
		}
		if orgID == "" {
			orgID = topOrgID
		}

		events = append(events, CallbackEvent{
			EventID:   eventID,
			EventType: eventType,
			Category:  category,
			Data:      parsedData,
			RawData:   rawData,
			AppID:     appID,
			OrgID:     orgID,
		})
	}

	return events, nil
}

func DecryptCallbackPayload(encryptedData, encodingKey string, knownAppID string) (*DecryptedCallbackPayload, error) {
	encryptedData = strings.TrimSpace(encryptedData)

	if strings.HasPrefix(encryptedData, "{") {
		var wrapper map[string]interface{}
		if err := json.Unmarshal([]byte(encryptedData), &wrapper); err != nil {
			return nil, fmt.Errorf("parsing encrypted callback JSON: %w", err)
		}
		if dataEncrypt, ok := wrapper["dataEncrypt"].(string); ok && dataEncrypt != "" {
			encryptedData = dataEncrypt
		}
	}

	aesKey, err := decodeAESKey(encodingKey)
	if err != nil {
		return nil, fmt.Errorf("decoding AES key: %w", err)
	}

	iv := aesKey[:16]

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		decoded, err2 := base64.RawStdEncoding.DecodeString(encryptedData)
		if err2 != nil {
			return nil, fmt.Errorf("base64 decoding encrypted data: %w", err)
		}
		ciphertext = decoded
	}

	raw, err := aesCBCDecrypt(ciphertext, aesKey, iv)
	if err != nil {
		return nil, fmt.Errorf("AES-CBC decryption: %w", err)
	}

	raw, err = pkcs7Unpad(raw)
	if err != nil {
		return nil, fmt.Errorf("PKCS7 unpadding: %w", err)
	}

	if len(raw) < 20 {
		return nil, fmt.Errorf("decrypted data too short: %d bytes (minimum 20)", len(raw))
	}

	randomBytes := raw[:16]
	eventsLen := binary.BigEndian.Uint32(raw[16:20])

	totalAfterHeader := len(raw) - 20
	if uint32(totalAfterHeader) < eventsLen {
		return nil, fmt.Errorf("events length %d exceeds available data %d", eventsLen, totalAfterHeader)
	}

	eventsBytes := raw[20+totalAfterHeader-int(eventsLen):]
	middleBytes := raw[20:20+totalAfterHeader-int(eventsLen)]

	var eventsData []interface{}
	if err := json.Unmarshal(eventsBytes, &eventsData); err != nil {
		var singleEvent interface{}
		if err2 := json.Unmarshal(eventsBytes, &singleEvent); err2 != nil {
			return nil, fmt.Errorf("parsing events JSON: %w", err)
		}
		eventsData = []interface{}{singleEvent}
	}

	middleStr := string(middleBytes)
	orgID, appID := splitOrgAppID(middleStr, knownAppID)

	events := []CallbackEvent{}
	for _, entry := range eventsData {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		eventType := ""
		if t, ok := entryMap["eventType"].(string); ok {
			eventType = t
		} else if t, ok := entryMap["type"].(string); ok {
			eventType = t
		}

		category := CallbackEventTypes[eventType]
		rawData, _ := entryMap["data"].(map[string]interface{})
		parsedData := parseEventData(eventType, rawData)

		eventID := ""
		if id, ok := entryMap["eventId"].(string); ok {
			eventID = id
		} else if id, ok := entryMap["id"].(string); ok {
			eventID = id
		} else if id, ok := entryMap["eventId"].(float64); ok {
			eventID = fmt.Sprintf("%v", id)
		} else if id, ok := entryMap["id"].(float64); ok {
			eventID = fmt.Sprintf("%v", id)
		}

		entryAppID, _ := entryMap["appId"].(string)
		entryOrgID, _ := entryMap["orgId"].(string)
		if entryAppID == "" {
			entryAppID = appID
		}
		if entryOrgID == "" {
			entryOrgID = orgID
		}

		events = append(events, CallbackEvent{
			EventID:   eventID,
			EventType: eventType,
			Category:  category,
			Data:      parsedData,
			RawData:   rawData,
			AppID:     entryAppID,
			OrgID:     entryOrgID,
		})
	}

	return &DecryptedCallbackPayload{
		Random: string(randomBytes),
		OrgID:  orgID,
		AppID:  appID,
		Events: events,
		Length: eventsLen,
	}, nil
}

func VerifyCallbackSignature(timestamp, nonce, signature, encodingKey string, dataEncrypt string, callbackToken string) bool {
	token := callbackToken
	if token == "" {
		token = encodingKey
	}

	params := []string{token, timestamp, nonce, dataEncrypt}
	sort.Strings(params)
	joined := strings.Join(params, "")

	computed := sha1.Sum([]byte(joined))
	computedHex := hex.EncodeToString(computed[:])

	return computedHex == signature
}

func GetCallbackEventTypes() map[string]string {
	return CallbackEventTypes
}

func decodeAESKey(encodingKey string) ([]byte, error) {
	padLen := (4 - len(encodingKey)%4) % 4
	padded := encodingKey + strings.Repeat("=", padLen)
	aesKey, err := base64.StdEncoding.DecodeString(padded)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding encoding key: %w", err)
	}
	if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
		return nil, fmt.Errorf("invalid AES key length: %d bytes (expected 16, 24, or 32)", len(aesKey))
	}
	return aesKey, nil
}

func aesCBCDecrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating AES cipher: %w", err)
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length %d is not a multiple of block size %d", len(data), aes.BlockSize)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(data))
	mode.CryptBlocks(dst, data)

	return dst, nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data for PKCS7 unpadding")
	}
	padLen := int(data[len(data)-1])
	if padLen < 1 || padLen > 32 {
		return nil, fmt.Errorf("invalid PKCS7 padding value: %d", padLen)
	}
	if padLen > len(data) {
		return nil, fmt.Errorf("PKCS7 padding %d exceeds data length %d", padLen, len(data))
	}
	for i := 0; i < padLen; i++ {
		if data[len(data)-1-i] != byte(padLen) {
			return nil, fmt.Errorf("invalid PKCS7 padding bytes at position %d", len(data)-1-i)
		}
	}
	return data[:len(data)-padLen], nil
}

func splitOrgAppID(middleStr string, knownAppID string) (string, string) {
	if middleStr == "" {
		return "", ""
	}
	if knownAppID != "" && strings.HasSuffix(middleStr, knownAppID) {
		return middleStr[:len(middleStr)-len(knownAppID)], knownAppID
	}
	if knownAppID != "" {
		idx := strings.Index(middleStr, knownAppID)
		if idx >= 0 {
			return middleStr[:idx], knownAppID
		}
	}
	return middleStr, ""
}

var callbackFieldMaps = map[string]map[string]string{
	"account_subscribe":    {"staffId": "staff_id", "createTime": "create_time"},
	"account_unsubscribe":  {"staffId": "staff_id", "createTime": "create_time"},
	"staff_info":           {"staffId": "staff_id", "name": "name", "mobile": "mobile", "state": "state", "sex": "sex", "email": "email", "employId": "employee_id", "avatarId": "avatar_id", "timestamp": "timestamp"},
	"staff_modify":         {"staffId": "staff_id", "timestamp": "timestamp"},
	"staff_create":         {"staffId": "staff_id", "timestamp": "timestamp"},
	"staff_delete":         {"staffId": "staff_id", "timestamp": "timestamp"},
	"telephone_track":      {"transactionId": "transaction_id", "attach": "attach", "confirmType": "confirm_type", "timestamp": "timestamp"},
	"dept_create":          {"deptId": "dept_id", "timestamp": "timestamp"},
	"dept_modify":          {"deptId": "dept_id", "timestamp": "timestamp"},
	"dept_delete":          {"deptId": "dept_id", "timestamp": "timestamp"},
	"app_install_org":      {"orgId": "org_id", "orgName": "org_name", "timestamp": "timestamp"},
	"app_uninstall_org":    {"orgId": "org_id", "orgName": "org_name", "timestamp": "timestamp"},
	"ua_cert_create":       {"staffId": "staff_id", "deviceId": "device_id", "uaCert": "ua_cert", "timestamp": "timestamp"},
	"ua_cert_delete":       {"staffId": "staff_id", "deviceId": "device_id", "timestamp": "timestamp"},
	"report_location":      {},
	"user_logout":          {"staffId": "staff_id", "deviceId": "device_id", "timestamp": "timestamp"},
	"data_scope":           {"deptIds": "dept_ids", "timestamp": "timestamp"},
	"bot_private_message":  {"from": "from_id", "entryId": "entry_id", "msgType": "msg_type", "msgData": "msg_data", "msgId": "msg_id", "referenceMsg": "reference_msg"},
	"bot_group_message":    {"from": "from_id", "entryId": "entry_id", "msgType": "msg_type", "msgData": "msg_data", "groupId": "group_id", "fromType": "from_type", "groupName": "group_name", "botCreator": "bot_creator", "msgId": "msg_id", "botId": "bot_id", "isAtMe": "is_at_me", "isAtAll": "is_at_all", "referenceMsg": "reference_msg"},
	"wb_visible_config":    {"entryId": "entry_id", "departmentIds": "department_ids", "staffIds": "staff_ids", "timestamp": "timestamp", "isTestModeOn": "is_test_mode_on"},
	"group_create_approve": {"applyRequestId": "apply_request_id", "groupId": "group_id", "timestamp": "timestamp"},
	"schedule_modify":      {"primaryScheduleId": "primary_schedule_id", "scheduleId": "schedule_id", "summary": "summary", "description": "description", "operationType": "operation_type", "currentTime": "current_time", "repeatType": "repeat_type", "expireDateType": "expire_date_type", "allDay": "all_day", "rule": "rule", "ruleStartTime": "rule_start_time", "ruleEndTime": "rule_end_time", "startTime": "start_time", "endTime": "end_time", "operator": "operator", "attendees": "attendees", "timestamp": "timestamp"},
	"schedule_delete":      {"primaryScheduleId": "primary_schedule_id", "scheduleId": "schedule_id", "summary": "summary", "description": "description", "operationType": "operation_type", "currentTime": "current_time", "repeatType": "repeat_type", "expireDateType": "expire_date_type", "allDay": "all_day", "rule": "rule", "ruleStartTime": "rule_start_time", "ruleEndTime": "rule_end_time", "startTime": "start_time", "endTime": "end_time", "operator": "operator", "timestamp": "timestamp"},
	"tag_member":           {"tagId": "tag_id", "timestamp": "timestamp"},
}

func parseEventData(eventType string, rawData map[string]interface{}) map[string]interface{} {
	if rawData == nil {
		return nil
	}

	fieldMap, ok := callbackFieldMaps[eventType]
	if !ok {
		return rawData
	}

	if eventType == "telephone_track" {
		return parseTelephoneTrackData(rawData)
	}

	if eventType == "report_location" {
		result := map[string]interface{}{}
		if loc, ok := rawData["locationInfo"]; ok {
			result["location_info"] = loc
		}
		for k, v := range rawData {
			if k != "locationInfo" {
				result[k] = v
			}
		}
		return result
	}

	if len(fieldMap) == 0 {
		return rawData
	}

	result := map[string]interface{}{}
	for apiKey, goKey := range fieldMap {
		if val, ok := rawData[apiKey]; ok {
			result[goKey] = val
		}
	}
	for k, v := range rawData {
		if _, mapped := fieldMap[k]; !mapped {
			result[k] = v
		}
	}

	return result
}

func parseTelephoneTrackData(rawData map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"transaction_id": rawData["transactionId"],
		"attach":         rawData["attach"],
		"confirm_type":   rawData["confirmType"],
		"timestamp":      rawData["timestamp"],
	}

	if caller, ok := rawData["caller"].(map[string]interface{}); ok {
		result["caller"] = parseTelephoneTrackCaller(caller)
	}
	if callee, ok := rawData["callee"].(map[string]interface{}); ok {
		result["callee"] = parseTelephoneTrackCaller(callee)
	}

	for k, v := range rawData {
		if k != "transactionId" && k != "attach" && k != "confirmType" && k != "timestamp" && k != "caller" && k != "callee" {
			result[k] = v
		}
	}

	return result
}

func parseTelephoneTrackCaller(raw map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"staff_id": raw["staffId"],
	}
	if phone, ok := raw["mobilePhone"].(map[string]interface{}); ok {
		result["country_code"] = phone["countryCode"]
		result["number"] = phone["number"]
	}
	return result
}