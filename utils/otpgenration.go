package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendOtp(to string,otp string) error{
	m:=gomail.NewMessage()
	m.SetHeader("From","pranayrohith.niatian@gmail.com")
	m.SetHeader("To",to)
	m.SetHeader("Subject","Your login otp")
	msg:=fmt.Sprintf("Your OTP is {%v}. Valid for 5 minutes. Do not share with anyone.",otp)
	m.SetBody("text/plain",msg)


	d:=gomail.NewDialer("smtp.gmail.com",587,"pranayrohith.niatian@gmail.com","xpfc pubn asps lmpx")
	err:=d.DialAndSend(m)
	if err!=nil{                         
		fmt.Println("sendotp error",err)
	}
	return err 
}






