package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/arcanericky/opiekey"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var versionText = "unknown"

func getChallengeResponse(seq int, seed, passPhrase string, outputHex bool, codeCount int, alg opiekey.Algorithm) []string {
	var result []string

	// set response type
	computeRoutine := opiekey.ComputeWordResponse
	if outputHex {
		computeRoutine = opiekey.ComputeHexResponse
	}

	if codeCount > 1 {
		// multiple codes
		for i := seq - codeCount + 1; i <= seq; i++ {
			result = append(result, fmt.Sprintf("%d: %s", i,
				computeRoutine(i, seed, passPhrase, alg)))
		}
	} else {
		// single code
		result = append(result, computeRoutine(seq, seed, passPhrase, alg))
	}

	return result
}

func getPassPhrase(args []string, alg opiekey.Algorithm) string {
	passPhrase := ""
	if len(args) == 3 {
		passPhrase = args[2]
	} else {
		passPhrase = os.Getenv("OPIE_PASSPHRASE")

		if passPhrase == "" {
			fmt.Println("Using the", alg, "algorithm to compute response.")
			fmt.Println("Reminder: Don't use opiekey from telnet or dial-in sessions.")
			fmt.Print("Enter secret pass phrase: ")
			bytePassPhrase, _ := terminal.ReadPassword(int(syscall.Stdin))
			passPhrase = string(bytePassPhrase)
			fmt.Println("")
		}
	}

	return passPhrase
}

func getHashAlgorithm(md4, sha1 bool) opiekey.Algorithm {
	// sha1 > md4 > md5
	alg := opiekey.MD5
	if sha1 {
		alg = opiekey.SHA1
	} else if md4 {
		alg = opiekey.MD4
	}

	return alg
}

func outputChallengeResponse(md4, sha1, hex bool, codeCount int, args []string) {
	// get sequence number
	seq, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	alg := getHashAlgorithm(md4, sha1)

	for _, i := range getChallengeResponse(seq, args[1], getPassPhrase(args, alg),
		hex, codeCount, alg) {
		fmt.Println(i)
	}
}

func execute() {
	const cmdName = "opiekey"
	var hex, md4, md5, sha1 bool
	var codeCount int

	rootCmd := &cobra.Command{
		Use:   cmdName + " sequence_number seed passphrase",
		Short: cmdName,
		Long: cmdName + ` - Program for computing responses to OTP challenges

opiekey takes the optional count of the number of responses to print
along with a (maximum) sequence number, seed and optional secret pass
phrase as command line args then produces an OPIE response as six
words or hexadecimal numbers. If the OPIE_PASSPHRASE environment
variable is set it will be used for the secret pass phrase. If no
secret pass phrase was specified as a command line argument or
environment variable, the program will prompt for it.`,
		Version: versionText + " " + runtime.GOOS + "/" + runtime.GOARCH,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 && len(args) != 3 {
				return fmt.Errorf("accepts 2 or 3 args, received %d", len(args))
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			outputChallengeResponse(md4, sha1, hex, codeCount, args)
		},
	}

	rootCmd.Flags().BoolVarP(&md5, "md5", "5", true, "selects MD5 as the response generation algorithm")
	rootCmd.Flags().BoolVarP(&md4, "md4", "4", false, "selects MD4 as the response generation algorithm")
	rootCmd.Flags().BoolVarP(&sha1, "sha1", "s", false, "selects SHA1 as the response generation algorithm")
	rootCmd.Flags().BoolVarP(&hex, "hex", "x", false, "output the OTPs as hexadecimal numbers instead of six words")
	rootCmd.Flags().IntVarP(&codeCount, "number", "n", 1, "the number of one time access passwords to print")

	rootCmd.Execute()
}

func main() {
	execute()
}
