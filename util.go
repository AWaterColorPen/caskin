package caskin

import (
    "github.com/ahmetb/go-linq/v3"
)

func Filter(e ienforcer, u User, d Domain, action Action, fn func() Object, source interface{}) interface{} {
    linq.From(source).Where(func(v interface{}) bool {
        return Check(e, u, d, action, fn, v.(entry))
    }).ToSlice(&source)
    return source
}

func Check(e ienforcer, u User, d Domain, action Action, fn func() Object, one entry) bool {
    if !one.IsObject() {
        return true
    }

    o := fn()
    _ = o.Decode(one.GetObject())
    ok, _ := e.Enforce(u, o, d, action)
    return ok
}

func isValid(e entry) error {
    if e.GetID() == 0 {
        return ErrEmptyID
    }

    return  nil
}