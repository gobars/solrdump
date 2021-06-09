# README

Export documents from a SOLR index as JSON, fast and simply from the command line.

## Requirements

SOLR 4.7 or higher, since the cursor mechanism
introduced ([2014-02-25](https://archive.apache.org/dist/lucene/solr/4.7.0/)) &mdash; see
also [efficient deep paging with cursors](https://solr.pl/en/2014/03/10/solr-4-7-efficient-deep-paging/).

## Installation

1. Via debian or rpm [package](https://github.com/gobars/solrdump/releases).
2. Or via go tool: `go install github.com/gobars/solrdump@latest`

## Features

1. `_version_` deleted from the result
2. `output="http://127.0.0.1:9092/zz/docs?routing=@path"` `@path` in the url will be evaluated by `{GjsonPath}`
   , [Syntax](https://github.com/bingoohuang/jj/blob/master/SYNTAX.md)
3. `output="http://127.0.0.1:9092/zz/_bulk?routing=@path"` will change to bulk mode for elasticsearch automatically.

## Usage

```shell
Usage of solrdump (0.1.2):
  -max int       Max number of rows (default 100)
  -q string      SOLR query (default "*:*")
  -rows int      Number of rows returned per request (default 100)
  -server string SOLR server with index name, eg. localhost:8983/solr/example
  -version       Show version and exit
  -remove-fields Remove fields, _version_ defaulted
  -output        Output file, or http url, or noop
  -v             Verbose, -v -vv -vvv
```

### Print docs

```shell
$ solrdump -server 192.168.2.6:8983/solr/zz -max 3          
2021/06/03 17:57:50 http://192.168.2.6:8983/solr/zz/select?cursorMark=%2A&fl=&q=%2A%3A%2A&rows=3&sort=id+asc&wt=json
{"dataType":"INTEGRATION","id":"000007c8-3d83-47c7-b9f0-1e0d15670599","createdDate":"2020-08-06T06:49:37Z"}
{"dataType":"MANUAL","id":"00004d53-d76d-43c3-906d-90ff475bd1a2","createdDate":"2021-05-10T08:14:14Z"}
{"dataType":"MANUAL","id":"000070fe-309f-4755-998e-2445cc66ef9f","createdDate":"2021-05-10T08:14:14Z"}
2021/06/03 17:57:50 fetched 3/509309 docs
```

### Write to elastic search

```sh
$ solrdump -server 192.168.2.6:8983/solr/licenseIndex -max 10 -vv -output "192.168.2.8:9202/license/docs?routing=@holderIdentityNum.0"
2021/06/08 14:30:53 started
2021/06/08 14:30:58 fetched 10/509311 docs
2021/06/08 14:30:58 solr query: "http://192.168.2.6:8983/solr/licenseIndex/select?cursorMark=*&fl=&q=*:*&rows=10&sort=id asc&wt=json"
2021/06/08 14:30:58 http uri: http://192.168.2.8:9202/license/docs?routing=441421201510165436
2021/06/08 14:30:58 sent cost: 5.175523ms status: 201, body: {"_index":"license","_type":"docs","_id":"rB1R6nkBuRYpTrL3HTzU","_version":1,"result":"created","_shards":{"total":1,"successful":1,"failed":0},"_seq_no":30989363,"_primary_term":3}
2021/06/08 14:31:03 process 10 docs, rate 0.999881 docs/s, cost 10.001191743s
```

then [elasticsearch query on 441421201510165436](http://192.168.2.8:9202/license/_search?routing=441421201510165436&q=holderIdentityNum:441421201510165436)
can be performed.

## Resources

1. [o19s/solr-to-es](https://github.com/o19s/solr-to-es)
2. [solr cursor select query](https://github.com/frizner/glsolr)
3. [frizner/solrdump](https://github.com/frizner/solrdump)
4. [hectorcorrea/solr-for-newbies](https://github.com/hectorcorrea/solr-for-newbies)
5. [online json-generator](https://www.json-generator.com)

## Pagination of Results

* https://cwiki.apache.org/confluence/display/solr/Pagination+of+Results

Requesting large number of documents from SOLR can lead to *Deep Paging*
problems:

> When you wish to fetch a very large number of sorted results from Solr to
> feed into an external system, using very large values for the start or rows
> parameters can be very inefficient.

See also: *Fetching A Large Number of Sorted Results: Cursors*

> As an alternative to increasing the "start" parameter to request subsequent
> pages of sorted results, Solr supports using a "Cursor" to scan through
> results. Cursors in Solr are a logical concept, that doesn't involve caching
> any state information on the server. Instead the sort values of the last
> document returned to the client are used to compute a "mark" representing a
> logical point in the ordered space of sort values.

`http://192.168.2.6:8983/solr/zz/select?q=*:*&wt=json&cursorMark=*&sort=id asc`

```json
{
  "responseHeader": {
    "status": 0,
    "QTime": 13,
    "params": {
      "q": "*:*",
      "cursorMark": "*",
      "sort": "id asc",
      "wt": "json"
    }
  },
  "response": {
    "numFound": 509309,
    "start": 0,
    "docs": [
      {
        "name": "测试",
        "id": "00013dd9-7326-43d7-977d-60cdab8deb95",
        "createdBy": "4f070706-c30c-481c-a837-9f39136c62de",
        "createdDate": "2021-05-10T08:14:14Z",
        "lastModifiedBy": "4f070706-c30c-481c-a837-9f39136c62de",
        "lastModifiedDate": "2021-05-28T04:33:22Z",
        "_version_": 1700975244065374200
      }
    ]
  },
  "nextCursorMark": "AoE/BTAwMDEzZGQ5LTczMjYtNDNkNy05NzdkLTYwY2RhYjhkZWI5NQ=="
}
```

## Docker

### Solr Docker

1. `docker pull geerlingguy/solr:4.10.4` [geerlingguy/solr](https://hub.docker.com/r/geerlingguy/solr)
2. `docker run -d --name=solr -p 8983:8983 geerlingguy/solr:4.10.4 /opt/solr/bin/solr start -p 8983 -f`
3. http://127.0.0.1:8983/solr/#/
   ![img.png](_images/img.png)
4. POST requests
    ```http
    POST /solr/collection1/update?commitWithin=1000 HTTP/1.1
    Host: 127.0.0.1:8983
    Content-Type: application/json
    Content-Length: 814
    
    [
        {
            "id": "60bad0a3f46639d20ed3e854",
            "index_i": 0,
            "guid_s": "5e0abb5b-d324-4637-a11e-98a33afc2de1",
            "isActive_b": true,
            "balance_s": "$2,107.27",
            "picture_s": "http://placehold.it/32x32",
            "age_i": 25,
            "eyeColor_s": "green",
            "name_s": "Aline Wiley",
            "gender_s": "female",
            "company_s": "EURON",
            "email_s": "alinewiley@euron.com",
            "phone_s": "+1 (844) 431-2077",
            "address_s": "238 Alabama Avenue, Fillmore, South Dakota, 6551",
            "about_s": "Dolor adipisicing duis anim anim veniam nulla nostrud nulla",
            "registered_s": "2017-11-29T06:18:46 -08:00",
            "greeting_s": "Hello, Aline Wiley! You have 1 unread messages.",
            "favoriteFruit_s": "apple"
        }
    ]
   
    HTTP/1.1 200 OK
    Content-Type: text/plain;charset=UTF-8
    Transfer-Encoding: chunked
    
    {"responseHeader":{"status":0,"QTime":3}}
    ```
5. QUERY: http://localhost:8983/solr/collection1/select?q=*%3A*&rows=1000&wt=json&indent=true

### Elasticsearch docker

1. `docker pull elasticsearch:7.13.1`
   , [docker hub](https://hub.docker.com/_/elasticsearch?tab=description&page=1&ordering=last_updated)
2. `$ docker run -d --name elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.13.1`
3. [chrome extension ElasticSearch Head](https://chrome.google.com/webstore/detail/elasticsearch-head/ffmkiejjmecolpfloofpjologoblkegm)
   , type in `http://127.0.0.1:9200/`

## FAQ

1. [How to insert data to solr collection using Postman?](https://stackoverflow.com/a/49179604)
    - POST request to http://localhost:8983/solr/[collection_name]/update?commitWithin=1000
    - Add the header Content-Type: application/json
    - Make sure to create a json array for multiple documents

    ```json
    [ { "name": "John", "age": 30, "cars": "BMW" }, { "name": "Harry", "age": 30, "cars": "BMW" }, { "name": "Pinku", "age": 30, "cars": "BMW" } ]
    ```
2. [Deleting documents in SOLR](https://gist.github.com/CesarCapillas/a796c0e7cba10ac02213c7f3485d6e90#file-delete-by-id-sh)
    ```sh
    $ curl -X POST "http://127.0.0.1:8983/solr/collection1/update?commit=true&wt=json" -H "Content-Type: text/xml" --data-binary "<delete><id>60bad0a3f46639d20ed3e855</id></delete>"
    {"responseHeader":{"status":0,"QTime":10}}
    ```
3. [ElasticSearch bulk api](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html)
    ```http
    POST /zz/_bulk HTTP/1.1
    Host: 192.168.2.8:9202
    Content-Type: application/json
    Content-Length: 214
    
    { "index" : { "_type" : "docs","_routing":1 } }
    { "zzCode": 1, "field1" : "value10","field2" : "value20" }
    { "index" : { "_type" : "docs","_routing":2 } }
    { "zzCode": 2, "field1" : "value11","field2" : "value21" }
 
    ```

   response:

    ```json
    {"took":95,"errors":false,"items":[{"index":{"_index":"zz","_type":"docs","_id":"jSQE73kBuRYpTrL3mtJh","_version":1,"result":"created","_shards":{"total":2,"successful":2,"failed":0},"_seq_no":0,"_primary_term":1,"status":201}},{"index":{"_index":"zz","_type":"docs","_id":"jiQE73kBuRYpTrL3mtJh","_version":1,"result":"created","_shards":{"total":2,"successful":2,"failed":0},"_seq_no":0,"_primary_term":1,"status":201}}]}
    ```

   performance pk, 10000 dos, bulk mode use 1.7s other than 19s in normal one by one insert mode:

   ```sh
    # solrdump -server 192.168.2.6:8983/solr/licenseIndex -output "192.168.2.8:9202/license/docs?routing=@holderIdentityNum.0" -max 10000
    2021/06/09 13:20:05 started
    2021/06/09 13:20:15 fetched 10000/509311 docs
    2021/06/09 13:20:24 process 10000 docs, rate 517.326652 docs/s, cost 19.330146562s
    # solrdump -server 192.168.2.6:8983/solr/licenseIndex -output "192.168.2.8:9202/license/_bulk?routing=@holderIdentityNum.0" -max 10000
    2021/06/09 13:24:35 started
    2021/06/09 13:24:37 fetched 10000/509311 docs
    2021/06/09 13:24:37 process 10000 docs, rate 5784.721384 docs/s, cost 1.728691727s
    ```
