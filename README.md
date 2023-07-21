# Stylezine

This is an e-commerce web application created with Go's backend framework Gin and PostgreSQL. The application consists of both admin and user sections, offering different functionalities to each role. Administrators have the ability to manage products, users, orders, and more. Users can view products, add them to their wish list or cart, make purchases, and apply applicable coupons. Authentication and authorization are handled using JSON Web Tokens (JWT), and the PostgreSQL database is simplified with GORM. The app also includes payment integration with Razorpay.

## Features

- **Admin Functionality**
  - Product Management: Administrators can create, update, and delete products.
  - User Management: Administrators have the authority to manage users, including account creation, updates, and deletions.
  - Order Management: Administrators can view and manage user orders.

- **User Functionality**
  - Product Viewing: Users can browse and view available products.
  - Wish List: Users can add products to their wish list for future reference.
  - Cart: Users can add products to their shopping cart for purchase.
  - Coupon Application: Users can apply applicable coupons to their cart.

## Technologies Used

- Go - Backend programming language
- Gin - Go web framework
- PostgreSQL - Relational database management system
- GORM - Object-relational mapping library for Go
- JWT - JSON Web Token library for Go
- Razorpay - Payment gateway integration platform

## Installation

To run this e-commerce web app locally, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/ajalck/Stylezine_e-commerce_project.git

2. Navigate to the project directory:

   cd Stylezine_e-commerce_project.git

3. Run the application:

   go run main.go

4. The application will be accessible at `http://localhost:5050`.



