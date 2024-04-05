package acquirer

import (
	"fmt"
	"testing"
)

func TestSqliteStorage(t *testing.T) {

	info := &DBMetaInfo{
		InfoHash: "XXXXXXXXXX",
		Metadata: "xxxxxxxxxxx",
	}

	DataStorage.Put(info)
	info = DataStorage.Get("XXXXXXXXXX")
	fmt.Println(info.Id, info.InfoHash, info.Metadata)
}
