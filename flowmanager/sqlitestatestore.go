package flowmanager

import (
	"database/sql"
	"math/big"
	"os"
	"path"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	"github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStateStore struct {
	db *sql.DB
	tx *sql.Tx

	latestMBlock *protos.MBlock
	nextMBlock   *protos.MBlock
	finalized    bool
}

func NewSQLiteStateStore(dbFile string) (*SQLiteStateStore, error) {
	err := os.MkdirAll(path.Dir(dbFile), 0700)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFile+"?journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	s := &SQLiteStateStore{
		db: db,
	}

	err = s.ensureSchemas()
	if err != nil {
		s.Close()
		return nil, err
	}

	s.tx, err = s.db.Begin()
	if err != nil {
		s.Close()
		return nil, err
	}

	s.latestMBlock, err = s.LatestMBlock()
	if err != nil {
		s.Close()
		return nil, err
	}

	err = s.ResetPending()
	if err != nil {
		s.Close()
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStateStore) ensureSchemas() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS mblock (
	height	INTEGER	PRIMARY KEY,
	hash	BLOB	UNIQUE NOT NULL,
	payload	BLOB
)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS last_block_report (
	chain_id	TEXT	PRIMARY KEY,
	payload		BLOB
)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS workflow (
	id		BLOB	PRIMARY KEY,
	payload	BLOB
)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS monitor (
	workflow_id	BLOB	NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
	chain_id	TEXT	NOT NULL,
	payload		BLOB
);
CREATE INDEX IF NOT EXISTS idx_monitor_workflow_id ON monitor(workflow_id);
CREATE INDEX IF NOT EXISTS idx_monitor_chain_id ON monitor(chain_id);
`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS action (
	workflow_id	BLOB	NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
	action_name	TEXT	NOT NULL,
	payload		BLOB,
	PRIMARY KEY(workflow_id, action_name)
)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS event_action (
	workflow_id	BLOB	NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
	event_name	TEXT	NOT NULL,
	action_name	TEXT	NOT NULL,
	PRIMARY KEY(workflow_id, event_name, action_name)
)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
CREATE TABLE IF NOT EXISTS action_state (
	workflow_id	BLOB	NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
	action_name	TEXT	NOT NULL,
	event_name	TEXT	NOT NULL,
	payload		BLOB,
	PRIMARY KEY(workflow_id, action_name, event_name)
)`)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStateStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStateStore) LastReportedBlockByChain(chain string) (*protos.Block, error) {
	row := s.tx.QueryRow(`
SELECT payload
FROM last_block_report
WHERE chain_id = ?
`, chain)

	var payload []byte
	err := row.Scan(&payload)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	block := new(protos.Block)
	err = proto.Unmarshal(payload, block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s *SQLiteStateStore) WorkflowByID(id WorkflowID) (*workflow.WorkflowDecl, error) {
	row := s.tx.QueryRow(`
SELECT payload
FROM workflow
WHERE id = ?
`, id.UnsafeBytes())

	var payload []byte
	err := row.Scan(&payload)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	wf, err := workflow.NewByteDecoder(payload).DecodeWorkflow()
	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (s *SQLiteStateStore) MonitorsByChain(chain string) ([]*grpc.Monitor, error) {
	rows, err := s.tx.Query(`
SELECT workflow_id, payload
FROM monitor
WHERE chain_id = ?
`, chain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []*grpc.Monitor

	for rows.Next() {
		var wfid []byte
		var payload []byte

		if err = rows.Scan(&wfid, &payload); err != nil {
			return nil, err
		}

		monitor, err := workflow.NewByteDecoder(payload).DecodeMonitorDecl()
		if err != nil {
			return nil, err
		}

		monitors = append(monitors, &grpc.Monitor{
			WorkflowId:  wfid,
			MonitorName: monitor.Name().Id(),
			Monitor:     payload,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (s *SQLiteStateStore) ActionStatesByEvent(wfid WorkflowID, eventName string) ([]*ActionState, error) {
	rows, err := s.tx.Query(`
SELECT
	e.action_name,
	a.payload action_payload,
	s.event_name,
	s.payload state_payload
FROM event_action e
	JOIN action a ON e.workflow_id = a.workflow_id AND e.action_name = a.action_name
	LEFT JOIN action_state s ON e.workflow_id = s.workflow_id AND e.action_name = s.action_name
WHERE e.workflow_id = ? AND e.event_name = ?
ORDER BY e.rowid
`, wfid.UnsafeBytes(), eventName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*ActionState
	var currentAction *ActionState

	for rows.Next() {
		var actionName string
		var eventName sql.NullString
		var actionPayload, statePayload []byte

		if err = rows.Scan(&actionName, &actionPayload, &eventName, &statePayload); err != nil {
			return nil, err
		}

		if currentAction == nil || actionName != currentAction.Action.Name().Id() {
			if currentAction != nil {
				states = append(states, currentAction)
			}

			action, err := workflow.NewByteDecoder(actionPayload).DecodeActionDecl()
			if err != nil {
				return nil, err
			}

			currentAction = &ActionState{
				Action: action,
				State:  make(map[string]*workflow.ObjConst),
			}
		}

		if eventName.Valid {
			state, err := workflow.NewByteDecoder(statePayload).DecodeObjConst()
			if err != nil {
				return nil, err
			}

			currentAction.State[eventName.String] = state
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if currentAction != nil {
		states = append(states, currentAction)
	}

	return states, nil
}

func scanMBlock(row *sql.Row) (*protos.MBlock, error) {
	var payload []byte
	err := row.Scan(&payload)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	mblock := new(protos.MBlock)
	err = proto.Unmarshal(payload, mblock)
	if err != nil {
		return nil, err
	}

	return mblock, nil
}

func (s *SQLiteStateStore) LatestMBlock() (*protos.MBlock, error) {
	row := s.tx.QueryRow(`
SELECT payload
FROM mblock
ORDER BY height DESC
LIMIT 1
`)

	return scanMBlock(row)
}

func (s *SQLiteStateStore) MBlockByHeight(height int64) (*protos.MBlock, error) {
	row := s.tx.QueryRow(`
SELECT payload
FROM mblock
WHERE height = ?
`, height)

	return scanMBlock(row)
}

func (s *SQLiteStateStore) PutBlockReport(block *protos.Block) error {
	if s.finalized {
		return ErrFinalized
	}

	payload, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`
INSERT OR REPLACE INTO last_block_report(chain_id, payload)
VALUES (?, ?)
`, block.Chain, payload)
	if err != nil {
		return err
	}

	s.nextMBlock.Events = append(s.nextMBlock.Events, &protos.MEvent{
		Body: &protos.MEvent_ReportBlock{
			ReportBlock: block,
		},
	})

	return nil
}

func (s *SQLiteStateStore) PutMonitorEvent(event *protos.MonitorEvent) error {
	if s.finalized {
		return ErrFinalized
	}

	s.nextMBlock.Events = append(s.nextMBlock.Events, &protos.MEvent{
		Body: &protos.MEvent_MonitorEvent{
			MonitorEvent: event,
		},
	})

	return nil
}

func (s *SQLiteStateStore) PutActionEvent(event *protos.ActionEvent) error {
	if s.finalized {
		return ErrFinalized
	}

	s.nextMBlock.Events = append(s.nextMBlock.Events, &protos.MEvent{
		Body: &protos.MEvent_ActionEvent{
			ActionEvent: event,
		},
	})

	return nil
}

func (s *SQLiteStateStore) PutActionEventState(wfid WorkflowID, actionName string, eventName string, state *workflow.ObjConst) error {
	if s.finalized {
		return ErrFinalized
	}

	payload, err := workflow.EncodeObjConst(state)
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`
INSERT OR REPLACE INTO action_state(workflow_id, action_name, event_name, payload)
VALUES (?, ?, ?, ?)
`, wfid.UnsafeBytes(), actionName, eventName, payload)

	return err
}

func (s *SQLiteStateStore) ResetActionState(wfid WorkflowID, actionName string) error {
	if s.finalized {
		return ErrFinalized
	}

	_, err := s.tx.Exec(`
DELETE FROM action_state
WHERE workflow_id = ? AND action_name = ?
`, wfid.UnsafeBytes(), actionName)

	return err
}

func (s *SQLiteStateStore) PutWorkflow(id WorkflowID, wf *workflow.WorkflowDecl) error {
	if s.finalized {
		return ErrFinalized
	}

	payload, err := workflow.EncodeWorkflow(wf)
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`
INSERT INTO workflow(id, payload)
VALUES (?, ?)
`, id.UnsafeBytes(), payload)
	if err != nil {
		return err
	}

	insertMonitorStmt, err := s.tx.Prepare(`
INSERT INTO monitor(workflow_id, chain_id, payload)
VALUES (?, ?, ?)
`)
	if err != nil {
		return err
	}

	insertActionStmt, err := s.tx.Prepare(`
INSERT INTO action(workflow_id, action_name, payload)
VALUES (?, ?, ?)
`)
	if err != nil {
		return err
	}

	insertEventActionStmt, err := s.tx.Prepare(`
INSERT INTO event_action(workflow_id, event_name, action_name)
VALUES (?, ?, ?)
`)
	if err != nil {
		return err
	}

	for _, monitor := range wf.MonitorDecls() {
		payload, err = workflow.EncodeMonitorDecl(monitor)
		if err != nil {
			return err
		}

		_, err = insertMonitorStmt.Exec(id.UnsafeBytes(), monitor.Chain(), payload)
		if err != nil {
			return err
		}
	}

	for _, action := range wf.ActionDecls() {
		payload, err = workflow.EncodeActionDecl(action)
		if err != nil {
			return err
		}

		_, err = insertActionStmt.Exec(id.UnsafeBytes(), action.Name().Id(), payload)
		if err != nil {
			return err
		}

		for _, event := range action.TriggerEvents() {
			_, err = insertEventActionStmt.Exec(id.UnsafeBytes(), event, action.Name().Id())
			if err != nil {
				return err
			}
		}
	}

	s.nextMBlock.Events = append(s.nextMBlock.Events, &protos.MEvent{
		Body: &protos.MEvent_DeployWorkflow{
			DeployWorkflow: &protos.DeployWorkflow{
				WorkflowId: id.Bytes(),
				Payload:    payload,
			},
		},
	})

	return nil
}

func (s *SQLiteStateStore) DeleteWorkflow(id WorkflowID) error {
	if s.finalized {
		return ErrFinalized
	}

	_, err := s.tx.Exec(`
DELETE FROM workflow
WHERE id = ? 
`, id.UnsafeBytes())
	if err != nil {
		return err
	}

	s.nextMBlock.Events = append(s.nextMBlock.Events, &protos.MEvent{
		Body: &protos.MEvent_UndeployWorkflow{
			UndeployWorkflow: &protos.UndeployWorkflow{
				WorkflowId: id.Bytes(),
			},
		},
	})

	return nil
}

func (s *SQLiteStateStore) FinalizeMBlock() (*protos.MBlock, error) {
	// TODO - hash
	s.nextMBlock.Hash = big.NewInt(s.nextMBlock.Height).Bytes()
	s.finalized = true
	return s.nextMBlock, nil
}

func (s *SQLiteStateStore) CommitPendingMBlock(block *protos.MBlock) error {
	if block != s.nextMBlock {
		return ErrUnknownMBlock
	}

	payload, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`
INSERT INTO mblock(height, hash, payload)
VALUES (?, ?, ?)
`, block.Height, block.Hash, payload)
	if err != nil {
		return err
	}

	err = s.tx.Commit()
	if err != nil {
		return err
	}
	s.tx = nil

	s.latestMBlock = block

	return s.ResetPending()
}

func (s *SQLiteStateStore) ResetPending() error {
	if s.tx != nil {
		s.tx.Rollback()
		s.tx = nil
	}

	s.nextMBlock = new(protos.MBlock)
	if s.latestMBlock != nil {
		s.nextMBlock.ParentHash = s.latestMBlock.Hash
		s.nextMBlock.Height = s.latestMBlock.Height + 1
	}

	s.finalized = false

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx

	return nil
}
