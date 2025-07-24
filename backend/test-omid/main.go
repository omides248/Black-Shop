package main

import "fmt"

type ProductStorer interface {
	Find(id int) string
}

type TestStore struct{}

func (s TestStore) Find(id int) string {
	return fmt.Sprintf("محصول تستی شماره %d از انبار موقت", id)
}

type ProductionStore struct{}

func (s ProductionStore) Find(id int) string {
	return fmt.Sprintf("محصول واقعی شماره %d از دیتابیس اصلی", id)
}

type Employee struct {
	Store ProductStorer
}

func (e Employee) GetProductInfo(productID int) {
	productName := e.Store.Find(productID)
	fmt.Println("نتیجه جستجو:", productName)
}

func main() {
	testDB := TestStore{}
	productionDB := ProductionStore{}

	employee := Employee{}

	fmt.Println("--- فاز اول: کار با انبار تستی ---")

	employee.Store = testDB
	employee.GetProductInfo(101)

	fmt.Println("\n--- فاز دوم: کار با انبار اصلی ---")

	employee.Store = productionDB
	employee.GetProductInfo(202)
}
