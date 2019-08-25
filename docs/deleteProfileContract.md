[Back](../README.md)

# Delete Profile

Delete a Profile with a given ID

**URL** : `api/v1/profiles/{id}`

**Method** : `DELETE`

**Auth required** : NO

**Permissions required** : NO

**Data example** ID parameter required. No payload required.

## Success Response

**Condition** : If profile exists

**Code** : `200 OK`

## Error Responses

**Condition** : If profile does not exist

**Code** : `404 NOT FOUND`

**Content** : 
```json
{
    "status": "Not Found",
    "message": "Profile not found",
    "error": "Could not find profile associated with: 123"
}
```

### Or

**Condition** : If the server chokes on something

**Code** : `500 INTERNAL SERVER ERROR`

**Content**
```json
{
    "status": "Internal Server Error",
    "message": "DeleteProfile failed",
    "error": "puke"
}
```