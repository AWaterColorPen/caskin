# caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)


## TODO List

- feature
  - [ ] (p2) web feature, frontend menu and sub function, backend API, RestFUL API
  - [ ] (p2) abstract features for just like web feature and other feature
  - [ ] (p3) go 1.16
  - [ ] (p3) if object type == object, it will not recover it by admin user now. we want to support it by a special API
  - [ ] (p1) fen kai role and object's logic code
  - [ ] (p1) can't update or delete or recover or create root object exclude superadmin
- bug
  - [ ] (p2) fix issue when modify object data, it will not just filter. it should return no permission
  

- [ ] role的测试
- [ ] 综合多种情况进行问题位置的判断
- [ ] object_data的测试()
- [ ] backend 数据表 每一条数据独占一个objet数据组，所以需要联动修改 ()

- ObjectData是真实的数据，而不是数据组
- Object是数据组的概念
