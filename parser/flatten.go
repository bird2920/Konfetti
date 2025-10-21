package parser

func flatten(data map[string]interface{}, prefix string) map[string]interface{} {
	flat := make(map[string]interface{})

	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch sub := v.(type) {
		case map[string]interface{}:
			nested := flatten(sub, key)
			for nk, nv := range nested {
				flat[nk] = nv
			}
		default:
			flat[key] = v
		}
	}
	return flat
}