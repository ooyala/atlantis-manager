package crypto

import (
	"atlantis/crypto"
	"atlantis/manager/rpc/types"
	"encoding/json"
)

func EncryptDependerEnvData(data *types.DependerEnvData) error {
	// encrypt DataMap and nil out DataMap
	// convert to JSON
	jsonBytes, err := json.Marshal(data.DataMap)
	if err != nil {
		return err
	}
	// encrypt into Data
	data.EncryptedData = string(crypto.Encrypt(jsonBytes))
	// nil out DataMap
	data.DataMap = nil
	return nil
}

func DecryptDependerEnvData(data *types.DependerEnvData) error {
	// decrypt Data to DataMap
	// decrypt Data
	decryptedBytes := crypto.Decrypt([]byte(data.EncryptedData))
	data.DataMap = map[string]interface{}{}
	// Unmarshal JSON to DataMap
	if err := json.Unmarshal(decryptedBytes, &data.DataMap); err != nil {
		return err
	}
	return nil
}
