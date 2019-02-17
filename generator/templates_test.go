package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for i := 1; i <= 2; i++ {
		inpath := fmt.Sprintf("testdata/%d/in.yaml", i)
		t.Run(inpath, func(t2 *testing.T) {
			in, err := ioutil.ReadFile(inpath)
			if err != nil {
				t.Fatal(err)
			}

			cfg, err := load([]byte(in))
			if err != nil {
				t.Fatal(err)
			}
			out, err := exec(cfg)
			if err != nil {
				t.Fatal(err)
			}
			for city, content := range out {
				outpath := fmt.Sprintf("testdata/%d/%s.md", i, strings.ToLower(city))
				expected, err := ioutil.ReadFile(outpath)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(content, expected) {
					t.Errorf("%s: actual: %s, expected: %s", outpath, content, expected)
				}
			}
		})
	}
}
