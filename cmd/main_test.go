package main

import (
	"os"
	"testing"

	"github.com/arcanericky/opiekey"
)

const (
	ROOTCMD    = "opiekey"
	SEQ        = "499"
	SEED       = "testseed"
	PASSPHRASE = "testpassphrase"
)

func TestMain(t *testing.T) {
	savedArgs := os.Args

	for _, args := range [][]string{
		// no options
		{ROOTCMD},
		// prompt for passphrase
		{ROOTCMD, SEQ, SEED},
		// including passphrase
		{ROOTCMD, SEQ, SEED, PASSPHRASE},
		// hex output
		{ROOTCMD, "-x", SEQ, SEED, PASSPHRASE},
		// multiple output responses
		{ROOTCMD, "-n", "5", SEQ, SEED, PASSPHRASE},
		// md4
		{ROOTCMD, "-4", SEQ, SEED, PASSPHRASE},
		// sha1
		{ROOTCMD, "-s", SEQ, SEED, PASSPHRASE},
		// invalid sequence number
		{ROOTCMD, "abc", SEED, PASSPHRASE},
	} {
		os.Args = args
		main()
	}

	os.Args = savedArgs

}

func TestGetChallengeResponse(t *testing.T) {
	type testItem struct {
		seq           int
		seed          string
		passPhrase    string
		hex           bool
		responseCount int
		alg           opiekey.Algorithm
	}

	for i, testParms := range []struct {
		seq             int
		seed            string
		passPhrase      string
		hex             bool
		responseCount   int
		alg             opiekey.Algorithm
		expectedResults []string
	}{
		{499, SEED, PASSPHRASE, false, 1, opiekey.MD5, []string{"LAIR FUME GIBE FROM JIG COP"}},
		{499, SEED, PASSPHRASE, false, 2, opiekey.MD5, []string{
			"498: SAVE CULT FEET WOOL KILL HATE",
			"499: LAIR FUME GIBE FROM JIG COP",
		}},
		{499, SEED, PASSPHRASE, false, 1, opiekey.MD4, []string{"BOSE DOOM LOS MADE LUSH NONE"}},
		{499, SEED, PASSPHRASE, false, 2, opiekey.MD4, []string{
			"498: TOW HURD GENE PIT SING WORM",
			"499: BOSE DOOM LOS MADE LUSH NONE",
		}},
		{499, SEED, PASSPHRASE, true, 1, opiekey.MD5, []string{"A9B0 E62C 4362 0217"}},
		{499, SEED, PASSPHRASE, true, 2, opiekey.MD5, []string{
			"498: D72D ADFB FEBA 6929",
			"499: A9B0 E62C 4362 0217",
		}},
		{499, SEED, PASSPHRASE, true, 1, opiekey.MD4, []string{"5D6E A896 DADB 4D8A"}},
		{499, SEED, PASSPHRASE, true, 2, opiekey.MD4, []string{
			"498: 41F3 A22A 1A0D CBFB",
			"499: 5D6E A896 DADB 4D8A",
		}},
	} {
		result := getChallengeResponse(testParms.seq, testParms.seed, testParms.passPhrase,
			testParms.hex, testParms.responseCount, testParms.alg)
		if len(result) != testParms.responseCount {
			t.Errorf("Result length not correct for test %d", i)
		}

		if len(testParms.expectedResults) > 0 {
			for j, expectedResult := range result {
				if expectedResult != testParms.expectedResults[j] {
					t.Errorf("Challenge output incorrect. Expected: %s. Generated: %s",
						expectedResult, testParms.expectedResults[j])
				}
			}
		}
	}
}

func TestGetHashAlgorithm(t *testing.T) {
	for i, testParms := range []struct {
		md4            bool
		sha1           bool
		expectedResult opiekey.Algorithm
	}{
		{false, false, opiekey.MD5},
		{true, false, opiekey.MD4},
		{false, true, opiekey.SHA1},
		{true, true, opiekey.SHA1},
	} {
		result := getHashAlgorithm(testParms.md4, testParms.sha1)

		if result != testParms.expectedResult {
			t.Errorf("Hash algorithm failed for test %d", i)
		}
	}
}

func TestGetPassPhrase(t *testing.T) {
	const (
		CLIPASSPHRASE   = "clipassphrase"
		ENVPASSPHRASE   = "envpassphrase"
		STDINPASSPHRASE = ""
	)

	// prompt on no environment or cli passphrase
	args := []string{"499", "testseed"}
	result := getPassPhrase(args, opiekey.MD5)
	if result != STDINPASSPHRASE {
		t.Errorf("Getting passphrase from stdin")
	}

	// env pass phrase
	os.Setenv("OPIE_PASSPHRASE", ENVPASSPHRASE)
	result = getPassPhrase(args, opiekey.MD5)
	if result != ENVPASSPHRASE {
		t.Errorf("Getting passphrase from environment")
	}

	// cli pass phrase overrides all
	args = append(args, CLIPASSPHRASE)
	result = getPassPhrase(args, opiekey.MD5)
	if result != CLIPASSPHRASE {
		t.Errorf("Getting passphrase from command line")
	}
}
