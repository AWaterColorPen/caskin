package caskin

import (
	"testing"
)

// testUser is a minimal User implementation for testing.
type testUser struct {
	ID   uint64
	Code string
}

func (u *testUser) GetID() uint64      { return u.ID }
func (u *testUser) SetID(id uint64)    { u.ID = id }
func (u *testUser) Encode() string     { return u.Code }
func (u *testUser) Decode(s string) error {
	u.Code = s
	return nil
}

// testRole is a minimal Role implementation for testing.
type testRole struct {
	ID       uint64
	Code     string
	ObjectID uint64
	DomainID uint64
}

func (r *testRole) GetID() uint64         { return r.ID }
func (r *testRole) SetID(id uint64)       { r.ID = id }
func (r *testRole) Encode() string        { return r.Code }
func (r *testRole) Decode(s string) error { r.Code = s; return nil }
func (r *testRole) GetObjectID() uint64   { return r.ObjectID }
func (r *testRole) SetObjectID(id uint64) { r.ObjectID = id }
func (r *testRole) GetDomainID() uint64   { return r.DomainID }
func (r *testRole) SetDomainID(id uint64) { r.DomainID = id }

// testObject is a minimal Object implementation for testing.
type testObject struct {
	ID       uint64
	Code     string
	ParentID uint64
	DomainID uint64
	ObjType  string
}

func (o *testObject) GetID() uint64         { return o.ID }
func (o *testObject) SetID(id uint64)       { o.ID = id }
func (o *testObject) Encode() string        { return o.Code }
func (o *testObject) Decode(s string) error { o.Code = s; return nil }
func (o *testObject) GetParentID() uint64   { return o.ParentID }
func (o *testObject) SetParentID(id uint64) { o.ParentID = id }
func (o *testObject) GetDomainID() uint64   { return o.DomainID }
func (o *testObject) SetDomainID(id uint64) { o.DomainID = id }
func (o *testObject) GetObjectType() string { return o.ObjType }

// testDomain is a minimal Domain implementation for testing.
type testDomain struct {
	ID   uint64
	Code string
}

func (d *testDomain) GetID() uint64      { return d.ID }
func (d *testDomain) SetID(id uint64)    { d.ID = id }
func (d *testDomain) Encode() string     { return d.Code }
func (d *testDomain) Decode(s string) error {
	d.Code = s
	return nil
}

// prefixObject decodes only strings starting with "prefix:".
type prefixObject struct {
	ID       uint64
	Name     string
	ParentID uint64
	DomainID uint64
}

func (o *prefixObject) GetID() uint64         { return o.ID }
func (o *prefixObject) SetID(id uint64)       { o.ID = id }
func (o *prefixObject) Encode() string        { return "prefix:" + o.Name }
func (o *prefixObject) Decode(s string) error {
	if len(s) < 7 || s[:7] != "prefix:" {
		return ErrInValidObject
	}
	o.Name = s[7:]
	return nil
}
func (o *prefixObject) GetParentID() uint64   { return o.ParentID }
func (o *prefixObject) SetParentID(id uint64) { o.ParentID = id }
func (o *prefixObject) GetDomainID() uint64   { return o.DomainID }
func (o *prefixObject) SetDomainID(id uint64) { o.DomainID = id }
func (o *prefixObject) GetObjectType() string { return "prefix" }

func TestRegister_NoOptions(t *testing.T) {
	Register[*testUser, *testRole, *testObject, *testDomain]()
	f := DefaultFactory()
	if f == nil {
		t.Fatal("DefaultFactory() returned nil")
	}

	// Basic decode should work
	u, err := f.User("alice")
	if err != nil {
		t.Fatalf("User decode failed: %v", err)
	}
	if u.Encode() != "alice" {
		t.Fatalf("expected 'alice', got %q", u.Encode())
	}

	// NamedObject built-in should be in the candidate list
	obj, err := f.Object("any_string")
	if err != nil {
		t.Fatalf("Object decode failed (NamedObject fallback): %v", err)
	}
	if obj.Encode() != "any_string" {
		t.Fatalf("expected 'any_string', got %q", obj.Encode())
	}
}

// selectiveObject only decodes strings starting with "sel:".
type selectiveObject struct {
	ID       uint64
	Code     string
	ParentID uint64
	DomainID uint64
}

func (o *selectiveObject) GetID() uint64         { return o.ID }
func (o *selectiveObject) SetID(id uint64)       { o.ID = id }
func (o *selectiveObject) Encode() string        { return "sel:" + o.Code }
func (o *selectiveObject) Decode(s string) error {
	if len(s) < 4 || s[:4] != "sel:" {
		return ErrInValidObject
	}
	o.Code = s[4:]
	return nil
}
func (o *selectiveObject) GetParentID() uint64   { return o.ParentID }
func (o *selectiveObject) SetParentID(id uint64) { o.ParentID = id }
func (o *selectiveObject) GetDomainID() uint64   { return o.DomainID }
func (o *selectiveObject) SetDomainID(id uint64) { o.DomainID = id }
func (o *selectiveObject) GetObjectType() string { return "selective" }

func TestRegister_WithObject(t *testing.T) {
	// Use selectiveObject as the primary type (only accepts "sel:" prefix).
	// Add prefixObject via WithObject (only accepts "prefix:" prefix).
	// NamedObject (built-in) will be fallback for anything else.
	Register[*testUser, *testRole, *selectiveObject, *testDomain](
		WithObject(&prefixObject{}),
	)
	f := DefaultFactory()

	// "sel:hello" should be decoded by selectiveObject (primary)
	obj, err := f.Object("sel:hello")
	if err != nil {
		t.Fatalf("Object decode 'sel:hello' failed: %v", err)
	}
	if obj.GetObjectType() != "selective" {
		t.Fatalf("expected object type 'selective', got %q", obj.GetObjectType())
	}

	// "prefix:foo" should be decoded by prefixObject (extra candidate)
	obj2, err := f.Object("prefix:foo")
	if err != nil {
		t.Fatalf("Object decode 'prefix:foo' failed: %v", err)
	}
	if obj2.GetObjectType() != "prefix" {
		t.Fatalf("expected object type 'prefix', got %q", obj2.GetObjectType())
	}
	if obj2.Encode() != "prefix:foo" {
		t.Fatalf("expected 'prefix:foo', got %q", obj2.Encode())
	}

	// "plain_string" should fallback to NamedObject (built-in)
	obj3, err := f.Object("plain_string")
	if err != nil {
		t.Fatalf("Object decode fallback to NamedObject failed: %v", err)
	}
	if obj3.Encode() != "plain_string" {
		t.Fatalf("expected 'plain_string', got %q", obj3.Encode())
	}
}

func TestRegister_WithoutBuiltins(t *testing.T) {
	Register[*testUser, *testRole, *testObject, *testDomain](
		WithoutBuiltins(),
	)
	f := DefaultFactory()

	// testObject accepts anything, so it will still decode
	obj, err := f.Object("anything")
	if err != nil {
		t.Fatalf("Object decode failed without builtins: %v", err)
	}
	if obj.Encode() != "anything" {
		t.Fatalf("expected 'anything', got %q", obj.Encode())
	}
}

func TestDecode_ErrorWrapping(t *testing.T) {
	Register[*testUser, *testRole, *testObject, *testDomain](
		WithoutBuiltins(),
		WithObject(&prefixObject{}),
	)
	f := DefaultFactory()

	// testObject decodes anything, so this won't fail.
	// Let's test with a type that rejects input.
	// Create a factory with only prefixObject (no testObject).
	// We can't easily do that with the current generic Register, but we can
	// test the error message format by checking that a failed decode includes
	// the type info.
	_ = f // Just ensure compilation works; the detailed error test is below.

	// Direct test of decode function with a restrictive candidate list.
	candidates := []Object{&prefixObject{}}
	_, err := decode[Object]("no_prefix_here", candidates)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	errMsg := err.Error()
	if !contains(errMsg, "no registered factory for") {
		t.Fatalf("error message missing expected prefix: %v", errMsg)
	}
	if !contains(errMsg, "prefixObject") {
		t.Fatalf("error message should mention candidate type: %v", errMsg)
	}
	if !contains(errMsg, ErrInValidObject.Error()) {
		t.Fatalf("error message should contain wrapped error: %v", errMsg)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
