package holder_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tokencard/contracts/pkg/bindings"
	"github.com/tokencard/contracts/pkg/bindings/mocks"
	. "github.com/tokencard/contracts/test/shared"
)

func init() {
	TestRig.AddCoverageForContracts(
		"../../build/holder/combined.json",
		"../../contracts")

	TestRig.AddCoverageForContracts(
		"../../build/tkn/combined.json",
		"../../contracts")
}

func TestHolderSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contract Suite")
}

var _ = BeforeEach(func() {
	err := InitializeBackend()
	Expect(err).ToNot(HaveOccurred())

	var tx *types.Transaction
	ERC20Contract1Address, tx, ERC20Contract1, err = mocks.DeployToken(RandomAccount.TransactOpts(), Backend)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

	ERC20Contract2Address, tx, ERC20Contract2, err = mocks.DeployToken(RandomAccount.TransactOpts(), Backend)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

	TKNBurnerAddress, tx, TKNBurner, err = bindings.DeployToken(Owner.TransactOpts(), Backend)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

	TokenHolderAddress, tx, TokenHolder, err = bindings.DeployHolder(Owner.TransactOpts(), Backend, TKNBurnerAddress)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

})

var _ = AfterSuite(func() {
	TestRig.ExpectMinimumCoverage("holder.sol", 0.00)
	TestRig.ExpectMinimumCoverage("tkn.sol", 10.00)
	TestRig.PrintGasUsage(os.Stdout)
})

func isSuccessful(tx *types.Transaction) bool {
	r, err := Backend.TransactionReceipt(context.Background(), tx.Hash())
	Expect(err).ToNot(HaveOccurred())
	return r.Status == types.ReceiptStatusSuccessful
}

func isGasExhausted(tx *types.Transaction, gasLimit uint64) bool {
	r, err := Backend.TransactionReceipt(context.Background(), tx.Hash())
	Expect(err).ToNot(HaveOccurred())
	if r.Status == types.ReceiptStatusSuccessful {
		return false
	}
	return r.GasUsed == gasLimit
}

var _ = AfterEach(func() {
	td := CurrentGinkgoTestDescription()
	if td.Failed {
		fmt.Fprintf(GinkgoWriter, "\nLast Executed Smart Contract Line for %s:%d\n", td.FileName, td.LineNumber)
		fmt.Fprintln(GinkgoWriter, TestRig.LastExecuted())
	}
	err := Backend.Close()
	Expect(err).ToNot(HaveOccurred())
})
