package jsonconfig

type mockSecretManager struct {
	ValuesPool map[string]string
}

func (m *mockSecretManager) GetByName(keyName string) (string, bool) {
	result, isExists := m.ValuesPool[keyName]

	return result, isExists
}

func (m *mockSecretManager) GetByNameAndPath(keyName string) (string, bool) {
	result, isExists := m.ValuesPool[keyName]

	return result, isExists
}
