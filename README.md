# Caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)

Caskin is a multi-domain rbac library for Golang projects. It develops base on [caskin](https://github.com/casbin/casbin) 

## Design Goal

### 1. User: the person - 真实用户

### 2. Role: the group of user - 角色=一组用户

### 3. Object: the resource or resource group of authorization node  - 资源=权限节点/权限组

### 4. Domain: organization - 域=组织

| User - 用户 | Role - 角色 | Object -  | Domain - Organization |
| --- | --- | --- | --- |
| Header | Title | Header | school-1 |
| Paragraph | Text | Header | school-2 |

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

1. 
- docs
  - overall design
  - get start

```
project
│   README.md
│   file001.txt    
│
└───folder1
│   │   file011.txt
│   │   file012.txt
│   │
│   └───subfolder1
│       │   file111.txt
│       │   file112.txt
│       │   ...
│   
└───folder2
    │   file021.txt
    │   file022.txt
```