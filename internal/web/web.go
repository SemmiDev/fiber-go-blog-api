package web

import (
	"github.com/SemmiDev/fiber-go-blog/internal/config"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

func GetUrlPathInt64(query string) (int64, error) {
	i, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return 0, constant.ErrUrlPathParameter
	}
	return i, nil
}

func GetPagination(c *fiber.Ctx) (limit, offset int, err error) {
	limitQuery := c.Query("limit")
	offsetQuery := c.Query("offset")

	log.Println(limitQuery)
	log.Println(offsetQuery)

	if limitQuery == "" {
		limit = config.Cfg().PaginationLimit
	} else {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			return 0, 0, err
		} else if limit < 0 {
			return 0, 0, constant.ErrUrlQueryParameter
		} else if limit > config.Cfg().PaginationLimit {
			limit = config.Cfg().PaginationLimit
		}
	}

	if offsetQuery == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			return 0, 0, err
		} else if offset < 0 {
			return 0, 0, constant.ErrUrlQueryParameter
		}
	}

	return
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}