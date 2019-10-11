package opiekey

import (
	"testing"
)

func TestOPIEKey(t *testing.T) {
	for i, item := range []struct {
		seq        int
		seed       string
		passphrase string
		alg        Algorithm
		wordResult string
		hexResult  string
	}{
		{3, "az3817", "d0g x h0us3", MD5, "HAM LINT KIN LACE EDNA BEET", "18D6 0488 D477 AAAB"},
		{3, "az3817", "d0g x h0us3", MD4, "BOGY WING PEG HYDE GUSH SO", "5AFF 84CB 4EC9 187A"},

		{9995, "doggie", "get in the doghouse", MD5, "DOUR DIME RACY LAYS BOO NET", "760E 672D D5D0 8457"},

		{1337, "challenge", "this is an awesome passphrase", MD5, "TOIL TEAM ANNE FUR SUP THEY", "ED1D 0534 8AA3 EDD4"},

		{500, "testseed", "testpassphrase", MD5, "DEED WOLF LOAN HIND INCA HYMN", "719F A2C4 CC39 E73B"},
		{500, "testseed", "testpassphrase", MD4, "DEAD SONG SCAN LAM NICK AUTO", "70BC 475D 918C 449E"},

		{500, "testseed", "testpassphrase", SHA1, "TWIT GOT DOSE SURE HOOK CURB", "F1E2 EDD6 F2B9 A8DB"},
	} {
		_ = item.alg.String()

		wordResponse := ComputeWordResponse(item.seq, item.seed, item.passphrase, item.alg)
		if wordResponse != item.wordResult {
			t.Errorf("Fail on entry %d", i)
		}

		hexResponse := ComputeHexResponse(item.seq, item.seed, item.passphrase, item.alg)
		if hexResponse != item.hexResult {
			t.Errorf("Fail on entry %d", i)
		}
	}
}
