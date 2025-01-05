package emailbodyprovider

const DSRequestPasswordReset = `
		<!DOCTYPE html>
			<html xmlns="http://www.w3.org/1999/xhtml">
			<head>
			    <meta charset="utf-8">
			    <title>Password Reset Request</title>
			</head>
			<body>
			    <div>
			        <p>We understand how this can be annoying when one forgets their password, ugh. Great news, we've just received a request to reset the password and our team is on standby for the username: <strong>%s </strong>.</p>
			        <p>If you didn't initiate this request, no worries - simply disregard this email.</p>
			        <p>Excitingly, you have the opportunity to reset your password by clicking the link below. But remember, the link expires in just one hour, so let's get that password reset as soon as possible!</p>
			        <p>🚀%s</a></p>
			    </div>
			</body>
			</html>`

const younifiedPasswordReset = `<p> Hello </p> <p>We've received a request to reset the password for the username: <b>%s</b></p><p>If you didn't make this request, please disregard this email.</p><p> You can reset your password by clicking the link below: </p> <p><i> Expires in one hour! </i></p><p>%s</p>`
