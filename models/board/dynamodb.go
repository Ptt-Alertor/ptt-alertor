package board

import (
	"encoding/json"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/Ptt-Alertor/ptt-alertor/models/article"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
)

const tableName string = "boards"

// column: Board, Articles
type DynamoDB struct {
}

func (DynamoDB) GetArticles(boardName string) (articles article.Articles) {
	dynamo := dynamodb.New(session.New())
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Board": {
				S: aws.String(boardName),
			},
		},
	})
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Find Board Failed")
		return
	}

	if len(result.Item) == 0 {
		log.WithField("board", boardName).Warn("Board Not Found")
		return
	}

	articlesJSON := aws.StringValue(result.Item["Articles"].S)

	if articlesJSON != "" {
		err = json.Unmarshal([]byte(articlesJSON), &articles)
		if err != nil {
			myutil.LogJSONDecode(err, articlesJSON)
		}
	}
	return articles
}

func (DynamoDB) Save(boardName string, articles article.Articles) error {
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
		return err
	}

	dynamo := dynamodb.New(session.New())
	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Board": {
				S: aws.String(boardName),
			},
			"Articles": {
				S: aws.String(string(articlesJSON)),
			},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Save Board Failed")
	}
	return err
}

func (DynamoDB) Delete(boardName string) error {
	dynamo := dynamodb.New(session.New())
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Board": {
				S: aws.String(boardName),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Delete Board Failed")
	}

	return err
}
