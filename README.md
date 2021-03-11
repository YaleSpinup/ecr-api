# ecr-api

This API provides simple restful API access to a service.

## Endpoints

```
GET /v1/ecr/ping
GET /v1/ecr/version
GET /v1/ecr/metrics

GET    /v1/ecr/{account}/repositories
POST   /v1/ecr/{account}/repositories/{group}
GET    /v1/ecr/{account}/repositories/{group}
GET    /v1/ecr/{account}/repositories/{group}/{name}
PUT    /v1/ecr/{account}/repositories/{group}/{name}
DELETE /v1/ecr/{account}/repositories/{group}/{name}

GET    /v1/ecr/{account}/repositories/{group}/{name}/images
GET    /v1/ecr/{account}/repositories/{group}/{name}/images/{tag}

GET    /v1/ecr/{account}/repositories/{group}/{name}/users
POST   /v1/ecr/{account}/repositories/{group}/{name}/users
GET    /v1/ecr/{account}/repositories/{group}/{name}/users/{user}
PUT    /v1/ecr/{account}/repositories/{group}/{name}/users/{user}
DELETE /v1/ecr/{account}/repositories/{group}/{name}/users/{user}
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed via the `X-Auth-Token` header.

## Usage

### Repositories

#### Create a repository

POST `/v1/ecr/{account}/repositories/{group}`

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | create a repository             |
| **400 Bad Request**           | badly formed request            |
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

##### Example create request body

```json
{
    "RepositoryName": "myAwesomeRepository",
    "Groups": ["spindev-000001", "spindev-000002"],
    "ScanOnPush": "true",
    "Tags": [
        {
            "Key": "CreatedBy",
            "Value": "cf322"
        },
        {
            "Key": "spinup:spaceid",
            "Value": "spincool-00001"
        }
    ]
}
```

##### Example create response body

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
    "Groups": ["spindev-000001", "spindev-000002"],
    "KmsKeyId": "",
    "ScanOnPush": "true",
    "ImageTagMutability": "MUTABLE",
    "RegistryId": "0123456789",
    "RepositoryArn": "arn:aws:ecr:us-east-1:0123456789:repository/spindev-00001/myAwesomeRepository",
    "RepositoryName": "spindev-00001/camdenstestrepo02",
    "RepositoryUri": "0123456789.dkr.ecr.us-east-1.amazonaws.com/spindev-00001/myAwesomeRepository",
    "Tags": [
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        },
        {
            "Key": "CreatedBy",
            "Value": "santa"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "Name",
            "Value": "spindev-00001/myAwesomeRepository"
        }
    ]
}
```

#### List Repositories

GET `/v1/ecr/{account}/repositories`

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | return the list of repositories |
| **400 Bad Request**           | badly formed request            |
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

#### Example list response

```json
[
    "spindev-00001/dasher",
    "spindev-00001/dancer",
    "spindev-00002/prancer",
    "spindev-00002/vixen"
    "spindev-00003/comet",
    "spindev-00003/cupid",
    "spindev-00004/donner",
    "spindev-00005/blitzen",
    "spindev-00006/rudolph",
]
```

#### List Repositories by group id

GET `/v1/ecr/{account}/repositories/{group}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return the list of repositories  |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account not found                |
| **500 Internal Server Error** | a server error occurred          |

##### Example list by group response

```json
[
    "spindev-00006/rudolph"
]
```

#### Get details about a Repository

GET `/v1/ecr/{account}/repositories/{group}/{id}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return details of a repository   |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or repository not found  |
| **500 Internal Server Error** | a server error occurred          |

##### Example show response

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
    "Groups": ["spindev-000001", "spindev-000002"],
    "KmsKeyId": "",
    "ScanOnPush": "true",
    "ImageTagMutability": "MUTABLE",
    "RegistryId": "0123456789",
    "RepositoryArn": "arn:aws:ecr:us-east-1:0123456789:repository/spindev-00001/myAwesomeRepository",
    "RepositoryName": "spindev-00001/camdenstestrepo02",
    "RepositoryUri": "0123456789.dkr.ecr.us-east-1.amazonaws.com/spindev-00001/myAwesomeRepository",
    "Tags": [
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        },
        {
            "Key": "CreatedBy",
            "Value": "santa"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "Name",
            "Value": "spindev-00001/myAwesomeRepository"
        }
    ]
}
```

#### Update a repository

PUT `/1/ecr/{account}/repositories/{group}/{id}`

##### Example update request body

```json
{
    "ScanOnPush": "false",
    "Groups": ["spindev-000001", "spindev-000002", "spindev-000003"],
    "Tags": [
        {
            "Key": "Application",
            "Value": "MyAwesomeCloudApp"
        }
    ]
}
```

##### Example update response body

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
    "Groups": ["spindev-000001", "spindev-000002", "spindev-000003"],
    "KmsKeyId": "",
    "ScanOnPush": "false",
    "ImageTagMutability": "MUTABLE",
    "RegistryId": "0123456789",
    "RepositoryArn": "arn:aws:ecr:us-east-1:0123456789:repository/spindev-00001/myAwesomeRepository",
    "RepositoryName": "spindev-00001/camdenstestrepo02",
    "RepositoryUri": "0123456789.dkr.ecr.us-east-1.amazonaws.com/spindev-00001/myAwesomeRepository",
    "Tags": [
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        },
        {
            "Key": "CreatedBy",
            "Value": "santa"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "Name",
            "Value": "spindev-00001/myAwesomeRepository"
        },
        {
            "Key": "Application",
            "Value": "MyAwesomeCloudApp"
        }
    ]
}
```

#### Delete a repository and all images

*NOTE:* deleting a repository does not currently cleanup users, this must be done first.

DELETE `/v1/ecr/{account}/repositories/{group}/{id}`

| Response Code                 | Definition                               |
| ----------------------------- | -----------------------------------------|
| **200 Submitted**             | delete request is submitted              |
| **400 Bad Request**           | badly formed request                     |
| **403 Forbidden**             | bad token or fail to assume role         |
| **404 Not Found**             | account or repository not found          |
| **409 Conflict**              | repository is not in the available state |
| **500 Internal Server Error** | a server error occurred                  |

### Images

#### List images in a repository

GET `/v1/ecr/{account}/repositories/{group}/{id}/images`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return list of images            |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or repository not found  |
| **500 Internal Server Error** | a server error occurred          |

##### Example response body

```json
[
    {
        "ArtifactMediaType": "application/vnd.docker.container.image.v1+json",
        "ImageDigest": "sha256:ac81321d3627bcde149b383220b16dabc590f2d247f4c72c64cb14f58e7fb9c2",
        "ImageManifestMediaType": "application/vnd.docker.distribution.manifest.v2+json",
        "ImagePushedAt": "2020-12-14T15:56:11Z",
        "ImageScanFindingsSummary": {
            "FindingSeverityCounts": {},
            "ImageScanCompletedAt": "2020-12-14T16:02:36Z",
            "VulnerabilitySourceUpdatedAt": "2020-11-04T01:21:09Z"
        },
        "ImageScanStatus": {
            "Description": "The scan was completed successfully.",
            "Status": "COMPLETE"
        },
        "ImageSizeInBytes": 16093514,
        "ImageTags": [
            "latest"
        ],
        "RegistryId": "0123456789",
        "RepositoryName": "spindev-00001/myAwesomeRepository"
    }
]
```

#### Get details about an image tag

This gets the image scanning results for an image tag.

GET `/v1/ecr/{account}/repositories/{group}/{id}/images/{tag}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return scanning results          |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or repository not found  |
| **500 Internal Server Error** | a server error occurred          |

##### Example response body

```json
{
    "FindingSeverityCounts": {
        "HIGH": 1,
        "INFORMATIONAL": 61,
        "LOW": 17,
        "MEDIUM": 22,
        "UNDEFINED": 5
    },
    "Findings": [
        {
            "Attributes": [
                {
                    "Key": "package_version",
                    "Value": "2.28-10"
                },
                {
                    "Key": "package_name",
                    "Value": "glibc"
                },
                {
                    "Key": "CVSS2_VECTOR",
                    "Value": "AV:N/AC:M/Au:N/C:N/I:N/A:C"
                },
                {
                    "Key": "CVSS2_SCORE",
                    "Value": "7.1"
                }
            ],
            "Description": "The iconv feature in the GNU C Library (aka glibc or libc6) through 2.32, when processing invalid multi-byte input sequences in the EUC-KR encoding, may have a buffer over-read.",
            "Name": "CVE-2019-25013",
            "Severity": "HIGH",
            "Uri": "https://security-tracker.debian.org/tracker/CVE-2019-25013"
        },
        {
            "Attributes": [
                {
                    "Key": "package_version",
                    "Value": "7.64.0-4+deb10u1"
                },
                {
                    "Key": "package_name",
                    "Value": "curl"
                },
                {
                    "Key": "CVSS2_VECTOR",
                    "Value": "AV:N/AC:L/Au:N/C:N/I:P/A:N"
                },
                {
                    "Key": "CVSS2_SCORE",
                    "Value": "5"
                }
            ],
            "Description": "curl 7.41.0 through 7.73.0 is vulnerable to an improper check for certificate revocation due to insufficient verification of the OCSP response.",
            "Name": "CVE-2020-8286",
            "Severity": "MEDIUM",
            "Uri": "https://security-tracker.debian.org/tracker/CVE-2020-8286"
        },
        ...
    ],
    "ImageScanCompletedAt": "2021-03-11T17:27:30Z",
    "VulnerabilitySourceUpdatedAt": "2021-03-11T07:49:36Z"
}
```

### Users

Repository users are created in the same account as the repository.  An account is "bootstrapped" by
the create action if its not already prepared to contain repository users.  This bootstrapping creates
a shared role and group for admin users.  The role grants a user access to a repository based on the tags
`spinup:org`, `spinup:spaceid` and `ResourceName` (the repository name).

#### List all users for a repository

GET    /v1/ecr/{account}/repositories/{group}/{name}/users

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return the list of users         |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account not found                |
| **500 Internal Server Error** | a server error occurred          |

##### Example list users response

```json
[
    "user1"
]
```

#### Create a user

POST   /v1/ecr/{account}/repositories/{group}/{name}/users

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | create a repository             |
| **400 Bad Request**           | badly formed request            |
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

##### Example create user request body

```json
{
    "username": "user1",
    "tags": [
        {
            "key": "application",
            "value": "myapp"
        }
    ],
    "groups": [
        "SpinupECRAdminGroup"
    ]
}
```

##### Example create user response body

```json
{
    "UserName": "user1",
    "AccessKeys": [],
    "Groups": [
        "SpinupECRAdminGroup"
    ],
    "Tags": [
        {
            "Key": "application",
            "Value": "myapp"
        },
        {
            "Key": "ResourceName",
            "Value": "spindev-00001-myAwesomeRepository-user1"
        },
        {
            "Key": "ResourceName",
            "Value": "spindev-00001/myAwesomeRepository"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        }
    ]
}
```

#### Get details about a user

GET    /v1/ecr/{account}/repositories/{group}/{name}/users/{user}

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return details of a user         |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or user not found        |
| **500 Internal Server Error** | a server error occurred          |

##### Example show user response

```json
{
    "UserName": "user1",
    "AccessKeys": [],
    "Groups": [
        "SpinupECRAdminGroup"
    ],
    "Tags": [
        {
            "Key": "application",
            "Value": "myapps"
        },
        {
            "Key": "ResourceName",
            "Value": "spindev-00001-myAwesomeRepository-user1"
        },
        {
            "Key": "Name",
            "Value": "spindev-00001/myAwesomeRepository"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        }
    ]
}
```

#### Update a user

A user's tags and/or its access key can be updated.  The operations are independent,
and can occur in the same request, or individually.  If the access key is reset,
a new access key will be returned with the response.

PUT /v1/ecr/{account}/repositories/{group}/{name}/users/{user}

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | updated the user                |
| **400 Bad Request**           | badly formed request            |
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

##### Example update user request body

```json
{
    "resetkey": true,
    "tags": [
        {
            "key": "application",
            "value": "myapp123"
        }
    ],
}
```

##### Example update user response

```json
{
    "UserName": "user1",
    "AccessKey": {
        "AccessKeyId": "AAAAABBBBBCCCCCDDDDDEEEEEFFFFF",
        "CreateDate": "2021-02-03T22:37:30Z",
        "SecretAccessKey": "gxyz1234567890abcdefghijklmnop",
        "Status": "Active",
        "UserName": "spincool-00001-testrepo1-user1"
    },
    "DeletedAccessKeys": [
        "QQQQQRRRRRSSSSSTTTTTUUUUVVVV"
    ],
    "Tags": [
        {
            "Key": "application",
            "Value": "myapp123"
        },
        {
            "Key": "ResourceName",
            "Value": "spindev-00001-myAwesomeRepository-user1"
        },
        {
            "Key": "Name",
            "Value": "spindev-00001/myAwesomeRepository"
        },
        {
            "Key": "spinup:org",
            "Value": "spindev"
        },
        {
            "Key": "spinup:spaceid",
            "Value": "spindev-00001"
        }
    ]
}
```

#### Delete a user

DELETE /v1/ecr/{account}/repositories/{group}/{name}/users/{user}

| Response Code                 | Definition                               |
| ----------------------------- | -----------------------------------------|
| **200 Submitted**             | delete request is submitted              |
| **400 Bad Request**           | badly formed request                     |
| **403 Forbidden**             | bad token or fail to assume role         |
| **404 Not Found**             | account or user not found                |
| **409 Conflict**              | user is not in the available state       |
| **500 Internal Server Error** | a server error occurred                  |

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2020 Yale University
