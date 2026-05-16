package builtin

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	runtimetool "lattice-coding/internal/runtime/tool"
)

func recordFileReadState(manager runtimetool.FileReadStateManager, absPath string, info os.FileInfo, content []byte) string {
	checksum := fileChecksum(content)
	manager.MarkRead(runtimetool.ReadFileState{
		Path:       absPath,
		Version:    info.ModTime().UTC().Format("2006-01-02T15:04:05.000000000Z07:00"),
		MTime:      info.ModTime().UnixNano(),
		Checksum:   checksum,
		LastReadAt: time.Now().Unix(),
	})
	return checksum
}

func fileChecksum(content []byte) string {
	checksum := sha256.Sum256(content)
	return hex.EncodeToString(checksum[:])
}
