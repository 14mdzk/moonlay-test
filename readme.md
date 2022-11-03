
# Moonlay Technology Practical Test
An assessment needed to take as a Backend Engineer


## API Reference

### **List**
#### Get all list except sub list

```http
  GET /list
```

#### Get selected list with related sub list if exist

```http
  GET /list/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `integer` | **Required** |

#### Create A List

```http
  POST /list
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
|**Body**| ||
| `title`      | `string` | **Required**, **Max: 100** |
| `description` | `string` | **Required**, **Max: 1000**|
| `file` | `file` | **Optional**, **Type:pdf,txt**|

#### Update A List
```http
  PUT /list/${id}
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
|`id`| `integer`| **Required**|
|**Body**| ||
| `title`      | `string` | **Required**, **Max: 100** |
| `description` | `string` | **Required**, **Max: 1000**|
| `file` | `file` | **Optional**, **Type:pdf,txt**|

#### Delete A List
```http
  Delete /list/${id}
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
|`id`| `integer`| **Required**|


### **Sub List**
#### Get All Sub List related to given id

```http
  GET /list/${id}/sub_list
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `integer` | **Required** |

#### Get Sub List related to given id and sub_id

```http
  GET /list/${id}/sub_list/${sub_id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `integer` | **Required** |
| `sub_id`      | `integer` | **Required** |

#### Create A List

```http
  POST /list/${id}/sub_list
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
| `id`      | `integer` | **Required** |
|**Body**| ||
| `title`      | `string` | **Required**, **Max: 100** |
| `description` | `string` | **Required**, **Max: 1000**|
| `file` | `file` | **Optional**, **Type:pdf,txt**|

#### Update A List
```http
  PUT /list/${id}/sub_list/${sub_id}
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
|`id`| `integer`| **Required**|
| `sub_id`      | `integer` | **Required** |
|**Body**| ||
| `title`      | `string` | **Required**, **Max: 100** |
| `description` | `string` | **Required**, **Max: 1000**|
| `file` | `file` | **Optional**, **Type:pdf,txt**|

#### Delete A List
```http
  Delete /list/${id}/sub_list/${sub_id}
```
| Parameter | Type     | Description                       |
| -------- | :------- | :-------------------------------- |
|`id`| `integer`| **Required**|
|`sub_id`| `integer`| **Required**|