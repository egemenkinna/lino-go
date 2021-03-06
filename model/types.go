package model

type Permission int
type DetailType int

const (
	InvalidSeqErrCode = 154

	// Different permission levels
	UnknownPermission          = Permission(0)
	AppPermission              = Permission(1)
	TransactionPermission      = Permission(2)
	ResetPermission            = Permission(3)
	GrantAppPermission         = Permission(4)
	PreAuthorizationPermission = Permission(5)

	// Different possible incomes
	TransferIn           = DetailType(0)
	DonationIn           = DetailType(1)
	ClaimReward          = DetailType(2)
	ValidatorInflation   = DetailType(3)
	DeveloperInflation   = DetailType(4)
	InfraInflation       = DetailType(5)
	VoteReturnCoin       = DetailType(6)
	DelegationReturnCoin = DetailType(7)
	ValidatorReturnCoin  = DetailType(8)
	DeveloperReturnCoin  = DetailType(9)
	InfraReturnCoin      = DetailType(10)
	ProposalReturnCoin   = DetailType(11)
	GenesisCoin          = DetailType(12)
	ClaimInterest        = DetailType(13)

	// Different possible outcomes
	TransferOut      = DetailType(20)
	DonationOut      = DetailType(21)
	Delegate         = DetailType(22)
	VoterDeposit     = DetailType(23)
	ValidatorDeposit = DetailType(24)
	DeveloperDeposit = DetailType(25)
	InfraDeposit     = DetailType(26)
	ProposalDeposit  = DetailType(27)
)
