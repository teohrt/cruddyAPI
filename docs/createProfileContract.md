[Back](../README.md)

# Create Profile

Create a Profile for the user if one does not already exist. Profiles and emails are 1 to 1.

**URL** : `api/v1/cruddyAPI/profiles`

**Method** : `POST`

**Auth required** : NO

**Permissions required** : NO

**Data example** Email is the only required field

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

**Condition** : If everything is OK and a Profile didn't exist for this User.

**Code** : `201 CREATED`

**Content example**

```json
{
    "ProfileID": 123,
}
```

## Error Responses

**Condition** : If profile already exists for User.

**Code** : `400 BAD REQUEST`

**Headers** : `Location: http://testserver/api/profiles/123/`

**Content** : 
```json
{
    "status": "Bad Request",
    "message": "Profile already exists",
    "error": "Can not create profile. Already exists"
}
```

### Or

**Condition** : If required fields are missed.

**Code** : `400 BAD REQUEST`

**Content example**

```json
{
    "status": "Bad Request",
    "message": "Profile validation failed",
    "error": "Key: 'ProfileData.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```