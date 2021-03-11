# caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)


## TODO List

- feature
  - [ ] (p2) web feature, frontend menu and sub function, backend API
  - [ ] (p2) abstract features for just like web feature and other feature
  - [ ] (p3) go 1.16
- bug
  - [ ] (p1) fix `Role` and `Object` 's `Create`, `Recover`, `Update` API, it should update parent's `g/g2` casbin policy
  - [x] (p1) fix issue when modify `Role` and `Object`, it should check old item's Parent's write permission
  - [ ] (p1) fix issue when modify object data, it should check relate `Object.GetObjectType`
  - [x] (p1) fix issue when remove `Role` or `Object`, it should remove all its son node
  - [x] (p1) fix issue when modify policies, it does not successfully update the policies
- unit test
  - [x] (p0) unit for executor domain API
  - [x] (p0) unit for executor user API
  - [ ] (p0) unit for executor user role API
  - [x] (p0) unit for executor superadmin API
  - [ ] (p0) unit for executor role API
  - [ ] (p0) unit for executor object API
  - [x] (p0) unit for executor policy API
  