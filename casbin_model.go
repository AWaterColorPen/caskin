package caskin

import (
    "bytes"
    "net/http"

    "github.com/casbin/casbin/v2/model"
)

func CasbinModel() (model.Model, error) {
    return getCasbinModelFromUrl(ModelConfPathSuperadmin)
}

func CasbinModelNoSuperadmin() (model.Model, error) {
    return getCasbinModelFromUrl(ModelConfPathNoSuperadmin)
}

func getCasbinModelFromUrl(url string) (model.Model, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
    buf := &bytes.Buffer{}
    _, err = buf.ReadFrom(resp.Body)
    if err != nil {
        return nil, err
    }

    return model.NewModelFromString(buf.String())
}
