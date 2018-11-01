/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package p2p

import (
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestBlockResponseHandler_handle(t *testing.T) {
	blkNo  := uint64(100)
	prevHash := dummyBlockHash
	inputHashes := make([]message.BlockHash,len(sampleBlks))
	inputBlocks := make([]*types.Block,len(sampleBlks))
	for i, hash := range sampleBlks {
		inputHashes[i] = hash
		inputBlocks[i] = &types.Block{Hash:hash, Header:&types.BlockHeader{PrevBlockHash:prevHash, BlockNo:blkNo}}
		blkNo++
		prevHash = hash
	}

	tests := []struct {
		name string

		receiver ResponseReceiver
		consume bool
		callSM bool
	}{
		// 1. not exist receiver and consumed message
		//{"Tnothing",nil, true},
		// 2. exist receiver and consume successfully
		{"TexistAndConsume", func(msg Message, body proto.Message) bool {
			return true
		}, true, false},
		// 2. exist receiver but not consumed
		{"TExistWrong", dummyResponseReceiver, false, true},
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPM := new(MockPeerManager)
			mockPeer := new(MockRemotePeer)
			mockPeer.On("ID").Return(dummyPeerID)
			mockPeer.On("consumeRequest", mock.AnythingOfType("p2p.MsgID"))
			mockActor := new(MockActorService)
			mockSM := new(MockSyncManager)
			mockSM.On("HandleGetBlockResponse",mockPeer, mock.Anything, mock.AnythingOfType("*types.GetBlockResponse"))

			mockPeer.On("GetReceiver", mock.AnythingOfType("p2p.MsgID")).Return(test.receiver)
			msg := &V030Message{subProtocol:GetBlocksResponse, id: sampleMsgID}
			body := &types.GetBlockResponse{Blocks:make([]*types.Block,2)}
			h := newBlockRespHandler(mockPM, mockPeer, logger, mockActor, mockSM)
			h.handle(msg, body)
			if  test.consume {
				mockSM.AssertNumberOfCalls(t, "HandleGetBlockResponse", 0)
			} else {
				mockSM.AssertNumberOfCalls(t, "HandleGetBlockResponse", 1)
			}
		})
	}
}
