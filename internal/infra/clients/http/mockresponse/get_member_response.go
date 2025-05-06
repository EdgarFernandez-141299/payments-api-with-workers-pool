package mockresponse

const SuccessMemberResponseMock = `{
	"message": "Successfully retrieved member",
	"data": {
		"id": "123",
		"first_name": "John",
		"last_name": "Doe",
		"birth_date": "1990-01-01",
		"phones": [
			{
				"phone_type": "mobile",
				"number": "1234567890",
				"is_default": true,
				"country_id": 1,
				"created_at": "2023-01-01",
				"is_favorite": true
			}
		],
		"emails": [
			{
				"email": "john.doe@example.com",
				"is_default": true,
				"created_at": "2023-01-01",
				"is_coowner": false,
				"is_favorite": true
			}
		]
	}
}`

const MemberNotFoundResponseMock = `{
	"message": "Member found",
	"data": {
		"id": "456", 
		"first_name": "Jane",
		"last_name": "Smith",
		"birth_date": "1985-05-05",
		"phones": [],
		"emails": []
	}
42|}`

const SuccessUserProfileInfoResponseMock = `{
	"id": "123",
	"first_name": "John",
	"last_name": "Doe",
	"email": "john.doe@example.com",
	"enterprise_id": "456"
}`

const UserProfileInfoNotFoundResponseMock = `{
    "code": 404,
    "message": "member not found with ID: 123"
}`
