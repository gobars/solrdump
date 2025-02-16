package main

import (
	"bytes"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bingoohuang/gg/pkg/vars"
	"github.com/bingoohuang/jj"
)

func (a *Arg) createBulkOutput(uri string) func(doc []byte) {
	if !strings.Contains(uri, "/_bulk") {
		return nil
	}

	// support es bulk mode
	docCh := make(chan []byte, a.Bulk)
	fn := func(doc []byte) { docCh <- doc }

	var wg sync.WaitGroup
	wg.Add(1)
	go a.elasticSearchBulk(uri, docCh, &wg)

	a.closers = append(a.closers, closeFn(func() { close(docCh); wg.Wait() }))

	return fn
}

func (a *Arg) elasticSearchBulk(uri string, docCh chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	u, _ := url.Parse(uri)
	query := u.Query()

	var routingExpr vars.Subs
	if routing := query.Get("routing"); routing != "" {
		query.Del("routing")
		routingExpr = vars.ParseExpr(routing)
	}

	u.RawQuery = query.Encode()
	uri = u.String()
	b := &bytes.Buffer{}

	for {
		ok := a.numOrTicker(b, docCh, routingExpr)
		outputHttp(uri, b.Bytes(), a.Verbose, a.printer)
		if !ok {
			return
		}
	}
}

func (a *Arg) numOrTicker(b *bytes.Buffer, docCh chan []byte, routingExpr vars.Subs) (continued bool) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	num := 0
	b.Reset()

	for {
		select {
		case <-ticker.C:
			if num > 0 {
				return true
			}

		case doc, ok := <-docCh:
			if !ok {
				return false
			}

			var bulkFirstLine []byte
			if len(routingExpr) > 0 {
				evalResult, err := routingExpr.Eval(&JsonValue{Doc: doc})
				if err != nil {

				}
				routing := evalResult.(string)
				bulkFirstLine = []byte(`{"index":{"_type":"docs","` + a.Routing + `":"` + routing + `"}}`)
			} else {
				bulkFirstLine = []byte(`{"index":{"_type":"docs"}}`)
			}

			if lc := jj.GetBytes(doc, "licenseCode").String(); lc != "" { // licenseCode 是唯一的，作为唯一表示。和更新操作保持一致
				bulkFirstLine, _ = jj.SetBytes(bulkFirstLine, "index._id", lc)
			}

			b.Write(bulkFirstLine)
			b.Write([]byte("\n"))
			b.Write(doc)
			b.Write([]byte("\n"))
			if num++; num >= a.Bulk {
				return true
			}
		}
	}
}
