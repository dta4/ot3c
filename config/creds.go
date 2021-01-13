package config

import huaweiotc "github.com/huaweicloud/golangsdk"

//APIKey is the APIkey for the Enterprise Dashboard
var APIKey string
var OTCCreds huaweiotc.AKSKAuthOptions

//OTCUsingOTP prompts user for OTP on login
var OTCUsingOTP bool
