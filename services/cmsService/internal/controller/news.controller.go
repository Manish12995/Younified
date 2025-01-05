package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"
	"younified-backend/contracts/cms/model"
	"younified-backend/providers/aws"
	"younified-backend/providers/database"
	"younified-backend/providers/graphqlclient"
	"younified-backend/services/cmsService/internal/repository"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CmsController struct {
	CMSRepository       *repository.MongoCommsRepository
	dbManager           *database.DBManager
	graphqlManager      *graphqlclient.Graph
	NewsRedisRepository *repository.RedisUserRepository
	awsProvider         *aws.AWSProvider
}

var Response string = "Operation Successful"

func NewCMSController(dbManager *database.DBManager, graphqlManager *graphqlclient.Graph, redisClient *database.RedisClient, awsProvider *aws.AWSProvider) *CmsController {
	if dbManager == nil {
		panic("dbManager cannot be nil")
	}

	return &CmsController{
		CMSRepository:       repository.NewMongoCommsRepository(dbManager, "unified_base"),
		dbManager:           dbManager,
		graphqlManager:      graphqlManager,
		NewsRedisRepository: repository.NewRedisUserRepository(redisClient),
		awsProvider:         awsProvider,
	}
}

func (c *CmsController) GetAllNewsPosts(ctx context.Context, unionID primitive.ObjectID, page int, limit int) (*model.Report, error) {
	if unionID.IsZero() {
		err := fmt.Errorf("unionID is required to process")
		log.Tracef("Invalid unionID %v", unionID)
		return nil, err
	}
	filter := bson.M{"deleted": false}
	sort := bson.M{"pinned": -1}
	NewsReport := model.Report{}
	Data, Total, err := c.CMSRepository.GetAllNewsPosts(ctx, unionID.Hex(), filter, page, limit, sort)
	if err != nil {
		err = fmt.Errorf("could not find newsfeed")
		return nil, err
	}
	NewsReport.Data, NewsReport.Total = Data, Total

	//get comments and attacch to News
	var comments = []*model.Comment{}
	for _, news := range NewsReport.Data {
		log.Tracef("Get comment for %+v", news.ID)
		comments, err = c.CMSRepository.RetrieveComments(ctx, bson.M{"deleted": false}, unionID.Hex(), news.ID.Hex(), 1, 5)
		if err != nil {
			err = fmt.Errorf("could not get the comments due to %v", err)
			return nil, err
		}
		news.CommentCount = int64(len(comments))
		news.Comments = append(news.Comments, comments...)
	}
	return &NewsReport, nil
}

func (c *CmsController) GetComments(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, page int, limit int) ([]*model.Comment, error) {
	filter := bson.M{"deleted": false}
	return c.CMSRepository.RetrieveComments(ctx, filter, unionID.Hex(), newsID.Hex(), page, limit)
}

func (c *CmsController) DeleteNews(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID) (*string, error) {
	filter := bson.M{"_id": newsID}
	update := bson.M{"deleted": true}
	_, err := c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
	if err != nil {
		return nil, err
	}
	return &Response, nil
}

func (c *CmsController) LikeNewsItem(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, userID primitive.ObjectID) (*model.News, error) {
	filter := bson.M{"_id": newsID}
	news, err := c.CMSRepository.GetOneNews(ctx, unionID.Hex(), newsID)
	if err != nil {
		return nil, err
	}
	if slices.Contains(news.Likes, userID) {
		likeIndex := slices.Index(news.Likes, userID)
		updated := slices.Delete(news.Likes, likeIndex, likeIndex+1)
		result, er := c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, bson.M{"$set": bson.M{"likes": updated}})
		if err != nil {
			return nil, er
		}
		return result, nil
	} else {
		filter := bson.M{"_id": newsID}
		update := bson.M{"$addToSet": bson.M{"likes": userID}}
		result, er := c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
		if er != nil {
			return nil, err
		}
		return result, nil
	}
}

func (c *CmsController) AddComment(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, input model.Comment) (*model.Comment, error) {
	if unionID.IsZero() {
		err := fmt.Errorf("unionID is required to process")
		return nil, err
	}
	input.Likes = []primitive.ObjectID{}
	input.Dislikes = []primitive.ObjectID{}
	return c.CMSRepository.AddComment(ctx, unionID.Hex(), newsID, &input)
}

func (c *CmsController) LikeButtonToggle(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, input bool) (*model.News, error) {
	if unionID.IsZero() || newsID.IsZero() {
		err := fmt.Errorf("unionID and newsId are required")
		return nil, err
	}
	filter := bson.M{"_id": newsID}
	update := bson.M{"$set": bson.M{"showLikes": input}}
	return c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
}

func (c *CmsController) CommentButtonToggle(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, input bool) (*model.News, error) {
	if unionID.IsZero() || newsID.IsZero() {
		err := fmt.Errorf("unionID and newsId are required")
		return nil, err
	}
	filter := bson.M{"_id": newsID}
	update := bson.M{"$set": bson.M{"showComments": input}}
	return c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
}

func (c *CmsController) MakeNewsPrivate(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, input bool) (*model.News, error) {
	if unionID.IsZero() || newsID.IsZero() {
		err := fmt.Errorf("unionID and newsId are required")
		return nil, err
	}
	filter := bson.M{"_id": newsID}
	update := bson.M{"$set": bson.M{"private": input}}
	return c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
}

func (c *CmsController) ShowPinOption(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, input bool) (*model.News, error) {
	if unionID.IsZero() || newsID.IsZero() {
		err := fmt.Errorf("unionID and newsId are required")
		return nil, err
	}
	filter := bson.M{"_id": newsID}
	update := bson.M{"$set": bson.M{"show": input}}
	return c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
}

func (c *CmsController) PinNewsPost(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID) (*model.News, error) {
	alreadyPinned := false

	items, _, err := c.CMSRepository.GetAllNewsPosts(ctx, unionID.Hex(), bson.M{"pinned": true, "deleted": false}, 1, 20, bson.M{"createdOn": -1})
	if err != nil {
		log.Errorf("Could not pin news item %+v", err)
		return nil, err
	}
	if len(items) > 9 {
		log.Errorf("Maximum number of pinned document exceded")
		err = fmt.Errorf("Maximum number of pinned document exceded")
		return nil, err
	}

	for _, v := range items {
		item := *v
		if item.ID == newsID {
			alreadyPinned = true
			break
		}
	}
	//setting filter and update
	filter := bson.M{"_id": newsID}
	update := bson.M{"$set": bson.M{"pinned": true, "pinnedAt": time.Now(), "show": true}}

	if alreadyPinned {
		update = bson.M{"$set": bson.M{"pinned": false}}
		unPinnedNews, er := c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
		if er != nil {
			log.Errorf("Could not unpin document %+v", err)
			return nil, er
		}
		return unPinnedNews, nil
	}

	pinnedNews, erro := c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
	if erro != nil {
		log.Errorf("Could not pin news item %+v", err)
		return nil, err
	}
	return pinnedNews, nil
}

func (c *CmsController) LikeComment(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, commentID primitive.ObjectID, userID primitive.ObjectID) (*model.Comment, error) {
	filter := bson.M{"_id": commentID}
	comment, err := c.CMSRepository.GetOneComment(ctx, unionID.Hex(), newsID, commentID)
	if err != nil {
		return nil, err
	}
	if slices.Contains(comment.Likes, userID) {
		likeIndex := slices.Index(comment.Likes, userID)
		updated := slices.Delete(comment.Likes, likeIndex, likeIndex+1)
		result, er := c.CMSRepository.UpdateOneComment(ctx, unionID.Hex(), newsID, filter, bson.M{"$set": bson.M{"likes": updated}})
		if er != nil {
			return nil, fmt.Errorf("could not find comment%v", err)
		}
		return result, nil
	} else {
		filter := bson.M{"_id": comment}
		update := bson.M{"$addToSet": bson.M{"likes": userID}}
		result, er := c.CMSRepository.UpdateOneComment(ctx, unionID.Hex(), newsID, filter, update)
		if er != nil {
			return nil, fmt.Errorf("could not find comment%v", er)
		}
		return result, nil
	}
}

func (c *CmsController) DeleteComment(ctx context.Context, unionID primitive.ObjectID, newsID primitive.ObjectID, commentID primitive.ObjectID) (*string, error) {
	filter := bson.M{"_id": commentID}
	update := bson.M{"$set": bson.M{"deleted": true}}
	_, err := c.CMSRepository.UpdateOneComment(ctx, unionID.Hex(), newsID, filter, update)
	if err != nil {
		err = fmt.Errorf("could not delete the comment")
		return nil, err
	}
	var response string = "Deleted Successfully"
	return &response, nil
}

func (c *CmsController) CreateNews(ctx context.Context, unionID primitive.ObjectID, input model.News, images []*string, documents []*model.Document, category string) (*model.News, error) {
	if input.CreatedOn.IsZero() {
		input.CreatedOn = time.Now()
	}

	input.Likes = []primitive.ObjectID{}
	input.Dislikes = []primitive.ObjectID{}
	input.Category = "news"
	input.Images = append(input.Images, *images[0])
	input.Documents = append(input.Documents, documents...)

	news, err := c.CMSRepository.CreateNews(ctx, unionID.Hex(), input)
	if err != nil {
		return nil, err
	}

	// urls, err := c.uploadMedia(ctx, images, "img", unionID, category, news.ID)
	// fmt.Println(urls)
	// if err != nil {
	// 	logrus.Tracef("could not upload news post image on bucket: %v", err)
	// 	err = fmt.Errorf("news created with errors in uploading image %v", err)
	// 	return nil, err
	// }

	// if len(urls) > 0 {
	// 	filter := bson.M{"_id": news.ID}
	// 	update := bson.M{"$set": bson.M{"images": urls}}
	// 	news, err = c.CMSRepository.UpdateNews(ctx, unionID.Hex(), filter, update)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	return news, nil
}

// uploadMedia uploads images to S3 and updates the database
func (c *CmsController) uploadMedia(ctx context.Context, files []*graphql.Upload, mediaType string, unionID primitive.ObjectID, category string, newsID primitive.ObjectID) ([]string, error) {
	if len(files) == 0 {
		return nil, nil
	}
	var urls []string
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		filename := uuid.New()
		dir := unionID.Hex() + "/" + category + "/" + filename.String() + extension

		var buf bytes.Buffer

		// Copy the file content to the buffer
		_, err := io.Copy(&buf, file.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read upload: %w", err)
		}
		url, err := c.awsProvider.UploadToS3(ctx, os.Getenv("AWS_S3_BUCKET"), os.Getenv("AWS_REGION"), dir, buf.Bytes())
		// url, err := data.SingleUploadS3(unionID.Hex(), "news", *file)
		if err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", mediaType, err)
		}
		urls = append(urls, *url)
	}
	return urls, nil
}
