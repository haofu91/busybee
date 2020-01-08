package storage

import (
	"github.com/deepfabric/busybee/pkg/pb/rpcpb"
	"github.com/fagongzi/goetty"
	"github.com/fagongzi/log"
	"github.com/fagongzi/util/protoc"
)

var (
	bizCodec = &codec{}
	decoder  = goetty.NewIntLengthFieldBasedDecoder(bizCodec)
	encoder  = goetty.NewIntLengthFieldBasedEncoder(bizCodec)
)

type codec struct {
}

func (c *codec) Decode(in *goetty.ByteBuf) (bool, interface{}, error) {
	v, err := in.ReadByte()
	if err != nil {
		return true, nil, err
	}

	var msg protoc.PB
	data := in.GetMarkedRemindData()
	rpcType := rpcpb.Type(v)

	switch rpcType {
	case rpcpb.Set:
		msg = rpcpb.AcquireSetRequest()
	case rpcpb.Get:
		msg = rpcpb.AcquireGetRequest()
	case rpcpb.Delete:
		msg = rpcpb.AcquireDeleteRequest()
	case rpcpb.BMCreate:
		msg = rpcpb.AcquireBMCreateRequest()
	case rpcpb.BMAdd:
		msg = rpcpb.AcquireBMAddRequest()
	case rpcpb.BMRemove:
		msg = rpcpb.AcquireBMRemoveRequest()
	case rpcpb.BMClear:
		msg = rpcpb.AcquireBMClearRequest()
	case rpcpb.BMContains:
		msg = rpcpb.AcquireBMContainsRequest()
	case rpcpb.BMDel:
		msg = rpcpb.AcquireBMDelRequest()
	case rpcpb.BMCount:
		msg = rpcpb.AcquireBMCountRequest()
	case rpcpb.BMRange:
		msg = rpcpb.AcquireBMRangeRequest()
	case rpcpb.StartingInstance:
		msg = rpcpb.AcquireStartingInstanceRequest()
	case rpcpb.StartedInstance:
		msg = rpcpb.AcquireStartedInstanceRequest()
	case rpcpb.CreateInstanceStateShard:
		msg = rpcpb.AcquireCreateInstanceStateShardRequest()
	case rpcpb.UpdateInstanceStateShard:
		msg = rpcpb.AcquireUpdateInstanceStateShardRequest()
	case rpcpb.RemoveInstanceStateShard:
		msg = rpcpb.AcquireRemoveInstanceStateShardRequest()
	case rpcpb.StepInstanceStateShard:
		msg = rpcpb.AcquireStepInstanceStateShardRequest()
	case rpcpb.QueueAdd:
		msg = rpcpb.AcquireQueueAddRequest()
	case rpcpb.QueueFetch:
		msg = rpcpb.AcquireQueueFetchRequest()
	default:
		log.Fatalf("BUG: not support msg type %d", v)
	}

	protoc.MustUnmarshal(msg, data)
	in.MarkedBytesReaded()
	return true, msg, nil
}

func (c *codec) Encode(data interface{}, out *goetty.ByteBuf) error {
	var t rpcpb.Type
	switch data.(type) {
	case *rpcpb.SetResponse:
		t = rpcpb.Set
	case *rpcpb.GetResponse:
		t = rpcpb.Get
	case *rpcpb.DeleteResponse:
		t = rpcpb.Delete
	case *rpcpb.BMCreateResponse:
		t = rpcpb.BMCreate
	case *rpcpb.BMAddResponse:
		t = rpcpb.BMAdd
	case *rpcpb.BMRemoveResponse:
		t = rpcpb.BMRemove
	case *rpcpb.BMClearResponse:
		t = rpcpb.BMClear
	case *rpcpb.BMContainsResponse:
		t = rpcpb.BMContains
	case *rpcpb.BMDelResponse:
		t = rpcpb.BMDel
	case *rpcpb.BMCountResponse:
		t = rpcpb.BMCount
	case *rpcpb.BMRangeResponse:
		t = rpcpb.BMRange
	case *rpcpb.StartingInstanceResponse:
		t = rpcpb.StartingInstance
	case *rpcpb.StartedInstanceResponse:
		t = rpcpb.StartedInstance
	case *rpcpb.CreateInstanceStateShardResponse:
		t = rpcpb.CreateInstanceStateShard
	case *rpcpb.UpdateInstanceStateShardResponse:
		t = rpcpb.UpdateInstanceStateShard
	case *rpcpb.RemoveInstanceStateShardResponse:
		t = rpcpb.RemoveInstanceStateShard
	case *rpcpb.StepInstanceStateShardResponse:
		t = rpcpb.StepInstanceStateShard
	case *rpcpb.QueueAddResponse:
		t = rpcpb.QueueAdd
	case *rpcpb.QueueFetchResponse:
		t = rpcpb.QueueFetch
	default:
		log.Fatalf("BUG: not support msg type %T", data)
	}

	m := data.(protoc.PB)
	size := m.Size()
	out.WriteInt(size + 1)
	out.WriteByte(byte(t))

	if size > 0 {
		index := out.GetWriteIndex()
		out.Expansion(size)
		protoc.MustMarshalTo(m, out.RawBuf()[index:index+size])
		out.SetWriterIndex(index + size)
	}

	return nil
}