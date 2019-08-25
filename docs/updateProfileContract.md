[Back](../README.md)

# Update Profile

Update existing Profiles.

**URL** : `api/v1/profiles/{id}`

**Method** : `PUT`

**Auth required** : NO

**Permissions required** : NO

**Data example** Email is required, firstName and lastName fields are alpha

```json
{
	"firstName" : "Trace",
	"lastName" : "Ohrt",
	"address" : {
		"street": "1600 Pennsylvania Ave NW",
		"city": "Washington",
		"state" : "DC",
		"zipcode" : "20500"
	},
	"email" : "teohrt18@gmail.com"
}
```

## Success Response

**Condition** : If request is good and profile was updated

**Code** : `200 OK`

## Error Responses

**Condition** : If the email field is modified

**Code** : `400 BAD REQUEST`

**Headers** : `Location: http://testserver/api/profiles/123/`

**Response** : 
```json
{
    "status": "Bad Request",
    "message": "UpdateProfile failed: attempted to change email",
    "error": "Email inconsistent with ProfileID"
}
```

### Or

**Condition** : If fields are determined to be invalid.

**Code** : `400 BAD REQUEST`

**Response**

```json
{
    "status": "Bad Request",
    "message": "Profile validation failed",
    "error": "Key: 'ProfileData.LastName' Error:Field validation for 'LastName' failed on the 'alpha' tag"
}
```