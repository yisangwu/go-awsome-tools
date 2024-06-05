package main

// 待完善

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	// 源 Elasticsearch 配置
	sourceConfig := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	// 目标 Elasticsearch 配置
	targetConfig := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	// 创建源 Elasticsearch 客户端
	sourceClient, err := elasticsearch.NewClient(sourceConfig)
	if err != nil {
		log.Fatalf("Failed to create source Elasticsearch client: %s", err)
	}

	// 创建目标 Elasticsearch 客户端
	targetClient, err := elasticsearch.NewClient(targetConfig)
	if err != nil {
		log.Fatalf("Failed to create target Elasticsearch client: %s", err)
	}

	// 获取源索引中最新文档的时间戳
	sourceTimestamp, err := getSourceTimestamp(sourceClient, "source_index", "timestamp")
	if err != nil {
		log.Fatalf("Failed to get source timestamp: %s", err)
	}

	// 从源索引中获取时间戳之后的所有文档
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": sourceTimestamp,
				},
			},
		},
	}

	scroll := elasticsearch.NewScrollService(sourceClient)
	scroll.Body = bytes.NewReader(serializeQuery(query))

	for {
		response, err := scroll.Do(context.Background())
		if err != nil {
			log.Fatalf("Failed to execute scroll: %s", err)
		}

		if len(response.Hits.Hits) == 0 {
			break // 没有更多结果了
		}

		// 将结果写入目标索引
		bulkRequest := make([]map[string]interface{}, 0)
		for _, hit := range response.Hits.Hits {
			doc := map[string]interface{}{
				"index": map[string]interface{}{
					"_index": hit.Index,
				},
			}
			bulkRequest = append(bulkRequest, doc)
			bulkRequest = append(bulkRequest, hit.Source)
		}

		targetResponse, err := bulkInsert(targetClient, bulkRequest)
		if err != nil {
			log.Fatalf("Failed to bulk insert into target index: %s", err)
		}
		if targetResponse.IsError() {
			log.Fatalf("Error response: %s", targetResponse.String())
		}
	}

	log.Println("Data transfer complete!")
}

// 获取源索引中最新文档的时间戳
func getSourceTimestamp(client *elasticsearch.Client, index, timestampField string) (string, error) {
	query := map[string]interface{}{
		"sort": []map[string]interface{}{
			{
				timestampField: map[string]string{
					"order": "desc",
				},
			},
		},
		"size": 1,
	}

	response, err := client.Search(
		client.Search.WithIndex(index),
		client.Search.WithBody(bytes.NewReader(serializeQuery(query))),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return "", nil
	}

	source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
	return source[timestampField].(string), nil
}

// 将批量请求写入目标索引
func bulkInsert(client *elasticsearch.Client, request []map[string]interface{}) (*esapi.Response, error) {
	body := bytes.NewBufferString("")
	for _, req := range request {
		header, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		body.Write(header)
		body.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body: bytes.NewReader(body.Bytes()),
	}
	return req.Do(context.Background(), client)
}

// 序列化查询
func serializeQuery(query map[string]interface{}) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(query)
	return buf.Bytes()
}
