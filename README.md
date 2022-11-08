# Caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)

Caskin is a multi-domain rbac library for Golang projects. It develops base on [caskin](https://github.com/casbin/casbin) 

## Introduction

### Example

## Documentation

1. [Configuration](./docs/configuration.md) to configure caskin instance and dictionary.

## Design Goal

### 1. User: the person - 真实用户

real person for school business. for example:
- student
- parent
- teacher
  - class teacher
  - subject teacher
- property management company
  - dorm supervisor
  - teaching building manager
  - guard

### 2. Role: the group of user - 角色=一组用户

role group for school business. for example:

- student
  - junior student
    - junior class 1
  - senior student 
- parent
- teacher
  - class teacher
  - subject teacher
- property management company
  - dorm supervisor
  - teaching building manager
  - guard

### 3. Object: the resource or resource group of authorization node  - 资源=权限节点/权限组

#### feature resource

authorization action = read

- school gateway
- playground
- dining room

#### data resource

authorization action = read / write / manage

- course management
  - class 1 course management
  - class 2 course management
- student management
  - class 1 student management
  - class 2 student management
- teacher management
- dorm management
  - room 1
  - idle room
- teaching building management
  - junior building
    - classroom 1
  - senior building
  - administration building
  - idle building

### 4. Domain: organization - 域=组织

every student / parent / teacher / property management company can be working in one or more schools

| Domain   |
|----------|
| school-1 |
| school-2 |

## Feature

1. Caskin focus on management of authorization business.
   - casbin is easy to access control. but it is bad at access controlling on authorization management.
2. Abstracting `User`, `Role`, `Object`, `Domain` as interface, it is easy for business logic code to use.
3. Facade pattern, provide management APIs for overall authorization business
   - management [user](https://github.com/AWaterColorPen/caskin/blob/main/executor_user.go)
   - management relation of [user and domain](https://github.com/AWaterColorPen/caskin/blob/main/executor_user_domain.go)
   - management relation of [user and role](https://github.com/AWaterColorPen/caskin/blob/main/executor_user_role.go)
   - management [superadmin](https://github.com/AWaterColorPen/caskin/blob/main/executor_superadmin.go)
   - management [role](https://github.com/AWaterColorPen/caskin/blob/main/executor_role.go)
   - management [object](https://github.com/AWaterColorPen/caskin/blob/main/executor_object.go)
   - management resource group [object data](https://github.com/AWaterColorPen/caskin/blob/main/executor_object_data.go)

## Getting Started

## License

See the [License File](./LICENSE).