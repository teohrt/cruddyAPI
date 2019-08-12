# cruddyAPI
![cruddyAPI](cruddyAPI.jpeg)

The purpose of this RESTful CRUD API is to manage profile data.

## Routes
* Create profile
    * POST https://aws-endpoint/v1/cruddyAPI/profiles 
* Read / retrieve profile
    * GET https://aws-endpoint/v1/cruddyAPI/profiles/{id} 
* Update profile
    * PUT https://aws-endpoint/v1/cruddyAPI/profiles/{id}
* Delete profile
    * DELETE https://aws-endpoint/v1/cruddyAPI/profiles/{id}

## Development
### Requirements
* [Golang](https://golang.org/dl/) >= 1.11
* [Terraform](https://www.terraform.io/downloads.html) >= 0.11.11
* [AWS](https://aws.amazon.com/) Credentials
    * Infrastructure costs fit comfortably within AWS' Free-Tier

```bash
git clone https://github.com/teohrt/cruddyAPI.git
cd cruddyAPI
make deploy
```

## Licence
See LICENSE.
