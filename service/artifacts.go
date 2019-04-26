package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
)

type ammoPath struct {
	url    string
	bucket string
	key    string
}

func (a ammoPath) String() string {
	return fmt.Sprintf("%v/%v/%v", a.url, a.bucket, a.key)
}

// UploadAmmo - Заливает в S3 файл с патронами и создает в базе соответствующую запись
func (s *OverloadService) UploadAmmo(req *overload.UploadAmmoRequest) (*overload.UploadAmmoResponse, error) {
	file := strings.NewReader(req.File)
	s3UploadResp, err := s.S3.UploadAmmo(context.Background(), req.Name, file)
	if err != nil {
		log.Println(err)
	}
	url := s3UploadResp.Location
	fmt.Printf("uploaded %v", url)

	// TODO: Author когда прикрутим авторизацию
	ammo := &models.Ammo{
		Url:    s.S3.Client.Endpoint,
		Bucket: s.S3.GetAmmoBucket(),
		// в качестве ключа откусываем от полученного урла последний элемент, предварительно убрав вероятный trailing slash
		Key: url[strings.LastIndex(strings.TrimRight(url, "/"), "/")+1:],
	}
	err = s.Db.Insert(ammo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to create new ammo: %v", err))
	}

	return &overload.UploadAmmoResponse{Url: url}, nil
}

// UploadArtifact - Заливает в S3 файл с артефактом
// добавляет номер джобы и слэш в качестве folder
func (s *OverloadService) UploadArtifact(req *overload.UploadArtifactRequest) (*overload.UploadArtifactResponse, error) {
	file := strings.NewReader(req.File)
	s3UploadResp, err := s.S3.UploadArtifact(context.Background(), req.Job+"/"+req.Name, file)
	if err != nil {
		log.Println(err)
	}
	url := s3UploadResp.Location
	fmt.Printf("uploaded %v", url)

	return &overload.UploadArtifactResponse{Url: url}, nil
}

func (s *OverloadService) ListAmmo(req *overload.ListAmmoRequest) (*overload.ListAmmoResponse, error) {

	var (
		ammoKeys     []string
		ammoFromDB   []*models.Ammo
		preparedAmmo []*overload.Ammo
	)
	ammoByPath := make(map[string]*models.Ammo)

	ammoFromS3, err := s.S3.ListAmmo(context.Background())
	if err != nil {
		log.Println(err)
	}

	for _, as3 := range ammoFromS3 {
		ammoKeys = append(ammoKeys, *as3.Key)
	}

	err = s.Db.
		Model(&ammoFromDB).
		Where("url = ?", s.S3.Client.Endpoint).
		Where("bucket = ?", s.S3.GetAmmoBucket()).
		Where("key in (?)", ammoKeys).
		Select()

	for _, adb := range ammoFromDB {
		path := fmt.Sprint(ammoPath{
			s.S3.Client.Endpoint,
			s.S3.GetAmmoBucket(),
			adb.Key,
		})
		ammoByPath[path] = adb
	}

	for _, as3 := range ammoFromS3 {
		path := fmt.Sprint(ammoPath{
			s.S3.Client.Endpoint,
			s.S3.GetAmmoBucket(),
			*as3.Key,
		})
		ammo := &overload.Ammo{
			Etag:         *as3.ETag,
			Key:          *as3.Key,
			LastModified: as3.LastModified.String(),
			Size_:        *as3.Size,
			Path:         path,
		}

		if adb, ok := ammoByPath[path]; ok {
			ammo.Author = adb.Author
			ammo.LastUsed = adb.LastUsed
			ammo.Type = adb.Type
		}

		preparedAmmo = append(preparedAmmo, ammo)
	}

	return &overload.ListAmmoResponse{
		Ammo: preparedAmmo,
	}, nil
}

func (s *OverloadService) DeleteAmmo(req *overload.DeleteAmmoRequest) (*overload.DeleteAmmoResponse, error) {
	key := req.Key
	if key == "" {
		return nil, errors.New("Key param is mandatory")
	}
	err := s.S3.DeleteAmmo(context.Background(), key)
	if err != nil {
		log.Println(err)
		return nil, errors.New(fmt.Sprintf("Failed to delete ammo: %v", err))
	}
	return &overload.DeleteAmmoResponse{}, err
}

func (s *OverloadService) ListArtifacts(req *overload.ListArtifactsRequest) (*overload.ListArtifactsResponse, error) {
	var (
		preparedArtifacts []*overload.Artifact
		folderPrefix      string
	)
	folderPrefix = strconv.Itoa(int(req.Job)) + "/"
	artifactsFromS3, err := s.S3.ListArtifacts(context.Background(), folderPrefix)
	if err != nil {
		log.Println(err)
	}

	for _, as3 := range artifactsFromS3 {
		path := fmt.Sprint(ammoPath{
			s.S3.Client.Endpoint,
			s.S3.GetArtifactsBucket(),
			*as3.Key,
		})
		artifact := &overload.Artifact{
			Etag:         *as3.ETag,
			Key:          *as3.Key,
			LastModified: as3.LastModified.String(),
			Size_:        *as3.Size,
			Path:         path,
			Job:          req.Job,
		}
		preparedArtifacts = append(preparedArtifacts, artifact)
	}

	return &overload.ListArtifactsResponse{
		Artifacts: preparedArtifacts,
	}, nil
}
