package provider

func removeEmptyEntries(objectList []interface{}) []interface{} {
    for _, obj := range objectList {
        m, ok := obj.(map[string]interface{})
        if !ok {
            continue // skip if not a map
        }
        for k, v := range m {
            if v == nil || v == "" { // check for empty value
                delete(m, k) // delete the key with empty value
            }
        }
    }

    return objectList
}

func validateIdName(objectList []interface{}, id string, name string) {
    for _, obj := range objectList {
        _id := obj.Get(id)
        _name := obj.Get(name)

        if _id == "" && _name == "" {
            return diag.Errorf("Either %s or %s must be provided", id, name)
        }
    }
}

func validateIdName(id string, name string) {
    if _id == "" && _name == "" {
        return diag.Errorf("Either %s or %s must be provided", id, name)
    }
}