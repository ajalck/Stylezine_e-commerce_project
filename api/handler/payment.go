package handler

// import (
// 	"ajalck/e_commerce/domain"
// 	"ajalck/e_commerce/utils"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	razorpay "github.com/razorpay/razorpay-go"
// )

// type Home struct {
// 	userid      string
// 	Name        string
// 	total_price int
// 	Amount      int
// 	OrderId     string
// 	Email       string
// 	Contact     string
// }

// func (cr *UserHandler) RazorPay(c *gin.Context) {
// 	//email := c.Query("email")
// 	email := "ack6627@gmail.com"
// 	//coupon := c.Query("coupon")
// 	coupon := "YATSTJ61"
// 	address_id := 1
// 	user, err := cr.AuthService.FindUser(email)
// 	if err != nil {
// 		c.HTML(400, "failed to find user", nil)
// 		c.Abort()
// 		return
// 	}

// 	Coupn, _ := cr.UserService.FindCoupon(coupon)

// 	if Coupn.Validity < time.Now().Unix() {
// 		Coupn.Discount = 0
// 		Coupn.Coupon = "no value"
// 		// res := response.ErrorResponse("your coupon expired", "your coupon expired", "enter anothere coupon")
// 		// //c.Writer.WriteHeader(300)
// 		// utils.ResponseJSON(c, res)
// 	}

// 	sum, err := cr.UserService.FindTheSumOfCart(int(user.ID))
// 	fmt.Println(sum)
// 	if err != nil {

// 		c.HTML(300, "fdisaoifdjaio", sum)
// 	}
// 	client := razorpay.NewClient("rzp_test_0h8oXuKI0kORyw", "McmgVREukL239BhjpTuS4j3t")
// 	razorpaytotal := (sum - Coupn.Discount) * 100
// 	data := map[string]interface{}{
// 		"amount":   razorpaytotal,
// 		"currency": "INR",
// 	}
// 	body, err := client.Order.Create(data, nil)
// 	if err != nil {
// 		c.HTML(422, "failed to create order", nil)
// 	}
// 	value := fmt.Sprint(body["id"])
// 	user_id := fmt.Sprint(user.ID)
// 	Home := Home{
// 		userid:      user_id,
// 		Name:        user.First_Name,
// 		total_price: sum,
// 		Amount:      razorpaytotal,
// 		OrderId:     value,
// 		Email:       user.Email,
// 		Contact:     user.Password,
// 	}
// 	order := domain.Order{
// 		Created_at:      time.Now(),
// 		User_Id:         user.ID,
// 		Order_Id:        value,
// 		Applied_Coupons: Coupn.Coupon,
// 		Discount:        uint(Coupn.Discount),
// 		Total_Amount:    uint(sum),
// 		Balance_Amount:  sum - Coupn.Discount,
// 		PaymentMethod:   "razorpay",
// 		Payment_Status:  "incomplete",
// 		Order_Status:    "order_placed",
// 		Address_Id:      uint(address_id),
// 	}
// 	err = cr.UserService.CreateOrder(order)
// 	if err != nil {
// 		c.HTML(422, "faile to create order", nil)
// 	}
// 	//c.HTML(200, "success ", Home)
// 	c.HTML(200, "app.html", Home)

// }
// func (cr *UserHandler) Payment_Success(c *gin.Context) {
// 	payment_id := c.Query("paymentid")
// 	orderid := c.Query("orderid")
// 	orderid = strings.Trim(orderid, " ")
// 	//signature := c.Query("signature")
// 	order, err := cr.UserService.SearchOrder(orderid)
// 	if err != nil {
// 		res := response.ErrorResponse("order is not successfull", err.Error(), "order is failed")
// 		c.Writer.WriteHeader(422)
// 		utils.ResponseJSON(c, res)
// 		return
// 	}
// 	err = cr.UserService.UpdateOrders(payment_id, orderid)
// 	if err != nil {
// 		res := response.ErrorResponse("order is no is no updated", err.Error(), "order updation problem")
// 		c.Writer.WriteHeader(422)
// 		utils.ResponseJSON(c, res)
// 	}
// 	cart, err := cr.UserService.ListViewCart(order.User_Id)
// 	if err != nil {
// 		res := response.ErrorResponse("error in list cart", err.Error(), "update orders listcart")
// 		utils.ResponseJSON(c, res)
// 		return
// 	}
// 	for _, list := range cart {
// 		err = cr.UserService.Insert_To_My_Order(list, orderid)
// 		if err != nil {
// 			res := response.ErrorResponse("error in insert into order", err.Error(), "insert into myorder")
// 			c.Writer.WriteHeader(422)
// 			utils.ResponseJSON(c, res)
// 			return
// 		}
// 	}
// 	//clear the cart
// 	err = cr.UserService.ClearCart(order.User_Id)
// 	if err != nil {
// 		res := response.ErrorResponse("error in clearCArt", err.Error(), "error in clear cart")
// 		c.Writer.WriteHeader(422)
// 		utils.ResponseJSON(c, res)
// 		return
// 	}
// 	//c.HTML(200, "success.html", "success")
// 	res := response.SuccessResponse(true, "payment success", "payment success")
// 	c.Writer.WriteHeader(200)
// 	utils.ResponseJSON(c, res)

// }
