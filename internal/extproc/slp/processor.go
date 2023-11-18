package slp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct {
		confPath string
	}
)

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	cmd := exec.Command("slp", "my", "--config", p.confPath, "--output", "standard", "--format", "tsv", "--file", bodyPath)

	res, err := cmd.Output()
	if err != nil {
		log.Println(string(res))
		return nil, fmt.Errorf("external process aborted: %w", err)
	}

	cmd = exec.Command("pt-query-digest", "--type", "slowlog", bodyPath)
	ptResult, err := cmd.Output()
	if err != nil {
		log.Printf("failed to execute pt-query-digest: %v", err)
		return nil, nil //本物とは関係ない処理なのでエラーが出ても通す。
	}
	err = savePtResult("data/pt-"+snapshot.ID, ptResult)
	if err != nil {
		log.Printf("failed to save pt-query-digest result: %v", err)
		return nil, nil //本物とは関係ない処理なのでエラーが出ても通す。
	}

	return io.NopCloser(bytes.NewBuffer(res)), nil
}

func savePtResult(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file for pt result: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write pt result to file: %w", err)
	}

	return nil
}
