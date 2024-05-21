package controllers

import (
	"activity-tracker-api/models/activity"
	"activity-tracker-api/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"sort"
	"sync"
	"time"
)

type ActivityController struct {
	ActivityRepo    storage.ActivityRepository
	S3Client        *s3.Client
	S3PresignClient *s3.PresignClient
}

func (controller *ActivityController) PostActivityHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	jsonPart := form.Value["activity"]
	if len(jsonPart) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("JSON data is missing.")
	}

	var request activity.Activity
	if err := json.Unmarshal([]byte(jsonPart[0]), &request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON data.")
	}

	activityUUID := uuid.New().String()
	imagePart := form.File["image"]
	if len(imagePart) != 0 {
		bucketName := os.Getenv("S3_BUCKET_NAME")
		imageFile := imagePart[0]
		file, err := imageFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to open image file.")
		}

		defer file.Close()

		imageKey := "images/" + activityUUID

		_, err = controller.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(imageKey),
			Body:        file,
			ContentType: aws.String(imageFile.Header.Get("Content-Type")),
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to upload image to S3.")
		}

		presignedS3Url, err := generatePresignedURL(controller.S3PresignClient, bucketName, imageKey, 60)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate presigned URL.")
		}

		request.Id = activityUUID
		request.ImageUrl = presignedS3Url
	}

	dbActivity := request.ToDbActivity(activityUUID, userId)

	err = controller.ActivityRepo.Insert(&dbActivity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Can't save activity.")
	}

	return c.Status(fiber.StatusOK).JSON(request)
}

func (controller *ActivityController) GetActivityHandler(c *fiber.Ctx) error {
	activityId := c.Params("id")

	dbActivity, err := controller.ActivityRepo.GetByID(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	imageKey := fmt.Sprintf("images/%s", dbActivity.Id)
	bucketName := os.Getenv("S3_BUCKET_NAME")
	url, err := generatePresignedURL(controller.S3PresignClient, bucketName, imageKey, 2*60)

	return c.Status(fiber.StatusOK).JSON(dbActivity.ToActivity(url))
}

func (controller *ActivityController) DeleteActivityHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)
	activityId := c.Params("id")

	err := controller.ActivityRepo.Delete(activityId, userId)
	if err != nil {
		log.Debug(err)
		return c.Status(fiber.StatusBadRequest).SendString("Can't delete activity.")
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (controller *ActivityController) GetActivitiesHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)

	dbActivities, err := controller.ActivityRepo.GetForUserId(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Can't get activities.")
	}

	var wg sync.WaitGroup
	urlChan := make(chan struct {
		index int
		url   string
		err   error
	}, len(dbActivities))

	for i, dbActivity := range dbActivities {
		wg.Add(1)

		go func(i int, dbActivity activity.DbActivity) {
			defer wg.Done()
			imageKey := fmt.Sprintf("images/%s", dbActivity.Id)
			bucketName := os.Getenv("S3_BUCKET_NAME")
			url, err := generatePresignedURL(controller.S3PresignClient, bucketName, imageKey, 15*60)
			urlChan <- struct {
				index int
				url   string
				err   error
			}{index: i, url: url, err: err}
		}(i, dbActivity)
	}

	go func() {
		wg.Wait()
		close(urlChan)
	}()

	presignedUrls := make(map[int]string, len(dbActivities))
	for result := range urlChan {
		if result.err != nil {
			continue
		}

		presignedUrls[result.index] = result.url
	}

	var activities []activity.Activity
	for i, dbActivity := range dbActivities {
		imageUrl := presignedUrls[i]
		activities = append(activities, dbActivity.ToActivity(imageUrl))
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].StartTimestamp > activities[j].StartTimestamp
	})

	return c.Status(fiber.StatusOK).JSON(activities)
}

func generatePresignedURL(presignClient *s3.PresignClient, bucketName, imageKey string, expirationSeconds int64) (string, error) {
	presignedS3Request, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expirationSeconds * int64(time.Second))
	})

	if err != nil {
		return "", err
	}

	return presignedS3Request.URL, err
}
