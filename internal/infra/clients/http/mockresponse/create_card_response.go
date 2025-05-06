package mockresponse

const SuccessCardResponseMock = `{
  "data": {
    "id": "12345",
    "user_id": "67890",
    "card_holder": "John Doe",
    "card_holder_dni": "987654321",
    "company": "Visa",
    "last_four": "1234",
    "first_six": "123456",
    "expiration_date": "12/24",
    "is_valid": true,
    "is_expired": false,
    "verified_by": "System",
    "verified_with_transaction_id": "abc123",
    "verified_at": "2023-10-12T15:00:00Z",
    "last_used": "2023-10-15T09:30:00Z",
    "created_at": "2023-01-01T08:00:00Z",
    "updated_at": "2023-10-01T10:00:00Z",
    "deleted_at": "",
    "bank_name": "Bank of Somewhere",
    "country_iso": "US",
    "card_type": "Credit",
    "source": "Online",
    "zip_code": "12345",
    "vault": "VaultService"
  }
}`

const EmptyDataMock = `{
  "data": {}
}`
