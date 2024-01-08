package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

//Policy[1]->Lines[*]->Risk[*]->Coverages[*]
//Policy[1]->Lines[*]->Coverages[*]
//Policy[1]->Transactions[*]

type Transaction struct {
	Id             string
	TypeCode       string
	ImageType      string
	Status         string
	State          string
	EffectiveDate  time.Time
	CreatedDate    time.Time
	ExpirationDate time.Time
	IssuedDate     time.Time
	TermFactor     float64
}

type Coverage struct {
	indicator           bool
	included            bool
	deleted             bool
	modified            bool
	deductible          int
	level               string // level-[Risk/Line]
	coverageType        string //coverageType-[Premium,Tax,Fees,Surcharge],
	code                string
	formulae            string
	statCode            string
	currentLimit        string
	priorLimit          string
	baseRate            string
	ilf                 string
	addedDate           time.Time
	effectiveDate       time.Time
	expirationDate      time.Time
	deletedDate         time.Time
	termPremium         decimal.Decimal
	priorTermPremium    decimal.Decimal
	changePremium       decimal.Decimal
	writtenPremium      decimal.Decimal
	priorWrittenPremium decimal.Decimal
	termFactor          decimal.Decimal
	drf                 decimal.Decimal
	stringBuilder       strings.Builder
	//WrittenPremium=Prior Written + Term Charge {[Term Premium -Prior Term Premium] *DRF}
	//For NB , Prior Written & Prior Term Premium[WrittenPremium-ChangePremium] = 0
}

type Risk struct {
	Id                   string
	Indicator            bool
	Included             bool
	Deleted              bool
	GUID                 string
	AssociatedLocationID string
	StringBuilderRisk    strings.Builder
	coverages            []Coverage
}

type Line struct {
	Id        string
	TypeLOB   string
	coverages []Coverage
}

type Policy struct {
	number         string //GUID  (UGH!)
	effectiveDate  time.Time
	expirationDate time.Time
	lines          []Line
	transactions   []Transaction
}

type PolicyOption func(*Policy)

func NewPolicy(number string, opts ...PolicyOption) *Policy {
	p := &Policy{
		number:       number,
		lines:        make([]Line, 0),
		transactions: make([]Transaction, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
func SetEffectiveDate(d time.Time) PolicyOption {
	return func(p *Policy) {
		p.effectiveDate = d
	}
}

func SetExpirationDate(d time.Time) PolicyOption {
	return func(p *Policy) {
		p.expirationDate = d
	}
}

const shortFormTime = "2006-01-02"

var LineObject Line //Collection of Line Objects

var TransactionObject Transaction //Collection of Transaction Objects

func main() {
	now := time.Now()
	aPolicy := NewPolicy("PackagePolicy1P2024",
		SetEffectiveDate(now),
		SetExpirationDate(addMonth(now, 12)),
		WithLines(LineObject),
		WithTransactions(TransactionObject))

	fmt.Println("Policy Number :- " + aPolicy.number)
	fmt.Println("Policy Effective Date :- " + aPolicy.effectiveDate.Format(shortFormTime))
	fmt.Println("Policy Expiration Date :- " + aPolicy.expirationDate.Format(shortFormTime))
	for _, line := range aPolicy.lines {
		fmt.Println("Line ID - " + line.Id + " Type - " + line.TypeLOB)
	}
	for _, transaction := range aPolicy.transactions {
		fmt.Println("Transaction ID - " + transaction.Id + " Type - " + transaction.TypeCode)
	}

}

func addMonth(t time.Time, month int) time.Time {
	return t.AddDate(0, month, 0)
}

func WithTransactions(Transactions Transaction) PolicyOption {
	return func(p *Policy) {
		transactionscsv := readTransactionObjectsCSVFile("Objects/Transaction.csv")
		for _, transaction := range transactionscsv {
			//	fmt.Println("Line ID -" + line.TypeLOB)
			//	p.lines = append(p.lines, line)
			p.transactions = append(p.transactions, transaction)

		}
	}
}

func WithLines(lines Line) PolicyOption {
	return func(p *Policy) {
		linescsv := readLineObjectsCSVFile("Objects/Line.csv")
		for _, line := range linescsv {
			//	fmt.Println("Line ID -" + line.TypeLOB)
			p.lines = append(p.lines, line)

		}
	}
}

func readLineObjectsCSVFile(filePath string) (lines []Line) {
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()
		defer f.Close()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		checkError("Some other error occurred", err)

		// Handle first row case
		if isFirstRow {
			isFirstRow = false

			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}

			// Skip next code
			continue
		}

		// Create new coverage and add to persons array
		lines = append(lines, Line{
			TypeLOB: record[headerMap["Type"]],
			Id:      record[headerMap["ID"]],
		})
	}
	return

}
func readTransactionObjectsCSVFile(filePath string) (transactions []Transaction) {
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()
		defer f.Close()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		checkError("Some other error occurred", err)

		// Handle first row case
		if isFirstRow {
			isFirstRow = false

			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}

			// Skip next code
			continue
		}

		// Create new coverage and add to persons array
		transactions = append(transactions, Transaction{
			TypeCode: record[headerMap["Type"]],
			Id:       record[headerMap["ID"]],
		})
	}
	return

}

func checkError(message string, err error) {
	// Error Logging
	if err != nil {
		//log.Fatal(message, err)
		//log := zerolog.New(os.Stdout).With().Logger()
		//log.Debug().Str(" checkError %s ", message)
	}
}
