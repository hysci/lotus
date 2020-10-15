package build

import (
	"sort"

	"github.com/filecoin-project/lotus/node/modules/dtypes"
)

type DrandEnum int

func DrandConfigSchedule() dtypes.DrandSchedule {
	out := dtypes.DrandSchedule{}
	for start, config := range DrandSchedule {
		out = append(out, dtypes.DrandPoint{Start: start, Config: DrandConfigs[config]})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Start < out[j].Start
	})

	return out
}

const (
	DrandMainnet DrandEnum = iota + 1
	DrandTestnet
	DrandDevnet
	DrandLocalnet
	DrandIncentinet
)

var DrandConfigs = map[DrandEnum]dtypes.DrandConfig{
	DrandMainnet: {
		Servers: []string{
			"https://api.drand.sh",
			"https://api2.drand.sh",
			"https://api3.drand.sh",
		},
		Relays: []string{
			"/dnsaddr/api.drand.sh/",
			"/dnsaddr/api2.drand.sh/",
			"/dnsaddr/api3.drand.sh/",
		},
		ChainInfoJSON: `{"public_key":"868f005eb8e6e4ca0a47c8a77ceaa5309a47978a7c71bc5cce96366b5d7a569937c529eeda66c7293784a9402801af31","period":30,"genesis_time":1595431050,"hash":"8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce","groupHash":"176f93498eac9ca337150b46d21dd58673ea4e3581185f869672e59fa4cb390a"}`,
	},
	DrandTestnet: {
		Servers: []string{
			"http://193.8.82.144:9090",
			"http://118.123.241.64:9090",
			"http://118.123.244.155:9090",
			"http://107.148.247.35:9090",
		},
		Relays:        []string{},
		ChainInfoJSON: `{"public_key":"a637ff023ec222c3a40f8192071f501a87b0b7bf024e286cae79104d44f7c3a77bc5c8cef8af5c002f85e7dc19b9062e","period":30,"genesis_time":1602770640,"hash":"e23455c5a6ce6d7248bf4e7425e398b3716fdb0a398ee27a582ee10c5819d3d6","groupHash":"067a43e171316341854eaae8de56db54135a2b6b107ad73a8fb7eb8299718c48"}`,
	},
	DrandDevnet: {
		Servers: []string{
			"https://dev1.drand.sh",
			"https://dev2.drand.sh",
		},
		Relays: []string{
			"/dnsaddr/dev1.drand.sh/",
			"/dnsaddr/dev2.drand.sh/",
		},
		ChainInfoJSON: `{"public_key":"8cda589f88914aa728fd183f383980b35789ce81b274e5daee1f338b77d02566ef4d3fb0098af1f844f10f9c803c1827","period":25,"genesis_time":1595348225,"hash":"e73b7dc3c4f6a236378220c0dd6aa110eb16eed26c11259606e07ee122838d4f","groupHash":"567d4785122a5a3e75a9bc9911d7ea807dd85ff76b78dc4ff06b075712898607"}`,
	},
	DrandIncentinet: {
		Servers: []string{
			"http://193.8.82.144:9090",
			"http://118.123.241.64:9090",
			"http://118.123.244.155:9090",
			"http://107.148.247.35:9090",
		},
		Relays:        []string{},
		ChainInfoJSON: `{"public_key":"a637ff023ec222c3a40f8192071f501a87b0b7bf024e286cae79104d44f7c3a77bc5c8cef8af5c002f85e7dc19b9062e","period":30,"genesis_time":1602770640,"hash":"e23455c5a6ce6d7248bf4e7425e398b3716fdb0a398ee27a582ee10c5819d3d6","groupHash":"067a43e171316341854eaae8de56db54135a2b6b107ad73a8fb7eb8299718c48"}`,
	},
}
