package controller

import (
	"context"
	"fmt"
	"time"
	"younified-backend/contracts/cms/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *CmsController) GetBlogPosts(ctx context.Context) ([]*model.Blog, error) {
	filter := bson.M{"deleted": false}
	blogs, err := c.CMSRepository.GetBlogs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("No Blog posts found%v", err)
	}
	return blogs, nil
}

func (c *CmsController) GetOneBlogPost(ctx context.Context, blogID primitive.ObjectID) (*model.Blog, error) {
	filter := bson.M{"_id": blogID, "deleted": false}
	blog, err := c.CMSRepository.GetOneBlog(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not find blog post%v", err)
	}
	return blog, nil
}
func (c *CmsController) CreateBlogPost(ctx context.Context, unionID primitive.ObjectID, input model.Blog, image []*string) (*model.Blog, error) {
	if input.CreatedOn.IsZero() {
		input.CreatedOn = time.Now()
	}

	blog, err := c.CMSRepository.CreateBlog(ctx, input)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (c *CmsController) UpdateBlogPost(ctx context.Context, unionID primitive.ObjectID, blogID primitive.ObjectID, input model.Blog) (*model.Blog, error) {
	filter := bson.M{"_id": blogID, "deleted": false}
	updateDoc := bson.M{
		"$set": input,
	}
	blog, err := c.CMSRepository.UpdateBlog(ctx, filter, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("could not find blog post%v", err)
	}
	return blog, nil
}

func (c *CmsController) DeleteBlogPost(ctx context.Context, unionID primitive.ObjectID, blogID primitive.ObjectID) (*string, error) {
	filter := bson.M{"_id": blogID, "deleted": false}
	updateDoc := bson.M{
		"$set": bson.M{"deleted": true},
	}
	_, err := c.CMSRepository.UpdateBlog(ctx, filter, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("could not find blog post to delete%v", err)
	}
	return &Response, nil
}

func (c *CmsController) FeatureBlogPost(ctx context.Context, unionID primitive.ObjectID, blogID primitive.ObjectID, featured bool) (*string, error) {
	filter := bson.M{"_id": blogID, "deleted": false}
	updateDoc := bson.M{
		"$set": bson.M{"featured": featured},
	}
	_, err := c.CMSRepository.UpdateBlog(ctx, filter, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("could not find blog post%v", err)
	}
	return &Response, nil
}
