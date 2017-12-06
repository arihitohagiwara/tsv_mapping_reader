package tsv_mapping_reader

import (
	"os/exec"
	"testing"
	"time"
)

func TestNewMapping(t *testing.T) {
	// ファイルがない場合
	data, err := NewMapping("./testdata/nofile.tsv", 0)
	if err == nil {
		t.Errorf("TestNewMapping\n")
	}
	data, err = NewMapping("./testdata/maptest.tsv", 0)
	if data == nil || err != nil {
		t.Errorf("TestNewMapping %v\n", err)
	}
	// ファイルがある場合(スケジュール登録あり)
	data, err = NewMapping("./testdata/maptest.tsv", 1)
	if data == nil || err != nil {
		t.Errorf("TestNewMapping %v\n", err)
	}
	// scheduleを停めておく
	data.SchedCancel()
}

func TestGet(t *testing.T) {
	data, err := NewMapping("./testdata/maptest.tsv", 0)
	if data == nil || err != nil {
		t.Errorf("TestNewMapping %v\n", err)
	}
	// データがある場合
	v, ok := data.Get("hogehoge")
	if v != "1" || ok != true {
		t.Errorf("TestGet %v\n", v)
	}
	// データがない場合
	v, ok = data.Get("nashi")
	if ok != false {
		t.Errorf("TestGet %v\n", ok)
	}
}

func TestFileModifiedCheck(t *testing.T) {
	data, _ := NewMapping("./testdata/maptest.tsv", 0)
	// init時に読んでいるので更新なし
	actual := data.fileModifiedCheck()
	expected := false
	if actual != expected {
		t.Errorf("got fileModifiedCheck %v\nwant %v\n", actual, expected)
	}
	// touchしてファイルを更新
	exec.Command("touch", "./testdata/maptest.tsv").Output()
	actual = data.fileModifiedCheck()
	expected = true
	if actual != expected {
		t.Errorf("got fileModifiedCheck %v\nwant %v\n", actual, expected)
	}
}

func TestReadMapping(t *testing.T) {
	data, _ := NewMapping("./testdata/maptest.tsv", 0)
	// init時に読んでいるので更新なし
	actual := data.readMapping()
	expected := false
	if actual != expected {
		t.Errorf("got readMapping %v\nwant %v\n", actual, expected)
	}
	// touchしてファイルを更新
	time.Sleep(1000000000)
	exec.Command("touch", "testdata/maptest.tsv").Output()
	actual = data.readMapping()
	expected = true
	if actual != expected {
		t.Errorf("got readMapping %v\nwant %v\n", actual, expected)
	}
}
