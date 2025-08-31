package enums

type LoanStatus int

const (
	LoanStatusProposed LoanStatus = iota + 1
	LoanStatusApproved
	LoanStatusRejected
	LoanStatusInvested
	LoanStatusDisbursed
)

func (ls LoanStatus) String() string {
	switch ls {
	case LoanStatusProposed:
		return "PROPOSED"
	case LoanStatusApproved:
		return "APPROVED"
	case LoanStatusRejected:
		return "REJECTED"
	case LoanStatusInvested:
		return "INVESTED"
	case LoanStatusDisbursed:
		return "DISBURSED"
	default:
		return "UNKNOWN"
	}
}

func (ls LoanStatus) Int() int {
	return int(ls)
}

func LoanStatusFromString(status string) LoanStatus {
	switch status {
	case "PROPOSED":
		return LoanStatusProposed
	case "APPROVED":
		return LoanStatusApproved
	case "REJECTED":
		return LoanStatusRejected
	case "INVESTED":
		return LoanStatusInvested
	case "DISBURSED":
		return LoanStatusDisbursed
	default:
		return LoanStatusProposed
	}
}

func LoanStatusFromInt(status int) LoanStatus {
	switch status {
	case 1:
		return LoanStatusProposed
	case 2:
		return LoanStatusApproved
	case 3:
		return LoanStatusRejected
	case 4:
		return LoanStatusInvested
	case 5:
		return LoanStatusDisbursed
	default:
		return LoanStatusProposed
	}
}

func GetAllLoanStatuses() []LoanStatus {
	return []LoanStatus{
		LoanStatusProposed,
		LoanStatusApproved,
		LoanStatusRejected,
		LoanStatusInvested,
		LoanStatusDisbursed,
	}
}

func GetLoanStatusMap() map[int]string {
	return map[int]string{
		1: "PROPOSED",
		2: "APPROVED",
		3: "REJECTED",
		4: "INVESTED",
		5: "DISBURSED",
	}
}
