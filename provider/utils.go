package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func validateListIdName(objectList []interface{}, id string, name string) diag.Diagnostics {
	for _, obj := range objectList {
		m, ok := obj.(*schema.ResourceData)
		if !ok {
			continue // skip if not a map
		}

		_id := m.Get(id)
		_name := m.Get(name)

		if _id == "" && _name == "" {
			return diag.Errorf("Either %s or %s must be provided", id, name)
		}
	}

	return nil
}

func validateIdName(id string, name string) diag.Diagnostics {
	if id == "" && name == "" {
		return diag.Errorf("Either %s or %s must be provided", id, name)
	}

	return nil
}
