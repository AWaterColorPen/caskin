# caskin 现代化迭代设计文档

**日期：** 2026-03-10  
**目标：** 分阶段对 caskin 进行现代化维护和文档建设，保持向后兼容，不让存量用户用不了

---

## 背景

caskin 是一个 Go 的多域 RBAC 权限管理库，基于 casbin 开发。最后一次提交是 2023 年 5 月，处于维护停滞状态。依赖版本较旧（Go 1.20），文档不完善，示例缺乏。

## 目标

1. **现代化维护** — Go 版本 + 依赖升级，使用新语法，分阶段进行避免存量用户不兼容
2. **文档建设** — 使用者文档（快速上手、API 说明）+ 贡献者文档（架构设计）

## 原则

- 分阶段渐进，**每个 Phase 切换前等用户确认**
- 每次有意义的改动单独 PR + commit，保持可回滚
- Breaking change 单独 PR，附 CHANGELOG 说明
- 语法现代化和依赖升级分开 PR

---

## Phase 1：地基（2-3 周）

**目标：让项目跑起来、看得懂**

### 任务清单

- [x] 升级 Go 到最新版（1.24）
- [x] 依赖升级（只升无破坏性变更的小版本/patch）
- [x] 完善 README：补充真实可运行的 Quick Start 示例（PR #26，已合并）
- [x] 使用新 Go 语法重构（PR #27）：
  - 完全移除 `go-linq` 和 `golang.org/x/exp`，替换为 stdlib
  - `sort.Slice` → `slices.SortFunc`
  - `constraints.Ordered` → `cmp.Ordered`
  - BFS 循环加注释增强可读性
  - `distinct()` 泛型函数复用
- [x] 跑通所有现有测试，修复失败的（所有测试通过）
- [x] 补充 godoc 注释（核心 API）（PR #28，已合并 2026-04-14）

**✅ Phase 1 全部完成！等用户确认进入 Phase 2**

---

## Phase 2：文档建设（2-3 周）

**目标：让新用户能快速上手，让贡献者能参与**

### 任务清单

- [x] 重写 Getting Started（step-by-step，有完整可运行代码）（PR #29）
- [x] API 文档：每个方法说明参数、返回值、使用场景（PR #30，docs/api-reference.md）
- [x] 增加常见使用场景示例（PR #31，docs/use-cases.md）：
  - 多域管理
  - 角色继承
  - 权限检查
  - 前端/后端权限分离
  - ⚠️ Review 修复（2026-04-22）：4 个编译阻断 bug 已修复（见下方经验积累）
  - ⚠️ Review Round 2（2026-04-23）：2 个新编译阻断 bug 已修复（见下方经验积累）
  - ✅ Review Round 3（2026-04-24）：Round 2 两个 bug 已修复 — `CheckObject` 全部替换为 `GetObject` 过滤模式，`GetDomain(superadmin)` 改为 `GetDomain()`
- [x] 架构说明文档（给贡献者看，docs/architecture.md）
- [x] CONTRIBUTING.md

**✅ Phase 2 全部完成！等用户确认进入 Phase 3**

---

## Review 经验积累

### 文档示例代码的常见陷阱（2026-04-22，PR #31）

1. **`caskin.New` 无 functional option 变体**：只接受 `*Options` struct，文档示例不要使用不存在的 `WithDB`/`WithDictionary` 等封装函数。
2. **`IService` 不暴露 `GetEnforcer()`**：需要低级别 `caskin.Check` 时，示例应使用 `svc.CheckObject()` 替代，或说明需要持有 concrete `*server` 类型。
3. **`caskin.Object` 接口方法名**：是 `GetObjectType()`，不是 `GetType()`，写示例前需对照 `schema.go` 确认接口定义。
4. **`caskin.Register[...]()` 必须在 `New` 前调用**：每个新的示例/测试 setup 函数都要检查是否有这行，否则运行时 panic。

### Review 经验积累

**caskin API 文档示例注意事项（2026-04-22 from PR #31）：**

1. **必须调用 `caskin.Register[U,R,O,D]()`** — 在 `caskin.New` 之前调用，否则 factory 不知道具体类型，运行时 panic
2. **`caskin.New` 接受 `*Options` 结构体** — 无 `WithDB` / `WithDictionary` / `DefaultModelText()` 等 functional options；正确写法：`caskin.New(&caskin.Options{DB: dbOption, Dictionary: &caskin.DictionaryOption{Dsn: "configs/caskin.toml"}})`
3. **`IService` 不暴露 `GetEnforcer()`** — 使用 `svc.CheckObject(user, domain, obj, action) == nil` 替代 `caskin.Check(svc.GetEnforcer(), ...)`；`caskin.Check` 需要具体的 enforcer，不能通过 `IService` 调用
4. **`caskin.Object` 接口方法是 `GetObjectType()`** — 不是 `GetType()`

**caskin API 文档示例追加注意事项（2026-04-23 from PR #31 Round 2）：**

5. **`IService` 没有 `CheckObject(user, domain, obj, action)` 方法** — `CheckObject` 只在 `*server` struct 上，**不在 `IService` 接口里**。所有 `svc.CheckObject(...)` 调用都会编译失败。对 `ObjectData` 类型可用 `svc.CheckModifyObjectData(user, domain, objData)` 等方法；对纯 `Object` 的权限检查需确认正确的公开 API。
6. **`IService.GetDomain()` 无参数** — 接口签名是 `GetDomain() ([]Domain, error)`，不接受任何参数。要列出特定用户所在的域，使用 `GetDomainByUser(user User) ([]Domain, error)`。

**caskin API 文档示例追加注意事项（2026-04-24 from PR #31 Round 3）：**

7. **`IService` 无 `CheckObject` — 使用 `GetObject` 查询模式替代** — caskin 设计为查询导向：`GetObject(user, domain, action)` 只返回调用者有权限访问的对象。判断权限的正确方式是：取得列表后检查目标 object 是否在列表中（`containsObj`）。`CheckObject` 存在于 `*server` struct 上，但不在 `IService` 接口中暴露。

---

## Phase 3：现代化深化（3-4 周）✅ 已完成（2026-05-04 结案）

**目标：主要依赖全面更新，代码质量提升**

### 任务清单

- [x] 升级 casbin/v2 v2.69 → v2.135（commit 210e8b3）
- [x] 升级 gorm v1.31.1、gorm-drivers v1.6.0、redis-watcher v2.8.0、go-redis v9.18.0（commit f41873d）
- [x] 提升测试覆盖率至 82.8%（主包），超 80% 目标（commit f41873d）
- [x] 消除 golang-jwt/v4 CVE 安全漏洞（go-mssqldb v0.17→v1.9.5，commit ef1b692）

**✅ Phase 3 全部完成！已打 v0.3.0 tag**

---

## Phase 4：评估（2026-05-06）

**状态：等待用户决策**

### 可选方向

#### 候选 A：迁移 casbin/v3（较大重构，非紧迫）
- **前提**：gorm-adapter/v3 v3.41+ 已迁移至 casbin/v3，若要升级 gorm-adapter 必须同步迁移
- **估计工作量**：中等（需全面替换 casbin/v2 API 调用、gorm-adapter 联动）
- **优先级建议**：低，当前 casbin/v2 v2.135 安全且功能完整
- **建议**：除非 casbin/v2 被宣布 EOL 或出现新 CVE，否则不急于迁移

#### 候选 B：降频维护（当前推荐）
- Phase 1-3 成果：Go 1.24、全面依赖升级、文档建设、测试覆盖 82.8%、CVE 清零
- 项目已达到健康可用状态
- **建议**：进入 "安全守护" 节奏——只在出现新安全漏洞或重大 Go 版本变化时响应

**卡点：需用户确认 Phase 4 方向

---

## 推进节奏

- 每 **2-3 天** 自动推进一个小任务
- 完成有意义节点后在对话里发进度报告
- Phase 切换必须等用户确认
- 所有改动都 commit，保持可回滚
