
[1] - Singleton
[*] - Multi Instance

//Policy[1]->Lines[*]->Risk[*]->Coverages[*]
//Policy[1]->Lines[*]->Coverages[*]
//Policy[1]->Transactions[*]

Completed:
-Policy Struct Initialize
-Add Lines/Transactions Dynamically



To be done:
-Add Risks to Line and print them
Risk1andRisk2 should be appended to Line1 based on LineID column
Risk3andRisk4 should be appended to Line2 based on LineID column
Refer-AutoRiskData.csv-
Name,ID,Include,Deleted,Limit,LineID,LocationID
Risk1,1,1,0,25,1,Location1
Risk2,2,1,0,22,1,Location1
Risk3,3,1,0,21,2,Location2
Risk4,4,1,0,15,2,Location2

Ignore below:Self Notes

Constants
const shortFormTime = "2006-01-02"

main.go
//var LineObject Line //Collection of Line Objects
//var TransactionObject Transaction //Collection of Transaction Objects

coverage
	//WrittenPremium=Prior Written + Term Charge {[Term Premium -Prior Term Premium] *DRF}
	//For NB , Prior Written & Prior Term Premium[WrittenPremium-ChangePremium] = 0


Append lines
    func WithLines(lines []*Line) PolicyOption { //func WithLines(lines Line) PolicyOption {
	return func(p *Policy) {
		//	linescsv := readLineObjectsCSVFile("Objects/Line.csv")
		for _, line := range lines { //for _, line := range linescsv {
			//	fmt.Println("Line ID -" + line.TypeLOB)
			p.lines = append(p.lines, line)

		}
	}
}


