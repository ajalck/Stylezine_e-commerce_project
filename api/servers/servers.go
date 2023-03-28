package servers

import (
	"ajalck/e_commerce/api/handler"
	"ajalck/e_commerce/api/middleware"
	"ajalck/e_commerce/payment"

	"github.com/gin-gonic/gin"
)

func UserServer(routes *gin.Engine,
	userHandler handler.UserHandler,
	userAuthHandler handler.UserAuthHandler,
	Middleware middleware.Middleware,
	payment payment.Payment) {

	routes.LoadHTMLGlob("payment/templates/*.html")

	user := routes.Group("/user")
	{
		user.POST("/signup", userHandler.CreateUser)
		user.POST("/login", userAuthHandler.UserSignin)
		user.GET("/listproducts/:page/:records", userHandler.ListProducts)
		user.GET("/checkout/razorpay", payment.RazorPay)
		user.GET("/payment-success", payment.PaymentStatus)

		user.Use(Middleware.AuthorizeJWT)
		{
			wishlist := user.Group("/wishlist")
			{
				wishlist.POST("/add/:productid", userHandler.AddWishlist)
				wishlist.GET("/view/:page/:records", userHandler.ViewWishList)
				{
					wishlist.POST("/addcart/:productid", userHandler.AddCart)
				}
				wishlist.DELETE("/remove/:productid", userHandler.DeleteWishList)
			}
			cart := user.Group("/cart")
			{
				cart.POST("/add/:productid", userHandler.AddCart)
				cart.GET("/view/:page/:records", userHandler.ViewCart)
				cart.DELETE("/remove/:productid", userHandler.DeleteCart)
			}
			coupon := user.Group("/coupon")
			{
				coupon.GET("/listcoupon/:productid", userHandler.ListCoupon)
				coupon.POST("/applycoupon/:cartid/:orderid/:couponid", userHandler.ApplyCoupon)
				coupon.DELETE("/cancelcoupon/:cartid/:orderid/:couponid", userHandler.CancelCoupon)
			}

			shipping := user.Group("/shipping")
			{
				shipping.POST("/adddetails", userHandler.AddShippingDetails)
				shipping.GET("/listdetails", userHandler.ListShippingDetails)
				shipping.DELETE("/deletedetails/:addressid", userHandler.DeleteShippingDetails)
			}
			order := user.Group("/order")
			{
				order.POST("/checkout/:cartid/:productid/:shippingid", userHandler.CheckOut)
				order.GET("/ordersummery", userHandler.OrderSummery)
				order.PATCH("/cancel/:orderid", userHandler.CancelOrder)
			}
		}
	}
}
func AdminServer(routes *gin.Engine,
	adminHandler handler.AdminHandler,
	adminAuthHandler handler.AdminAuthHandler,
	Middleware middleware.Middleware) {

	admin := routes.Group("/admin")
	{
		registration := admin.Group("/registration")
		{
			registration.POST("/signup", adminHandler.CreateAdmin)
			registration.POST("/login", adminAuthHandler.AdminSignin)
		}
		admin.Use(Middleware.AuthorizeJWT)
		{
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
				category.GET("/list", adminHandler.ListCategory)
				category.PATCH("/edit", adminHandler.EditCategory)
				category.DELETE("/delete", adminHandler.DeleteCategory)
			}

			brand := admin.Group("/brandManagement")
			{
				brand.POST("/add", adminHandler.AddBrand)
				brand.GET("/list", adminHandler.ListBrands)
				brand.PATCH("/edit", adminHandler.EditBrand)
				brand.DELETE("/delete", adminHandler.DeleteBrand)
			}

			products := admin.Group("/productManagement")
			{
				products.POST("/add", adminHandler.AddProducts)
				products.GET("/list/:page/:records", adminHandler.ListProducts)
				products.PATCH("/edit", adminHandler.EditProducts)
				products.DELETE("/delete", adminHandler.DeleteProducts)
			}
			coupon := admin.Group("/coupon")
			{
				coupon.POST("/add", adminHandler.AddCoupon)
				coupon.GET("/list/:page/:records", adminHandler.ListCoupon)
				coupon.DELETE("/delete/:couponid", adminHandler.DeleteCoupon)
			}
			admin.GET("/sales_report", adminHandler.SalesReport)
		}

	}

}
