package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}
}

// Finds whether a email is already in the database or not
func (c *adminDatabase) FindByEmail(ctx context.Context, Email string) (domain.Admin, error) {
	var admin domain.Admin
	query := `SELECT * FROM admins WHERE email=$1`
	c.DB.Raw(query, Email).Scan(&admin)
	//_ = c.DB.Where("Email=?", Email).Find(&admin)
	if admin.ID == 0 {
		return domain.Admin{}, errors.New("invalid Email")
	}
	return admin, nil
}

// Adminsign-up
func (c *adminDatabase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
	query := `INSERT INTO admins(first_name,last_name,email,phone_num,password)VALUES($1,$2,$3,$4,$5)`
	if err := c.DB.Raw(query, admin.FirstName, admin.LastName, admin.Email, admin.PhoneNum, admin.Password).Error; err != nil {
		return err
	}
	return nil
	// err := c.DB.Create(&admin).Error
	// return err
}

// List all users
func (c *adminDatabase) ListUsers(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseUsers, error) {
	offset := pagination.Offset
	limit := pagination.Limit
	var users []utils.ResponseUsers
	query := `SELECT id,first_name,last_name,email,phone_num,block from users LIMIT $1 OFFSET $2`
	err := c.DB.Raw(query, limit, offset).Scan(&users).Error
	if err != nil {
		return users, errors.New("failed to retrieve all the users")
	}
	return users, nil
}

// Manage the access of users

func (c *adminDatabase) AccessHandler(ctx context.Context, id string, access bool) error {
	err := c.DB.Model(&domain.Users{}).Where("id=?", id).UpdateColumn("block", access).Error
	if err != nil {
		return errors.New("failed to update")
	}
	return nil
}

// Admin dashboard

func (c *adminDatabase) Dashboard(ctx context.Context) (utils.ResponseWidgets, error) {
	var responseWidgets utils.ResponseWidgets
	if err := c.DB.Model(&domain.Users{}).Select("count(users)").Where("block='f'").Scan(&responseWidgets.ActiveUsers).Error; err != nil {
		return responseWidgets, err
	}
	if err := c.DB.Model(&domain.Users{}).Select("count(users)").Where("block='t'").Scan(&responseWidgets.BlockedUsers).Error; err != nil {
		return responseWidgets, err
	}
	if err := c.DB.Model(&domain.Products{}).Select("count(products)").Where("deleted_at is null").Scan(&responseWidgets.Products).Error; err != nil {
		return responseWidgets, err
	}
	if err := c.DB.Model(&domain.OrderDetails{}).Select("count(order_details)").Where("order_status_id=?", 1).Scan(&responseWidgets.Pendingorders).Error; err != nil {
		return responseWidgets, err
	}
	if err := c.DB.Model(&domain.OrderDetails{}).Select("count(order_details)").Where("order_status_id=?", 7).Scan(&responseWidgets.ReturnRequests).Error; err != nil {
		return responseWidgets, err
	}

	return responseWidgets, nil
}

// Sales report
func (c *adminDatabase) SalesReport(reqData utils.SalesReport) ([]utils.ResponseSalesReport, error) {
	var salesreport []utils.ResponseSalesReport
	if reqData.Frequency == "MONTHLY" {
		result := c.DB.Model(&domain.Order{}).Where("EXTRACT(YEAR FROM orders.placed_date) = ? AND EXTRACT(MONTH FROM orders.placed_date) = ?", reqData.Year, reqData.Month).
			Joins("JOIN order_details od on orders.id=od.order_id").
			Joins("JOIN product_details pd on pd.id=od.product_detail_id").
			Joins("JOIN products p on p.id=pd.product_id").
			Joins("JOIN payment_modes pm on pm.id=orders.payment_id").
			Joins("JOIN users u on orders.user_id=u.id").
			Joins("JOIN order_statuses os on os.id=od.order_status_id").
			// Joins("JOIN discounts d on d.id=pd.discount_id").
			Select("u.id as userid,u.first_name,u.email,od.product_detail_id as productdetailid,p.product_name as productname,od.quantity,orders.id as orderid,orders.placed_date,pm.mode as paymentmode,p.price as price,p.discount_price as discountprice,os.status as orderstatus").
			Order("orders.placed_date DESC").Scan(&salesreport)
		if result.Error != nil {
			return salesreport, result.Error
		}
	}
	if reqData.Frequency == "YEARLY" {
		result := c.DB.Model(&domain.Order{}).Where("EXTRACT(YEAR FROM orders.placed_date) = ?", reqData.Year).
			Joins("JOIN order_details od on orders.id=od.order_id").
			Joins("JOIN product_details pd on pd.id=od.product_detail_id").
			Joins("JOIN products p on p.id=pd.product_id").
			Joins("JOIN payment_modes pm on pm.id=orders.payment_id").
			Joins("JOIN users u on orders.user_id=u.id").
			Joins("JOIN order_statuses os on os.id=od.order_status_id").
			// Joins("JOIN discounts d on d.id=pd.discount_id").
			Select("u.id as userid,u.first_name,u.email,od.product_detail_id as productdetailid,p.product_name as productname,od.quantity,orders.id as orderid,orders.placed_date,pm.mode as paymentmode,p.price as price,p.discount_price as discountprice,os.status as orderstatus").
			Order("orders.placed_date DESC").Scan(&salesreport)
		if result.Error != nil {
			return salesreport, result.Error
		}
	}
	return salesreport, nil
}

// adding a new coupon
func (c *adminDatabase) AddCoupon(ctx context.Context, coupon domain.Coupon) error {
	if err := c.DB.Create(&coupon).Error; err != nil {
		return err
	}
	return nil
}

// view all coupons available
func (c *adminDatabase) GetAllCoupons(ctx context.Context, pagination utils.Pagination) ([]domain.Coupon, error) {
	var coupons []domain.Coupon
	if err := c.DB.Limit(int(pagination.Limit)).Offset(int(pagination.Offset)).Find(&coupons).Error; err != nil {
		return coupons, err
	}
	return coupons, nil
}

// updating coupons
func (c *adminDatabase) UpdateCoupon(ctx context.Context, coupon domain.Coupon, couponId string) error {
	if err := c.DB.Where("id=?", couponId).Updates(&coupon).Error; err != nil {
		return err
	}
	return nil
}

// deleting coupons
func (c *adminDatabase) DeleteCoupon(ctx context.Context, couponId string) error {
	if err := c.DB.Where("id=?", couponId).Delete(&domain.Coupon{}).Error; err != nil {
		return err
	}
	return nil
}
