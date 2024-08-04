package gnosis

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSDK_GetTopDelegates(t *testing.T) {
	s := &SDK{
		commonHTTP: &commonHTTP{},
	}

	rm := `{"name":"split-delegation","network":"1","params":{"backendUrl":"https://delegate-api.gnosisguild.org","strategies":[{"name":"erc20-balance-of","params":{"symbol":"SAFE","address":"0x5aFE3855358E112B5647B952709E6165e1c1eEEe","decimals":18},"network":1},{"name":"safe-vested","params":{"symbol":"SAFE (vested)","claimDateLimit":"2022-12-27T10:00:00+00:00","allocationsSource":"https://safe-claiming-app-data.gnosis-safe.io/allocations/1/snapshot-allocations-data.json"},"network":1},{"name":"contract-call","params":{"symbol":"SAFE (locked)","address":"0x0a7CB434f96f65972D46A5c1A64a9654dC9959b2","decimals":18,"methodABI":{"name":"getUserTokenBalance","type":"function","inputs":[{"name":"holder","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint96","internalType":"uint96"}],"stateMutability":"view"}},"network":1}],"totalSupply":1000000000}}`
	got, err := s.GetTopDelegates(nil, TopDelegatesRequest{
		Dao:      "safe.ggtest.eth",
		Strategy: json.RawMessage(rm),
		By:       "power",
		Limit:    18,
		Offset:   0,
	})

	if err != nil {
		t.Errorf("SDK.GetTopDelegates() error = %v", err)
		return
	}

	fmt.Println(got)
}

func TestSDK_GetDelegateProfile(t *testing.T) {
	s := &SDK{
		commonHTTP: &commonHTTP{},
	}

	rm := `{"name":"split-delegation","network":"1","params":{"backendUrl":"https://delegate-api.gnosisguild.org","strategies":[{"name":"erc20-balance-of","params":{"symbol":"SAFE","address":"0x5aFE3855358E112B5647B952709E6165e1c1eEEe","decimals":18},"network":1},{"name":"safe-vested","params":{"symbol":"SAFE (vested)","claimDateLimit":"2022-12-27T10:00:00+00:00","allocationsSource":"https://safe-claiming-app-data.gnosis-safe.io/allocations/1/snapshot-allocations-data.json"},"network":1},{"name":"contract-call","params":{"symbol":"SAFE (locked)","address":"0x0a7CB434f96f65972D46A5c1A64a9654dC9959b2","decimals":18,"methodABI":{"name":"getUserTokenBalance","type":"function","inputs":[{"name":"holder","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint96","internalType":"uint96"}],"stateMutability":"view"}},"network":1}],"totalSupply":1000000000}}`
	got, err := s.GetDelegateProfile(nil, DelegateProfileRequest{
		Dao:      "safe.ggtest.eth",
		Strategy: json.RawMessage(rm),
		Address:  "0x7697cAB0e123c68d27d7D5A9EbA346d7584Af888",
	})

	if err != nil {
		t.Errorf("SDK.GetTopDelegates() error = %v", err)
		return
	}

	marshal, err := json.Marshal(got)
	fmt.Println(string(marshal))
}
