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

const shortFormTime = "2006-01-02"

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
}

type Line struct {
	Id        string
	TypeLOB   string
	coverages []*Coverage //	coverages []Coverage
	risks     []*Risk
}

type Risk struct {
	Id                   string
	Indicator            bool
	Included             bool
	Deleted              bool
	GUID                 string
	AssociatedLocationID string
	StringBuilderRisk    strings.Builder
	coverages            []*Coverage //coverages []Coverage
}

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

type Location struct {
	Id string
}
type Policy struct {
	number         string //GUID  (UGH!)
	inceptionDate  time.Time
	effectiveDate  time.Time
	expirationDate time.Time
	lines          []*Line        //lines          []Line
	transactions   []*Transaction //transactions   []Transaction
}

type PolicyOption func(*Policy)

type LineOption func(*Line)

func NewPolicy(number string, opts ...PolicyOption) *Policy {
	p := &Policy{
		number:       number,
		lines:        make([]*Line, 0),        //lines:        make([]Line, 0),
		transactions: make([]*Transaction, 0), //	transactions: make([]Transaction, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func SetInceptionDate(d time.Time) PolicyOption {
	return func(p *Policy) {
		p.inceptionDate = d
	}
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

func WithTransactions(transactions []*Transaction) PolicyOption {
	return func(p *Policy) {
		for _, transaction := range transactions {
			p.transactions = append(p.transactions, transaction)
		}
	}
}

func WithLines(lines []*Line) PolicyOption {
	return func(p *Policy) {
		for _, line := range lines {
			p.lines = append(p.lines, line)
		}
	}
}

func WithLinesRisk(risks []*Risk) LineOption {
	return func(line *Line) {
		//	for _, risk := range risks {
		//		p.lines = append(p.lines, risk)
		//	}
	}
}

func main() {
	lines, lineerr := readLineObjectsCSVFile("Objects/Line.csv")
	transactions, txerr := readTransactionObjectsCSVFile("Objects/Transaction.csv")
	risks, riskerr := readRiskObjectsCSVFile("Objects/AutoRiskData.csv")
	for _, risk := range risks {
		fmt.Println("Risk ID - " + risk.Id)
	}
	if lineerr != nil {
		panic(lineerr)
	}
	if txerr != nil {
		panic(txerr)
	}
	if riskerr != nil {
		panic(riskerr)
	}
	now := time.Now()
	aPolicy := NewPolicy("PackagePolicy1P2024",
		SetInceptionDate(now),
		SetEffectiveDate(now),
		SetExpirationDate(addMonth(now, 12)),
		WithLines(lines),               //WithLines(LineObject),
		WithTransactions(transactions)) //WithTransactions(TransactionObject)

	fmt.Println("Policy Number :- " + aPolicy.number)
	fmt.Println("Policy Inception Date :- " + aPolicy.inceptionDate.Format(shortFormTime))
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

func readLineObjectsCSVFile(filePath string) (lines []*Line, err error) { //func readLineObjectsCSVFile(filePath string) (lines []Line) {
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)
	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		//	checkError("Some other error occurred", err)

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
		lines = append(lines, &Line{ //	lines = append(lines, Line{
			TypeLOB: record[headerMap["Type"]],
			Id:      record[headerMap["ID"]],
		})
	}
	return lines, err
	//	return

}

func readRiskObjectsCSVFile(filePath string) (risks []*Risk, err error) { //func readLineObjectsCSVFile(filePath string) (lines []Line) {
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)
	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		//	checkError("Some other error occurred", err)

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
		risks = append(risks, &Risk{ //	lines = append(lines, Line{
			GUID: record[headerMap["Name"]],
			Id:   record[headerMap["ID"]],
		})
	}
	return risks, err
	//	return

}

func readTransactionObjectsCSVFile(filePath string) (transactions []*Transaction, err error) { //func readTransactionObjectsCSVFile(filePath string) (transactions []Transaction)
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)
	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		//	checkError("Some other error occurred", err)

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
		transactions = append(transactions, &Transaction{ //	transactions = append(transactions, Transaction{
			TypeCode: record[headerMap["Type"]],
			Id:       record[headerMap["ID"]],
		})
	}
	return transactions, err
	//return

}

/*

func checkError(message string, err error) {
	// Error Logging
	if err != nil {
		//log.Fatal(message, err)
		//log := zerolog.New(os.Stdout).With().Logger()
		//log.Debug().Str(" checkError %s ", message)
	}
} */
