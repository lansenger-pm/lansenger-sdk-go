package lansenger

import (
	"context"
)

// NoticeSendParams carries the fields for SendNotice (通知系统 /v1/send).
//
// Flag-style fields use *int: nil means "omit" (server default applies),
// 1=yes, 0=no. One of CreateMobile / CreateUserID is required.
type NoticeSendParams struct {
	Title       string
	ContentType int    // 1=text (Content required), 2=link (NoticeLink required)
	AccountCode string // official account CODE — see FetchNoticeAccounts
	UserType    int    // 1=phone (ReleasePhones), 2=staff/department (ReleaseRange)

	Content        string
	NoticeLink     string
	NoticeLocation string
	Latitude       float64
	LatitudeSet    bool // include Latitude even when its value is zero
	Longitude      float64
	LongitudeSet   bool // include Longitude even when its value is zero

	ReleasePhones []string                 // max 10
	CCPhones      []string                 // max 10
	ReleaseRange  []map[string]interface{} // items {objId, objName, objType: 1=staff|2=department}, max 200
	CCStaffIDs    []string                 // max 200

	CreateMobile string // operator mobile (UserType=1)
	CreateUserID string // creator staff ID (UserType=2)

	ResourceList []map[string]interface{}
	ExtendID     string

	ConfirmFlag   *int
	ForwardFlag   *int
	ReplyFlag     *int
	AnonymousFlag *int

	RemindStatus               *int
	RemindMsgType              string // comma-separated: mobile,sms,app
	AtOnceFlag                 *int
	RemindAfterType            string // never / unOperate / count
	RemindMaxCount             *int
	RemindIntervalTime         *int
	RemindIntervalTimeDuration string // minutes / hour / day
	RemindRangeType            string // all / receiver / partialRemind / notReminder
	RemindRangeStaffIDs        []string

	UserToken string
}

// SendNotice sends a notice via an official account
// (通知系统 /xtra/notice/server/openapi/v1/send).
//
// The notice module has no revoke/delete interface. Paths carry a /server
// segment (production stage; dev/test environments omit it).
func (c *LansengerClient) SendNotice(ctx context.Context, p *NoticeSendParams) (*NoticeSendResult, error) {
	if p == nil {
		return &NoticeSendResult{Success: false, Error: "params is required"}, nil
	}

	userToken := p.UserToken
	if userToken == "" {
		userToken = getDefaultUserToken()
	}
	createUserID := p.CreateUserID
	if createUserID == "" {
		createUserID = getDefaultUserID()
	}

	if p.Title == "" {
		return &NoticeSendResult{Success: false, Error: "title is required"}, nil
	}
	if p.ContentType != 1 && p.ContentType != 2 {
		return &NoticeSendResult{Success: false, Error: "content_type must be 1 (text) or 2 (link)"}, nil
	}
	if p.ContentType == 1 && p.Content == "" {
		return &NoticeSendResult{Success: false, Error: "content is required when content_type is 1"}, nil
	}
	if p.ContentType == 2 && p.NoticeLink == "" {
		return &NoticeSendResult{Success: false, Error: "notice_link is required when content_type is 2"}, nil
	}
	if p.AccountCode == "" {
		return &NoticeSendResult{Success: false, Error: "account_code is required"}, nil
	}
	if p.UserType != 1 && p.UserType != 2 {
		return &NoticeSendResult{Success: false, Error: "user_type must be 1 (phone) or 2 (openid)"}, nil
	}
	if p.UserType == 1 {
		if len(p.ReleasePhones) == 0 {
			return &NoticeSendResult{Success: false, Error: "release_phones is required when user_type is 1"}, nil
		}
		if len(p.ReleasePhones) > 10 || len(p.CCPhones) > 10 {
			return &NoticeSendResult{Success: false, Error: "release_phones and cc_phones allow at most 10 numbers"}, nil
		}
		if p.CreateMobile == "" && createUserID == "" {
			return &NoticeSendResult{Success: false, Error: "create_mobile or create_user_id is required when user_type is 1"}, nil
		}
	}
	if p.UserType == 2 {
		if len(p.ReleaseRange) == 0 {
			return &NoticeSendResult{Success: false, Error: "release_range is required when user_type is 2"}, nil
		}
		if len(p.ReleaseRange) > 200 || len(p.CCStaffIDs) > 200 {
			return &NoticeSendResult{Success: false, Error: "release_range and cc_staff_ids allow at most 200 entries"}, nil
		}
		if createUserID == "" && p.CreateMobile == "" {
			return &NoticeSendResult{Success: false, Error: "create_mobile or create_user_id is required when user_type is 2"}, nil
		}
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "notices", "send", token, WithUserToken(userToken))

	body := map[string]interface{}{
		"title":       p.Title,
		"contentType": p.ContentType,
		"accountCode": p.AccountCode,
		"userType":    p.UserType,
	}
	if p.Content != "" {
		body["content"] = p.Content
	}
	if p.NoticeLink != "" {
		body["noticeLink"] = p.NoticeLink
	}
	if p.NoticeLocation != "" {
		body["noticeLocation"] = p.NoticeLocation
	}
	if p.LatitudeSet || p.Latitude != 0 {
		body["latitude"] = p.Latitude
	}
	if p.LongitudeSet || p.Longitude != 0 {
		body["longitude"] = p.Longitude
	}
	if p.UserType == 1 {
		// 实测（stage 2026-09-17）：range 对象内缺失或为 null 的 ccRangeList
		// 均可能触发服务端 NPE（errCode=-1），强制下发且 nil 归一为空数组。
		ccPhones := p.CCPhones
		if ccPhones == nil {
			ccPhones = []string{}
		}
		body["phoneUserRange"] = map[string]interface{}{
			"releaseRangeList": p.ReleasePhones,
			"ccRangeList":      ccPhones,
		}
	}
	if p.UserType == 2 {
		ccStaffIDs := p.CCStaffIDs
		if ccStaffIDs == nil {
			ccStaffIDs = []string{}
		}
		body["openUserRange"] = map[string]interface{}{
			"releaseRangeList": p.ReleaseRange,
			"ccRangeList":      ccStaffIDs,
		}
	}
	if p.CreateMobile != "" {
		body["createMobile"] = p.CreateMobile
	}
	if createUserID != "" {
		body["createUserId"] = createUserID
	}
	if len(p.ResourceList) > 0 {
		body["resourceList"] = p.ResourceList
	}
	if p.ExtendID != "" {
		body["extendId"] = p.ExtendID
	}
	if p.ConfirmFlag != nil {
		body["confirmFlag"] = *p.ConfirmFlag
	}
	if p.ForwardFlag != nil {
		body["forwardFlag"] = *p.ForwardFlag
	}
	if p.ReplyFlag != nil {
		body["replyFlag"] = *p.ReplyFlag
	}
	if p.AnonymousFlag != nil {
		body["anonymousFlag"] = *p.AnonymousFlag
	}
	if p.RemindStatus != nil {
		body["remindStatus"] = *p.RemindStatus
	} else {
		// 实测（stage 2026-09-17）：缺失 remindStatus 时服务端抛
		// errCode=-1 unknown exception（NPE），兜底为 0（不提醒）。
		body["remindStatus"] = 0
	}
	if p.RemindMsgType != "" {
		body["remindMsgType"] = p.RemindMsgType
	}
	if p.AtOnceFlag != nil {
		body["atOnceFlag"] = *p.AtOnceFlag
	}
	if p.RemindAfterType != "" {
		body["remindAfterType"] = p.RemindAfterType
	}
	if p.RemindMaxCount != nil {
		body["remindMaxCount"] = *p.RemindMaxCount
	}
	if p.RemindIntervalTime != nil {
		body["remindIntervalTime"] = *p.RemindIntervalTime
	}
	if p.RemindIntervalTimeDuration != "" {
		body["remindIntervalTimeDuration"] = p.RemindIntervalTimeDuration
	}
	if p.RemindRangeType != "" {
		body["remindRangeType"] = p.RemindRangeType
	}
	if len(p.RemindRangeStaffIDs) > 0 {
		body["remindRangeStaffIds"] = p.RemindRangeStaffIDs
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &NoticeSendResult{Success: false, Error: err.Error()}, nil
	}

	res := &NoticeSendResult{Success: true, RawResponse: result}
	data := extractData(result)
	if data != nil {
		res.NoticeCode = strFromMap(data, "code")
		res.NoticeID = intFromMap(data, "id")
		res.Title = strFromMap(data, "title")
		res.NoticeType = intFromMap(data, "noticeType")
		res.ContentType = intFromMap(data, "contentType")
		res.ContentAbstract = strFromMap(data, "contentAbstract")
		res.NoticeLink = strFromMap(data, "noticeLink")
		res.NoticeStatus = intFromMap(data, "noticeStatus")
		res.ConfirmStatus = intFromMap(data, "confirmStatus")
		res.PublishTime = int64(intFromMap(data, "publishTime"))
		res.PublishUserID = strFromMap(data, "publishUserId")
		res.PublishUserName = strFromMap(data, "publishUserName")
	}
	return res, nil
}

// FetchNoticeAccounts lists the official accounts of an organization
// (通知系统 /xtra/notice/server/openapi/v1/notice/account).
// Each account's "code" field is the AccountCode needed by SendNotice.
func (c *LansengerClient) FetchNoticeAccounts(ctx context.Context, orgID, userToken string) (*NoticeAccountListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "notices", "accounts_fetch", token, WithUserToken(userToken))

	body := map[string]interface{}{}
	if orgID != "" {
		body["orgId"] = orgID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &NoticeAccountListResult{Success: false, Error: err.Error()}, nil
	}

	res := &NoticeAccountListResult{Success: true, RawResponse: result}
	if list := extractDataArray(result); list != nil {
		res.Accounts = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Accounts = append(res.Accounts, m)
			}
		}
		res.Total = len(res.Accounts)
	}
	return res, nil
}
