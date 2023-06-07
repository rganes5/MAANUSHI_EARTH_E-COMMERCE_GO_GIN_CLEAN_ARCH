package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) {
	//sets up a route group for the "/user" endpoint
	signup := api.Group("/user")
	{
		signup.POST("/signup", userHandler.UserSignUp)
		signup.POST("/signup/otp/verify", userHandler.SignupOtpverify)
	}
	login := api.Group("/user")
	{
		login.POST("/login", userHandler.LoginHandler)

	}
	forgotPassword := api.Group("/user/forgot/password")
	{
		forgotPassword.POST("/", userHandler.ForgotPassword)
		forgotPassword.PATCH("/otp/verify", userHandler.ForgotPasswordOtpVerify)
	}
	home := login.Group("/")
	{
		//AuthorizationMiddleware as middleware to perform authorization checks for users accessing the "/user" endpoint.
		home.Use(middleware.AuthorizationMiddleware("user"))
		home.POST("/logout", userHandler.LogoutHandler)
		category := home.Group("/category")
		{
			category.GET("/listall", productHandler.ListCategories)
			// category.GET("/list/categorybased/products", productHandler.ListCategoryBasedProducts)
		}
		product := home.Group("/products")
		{
			product.GET("/listall", productHandler.ListProducts)
			product.GET("/listbasedoncategory/:category_id", productHandler.ListProductsBasedOnCategory)
			product.GET("/findproductdetails/:product_id", productHandler.ListProductDetailsById)
			product.GET("/findproductanddetails/:product_id", productHandler.ListProductAndDetailsById)
		}
		cart := home.Group("/cart")
		{
			cart.POST("/add/:product_id", cartHandler.AddToCart)
			cart.POST("/remove/:product_id", cartHandler.RemoveFromCart)
			cart.GET("/list", cartHandler.ListCart)
		}
		checkout := home.Group("/checkout")
		{
			checkout.POST("/placeorder", orderHandler.PlaceNewOrder)
		}
		orders := home.Group("/orders")
		{
			orders.GET("/list/all", orderHandler.ListOrders)
			orders.GET("/list/details/:order_id", orderHandler.ListOrderDetails)
			orders.POST("/cancel/:order_details_id", orderHandler.CancelOrder)
		}
		userprofile := home.Group("/profile")
		{
			userprofile.GET("/home", userHandler.HomeHandler)
			userprofile.PATCH("/edit/profile", userHandler.UpdateProfile)
			userprofile.POST("/add/address", userHandler.AddAddress)
			userprofile.GET("/list/address", userHandler.ListAddress)
			userprofile.PATCH("/edit/address/:address_id", userHandler.UpdateAddress)
			userprofile.POST("/delete/address/:address_id", userHandler.DeleteAddress)
		}
	}

}
