package consensus

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const consensusFile = "./database/consensus/consensus.data"

type AuditLog struct {
	Identity           string
	Chunk              string
	TimeInMilliseconds int32
	Result             string
	Message            string
	NextAction         string
}

type Audit struct {
	Identity      string
	Successes     int32
	Offlines      int32
	Fails         int32
	Status        string
	TotalCount    int32
	LastAuditTime string
}

type ConsensusData struct {
	Data []Audit
	Logs []AuditLog
}

func (cs *ConsensusData) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(cs)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(consensusFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (cs *ConsensusData) LoadFile() error {
	if _, err := os.Stat(consensusFile); os.IsNotExist(err) {
		return err
	}

	var consensusData ConsensusData

	fileContent, err := ioutil.ReadFile(consensusFile)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&consensusData)
	if err != nil {
		return err
	}

	cs.Data = consensusData.Data

	return nil
}

func CreateAuditLog(Identity string,
	Chunk string,
	TimeInMilliseconds int32,
	Result string,
	Message string,
	NextAction string) AuditLog {
	consensus := AuditLog{
		Identity:           Identity,
		Chunk:              Chunk,
		TimeInMilliseconds: TimeInMilliseconds,
		Result:             Result,
		Message:            Message,
		NextAction:         NextAction,
	}

	return consensus
}

func CreateAudit(Identity string,
	Successes int32,
	Offlines int32,
	Fails int32,
	Status string,
	TotalCount int32,
	LastAuditTime string) Audit {
	consensus := Audit{
		Identity:      Identity,
		Successes:     Successes,
		Offlines:      Offlines,
		Fails:         Fails,
		Status:        Status,
		TotalCount:    TotalCount,
		LastAuditTime: LastAuditTime,
	}

	return consensus
}

func CreateConsensusData() (*ConsensusData, error) {
	consensusData := ConsensusData{}
	err := consensusData.LoadFile()

	return &consensusData, err
}

func (cs *ConsensusData) AddAudit(audit Audit) []Audit {

	cs.Data = append(cs.Data, audit)

	return cs.Data
}

func (cs *ConsensusData) AddAuditLog(auditLog AuditLog) []AuditLog {

	cs.Logs = append(cs.Logs, auditLog)

	return cs.Logs
}
