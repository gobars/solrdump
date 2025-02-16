# Recreates the bibdata core empty
solr delete -c bibdata
solr create -c bibdata


# Add support for multi-value text in english dynamic fields (*_txts_en)
curl -X POST -H 'Content-type:application/json' --data-binary '{
  "add-dynamic-field":{
     "name":"*_txts_en",
     "type":"text_en",
     "multiValued":true}
}' http://localhost:8983/solr/bibdata/schema


# Copy title to a string field (for sorting)
curl -X POST -H 'Content-type:application/json' --data-binary '{
  "add-copy-field":[
    {
      "source":"title_txt_en",
      "dest":[ "title_s" ]
    }
  ]
}' http://localhost:8983/solr/bibdata/schema


# Create a combined "authors_all_txts_en" field (for searching)
curl -X POST -H 'Content-type:application/json' --data-binary '{
  "add-copy-field":[
    {
      "source":"author_txt_en",
      "dest":[ "authors_all_txts_en" ]
    },
    {
      "source":"authors_other_txts_en",
      "dest":[ "authors_all_txts_en" ]
    }
  ]
}' http://localhost:8983/solr/bibdata/schema


# Copy subject to a (multi-value) string field (for faceting)
curl -X POST -H 'Content-type:application/json' --data-binary '{
  "add-copy-field":[
    {
      "source":"subjects_txts_en",
      "dest": "subjects_s",
      "maxChars": "100"
    }
  ]
}' http://localhost:8983/solr/bibdata/schema


# Copy publisher to a (single-value) string field (for faceting)
curl -X POST -H 'Content-type:application/json' --data-binary '{ 
  "add-copy-field":[
    { 
      "source":"publisher_txt_en", 
      "dest": "publisher_s", 
      "maxChars": "100"
    } 
  ]
}' http://localhost:8983/solr/bibdata/schema
