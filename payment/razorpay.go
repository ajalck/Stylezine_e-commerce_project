package payment

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

type Pay struct {
	db *gorm.DB
}
type Payment interface {
	RazorPay(c *gin.Context)
	PaymentStatus(c *gin.Context)
}

func PaymentHandler(db *gorm.DB) Payment {
	return &Pay{db: db}
}

type Home struct {
	userid      string
	Name        string
	total_price int
	Amount      int
	OrderId     string
	Email       string
	Contact     string
}

// @Summary Payment
// @ID user payment
// @Tags User Payment
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/checkout/razorpay [get]
func (p *Pay) RazorPay(c *gin.Context) {
	// user_id := c.Writer.Header().Get("user_id")
	// fmt.Println(user_id)
	user_id := "GyhWnmxVjY"
	orders := domain.Order{}
	result := p.db.Table("orders").Where(&domain.Order{User_ID: user_id, Order_Status: "pending"}).First(&orders)
	if result.Error != nil {
		response := utils.ErrorResponse("Payment Unsuccessfull", result.Error.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	// requestURL := "/user/checkout/razorpay"
	// http.Redirect(c.Writer, c.Request, requestURL, http.StatusSeeOther)
	client := razorpay.NewClient("rzp_test_qCClRw9jUfHAD8", "UidHSAWYPWpHdX43nKZfsDku")
	data := map[string]interface{}{
		"amount":   orders.Final * 100,
		"currency": "INR",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(422, gin.H{"Failed to create order": err})
		c.HTML(422, "failed to create order", nil)
		return
	}
	fmt.Println(body)
	value := fmt.Sprint(body["id"])
	// var user domain.User
	// result = p.db.Table("users").Where("user_id", user_id).First(user)
	// if result.Error != nil {
	// 	response := utils.ErrorResponse("Payment Unsuccessfull", result.Error.Error(), nil)
	// 	c.Writer.WriteHeader(400)
	// 	utils.ResponseJSON(c, response)
	// 	return
	// }
	p.db.Table("orders").Where("user_id", user_id).Update("razor_pay_order_id", value)

	Home := Home{
		userid:      user_id,
		Name:        "Ajal",
		total_price: int(orders.Grand_Total),
		Amount:      int(orders.Final),
		OrderId:     value,
		Email:       "email",
		Contact:     "phone",
	}
	fmt.Println(Home)
	//c.HTML(200, "success ", Home)
	c.HTML(200, "checkout.html", Home)

}

// @Summary Payment
// @ID user payment_Status
// @Tags User Payment
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/payment-success [get]
func (p *Pay) PaymentStatus(c *gin.Context) {
	c.HTML(200, "success.html", nil)
	payment_id := c.Query("paymentid")
	orderid := c.Query("orderid")
	orderid = strings.Trim(orderid, " ")
	result := p.db.Where("razor_pay_order_id", orderid).Updates(&domain.Order{Payment_ID: payment_id, Mode_of_Payment: "razorpay", Order_Status: "succeeded", Payment_Status: "succeeded"})
	if result.Error != nil {
		response := utils.ErrorResponse("can't update order table", result.Error.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
}
