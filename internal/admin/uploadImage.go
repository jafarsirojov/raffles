package admin

import (
	"crm/internal/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
)

func (s *service) uploadImages(files []multipart.File) (newImagesName []string, err error) {
	for i, _ := range files {

		newUUID := uuid.NewString()
		filename := newUUID + ".png"

		newFile, err := os.Create(structs.ImagePath + filename)
		if err != nil {
			s.logger.Error("internal.admin.uploadImages os.Create", zap.Error(err))
			return newImagesName, err
		}

		fileBytes, err := io.ReadAll(files[i])
		if err != nil {
			s.logger.Error("internal.admin.uploadImages io.ReadAll", zap.Error(err))
			newFile.Close()
			return newImagesName, err
		}

		_, err = newFile.Write(fileBytes)
		if err != nil {
			s.logger.Error("internal.admin.uploadImages io.ReadAll", zap.Error(err))
			newFile.Close()
			return newImagesName, err
		}

		newImagesName = append(newImagesName, filename)

		newFile.Close()
	}

	return newImagesName, nil
}
