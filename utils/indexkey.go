package utils

// key schemes
// primary document: {entityName}#{pk}
// indexed document: {entityName}@{indexName}:{indexValue}#pk
var DocumentSeqDelimiter = []byte(string('#'))
var DocumentIndexKeyDelimiter = []byte(string('@'))
var DocumentIndexValueDelimiter = []byte(string(':'))
var DocumentHeightIndex = []byte(string('h'))

func BuildDocumentKey(entityName, pk []byte) []byte {
	return ConcatBytes(
		entityName,
		DocumentSeqDelimiter,
		pk,
	)
}

func BuildIteratorPrefix(entityName, indexName, indexKey []byte) []byte {
	return ConcatBytes(
		entityName,
		DocumentIndexKeyDelimiter,
		indexName,
		DocumentIndexValueDelimiter,
		indexKey,
		DocumentSeqDelimiter,
	)
}

func BuildIndexedDocumentKey(entityName, indexName, indexKey, pk []byte) []byte {
	return ConcatBytes(
		BuildIteratorPrefix(entityName, indexName, indexKey),
		pk,
	)
}
