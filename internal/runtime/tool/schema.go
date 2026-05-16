package tool

type Schema map[string]interface{}

func ObjectSchema(properties map[string]interface{}, required ...string) Schema {
	schema := Schema{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func StringSchema(description string) map[string]interface{} {
	schema := map[string]interface{}{"type": "string"}
	if description != "" {
		schema["description"] = description
	}
	return schema
}

func BooleanSchema(description string) map[string]interface{} {
	schema := map[string]interface{}{"type": "boolean"}
	if description != "" {
		schema["description"] = description
	}
	return schema
}

func NumberSchema(description string) map[string]interface{} {
	schema := map[string]interface{}{"type": "number"}
	if description != "" {
		schema["description"] = description
	}
	return schema
}
