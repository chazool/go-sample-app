package handler

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func DocumentHandler(file string) func(*fiber.Ctx) error {
	var filePath, templateName string

	switch file {
	case constant.Doc:
		filePath = constant.Basepath + constant.DocumentHtml
		templateName = constant.DocumentHtml
	case constant.Static:
		filePath = constant.Basepath + constant.StaticHtml
		templateName = constant.StaticHtml
	}

	return func(ctx *fiber.Ctx) (err error) {
		html, err := os.ReadFile(strings.TrimSpace(filePath))
		if err != nil {
			utils.Logger.Debug(fmt.Sprintf(constant.FileReadError, filePath), zap.Error(err))
		}

		temp, err := template.New(templateName).Parse(string(html))

		if err != nil {
			utils.Logger.Debug(constant.HTMLTempPassError, zap.Error(err))
		}
		ctx.Type(constant.HTML)
		return temp.Execute(ctx, nil)
	}
}
