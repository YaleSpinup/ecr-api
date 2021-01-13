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
```

## Authentication

Authentication is accomplished via an encrypted pre-shared key passed via the `X-Auth-Token` header.

## Usage

### Create a repository

POST `/v1/ecr/{account}/repositories/{group}`

| Response Code                 | Definition                      |
| ----------------------------- | --------------------------------|
| **200 OK**                    | create a repository             |
| **400 Bad Request**           | badly formed request            |
| **404 Not Found**             | account not found               |
| **500 Internal Server Error** | a server error occurred         |

#### Example create request body

```json
{
    "RepositoryName": "myAwesomeRepository",
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

#### Example create response body

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
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

### List Repositories

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

### List Repositories by group id

GET `/v1/ecr/{account}/repositories/{group}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return the list of repositories  |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account not found                |
| **500 Internal Server Error** | a server error occurred          |

#### Example list by group response

```json
[
    "spindev-00006/rudolph"
]
```

### Get details about a Repository

GET `/v1/ecr/{account}/repositories/{group}/{id}`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return details of a repository   |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or repository not found  |
| **500 Internal Server Error** | a server error occurred          |

#### Example show response

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
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

### Update a repository

PUT `/1/ecr/{account}/repositories/{group}/{id}`

#### Example update request body

```json
{
    "ScanOnPush": "false",
    "Tags": [
        {
            "Key": "Application",
            "Value": "MyAwesomeCloudApp"
        }
    ]
}
```

#### Example update response body

```json
{
    "CreatedAt": "2020-12-14T15:34:18Z",
    "EncryptionType": "AES256",
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

### Delete a repository and all images

DELETE `/v1/ecr/{account}/repositories/{group}/{id}`

| Response Code                 | Definition                               |
| ----------------------------- | -----------------------------------------|
| **200 Submitted**             | delete request is submitted              |
| **400 Bad Request**           | badly formed request                     |
| **403 Forbidden**             | bad token or fail to assume role         |
| **404 Not Found**             | account or repository not found          |
| **409 Conflict**              | repository is not in the available state |
| **500 Internal Server Error** | a server error occurred                  |

### List images in a repository

GET `/v1/ecr/{account}/repositories/{group}/{id}/images`

| Response Code                 | Definition                       |
| ----------------------------- | ---------------------------------|
| **200 OK**                    | return details of a repository   |
| **400 Bad Request**           | badly formed request             |
| **403 Forbidden**             | bad token or fail to assume role |
| **404 Not Found**             | account or repository not found  |
| **500 Internal Server Error** | a server error occurred          |

#### Example response body

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

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright © 2020 Yale University
