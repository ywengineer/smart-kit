package inf

import "errors"

var ErrExpiredSub = errors.New("ErrExpiredSub")
var ErrFail = errors.New("ErrFail")
var ErrNoChannel = errors.New("ErrNoChannel")
var ErrNoValidator = errors.New("ErrNoValidator")
var IncompletePurchase = errors.New("IncompletePurchase")
var OtherAppPurchase = errors.New("OtherAppPurchase")
