package main

import (
	"fmt"
	"net/http"
	"time"
)

// User 表示系统中的用户
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

// Product 表示电子商务系统中的商品
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
}

// Order 表示用户订单
type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Products  []OrderItem `json:"products"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

// OrderItem 表示订单中的商品项
type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// NewUser 创建一个新用户
func NewUser(username, email string) *User {
	return &User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		Active:    true,
	}
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	// 这里应该有数据库查询逻辑
	// 为了示例，我们直接返回一个模拟的用户
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}
	
	return &User{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		Active:    true,
	}, nil
}

// CreateOrder 创建一个新订单
func CreateOrder(userID int, items []OrderItem) (*Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	
	if len(items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}
	
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	
	return &Order{
		UserID:    userID,
		Products:  items,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}, nil
}

// HandleUserRequest 处理用户请求的HTTP处理函数
func HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	userID := 1 // 在实际应用中，这应该从请求中提取
	
	user, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 在实际应用中，这里应该将用户数据序列化为JSON
	fmt.Fprintf(w, "User: %s (%s)", user.Username, user.Email)
}

func main() {
	// 设置HTTP路由
	http.HandleFunc("/user", HandleUserRequest)
	
	// 启动HTTP服务器
	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
} 