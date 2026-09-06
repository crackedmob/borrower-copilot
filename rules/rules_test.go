package rules
 
import (
	"fmt"
	"testing"
)
 
func TestPersonas(t *testing.T) {
	priya := Input{
		NetMonthlyIncome: 110000, IncomeType: IncomeSalaried,
		ExistingEMIs: 14000, HouseholdExpenses: 28000,
		RequestedAmount: 800000, Purpose: PurposeWedding, Age: 29,
		CreditScore: CreditScore{Known: true, Score: 780},
		IncomeStableYears: 5, EmergencySavingsMonths: 2,
		QuestionsAnswered: 4, QuestionsApplicable: 6,
	}
	ravi := Input{
		NetMonthlyIncome: 35000, UndeclaredCashIncome: 25000,
		IncomeType: IncomeSelfEmployed, ExistingEMIs: 0, HouseholdExpenses: 20000,
		RequestedAmount: 1500000, Purpose: PurposeBusiness, Age: 42,
		CreditScore: CreditScore{Known: false},
		HasCollateral: true, CollateralValue: 4500000,
		HasCoApplicant: true, CoApplicantIncome: 18000, EmergencySavingsMonths: 1,
		QuestionsAnswered: 4, QuestionsApplicable: 6,
	}
	anita := Input{
		NetMonthlyIncome: 28000, IncomeType: IncomeInformal,
		ExistingEMIs: 8000, HouseholdExpenses: 15000,
		RequestedAmount: 150000, Purpose: PurposeVehicle, Age: 35,
		CreditScore: CreditScore{Known: false}, RecentBounce: true,
		EmergencySavingsMonths: 0,
		QuestionsAnswered: 2, QuestionsApplicable: 6,
	}
 
	for _, tc := range []struct {
		name string
		in   Input
	}{{"Priya", priya}, {"Ravi", ravi}, {"Anita", anita}} {
		out := Calculate(tc.in)
		fmt.Printf("\n=== %s ===\nverdict=%s product=%s\nlenderSanction=%.0f safeAmount=%.0f recommended=%.0f\nrate=%.2f-%.2f%% APR=%.2f-%.2f%% tenure=%.1fy\nEMIceiling=%.0f stressEMI(+2pp)=%.0f stressIncomeDropShare=%.1f%%\nprojectedFOIR=%.1f%% confidenceLow=%v\n",
			tc.name, out.Verdict, out.Product,
			out.LenderSanctionAmount, out.BorrowerSafeAmount, out.RecommendedAmount,
			out.RateLow, out.RateHigh, out.APRLow, out.APRHigh, out.TenureYears,
			out.EMICeiling, out.StressEMI, out.StressIncomeDrop,
			out.ProjectedFOIR, out.ConfidenceLow)
		for _, e := range out.Explanations {
			fmt.Println("  -", e)
		}
	}
}