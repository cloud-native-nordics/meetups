package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	isTesting = true
	for i := 1; i <= 2; i++ {
		indir := fmt.Sprintf("testdata/%d", i)
		t.Run(indir, func(t2 *testing.T) {
			speakersFile := filepath.Join(indir, "speakers.yaml")
			companiesFile := filepath.Join(indir, "companies.yaml")

			cfg, err := load(companiesFile, speakersFile, indir)
			if err != nil {
				t.Fatal(err)
			}
			if err := update(cfg); err != nil {
				t.Fatal(err)
			}
			out, err := exec(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if err := validate(out, indir); err != nil {
				t.Errorf("generation and validation failed: %v", err)
			}
			//t.Log(out)
			*dryRun = false
			if err := apply(out, indir); err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
