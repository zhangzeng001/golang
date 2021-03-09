package elasPkg

import (
	"context"
	"fmt"
)

func (e *EsObj)Exputdata(data interface{})error{
	indexName := e.IndexName
	_, err := e.client.Index().
		Index(indexName).
		//Id("111").
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Println(err,data)
		return err
	}
	return nil
}


