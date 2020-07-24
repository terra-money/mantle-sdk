package utils

import (
	"testing"

	types "github.com/terra-project/mantle/types"
)

func TestUnmarshalBlockResponseLCD(t *testing.T) {
	var blockResponse = `{"block_meta":{"block_id":{"hash":"4BED0CF615F23536A4259AB2D261377D56CE968D2B10ECA14FE4C2A0F600B53B","parts":{"total":"1","hash":"F2460B8C417E3053D2D4E6589EB7400B05DABA55CAEE978D0925F4D610640D2B"}},"header":{"version":{"block":"10","app":"0"},"chain_id":"columbus-3","height":"1","time":"2019-12-13T16:42:15Z","num_txs":"0","total_txs":"0","last_block_id":{"hash":"","parts":{"total":"0","hash":""}},"last_commit_hash":"","data_hash":"","validators_hash":"B736F13D962B7686D50FE8577EAE2024B80D0B639E41828A43D403FF998EA3DB","next_validators_hash":"B736F13D962B7686D50FE8577EAE2024B80D0B639E41828A43D403FF998EA3DB","consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F","app_hash":"","last_results_hash":"","evidence_hash":"","proposer_address":"41DD2D4597CA968941F6B92E73C6EDED65ADE1D1"}},"block":{"header":{"version":{"block":"10","app":"0"},"chain_id":"columbus-3","height":"1","time":"2019-12-13T16:42:15Z","num_txs":"0","total_txs":"0","last_block_id":{"hash":"","parts":{"total":"0","hash":""}},"last_commit_hash":"","data_hash":"","validators_hash":"B736F13D962B7686D50FE8577EAE2024B80D0B639E41828A43D403FF998EA3DB","next_validators_hash":"B736F13D962B7686D50FE8577EAE2024B80D0B639E41828A43D403FF998EA3DB","consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F","app_hash":"","last_results_hash":"","evidence_hash":"","proposer_address":"41DD2D4597CA968941F6B92E73C6EDED65ADE1D1"},"data":{"txs":null},"evidence":{"evidence":null},"last_commit":{"block_id":{"hash":"","parts":{"total":"0","hash":""}},"precommits":null}}}`
	var block = types.Block{}

	UnmarshalBlockResponseFromLCD([]byte(blockResponse), &block)

	t.Log(block)
}
