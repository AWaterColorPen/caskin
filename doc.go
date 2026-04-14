// Package caskin provides a multi-domain RBAC (Role-Based Access Control)
// permission management library built on top of casbin.
//
// # Overview
//
// caskin extends casbin's RBAC model to support multiple isolated domains,
// where each domain has its own set of roles, objects, and policies. This
// allows you to build multi-tenant applications where each tenant (domain)
// has fully isolated permissions.
//
// # Core Concepts
//
//   - User: An entity that can be assigned roles within a domain.
//   - Domain: An isolated permission scope (e.g., a tenant or organization).
//   - Role: A named set of permissions within a domain; roles can be hierarchical.
//   - Object: A resource that can be accessed; objects form a tree hierarchy.
//   - ObjectData: A domain-specific data item associated with an Object type.
//   - Policy: A tuple of (Role, Object, Domain, Action) granting a role
//     permission to perform an action on an object within a domain.
//   - Action: One of read, write, or manage.
//
// # Quick Start
//
//	// 1. Define your own User, Role, Object, Domain types (implement the
//	//    corresponding interfaces), then register them with caskin:
//	caskin.Register[*MyUser, *MyRole, *MyObject, *MyDomain]()
//
//	// 2. Create a service instance:
//	svc, err := caskin.New(&caskin.Options{
//	    DB: &caskin.DBOption{DSN: "..."},
//	})
//
//	// 3. Use the service to manage permissions:
//	svc.CreateDomain(domain)
//	svc.CreateRole(admin, domain, role)
//	svc.AddUserRole(admin, domain, []*caskin.UserRolePair{{User: user, Role: role}})
//
// # Superadmin
//
// caskin has a special "superadmin" role that transcends all domain
// boundaries. Superadmins can manage any domain and bypass all permission
// checks. Use [IBaseService.AddSuperadmin] / [IBaseService.DeleteSuperadmin]
// to manage superadmin users.
//
// # Directory
//
// Objects can be organized into tree-structured directories. The
// [IDirectoryService] provides operations to create, move, copy, and
// delete directory nodes and their contents.
package caskin
