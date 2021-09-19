package sms

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"strings"
	"template/internal/constant"
	"template/internal/constant/errors"
	"template/internal/constant/model"
	"template/internal/module/notification/sms"
)
//SmsHandler contains all handler interfaces
type SmsHandler interface {
	MiddleWareValidateSmsMessage(c *gin.Context)
	SendSmsMessage(c *gin.Context)
	GetCountUnreadSMsMessages(c *gin.Context)
}
//smsHandler implements sms servicea and golang validator object
type smsHandler struct {
	smsUseCase        sms.Usecase
	validate            *validator.Validate
}

//NewSmsHandler  initializes notification services and golang validator
func NewSmsHandler(s sms.Usecase, valid *validator.Validate) SmsHandler {
	return &smsHandler{smsUseCase: s, validate:    valid,}
}
//MiddleWareValidateSmsMessage binds sms data SMS struct
func (n smsHandler) MiddleWareValidateSmsMessage(c *gin.Context) {
	sms := model.SMS{}
	err := c.Bind(&sms)
	if err != nil {
		errValue := errors.ErrorModel{
			ErrorCode:        strconv.Itoa(errors.StatusCodes[errors.ErrInvalidRequest]),
			ErrorDescription: errors.Descriptions[errors.ErrInvalidRequest],
			ErrorMessage:     errors.ErrInvalidRequest.Error(),
		}
		constant.ResponseJson(c, errValue, errors.StatusCodes[errors.ErrInvalidRequest])
	}
	errV := constant.StructValidator(sms, n.validate)
	if errV != nil {
		constant.ResponseJson(c, errV, errors.StatusCodes[errors.ErrorUnableToBindJsonToStruct])
	}
	c.Set("x-sms", sms)
	c.Next()
}
//SendSmsMessage  sends sms message to a user via phone number
func (n smsHandler) SendSmsMessage(c *gin.Context) {
	sms := c.MustGet("x-sms").(model.SMS)
	// TODO:01 sms notification code put here
	_, err := SendSmsMessage(sms)
	if err != nil {
		errValue := errors.ErrorModel{
			ErrorCode:        strconv.Itoa(errors.StatusCodes[errors.ErrUnableToSendSmsMessage]),
			ErrorDescription: errors.Descriptions[errors.ErrUnableToSendSmsMessage],
			ErrorMessage:     errors.ErrUnableToSendSmsMessage.Error(),
		}
		constant.ResponseJson(c, errValue, errors.StatusCodes[errors.ErrorUnableToConvert])
	}
	// TODO:02 sms notification data store in the database put here
	data,errData:=n.smsUseCase.SendSmsMessage(sms)
	if errData!= nil {
		code, _ :=strconv.Atoi(errData.ErrorCode)
		constant.ResponseJson(c, *errData, code)
	}
	constant.ResponseJson(c, *data, data.Code)
}
//GetCountUnreadSMsMessages counts unread sms message
func (n smsHandler) GetCountUnreadSMsMessages(c *gin.Context) {
	count:=n.smsUseCase.GetCountUnreadSmsMessages()
	constant.ResponseJson(c, map[string]interface{}{"count": count}, http.StatusOK)
}
//SendSmsMessage sends sms message via phone number
func SendSmsMessage(sms model.SMS) (interface{}, error) {
	reqString := fmt.Sprintf(`
		{
			"type":"text",
			"auth":{"username":"%s", "password":"%s"},
			"sender":"%s",
			"receiver":"%s",
			"dcs":"GSM",
			"text":"%s",
			"dlrMask":3,
			"dlrUrl":"%s"
        }
	`, sms.User, sms.Password, sms.SenderId, sms.ReceiverPhone, sms.Body, sms.CallBackUrl)
	requestBody := strings.NewReader(reqString)
	// post some data
	res, err := http.Post(sms.ApiGateWay, "application/json; charset=UTF-8", requestBody)
	if err != nil {
		return nil,errors.ErrUnableToSendSmsMessage
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, errors.ErrUnableToSendSmsMessage
	}
	// read response data
	var smsResponseData interface{}
	err = json.NewDecoder(res.Body).Decode(&smsResponseData)
	if err != nil {
		return nil, errors.ErrUnableToSendSmsMessage
	}
	err = res.Body.Close()
	if err != nil {
		return nil, errors.ErrUnableToSendSmsMessage
	}
	return smsResponseData, nil
}