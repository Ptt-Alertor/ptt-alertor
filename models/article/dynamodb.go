package article

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/myutil"
)

// table: code, board, content
type DynamoDB struct{}

const tableName string = "articles"

func (DynamoDB) Find(code string) Article {
	dynamo := dynamodb.New(session.New())
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Code": {
				S: aws.String(code),
			},
		},
	})
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Find Article Failed")
		return Article{}
	}

	if len(result.Item) == 0 {
		log.WithField("code", code).Warn("Article Not Found")
		return Article{}
	}

	a := Article{
		Code:   aws.StringValue(result.Item["Code"].S),
		Title:  aws.StringValue(result.Item["Title"].S),
		Link:   aws.StringValue(result.Item["Link"].S),
		Date:   aws.StringValue(result.Item["Date"].S),
		Author: aws.StringValue(result.Item["Author"].S),
		Board:  aws.StringValue(result.Item["Board"].S),
	}
	if err := dynamodbattribute.Unmarshal(result.Item["ID"], &a.ID); err != nil {
		log.WithFields(log.Fields{
			"code": code,
			"id":   result.Item["ID"],
		}).WithError(err).Warn("Article ID Unmarshal Failed")
	}
	if err := dynamodbattribute.Unmarshal(result.Item["PushSum"], &a.PushSum); err != nil {
		log.WithFields(log.Fields{
			"code":    code,
			"pushSum": result.Item["PushSum"],
		}).WithError(err).Warn("Article PushSum Unmarshal Failed")
	}
	if a.LastPushDateTime, err = time.Parse(time.RFC3339, aws.StringValue(result.Item["LastPushDateTime"].S)); err != nil {
		log.WithFields(log.Fields{
			"code":             code,
			"lastPushDateTime": result.Item["LastPushDateTime"],
		}).WithError(err).Warn("Article LastPushDateTime Unmarshal Failed")
	}
	comments := aws.StringValue(result.Item["Comments"].S)
	if err = json.Unmarshal([]byte(comments), &a.Comments); err != nil {
		log.WithFields(log.Fields{
			"code":     code,
			"comments": result.Item["Comments"],
		}).Warn("Article Comments Unmarshal Failed")
		myutil.LogJSONDecode(err, comments)
	}

	return a
}

func (DynamoDB) Save(a Article) error {
	commentsJSON, err := json.Marshal(a.Comments)
	if err != nil {
		myutil.LogJSONEncode(err, a)
		return err
	}
	dynamo := dynamodb.New(session.New())
	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(strconv.Itoa(a.ID)),
			},
			"Code": {
				S: aws.String(a.Code),
			},
			"Title": {
				S: aws.String(a.Title),
			},
			"Link": {
				S: aws.String(a.Link),
			},
			"Date": {
				S: aws.String(a.Date),
			},
			"Author": {
				S: aws.String(a.Author),
			},
			"Comments": {
				S: aws.String(string(commentsJSON)),
			},
			"LastPushDateTime": {
				S: aws.String(a.LastPushDateTime.Format(time.RFC3339)),
			},
			"Board": {
				S: aws.String(a.Board),
			},
			"PushSum": {
				N: aws.String(strconv.Itoa(a.PushSum)),
			},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Save Article Failed")
	}
	return err
}

func (DynamoDB) Delete(code string) error {
	dynamo := dynamodb.New(session.New())
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Code": {
				S: aws.String(code),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error("DynamoDB Delete Article Failed")
	}
	return err
}
