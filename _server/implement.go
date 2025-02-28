package main

import (
	"context"
	"fmt"

	"github.com/axiaoxin-com/base64Captcha"
	"github.com/axiaoxin-com/base64Captcha/_server/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type implement struct {
	pb.UnimplementedVcodeServiceServer
}

// Generate Generate方法实现
func (i *implement) Generate(ctx context.Context, req *pb.GenerateReq) (*pb.GenerateRsp, error) {
	height := req.GetHeight()
	if height == 0 {
		height = 40
	}
	width := req.GetWidth()
	if width == 0 {
		width = 140
	}
	length := req.GetLength()
	if length == 0 {
		length = 5
	}

	lang := "zh"
	switch req.GetLang() {
	case pb.LangEnum_Latin:
		lang = "latin"
	case pb.LangEnum_Ko:
		lang = "ko"
	case pb.LangEnum_Jp:
		lang = "jp"
	case pb.LangEnum_Ru:
		lang = "ru"
	case pb.LangEnum_Th:
		lang = "th"
	case pb.LangEnum_Greek:
		lang = "greek"
	case pb.LangEnum_Arabic:
		lang = "arabic"
	case pb.LangEnum_Hebrew:
		lang = "hebrew"
	}

	var driver base64Captcha.Driver
	switch req.GetDriver() {
	case pb.DriverEnum_String:
		driver = base64Captcha.NewDriverString(
			int(height),
			int(width),
			0,
			base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowHollowLine,
			int(length),
			base64Captcha.TxtAlphabet+base64Captcha.TxtNumbers,
			nil,
			nil,
			nil,
		)
	case pb.DriverEnum_Math:
		driver = base64Captcha.NewDriverMath(int(height), int(width), 1, base64Captcha.OptionShowHollowLine, nil, nil, nil)
	case pb.DriverEnum_Chinese:
		driver = base64Captcha.NewDriverChinese(
			int(height),
			int(width),
			40,
			base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowHollowLine,
			int(length),
			base64Captcha.TxtChineseCharaters,
			nil,
			nil,
			[]string{"wqy-microhei.ttc", "LXGWWenKai-Regular.ttf"},
		)
	case pb.DriverEnum_Audio:
		driver = base64Captcha.NewDriverAudio(int(length), lang)
	case pb.DriverEnum_Language:
		driver = base64Captcha.NewDriverLanguage(int(height), int(width), 0, base64Captcha.OptionShowHollowLine, int(length), nil, nil, nil, lang)
	default:
		maxSkew := 0.7
		if req.GetMaxSkew() > 0 {
			maxSkew = req.GetMaxSkew()
		}
		dotCount := 80
		if req.GetDotCount() > 0 {
			dotCount = int(req.GetDotCount())
		}
		driver = base64Captcha.NewDriverDigit(int(height), int(width), int(length), maxSkew, dotCount)
	}

	captcha := base64Captcha.NewCaptcha(driver, RedisStore)
	id, data, err := captcha.Generate()
	if err != nil {
		return nil, status.Error(codes.Internal, "captcha generate error:"+err.Error())
	}

	rsp := pb.GenerateRsp{
		Data: data,
		ID:   id,
	}
	return &rsp, nil
}

// Verify Verify方法实现
func (i *implement) Verify(ctx context.Context, req *pb.VerifyReq) (*pb.VerifyRsp, error) {
	rsp := pb.VerifyRsp{}
	if req.GetAnswer() == "" {
		return nil, status.Error(codes.InvalidArgument, "Answer is required")
	}

	ok := RedisStore.Verify(req.GetID(), req.GetAnswer(), true)
	rsp = pb.VerifyRsp{OK: ok}
	return &rsp, nil
}

// GenRawCode GenRawCode方法实现
func (i *implement) GenRawCode(ctx context.Context, req *pb.GenRawCodeReq) (*pb.GenerateRsp, error) {
	idTag := req.GetIDTag()
	if idTag == "" {
		return nil, status.Error(codes.InvalidArgument, "IDTag is required")
	}
	if len(idTag) > 64 {
		return nil, status.Error(codes.InvalidArgument, "IDTag max length is 64")
	}
	length := req.GetLength()
	if length == 0 {
		length = 5
	}
	driver := base64Captcha.NewDriverString(1, 1, 0, 0, int(length), base64Captcha.TxtSimpleCharaters, nil, nil, nil)
	captcha := base64Captcha.NewCaptcha(driver, RedisStore)
	qid, _, answer := captcha.Driver.GenerateIdQuestionAnswer()
	id := fmt.Sprintf("%v:%v", idTag, qid)
	if err := captcha.Store.Set(id, answer); err != nil {
		return nil, status.Error(codes.Internal, "store set error:"+err.Error())
	}
	rsp := pb.GenerateRsp{
		Data: answer,
		ID:   id,
	}
	return &rsp, nil
}
