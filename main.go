package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// ---------- Модели ----------

type Supplier struct {
	ID      int    `json:"id" db:"id_Поставщика"`
	Name    string `json:"name" db:"Название"`
	Phone   string `json:"phone" db:"Телефон"`
	Address string `json:"address" db:"Адрес"`
	Email   string `json:"email" db:"Email"`
}

type Supply struct {
	ID         int       `json:"id" db:"id_Поставки"`
	SupplierID int       `json:"supplier_id" db:"id_Поставщика"`
	SupplyDate time.Time `json:"supply_date" db:"Дата поставки"`
}

type Material struct {
	ID       int    `json:"id" db:"id_Материала"`
	Name     string `json:"name" db:"Название"`
	Category string `json:"category" db:"Категория"`
}

type SupplyDetail struct {
	ID         int `json:"id" db:"id_Записи"`
	SupplyID   int `json:"supply_id" db:"id_Поставки"`
	MaterialID int `json:"material_id" db:"id_Материала"`
	Quantity   int `json:"quantity" db:"Количество"`
	Price      int `json:"price" db:"Цена"`
}

// ---------- Репозиторий ----------

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// ----- Поставщики -----

func (r *Repository) GetSuppliers() ([]Supplier, error) {
	var suppliers []Supplier
	err := r.db.Select(&suppliers, "SELECT * FROM `Поставщик`")
	return suppliers, err
}

func (r *Repository) GetSupplierByID(id int) (*Supplier, error) {
	var supplier Supplier
	err := r.db.Get(&supplier, "SELECT * FROM `Поставщик` WHERE `id_Поставщика` = ?", id)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *Repository) CreateSupplier(supplier *Supplier) error {
	query := "INSERT INTO `Поставщик` (`id_Поставщика`, `Название`, `Телефон`, `Адрес`, `Email`) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, supplier.ID, supplier.Name, supplier.Phone, supplier.Address, supplier.Email)
	return err
}

func (r *Repository) UpdateSupplier(supplier *Supplier) error {
	query := "UPDATE `Поставщик` SET `Название` = ?, `Телефон` = ?, `Адрес` = ?, `Email` = ? WHERE `id_Поставщика` = ?"
	_, err := r.db.Exec(query, supplier.Name, supplier.Phone, supplier.Address, supplier.Email, supplier.ID)
	return err
}

func (r *Repository) DeleteSupplier(id int) error {
	_, err := r.db.Exec("DELETE FROM `Поставщик` WHERE `id_Поставщика` = ?", id)
	return err
}

// ----- Поставки -----

func (r *Repository) GetSupplies() ([]Supply, error) {
	var supplies []Supply
	err := r.db.Select(&supplies, "SELECT * FROM `Поставка`")
	return supplies, err
}

func (r *Repository) GetSupplyByID(id int) (*Supply, error) {
	var supply Supply
	err := r.db.Get(&supply, "SELECT * FROM `Поставка` WHERE `id_Поставки` = ?", id)
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *Repository) GetSuppliesBySupplier(supplierID int) ([]Supply, error) {
	var supplies []Supply
	err := r.db.Select(&supplies, "SELECT * FROM `Поставка` WHERE `id_Поставщика` = ?", supplierID)
	return supplies, err
}

func (r *Repository) CreateSupply(supply *Supply) error {
	query := "INSERT INTO `Поставка` (`id_Поставки`, `id_Поставщика`, `Дата поставки`) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, supply.ID, supply.SupplierID, supply.SupplyDate)
	return err
}

// ----- Материалы -----

func (r *Repository) GetMaterials() ([]Material, error) {
	var materials []Material
	err := r.db.Select(&materials, "SELECT * FROM `Материал`")
	return materials, err
}

func (r *Repository) GetMaterialByID(id int) (*Material, error) {
	var material Material
	err := r.db.Get(&material, "SELECT * FROM `Материал` WHERE `id_Материала` = ?", id)
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *Repository) CreateMaterial(material *Material) error {
	query := "INSERT INTO `Материал` (`id_Материала`, `Название`, `Категория`) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, material.ID, material.Name, material.Category)
	return err
}

func (r *Repository) UpdateMaterial(material *Material) error {
	query := "UPDATE `Материал` SET `Название` = ?, `Категория` = ? WHERE `id_Материала` = ?"
	_, err := r.db.Exec(query, material.Name, material.Category, material.ID)
	return err
}

func (r *Repository) DeleteMaterial(id int) error {
	_, err := r.db.Exec("DELETE FROM `Материал` WHERE `id_Материала` = ?", id)
	return err
}

// ----- Детали поставки -----

func (r *Repository) GetSupplyDetails() ([]SupplyDetail, error) {
	var details []SupplyDetail
	err := r.db.Select(&details, "SELECT * FROM `Детали поставки`")
	return details, err
}

func (r *Repository) GetSupplyDetailByID(id int) (*SupplyDetail, error) {
	var detail SupplyDetail
	err := r.db.Get(&detail, "SELECT * FROM `Детали поставки` WHERE `id_Записи` = ?", id)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *Repository) GetDetailsBySupply(supplyID int) ([]SupplyDetail, error) {
	var details []SupplyDetail
	err := r.db.Select(&details, "SELECT * FROM `Детали поставки` WHERE `id_Поставки` = ?", supplyID)
	return details, err
}

func (r *Repository) CreateSupplyDetail(detail *SupplyDetail) error {
	query := "INSERT INTO `Детали поставки` (`id_Записи`, `id_Поставки`, `id_Материала`, `Количество`, `Цена`) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, detail.ID, detail.SupplyID, detail.MaterialID, detail.Quantity, detail.Price)
	return err
}

func (r *Repository) UpdateSupplyDetail(detail *SupplyDetail) error {
	query := "UPDATE `Детали поставки` SET `id_Поставки` = ?, `id_Материала` = ?, `Количество` = ?, `Цена` = ? WHERE `id_Записи` = ?"
	_, err := r.db.Exec(query, detail.SupplyID, detail.MaterialID, detail.Quantity, detail.Price, detail.ID)
	return err
}

func (r *Repository) DeleteSupplyDetail(id int) error {
	_, err := r.db.Exec("DELETE FROM `Детали поставки` WHERE `id_Записи` = ?", id)
	return err
}

// CreateSupplyWithDetails создаёт поставку и добавляет связанные детали (материалы).
func (r *Repository) CreateSupplyWithDetails(supply *Supply, details []SupplyDetail) error {
	query := "INSERT INTO `Поставка` (`id_Поставки`, `id_Поставщика`, `Дата поставки`) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, supply.ID, supply.SupplierID, supply.SupplyDate)
	if err != nil {
		return err
	}

	detailQuery := "INSERT INTO `Детали поставки` (`id_Записи`, `id_Поставки`, `id_Материала`, `Количество`, `Цена`) VALUES (?, ?, ?, ?, ?)"
	for _, d := range details {
		_, err := r.db.Exec(detailQuery, d.ID, supply.ID, d.MaterialID, d.Quantity, d.Price)
		if err != nil {
			return err
		}
	}

	return nil
}

// Обновление поставки
func (r *Repository) UpdateSupply(supply *Supply) error {
	query := "UPDATE `Поставка` SET `id_Поставщика` = ?, `Дата поставки` = ? WHERE `id_Поставки` = ?"
	_, err := r.db.Exec(query, supply.SupplierID, supply.SupplyDate, supply.ID)
	return err
}

// Удаление поставки вместе с её деталями
func (r *Repository) DeleteSupply(id int) error {
	// Сначала удалим связанные детали поставки для безопасности
	_, err := r.db.Exec("DELETE FROM `Детали поставки` WHERE `id_Поставки` = ?", id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM `Поставка` WHERE `id_Поставки` = ?", id)
	return err
}

// Получить материалы для поставки
func (r *Repository) GetSupplyMaterials(supplyID int) ([]Material, error) {
	query := "SELECT m.* FROM `Материал` m JOIN `Детали поставки` d ON m.`id_Материала` = d.`id_Материала` WHERE d.`id_Поставки` = ?"
	var materials []Material
	err := r.db.Select(&materials, query, supplyID)
	return materials, err
}

// ---------- Обработчики ----------

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// ----- Поставщики -----

func GetSuppliers(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		suppliers, err := repo.GetSuppliers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		if len(suppliers) == 0 {
			w.Write([]byte("Поставщиков нет"))
			return
		}
		for _, s := range suppliers {
			w.Write([]byte(fmt.Sprintf("ID: %d, Поставщик: %s, Телефон: %s, Адрес: %s, Email: %s\n",
				s.ID, s.Name, s.Phone, s.Address, s.Email)))
		}
	}
}

func GetSupplierByID(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		supplier, err := repo.GetSupplierByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставщик с ID %d не найден", id)))
			return
		}
		w.Write([]byte(fmt.Sprintf("ID: %d, Поставщик: %s, Телефон: %s, Адрес: %s, Email: %s",
			supplier.ID, supplier.Name, supplier.Phone, supplier.Address, supplier.Email)))
	}
}

func CreateSupplier(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supplier Supplier

		err := json.NewDecoder(r.Body).Decode(&supplier)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}

		err = repo.CreateSupplier(&supplier)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка создания поставщика: %v", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Поставщик создан!"))
	}
}

func UpdateSupplier(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}

		supplier, err := repo.GetSupplierByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставщик с ID %d не найден", id)))
			return
		}

		// Читаем новые данные из тела запроса
		err = json.NewDecoder(r.Body).Decode(&supplier)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}
		supplier.ID = id // сохраняем ID из URL

		err = repo.UpdateSupplier(supplier)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка обновления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Поставщик обновлён!"))
	}
}

func DeleteSupplier(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		_, err = repo.GetSupplierByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставщик с ID %d не найден", id)))
			return
		}
		err = repo.DeleteSupplier(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка удаления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Поставщик удалён!"))
	}
}

// ----- Поставки -----

func GetSupplies(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supplies, err := repo.GetSupplies()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		if len(supplies) == 0 {
			w.Write([]byte("Поставок нет"))
			return
		}
		for _, s := range supplies {
			w.Write([]byte(fmt.Sprintf("ID: %d, ID Поставщика: %d, Дата поставки: %s\n",
				s.ID, s.SupplierID, s.SupplyDate.Format("2006-01-02"))))
		}
	}
}

func GetSupplyByID(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		supply, err := repo.GetSupplyByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставка с ID %d не найдена", id)))
			return
		}
		w.Write([]byte(fmt.Sprintf("ID: %d, ID Поставщика: %d, Дата поставки: %s",
			supply.ID, supply.SupplierID, supply.SupplyDate.Format("2006-01-02"))))
	}
}

func CreateSupply(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supply Supply

		err := json.NewDecoder(r.Body).Decode(&supply)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}

		err = repo.CreateSupply(&supply)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка создания поставки: %v", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Поставка создана!"))
	}
}

func UpdateSupply(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}

		supply, err := repo.GetSupplyByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставка с ID %d не найдена", id)))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&supply)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}
		supply.ID = id

		err = repo.UpdateSupply(supply)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка обновления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Поставка обновлена!"))
	}
}

func DeleteSupply(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		_, err = repo.GetSupplyByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Поставка с ID %d не найдена", id)))
			return
		}
		err = repo.DeleteSupply(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка удаления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Поставка удалена!"))
	}
}

// ----- Материалы -----

func GetMaterials(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		materials, err := repo.GetMaterials()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		if len(materials) == 0 {
			w.Write([]byte("Материалов нет"))
			return
		}
		for _, m := range materials {
			w.Write([]byte(fmt.Sprintf("ID: %d, Название: %s, Категория: %s\n",
				m.ID, m.Name, m.Category)))
		}
	}
}

func GetMaterialByID(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		material, err := repo.GetMaterialByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Материал с ID %d не найден", id)))
			return
		}
		w.Write([]byte(fmt.Sprintf("ID: %d, Название: %s, Категория: %s",
			material.ID, material.Name, material.Category)))
	}
}

func CreateMaterial(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var material Material

		err := json.NewDecoder(r.Body).Decode(&material)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}

		err = repo.CreateMaterial(&material)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка создания материала: %v", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Материал создан!"))
	}
}

func UpdateMaterial(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}

		material, err := repo.GetMaterialByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Материал с ID %d не найден", id)))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&material)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}
		material.ID = id

		err = repo.UpdateMaterial(material)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка обновления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Материал обновлён!"))
	}
}

func DeleteMaterial(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		_, err = repo.GetMaterialByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Материал с ID %d не найден", id)))
			return
		}
		err = repo.DeleteMaterial(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка удаления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Материал удалён!"))
	}
}

// ----- Детали поставки -----

func GetSupplyDetails(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, err := repo.GetSupplyDetails()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		if len(details) == 0 {
			w.Write([]byte("Деталей поставок нет"))
			return
		}
		for _, d := range details {
			w.Write([]byte(fmt.Sprintf("ID записи: %d, ID Поставки: %d, ID Материала: %d, Количество: %d, Цена: %d\n",
				d.ID, d.SupplyID, d.MaterialID, d.Quantity, d.Price)))
		}
	}
}

func GetSupplyDetailByID(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		detail, err := repo.GetSupplyDetailByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Запись с ID %d не найдена", id)))
			return
		}
		w.Write([]byte(fmt.Sprintf("ID записи: %d, ID Поставки: %d, ID Материала: %d, Количество: %d, Цена: %d",
			detail.ID, detail.SupplyID, detail.MaterialID, detail.Quantity, detail.Price)))
	}
}

func CreateSupplyDetail(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var detail SupplyDetail
		err := json.NewDecoder(r.Body).Decode(&detail)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}
		err = repo.CreateSupplyDetail(&detail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка создания детали поставки: %v", err)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Деталь поставки добавлена!"))
	}
}

func UpdateSupplyDetail(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		detail, err := repo.GetSupplyDetailByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Запись с ID %d не найдена", id)))
			return
		}
		err = json.NewDecoder(r.Body).Decode(&detail)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Некорректный JSON в запросе: %v", err)))
			return
		}
		detail.ID = id
		err = repo.UpdateSupplyDetail(detail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка обновления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Деталь поставки обновлена!"))
	}
}

func DeleteSupplyDetail(repo *Repository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный ID"))
			return
		}
		_, err = repo.GetSupplyDetailByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Запись с ID %d не найдена", id)))
			return
		}
		err = repo.DeleteSupplyDetail(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Ошибка удаления: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Деталь поставки удалена!"))
	}
}

// ---------- Главная функция ----------

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/furniture_factory?parseTime=true"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	log.Println("Подключено к MySQL!")

	repo := NewRepository(db)
	r := mux.NewRouter()

	// Поставщики
	r.HandleFunc("/suppliers", GetSuppliers(repo)).Methods("GET")
	r.HandleFunc("/suppliers/{id}", GetSupplierByID(repo)).Methods("GET")
	r.HandleFunc("/suppliers", CreateSupplier(repo)).Methods("POST")
	r.HandleFunc("/suppliers/{id}", UpdateSupplier(repo)).Methods("PUT")
	r.HandleFunc("/suppliers/{id}", DeleteSupplier(repo)).Methods("DELETE")

	// Поставки
	r.HandleFunc("/supplies", GetSupplies(repo)).Methods("GET")
	r.HandleFunc("/supplies/{id}", GetSupplyByID(repo)).Methods("GET")
	r.HandleFunc("/supplies", CreateSupply(repo)).Methods("POST")
	r.HandleFunc("/supplies/{id}", UpdateSupply(repo)).Methods("PUT")
	r.HandleFunc("/supplies/{id}", DeleteSupply(repo)).Methods("DELETE")

	// Материалы
	r.HandleFunc("/materials", GetMaterials(repo)).Methods("GET")
	r.HandleFunc("/materials/{id}", GetMaterialByID(repo)).Methods("GET")
	r.HandleFunc("/materials", CreateMaterial(repo)).Methods("POST")
	r.HandleFunc("/materials/{id}", UpdateMaterial(repo)).Methods("PUT")
	r.HandleFunc("/materials/{id}", DeleteMaterial(repo)).Methods("DELETE")

	// Детали поставки
	r.HandleFunc("/supply-details", GetSupplyDetails(repo)).Methods("GET")
	r.HandleFunc("/supply-details/{id}", GetSupplyDetailByID(repo)).Methods("GET")
	r.HandleFunc("/supply-details", CreateSupplyDetail(repo)).Methods("POST")
	r.HandleFunc("/supply-details/{id}", UpdateSupplyDetail(repo)).Methods("PUT")
	r.HandleFunc("/supply-details/{id}", DeleteSupplyDetail(repo)).Methods("DELETE")

	log.Println("Сервер запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
