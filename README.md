# caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)


## TODO List

- feature
  - [x] redesign `DomainCreator` to help `CreateDomain` `RecoverDomain` to be reentrant API
  - [ ] (p2) web feature, frontend menu and sub function, backend API
  - [ ] (p2) abstract features for just like web feature and other feature
  - [ ] (p3) go 1.16
- bug
  - [x] (p1) fix issue when remove `Role`, it should remove children role's g and parent role's g
  - [ ] (p1) fix issue when modify `Role` and `Object`, it should check old item's Parent's write permission
  - [ ] (p1) fix issue when modify object data, it should check relate `Object.GetObjectType`
  - [ ] (p1) fix issue when remove `Role` or `Object`, it should remove all its son node
  - [x] (p1) fix issue when modify policies, it uses []byte as map's key
  - [x] (p1) fix issue when modify policies, it does not successfully update the policies
  - [x] (p1) different API for `CreateDomain` and `ReInitializeDomain`
- unit test
  - [x] unit test get `casbin.Model` with cache. it should not create new one per unit test
  - [x] create domain in unit test
  - [x] add users: superadmin, domain-admin, domain-member
  - [x] wrap a `Stage` with `Caskin` instance and initialized `User`, `Domain`, `Role`, `Object`
  - [ ] (p0) unit for executor domain API
  - [x] (p0) unit for executor user API
  - [ ] (p0) unit for executor user role API
  - [x] (p0) unit for executor superadmin API
  - [ ] (p0) unit for executor role API
  - [ ] (p0) unit for executor object API
  - [ ] (p0) unit for executor policy API
  