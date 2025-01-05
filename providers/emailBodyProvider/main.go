package emailbodyprovider

import "fmt"

func GetResetPasswordBody(unionID string, username string, link string) string {
	var contentBody string
	if unionID == "downsyndrome" {
		contentBody = DSRequestPasswordReset
	} else {
		contentBody = younifiedPasswordReset
	}
	return fmt.Sprintf(contentBody, username, link)

}
