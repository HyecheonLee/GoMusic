package dblayer

import (
	"backend/src/models"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type DBORM struct {
	*gorm.DB
}

func NewORM(dbname, con string) (*DBORM, error) {
	db, err := gorm.Open(dbname, con)
	return &DBORM{
		DB: db,
	}, err
}

func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error
}
func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error
}
func (db *DBORM) GetCustomerByName(firstname string, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error
}
func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error
}
func (db *DBORM) GetProduct(id int) (product models.Product, error error) {
	return product, db.First(&product, id).Error
}
func (db *DBORM) AddUser(customer models.Customer) (models.Customer, error) {
	//pass received password by reference so that we can change it to its hashed version
	hashPassword(&customer.Pass)
	customer.LoggedIn = true
	err := db.Create(&customer).Error
	customer.Pass = ""
	return customer, err
}
func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {
	//Obtain a *gorm.DB object representing our customer's row
	result := db.Table("Customers").Where(&models.Customer{Email: email})
	//Retrieve the data for the customer with the passed email
	err = result.First(&customer).Error
	if err != nil {
		return customer, err
	}
	//Compare the saved hashed password with the password provided by the user trying to sign in
	if !checkPassword(customer.Pass, pass) {
		//If failed, returns an error
		return customer, ErrINVALIDPASSWORD
	}
	//set customer pass to empty because we don't need to share this information again
	customer.Pass = ""
	//update the loggedin field
	err = result.Update("loggedin", 1).Error
	if err != nil {
		return customer, err
	}
	//return the new customer row
	return customer, result.Find(&customer).Error
}
func (db *DBORM) SignOutUserById(id int) error {
	//Create a customer Go struct with the provided if
	customer := models.Customer{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	//Update the customer row to reflect the fact that the customer is not logged in
	return db.Table("Customers").Where(&customer).Update("loggedin", 0).Error
}
func (db *DBORM) GetCustomerOrdersByID(id int) (orders []models.Order, err error) {
	return orders, db.Table("orders").
		Select("*").
		Joins("join customers on customers.id = customer_id").
		Joins("join products on products.id = product_id").
		Where("customer_id=?", id).
		Scan(&orders).Error

}
func checkPassword(existingHash, incomingPass string) bool {
	//this method will return an error if the hash does not match the provided password string
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
func hashPassword(s *string) error {
	if s == nil {
		return errors.New("Reference provided for hashing password is nil")
	}

	//convert password string to byte slice so that we can use it with the bcrypt package
	sBytes := []byte(*s)

	//Obtain hashed password via the GenerateFromPassword() method
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//update password string with the hashed version
	*s = string(hashedBytes[:])
	return nil
}
