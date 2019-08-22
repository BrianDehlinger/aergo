package system

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/aergoio/aergo/state"
	"github.com/aergoio/aergo/types"
)

var ErrTxSystemOperatorIsNotSet = errors.New("operator is not set")

func ValidateSystemTx(account []byte, txBody *types.TxBody, sender *state.V,
	scs *state.ContractState, blockNo uint64) (*SystemContext, error) {
	var ci types.CallInfo
	if err := json.Unmarshal(txBody.Payload, &ci); err != nil {
		return nil, types.ErrTxInvalidPayload
	}
	context := &SystemContext{Call: &ci, Sender: sender, BlockNo: blockNo, op: types.GetOpSysTx(ci.Name), scs: scs, txBody: txBody}

	switch context.op {
	case types.Opstake:
		if sender != nil && sender.Balance().Cmp(txBody.GetAmountBigInt()) < 0 {
			return nil, types.ErrInsufficientBalance
		}
		staked, err := validateForStaking(account, txBody, scs, blockNo)
		if err != nil {
			return nil, err
		}
		context.Staked = staked
	case types.OpvoteBP:
		staked, oldvote, err := validateForVote(account, txBody, scs, blockNo, []byte(context.op.ID()))
		if err != nil {
			return nil, err
		}
		context.Staked = staked
		context.Vote = oldvote
	case types.Opunstake:
		staked, err := validateForUnstaking(account, txBody, scs, blockNo)
		if err != nil {
			return nil, err
		}
		context.Staked = staked
	case types.OpcreateProposal:
		staked, err := checkStakingBefore(account, scs)
		if err != nil {
			return nil, err
		}
		_, err = checkOperator(scs, sender.ID())
		if err != nil {
			return nil, err
		}
		id, err := parseIDForProposal(&ci)
		if err != nil {
			return nil, err
		}
		proposal, err := getProposal(scs, id)
		if err != nil {
			return nil, err
		}
		if proposal != nil {
			return nil, fmt.Errorf("already created proposal id: %s", proposal.ID)
		}
		if len(ci.Args) != 3 {
			return nil, fmt.Errorf("the request should be have 3 arguments: %d", len(ci.Args))
		}
		max, ok := ci.Args[1].(string)
		if !ok {
			return nil, fmt.Errorf("could not parse the max")
		}
		multipleChoice, err := strconv.ParseUint(max, 10, 32)
		if err != nil {
			return nil, err
		}
		desc, ok := ci.Args[2].(string)
		if !ok {
			return nil, fmt.Errorf("could not parse the desc")
		}
		context.Staked = staked
		context.Proposal = &Proposal{
			ID:             id,
			Blockfrom:      0,
			Blockto:        0,
			MultipleChoice: uint32(multipleChoice),
			Description:    desc,
		}
	case types.OpvoteProposal:
		id, err := parseIDForProposal(&ci)
		if err != nil {
			return nil, err
		}
		proposal, err := getProposal(scs, id)
		if err != nil {
			return nil, err
		}
		if proposal == nil {
			return nil, fmt.Errorf("the proposal is not created (%s)", id)
		}
		if blockNo < proposal.Blockfrom {
			return nil, fmt.Errorf("the voting begins at %d", proposal.Blockfrom)
		}
		if proposal.Blockto != 0 && blockNo > proposal.Blockto {
			return nil, fmt.Errorf("the voting was already done at %d", proposal.Blockto)
		}
		candis := ci.Args[1:]
		if int64(len(candis)) > int64(proposal.MultipleChoice) {
			return nil, fmt.Errorf("too many candidates arguments (max : %d)", proposal.MultipleChoice)
		}
		sort.Slice(proposal.Candidates, func(i, j int) bool {
			return proposal.Candidates[i] <= proposal.Candidates[j]
		})
		if len(proposal.Candidates) != 0 {
			for _, c := range candis {
				candidate, ok := c.(string)
				if !ok {
					return nil, fmt.Errorf("include invalid candidate")
				}
				i := sort.SearchStrings(proposal.Candidates, candidate)
				if i < len(proposal.Candidates) && proposal.Candidates[i] == candidate {
					//fmt.Printf("Found %s at index %d in %v.\n", x, i, a)
				} else {
					return nil, fmt.Errorf("candidate should be in %v", proposal.Candidates)
				}
			}
		}

		staked, oldvote, err := validateForVote(account, txBody, scs, blockNo, proposal.GetKey())
		if err != nil {
			return nil, err
		}
		context.Proposal = proposal
		context.Staked = staked
		context.Vote = oldvote
	case types.OpaddOperator,
		types.OpremoveOperator:
		if err := checkOperatorArg(context, &ci); err != nil {
			return nil, err
		}
		operators, err := checkOperator(scs, sender.ID())
		if err != nil &&
			err != ErrTxSystemOperatorIsNotSet {
			return nil, err
		}
		operatorAddr := types.ToAddress(context.Args[0])
		if context.op == types.OpaddOperator {
			if operators.IsExist(operatorAddr) {
				return nil, fmt.Errorf("already exist operator: %s", ci.Args[0])
			}
			operators = append(operators, operatorAddr)
		} else if context.op == types.OpremoveOperator {
			if !operators.IsExist(sender.ID()) {
				return nil, fmt.Errorf("operator is not exist : %s", ci.Args[0])
			}
			for i, v := range operators {
				if bytes.Equal(v, operatorAddr) {
					operators = append(operators[:i], operators[i+1:]...)
					break
				}
			}
		}
		context.Operators = operators
	default:
		return nil, types.ErrTxInvalidPayload
	}
	return context, nil
}

func checkStakingBefore(account []byte, scs *state.ContractState) (*types.Staking, error) {
	staked, err := getStaking(scs, account)
	if err != nil {
		return nil, err
	}
	if staked.GetAmountBigInt().Cmp(new(big.Int).SetUint64(0)) == 0 {
		return nil, fmt.Errorf("not staking before")
	}
	return staked, nil
}

func validateForStaking(account []byte, txBody *types.TxBody, scs *state.ContractState, blockNo uint64) (*types.Staking, error) {
	staked, err := getStaking(scs, account)
	if err != nil {
		return nil, err
	}
	if staked.GetAmount() != nil && staked.GetWhen()+StakingDelay > blockNo {
		return nil, types.ErrLessTimeHasPassed
	}
	toBe := new(big.Int).Add(staked.GetAmountBigInt(), txBody.GetAmountBigInt())
	stakingMin, err := getStakingMinimum(scs)
	if err != nil {
		return nil, err
	}
	if stakingMin.Cmp(toBe) > 0 {
		return nil, types.ErrTooSmallAmount
	}
	return staked, nil
}

func validateForVote(account []byte, txBody *types.TxBody, scs *state.ContractState, blockNo uint64, voteKey []byte) (*types.Staking, *types.Vote, error) {
	staked, err := checkStakingBefore(account, scs)
	if err != nil {
		return nil, nil, types.ErrMustStakeBeforeVote
	}
	oldvote, err := GetVote(scs, account, voteKey)
	if err != nil {
		return nil, nil, err
	}
	if oldvote.Amount != nil && staked.GetWhen()+VotingDelay > blockNo {
		return nil, nil, types.ErrLessTimeHasPassed
	}
	return staked, oldvote, nil
}

func validateForUnstaking(account []byte, txBody *types.TxBody, scs *state.ContractState, blockNo uint64) (*types.Staking, error) {
	staked, err := checkStakingBefore(account, scs)
	if err != nil {
		return nil, types.ErrMustStakeBeforeUnstake
	}
	if staked.GetAmountBigInt().Cmp(txBody.GetAmountBigInt()) < 0 {
		return nil, types.ErrExceedAmount
	}
	if staked.GetWhen()+StakingDelay > blockNo {
		return nil, types.ErrLessTimeHasPassed
	}
	toBe := new(big.Int).Sub(staked.GetAmountBigInt(), txBody.GetAmountBigInt())
	stakingMin, err := getStakingMinimum(scs)
	if err != nil {
		return nil, err
	}
	if toBe.Cmp(big.NewInt(0)) != 0 && stakingMin.Cmp(toBe) > 0 {
		return nil, types.ErrTooSmallAmount
	}
	return staked, nil
}

func parseIDForProposal(ci *types.CallInfo) (string, error) {
	//length should be checked before this function
	id, ok := ci.Args[0].(string)
	if !ok || len(id) < 1 || !isValidID(id) {
		return "", fmt.Errorf("args[%d] invalid id", 0)
	}
	return id, nil
}

func checkOperatorArg(context *SystemContext, ci *types.CallInfo) error {
	if len(ci.Args) != 1 { //args[0] : operator address
		return fmt.Errorf("invalid argument count %s : %s", ci.Name, ci.Args)
	}
	arg, ok := ci.Args[0].(string)
	if !ok {
		return fmt.Errorf("invalid string in the argument: %s", ci.Args)
	}
	address := types.ToAddress(arg)
	if len(address) == 0 {
		return fmt.Errorf("invalid address: %s", ci.Args[0])
	}
	context.Args = append(context.Args, arg)
	return nil
}

func checkOperator(scs *state.ContractState, address []byte) (Operators, error) {
	ops, err := getOperators(scs)
	if err != nil {
		return nil, fmt.Errorf("could not get admin in enterprise contract")
	}
	if ops == nil {
		return nil, ErrTxSystemOperatorIsNotSet
	}
	if i := bytes.Index(bytes.Join(ops, []byte("")), address); i == -1 && i%types.AddressLength != 0 {
		return nil, fmt.Errorf("operator address not matched")
	}
	return ops, nil
}
