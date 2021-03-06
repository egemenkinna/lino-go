// Pacakge broadcast includes the functionalities to broadcast
// all kinds of transactions to blockchain.
package broadcast

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/lino-network/lino-go/errors"
	"github.com/lino-network/lino-go/model"
	"github.com/lino-network/lino-go/transport"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Broadcast is a wrapper of broadcasting transactions to blockchain.
type Broadcast struct {
	transport *transport.Transport
}

// NewBroadcast returns an instance of Broadcast.
func NewBroadcast(transport *transport.Transport) *Broadcast {
	return &Broadcast{
		transport: transport,
	}
}

//
// Account related tx
//

// Register registers a new user on blockchain.
// It composes RegisterMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Register(ctx context.Context, referrer, registerFee, username, resetPubKeyHex,
	transactionPubKeyHex, appPubKeyHex, referrerPrivKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	resetPubKey, err := transport.GetPubKeyFromHex(resetPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHex("Register: failed to get Reset pub key").AddCause(err)
	}
	txPubKey, err := transport.GetPubKeyFromHex(transactionPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHex("Register: failed to get Tx pub key").AddCause(err)
	}
	appPubKey, err := transport.GetPubKeyFromHex(appPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHex("Register: failed to get App pub key").AddCause(err)
	}

	msg := model.RegisterMsg{
		Referrer:             referrer,
		RegisterFee:          registerFee,
		NewUser:              username,
		NewResetPubKey:       resetPubKey,
		NewTransactionPubKey: txPubKey,
		NewAppPubKey:         appPubKey,
	}
	return broadcast.broadcastTransaction(ctx, msg, referrerPrivKeyHex, seq, "", false)
}

// Transfer sends a certain amount of LINO token from the sender to the receiver.
// It composes TransferMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Transfer(ctx context.Context, sender, receiver, amount, memo,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.TransferMsg{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
		Memo:     memo,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Follow creates a social relationship between follower and followee.
// It composes FollowMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Follow(ctx context.Context, follower, followee,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.FollowMsg{
		Follower: follower,
		Followee: followee,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Unfollow revokes the social relationship between follower and followee.
// It composes UnfollowMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Unfollow(ctx context.Context, follower, followee,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.UnfollowMsg{
		Follower: follower,
		Followee: followee,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Claim claims rewards of a certain user.
// It composes ClaimMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Claim(ctx context.Context, username,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ClaimMsg{
		Username: username,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// UpdateAccount updates account related info in jsonMeta which are not
// included in AccountInfo or AccountBank.
// It composes UpdateAccountMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) UpdateAccount(ctx context.Context, username, jsonMeta,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.UpdateAccountMsg{
		Username: username,
		JSONMeta: jsonMeta,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Recover recovers all keys of a user in case of losing or compromising.
// It composes RecoverMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Recover(ctx context.Context, username, newResetPubKeyHex,
	newTransactionPubKeyHex, newAppPubKeyHex, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	resetPubKey, err := transport.GetPubKeyFromHex(newResetPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHexf("Recover: failed to get Reset pub key").AddCause(err)
	}
	txPubKey, err := transport.GetPubKeyFromHex(newTransactionPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHexf("Recover: failed to get Tx pub key").AddCause(err)
	}
	appPubKey, err := transport.GetPubKeyFromHex(newAppPubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHexf("Recover: failed to get App pub key").AddCause(err)
	}

	msg := model.RecoverMsg{
		Username:             username,
		NewResetPubKey:       resetPubKey,
		NewTransactionPubKey: txPubKey,
		NewAppPubKey:         appPubKey,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// Post related tx
//

// CreatePost creates a new post on blockchain.
// It composes CreatePostMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) CreatePost(ctx context.Context, author, postID, title, content,
	parentAuthor, parentPostID, sourceAuthor, sourcePostID, redistributionSplitRate string,
	links map[string]string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	var mLinks []model.IDToURLMapping
	if links == nil || len(links) == 0 {
		mLinks = nil
	} else {
		for k, v := range links {
			mLinks = append(mLinks, model.IDToURLMapping{k, v})
		}
	}

	msg := model.CreatePostMsg{
		Author:       author,
		PostID:       postID,
		Title:        title,
		Content:      content,
		ParentAuthor: parentAuthor,
		ParentPostID: parentPostID,
		SourceAuthor: sourceAuthor,
		SourcePostID: sourcePostID,
		Links:        mLinks,
		RedistributionSplitRate: redistributionSplitRate,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// CreatePost creates a new post on blockchain.
// It composes CreatePostMsg and then broadcasts the transaction to blockchain return when checkTx pass.
func (broadcast *Broadcast) CreatePostSync(ctx context.Context, author, postID, title, content,
	parentAuthor, parentPostID, sourceAuthor, sourcePostID, redistributionSplitRate string,
	links map[string]string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	var mLinks []model.IDToURLMapping
	if links == nil || len(links) == 0 {
		mLinks = nil
	} else {
		for k, v := range links {
			mLinks = append(mLinks, model.IDToURLMapping{k, v})
		}
	}

	msg := model.CreatePostMsg{
		Author:       author,
		PostID:       postID,
		Title:        title,
		Content:      content,
		ParentAuthor: parentAuthor,
		ParentPostID: parentPostID,
		SourceAuthor: sourceAuthor,
		SourcePostID: sourcePostID,
		Links:        mLinks,
		RedistributionSplitRate: redistributionSplitRate,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", true)
}

// Donate adds a money donation to a post by a user.
// It composes DonateMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Donate(ctx context.Context, username, author,
	amount, postID, fromApp, memo string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DonateMsg{
		Username: username,
		Amount:   amount,
		Author:   author,
		PostID:   postID,
		FromApp:  fromApp,
		Memo:     memo,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Donate adds a money donation to a post by a user.
// It composes DonateMsg and then broadcasts the transaction to blockchain return after pass checkTx.
func (broadcast *Broadcast) DonateSync(ctx context.Context, username, author,
	amount, postID, fromApp, memo string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DonateMsg{
		Username: username,
		Amount:   amount,
		Author:   author,
		PostID:   postID,
		FromApp:  fromApp,
		Memo:     memo,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", true)
}

// ReportOrUpvote adds a report or upvote action to a post.
// It composes ReportOrUpvoteMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ReportOrUpvote(ctx context.Context, username, author,
	postID string, isReport bool, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ReportOrUpvoteMsg{
		Username: username,
		Author:   author,
		PostID:   postID,
		IsReport: isReport,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// DeletePost deletes a post from the blockchain. It doesn't actually
// remove the post from the blockchain, instead it sets IsDeleted to true
// and clears all the other data.
// It composes DeletePostMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DeletePost(ctx context.Context, author, postID,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DeletePostMsg{
		Author: author,
		PostID: postID,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// View increases the view count of a post by one.
// It composes ViewMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) View(ctx context.Context, username, author, postID,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ViewMsg{
		Username: username,
		Author:   author,
		PostID:   postID,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// UpdatePost updates post info with new data.
// It composes UpdatePostMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) UpdatePost(ctx context.Context, author, title, postID, content string,
	links map[string]string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	var mLinks []model.IDToURLMapping
	if links == nil || len(links) == 0 {
		mLinks = nil
	} else {
		for k, v := range links {
			mLinks = append(mLinks, model.IDToURLMapping{k, v})
		}
	}

	msg := model.UpdatePostMsg{
		Author:  author,
		PostID:  postID,
		Title:   title,
		Content: content,
		Links:   mLinks,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// Validator related tx
//

// ValidatorDeposit deposits a certain amount of LINO token for a user
// in order to become a validator. Before becoming a validator, the user
// has to be a voter.
// It composes ValidatorDepositMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ValidatorDeposit(ctx context.Context, username, deposit,
	validatorPubKey, link, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	valPubKey, err := transport.GetPubKeyFromHex(validatorPubKey)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHexf("ValidatorDeposit: failed to get Val pub key").AddCause(err)
	}
	msg := model.ValidatorDepositMsg{
		Username:  username,
		Deposit:   deposit,
		ValPubKey: valPubKey,
		Link:      link,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ValidatorWithdraw withdraws part of LINO token from a validator's deposit,
// while still keep being a validator.
// It composes ValidatorDepositMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ValidatorWithdraw(ctx context.Context, username, amount,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ValidatorWithdrawMsg{
		Username: username,
		Amount:   amount,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ValidatorRevoke revokes all deposited LINO token of a validator
// so that the user will not be a validator anymore.
// It composes ValidatorRevokeMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ValidatorRevoke(ctx context.Context, username,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ValidatorRevokeMsg{
		Username: username,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// Vote related tx
//

// StakeIn deposits a certain amount of LINO token for a user
// in order to become a voter.
// It composes StakeInMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) StakeIn(ctx context.Context, username, deposit,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.StakeInMsg{
		Username: username,
		Deposit:  deposit,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// StakeOut withdraws part of LINO token from a voter's deposit.
// It composes StakeOutMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) StakeOut(ctx context.Context, username, amount,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.StakeOutMsg{
		Username: username,
		Amount:   amount,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// Delegate delegates a certain amount of LINO token of delegator to a voter, so
// the voter will have more voting power.
// It composes DelegateMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) Delegate(ctx context.Context, delegator, voter, amount,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DelegateMsg{
		Delegator: delegator,
		Voter:     voter,
		Amount:    amount,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// DelegatorWithdraw withdraws part of delegated LINO token of a delegator
// to a voter, while the delegation still exists.
// It composes DelegatorWithdrawMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DelegatorWithdraw(ctx context.Context, delegator, voter, amount,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DelegatorWithdrawMsg{
		Delegator: delegator,
		Voter:     voter,
		Amount:    amount,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ClaimInterest claims interest of a certain user.
// It composes ClaimInterestMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ClaimInterest(ctx context.Context, username,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ClaimInterestMsg{
		Username: username,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// Developer related tx
//

// DeveloperRegsiter registers a developer with a certain amount of LINO token on blockchain.
// It composes DeveloperRegisterMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DeveloperRegister(ctx context.Context, username, deposit, website,
	description, appMetaData, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DeveloperRegisterMsg{
		Username:    username,
		Deposit:     deposit,
		Website:     website,
		Description: description,
		AppMetaData: appMetaData,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// DeveloperUpdate updates a developer  info on blockchain.
// It composes DeveloperUpdateMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DeveloperUpdate(ctx context.Context, username, website,
	description, appMetaData, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DeveloperUpdateMsg{
		Username:    username,
		Website:     website,
		Description: description,
		AppMetaData: appMetaData,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// DeveloperRevoke reovkes all deposited LINO token of a developer
// so the user will not be a developer anymore.
// It composes DeveloperRevokeMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DeveloperRevoke(ctx context.Context, username,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.DeveloperRevokeMsg{
		Username: username,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// GrantPermission grants a certain (e.g. App) permission to
// an authorized app with a certain period of time.
// It composes GrantPermissionMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) GrantPermission(ctx context.Context, username, authorizedApp string,
	validityPeriodSec int64, grantLevel model.Permission, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.GrantPermissionMsg{
		Username:          username,
		AuthorizedApp:     authorizedApp,
		ValidityPeriodSec: validityPeriodSec,
		GrantLevel:        grantLevel,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// PreAuthorizationPermission grants a PreAuthorization permission to
// an authorzied app with a certain period of time.
// It composes PreAuthorizationMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) PreAuthorizationPermission(ctx context.Context, username, authorizedApp string,
	validityPeriodSec int64, amount string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.PreAuthorizationMsg{
		Username:          username,
		AuthorizedApp:     authorizedApp,
		ValidityPeriodSec: validityPeriodSec,
		Amount:            amount,
	}

	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// RevokePermission revokes the permission given previously to a app.
// It composes RevokePermissionMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) RevokePermission(ctx context.Context, username, pubKeyHex string,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	pubKey, err := transport.GetPubKeyFromHex(pubKeyHex)
	if err != nil {
		return nil, errors.FailedToGetPubKeyFromHex("Register: failed to get pub key").AddCause(err)
	}

	msg := model.RevokePermissionMsg{
		Username: username,
		PubKey:   pubKey,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// infra related tx
//

// ProviderReport reports infra usage of a infra provider in order to get infra inflation.
// It composes ProviderReportMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ProviderReport(ctx context.Context, username string, usage int64,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ProviderReportMsg{
		Username: username,
		Usage:    usage,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// proposal related tx
//

// ChangeEvaluateOfContentValueParam changes EvaluateOfContentValueParam with new value.
// It composes ChangeEvaluateOfContentValueParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeEvaluateOfContentValueParam(ctx context.Context, creator string,
	parameter model.EvaluateOfContentValueParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeEvaluateOfContentValueParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeGlobalAllocationParam changes GlobalAllocationParam with new value.
// It composes ChangeGlobalAllocationParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeGlobalAllocationParam(ctx context.Context, creator string,
	parameter model.GlobalAllocationParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeGlobalAllocationParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeInfraInternalAllocationParam changes InfraInternalAllocationParam with new value.
// It composes ChangeInfraInternalAllocationParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeInfraInternalAllocationParam(ctx context.Context, creator string,
	parameter model.InfraInternalAllocationParam,
	reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeInfraInternalAllocationParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeVoteParam changes VoteParam with new value.
// It composes ChangeVoteParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeVoteParam(ctx context.Context, creator string,
	parameter model.VoteParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeVoteParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeProposalParam changes ProposalParam with new value.
// It composes ChangeProposalParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeProposalParam(ctx context.Context, creator string,
	parameter model.ProposalParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeProposalParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeDeveloperParam changes DeveloperParam with new value.
// It composes ChangeDeveloperParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeDeveloperParam(ctx context.Context, creator string,
	parameter model.DeveloperParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeDeveloperParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeValidatorParam changes ValidatorParam with new value.
// It composes ChangeValidatorParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeValidatorParam(ctx context.Context, creator string,
	parameter model.ValidatorParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeValidatorParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeBandwidthParam changes BandwidthParam with new value.
// It composes ChangeBandwidthParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeBandwidthParam(ctx context.Context, creator string,
	parameter model.BandwidthParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeBandwidthParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangeAccountParam changes AccountParam with new value.
// It composes ChangeAccountParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangeAccountParam(ctx context.Context, creator string,
	parameter model.AccountParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangeAccountParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// ChangePostParam changes PostParam with new value.
// It composes ChangePostParamMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) ChangePostParam(ctx context.Context, creator string,
	parameter model.PostParam, reason string, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.ChangePostParamMsg{
		Creator:   creator,
		Parameter: parameter,
		Reason:    reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// DeletePostContent deletes the content of a post on blockchain, which is used
// for content censorship.
// It composes DeletePostContentMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) DeletePostContent(ctx context.Context, creator, postAuthor,
	postID, reason, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	permlink := string(string(postAuthor) + "#" + postID)
	msg := model.DeletePostContentMsg{
		Creator:  creator,
		Permlink: permlink,
		Reason:   reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// VoteProposal adds a vote to a certain proposal with agree/disagree.
// It composes VoteProposalMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) VoteProposal(ctx context.Context, voter, proposalID string,
	result bool, privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.VoteProposalMsg{
		Voter:      voter,
		ProposalID: proposalID,
		Result:     result,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

// UpgradeProtocol upgrades the protocol.
// It composes UpgradeProtocolMsg and then broadcasts the transaction to blockchain.
func (broadcast *Broadcast) UpgradeProtocol(ctx context.Context, creator, link, reason string,
	privKeyHex string, seq int64) (*model.BroadcastResponse, error) {
	msg := model.UpgradeProtocolMsg{
		Creator: creator,
		Link:    link,
		Reason:  reason,
	}
	return broadcast.broadcastTransaction(ctx, msg, privKeyHex, seq, "", false)
}

//
// internal helper functions
//
func (broadcast *Broadcast) broadcastTransaction(ctx context.Context, msg model.Msg, privKeyHex string,
	seq int64, memo string, checkTxOnly bool) (*model.BroadcastResponse, error) {
	broadcastResp := &model.BroadcastResponse{}

	var res interface{}
	var err error
	finishChan := make(chan bool)
	go func() {
		res, err = broadcast.transport.SignBuildBroadcast(msg, privKeyHex, seq, memo, checkTxOnly)
		finishChan <- true
	}()

	select {
	case <-finishChan:
		break
	case <-ctx.Done():
		return nil, errors.Timeoutf("msg timeout: %v", msg).AddCause(ctx.Err())
	}

	if err != nil {
		return nil, errors.FailedToBroadcast(err.Error())
	}

	if checkTxOnly {
		res, ok := res.(*ctypes.ResultBroadcastTx)
		if !ok {
			return nil, errors.FailedToBroadcast("error to parse the broadcast response")
		}
		code := retrieveCodeFromBlockChainCode(res.Code)
		if err == nil && code == model.InvalidSeqErrCode {
			return nil, errors.InvalidSequenceNumber("invalid seq").AddBlockChainCode(res.Code).AddBlockChainLog(res.Log)
		}

		if res.Code != uint32(0) {
			return nil, errors.CheckTxFail("CheckTx failed!").AddBlockChainCode(res.Code).AddBlockChainLog(res.Log)
		}
		if res.Code != uint32(0) {
			return nil, errors.DeliverTxFail("DeliverTx failed!").AddBlockChainCode(res.Code).AddBlockChainLog(res.Log)
		}
		commitHash := hex.EncodeToString(res.Hash)
		broadcastResp.CommitHash = strings.ToUpper(commitHash)
	} else {
		res, ok := res.(*ctypes.ResultBroadcastTxCommit)
		if !ok {
			return nil, errors.FailedToBroadcast("error to parse the broadcast response")
		}
		code := retrieveCodeFromBlockChainCode(res.CheckTx.Code)
		if err == nil && code == model.InvalidSeqErrCode {
			return nil, errors.InvalidSequenceNumber("invalid seq").AddBlockChainCode(res.CheckTx.Code).AddBlockChainLog(res.CheckTx.Log)
		}

		if res.CheckTx.Code != uint32(0) {
			return nil, errors.CheckTxFail("CheckTx failed!").AddBlockChainCode(res.CheckTx.Code).AddBlockChainLog(res.CheckTx.Log)
		}
		if res.DeliverTx.Code != uint32(0) {
			return nil, errors.DeliverTxFail("DeliverTx failed!").AddBlockChainCode(res.DeliverTx.Code).AddBlockChainLog(res.DeliverTx.Log)
		}
		commitHash := hex.EncodeToString(res.Hash)
		broadcastResp.CommitHash = strings.ToUpper(commitHash)
	}

	return broadcastResp, nil
}

func retrieveCodeFromBlockChainCode(bcCode uint32) uint32 {
	return bcCode & 0xff
}
