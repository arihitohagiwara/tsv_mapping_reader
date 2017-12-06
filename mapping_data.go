package tsv_mapping_reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// MappingData 管理用構造体
type MappingData struct {
	mappingFilePath string
	modifiedTime    int64
	checkDulation   time.Duration
	mapData         map[string]string
	sched           *time.Timer
}

// NewMapping 各リソースの初期化
func NewMapping(mapFile string, dulation int64) (*MappingData, error) {
	data := &MappingData{
		mappingFilePath: mapFile,
		modifiedTime:    0,
		checkDulation:   time.Duration(dulation),
		mapData:         make(map[string]string),
		sched:           nil,
	}
	ret := data.readMapping()
	if ret == false {
		return nil, fmt.Errorf("newMapping initialize error")
	}
	if dulation > 0 {
		data.sched = time.AfterFunc(data.checkDulation*time.Second, data.ReadMappingSched)
	}
	return data, nil
}

// ReadMappingSched tsvファイルを定期的にmapに読み込む
func (data *MappingData) ReadMappingSched() {
	data.readMapping()
	data.sched = time.AfterFunc(data.checkDulation*time.Second, data.ReadMappingSched)
}

// SchedCancel ファイルを読み込むスケジュールをキャンセルする。
func (data *MappingData) SchedCancel() {
	if data.sched != nil {
		data.sched.Stop()
	}
	data.sched = nil
}

func (data *MappingData) Get(key string) (string, bool) {
	v, ok := data.mapData[key]
	return v, ok
}

// ファイルの更新があったかチェックする
func (data *MappingData) fileModifiedCheck() bool {
	fileInfo, err := os.Stat(data.mappingFilePath)
	if err != nil {
		fmt.Errorf("MappingData file not found error:", err.Error(), " file:", data.mappingFilePath)
		return false
	}
	modTime := fileInfo.ModTime().Unix()
	if modTime > data.modifiedTime {
		fmt.Println("MappingData file is modified file:", data.mappingFilePath)
		data.modifiedTime = modTime
		return true
	}
	fmt.Println("MappingData file is not modified file:", data.mappingFilePath)
	return false
}

// tsv形式のファイルを読み込んでmapに詰める
func (data *MappingData) readMapping() bool {
	if data.fileModifiedCheck() == false {
		return false
	}
	fd, err := os.Open(data.mappingFilePath)
	defer fd.Close()
	if err != nil {
		fmt.Errorf("MappingData file open error :", data.mappingFilePath)
		return false
	}
	reader := csv.NewReader(fd)
	reader.Comma = '\t'
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない！
	tmpMap := make(map[string]string)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Errorf("MappingData file Read failed file:", data.mappingFilePath, " err:", err)
			return false
		}
		if len(record) <= 1 {
			fmt.Errorf("MappingData file format error:", data.mappingFilePath)
			return false
		}
		if strings.HasPrefix(record[0], "#") == true {
			continue
		}
		tmpMap[record[0]] = record[1]
	}
	// 作成に成功した場合のみ上書きする
	data.mapData = tmpMap
	return true
}
