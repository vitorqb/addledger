package printer_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	. "github.com/vitorqb/addledger/internal/printer"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

func RunTest(
	t *testing.T,
	name string,
	transaction journal.Transaction,
	numLineBreaksBefore int,
	numLineBreaksAfter int,
	expectedOutput string,
) {
	t.Run(name, func(t *testing.T) {
		var buf bytes.Buffer
		printerInstance := New(numLineBreaksBefore, numLineBreaksAfter)
		err := printerInstance.Print(&buf, transaction)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedOutput, buf.String())

	})
}

func TestPrinter_Print(t *testing.T) {

	simpleTransaction := *tu.Transaction_1(t)
	RunTest(
		t,
		"Simple",
		simpleTransaction,
		0,
		0,
		"1993-11-23 Description1\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2",
	)

	RunTest(
		t,
		"Simple (with empty lines)",
		simpleTransaction,
		2,
		1,
		"\n\n1993-11-23 Description1\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2\n",
	)

	noCommodityTransaction := *tu.Transaction_1(t)
	noCommodityTransaction.Posting[0].Ammount.Commodity = ""
	RunTest(
		t,
		"No commodity",
		noCommodityTransaction,
		0,
		0,
		"1993-11-23 Description1\n    ACC1    12.2\n    ACC2    EUR -12.2",
	)

	fourPostingsTransaction := *tu.Transaction_1(t)
	fourPostingsTransaction.Posting = append(fourPostingsTransaction.Posting, fourPostingsTransaction.Posting[0], fourPostingsTransaction.Posting[1])
	RunTest(
		t,
		"Four postings",
		fourPostingsTransaction,
		0,
		0,
		"1993-11-23 Description1\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2",
	)

}
