# Configuration

- [Service configuration](#service-configuration)
  - [Default superadmin domain name](#default-superadmin-domain-name)
  - [Default superadmin role name](#default-superadmin-role-name)
  - [Database](#database)
  - [Dictionary](#dictionary)
  - [Watcher](#watcher)
- [Dictionary configuration](#dictionary-configuration)
  - [Creator](#creator)
  - [Feature](#feature)

## Service configuration

Configuration for creating a new caskin service instance.

#### Example:

```json
{
  "default_superadmin_domain_name": "",
  "default_superadmin_role_name": "",
  "db": {},
  "dictionary": {},
  "watcher": {},
}
```

### Default superadmin domain name

`default_superadmin_domain_name` **(optional)** is a string for setting superadmin domain name. The default value is `"superadmin_domain"`

#### Example:

```json
{
  "default_superadmin_domain_name": "超级管理域"
}
```

### Default superadmin role name

`default_superadmin_role_name` **(optional)** is a string for setting superadmin role name. The default value is `"superadmin_role"`

#### Example:

```json
{
  "default_superadmin_role_name": "超级管理员"
}
```

### Database

`db` **(required)** is a structure for defining **caskin** and [caskin](https://github.com/casbin/casbin) metadata database client.

1. `dsn` is the path to database.
2. `type` is the type of database.

#### Supported database client type:

| type         | [dsn relate package version](../go.mod) |
|--------------|-----------------------------------------|
| `sqlite`     | gorm.io/driver/sqlite v1.3.6            |
| `mysql`      | gorm.io/driver/mysql v1.3.6             |
| `postgres`   | gorm.io/driver/postgres v1.3.6          |

#### Example:

```json
{
  "db": {
    "dsn": "./sqlite.db",
    "type": "sqlite"
  }
}
```

### Dictionary

`dictionary` **(required)** is the option for loading [Dictionary configuration](#dictionary-configuration).

1. `dsn` is the path to load dictionary configuration.
2. `type` is dictionary adaptor type.

#### Supported dictionary option adaptor type:

| adaptor type | description                             |
|--------------|-----------------------------------------|
| ` `          | default type is `FILE`                  |
| `FILE`       | load dictionary configuration from file |

#### Example:

```json
{
  "dictionary": {
    "dsn": "caskin.toml",
    "type": "FILE"
  }
}
```

### Watcher

`watcher` **(optional)** is an option for setting to keep consistence between multiple [caskin](https://github.com/casbin/casbin) enforcer instances.

1. `type` **(required)** is the type of watcher.
2. `address` **(optional)** is the address of watcher.
3. `password` **(optional)** is the password of watcher.
4. `channel` **(optional)** is the channel key of watcher.
5. `auto_load` **(optional)** is an autoload time interval setting when `type` is default.

#### Supported watcher type:

| type    | description                 |
|---------|-----------------------------|
| ` `     | default type is no watcher. |
| `redis` | redis watcher.              |

#### Example:

```json
{
  "watcher": {
    "type": "redis",
    "address": "localhost:6379",
    "password": "",
    "channel": "caskin",
    "auto_load": 0 
  }
}
```

## Dictionary configuration

### Creator

### Feature
