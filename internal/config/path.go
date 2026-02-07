package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Get retrieves a configuration value by key path (e.g., "ui.color", "default.branch")
func (cm *Manager) Get(key string) (any, error) {
	sanitized, err := sanitizeConfigPath(key)
	if err != nil {
		return nil, err
	}
	return cm.getValueByPath(cm.config, sanitized)
}

// Set sets a configuration value by key path
func (cm *Manager) Set(key string, value any) error {
	sanitized, err := sanitizeConfigPath(key)
	if err != nil {
		return err
	}
	if err := cm.setValueByPath(cm.config, sanitized, value); err != nil {
		return err
	}
	if err := cm.config.Validate(); err != nil {
		return err
	}
	return cm.Save()
}

func sanitizeConfigPath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("config path cannot be empty")
	}
	parts := strings.Split(trimmed, ".")
	for idx, part := range parts {
		if part == "" {
			return "", fmt.Errorf("config path segment %d is empty", idx+1)
		}
		if !configPathSegmentRe.MatchString(part) {
			return "", fmt.Errorf("config path segment %q contains invalid characters", part)
		}
	}
	return trimmed, nil
}

// List returns all configuration keys and values
func (cm *Manager) List() map[string]any {
	result := make(map[string]any)
	cm.flattenConfig(cm.config, "", result)
	return result
}

// getValueByPath retrieves a value using dot notation path
func (cm *Manager) getValueByPath(obj any, path string) (any, error) {
	parts := strings.Split(path, ".")
	current := reflect.ValueOf(obj)

	for _, part := range parts {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}

		switch current.Kind() {
		case reflect.Struct:
			field, found := cm.findFieldByYamlTag(current.Type(), current, part)
			if !found {
				return nil, fmt.Errorf("field '%s' not found", part)
			}
			current = field

		case reflect.Map:
			mapValue := current.MapIndex(reflect.ValueOf(part))
			if !mapValue.IsValid() {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
			current = mapValue

		default:
			return nil, fmt.Errorf("cannot navigate into %s", current.Kind())
		}
	}

	return current.Interface(), nil
}

// findFieldByYamlTag finds a struct field by its YAML tag or field name
func (cm *Manager) findFieldByYamlTag(structType reflect.Type, structValue reflect.Value, tagName string) (reflect.Value, bool) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Check YAML tag first
		yamlTag := field.Tag.Get("yaml")
		if yamlTag != "" {
			yamlName := strings.Split(yamlTag, ",")[0]
			if yamlName == tagName {
				return fieldValue, true
			}
		}

		// Fallback to field name (case-insensitive)
		if strings.EqualFold(field.Name, tagName) {
			return fieldValue, true
		}
	}
	return reflect.Value{}, false
}

// setValueByPath sets a value using dot notation path
func (cm *Manager) setValueByPath(obj any, path string, value any) error {
	parts := strings.Split(path, ".")
	current := reflect.ValueOf(obj)

	parent, err := cm.navigateToParent(current, parts)
	if err != nil {
		return err
	}

	return cm.setFinalValue(parent, parts[len(parts)-1], value)
}

// navigateToParent navigates to the parent of the target field
func (cm *Manager) navigateToParent(current reflect.Value, parts []string) (reflect.Value, error) {
	for i, part := range parts[:len(parts)-1] {
		var err error
		current, err = cm.navigateOneLevel(current, part, parts[:i+1])
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return current, nil
}

// navigateOneLevel navigates one level into a struct or map
func (cm *Manager) navigateOneLevel(current reflect.Value, part string, pathSoFar []string) (reflect.Value, error) {
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	switch current.Kind() {
	case reflect.Struct:
		field, found := cm.findFieldByYamlTag(current.Type(), current, part)
		if !found {
			return reflect.Value{}, fmt.Errorf("field '%s' not found", strings.Join(pathSoFar, "."))
		}
		return field, nil

	case reflect.Map:
		mapValue := current.MapIndex(reflect.ValueOf(part))
		if !mapValue.IsValid() {
			return reflect.Value{}, fmt.Errorf("key '%s' not found", strings.Join(pathSoFar, "."))
		}
		return mapValue, nil

	default:
		return reflect.Value{}, fmt.Errorf("cannot navigate into %s", current.Kind())
	}
}

// setFinalValue sets the final value in the target location
func (cm *Manager) setFinalValue(current reflect.Value, lastPart string, value any) error {
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	switch current.Kind() {
	case reflect.Struct:
		return cm.setStructField(current, lastPart, value)
	case reflect.Map:
		return cm.setMapValue(current, lastPart, value)
	default:
		return fmt.Errorf("cannot set value in %s", current.Kind())
	}
}

// setStructField sets a field value in a struct
func (cm *Manager) setStructField(current reflect.Value, fieldName string, value any) error {
	field, found := cm.findFieldByYamlTag(current.Type(), current, fieldName)
	if !found || !field.CanSet() {
		return fmt.Errorf("field '%s' not found or cannot be set", fieldName)
	}

	newValue := reflect.ValueOf(value)
	if !newValue.Type().ConvertibleTo(field.Type()) {
		return fmt.Errorf("cannot convert %s to %s", newValue.Type(), field.Type())
	}

	field.Set(newValue.Convert(field.Type()))
	return nil
}

// setMapValue sets a value in a map
func (cm *Manager) setMapValue(current reflect.Value, key string, value any) error {
	if current.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("map key must be string")
	}

	newValue := reflect.ValueOf(value)
	if !newValue.Type().ConvertibleTo(current.Type().Elem()) {
		return fmt.Errorf("cannot convert %s to %s", newValue.Type(), current.Type().Elem())
	}

	current.SetMapIndex(reflect.ValueOf(key), newValue.Convert(current.Type().Elem()))
	return nil
}

// flattenConfig converts nested config to flat key-value pairs
func (cm *Manager) flattenConfig(obj any, prefix string, result map[string]any) {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Struct:
		cm.flattenStruct(value, prefix, result)
	case reflect.Map:
		cm.flattenMap(value, prefix, result)
	}
}

func (cm *Manager) flattenStruct(value reflect.Value, prefix string, result map[string]any) {
	structType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := structType.Field(i)

		fieldName := fieldType.Name
		if yamlTag := fieldType.Tag.Get("yaml"); yamlTag != "" {
			fieldName = strings.Split(yamlTag, ",")[0]
		}
		key := fieldName
		if prefix != "" {
			key = prefix + "." + fieldName
		}
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Map && field.Type().Elem().Kind() != reflect.Interface) {
			cm.flattenConfig(field.Interface(), key, result)
		} else {
			result[key] = field.Interface()
		}
	}
}

func (cm *Manager) flattenMap(value reflect.Value, prefix string, result map[string]any) {
	for _, mapKey := range value.MapKeys() {
		mapValue := value.MapIndex(mapKey)
		key := mapKey.String()
		if prefix != "" {
			key = prefix + "." + mapKey.String()
		}
		result[key] = mapValue.Interface()
	}
}
