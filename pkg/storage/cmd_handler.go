package storage

import (
	"github.com/deepfabric/beehive/pb/raftcmdpb"
	"github.com/deepfabric/beehive/raftstore"
	"github.com/deepfabric/busybee/pkg/pb/rpcpb"
	"github.com/fagongzi/goetty"
	"github.com/fagongzi/log"
	"github.com/fagongzi/util/protoc"
)

func (h *beeStorage) init() {
	h.AddReadFunc("get", uint64(rpcpb.Get), h.get)
	h.AddWriteFunc("allocid", uint64(rpcpb.AllocID), h.allocID)
	h.AddWriteFunc("resetid", uint64(rpcpb.ResetID), h.resetID)
	h.AddReadFunc("bm-contains", uint64(rpcpb.BMContains), h.bmcontains)
	h.AddReadFunc("bm-count", uint64(rpcpb.BMCount), h.bmcount)
	h.AddReadFunc("bm-range", uint64(rpcpb.BMRange), h.bmrange)

	h.AddWriteFunc("starting-instance", uint64(rpcpb.StartingInstance), h.startingInstance)
	h.AddWriteFunc("started-instance", uint64(rpcpb.StartedInstance), h.startedInstance)
	h.AddWriteFunc("stop-instance", uint64(rpcpb.StopInstance), h.stopInstance)
	h.AddWriteFunc("create-state", uint64(rpcpb.CreateInstanceStateShard), h.createState)
	h.AddWriteFunc("update-state", uint64(rpcpb.UpdateInstanceStateShard), h.updateState)
	h.AddWriteFunc("remove-state", uint64(rpcpb.RemoveInstanceStateShard), h.removeState)
	h.AddWriteFunc("queue-fetch", uint64(rpcpb.QueueFetch), h.queueFetch)

	h.runner.RunCancelableTask(h.handleShardCycle)
}

func (h *beeStorage) BuildRequest(req *raftcmdpb.Request, cmd interface{}) error {
	switch cmd.(type) {
	case *rpcpb.SetRequest:
		msg := cmd.(*rpcpb.SetRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.Set)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseSetRequest(msg)
	case *rpcpb.GetRequest:
		msg := cmd.(*rpcpb.GetRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.Get)
		req.Type = raftcmdpb.Read
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseGetRequest(msg)
	case *rpcpb.DeleteRequest:
		msg := cmd.(*rpcpb.DeleteRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.Delete)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseDeleteRequest(msg)
	case *rpcpb.AllocIDRequest:
		msg := cmd.(*rpcpb.AllocIDRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.AllocID)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseAllocIDRequest(msg)
	case *rpcpb.ResetIDRequest:
		msg := cmd.(*rpcpb.ResetIDRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.ResetID)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseResetIDRequest(msg)
	case *rpcpb.BMCreateRequest:
		msg := cmd.(*rpcpb.BMCreateRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMCreate)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMCreateRequest(msg)
	case *rpcpb.BMAddRequest:
		msg := cmd.(*rpcpb.BMAddRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMAdd)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMAddRequest(msg)
	case *rpcpb.BMRemoveRequest:
		msg := cmd.(*rpcpb.BMRemoveRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMRemove)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMRemoveRequest(msg)
	case *rpcpb.BMClearRequest:
		msg := cmd.(*rpcpb.BMClearRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMClear)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMClearRequest(msg)
	case *rpcpb.BMContainsRequest:
		msg := cmd.(*rpcpb.BMContainsRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMContains)
		req.Type = raftcmdpb.Read
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMContainsRequest(msg)
	case *rpcpb.BMCountRequest:
		msg := cmd.(*rpcpb.BMCountRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMCount)
		req.Type = raftcmdpb.Read
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMCountRequest(msg)
	case *rpcpb.BMRangeRequest:
		msg := cmd.(*rpcpb.BMRangeRequest)
		req.Key = KVKey(msg.Key)
		req.CustemType = uint64(rpcpb.BMRange)
		req.Type = raftcmdpb.Read
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseBMRangeRequest(msg)
	case *rpcpb.StartingInstanceRequest:
		msg := cmd.(*rpcpb.StartingInstanceRequest)
		req.Key = StartedInstanceKey(msg.Instance.Snapshot.ID)
		req.CustemType = uint64(rpcpb.StartingInstance)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseStartingInstanceRequest(msg)
	case *rpcpb.StartedInstanceRequest:
		msg := cmd.(*rpcpb.StartedInstanceRequest)
		req.Key = StartedInstanceKey(msg.WorkflowID)
		req.CustemType = uint64(rpcpb.StartedInstance)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseStartedInstanceRequest(msg)
	case *rpcpb.StopInstanceRequest:
		msg := cmd.(*rpcpb.StopInstanceRequest)
		req.Key = StartedInstanceKey(msg.WorkflowID)
		req.CustemType = uint64(rpcpb.StopInstance)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseStopInstanceRequest(msg)
	case *rpcpb.CreateInstanceStateShardRequest:
		msg := cmd.(*rpcpb.CreateInstanceStateShardRequest)
		req.Key = InstanceShardKey(msg.State.WorkflowID, msg.State.Start, msg.State.End)
		req.CustemType = uint64(rpcpb.CreateInstanceStateShard)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseCreateInstanceStateShardRequest(msg)
	case *rpcpb.UpdateInstanceStateShardRequest:
		msg := cmd.(*rpcpb.UpdateInstanceStateShardRequest)
		req.Key = InstanceShardKey(msg.State.WorkflowID, msg.State.Start, msg.State.End)
		req.CustemType = uint64(rpcpb.UpdateInstanceStateShard)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseUpdateInstanceStateShardRequest(msg)
	case *rpcpb.RemoveInstanceStateShardRequest:
		msg := cmd.(*rpcpb.RemoveInstanceStateShardRequest)
		req.Key = InstanceShardKey(msg.WorkflowID, msg.Start, msg.End)
		req.CustemType = uint64(rpcpb.RemoveInstanceStateShard)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseRemoveInstanceStateShardRequest(msg)
	case *rpcpb.QueueAddRequest:
		msg := cmd.(*rpcpb.QueueAddRequest)
		req.Key = msg.Key
		req.CustemType = uint64(rpcpb.QueueAdd)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseQueueAddRequest(msg)
	case *rpcpb.QueueFetchRequest:
		msg := cmd.(*rpcpb.QueueFetchRequest)
		req.Key = msg.Key
		req.CustemType = uint64(rpcpb.QueueFetch)
		req.Type = raftcmdpb.Write
		req.Cmd = protoc.MustMarshal(msg)
		rpcpb.ReleaseQueueFetchRequest(msg)
	default:
		log.Fatalf("not support request %+v(%+T)", cmd, cmd)
	}

	return nil
}

func (h *beeStorage) Codec() (goetty.Decoder, goetty.Encoder) {
	return nil, nil
}

func (h *beeStorage) AddReadFunc(cmd string, cmdType uint64, cb raftstore.ReadCommandFunc) {
	h.store.RegisterReadFunc(cmdType, cb)
}

func (h *beeStorage) AddWriteFunc(cmd string, cmdType uint64, cb raftstore.WriteCommandFunc) {
	h.store.RegisterWriteFunc(cmdType, cb)
}

func (h *beeStorage) AddLocalFunc(cmd string, cmdType uint64, cb raftstore.LocalCommandFunc) {
	h.store.RegisterLocalFunc(cmdType, cb)
}

func (h *beeStorage) WriteBatch() raftstore.CommandWriteBatch {
	return newBatch(h, newKVBatch(), newBitmapBatch(), newQueueBatch())
}

func (h *beeStorage) getValue(shard uint64, key []byte) ([]byte, error) {
	value, err := h.getStore(shard).Get(key)
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, nil
	}

	return value[1:], nil
}

func (h *beeStorage) getValueWithPrefix(shard uint64, key []byte) ([]byte, error) {
	value, err := h.getStore(shard).Get(key)
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, nil
	}

	return value, nil
}
