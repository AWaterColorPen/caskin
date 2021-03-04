# caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)

> 

## TODO List

- feature
  - [ ] redesign `DomainCreator` to help `CreateDomain` `RecoverDomain` to be reentrant API
- bug
  - [ ] fix issue when modify `Role` and `Object`, it should check old item's Parent's write permission
  - [ ] fix issue when modify object data, it should check relate `Object.GetObjectType`
- unit test
  - [ ] unit test get `casbin.Model` with cache. it should not create new one per unit test
