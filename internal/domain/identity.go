package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func StableID(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(strings.TrimSpace(part)))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func AttachmentChecksum(name, path string) string {
	return StableID(name, path)
}

func RecordID(machineID, title string) string {
	return "record-" + StableID(machineID, title)
}

func WorkflowID(recordID, kind string) string {
	return "workflow-" + StableID(recordID, kind)
}

func EventID(recordID, action string, createdAt int64) string {
	return "event-" + StableID(recordID, action, string(rune(createdAt)))
}
