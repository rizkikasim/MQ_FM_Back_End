package categoriesadminusecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	categoriesadminentity "mqfm_backend/internal/domain/entities/podcast/categories/admin"
	categoriesadminservice "mqfm_backend/internal/application/services/podcast/categories/admin"
	req "mqfm_backend/internal/domain/request/podcast/categories/admin"

)

type CategoriesAdminUsecase struct {
	service *categoriesadminservice.CategoriesAdminService
}

func NewCategoriesAdminUsecase(service *categoriesadminservice.CategoriesAdminService) *CategoriesAdminUsecase {
	return &CategoriesAdminUsecase{service: service}
}

// ---------------- CREATE ----------------
func (u *CategoriesAdminUsecase) Create(request req.CreateCategoryAdminRequest) (*categoriesadminentity.CategoryAdmin, error) {
	imageFile := ""
	if request.Image != nil {
		uploadDir := "storage/categories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("gagal membuat folder upload: %v", err)
		}

		filename := fmt.Sprintf("category_%s", request.Image.Filename)
		filePath := filepath.Join(uploadDir, filename)

		src, err := request.Image.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file upload: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file upload: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis file upload: %v", err)
		}

		imageFile = filename
	}

	category, err := u.service.Create(
		request.Name,
		request.Description,
		imageFile,
	)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// ---------------- GET ALL ----------------
func (u *CategoriesAdminUsecase) GetAll() ([]*categoriesadminentity.CategoryAdmin, error) {
	return u.service.GetAll()
}

// ---------------- GET BY ID ----------------
func (u *CategoriesAdminUsecase) GetByID(request req.GetCategoryAdminRequest) (*categoriesadminentity.CategoryAdmin, error) {
	return u.service.GetByID(request.ID)
}

// ---------------- UPDATE ----------------
func (u *CategoriesAdminUsecase) Update(request req.UpdateCategoryAdminRequest) (*categoriesadminentity.CategoryAdmin, error) {
	newImage := ""

	if request.NewImage != nil {
		uploadDir := "storage/categories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("gagal membuat folder upload: %v", err)
		}

		filename := fmt.Sprintf("category_%d_%s", request.ID, request.NewImage.Filename)
		filePath := filepath.Join(uploadDir, filename)

		src, err := request.NewImage.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file upload: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file upload: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("gagal menulis file upload: %v", err)
		}

		newImage = filename
	}

	updatedCategory, err := u.service.Update(
		request.ID,
		request.NewName,
		request.NewDesc,
		newImage,
	)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

// ---------------- DELETE ----------------
func (u *CategoriesAdminUsecase) Delete(request req.DeleteCategoryAdminRequest) error {
	return u.service.Delete(request.ID)
}
