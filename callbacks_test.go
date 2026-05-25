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
	"strings"
	"testing"
)

const testAESKeyB64 = "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY="

func testEncryptPayload(eventsJSON string, orgID string, appID string, encodingKey string) string {
	aesKey, _ := base64.StdEncoding.DecodeString(encodingKey)
	iv := aesKey[:16]

	eventsBytes := []byte(eventsJSON)
	eventsLen := make([]byte, 4)
	binary.BigEndian.PutUint32(eventsLen, uint32(len(eventsBytes)))

	randomBytes := []byte("random16bytes!!!")
	plaintext := append(randomBytes, eventsLen...)
	plaintext = append(plaintext, []byte(orgID)...)
	plaintext = append(plaintext, []byte(appID)...)
	plaintext = append(plaintext, eventsBytes...)

	padLen := 32 - (len(plaintext) % 32)
	for i := 0; i < padLen; i++ {
		plaintext = append(plaintext, byte(padLen))
	}

	block, _ := aes.NewCipher(aesKey)
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func TestDecryptCallbackPayload(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_modify","data":{"staffId":"s001"},"eventId":"e1"}]`
	orgID := "3211264"
	appID := "2285568-12042496"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OrgID != orgID {
		t.Errorf("expected orgId=%s, got %s", orgID, result.OrgID)
	}
	if result.AppID != appID {
		t.Errorf("expected appId=%s, got %s", appID, result.AppID)
	}
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result.Events))
	}

	event := result.Events[0]
	if event.EventType != "staff_modify" {
		t.Errorf("expected eventType=staff_modify, got %s", event.EventType)
	}
	if event.Category != "staff" {
		t.Errorf("expected category=staff, got %s", event.Category)
	}
	if event.Data["staff_id"] != "s001" {
		t.Errorf("expected data.staff_id=s001, got %v", event.Data["staff_id"])
	}
}

func TestDecryptCallbackPayloadJSONWrapper(t *testing.T) {
	eventsJSON := `[{"eventType":"bot_private_message","data":{"from":"staff1","entryId":"entry1","msgType":"text","msgData":{"content":"hello"}}}]`
	orgID := "org123"
	appID := "app456"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	wrapper := map[string]interface{}{
		"dataEncrypt": encrypted,
	}
	wrapperJSON, _ := json.Marshal(wrapper)

	result, err := DecryptCallbackPayload(string(wrapperJSON), testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OrgID != orgID {
		t.Errorf("expected orgId=%s, got %s", orgID, result.OrgID)
	}
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result.Events))
	}

	event := result.Events[0]
	if event.EventType != "bot_private_message" {
		t.Errorf("expected eventType=bot_private_message, got %s", event.EventType)
	}
	if event.Data["from_id"] != "staff1" {
		t.Errorf("expected data.from_id=staff1, got %v", event.Data["from_id"])
	}
}

func TestDecryptCallbackPayloadMultipleEvents(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_modify","data":{"staffId":"s001"},"eventId":"e1"},{"eventType":"dept_create","data":{"deptId":"d001"},"eventId":"e2"}]`
	orgID := "3211264"
	appID := "2285568-12042496"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result.Events))
	}

	if result.Events[0].EventType != "staff_modify" {
		t.Errorf("expected first event=staff_modify, got %s", result.Events[0].EventType)
	}
	if result.Events[1].EventType != "dept_create" {
		t.Errorf("expected second event=dept_create, got %s", result.Events[1].EventType)
	}
	if result.Events[1].Data["dept_id"] != "d001" {
		t.Errorf("expected dept_id=d001, got %v", result.Events[1].Data["dept_id"])
	}
}

func TestDecryptCallbackPayloadTelephoneTrack(t *testing.T) {
	eventsJSON := `[{"eventType":"telephone_track","data":{"transactionId":"tx1","caller":{"staffId":"staff1","mobilePhone":{"countryCode":"86","number":"13800138000"}},"callee":{"staffId":"staff2","mobilePhone":{"countryCode":"86","number":"13900139000"}},"confirmType":"1","timestamp":"1234567890"}}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	event := result.Events[0]
	if event.EventType != "telephone_track" {
		t.Errorf("expected eventType=telephone_track, got %s", event.EventType)
	}
	if event.Data["transaction_id"] != "tx1" {
		t.Errorf("expected transaction_id=tx1, got %v", event.Data["transaction_id"])
	}

	caller, ok := event.Data["caller"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected caller to be map, got %T", event.Data["caller"])
	}
	if caller["staff_id"] != "staff1" {
		t.Errorf("expected caller.staff_id=staff1, got %v", caller["staff_id"])
	}
	if caller["country_code"] != "86" {
		t.Errorf("expected caller.country_code=86, got %v", caller["country_code"])
	}

	callee, ok := event.Data["callee"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected callee to be map, got %T", event.Data["callee"])
	}
	if callee["staff_id"] != "staff2" {
		t.Errorf("expected callee.staff_id=staff2, got %v", callee["staff_id"])
	}
}

func TestDecryptCallbackPayloadInvalidKey(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_modify"}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	_, err := DecryptCallbackPayload(encrypted, "invalidkey", "")
	if err == nil {
		t.Error("expected error for invalid key, got nil")
	}
}

func TestDecryptCallbackPayloadWrongKey(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_modify"}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	wrongKey := base64.StdEncoding.EncodeToString([]byte("wrongkey1234567890wrongkey12345"))
	_, err := DecryptCallbackPayload(encrypted, wrongKey, "")
	if err == nil {
		t.Error("expected error for wrong decryption key")
	}
}

func TestVerifyCallbackSignatureSHA1(t *testing.T) {
	token := "test_token_123"
	timestamp := "1234567890"
	nonce := "nonce123"
	dataEncrypt := "encrypted_data_here"

	params := []string{token, timestamp, nonce, dataEncrypt}
	sorted := make([]string, len(params))
	copy(sorted, params)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	joined := strings.Join(sorted, "")
	computed := sha1.Sum([]byte(joined))
	signature := hex.EncodeToString(computed[:])

	valid := VerifyCallbackSignature(timestamp, nonce, signature, token, dataEncrypt, "")
	if !valid {
		t.Errorf("expected signature to be valid, got invalid")
	}

	valid2 := VerifyCallbackSignature(timestamp, nonce, "wrong_signature", token, dataEncrypt, "")
	if valid2 {
		t.Errorf("expected wrong signature to be invalid, got valid")
	}
}

func TestVerifyCallbackSignatureWithCallbackToken(t *testing.T) {
	encodingKey := "encoding_key_abc"
	callbackToken := "callback_token_xyz"
	timestamp := "1234567890"
	nonce := "nonce123"
	dataEncrypt := ""

	params := []string{callbackToken, timestamp, nonce, dataEncrypt}
	sorted := make([]string, len(params))
	copy(sorted, params)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	joined := strings.Join(sorted, "")
	computed := sha1.Sum([]byte(joined))
	signature := hex.EncodeToString(computed[:])

	valid := VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, callbackToken)
	if !valid {
		t.Errorf("expected signature with callback_token to be valid")
	}

	validNoToken := VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, "")
	if validNoToken {
		t.Errorf("expected signature without callback_token to be invalid (uses encoding_key)")
	}
}

func TestParseCallbackPayloadQueryString(t *testing.T) {
	events, err := ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "staff_modify" {
		t.Errorf("expected eventType=staff_modify, got %s", events[0].EventType)
	}
	if events[0].Category != "staff" {
		t.Errorf("expected category=staff, got %s", events[0].Category)
	}
}

func TestParseCallbackPayloadMultipleEvents(t *testing.T) {
	events, err := ParseCallbackPayload("eventType=staff_modify&eventType=dept_create&staffId=s001&deptId=d001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].EventType != "staff_modify" {
		t.Errorf("expected first event=staff_modify, got %s", events[0].EventType)
	}
	if events[1].EventType != "dept_create" {
		t.Errorf("expected second event=dept_create, got %s", events[1].EventType)
	}
	if events[1].Category != "department" {
		t.Errorf("expected category=department, got %s", events[1].Category)
	}
}

func TestParseCallbackPayloadPlainJSON(t *testing.T) {
	input := `{"events":[{"eventType":"staff_modify","data":{"staffId":"s001"},"eventId":"e1"}],"orgId":"org1","appId":"app1"}`
	events, err := ParseCallbackPayload(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "staff_modify" {
		t.Errorf("expected eventType=staff_modify, got %s", events[0].EventType)
	}
	if events[0].OrgID != "org1" {
		t.Errorf("expected orgId=org1, got %s", events[0].OrgID)
	}
	if events[0].AppID != "app1" {
		t.Errorf("expected appId=app1, got %s", events[0].AppID)
	}
	if events[0].Data["staff_id"] != "s001" {
		t.Errorf("expected staff_id=s001, got %v", events[0].Data["staff_id"])
	}
}

func TestParseCallbackPayloadEncryptedJSONError(t *testing.T) {
	input := `{"dataEncrypt":"someencrypteddata","timestamp":"123"}`
	_, err := ParseCallbackPayload(input)
	if err == nil {
		t.Error("expected error for encrypted data without encoding_key")
	}
	if !strings.Contains(err.Error(), "encoding_key") {
		t.Errorf("expected error to mention encoding_key, got: %v", err)
	}
}

func TestParseCallbackPayloadInvalid(t *testing.T) {
	_, err := ParseCallbackPayload("not=valid=query=string===extra")
	if err != nil {
		t.Logf("ParseCallbackPayload handled invalid input: %v", err)
	}
}

func TestPKCS7Unpad(t *testing.T) {
	data := []byte("hello world\x05\x05\x05\x05\x05")
	result, err := pkcs7Unpad(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(result))
	}

	invalidData := []byte("hello\x03\x04\x03")
	_, err = pkcs7Unpad(invalidData)
	if err == nil {
		t.Error("expected error for invalid padding")
	}
}

func TestDecodeAESKey(t *testing.T) {
	key, err := decodeAESKey(testAESKeyB64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(key))
	}

	shortKey := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	key2, err := decodeAESKey(shortKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(key2) != 32 {
		t.Errorf("expected 32-byte key, got %d", len(key2))
	}

	_, err = decodeAESKey("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestSplitOrgAppID(t *testing.T) {
	orgID, appID := splitOrgAppID("32112642285568-12042496", "2285568-12042496")
	if orgID != "3211264" {
		t.Errorf("expected orgId=3211264, got %s", orgID)
	}
	if appID != "2285568-12042496" {
		t.Errorf("expected appId=2285568-12042496, got %s", appID)
	}

	orgID2, appID2 := splitOrgAppID("32112642285568-12042496", "")
	if appID2 != "" {
		t.Errorf("expected empty appId without knownAppID, got %s", appID2)
	}
	if orgID2 != "32112642285568-12042496" {
		t.Errorf("expected full middle as orgId, got %s", orgID2)
	}

	orgID3, appID3 := splitOrgAppID("", "")
	if orgID3 != "" || appID3 != "" {
		t.Errorf("expected empty strings for empty input")
	}

	orgID4, appID4 := splitOrgAppID("org123app456", "app456")
	if orgID4 != "org123" || appID4 != "app456" {
		t.Errorf("expected orgId=org123 appId=app456, got orgId=%s appId=%s", orgID4, appID4)
	}

	orgID5, appID5 := splitOrgAppID("prefixapp456middleapp456suffix", "app456")
	if appID5 != "app456" {
		t.Errorf("expected appId=app456, got %s", appID5)
	}
	if orgID5 != "prefix" {
		t.Errorf("expected orgId=prefix, got %s", orgID5)
	}
}

func TestGetCallbackEventTypes(t *testing.T) {
	types := GetCallbackEventTypes()
	if len(types) == 0 {
		t.Error("expected non-empty callback event types")
	}
	if types["staff_modify"] != "staff" {
		t.Errorf("expected staff_modify=staff, got %s", types["staff_modify"])
	}
	if types["bot_private_message"] != "bot" {
		t.Errorf("expected bot_private_message=bot, got %s", types["bot_private_message"])
	}
}

func TestDecryptAndVerifyFlow(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_modify","data":{"staffId":"s001"},"eventId":"e1"}]`
	orgID := "3211264"
	appID := "2285568-12042496"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	token := testAESKeyB64
	timestamp := "1234567890"
	nonce := "nonce123"
	dataEncrypt := encrypted

	params := []string{token, timestamp, nonce, dataEncrypt}
	sorted := make([]string, len(params))
	copy(sorted, params)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	joined := strings.Join(sorted, "")
	computed := sha1.Sum([]byte(joined))
	signature := hex.EncodeToString(computed[:])

	valid := VerifyCallbackSignature(timestamp, nonce, signature, testAESKeyB64, dataEncrypt, "")
	if !valid {
		t.Errorf("expected signature to be valid in full decrypt+verify flow")
	}

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected decrypt error: %v", err)
	}
	if result.OrgID != orgID {
		t.Errorf("expected orgId=%s, got %s", orgID, result.OrgID)
	}
}

func TestDecryptCallbackPayloadWithKnownAppID(t *testing.T) {
	eventsJSON := `[{"eventType":"staff_info","data":{"staffId":"s001","name":"张三","mobile":"13800138000"}}]`
	orgID := "3211264"
	appID := "2285568-12042496"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AppID != appID {
		t.Errorf("expected appId=%s, got %s", appID, result.AppID)
	}
	if result.OrgID != orgID {
		t.Errorf("expected orgId=%s, got %s", orgID, result.OrgID)
	}

	event := result.Events[0]
	if event.Data["staff_id"] != "s001" {
		t.Errorf("expected staff_id=s001, got %v", event.Data["staff_id"])
	}
	if event.Data["name"] != "张三" {
		t.Errorf("expected name=张三, got %v", event.Data["name"])
	}
}

func TestFieldMapParsing(t *testing.T) {
	rawData := map[string]interface{}{
		"staffId":  "s001",
		"mobile":   "13800138000",
		"state":    "active",
		"timestamp": "1234567890",
	}
	parsed := parseEventData("staff_modify", rawData)
	if parsed["staff_id"] != "s001" {
		t.Errorf("expected staff_id=s001, got %v", parsed["staff_id"])
	}
	if parsed["timestamp"] != "1234567890" {
		t.Errorf("expected timestamp=1234567890, got %v", parsed["timestamp"])
	}

	rawData2 := map[string]interface{}{
		"from":    "staff1",
		"msgType": "text",
		"msgData": map[string]interface{}{"content": "hello"},
	}
	parsed2 := parseEventData("bot_private_message", rawData2)
	if parsed2["from_id"] != "staff1" {
		t.Errorf("expected from_id=staff1, got %v", parsed2["from_id"])
	}
	if parsed2["msg_type"] != "text" {
		t.Errorf("expected msg_type=text, got %v", parsed2["msg_type"])
	}

	rawData3 := map[string]interface{}{
		"eventType": "unknown_event",
		"someField": "someValue",
	}
	parsed3 := parseEventData("unknown_event", rawData3)
	if parsed3["someField"] != "someValue" {
		t.Errorf("expected unknown event data to pass through unchanged")
	}
}

func TestDecryptCallbackPayloadReportLocation(t *testing.T) {
	eventsJSON := `[{"eventType":"report_location","data":{"locationInfo":{"lat":"39.9","lng":"116.4"}}}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	event := result.Events[0]
	if event.EventType != "report_location" {
		t.Errorf("expected eventType=report_location, got %s", event.EventType)
	}
	locInfo, ok := event.Data["location_info"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected location_info to be map")
	}
	if locInfo["lat"] != "39.9" {
		t.Errorf("expected lat=39.9, got %v", locInfo["lat"])
	}
}

func TestAESCBCDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	iv := key[:16]
	plaintext := []byte("hello world AES test data here!!")

	block, _ := aes.NewCipher(key)
	padLen := 16 - (len(plaintext) % 16)
	padded := append(plaintext, bytesRepeat(byte(padLen), padLen)...)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)

	decrypted, err := aesCBCDecrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unpadded, err := pkcs7Unpad(decrypted)
	if err != nil {
		t.Fatalf("unexpected unpad error: %v", err)
	}

	if string(unpadded) != "hello world AES test data here!!" {
		t.Errorf("expected 'hello world AES test data here!!', got '%s'", string(unpadded))
	}
}

func bytesRepeat(b byte, n int) []byte {
	result := make([]byte, n)
	for i := range result {
		result[i] = b
	}
	return result
}

func TestDecryptCallbackPayloadSingleEventNotArray(t *testing.T) {
	eventsJSON := `{"eventType":"staff_modify","data":{"staffId":"s001"}}`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event from single event JSON, got %d", len(result.Events))
	}
	if result.Events[0].EventType != "staff_modify" {
		t.Errorf("expected eventType=staff_modify, got %s", result.Events[0].EventType)
	}
}

func BenchmarkDecryptCallbackPayload(b *testing.B) {
	eventsJSON := `[{"eventType":"staff_modify","data":{"staffId":"s001"},"eventId":"e1"}]`
	orgID := "3211264"
	appID := "2285568-12042496"
	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	for i := 0; i < b.N; i++ {
		_, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestDecryptCallbackPayloadWith16ByteKey(t *testing.T) {
	key16 := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef"))
	eventsJSON := `[{"eventType":"staff_modify","data":{"staffId":"s001"}}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, key16)

	result, err := DecryptCallbackPayload(encrypted, key16, appID)
	if err != nil {
		t.Fatalf("unexpected error with 16-byte key: %v", err)
	}
	if result.OrgID != orgID {
		t.Errorf("expected orgId=%s, got %s", orgID, result.OrgID)
	}
}

func TestDecryptCallbackPayloadScheduleModify(t *testing.T) {
	eventsJSON := `[{"eventType":"schedule_modify","data":{"primaryScheduleId":"ps1","scheduleId":"s1","summary":"Meeting","operationType":"create","timestamp":"12345"}}]`
	orgID := "org1"
	appID := "app1"

	encrypted := testEncryptPayload(eventsJSON, orgID, appID, testAESKeyB64)

	result, err := DecryptCallbackPayload(encrypted, testAESKeyB64, appID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	event := result.Events[0]
	if event.EventType != "schedule_modify" {
		t.Errorf("expected eventType=schedule_modify, got %s", event.EventType)
	}
	if event.Data["primary_schedule_id"] != "ps1" {
		t.Errorf("expected primary_schedule_id=ps1, got %v", event.Data["primary_schedule_id"])
	}
	if event.Data["operation_type"] != "create" {
		t.Errorf("expected operation_type=create, got %v", event.Data["operation_type"])
	}
}

func verifySigHelper(token, timestamp, nonce, dataEncrypt string) string {
	params := []string{token, timestamp, nonce, dataEncrypt}
	sorted := make([]string, len(params))
	copy(sorted, params)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	joined := strings.Join(sorted, "")
	computed := sha1.Sum([]byte(joined))
	return hex.EncodeToString(computed[:])
}

func TestVerifyCallbackSignatureEmptyDataEncrypt(t *testing.T) {
	token := "mytoken"
	timestamp := "1234567890"
	nonce := "nonce123"
	dataEncrypt := ""
	signature := verifySigHelper(token, timestamp, nonce, dataEncrypt)

	valid := VerifyCallbackSignature(timestamp, nonce, signature, "mytoken", dataEncrypt, "")
	if !valid {
		t.Errorf("expected valid signature with empty dataEncrypt")
	}
}

func TestParseCallbackPayloadBotGroupMessage(t *testing.T) {
	input := `{"events":[{"eventType":"bot_group_message","data":{"from":"staff1","groupId":"g1","msgType":"text","msgData":{"content":"hello"},"isAtMe":true,"isAtAll":false}}]}`
	events, err := ParseCallbackPayload(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "bot_group_message" {
		t.Errorf("expected bot_group_message, got %s", event.EventType)
	}
	if event.Data["from_id"] != "staff1" {
		t.Errorf("expected from_id=staff1, got %v", event.Data["from_id"])
	}
	if event.Data["group_id"] != "g1" {
		t.Errorf("expected group_id=g1, got %v", event.Data["group_id"])
	}
	if event.Data["is_at_me"] != true {
		t.Errorf("expected is_at_me=true, got %v", event.Data["is_at_me"])
	}
}

// Silence the unused import warning for fmt
var _ = fmt.Sprintf