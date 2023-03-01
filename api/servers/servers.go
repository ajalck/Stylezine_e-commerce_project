package servers

import (
	"ajalck/e_commerce/api/handler"
	"ajalck/e_commerce/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserServer(routes *gin.Engine,
	userHandler handler.UserHandler,
	userAuthHandler handler.UserAuthHandler,
	userMiddleware middleware.UserMiddleware) {
	user := routes.Group("/user")
	{
		user.POST("/signup", userHandler.CreateUser)
		user.POST("/login", userAuthHandler.UserSignin)
		user.GET("/listproducts/:page/:records", userHandler.ListProducts)

		user.Use(userMiddleware.AuthorizeJWT)
		{
			wishlist := user.Group("/wishlist")
			{
				wishlist.POST("/add/:userid/:productid", userHandler.AddWishlist)
				wishlist.GET("/view/:userid/:page/:records", userHandler.ViewWishList)
				wishlist.DELETE("/add/:userid/:productid", userHandler.DeleteWishList)
			}
		}

	}

}
func AdminServer(routes *gin.Engine,
	adminHandler handler.AdminHandler,
	adminAuthHandler handler.AdminAuthHandler) {

	admin := routes.Group("/admin")
	{
		registration := admin.Group("/registration")
		{
			registration.POST("/signup", adminHandler.CreateAdmin)
			registration.POST("/login", adminAuthHandler.AdminSignin)
		}

		users := admin.Group("/userManagement")
		{
			users.GET("/listusers/:page/:records", adminHandler.ListUsers)
			users.GET("/viewuser/:id", adminHandler.ViewUser)
			users.PATCH("/blockuser/:id", adminHandler.BlockUser)
			users.PATCH("/unblockuser/:id", adminHandler.UnblockUser)
			users.GET("/list/blockedusers/:page/:records", adminHandler.ListBlockedUsers)
			users.GET("/list/activeusers/:page/:records", adminHandler.ListActiveUsers)
		}

		category := admin.Group("/categoryManagement")
		{
			category.POST("/add", adminHandler.AddCategory)
			category.PATCH("/edit", adminHandler.EditCategory)
			category.DELETE("/delete", adminHandler.DeleteCategory)
		}

		brand := admin.Group("/brandManagement")
		{
			brand.POST("/add", adminHandler.AddBrand)
			brand.PATCH("/edit", adminHandler.EditBrand)
			brand.DELETE("/delete", adminHandler.DeleteBrand)
		}

		products := admin.Group("/productManagement")
		{
			products.POST("/add", adminHandler.AddProducts)
			products.PATCH("/edit", adminHandler.EditProducts)
			products.DELETE("/delete", adminHandler.DeleteProducts)
		}
	}

}
