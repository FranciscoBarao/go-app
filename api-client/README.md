# API client

Bruno collections for local HTTP against the running services. Each subdirectory is its own collection. 

Install [Bruno](https://www.usebruno.com). Open a collection folder (the directory that contains `bruno.json`), not the `api-client/` parent.

| Collection      | Open                         | Compose port |
| --------------- | ---------------------------- | ------------ |
| Catalog         | `api-client/catalog`         | 8081         |
| User Management | `api-client/user-management` | 8080         |

Select the `local` environment in each collection. `token` is empty in git and listed as a secret. Do not commit a filled value.

## Authorization

User Management issues the OAuth token. Catalog does not share that collection's environment.

1. Run **Register** (optional) then **Login**, or **Client credentials**. Login and client-credentials scripts store `access_token` in User Management env `token`.
2. Copy `token` into Catalog's `local` environment.
3. On a Catalog write (or Rate), enable the disabled `Authorization: Bearer {{token}}` header, or switch that request's auth from none to inherit (collection Bearer).

Public GETs send no Authorization header.


## QUERY

List resources accept HTTP `QUERY` with a JSON body and return the same envelope as GET. Those requests use Bruno's custom-method `http` block. If the installed Bruno version cannot send `QUERY`, use the matching **List** GET request instead.
