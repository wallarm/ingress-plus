package nginx

import (
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

var configMap = api.ConfigMap{
	ObjectMeta: api.ObjectMeta{
		Name:      "test",
		Namespace: "default",
	},
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
}
var ingress = extensions.Ingress{
	ObjectMeta: api.ObjectMeta{
		Name:      "test",
		Namespace: "kube-system",
	},
	TypeMeta: unversioned.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
}

//
// GetMapKeyAsBool
//
func TestGetMapKeyAsBool(t *testing.T) {
	configMap := configMap
	configMap.Data = map[string]string{
		"key": "True",
	}

	b, err := GetMapKeyAsBool(configMap.Data, "key", &configMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if b != true {
		t.Errorf("Result should be true")
	}
}

func TestGetMapKeyAsBoolNotFound(t *testing.T) {
	configMap := configMap
	configMap.Data = map[string]string{}

	_, err := GetMapKeyAsBool(configMap.Data, "key", &configMap)
	if err != ErrorKeyNotFound {
		t.Errorf("ErrorKeyNotFound was expected, got: %v", err)
	}
}

func TestGetMapKeyAsBoolErrorMessage(t *testing.T) {
	cfgm := configMap
	cfgm.Data = map[string]string{
		"key": "string",
	}

	// Test with configmap
	_, err := GetMapKeyAsBool(cfgm.Data, "key", &cfgm)
	if err == nil {
		t.Error("An error was expected")
	}
	expected := `ConfigMap default/test 'key' contains invalid bool: strconv.ParseBool: parsing "string": invalid syntax, ignoring`
	if err.Error() != expected {
		t.Errorf("The error message does not match expectations:\nGot: %v\nExpected: %v", err, expected)
	}

	// Test with ingress object
	ingress := ingress
	ingress.Annotations = map[string]string{
		"key": "other_string",
	}
	_, err = GetMapKeyAsBool(ingress.Annotations, "key", &ingress)
	if err == nil {
		t.Error("An error was expected")
	}
	expected = `Ingress kube-system/test 'key' contains invalid bool: strconv.ParseBool: parsing "other_string": invalid syntax, ignoring`
	if err.Error() != expected {
		t.Errorf("The error message does not match expectations:\nGot: %v\nExpected: %v", err, expected)
	}
}

//
// GetMapKeyAsInt
//
func TestGetMapKeyAsInt(t *testing.T) {
	configMap := configMap
	configMap.Data = map[string]string{
		"key": "123456789",
	}

	i, err := GetMapKeyAsInt(configMap.Data, "key", &configMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var expected int64 = 123456789
	if i != expected {
		t.Errorf("Unexpected return value:\nGot: %v\nExpected: %v", i, expected)
	}
}

func TestGetMapKeyAsIntNotFound(t *testing.T) {
	configMap := configMap
	configMap.Data = map[string]string{}

	_, err := GetMapKeyAsInt(configMap.Data, "key", &configMap)
	if err != ErrorKeyNotFound {
		t.Errorf("ErrorKeyNotFound was expected, got: %v", err)
	}
}

func TestGetMapKeyAsIntErrorMessage(t *testing.T) {
	cfgm := configMap
	cfgm.Data = map[string]string{
		"key": "string",
	}

	// Test with configmap
	_, err := GetMapKeyAsInt(cfgm.Data, "key", &cfgm)
	if err == nil {
		t.Error("An error was expected")
	}
	expected := `ConfigMap default/test 'key' contains invalid integer: strconv.ParseInt: parsing "string": invalid syntax, ignoring`
	if err.Error() != expected {
		t.Errorf("The error message does not match expectations:\nGot: %v\nExpected: %v", err, expected)
	}

	// Test with ingress object
	ingress := ingress
	ingress.Annotations = map[string]string{
		"key": "other_string",
	}
	_, err = GetMapKeyAsInt(ingress.Annotations, "key", &ingress)
	if err == nil {
		t.Error("An error was expected")
	}
	expected = `Ingress kube-system/test 'key' contains invalid integer: strconv.ParseInt: parsing "other_string": invalid syntax, ignoring`
	if err.Error() != expected {
		t.Errorf("The error message does not match expectations:\nGot: %v\nExpected: %v", err, expected)
	}
}
