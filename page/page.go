package page

import (
	"strconv"

	"github.com/hyacinthus/x/xerr"
	"github.com/labstack/echo/v4"
)

// TODO: 无差别通过中间件分页在性能上不划算，因为需要分页的接口比例不高。此方式已废弃，不要再使用。

// Parse 帮助解析分页中间件的计算值
func Parse(c echo.Context) (offset, limit int) {
	offset = c.Get("offset").(int)
	limit = c.Get("limit").(int)
	return
}

// Middleware 获得页码，每页条数，Echo中间件。
func Middleware(defaultSize int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			var page, pageSize int
			// 获得页码
			if c.QueryParam("page") == "" {
				page = 1
			} else {
				if page, err = strconv.Atoi(c.QueryParam("page")); err != nil {
					return xerr.New(400, "InvalidPage", "请在URL中提供合法的页码")
				}
			}
			// 获得每页条数
			if c.QueryParam("per_page") == "" {
				pageSize = defaultSize
			} else {
				if pageSize, err = strconv.Atoi(c.QueryParam("per_page")); err != nil {
					return xerr.New(400, "InvalidPage", "请在URL中提供合法的每页条数")
				}
			}
			// 设置查询数据时的 offset 和 limit
			c.Set("page", page)
			c.Set("offset", (page-1)*pageSize)
			c.Set("limit", pageSize)
			// 设置返回的Header
			c.Response().Header().Set("X-Page-Num", strconv.Itoa(page))
			c.Response().Header().Set("X-Page-Size", strconv.Itoa(pageSize))
			return next(c)
		}
	}
}
