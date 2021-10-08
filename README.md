# Caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)

Caskin is a multi-domain rbac library for Golang projects. It develops base on [caskin](https://github.com/casbin/casbin) 

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
