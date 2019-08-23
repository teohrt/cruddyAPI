# Get Profile

Retrieve Profile data for a user specified by ProfileID.

**URL** : `api/v1/cruddyAPI/profiles/{id}`

**Method** : `GET`

**Auth required** : NO

**Permissions required** : NO

**Data example** ID parameter required. No payload required.

## Success Response

**Condition** : If profile exists

**Code** : `200 OK`

**Content example**

```json
{
	"firstName" : "Trace",
	"lastName" : "Ohrt",
	"address" : {
		"street": "175 Calvert Dr",
		"city": "Cupertino",
		"state" : "California",
		"zipcode" : "95014"
	},
	"email" : "teohrt18@gmail.com"
}
```

## Error Responses

**Condition** : If profile does not exist

**Code** : `404 NOT FOUND`

**Content** : 
```json
{
    "status": "Not Found",
    "message": "Profile not found",
    "error": "Could not find profile associated with: 333"
}
```

### Or

**Condition** : If the server chokes on something

**Code** : `500 INTERNAL SERVER ERROR`

**Content**
```json
{
    "status": "Internal Server Error",
    "message": "Get profile failed",
    "error": "RequestError: send request failed\ncaused by: Post http://localhost:8000/: dial tcp [::1]:8000: connect: connection refused"
}
```