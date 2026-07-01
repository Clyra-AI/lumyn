package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
)

func yamlScalar(value string) string {
	if value == "" {
		return `""`
	}
	if strings.ContainsAny(value, "\n\r\t#:") || strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") {
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}
	return value
}

func parseYAMLValue(value string) string {
	value = stripYAMLInlineComment(value)
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return value
}

func stripYAMLInlineComment(value string) string {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	for index, char := range value {
		switch {
		case escaped:
			escaped = false
		case char == '\\' && inDoubleQuote:
			escaped = true
		case char == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case char == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case char == '#' && !inSingleQuote && !inDoubleQuote:
			if index > 0 && (value[index-1] == ' ' || value[index-1] == '\t') {
				return value[:index]
			}
		}
	}
	return value
}

func yamlInlineValueHasEntries(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, `"'`))
	switch value {
	case "", "{}", "[]", "null", "~":
		return false
	default:
		return true
	}
}

func yamlLocalComponentRef(value, component string) (string, bool) {
	prefix := "#/components/" + component + "/"
	for _, candidate := range []string{parseYAMLValue(value), yamlInlineRefValue(value)} {
		if strings.HasPrefix(candidate, prefix) {
			return unescapeJSONPointer(strings.TrimPrefix(candidate, prefix)), true
		}
	}
	return "", false
}

func yamlInlineRefValue(value string) string {
	value = strings.TrimSpace(stripYAMLInlineComment(value))
	value = strings.Trim(value, "{} ")
	for _, entry := range strings.Split(value, ",") {
		key, refValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok || strings.Trim(key, `"'`) != "$ref" {
			continue
		}
		return parseYAMLValue(strings.Trim(refValue, "{} "))
	}
	return ""
}

func parseYAMLParameterLine(parameters []openAPIParameter, current **openAPIParameter, trimmed, key, value, pointer string, lineNo int, parameterComponents map[string]openAPIParameter) []openAPIParameter {
	if strings.HasPrefix(trimmed, "- ") {
		item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		if strings.HasPrefix(item, "{") {
			if refName, ok := yamlLocalComponentRef(item, "parameters"); ok {
				if referenced, ok := parameterComponents[refName]; ok {
					parameters = append(parameters, referenced)
				} else {
					parameters = append(parameters, openAPIParameter{
						Name:    yamlInlineRefValue(item),
						Pointer: pointer,
						Line:    lineNo,
					})
				}
				*current = &parameters[len(parameters)-1]
				return parameters
			}
			parameter := yamlInlineParameter(item, pointer)
			parameter.Line = lineNo
			parameters = append(parameters, parameter)
			*current = &parameters[len(parameters)-1]
			return parameters
		}
		parameterKey, parameterValue, ok := yamlKeyValue(item)
		parameter := openAPIParameter{
			Pointer: pointer,
			Line:    lineNo,
		}
		if ok {
			if parameterKey == "$ref" {
				if refName, ok := yamlLocalComponentRef(parameterValue, "parameters"); ok {
					if referenced, ok := parameterComponents[refName]; ok {
						parameters = append(parameters, referenced)
					} else {
						parameters = append(parameters, openAPIParameter{
							Name:    parameterValue,
							Pointer: pointer,
							Line:    lineNo,
						})
					}
					*current = &parameters[len(parameters)-1]
					return parameters
				}
			}
			switch parameterKey {
			case "name":
				parameter.Name = parameterValue
			case "in":
				parameter.In = parameterValue
			case "description":
				parameter.Description = parameterValue
			}
		}
		parameters = append(parameters, parameter)
		*current = &parameters[len(parameters)-1]
		return parameters
	}
	if *current == nil {
		return parameters
	}
	switch key {
	case "name":
		(*current).Name = value
	case "in":
		(*current).In = value
	case "description":
		(*current).Description = value
	}
	return parameters
}

func yamlComponentParameters(data []byte) map[string]openAPIParameter {
	parameters := map[string]openAPIParameter{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	parametersIndent := -1
	itemIndent := -1
	currentItem := ""
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" {
			continue
		}
		indent := leadingSpaces(scanner.Text())
		key, value, ok := yamlKeyValue(trimmed)
		if !ok {
			continue
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
		}
		if parametersIndent >= 0 && indent <= parametersIndent && key != "parameters" {
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
		}
		if itemIndent >= 0 && indent <= itemIndent && key != currentItem {
			itemIndent = -1
			currentItem = ""
		}
		if key == "components" && indent == 0 {
			componentsIndent = indent
			parametersIndent = -1
			itemIndent = -1
			currentItem = ""
			continue
		}
		if componentsIndent >= 0 && indent > componentsIndent && key == "parameters" {
			parametersIndent = indent
			itemIndent = -1
			currentItem = ""
			continue
		}
		if parametersIndent >= 0 && indent > parametersIndent && (itemIndent < 0 || indent <= itemIndent) {
			currentItem = key
			itemIndent = indent
			parameters[currentItem] = yamlInlineParameter(value, "#/components/parameters/"+escapeJSONPointer(currentItem))
			continue
		}
		if currentItem == "" || indent <= itemIndent {
			continue
		}
		parameter := parameters[currentItem]
		switch key {
		case "name":
			parameter.Name = value
		case "in":
			parameter.In = value
		case "description":
			parameter.Description = value
		}
		parameters[currentItem] = parameter
	}
	return parameters
}

func yamlInlineParameter(value, pointer string) openAPIParameter {
	parameter := openAPIParameter{Pointer: pointer}
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	for _, entry := range strings.Split(value, ",") {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if !ok {
			continue
		}
		switch key {
		case "name":
			parameter.Name = entryValue
		case "in":
			parameter.In = entryValue
		case "description":
			parameter.Description = entryValue
		}
	}
	return parameter
}

func yamlInlineParameterSequence(value, pointer string, parameterComponents map[string]openAPIParameter) []openAPIParameter {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
		return nil
	}
	value = strings.TrimSpace(strings.Trim(value, "[]"))
	if value == "" {
		return nil
	}
	parameters := []openAPIParameter{}
	for _, item := range splitYAMLFlowEntries(value) {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if refName, ok := yamlLocalComponentRef(item, "parameters"); ok {
			if referenced, ok := parameterComponents[refName]; ok {
				parameters = append(parameters, referenced)
			} else {
				parameters = append(parameters, openAPIParameter{
					Name:    yamlInlineRefValue(item),
					Pointer: pointer,
				})
			}
			continue
		}
		parameters = append(parameters, yamlInlineParameter(item, pointer))
	}
	return parameters
}

func splitYAMLFlowEntries(value string) []string {
	entries := []string{}
	start := 0
	depth := 0
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	for index, char := range value {
		switch {
		case escaped:
			escaped = false
		case char == '\\' && inDoubleQuote:
			escaped = true
		case char == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case char == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case inSingleQuote || inDoubleQuote:
			continue
		case char == '{' || char == '[':
			depth++
		case char == '}' || char == ']':
			if depth > 0 {
				depth--
			}
		case char == ',' && depth == 0:
			entries = append(entries, value[start:index])
			start = index + 1
		}
	}
	entries = append(entries, value[start:])
	return entries
}

func yamlComponentContentSchemaRefs(data []byte, group string) map[string]bool {
	refs := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	componentsIndent := -1
	groupIndent := -1
	itemIndent := -1
	contentIndent := -1
	mediaIndent := -1
	currentItem := ""
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" {
			continue
		}
		indent := leadingSpaces(scanner.Text())
		key, value, ok := yamlKeyValue(trimmed)
		if !ok {
			continue
		}
		if componentsIndent >= 0 && indent <= componentsIndent && key != "components" {
			componentsIndent = -1
			groupIndent = -1
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if groupIndent >= 0 && indent <= groupIndent && key != group {
			groupIndent = -1
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if itemIndent >= 0 && indent <= itemIndent && key != currentItem {
			itemIndent = -1
			contentIndent = -1
			mediaIndent = -1
			currentItem = ""
		}
		if contentIndent >= 0 && indent <= contentIndent && key != "content" {
			contentIndent = -1
			mediaIndent = -1
		}
		if mediaIndent >= 0 && indent <= mediaIndent {
			mediaIndent = -1
		}
		if key == "components" && indent == 0 {
			componentsIndent = indent
			continue
		}
		if componentsIndent >= 0 && indent > componentsIndent && key == group {
			groupIndent = indent
			for item, hasSchema := range yamlInlineComponentContentSchemaRefs(value) {
				if hasSchema {
					refs[item] = true
				}
			}
			continue
		}
		if groupIndent >= 0 && indent > groupIndent && (itemIndent < 0 || indent <= itemIndent) {
			currentItem = key
			itemIndent = indent
			contentIndent = -1
			mediaIndent = -1
			continue
		}
		if currentItem == "" {
			continue
		}
		if itemIndent >= 0 && indent > itemIndent && key == "content" {
			contentIndent = indent
			mediaIndent = -1
			continue
		}
		if contentIndent >= 0 && indent > contentIndent && strings.Contains(key, "/") {
			mediaIndent = indent
			if yamlInlineValueHasSchema(value) {
				refs[currentItem] = true
			}
			continue
		}
		if mediaIndent >= 0 && indent > mediaIndent && key == "schema" {
			refs[currentItem] = true
		}
	}
	return refs
}

func yamlInlineComponentContentSchemaRefs(value string) map[string]bool {
	refs := map[string]bool{}
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return refs
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		itemName, itemValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && yamlInlineContentContainerHasSchema(itemValue) {
			refs[itemName] = true
		}
	}
	return refs
}

func yamlInlineValueHasSchema(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	schemaValue, ok := yamlInlineMapValue(value, "schema")
	return ok && yamlInlineValueHasEntries(schemaValue)
}

func yamlInlineContentContainerHasSchema(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && strings.Trim(key, `"'`) == "content" {
			return yamlInlineContentValueHasSchema(entryValue)
		}
	}
	return false
}

func yamlInlineContentValueHasSchema(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && strings.Contains(strings.Trim(key, `"'`), "/") && yamlInlineValueHasSchema(entryValue) {
			return true
		}
	}
	return false
}

func yamlInlineMapHasKey(value, expected string) bool {
	_, ok := yamlInlineMapValue(value, expected)
	return ok
}

func yamlInlineMapValueEquals(value, expected, want string) bool {
	got, ok := yamlInlineMapValue(value, expected)
	return ok && strings.EqualFold(strings.Trim(got, `"'`), want)
}

func yamlInlineMapValueNonEmpty(value, expected string) bool {
	got, ok := yamlInlineMapValue(value, expected)
	return ok && strings.Trim(got, `"'`) != ""
}

func yamlInlineSecuritySchemesHaveType(value string) bool {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	if value == "" {
		return false
	}
	for _, entry := range splitYAMLFlowEntries(value) {
		_, schemeValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && yamlInlineMapValueNonEmpty(schemeValue, "type") {
			return true
		}
	}
	return false
}

func yamlInlineMapValue(value, expected string) (string, bool) {
	value = strings.TrimSpace(strings.Trim(value, "{} "))
	for _, entry := range splitYAMLFlowEntries(value) {
		key, entryValue, ok := yamlKeyValue(strings.TrimSpace(entry))
		if ok && strings.Trim(key, `"'`) == expected {
			return entryValue, true
		}
	}
	return "", false
}

func yamlKeyValue(trimmed string) (string, string, bool) {
	index := yamlMappingSeparator(trimmed)
	if index < 0 {
		return "", "", false
	}
	key := strings.Trim(strings.TrimSpace(trimmed[:index]), `"'`)
	value := parseYAMLValue(trimmed[index+1:])
	return key, value, true
}

func yamlMappingSeparator(trimmed string) int {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	for index, char := range trimmed {
		switch {
		case escaped:
			escaped = false
		case char == '\\' && inDoubleQuote:
			escaped = true
		case char == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case char == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case char == ':' && !inSingleQuote && !inDoubleQuote:
			if index+1 == len(trimmed) || isYAMLMappingValueStart(trimmed[index+1]) {
				return index
			}
		}
	}
	return -1
}

func isYAMLMappingValueStart(char byte) bool {
	switch char {
	case ' ', '\t', '{', '[', '"', '\'':
		return true
	default:
		return false
	}
}

func leadingSpaces(value string) int {
	return len(value) - len(strings.TrimLeft(value, " "))
}

func is2xxStatusKey(key string) bool {
	key = strings.Trim(key, `"'`)
	return len(key) == 3 && key[0] == '2'
}

func escapeJSONPointer(value string) string {
	value = strings.ReplaceAll(value, "~", "~0")
	value = strings.ReplaceAll(value, "/", "~1")
	return value
}

func unescapeJSONPointer(value string) string {
	value = strings.ReplaceAll(value, "~1", "/")
	value = strings.ReplaceAll(value, "~0", "~")
	return value
}
