// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package rewrite

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func Test_New(t *testing.T) {
	// Test with no config
	m := New()

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}

	// Test with config
	m = New(Config{
		Rules: map[string]string{
			"/old": "/new",
		},
	})

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}

	// Test with full config
	m = New(Config{
		Filter: func(*fiber.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	})

	if m == nil {
		t.Error("Expected middleware to be returned, got nil")
	}
}

func Test_Rewrite(t *testing.T) {
	// Case 1: filter function always returns true
	app := fiber.New()
	app.Use(New(Config{
		Filter: func(*fiber.Ctx) bool {
			return true
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	}))

	app.Get("/old", func(c *fiber.Ctx) error {
		return c.SendString("Rewrite Successful")
	})

	req, _ := http.NewRequest("GET", "/old", nil)
	resp, err := app.Test(req)
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)

	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "Rewrite Successful", bodyString)

	// Case 2: filter function always returns false
	app = fiber.New()
	app.Use(New(Config{
		Filter: func(*fiber.Ctx) bool {
			return false
		},
		Rules: map[string]string{
			"/old": "/new",
		},
	}))

	app.Get("/new", func(c *fiber.Ctx) error {
		return c.SendString("Rewrite Successful")
	})

	req, _ = http.NewRequest("GET", "/old", nil)
	resp, err = app.Test(req)
	body, _ = io.ReadAll(resp.Body)
	bodyString = string(body)

	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "Rewrite Successful", bodyString)

	// Case 3: check for captured tokens in rewrite rule
	app = fiber.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	req, _ = http.NewRequest("GET", "/users/123/orders/456", nil)
	resp, err = app.Test(req)
	body, _ = io.ReadAll(resp.Body)
	bodyString = string(body)

	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "User ID: 123, Order ID: 456", bodyString)

	// Case 4: Send non-matching request, handled by default route
	app = fiber.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req, _ = http.NewRequest("GET", "/not-matching-any-rule", nil)
	resp, err = app.Test(req)
	body, _ = io.ReadAll(resp.Body)
	bodyString = string(body)

	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)
	utils.AssertEqual(t, "OK", bodyString)

	// Case 4: Send non-matching request, with no default route
	app = fiber.New()
	app.Use(New(Config{
		Rules: map[string]string{
			"/users/*/orders/*": "/user/$1/order/$2",
		},
	}))

	app.Get("/user/:userID/order/:orderID", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("User ID: %s, Order ID: %s", c.Params("userID"), c.Params("orderID")))
	})

	req, _ = http.NewRequest("GET", "/not-matching-any-rule", nil)
	resp, err = app.Test(req)
	body, _ = io.ReadAll(resp.Body)
	bodyString = string(body)

	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)
}
